package main

use "std::io"

fn failure() error {
    return error("failure")
}

fn tuple_failure() (i32, error) {
    return 0, error("failure")
}

fn worker() {
    defer io::println("worker cleanup")
    if true {
        defer io::println("completed scope")
    }
    if true {
        defer io::println("inner cleanup")
        failure() @fail
    }
}

fn main() i32 {
    defer io::println("caller cleanup")
    worker()
    return 0
}
