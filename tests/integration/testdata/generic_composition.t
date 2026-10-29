package main

use "std::io"
use "std::libc"

struct Pair<A, B> {
    first A
    second B
}

fn identity<T>(value T) T {
    return value
}

fn make_pair<A, B>(first A, second B) Pair<A, B> {
    return Pair<A, B>{ first: identity<A>(first), second: identity<B>(second) }
}

fn countdown<T>(value T, remaining i32) T {
    if remaining == 0 {
        return value
    }
    return countdown<T>(value, remaining - 1)
}

fn main() i32 {
    pair := make_pair<i32, string>(42, "hello")
    libc::printf("%d\n", pair.first)
    io::println(pair.second)
    result := countdown<i32>(5, 150)
    libc::printf("%d\n", result)
    return 0
}
