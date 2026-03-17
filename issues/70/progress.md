# Issue #70 — C backend: control flow (if/elif/else, for, while)

## Acceptance criteria

- [x] `if` blocks with no `else` emit correct C `if` statement
- [x] `if`/`else` blocks emit correct C `if`/`else`
- [x] `elif` chains emit correct C `else if` chains
- [x] `for` loops (init; cond; update) emit correct C `for` loop
- [x] `while` loops emit correct C `while` loop
- [x] New `TestForLoop` integration test passes
- [x] New `TestWhileLoop` integration test passes
- [x] All previously passing tests continue to pass

## Decisions

- All C emission logic already exists in stmts.go from issue #69. This issue only adds test programs and test functions.
- Three separate test programs: for_loop.t, while_loop.t, cond.t (one per control flow construct).
- cond.t uses a sign() function to exercise all three branches: if, elif, else.
- Added TestCondStatement as a bonus (covers if/elif/else, not in original issue criteria but clearly needed).

## Last checkpoint

COMPLETE. All criteria satisfied plus bonus TestCondStatement. One bug found and fixed: emitForUpdate only handled KIND_ASSIGNMENT_STMT but the parser emits for-loop updates as KIND_VAR_STMT. Fixed by adding KIND_VAR_STMT case and extracting a emitVarTarget() helper that handles KIND_VAR_ID_STMT nodes (which are not expressions and cannot go through emitExpr). Also fixed emitVarStmt reassignment to use the same helper.
