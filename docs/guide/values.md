# Variables, constants, and basic types

Prerequisite: [programs and imports](programs.md).

Use `:=` to declare a variable and `=` to change it. An optional type follows
the variable name, without a colon. `const` declares a compile-time value that
cannot be reassigned.

```telia
package main

use "std::io"

const VotingAge i32 = 18

fn main() i32 {
    age i32 := 17
    age = age + 1
    name := "Ada"
    eligible := age >= VotingAge
    if eligible {
        io::println(name)
        io::println("can vote")
    }
    return 0
}
```

Expected output:

```text
Ada
can vote
```

An inferred variable still has a fixed type: `name` is a `string`, not a container
that can later hold an integer. Use explicit types at boundaries where the width
matters, especially C interfaces and collection indices.

| Type | Use |
| --- | --- |
| `bool` | `true` or `false` |
| `i8` … `i128` | Signed integers of the named width |
| `u8` … `u128` | Unsigned integers of the named width |
| `int`, `uint` | Default 32-bit signed/unsigned integer types |
| `f32`, `f64`, `float` | Floating-point values; `float` is 32-bit |
| `string` | Read-only borrowed byte view, not owned text |

The widths listed above are the compiler's current definitions, not C's native
`int`/`long` rules. C has [separate scalar types](../stdlib/ffi-c.md).

A literal such as `"Ada"` stays alive for the whole program. Copying a `string`
copies the view, not its bytes. Later chapters explain allocated text.

**Common mistakes:** redeclaring a name with `:=` in the same scope; using
`var`, `let`, or `name: Type` syntax from another language; assuming all integer
types freely convert to one another. Use [explicit conversions](../reference/conversions.md)
when a conversion is needed.

**Try it:** change `age` to 16 and observe that nothing is printed.

Reference: [values and types](../reference/values.md),
[constants](../reference/constants-and-bitwise.md).
Next: [operators and control flow](control-flow.md).
