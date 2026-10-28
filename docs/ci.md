# Continuous integration

The CI workflow runs separate build, full test suite, formatting, and lint checks
on Ubuntu 24.04 for pull requests and pushes to `master`, or by manual dispatch.

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
