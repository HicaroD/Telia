# Operators and control flow

Prerequisite: [variables and basic types](values.md).

`if`, `elif`, and `else` select a branch. Conditions must be boolean. Use `and`,
`or`, and `not` for logic, not C's `&&`, `||`, and `!`.

```telia
package main

use "std::io"

fn main() i32 {
    total := 0
    for i := 0; i < 4; i = i + 1 {
        if i == 1 {
            continue
        }
        total = total + i
    }
    while total > 0 {
        io::println("counting")
        total = total - 1
        if total == 3 {
            break
        }
    }
    return 0
}
```

Expected output:

```text
counting
counting
```

The `for` initializer runs once; the condition runs before each iteration; the
update runs after the body (including after `continue`). Here the sum is five
because the iteration for one is skipped. The `while` loop then stops at three.
`break` exits the nearest loop; `continue` skips to its next iteration.

Arithmetic uses `+`, `-`, `*`, and `/`. Comparisons use `==`, `!=`, `<`, `<=`,
`>`, and `>=`. Write `i = i + 1`; there is no `++` operator. Use parentheses
when a mixed expression would be hard to read. There is no `%` operator.

`and` and `or` short-circuit: the right expression runs only if needed. Bitwise
`&`, `|`, `^`, `~`, `<<`, and `>>` operate on integers, not booleans.

A name declared in a loop body cannot be used after the loop. A loop initializer
also has loop-local scope. Later you will use `for value in values` for arrays
and slices.

**Try it:** remove the `continue` branch and work out how many lines will print.

Reference: [expressions and statements](../reference/expressions.md),
[bitwise operators and precedence](../reference/constants-and-bitwise.md#bitwise-operators).
Next: [functions](functions.md).
