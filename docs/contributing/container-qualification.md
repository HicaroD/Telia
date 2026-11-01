# Map and set qualification (#109 / #116)

This is an evidence ledger, not a declaration that #109 is complete. The local
qualification branch descends from integration `901bd23` (the completed #115
merge). No public hashing policy or container API was added by qualification.

## Acceptance audit

The rows follow **every parent acceptance criterion in order**. Evidence names
are executable tests or fixtures, not inferred success from completed siblings.
Platform execution and final verification remain separate gates below.

| #109 criterion | Evidence / disposition |
| --- | --- |
| 1. Automatic eligible keys | `maps_keys.t`, `sets_keys.t`, `hashing.t`: bool, every integer width including i/u128 and native aliases, floats, strings, pointer/null identity, arrays/tuples, empty/nested generic records, payload/payload-free and concrete error enums. No callbacks or annotations. |
| 2. Complete eligibility matrix / diagnostics | [Equality matrix](../reference/expressions.md#equality-and-key-eligibility), `TestEqualityRejectsUnsupportedNestedComponents`, `TestHashKeyDiagnostics`, `TestMapsRejectInvalidKeysAtInstantiation`, `TestSetsRejectInvalidElementsAtInstantiation`. Foreign/opaque/error/function/slice/owner policies are explicit; unused constructors are checked. |
| 3. Structural equality / active storage / evaluation | `TestStructuralEquality`, `TestStructuralEqualityWithNativeSanitizers`, `equality.t` and native poisoned-padding/inactive-payload fixtures; complement, operand order and single evaluation. |
| 4. Edge keys including nonreflexive floats | Map/set key and clone fixtures cover signed zero, distinct NaNs through growth/clone/iteration, infinities, empty/embedded-zero/malformed byte strings and separate allocations, null/identity pointers, empty/nested generics. `containers_model.t` explicitly models NaNs as an unreachable multiset, not a reflexive reference-map key. |
| 5. Full API / inference / zero / transfer | `maps.t`, `sets.t`, `maps_arena.t`, `sets_arena.t`, `maps_checked.t`, published examples, `TestContainerMultiResultReassignment`. A silently skipped, already-accepted multi-result reassignment was fixed; repeated `next`, `get`, `remove` and full fallible results execute once. Error-handled reassignments remain explicitly unsupported; fresh checked declarations work. |
| 6. Stored identity / explicit resource cleanup | `maps_ownership.t`, `sets_ownership.t`, `maps_clone.t`, `examples/maps.t`, `examples/sets.t`: replacement/duplicate/removal identity, nested owners and one cleanup owner for shallow clones. |
| 7. Direct / explicit iteration and diagnostics | `TestContainerLoops*`, `container_loops.t`, reassignment fixture, existing array/slice and loop-control regressions: empty/exhausted, nested, single source evaluation, copied bindings, get_ptr updates, break/continue/return/@try/defer, invalid arity/temporary/const/lookalike diagnostics. |
| 8. Models / collisions / deletion / capacity | Existing 3,000-step map/set models plus `containers_model.t`: three reproducible seeds (19/38/57), 2,000 combined operations per seed per allocator, clone/reserve/clear/reuse and NaNs. Native runs repeat with normal hashes and forced identical fingerprints/final-group starts, long triangular probe chains, tombstones and wraparound. The internal seam is never an application hash callback. |
| 9. SipHash reference / native arithmetic | `TestSipHash13ReferenceVectors`, `TestSemanticHashing`, `TestHashSeedNative`, `TestHashSeedAdapters`: independent vectors, unaligned/boundary/multiblock streams, unsigned fixed-width arithmetic, adapters and UBSan. See [hashing](hashing.md). |
| 10. Entropy / partial allocation failures | `maps_failure.t`, `sets_failure.t`, clone/arena fixtures under `maps_control.h`: classified entropy failure, first/second buffers, clone/growth/tombstone cleanup, overflow, source preservation, outstanding pointer/iterator validity, allocation counters and leak/double-free checks. Forced-collision runs repeat these failures. |
| 11. Heap / arena ownership / alignment | Map/set arena, ownership and clone fixtures, existing collection/list regressions: aligned storage, retained allocator, old arena buffers, repeated free/reuse, null selects heap, reset/destruction and transfer. Partial arena allocation is logically atomic, not rolled back. |
| 12. Supported native targets / optimizations / sanitizers | **Configured, not all executed:** `.github/workflows/native.yml` has Linux amd64, macOS amd64/arm64, Windows amd64 UCRT64 MinGW-w64. Container steps run debug/optimized normal/forced-collision fixtures; non-Windows defaults to UBSan. Native CI on the exact pushed revision must still succeed. No ASan pass is claimed. |
| 13. Reproducible performance / memory baseline | Native benchmark below covers both wrappers, all six operation families, three key families and three sizes. Actual raw measurements are retained; there is no pre-feature map/set baseline or cross-language parity claim. |
| 14. Authoritative docs / runnable examples | Package references/index, equality/generics, iteration and ownership rules are linked from `docs/agents.md`. Qualification corrected the engine's stale “direct loops later” statement and documented repeated explicit-iterator assignment. Examples use singleton output or commutative totals, never unspecified order. |
| 15. Required verification | Commands and actual outcomes are recorded below. Platform CI remains a delivery gate regardless of local results. |

## #116 deliverable audit

- Combined adversarial/reference models: new cross-wrapper finite-key/NaN model
  composes heap/arena allocation, clones and clear/reservation/churn; existing
  focused failure models remain part of the same native qualification command.
- Progress/invariants: forced final-group collisions exercise wraparound and long
  probes, with bounded process timeouts. Failure headers count each live heap
  allocation and abort on leaks, while UBSan checks generated native accesses.
- Supported targets: workflow configured; actual remote execution is pending.
- Benchmarks: commands, scope, baseline, environment and actual measurements below.
- Docs/lifetimes: coherent source contracts audited; runnable examples checked.
- Full checks: local results below, not a substitute for pending native CI.

## Benchmark method and baseline

```sh
TELIA_CONTAINER_BENCH=1 GOTOOLCHAIN=go1.22.1 go test ./tests/integration \
  -run '^TestContainerBenchmarks$' -count=1 -v
```

This opt-in test compiles six native binaries at **`-O3`**, then times only native
operations (not compiler startup, process startup or key construction). Each size
runs 32 repetitions. Timings are process CPU time from C `clock()`, expressed in
nanoseconds per entry; small cases are sensitive to clock resolution, thermal
state and instrumentation overhead. A checksum consumed by native output keeps
lookups/iteration observable. All table allocations must be freed at exit.

Keys: u64 sequential integers; distinct 12-byte `key-00000000` strings in stable
static buffers; records `{index u64, label string}` with the literal `record-key`.
Map values are u64. Sizes: 64, 1,024, 8,192 live entries. Heap mode only; arena
retained-space costs are documented by contract but not measured by this harness.
The benchmark-only seed is fixed to sixteen `0x71` bytes, avoiding OS entropy and
seed variance; production still uses OS randomness. Default collision behavior is
measured, not the forced-collision test policy.

Phases: insertion into a pre-reserved empty table; successful lookup; full direct
iteration; removal of every entry; insert/remove pairs on the emptied tombstone
table; and insertion from a zero-capacity table including growth/rehashing.
Churn is **ns/pair**, not ns per individual operation. Reservation and key setup
are outside reserved-insert timing; growth includes table allocation but uses
the fixed seed seam. Memory is peak simultaneously live **requested table bytes**
over all phases, including old/new buffers during growth, not RSS, allocator
metadata, key buffers, stack or compiler memory. The memory-counting allocator
adds overhead and preserves `max_align_t` alignment.

This is the first recorded Telia container benchmark baseline, not a claimed
speedup over an unmeasured historical version. Reserved insertion is the measured
within-revision baseline for assessing whole-table growth cost. There are no
throughput targets, Go/Rust parity claims, or API redesigns based on these numbers.
Raw rows preserve actual precision, not measurement certainty.

### Actual baseline (2026-11-01)

[Raw CSV](container-benchmark-baseline.csv): one complete native invocation, all
108 rows, before concurrent full-suite/Rosetta work. Values below are ns/entry,
except churn ns/pair; peak is bytes. These are observations, not stable targets.

| Container | Key | Size | Insert | Lookup | Remove | Iterate | Growth | Churn pair | Peak bytes |
| --- | --- | --- | --- | --- | --- | --- | --- | --- | --- |
| maps | integer | 64 | 118.164 | 40.039 | 36.133 | 16.113 | 170.410 | 55.176 | 3264 |
| maps | integer | 1024 | 91.705 | 27.771 | 28.198 | 4.242 | 149.109 | 40.009 | 52224 |
| maps | integer | 8192 | 50.095 | 29.915 | 30.258 | 3.292 | 104.206 | 21.618 | 417792 |
| sets | integer | 64 | 121.582 | 44.434 | 49.316 | 11.719 | 208.008 | 80.566 | 3264 |
| sets | integer | 1024 | 101.318 | 37.537 | 38.361 | 4.211 | 184.357 | 52.338 | 52224 |
| sets | integer | 8192 | 51.701 | 31.334 | 32.131 | 3.078 | 108.887 | 23.026 | 417792 |
| maps | string | 64 | 248.535 | 131.348 | 131.348 | 15.625 | 493.164 | 227.051 | 4800 |
| maps | string | 1024 | 215.790 | 111.481 | 112.091 | 3.967 | 423.004 | 174.255 | 76800 |
| maps | string | 8192 | 107.189 | 60.802 | 60.074 | 2.453 | 211.327 | 83.172 | 614400 |
| sets | string | 64 | 250.488 | 132.812 | 128.906 | 14.648 | 498.047 | 230.469 | 4800 |
| sets | string | 1024 | 216.766 | 110.382 | 113.312 | 3.174 | 417.603 | 171.509 | 76800 |
| sets | string | 8192 | 106.346 | 60.272 | 60.654 | 2.098 | 210.823 | 83.179 | 614400 |
| maps | record | 64 | 248.047 | 137.695 | 124.512 | 11.719 | 449.707 | 215.820 | 6336 |
| maps | record | 1024 | 204.468 | 107.880 | 109.528 | 4.120 | 396.454 | 160.797 | 101376 |
| maps | record | 8192 | 104.401 | 62.698 | 62.027 | 2.514 | 204.315 | 79.082 | 811008 |
| sets | record | 64 | 242.188 | 130.371 | 135.254 | 12.207 | 456.055 | 225.098 | 6336 |
| sets | record | 1024 | 210.907 | 111.725 | 112.762 | 3.387 | 406.525 | 163.910 | 101376 |
| sets | record | 8192 | 104.893 | 61.275 | 61.184 | 2.083 | 204.842 | 79.304 | 811008 |

## Local environment and verification

Host: Apple M2, macOS Darwin 25.5.0 arm64; Apple Clang 17.0.0
(`clang-1700.6.3.2`), project Go toolchain `go1.22.1`. Host default Go is 1.26.10;
the required suite uses the project version explicitly. `config/env` is ignored
and points to this qualification checkout's `base/std`.

ASan is intentionally **not retried**: prior minimal-C host probes hung. Local
native sanitizer evidence is UBSan plus deterministic allocation accounting,
not ASan or general leak-sanitizer coverage. No Linux/Windows native success is
inferred from adapter simulations or from workflow configuration.

## Coordinator handoff / remaining gate

After reviewing and pushing the qualification branch, open a PR targeting the
integration branch. The existing `pull_request` trigger must execute the Native
vendor packages matrix on the reviewed SHA. Retain all four target results,
including the “Seeded semantic hashing and native entropy” and “Map and set
qualification” steps. Each must finish successfully; verify Windows selects
UCRT64 MinGW-w64 and macOS has both Intel and ARM runners. The general CI full
suite/docs/format/pinned lint must also pass on that SHA. Configuring a matrix
is not completion of the supported-target acceptance criterion. No workflows
were dispatched, PRs/issues changed, branches pushed, or integration merges made
by this qualification session.
