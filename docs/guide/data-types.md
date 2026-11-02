# Structs, enums, and matching

Prerequisite: [functions](functions.md).

A struct groups named fields. An enum selects one named variant, optionally
carrying fields of its own. Use `match` to inspect a variant and produce a value.

```telia
package main

use "std::io"

struct User {
    name string
    age i32
}

enum Greeting {
    Anonymous
    Named {
        name string
    }
}

fn main() i32 {
    user := User{name: "Ada", age: 18}
    greeting := Greeting::Named{name: user.name}
    text := match greeting {
        Greeting::Anonymous => "Hello, stranger"
        Greeting::Named{name} => name
    }
    io::println(text)
    return 0
}
```

Expected output:

```text
Ada
```

Construct structs with `Type{field: value}` and access fields with `value.field`.
Omitted struct fields are zero-initialized. Enum payload construction instead
requires every payload field exactly once. `Greeting::Named{name}` in a pattern
binds a copy of the field to the name `name`.

Enum matches must cover every variant or end with `_ => ...`. A multi-statement
arm uses a block and `yield` to provide its result, not `return` (which returns
from the function). Payload enums have automatic equality when all payload
fields support it; see [equality rules](../reference/expressions.md#equality-and-key-eligibility).

**Common mistake:** constructing a literal directly in a control-flow header
without parentheses. Use `match (Greeting::Named{name: "Ada"}) { ... }` because
the unparenthesized brace is otherwise the body's opening brace.

**Try it:** construct `Greeting::Anonymous` instead.

Reference: [structs and types](../reference/values.md),
[literal rules](../reference/literals.md), [enums and matching](../reference/enums.md).
Next: [collections](collections.md).
