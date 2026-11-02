package main

use "std::maps"
use "std::errors"
use "std::testing"
use "std::io"

@[default_cc="c"]
extern {
    fn _test_fail_after(count i64)
    fn _test_live_allocations() i64
}

fn propagated() (i32, bool, error) {
    map := maps::create<i32, i32>()
    defer maps::free(&map)
    maps::reserve(&map, 7) @try
    key i32 := 0
    while key < 7 {
        maps::insert(&map, key, key) @try
        key = key + 1
    }
    _test_fail_after(1)
    previous, replaced := maps::insert(&map, 8, 80) @try
    return previous, replaced, nil
}

fn caught_failure() error {
    map := maps::create<i32, i32>()
    defer maps::free(&map)
    _test_fail_after(0)
    old, updated := maps::insert(&map, 1, 1) @catch err {
        return err
    }
    return error("unreachable catch success")
}

fn checked() error {
    map := maps::create<i32, i32>()
    defer maps::free(&map)
    initial, initially_replaced := maps::insert(&map, 1, 10) @fail
    testing::assert(initial == 0 and not initially_replaced)
    caught_old, caught_replaced := maps::insert(&map, 1, 20) @catch err {
        return err
    }
    testing::assert(caught_old == 10 and caught_replaced)
    old, updated := maps::insert(&map, 1, 30) @catch err {
        return err
    }
    testing::assert(old == 20 and updated)
    previous, replaced, err := propagated()
    testing::assert(err.kind == errors::Error::AllocationFailed)
    testing::assert(previous == 0 and not replaced)
    // @try unwound the nested map's committed buffers and its partial growth.
    testing::assert(_test_live_allocations() == 2)
    caught_error := caught_failure()
    testing::assert(caught_error.kind == errors::Error::AllocationFailed)
    testing::assert(_test_live_allocations() == 2)
    _test_fail_after(0)
    ignored, changed := maps::insert(&map, 2, 2) @catch allocation_error {
        return allocation_error
    }
    // This insert fits: @catch success must extract the two non-error fields.
    testing::assert(ignored == 0 and not changed)
    return nil
}

fn main() i32 {
    checked() @fail
    testing::assert(_test_live_allocations() == 0)
    io::println("checked ok")
    return 0
}
