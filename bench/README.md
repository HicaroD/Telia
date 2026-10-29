# Performance investigations

Find Telia bottlenecks, explain why they occur, and validate improvements. This
is not a language comparison suite. Historical C/Go Fibonacci examples and
`imgs/` are retained, but are not used by the runner.

## Run

Requires Go, a C compiler available as `cc`, and Python 3; no Python packages.
Run from any directory:

```sh
./bench/bench.sh
./bench/bench.sh -runs 10 -sizes 100,500,1000,2000 -profile-time 10s
```

The shell runner builds `profile.go` and creates a unique `results/<UTC>-<PID>/`:

- `results.json`: every sample, commit/worktree state, toolchain and machine details.
- `functions_*.t`: generated scaling workloads, with all functions called and a
  checked result. Optimizing the final executable is fine: this tests compilation.
- `phases.svg`: phase medians, total ranges, toolchain time, allocation volume.
- `cpu.pprof` / `cpu.txt`: frontend-only CPU sampling, with parse/sema labels.
- `allocs.pprof` / `allocs.txt`: sampled cumulative Go allocation volume.
- `toolchain.jsonl`: commands from the last build, including target probes.

The growing workload stresses parsing, symbol lookup, function checking, and
C emission. It does **not** represent imports, generics, collections, strings,
or runtime behavior; add targeted workloads when investigating those areas.

## Interpret correctly

Each size gets one discarded warmup and then measured in-process builds. A GC
runs before each build, outside timing; filesystem/toolchain caches remain warm.
The harness bypasses CLI configuration/startup. It checks the executable exits
successfully outside build timing. No execution-speed claim is made.

`backend_ms` includes C generation, target probes, compilation/linking, and output
deployment. `toolchain_ms` measures subprocess wall time inside a `cc` wrapper,
including probes; it is a **subset**, not another additive phase. The difference
also includes wrapper startup and filesystem/process overhead: do not call it
pure C-emission time. Child compiler CPU/memory is not in Go profiles.

`allocated_bytes` counts total Go bytes allocated during a build, not peak memory
or native compiler memory. Allocation profiles cover the entire harness process
and use Go's default sampling rate. CPU profiles are collected separately by
repeating the largest frontend workload; they do not distort build samples.

## Diagnose why, then test how

```sh
go tool pprof -http=127.0.0.1:8080 bench/results/<run>/profile bench/results/<run>/cpu.pprof
go tool pprof -top -cum bench/results/<run>/profile bench/results/<run>/cpu.pprof
go tool pprof -list 'parse|check' bench/results/<run>/profile bench/results/<run>/cpu.pprof
go tool pprof -top -alloc_space bench/results/<run>/profile bench/results/<run>/allocs.pprof
```

1. Locate the dominant phase and check how it scales with workload size.
2. Inspect cumulative call stacks and source lines, not just flat rankings.
3. Record evidence, cause, proposed fix, and tradeoffs in `analysis.md`.
4. Make one targeted change, run correctness tests, and repeat identical settings.
5. Compare per-phase distributions and allocation volume, not single timings.

Generate a before/after chart of parse time and allocation volume:

```sh
python3 bench/charts.py bench/results/<after>/results.json \
  --baseline bench/results/<before>/results.json
```

This writes `comparison.svg` into the after-run directory. Compare the same
workload sizes and settings on the same machine, without competing benchmark jobs.

If the external toolchain dominates, inspect its recorded commands and generated
C before optimizing Telia's frontend. If a runtime workload is slow, use native
profiling (Instruments on macOS, `perf` on Linux) and inspect the lowering and
standard-library implementation. Do not substitute Go compiler profiles for
generated-program profiles.

Local results are ignored by Git. Keep reviewed findings in `analysis.md` and
explicitly publish selected raw results/charts when a finding needs a baseline.
