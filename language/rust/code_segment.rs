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

    let ttl = self.positive_min_ttl.max(ttl);
    let valid_until = now + ttl;
    let lookup = Lookup::new_with_deadline(Arc::new(rdatas), valid_until);
    self.cache.insert(
        query,
        LruValue {
            lookup: Some(lookup.clone()),
            valid_until,
        },  
        );  

    lookup
} 
