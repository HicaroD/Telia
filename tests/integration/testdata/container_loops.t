package main

use "std::maps"
use "std::sets"
use "std::testing"

@[default_cc="c"]
extern {
    fn _test_fail_after(count i64)
    fn _test_live_allocations() i64
    fn _map_seed_failure(fail bool)
    fn _map_seed_calls() i64
}

type Scores = maps::Map<i32, i32>

struct Box {
    scores Scores
}

fn source(box *Box, calls *i32) *Box {
    *calls = *calls + 1
    return box
}

fn bump(value *i32) {
    *value = *value + 1
}

fn loop_index(calls *i32) u64 {
    *calls = *calls + 1
    return 0
}

fn failure() error {
    return error("expected")
}

fn propagate(scores *Scores, cleaned *i32) error {
    defer bump(cleaned)
    for key, value in *scores {
        defer bump(cleaned)
        failure() @try
    }
    return nil
}

fn early(scores *Scores, cleaned *i32) i32 {
    defer bump(cleaned)
    for key, value in *scores {
        defer bump(cleaned)
        return value
    }
    return 0
}

fn run() error {
    box := Box{}
    defer maps::free(&box.scores)
    elements := sets::create<i32>()
    defer sets::free(&elements)
    _test_fail_after(0)
    _map_seed_failure(true)
    for key, value in box.scores {
        testing::assert(false)
    }
    for element in elements {
        testing::assert(false)
    }
    testing::assert(_test_live_allocations() == 0 and _map_seed_calls() == 0)
    _test_fail_after(-1)
    _map_seed_failure(false)
    for i i32 := 1; i <= 20; i = i + 1 {
        maps::insert(&box.scores, i, i * 10) @try
        sets::insert(&elements, i) @try
    }
    // Tombstones are skipped and a reinserted live entry is visited once.
    removed_key, removed_value, removed := maps::remove(&box.scores, 7)
    stored, set_removed := sets::remove(&elements, 7)
    maps::insert(&box.scores, 7, 70) @try
    sets::insert(&elements, 7) @try
    live := _test_live_allocations()
    seeds := _map_seed_calls()
    _test_fail_after(0)
    _map_seed_failure(true)
    calls i32 := 0
    cleaned i32 := 0
    visits i32 := 0
    sum i32 := 0
    seen bool[21] := .{false, false, false, false, false, false, false, false, false, false, false, false, false, false, false, false, false, false, false, false, false}
    for key, value in source(&box, &calls).scores {
        defer bump(&cleaned)
        visits = visits + 1
        sum = sum + value
        index := narrow<u64>(key) @try
        testing::assert(not seen[index])
        seen[index] = true
        nested i32 := 0
        for element in elements {
            defer bump(&nested)
            if element != key {
                continue
            }
            testing::assert(sets::contains(&elements, element))
        }
        testing::assert(nested == 20)
        pointer := maps::get_ptr(&box.scores, key)
        *pointer = *pointer + 1
        key = 99
        value = 99
        continue
    }
    testing::assert(calls == 1 and visits == 20 and cleaned == 20 and sum == 2100)
    // Index expressions also borrow live storage and evaluate only once.
    pointer := &box.scores
    raw rawptr := pointer
    views := slice_from_raw<Scores>(raw, 1)
    indexed_sum i32 := 0
    for key, value in views[loop_index(&calls)] {
        indexed_sum = indexed_sum + value
    }
    testing::assert(calls == 2 and indexed_sum == 2120)
    testing::assert(_test_live_allocations() == live and _map_seed_calls() == seeds)
    stops i32 := 0
    for element in elements {
        defer bump(&stops)
        break
    }
    testing::assert(stops == 1)
    for key, value in box.scores {
        testing::assert(value == key * 10 + 1)
    }
    returned := early(&box.scores, &cleaned)
    testing::assert(returned > 0 and cleaned == 22)
    err := propagate(&box.scores, &cleaned)
    testing::assert(err != nil and cleaned == 24)
    // Explicit iterators and direct loops agree on the same live storage.
    iterator := maps::iter(&box.scores)
    explicit_sum i32 := 0
    while true {
        key, value, found := maps::next(&iterator)
        if not found {
            break
        }
        explicit_sum = explicit_sum + value
    }
    testing::assert(explicit_sum == 2120)
    _test_fail_after(-1)
    _map_seed_failure(false)
    return nil
}

fn main() i32 {
    run() @fail
    testing::assert(_test_live_allocations() == 0)
    return 0
}
