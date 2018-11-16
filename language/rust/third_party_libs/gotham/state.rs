use std::any::{Any, TypeId};
use std::collections::HashMap;

//dynamic type hashmap
pub struct State {
    data: HashMap<TypeId, Box<Any + Send>>,
}

//Most types implement Any. However, any type which 
//contains a non-'static reference does not
pub trait StateData: Any + Send {}

impl State {
    pub(crate) fn new() -> State {
        State {
            data: HashMap::new(),
        }
    }

    pub fn put<T>(&mut self, t: T)
    where
        T: StateData,
    {
        let type_id = TypeId::of::<T>();
        trace!(" inserting record to state for type_id `{:?}`", type_id);
        self.data.insert(type_id, Box::new(t));
    }

    pub fn has<T>(&self) -> bool
    where
        T: StateData,
    {
        let type_id = TypeId::of::<T>();
        self.data.get(&type_id).is_some()
    }

    pub fn try_borrow<T>(&self) -> Option<&T>
    where
        T: StateData,
    {
        let type_id = TypeId::of::<T>();
        trace!(" borrowing state data for type_id `{:?}`", type_id);
        self.data.get(&type_id).and_then(|b| b.downcast_ref::<T>())
    }

    pub fn borrow<T>(&self) -> &T
    where
        T: StateData,
    {
        self.try_borrow()
            .expect("required type is not present in State container")
    }

    pub fn try_borrow_mut<T>(&mut self) -> Option<&mut T>
    where
        T: StateData,
    {
        let type_id = TypeId::of::<T>();
        trace!(" mutably borrowing state data for type_id `{:?}`", type_id);
        self.data
            .get_mut(&type_id)
            .and_then(|b| b.downcast_mut::<T>())
    }

    pub fn borrow_mut<T>(&mut self) -> &mut T
    where
        T: StateData,
    {
        self.try_borrow_mut()
            .expect("required type is not present in State container")
    }

    pub fn try_take<T>(&mut self) -> Option<T>
    where
        T: StateData,
    {
        let type_id = TypeId::of::<T>();
        self.data
            .remove(&type_id)
            .and_then(|b| b.downcast::<T>().ok())
            .map(|b| *b)
    }

    pub fn take<T>(&mut self) -> T
    where
        T: StateData,
    {
        self.try_take()
            .expect("required type is not present in State container")
    }
}
