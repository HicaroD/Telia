package main

use "std::io"
use "std::libc"

fn load<T>(ptr *T) T {
    return *ptr
}

fn maybe<T>(value T) (T, error) {
    return value, nil
}

fn roundtrip<T>(value T) T {
    result T := value
    for i := 0; i < 3; i = i + 1 {
        result = load<T>(&value)
    }
    counter := 1
    while counter > 0 {
        counter = counter - 1
    }
    if counter == 0 {
        copy := maybe<T>(result) @catch err {
            return value
        }
        defer io::println("cleanup")
        return copy
    } else {
        return value
    }
}

fn main() i32 {
    number := roundtrip<i32>(42)
    text := roundtrip<string>("hello")
    libc::printf("%d\n", number)
    io::println(text)
    return 0
}
