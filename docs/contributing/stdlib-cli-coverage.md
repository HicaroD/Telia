# Standard-library CLI expansion coverage audit

Scope: [spec #97](https://github.com/HicaroD/Telia/issues/97), delivered by
#98–#105. This maps every parent acceptance criterion to repository evidence;
test presence is not a claim that every platform has passed.

All named tests are in `tests/integration/`. Public contracts are documented in
[Unicode](../stdlib/unicode.md), [strings and memory](../stdlib/strings-and-memory.md),
[conversions](../stdlib/string-conversions.md), [process](../stdlib/process.md), [files](../stdlib/files.md), and
[standard I/O](../stdlib/standard-io.md). The [CLI guide](../cookbook/command-line-programs.md) is the
coherent entry point and includes runnable commands and lifetime guidance.

| Parent acceptance criterion | Implementation / documentation evidence | Automated test evidence |
| --- | --- | --- |
| Resolved public contracts documented | Guides above cover signatures, defaults, ownership, encodings, and limitations | API-level fixtures in the suites below |
| Validate, decode, iterate, count, trim without native Unicode dependency | `base/std/unicode/utf8/utf8.t`, `base/std/strings/trim.t`; bundled data | `TestUTF8Validation`, `TestUTF8DecodingAndCounting`, `TestUTF8BorrowedIteration`, `TestStringTrimUnicodeWhitespace` |
| Reject malformed UTF-8, overlong sequences, surrogates, out-of-range points | Strict scalar-value decoder | `TestUTF8StrictValidation`, `TestUTF8ScalarBoundaries`, `TestStringTrimRejectsMalformedInterior` |
| Exact byte search/prefix/suffix/split, empty and invalid-byte behavior | `base/std/strings/search.t`, `split.t`; string guide | `TestStringFindExactBytes`, `TestStringExactPredicates`, `TestStringSplitPreservesFields` |
| Borrowed trim/split lifetimes and list-only cleanup | String and collection ownership contracts | `TestStringUtilitiesOwnershipAndAllocationFailures`, `TestStringUtilitiesWithNativeSanitizers` |
| Strict i64/u64/bool conversion, limits, overflow, round trips | `base/std/strconv/strconv.t`; conversion guide | `TestStrconvIntegerBoundaries`, `TestStrconvUnsignedParsing`, `TestStrconvBooleanParsing`, `TestStrconvUnsignedFormatting`, `TestStrconvSignedAndBooleanFormatting` |
| Preserve existing byte strings, embedded zeros, memory/collections | Existing representation and cleanup APIs unchanged | `TestStringOutputPreservesBytes`, `TestStringEqualityComparesBytes`, `TestStringSlicesUseCheckedByteRanges`, `TestCollectionsWithNativeSanitizers` |
| Process arguments/environment ownership, encoding, missing vs empty | `base/std/process/`; process guide; native-wide Windows boundary | `TestProcessArgumentsPreserveValuesAndOrder`, `TestProcessArgumentBoundsAndCleanup`, `TestProcessEnvironmentMissingEmptyAndPresent`, `TestProcessEnvironmentUnicodeNamesAndNativeComparison` |
| Files/stdin/stdout/stderr errors, EOF, partial I/O, cleanup | `base/std/fs/`, `base/std/io/`; file/stream guides | `TestFSBinaryRoundTripAndTruncate`, `TestFSShortIOFailuresAndCleanup`, `TestIOStdinThroughEOF`, `TestIOEmptyEOFAndEmptyWriteFlush`, `TestIOSeparateStreamsAndPrintln`, `TestIOFailuresAndCleanup`, `TestIOStderrFailures`, `TestIOBrokenPipeIsChecked` |
| ASCII, multilingual, emoji, whitespace, malformed/empty/zero/split/allocation/cleanup coverage | UTF-8, utility, ownership, and CLI fixtures | Suites above; `TestStrconvOwnershipAndAllocationFailures`, `TestProcessArgumentAllocationFailureCleanup`, `TestProcessEnvironmentOwnershipAndAllocationFailures`, `TestTextCLIAllocationFailuresAndSanitizers` |
| OS-facing integration and supported-platform behavior | Native matrix: Linux amd64, macOS amd64/arm64, Windows amd64/MinGW-w64 | `TestProcessRejectsInvalidNativeUTF8`, `TestProcessWindowsRejectsInvalidWideArguments`, `TestProcessEnvironmentRejectsInvalidText`, `TestProcessEnvironmentWindowsFailures`, `TestFSRejectsInvalidPathsAndReportsOpenFailures`, `TestIONativeTextModeRestored`; CLI Unicode file path and OS failure tests |
| Debug/release, sanitizer/allocation checks, full Go suite | Explicit `-O0`/`-O3` loops; native sanitizer hooks; CI `make test` | Focused suites in `native.yml`; full suite in `ci.yml`. Results/gaps below |
| Unicode provenance/version, API docs, ownership, limitations, examples | Unicode 17.0.0 data, Unicode License V3, upstream URL and update steps in Unicode guide; individual runnable examples | `TestExamplesUsePublishedStandardLibraryAPIs` plus individual process/file/I/O example tests |
| Runnable combined CLI without local C bindings | `examples/text_cli.t`; CLI guide | `TestTextCLITrimsUnicodeLinesFromStdin`, `TestTextCLIFileInputEnvironmentDefaultAndOverride`, `TestTextCLICountsCodePointsNotBytesOrGraphemes`, `TestTextCLIRejectsInvalidOptions`, `TestTextCLIValidatesAllInputBeforeOutput`, `TestTextCLIEmptyInputAndOSFailures` |

## Verification record and remaining gaps

- Local macOS arm64: the focused `TestTextCLI` suite passed at `-O0` and `-O3`,
  including deterministic allocation failure through first success on stdin,
  environment, formatting, and file paths, with tracked leak detection and UBSan.
- AddressSanitizer plus UBSan was attempted locally. Instrumented subprocesses
  timed out before producing output, including the independently rerun existing
  `TestProcessEnvironmentOwnershipAndAllocationFailures`. This is an unresolved
  local ASan execution gap, not a passing sanitizer result.
- `.github/workflows/native.yml` runs the new `TestTextCLI` suite on all four native
  targets, with `address,undefined` requested on non-Windows. Existing UTF-8,
  string, conversion, process, file, and stream suites remain wired. Windows uses
  failure injection and allocation tracking, not unsupported MinGW sanitizers.
- `.github/workflows/ci.yml` runs `make test` (`go test ./...`). The full suite is
  intentionally not run locally for this task, per stakeholder instruction.
- The latest observed successful CI/native runs were for earlier commit
  `70768d74b66615c01824d5560578846637c1a661`:
  [CI](https://github.com/HicaroD/Telia/actions/runs/41224190563) and
  [native](https://github.com/HicaroD/Telia/actions/runs/41224190553).
  They do **not** verify these changes or the intervening local commits.
  New full-suite and native-platform results are deferred to post-push CI;
  cross-platform success remains unverified until those runs complete.

No new library capability, language change, external dependency, or unrelated
vendor change is introduced by this final slice. Issue closure records delivery
with the explicitly deferred verification above, not a claim that pending checks
have passed.
