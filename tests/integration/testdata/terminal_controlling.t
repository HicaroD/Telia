package main

use "std::terminal"

fn main() i32 {
    device := terminal::open_controlling(true) @fail
    defer terminal::close(&device)
    attached := terminal::is_terminal(&device) @fail
    if not attached or not device.owned {
        return 1
    }
    alias := terminal::borrow_native(device.native)
    terminal::close(&device) @fail
    invalid, invalid_err := terminal::is_terminal(&alias)
    if invalid_err == nil or invalid {
        return 2
    }
    output := terminal::stdout()
    ready u8[1] := .{82}
    terminal::write_all(&output, ready[:]) @fail
    return 0
}
