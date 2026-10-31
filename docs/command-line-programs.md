# Unicode-aware command-line programs

[`examples/text_cli.t`](../examples/text_cli.t) combines the published standard
library APIs without application-local C bindings:

```text
text_cli [INPUT|- [REPEAT [COUNTS]]]
```

- `INPUT` defaults to `-` (stdin); otherwise it is a UTF-8 file path.
- `REPEAT` is a strict signed decimal integer from 1 through 100. When omitted,
  `TELIA_TEXT_REPEAT` supplies its default, or `1` if the variable is missing.
  Present-but-empty defaults are invalid. An explicit argument bypasses lookup.
- `COUNTS` is exactly `true` or `false` (default `false`). True prints each trimmed
  line's code-point count instead of its text, using canonical unsigned decimal.
- Input is validated in full as UTF-8 before output. Split on LF, trim Unicode
  whitespace from each field, and print each result `REPEAT` times with an LF.
  Empty fields are preserved: empty input produces one empty line, and a trailing
  LF produces a final empty line. CRLF works because trimming removes the CR.
- Exit status is zero on success, one on handled errors. Diagnostics go to stderr.
  A failed diagnostic still returns one. Writes are checked and flushed, but an
  output failure can leave a partial result; there is no rollback.

## Run it

From a configured development checkout (use `telia` instead of `./telia` for an
installed compiler):

```sh
./telia build examples/text_cli.t -o bin/text_cli
printf '  世界 😀  ' | ./bin/text_cli
# 世界 😀
printf '  世界 😀  ' | ./bin/text_cli - +2 true
# 4
# 4
printf '  café  ' > input.txt
TELIA_TEXT_REPEAT=2 ./bin/text_cli input.txt
# café
# café
./bin/text_cli input.txt bad
# stderr: repeat must be an integer from 1 to 100; exit 1
```

In PowerShell, build to `bin/text_cli.exe`, set
`$env:TELIA_TEXT_REPEAT = '2'`, and run `./bin/text_cli.exe input.txt`.
Shell pipelines can add newlines or transcode text; use a UTF-8 file when exact
bytes matter. The tests supply subprocess bytes directly on every platform.

## Bytes, text, and ownership

[File I/O](files.md) and [standard-stream I/O](standard-io.md) read arbitrary bytes,
including embedded zero and malformed UTF-8. This application explicitly chooses
strict text processing at `utf8::validate`; embedded zero remains a valid code
point. Counts are not byte lengths or user-perceived grapheme counts. No automatic
normalization or case folding occurs. Whitespace uses the bundled Unicode
**17.0.0** `White_Space` property; see [provenance, license, and update procedure](unicode.md#whitespace-data).

The [argument snapshot and environment result](process.md) own their storage.
The input owner stays alive until processing returns. Split owns only its list
backing storage; its fields and trimmed strings borrow the input. The list is
freed before the input, never the fields individually. Formatted counts are
freed once per call to `print_line`, not by defers accumulated in a loop.
Empty local owners are zero-initialized, not constructed from borrowed literals.
Owner assignments are explicit transfers: the old shallow copy is never freed.

All acquired owners have ordinary deferred cleanup and all application errors
return normally: the example does not depend on `@fail` unwinding. File APIs close
their own handles; standard streams remain process-owned. Whole-input buffering,
split storage, and up to 100 repetitions make this a teaching example, not a
streaming or resource-quota framework.

See also [string views and utilities](strings-and-memory.md),
[checked conversions](string-conversions.md), and the
[expansion coverage audit](stdlib-cli-coverage.md).

## Focused verification

```sh
go test ./tests/integration -run '^TestTextCLI' -count=1 -timeout=5m
TELIA_TEST_SANITIZERS=address,undefined go test ./tests/integration -run '^TestTextCLIAllocationFailuresAndSanitizers$' -count=1 -timeout=5m
```

Tests build at `-O0` and `-O3`, launch controlled subprocesses with explicit
arguments/environment/stdin, and compare stdout, stderr, and exit status. Native
test-only hooks fail each allocation up to successful completion, abort on leaked
storage, and inject read/write/flush failures. They do not alter the example's
published API boundary. Non-Windows fixtures default to UndefinedBehaviorSanitizer;
Windows runs allocation/failure checks without unsupported MinGW sanitizers.
