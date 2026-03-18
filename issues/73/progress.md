# Issue #73 — C backend: tuples (multiple return values)

## Acceptance criteria

- [x] Tuple return types emit a `typedef struct { T0 _0; T1 _1; ... } _tuple_<signature>;` in the generated file header
- [x] Identical tuple shapes produce only one typedef (deduplication by type signature)
- [x] Functions returning tuples use the correct typedef'd struct as their C return type
- [x] Calling a tuple-returning function and assigning to multiple variables correctly unpacks `result._0`, `result._1`, etc.
- [x] Tuple literals (e.g. `(a, b)`) construct the correct struct value
- [x] New `TestMultipleReturnValues` integration test passes
- [x] All previously passing tests continue to pass

## Decisions

### Tuple type registry on `ast.Program`
Rather than discovering tuple shapes during codegen (which would require a full AST scan or a `map[string]bool` on `CCodegen`), we added `TupleTypes []*ast.ExprType` to `ast.Program`. Sema's `Check()` pass populates it via `registerTupleType()` — called from `checkFnDecl` (after `RetType` is resolved) and `inferTupleExprTypeWithContext` (after a tuple literal is resolved). Deduplication uses a canonical key: comma-joined element type strings. This gives codegen an exact, pre-computed, pre-deduplicated list — no map on codegen, no discovery pass.

### `sema.program` field
`registerTupleType` needs access to `ast.Program`, but the deep call chain (`checkFile → checkFnDecl`) doesn't thread `program` through. Rather than threading it 4 levels deep, we stored `program *ast.Program` on the `sema` struct and set it at the start of `Check()`. This is consistent with how `s.pkg` and `s.file` are stored.

### `tupleTypedefName` sanitization
C type strings like `int32_t` are already identifier-safe. Pointer types (`void *`, `char *`) and types with spaces are sanitized via `strings.NewReplacer(" ", "_", "*", "ptr")`. This produces names like `_Tuple_int32_t_int32_t` and `_Tuple_int32_t_int64_t`.

### Two paths in `emitVarStmt`
Multi-name declarations come in two shapes:
- `a, b := fnCall()` — `Expr` is `KIND_FN_CALL`; needs a temp variable + field unpacking.
- `c i32, d i64 := 10, 20` — `Expr` is `KIND_TUPLE_LITERAL_EXPR`; `TupleExpr.Type` is nil (sema handles it per-element); emit each var directly from its paired expression.

These are now handled as two distinct `else if` / `else` branches in `emitVarStmt`.

## Last checkpoint

All 7 acceptance criteria satisfied. Full test suite green (12/12 integration tests, sema, parser, lexer). Ready to close.
