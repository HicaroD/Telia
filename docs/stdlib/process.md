# Process arguments and environment

For child processes, binary pipes, and termination/reaping, see
[subprocesses](native-cli-apis.md#subprocesses). For `chdir(directory string) error`,
see [working-directory changes](native-cli-apis.md#paths-and-current-directory).

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

[`examples/process_args.t`](../../examples/process_args.t) prints each argument on
its own line and releases its snapshot on success and handled errors:

```sh
telia build examples/process_args.t
./process_args '' 'two words' '世界' '😀'
```

## Environment lookup

`lookup_env(name string)` returns `(strings::Owned, bool, error)`:

| Native value | Owner text | Found | Error |
| --- | --- | --- | --- |
| Missing | Empty | `false` | `nil` |
| Present but empty | Empty | `true` | `nil` |
| Present and nonempty | Independent UTF-8 copy | `true` | `nil` |
| Invalid text, allocation, or native lookup failure | Empty | `false` | Non-nil |

Release the owner with `strings::free`, including empty results. Its view remains
valid until cleanup and does not borrow native environment storage. Copies of
`strings::Owned` are shallow aliases, not independent owners.

Names must be nonempty UTF-8 without `=` or embedded zeros. Empty names and `=`
produce `InvalidEnvironmentName`; zeros produce `EmbeddedZero`. Malformed names
or values produce `InvalidUTF8`, never replacement characters. Allocation/size
failures use `AllocationFailed`/`SizeOverflow`; native Windows lookup failures
use `EnvironmentLookupFailed`. Partial storage is released before returning.

POSIX lookup uses native `getenv`, validates the value bytes as UTF-8, and is
case-sensitive. Windows uses `GetEnvironmentVariableW` with strict UTF-8/UTF-16
conversion, independent of the ANSI code page, and native case-insensitive name
comparison. No normalization or additional name folding is performed.
Environment mutation/enumeration is not provided. Native code must not mutate
the environment concurrently with lookup: POSIX native storage is borrowed
during copying, and a Windows size/read race can return `EnvironmentLookupFailed`
rather than a partial value. The existing borrowed `std::libc` binding is unchanged.

[`examples/process_env.t`](../../examples/process_env.t) prints an environment
greeting, using a default only when the variable is missing (not when empty):

```sh
telia build examples/process_env.t
TELIA_GREETING='Hello 世界 😀' ./process_env
TELIA_GREETING='' ./process_env
```

PowerShell: `$env:TELIA_GREETING = 'Hello 世界 😀'; .\process_env.exe`.

See [strings and memory](strings-and-memory.md) for borrowed views and
[Unicode](unicode.md) for explicit Unicode processing.
