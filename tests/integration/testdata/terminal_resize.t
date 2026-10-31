package main

use "std::terminal"

fn main() i32 {
    input := terminal::stdin()
    output := terminal::stdout()
    saved := terminal::apply_mode(&input, terminal::raw_mode()) @fail
    defer terminal::restore(&saved)
    initial := terminal::wait(&input, terminal::WaitOptions{input: true, resize: true}) @fail
    if not initial.timed_out {
        return 1
    }
    ready u8[1] := .{82}
    terminal::write_all(&output, ready[:]) @fail
    resized := terminal::wait(&input, terminal::WaitOptions{input: true, resize: true, timeout_ms: 2000}) @fail
    if not resized.resized or resized.timed_out or resized.input {
        return 2
    }
    dimensions := terminal::size(&input) @fail
    if dimensions.columns != 100 or dimensions.rows != 40 {
        return 3
    }
    terminal::write_all(&output, ready[:]) @fail
    incoming := terminal::wait(&input, terminal::WaitOptions{input: true, resize: true, timeout_ms: 2000}) @fail
    if not incoming.input or incoming.resized or incoming.timed_out {
        return 4
    }
    buffer := zero<u8[1]>()
    read := terminal::read(&input, buffer[:]) @fail
    if read.count != 1 or buffer[0] != 3 {
        return 5
    }
    return 0
}
