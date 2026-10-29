package main

use "std::io"
use "std::libc"

struct Box<T> {
    value T
}

fn numbered<T>(value T) (T, i32) {
    return value, 1
}

fn main() i32 {
    number, n := numbered<Box<i32>>(Box<i32>{ value: 42 })
    text, m := numbered<Box<string>>(Box<string>{ value: "hello" })
    libc::printf("%d\n", number.value)
    io::println(text.value)
    return 0
}
