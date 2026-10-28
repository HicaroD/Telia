# Header-owned opaque pointers

Declare an incomplete C type inside a header-owned `extern` block:

```rust
@[header="library/window.h", default_cc="c"]
extern {
    @[link_name="LibraryWindow"]
    opaque Window

    @[link_name="library_window_create"]
    fn create(value u32) *Window
}
```

`header` includes the actual C header before generated type and function
declarations. It uses the native C compiler's include search paths; absolute
header paths are also accepted. `link_name` selects the header's type name;
without it, the Telia name is used. Multiple attributes may share one `@[...]`.

Opaque types have nominal identity. Typed pointers can cross the C boundary,
including through aliases, without replacing their pointees with `void`.
Opaque values cannot be constructed, sized, passed/returned by value, copied
out of pointers, or accessed through fields. Use C-declared accessors instead.
Pointers do not imply ownership or automatic destruction: follow the library's
allocation, lifetime, and thread-safety rules.

The regression fixture uses an incomplete typedef in `opaque.h` and a real
definition in separately compiled `opaque.c`. It exercises typed arguments and
returns, links, and executes with incompatible-pointer diagnostics promoted to
errors on the local macOS arm64 host. This verifies those pointer interactions,
not foreign record layout or untested platforms.

## Further interoperability

See [header-owned C interoperability](foreign-interop.md) for records, unions,
enums, constants, fixed-array fields, pointer qualifiers, callbacks, C scalars,
variadic validation, and borrowed native `va_list` forwarding. These extensions
have local C fixtures, not a completed cross-platform ABI guarantee. Native
`va_list` lifecycle operations are implemented too; broader declarators and full
target verification remain pending regardless of tracking-issue closure.
