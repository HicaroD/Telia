# Issue #69 — C backend: functions, basic expressions, and hello world

## Acceptance criteria

- [x] `TestCompileHelloWorld` passes
- [x] `TestCompileFibonacci` passes
- [x] `TestCompileCalculator` passes
- [x] All existing sema/parser error tests still pass
- [x] Forward declarations are emitted for all functions before any body
- [x] `extern` blocks with `link_name` emit the correct C function name at call sites
- [x] Variadic function calls compile and run correctly
- [x] Variable declarations use the correct C type from `emitCType`
- [x] `BackendType` on `VarIdStmt` and `Param` is set to `*CVariable`

## Decisions

- `c_int`/`c_size_t` are fully resolved by sema to `i32`/`u32` — `emitCType` needs no changes for them.
- `EXPR_TYPE_ALIAS` in `emitCType` panics as a safety net (should never be reached post-sema).
- Name mangling: root `main` package functions use bare names; all other packages use `<pkg>__<fn>`.
- Known-safe libc symbols (printf, puts, malloc, etc.) are skipped in extern forward declarations — already covered by preamble headers.
- `generatePackage` uses `pkg.Processed` guard for DFS deduplication, same as old LLVM backend.
- `currentPkg *ast.Package` added to `CCodegen` to track package context during emission.
- `tmpCnt int` added for unique temp variable names.

## Last checkpoint

COMPLETE. All 9 acceptance criteria satisfied. One unexpected issue found and fixed during implementation: `base/std/libc/stdarg.t` contains LLVM intrinsic externs (`llvm.va_start`, `llvm.va_end`, `llvm.va_copy`) and `base/runtime/mem.t` contains `llvm.memcpy.p0.p0.i64`. These are not valid C identifiers — fixed by adding `isLLVMIntrinsic()` check (any name containing `.`) to silently skip them in `emitExternDecl`.
