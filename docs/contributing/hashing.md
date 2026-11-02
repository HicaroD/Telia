# Internal semantic hashing

This compiler/library boundary prepares automatic keys for containers. It is
**not an application API**, customization protocol, randomness package, or
persistent serialization format. There are no user hash callbacks. Heap
[maps](../stdlib/maps.md) consume this boundary through the
[shared engine](maps-engine.md). Applications should not depend on these reserved
internal names.

## Handoff interface

- `__key_hash<T>(value T, k0 u64, k1 u64) u64` requires one concrete type argument.
  Generic library code can pass its specialized key parameter. Sema applies
  `ast.EqualityEligibility(T, true)` before lowering; nested failures identify
  the offending component. The value and then both seed words are evaluated
  once, in source order, at the expression site.
- `__key_seed() (u64[2], error)` obtains sixteen OS entropy bytes, without heap
  allocation. Success returns two seed words and `nil`. Failure returns a zero
  seed and a classified `std::errors::Error::RandomnessFailed`, never a weak
  fallback. `@try`/`@fail` work on the checked result. Callers must not use the
  zero failure result as a seed. A successful all-zero random seed is valid.
- Future containers should acquire and retain their own secret seed lazily,
  before committing first storage; no process-global seed or initialization is
  required. Empty construction/reads need not call this boundary. Container
  allocation, growth, atomicity and lifetime policies are defined by the engine.

Both operations lower directly to private generated C. No new standard-library
package or public custom hasher is provided. Native overrides are build-time test
seams, not Telia-level callbacks.

## Encoding and equality

