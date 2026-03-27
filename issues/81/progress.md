# Issue #81: Error handling: nil comparison for manual error handling

## Acceptance criteria

- [x] `err != nil` type-checks in sema (no diagnostic error)
- [x] `err == nil` type-checks in sema
- [x] `err != nil` emits `err.msg != NULL` in C
- [x] `err == nil` emits `err.msg == NULL` in C
- [x] Integration test: `error_return.t` passes (unpack + nil check + field access)
- [x] Sema test: nil comparison with error type passes (including `==`, reversed order)

## Decisions

- Core implementation was carried over from #80 (error type, constructor, field access)
- `isErrorNilComparison()` in sema handles both orderings (`err op nil` and `nil op err`)
- Codegen emits `err.msg != NULL` / `err.msg == NULL` because C cannot compare structs by value
- Added 3 additional sema test cases: `err == nil`, `nil != err`, `nil == err`
- Existing integration test `error_return.t` covers the full unpack + nil check + field access flow

## Last checkpoint

All acceptance criteria complete. Sema tests cover all 4 operator/ordering combinations.
Integration test `error_return.t` passes end-to-end.
