package main

use "std::maps"
use "std::memory"
use "std::testing"
use "std::strings"

@[default_cc="c"]
extern {
    fn _map_nan() f64
    fn _map_seed_failure(fail bool)
    fn _test_fail_after(count i64)
}

fn nan_entries(map *maps::Map<f64, u64>) {
    testing::assert(maps::count(map) == 100)
    seen bool[100] := .{}
    cursor := maps::iter(map)
    total u64 := 0
    while true {
        key, value, found := maps::next(&cursor)
        if not found {
            break
        }
        testing::assert(key != key and value < 100 and not seen[value])
        seen[value] = true
        total = total + 1
    }
    testing::assert(total == 100)
    return
}

fn nan_clones(arena *memory::Arena) {
    source := maps::create<f64, u64>()
    defer maps::free(&source)
    index u64 := 0
    while index < 100 {
        maps::insert(&source, _map_nan(), index) @fail
        index = index + 1
    }
    heap := maps::clone(&source) @fail
    defer maps::free(&heap)
    selected := maps::clone_in(&source, arena) @fail
    defer maps::free(&selected)
    nan_entries(&source)
    nan_entries(&heap)
    nan_entries(&selected)
    maps::free(&source)
    maps::reserve(&heap, 1000) @fail
    maps::reserve(&selected, 1000) @fail
    nan_entries(&heap)
    nan_entries(&selected)
    maps::clear(&selected)
    testing::assert(maps::count(&selected) == 0)
    return
}

fn stored_identity(arena *memory::Arena) {
    first := strings::clone("same") @fail
    defer strings::free(&first)
    equal := strings::clone("same") @fail
    defer strings::free(&equal)
    source := maps::create<string, u64>()
    defer maps::free(&source)
    maps::insert(&source, strings::view(&first), 1) @fail
    maps::insert(&source, strings::view(&equal), 2) @fail
    heap := maps::clone(&source) @fail
    defer maps::free(&heap)
    selected := maps::clone_in(&source, arena) @fail
    defer maps::free(&selected)
    maps::free(&source)
    maps::reserve(&heap, 100) @fail
    maps::reserve(&selected, 100) @fail
    key, value, found := maps::remove(&heap, strings::view(&equal))
    testing::assert(found and value == 2 and strings::c_pointer(key) == strings::c_pointer(strings::view(&first)))
    arena_key, arena_value, arena_found := maps::remove(&selected, strings::view(&equal))
    testing::assert(arena_found and arena_value == 2 and strings::c_pointer(arena_key) == strings::c_pointer(key))
    return
}

fn main() i32 {
    arena := memory::arena_create() @fail
    defer memory::arena_destroy(&arena)
    nan_clones(&arena)
    stored_identity(&arena)
    empty := maps::create_in<u64, u64>(&arena)
    maps::reserve(&empty, 100) @fail
    maps::clear(&empty)
    _map_seed_failure(true)
    _test_fail_after(0)
    // Cloning an empty but reserved map needs neither seed nor allocation.
    empty_copy := maps::clone_in(&empty, &arena) @fail
    testing::assert(maps::count(&empty_copy) == 0 and maps::capacity(&empty_copy) == 0)
    maps::free(&empty_copy)
    maps::free(&empty)
    _map_seed_failure(false)
    _test_fail_after(-1)
    source := maps::create_in<u64, u64>(&arena)
    defer maps::free(&source)
    maps::insert(&source, 7, 42) @fail
    pointer := maps::get_ptr(&source, 7)
    copy := maps::clone(&source) @fail
    defer maps::free(&copy)
    value, found := maps::get(&copy, 7)
    testing::assert(found and value == 42)
    testing::assert(maps::get_ptr(&copy, 7) != pointer)
    copied_pointer := maps::get_ptr(&copy, 7)
    *copied_pointer = 99
    testing::assert(*pointer == 42)
    maps::free(&source)
    testing::assert(*maps::get_ptr(&copy, 7) == 99)
    return 0
}
