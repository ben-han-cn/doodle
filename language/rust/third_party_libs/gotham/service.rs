pub(crate) struct GothamService<T>
where
    T: NewHandler + 'static,
{
    handler: Arc<T>,
}

impl<T> GothamService<T>
where
    T: NewHandler + 'static,
{
    pub(crate) fn new(handler: T) -> GothamService<T> {
        GothamService {
            handler: Arc::new(handler),
        }
    }

    pub(crate) fn connect(&self, client_addr: SocketAddr) -> ConnectedGothamService<T> {
        ConnectedGothamService {
            client_addr,
            handler: self.handler.clone(),
        }
    }
}

//ConnectedGothamService make sure connect is invoked on GothamService 
pub(crate) struct ConnectedGothamService<T>
where
    T: NewHandler + 'static,
{
    handler: Arc<T>,
    client_addr: SocketAddr,
}

impl<T> Service for ConnectedGothamService<T>
where
    T: NewHandler,
{
    type ReqBody = Body; // required by hyper::server::conn::Http::serve_connection()
    type ResBody = Body; // has to impl Payload...
    type Error = failure::Compat<failure::Error>; // :Into<Box<StdError + Send + Sync>>
    type Future = Box<Future<Item = Response<Self::ResBody>, Error = Self::Error> + Send>;

    fn call(&mut self, req: Request<Self::ReqBody>) -> Self::Future {
        let mut state = State::new();
        put_client_addr(&mut state, self.client_addr);
        let (
            request::Parts {
                method,
                uri,
                version,
                headers,
                //extensions?
                ..  
            },  
            body,
        ) = req.into_parts();
        state.put(RequestPathSegments::new(uri.path()));
        state.put(method);
        state.put(uri);
        state.put(version);
        state.put(headers);
        state.put(body);

        {   
            let request_id = set_request_id(&mut state);
            debug!(
                "[DEBUG][{}][Thread][{:?}]",
                request_id,
                thread::current().id(),
            );  
        };  

        trap::call_handler(&*self.handler, AssertUnwindSafe(state))
    }   
}


//catch_unwind a way to handle panic
pub(super) fn call_handler<'a, T>(
    t: &T,
    state: AssertUnwindSafe<State>,
) -> Box<Future<Item = Response<Body>, Error = CompatError> + Send + 'a>
where
    T: NewHandler + 'a,
{
    let res = catch_unwind(move || {
        // Hyper doesn't allow us to present an affine-typed `Handler` interface directly. We have
        // to emulate the promise given by hyper's documentation, by creating a `Handler` value and
        // immediately consuming it.
        t.new_handler()
            .into_future()
            .map_err(|e| failure::Error::from(e).compat())
            .and_then(move |handler| {
                let AssertUnwindSafe(state) = state;

                handler.handle(state).then(move |result| match result {
                    Ok((_state, res)) => future::ok(res),
                    Err((state, err)) => finalize_error_response(state, err),
                })
            })
    });

    if let Ok(f) = res {
        return Box::new(
            UnwindSafeFuture::new(f)
                .catch_unwind()
                .then(finalize_catch_unwind_response), // must be Future<Item = impl Payload>
        );
    }

    Box::new(finalize_panic_response())
}
