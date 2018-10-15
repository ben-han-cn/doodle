//interface:
pub struct Container<T>

impl<T> Default for Container<T> {}

impl<T> Container<T> {
    //constructor
    pub fn new() -> Container<T> {}
    pub fn with_capacity(capacity: usize) -> Container<T> {}

    pub fn clear(&mut self) {}
    pub fn len(&self) -> usize {}

    pub fn iter(&self) -> Iter<T> {}
    pub fn iter_mut(&mut self) -> IterMut<T> {}
    pub fn get(&self, key: usize) -> Option<&T> {}
    pub fn get_mut(&mut self, key: usize) -> Option<&mut T> {}
    pub fn insert(&mut self, val: T) -> usize {}
    pub fn remove(&mut self, key: usize) -> T {}
    pub fn contains(&self, key: usize) -> bool {}

    pub fn retain<F>(&mut self, mut f: F) 
        where F: FnMut(usize, &mut T) -> bool {
    }
}


impl<T> ops::Index<usize> for Container<T> {
    type Output = T;
    fn index(&self, key: usize) -> &T {
    }
}

impl<T> ops::IndexMut<usize> for Container<T> {
    fn index_mut(&mut self, key: usize) -> &mut T {
    }
}

impl<'a, T> IntoIterator for &'a Container<T> {
    type Item = (usize, &'a T);
    type IntoIter = Iter<'a, T>;
    fn into_iter(self) -> Iter<'a, T> {
    }
}

impl<'a, T> IntoIterator for &'a mut Container<T> {
    type Item = (usize, &'a mut T);
    type IntoIter = IterMut<'a, T>;
    fn into_iter(self) -> IterMut<'a, T> {
    }
}

impl<T: fmt::Debug> fmt::Debug for Container<T> {
    fn fmt(&self, fmt: &mut fmt::Formatter) -> fmt::Result {
    }
}

impl<'a, T> Iterator for Iter<'a, T> {
    type Item = (usize, &'a T);
    fn next(&mut self) -> Option<(usize, &'a T)> {
    }
}

impl<'a, T> Iterator for IterMut<'a, T> {
    type Item = (usize, &'a mut T);
    fn next(&mut self) -> Option<(usize, &'a mut T)> {
    }
}


//optional
impl <T:Clone> Clone for Contaienr<T> {
    fn clone(&self) -> Contaienr<T> {
    }
}

impl AsRef<[u8]> for Contaienr<T> {
    #[inline]
    fn as_ref(&self) -> &[u8] {
    }
}

//this is normally for buffer, which use slice 
//as underlaying data strucutre
impl ops::Deref for Contaienr {
    type Target = [u8];

    #[inline]
    fn deref(&self) -> &[u8] {
    }
}

impl From<Vec<u8>> for Container {
    fn from(src: Vec<u8>) -> Bytes {
        BytesMut::from(src).freeze()
    }
}

impl <T> Borrow<[T]> for Container<T> {
    fn borrow(&self) -> &[T] {
    }
}

impl<K, V> DoubleEndedIterator for IntoIter<K, V> {
    fn next_back(&mut self) -> Option<(K, V)> {
    }
}

impl<K, V> ExactSizeIterator for IntoIter<K, V> {
    fn len(&self) -> usize {
    }
}
