#Architecutre
1. Futures
A computation maybe blocked, and will be finished in the futrue.
```rust
trait Future {
    type Output;
    fn poll(self: Pin<&mut Self>, cx: &mut Context<'_>) -> Poll<Self::Output>;
}
```

2. Task
A future with extra info, which will be scheduled by executor.
It will be waked when the event the future waiting for is arrived.
Then it will make itself rescheduled by executor. Therefore, like spawner,
task knows how to send a future to the executor.
```
struct Task {
    /// In-progress future that should be pushed to completion.
    ///
    /// The `Mutex` is not necessary for correctness, since we only have
    /// one thread executing tasks at once. However, Rust isn't smart
    /// enough to know that `future` is only mutated from one thread,
    /// so we need use the `Mutex` to prove thread-safety. A production
    /// executor wouild not need this, and could use `UnsafeCell` instead.
    future: Mutex<Option<BoxFuture<'static, ()>>>,

    /// Handle to place the task itself back onto the task queue.
    task_sender: SyncSender<Arc<Task>>,
} 
```

3. Spawner
A middle man who knows where is executor, and send task to him
```rust
struct Spawner {
    task_sender: SyncSender<Arc<Task>>,
}

impl Spawner {
    fn spawn(&self, future: impl Future<Output = ()> + 'static + Send) {
        let future = future.boxed();
        let task = Arc::new(Task {
            future: Mutex::new(Some(future)),
            task_sender: self.task_sender.clone(),
        });
        self.task_sender.send(task).expect("too many tasks queued");
    }
}
```
4. Executor
Save all the scheduleable task in a queue. 
Pop task from ready queue, and invoke poll of the future belongs to the task.
If future return not ready, it will be store back to the task. A waker is created
by executor and put it into context, which could be feteched in poll function, 
future could register the waker to reactor, reactor will wake the waker, the waker
will notify the task, and task will send itself to the executor and to be polled 
again.
```
struct Executor {
    ready_queue: Receiver<Arc<Task>>,
}

impl Executor {
    fn run(&self) {
        while let Ok(task) = self.ready_queue.recv() {
            // Take the future, and if it has not yet completed (is still Some),
            // poll it in an attempt to complete it.
            let mut future_slot = task.future.lock().unwrap();
            if let Some(mut future) = future_slot.take() {
                // Create a `LocalWaker` from the task itself
                let waker = waker_ref(&task);
                let context = &mut Context::from_waker(&*waker);
                // `BoxFuture<T>` is a type alias for
                // `Pin<Box<dyn Future<Output = T> + Send + 'static>>`.
                // We can get a `Pin<&mut dyn Future + Send + 'static>`
                // from it by calling the `Pin::as_mut` method.
                if let Poll::Pending = future.as_mut().poll(context) {
                    // We're not done processing the future, so put it
                    // back in its task to be run again in the future.
                    *future_slot = Some(future);
                }
            }
        }
    }
}
```
5. Reactor
Register an event with a waker, and when the event arrive, wake the waker.

6. Stream
Similar with future, but could generate several value before completing. Stream like a
async iterator.
```
trait Stream {
    /// The type of the value yielded by the stream.
    type Item;

    /// Attempt to resolve the next item in the stream.
    /// Retuns `Poll::Pending` if not ready, `Poll::Ready(Some(x))` if a value
    /// is ready, and `Poll::Ready(None)` if the stream has completed.
    fn poll_next(self: Pin<&mut Self>, cx: &mut Context<'_>) -> Poll<Option<Self::Item>>;
}
```

#async/await
1. async
There are two kinds of async, both of them when execute will generate a type 
implement future trait
  * async fn 
  * async block
2. .await
await like a method of a future, it will schedule the future and wait for it
to complete.
3. lifetime
```rust
// This function:
async fn foo(x: &u8) -> u8 { *x }

// Is equivalent to this function:
fn foo_expanded<'a>(x: &'a u8) -> impl Future<Output = u8> + 'a {
    async move { *x }
}
```
async block move the argument into the future, extends its lifetime to math
the future returned, but in the sametime, avoid to share the variable.
```
async fn borrow_x(x: &u8) -> u8 {}
fn foo() -> impl Future<Output=u8> {
    async {
        let x = 5;
        borrow_x(&x).await
    }
}
```
async move does the same thing.
