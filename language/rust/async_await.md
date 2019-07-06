# async keyword
```rust
async fn read_and_input() -> Result<(), ::std::io::Error> {}
fn read_and_input() -> impl Future<Output=Result<(), ::std::io::Error>>
```
