# `std::slices`

Source: [`slices.t`](../../base/std/slices/slices.t).
Import `std::slices` for recoverable element access on borrowed slices.

| API | Result | Contract |
| --- | --- | --- |
| `length<T>(values T[])` | `u64` | Element count; no allocation, including if `len` is shadowed |
| `get<T>(values T[], index u64)` | `(T, error)` | Shallow element copy; zero value and `OutOfBounds` if invalid |
| `set<T>(values T[], index u64, value T)` | `error` | Overwrite one element; `OutOfBounds` changes nothing |

These APIs neither allocate nor take ownership. The caller must keep the backing
storage alive. `set` does not clean up the overwritten element's resources.

Fragment inside a function returning `error`:

```rust
values i32[2] := .{10, 20}
slices::set<i32>(values[:], 1, 42) @try
value := slices::get<i32>(values[:], 1) @try
```

Ordinary `values[index]` instead panics when out of bounds. A bounds check does
not validate a pointer or extend a view's lifetime. See
[arrays and slice semantics](../reference/collections.md).
