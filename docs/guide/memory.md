# Pointers, ownership, and cleanup

Prerequisite: [collections](collections.md).

A pointer refers to another value's storage. `&value` takes its address; `*ptr`
accesses the pointed-to value. A pointer does not keep storage alive or grant
ownership. Telia does not have a borrow checker or automatic resource cleanup.

```telia
package main

use "std::io"
use "std::strings"

fn increment(value *i32) {
    *value = *value + 1
}

fn run() error {
    count i32 := 0
    increment(&count)
    text := strings::clone("owned text") @try
    defer strings::free(&text)
    if count == 1 {
        io::println(strings::view(&text))
    }
    return nil
}

fn main() i32 {
    run() @fail
    return 0
}
```

Expected output:

```text
owned text
```

`strings::clone` allocates an independently owned copy. The `Owned` record is
the cleanup handle; `strings::view` borrows its bytes as a `string`.
`defer strings::free(&text)` frees and empties the handle when `run` exits.

The error operators are explained next: `@try` returns failures while running
registered cleanup. `@fail` reports a fatal error without running pending
cleanup; placing it in `main` lets `run` clean up before fatal reporting.

`defer` is **block-scoped**, runs in reverse registration order, and reads its
variables when cleanup executes. It does not freeze arguments when registered.
Returns, `break`, and `continue` run cleanup for scopes they exit.

Copying an owner is not cloning its allocation. Keep one cleanup owner; freeing
two shallow copies can double-free. Borrowed strings and slices become invalid
when their storage is freed or reset. Do not return `strings::view(&text)` from
the example: its defer would free the bytes before the caller could use them.

To return an owner with a registered defer, explicitly transfer it with
`memory::take<strings::Owned>(&text)` after importing `std::memory`. This empties
the source, making its pending cleanup safe. Arenas provide another explicit
strategy for values that share a lifetime.

**Try it:** add a second `io::println` in a deferred statement before the free
registration; observe the reverse cleanup order.

Reference: [memory semantics](../reference/memory.md),
[ownership contracts](../stdlib/strings-and-memory.md#ownership-conventions-and-transfer).
Next: [errors](errors.md).
