/*
poll instead of push, make backpress handling more free, since only when work is done, 
future will read the data again, instead of driver read too many data, and future cann't 
handle it, then the data has to be discard.
*/


//future call UdpSocket poll_recv
//PollEvented return not ready and store current task associate with the socket read status
//epoll notified by os that the socket is ready
//reactor get the task and notify the task

pub fn poll_recv(&mut self, buf: &mut [u8]) -> Poll<usize, io::Error> {
    try_ready!(self.io.poll_read_ready(mio::Ready::readable()));
    //
}

//struct PollEvented
impl PollEvented {
    pub fn poll_read_ready(&self, mask: mio::Ready) -> Poll<mio::Ready, io::Error> {
        assert!(!mask.is_writable(), "cannot poll for write readiness");
        poll_ready!(
            self, mask, read_readiness, take_read_ready,
            self.inner.registration.poll_read_ready()
            )
    }
}

impl Registration {
    pub fn poll_read_ready(&self) -> Poll<mio::Ready, io::Error> {
        self.poll_ready(Direction::Read, true, || task::current())
            .map(|v| match v {
                Some(v) => Async::Ready(v),
                _ => Async::NotReady,
            })
    }
}

impl Reactor {
    fn dispatch(&self, token: mio::Token, ready: mio::Ready) {
        //....
        task.notify();
    }
}
