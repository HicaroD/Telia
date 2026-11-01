# `std::sets`

Import `std::sets`. Source: [`sets.t`](../../base/std/sets/sets.t).
`Set<T>` is a heap- or arena-backed generic owner using the same
[table engine](../contributing/maps-engine.md) as maps. Elements automatically
use compiler-defined [key eligibility, hashing, and equality](../reference/expressions.md#equality-and-key-eligibility).
The complete eligibility matrix there applies to set elements: no annotations,
traits, or callbacks are required. Unsupported types, including nested fields
and enum payloads, are diagnosed even for unused constructors or zero handles.
Known standard-library owning handles are not eligible elements; pointers to
them are eligible by pointer identity, not by their contents.

## API

Generic inference works from arguments and expected results:
`names sets::Set<string> := sets::create()` and
`sets::insert(&names, "alice")`. Explicit `<string>` remains supported.
Representation fields, `Unit`, iterator fields, and `std::internal::swiss` are
implementation details, not application APIs. No dummy values are exposed by
set operations.

| Operation | Result / behavior |
| --- | --- |
| `create<T>()` | Empty heap `Set<T>`; no allocation or randomness |
| `create_in<T>(arena *memory::Arena)` | Empty selected-allocator set; null selects heap |
| `insert(&set, element)` | `(inserted bool, error)`; preserve an equal stored element |
| `contains(&set, element)` | Membership `bool` |
| `remove(&set, element)` | `(stored_element T, found bool)` |
| `count(&set)` | `u64` live entry count |
| `capacity(&set)` | `u64` total entries supportable without increasing bucket count |
| `reserve(&set, requested u64)` | `error`; ensure capacity for at least this many **total** entries |
| `clear(&set)` | Remove elements, retain storage, seed, and allocator; no element cleanup |
| `free(&set)` | Release heap backing or forget arena backing; reset count/capacity/seed, retain allocator; no element cleanup |
| `clone(&set)` | `(Set<T>, error)`; independent heap backing, shallow elements |
| `clone_in(&set, arena)` | `(Set<T>, error)`; independent selected-allocator backing, shallow elements |
| `iter(&set)` | Allocation-free `Iterator<T>` borrowing the live set handle |
| `next(&iterator)` | `(element T, found bool)`; shallow copies, false at exhaustion |

Duplicate insertion is `false, nil`: it does not allocate, acquire randomness,
replace the actual stored element, or consume the incoming equal element.
Failed insertion always returns false with an error. Removal returns the
**actual stored element**, not the lookup argument, so its identity remains
available for explicit cleanup. Missing removal returns zero/false. Zero and
exhausted iterators return zero/false, including on repeated calls.

`zero<Set<T>>()` is immediately usable in heap mode. `memory::take(&set)`
transfers the handle and leaves such a usable zero source. Repeated clear/free
of the same handle is harmless. Direct `for ... in set` is not supported yet;
use explicit iteration.

Signed floating zero is one element; inserting an equal opposite-sign zero
preserves the sign originally stored. NaN does not equal itself, including
inside aggregates: repeated insertion stores separate entries, lookup/removal
by NaN misses, and all live NaNs survive growth/cloning and appear in iteration.
Clear/free discard them without element cleanup.

## Ownership and borrowing

Elements and handle assignments are shallow copies. Choose exactly one table
lifecycle owner; stale handle copies are not updated by mutation/free and must
not be independently freed. String elements borrow bytes that must stay alive
and unchanged while stored, including embedded zeros and malformed UTF-8.
Pointers compare identity without reading pointees. Do not mutate any stored
element's hash/equality-relevant contents; there is no mutable-element API.

An eligible user record can still represent an explicitly owned resource (for
example a raw pointer). The set does not infer ownership or destroy resources.
The caller cleans up unused duplicates, failed insertion inputs, and removed
elements. Before clear/free, explicitly release any resources owned by stored
elements. Iteration copies owners shallowly: destroy each resource only once,
immediately discard the stale stored entries, and do not use those copies again.
See [ownership conventions](strings-and-memory.md#ownership-conventions-and-transfer).

Clones duplicate only table backing, not resources referenced by elements.
Strings still borrow the same bytes; pointers still designate the same objects.
Keep one resource owner alive until all shallow copies are discarded. Never
destroy shared resources once per cloned set. Clones preserve actual stored
elements and every NaN, but not capacity, layout, seed, or iteration order.
`clone` always selects heap; `clone_in` selects its argument, not the source
allocator. Empty clones require no allocation/seed, even from reserved sources.

Explicit iterators borrow the set handle at its current live address and visit
every live element once in unspecified order. No insertion (even duplicates),
removal, reservation, clear, free, transfer, or structural mutation is permitted
during iteration. This is a caller contract, not borrow checking. Assigning to
a returned element copy does not modify the set.

Construction creates no borrows. Reads, failed operations, no-op reserve, and
cloning preserve source iterators and backing. Mutating/freeing an independent
clone does not invalidate source table backing (shared resources have separate
lifetimes). New insertion may rebuild/grow; a capacity-increasing reserve
rebuilds. Both invalidate backing borrows. Remove invalidates the removed entry;
clear/free invalidate all entries and iterators. Taking the handle preserves
storage but invalidates iterators borrowing the old handle. There is no public
entry-pointer operation for sets.

## Arenas and checked failure

The arena must outlive its set backing and remain at the same address. Growth
allocates aligned fresh buffers and leaves old buffers arena-owned until reset
or destruction. Free forgets backing but retains the arena association; later
insert/reserve uses that arena. Free all associated set handles before arena
reset/destruction; never use invalidated storage. Reuse of an emptied handle
requires the arena object itself to remain live at the same address. Independent
heap clones survive reset, but their shallow elements may still borrow arena
resources and thus become invalid.

Empty construction/cloning, reads, remove, iteration, clear/free, and reserve(0)
need neither allocation nor randomness. Storage-demanding operations lazily
acquire a secret OS-generated seed, retained through clear/growth and forgotten
by free. There is no weak fallback and no stable hash/iteration-order contract.
Target support follows the [hashing boundary](../contributing/hashing.md#entropy-and-verification).

Seed failure forwards `errors::Error::RandomnessFailed`; allocation failure is
`AllocationFailed`; checked capacity/byte-size overflow is `SizeOverflow`.
Arena errors are forwarded. Compare error identities, not messages. Failed
insert/reserve preserves contents, count, capacity, seed, and valid borrows;
partial heap buffers are released. Failed clone returns an empty zero-capacity,
cleanup-safe handle retaining the selected allocator and leaves the source
unchanged. Arena failure is logically atomic, not a rollback: partial allocation
can consume arena space/blocks until reset/destruction. Tombstone cleanup may
allocate even without increasing bucket count. Duplicates never allocate.

## Runnable example

```telia
package main

use "std::sets"
use "std::io"

struct Position {
    x i32
    y i32
}

fn run() error {
    visited sets::Set<Position> := sets::create()
    defer sets::free(&visited)
    sets::reserve(&visited, 20) @try
    added := sets::insert(&visited, Position{x: 3, y: 7}) @try
    duplicate := sets::insert(&visited, Position{x: 3, y: 7}) @try
    if added and not duplicate {
        io::println("one position")
    }
    copy := sets::clone(&visited) @try
    defer sets::free(&copy)
    sets::clear(&visited)
    cursor := sets::iter(&copy)
    while true {
        position, found := sets::next(&cursor)
        if not found {
            break
        }
        if position.x == 3 and position.y == 7 {
            io::println("3,7")
        }
    }
    stored, removed := sets::remove(&copy, Position{x: 3, y: 7})
    if removed and sets::count(&copy) == 0 {
        io::println("empty")
    }
    return nil
}

fn main() i32 {
    run() @fail
    return 0
}
```

Expected output:

```text
one position
3,7
empty
```

See [`examples/sets.t`](../../examples/sets.t) for separately owned string
storage shared by heap/arena shallow clones, with explicit cleanup order.
