# Values and types

## Variables and scope

These are statement fragments:

```rust
count := 1
wide i64 := 1
count = count + 1
left i32, right i64 := 10, 20
```

`:=` declares names; `=` assigns existing writable destinations. An annotation
is `name Type`, not `name: Type`. The right-hand side must provide the same
number of values as the destinations. Inference fixes a variable's type; it
does not make it dynamically typed. Use `zero<T>()` for an explicit zero value.

Blocks introduce scope. A local name cannot be declared twice in the same
scope; nested scopes can shadow enclosing locals. Loop-body declarations and
loop initializer variables do not escape the loop. Constants have no writable
or addressable storage; see [constant expressions](constants-and-bitwise.md).

## Primitive types

| Type | Meaning |
| --- | --- |
| `bool` | Boolean, with literals `true` and `false` |
| `i8`, `i16`, `i32`, `i64`, `i128` | Signed integer of the specified bit width |
| `u8`, `u16`, `u32`, `u64`, `u128` | Unsigned integer of the specified bit width |
| `int`, `uint` | 32-bit signed and unsigned default integer types |
| `f32`, `f64`, `float` | 32-bit, 64-bit, and default 32-bit floating point |
| `string` | Read-only borrowed byte pointer and length |
| `error` | Presence, optional classified identity/payload, borrowed message |
| `rawptr` | Untyped object-storage pointer; no lifetime or bounds information |

Do not confuse Telia primitives with [C-native scalar aliases](../stdlib/ffi-c.md).
Compatibility of literal values does not imply arbitrary implicit conversions
between typed variables. Prefer explicit [widening and checked narrowing](conversions.md)
at integer boundaries. Float, boolean, enum, and pointer conversions are not
provided by the integer conversion operations.

`nil` supplies a null pointer or successful error in an appropriate context.
It alone cannot infer a generic pointee type. Functions with no result omit
their result annotation; do not use bare `void` as a generic value argument.

## Structs and aliases

```rust
type Count = u64

struct Point {
    x i32
    y i32
}

point := Point{x: 10}
point.y = 20
empty := Point{}
```

Aliases are transparent, not new nominal types. Struct construction uses named
fields; omitted fields are zero-initialized. Duplicate/unknown fields are
invalid. Field access composes with indexing and calls. Fields of temporary
records may be read, but not assigned or addressed.

Native Telia records use `struct`; `union` is a foreign-record facility.
Pointer-recursive records are possible, but infinitely sized by-value layouts
are rejected. A record owning another resource does not acquire a destructor.

## Aggregate and pointer types

- `*T` is a typed pointer; `const *T` prevents writes through that pointer.
- `T[N]` is an inline fixed-length array; `T[]` is a borrowed writable slice.
- `(A, B)` is a tuple, including for multiple function results.
- `fn(value T) R` is a concrete function signature, not a capturing closure.
- Structs, enums, and functions may have [type parameters](generics.md).

See [collections](collections.md) for suffix order, array lengths, and bounds;
[enums](enums.md) for nominal identity and payloads; and
[C interoperability](foreign-interop.md) for qualifier and callback restrictions.

## Copying and zero values

Assignments, arguments, and results copy values. Structs and arrays copy their
inline fields/elements. Pointers, strings, slices, and owning handles remain
shallow: their copied addresses still refer to the same backing storage.
Nothing is implicitly deep-cloned, reference-counted, or freed.

`zero<T>()` returns a zero-initialized concrete value without acquiring storage
for referenced resources. A zero value is not necessarily a valid state for an
arbitrary foreign resource. Never assume zeroing a live owner releases it.
See [memory and lifetimes](memory.md).
