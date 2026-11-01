# Checked standard-stream byte I/O

Import `std::io`:

| API | Result | Contract |
| --- | --- | --- |
| `read_stdin()` | `(strings::Owned, error)` | Read through EOF; caller frees the owner with `strings::free` |
| `write_stdout(data string)` | `error` | Write all borrowed bytes to stdout, then flush |
| `write_stderr(data string)` | `error` | Write all borrowed bytes to stderr, then flush |
| `println(text string)` | none | Existing unchecked output with a trailing newline; unchanged |

These checked operations preserve arbitrary bytes, including embedded zeros,
invalid UTF-8, CR/LF, and control-Z. They neither validate Unicode nor append a
newline. To interpret input as text, explicitly validate
`utf8::validate(strings::view(&owner))`. Stdout and stderr remain separate streams;
there is no ordering guarantee between them when an external consumer merges them.

## Ownership and failure

Stdin is read without seeking or assuming a size. Short reads continue until EOF;
immediate EOF succeeds with an empty owner. Views borrow the returned owner.
`strings::free` releases and clears it and is safe to repeat. Every failed read
releases partial storage and returns an empty owner. Input already consumed is
not restored, including on allocation failure.

Writes borrow input only for the call, handle short writes, and flush before
returning. Empty writes also flush. A flush is not a durable-storage guarantee.
Output already emitted cannot be rolled back; failure may leave partial output.
Unsent buffered output may be discarded on failure to prevent a later implicit
flush from terminating the process on a broken pipe.

Errors use `std::errors::Error` identities through `err.kind`:

- `AllocationFailed`, `SizeOverflow`: input buffer allocation or size growth failed.
- `StreamReadFailed`: native read failed, or made no progress without EOF.
- `StreamWriteFailed`: native write failed or made no progress; on POSIX, signal
  protection setup/restoration failure is also reported this way.
- `StreamFlushFailed`: flushing failed.
- `StreamModeFailed`: Windows binary-mode setup or restoration failed.

An earlier operation error takes precedence over a subsequent flush or mode-restore
error. No OS error numbers are exposed. Standard streams are process-owned and
are **never closed** by these APIs. Native EOF/error indicators are not cleared;
retries after EOF/errors are not a stream-reset facility.

## Native interaction and platforms

On POSIX, no text-mode conversion is performed. Checked writes temporarily block
SIGPIPE in the calling thread, consume only a newly pending SIGPIPE, and restore
the original signal mask. Existing pending SIGPIPE and the process-wide handler
are preserved. This protection does not change unchecked `println` or native C
output behavior.

On Windows (MinGW-w64), each operation temporarily selects `_O_BINARY` for its
standard-stream descriptor and restores the previous mode before returning.
Output buffered by earlier native calls is flushed **before** changing mode, so
those bytes retain their prior text semantics. Checked output is flushed before
restoring mode. If restoration itself fails, the descriptor may remain binary.
Read stdin before mixing in native text-mode input: translation already performed
by an earlier read cannot be undone. Binary mode preserves bytes; it does not
provide terminal editing, Unicode console conversion, or terminal control.

Sequential checked writes and `println` compose in call order on stdout:
`write_stdout("prefix:")` followed by `println("line")` emits `prefix:line\n`.
`println` keeps its existing signature and unchecked semantics. Native C output
after a successful checked call uses its original descriptor mode. Serialize
access to each stream when using native threads: these APIs do not promise an
atomic multi-call transaction or synchronized Windows mode changes.

Support targets are Linux amd64, macOS amd64/arm64, and Windows amd64 with
MinGW-w64. Focused native CI runs the byte and failure-path tests at `-O0` and
`-O3`; sanitizer checks apply where the toolchain supports them. CI configuration
does not itself establish a passing platform run.

## Runnable example

[stdio.t](../../examples/stdio.t) copies stdin through EOF to stdout and reports
failures on stderr. It buffers the entire input, so it is not a streaming copy.

```sh
telia build examples/stdio.t -o stdio-example
./stdio-example < input.bin > output.bin
```

On Windows, use `stdio-example.exe < input.bin > output.bin` from `cmd.exe` to
avoid shell text conversion. Line readers, public stream handles, and file APIs
are not part of this module; see [whole-file I/O](files.md) for files.
