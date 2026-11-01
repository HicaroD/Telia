# Shared table engine

[`std::internal::swiss`](../../base/std/internal/swiss/swiss.t) is implementation
plumbing, not a supported application package. Maps wrap `Table<K, V>`; sets
can reuse it with a concrete nonzero-sized empty value record (never assume a
zero-sized backend representation). There are no public hash callbacks.

## Storage and probing invariants

- Two separately allocated, equally sized buffers: byte controls and naturally
  aligned `Entry<K, V>` objects. Only live entry slots are initialized/read.
- Physical buckets are a power of two, at least eight. Group width is **8**.
  Public capacity is `buckets / 8 * 7` entries, a **7/8 occupied load bound**.
- Controls `0..127` are the low seven hash bits; `128` is empty and `129`
  deleted. No raw object memory is hashed. High hash bits select the initial
  aligned group modulo the power-of-two group count.
- Probe groups advance by 1, 2, 3, ... modulo group count (triangular probing),
  visiting every group once. Scalar lane scans are used, not SIMD instructions.
  Every lane in a candidate group is searched before stopping at an empty lane.
- Removal always leaves a tombstone, never an empty hole, preserving lookup
  reachability through collision chains. `occupied` counts live plus deleted
  slots. Reusing a tombstone does not increase it. A new insert checks both live
  capacity and occupied capacity; exhausted occupied allowance rebuilds at the
  same bucket count, discarding tombstones. Clear resets all controls to empty.
- `find` performs semantic equality after fingerprint matching. `vacant` does
  **not** compare keys: rehash directly copies every live entry into a vacant
  slot, including each NaN entry. It cannot call public insert, which would
  change replacement behavior and could coalesce entries under other policies.

## Internal operation contract

`Table<K, V>` zero is heap mode. Its fields are private-by-convention. The
`arena` pointer is the internal allocation-selection boundary used by map
constructors and clones; null selects heap. `allocate<T>`
checks byte arithmetic and selects heap or aligned arena allocation. `release`
frees heap buffers or forgets arena buffers, then zeroes everything **except**
the allocator pointer. Clear retains allocator, storage, and seed. No element
destruction occurs. Arena callers must keep the arena address/lifetime
valid; failed partial arena allocation may consume arena space but must not
publish a changed logical table.

- `find(table, key)` returns slot/found; missing slot is not meaningful.
- `reserve(table, requested)` ensures total live-entry capacity, with a no-op
  when sufficient. `prepare` ensures one further insertion can commit.
- `insert_new(table, key, value)` requires the caller has already rejected
  equal-key replacement, and performs checked preparation then an infallible
  shallow insertion. Failure does not change the table or consume inputs.
- `rebuild(table, buckets)` requires a valid power-of-two bucket count, at least
  eight, large enough for all live entries. It checks byte sizes before seed or
  allocation work. Seed acquisition happens only for a table without buffers;
  both seed and storage publish only after success. Both heap buffers are
  independently failure-tested. Copying into the fresh table cannot fail.
- `vacant(table, hash)` requires allocated controls and at least one reusable
  slot. Fingerprinting/probing use bounded arithmetic, with no dependency on
  wrapping Telia arithmetic. Rehash reuses the retained seed.
- Map replacement touches only the stored value; removal returns the stored
  entry. Iterators retain a pointer to the map handle, scanning controls and
  reading live entries lazily; updating values via get_ptr is permitted.

`clone_in(source, arena)` initializes an empty destination with the selected
allocator and reserves for the source's live count. Empty clones need no storage
or seed. Nonempty clones acquire their own seed, then `copy_live` rehashes semantic
keys into vacant slots **without equality**. Rebuild shares `copy_live`, avoiding
any NaN-coalescing or stored-key replacement path. Allocation failure returns an
empty cleanup-safe destination and preserves the source; partial heap buffers
are released, partial arena allocations remain arena-owned. Clone capacity is
not a copy of source reservation. No set/direct-loop placeholders are exposed.

Set handoff: wrap `Table<T, ConcreteEmptyRecord>` and delegate `clone_in`,
`reserve`, `clear`, and `release`; initialize `.arena` for `create_in`. The record
must have a verified nonzero C representation. Insert must use `find` first to
preserve stored-element identity, then `insert_new` only for absent elements.
Do not clone or rebuild via equality-based public insertion.

## Native test seams and evidence

`_TELIA_KEY_HASH(value)` defaults to identity around the completed SipHash
result. A native forced-include test header may replace it to force identical
fingerprints and final-group starts. `_TELIA_ENTROPY` controls seed failures;
`_TELIA_MALLOC`/`_TELIA_FREE` count allocations and inject first/second-buffer
failures. These are build-time internal seams, never Telia-level callbacks.

`TestMapsNative` runs lifecycle, 3,000-operation reference-model/churn, failure,
key-matrix/NaN, heap/arena clone and allocation-failure, alignment/reset/transfer,
owning-resource, and checked multi-result error fixtures at
`-O0`/`-O3`, with normal and
forced-collision hashes, and UBSan. `TestMapsRejectInvalidKeysAtInstantiation`
checks unsupported-key diagnostics even for unused constructors. The key
classifier recognizes the actual `std::maps::Map` declaration, not unrelated
user records named Map. Broad target qualification/benchmarks remain separate.
