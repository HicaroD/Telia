# Issue #83: @catch operator (inline error handling)

## Acceptance criteria
- [x] Sema validates @catch only on error-returning functions
- [x] Sema rejects @catch on non-error functions with diagnostic
- [x] Error variable is scoped to catch handler block only
- [x] Handler block return type must match enclosing function's return type
- [x] @catch on (T, error) emits: if error -> handler, else -> extract value
- [x] @catch on (T,...,error) emits: if error -> handler, else -> extract values (same pattern)
- [x] @catch on error-only emits: if error -> handler, else -> continue
- [x] Integration test: error_catch_tuple.t passes
- [x] Integration test: error_catch_void.t passes
- [x] Sema test: @catch on non-error fn produces diagnostic

## Decisions
- Handler block return type must match enclosing function's return type strictly (user confirmed)
- Follow @fail patterns in sema (checkFnCall) and codegen (emitAtFailBareCall)
- AST/parser already complete (AT_OPERATOR_CATCH, CatchAtOperator, parseCatchOperator)
- Added returnTy parameter to checkFnCall to enable handler block return type validation
- Added @catch handling in checkVar for tuple unwrapping (same as @fail)
- Variable declared before if/else in codegen to avoid scoping issues

## Changes made
- `internal/sema/sema.go`: Added @catch validation in checkFnCall, threaded returnTy through checkFnCall/checkVar
- `internal/codegen/c/stmts.go`: Added emitAtCatchBareCall, isAtCatchTupleCall, @catch handling in emitVarStmt
- `internal/sema/sema_test.go`: Added 7 @catch sema test cases
- `tests/integration/testdata/error_catch_void.t`: New integration test
- `tests/integration/testdata/error_catch_tuple.t`: New integration test (happy path)
- `tests/integration/testdata/error_catch_tuple_fail.t`: New integration test (error path)
- `tests/integration/compile_test.go`: Added TestErrorCatchVoid, TestErrorCatchTuple, TestErrorCatchTupleFail

## Last checkpoint
- All acceptance criteria met
- All tests passing
- Ready for review
