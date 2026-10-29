# Named literals

Construct structs, foreign records/unions, and enum payloads with braces directly
after the type or variant name:

```rust
point := Point{x: 10, y: 20}
box := Box<i32>{value: 42}
record := ffi::Record{value: 42}
item := Option<i32>::Some{value: 42}
```

The same notation applies to enum payload patterns:

```rust
value := match item {
    Option<i32>::Some{value: result} => result
    Option<i32>::None => 0
}
```

Patterns support field-name bindings (`{value}`), renaming (`{value: result}`),
and `..` to omit remaining fields. Construction and initialization semantics
are unchanged. The older dotted spelling is still accepted for compatibility,
but examples and libraries use the undotted form.

## Control-flow headers

In an `if`, `elif`, `while`, `for ... in`, or `match` header, an unparenthesized
brace after a name starts the body. Parenthesize a named literal used directly
in that header:

```rust
value := match (Option<i32>::Some{value: 42}) {
    Option<i32>::Some{value} => value
    Option<i32>::None => 0
}
```

Call arguments, parentheses, and index/slice expressions already delimit an
expression, so literals inside them need no additional parentheses:

```rust
if contains(Point{x: 10, y: 20}) {
    // ...
}
```

The same body-brace rule applies to the update clause of a three-clause `for`.
Its initializer and condition are delimited by semicolons instead.

## Arrays

Contextual array literals retain their dot and require an expected array type:

```rust
values i32[3] := .{10, 20, 30}
empty i32[3] := .{}
```
