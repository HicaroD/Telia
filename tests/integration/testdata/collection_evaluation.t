package main

use "std::io"
use "std::libc"

fn source(values i32[]) i32[] {
    io::println("source")
    return values
}

fn pick_index() u64 {
    io::println("index")
    return 0
}

fn main() i32 {
    values i32[1] := .{42}
    view := values[:]
    libc::printf("%d\n", source(view)[pick_index()])
    part := source(view)[pick_index():1] @fail
    libc::printf("%d\n", part[0])
    return 0
}
