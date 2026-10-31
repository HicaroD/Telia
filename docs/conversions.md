# Explicit conversions

These compiler-supported operations work in ordinary Telia code without native
bindings, including in generic functions. Each takes one destination type and
one value; ordinary declarations may shadow their names. Conversions do not
change the existing implicit numeric compatibility rules.

## Integer widening

```rust
small u32 := 4294967295
large := widen<i64>(small)
```

`widen<T>(value)` returns `T`. It accepts integer types only and must preserve
**every** value of the source type, not just the current value:

| Source → destination | Allowed |
| --- | --- |
| Same signedness | Destination width is at least the source width |
| Unsigned → signed | Destination is strictly wider |
| Signed → unsigned | Never; use checked conversion |

Aliases and C scalar aliases use their resolved widths and signedness. Identity
conversions are allowed. Integer literals without an explicit source type use
their normal inferred type; type a variable first when its range matters.

## Checked integer conversion

```rust
index u64 := 255
byte := narrow<u8>(index) @try
```

`narrow<T>(value)` returns `(T, error)`. All integer-to-integer combinations are
accepted, including signedness changes, widening, identity, and 128-bit types.
The source expression is evaluated once. If the value fits the destination,
the result is exact and the error is nil. Otherwise the result is zero and the
error has identity `errors::Error::ConversionOutOfRange`. No truncation,
wrapping, or saturating conversion occurs. Handle it with `@try`, `@catch`, or
explicit tuple assignment; `@fail` remains an explicitly fatal option.

Neither operation converts floating-point, boolean, enum, or pointer values.

## Low-level object-pointer conversion

```rust
value i32 := 42
storage := ptr_cast<rawptr>(&value)
pointer := ptr_cast<*i32>(storage)
readonly := ptr_cast<const *i32>(pointer)
```

`ptr_cast<T>(value)` returns `T`. Both source and destination must be object
pointers or `rawptr`. It preserves the address, including null, and performs no
allocation, dereference, validation, ownership transfer, or lifetime extension.
It cannot remove immediate pointee constness; const pointers cannot be erased
to mutable `rawptr`. Nested pointer sources cannot be reinterpreted: their
pointee type/qualifiers must match, apart from adding immediate constness.
Integer addresses, callbacks, and C-string values are not accepted.

The caller must ensure alignment, live storage, a compatible stored value and
aliasing discipline before dereferencing. A cast does not make an invalid
pointer valid, authorize mutation, or make interior/arena pointers heap-freeable.
