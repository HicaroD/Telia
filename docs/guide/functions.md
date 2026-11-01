# Functions and multiple results

Prerequisite: [control flow](control-flow.md).

A function names its parameters and places its result type after `)`. Omit the
result type for a function that returns no value. A tuple result describes
multiple returned values.

```telia
package main

use "std::io"

fn divide(value i32, divisor i32) (i32, bool) {
    if divisor == 0 {
        return 0, false
    }
    return value / divisor, true
}

fn main() i32 {
    quotient, ok := divide(12, 3)
    if ok and quotient == 4 {
        io::println("four")
    }
    return 0
}
```

Expected output:

```text
four
```

`quotient, ok := ...` receives both results in order. The number of destinations
must match the number of results. Functions pass ordinary values by copy. Copies
of pointer-like values still refer to the same storage; passing a value does not
deep-clone its resources.

To share writable state explicitly, pass a pointer, as described in the
[memory chapter](memory.md). Do not assume that assigning to a parameter changes
the caller's variable.

Functions can also be passed as callbacks using a concrete `fn(...) Result`
type. These are not capturing closures; see the [function reference](../reference/functions.md).

**Try it:** call `divide(12, 0)` and add an `else` branch that prints a message.
Then learn how to test this function in the [testing guide](../getting-started/testing.md).

Next: [structs and enums](data-types.md).
