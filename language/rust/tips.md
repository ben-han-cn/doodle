# struct literal operator ".."
Create a copy of a struct with one or more fields different.
```rust
#[derive(Default)]
sturct Foo {
    x: i32,
    y: i32,
}

let a = Foo{x: 1, y: 2};
let b = Foo{x:2, ..a};
let c = Foo{x:2, ..Default::default()};
```
# padding format operator
 * :> left-pad
 * :< right-pad
 * :^ center-pad
```rust
let title = "SCORES";

let player1 = "first player:";
let player2 = "second player:";
let player3 = "third player:";

let score1 = 100;
let score2 = 1000;
let score3 = 10000;

println!("{:_^20}", title);
println!("{:<14} {:>5}", player1, score1);
println!("{:<14} {:>5}", player2, score2);
println!("{:<14} {:>5}", player3, score3);
```
# switch to nightly version 
```shell
rustup update nightly
rustup override set nightly
```
