# Language reference

Use the [guide](../guide/README.md) to learn concepts in order. This reference
describes implemented rules; it is not a design proposal. Library-specific
behavior belongs in the [standard library reference](../stdlib/README.md).

| Topic | Reference |
| --- | --- |
| Files, packages, imports, names | [Program structure](programs.md) |
| Variables, primitives, structs, aliases, copying | [Values and types](values.md) |
| Operators, conditions, loops, evaluation | [Expressions and statements](expressions.md) |
| Parameters, results, callbacks | [Functions](functions.md) |
| Constants and operator precedence | [Constants and bitwise expressions](constants-and-bitwise.md) |
| Record and array construction | [Literals](literals.md) |
| Arrays, slices, indexing, iteration | [Collections](collections.md) |
| Integer and pointer conversions | [Conversions](conversions.md) |
| Variants, payloads, match patterns, guards | [Enums and matching](enums.md) |
| Pointers, lifetimes, copying, defer | [Memory](memory.md) |
| Error values and propagation | [Errors](errors.md) |
| Type parameters, inference, specialization | [Generics](generics.md) |
| Headers, ABI, callbacks, C scalars, variadics | [C interoperability](foreign-interop.md) |
| Header-owned incomplete types | [Opaque pointers](foreign-opaque.md) |
| Build/test commands and configuration | [Compiler commands](commands.md) |

[grammar.ebnf](grammar.ebnf) is a compact syntax outline, not a complete parser
specification or a guarantee that every syntactically expressible type works in
every backend context. Prose contracts, compiler checks, and integration tests
describe the restrictions. Report disagreements rather than silently choosing a
new rule.

Unless identified as a complete program, examples in this section are fragments
requiring surrounding declarations and imports. Diagnostics and fatal examples
are not programs expected to run successfully.
