# Constants and bitwise expressions

## Named constants

Use `const Name [Type] = expression` at package scope or inside a block:

```rust
const DefaultWidth i32 = 800
const Answer = 40 + 2

type WindowFlags = u64
const WindowFullscreen WindowFlags = 8
const WindowResizable WindowFlags = 32

fn main() {
    const Options WindowFlags = WindowFullscreen | WindowResizable
    enabled := Options & WindowFullscreen != 0
}
```

Initializers are evaluated during semantic checking, not at runtime. Supported
scalar expressions use literals, other constants, arithmetic, comparisons,
logical operators, and the bitwise operators listed below. Runtime function
calls, variables, and aggregate initializers are not constant expressions.

An explicit type fixes the constant's type. Without one, scalar types are
inferred as for literal expressions; compatible expected types can specialize
individual uses without changing the declaration. Integer defaults follow the
existing `int` type. Values must fit their declared or contextual type. Integer
arithmetic is checked for overflow during evaluation, including intermediate
results; division by zero and cyclic dependencies are errors.

Constants cannot be reassigned and have no addressable storage. Copy a value
into a variable with `:=` if mutable or addressable storage is needed. Local
constants follow block scope. Generic function specializations have independent
local constant state.

Package constants use ordinary imported-name access:

```rust
use "pkg::settings"

const Width i32 = settings::DefaultWidth
```

Package constants may refer to later declarations. Names and imports in an
initializer resolve in its declaring file and scope.

## Bitwise operators

`&`, `|`, and `^` perform integer AND, OR, and XOR. Prefix `~` complements an
integer within its width. Operands must have the same integer type, except that
compatible literal/constant operands may adopt a contextual type. Boolean
operands are rejected; use `and`, `or`, and `not` for boolean logic.

Precedence, strongest to weakest:

1. Unary operators, including `~`.
2. `*`, `/`.
3. `+`, `-`.
4. `<<`, `>>`.
5. `&`.
6. `^`.
7. `|`.
8. Comparisons.
9. `and`.
10. `or`.

Thus `options & mask != 0` means `(options & mask) != 0`, unlike C's
precedence. Prefix `&value` remains address-of.

`<<` and `>>` require integer operands. The left operand determines the result
type and width; the count may have a different integer type. Counts must be
nonnegative and smaller than that width. Invalid constant counts are diagnosed;
invalid dynamic counts panic. Left shifts discard bits beyond the width,
including for signed integers. Signed right shifts replicate the sign bit:
`-8 >> 1` is `-4`. Unsigned right shifts fill with zeroes.

```rust
type WindowFlags = u64
const Fullscreen WindowFlags = 1 << 0
const Resizable WindowFlags = 1 << 5
const HighBit WindowFlags = 1 << 63
```

Shifts are left-associative; addition binds more tightly. Constant evaluation
and generated C use the same fixed-width rules. Runtime tests exercise signed
and unsigned boundaries, including 128-bit values, under UBSan on macOS arm64.

## Current delivery boundary

The constants/bitwise foundation and the first [enum slice](enums.md) of issue
#93 are implemented, along with [header-owned opaque pointers](foreign-opaque.md).
Further [C interoperability](foreign-interop.md) includes header-supplied
constants, records/enums, pointer qualifiers, callbacks, target-aware C scalars,
and native borrowed/owning `va_list` operations. Header constants remain C expressions,
not inputs to Telia's compile-time evaluator. Array lengths still require literal
lengths. Standard-library declarations now use pointer-level const qualification.
