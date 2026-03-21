# #82: @fail operator

## Goal
Implement the `@fail` post-call operator that checks if an error-returning function produced an error, and if so, panics with the error message. On success, extracts non-error values from the tuple.

## Approach
TDD with RED→GREEN vertical slices.

## Status: DONE ✅

### Cycle 1: Sema validation — @fail on non-error function
- **RED**: Added test `@fail on non-error function` to `TestErrorHandling` (sema_test.go) — expects diagnostic
- **GREEN**: Added `returnsError()` helper + @fail check in `checkFnCall` (sema.go:1150-1165)

### Cycle 2: Sema validation — @fail on error functions
- **RED**: Added test `@fail on error-only function` — expects pass
- **GREEN**: Test already passed (returnsError handles error-only and tuple-with-error)
- **RED**: Added test `@fail on (i32, error) function` — expects pass
- **GREEN**: Test passed

### Cycle 3: Codegen — @fail on tuple-returning function (success path)
- **RED**: `TestErrorFailTuple` already failing — C compiler error: `initializing 'int32_t' with '_Tuple_int32_t__Error'`
- **GREEN**: 
  - Added `_panic` function to C preamble (codegen.go:26-29)
  - Added `isAtFailTupleCall()` and `getFullFnRetType()` helpers (stmts.go)
  - Added @fail handling in `emitVarStmt` single-var path (stmts.go:99-130): stores full tuple in temp, checks error field, extracts non-error values
  - Fixed sema: skip `checkVarExpr` for @fail to prevent overwriting unwrapped type (sema.go:642-652)

### Cycle 4: Codegen — @fail panic path (error != nil)
- **RED**: Added `TestErrorFailPanic` — expects non-zero exit + stderr "it failed"
- **GREEN**: Test passed (codegen from cycle 3 handles panic path)

### Cycle 5: Codegen — bare @fail on error-only function
- **RED**: Added `error_fail_void.t` + `TestErrorFailVoidPanic` — expects non-zero exit + stderr "boom"
- **GREEN**: Added `emitAtFailBareCall()` (stmts.go) + @fail branch in KIND_FN_CALL handler (stmts.go:65-70)

### Cycle 6: Codegen — bare @fail success path
- **RED**: Added `error_fail_void_ok.t` + `TestErrorFailVoidOK` — expects "ok\n" (no panic)
- **GREEN**: Test passed (emitAtFailBareCall correctly skips panic when msg is NULL)

## Files changed
- `internal/sema/sema.go` — returnsError(), unwrapErrorFromTuple(), @fail in checkFnCall, @fail in checkVar
- `internal/sema/sema_test.go` — 3 @fail sema tests
- `internal/codegen/c/codegen.go` — _panic in preamble, removed debug print
- `internal/codegen/c/stmts.go` — @fail in emitVarStmt, emitAtFailBareCall(), KIND_FN_CALL handler
- `tests/integration/compile_test.go` — 4 integration tests
- `tests/integration/testdata/error_fail_tuple.t` — success-path test source
- `tests/integration/testdata/error_fail_panic.t` — panic-path test source (tuple)
- `tests/integration/testdata/error_fail_void.t` — panic-path test source (error-only)
- `tests/integration/testdata/error_fail_void_ok.t` — success-path test source (error-only)
- `examples/error.t` — updated to use @fail
