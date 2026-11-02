package main

use "std::maps"
use "std::sets"

struct Key {
    empty i32[0]
    marker i32
}

enum Choice {
    Empty
    Value {
        key Key
    }
}

fn next(trace *i32, digit i32) i32[0] {
    *trace = *trace * 10 + digit
    return zero<i32[0]>()
}

fn pair(marker i32) (i32[0], i32) {
    return zero<i32[0]>(), marker
}

// Exercise both semantic hashing and equality through public container APIs.
fn check_keys<K>(first K, equivalent K) error {
    map := maps::create<K, i32>()
    defer maps::free(&map)
    old, replaced := maps::insert(&map, first, 7) @try
    if replaced or old != 0 {
        return error("first map insert")
    }
    previous, updated := maps::insert(&map, equivalent, 9) @try
    value, found := maps::get(&map, first)
    if not updated or previous != 7 or not found or value != 9 or maps::count(&map) != 1 {
        return error("equivalent map key")
    }
    stored, removed, present := maps::remove(&map, equivalent)
    if not present or stored != first or removed != 9 or maps::contains(&map, first) {
        return error("remove map key")
    }

    set := sets::create<K>()
    defer sets::free(&set)
    inserted := sets::insert(&set, first) @try
    duplicate := sets::insert(&set, equivalent) @try
    if not inserted or duplicate or sets::count(&set) != 1 or not sets::contains(&set, equivalent) {
        return error("equivalent set key")
    }
    element, existed := sets::remove(&set, equivalent)
    if not existed or element != first or sets::contains(&set, first) {
        return error("remove set key")
    }
    return nil
}

fn main() i32 {
    a := zero<i32[0]>()
    b := zero<i32[0]>()
    // Representation bytes are not elements and must not affect equality/hash.
    pointer := &a
    bytes := slice_from_raw<u8>(pointer, sizeof<i32[0]>())
    bytes[0] = 255
    if not (a == b) or a != b {
        return 1
    }
    trace i32 := 0
    equal := next(&trace, 1) == next(&trace, 2)
    if not equal or trace != 12 {
        return 2
    }
    trace = 0
    different := next(&trace, 1) != next(&trace, 2)
    if different or trace != 12 {
        return 3
    }
    // A constant result must still obey surrounding short-circuit evaluation.
    trace = 0
    if false and next(&trace, 1) == next(&trace, 2) {
        return 4
    }
    if trace != 0 {
        return 5
    }

    key := Key{empty: a, marker: 7}
    equivalent := Key{empty: b, marker: 7}
    other := Key{empty: b, marker: 8}
    if not (key == equivalent) or key != equivalent or key == other or not (key != other) {
        return 6
    }
    values Key[2] := .{key, equivalent}
    copy := values
    if not (values == copy) or values != copy {
        return 7
    }
    copy[1] = other
    if values == copy or not (values != copy) {
        return 8
    }
    choice := Choice::Value{key: key}
    same_choice := Choice::Value{key: equivalent}
    other_choice := Choice::Value{key: other}
    if not (choice == same_choice) or choice != same_choice or choice == other_choice {
        return 9
    }
    if choice == Choice::Empty or not (choice != other_choice) {
        return 10
    }
    if not (pair(7) == pair(7)) or pair(7) != pair(7) or pair(7) == pair(8) or not (pair(7) != pair(8)) {
        return 11
    }

    check_keys(a, b) @fail
    check_keys(key, equivalent) @fail
    check_keys(choice, same_choice) @fail
    check_keys(values, values) @fail
    check_keys(pair(7), pair(7)) @fail
    return 0
}
