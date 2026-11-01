# Memory and lifetimes

Telia has explicit resource management. It does not enforce ownership, moves,
borrowing, reference counting, or automatic destruction.

## Pointers and views

`&value` takes an address and `*pointer` accesses the pointee. The caller must
ensure a pointer is valid, aligned, and alive before dereferencing. `rawptr`
does not carry an element type or size. A checked index cannot make invalid
backing storage safe.

`const *T` prohibits mutation through that pointer; it does not freeze storage
against other aliases or native code. Pointer qualification and conversion
restrictions are described in [conversions](conversions.md) and
[C interoperability](foreign-interop.md).

Strings and slices are borrowed views. Copying a view does not copy its storage
or extend its lifetime. Literals have static storage. Local inline storage,
heap owners, arenas, and native allocations have different lifetimes; keep each
view within the lifetime of its actual backing storage.

## Owners and transfers

An owning record is still an ordinary copyable value. Keep exactly one cleanup
owner for a live resource. A shallow owner copy is not a second allocation.
`memory::take<T>(&source)` returns its value and zeroes the source, without
allocation or cleanup. It is appropriate only when zero is an inert valid state
for that type. It cannot repair other shallow copies or extend an arena's life.

See [ownership conventions](../stdlib/strings-and-memory.md#ownership-conventions-and-transfer)
for strings, lists, arenas, transfer, and invalidation. Native file, child, and
terminal handles have their own close/restore contracts; do not generalize a
transfer pattern to an arbitrary foreign handle.

## Deferred cleanup

`defer` registers a statement in the current block when execution reaches it.
Normal block exit executes registered statements in reverse registration order.
`return`, `break`, `continue`, and error propagation run cleanup for scopes they
exit, innermost first. A return or match `yield` copies its result before cleanup.

Deferred expressions read current variables at cleanup time, not frozen copies
of the original call arguments. A statement after a taken return is never
registered. Do not return a borrowed view or error message whose owner the
cleanup frees: the copied pointer would dangle.

`@fail`, process exit, and native crashes do not provide scoped cleanup.
Use `@try` or `@catch` with ordinary returns for cleanup-safe failure paths.
Deferred cleanup must handle its own failures; `@try` is not allowed there.
See [error handling](errors.md#handling-and-discarding).

## Low-level storage operations

| Operation | Result and obligation |
| --- | --- |
| `sizeof<T>()` | `u64` concrete backend layout size |
| `alignof<T>()` | `u64` concrete backend alignment |
| `zero<T>()` | Zero-initialized concrete value; no resource cleanup |
| `slice_from_raw<T>(storage rawptr, count u64)` | Borrowed writable slice; caller provides valid aligned initialized storage and checks size arithmetic |
| `slice_data<T>(values T[])` | Borrowed raw backing pointer; does not confer heap ownership |

These operations allocate nothing. Except reserved `sizeof`, their names can be
shadowed by ordinary declarations. Only a heap allocation's base pointer can
be freed with the matching heap allocator; array, interior, subslice, and arena
pointers are not individually heap-freeable. See [memory APIs](../stdlib/memory.md).
