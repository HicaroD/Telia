package main

use "std::lists"
use "std::memory"
use "std::io"
use "std::libc"

@[default_cc="c"]
extern {
    fn _test_fail_after(count i64)
    fn _test_live_allocations() i64
}

fn heap_failure() {
    values := lists::create<i32>()
    defer lists::free<i32>(&values)
    initial i32[2] := .{1, 2}
    lists::append_slice<i32>(&values, initial[:]) @fail
    saved := lists::view<i32>(&values)
    _test_fail_after(0)
    lists::reserve<i32>(&values, 1024) @catch err {
        io::println("reserve failed")
    }
    lists::append<i32>(&values, 3) @catch err {
        io::println("append failed")
    }
    copied, err := lists::clone<i32>(&values)
    if err != nil {
        if lists::count<i32>(&copied) != 0 {
            libc::exit(2)
        }
        io::println("clone failed")
    }
    lists::free<i32>(&copied)
    lists::reserve<i32>(&values, 18446744073709551615) @catch err {
        io::println("overflow rejected")
    }
    libc::printf("%llu\n%d\n", lists::count<i32>(&values), saved[1])
    _test_fail_after(-1)
    return
}

fn arena_failure() {
    arena := memory::arena_create() @fail
    defer memory::arena_destroy(&arena)
    values := lists::create_in<i64>(&arena)
    defer lists::free<i64>(&values)
    lists::append<i64>(&values, 7) @fail
    saved := lists::view<i64>(&values)
    _test_fail_after(0)
    lists::reserve<i64>(&values, 10000) @catch err {
        io::println("arena failed")
    }
    libc::printf("%lld\n", saved[0])
    lists::free<i64>(&values)
    memory::arena_reset(&arena)
    lists::append<i64>(&values, 9) @fail
    _test_fail_after(-1)
    return
}

fn main() i32 {
    heap_failure()
    arena_failure()
    if _test_live_allocations() != 0 {
        return 1
    }
    io::println("released")
    return 0
}
