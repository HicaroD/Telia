package main

use "std::libc"

struct Box<T> {
    value T
}

fn zero<T>() T {
    return T{}
}

fn main() i32 {
    value := zero<Box<i32>>()
    libc::printf("%d\n", value.value)
    return 0
}
