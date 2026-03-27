# Issue #84: Error handling: error as parameter and struct field type

## Acceptance criteria
- [x] `fn foo(err error)` parses and compiles correctly
- [x] `struct Result { err error }` parses and compiles correctly
- [x] Passing error values to error-typed parameters works
- [x] Accessing error struct fields works
- [x] Integration test demonstrates error as parameter type
- [x] Integration test demonstrates error as struct field type

## Decisions
- Use TDD with small vertical slices: add one failing behavior test, make the minimal fix, then repeat
- Nested field access like `result.err.msg` is the missing behavior; resolve field access one segment at a time in sema and annotate each access node
- Public verification lives in integration tests `error_param.t` and `error_struct_field.t`; sema tests cover the same behaviors more cheaply

## Last checkpoint
- Done. Red-green cycle 1: added sema coverage for `result.err.msg`, saw failure `other field access type was found during sema`, then implemented recursive field-access resolution. Red-green cycle 2: added integration coverage for error-typed struct fields and error-typed parameters. Full suite passes with `go test ./...`.
