package main

use "std::strings"
use "std::lists"
use "std::io"
use "std::libc"

fn main() i32 {
    owner, owner_err := strings::clone("　é,世界,😀 ")
    if owner_err != nil {
        return 1
    }
    defer strings::free(&owner)
    text, trim_err := strings::trim_space(strings::view(&owner))
    if trim_err != nil {
        return 2
    }
    offset, found := strings::find(text, "世界")
    if found {
        libc::printf("世界 at byte %llu\n", offset)
    }
    fields, split_err := strings::split(text, ",")
    if split_err != nil {
        return 3
    }
    defer lists::free<string>(&fields)
    index u64 := 0
    while index < lists::count<string>(&fields) {
        field, err := lists::get<string>(&fields, index)
        if err != nil {
            return 4
        }
        io::println(field)
        index = index + 1
    }
    return 0
}
