package main

use "std::sets"
use "std::testing"
use "std::errors"

@[default_cc="c"]
extern {
    fn _test_fail_after(count i64)
    fn _test_live_allocations() i64
    fn _map_seed_failure(fail bool)
    fn _map_seed_calls() i64
}

fn verify(set *sets::Set<u64>, expected u64, capacity u64) {
    testing::assert(sets::count(set) == expected and sets::capacity(set) == capacity)
    key u64 := 0
    while key < expected {
        testing::assert(sets::contains(set, key))
        key = key + 1
    }
    return
}

fn main() i32 {
    set := sets::create<u64>()
    _map_seed_failure(true)
    _test_fail_after(0)
    testing::assert(not sets::contains(&set, 1))
    missing, found := sets::remove(&set, 1)
    testing::assert(not found and missing == 0)
    cursor := sets::iter(&set)
    element, more := sets::next(&cursor)
    testing::assert(not more and element == 0)
    sets::clear(&set)
    sets::free(&set)
    sets::reserve(&set, 0) @fail
    empty := sets::clone(&set) @fail
    sets::free(&empty)
    testing::assert(_map_seed_calls() == 0 and _test_live_allocations() == 0)
    overflow := sets::reserve(&set, 18446744073709551615)
    testing::assert(overflow.kind == errors::Error::SizeOverflow and _map_seed_calls() == 0)
    added, seed_error := sets::insert(&set, 1)
    testing::assert(not added and seed_error.kind == errors::Error::RandomnessFailed)
    seed_reserve := sets::reserve(&set, 7)
    testing::assert(seed_reserve.kind == errors::Error::RandomnessFailed)
    verify(&set, 0, 0)
    _map_seed_failure(false)
    attempt i64 := 0
    while attempt < 2 {
        _test_fail_after(attempt)
        inserted, err := sets::insert(&set, 1)
        testing::assert(not inserted and err.kind == errors::Error::AllocationFailed)
        verify(&set, 0, 0)
        testing::assert(_test_live_allocations() == 0)
        attempt = attempt + 1
    }
    _test_fail_after(-1)
    key u64 := 0
    while key < 7 {
        sets::insert(&set, key) @fail
        key = key + 1
    }
    kept_capacity := sets::capacity(&set)
    baseline := _test_live_allocations()
    _map_seed_failure(true)
    failed_seed_clone, clone_seed_error := sets::clone(&set)
    testing::assert(clone_seed_error.kind == errors::Error::RandomnessFailed)
    verify(&failed_seed_clone, 0, 0)
    sets::free(&failed_seed_clone)
    verify(&set, 7, kept_capacity)
    _map_seed_failure(false)
    seeds := _map_seed_calls()
    attempt = 0
    while attempt < 2 {
        _test_fail_after(attempt)
        // A live iterator remains valid after failed structural operations.
        preserved := sets::iter(&set)
        reserve_error := sets::reserve(&set, 100)
        testing::assert(reserve_error.kind == errors::Error::AllocationFailed)
        verify(&set, 7, kept_capacity)
        testing::assert(_test_live_allocations() == baseline)
        _test_fail_after(attempt)
        inserted, insert_error := sets::insert(&set, 8)
        testing::assert(not inserted and insert_error.kind == errors::Error::AllocationFailed)
        testing::assert(not sets::contains(&set, 8))
        verify(&set, 7, kept_capacity)
        testing::assert(_test_live_allocations() == baseline)
        visits u64 := 0
        while true {
            stored, has := sets::next(&preserved)
            if not has {
                break
            }
            testing::assert(stored < 7)
            visits = visits + 1
        }
        testing::assert(visits == 7 and _map_seed_calls() == seeds)
        _test_fail_after(attempt)
        failed, clone_error := sets::clone(&set)
        testing::assert(clone_error.kind == errors::Error::AllocationFailed)
        verify(&failed, 0, 0)
        sets::free(&failed)
        sets::free(&failed)
        verify(&set, 7, kept_capacity)
        testing::assert(_test_live_allocations() == baseline)
        seeds = _map_seed_calls()
        attempt = attempt + 1
    }
    _map_seed_failure(true)
    _test_fail_after(0)
    duplicate := sets::insert(&set, 3) @fail
    testing::assert(not duplicate)
    sets::reserve(&set, kept_capacity) @fail
    byte_overflow := sets::reserve(&set, 9223372036854775807)
    testing::assert(byte_overflow.kind == errors::Error::SizeOverflow)
    testing::assert(_map_seed_calls() == seeds)
    verify(&set, 7, kept_capacity)
    sets::free(&set)
    sets::free(&set)
    testing::assert(_test_live_allocations() == 0)
    _map_seed_failure(false)
    _test_fail_after(-1)
    key = 0
    while key < 7 {
        sets::insert(&set, key) @fail
        key = key + 1
    }
    key = 1
    while key < 7 {
        stored, removed := sets::remove(&set, key)
        testing::assert(removed and stored == key)
        key = key + 1
    }
    attempt = 0
    while attempt < 2 {
        _test_fail_after(attempt)
        inserted, err := sets::insert(&set, 100)
        testing::assert(not inserted and err.kind == errors::Error::AllocationFailed)
        testing::assert(sets::count(&set) == 1 and sets::capacity(&set) == 7)
        testing::assert(sets::contains(&set, 0) and not sets::contains(&set, 100))
        testing::assert(_test_live_allocations() == 2)
        attempt = attempt + 1
    }
    _test_fail_after(-1)
    sets::insert(&set, 100) @fail
    testing::assert(sets::count(&set) == 2 and sets::capacity(&set) == 7)
    sets::clear(&set)
    seeds = _map_seed_calls()
    _map_seed_failure(true)
    _test_fail_after(0)
    sets::insert(&set, 5) @fail
    testing::assert(_map_seed_calls() == seeds)
    // Empty, previously reserved clones still require neither allocation nor seed.
    sets::clear(&set)
    empty_reserved := sets::clone(&set) @fail
    verify(&empty_reserved, 0, 0)
    sets::free(&empty_reserved)
    sets::free(&set)
    testing::assert(_test_live_allocations() == 0)
    return 0
}
