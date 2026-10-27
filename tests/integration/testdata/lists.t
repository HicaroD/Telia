package main

use "std::lists"
use "std::libc"

fn main() i32 {
    values := lists::create<i32>()
    defer lists::free<i32>(&values)
    lists::append<i32>(&values, 10) @fail
    lists::append<i32>(&values, 20) @fail
    lists::append<i32>(&values, 30) @fail
    view := lists::view<i32>(&values)
    view[1] = 77
    libc::printf("%llu\n", lists::count<i32>(&values))
    for value in view {
        libc::printf("%d\n", value)
    }
    return 0
}
