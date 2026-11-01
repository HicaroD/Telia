# Errors and propagation

Prerequisite: [ownership and cleanup](memory.md).

`error` is an ordinary value. `nil` means success; `error("message")` constructs
a failure. A function often returns `(Value, error)` so callers can handle the
failure or propagate it.

```telia
package main

use "std::io"
use "std::strconv"

fn parse_positive(text string) (i64, error) {
    number := strconv::parse_i64(text) @try
    if number <= 0 {
        return 0, error("expected a positive number")
    }
    return number, nil
}

fn main() i32 {
    number, err := parse_positive("-2")
    if err != nil {
        io::println(err.msg)
        return 0
    }
    if number > 0 {
        io::println("positive")
    }
    return 0
}
```

Expected output:

```text
expected a positive number
```

Choose how to handle each call:

| Form | Behavior on failure |
| --- | --- |
| `value, err := operation()` | Inspect and handle `err` explicitly |
| `value := operation() @try` | Return the error from this function, running defers |
| `operation() @catch err { ... }` | Execute your handler |
| `operation() @fail` | Print the message and exit with status 1, skipping pending defers |

`@try` requires an enclosing result of `error` or a tuple ending in `error`.
When a catch's success value is consumed, every handler path must return from
the enclosing function. Merely setting a local variable cannot replace the
failed call's result.

For stable failure categories, use an `enum Name: error` and explicitly construct
`error(Name::Variant, "message")`. Compare `err.kind`, not human-readable text.
An error match requires a final `_` arm because `error` can contain other
categories, unclassified failures, or `nil`.

Error messages and borrowed payload fields must remain alive until the error is
finished being used. Returning an error does not extend the lifetime of its text.
Even `error("")` is a failure; only `nil` is success.

**Try it:** change the input to `"2"`, then to `"not a number"`.

Reference: [errors](../reference/errors.md), [standard identities](../stdlib/errors.md).
Next: [generics](generics.md).
