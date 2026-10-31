# Strings and explicit memory management

## Values and lifetime

`string` is a read-only borrowed view: a byte pointer and a byte length. Copying
or passing it copies the view, not the bytes. String literals have static
storage. A zero-initialized string, including an omitted struct field, is empty.
An empty view may have a null pointer; it compares equal to every other empty
view and can be printed or sliced without dereferencing that pointer.

Text is UTF-8 by convention. Use [std::unicode::utf8](unicode.md) for strict validation,
decoding, counting, and borrowed iteration. No operation here validates UTF-8, normalizes text,
or folds case. `==` and `!=` compare byte lengths and contents. Embedded zero
bytes are ordinary string data. `io::println` writes every byte and then one
newline, including for empty strings. Runtime panic and error messages are also
length-aware. `error("")` is an error; it is distinct from `nil`.

Read-only APIs do not enforce memory safety. Low-level pointers, aliases, and C
code can still modify storage. Keep every borrowed view within its storage's
lifetime. Returning a borrowed view does not extend that lifetime.

The current C backend and libc bindings require a 64-bit `size_t`, checked by
the generated C. Public lengths, sizes, offsets, and alignments use `u64`.

## Borrowed-string operations

Import `std::strings`.

| API | Result | Allocation and lifetime |
| --- | --- | --- |
| `len(text string)` | `u64` byte length | No allocation |
| `slice(text string, start u64, end u64)` | `(string, error)` | No allocation; borrows the original storage |
| `from_c(text c::ConstCString)` | `string` | Scans to the first zero; borrows the C storage |

Slicing uses an inclusive start and exclusive end. `start > end` or
`end > len(text)` returns an error and an empty result. Empty ranges, including
`[len, len)`, are valid. A slice may split a UTF-8 sequence.

`from_c` requires a valid non-null, zero-terminated pointer. The caller is
responsible for its validity and lifetime. It does not copy or free storage.

## Individually owned text

`strings::Owned` owns a heap allocation backed by libc `malloc`/`free`.

| API | Result | Ownership and allocation |
| --- | --- | --- |
| `clone(text string)` | `(Owned, error)` | Copies bytes into independently owned storage; empty text needs no allocation |
| `concat(left string, right string)` | `(Owned, error)` | Copies both inputs into new owned storage; two empty inputs need no allocation |
| `view(owner *Owned)` | `string` | Borrows the owner's storage; no allocation |
| `free(owner *Owned)` | Nothing | Releases storage and resets this owner to empty; empty/repeated cleanup is a no-op |

Operations do not edit or resize existing text. Allocation failures and size
overflow return errors with empty results. Inputs remain valid when creation
fails. Error messages use static storage and need no cleanup.

Assignment and parameter passing retain ordinary shallow-copy behavior. Copying
an `Owned` duplicates its pointer, not its allocation. Arrange one lifecycle
owner. Freeing one copy does not update other copies or borrowed views, and
freeing a stale copy can double-free. Use `clone(view(&owner))` for an independent
allocation. No moves, non-copyable types, or borrow checking are enforced.

```rust
text := strings::clone("Hello") @fail
defer strings::free(&text)
view := strings::view(&text)
part := strings::slice(view, 0, 5) @fail
io::println(part)
```

## C strings

Import `std::ffi::c` for `c::CString` (mutable C characters) and
`c::ConstCString` (read-only C characters). These are zero-terminated character
pointers, separate from `string`, and do not imply ownership. A literal in a C-string context denotes static
zero-terminated storage; embedded zeros in such a literal are rejected. Existing
`string` values require an explicit conversion.

`cstring` is no longer a built-in type or reserved word. Replace `cstring` with
`c::CString` and `const cstring` with `c::ConstCString`. Aliases preserve
these types; `*c::ConstCString` is a writable pointer output slot, while
`const *c::ConstCString` also protects that slot. No lifetime is extended.
The C-text types are not exported from `std::libc`; that package contains actual
C standard-library function bindings.

| API | Result | Ownership and allocation |
| --- | --- | --- |
| `to_c(text string)` | `(COwned, error)` | Allocates a copy plus a trailing zero, even for empty text |
| `c_view(owner *COwned)` | `c::CString` | Borrows the conversion buffer; no allocation |
| `free_c(owner *COwned)` | Nothing | Releases the buffer and resets this owner; empty/repeated cleanup is a no-op |
| `to_c_in(arena *memory::Arena, text string)` | `(c::CString, error)` | Copies plus a trailing zero into arena-owned storage |

