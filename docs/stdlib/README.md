# Standard library

Import packages with `use "std::name"`. This index covers all application-facing
packages under [`base/std`](../../base/std). Nested packages are separate imports.
The linked contracts define results, errors, allocation, ownership, and platform
restrictions; do not infer those from a signature alone.

| Package | Purpose and API reference |
| --- | --- |
| `std::errors` | [Classified failure identities](errors.md) |
| `std::ffi::c` | [C-native scalar, C-string, and variadic types](ffi-c.md) |
| `std::fs` | [Whole files, metadata, bounded reads, and atomic saves](fs.md) |
| `std::io` | [Standard streams and printing](standard-io.md) |
| `std::libc` | [Low-level C library declarations](libc.md) |
| `std::lists` | [Growable generic lists](lists.md) |
| `std::maps` | [Heap hash maps and explicit iteration](maps.md) |
| `std::math` | [Native floating-point math](math.md) |
| `std::memory` | [Heap allocation, arenas, and transfer](memory.md) |
| `std::path` | [Path validation and canonicalization](path.md) |
| `std::process` | [Arguments/environment and package navigation](process.md), [children](native-cli-apis.md#subprocesses), [working directory](native-cli-apis.md#paths-and-current-directory) |
| `std::slices` | [Borrowed-slice helpers](slices.md) |
| `std::strconv` | [Checked text parsing and formatting](string-conversions.md) |
| `std::strings` | [Borrowed/owned text and C strings](strings-and-memory.md), [byte copying](native-cli-apis.md#bytes-and-strings) |
| `std::terminal` | [Handles, modes, readiness, and control sequences](terminal.md) |
| `std::testing` | [Assertions and named tests](testing.md) |
| `std::unicode` | [Code-point display width](native-cli-apis.md#unicode-display-cell-width) |
| `std::unicode::utf8` | [Strict validation, decoding, and iteration](unicode.md) |

For SDL and other vendor packages, see [native libraries](native-libraries.md).
Those packages are not part of `std` and are not exhaustively catalogued here.

## Shared contracts

- `string` and slices borrow storage. Owner handles are shallow-copyable; choose
  one cleanup owner. Read [ownership conventions](strings-and-memory.md#ownership-conventions-and-transfer)
  before retaining views or returning an owner with a defer.
- Errors are checked values, not exceptions. Use `@try` for cleanup-safe
  propagation; `@fail` skips pending defers. Compare classified `err.kind`, not
  messages. See [language errors](../reference/errors.md).
- Bytes are not automatically valid UTF-8. A `string` may contain embedded zeros
  or malformed UTF-8 unless an API explicitly validates it.
- Public lengths, offsets, sizes, and alignments use `u64`; native C scalar types
  follow the target instead. See [conversions](../reference/conversions.md).
- OS-facing support is documented for Linux amd64, macOS amd64/arm64, and Windows
  amd64 with MinGW-w64. Platform-specific limitations remain part of each API's
  contract; a CI matrix is not evidence that a particular commit passed it.

## Application APIs versus implementation helpers

The compiler does not provide a visibility boundary for every library helper.
Names being accessible in source does **not** make them application APIs.
Use the operations documented by these package pages. Native adapter functions
(`*_into`, `native_*`, `*_state`, and runtime link names), representation fields,
generated Unicode classifiers, list buffer allocation helpers, parsing internals,
and terminal encoding internals are implementation plumbing unless explicitly
documented otherwise. `strings::from_c`, `slices::length`, and other explicitly
listed low-level operations remain application APIs.

`std::internal::swiss` is the shared [table engine](../contributing/maps-engine.md),
not an application API. Map representation and iterator fields are also internal.

| Internal package | Contributor contract |
| --- | --- |
| `std::internal::swiss` | [Shared map/set table engine](../contributing/maps-engine.md) (not an application API) |

The sources linked from each page help audit behavior; they do not replace the
contracts or authorize calling undocumented helpers. Contributors adding a new
application API must add it to the relevant reference and tests.
