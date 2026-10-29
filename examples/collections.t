package main

use "std::lists"
use "std::memory"
use "std::libc"
use "std::io"

struct Record {
    name string
    samples i32[3]
}

fn main() i32 {
    samples i32[3] := .{10, 20, 30}
    borrowed := samples[:]
    middle := borrowed[1:3] @fail
    middle[0] = 42

    records := lists::create<Record>()
    defer lists::free<Record>(&records)
    record := Record{name: "samples", samples: samples}
    lists::append<Record>(&records, record) @fail
    view := lists::view<Record>(&records)
    for item in view {
        io::println(item.name)
        for sample in item.samples {
            libc::printf("%d\n", sample)
        }
    }

    arena := memory::arena_create() @fail
    defer memory::arena_destroy(&arena)
    numbers := lists::create_in<i32>(&arena)
    defer lists::free<i32>(&numbers)
    lists::append_slice<i32>(&numbers, middle) @fail
    libc::printf("%llu\n", lists::count<i32>(&numbers))
    return 0
}
