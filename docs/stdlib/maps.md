# `std::maps`

Import `std::maps`. Source: [`maps.t`](../../base/std/maps/maps.t).
`Map<K, V>` is a heap-backed generic owner. Keys use compiler-defined
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
| `free(&map)` | Release backing storage, reset count/capacity; no element cleanup |
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
Repeated `clear`/`free` of the same handle is harmless. Arena constructors,
cloning, sets, and direct `for ... in map` are not provided yet; use explicit
iteration below.

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
applies to NaN-containing aggregates. Such entries survive growth, appear in
iteration, and are discarded by clear/free. Unknown enum tags use only their
tag, never inactive payload bytes.

Owned values (for example nested lists or `strings::Owned`) require explicit
cleanup. Retrieve a displaced value or remove an entry and clean up the returned
owners. Before clear/free, iterate to release any element resources you own.
Iteration copies owners shallowly: cleaning up each copy once is appropriate
only when immediately discarding the corresponding stored entries, with no
subsequent use of those stale handles. See
[ownership conventions](strings-and-memory.md#ownership-conventions-and-transfer).

## Borrowing and invalidation

Reads, iteration, failed operations, and no-op reservation preserve entry
pointers. Replacement updates the existing value in place, preserving its
address and stored key. A successful new insertion may grow or rebuild to
remove tombstones, invalidating **all** entry pointers; otherwise existing
entries stay in place. A successful reserve that increases capacity invalidates
all entry pointers. Remove invalidates only the removed entry; it does not
shift other entries. Clear/free invalidate all entries. Taking a map preserves
entry storage, but an iterator still borrows its original map handle.

Explicit iterators visit every live entry once, with unspecified order, and
allocate nothing. The borrowed map handle must remain at the same live address.
Do not insert (including replacement), remove, reserve, clear, free, transfer,
or otherwise structurally mutate it during iteration. Updating existing values
through `get_ptr` is permitted; subsequent `next` reads the current value.
Assigning a returned key/value copy does not modify the stored entry.

## Allocation and checked failure

Empty construction, reads, removal, iteration, clear/free, and reserve(0) need
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
    iterator := maps::iter(&scores)
    while true {
        key, score, found := maps::next(&iterator)
        if not found {
            break
        }
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
