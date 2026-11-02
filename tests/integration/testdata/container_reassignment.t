package main

use "std::maps"
use "std::sets"
use "std::testing"

fn checked() error {
    map := maps::create<i32, i32>()
    defer maps::free(&map)
    old, replaced, err := maps::insert(&map, 1, 10)
    testing::assert(err == nil)
    old, replaced, err = maps::insert(&map, 1, 40)
    testing::assert(err == nil and old == 10 and replaced)
    value, found := maps::get(&map, 9)
    value, found = maps::get(&map, 1)
    testing::assert(value == 40 and found)
    cursor := maps::iter(&map)
    key, current, live := maps::next(&cursor)
    testing::assert(live and key == 1 and current == 40)
    key, current, live = maps::next(&cursor)
    testing::assert(not live and key == 0 and current == 0)
    key, current, live = maps::remove(&map, 1)
    testing::assert(live and key == 1 and current == 40)
    set := sets::create<i32>()
    defer sets::free(&set)
    sets::insert(&set, 7) @try
    iterator := sets::iter(&set)
    element, exists := sets::next(&iterator)
    testing::assert(exists and element == 7)
    element, exists = sets::next(&iterator)
    testing::assert(not exists and element == 0)
    element, exists = sets::remove(&set, 7)
    testing::assert(exists and element == 7)
    return nil
}

fn main() {
    checked() @fail
    return
}
