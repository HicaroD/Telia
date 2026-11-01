# `std::fs`

Import `std::fs`. File contents are arbitrary bytes; paths must be nonempty
UTF-8 with no embedded zero. All lengths and offsets below are `u64`.

| Task | APIs and authoritative contracts |
| --- | --- |
| Whole-file bytes | [`read_file(path string) (strings::Owned, error)`, `write_file(path string, data string) error`](files.md) |
| Identity, type, size, time, permissions | [`stat`, `lstat`, `same_file`, `Metadata`, `FileType`](native-cli-apis.md#metadata-and-errors) |
| File handles and bounded input | [`open_read`, `read`, `close`, `read_preview`, `File`](native-cli-apis.md#bounded-reads-and-file-handles) |
| Atomic replacement and primitives | [`write_atomic`, `create_temp`, `write`, `sync`, `preserve_permissions`, `rename`, `remove`, `sync_parent`](native-cli-apis.md#atomic-saving) |

Whole-file and preview results own their returned bytes; free them with
`strings::free`. Views borrow those bytes. `File` is a live owning handle, not a
copyable independent file. `close` clears it and can be repeated, but cleanup
errors require explicit handling if they matter. Reading and writing may have
observable partial effects even on failure.

This fragment assumes `std::strings` and an enclosing function returning `error`:

```rust
contents := fs::read_file("input.txt") @try
defer strings::free(&contents)
```

Check the atomic-save boolean as well as the error: a failure after replacement
does not mean the old file survived. On Windows, parent-directory sync is
unsupported, so `(true, Unsupported)` is a documented successful replacement
without that durability guarantee. Never blindly retry it.

Source files: [whole-file I/O](../../base/std/fs/fs.t),
[metadata](../../base/std/fs/metadata.t), [handles](../../base/std/fs/file.t),
[atomic saving](../../base/std/fs/atomic.t). Native adapters and `cleanup_temp`
are implementation helpers, not a substitute for a complete cleanup algorithm.
