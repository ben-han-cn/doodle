/*
build the composite future, use tokio::run to run it
this future has no difference from runtime spawned future
run will block the current thread(main thread)

worker runs in each thread which composit the thread pool
worker has a deque to store all the task

task is a wrapper of boxed future

timer is a seperate thread which get a ref of current task
task will be notified when deadline or interval arrived

reactor
Establishing a TCP connection usually cannot be completed immediately.
[`TcpStream::connect`] does not block the current thread. Instead, it
returns a [future][connect-future] that resolves once the TCP connection has
been established. The connect future itself has no way of knowing when the
TCP connection has been established.

Before returning the future, [`TcpStream::connect`] registers the socket
with a reactor. This registration process, handled by [`Registration`], is
what links the [`TcpStream`] with the [`Reactor`] instance. At this point,
the reactor starts listening for connection events from the operating system
for that socket.

Once the connect future is passed to [`tokio::run`], it is spawned onto a
thread pool. The thread pool waits until it is notified that the connection
has completed.

When the TCP connection is established, the reactor receives an event from
the operating system. It then notifies the thread pool, telling it that the
connect future can complete. At this point, the thread pool will schedule
the task to run on one of its worker threads. This results in the `and_then`
closure to get executed.

 */
pub fn run<F>(future: F)
where
    F: Future<Item = (), Error = ()> + Send + 'static,
{
    let mut entered = enter().expect("nested tokio::run");
    let mut runtime = Runtime::new().expect("failed to start new Runtime");
    runtime.spawn(future);
    entered
        .block_on(runtime.shutdown_on_idle())
        .expect("shutdown cannot error")
}

pub struct Runtime {
    inner: Option<Inner>,
}

#[derive(Debug)]
struct Inner {
    reactor: Handle,
    pool: threadpool::ThreadPool,
}

impl Runtime {
    pub fn spawn<F>(&mut self, future: F) -> &mut Self
    where
        F: Future<Item = (), Error = ()> + Send + 'static,
    {
        //== self.inner_mut().pool.spawn(future)
        self.inner_mut().pool.sender().spawn(future).unwrap();
        self
    }
}

impl Enter {
    pub fn block_on<F: Future>(&mut self, f: F) -> Result<F::Item, F::Error> {
        futures::executor::spawn(f).wait_future()
    }
}

pub struct ThreadPool {
    pub(crate) inner: Option<Sender>,
}

impl ThreadPool {
    pub fn spawn<F>(&self, future: F)
    where
        F: Future<Item = (), Error = ()> + Send + 'static,
    {
        self.sender().spawn(future).unwrap();
    }
}

pub struct Sender {
    pub(crate) inner: Arc<Pool>,
}

impl<'a> tokio_executor::Executor for &'a Sender {
    fn spawn(
        &mut self,
        future: Box<Future<Item = (), Error = ()> + Send>,
    ) -> Result<(), SpawnError> {
        //check pool is over use
        self.prepare_for_spawn()?;
        let task = Arc::new(Task::new(future));
        self.inner.submit_to_random(task, &self.inner);
        Ok(())
    }
}

//spawn new task in future
//each thread has a default executor, which should be set by invoke
//with_default, threadpool invoke this function in do_run in file worker/mod.rs
use tokio::executor;
pub fn spawn<T>(future: T)
where
    T: Future<Item = (), Error = ()> + Send + 'static,
{
    DefaultExecutor::current().spawn(Box::new(future)).unwrap()
}
