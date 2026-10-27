package main

use "std::libc"

fn main() i32 {
    values i32[3] := .{10, 20, 30}
    copied := values
    values[1] = 99
    empty i32[3] := .{}
    libc::printf("%llu\n", len(values))
    libc::printf("%d\n", copied[1])
    libc::printf("%d\n", empty[2])
    libc::printf("%d\n", values[1])
    return 0
}
