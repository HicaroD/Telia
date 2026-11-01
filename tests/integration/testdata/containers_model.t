package main

use "std::maps"
use "std::sets"
use "std::memory"
use "std::testing"
use "std::errors"

@[default_cc="c"]
extern {
    fn _map_nan() f64
    fn _qualification_float(value u64) f64
    fn _qualification_index(value f64) u64
    fn _test_live_allocations() i64
    fn _test_fail_after(count i64)
    fn _map_seed_failure(fail bool)
}

// Finite keys are modeled by indexed slots; NaNs are a separate multiset:
// every insert adds an unreachable entry, removed only by clear/free.
fn verify(map *maps::Map<f64, u64>, set *sets::Set<f64>, present bool[], values u64[], nans u64) {
    expected := nans
    key u64 := 0
    while key < 64 {
        value, found := maps::get(map, _qualification_float(key))
        testing::assert(found == present[key])
        testing::assert(sets::contains(set, _qualification_float(key)) == found)
        if found {
            testing::assert(value == values[key])
            expected = expected + 1
        }
        key = key + 1
    }
    testing::assert(maps::count(map) == expected and sets::count(set) == expected)
    map_seen bool[64] := .{}
    set_seen bool[64] := .{}
    map_nans u64 := 0
    set_nans u64 := 0
    for stored, value in *map {
        if stored != stored {
            map_nans = map_nans + 1
        } else {
            index := _qualification_index(stored)
            testing::assert(index < 64 and present[index] and not map_seen[index])
            testing::assert(value == values[index])
            map_seen[index] = true
        }
    }
    cursor := sets::iter(set)
    stored, found := sets::next(&cursor)
    while found {
        if stored != stored {
            set_nans = set_nans + 1
        } else {
            index := _qualification_index(stored)
            testing::assert(index < 64 and present[index] and not set_seen[index])
            set_seen[index] = true
        }
        stored, found = sets::next(&cursor)
    }
    testing::assert(map_nans == nans and set_nans == nans)
    key = 0
    while key < 64 {
        testing::assert(map_seen[key] == present[key] and set_seen[key] == present[key])
        key = key + 1
    }
    nan := _map_nan()
    testing::assert(not maps::contains(map, nan) and not sets::contains(set, nan))
    removed_key, removed_value, removed := maps::remove(map, nan)
    removed_element, removed_set := sets::remove(set, nan)
    testing::assert(not removed and not removed_set)
    return
}

fn run(arena *memory::Arena, seed u64) {
    map := maps::create_in<f64, u64>(arena)
    set := sets::create_in<f64>(arena)
    defer maps::free(&map)
    defer sets::free(&set)
    present bool[64] := .{}
    values u64[64] := .{}
    nans u64 := 0
    state := seed
    step u64 := 0
    while step < 2000 {
        state = state * 251 + 17
        state = state - state / 65521 * 65521
        key := state - state / 64 * 64
        op := state / 64
        op = op - op / 8 * 8
        if op < 3 {
            old, replaced := maps::insert(&map, _qualification_float(key), step) @fail
            added := sets::insert(&set, _qualification_float(key)) @fail
            testing::assert(replaced == present[key] and added != present[key])
            if replaced {
                testing::assert(old == values[key])
            }
            present[key] = true
            values[key] = step
        } elif op == 3 {
            stored, value, found := maps::remove(&map, _qualification_float(key))
            element, removed := sets::remove(&set, _qualification_float(key))
            testing::assert(found == present[key] and removed == found)
            if found {
                testing::assert(stored == _qualification_float(key) and element == stored and value == values[key])
            }
            present[key] = false
        } elif op == 4 {
            maps::insert(&map, _map_nan(), step) @fail
            sets::insert(&set, _map_nan()) @fail
            nans = nans + 1
        } elif op == 5 {
            maps::reserve(&map, key + 128) @fail
            sets::reserve(&set, key + 128) @fail
        } elif op == 6 {
            // Cross allocator clones must rehash NaNs without equality.
            copy := maps::clone(&map) @fail
            other := sets::clone_in(&set, arena) @fail
            verify(&copy, &other, present[:], values[:], nans)
            maps::clear(&copy)
            sets::free(&other)
            verify(&map, &set, present[:], values[:], nans)
            maps::free(&copy)
        } else {
            verify(&map, &set, present[:], values[:], nans)
        }
        if arena == nil and step == 500 {
            // Failure composition while both finite entries and NaNs are live.
            // A failed clone must not corrupt either source or its live borrows.
            maps::insert(&map, 0.0, step) @fail
            sets::insert(&set, 0.0) @fail
            present[0] = true
            values[0] = step
            testing::assert(nans > 0)
            baseline := _test_live_allocations()
            pointer := maps::get_ptr(&map, 0.0)
            attempt i64 := 0
            while attempt < 2 {
                cursor := sets::iter(&set)
                _test_fail_after(attempt)
                failed_map, map_error := maps::clone(&map)
                testing::assert(map_error.kind == errors::Error::AllocationFailed)
                testing::assert(maps::capacity(&failed_map) == 0)
                maps::free(&failed_map)
                _test_fail_after(attempt)
                failed_set, set_error := sets::clone(&set)
                testing::assert(set_error.kind == errors::Error::AllocationFailed)
                testing::assert(sets::capacity(&failed_set) == 0)
                sets::free(&failed_set)
                testing::assert(_test_live_allocations() == baseline)
                testing::assert(maps::get_ptr(&map, 0.0) == pointer)
                testing::assert(*pointer == step)
                visits u64 := 0
                while true {
                    item, found := sets::next(&cursor)
                    if not found {
                        break
                    }
                    visits = visits + 1
                }
                testing::assert(visits == sets::count(&set))
                attempt = attempt + 1
            }
            _test_fail_after(-1)
            _map_seed_failure(true)
            seed_map, map_error := maps::clone(&map)
            seed_set, set_error := sets::clone(&set)
            testing::assert(map_error.kind == errors::Error::RandomnessFailed)
            testing::assert(set_error.kind == errors::Error::RandomnessFailed)
            maps::free(&seed_map)
            sets::free(&seed_set)
            _map_seed_failure(false)
            verify(&map, &set, present[:], values[:], nans)
            testing::assert(_test_live_allocations() == baseline)
        }
        if step == 999 {
            maps::clear(&map)
            sets::clear(&set)
            present = zero<bool[64]>()
            nans = 0
        }
        step = step + 1
    }
    verify(&map, &set, present[:], values[:], nans)
    return
}

fn main() {
    seed u64 := 1
    while seed < 4 {
        run(nil, seed * 19)
        arena := memory::arena_create() @fail
        run(&arena, seed * 19)
        memory::arena_destroy(&arena)
        testing::assert(_test_live_allocations() == 0)
        seed = seed + 1
    }
    return
}
