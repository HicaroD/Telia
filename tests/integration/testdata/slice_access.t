package main

use "std::slices"
use "std::io"
use "std::libc"

fn main() i32 {
    values i32[1] := .{}
    view := values[:]
    slices::set<i32>(view, 0, 42) @fail
    value := slices::get<i32>(view, 0) @fail
    libc::printf("%d\n", value)
    missing, err := slices::get<i32>(view, 1)
    if err != nil {
        io::println("missing")
    }
    slices::set<i32>(view, 1, 99) @catch problem {
        io::println("unchanged")
    }
    libc::printf("%d\n", values[0])
    return 0
}
