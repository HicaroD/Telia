# The Telia Programming Language

Telia is a general-purpose programming language. Heavily inspired on the Go
Programming Language, Telia is simple in its essence.

```rust
package main

use "std::io"

fn main() i32 {
    io::println("Hello, world!")
    return 0
}
```

> [!WARNING]
> Telia compiler is heavily imature and not ready to use.

## Compiler diagnostics

`telia build` reports errors to stderr with source highlights, related locations,
and final error/warning counts. Independent errors are collected together; failed
compilations exit nonzero and do not run code generation.

Use `--color=auto|always|never` (default: `auto`). Automatic color is disabled for
redirected stderr and when `NO_COLOR` is set. Explicit `always` or `never` overrides
`NO_COLOR`.
