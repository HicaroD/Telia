package main

use "std::sets"
use "std::strings"
use "std::memory"
use "std::testing"

@[default_cc="c"]
extern {
    fn _test_live_allocations() i64
}

// Eligible records may contain pointers to explicitly owned resources; they
// compare pointer identity, not pointee contents. Known std owners themselves
// are deliberately not eligible elements.
struct Resource {
    allocation rawptr
}

fn run() error {
    first := strings::clone("same") @try
    defer strings::free(&first)
    incoming := strings::clone("same") @try
    defer strings::free(&incoming)
    lookup := strings::clone("same") @try
    defer strings::free(&lookup)
    testing::assert(strings::c_pointer(strings::view(&first)) != strings::c_pointer(strings::view(&incoming)))
    source := sets::create<string>()
    defer sets::free(&source)
    sets::insert(&source, strings::view(&first)) @try
    duplicate := sets::insert(&source, strings::view(&incoming)) @try
    testing::assert(not duplicate)
    arena := memory::arena_create() @try
    defer memory::arena_destroy(&arena)
    heap := sets::clone(&source) @try
    defer sets::free(&heap)
    selected := sets::clone_in(&source, &arena) @try
    defer sets::free(&selected)
    strings::free(&incoming)
    sets::reserve(&source, 100) @try
    cursor := sets::iter(&source)
    stored, found := sets::next(&cursor)
    testing::assert(found and strings::c_pointer(stored) == strings::c_pointer(strings::view(&first)))
    removed, has := sets::remove(&source, strings::view(&lookup))
    testing::assert(has and strings::c_pointer(removed) == strings::c_pointer(stored))
    sets::free(&source)
    sets::reserve(&heap, 100) @try
    heap_stored, heap_found := sets::remove(&heap, strings::view(&lookup))
    arena_stored, arena_found := sets::remove(&selected, strings::view(&lookup))
    testing::assert(heap_found and arena_found)
    testing::assert(strings::c_pointer(heap_stored) == strings::c_pointer(stored))
    testing::assert(strings::c_pointer(arena_stored) == strings::c_pointer(stored))
    resources := sets::create<Resource>()
    defer sets::free(&resources)
    index u64 := 0
    while index < 40 {
        allocation := memory::alloc(sizeof<u64>()) @try
        values := slice_from_raw<u64>(allocation, 1)
        values[0] = index
        added := sets::insert(&resources, Resource{allocation: allocation}) @try
        duplicate_resource := sets::insert(&resources, Resource{allocation: allocation}) @try
        testing::assert(added and not duplicate_resource)
        index = index + 1
    }
    // Shallow clones share resources. Discard clones without destroying them.
    resource_copy := sets::clone_in(&resources, &arena) @try
    sets::clear(&resource_copy)
    sets::free(&resource_copy)
    resource_cursor := sets::iter(&resources)
    seen bool[40] := .{}
    while true {
        resource, more := sets::next(&resource_cursor)
        if not more {
            break
        }
        values := slice_from_raw<u64>(resource.allocation, 1)
        testing::assert(values[0] < 40 and not seen[values[0]])
        seen[values[0]] = true
        memory::free(resource.allocation)
    }
    sets::clear(&resources)
    // Neither clear nor free destroys shallow elements.
    allocation := memory::alloc(8) @try
    sets::insert(&resources, Resource{allocation: allocation}) @try
    sets::clear(&resources)
    sets::insert(&resources, Resource{allocation: allocation}) @try
    sets::free(&resources)
    memory::free(allocation)
    return nil
}

fn main() i32 {
    run() @fail
    testing::assert(_test_live_allocations() == 0)
    return 0
}