For a fixed concrete key type and seed, equal values hash equally, following the
[eligibility matrix](../reference/expressions.md#equality-and-key-eligibility).
There is no cross-type hash identity or collision-free promise. Hashes and future
iteration order are nonpersistent: compiler/runtime versions, targets and seeds
may change them. Do not save hashes to disk, transmit them as identifiers, or
expose container seeds to untrusted input.

The compiler streams semantic components into SipHash-1-3:

- Integers (including signed/unsigned 128-bit and C-native scalars) encode their
  value in fixed-width little-endian bytes. Signed conversion is explicitly
  unsigned before shifting. Booleans encode their numeric value.
- Floats use their IEEE bits, but both signs of numeric zero encode positive
  zero. NaNs and infinities are safe to hash; NaN retains ordinary unequal-to-self
  equality, so no equal-NaN promise is invented.
- Strings encode an eight-byte length followed by exactly that many bytes,
  including embedded zeros and invalid UTF-8. This distinguishes adjacent
  variable-length fields. Empty strings never dereference their backing pointer.
- Object pointers, `rawptr`, C string pointers, and inferred error kinds encode
  identity through `uintptr_t`; null is valid. They never read pointees.
- Arrays, tuples and records stream elements/fields in declaration order for
  their known concrete shape. Empty records/arrays contribute no bytes. No
  padding, dummy fields, string backing addresses, or ownership storage is read.
- Enums encode their tag and only the selected variant's fields. An undeclared
  payload-enum tag encodes **only its tag**, even with invalid inactive payload
  bytes. Header-owned payload-free enums use their native integer value/width.

`runtime_hash.h` implements unsigned 64-bit wrapping additions and nonzero,
fixed rotations. Streaming uses byte loads and explicit shifts, not unaligned
loads or host-endian casts. Telia runtime arithmetic is not assumed to wrap.

## Entropy and verification

`runtime_seed.h` uses Linux `getrandom` (retry interruptions and short reads;
reject zero/error), macOS `getentropy` (retry interruptions), and Windows
`BCryptGenRandom` with the system-preferred provider (MinGW-w64). The Windows
adapter resolves bcrypt from System32 at runtime and checks module, symbol and status failures,
so applications need no additional link flags. Unsupported systems fail closed.
There is no time, address, fixed-key, or noncryptographic fallback.

Native tests can define `_TELIA_ENTROPY(bytes, size)` before runtime emission:
return zero only after filling all requested bytes, nonzero for failure. This
internal seam tests success and partial-write failures deterministically; failed
bytes are never published. Adapter tests additionally inject individual OS calls
to exercise interruption, partial read, zero read and native failures. These
simulations are not claims of native execution on the simulated platform.

Run from the repository root with the ignored `config/env` set for this checkout:

```sh
go test ./internal/codegen/c ./internal/sema ./tests/integration -run 'TestSipHash13|TestHash|TestSemanticHashing' -count=1
```

Tests cover independent rust-siphash SipHash-1-3 vectors (key bytes `00..0f`,
message bytes `00..length-1`), split writes, unaligned inputs, block boundaries,
multi-block inputs, semantic generated-C fixtures at `-O0`/`-O3`, deterministic
entropy errors and UBSan on non-Windows hosts. On macOS with an x86-64 compiler
and Rosetta, `TELIA_TEST_HASH_AMD64=1 go test ./internal/codegen/c -run
'TestSipHash13|TestHash' -count=1` additionally executes x86-64 native harnesses.
The native CI matrix runs the tests on all four supported targets; configuration
is not evidence that those runs have completed. ASan requires a working host
runtime and is not implied by the default UBSan checks.

## Local implementation evidence (2026-11-01)

Checkout: `feat/111-seeded-hashing`, based on integration `925bc40` with #110.
Host: macOS arm64, Apple Clang 17.0.0 (`clang-1700.6.3.2`). The ignored
`config/env` points at this worktree's `base/std`.

Completed checks:

- `make build`, `python3 scripts/check-docs.py`, `make fmt-check`: passed.
- `go test ./internal/codegen/c ./internal/sema ./tests/integration -run
  'TestSipHash13|TestHash|TestSemanticHashing' -count=1`: passed, including
  native entropy, adapter simulations, debug/optimized fixtures and UBSan.
- `TELIA_TEST_HASH_AMD64=1 go test ./internal/codegen/c -run
  'TestSipHash13|TestHash' -count=1`: passed as x86-64 native harnesses under
  Rosetta, including actual macOS entropy and UBSan vectors. This is not a full
  macOS amd64 compiler/integration-matrix qualification.
- `GOTOOLCHAIN=go1.22.1 go test ./... -timeout=30m`: passed; integration completed
  in 523 seconds. An earlier `go test ./...` using host Go 1.26.10 exhausted the
  default ten-minute limit late in integration tests; it is not counted as a pass.
- `go test ./tests/integration -run '^TestDocumentationExamples$' -count=1
  -timeout=10m`: passed.
- `GOTOOLCHAIN=go1.23.12 go run
  github.com/golangci/golangci-lint/cmd/golangci-lint@v1.64.8 run --config
  .golangci.yml`: passed. Plain `make lint` was blocked by an installed v2 tool
  rejecting the repository's v1 config; the pinned command above is the result
  used here.

No Linux or Windows native runner was available locally (Docker's daemon was
not running; no MinGW compiler was installed). Their adapter simulations passed,
and native CI is configured, but real platform qualification remains pending.
ASan was not rerun because prior host probes hung even for minimal C; only the
actual UBSan results above are claimed. No background suite remains pending.

## Windows CI harness correction (2026-11-01)

PR #117's first native CI run (`41362041449`, Windows job `125014272346`)
failed in `TestHashSeedAdapters/linux`: the simulation removed `_WIN32`
before including MinGW's `sys/types.h`, which rejects non-Win32 targets.
The same job's native hashing integration selection passed; that does not make
the job or the full Windows matrix a pass.

The harness now loads `sys/types.h` and `wchar.h` under the host target macros
before selecting any simulated adapter. Guarded forwarding headers reproduce
the invalid include ordering on any host. All Linux, macOS and Windows adapter
control-flow checks remain enabled at `-O0` and `-O3`; native entropy code is
unchanged. The regression failed before the include-order fix and passed after
it on macOS arm64. Focused hashing tests, build, documentation examples, link
checks and formatting also passed locally. No MinGW cross-compiler was
available; native Windows verification of this correction remains pending CI.
ASan was not retried because of the recorded host hangs.
