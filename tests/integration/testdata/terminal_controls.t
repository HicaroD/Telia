package main

use "std::terminal"

fn main() i32 {
    buffer := zero<u8[64]>()
    count, err := terminal::encode_control(terminal::Control{kind: terminal::ControlKind::CursorPosition, first: 12, second: 34}, buffer[:])
    if err != nil or count != 8 {
        return 1
    }
    expected u8[8] := .{27, 91, 49, 50, 59, 51, 52, 72}
    i u64 := 0
    while i < 8 {
        if buffer[i] != expected[i] {
            return 2
        }
        i = i + 1
    }
    return 0
}
