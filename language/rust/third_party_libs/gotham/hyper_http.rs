

pub fn serve_connection<S, I, Bd>(&self, io: I, service: S) -> Connection<I, S>
    where
        S: Service<ReqBody=Body, ResBody=Bd>,
        S::Error: Into<Box<::std::error::Error + Send + Sync>>,
        S::Future: Send + 'static,
        Bd: Payload,
        I: AsyncRead + AsyncWrite,
    {
        let either = if !self.http2 {
            let mut conn = proto::Conn::new(io);
            if !self.keep_alive {
                conn.disable_keep_alive();
            }
            conn.set_flush_pipeline(self.pipeline_flush);
            if let Some(max) = self.max_buf_size {
                conn.set_max_buf_size(max);
            }
            let sd = proto::h1::dispatch::Server::new(service);
            Either::A(proto::h1::Dispatcher::new(sd, conn))
        } else {
            let rewind_io = Rewind::new(io);
            let h2 = proto::h2::Server::new(rewind_io, service, self.exec.clone());
            Either::B(h2)
        };

        Connection {
            conn: Some(either),
        }
    }


pub trait Service {
    type ReqBody: Payload;
    type ResBody: Payload;
    type Error: Into<Box<StdError + Send + Sync>>;
    type Future: Future<Item=Response<Self::ResBody>, Error=Self::Error>;
    fn call(&mut self, req: Request<Self::ReqBody>) -> Self::Future;
}
