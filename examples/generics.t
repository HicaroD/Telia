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
    return Pair<A, B>{
        first: identity(first),
        second: identity(second)
    }
}

fn main() i32 {
    number i32 := 42
    entry := make_pair<string, _>("count", number)
    io::println(entry.first)
    libc::printf("%d\n", entry.second)
    return 0
}
