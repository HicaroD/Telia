package main

use "std::maps"
use "std::memory"
use "std::io"
use "std::libc"

struct Position {
    x i32
    y i32
}

fn run() error {
    scores maps::Map<string, i32> := maps::create()
    defer maps::free(&scores)
    maps::reserve(&scores, 20) @try
    old, replaced := maps::insert(&scores, "alice", 10) @try
    previous, updated := maps::insert(&scores, "alice", 20) @try
    inserted := maps::insert_if_absent(&scores, "alice", 99) @try
    value := maps::get_ptr(&scores, "alice")
    *value = *value + 5
    iterator := maps::iter(&scores)
    while true {
        key, current, found := maps::next(&iterator)
        if not found {
            break
        }
        io::println(key)
        libc::printf("%d\n", current)
    }
    stored_key, removed_value, removed := maps::remove(&scores, "alice")
    positions := maps::create<Position, string>()
    defer maps::free(&positions)
    maps::insert(&positions, Position{x: 3, y: 7}, "visited") @try
    label, found := maps::get(&positions, Position{x: 3, y: 7})
    io::println(label)
    arena := memory::arena_create() @try
    defer memory::arena_destroy(&arena)
    selected := maps::clone_in(&positions, &arena) @try
    defer maps::free(&selected)
    independent := maps::clone(&selected) @try
    defer maps::free(&independent)
    maps::clear(&selected)
    copied_label, copied_found := maps::get(&independent, Position{x: 3, y: 7})
    io::println(copied_label)
    destination := memory::take(&scores)
    defer maps::free(&destination)
    return nil
}

fn main() i32 {
    run() @fail
    return 0
}
