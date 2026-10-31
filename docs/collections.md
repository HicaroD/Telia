# Arrays, borrowed slices, and lists

Arrays and slices are language types. Lists are ordinary generic Telia code in
[`std::lists`](../base/std/lists/lists.t), not a compiler-recognized collection.
The same language and memory operations are available to application libraries.
Maps and sets are not part of this milestone.

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

Import `std::lists`. A `List<T>` stores a backing slice, a live count, and an
optional arena pointer. Its zero value is an empty heap-backed list.

```rust
items := lists::create<i32>()
defer lists::free<i32>(&items)
lists::append<i32>(&items, 42) @fail
value := lists::get<i32>(&items, 0) @fail
```

| API | Result and behavior |
| --- | --- |
| `create<T>()` | Empty heap-backed `List<T>`; no allocation |
| `create_in<T>(arena *memory::Arena)` | Empty arena-backed `List<T>`; no allocation |
| `count<T>(list *List<T>)` | `u64` live element count |
| `capacity<T>(list *List<T>)` | `u64` element capacity |
| `reserve<T>(list *List<T>, requested u64)` | `error`; ensure capacity without changing count |
| `append<T>(list *List<T>, value T)` | `error`; append a shallow value copy |
| `append_slice<T>(list *List<T>, values T[])` | `error`; append shallow copies of a slice |
| `view<T>(list *List<T>)` | Writable `T[]` covering live elements only |
| `get<T>(list *List<T>, index u64)` | `(T, error)`; zero value on invalid index |
| `set<T>(list *List<T>, index u64, value T)` | `error`; overwrite an existing element |
| `pop<T>(list *List<T>)` | `(T, bool)`; remove last, or zero/false when empty |
| `remove<T>(list *List<T>, index u64)` | `(T, error)`; ordered removal, returning the displaced element |
| `remove_unordered<T>(list *List<T>, index u64)` | `(T, error)`; fill the gap with the last element |
| `clear<T>(list *List<T>)` | Set count to zero, retaining capacity and allocator |
| `clone<T>(list *List<T>)` | `(List<T>, error)`; independent heap backing, shallow element copies |
| `clone_in<T>(list *List<T>, arena *memory::Arena)` | Same, with independent backing in the selected arena |
| `free<T>(list *List<T>)` | Release heap backing, or forget arena backing; reset count/capacity, retain allocator |

The arena must outlive its lists. A null arena pointer selects the heap.
Arena-backed growth allocates a new contiguous buffer and copies live elements;
old buffers remain until arena reset/destruction. Do not individually free arena
storage. Before arena reset, discard/free list handles; reset invalidates all
their storage. An emptied list retaining that arena can then be reused.

Capacity and byte-size arithmetic is checked. Failed reserve/append leaves the
original list, its capacity, and its views unchanged. Failed clone returns an
empty handle. Appending a live view of the same list is supported, including
when reallocation is needed; input elements are read before old storage is freed.
Growth invalidates prior element pointers and views. Removal/clear invalidate
the affected logical elements even when the backing allocation remains.

`set` overwrites rather than returns the old value. Retrieve and clean up an
owned old value before replacing it. Removal/pop return shallow copies so
callers can manage their resources explicitly.

### Ownership and cleanup

Copying a list copies its storage handle, **not** its allocation. Arrange one
lifecycle owner; freeing one copy does not update any other copy. `clone` copies
backing storage, not resources owned by elements. Arrays likewise copy element
handles shallowly; a copied borrowed string still borrows the same bytes.

Neither `clear` nor `free` runs element cleanup. For nested owning elements,
release their resources explicitly before clearing/freeing the outer list.
Repeated cleanup of the same emptied list is harmless; cleanup of a stale
shallow copy can double-free. Lifetimes and ownership remain programmer-managed.
Do not modify the list representation fields directly.

Transfer a list with `memory::take<lists::List<T>>(&source)` instead of a shallow
copy followed by manual zeroing. The source becomes an empty heap-mode list,
safe to free or reuse; the destination retains the original arena association.
Moving an arena-backed list does not extend the arena's lifetime. See the
[ownership conventions](strings-and-memory.md#ownership-conventions-and-transfer).

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
