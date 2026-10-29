# Enums and matching

See [literal syntax](literals.md) for construction syntax and control-flow
header rules.

Enums have nominal identity and an optional integer tag representation:

```rust
const FirstTag u8 = 7
enum Event u8 {
    Quit = FirstTag
    Resize
}

fn identity(event Event) Event {
    return event
}

event := identity(Event::Quit)
different := event != Event::Resize
```

Without a representation, tags use `int`. The first implicit tag is zero;
each later implicit tag is the previous tag plus one. Explicit tags use the
existing scalar constant-expression evaluator. Every tag must fit the integer
representation; implicit increments cannot wrap. Names and tags must be unique.
Representation aliases are resolved in the enum's declaring file.

Construction uses `Enum::Variant`, including `package::Enum::Variant` for
imports. Values can be assigned, passed, returned, and compared with `==` and
`!=` within the same enum type. Enums are distinct from integers and other
enums, even when their representations and tags coincide. Flags remain integer
aliases plus constants, not enum arithmetic.

Generic declarations and named payloads use the existing `<T>` type arguments:

```rust
enum Option<T> {
    None
    Some {
        value T
    }
}
item := Option<i32>::Some{value: 42}

value := match item {
    Option<i32>::None => 0
    Option<i32>::Some{value} => {
        value = value + 1
        yield value
    }
}
```

Payload construction initializes every field exactly once by name. Field
expressions are evaluated once in source order. Positional payloads are not
supported. Generic specializations keep independent payload types and layouts;
construction and matching also work inside generic functions.

Match patterns use named fields (`{value}`), optional renaming
(`{value: result}`), and explicit `..` for omitted fields. Arms use `=>`;
block results use `yield`. Every variant must be covered, or the final arm must
be a `_` fallback. Duplicate variants, unreachable arms, and omitted fields
without `..` are rejected. Every completing block path must yield a value.
All arm results must have compatible types.

The matched expression is evaluated once, and only the selected arm executes.
Matching remains at its expression evaluation site, including short-circuit
operators and loop conditions. Nested matches have independent yield targets.
Yield saves the result before running the arm's deferred cleanup; it does not
run cleanup for enclosing scopes that remain active.

**Bindings copy fields using ordinary assignment semantics.** Reassigning a
binding does not change the original enum field. Copying a pointer still refers
to the same pointee; copying does not deep-clone the pointee or transfer its
ownership. There is no implicit borrow. `yield` cannot escape deferred cleanup
or a called function into the enclosing match.

Payload-free enums support `==` and `!=`. Payload enums do not have implicit
equality: inspect them with `match` and explicitly compare relevant fields.

The current C backend uses the selected integer representation for payload-free
values and private tagged unions for payload enums. Match expressions currently
use GCC/Clang statement expressions to preserve their evaluation site. These
are internal lowering details, **not a foreign C enum ABI guarantee**.
Header-owned foreign enums and their required fallback matches are described in
[C interoperability](foreign-interop.md), along with records/constants,
qualifiers, callbacks, convention macros, C scalars, and native `va_list` lifecycle.
Integration tests compile and execute generated C on the local macOS arm64
host; this does not establish compatibility on other targets.
