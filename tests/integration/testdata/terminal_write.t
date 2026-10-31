package main

use "std::terminal"

fn main() i32 {
    output := terminal::stdout()
    bytes u8[6] := .{65, 0, 66, 10, 195, 169}
    terminal::write_all(&output, bytes[:]) @fail
    return 0
}
