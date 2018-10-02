//dynamic array but remove won't free elem
//use next and vacant elem to build a free list

pub struct Slab<T> {
    entries: Vec<Entry<T>>,
    len: usize,
    // next available slot in the slab. Set to the slab's
    // capacity when the slab is full.
    next: usize,
}

pub fn insert(&mut self, val: T) -> usize {
    let key = self.next;
    self.insert_at(key, val);
    key
}

fn insert_at(&mut self, key: usize, val: T) {
    self.len += 1;

    if key == self.entries.len() {
        self.entries.push(Entry::Occupied(val));
        self.next = key + 1;
    } else {
        let prev = mem::replace(
            &mut self.entries[key],
            Entry::Occupied(val));
        match prev {
            Entry::Vacant(next) => {
                self.next = next;
            }
            _ => unreachable!(),
        }
    }
}


pub fn remove(&mut self, key: usize) -> T {
    // Swap the entry at the provided value
    let prev = mem::replace(
        &mut self.entries[key],
        Entry::Vacant(self.next));

    match prev {
        Entry::Occupied(val) => {
            self.len -= 1;
            self.next = key;
            val
        }
        _ => {
            // Woops, the entry is actually vacant, restore the state
            self.entries[key] = prev;
            panic!("invalid key");
        }
    }
}
