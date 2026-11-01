# `std::maps`

Import `std::maps`. Source: [`maps.t`](../../base/std/maps/maps.t).
`Map<K, V>` is a heap- or arena-backed generic owner. Keys use compiler-defined
[semantic hashing and equality](../reference/expressions.md#equality-and-key-eligibility),
without callbacks or annotations. Invalid keys are diagnosed when the map type
is instantiated, including nested unsupported fields or enum payloads.

## API

Type arguments can be inferred from arguments and expected result types. For
example, `scores maps::Map<string, i32> := maps::create()` and
`maps::insert(&scores, "alice", 10)` work; explicit `<string, i32>` is also valid.
Representation fields, iterator fields, and `std::internal::swiss` are internal.

| Operation | Result / behavior |
| --- | --- |
| `create<K, V>()` | Empty `Map<K, V>`; no allocation or randomness |
| `create_in<K, V>(arena)` | Empty selected-allocator map; null arena selects heap |
| `insert(&map, key, value)` | `(old V, replaced bool, error)`; insert or replace only the value |
| `insert_if_absent(&map, key, value)` | `(inserted bool, error)`; existing entry stays unchanged |
| `get(&map, key)` | `(V, bool)`; shallow value copy and found flag |
| `get_ptr(&map, key)` | Borrowed `*V`, or null when missing |
| `contains(&map, key)` | `bool` membership |
| `remove(&map, key)` | `(stored_key K, stored_value V, found bool)` |
| `count(&map)` | `u64` live entry count |
| `capacity(&map)` | `u64` total entries supportable without increasing bucket count |
| `reserve(&map, requested u64)` | `error`; ensure room for at least this many **total** entries |
| `clear(&map)` | Remove entries, retain storage; no element cleanup |
| `free(&map)` | Release heap backing or forget arena backing; reset count/capacity, retain allocator; no element cleanup |
| `clone(&map)` | `(Map<K, V>, error)`; independent heap backing, shallow entries |
| `clone_in(&map, arena)` | `(Map<K, V>, error)`; independent selected-allocator backing, shallow entries |
| `iter(&map)` | Allocation-free `Iterator<K, V>` borrowing the live map handle |
| `next(&iterator)` | `(K, V, bool)`; shallow copies, false at exhaustion |

Missing results are zero values with false flags. Exhausted iterators keep
returning zero/zero/false; a zero iterator is exhausted. Failed insertions return
zero/false/error; failed insert-if-absent returns false/error. A duplicate
insert-if-absent is an ordinary false/success result.

Replacement preserves the **original stored key** and returns the old value.
The incoming equal key is unused and remains the caller's responsibility.
Rejected insert-if-absent leaves both inputs with the caller. Removal returns
the actual stored key, not the lookup argument. No operation implicitly destroys
keys, values, displaced values, or unused inputs.

`zero<Map<K, V>>()` is immediately usable in heap mode.
`memory::take(&map)` transfers the handle and leaves a usable zero source.
Repeated `clear`/`free` of the same handle is harmless. See [sets](sets.md) for
element-only containers. Prefer `for key, value in map` for borrowed iteration;
explicit iterators remain available.

Clones preserve every live entry, including stored-key identity and NaNs, but
do not promise the source's capacity, bucket layout, seed, or iteration order.
An empty source clones without allocation or randomness, even when it has
reserved storage. `clone` always selects heap, regardless of the source allocator;
`clone_in` selects its argument, not the source allocator. Failed clones return
an empty, zero-capacity, cleanup-safe handle retaining that selected allocator.

## Keys and ownership

Keys/values are shallow copies. Assignment copies the map handle, not its
storage. Choose exactly one lifecycle owner: stale copies are not updated by
mutation/free and must not be freed. This is a contract, not borrow checking.

String keys borrow bytes which must remain alive and unchanged while stored;
embedded zeros and malformed UTF-8 are ordinary bytes. Pointer keys compare
identity, including null, without reading pointees. Do not mutate a stored key
in a way affecting hash/equality. There is no mutable-key API.

Signed floating zero is one key. NaN does not equal itself: repeated NaN-key
insertion creates separate entries, and lookup/removal by NaN misses. This also
applies to NaN-containing aggregates. Such entries survive growth/cloning, appear in
iteration, and are discarded by clear/free. Unknown enum tags use only their
tag, never inactive payload bytes.

Owned values (for example nested lists or `strings::Owned`) require explicit
cleanup. Retrieve a displaced value or remove an entry and clean up the returned
owners. Before clear/free, iterate to release any element resources you own.
Iteration copies owners shallowly: cleaning up each copy once is appropriate
only when immediately discarding the corresponding stored entries, with no
subsequent use of those stale handles. See
[ownership conventions](strings-and-memory.md#ownership-conventions-and-transfer).

Cloning duplicates only table storage, **not** resources owned by keys/values.
Cloned strings still borrow the same bytes; cloned nested owning handles share
their underlying resource. Keep one resource owner alive until every copy is
discarded. Do not independently destroy shared elements from both maps.

## Arenas

The selected arena must outlive the map's storage and remain at the same address;
do not transfer/move it while maps refer to it. Growth allocates naturally aligned
fresh entry storage and control bytes, leaving old allocations in the arena until
reset/destruction. They are never individually freed. This does not extend the
public lifetime of entry pointers invalidated by growth.

Clear retains capacity, allocator, and seed. Free forgets all storage and seed but
retains the arena pointer, so subsequent insert/reserve still uses that arena.
Before resetting/destroying an arena, free all its map handles (and release any
separately owned element resources as needed). Reset/destruction invalidates all
arena-backed storage and borrows, not independent heap clones; shallow elements
in those clones may still borrow arena resources. Do not read or mutate a map
with backing invalidated by reset. An emptied map can be reused after arena
reset/destruction if the arena object itself is still live at the same address.

## Borrowing and invalidation

Reads, iteration, failed operations, and no-op reservation preserve entry
pointers. Construction creates no borrows. Cloning, successful or failed,
preserves all source entry pointers and iterators; mutating/freeing a clone does
not invalidate source table storage (shared shallow element resources are a
separate lifetime obligation). Replacement updates the existing value in place,
preserving its address and stored key. A successful new insertion may grow or rebuild to
remove tombstones, invalidating **all** entry pointers; otherwise existing
entries stay in place. A successful reserve that increases capacity invalidates
all entry pointers. Remove invalidates only the removed entry; it does not
shift other entries. Clear/free invalidate all entries. Taking a map preserves
entry storage, but an iterator still borrows its original map handle.

Direct loops and explicit iterators visit every live entry once, with unspecified order, and
allocate nothing. The borrowed map handle must remain at the same live address.
Do not insert (including replacement), remove, reserve, clear, free, transfer,
or otherwise structurally mutate it during iteration. Updating existing values
through `get_ptr` is permitted; subsequent `next` reads the current value.
Assigning a returned key/value copy does not modify the stored entry.
Clear/free invalidate iterators as well as entries. Reads, no-op reserve, failed
operations, and cloning preserve iterators; the prohibition on structural
operations during iteration still applies. Transferring the handle leaves entry
pointers valid but invalidates iterators borrowing the old handle.

Direct loops require exactly two binding names and addressable, mutable storage;
store returned temporary maps in a variable before iterating. The source is
evaluated once without creating another owner. See
[collection iteration](../reference/collections.md#maps-and-sets) for lifetime
requirements and normal loop cleanup/control flow.

## Allocation and checked failure

Empty construction/cloning, reads, removal, iteration, clear/free, and reserve(0) need
neither allocation nor randomness. Storage-demanding operations acquire a
secret per-table OS seed lazily. It is retained across clear/growth, forgotten
by free, and never exposed as an application hashing contract. Hashes/order are
not stable serialization identifiers. OS support follows the
[internal hashing boundary](../contributing/hashing.md#entropy-and-verification).

Seed failure forwards `errors::Error::RandomnessFailed`, with no weak fallback.
Allocation failure is `AllocationFailed`; checked capacity/byte arithmetic
overflow is `SizeOverflow`. Compare classified identities rather than messages.
Growth builds both new buffers before publishing changes; failed reserve/insert
preserves entries, count, capacity, seed, and every previously valid borrow, and
releases partial heap buffers. Existing-key replacement and duplicate
insert-if-absent cannot allocate or acquire randomness. Tombstone cleanup can
allocate even when the bucket count need not grow.

Arena errors are forwarded from `memory::arena_alloc`, including alignment and
size errors. Failure atomicity is **logical**, not an arena rollback: a partial
allocation can consume arena space or append a block even though insert/reserve
leaves the map unchanged or clone returns an empty handle. This space remains
arena-owned until reset/destruction. Heap partial allocations are released.

## Runnable example

```telia
package main

use "std::maps"
use "std::io"
use "std::libc"

fn run() error {
    scores := maps::create<string, i32>()
    defer maps::free(&scores)
    old, replaced := maps::insert(&scores, "alice", 10) @try
    previous, updated := maps::insert(&scores, "alice", 20) @try
    value := maps::get_ptr(&scores, "alice")
    *value = *value + 5
    for key, score in scores {
        io::println(key)
        libc::printf("%d\n", score)
    }
    key, removed_value, removed := maps::remove(&scores, "alice")
    libc::printf("%llu\n", maps::count(&scores))
    return nil
}

fn main() i32 {
    run() @fail
    return 0
}
```

Expected output:

```text
alice
25
0
```

See also [`examples/maps.t`](../../examples/maps.t) and the
[engine contract](../contributing/maps-engine.md).

For manual stepping, retain the explicit API:

```rust
iterator := maps::iter(&scores)
while true {
    key, score, found := maps::next(&iterator)
    if not found {
        break
    }
    // Use key and score; the same borrowing restrictions apply.
}
```

Existing bindings can also be updated with `key, score, found = maps::next(&iterator)`;
the call executes once before its results are assigned. To reassign a fallible
multi-result call, bind its full result including `error` and check that error.
Error-handled (`@try`/`@fail`/`@catch`) multi-result reassignments are currently
diagnosed; use fresh `:=` bindings for those forms.

## Shallow ownership and arena clone example

One separately owned string is shared by three independently managed tables.
Discard all tables before releasing that string; do not free each copied value.

```telia
package main

use "std::maps"
use "std::memory"
use "std::strings"
use "std::io"

fn run() error {
    owner := strings::clone("shared") @try
    defer strings::free(&owner)
    arena := memory::arena_create() @try
    defer memory::arena_destroy(&arena)
    original := maps::create_in<u64, strings::Owned>(&arena)
    defer maps::free(&original)
    maps::insert(&original, 1, owner) @try
    heap := maps::clone(&original) @try
    defer maps::free(&heap)
    selected := maps::clone_in(&original, &arena) @try
    defer maps::free(&selected)
    maps::clear(&original)
    maps::free(&selected)
    value, found := maps::get(&heap, 1)
    io::println(strings::view(&value))
    maps::free(&original)
    memory::arena_reset(&arena)
    // The heap table and separately heap-owned string survive arena reset.
    value, found = maps::get(&heap, 1)
    io::println(strings::view(&value))
    return nil
}

fn main() i32 {
    run() @fail
    return 0
}
```

Expected output:

```text
shared
shared
```
