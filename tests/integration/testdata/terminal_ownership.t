package main

use "std::terminal"

fn main() i32 {
    output := terminal::stdout()
    borrowed := terminal::borrow_native(output.native)
    terminal::close(&borrowed) @fail
    terminal::close(&borrowed) @fail
    bytes u8[2] := .{111, 107}
    count, err := terminal::write(&borrowed, bytes[:])
    if err == nil or count != 0 {
        return 1
    }
    terminal::write_all(&output, bytes[:]) @fail
    return 0
}
