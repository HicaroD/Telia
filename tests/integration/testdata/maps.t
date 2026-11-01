package main

use "std::maps"
use "std::memory"
use "std::libc"
use "std::io"

fn check(ok bool) {
    if not ok {
        libc::exit(7)
    }
    return
}

fn main() i32 {
    values := zero<maps::Map<string, i32>>()
    defer maps::free(&values)
    check(maps::count(&values) == 0 and maps::capacity(&values) == 0)
    missing, found := maps::get(&values, "absent")
    check(not found and missing == 0)
    check(maps::get_ptr(&values, "absent") == nil)
    cursor := maps::iter(&values)
    k, v, more := maps::next(&cursor)
    check(not more and k == "" and v == 0)
    old, replaced := maps::insert(&values, "one", 10) @fail
    check(old == 0 and not replaced)
    previous, updated := maps::insert(&values, "one", 20) @fail
    check(previous == 10 and updated)
    inserted := maps::insert_if_absent(&values, "one", 99) @fail
    check(not inserted)
    pointer := maps::get_ptr(&values, "one")
    cursor = maps::iter(&values)
    *pointer = 30
    key, value, exists := maps::next(&cursor)
    check(exists and key == "one" and value == 30)
    exhausted_key, exhausted_value, exhausted := maps::next(&cursor)
    check(not exhausted and exhausted_key == "" and exhausted_value == 0)
    stored, removed_value, removed := maps::remove(&values, "one")
    check(removed and stored == "one" and removed_value == 30)
    absent_key, absent_value, absent := maps::remove(&values, "one")
    check(not absent and absent_key == "" and absent_value == 0)
    maps::reserve(&values, 100) @fail
    check(maps::capacity(&values) >= 100)
    kept := maps::capacity(&values)
    maps::clear(&values)
    check(maps::capacity(&values) == kept)
    destination := memory::take(&values)
    maps::free(&destination)
    maps::free(&destination)
    check(maps::capacity(&values) == 0)
    maps::insert(&values, "again", 1) @fail
    check(maps::contains(&values, "again"))
    io::println("maps ok")
    return 0
}
