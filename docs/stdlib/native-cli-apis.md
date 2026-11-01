# Files, child processes, and terminal cell widths

These public APIs support binary previews, atomic saves, search tools such as
`rg`, and clipboard programs without application-local C bindings. Targets:
Linux amd64, macOS amd64/arm64, Windows amd64/MinGW-w64.

## Bytes and strings

| API | Result | Contract |
| --- | --- | --- |
| `strings::byte_at(text, index)` | `(u8, error)` | Checked byte access; zero and `OutOfBounds` on failure |
| `strings::to_bytes(text)` | `(lists::List<u8>, error)` | Independent mutable copy; release with `lists::free<u8>` |
| `strings::from_bytes(data u8[])` | `(strings::Owned, error)` | Independent copy; release with `strings::free` |
| `slices::length<T>(data T[])` | `u64` | Slice length, including where `len` is shadowed |

Empty input succeeds. Embedded zeros and invalid UTF-8 are preserved. Changing or
freeing the source does not affect the copy. Failed conversions release partial
allocations and return empty owners. No mutable string-storage alias is exposed.

## Metadata and errors

`fs::stat(path)` and `fs::lstat(path)` return `(fs::Metadata, error)`. `stat`
follows symlinks; `lstat` describes the link itself. Metadata fields:

- `kind`: `fs::FileType::{Regular, Directory, Symlink, Other}`. Windows `lstat`
  classifies reparse points as `Symlink`.
- `volume`, `identity`: POSIX device/inode or Windows volume serial/file index.
  `fs::same_file(left, right)` compares both; use only successful metadata.
  Hard links compare equal. IDs can be reused after deletion, not persistent IDs.
- `size`: byte size, meaningful for regular files.
- `modified_seconds`: signed Unix seconds; `modified_nanoseconds`: fractional
  nanoseconds (`0..999999999`). Actual precision depends on the filesystem;
  Windows FILETIME has 100 ns resolution. No floating-point conversion.
- `permissions`: POSIX mode bits (`07777`), or Windows file attributes.

`NotFound` and `PermissionDenied` distinguish missing paths and access failures
from operation failures such as `MetadataFailed`. Existing `read_file` and
`write_file` now report these identities instead of collapsing all open failures
into `FileOpenFailed`. POSIX non-directory components also yield `NotFound`.
Match `err.kind`, not messages. Runtime errors are checked `error` values, not
compiler diagnostics. Paths remain nonempty UTF-8 without embedded zeros;
arbitrary non-UTF-8 POSIX filenames are unsupported.

## Paths and current directory

- `path::canonicalize(path)` → `(strings::Owned, error)`: existing absolute path,
  resolving symlinks and parent components.
- `path::canonicalize_for_create(path)` → `(strings::Owned, error)`: the same for
  existing paths; otherwise resolve an **existing directory parent** and append
  the missing final filename. Missing intermediate directories fail. Dangling
  symlinks/reparse points fail rather than becoming new regular files.
- `process::chdir(path)` → `error`: change the process-global working directory.

