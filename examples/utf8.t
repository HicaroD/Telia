package main

use "std::unicode::utf8"
use "std::strings"
use "std::io"
use "std::libc"

fn main() i32 {
    text := "Hi, 世界😀"
    total := utf8::count(text) @fail
    libc::printf("%llu bytes, %llu code points\n", strings::len(text), total)
    iter := utf8::iterator(text) @fail
    remaining := total
    while remaining > 0 {
        point, offset, found, err := utf8::next(&iter)
        if err != nil {
            io::println("invalid UTF-8")
            return 1
        }
        if not found {
            return 2
        }
        libc::printf("%u at byte %llu\n", point, offset)
        remaining = remaining - 1
    }
    return 0
}
