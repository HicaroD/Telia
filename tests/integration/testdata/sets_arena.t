package main

use "std::sets"
use "std::memory"
use "std::testing"
use "std::errors"
use "std::internal::swiss"

@[default_cc="c"]
extern {
    fn _test_fail_after(count i64)
    fn _test_live_allocations() i64
    fn _test_is_aligned(ptr rawptr, alignment u64) bool
}

struct Aligned {
    wide u128
    number u64
}

fn lifecycle() {
    arena := memory::arena_create() @fail
    defer memory::arena_destroy(&arena)
    sentinel := memory::arena_alloc(&arena, 1, 1) @fail
    bytes := slice_from_raw<u8>(sentinel, 1)
    bytes[0] = 91
    set sets::Set<Aligned> := sets::create_in(&arena)
    defer sets::free(&set)
    key u64 := 0
    while key < 1000 {
        sets::insert(&set, Aligned{wide: 17, number: key}) @fail
        testing::assert(sets::contains(&set, Aligned{wide: 17, number: key}))
        testing::assert(bytes[0] == 91)
        key = key + 1
    }
    // Internal representation check complements public aligned-element behavior.
    testing::assert(_test_is_aligned(slice_data<swiss::Entry<Aligned, sets::Unit>>(set.table.entries), alignof<Aligned>()))
    heap := sets::clone(&set) @fail
    selected := sets::clone_in(&set, &arena) @fail
    sets::clear(&selected)
    testing::assert(sets::count(&set) == 1000 and sets::count(&heap) == 1000)
    sets::free(&selected)
    kept := sets::capacity(&set)
    sets::clear(&set)
    _test_fail_after(0)
    sets::insert(&set, Aligned{wide: 2, number: 3}) @fail
    testing::assert(sets::capacity(&set) == kept)
    _test_fail_after(-1)
    allocations := _test_live_allocations()
    sets::free(&set)
    sets::free(&set)
    testing::assert(_test_live_allocations() == allocations)
    memory::arena_reset(&arena)
    testing::assert(sets::contains(&heap, Aligned{wide: 17, number: 3}))
    sets::free(&heap)
    // Retained arena storage allows reuse even with heap allocation disabled.
    _test_fail_after(0)
    sets::insert(&set, Aligned{wide: 5, number: 6}) @fail
    moved := memory::take(&set)
    sets::free(&set)
    testing::assert(sets::count(&moved) == 1)
    sets::free(&moved)
    memory::arena_destroy(&arena)
    memory::arena_destroy(&arena)
    testing::assert(_test_live_allocations() == 0)
    inserted, err := sets::insert(&set, Aligned{})
    testing::assert(not inserted and err.kind == errors::Error::AllocationFailed)
    _test_fail_after(-1)
    sets::insert(&set, Aligned{}) @fail
    sets::free(&set)
    return
}

fn failures() {
    attempt i64 := 0
    while attempt < 2 {
        arena := memory::arena_create() @fail
        set := sets::create_in<u64>(&arena)
        sets::insert(&set, 1) @fail
        kept := sets::capacity(&set)
        baseline := _test_live_allocations()
        _test_fail_after(attempt)
        err := sets::reserve(&set, 10000)
        testing::assert(err.kind == errors::Error::AllocationFailed)
        testing::assert(sets::count(&set) == 1 and sets::capacity(&set) == kept)
        testing::assert(sets::contains(&set, 1))
        testing::assert(_test_live_allocations() == baseline + attempt)
        sets::free(&set)
        memory::arena_destroy(&arena)
        testing::assert(_test_live_allocations() == 0)
        _test_fail_after(-1)
        attempt = attempt + 1
    }
    // Large eligible array elements force fresh arena storage during growth.
    arena := memory::arena_create() @fail
    set := sets::create_in<u8[8192]>(&arena)
    key u8 := 0
    while key < 7 {
        element := zero<u8[8192]>()
        element[0] = key
        sets::insert(&set, element) @fail
        key = key + 1
    }
    absent := zero<u8[8192]>()
    absent[0] = 8
    _test_fail_after(0)
    added, err := sets::insert(&set, absent)
    testing::assert(not added and err.kind == errors::Error::AllocationFailed)
    testing::assert(sets::capacity(&set) == 7 and sets::count(&set) == 7)
    testing::assert(not sets::contains(&set, absent))
    sets::free(&set)
    memory::arena_destroy(&arena)
    _test_fail_after(-1)
    return
}

fn clone_failures() {
    source := sets::create<u64>()
    key u64 := 0
    while key < 5000 {
        sets::insert(&source, key) @fail
        key = key + 1
    }
    kept := sets::capacity(&source)
    attempt i64 := 0
    while attempt < 2 {
        arena := memory::arena_create() @fail
        baseline := _test_live_allocations()
        _test_fail_after(attempt)
        copy, err := sets::clone_in(&source, &arena)
        testing::assert(err.kind == errors::Error::AllocationFailed)
        testing::assert(sets::count(&copy) == 0 and sets::capacity(&copy) == 0)
        testing::assert(sets::contains(&source, 3))
        testing::assert(sets::count(&source) == 5000 and sets::capacity(&source) == kept)
        testing::assert(_test_live_allocations() == baseline + attempt)
        sets::free(&copy)
        sets::free(&copy)
        _test_fail_after(-1)
        prime := memory::arena_alloc(&arena, 4096, 8) @fail
        memory::arena_reset(&arena)
        _test_fail_after(0)
        sets::insert(&copy, 9) @fail
        sets::free(&copy)
        memory::arena_destroy(&arena)
        testing::assert(_test_live_allocations() == baseline)
        _test_fail_after(-1)
        attempt = attempt + 1
    }
    sets::free(&source)
    return
}

fn main() i32 {
    lifecycle()
    failures()
    clone_failures()
    heap sets::Set<u64> := sets::create_in(nil)
    sets::insert(&heap, 1) @fail
    copy := sets::clone_in(&heap, nil) @fail
    sets::free(&heap)
    testing::assert(sets::contains(&copy, 1))
    sets::free(&copy)
    testing::assert(_test_live_allocations() == 0)
    return 0
}
