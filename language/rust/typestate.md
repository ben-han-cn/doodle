# Typestate Pattern 
is an API design pattern that encodes information about an object's 
run-time state in its compile-time type
1. Operations on an object that are only available when the object
   is in certain states
2. A way of encoding these state at the type level, such that attempts
   to use the operations in the wrong state fail to compile
3. State transition operations that change the type-level state of 
   objects instead of chaning run-time dynamic state which make the
   operations in the previous state are no longer possible.

# Benefits
1. It moves certian types of error from run-time to compile time
2. It can eliminate run-time checks, make code faster/smaller

# Examples
1. Use state type as type parameter with bounded trait 
1. Define operation in each state
```rust
// S is the state parameter. We require it to impl
// our ResponseState trait (below) to prevent users
// from trying weird types like HttpResponse<u8>.
struct HttpResponse<S: ResponseState> {
    // This is the same field as in the previous example.
    state: Box<ActualResponseState>,
    // This reassures the compiler that the parameter
    // gets used.
    marker: std::marker::PhantomData<S>,
}

// State type options.
enum Start {} // expecting status line
enum Headers {} // expecting headers or body

trait ResponseState {}
impl ResponseState for Start {}
impl ResponseState for Headers {}

impl HttpResponse<Start> {
    fn new() -> Self {
        // ...
    }

    fn status_line(self, code: u8, message: &str)
        -> HttpResponse<Headers>
    {
        // ...
    }
}

/// Operations that are valid only in Headers state.
impl HttpResponse<Headers> {
    fn header(&mut self, key: &str, value: &str) {
        // ...
    }

    fn body(self, contents: &str) {
        // ...
    }
}
```
3. Define operation which is valid in every state
```rust
/// These operations are available in any state.
impl<S> HttpResponse<S> {
    fn bytes_so_far(&self) -> usize { /* ... */ }
}
```
4. State type could be modeled as struct when there are
special state at that state
```rust
// Similar to before:
struct HttpResponse<S: ResponseState> {
    // This is the same field as in the previous example.
    state: Box<ActualResponseState>,
    // Instead of PhantomData<S>, we store an actual copy.
    extra: S,
}

struct Headers {
    response_code: u8,
}

trait ResponseState {}
impl ResponseState for Start {}
impl ResponseState for Headers {}
```
