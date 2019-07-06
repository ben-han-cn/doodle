# Iterator
Consume object during traversing, return None when
it's exhausted
```rust
trait Iterator {
    type Item;
    fn next(&mut self) -> Option<Self::Item>;

    //other default impl
}

trait DoubleEndedIterator: Iterator {
    fn next_back(&mut self) -> Option<Self::Item>;
    //...
}

trait ExactSizeIterator: Iterator {
    fn len(&self) -> usize { ... }
    fn is_empty(&self) -> bool { ... }
}
```

# Container
```rust
pub trait IntoIterator {
    type Item;
    type IntoIter: Iterator<Item = Self::Item>;
    fn into_iter(self) -> Self::IntoIter;
}

pub trait FromIterator<A> {
    fn from_iter<T>(iter: T) -> Self
    where
        T: IntoIterator<Item = A>;
}
```

Iterator self is also implement IntoIterator
```rust
impl<I: Iterator> IntoIterator for I {
    type Item = I::Item;
    type IntoIter = I;

    fn into_iter(self) -> I {
        self
    }
}
```

# Iterator adaptor
Take an Iterator return another Iterator
```rust
pub struct Double<I>(I);

impl<I: Iterator<Item = u32>> Iterator for Double<I> {
    type Item = u32;
    fn next(&mut self) -> Option<u32> {
        self.0.next().map(|x| x * 2)
    }   
}

fn double<I: IntoIterator<Item = u32>>(iter: I) -> Double<I::IntoIter> {
    Double(iter.into_iter())
}
```

# Syntax sugar
