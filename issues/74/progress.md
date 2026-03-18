# Issue #74 — C backend: defer statement

## Acceptance criteria

- [x] Deferred statements are emitted in reverse order before every `return` in the block
- [x] `defer` in nested scopes (e.g. inside an `if`) is handled correctly
- [x] Functions with no `defer` statements are unaffected
- [x] New `TestDefer` integration test passes and verifies LIFO execution order
- [x] All previously passing tests continue to pass

## Decisions

### Core codegen was already correct
`emitBlock` already flushed `DeferStack` in LIFO order before every `KIND_RETURN_STMT`. The only work needed was the bug fix and documentation.

### Bug fix: void functions / fall-through blocks silently dropped deferred statements
`emitBlock` only flushed on `KIND_RETURN_STMT`. A void function with no explicit `return` (or a nested `if`-block that falls through without a `return`) never triggered the flush — deferred statements were silently lost.

**Fix:** Extract flush logic into `flushDeferStack()`. Add a second flush call at the end of `emitBlock`, guarded by `!block.FoundReturn` (set by the parser when a `return` statement was parsed in this block). This avoids double-emission for blocks that already ended with an explicit return, with zero new fields and zero AST mutation.

### Design documentation
Telia's `defer` is **block-scoped**, not function-scoped (unlike Go). A deferred statement runs before the return of the block it appears in. This is documented on `BlockStmt.DeferStack` (`internal/ast/stmt.go`) and on `emitBlock` (`internal/codegen/c/stmts.go`).

### Nested scope test
`nested(1)` exercises the fall-through fix: `defer libc::printf("inner\n")` is inside an `if x > 0` block that has no `return`. With the fix, the `if`-block's DeferStack is flushed at fall-through, printing `inner` before control returns to the outer block. The outer `defer libc::printf("outer\n")` is then flushed before `return 0`, printing `outer`.

## Last checkpoint

All 5 acceptance criteria satisfied. Full suite green (13/13 integration tests). Ready to close.
