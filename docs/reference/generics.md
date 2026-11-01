# Compile-time generics

Named literals use `Type{field: value}`; see [literal syntax](literals.md) for
control-flow header rules and compatibility with the old dotted spelling.

Structs, enums, and functions can declare type parameters with angle brackets.
Type arguments can be explicit, inferred, or partially explicit using `_`.

```rust
struct Pair<A, B> {
    first A
    second B
}

fn identity<T>(value T) T {
    return value
}

fn make_pair<A, B>(first A, second B) Pair<A, B> {
    return Pair<A, B>{ first: first, second: second }
}

entry := make_pair<string, i32>("count", identity<i32>(42))
```

## Inference and partial arguments

Existing explicit applications remain valid and retain their meaning. When
arguments are omitted, the compiler combines information from value arguments
and expected types (annotated initializers, assignments, returns, and surrounding
function parameters):

```rust
number i32 := 42
copy := identity(number)                      // T = i32
entry := make_pair("count", number)           // A = string, B = i32
partial := make_pair<string, _>("count", number)
wide i64 := identity(42)                      // expected result supplies i64
```

Partial lists must include every position: use `<string, _>`, not `<string>`,
for a declaration with two parameters. Explicit positions are fixed; `_` asks
the compiler to infer that position. A completely omitted list infers all
positions. Placeholders are application syntax, not standalone types or missing
type arguments in variable annotations or declaration signatures.

Inference matches pointers, arrays, slices, tuples, callback signatures, and
nested generic types structurally. Concrete evidence is collected before
existing numeric literal defaults are applied. Conflicting concrete types are
errors; inference does not silently widen numbers. `nil` alone cannot determine
a pointee type.

Constructors use fields, payloads, and expected destination types:

```rust
entry := Pair{first: "count", second: number}
empty Pair<string, i32> := Pair{}
item := Option<_>::Some{value: number}
none Option<i32> := Option::None
```

Here `Option<T>` is the enum from [enums](enums.md). Enum patterns can omit
arguments already determined by the matched value: `Option::Some{value}` and
`Option::None`. The enum's declaration must still match; a general `error`
value does not reveal an error enum's generic arguments.

A generic function can become a concrete callback through an expected signature
or explicit specialization:

```rust
callback fn(value i32) i32 := identity
explicit := identity<i32>
partial_callback fn(value i32) i32 := identity<_>
```

Callback inference preserves calling-convention requirements. It introduces no
runtime-polymorphic function values and does not relax foreign ABI restrictions.

Without enough evidence, supply explicit arguments: an unannotated `Pair{}`, a
payload-free `Option::None`, a result-only generic call, or a bare generic
function value cannot invent unknown types. Inference uses local type information,
not operations inside a function body.

### Refactoring hints

The compiler emits an informational hint when explicit arguments can be removed
without changing the canonical concrete specialization. If complete omission is
unsafe, it may suggest replacing some positions with `_` instead. For example,
`identity<i64>(42)` must retain its type when omission would default to `int`.

Hints do not affect build success and never edit source. Explicit spelling
remains supported. Hint analysis is read-only and conservative: uncertain rewrites
and specialization-dependent suggestions inside generic templates stay silent.
Overlapping suggestions are avoided, as is circular evidence from an enclosing
type that was inferred from the very arguments being considered for removal.

## Qualified and nested applications

Use the same notation after a package qualifier:

```rust
value := helpers::identity<i32>(42)
entry helpers::Pair<string, i32> := helpers::Pair<string, i32>{
    first: "count",
    second: value
}
```

Applications nest: `Box<Pair<string, i32>>`. Parameters can occur in fields,
function parameters, return types, pointers, tuples, local declarations, and
arguments to other generics. A parameter bound to a struct can also construct a
struct literal, as in `return T{}`. That operation is rejected if `T` is
instantiated with a non-struct type.

## Type arguments and scope

Arguments must be concrete value types: numeric types, `bool`, `string`,
imported C-string types, `rawptr`, typed pointers, structs, tuples, `error`, and fully applied
generic types, including arrays and slices such as `i32[3]` and `i32[]`.
Bare `void`, literals, `nil`, and generic declarations missing
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

Structural `==` and `!=` use the resolved concrete fields of generic records,
enums, arrays, and tuples. A generic comparison such as `return left == right`
is checked when specialized: `T = Pair<string, i32>` can be eligible, while a
slice-containing specialization is rejected with its unsupported component path.
Every enum payload alternative participates in eligibility. Automatic key
eligibility uses the same recursive classification with separate exclusions for
function values and standard-library owning handles. No traits, annotations,
custom equality, hashing operation, or map/set API is introduced here. See
[equality](expressions.md#equality-and-key-eligibility).

In particular, copying `Pair<string, i32>` copies the borrowed string view, not
its bytes. Storing an owning handle does not clone its storage or arrange its
cleanup. Lifetimes, resource ownership, and explicit cleanup remain the
programmer's responsibility; `defer` works inside generic functions as usual.

See [strings and explicit memory management](../stdlib/strings-and-memory.md) for those
contracts, [collections](collections.md) for generic element storage, and
[the runnable generic example](../../examples/generics.t).

```sh
make build
./telia build examples/generics.t -debug -o bin/generics-debug
./bin/generics-debug
./telia build examples/generics.t -release -o bin/generics-release
./bin/generics-release
```
