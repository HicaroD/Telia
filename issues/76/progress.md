# Issue #76: C backend — package build mode smoke-test

## Acceptance criteria

- [x] `telia build ./mypackage/` compiles a multi-file package to a working executable
- [x] The executable produces the expected output
- [x] Imported packages (e.g. `std::io`, `std::libc`) resolve and link correctly in package build mode
- [x] All previously passing tests continue to pass

## Decisions

- Used `tests/integration/testdata/pkg_smoke/` as the test fixture (not `examples/sample/`, which has dangling function references).
- `pkg_smoke/main.t` imports `std::io` and `pkg::greet`; `greet/greet.t` imports `std::io` — covers cross-package stdlib resolution.
- Added `CompilePackage(dirPath string)` to `tests/compiler/compiler.go` alongside the existing `CompileFile`; it calls `compilePackagePipeline` which uses `ParsePackageAsProgram` instead of `ParseFileAsProgram`.
- Expected output: `"Hello from greet package!\nHello from main package!\n"` — greet prints first (called first in main), then main's own println.

## Last checkpoint

Complete. All four acceptance criteria satisfied. Full suite (`go test ./...`) passes.
