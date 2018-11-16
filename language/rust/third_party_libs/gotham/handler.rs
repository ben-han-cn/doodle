pub type HandlerFuture =
    Future<Item = (State, Response<Body>), Error = (State, HandlerError)> + Send;

pub trait Handler: Send {
    fn handle(self, state: State) -> Box<HandlerFuture>;
}


pub trait NewHandler: Send + Sync + RefUnwindSafe {
    type Instance: Handler + Send;
    fn new_handler(&self) -> Result<Self::Instance>;
}
