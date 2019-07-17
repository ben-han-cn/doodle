# Kind 
## Declarative macros
``` rust
macro_rules! $name {
    $rule0 ;
    $rule1 ;
    // …
    $ruleN ;
}
rule: ($pattern) => {$expansion}
```
1. Capture:
* item: an item, like a function, struct, module, etc.
* block: a block (i.e. a block of statements and/or an expression, surrounded by braces)
* stmt: a statement
* pat: a pattern
* expr: an expression
* ty: a type
* ident: an identifier
* path: a path 
* meta: a meta item; the things that go inside #[...] and #![...] attributes
* tt: a single token tree
* vis: 
* lifetime
* literal:

2. repetitions:
`$($Variant:ident),* // one or more identity joined by comma => a, b, c`
`$($Variant:ident,)* //one or more identity end with comma => a, b, c,`
```rust
{ $($Variant:ident),* $(,)? }
//one or more identity joined by comma, each identity may end with comma 
//match  a, b, c, or a, b, c
```

3. example
```rust
#[macro_export]
macro_rules! vec {
  ($($x:expr),*) => { 
     {
        let mut v = Vec::new();
        $(v.push($x);)*
        v
    }
  };
}
```

## Procedural macro 
1. Function like 
2. Attribute like
3. Drive
```
debug:
#![feature(trace_macros)]
macro_rules! rpn { /* ... */ }
fn main() {
  trace_macros!(true);
  let e = rpn!(2 3 7 + 4 *);
  trace_macros!(false);
  println!("{}", e);
}
```
