pub trait Stack<T> {
    type Value;
    fn make(&self, target: T) -> Result<Self::Value, Self::Error>;
}
//stack :: T -> Value

pub trait Layer<T, U, S: super::Stack<U>> {
    type Value;

    fn bind(&self, next: S) -> Stack<T, Value = Self::Value, Error = Self::Error>;
}
//bind :: Stack<U> -> Stack<T>; change the param, but the return value of stack isn't changed


pub trait Recognize<Request> {
    type Target: Clone + Eq + Hash;
    fn recognize(&self, req: &Request) -> Option<Self::Target>;
}

//Rec::recognize: Request -> Target
//stack::make:    Target -> Service<Request>
//service::call   Request -> Response
pub struct Router<Req, Rec, Stk>
where
    Rec: Recognize<Req>,
    Stk: stack::Stack<Rec::Target>,
    Stk::Value: svc::Service<Req>,
{
    inner: Arc<Inner<Req, Rec, Stk>>,
}

