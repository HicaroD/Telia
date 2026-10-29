# Bottleneck investigations

Record measured bottlenecks here, distinguishing observations from hypotheses.
No optimization is justified solely by an expensive-looking source pattern.

## Initial compiler investigation

Baseline: `results/20261029T215455Z-70229/`, based on commit
`972d552c34868e2cd87e1721725291498e5f56b9` with the new profiling harness in the
worktree. Apple M2, macOS/arm64, Go 1.26.10, Apple Clang 17, release `-O3`.
Five samples per size after one warmup; frontend CPU sampled separately for 3s.
All generated executables passed their result checks. These are exploratory
measurements, not optimization results.

### Phase baseline

Medians, milliseconds except allocation volume:

| Functions | Parse | Semantic analysis | Backend total | Toolchain subset | Go allocated/build |
| ---: | ---: | ---: | ---: | ---: | ---: |
| 100 | 4.53 | 0.38 | 178.28 | 169.74 | 2.61 MiB |
| 500 | 19.34 | 0.89 | 217.51 | 210.28 | 12.52 MiB |
| 1,000 | 33.19 | 1.64 | 272.67 | 263.88 | 25.14 MiB |

See the generated `phases.svg` for sample ranges. Phase medians do not necessarily
sum to the median total. This workload shows roughly proportional allocation
growth, not evidence of a quadratic compiler bottleneck.

### 1. External compilation dominates build latency

- **Evidence:** toolchain median is about 264ms at 1,000 functions, compared
  with about 35ms for parsing and semantics combined. The last build's toolchain
  log contains one `cc -O3 ... -lm` invocation, not repeated target probes.
- **Why:** every build emits C and invokes a native optimizing compiler and
  linker. `CCodegen.Generate` compiles the emitted translation unit each time;
  the harness does not benefit from a compiled-output cache. The toolchain's
  internal CPU hotspots were not profiled, so optimization versus linking cost
  is not yet established.
- **How to investigate:** retain one emitted C file in a controlled diagnostic
  run, then measure compile-only and link-only commands separately. Compare debug
  and release latency as a product tradeoff, not a like-for-like speedup.
  Evaluate content-addressed output caching only after measuring repeated-build
  frequency; cache keys must include generated C, compiler identity, flags,
  headers, and native dependencies.
- **Expected impact:** frontend-only changes have limited influence on total
  release build time here. A native cache could help unchanged builds but would
  not make genuinely changed cold builds faster. No caching change is implemented.

### 2. Repeated token creation is the largest measured allocation source

- **Evidence:** the sampled cumulative allocation profile attributes **77.37%**
  of allocated bytes directly to `Lexer.Next`, and **9.81%** to
  `Collector.RegisterSpan`. Source-level `pprof -list` identifies
  `tok := &token.Token{}` at `internal/lexer/lexer.go:126` as the allocation site.
  Peek paths cumulatively account for substantial allocation volume; these
  cumulative percentages overlap and must not be added.
- **Why:** `Peek`, `Peek1`, and `PeekN` call `Next` and restore the lexer position.
  Every speculative visit creates another token and registers its span; consuming
  the same token later lexes it again. The precedence parser repeatedly asks
  about upcoming tokens. Span registration also performs map writes, including
  repeated writes for the same source offset.
- **How:** first count consumed versus speculative `Next` calls in a diagnostic
  build. Prototype a lookahead buffer or token cache so repeated peeks can reuse
  lexed tokens. Preserve newline mode, brace-depth transitions, diagnostic
  behavior, and token lifetime/identity expectations. Investigate avoiding
  redundant span registration separately, with diagnostic tests.
- **Validation:** compare parse time and allocated bytes/build on the same
  scaling workload, then test lexer/parser/diagnostic behavior and full integration
  tests. Add lookahead-heavy and error-heavy inputs before accepting the change.
  The allocation percentages are not predicted runtime speedups.

### Profile interpretation and limits

The short CPU sample has substantial Go runtime/OS activity (`madvise`, scheduler,
and synchronization frames). It is not sufficient to assign a precise CPU
speedup to token caching. The allocation result is a clearer first lead; confirm
with longer repeat profiles and source-level stacks. Allocation totals in pprof
are cumulative across hundreds of frontend passes, not per-build or peak memory.

Imports, generics, standard-library-heavy builds, real applications, and generated
program runtime remain uninvestigated. At this baseline no compiler optimization
had been applied; the lookahead experiment below follows up on that finding.

## Buffered lookahead experiment

### Change and correctness

`Lexer` now buffers upcoming tokens on demand. `PeekN` extends the buffer only
as far as requested; `Next` consumes a buffered entry, restoring its ending
position and applying brace-depth changes. Speculative diagnostics are stored
with their token and reported once on consumption. A newline-mode change
invalidates unconsumed entries and restarts from the consumed cursor.

