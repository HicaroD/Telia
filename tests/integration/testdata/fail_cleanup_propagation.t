package main

use "std::io"

fn failure() (i32, error) {
    return 0, error("failure")
}

fn worker() error {
    defer io::println("worker first")
    defer io::println("worker second")
    if true {
        defer io::println("inner cleanup")
        value := failure() @catch err {
            return err
        }
    }
    return nil
}

fn run() error {
    defer io::println("caller cleanup")
    worker() @catch err {
        return err
    }
    return nil
}

fn main() i32 {
    run() @fail
    return 0
}
