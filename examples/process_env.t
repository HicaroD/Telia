package main

use "std::process"
use "std::strings"
use "std::io"

fn main() i32 {
    value, found, err := process::lookup_env("TELIA_GREETING")
    if err != nil {
        io::write_stderr("could not read TELIA_GREETING\n")
        return 1
    }
    defer strings::free(&value)
    text := "Hello from Telia"
    if found {
        text = strings::view(&value)
    }
    write_err := io::write_stdout(text)
    if write_err != nil {
        return 1
    }
    return 0
}
