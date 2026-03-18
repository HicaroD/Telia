# Issue #75 — C backend: float arithmetic

## Acceptance criteria

- [x] `f32`/`float` type maps to C `float`; `f64` maps to C `double`
- [x] Float literals emit with the correct C suffix (`1.5f` for f32, `1.5` for f64)
- [x] Float arithmetic (`+`, `-`, `*`, `/`) emits correct C expressions
- [x] Float comparisons (`==`, `!=`, `<`, `<=`, `>`, `>=`) emit correct C expressions
- [x] Float parameters and return types use the correct C types in function signatures
- [x] New `TestFloatArithmetic` integration test passes
- [x] All previously passing tests continue to pass

## Decisions

### Sema fix 1 — binary expression type context not propagated
`inferBinaryExprType` accepted `expectedType` but only used it for a strict `Equals` check at the end (sema.go ~1497), never passing it into `ensureBinaryOperatorsAreTheSame`. This meant `a f64 := 1.5 + 0.5` failed (`FLOAT_TYPE != F64_TYPE`) even though single literals coerce fine. Fixed by applying the same `IsCompatibleWith` + coercion path that `inferBasicExprTypeWithContext` uses for single literals — guarded by `!resultBasic.Explicit` so it only fires for untyped literal results.

### Sema fix 1b — mixed-type binary operands (explicit + literal)
`ensureBinaryOperatorsAreTheSame` (sema.go ~1575) only propagated context between operands when `lhsHasContext XOR rhsHasContext` (i.e., one side was a resolved variable). It did not handle `explicit_var OP literal` where the literal's untyped kind differed from the var's explicit kind (e.g. `f32_var + 1.5`). Fixed by coercing the non-explicit literal side to match the explicit side — **only when the non-explicit side's AST node is a `KIND_LITERAL_EXPR`**, to avoid corrupting shared type nodes of named variables.

### Sema fix 2 — `@c` variadic params enforced as `int`
`checkFnCallArgs` ignored `ParamAttributes.C` and type-checked all variadic args against `variadicParam.Type` (`int` for `printf`). The `@c` annotation existed in the AST and parser but had no sema implementation. Fixed by using `inferExprTypeWithoutContext` (resolve without type constraint) for `@c` variadic args instead of `inferExprTypeWithContext`. This keeps field accesses, identifiers, and other complex expressions fully resolved while skipping the `int` constraint.

## Last checkpoint

All 7 acceptance criteria satisfied. Full test suite passes (`internal/sema`, `internal/lexer`, `internal/parser`, `tests/integration` — 11/11 tests green). Ready to close.
