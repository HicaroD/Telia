# Checked text conversions

Import `std::strconv` for decimal integer and boolean parsing and formatting:

```rust
fn parse_i64(text string) (i64, error)
fn parse_u64(text string) (u64, error)
fn parse_bool(text string) (bool, error)
fn format_i64(value i64) (strings::Owned, error)
fn format_u64(value u64) (strings::Owned, error)
fn format_bool(value bool) (strings::Owned, error)
```

## Parsing

Parsing borrows its input, allocates nothing, and requires the entire input to
match. Signed integers permit one optional ASCII `+` or `-`; unsigned integers
permit ASCII digits only. At least one digit is required. Leading zeros and
signed negative zero are accepted. The complete `i64` and `u64` ranges are
supported, including `-9223372036854775808` and `18446744073709551615`.

Whitespace, trailing data, embedded zero bytes, non-ASCII digits, separators,
nondecimal bases, and malformed UTF-8 bytes are rejected. There is no implicit
trimming or locale dependence. Use `strings::trim_space(text)` explicitly for
validated Unicode whitespace trimming before parsing.

Integer failures return zero and an error classified as
`errors::Error::InvalidNumber` for missing digits or invalid syntax, or
`errors::Error::ConversionOutOfRange` for overflow. Scanning stops at the first
invalid byte or overflowing digit, so overflow may be reported before later
invalid syntax. Arithmetic never wraps during accumulation, even on long inputs.

Booleans accept exactly `true` and `false`. Invalid input returns false with
`errors::Error::InvalidBoolean`; a successfully parsed `false` has no error.
Check the error rather than using the value as a success indicator. Compare
classified identities through `err.kind`.

## Formatting and ownership

Integer output is canonical decimal: no leading plus, unnecessary zeros,
whitespace, or locale-dependent punctuation. Boolean output is lowercase.
Formatting makes an independently owned copy, including boolean output.
Allocation failure returns an empty owner with
`errors::Error::AllocationFailed`, releasing any temporary storage.

Use `strings::view(&owner)` while the owner is alive, and
`strings::free(&owner)` exactly once per owned allocation. Cleanup empties the
owner and repeated cleanup of that same handle is harmless. Copies of an owner
are shallow, not independent owners; do not free both copies. See
[strings and memory](strings-and-memory.md).

```rust
value, err := strconv::parse_i64("+0042")
if err != nil {
    return 1
}
owner, format_err := strconv::format_i64(value)
if format_err != nil {
    return 2
}
defer strings::free(&owner)
io::println(strings::view(&owner)) // 42
```

Floating-point conversions, other bases, and arena formatting are not provided.
These text APIs are distinct from compiler-supported
[integer type conversions](../reference/conversions.md).

## Runnable example

From a configured source checkout:

```sh
make build
./telia build examples/strconv.t -o strconv-example
./strconv-example
```

The example explicitly trims Unicode whitespace, parses signed/unsigned integers
and a boolean, formats them, and cleans up each owner. Output:

```text
-42
18446744073709551615
true
```
