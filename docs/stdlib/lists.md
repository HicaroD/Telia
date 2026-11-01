# `std::lists`

Source: [`lists.t`](../../base/std/lists/lists.t). Import `std::lists`.
`List<T>` is an ordinary generic Telia record, not a built-in container.

Fragment inside a function returning `error`:

```rust
items := lists::create<i32>()
defer lists::free<i32>(&items)
lists::append<i32>(&items, 42) @try
first := lists::get<i32>(&items, 0) @try
```

## API

A `List<T>` stores a backing slice, live count, and optional arena pointer.
Its zero value is an empty heap-backed list. Do not mutate representation fields
or call the internal `allocate_buffer` helper from applications.

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
| `clone_in<T>(list *List<T>, arena *memory::Arena)` | `(List<T>, error)`; independent backing in the selected arena |
| `free<T>(list *List<T>)` | Release heap backing, or forget arena backing; reset count/capacity, retain allocator |

## Allocation and failure

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
Structural mutation during borrowed-view iteration is prohibited by contract.

Invalid indexed access uses `OutOfBounds`; capacity arithmetic and allocation
failures use `SizeOverflow` and `AllocationFailed` (or a forwarded arena error).
`set` overwrites rather than returns the old value. Retrieve and clean up an
owned old value before replacing it. Removal/pop return shallow copies so
callers can manage their resources explicitly.

## Ownership and cleanup

Copying a list copies its storage handle, **not** its allocation. Arrange one
lifecycle owner; freeing one copy does not update any other copy. `clone` copies
backing storage, not resources owned by elements. Arrays likewise copy element
handles shallowly; a copied borrowed string still borrows the same bytes.

Neither `clear` nor `free` runs element cleanup. For nested owning elements,
release their resources explicitly before clearing/freeing the outer list.
Repeated cleanup of the same emptied list is harmless; cleanup of a stale
shallow copy can double-free. Lifetimes and ownership remain programmer-managed.

Transfer a list with `memory::take<lists::List<T>>(&source)` instead of a shallow
copy followed by manual zeroing. The source becomes an empty heap-mode list,
safe to free or reuse; the destination retains the original arena association.
Moving an arena-backed list does not extend the arena's lifetime. See the
[ownership conventions](strings-and-memory.md#ownership-conventions-and-transfer).

Runnable example: [`examples/collections.t`](../../examples/collections.t).
