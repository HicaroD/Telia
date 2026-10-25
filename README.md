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

## Strings and memory

`string` is a read-only borrowed byte view with an explicit length. Text is
UTF-8 by convention; embedded zero bytes are preserved. `==` and `!=` compare
bytes. Use `std::strings` for byte length, checked slices, owning copies, and
explicit C-string conversions.

```rust
text := strings::clone("Hello") @fail
defer strings::free(&text)
part := strings::slice(strings::view(&text), 0, 5) @fail
io::println(part)
```

Heap owners require explicit cleanup. `std::memory` also provides arenas for
allocations that share a lifetime. Copies are shallow: copying an owner or an
arena does not duplicate storage or create an independent owner.

See the [strings and memory reference](docs/strings-and-memory.md) and the
runnable [strings example](examples/strings.t).
