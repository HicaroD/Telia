# `std::memory`

Source: [`memory.t`](../../base/std/memory/memory.t).
Import `std::memory`. This package provides uninitialized heap storage, arenas,
and explicit transfer; it does not manage arbitrary pointers automatically.

| API | Result | Contract |
| --- | --- | --- |
| `alloc(size u64)` | `(rawptr, error)` | Uninitialized heap allocation; zero size succeeds with null |
| `free(ptr rawptr)` | None | Release a matching heap allocation base; null is harmless; does not clear the caller's pointer |
| `arena_create()` | `(Arena, error)` | Empty lazy arena; currently no allocation or failure |
| `arena_alloc(arena *Arena, size u64, alignment u64)` | `(rawptr, error)` | Arena-owned storage; alignment must be a supported nonzero power of two |
| `arena_reset(arena *Arena)` | None | Invalidate all allocations/views, retaining backing blocks |
| `arena_destroy(arena *Arena)` | None | Invalidate allocations/views, release blocks, empty the arena |
| `take<T>(source *T)` | `T` | Copy out the value and replace it with `zero<T>()`; no allocation/cleanup |

Allocation failures return null and an error. Arena alignment, overflow, and
growth failures preserve earlier allocations. Heap allocations must be freed
once; arena allocations must **not** be individually freed. Arena growth does
not move earlier allocations. Reset/destruction invalidates every borrowed view.
An arena has one lifecycle owner despite being shallow-copyable.

This fragment belongs in a function returning `error`:

```rust
arena := memory::arena_create() @try
defer memory::arena_destroy(&arena)
storage := memory::arena_alloc(&arena, 64, 8) @try
```

The storage is not initialized. A typed/raw slice view also requires appropriate
alignment, checked size arithmetic, initialization, and a valid lifetime.
See [raw-storage operations](../reference/memory.md#low-level-storage-operations).

`take` requires a live mutable non-null pointer and a type whose zero state is
inert and cleanup-safe. It does not free an overwritten destination, repair
other copies, or move external references to the source's address. In particular,
do not move an arena while lists or maps still point to the old arena value.
See [transfer contracts and examples](strings-and-memory.md#ownership-conventions-and-transfer)
and [detailed arena behavior](strings-and-memory.md#arenas).
