# Implicit task
```rust
fn poll(&mut self) -> Poll<Self::Item, Self::Error>;
```
poll takes no arguments, and returns either Ready, NotReady, or an error. 
From the signature, it looks as though the only way to complete a Future 
is to continually call poll in a loop until it returns Ready. This is 
incorrect.

Furthermore, when implementing Future, it's necessary to schedule a wakeup 
when you return NotReady, or else you might never be polled again. The trick 
to the API is that there is an implicit Task argument to Future::poll which 
is passed using thread-local storage. This Task argument must be notifyd by 
the Future in order for the task to awaken and poll again. It's essential to 
remember to do this, and to ensure that the right schedulings occur on every 
code path.

# Reasoning about code
The implicit Task argument makes it difficult to tell which functions aside 
from poll will use the Task handle to schedule a wakeup. That has a few 
implications.

First, it's easy to accidentally call a function that will attempt to access 
the task outside of a task context. Doing so will result in a panic, but it 
would be better to detect this mistake statically.

Second, and relatedly, it is hard to audit a piece of code for which calls might 
involve scheduling a wakeup--something that is critical to get right in order 
to avoid "lost wakeups", which are far harder to debug than an explicit panic.
