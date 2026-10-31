# Strict UTF-8

Import `std::unicode::utf8`. These allocation-free operations implement Unicode
scalar-value UTF-8 (RFC 3629). They reject malformed input rather than replacing
it: invalid leading/continuation bytes, truncated sequences, overlong encodings,
surrogates (`U+D800`–`U+DFFF`), and values above `U+10FFFF`.
Embedded zero bytes represent `U+0000` and do not terminate text. Noncharacters
are valid scalar values. No external Unicode dependency or versioned character
property data is needed for UTF-8 encoding rules.

| API | Result | Contract |
| --- | --- | --- |
| `validate(text string)` | `error` | Validates the entire input; empty input succeeds |
| `decode(text string, offset u64)` | `(u32, u64, error)` | One code point and its byte width, starting at a byte offset |
| `count(text string)` | `(u64, error)` | Number of code points; validates the entire input |
| `iterator(text string)` | `(Iterator, error)` | Validates the entire input before returning a borrowed iterator |
| `next(iter *Iterator)` | `(u32, u64, bool, error)` | Code point, starting byte offset, found flag, error; advances on success |

`decode` validates only the sequence at the supplied offset, not unrelated bytes.
An offset at/beyond the byte length (including an empty input) or at a continuation
byte is an error. Its failure results are zero code point and zero width. `count`
returns zero on failure, even after a valid prefix. Error messages are static,
need no cleanup, and are not stable identifiers.

Iterator creation returns an empty iterator on failure, exposing no valid prefix.
Clean exhaustion returns `(0, 0, false, nil)`, including on repeated calls and on a
zero-initialized iterator. A decoding error returns `(0, 0, false, error)` without
advancing. `U+0000` is a value with `found = true`, not exhaustion.

An iterator borrows its source. Keep that storage alive and unchanged while using
it; no copy or ownership transfer occurs. Iterator copies share source storage but
advance independently. Treat its representation fields as implementation details;
use `iterator` and `next`, not direct mutation. No iterator cleanup is needed.

Offsets and widths are bytes, compatible with `strings::slice`. Existing
`strings::len` still counts bytes and byte-oriented string/FFI behavior is unchanged.
Code points are not grapheme clusters. Normalization, casing, grapheme operations,
and Unicode whitespace classification are outside this API.

The source-level `byte_into` and `read_byte` functions are implementation plumbing,
not application APIs. The private C `_telia_string_byte` helper bounds-checks the
index, returns a byte widened to `u32` (zero on error), and never exposes storage
or allocates. All UTF-8 decoding and iteration logic is ordinary Telia code.

## Example and verification

Build with `./telia build examples/utf8.t -o bin/utf8`, then run `./bin/utf8`
(see the CLI help for local build configuration).
The example counts and iterates multilingual text with supplementary-plane emoji.

```sh
go test ./tests/integration -run TestUTF8 -count=1
go test ./...
```

Fixtures run at `-O0` and `-O3` with subprocess timeouts. On supported non-Windows
hosts, malformed-byte tests and native sanitizer fixtures use UndefinedBehaviorSanitizer.
Set `TELIA_TEST_SANITIZERS=address,undefined` for the native sanitizer fixtures when
the host provides a working AddressSanitizer runtime. MinGW-w64 sanitizer tests
are skipped; ordinary debug/release behavior remains tested on Windows.
