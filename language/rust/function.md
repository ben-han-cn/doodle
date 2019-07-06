# Static dispatching
## Method lookup
1. Inherent methods
Any methods defined within an impl like ``impl foo``
2. Extension methods
  * Derived from imported traits. So the trait must be import first.
  * Trait is the rust way to implement function overloading
  * when two traits has same prototype, use <F as XTrait>::foo to 
    specify which function to call
  * Rust support dispatch function based on return type 
```rust
//self is the iterator, the implementation is from iter which is
//dispatched base on B
fn collect<B: FromIterator<Self::Item>>(self) -> B where Self: Sized {
    FromIterator::from_iter(self)
}
```

## Generic method
Compiler creates a specialised version of the generic function
for every type used with it.
```rust
fn func<T: Foo>(x: &T);
fn func<T>(x: T) where T: Foo;
fn func(x: impl Foo);
```
Note: Type in generic function is sized by default. 

# Dynamic dispatching
dispatch based on self (parameter)
Function use trait object
```rust
fn func(x: &dyn Foo);
fn func(x: &Box<dyn Foo>);
```

#function compose
```rust
macro_rules! compose {
    ( $last:expr ) => { $last };
    ( $head:expr, $($tail:expr), +) => {
        compose_two($head, compose!($($tail),+))
    };
}

fn compose_two<A, B, C, G, F>(f: F, g: G) -> impl Fn(A) -> C
where
    F: Fn(A) -> B,
    G: Fn(B) -> C,
{
    move |x| g(f(x))
}

fn main() {
    let add = |x| x + 2;
    let multiply = |x| x * 2;
    let divide = |x| x / 2;
    let intermediate = compose!(add, multiply, divide);

    let subtract = |x| x - 2;
    let finally = compose!(intermediate, subtract);

    println!("Result is {}", finally(10));
}
```
