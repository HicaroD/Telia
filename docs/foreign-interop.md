# Header-owned C interoperability

## Language boundary

Native strings, arrays, structs, tagged enums, generics, and function values keep
Telia semantics. Typed pointers and read-only access remain available to native
systems code; being low-level does not make a construct C-specific.

C text, ABI scalars, and variadic lifecycle operations are imported through
`std::ffi::c`. Actual C standard-library functions (`malloc`, `printf`, `getenv`,
`atoi`, and others) live in `std::libc`, whose declarations use the imported C
types. Applications should import these bindings rather than redeclaring libc
functions in local `extern` blocks. Header-owned unions,
opaque types, and initializer-free constants belong to `extern` declarations.
C callback contracts require explicit ABI annotations. `@c ...` and `va_list`
lifecycle operations remain foreign facilities, not a general ownership model.
The C backend's implementation of native runtime helpers does not turn native
`string` or function values into public C ABI types.

Migration:

- Replace `cstring` with `c::CString` and `const cstring` with
  `c::ConstCString`, importing `std::ffi::c` where needed.
- Annotate C callback aliases and their implementations with `@[cc="c"]` or
  the matching `cc_macro`. Extern defaults apply to prototypes, not aliases.
- Move `libc::c_int`, `libc::c_size_t`, and other scalar references to `c::c_int`,
  `c::c_size_t`, etc.; the type names themselves are unchanged.
- Move `libc::c_va_list`, `libc::VaList`, and `libc::va_*` operations to `c::`.
  There are no compatibility re-exports from `std::libc`.

Foreign storage and variadic lifecycle types still require compiler support;
this boundary does not imply that all foreign metadata has been removed from
the compiler's shared type representation.

## Standard-library calls

Applications import existing libc bindings instead of declaring their own:

```rust
use "std::libc"

fn main() i32 {
    configured := libc::getenv("TELIA_FRAME_LIMIT")
    if configured != nil {
        limit := libc::atoi(configured)
        libc::printf("%d\n", limit)
    }
    return 0
}
```

Type inference means this example does not need to import `std::ffi::c` directly.
Explicit C types do: `getenv` returns `c::CString`, and `atoi` returns `c::c_int`.
`getenv` borrows environment storage; do not free it. `atoi` keeps C's conversion
contract, not a checked Telia parsing API. `puts` is exported under its standard
name; the existing `libc::println` alias remains available.

## Verification scope

The implementations tracked in issue #93 have real, separately compiled
C regressions on **macOS arm64 only**. CI is configured to run the fixtures on
Linux x64, macOS Intel, macOS Apple Silicon, and Windows x64 with MinGW-w64;
configuration is not a claim that those runs have passed.

## Records, unions, constants, and enums

```rust
@[header="library.h", default_cc="c"]
extern {
    @[link_name="LibraryRecord"]
    struct Record {
        value i32
        data u8[4]
    }

    @[link_name="LibraryValue"]
    union Value {
        integer i32
        real f64
    }

    @[link_name="LIBRARY_ENABLED"]
    const Enabled u64

    @[link_name="LibraryStatus"]
    enum Status {
        @[link_name="LIBRARY_READY"]
        Ready
    }
}
```

The header owns C type identity, layout, packing, enum tags, and constant values.
Telia does not emit replacement definitions for these records. Declared field
types are checked against the header by the C compiler, with field-level Telia
diagnostics on mismatch. An optional foreign enum representation annotation is
checked for size and signedness; omitting it uses the actual header enum type.

Record literals construct the actual C type; omitted fields, including fields
not exposed in the Telia declaration, are zero-initialized. Union literals must
select exactly one member. Reading an inactive union member is not made safe by
Telia: follow the C library's active-member contract.

```rust
record := Record{value: 42, data: .{10, 20, 30, 40}}
record.data[1] = 42
value := Value{integer: 42}
```

Fixed-array fields use inline C storage, not Telia wrapper arrays. Element
indexing remains bounds checked. Use an explicit element pointer such as
`&record.data[0]` for C pointer parameters. Whole-array assignment/copying and
implicit array decay are not supported. Multidimensional fields and
pointer-to-array foreign declarators require a C adapter for now. Ordinary
Telia arrays and records are not automatically foreign ABI types.

Foreign declarations also work in named extern scopes. Foreign constants have
an explicit scalar type and **no initializer**. Their values remain C expressions;
they cannot currently drive Telia's compile-time scalar evaluator, native enum
tags, or array lengths. Declare the correct header constant type rather than
relying on narrowing conversions.

