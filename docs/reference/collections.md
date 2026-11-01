# Arrays, borrowed slices, and lists

Arrays and slices are language types. Lists are ordinary generic Telia code in
[`std::lists`](../../base/std/lists/lists.t), not a compiler-recognized collection.
The same language and memory operations are available to application libraries.
Heap-backed [maps](../stdlib/maps.md) provide automatic semantic keys and explicit
iterators. Sets and direct map `for ... in` syntax are not provided yet.

Generic functions can infer element types from array and slice arguments, and
from array-literal elements when the parameter supplies the array shape.
Explicit applications remain valid; partial applications use `_` in every
inferred position. See [generic inference](generics.md#inference-and-partial-arguments).

## Arrays

```rust
values i32[3] := .{10, 20, 30}
empty i32[3] := .{}
copy := values
values[1] = 42
```

`T[N]` contains inline storage for exactly `N` elements. `N` must currently be a
nonnegative decimal integer literal. A nonempty literal must supply exactly `N`
elements; `.{}` zero-initializes every element. Contextual literals need an
expected array type, including in fields, parameters, and returns.

Arrays can be fields, generic type arguments, and elements of other collections.
Suffixes build types left to right: `i32[2][3]` is three arrays of two integers.
Use aliases to make nested layouts clearer. `*T[N]` is a pointer to an array.
Arrays are emitted as concrete C wrapper structs, so assignment, parameter
passing, and returning copy all elements by value, not C array-to-pointer decay.
Zero-length arrays use a one-byte placeholder in generated C and have length
zero. Empty structs likewise have a one-byte placeholder. Inline storage does
not imply stack allocation.

`len(values)` returns `u64`. Indexing reads or writes an element. Indices must be
`u64` values or integer literals representable as `u64`; other integer types are
not implicitly converted by indexing. Negative indices are rejected. Ordinary
out-of-bounds indexing panics in **both debug and release** builds. Low-level
raw-storage views still require valid storage; bounds checks are not ownership
or pointer-validity checks.

Field access, indexing, and calls compose: `items[i].path`,
`container.items[i].path`, and `get_items()[i].path` need no intermediate local.
An indexed field can be read, assigned, or addressed (`&items[i].path`) when its
receiver denotes writable/addressable storage. Fields of temporary record
values can be read but not assigned or addressed. Indexing still checks bounds
and evaluates its receiver and index once per access. Use
[`widen<u64>` or checked `narrow<u64>`](conversions.md) for other index types.

## Borrowed slices

```rust
view := values[:]
part := view[1:3] @fail
part[0] = 99
end := view[3:3] @fail
empty_view := zero<i32[]>()
```

`T[]` is a writable element pointer and a `u64` length. Copying a slice copies
only this view. Slicing allocates nothing and never extends storage lifetime.
`string` remains a separate read-only byte view; it does not implicitly convert
to a writable slice.

`values[:]` returns a full view and cannot fail. Other ranges return `(T[],
error)`, using the usual `@fail`, `@catch`, or explicit tuple assignment. The
start is inclusive and the end exclusive; omitted endpoints default to zero and
the original length. `start > end` or `end > len(values)` returns an empty view
and an error. Empty ranges at any valid boundary are allowed. Empty views may
have null data pointers. Arrays must have addressable storage to borrow them;
store a returned array in a variable before slicing it.

Import `std::slices` for recoverable indexed access:

| API | Result |
| --- | --- |
| `get<T>(values T[], index u64)` | `(T, error)`; zero value on invalid index |
| `set<T>(values T[], index u64, value T)` | `error`; invalid index changes nothing |

The same checked APIs work on arrays through `values[:]`.

## Iteration

```rust
for value in values {
    // value is a shallow value copy of the current element.
}
```

Arrays and slices support `for value in collection`. The source is evaluated
once. Array iteration uses the original storage; slice iteration snapshots its
pointer and length. Changing the loop variable does not update the element;
write through indexing instead. Element pointers can be obtained explicitly
with `&view[index]`.

All loop forms (`while`, three-part `for`, and array/slice iteration) support
`break` and `continue`. `break` exits the nearest enclosing loop; `continue`
skips the rest of its body and starts the next iteration. Three-part `for`
loops run their update step before testing the condition again; range loops
advance to the next element. Neither statement takes a label or value.

Both statements run registered `defer` cleanup in each exited scope,
innermost first, without running cleanup in scopes outside the loop body.
They also work inside match-arm blocks nested in a loop. Using either outside
a loop is an error. Deferred cleanup cannot jump to a loop outside the cleanup
itself, but may use `break` and `continue` in its own nested loops.

List iteration uses its ordinary borrowed view, not a special compiler hook:

```rust
view := lists::view<i32>(&items)
for value in view {
    // Do not append, remove, clear, reserve, or free items here.
}
```

Structural list mutation during iteration is prohibited by the API contract,
not enforced by lifetime analysis. Nested scopes and deferred cleanup work as
usual. Generic user-defined iteration protocols are not introduced here.

## Dynamic lists

Lists are library values, not language types. See [`std::lists`](../stdlib/lists.md)
for the complete API, allocation/failure contracts, and ownership rules.
The compiler supplies only the ordinary generic, slice, pointer, and memory
operations used by that implementation.

## Building other data structures in Telia

Generic library code has these reusable compiler operations:

| Operation | Result |
| --- | --- |
| `sizeof<T>()` | `u64` concrete C layout size |
| `alignof<T>()` | `u64` concrete C alignment |
| `zero<T>()` | The zero-initialized value of any concrete value type |
| `slice_from_raw<T>(storage rawptr, count u64)` | Writable `T[]` over caller-provided storage; no allocation or initialization |
| `slice_data<T>(values T[])` | Borrowed `rawptr` to its backing storage |

Except the existing reserved `sizeof`, these names may be shadowed by ordinary
declarations. Size and alignment use the C backend's actual specialized types.
The raw-view operations are low-level trusted operations: the caller must
provide suitably aligned, live storage for `count * sizeof<T>()` bytes, ensure
arithmetic does not overflow, and initialize each element before reading it.
They perform no allocation, pointer validation, ownership transfer, or cleanup.
Only allocation-base pointers may be passed to heap `memory::free`; a pointer
from an array, subslice, or arena is not individually heap-owned.

Combine those operations with `memory::alloc`, `memory::arena_alloc`, indexed
writes, and ordinary generic functions. The list implementation demonstrates
growth, overflow checks, and failure-atomic publication entirely in Telia.
There are no list-specific compiler operations, automatic destructors, borrowing
checks, or new trait/interface systems.

## Verification

```sh
go test ./...
make build
./telia build examples/collections.t -debug -o bin/collections-debug
./bin/collections-debug
./telia build examples/collections.t -release -o bin/collections-release
./bin/collections-release
```

Collection integration tests exercise generated C with UndefinedBehaviorSanitizer
and deterministic allocation failure/leak counters. Bounds-panics are checked
at both `-O0` and `-O3`. On a host with a working AddressSanitizer runtime:

```sh
TELIA_TEST_SANITIZERS=address,undefined go test ./tests/integration -run TestCollectionsWithNativeSanitizers -count=1
```
