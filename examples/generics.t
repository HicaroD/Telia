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
        first: identity<A>(first),
        second: identity<B>(second)
    }
}

fn main() i32 {
    entry := make_pair<string, i32>("count", 42)
    io::println(entry.first)
    libc::printf("%d\n", entry.second)
    return 0
}
