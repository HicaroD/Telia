# `std::path`

Source: [`path.t`](../../base/std/path/path.t). Import `std::path`.

| API | Result | Contract |
| --- | --- | --- |
| `validate(path string)` | `error` | Nonempty, valid UTF-8, no embedded zero; no existence check |
| `canonicalize(path string)` | `(strings::Owned, error)` | Absolute canonical existing path, resolving links and parent components |
| `canonicalize_for_create(path string)` | `(strings::Owned, error)` | Existing canonical path, or canonical existing parent plus one missing final component |

Validation reports `InvalidPath` for empty input, `EmbeddedZero` for zero bytes,
or `InvalidUTF8`. Canonicalization can report native lookup/access/path failures
and allocation errors; failures return empty owners. Keep borrowed inputs live
for the call. Successful output is independent owned storage; release it with
`strings::free`.

Fragment inside a function returning `error`, with `std::strings` also imported:

```rust
resolved := path::canonicalize(".") @try
defer strings::free(&resolved)
```

Missing intermediate directories and dangling links fail. Canonicalization is
not a sandbox, file reservation, or concurrency guarantee. POSIX arbitrary-byte
filenames are not supported: inputs must be UTF-8. Windows uses wide APIs and
native separators. See [path and working-directory details](native-cli-apis.md#paths-and-current-directory).
`resolve` and `canonical_into` are implementation helpers.
