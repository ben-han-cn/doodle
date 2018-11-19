pub trait Future {
    type Item;
    type Error;
    fn poll(&mut self) -> Poll<Self::Item, Self::Error>;

    fn wait(self) -> result::Result<Self::Item, Self::Error> {
        ::executor::spawn(self).wait_future()
    }
}

pub type Poll<T, E> = Result<Async<T>, E>;
#[derive(Copy, Clone, Debug, PartialEq)]
pub enum Async<T> {
    Ready(T),
    NotReady,
}

//it's Spawn not Task the unit to schedule/execute by executor
//task is just a handler(TaskUnpark) which could be used to notify
//some progress could be made which will give the future which resides
//in Spawn could be polled again
#[derive(Clone)]
pub struct Task {
    id: usize,
    unpark: TaskUnpark,
    events: UnparkEvents,
}

pub struct TaskUnpark {
    handle: NotifyHandle,
    id: usize,
}

pub struct NotifyHandle {
    inner: *mut UnsafeNotify,
}


pub unsafe trait UnsafeNotify: Notify {
    unsafe fn clone_raw(&self) -> NotifyHandle;
    unsafe fn drop_raw(&self);
}

pub trait Notify: Send + Sync {
    fn notify(&self, id: usize);

    fn clone_id(&self, id: usize) -> usize {
        id
    }

    fn drop_id(&self, id: usize) {
        drop(id);
    }
}


#[derive(Clone)]
pub struct UnparkEvents;

/////////////////////////////////////////////
pub struct BorrowedTask<'a> {
    id: usize,
    unpark: BorrowedUnpark<'a>,
    events: BorrowedEvents<'a>,
    map: &'a LocalMap, //local stoage
}

#[derive(Copy, Clone)]
pub struct BorrowedEvents<'a>(marker::PhantomData<&'a ()>);

#[derive(Copy, Clone)]
pub struct BorrowedUnpark<'a> {
    f: &'a Fn() -> NotifyHandle,
    id: usize,
}


pub struct UnparkEvents;

impl<'a> BorrowedEvents<'a> {
    pub fn new() -> BorrowedEvents<'a> {
        BorrowedEvents(marker::PhantomData)
    }

    pub fn to_owned(&self) -> UnparkEvents {
        UnparkEvents
    }
}

impl<'a> BorrowedUnpark<'a> {
    #[inline]
    pub fn new(f: &'a Fn() -> NotifyHandle, id: usize) -> BorrowedUnpark<'a> {
        BorrowedUnpark { f: f, id: id }
    }

    #[inline]
    pub fn to_owned(&self) -> TaskUnpark {
        let handle = (self.f)();
        let id = handle.clone_id(self.id);
        TaskUnpark { handle: handle, id: id }
    }
}

pub fn current() -> Task {
    with(|borrowed| {
        let unpark = borrowed.unpark.to_owned();
        let events = borrowed.events.to_owned();

        Task {
            id: borrowed.id,
            unpark: unpark,
            events: events,
        }
    })
}
