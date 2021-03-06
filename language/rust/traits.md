# types of trait
1. Marker traits
* Traits without any methods. like ```Copy, Send, Sync```
* They are used to simply mark a type as belongs to a 
  particular family for t gain some compile time guarantees.
2. Simple traits
3. Generic traits
4. Associated type traits
5. Inherited traits.

# two kind of methods in trait
1. Associate methods. ```self``` isn't in parameter but return
   part
2. Instance methods, ```self``` is in the parameter part, normally
   the first paremeter, therefore, these methods are only available
   on the instance of the type that are implementing the trait.

# return generic type with trait bound
```rust
trait Person {}
fn new_persion<T: Person>() -> T {} //won't compile
fn new_persion() -> impl Person {}
```

# trait object
obj ptr + vtable
normally, there is only one vtable
different from c++, vtable is create in trait object when
call dynamically dispatched function instead of stored
in each instance.

```shell
$ rustc --explain E0225
You attempted to use multiple types as bounds for a closure or trait object.
Rust does not currently support this. A simple example that causes this error:
```

```rust
fn main() {
    let _: Box<std::io::Read+std::io::Write>;
}
```

Builtin traits are an exception to this rule: it's possible to have bounds of
one non-builtin type, plus any number of builtin types. For example, the
following compiles correctly:
```rust
fn main() {
    let _: Box<std::io::Read+Copy+Sync>;
}
```

One way to work around this is to crate a new trait with these traits as super
trait
```rust
trait MyTrait: Any + From<String> + PartialOrd {}
Box<MyTrait>
```

# object safety: which trait could be put into a trait object
```rust 
pub struct TraitObject {
    pub data: *mut (),
    pub vtable: *mut (),
}
```
trait object is noramlly used as a way to do dynamic dispatch
trait object is compile time trick, so the restrictions are the set of things related
to compile time.
1. concrete type is erased
2. different vtable is generated for each concrete type
3. method in trait is all object safety then the trait is object safety, then it 
could be put into trait object. The following method isn't object safe
    1. use Self as paramter or return value
        ```rust
        fn foo(&self) -> Self; 
        let y = x.foo(); //concrete type is erased, type of y is unknown
        ```
    2. has generic method
        ```rust
        fn generic_method<A>(&self, value: A); //vtable is very hard to generate
        ```
    3. method without self paremter
        ```rust
        fn foo() -> u8 //no self, no vtable
        ```

