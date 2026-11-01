# Whole-file byte I/O

Import `std::fs`:

| API | Result | Contract |
| --- | --- | --- |
| `read_file(path string)` | `(strings::Owned, error)` | Read through EOF; caller frees the returned owner with `strings::free` |
| `write_file(path string, data string)` | `error` | Create or truncate, then write all borrowed bytes and close |

Both operations are binary-safe. Embedded zeros, invalid UTF-8, and newline bytes
are preserved without platform text-mode translation. **Reading bytes does not
validate Unicode text.** Use `utf8::validate(strings::view(&owner))` explicitly
before using the bytes as UTF-8. See [the runnable example](../../examples/files.t).

Paths must be nonempty, valid UTF-8, and contain no embedded zero bytes. Validation
happens before opening a file. POSIX uses these UTF-8 bytes directly; Windows uses
strict UTF-8-to-native-wide conversion, never the active ANSI code page. Arbitrary
non-UTF-8 POSIX paths are not supported. Relative paths use the process working
directory. Normal OS access rules and symlink resolution apply.

## Ownership and failures

Reads grow a checked buffer through EOF, without assuming a fixed file size or
requiring a seekable file. The returned owner independently owns its bytes; views
borrow that owner. Release it with `strings::free`, which clears it and is safe
to repeat. On **any** read failure, including close failure, partial storage is
released and the returned owner is empty. Empty files succeed with an empty owner.
All internally opened handles are closed on success and failure.

Errors use `std::errors::Error` identities, available through `err.kind`:

- `InvalidPath`: empty path.
- `EmbeddedZero`: zero byte in path (checked before UTF-8 validation).
- `InvalidUTF8`: malformed path encoding.
- `SizeOverflow`, `AllocationFailed`: path conversion or read buffer failure.
- `FileOpenFailed`, `FileReadFailed`, `FileWriteFailed`, `FileCloseFailed`:
  the corresponding native operation failed. A zero-progress operation that is
  not EOF also fails rather than looping forever.
- `NotFound`, `PermissionDenied`: missing/non-directory path components and
  access-denied open failures are distinct from other open failures.

Writes borrow `data` only for the call. Short reads/writes are handled internally.
If closing also fails after an earlier operation failed, the earlier error wins.
Native error numbers and detailed OS messages are not exposed by this initial API.

**Writes are not atomic:** opening truncates an existing destination, including
when `data` is empty. Failure may leave an empty or partially written file. No
rollback, atomic replacement, or durable-storage synchronization is promised.
Reads are not snapshots of concurrently modified files.

Public handles, append, directory operations, and terminal streams are outside
this whole-file API. For public handles, metadata, bounded previews, paths, and
atomic replacement, see [native CLI APIs](native-cli-apis.md).
Support targets are Linux amd64, macOS amd64/arm64, and Windows amd64
with MinGW-w64. Focused native CI covers these targets; local test success does
not establish a passing run on other platforms.

Run the example in a scratch directory (it creates/truncates `telia-file-example.bin`):

```sh
telia build /path/to/Telia/examples/files.t -o files-example
./files-example
```

On Windows, run `.\files-example.exe` after building.
