# Expressions and statements

## Operators

From strongest to weakest binding:

| Group | Operators |
| --- | --- |
| Postfix | Calls, `.field`, `[index]`, `[start:end]` |
| Unary | `not`, `~`, `-`, `&` (address), `*` (dereference) |
| Multiplicative | `*`, `/` |
| Additive | `+`, `-` |
| Shifts | `<<`, `>>` |
| Bitwise AND, XOR, OR | `&`, then `^`, then `\|` |
| Comparison | `==`, `!=`, `<`, `<=`, `>`, `>=` |
| Logical AND | `and` |
| Logical OR | `or` |

Parentheses override grouping. Comparisons are not chained; write
`low <= value and value < high`. Conditions use booleans, not integer truthiness.
`and` and `or` short-circuit; an unselected right-hand side is not evaluated.
Integer `/` is division, not a floating-point conversion. There are no `%`,
`++`, `--`, or compound-assignment operators. String concatenation is a library
operation, not `+`.

Use [explicit conversions](conversions.md) for typed integer boundaries.
[Bitwise rules](constants-and-bitwise.md#bitwise-operators) define operand
compatibility, signed shifts, and invalid shift counts. Do not assume runtime
integer overflow is a checked `error` or a portable wrapping operation: use
range checks where an algorithm requires them. Constant arithmetic is checked.

Receiver/index expressions are evaluated once per indexed access. Function
calls and error propagation stay at their expression evaluation site; an
unselected branch does not execute them. See [matching](enums.md) for arm order
and [errors](errors.md) for propagation inside expressions.

## Branches

```rust
if ready {
    work()
} elif retry {
    wait()
} else {
    stop()
}
```

`if` is a statement, not a ternary expression. `match` produces a value and
supports enum, integer, and error patterns. A literal appearing directly in a
control-flow header needs [parentheses](literals.md#control-flow-headers).

## Loops

```rust
while ready {
    work()
}

for index := 0; index < 10; index = index + 1 {
    work()
}

for item in values {
    inspect(item)
}
```

`while` tests before each iteration. Three-clause `for` runs its initializer
once, tests before each iteration, and runs its update after the body.
`for item in values` accepts arrays and slices, evaluates its source once, and
binds a shallow element copy. Lists iterate through their borrowed view; there
is no generic user-defined iteration protocol.

`break` exits the nearest loop; `continue` starts its next iteration (running a
three-clause update first). Neither accepts a label or result. Both run defers
for scopes they exit. See [collection iteration](collections.md#iteration).

## Return, yield, and defer

`return` exits the function, copying its results before deferred cleanup.
`yield` supplies a match-arm block's result, not a function result. Completing
match-arm paths must yield compatible values. `defer Statement` registers
block-scoped cleanup; it runs in reverse registration order on ordinary scope
exit. Deferred code cannot jump back out into an enclosing loop or match.
See [memory](memory.md) for lifetime implications and fatal-exit exceptions.
