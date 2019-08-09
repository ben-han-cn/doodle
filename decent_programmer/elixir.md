# function
- named funciton is identified by tuple of three module/func/arity
- named function is different from anonymous funciton
- closure can only be implement by anonymous function 
- invoke syntax is different `name_func(10)` vs `anonymous_func.(10)`

# atomic
- literal named constants
- :good vs :"good boy" 
- atom consists of text and value, `val = :someatom`, the text is saved
  in runtime atom table, the reference to it is saved in val
- module name is atom
- uppercase atom `AnAtom` is equal to `:"Elixir.AnAtom"` they are alias
- there is no dedicated boolean type, `true` just syntax sugar to `:true`

# binary & bitstring
- a chunk of bytes, `<<1, 2, 3>>`, 
- byte value bigger than 255 which be truncated `<<257>> == <<1>>`
- size of each value could be speicified `<<257::16>> == <<1, 1>>` 
- binary whose total size isn't a multiplier of 8 is called bitstring

# strings vs character list
- string are binary: `"good boy"`
- string has heredoc syntax with `"""`
- character list are list `'ABC' == [65, 66, 67]`

# runtime
- module name is coded to compiled file which including beam instruction
  `defmodule Geometry` => Elixir.Geometry.beam
- code path is used to search the module which isn't loaded in vm memory
```shell
idex -pa my/code/path -pa another/code/path
```
```elixir
:code.get_path
```
- apply is used to dynamic invoke function `apply(IO, :puts, ["hello world"])`
- every thing is excuted inside a process.

# abstract with module
```elixir
defmodule Fraction do
    defstruct a: nil, b: nil
    def new(a, b), do: %Fraction{a: a, b: b}

    def value(%Fraction{a: a, b: b}) do
        a/b
    end
end
one_half = %Fraction{a:1, b:2}
one_half = Fraction.new(1, 2)
```
- module is the combination of a structure and functions (use the struct as first parameter) 
  manipuate on it.
- struct is a map with a specical key `__struct__` whose value is the module name.
- module has attribute with syntax `@doc`, they are evaluated at compile time, normally be used
  as metadata or constant.
- `__info__/1` function is automatically injected into each Elixir module during compilation, it
  lists all exported functions of a module.

# protocol
- A protocol is a module in which you declare functions without implementing them
```elixir
defprotocol String.Chars do
    def to_string(thing)
end

defimpl String.Chars for: ToDoList do
    def to_string(thing) do
    end
end
```
- funciton should have at least one parameter which will be used to dispatch to right implementation
- Protocol implementation doesn't need to be part of any module.
- Inside the module, `defimpl` could omit the `for` part 
- Each implementing of the protocol will became the sub module of the protocol module
  
# distribute system

