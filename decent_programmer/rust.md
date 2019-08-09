# iterator

# builder
```rust
pub struct Builder<I, E = Exec> {
    incoming: I,
    protocol: Http_<E>,
}

impl<I, E> Builder<I, E> {
    pub fn new(incoming: I, protocol: Http_<E>) -> Self {
        Builder {
            incoming,
            protocol,
        }
    }

    /// Sets whether to use keep-alive for HTTP/1 connections.
    /// Default is `true`.
    pub fn http1_keepalive(mut self, val: bool) -> Self {
        self.protocol.keep_alive(val);
        self
    }
}
```


# template function


# build type relationship to help static check
```rust
pub trait Actor: Sized + 'static {
    type Context: ActorContext;

    fn start(self) -> Addr<Self>
    where
        Self: Actor<Context = Context<Self>>,
    {   
        Context::new().run(self)
    }   
}

pub struct Context<A>
where
    A: Actor<Context = Context<A>>,
{
    parts: ContextParts<A>,
    mb: Option<Mailbox<A>>,
}

pub struct Mailbox<A>
where
    A: Actor,
    A::Context: AsyncContext<A>,
{
    msgs: AddressReceiver<A>,
}

pub struct Addr<A: Actor> {
    tx: AddressSender<A>,
}

impl<A: Actor> Addr<A> {
    pub fn new(tx: AddressSender<A>) -> Addr<A> {
        Addr { tx }
    }

    pub fn send<M>(&self, msg: M) -> Request<A, M>
    where
        M: Message + Send,
        M::Result: Send,
        A: Handler<M>,
        A::Context: ToEnvelope<A, M>,
    {
    }
}
```
Actor a process which could handle specified message
Context mailbox and other realted state associate with actior
Address to identify a actior, and as the target for message passing.
The generic type using static type checking, at compile time to make sure
1 Messages which one actor cared about could be send to the actor
2 Two actor address is different and won't be misused

# runtime reflection
```rust
type AnyMap = HashMap<TypeId, Box<Any>>;

pub struct Registry {
    registry: Rc<RefCell<AnyMap>>,
}

thread_local! {
    static AREG: Registry = {
        Registry {
            registry: Rc::new(RefCell::new(AnyMap::new()))
        }
    };
}

impl Registry {
    pub fn get<A: ArbiterService + Actor<Context = Context<A>>>(&self) -> Addr<A> {
        let id = TypeId::of::<A>();
        if let Some(addr) = self.registry.borrow().get(&id) {
            if let Some(addr) = addr.downcast_ref::<Addr<A>>() {
                return addr.clone();
            }
        }
        let addr: Addr<A> = A::start_service();

        self.registry
            .borrow_mut()
            .insert(id, Box::new(addr.clone()));
        addr
    }
}
```
Not all type in rust implement Any.
RefCell make immutable reference could also change the state, which make state
share in same thread much easier 

# vec manipulate pattern
```rust
fn drop(&mut self) {
    //drain will clean the old the vec
    for (_, outbound) in self.outbound_streams.drain() {
        self.mutex.destroy_outbound(outbound);
    }
}

//pop from the vec and handle it, if failed push back
fn poll(&mut self) {
    for n in (0..self.outbound_streams.len()).rev() {
        let (user_data, mut outbound) = self.outbound_streams.swap_remove(n);
        match self.muxer.poll_outbound(&mut outbound) {
            Ok(_) => {}
            Err(_) => {
                self.outbound_streams.push((user_data, outbound));
            }
    }
}
```


# Gotchas
## mutable variable and mutable reference 
```rust
let mut a = 10;
let b = &mut a;
*b = 10;
let mut c = 20;
b = &mut c; //this is illegal
```
mutable variable means variable itself could be changed, mutable reference
means the reference could be used to change the value it point to.

## lifetime of temporary variable

## generic trait vs trait with associate type
Generic trait is not a type, it's a kind
Trait with associate type, is still a trait, the associate type
is like the funciton declared which must be specified during 
implementation. 
```rust
pub trait Person {
    fn get_name(&self) -> i32;
    fn from_str(&str) -> Self;
}

fn create_person_gen<T: Person>(name: String) -> T { 
    <T as Person>::from_str(name.as_ref())
}
```

For static dispatch, trait with associate type could be used as constraint for generic
type without specify associate type, since generic function will generate concreate function
at compiled time, the associate type will be inferred with concrete object invoke. 
For dynamic dispatch, which using trait object, the trait and the associate type has to
be specified.

## closure is struct with special method
1. `fn(A) -> B` a general function is pure function, which output only depends on input.
```rust
let obj1 = MyStruct::new("Hello", 15);
let obj2 = MyStruct::new("More Text", 10);
let fn1 = |x: &MyStruct| x.get_number() + 3;
```

1. `Fn(A) -> B` a closure/struct which normally has borrow or move the context around it, 
but it won't modify them 
```rust
let obj1 = MyStruct::new("Hello", 15);
let obj2 = MyStruct::new("More Text", 10);
// obj1 is borrowed by the closure immutably.
let closure1 = |x: &MyStruct| x.get_number() + obj1.get_number();
assert_eq!(closure1(&obj2), 25);


struct ClosureType<'a> {
    obj: &'a MyStruct,
}
impl<'a> ClosureType<'a> {
    fn call(&self, x: &MyStruct) {
        x.get_number() + self.obj.get_number()
    }
}
```

1. `FnMut(A) -> B` a closure/sturct which normally has mut borrow or move the context, it 
will change the context
```rust
let mut obj1 = MyStruct::new("Hello", 15);
let obj2 = MyStruct::new("More Text", 10);
// obj1 is borrowed by the closure mutably.
let mut closure3 = |x: &MyStruct| {
    obj1.inc_number();
    x.get_number() + obj1.get_number()
};


struct ClosureType<'a> {
    obj: &'a mut MyStruct,
}
impl<'a> ClosureType<'a> {
    fn call(&mut self, x: &MyStruct) {
        self.obj.inc_number();
        x.get_number() + self.obj.get_number()
    }
}

```

1. `FnOnce(A) -> B` a closure/struct which normally consume the context, can can only invoked
once
```rust
let obj1 = MyStruct::new("Hello", 15);
let obj2 = MyStruct::new("More Text", 10);
// obj1 is owned by the closure
let closure4 = |x: &MyStruct| {
    obj1.destructor();
    x.get_number()
};


struct ClosureType {
    obj: MyStruct,
}
impl<'a> ClosureType<'a> {
    fn call(self, x: &MyStruct) {
        self.obj.destructor();
        x.get_number() 
    }
}
```

## empty enum vs bottom type
In haskell, there is bottom type which represent a endless loop, it's the subtype of any other
type, in rust matching on a empty enum will generate `!` type, which is a builtin type and can
be coerces to any other type, empty enum will never get an instance of it, so match on it will
never happen, the match will never be executed.
```rust
pub enum Infallible {}
impl fmt::Display for Infallible {
    fn fmt(&self, _: &mut fmt::Formatter) -> fmt::Result {
        match *self {}
    }
}
```

## static lifetime
`'static` doesn't mean the object will live until the end of program, it means the lifetime of 
the object isn't restrict by other variable, `'a` means a variable couldn't live logger than '`a'
so normally it has reference to other variable whose lifetime is `'a`, from this point of view, 
static lifetime enforce the object includes no reference to other variables.  
