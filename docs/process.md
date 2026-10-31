# Process arguments

Import `std::process`. No change to a Telia program's `main` signature is needed.

| API | Result | Contract |
| --- | --- | --- |
| `args()` | `(Args, error)` | Independently owned, strict UTF-8 snapshot, excluding the executable name |
| `count(snapshot *Args)` | `u64` | Argument count; zero for an empty/freed handle |
| `get(snapshot *Args, index u64)` | `(string, error)` | Zero-based, checked access; string borrows the snapshot |
| `free(snapshot *Args)` | no result | Releases all snapshot storage, clears the handle; repeated calls are harmless |

Order, empty arguments, spaces, quotes, and Unicode are preserved. No Unicode
normalization, case folding, or option parsing is performed. Native arguments
cannot contain embedded zero characters. The executable name is neither returned
nor validated. A process launched without arguments yields an empty successful
snapshot. Zero-initialized `Args{}` is also an empty handle.

Every returned argument view remains valid until its snapshot is freed. Free the
handle, not individual strings. `Args` copies are shallow aliases, **not** new
owners: do not free both copies or use an alias/view after cleanup. Call `args()`
again for a separate owner. Keep the `state` field opaque and unchanged.

All-or-nothing creation returns an empty handle and an error on malformed text,
allocation failure, or allocation-size overflow. Partial storage is released.
Kinds are `errors::Error::InvalidUTF8`, `AllocationFailed`, and `SizeOverflow`.
Invalid access (including after cleanup) returns `""` and `OutOfBounds`.

## Native encoding and support

- Linux amd64 and macOS amd64/arm64: native argument bytes are copied and
  validated as UTF-8. Invalid bytes are rejected, never replaced.
- Windows amd64/MinGW-w64: a native `wmain` entry point captures UTF-16 arguments;
  strict `WideCharToMultiByte(CP_UTF8, WC_ERR_INVALID_CHARS, ...)` conversion
  rejects unpaired surrogates. It does not depend on the active ANSI code page.
  The C backend links with `-municode`; Telia source `main` remains unchanged.

The runtime borrows startup storage from the OS/CRT until process exit; snapshots
own their copies. Existing return values and deferred cleanup remain in Telia's
main function. The generated native test runner captures arguments too.

## Example

[`examples/process_args.t`](../examples/process_args.t) prints each argument on
its own line and releases its snapshot on success and handled errors:

```sh
telia build examples/process_args.t
./process_args '' 'two words' '世界' '😀'
```

See [strings and memory](strings-and-memory.md) for borrowed views and
[Unicode](unicode.md) for explicit Unicode processing. Environment lookup is not
part of this API yet.
