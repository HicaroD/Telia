# Low-level terminal applications

Import `std::terminal` for byte I/O, terminal modes, dimensions, input waiting,
resize notifications, and terminal control sequences. CLI argument parsing,
key/mouse decoding, text layout, widgets, and PTYs are not part of this package.

The native backend targets Linux, macOS, and Windows. It ships with the standard
library and does not need an external terminal library.

## Handles and ownership

| Function | Result |
| --- | --- |
| `stdin()`, `stdout()`, `stderr()` | Borrowed standard stream |
| `borrow_native(value c::c_intptr_t)` | Borrowed Unix descriptor or Windows `HANDLE` |
| `open_controlling(input bool)` | Owned controlling-terminal handle, or error |
| `close(handle *Handle)` | Invalidates the wrapper; closes only owned handles |
| `is_terminal(handle *Handle)` | `(bool, error)`; files and pipes return false |

On Windows, `input` chooses `CONIN$` versus `CONOUT$`; Unix opens `/dev/tty`
read/write. Opening fails when the process has no controlling terminal.
Windows native handles are **not** C runtime file descriptors.

Treat `Handle` and `SavedMode` as noncopyable owners: Telia does not enforce this.
Do not mutate their bookkeeping fields. Initialize handles through constructors;
a zero-initialized handle is closed. Closing a wrapper twice is harmless.
Borrowed wrappers do not extend the native resource's lifetime. Keep handles open
until all saved modes have been restored. Coordinate mode changes and use one
reader per input stream; serialize output when interleaving would corrupt it.

## Byte I/O

- `read(handle, buffer u8[])` returns `(ReadResult, error)` with `count` and `eof`.
- `write(handle, bytes u8[])` returns `(count u64, error)` and may write partially.
- `write_all(handle, bytes u8[])` writes the whole slice or returns an error.
- `write_text(handle, text string)` writes the string's bytes without allocation,
  NUL termination, encoding validation, or newline translation.

Ordinary I/O works on files and pipes as well as terminals. Empty reads return
zero bytes without EOF. Reads may split UTF-8 characters or escape sequences.
Redirected bytes, including embedded zeros and invalid UTF-8, are preserved.
Windows console input uses the Unicode console API and is encoded as UTF-8;
unpaired UTF-16 surrogates become the replacement character U+FFFD. Bytes that
do not fit in a read buffer remain on that `Handle` for the next read, and
input waiting reports them as ready. Keep using the same wrapper while reading;
independently borrowed wrappers do not share this buffered state. Files and pipes
still use byte I/O without encoding conversion.
Unix broken-pipe writes return an error rather than raising SIGPIPE from that
write; existing signal handlers and previously pending signals are preserved.
Native failures use `Error::SystemFailure`; OS-specific error details are not
currently exposed. Nonblocking descriptors can report a system error when they
would block; this package does not change their blocking flags.

These operations bypass C stdio buffering. Avoid mixing them with buffered libc
output on the same stream unless you explicitly flush the libc stream.

## Modes and restoration

```rust
input := terminal::stdin()
saved := terminal::apply_mode(&input, terminal::raw_mode()) @fail
defer terminal::restore(&saved)
```

`get_mode(handle)` returns configurable `echo`, `canonical` (line buffering), and
`signals` settings. `raw_mode()` provides `Mode{raw: true}` with those settings
disabled: Ctrl+C arrives as byte 3 rather than terminating the application.
`apply_mode` applies the raw preset first when requested, then the three explicit
settings, and returns a snapshot of the original native state. On Unix, raw mode
also disables output postprocessing, so a newline does not imply carriage return.

Restore snapshots in reverse application order. `restore` consumes and clears a
successful snapshot and is then safe to repeat. A failed restore retains the
snapshot for retry. Native user settings are restored, although kernels may
update bookkeeping flags such as macOS `PENDIN`. Cleanup is explicit, not a
process-wide exit or signal hook: forced termination cannot guarantee restoration.

Telia's `@fail` terminates without running defers. After acquiring terminal
state, use `@catch` and return errors so cleanup runs; put the fatal `@fail`
boundary outside the owning function, as the interactive examples do. Deferred
cleanup errors are not automatically propagated: explicitly check `restore` on
normal shutdown when you need to report restoration failures.

