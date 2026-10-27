package main

use "std::libc"

fn sum<T>(values T[]) T {
    result T := 0
    for value in values {
        result = result + value
        value = 99
    }
    return result
}

fn first<T>(values T[3]) T {
    return values[0]
}

fn main() i32 {
    values i32[3] := .{1, 2, 3}
    view := values[:]
    libc::printf("%d\n", sum<i32>(view))
    libc::printf("%d\n", first<i32>(values))
    empty i32[0] := .{}
    total := 0
    for value in empty {
        total = total + value
    }
    libc::printf("%d\n", total)
    return 0
}
