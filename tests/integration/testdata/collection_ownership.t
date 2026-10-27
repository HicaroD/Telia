package main

use "std::lists"
use "std::strings"
use "std::io"
use "std::libc"

fn main() i32 {
    values := lists::create<strings::Owned>()
    defer lists::free<strings::Owned>(&values)
    text := strings::clone("owned") @fail
    lists::append<strings::Owned>(&values, text) @fail
    view := lists::view<strings::Owned>(&values)
    element := &view[0]
    io::println(strings::view(element))
    strings::free(element)
    borrowed := strings::view(element)
    libc::printf("%llu\n", strings::len(borrowed))
    lists::clear<strings::Owned>(&values)
    return 0
}
