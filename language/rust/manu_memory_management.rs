enum Entry<T> {
    Free { next_free: Option<usize> },
    Occupied { generation: u64, value: T},
}


struct GenerationArena<T> {
    items: Vec<Entry<T>>,
    free_list_header: Option<usize>,
    generation: u64,
}

//delete elem will mark it as free, but the mem is keeped
//each delete will set the generation one more than the delete target
//after delete the elem, and insert new one in the same index, the new entry 
//will have new generation
//search will needs index and generation, generation which smaller than current
//generation means it doesn't exists any more.
//
//

//GenerationArena is wasteful if too many empty entry in the vec
struct DenseVecStorage<T> {
    lookup: Vec<u32>,
    data: Vec<T>,
}
