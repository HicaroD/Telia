# Calling C and native libraries

Prerequisite: [generics](generics.md), especially the ownership chapter.

Prefer a standard-library API when it covers your task. C interoperability is a
low-level boundary: declarations must match the native ABI, pointers must be
valid, and C's ownership rules still apply.

`std::libc` provides existing declarations for a small part of the C library:

```telia
package main

use "std::io"
use "std::libc"

fn main() i32 {
    length := libc::strings::strlen("Telia")
    if length == 5 {
        io::println("five C bytes")
    }
    return 0
}
```

Expected output:

```text
five C bytes
```

The literal can be used at this C-text boundary. An arbitrary borrowed Telia
`string` is not automatically a zero-terminated C string. Convert it with
`strings::to_c`, keep the resulting `COwned` alive for the call, and clean it up
with `strings::free_c`. Embedded zero bytes are rejected by that conversion.

Use `std::ffi::c` types for target-native C scalars; do not assume C `long` is
always `i64`. An `extern` block declares native functions and types. Header-owned
records, opaque types, callbacks, qualifiers, variadics, and linkage attributes
have additional rules in the [C reference](../reference/foreign-interop.md).

For packaged libraries such as SDL, use the documented vendor package rather
than inventing declarations or manually copying library files. See
[native packages](../stdlib/native-libraries.md) for build and deployment setup.

**Try it:** change the literal and update the expected byte count. Remember that
UTF-8 byte length is not code-point count.

Continue with the [cookbook](../cookbook/README.md) or use the
[language reference](../reference/README.md) for precise rules.
