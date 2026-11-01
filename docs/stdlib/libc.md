# `std::libc`

Import `std::libc` for low-level native C calls. Prefer `std::io`, `std::strings`,
`std::memory`, `std::process`, and `std::strconv` for length-aware, checked APIs.
These bindings expose C behavior, not a safe wrapper or a general C library.

Parameter/result aliases below come from [`std::ffi::c`](ffi-c.md).

| Signature (prefix with `libc::`) | Contract |
| --- | --- |
| `printf(format c::ConstCString, args @c ...) c::c_int` | C formatted stdout output; format must match promoted native arguments |
| `vprintf(format c::ConstCString, args c::c_va_list) c::c_int` | C formatting from a valid borrowed native argument list |
| `puts(text c::ConstCString) c::c_int` | Write C text and a newline |
| `println(format c::ConstCString) c::c_int` | Alias of C `puts`, not `io::println` |
| `getenv(name c::ConstCString) c::CString` | Borrowed native environment pointer or null; do not free |
| `atoi(text c::ConstCString) c::c_int` | Unchecked native integer parsing; no reliable error distinction |
| `malloc(size c::c_size_t) rawptr` | Uninitialized heap allocation or null |
| `calloc(num c::c_size_t, size c::c_size_t) rawptr` | Zeroed allocation or null |
| `realloc(ptr rawptr, new_size c::c_size_t) rawptr` | Resize C heap allocation; failure for nonzero size preserves old storage |
| `free(ptr rawptr)` | Release matching C allocation; null is harmless; does not reset caller |
| `exit(code c::c_int)` | Terminate the process; pending Telia defers do not run |

The named `strings` extern namespace adds a qualifier:

| Signature (prefix with `libc::strings::`) | Contract |
| --- | --- |
| `strcmp(str1 c::ConstCString, str2 c::ConstCString) c::c_int` | Negative/zero/positive lexicographic C-byte comparison |
| `strlen(str c::ConstCString) c::c_size_t` | Bytes before the first zero; requires valid terminated storage |
| `strcat(str1 c::CString, str2 c::ConstCString) c::CString` | Append to writable destination with sufficient capacity; caller ensures no invalid overlap |

Example expression: `libc::strings::strlen("Telia")` returns five. See the
[complete interoperability example](../guide/interop.md).

Text arguments require valid zero-terminated pointers; embedded zeros terminate
C text rather than remaining ordinary data. Native pointers and borrowed
environment storage must satisfy C's lifetime rules. Keep one cleanup owner for
allocated pointers, and never free stack, literal, arena, or interior pointers.
Zero-size allocation/reallocation and OS text behavior follow the native C
implementation; use higher-level APIs when you need portable checked contracts.

Formatted output uses C variadic promotions and format strings. A mismatch may
cause undefined behavior, not a Telia `error`. Native I/O failure indicators and
errno are not translated into classified errors. The `mem` extern namespace is
currently empty.

Sources: [`stdio.t`](../../base/std/libc/stdio.t),
[`stdlib.t`](../../base/std/libc/stdlib.t),
[`string.t`](../../base/std/libc/string.t).
