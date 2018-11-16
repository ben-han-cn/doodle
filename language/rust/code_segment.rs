fn insert(&mut self, query: Query, rdatas_and_ttl: Vec<(RData, u32)>, now: Instant-> Lookup {
    let len = rdatas_and_ttl.len();
    let (rdatas, ttl): (Vec<RData>, Duration) = rdatas_and_ttl.into_iter().fold(
        (Vec::with_capacity(len), self.positive_max_ttl),
        |(mut rdatas, mut min_ttl), (rdata, ttl)| {
            rdatas.push(rdata);
            let ttl = Duration::from_secs(ttl as u64);
            min_ttl = min_ttl.min(ttl);
            (rdatas, min_ttl)
        },  
        );  
} 


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

    fn capacity(&self) -> usize {
        self.data.capacity()
    }

    fn is_full(&self) -> bool {
        self.size == self.data.capacity()
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
