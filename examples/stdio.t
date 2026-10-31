package main

use "std::io"
use "std::strings"

fn main() i32 {
    bytes, read_err := io::read_stdin()
    if read_err != nil {
        io::write_stderr("stdin read failed\n") @fail
        return 1
    }
    defer strings::free(&bytes)

    write_err := io::write_stdout(strings::view(&bytes))
    if write_err != nil {
        io::write_stderr("stdout write failed\n") @fail
        return 1
    }
    return 0
}
