package main

use "std::io"
use "std::libc"

fn identity<T>(value T) T {
    copy T := value
    return copy
}

fn main() i32 {
    number := identity<i32>(42)
    text := identity<string>("hello")
    other := identity<i32>(7)
    libc::printf("%d\n", number)
    io::println(text)
    libc::printf("%d\n", other)
    return 0
}
