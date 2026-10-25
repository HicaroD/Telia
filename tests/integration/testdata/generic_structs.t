package main

use "std::io"
use "std::libc"

struct Box<T> {
    value T
}

fn main() i32 {
    number := Box<i32>.{ value: 42 }
    text := Box<string>.{ value: "hello" }
    other Box<i32> := Box<i32>.{ value: 7 }
    libc::printf("%d\n", number.value)
    io::println(text.value)
    libc::printf("%d\n", other.value)
    return 0
}
