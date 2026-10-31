package main

use "std::terminal"

fn main() i32 {
    output := terminal::stdout()
    bytes u8[1] := .{65}
    count, err := terminal::write(&output, bytes[:])
    if err == nil or count != 0 {
        return 1
    }
    return 0
}
