/*
  create executor 
  spwan/create several future
  compostion future and wait for it to finish
*/

//How future is transformed to task? 

let pool = CpuPool::new(4);
let a = pool.spawn(LongRunningFuture::new(time::Duration::from_secs(1), 2));
let b = pool.spawn(LongRunningFuture::new(time::Duration::from_secs(3), 2));
let c = a.join(b).map(|(v1, v2)| v1 + v2).wait().unwrap();



//CpuPool::spawn(future) --> CpuFuture
// create a channel
// create inner task with the tx side, poll the future, if ready, send data to tx
// create CpuFuture with rx side, CpuFuture is another future which keep poll data from rx
//
//
//composit future.wait
//block current thread until the future return Ready
//
//
//thread pool in CpuPool
//each thread will get a Run to execute, run will be invoked
//run will call inner future's poll, if future isn't ready, run will exit
//when the future is notified, the run will be re-execute by executor

impl CpuPool {
    pub fn spawn<F>(&self, f: F) -> CpuFuture<F::Item, F::Error>
        where F: Future + Send + 'static,
              F::Item: Send + 'static,
              F::Error: Send + 'static,
              {   
                  let (tx, rx) = channel();
                  let keep_running_flag = Arc::new(AtomicBool::new(false));
                  let sender = MySender {
                      fut: AssertUnwindSafe(f).catch_unwind(),
                      tx: Some(tx),
                      keep_running_flag: keep_running_flag.clone(),
                  };  
                  executor::spawn(sender).execute(self.inner.clone());
                  CpuFuture { inner: rx , keep_running_flag: keep_running_flag.clone() }
              } 
}

pub struct Spawn<T: ?Sized> {
    id: usize,
    data: LocalMap,
    obj: T,
}

pub fn spawn<T>(obj: T) -> Spawn<T> {
    Spawn {
        id: fresh_task_id(),
        obj: obj,
        data: local_map(),
    }   
}

impl<F: Future> Spawn<F> {
    pub fn execute(self, exec: Arc<Executor>)
        where F: Future<Item=(), Error=()> + Send + 'static,
    {
        exec.clone().execute(Run {
            spawn: spawn(Box::new(self.into_inner())),
            inner: Arc::new(RunInner {
                exec: exec,
                mutex: UnparkMutex::new()
            }),
        })
    }
}

pub struct Run {
    spawn: Spawn<Box<Future<Item = (), Error = ()> + Send>>,
    inner: Arc<RunInner>,
}

struct RunInner {
    mutex: UnparkMutex<Run>,
    exec: Arc<Executor>,
}

impl OldExecutor for Inner {
    fn execute(&self, run: Run) {
        self.send(Message::Run(run))
    }
}

impl Inner {
    fn send(&self, msg: Message) {
        self.tx.lock().unwrap().send(msg).unwrap();
    }

    fn work(&self, after_start: Option<Arc<Fn() + Send + Sync>>, before_stop: Option<Arc<Fn() + Send + Sync>>) {
        after_start.map(|fun| fun());
        loop {
            let msg = self.rx.lock().unwrap().recv().unwrap();
            match msg {
                Message::Run(r) => r.run(),
                Message::Close => break,
            }
        }
        before_stop.map(|fun| fun());
    }
}

impl Run {
    pub fn run(self) {
        let Run { mut spawn, inner } = self;

        unsafe {
            inner.mutex.start_poll();
            loop {
                match spawn.poll_future_notify(&inner, 0) {
                    Ok(Async::NotReady) => {}
                    Ok(Async::Ready(())) |
                    Err(()) => return inner.mutex.complete(),
                }
                let run = Run { spawn: spawn, inner: inner.clone() };
                match inner.mutex.wait(run) {
                    Ok(()) => return,            
                    Err(r) => spawn = r.spawn,  
                }
            }
        }
    }
}

impl Notify for RunInner {
    fn notify(&self, _id: usize) {
        match self.mutex.notify() {
            Ok(run) => self.exec.execute(run),
            Err(()) => {}
        }
    }
}

impl<T: ?Sized> Spawn<T> {
    pub fn poll_fn_notify<N, F, R>(&mut self,
                             notify: &N,
                             id: usize,
                             f: F) -> R
        where F: FnOnce(&mut T) -> R,
              N: Clone + Into<NotifyHandle>,
    {
        let mk = || notify.clone().into();
        //enter make current available to f
        self.enter(BorrowedUnpark::new(&mk, id), f)
    }

    pub fn poll_future_notify<N>(&mut self,
                                 notify: &N,
                                 id: usize) -> Poll<T::Item, T::Error>
        where N: Clone + Into<NotifyHandle>,
              T: Future,
    {
        self.poll_fn_notify(notify, id, |f| f.poll())
    }

    fn enter<F, R>(&mut self, unpark: BorrowedUnpark, f: F) -> R
        where F: FnOnce(&mut T) -> R
    {
        let borrowed = BorrowedTask {
            id: self.id,
            unpark: unpark,
            events: BorrowedEvents::new(),
            map: &self.data,
        };
        let obj = &mut self.obj;
        set(&borrowed, || f(obj))
    }
}

pub trait Future {
    fn wait(self) -> result::Result<Self::Item, Self::Error>
        where Self: Sized
        {
            ::executor::spawn(self).wait_future()
        }
}

impl<F: Future> Spawn<F> {
    pub fn wait_future(&mut self) -> Result<F::Item, F::Error> {
        ThreadNotify::with_current(|notify| {
            loop {
                match self.poll_future_notify(notify, 0)? {
                    Async::NotReady => notify.park(),
                    Async::Ready(e) => return Ok(e),
                }
            }
        })
    }
}
