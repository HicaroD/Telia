package main

use "std::maps"
use "std::testing"
use "std::io"
use "std::memory"

fn verify(map *maps::Map<u64, u64>, present bool[], model u64[]) {
    expected u64 := 0
    key u64 := 0
    while key < len(model) {
        value, found := maps::get(map, key)
        testing::assert(found == present[key])
        testing::assert(maps::contains(map, key) == found)
        if found {
            testing::assert(value == model[key])
            expected = expected + 1
        } else {
            testing::assert(value == 0)
        }
        key = key + 1
    }
    testing::assert(maps::count(map) == expected)
    testing::assert(maps::capacity(map) >= expected)
    seen bool[128] := .{}
    cursor := maps::iter(map)
    visited u64 := 0
    while true {
        stored, value, found := maps::next(&cursor)
        if not found {
            break
        }
        testing::assert(stored < 128)
        testing::assert(not seen[stored] and present[stored])
        testing::assert(value == model[stored])
        seen[stored] = true
        visited = visited + 1
    }
    testing::assert(visited == expected)
    k, v, found := maps::next(&cursor)
    testing::assert(not found and k == 0 and v == 0)
    return
}

fn run(arena *memory::Arena) {
    map := maps::create_in<u64, u64>(arena)
    defer maps::free(&map)
    present bool[128] := .{}
    model u64[128] := .{}
    // Guarantee multi-group growth before randomized churn.
    key u64 := 0
    while key < 128 {
        maps::insert(&map, key, key + 1) @fail
        present[key] = true
        model[key] = key + 1
        key = key + 1
    }
    verify(&map, present[:], model[:])
    state u64 := 19
    round u64 := 0
    while round < 3000 {
        state = state * 251 + 17
        state = state - state / 65521 * 65521
        key = state - state / 128 * 128
        operation := state / 128
        operation = operation - operation / 7 * 7
        if operation == 0 or operation == 1 {
            old, replaced := maps::insert(&map, key, round) @fail
            testing::assert(replaced == present[key])
            if replaced {
                testing::assert(old == model[key])
            } else {
                testing::assert(old == 0)
            }
            present[key] = true
            model[key] = round
        } elif operation == 2 {
            stored, value, removed := maps::remove(&map, key)
            testing::assert(removed == present[key])
            if removed {
                testing::assert(stored == key and value == model[key])
            } else {
                testing::assert(stored == 0 and value == 0)
            }
            present[key] = false
        } elif operation == 3 {
            added := maps::insert_if_absent(&map, key, round) @fail
            testing::assert(added != present[key])
            if added {
                present[key] = true
                model[key] = round
            }
        } elif operation == 4 {
            pointer := maps::get_ptr(&map, key)
            if present[key] {
                *pointer = round
                model[key] = round
            } else {
                testing::assert(pointer == nil)
            }
        } elif operation == 5 {
            maps::reserve(&map, key) @fail
            testing::assert(maps::capacity(&map) >= key)
        } else {
            verify(&map, present[:], model[:])
        }
        if round == 1500 {
            kept := maps::capacity(&map)
            maps::clear(&map)
            present = zero<bool[128]>()
            testing::assert(maps::capacity(&map) == kept)
        }
        if round == 1000 or round == 2000 {
            copy := maps::clone_in(&map, arena) @fail
            verify(&copy, present[:], model[:])
            maps::clear(&copy)
            verify(&map, present[:], model[:])
            maps::free(&copy)
            heap := maps::clone(&map) @fail
            verify(&heap, present[:], model[:])
            maps::free(&heap)
        }
        round = round + 1
    }
    verify(&map, present[:], model[:])
    // Churn a near-load-bound table to force tombstone cleanup and reuse.
    maps::free(&map)
    round = 0
    while round < 100 {
        key = 0
        while key < 7 {
            maps::insert(&map, round * 7 + key, key) @fail
            key = key + 1
        }
        key = 0
        while key < 7 {
            stored, value, found := maps::remove(&map, round * 7 + key)
            testing::assert(found and value == key)
            key = key + 1
        }
        testing::assert(maps::capacity(&map) == 7)
        round = round + 1
    }
    return
}

fn main() i32 {
    run(nil)
    arena := memory::arena_create() @fail
    defer memory::arena_destroy(&arena)
    run(&arena)
    io::println("model ok")
    return 0
}
