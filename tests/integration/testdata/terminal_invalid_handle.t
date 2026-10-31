package main

use "std::terminal"
use "std::ffi::c"

fn main() i32 {
    value c::c_intptr_t := -1
    invalid := terminal::borrow_native(value)
    result, err := terminal::wait(&invalid, terminal::WaitOptions{input: true})
    if err == nil or result.timed_out or result.input {
        return 1
    }
    return 0
}
