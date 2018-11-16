use futures::{Future, Stream};
use tokio::executor;
use tokio::runtime::{self, Runtime, TaskExecutor};


pub fn start_with_num_threads<NH, A>(addr: A, new_handler: NH, threads: usize)
where
    NH: NewHandler + 'static,
    A: ToSocketAddrs + 'static,
{
    let runtime = new_runtime(threads);
    start_on_executor(addr, new_handler, runtime.executor());
    runtime.shutdown_on_idle().wait().unwrap();
}

fn new_runtime(threads: usize) -> Runtime {
    runtime::Builder::new()
        .core_threads(threads)
        .name_prefix("gotham-worker-")
        .build()
        .unwrap()
}

pub fn start_on_executor<NH, A>(addr: A, new_handler: NH, executor: TaskExecutor)
where
    NH: NewHandler + 'static,
    A: ToSocketAddrs + 'static,
{
    executor.spawn(init_server(addr, new_handler));
}

pub fn init_server<NH, A>(addr: A, new_handler: NH) -> impl Future<Item = (), Error = ()> {
    let (listener, addr) = tcp_listener(addr);
    let protocol = Arc::new(Http::new());
    let gotham_service = GothamService::new(new_handler);

    listener
        .incoming()
        .map_err(|e| panic!("socket error = {:?}", e))
        .for_each(move |socket| {
            let service = gotham_service.connect(socket.peer_addr().unwrap());
            let handler = protocol.serve_connection(socket, service).then(|_| Ok(()));

            executor::spawn(handler);

            Ok(())
        })
}
