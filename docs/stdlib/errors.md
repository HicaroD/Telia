# `std::errors`

[`errors.t`](../../base/std/errors/errors.t) declares the payload-free
`enum Error: error` used by standard-library operations. The package allocates
nothing and has no cleanup function.

Construct a failure with `error(errors::Error::Variant, "context")`; inspect an
existing failure with `err.kind == errors::Error::Variant`. These are fragments
that assume `use "std::errors"`. An enum value does not implicitly become an
`error`. Messages are human-readable context, not stable identifiers.

| Identity | Meaning |
| --- | --- |
| `OutOfBounds`, `InvalidRange` | Invalid index or range |
| `SizeOverflow`, `AllocationFailed`, `InvalidAlignment` | Storage arithmetic, allocation, or arena alignment failure |
| `EmbeddedZero`, `EmptySeparator`, `InvalidUTF8` | Invalid C text, splitting separator, or UTF-8 |
| `ConversionOutOfRange`, `InvalidNumber`, `InvalidBoolean` | Numeric conversion or text parsing failed |
| `InvalidPath`, `NotFound`, `PermissionDenied`, `PathFailed` | Path validation, lookup/access, or canonicalization failed |
| `FileOpenFailed`, `FileReadFailed`, `FileWriteFailed`, `FileCloseFailed` | File operation failed |
| `FileSyncFailed`, `FileRenameFailed`, `FileRemoveFailed`, `MetadataFailed` | Persistence, replacement/removal, or metadata operation failed |
| `StreamReadFailed`, `StreamWriteFailed`, `StreamFlushFailed`, `StreamModeFailed` | Stream I/O, flush, or binary-mode setup/restoration failed |
| `InvalidEnvironmentName`, `EnvironmentLookupFailed` | Invalid environment lookup name or native lookup failure |
| `InvalidState`, `Unsupported` | Handle/state misuse or unsupported operation |
| `ProcessSpawnFailed`, `ProcessWaitFailed`, `ProcessTerminateFailed` | Child lifecycle operation failed |
| `InvalidCodePoint` | Value is not a Unicode scalar |
| `RandomnessFailed` | Checked OS secret-seed acquisition failed when a map first needs storage; no weak fallback |

This list is an identity catalogue, not a claim that each function can return
every identity. Consult the operation's page for failure results and effects
that may already have occurred. `std::terminal` also defines its own error enum.

Unclassified failures and `nil` both have `kind == nil`; use `err != nil` to
test presence. Matching an `error` always needs a final `_` fallback.
See [language errors](../reference/errors.md) for payload storage and lifetime.