# built-in trait:
```rust
 impl PartialOrd for Person {
    fn partial_cmp(&self, other: &Person) -> Option<Ordering> {
        Some(self.cmp(other))
    }
}
//-------------------------------------------------------
impl Ord for Person {
    fn cmp(&self, other: &Person) -> Ordering {
        self.height.cmp(&other.height)
    }
}
//-------------------------------------------------------
impl PartialEq for Book {
    fn eq(&self, other: &Book) -> bool {
        self.isbn == other.isbn
    }
}
//-------------------------------------------------------
impl Neg for Sign {
    type Output = Sign;
    fn neg(self) -> Sign {
        match self {
            Sign::Negative => Sign::Positive,
            Sign::Zero => Sign::Zero,
            Sign::Positive => Sign::Negative,
        }
    }
}
//-------------------------------------------------------
#[derive(Copy)]
struct Stats {
   frequencies: [i32; 100],
}

impl Clone for Stats {
    fn clone(&self) -> Stats { *self }
}
//-------------------------------------------------------

struct Point {
    x: i32,
    y: i32,
}

impl fmt::Debug for Point {
    fn fmt(&self, f: &mut fmt::Formatter) -> fmt::Result {
        write!(f, "Point {{ x: {}, y: {} }}", self.x, self.y)
    }
}

impl fmt::Display for Point {
    fn fmt(&self, f: &mut fmt::Formatter) -> fmt::Result {
        write!(f, "({}, {})", self.x, self.y)
    }
}

let origin = Point { x: 0, y: 0 };

println!("The origin is: {:?}", origin); //call debug
println!("The origin is: {}", origin);   //call display
//-------------------------------------------------------

impl Drop for HasDrop {
    fn drop(&mut self) {
        println!("Dropping!");
    }
}
//-------------------------------------------------------
struct Foo {
    name: String,
    age: u8,
    hobby: Option<Hobby>,
}


impl Default for Foo {
    fn default() -> Self {
        Foo { 
            name: "".to_owned(),
            age: 0,
            hobby: None
        }
    }
}

let f = Foo { name: "ben", ..Default::default()};
//-------------------------------------------------------
impl Error for SuperError {
    fn description(&self) -> &str {
        "I'm the superhero of errors"
    }

    fn cause(&self) -> Option<&Error> {
        Some(&self.side)
    }
}
//-------------------------------------------------------
impl Add for Point {
    type Output = Point;

    fn add(self, other: Point) -> Point {
        Point {
            x: self.x + other.x,
            y: self.y + other.y,
        }
    }
}
//-------------------------------------------------------
impl Index<Side> for Balance {
    type Output = Weight;
    fn index<'a>(&'a self, index: Side) -> &'a Weight {
        println!("Accessing {:?}-side of balance immutably", index);
        match index {
            Side::Left => &self.left,
            Side::Right => &self.right,
        }
    }
}
//-------------------------------------------------------
impl IndexMut<Side> for Balance {
    fn index_mut<'a>(&'a mut self, index: Side) -> &'a mut Weight {
        println!("Accessing {:?}-side of balance mutably", index);
        match index {
            Side::Left => &mut self.left,
            Side::Right => &mut self.right,
        }
    }
}
//-------------------------------------------------------
struct DerefExample<T> {
    value: T
}

impl<T> Deref for DerefExample<T> {
    type Target = T;

    fn deref(&self) -> &T {
        &self.value
    }
}

impl<T> DerefMut for DerefMutExample<T> {
    fn deref_mut(&mut self) -> &mut T {
        &mut self.value
    }
}

let x = DerefExample { value: 'a' };
assert_eq!('a', *x);
let mut x = DerefExample { value: 'a' };
*x = b;
assert_eq!('b', *x);
//-------------------------------------------------------
struct SortedVec<T>(Vec<T>);
impl<'a, T: Ord + Clone> From<&'a [T]> for SortedVec<T> {
    fn from(slice: &[T]) -> Self {
        let mut vec = slice.to_owned();
        vec.sort();
        SortedVec(vec)
    }
}

impl<T: Ord + Clone> From<Vec<T>> for SortedVec<T> {
    fn from(mut vec: Vec<T>) -> Self {
        vec.sort();
        SortedVec(vec)
    }
}

impl<T> AsRef<Vec<T>> for SortedVec<T> {
    fn as_ref(&self) -> &Vec<T> {
        &self.0
    }
}

impl<T> AsMut<Vec<T>> for SortedVec<T> {
    fn as_mut(&mut self) -> &mut Vec<T> {
        &mut self.0
    }
}

//SortedVec could be used as parameter to function manipulate_vector
fn manipulate_vector<T, V: AsRef<Vec<T>>>(vec: V) -> Result<usize, ()> 

#[derive(Clone, Eq)]
pub struct Label(Rc<[u8]>);

impl AsRef<[u8]> for Label {
    fn as_ref(&self) -> &[u8] {
        &self.0
    }   
}

impl Borrow<[u8]> for Label {
    fn borrow(&self) -> &[u8] {
        &self.0
    }   
}
//from signature, borrow is same with AsRef, which is &self -> &T
//borrow means &self and &T has same hash and ordering result
//borrow is mainly used by standard container in std lib like HashMap
//and BTreeMap, self and T could be used as search key, when self
//is used as the key for the containers
//-------------------------------------------------------
pub trait FnOnce<Args> {
    type Output;
    extern "rust-call" fn call_once(self, args: Args) -> Self::Output;
}

pub trait FnMut<Args>: FnOnce<Args> {
    extern "rust-call" fn call_mut(&mut self, args: Args) -> Self::Output;
}

pub trait Fn<Args>: FnMut<Args> {
    extern "rust-call" fn call(&self, args: Args) -> Self::Output;
}
//-------------------------------------------------------
pub trait Extend<A> {
    fn extend<T: IntoIterator<Item = A>>(&mut self, iter: T);
}
```