Call `prepare_output(handle)` before terminal controls. It returns another
restorable snapshot. On Windows this enables virtual-terminal output and UTF-8;
on Unix it validates the terminal and otherwise does nothing. Input mode changes
enable virtual-terminal input and UTF-8 on Windows. Console code pages are shared
console state: applications must coordinate these changes too. Unsupported
Windows console modes return `Error::Unsupported`.

## Waiting and resizing

```rust
event := terminal::wait(&input, terminal::WaitOptions{
    input: true,
    resize: true,
    timeout_ms: 100,
}) @fail
```

`WaitResult` contains `input`, `resized`, and `timed_out`. Input and resize may be
reported together; timeout is separate. Readable input can mean EOF: call `read`
to distinguish it. Choose `timeout_ms: -1` to block, `0` to poll, or a positive
millisecond deadline. At least one interest is required. A timeout does not
consume input, and a later read remains the application's responsibility.

`size(handle)` returns `Size{columns, rows}`. Resize waiting establishes its
baseline on the first call for that handle, then reports changed dimensions on
subsequent calls or during a blocked wait. It uses size polling at intervals of
up to 16 ms rather than installing process-wide signal handlers. Notifications
coalesce changes; a resize back to the previous dimensions between samples can
be missed. Resize interest on a nonterminal returns `NotTerminal`.

Windows console input waits require noncanonical mode; line-buffered console
waiting returns `Unsupported`. Windows pipe waits support anonymous/byte pipes;
message-mode and overlapped-handle I/O are outside this version's contract.

## Capabilities and output controls

`capabilities(handle)` returns `Capabilities`. Each feature is `Supported`,
`Unsupported`, or `Unknown`. Detection uses attachment and environment hints,
not probes that might consume the application's input. Files, pipes, and
`TERM=dumb` report unsupported controls. Recognized ANSI terminal names,
`256color`, and `COLORTERM=truecolor`/`24bit` supply positive hints; extensions
such as mouse tracking otherwise remain unknown. This is not terminfo or proof
of emulator behavior. Applications can override individual fields explicitly.

`send_control(handle, control, &caps)` rejects nonterminals and known unsupported
features; unknown features are attempted. The package does not apply `NO_COLOR`
or CLI presentation policy. For redirected output, decide your own fallback.

`encode_control(control, buffer u8[])` returns `(count, error)` without requiring
a terminal. It writes nothing on invalid arguments or insufficient capacity.
A 64-byte buffer accommodates every supported control. Append encoded bytes and
text to your own frame buffer, then call `write_all` for fewer writes.

| `ControlKind` | Fields |
| --- | --- |
| `CursorPosition` | `first`: row, `second`: column, both **1-based**, 1..65535 |
| `CursorUp`, `CursorDown`, `CursorForward`, `CursorBack` | `first`: distance, 1..65535 |
| `CursorVisible`, `Bold` | `enabled` |
| `ClearScreen`, `ClearLine`, `ResetStyle` | No arguments; clearing does not move the cursor |
| `Foreground`, `Background` | `first`: indexed color, 0..255 |
| `ForegroundRGB`, `BackgroundRGB` | `first`, `second`, `third`: RGB, each 0..255 |
| `AlternateScreen`, `BracketedPaste`, `MouseTracking` | `enabled` |

Mouse tracking toggles normal button tracking and SGR mouse encoding (1000 and
1006). Paste/mouse replies remain raw input bytes. Explicitly undo emulator
settings (alternate screen, cursor visibility, mouse, paste, styles) before
restoring native modes: native snapshots cannot capture emulator state.

## Examples and tests

- [`terminal_echo.t`](../examples/terminal_echo.t): pipe-friendly byte copying.
- [`terminal_raw.t`](../examples/terminal_raw.t): raw input inspection.
- [`terminal_screen.t`](../examples/terminal_screen.t): full-screen lifecycle and resizing.

Run `go test ./tests/integration -run '^TestTerminal'`. POSIX terminal tests use
Python 3's real PTY support, without touching the developer's terminal. The
sanitizer checks default to UBSan; `TELIA_TEST_SANITIZERS` selects alternatives
where their runtime is available. Windows console behavior must also be verified
on Windows; POSIX PTY tests are skipped there.
