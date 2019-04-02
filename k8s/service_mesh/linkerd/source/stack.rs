pub trait Stack<T> {
    type Error;
    type Value;
    fn make(&self, target: &T) -> Result<Self::Value, Self::Error>;

    fn push<U, L>(self, layer: L) -> L::Stack 
    where 
        L: Layer<U, T, Self>,
        Self: Sized,
    {
        layer.bind(self)
    }
}
//stack :: T -> Value
//push :: Stack<T> -> Layer<U, T, Self> -> Stack<U>

pub trait Layer<T, U, S: super::Stack<U>> {
    type Value;
    type Error;

    fn bind(&self, next: S) -> Stack<T, Value = Self::Value, Error = Self::Error>;
}
//bind :: Stack<U> -> Stack<T>; change the param, but the return value of stack isn't changed
//layer and stack is modeling a pipeline, each layer works on last layer's result


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

