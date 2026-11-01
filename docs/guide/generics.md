# Generic functions and types

Prerequisite: [errors](errors.md).

Generics reuse a declaration with different concrete types. Type parameters
appear in angle brackets; the compiler specializes the declaration for each
combination used by the program.

```telia
package main

use "std::io"

struct Pair<A, B> {
    first A
    second B
}

fn identity<T>(value T) T {
    return value
}

fn main() i32 {
    number i32 := 42
    pair := Pair<string, i32>{first: "answer", second: identity(number)}
    if pair.second == 42 {
        io::println(pair.first)
    }
    return 0
}
```

Expected output:

```text
answer
```

Here `identity(number)` infers `T = i32` from its argument. You can write
`identity<i32>(number)` explicitly. Inference can also use expected result types.
Use `_` for a position to infer in a partially explicit list, such as
`Pair<string, _>{first: "answer", second: number}`.

Insufficient or conflicting evidence is an error. The compiler will not invent
unknown types or silently widen conflicting concrete argument types. Supply
explicit arguments when a type appears only in the result and there is no
expected type.

Generics add no hidden allocation or deep copying. A generic record holding an
owner still needs explicit cleanup. Operations that depend on a type parameter
are checked when the template is specialized; a numeric operation is not a
promise that the function works for strings. There is no trait/interface
constraint system introduced by these generics.

**Try it:** replace `identity(number)` with `identity<i32>(number)`.

Reference: [generic inference, checking, and limits](../reference/generics.md).
Next: [C interoperability](interop.md).
