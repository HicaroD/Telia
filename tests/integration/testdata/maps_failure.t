package main

use "std::maps"
use "std::errors"
use "std::testing"
use "std::io"

@[default_cc="c"]
extern {
    fn _test_fail_after(count i64)
    fn _test_live_allocations() i64
    fn _map_seed_failure(fail bool)
    fn _map_seed_calls() i64
}

fn main() i32 {
    map := maps::create<u64, u64>()
    defer maps::free(&map)
    _map_seed_failure(true)
    _test_fail_after(0)
    v, found := maps::get(&map, 1)
    testing::assert(not found and v == 0)
    testing::assert(not maps::contains(&map, 1))
    k, removed, existed := maps::remove(&map, 1)
    testing::assert(not existed and k == 0 and removed == 0)
    cursor := maps::iter(&map)
    ik, iv, more := maps::next(&cursor)
    testing::assert(not more)
    maps::clear(&map)
    maps::reserve(&map, 0) @fail
    testing::assert(_map_seed_calls() == 0 and _test_live_allocations() == 0)
    empty_overflow := maps::reserve(&map, 18446744073709551615)
    testing::assert(empty_overflow.kind == errors::Error::SizeOverflow)
    testing::assert(_map_seed_calls() == 0)
    old, replaced, seed_error := maps::insert(&map, 1, 2)
    testing::assert(seed_error.kind == errors::Error::RandomnessFailed)
    testing::assert(old == 0 and not replaced and maps::capacity(&map) == 0)
    added, seed_error2 := maps::insert_if_absent(&map, 1, 2)
    testing::assert(not added and seed_error2.kind == errors::Error::RandomnessFailed)
    reserve_seed_error := maps::reserve(&map, 1)
    testing::assert(reserve_seed_error.kind == errors::Error::RandomnessFailed)
    _map_seed_failure(false)
    // First and second-buffer failures both leave a seedless, empty handle.
    attempt i64 := 0
    while attempt < 2 {
        _test_fail_after(attempt)
        allocation_error := maps::reserve(&map, 7)
        testing::assert(allocation_error.kind == errors::Error::AllocationFailed)
        testing::assert(maps::capacity(&map) == 0 and maps::count(&map) == 0)
        testing::assert(_test_live_allocations() == 0)
        attempt = attempt + 1
    }
    _test_fail_after(-1)
    key u64 := 0
    while key < 7 {
        maps::insert(&map, key, key + 10) @fail
        key = key + 1
    }
    pointer := maps::get_ptr(&map, 3)
    kept_capacity := maps::capacity(&map)
    kept_seeds := _map_seed_calls()
    kept_allocations := _test_live_allocations()
    _map_seed_failure(true)
    seed_copy, clone_seed_error := maps::clone(&map)
    testing::assert(clone_seed_error.kind == errors::Error::RandomnessFailed)
    testing::assert(maps::count(&seed_copy) == 0 and maps::capacity(&seed_copy) == 0)
    maps::free(&seed_copy)
    testing::assert(maps::get_ptr(&map, 3) == pointer and *pointer == 13)
    _map_seed_failure(false)
    // Both independently allocated clone buffers must be cleaned on failure.
    clone_attempt i64 := 0
    while clone_attempt < 2 {
        _test_fail_after(clone_attempt)
        failed_copy, clone_error := maps::clone(&map)
        testing::assert(clone_error.kind == errors::Error::AllocationFailed)
        testing::assert(maps::count(&failed_copy) == 0 and maps::capacity(&failed_copy) == 0)
        maps::free(&failed_copy)
        maps::free(&failed_copy)
        testing::assert(_test_live_allocations() == kept_allocations)
        testing::assert(maps::get_ptr(&map, 3) == pointer and *pointer == 13)
        testing::assert(maps::count(&map) == 7 and maps::capacity(&map) == kept_capacity)
        clone_attempt = clone_attempt + 1
    }
    kept_seeds = _map_seed_calls()
    attempt = 0
    while attempt < 2 {
        _test_fail_after(attempt)
        reserve_error := maps::reserve(&map, 100)
        testing::assert(reserve_error.kind == errors::Error::AllocationFailed)
        testing::assert(*pointer == 13 and maps::get_ptr(&map, 3) == pointer)
        testing::assert(maps::capacity(&map) == kept_capacity and maps::count(&map) == 7)
        testing::assert(_test_live_allocations() == kept_allocations)
        _test_fail_after(attempt)
        unused, changed, insert_error := maps::insert(&map, 8, 80)
        testing::assert(insert_error.kind == errors::Error::AllocationFailed)
        testing::assert(not changed and unused == 0 and not maps::contains(&map, 8))
        testing::assert(maps::get_ptr(&map, 3) == pointer and *pointer == 13)
        testing::assert(maps::capacity(&map) == kept_capacity and maps::count(&map) == 7)
        testing::assert(_test_live_allocations() == kept_allocations)
        _test_fail_after(attempt)
        inserted, absent_error := maps::insert_if_absent(&map, 8, 80)
        testing::assert(not inserted and absent_error.kind == errors::Error::AllocationFailed)
        testing::assert(_test_live_allocations() == kept_allocations)
        attempt = attempt + 1
    }
    _test_fail_after(0)
    // Existing-key updates, no-op reserve, and iteration cannot allocate/seed.
    previous, updated := maps::insert(&map, 3, 30) @fail
    testing::assert(previous == 13 and updated and *pointer == 30)
    duplicate := maps::insert_if_absent(&map, 3, 300) @fail
    testing::assert(not duplicate)
    maps::reserve(&map, kept_capacity) @fail
    cursor = maps::iter(&map)
    visits u64 := 0
    while true {
        entry_key, entry_value, has := maps::next(&cursor)
        if not has {
            break
        }
        visits = visits + 1
    }
    testing::assert(visits == 7 and _map_seed_calls() == kept_seeds)
    overflow := maps::reserve(&map, 18446744073709551615)
    testing::assert(overflow.kind == errors::Error::SizeOverflow)
    byte_overflow := maps::reserve(&map, 9223372036854775807)
    testing::assert(byte_overflow.kind == errors::Error::SizeOverflow)
    testing::assert(maps::get_ptr(&map, 3) == pointer and *pointer == 30)
    testing::assert(maps::capacity(&map) == kept_capacity)
    key = 0
    while key < 7 {
        actual, has := maps::get(&map, key)
        testing::assert(has)
        if key != 3 {
            testing::assert(actual == key + 10)
        }
        key = key + 1
    }
    maps::free(&map)
    maps::free(&map)
    testing::assert(_test_live_allocations() == 0)
    _test_fail_after(-1)
    maps::insert(&map, 1, 1) @fail
    // Force same-size tombstone cleanup with one outstanding live borrow.
    key = 2
    while key < 8 {
        maps::insert(&map, key, key) @fail
        key = key + 1
    }
    key = 2
    while key < 8 {
        dead_key, dead_value, dead_found := maps::remove(&map, key)
        testing::assert(dead_found)
        key = key + 1
    }
    survivor := maps::get_ptr(&map, 1)
    attempt = 0
    while attempt < 2 {
        _test_fail_after(attempt)
        tombstone_old, tombstone_replaced, tombstone_error := maps::insert(&map, 100, 100)
        testing::assert(tombstone_error.kind == errors::Error::AllocationFailed)
        testing::assert(not tombstone_replaced and tombstone_old == 0)
        testing::assert(maps::capacity(&map) == 7 and maps::count(&map) == 1)
        testing::assert(maps::get_ptr(&map, 1) == survivor and *survivor == 1)
        testing::assert(_test_live_allocations() == 2)
        attempt = attempt + 1
    }
    _test_fail_after(-1)
    maps::insert(&map, 100, 100) @fail
    testing::assert(maps::capacity(&map) == 7 and maps::count(&map) == 2)
    reused_seed := _map_seed_calls()
    maps::clear(&map)
    _map_seed_failure(true)
    _test_fail_after(0)
    maps::insert(&map, 2, 2) @fail
    testing::assert(_map_seed_calls() == reused_seed)
    io::println("failure ok")
    return 0
}
