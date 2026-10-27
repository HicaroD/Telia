package main

use "std::memory"
use "std::libc"

fn make_zero<T>() T {
    return zero<T>()
}

fn main() i32 {
    libc::printf("%llu\n", sizeof<i32>())
    libc::printf("%llu\n", alignof<i32>())
    storage := memory::alloc(sizeof<i32>() * 2) @fail
    defer memory::free(storage)
    values := slice_from_raw<i32>(storage, 2)
    values[0] = 42
    values[1] = make_zero<i32>()
    libc::printf("%d\n", values[0])
    libc::printf("%d\n", values[1])
    return 0
}
