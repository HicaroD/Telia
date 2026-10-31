package main

use "std::terminal"

fn main() i32 {
    input := terminal::stdin()
    output := terminal::stdout()
    poll := terminal::wait(&input, terminal::WaitOptions{input: true}) @fail
    if not poll.timed_out or poll.input or poll.resized {
        return 1
    }
    started u8[1] := .{83}
    terminal::write_all(&output, started[:]) @fail
    timed := terminal::wait(&input, terminal::WaitOptions{input: true, timeout_ms: 50}) @fail
    if not timed.timed_out or timed.input {
        return 2
    }
    ready u8[1] := .{82}
    terminal::write_all(&output, ready[:]) @fail
    result := terminal::wait(&input, terminal::WaitOptions{input: true, timeout_ms: -1}) @fail
    if not result.input or result.timed_out or result.resized {
        return 3
    }
    buffer := zero<u8[1]>()
    received := terminal::read(&input, buffer[:]) @fail
    if received.count != 1 or buffer[0] != 88 {
        return 4
    }
    return 0
}
