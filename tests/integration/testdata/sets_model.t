package main

use "std::sets"
use "std::memory"
use "std::testing"

fn verify(set *sets::Set<u64>, present bool[]) {
    expected u64 := 0
    key u64 := 0
    while key < len(present) {
        testing::assert(sets::contains(set, key) == present[key])
        if present[key] {
            expected = expected + 1
        }
        key = key + 1
    }
    testing::assert(sets::count(set) == expected and sets::capacity(set) >= expected)
    seen bool[128] := .{}
    cursor := sets::iter(set)
    visits u64 := 0
    while true {
        stored, found := sets::next(&cursor)
        if not found {
            break
        }
        testing::assert(stored < 128)
        testing::assert(present[stored] and not seen[stored])
        seen[stored] = true
        visits = visits + 1
    }
    testing::assert(visits == expected)
    return
}

fn run(arena *memory::Arena) {
    set := sets::create_in<u64>(arena)
    defer sets::free(&set)
    present bool[128] := .{}
    key u64 := 0
    while key < 128 {
        sets::insert(&set, key) @fail
        present[key] = true
        key = key + 1
    }
    verify(&set, present[:])
    state u64 := 19
    round u64 := 0
    while round < 3000 {
        state = state * 251 + 17
        state = state - state / 65521 * 65521
        key = state - state / 128 * 128
        operation := state / 128
        operation = operation - operation / 5 * 5
        if operation == 0 or operation == 1 {
            added := sets::insert(&set, key) @fail
            testing::assert(added != present[key])
            present[key] = true
        } elif operation == 2 {
            stored, removed := sets::remove(&set, key)
            testing::assert(removed == present[key])
            if removed {
                testing::assert(stored == key)
            } else {
                testing::assert(stored == 0)
            }
            present[key] = false
        } elif operation == 3 {
            sets::reserve(&set, key) @fail
            testing::assert(sets::capacity(&set) >= key)
        } else {
            verify(&set, present[:])
        }
        if round == 1500 {
            kept := sets::capacity(&set)
            sets::clear(&set)
            present = zero<bool[128]>()
            testing::assert(sets::capacity(&set) == kept)
        }
        if round == 1000 or round == 2000 {
            copy := sets::clone_in(&set, arena) @fail
            verify(&copy, present[:])
            sets::clear(&copy)
            verify(&set, present[:])
            sets::free(&copy)
            heap := sets::clone(&set) @fail
            verify(&heap, present[:])
            sets::free(&heap)
        }
        round = round + 1
    }
    verify(&set, present[:])
    sets::free(&set)
    round = 0
    while round < 100 {
        key = 0
        while key < 7 {
            sets::insert(&set, round * 7 + key) @fail
            key = key + 1
        }
        key = 0
        while key < 7 {
            stored, found := sets::remove(&set, round * 7 + key)
            testing::assert(found and stored == round * 7 + key)
            key = key + 1
        }
        testing::assert(sets::capacity(&set) == 7)
        round = round + 1
    }
    return
}

fn main() i32 {
    run(nil)
    arena := memory::arena_create() @fail
    defer memory::arena_destroy(&arena)
    run(&arena)
    return 0
}
