struct Ring {
    size: usize,
    data: Vec<Option<u32>>,
}

impl Ring {
    fn with_capacity(cap: usize) -> Self {
        Ring {
            size: 0,
            data: vec![None; cap],
        }
    }

    fn emplace(&mut self, offset: usize, val: u32) -> Option<u32> {
        self.size += 1;
        mem::replace(&mut self.data[offset], Some(val))
    }

    fn dispalce(&mut self, offset: usize) -> Option<u32> {
        let res = mem::replace(&mut self.data[offset], None);
        if res.is_some() {
            self.size -= 1;
        }
        res
    }
}
////////////////////////////////////////////////////////////
pub struct RlpStream {}
impl RlpStream {
    pub fn append<E: Encodable>(&mut self, value: &E) {}

    pub fn append_list<E, K>(&mut self, values: &[K])
    where
        E: Encodable,
        K: Borrow<E>,
    {
        for value in values {
            self.append(value.borrow());
        }
    }
}

////////////////////////////////////////////////////////////
pub struct NibbleSlice<'a> {
    pub data: &'a [u8],
    pub offset: usize,
}

impl<'a, 'view> NibbleSlice<'a>
where
    'a: 'view,
{
    pub fn new(bytes: &'a [u8]) -> NibbleSlice<'a> {}

    pub fn mid(&'view self, i: usize) -> NibbleSlice<'a> {
        NibbleSlice {
            data: self.data,
            offset: self.offset + i,
        }
    }
}
////////////////////////////////////////////////////////////
pub trait Query {
    type Item;
    fn decode(self, &[u8]) -> Self::Item;
}

impl<F, T> Query for F
where
    F: FnOnce(&[u8]) -> T,
{
    type Item = T;
    fn decode(self, value: &[u8]) -> Self::Item {
        (self)(value)
    }
}

pub struct TrieDB<'a> {
    db: &'a HashDB,
}

//use query to avoid return inner data
impl<'a> TrieDB<'a> {
    fn get_aux<Q: Query>(&self, path: &NibbleSlice, query: Q) -> Result<Option<Q::Item>> {
        let node_rpl = self.db.get(&hash);
        Ok(Some(query.decode(value)))
    }
}
