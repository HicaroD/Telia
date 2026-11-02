package main

use "std::sets"
use "std::testing"
use "std::io"
use "std::memory"

// Standard container eligibility must not constrain unrelated user names.
struct Set<T> {
    value T
}

// This checks the actual backend representation used by sets, not an assumed
// C extension for zero-sized records (see list_empty_struct.t).
fn representation() {
    testing::assert(sizeof<sets::Unit>() > 0)
    testing::assert(sizeof<sets::Set<i32>>() > 0)
    return
}

fn main() i32 {
    unrelated := zero<Set<i32[]>>()
    testing::assert(len(unrelated.value) == 0)
    representation()
    set sets::Set<i32> := sets::create()
    defer sets::free(&set)
    inserted := sets::insert(&set, 7) @fail
    testing::assert(inserted)
    duplicate := sets::insert(&set, 7) @fail
    testing::assert(not duplicate)
    testing::assert(sets::count(&set) == 1 and sets::contains(&set, 7))
    sets::reserve(&set, 100) @fail
    testing::assert(sets::capacity(&set) >= 100)
    cursor := sets::iter(&set)
    element, found := sets::next(&cursor)
    testing::assert(found and element == 7)
    exhausted, more := sets::next(&cursor)
    testing::assert(not more and exhausted == 0)
    exhausted_again, more_again := sets::next(&cursor)
    testing::assert(not more_again and exhausted_again == 0)
    cursor = zero<sets::Iterator<i32>>()
    zero_element, zero_found := sets::next(&cursor)
    testing::assert(not zero_found and zero_element == 0)
    missing, had := sets::remove(&set, 9)
    testing::assert(not had and missing == 0)
    removed, present := sets::remove(&set, 7)
    testing::assert(present and removed == 7 and sets::count(&set) == 0)
    sets::insert(&set, 8) @fail
    kept := sets::capacity(&set)
    sets::clear(&set)
    sets::clear(&set)
    testing::assert(sets::count(&set) == 0 and sets::capacity(&set) == kept)
    moved := memory::take(&set)
    sets::free(&moved)
    testing::assert(sets::capacity(&set) == 0)
    sets::insert(&set, 9) @fail
    sets::free(&set)
    sets::free(&set)
    zero_set := zero<sets::Set<i32>>()
    sets::insert(&zero_set, 1) @fail
    sets::free(&zero_set)
    io::println("sets ok")
    return 0
}
