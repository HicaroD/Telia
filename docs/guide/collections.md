# Arrays, slices, and lists

Prerequisite: [structs and enums](data-types.md).

An array stores a fixed number of elements inline. A slice borrows a region of
existing storage. A list owns growable backing storage. The following program
uses only an array and a slice, so it needs no allocation or cleanup.

```telia
package main

use "std::io"

fn main() i32 {
    values i32[3] := .{10, 20, 30}
    view := values[:]
    view[1] = 42
    for value in values {
        if value == 42 {
            io::println("changed through the slice")
        }
    }
    return 0
}
```

Expected output:

```text
changed through the slice
```

`i32[3]` has three integers. The contextual array literal `.{...}` needs an
expected type. `values[:]` creates a writable full slice without allocating.
Changing the slice's element changes the original array's storage.

Copying an array copies all elements. Copying a slice copies only its pointer
and length. A slice must not outlive its storage. Indexing is bounds-checked and
panics on an invalid index; `std::slices::get` is a recoverable alternative.
Indices use `u64` (or compatible integer literals), not arbitrary integer types.

Lists use `std::lists`. This fragment illustrates allocation; the next chapters
explain the cleanup and error operators:

```rust
items := lists::create<i32>()
defer lists::free<i32>(&items)
lists::append<i32>(&items, 42) @try
```

`<i32>` supplies a generic element type. `&items` passes its address. `defer`
arranges cleanup. `@try` returns an error to the caller if append fails. Use this
fragment inside a function returning `error`, with `std::lists` imported.

List growth can invalidate borrowed views; never structurally mutate a list
while iterating its view. Lists do not clean up resources owned by their elements.

**Try it:** copy the array before modifying `view`, and check the copy's second
element. It will still be 20.

Reference: [collection semantics](../reference/collections.md),
[lists](../stdlib/lists.md), [slices](../stdlib/slices.md).
Next: [memory and cleanup](memory.md).
