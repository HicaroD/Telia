package main

use "std::maps"
use "std::lists"
use "std::strings"
use "std::memory"
use "std::testing"
use "std::io"
use "std::ffi::c"

@[default_cc="c"]
extern {
    fn _test_live_allocations() i64
}

// Tests use the private text pointer helper only to distinguish equal views.
fn text_pointer(text string) c::CString {
    return strings::c_pointer(text)
}

fn run() error {
    map maps::Map<string, lists::List<strings::Owned>> := maps::create()
    defer maps::free(&map)
    original := strings::clone("same") @try
    defer strings::free(&original)
    incoming := strings::clone("same") @try
    defer strings::free(&incoming)
    lookup := strings::clone("same") @try
    defer strings::free(&lookup)
    values := lists::create<strings::Owned>()
    text := strings::clone("old") @try
    lists::append(&values, memory::take(&text)) @try
    old, replaced := maps::insert(&map, strings::view(&original), memory::take(&values)) @try
    testing::assert(not replaced and lists::count(&old) == 0)
    replacement := lists::create<strings::Owned>()
    new_text := strings::clone("new") @try
    lists::append(&replacement, memory::take(&new_text)) @try
    displaced, updated := maps::insert(&map, strings::view(&incoming), memory::take(&replacement)) @try
    testing::assert(updated)
    displaced_view := lists::view(&displaced)
    testing::assert(strings::view(&displaced_view[0]) == "old")
    strings::free(&displaced_view[0])
    lists::free(&displaced)
    unused := lists::create<strings::Owned>()
    unused_text := strings::clone("unused") @try
    lists::append(&unused, memory::take(&unused_text)) @try
    added := maps::insert_if_absent(&map, strings::view(&incoming), unused) @try
    testing::assert(not added)
    unused_view := lists::view(&unused)
    strings::free(&unused_view[0])
    lists::free(&unused)
    // Incoming equal keys remain caller-owned; the original key stays stored.
    cursor := maps::iter(&map)
    stored, copied, found := maps::next(&cursor)
    testing::assert(found and text_pointer(stored) == text_pointer(strings::view(&original)))
    testing::assert(text_pointer(stored) != text_pointer(strings::view(&incoming)))
    strings::free(&incoming)
    maps::reserve(&map, 100) @try
    removed_key, removed_value, removed := maps::remove(&map, strings::view(&lookup))
    testing::assert(removed and text_pointer(removed_key) == text_pointer(strings::view(&original)))
    removed_view := lists::view(&removed_value)
    testing::assert(strings::view(&removed_view[0]) == "new")
    strings::free(&removed_view[0])
    lists::free(&removed_value)
    // Clear/free never recursively destroy owned elements.
    owner := strings::clone("survives") @try
    owned_map := maps::create<u64, strings::Owned>()
    maps::insert(&owned_map, 1, owner) @try
    // Independent map buffers still share the one shallow resource owner.
    arena := memory::arena_create() @try
    defer memory::arena_destroy(&arena)
    heap_copy := maps::clone(&owned_map) @try
    arena_copy := maps::clone_in(&owned_map, &arena) @try
    copied_owner, copied_found := maps::get(&heap_copy, 1)
    testing::assert(copied_found and strings::view(&copied_owner) == "survives")
    maps::free(&heap_copy)
    maps::clear(&arena_copy)
    maps::free(&arena_copy)
    maps::clear(&owned_map)
    testing::assert(strings::view(&owner) == "survives")
    maps::insert(&owned_map, 2, owner) @try
    maps::free(&owned_map)
    strings::free(&owner)
    // Transfer preserves backing while leaving the source usable and cleanup-safe.
    destination := memory::take(&map)
    maps::free(&destination)
    maps::insert(&map, "zero source", lists::create<strings::Owned>()) @try
    maps::clear(&map)
    return nil
}

fn main() i32 {
    run() @fail
    testing::assert(_test_live_allocations() == 0)
    io::println("ownership ok")
    return 0
}
