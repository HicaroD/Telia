# Functions

```rust
fn add(left i32, right i32) i32 {
    return left + right
}

fn report(message string) {
    io::println(message)
    return
}

fn inspect(value i32) (i32, bool) {
    return value, value > 0
}
```

These declarations assume `std::io` is imported. Each parameter is written as
`name Type`; the result follows the closing parenthesis. No result annotation
means no value result. A tuple result returns multiple values in order.
Return values must match the declared result types; a no-result function cannot
return a value. Use explicit returns on all paths that need to supply a result.

Call functions with parentheses, including no-argument calls. Multiple results
may be received by a matching declaration/assignment or returned from a
compatible function. Error-ending results support [error operators](errors.md).

Parameters use value-copy semantics, including shallow copying for pointer-like
values and handles. Pass `*T` when the function must change the caller's `T`.
Passing an owner by value neither transfers its cleanup obligation nor clones
its allocation. There is no automatic destruction at return.

## Function values

```rust
callback fn(left i32, right i32) i32 := add
sum := callback(20, 22)
```

Concrete named functions can be stored, passed, and called through function
types. This is not a closure system: no anonymous capturing-function syntax is
provided. Generic functions require a concrete specialization, supplied
explicitly or by [inference](generics.md#inference-and-partial-arguments).

Calling conventions are part of the native boundary. Callback compatibility,
null callbacks, C aggregate restrictions, and header-owned prototypes follow
the [foreign callback rules](foreign-interop.md#callbacks-and-calling-conventions).
Do not replace object-pointer conversions with function-pointer casts.

## Variadic parameters

The parser recognizes typed variadic parameters as `name ...T`, but the C
function lowering does not pack them into a Telia collection. Use `T[]` for
variable-length application inputs rather than relying on that incomplete path.
Native C variadics use `name @c ...`. Native argument promotion,
`va_start`, `va_copy`, `va_arg`, borrowed `va_list`, and explicit `va_end`
have strict lifecycle/ABI rules; see
[variadics](foreign-interop.md#variadics-and-native-va_list-lifecycle).

## Tests

A named `test "description" { ... }` is a separate test declaration, not a
function with a special name. Only `*_test.t` files may declare tests.
See [unit testing](../getting-started/testing.md).