Consumed buffer slots are reclaimed without recycling token objects: pointers
already returned to parser/AST users remain valid. Lookahead beyond EOF reuses
the final entry instead of allocating an arbitrarily long suffix. Scanning a
buffered token registers its span once unless newline-mode invalidation requires
re-lexing. This avoids the repeated registrations caused by ordinary peeking.

Regression tests cover mixed/repeated lookahead, cursor/span preservation,
brace depth, newline-mode changes, token lifetime, EOF, and deferred/multiple
diagnostics. `go test ./...` passed, including the integration suite.

### First before/after measurement

Same machine/toolchain as above; 10 samples per size after a warmup, separate
10-second frontend CPU sampling. Only the lexer implementation changed in compiler
code; the C backend and workload generator are unchanged.

- Before: `results/20261029T220014Z-77862/`
- After: `results/20261029T220349Z-87468/`
- Chart: `results/20261029T220349Z-87468/comparison.svg`

All values below are medians. Total build time is the median of sample totals,
not the sum of phase medians.

| Functions | Parse before → after | Parse reduction | Go allocated/build before → after | Allocation reduction | Total build before → after |
| ---: | ---: | ---: | ---: | ---: | ---: |
| 100 | 4.19 → 1.66 ms | 60.4% | 2.61 → 1.02 MiB | 61.0% | 180.94 → 181.31 ms |
| 500 | 19.71 → 7.97 ms | 59.6% | 12.52 → 4.65 MiB | 62.9% | 236.81 → 231.96 ms |
| 1,000 | 33.01 → 14.35 ms | 56.5% | 25.14 → 9.42 MiB | 62.5% | 309.29 → 294.75 ms |
| 2,000 | 53.22 → 25.81 ms | 51.5% | 50.43 → 19.01 MiB | 62.3% | 448.92 → 428.52 ms |

Parsing became approximately 2.1–2.5× as fast. At 1,000 and 2,000 functions,
observed total build times decreased about 4.7% and 4.5%, respectively. At 100
functions there was no measurable overall improvement. Native toolchain medians
were 0.7–2.8% higher after the change, so total-build changes should not be
attributed entirely to the lexer or generalized to other machines/workloads.

Semantic analysis medians increased from 0.86 → 1.19ms, 1.43 → 1.99ms, and
2.29 → 3.16ms at the three larger sizes. No semantic-analysis code changed;
the cause (GC scheduling, object layout/lifetime, or ordinary measurement effects)
is not established. This is a small absolute regression, but must not be hidden
by the larger parsing improvement. Compare combined frontend costs too.

### Confirmation run

Repeated the same settings without further lexer changes:
`results/20261029T220433Z-88112/`. Its `comparison.svg` uses the same baseline.

| Functions | Parse median | Parse reduction vs baseline | Allocation reduction vs baseline | Total build median |
| ---: | ---: | ---: | ---: | ---: |
| 100 | 1.60 ms | 61.9% | 61.0% | 179.27 ms |
| 500 | 7.95 ms | 59.7% | 62.9% | 232.02 ms |
| 1,000 | 14.86 ms | 55.0% | 62.5% | 297.78 ms |
| 2,000 | 25.36 ms | 52.3% | 62.3% | 429.15 ms |

The repeat supports the parsing and allocation improvements. Total-build
reductions at 1,000/2,000 functions were about 3.7%/4.4%; external-toolchain
variation limits precision. Semantic-analysis medians remained higher at
1.17/2.01/3.11ms for 500/1,000/2,000 functions, so the small absolute regression
is reproducible and its cause remains open. Even including semantics, median
sampled combined frontend time is substantially lower; no claim of a
semantic-analysis improvement is made. Two post-change rounds versus one baseline
round are useful confirmation, not a randomized or confidence-interval study.

### Remaining bottlenecks and limits

In the new sampled allocation profile, direct token allocation is about 30.8%
and span registration about 31.3% of cumulative allocation volume. These shares
are not directly comparable absolute byte counts across profiles: a fixed-time
CPU profiling loop completes more iterations with a faster frontend. The
per-build allocation column above is the useful absolute comparison.

The native compiler/linker is still the dominant total-build cost. Span-map
storage and required token/AST allocation are next frontend leads, not permission
to remove diagnostic information or reuse token objects unsafely. CPU sampling
still has substantial runtime/OS activity; it does not establish precise
function-level CPU savings.

These synthetic workloads establish a benefit for function-heavy builds, not
all Telia programs. Imports, generic-heavy code, error-heavy parser workloads,
and runtime execution need separate performance investigations. No generated
program runtime speedup is claimed.

### Finding template

- **Workload and environment:** result directory, commit, parameters, machine.
- **Evidence:** phase distribution, scaling, profile stacks/source lines.
- **Why:** explain the repeated work, allocation, algorithm, or toolchain cost.
- **How:** smallest targeted improvement, expected impact, correctness risks.
- **Validation:** identical before/after measurements, tests, and limitations.
