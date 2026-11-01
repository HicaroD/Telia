# `std::ffi::c`

Sources: [`types.t`](../../base/std/ffi/c/types.t),
[`stdarg.t`](../../base/std/ffi/c/stdarg.t). Import `std::ffi::c`.
This package names ABI types; declarations allocate nothing and confer no
ownership by themselves.

## Native scalar aliases

| Aliases | Corresponding C types |
| --- | --- |
| `c_char`, `c_schar`, `c_uchar` | `char`, signed/unsigned `char` |
| `c_short`, `c_ushort` | signed/unsigned `short` |
| `c_int`, `c_uint` | signed/unsigned `int` |
| `c_long`, `c_ulong` | signed/unsigned `long` |
| `c_longlong`, `c_ulonglong` | signed/unsigned `long long` |
| `c_size_t`, `c_ptrdiff_t` | `size_t`, `ptrdiff_t` |
| `c_wchar_t` | `wchar_t` |
| `c_intptr_t`, `c_uintptr_t` | `intptr_t`, `uintptr_t` |
| `c_float`, `c_double` | `float`, `double` |

Widths and signedness follow the target ABI; do not substitute fixed-width
Telia types based on one machine. Aliases retain C scalar provenance for foreign
checks. See [C scalars](../reference/foreign-interop.md#c-scalars) for operations,
generics, and conversions.

## C text

`CString` is mutable zero-terminated C text; `ConstCString` is read-only through
that pointer. Neither stores a Telia byte length or owns its storage. Literal
C-text use is supported, but arbitrary Telia `string` views require explicit
conversion. Fragment with `std::strings` imported, in an error-returning function:

```rust
text := strings::to_c("native text") @try
defer strings::free_c(&text)
pointer := strings::c_view(&text)
```

`pointer` borrows the owner. Do not free it separately or retain it after owner
cleanup. See [C-string contracts](strings-and-memory.md#c-strings).

## Native argument lists

`c_va_list` is a borrowed native argument-list type. `VaList` is an explicitly
managed owner for native variadic traversal/copying. These are compiler-recognized
types, not interchangeable records or serializable pointers. `va_start`,
`va_copy`, `va_arg`, `va_borrow`, and `va_end` have lifecycle and context
restrictions described in [native variadics](../reference/foreign-interop.md#variadics-and-native-va_list-lifecycle).
Never invent a `va_list` layout or shallow-copy an owner as a substitute for
`va_copy`.
