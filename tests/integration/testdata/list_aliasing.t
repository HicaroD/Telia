package main

use "std::lists"
use "std::memory"
use "std::libc"

fn exercise(values *lists::List<i32>) {
    initial i32[3] := .{1, 2, 3}
    lists::append_slice<i32>(values, initial[:]) @fail
    view := lists::view<i32>(values)
    middle := view[1:3] @fail
    lists::append_slice<i32>(values, middle) @fail
    result := lists::view<i32>(values)
    for value in result {
        libc::printf("%d\n", value)
    }
    return
}

fn main() i32 {
    heap := lists::create<i32>()
    defer lists::free<i32>(&heap)
    exercise(&heap)
    arena := memory::arena_create() @fail
    defer memory::arena_destroy(&arena)
    local := lists::create_in<i32>(&arena)
    defer lists::free<i32>(&local)
    exercise(&local)
    return 0
}
