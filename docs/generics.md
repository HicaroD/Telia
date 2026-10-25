# Compile-time generics

Structs and functions can declare type parameters with angle brackets. Function
calls require explicit type arguments; there is no type-argument inference.

```rust
struct Pair<A, B> {
    first A
    second B
}

fn identity<T>(value T) T {
    return value
}

fn make_pair<A, B>(first A, second B) Pair<A, B> {
    return Pair<A, B>.{ first: first, second: second }
}

entry := make_pair<string, i32>("count", identity<i32>(42))
```

Use the same notation after a package qualifier:

```rust
value := helpers::identity<i32>(42)
entry helpers::Pair<string, i32> := helpers::Pair<string, i32>.{
    first: "count",
    second: value
}
```

Applications nest: `Box<Pair<string, i32>>`. Parameters can occur in fields,
function parameters, return types, pointers, tuples, local declarations, and
arguments to other generics. A parameter bound to a struct can also construct a
struct literal, as in `return T.{}`. That operation is rejected if `T` is
instantiated with a non-struct type.

## Type arguments and scope

Arguments must be concrete value types: numeric types, `bool`, `string`,
`cstring`, `rawptr`, typed pointers, structs, tuples, `error`, and fully applied
generic types. Bare `void`, literals, `nil`, and generic declarations missing
their arguments are invalid arguments. Aliases are transparent: an alias and
its underlying type select the same specialization.

Type parameters are local to their declaration. They may shadow an outer type
name. Duplicate parameters and reserved built-in type names are rejected.
Other names in a template resolve using the declaration's package and imports;
caller-local names do not capture them. Type arguments at a call or literal
resolve in the caller's context. Aliases retain their declaration context.

## Specialization and checking

The compiler substitutes types structurally and emits ordinary concrete C
structs and functions for the combinations used by the program. Repeated uses
of a declaration with the same canonical arguments reuse one specialization.
Different declarations, including same-named declarations in different
packages, have distinct identities. Templates retain their original syntax;
instances have independent scopes and semantic/backend state.

Template syntax, names, argument counts, and independent primitive operations
are checked without requiring an instance. Operations involving substituted
types are checked when specialized. For example:

```rust
fn double<T>(value T) T {
    return value + value
}

number := double<i32>(21)       // valid
text := double<string>("bad")   // Telia diagnostic: string addition is invalid
```

Recursive calls with the same arguments reuse the in-progress instance.
Pointer-recursive struct layouts are valid; infinitely sized by-value layouts
are rejected. Struct and tuple definitions are emitted in dependency order.

A compilation permits at most **128 active instantiation levels** and
**4,096 total specializations**, counting generic struct and function instances
together. Reusing an instance does not consume another slot. Exceeding a limit
reports a diagnostic instead of allowing unbounded expansion. Specialization
errors include template context and an instantiation chain.

Generic templates cannot carry C linkage attributes or serve as C prototypes.
Use concrete functions as interoperability wrappers.

## Memory and copying

Generics introduce no runtime type lookup, boxing, or hidden allocation.
Fields retain their concrete layout and alignment. Assignment and parameter
passing keep the ordinary copying behavior of the substituted type.

In particular, copying `Pair<string, i32>` copies the borrowed string view, not
its bytes. Storing an owning handle does not clone its storage or arrange its
cleanup. Lifetimes, resource ownership, and explicit cleanup remain the
programmer's responsibility; `defer` works inside generic functions as usual.

See [strings and explicit memory management](strings-and-memory.md) for those
contracts and [the runnable generic example](../examples/generics.t).

```sh
make build
./telia build examples/generics.t -debug -o bin/generics-debug
./bin/generics-debug
./telia build examples/generics.t -release -o bin/generics-release
./bin/generics-release
```