Foreign enums have no payloads or type parameters. They may contain unnamed
values, so every match requires a final `_` fallback, even when every declared
member has an arm. The private representation of native payload enums remains
unrelated to the foreign C enum ABI.

## Pointer qualifiers and C text

`const *T` protects the **immediate pointee**, not the pointer variable. Aliases
preserve the qualifier:

```rust
use "std::ffi::c"

type Text = c::ConstCString
fn read_text(text Text) {}
fn replace_text(output *Text) {}
fn inspect_texts(texts const *Text) {}
```

`c::ConstCString` is the NUL-terminated C-text type for a pointer to const
characters. `*Text` is a writable output slot holding such a pointer;
`const *Text` protects the slot too. A mutable pointer may gain immediate const,
but nested pointee qualifiers are not silently changed. Const-protected fields
cannot yield mutable addresses or mutable slice borrows. Copied mutable pointers
stored inside a const record still point to mutable storage.

The old parameter attribute has been removed and standard-library declarations
migrated. Imported `c::CString`, literal C text, and explicit owned Telia-string /
C-string conversion APIs remain available. Qualifiers do not imply allocation,
ownership, termination validation, or lifetime extension.

## C scalars

Public scalar aliases live in `std::ffi::c`: `c_char`, `c_schar`, `c_uchar`,
`c_short`, `c_ushort`, `c_int`, `c_uint`, `c_long`, `c_ulong`, `c_longlong`,
`c_ulonglong`, `c_size_t`, `c_ptrdiff_t`, `c_wchar_t`, `c_intptr_t`,
`c_uintptr_t`, `c_float`, and `c_double`.

The compiler queries its selected native GCC/Clang toolchain for scalar widths
and plain-char signedness. Generated C retains the actual C spelling, including
plain `char` and `long`. Build-time assertions reject ABI-changing flags that
disagree with the probed size or signedness instead of silently changing Telia
semantics. This is a native-toolchain strategy, not cross-compilation support.
Numeric values with the same width, signedness, and numeric kind are compatible
across the C boundary, including arithmetic results, indexed elements, and foreign
constants. For example, a Telia `f32` expression can be passed to a C `float`
parameter without using `c::c_float` in application code. This does not permit
implicit `f64`-to-`float` narrowing. Typed pointers (including address-of
expressions) and callback signatures retain their declared C scalar identities.

## Callbacks and calling conventions

```rust
@[cc="c"]
type Compare = fn(left rawptr, right rawptr) i32

@[cc="c"]
fn compare_items(left rawptr, right rawptr) i32 {
    return 0
}

callback Compare := compare_items
result := callback(left, right)
```

Bare, non-generic function names are function values. Native function types
(`type Compare = fn(...) Result`) remain native. Foreign callback aliases and
their Telia implementations must explicitly opt in with `@[cc="c"]` or the
matching `cc_macro`; a matching parameter list alone is not sufficient. Foreign
signatures and record fields reject native function types, including callback
output slots. The ABI also participates in generic specialization identity.
Callback parameters, record fields, returns, indirect calls, and `*Compare`
output slots are supported.
An output slot is a pointer to callback storage, not a second function type.
Callbacks can be initialized and compared with `nil`; check nullable callbacks
before invoking them. Calling a null or expired callback is not made safe.
Generic functions can contain independently specialized callback types; use a
non-generic wrapper when taking the address of a generic function itself.

`cc="c"` selects the native C convention; extern prototypes default to it, but
ordinary Telia functions and function types do not opt in implicitly. Other named
conventions are diagnosed rather than silently ignored. `cc_macro` places a
header-defined convention macro on a function or callback alias;
`default_cc` / `default_cc_macro` select defaults for extern members. Do not set
both alternatives on one declaration. Callback and function macro annotations
must match; the compiler does not guess whether different macros are equivalent.

There is no automatic captured closure, userdata allocation, or ownership
transfer. Keep callbacks and userdata alive until the library has stopped all
possible invocations; free userdata only under that library's ownership rules.
Foreign-thread callbacks run synchronously on the calling thread, with no
implicit dispatch or synchronization. Thread-affine APIs and shared state remain
the binding author's responsibility. The fixtures do not validate concurrent
callbacks, foreign-thread runtime behavior, or library-specific lifetime rules.

## Variadics and native `va_list` lifecycle

```rust
use "std::libc"
use "std::ffi::c"

fn forward(format c::ConstCString, args c::c_va_list) c::c_int {
    return libc::vprintf(format, args)
}
```

