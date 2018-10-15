#[macro_use]
extern crate futures;

use futures::{Future, Poll};

pub trait Service {
    type Request;
    type Response;
    type Error;
    type Future: Future<Item = Self::Response, Error = Self::Error>;

    fn ready(self) -> Ready<Self> where Self: Sized {
        Ready { inner: Some(self) }
    }

    fn poll_ready(&mut self) -> Poll<(), Self::Error>;
    fn call(&mut self, req: Self::Request) -> Self::Future;
}

pub struct Ready<T> {
    inner: Option<T>,
}

impl<T: Service> Future for Ready<T> {
    type Item = T;
    type Error = T::Error;

    fn poll(&mut self) -> Poll<T, T::Error> {
        match self.inner {
            Some(ref mut service) => {
                let _ = try_ready!(service.poll_ready());
            }
            None => panic!("called `poll` after future completed"),
        }

        Ok(self.inner.take().unwrap().into())
    }
}

impl<'a, S: Service + 'a> Service for &'a mut S {
    type Request = S::Request;
    type Response = S::Response;
    type Error = S::Error;
    type Future = S::Future;

    fn poll_ready(&mut self) -> Poll<(), S::Error> {
        (**self).poll_ready()
    }

    fn call(&mut self, request: S::Request) -> S::Future {
        (**self).call(request)
    }
}

impl<S: Service + ?Sized> Service for Box<S> {
    type Request = S::Request;
    type Response = S::Response;
    type Error = S::Error;
    type Future = S::Future;

    fn poll_ready(&mut self) -> Poll<(), S::Error> {
        (**self).poll_ready()
    }

    fn call(&mut self, request: S::Request) -> S::Future {
        (**self).call(request)
    }
}
