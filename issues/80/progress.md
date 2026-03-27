# Issue #80: Error handling: error type, constructor, and field access

## Acceptance criteria

- [x] `error` keyword produces correct token in lexer
- [x] Parser constructs error type in type position (EXPR_TYPE_BASIC with ERROR_TYPE)
- [x] Parser parses `error("msg")` as constructor expression
- [x] Sema rejects user definitions of `error` (struct, fn, type alias)
- [x] Codegen emits `_Error` struct typedef in C preamble
- [x] `error("msg")` compiles to `(_Error){"msg"}`
- [x] `err.msg` compiles to `err.msg` (struct by value)
- [x] `return nil` in error fn emits `return (_Error){NULL};`
- [x] `return val, nil` in `(T, error)` fn emits correct tuple literal
- [x] Integration test: `error_return.t` passes
- [x] Lexer, parser, sema tests pass

## Decisions

- `error` is a basic type keyword (`ERROR_TYPE`) — `IsBasicType()` returns true
- In AST: `error` is `EXPR_TYPE_BASIC` with `BasicType{Kind: token.ERROR_TYPE}` — no synthetic struct in AST
- `_Error` struct typedef exists ONLY in C codegen preamble (C implementation detail)
- `error("msg")` parsed as `KIND_LITERAL_EXPR` producing `(_Error){"msg"}` in C
- `nil` in error position: sema allows via `IsError()` check on BasicType
- `err != nil`: special case in sema `inferBinaryExprType` — error vs nullptr → bool
- Field access on error: special case in `getAccessedField` + `emitFieldAccess` (no StructDecl)
- `currentRetTy` on `CCodegen` to support `return nil` → `(_Error){NULL}`
- Tuple returns: `emitReturn` handles nullptr→error elements inline
- Error-nil comparison in C: emits `e.msg != NULL` (can't compare structs by value)

## Last checkpoint

All acceptance criteria complete. Codegen fully implements error type with `_Error` struct,
error constructor emission, field access, nil returns, tuple returns, and error-nil comparison.
Integration test `error_return.t` covers error functions with nil, tuple (i32, error) returns,
error-nil comparison, and error field access.
