package main

use "std::lists"
use "std::memory"
use "std::libc"
use "std::io"

struct Record {
    name string
    samples i32[3]
}

fn run() error {
    samples i32[3] := .{10, 20, 30}
    borrowed := samples[:]
    middle := borrowed[1:3] @catch err {
        return err
    }
    middle[0] = 42

    records := lists::create<Record>()
    defer lists::free<Record>(&records)
    record := Record{name: "samples", samples: samples}
    lists::append<Record>(&records, record) @catch err {
        return err
    }
    view := lists::view<Record>(&records)
    for item in view {
        io::println(item.name)
        for sample in item.samples {
            libc::printf("%d\n", sample)
        }
    }

    arena := memory::arena_create() @catch err {
        return err
    }
    defer memory::arena_destroy(&arena)
    numbers := lists::create_in<i32>(&arena)
    defer lists::free<i32>(&numbers)
    lists::append_slice<i32>(&numbers, middle) @catch err {
        return err
    }
    libc::printf("%llu\n", lists::count<i32>(&numbers))
    return nil
}

fn main() i32 {
    run() @fail
    return 0
}
