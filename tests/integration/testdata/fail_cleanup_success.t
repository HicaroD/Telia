package main

use "std::io"

fn success() error {
    return nil
}

fn tuple_success() (i32, error) {
    return 7, nil
}

fn worker() {
    defer io::println("worker cleanup")
    success() @fail
    tuple_success() @fail
    value := tuple_success() @fail
    if value == 7 {
        io::println("continued")
    }
}

fn main() i32 {
    defer io::println("caller cleanup")
    worker()
    return 0
}