Free path owners with `strings::free`; failures return empty owners.
Canonicalization is not a sandbox, reservation, or concurrency guarantee.
Coordinate `chdir` with other threads and relative-path users. Windows uses wide
APIs and native separators; canonical drive paths omit the extended `\\?\`
prefix, while UNC paths retain their leading `\\`.

## Bounded reads and file handles

| API | Result |
| --- | --- |
| `fs::open_read(path)` | `(fs::File, error)` |
| `fs::read(file *File, buffer u8[])` | `(count u64, eof bool, error)` |
| `fs::close(file *File)` | `error` |
| `fs::read_preview(path, limit u64)` | `(strings::Owned, error)` |

Inspect only the first `count` bytes. Reads may be short. EOF can accompany final
bytes or arrive on a later read; empty-buffer reads do not probe EOF. File reads
can block and are not snapshots. A live `File` owns a handle: do not copy it.
`close` clears the handle and is repeatable; later I/O returns `InvalidState`.

Previews allocate at most `limit` payload bytes and read at most `limit` logical
bytes; native stream buffering may read ahead by a fixed amount. There is no
lookahead or claim that a limit-sized preview is the whole file. Zero limit still
validates/opens the path. All handles close on success/failure; read, close, and
allocation failures return an empty owner.

## Atomic saving

`fs::write_atomic(destination, data)` → `(replaced bool, error)`:

1. Canonicalize the destination, following existing symlinks. Existing targets
   must be regular files; a missing final filename is allowed.
2. Exclusively create a same-directory temporary file.
3. Write all bytes, preserve existing POSIX permission bits, flush/synchronize,
   and close. Windows `ReplaceFileW` preserves the target DACL during replacement.
4. Replace by rename, then synchronize the parent directory where supported.

`replaced` becomes true **only after successful replacement**. `(true,
FileSyncFailed)` means new content is visible but directory durability was not
established. Windows directory synchronization returns `Unsupported`, so a
successful Windows save returns `(true, Unsupported)`; **do not retry blindly**.
Atomic visibility differs from crash durability and filesystem/hardware guarantees.
No macOS `F_FULLFSYNC` or device-power-loss guarantee is made.

Failures before rename preserve the original destination and attempt temporary
cleanup, which does not allocate. The first error wins over secondary close/removal
failures. External interference or OS cleanup errors can still leave temporary
files. Successful rename never triggers unlink of the old temporary pathname.
A failed native Windows replacement can require inspecting native filesystem state.

Replacement changes file identity; other hard-link aliases retain old content.
Ownership, POSIX ACLs, extended attributes, and conflict detection are not
preserved/promised. New POSIX files start private (at most `0600`, subject to umask); new Windows files use
native security defaults. Lower-level primitives are public too:

| API | Result / behavior |
| --- | --- |
| `fs::create_temp(destination)` | `(File, strings::Owned, error)`; exclusive temporary beside the supplied pathname; caller closes/removes/frees it |
| `fs::write(file, data string)` | `error`; write all bytes to an open temporary file |
| `fs::preserve_permissions(source, destination *File)` | `error`; copy POSIX mode bits; standalone Windows copying is `Unsupported` |
| `fs::sync(file)` | `error`; flush user-space buffers and synchronize file |
| `fs::rename(source, destination)` | `error`; native replacement; use same-directory files for atomic saves |
| `fs::remove(path)` | `error`; remove file/link, not directory; missing is `NotFound` |
| `fs::sync_parent(path)` | `error`; synchronize containing directory; Windows is `Unsupported` |

Custom algorithms arrange their own cleanup. `create_temp` uses the supplied
pathname's directory; unlike `write_atomic`, it does not resolve destination links.

## Subprocesses

`process::spawn(executable string, arguments string[], options SpawnOptions)`
returns `(Child, error)`. Arguments exclude argv[0]; the executable becomes argv[0].
No shell evaluation, splitting, or expansion occurs. Empty arguments, whitespace,
quotes, backslashes, and Unicode are preserved; zeros and malformed UTF-8 fail
before spawning. Search follows native rules (`posix_spawnp`/`CreateProcessW`).
Windows uses CRT-compatible quoting and native command-line limits; programs with
custom argument parsers may differ. Environment/current directory are inherited.

`SpawnOptions` has `pipe_stdin`, `pipe_stdout`, and `pipe_stderr`. False inherits
the stream; true creates a binary pipe. POSIX internal pipe ends are close-on-exec;
Windows restricts inherited handles. `rg`/clipboard tools are external dependencies.

| API | Result / behavior |
| --- | --- |
| `read_stdout(child, buffer u8[])`, `read_stderr(child, buffer)` | `(count u64, eof bool, would_block bool, error)`; nonblocking bounded read |
| `write_stdin(child, data string)` | `(count u64, would_block bool, error)`; nonblocking partial write |
| `close_stdin(child)` | `error`; signal EOF; repeatable |
| `poll(child)` | `(ExitStatus, done bool, error)`; nonblocking status/reap |
| `wait(child)` | `(ExitStatus, error)`; blocking wait/reap; cached on repetition |
| `terminate(child)` | `error`; POSIX SIGTERM; Windows forced termination with code 1 |
| `kill(child)` | `error`; POSIX SIGKILL; Windows forced termination with code 1 |
| `close_child(child)` | `error`; close pipes, force-terminate if running, reap, release; repeatable |

Empty buffers/data are no-ops on valid pipes. Output EOF is repeatable and distinct
from would-block, and does not establish process exit. No-pipe/closed-child access
returns `InvalidState`. Broken stdin returns `StreamWriteFailed`, not fatal SIGPIPE.
Do not copy a live owning `Child`. If reaping fails, `close_child` retains ownership
for retry; explicitly check cleanup when errors matter.
File/child handles are not internally synchronized; coordinate callers.

`ExitStatus.code` is unsigned; `signal` is a POSIX terminating signal or zero.
Nonzero exit codes are **not** API errors (`rg` uses 1 for no matches). Lifecycle
failures use `ProcessSpawnFailed`, `ProcessWaitFailed`, `ProcessTerminateFailed`.
Termination alone does not reap. Do not install a competing child reaper.

For `rg`, use arguments such as `.{"--line-number", "--", pattern, path}` and
capture both output streams. For clipboard tools (`pbcopy`, `wl-copy`, `xclip`,
etc.), advance by write counts, then close stdin. Pump stdin/stdout/stderr together
when pipes can fill. **Drain outputs before blocking in `wait`** to avoid deadlock.
Retry would-block from an application's event/timer loop, not a tight spin.
Readiness integration, deadlines, process groups, and descendant-tree termination
are not supplied. `close_child` manages only the direct child and can block while
native termination/reaping completes.

## Unicode display-cell width

`unicode::width(point u32)` → `(i32, error)`, using bundled Unicode **17.0.0** data,
independent of locale and without runtime/build Unicode dependencies:

- `0`: NUL; combining marks (`Mn`, `Mc`, `Me`); formats (`Cf`); trailing Hangul
  Jamo; emoji skin-tone modifiers.
- `1`: ordinary printable scalars, including East Asian ambiguous characters.
- `2`: printable East Asian `Wide`/`Fullwidth` characters.
- `-1`: other controls, line/paragraph separators, unassigned scalars, noncharacters.
  This is a valid classification, not an error.
- Surrogates/values above `U+10FFFF`: zero and `InvalidCodePoint`.

Nonprintable classification precedes zero/wide properties. Tabs are nonprintable,
not fixed-width. This deterministic code-point policy does not guarantee every
terminal/font: graphemes, ZWJ emoji, variation selectors, and tabs need higher-level
layout logic.

`base/std/unicode/width_data.t` records upstream URLs/SHA-256 hashes for Unicode
`UnicodeData.txt` and `EastAsianWidth.txt`, under Unicode License V3 in
`base/std/unicode/data/LICENSE-UNICODE.txt`. Regenerate with
`python3 scripts/generate-unicode-width.py` (`--check` verifies reproducibility).
Updating Unicode changes public behavior: review diffs and independent boundary
expectations, then run width/UTF-8 tests. Only regeneration needs network access.

## Verification

```sh
go test ./tests/integration -run '^TestEditor' -count=1
go test ./...
```

Tests use public APIs at `-O0`/`-O3`, real files/children, binary data, sparse-file
previews, allocation-failure cleanup, pre/post-replacement failures, backpressure,
EOF/would-block, exit status, and termination. Native CI covers supported targets;
platform skips are explicit. Local macOS results do not verify Windows/Linux.
Non-Windows fixtures default to UBSan; `TELIA_TEST_SANITIZERS=address,undefined`
requests ASan too. MinGW runs ordinary/failure-injection checks without sanitizers.

Earlier implementation-time local results and sanitizer limitations are preserved
in [historical verification notes](../contributing/verification-records.md).
Use current runs, not those notes, to assess the checkout being changed.
