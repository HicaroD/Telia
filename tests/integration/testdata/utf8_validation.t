package main

use "std::unicode::utf8"

fn main() i32 {
    if utf8::validate("") != nil {
        return 1
    }
    if utf8::validate("Hello\0world") != nil {
        return 2
    }
    return 0
}
