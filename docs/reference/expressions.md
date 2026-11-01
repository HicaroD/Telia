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

## Equality and key eligibility

`==` and `!=` require compatible operands of the same concrete type; aliases
are transparent and enums retain nominal identity. Both operands are evaluated
**once, left before right**, at the expression's evaluation site. Short-circuit
`and`/`or` still skip unselected expressions. `!=` is the logical complement of
`==`, including when a component is NaN.

Eligible native aggregates compare their components, not their object bytes.
Arrays compare elements, tuples compare positions, and structs compare declared
fields. Empty records and eligible zero-length arrays compare equal. Padding,
backend dummy fields, and inactive enum union storage never participate.
Payload enums first compare the variant, then only that variant's active fields.
Every possible component must support equality, even in an empty array or an
unselected enum alternative. Diagnostics identify the unsupported field,
tuple position, array element type, or variant payload path.

For eligible payload enums, an undeclared tag has no active payload: equal tags
compare equal without reading any payload storage, and distinct tags compare
unequal. This matches payload-free enum tag equality. In particular,
`zero<Enum>()` still zero-initializes its representation; it does **not** select
the first declared variant. If all variant tags are nonzero, two such zero
values compare equal by their undeclared zero tag. Automatic hashing hashes only
the tag for an undeclared tag, so equal values have equal hashes.

The compiler also classifies concrete types for automatic **key eligibility**.
The compiler provides internal seeded semantic hashing for these concrete types;
this is not yet a map/set API or public hashing operation. Hashes are not stable
serialization identifiers. Equality and key eligibility are distinct:

| Concrete type family | Ordinary equality | Key eligible |
| --- | --- | --- |
| `bool`, signed/unsigned integers (including 128-bit), C-native integer scalars | Value | Yes |
| `float`, `f32`, `f64`, C-native float scalars | Numeric: signed zeros equal; NaN unequal to everything, including itself; infinities compare normally | Yes |
| `string` | Byte length and contents, including embedded zeros and malformed UTF-8; all empty views equal | Yes |
| `c::CString`, `c::ConstCString` | Pointer identity, **not** zero-terminated contents; null allowed | Yes |
| Typed object pointers, `rawptr` | Pointer identity; null allowed; pointee contents ignored | Yes |
| Fixed native arrays, tuples, native structs | Componentwise, if every component supports equality | If every component is key eligible, subject to owner exclusions below |
| Native and foreign payload-free enums | Variant/tag equality within the same nominal enum type | Yes |
| Native payload enums, including concrete `: error` enums | Same variant and equal active fields, if **all** alternatives support equality | If **all** alternatives' fields are key eligible |
| Aliases, concrete generic applications | Resolved concrete type's rule | Resolved concrete type's rule |
| Native function values and foreign callbacks | Function identity (including null), with compatible signature and convention | No, including when nested |
| Slices | No whole-value equality | No |
| Header-owned foreign records/unions and inline foreign arrays | No whole-value equality; compare supported fields/indexed elements explicitly | No |
| Opaque foreign values | No by-value equality/storage; pointers use pointer identity | No by-value keys |
| General `error` | Presence checks against `nil` only; no whole-value equality | No |
| Inferred error identities (`err.kind`) | Classified identity only, ignoring message/payload; unclassified and success kinds are null | Yes |
| Borrowed `c::c_va_list`, owning `c::VaList` | No whole-value equality; existing lifecycle restrictions remain | No |
| `void`, unresolved generic declarations, bare `nil` | Not concrete value types for aggregate/key classification | No |

Pointer traversal stops at the pointer, so `*T` does not require `T` to support
equality or keys. For example, a pointer to a slice-containing record remains
eligible by identity. No rule dereferences pointers or deep-clones resources.

### Owning handles

The following **standard-library declarations**, including aliases and concrete
generic instances, are excluded as keys even when their fields are eligible:

- `std::strings::Owned`, `std::strings::COwned` ([text ownership](../stdlib/strings-and-memory.md#individually-owned-text)).
- `std::memory::Arena` ([arena lifetime](../stdlib/memory.md)).
- `std::lists::List<T>` ([list ownership](../stdlib/lists.md)).
- `std::fs::File` ([file cleanup](../stdlib/files.md)).
- `std::process::Args`, `std::process::Child` ([process ownership](../stdlib/process.md)).
- `std::terminal::Handle`, `std::terminal::SavedMode` ([terminal ownership](../stdlib/terminal.md#handles-and-ownership)).

The nominal exclusion also reserves `std::maps::Map` and `std::sets::Set` for
their future implementations; neither container is provided by this slice.
It propagates through aggregate components, not through pointer pointees.
The exclusion also applies to zero values and borrowed instances of these
declarations, such as a borrowed terminal `Handle`.
Unrelated user records with the same names are not excluded. This is not a new
ownership annotation or inferred ownership system.

Owners retain ordinary fieldwise equality **when their constituent fields support
it**. Thus string owners compare their string fields by content; arena/file/
process wrappers compare their pointer fields by identity; terminal wrappers
compare all declared bookkeeping fields, not just the native handle. Lists
currently contain a slice and therefore do not support whole-value equality.
Comparison never transfers ownership or changes cleanup obligations. Copying an
owner for comparison is still shallow; see [memory rules](memory.md).

Runnable example (no containers or custom comparison declarations required):

```telia
package main

use "std::io"

struct Pair<T> {
    value T
    label string
}

enum Result<T> {
    Empty
    Value {
        pair Pair<T>
    }
}

fn main() i32 {
    first := Result<i32>::Value{pair: Pair<i32>{value: 7, label: "count"}}
    second := Result<i32>::Value{pair: Pair<i32>{value: 7, label: "count"}}
    if first == second {
        io::println("equal")
    }
    if first != Result<i32>::Empty {
        io::println("not empty")
    }
    return 0
}
```

Expected output:

```text
equal
not empty
```

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
