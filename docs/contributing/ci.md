# Continuous integration

The CI workflow runs separate build, full test suite, documentation, formatting, and lint checks
on Ubuntu 24.04 for pull requests and pushes to `master`, or by manual dispatch.
The **Native vendor packages** workflow uses the same triggers for its four
platforms. Both workflows cancel superseded runs for the same branch or PR.
The general CI workflow owns the full Go suite; the native matrix runs the
platform-specific checks rather than repeating the full suite on Linux.

Version tag pushes run **Release** instead of either test workflow. Release
requires successful `master` push runs for the exact tagged commit and reuses
their SDL artifacts. Only versioned packaging, extracted-bundle smoke tests, and
Windows installer checks run again. Missing or expired artifacts block release;
there is no automatic expensive rebuild. See [publishing a release](releases.md).

## Run locally

Install Go (the version in `go.mod`), Git, Make, Python 3, and a C compiler
available as `cc` (GCC or Clang). On minimal Ubuntu installations, use
`sudo apt-get install build-essential`. GitHub's Ubuntu runner already includes
the toolchain; CI verifies that it can compile and execute a program.

From the repository root:

```sh
# Match CI even if a newer Go release is installed locally.
export GOTOOLCHAIN="go$(awk '$1 == "go" { print $2 }' go.mod)"

# Required even for builds: config/env is embedded in the compiler.
printf 'T_STD_PATH=%s/base/std\n' "$PWD" > config/env

make install-format-tools install-lint-tools
export PATH="$(go env GOPATH)/bin:$PATH" # use GOBIN instead if you have set it

make build
make test
make test-docs
make fmt-check
make lint
```

Tool versions are pinned in the Makefile: golines v0.12.2, golangci-lint v1.64.8,
and Go 1.23.12 for building those tools. Installation automatically downloads
that toolchain and requires network access. Project checks use Go from `go.mod`.

`make fmt-check` reports Go and golines formatting differences in tracked Go
files without rewriting them. Use `make fmt` to apply formatting. golines uses
its default 100-column target with `gofmt` as its base formatter.

`make test` includes integration tests that compile generated C and execute the
resulting binaries. SDL-specific tests require separately prepared native
dependencies and `TELIA_TEST_SDL=1`; those remain covered by the existing
`Native vendor packages` workflow. General CI adds no release or publishing steps.

`make test-docs` checks local Markdown links/anchors, tests the link checker, and
compiles/runs complete documentation examples at `-O0` and `-O3`. The Go example
tests also run as part of the full suite. See [example conventions](README.md#example-conventions).

The native matrix also runs `TestUTF8`, `TestString`, and published standard-library
examples on Linux amd64, macOS amd64/arm64, and Windows amd64 (MinGW-w64).
String utility fixtures run at `-O0`/`-O3`, with deterministic allocation failures
and borrowed-storage cleanup checks. Non-Windows hosts additionally run native
sanitizers (UndefinedBehaviorSanitizer by default; use
`TELIA_TEST_SANITIZERS=address,undefined` when AddressSanitizer is available).
Windows skips sanitizer-only tests but runs ordinary behavior/failure checks.
Matrix configuration alone is not evidence of a successful platform run.
The matrix also runs [semantic hashing and seed-adapter checks](hashing.md),
including reference vectors, compiled fixtures, deterministic entropy failures,
and non-Windows UBSan checks.

The same four-target matrix runs maps, sets, direct container loops and combined
reference-model qualification, including `-O0`/`-O3`, deterministic entropy and
partial allocation failures, forced-collision wraparound probing, and non-Windows
UBSan. Windows uses the selected UCRT64 MinGW-w64 GCC toolchain and runs ordinary
checks without sanitizer-only flags. See [qualification evidence](container-qualification.md)
for executed results, pending native CI, acceptance coverage, and benchmark commands.
The selection also includes collection bounds checks at both optimization levels:
Windows runs those without UBSan and skips only the dedicated collection sanitizer
suite. Collection text-output assertions accept mixed LF/CRLF line endings;
byte-preservation tests retain exact byte comparisons.