Both conversions reject embedded zeros rather than truncating text. Their
allocation and overflow failures return empty results. `COwned` has the same
shallow-copy responsibilities as `Owned`. An arena conversion remains valid
until the arena is reset or destroyed; do not free its pointer individually.

```rust
ctext := strings::to_c(view) @fail
defer strings::free_c(&ctext)
libc::printf("%s\n", strings::c_view(&ctext))
```

C bindings such as `printf`, `puts` (exposed as `libc::println`), `strlen`, and
`strcmp` take `c::ConstCString`. Passing borrowed `string` data to a C variadic call is
rejected; convert it first. Use `io::println` to print Telia strings directly.

## Heap allocation

Import `std::memory`.

| API | Result | Ownership and allocation |
| --- | --- | --- |
| `alloc(size u64)` | `(rawptr, error)` | Allocates uninitialized libc heap storage; zero size returns null with no error |
| `free(ptr rawptr)` | Nothing | Releases an allocation; null is a no-op; does not reset caller variables |

The caller owns successful non-null allocations and must free each exactly
once. Allocation failure returns null and an error. Use the owning string APIs
for text, which also track byte lengths and reset their owner values on cleanup.

## Arenas

An arena owns allocations with a shared lifetime. Creation is lazy: it allocates
no backing blocks and currently cannot fail. Blocks are added on demand from
the default heap, and existing allocations never move during growth.

| API | Result | Ownership and allocation |
| --- | --- | --- |
| `arena_create()` | `(Arena, error)` | Returns an empty arena; no allocation |
| `arena_alloc(arena *Arena, size u64, alignment u64)` | `(rawptr, error)` | Uses or grows backing blocks; the arena owns the result |
| `arena_reset(arena *Arena)` | Nothing | Invalidates all existing allocations/views and retains blocks for reuse |
| `arena_destroy(arena *Arena)` | Nothing | Invalidates views, releases all blocks, and empties this arena |
| `strings::clone_in(arena *Arena, text string)` | `(string, error)` | Copies bytes into arena storage; empty text needs no allocation |
| `strings::concat_in(arena *Arena, left string, right string)` | `(string, error)` | Copies both inputs into arena storage; two empty inputs need no allocation |

Alignment must be a nonzero power of two. A zero-size allocation returns null
without allocating, after alignment validation. Invalid alignment, size overflow,
and failed growth return errors without invalidating earlier allocations.
Destroying an empty arena is a no-op; it can be used again afterward.

There is no individual arena free. Arena copies are shallow and share backing
storage; arrange a single lifecycle owner. Resetting or destroying one copy does
not update other copies. The same lifetime rules apply to arena-backed C strings.

```rust
arena := memory::arena_create() @fail
defer memory::arena_destroy(&arena)
text := strings::clone_in(&arena, "Hello") @fail
greeting := strings::concat_in(&arena, text, ", Telia!") @fail
io::println(greeting)
```

The `*_into`, `*_state`, `free_text`, and `c_pointer` declarations in the library
sources are implementation plumbing for the private C helpers. Their output
parameters must be empty destinations; they do not clean up an existing value.
Use the owning/view APIs above in applications.

## Deferred cleanup

`defer` is block-scoped and runs in reverse declaration order. Normal block exit
runs that block's registered cleanup. A return evaluates and copies its result,
then runs registered cleanup in each exited scope, innermost first. Deferred
statements after an executed return are not registered. Deferred expressions
use their declaring variables at cleanup time; they do not capture argument
values when registered.

Do not return a view whose owner is freed by deferred cleanup. Copying the
return value preserves only the view, not the allocation. `@fail` terminates
the process on an error; it does not unwind deferred cleanup. Handle errors with
`@catch` and return normally when scoped cleanup must run.

## Verification

Run the compiler suite with `go test ./...`. Memory fixtures also run with
UndefinedBehaviorSanitizer and deterministic allocation-failure/leak counters.
On hosts with a working AddressSanitizer runtime, run:

```sh
TELIA_TEST_SANITIZERS=address,undefined go test ./tests/integration -run TestStringMemoryWithNativeSanitizers -count=1
```

The sanitizer subprocess has a timeout so a broken host runtime cannot hang
the test suite. LeakSanitizer availability depends on the host toolchain; the
allocation counters independently verify the exercised cleanup paths.
