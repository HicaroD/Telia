package main

use "std::libc"
use "std::io"

fn main() i32 {
    values i32[3] := .{10, 20, 30}
    view := values[:]
    part := view[1:3] @fail
    part[0] = 77
    libc::printf("%llu\n", len(view))
    libc::printf("%llu\n", len(part))
    libc::printf("%d\n", values[1])
    end := view[3:3] @fail
    libc::printf("%llu\n", len(end))
    bad, err := view[2:1]
    if err != nil {
        io::println("bad range")
    }
    return 0
}
