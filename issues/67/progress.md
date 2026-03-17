# Issue #67 — C backend: skeleton module, compiler wiring, and pure-Go build

## Acceptance criteria

- [x] `internal/codegen/c/` package exists with `NewCG`, `Generate`, and `ExePath` implemented (stub that emits an empty `.c` file)
- [x] `cmd/compiler/main.go` imports the C codegen instead of LLVM
- [x] `tests/compiler/compiler.go` imports and uses the C codegen
- [x] `internal/codegen/llvm/` directory is fully deleted
- [x] `third/go-llvm/` directory is fully deleted
- [x] `go.mod` and `go.sum` contain no LLVM or CGO entries
- [x] `Makefile` `make build` is `go build -o telia ./cmd/compiler` (no CGO flags, no build tags)
- [x] `Makefile` `make test` is `go test ./...` (no build tags)
- [x] `go build ./cmd/compiler` succeeds with no CGO environment variables set
- [x] `go test ./...` runs; all non-codegen tests (lexer, parser, sema) pass

## Decisions

- `Generate()` stub writes an empty `.c` file and returns `nil` with an empty `exePath` — integration tests fail gracefully with a "binary not found" style error, which is expected at this stage.
- `CVariable{CType, Name string}` is created in `types.go` as the placeholder that will replace `*llvm.Variable` in `BackendType` slots in future issues.
- `go.mod` and `go.sum` require no changes — the go-llvm bindings were a local vendored package (same module), not an external dependency.
- The new C backend package is named `c` (import path `internal/codegen/c`) to mirror the existing `llvm` package structure.

## Last checkpoint

COMPLETE. All 10 acceptance criteria satisfied. Integration tests fail gracefully with "exec: no command" (expected — stub produces no executable). Lexer, parser, and sema tests all pass.
