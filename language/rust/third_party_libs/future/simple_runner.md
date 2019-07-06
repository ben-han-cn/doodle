#Simple Runner
```rust
use std::time::Duration;
use std::thread;
use std::cell::RefCell;

thread_local!(static NOTIFY: RefCell<bool> = RefCell::new(true));

mod task {
    use crate::NOTIFY;
    pub struct Task();
    
    impl Task {
        pub fn notify(&self) {
            NOTIFY.with(|f| {
                *f.borrow_mut() = true
            })         
        }
    }
    
    pub fn current() -> Task {
        Task()
    }
}

fn run<F>(mut f: F)
where
    F: Future<Item = (), Error = ()>,
{
    loop {
        if NOTIFY.with(|n| {
            if *n.borrow() {
                *n.borrow_mut() = false;
                match f.poll() {
                    Ok(Async::Ready(_)) | Err(_) => return true,
                    Ok(Async::NotReady) => (),
                }
            }
            thread::sleep(Duration::from_millis(100));
            false
        }) { break }
    }
}
```
