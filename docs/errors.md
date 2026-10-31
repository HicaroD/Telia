# Errors

`error` is an open-ended value containing presence, an optional classified
identity and payload, and a borrowed `string` message. `nil` means success.
`error("")` is a failure, even when its message has no backing storage.

## Declaring and constructing failures

```rust
enum FileError: error {
    NotFound {
        path string
    }
    PermissionDenied
}

fn read_file(path string) (string, error) {
    return "", error(
        FileError::NotFound{path: path},
        "file does not exist",
    )
}
```

Only enums declared with `: error` can classify errors. Their ordinary enum
construction and payload validation rules still apply. Generic error enums are
supported; each concrete specialization has its own identity. Foreign enums
cannot be error enums. The marker uses the default integer tag representation.

Construction is explicit: a variant does not implicitly convert to `error`.
`error(message)` remains an unclassified failure. Messages may be literals or
other `string` expressions; they are context, not stable failure identifiers.
Constructor arguments and payload initializers run once in source order, at
the expression's original evaluation site.

## Identity and matching

```rust
if err.kind == FileError::NotFound {
    io::println(err.msg)
}

description := match err {
    FileError::NotFound{path} => path
    FileError::PermissionDenied => "permission denied"
    _ => err.msg
}
```

`err.kind` is read-only. Compare it with `==` or `!=`, another error's kind,
or a bare error-enum variant. Payload-bearing variants need no initializer in
these identity comparisons. `kind := err.kind` saves an inferred identity value
that supports the same comparisons. Payloads and messages do not affect kind
equality. Identically named enums/variants in different packages remain distinct.

Unclassified failures and success both have `kind == nil`; use `err != nil`
to distinguish failure from success. Classification does not change when the
message is reassigned.

An `error` match **always requires a final `_` fallback**, even if all variants
of an error enum are covered. Other error enums, unclassified failures, and
`nil` can reach that arm. Patterns may mix error enums and use ordinary named
payload bindings, renaming, and `..`. The scrutinee is evaluated once. Bindings
copy fields with normal assignment semantics; they do not mutate the error.

## Copies, storage, and lifetime

The current C backend stores payloads inline, using a union of the concrete
error-enum layouts in the compiled program. Every `error` therefore has room
and alignment for the largest of those layouts. This is an internal,
whole-program representation, not a stable C ABI or serialization format.

Constructing, copying, matching, passing, or returning an error does not allocate
payload storage and needs no implicit cleanup. Construction itself cannot fail
because of allocation failure; evaluating its arguments can still allocate or
fail according to those expressions' contracts. Existing errors preserve their
identity, payload, and message through assignment, arguments, returns, tuples,
and explicit propagation with `return err`.

By-value layout cycles are rejected. In particular, an error-enum payload cannot
contain `error` directly or through another inline aggregate. Pointer indirection
breaks the layout cycle, but does not grant ownership or extend lifetimes.

Copies are shallow for strings, slices, pointers, and owner records. Messages
and borrowed payload fields must not outlive their backing storage. There is
no borrow checker, reference counting, deep copying, or automatic destruction.
Arrange one cleanup owner for explicitly owned fields. Returning an error copies
its inline value **before** deferred cleanup, but does not preserve storage that
cleanup frees. Match `yield` similarly saves the result before arm cleanup.
See [strings and memory](strings-and-memory.md) for ownership and deferred cleanup.

## Handling and discarding

Explicit returns, `@catch`, and `@fail` remain available. When a catch's success
value is consumed, every handler path must return from the enclosing function;
falling through would otherwise expose a zero-initialized success value. A
handler cannot replace that value by merely assigning a local. Error-only and
statement-only catches may complete normally because no success value is used.

`@fail` is a fatal boundary, not error propagation. On error it prints the
length-aware message and exits with status 1, without running pending defers in
the current function or its callers. Completed scopes have already run their
cleanup. On success, execution continues normally. If deferred cleanup itself
uses `@fail` and fails, remaining cleanup does not run either.

Use `@catch` with `return err` to propagate errors through ordinary returns when
cleanup must run. Each function exited this way runs its registered defers,
innermost scope first and in reverse registration order within each scope.
Place `@fail` outside resource-owning functions so they can clean up before fatal
reporting:

```rust
package main

use "std::strings"
use "std::io"

fn run() error {
    text := strings::clone("Hello") @catch err {
        return err
    }
    defer strings::free(&text)

    part := strings::slice(strings::view(&text), 0, 5) @catch err {
        return err
    }
    io::println(part)
    return nil
}

fn main() i32 {
    run() @fail
    return 0
}
```

Here `run` has freed `text` before `main` reports any fatal error. Propagated
error messages and borrowed payloads must still outlive cleanup;
do not return an error whose message is backed by the allocation being freed.

As before, a call may silently discard its returned error; this change does not
introduce a mandatory error-handling policy. Wrapping and cause chains are not
provided.

## Standard-library identities

Import `std::errors` and compare with `errors::Error::Variant`:

| Variant | Meaning |
| --- | --- |
| `OutOfBounds` | Invalid string byte, slice/list element, or UTF-8 decode offset |
| `InvalidRange` | Invalid string or borrowed-slice range |
| `SizeOverflow` | Allocation-size, list-capacity, or list-length arithmetic overflow |
| `AllocationFailed` | Heap or arena backing allocation failed |
| `InvalidAlignment` | Arena alignment is zero or not a supported power of two |
| `EmbeddedZero` | C-text conversion would contain an embedded zero byte |
| `EmptySeparator` | String splitting received an empty separator |
| `InvalidUTF8` | Invalid/truncated UTF-8, including invalid continuations or code points |

These payload-free identities are shared across runtime helpers and standard
memory, string, slice, list, and UTF-8 APIs. APIs preserve identities when
forwarding another failure. Messages retain human-readable detail and are not
stable identifiers. Vendor-library wrappers keep their existing message-only
errors unless separately documented.