`libc::printf` already declares the `args @c ...` foreign variadic marker;
applications call it through the imported package, not a local redeclaration.

C variadic calls validate scalar/pointer/callback/payload-free-enum arguments.
Native C default promotions apply: narrow integers and booleans promote to
`int` where representable, and `f32` promotes to `double`. Aggregates and Telia
strings are rejected; use appropriate C storage and explicit conversions.
Format-string correctness and the callee's expected argument types remain the
caller's responsibility.

`c_va_list` lowers to the compiler's actual `<stdarg.h>` `va_list` parameter,
including native parameter adjustment. It is **not** `rawptr`. C-owned lists can
be borrowed by Telia callbacks and forwarded directly to C functions during
that call. Assignment copying, returning, retaining an address, constructing
ordinary borrowed-value storage, and generic borrowed-value specialization are
not supported. The original C owner remains responsible for its list's cleanup.

The following compiler-backed operations are exported by `std::ffi::c`:

| Operation | Contract |
| --- | --- |
| `va_start(args)` | Initialize a new local owner from this function's `args @c ...` marker. |
| `va_copy(borrowed)` | Initialize an independent owner without advancing the source cursor. |
| `va_arg<T>(&owner)` | Read the next value using its actual promoted C type. |
| `va_end(&owner)` | End this owner's native storage exactly once. |
| `va_borrow(&owner)` | Produce a call-scoped native borrowed list. |

```rust
use "std::ffi::c"

@[cc="c"]
fn first(marker i32, args @c ...) c::c_int {
    owner c::VaList := c::va_start(args)
    defer c::va_end(&owner)
    copy := c::va_copy(c::va_borrow(&owner))
    defer c::va_end(&copy)
    return c::va_arg<c::c_int>(&copy)
}
```

Both factories must directly initialize a **new local variable**. Owning
`c::VaList` storage is initialized in place using native `va_start`/`va_copy`;
it is never initialized by a struct assignment or raw-pointer cast. Owners are
noncopyable, cannot be returned or passed as ordinary values, and cannot expose
ordinary addresses. `&owner` is accepted specifically by the lifecycle intrinsics.
The operations themselves are not ordinary callback values. Generic function
specializations may use the operations with independently resolved read types.

Use `defer va_end(&owner)` for block-scoped cleanup, or explicitly end the owner
on every executed exit path. Double-end, reading/copying an ended or consumed
owner, and leaving its scope without `va_end` panic with a source location before
an invalid native operation. On a normal scope exit missing explicit cleanup,
the backend ends the native list before reporting the missing-cleanup panic.
Return and match `yield` run registered defers before this scope-exit check.
These dynamic lifetime checks use GCC/Clang cleanup attributes; they do not
unwind process termination or foreign `longjmp` across Telia frames.

Use a borrow directly in a call, for example
`libc::vprintf(format, c::va_borrow(&owner))`; it cannot be saved in an ordinary
local variable. Passing an owner's borrow to an ordinary callee conservatively
marks that owner consumed: only `va_end` remains valid afterward. The specific
operation `va_copy(va_borrow(&owner))` does **not** consume the source owner.
Consumption of raw C-borrowed parameters is still governed by the foreign API;
do not reuse one after forwarding unless the callee explicitly guarantees it.
No cursor is implicitly cloned.

`va_arg<T>` accepts promoted scalar/pointer/callback types, not aggregates,
enums, `bool`, narrow integers, or `f32`. Read default-promoted values using
`c::c_int` / `c::c_double`. The programmer must still supply the actual
argument type and count: reading beyond the list or using a mismatched type can
have C undefined behavior. Cursor operations in supported expressions preserve
source order, short-circuiting, loop reevaluation, and selected match execution.

`va_start` requires at least one fixed parameter and a supported final fixed
type unchanged by default promotions. It currently requires the native C
convention, not an unverified `cc_macro`; use a C adapter for convention-specific
variadic storage. The former raw-pointer LLVM variadic declarations have been
removed; they were not native C lifecycle operations.

The lifecycle fixtures include C-to-Telia variadic calls, C-borrowed cursor
copying, independent owners, typed reads, explicit/deferred cleanup, runtime
lifetime failures, generic specializations, and optimized UBSan execution on the
local macOS arm64 host. Windows sanitizer coverage is explicitly skipped because
that runtime is not configured; ordinary lifecycle fixtures remain in its matrix.

**Remaining:** broader complex declarators and actual results from the full
target matrix. Tracking-issue closure does not establish untested platform support.
