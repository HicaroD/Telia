package main

use "std::sets"
use "std::strings"
use "std::memory"
use "std::io"

fn run() error {
    // The owners outlive every table borrowing their bytes.
    original := strings::clone("shared") @try
    defer strings::free(&original)
    equal := strings::clone("shared") @try
    defer strings::free(&equal)
    arena := memory::arena_create() @try
    defer memory::arena_destroy(&arena)
    names := sets::create_in<string>(&arena)
    defer sets::free(&names)
    sets::insert(&names, strings::view(&original)) @try
    duplicate := sets::insert(&names, strings::view(&equal)) @try
    // An equal incoming element is not retained; it remains caller-owned.
    strings::free(&equal)
    heap := sets::clone(&names) @try
    defer sets::free(&heap)
    selected := sets::clone_in(&names, &arena) @try
    defer sets::free(&selected)
    sets::free(&names)
    sets::free(&selected)
    memory::arena_reset(&arena)
    // Independent heap backing survives reset. Its bytes are separately owned.
    cursor := sets::iter(&heap)
    while true {
        name, found := sets::next(&cursor)
        if not found {
            break
        }
        io::println(name)
    }
    stored, removed := sets::remove(&heap, "shared")
    if removed and not duplicate {
        io::println(stored)
    }
    return nil
}

fn main() i32 {
    run() @fail
    return 0
}
