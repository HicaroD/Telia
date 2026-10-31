package main

use "std::fs"
use "std::io"
use "std::strings"
use "std::unicode::utf8"

fn main() i32 {
    fs::write_file("telia-file-example.bin", "Hello, 世界!\n\0") @fail
    bytes := fs::read_file("telia-file-example.bin") @fail
    defer strings::free(&bytes)

    text := strings::view(&bytes)
    utf8::validate(text) @fail
    io::println("File bytes validated as UTF-8")

    return 0
}
