package main

use "std::lists"
use "std::libc"

fn main() i32 {
    values := lists::create<i32>()
    defer lists::free<i32>(&values)
    initial i32[4] := .{10, 20, 30, 40}
    lists::append_slice<i32>(&values, initial[:]) @fail
    copied := lists::clone<i32>(&values) @fail
    defer lists::free<i32>(&copied)
    ordered := lists::remove<i32>(&values, 1) @fail
    unordered := lists::remove_unordered<i32>(&values, 0) @fail
    last, found := lists::pop<i32>(&values)
    libc::printf("%d\n%d\n%d\n%d\n", ordered, unordered, last, found)
    remaining := lists::get<i32>(&values, 0) @fail
    libc::printf("%d\n", remaining)
    before := lists::capacity<i32>(&values)
    lists::clear<i32>(&values)
    missing, exists := lists::pop<i32>(&values)
    libc::printf("%d\n%d\n%d\n", missing, exists, before == lists::capacity<i32>(&values))
    view := lists::view<i32>(&copied)
    for value in view {
        libc::printf("%d\n", value)
    }
    lists::free<i32>(&values)
    return 0
}
