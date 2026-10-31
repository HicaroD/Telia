package main

use "std::process"
use "std::io"

fn main() i32 {
    snapshot, err := process::args()
    if err != nil {
        io::println("could not read arguments")
        return 1
    }
    defer process::free(&snapshot)
    index u64 := 0
    while index < process::count(&snapshot) {
        text, access_err := process::get(&snapshot, index)
        if access_err != nil {
            return 1
        }
        io::println(text)
        index = index + 1
    }
    return 0
}
