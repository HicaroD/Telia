# Issue #68 — C backend: type mapping and C compiler invocation

## Acceptance criteria

- [x] `emitCType` handles: `bool`, `i8`–`i128`, `u8`–`u128`, `int`, `uint`, `f32`/`float`, `f64`, `void`, `string`, `cstring`, `rawptr`, pointer types, struct types
- [x] C compiler is auto-detected: `cc` tried first, then `gcc`, then `clang`
- [x] If no compiler is found, error message is `"no C compiler found on PATH; install gcc or clang"`
- [x] Debug builds pass `-O0` to the C compiler
- [x] Release builds (`-release` flag) pass `-O3` to the C compiler
- [x] Generated `.c` file is deleted from temp dir after successful build
- [x] Generated `.c` file is kept when binary is built with `build-dev` (DEV mode)
- [x] An empty but valid `.c` file (with `#include` preamble only) compiles without errors via the detected compiler

## Decisions

- `emitCType` is a package-level function in `types.go` — pure type-string translation, no codegen state needed.
- Tuple types (`EXPR_TYPE_TUPLE`) panic with a clear message pointing to issue #73 — the hook exists but is deferred.
- Preamble includes: `<stdint.h>`, `<stdbool.h>`, `<string.h>`, `<stdio.h>`, `<stdlib.h>`.
- Cleanup (`os.RemoveAll`) happens after successful compilation, not before.
- `c.exePath` is set to the produced executable path.

## Last checkpoint

COMPLETE. All 8 acceptance criteria satisfied. Integration tests now fail with `Undefined symbols: _main` (expected — preamble-only .c has no main; addressed in #69). Lexer, parser, sema all pass.
