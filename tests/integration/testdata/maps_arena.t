package main

use "std::maps"
use "std::memory"
use "std::testing"
use "std::errors"

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
    map := maps::create_in<u64, Aligned>(&arena)
    defer maps::free(&map)
    key u64 := 0
    while key < 1000 {
        maps::insert(&map, key, Aligned{wide: 17, number: key}) @fail
        pointer := maps::get_ptr(&map, key)
        testing::assert(_test_is_aligned(pointer, alignof<Aligned>()))
        testing::assert(pointer.wide == 17 and pointer.number == key)
        testing::assert(bytes[0] == 91)
        key = key + 1
    }
    kept := maps::capacity(&map)
    maps::clear(&map)
    _test_fail_after(0)
    maps::insert(&map, 1, Aligned{wide: 2, number: 3}) @fail
    testing::assert(maps::capacity(&map) == kept)
    _test_fail_after(-1)
    allocations := _test_live_allocations()
    maps::free(&map)
    maps::free(&map)
    testing::assert(_test_live_allocations() == allocations)
    testing::assert(maps::count(&map) == 0 and maps::capacity(&map) == 0)
    // Free the handle before resetting the arena; never read invalidated entries.
    memory::arena_reset(&arena)
    _test_fail_after(0)
    // Reuse proves free kept the arena association rather than switching to heap.
    maps::insert(&map, 4, Aligned{wide: 5, number: 6}) @fail
    moved := memory::take(&map)
    maps::free(&map)
    testing::assert(maps::count(&moved) == 1)
    maps::free(&moved)
    memory::arena_destroy(&arena)
    memory::arena_destroy(&arena)
    testing::assert(_test_live_allocations() == 0)
    // The moved source is a usable heap-mode zero, not an arena-associated handle.
    old, replaced, err := maps::insert(&map, 1, Aligned{})
    testing::assert(err.kind == errors::Error::AllocationFailed and not replaced)
    _test_fail_after(-1)
    maps::insert(&map, 1, Aligned{}) @fail
    maps::free(&map)
    return
}

fn failures() {
    attempt i64 := 0
    while attempt < 2 {
        arena := memory::arena_create() @fail
        map := maps::create_in<u64, u64>(&arena)
        maps::insert(&map, 1, 10) @fail
        pointer := maps::get_ptr(&map, 1)
        kept := maps::capacity(&map)
        baseline := _test_live_allocations()
        _test_fail_after(attempt)
        // Both new buffers exceed arena block slack: fail either allocation.
        err := maps::reserve(&map, 10000)
        testing::assert(err.kind == errors::Error::AllocationFailed)
        testing::assert(maps::count(&map) == 1 and maps::capacity(&map) == kept)
        testing::assert(maps::get_ptr(&map, 1) == pointer and *pointer == 10)
        testing::assert(_test_live_allocations() == baseline + attempt)
        overflow := maps::reserve(&map, 18446744073709551615)
        byte_overflow := maps::reserve(&map, 9223372036854775807)
        testing::assert(overflow.kind == errors::Error::SizeOverflow)
        testing::assert(byte_overflow.kind == errors::Error::SizeOverflow)
        testing::assert(maps::get_ptr(&map, 1) == pointer and *pointer == 10)
        maps::free(&map)
        memory::arena_destroy(&arena)
        testing::assert(_test_live_allocations() == 0)
        _test_fail_after(-1)
        attempt = attempt + 1
    }
    // Large values force entry-buffer allocation on the eighth insertion.
    arena := memory::arena_create() @fail
    map := maps::create_in<u64, u8[8192]>(&arena)
    key u64 := 0
    while key < 7 {
        maps::insert(&map, key, zero<u8[8192]>()) @fail
        key = key + 1
    }
    pointer := maps::get_ptr(&map, 3)
    _test_fail_after(0)
    old, replaced, err := maps::insert(&map, 8, zero<u8[8192]>())
    testing::assert(err.kind == errors::Error::AllocationFailed and not replaced)
    testing::assert(maps::capacity(&map) == 7 and maps::count(&map) == 7)
    testing::assert(maps::get_ptr(&map, 3) == pointer and not maps::contains(&map, 8))
    added, absent_error := maps::insert_if_absent(&map, 8, zero<u8[8192]>())
    testing::assert(not added and absent_error.kind == errors::Error::AllocationFailed)
    testing::assert(maps::get_ptr(&map, 3) == pointer)
    maps::free(&map)
    memory::arena_destroy(&arena)
    _test_fail_after(-1)
    return
}

fn clone_failures() {
    source := maps::create<u64, u64>()
    key u64 := 0
    while key < 5000 {
        maps::insert(&source, key, key + 1) @fail
        key = key + 1
    }
    pointer := maps::get_ptr(&source, 3)
    kept := maps::capacity(&source)
    attempt i64 := 0
    while attempt < 2 {
        arena := memory::arena_create() @fail
        baseline := _test_live_allocations()
        _test_fail_after(attempt)
        copy, err := maps::clone_in(&source, &arena)
        testing::assert(err.kind == errors::Error::AllocationFailed)
        testing::assert(maps::count(&copy) == 0 and maps::capacity(&copy) == 0)
        testing::assert(maps::get_ptr(&source, 3) == pointer and *pointer == 4)
        testing::assert(maps::count(&source) == 5000 and maps::capacity(&source) == kept)
        testing::assert(_test_live_allocations() == baseline + attempt)
        maps::free(&copy)
        maps::free(&copy)
        _test_fail_after(-1)
        prime := memory::arena_alloc(&arena, 4096, 8) @fail
        memory::arena_reset(&arena)
        // Failed clone retains the requested allocator, even through free.
        _test_fail_after(0)
        maps::insert(&copy, 9, 10) @fail
        maps::free(&copy)
        memory::arena_destroy(&arena)
        testing::assert(_test_live_allocations() == baseline)
        _test_fail_after(-1)
        attempt = attempt + 1
    }
    maps::free(&source)
    return
}

fn main() i32 {
    lifecycle()
    failures()
    clone_failures()
    heap := maps::create_in<u64, u64>(nil)
    maps::insert(&heap, 1, 2) @fail
    copy := maps::clone_in(&heap, nil) @fail
    maps::free(&heap)
    testing::assert(*maps::get_ptr(&copy, 1) == 2)
    maps::free(&copy)
    testing::assert(_test_live_allocations() == 0)
    return 0
}
