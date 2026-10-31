package main

use "std::terminal"
use "std::ffi::c"

@[header="terminal_console.h", default_cc="c"]
extern {
    fn terminal_test_console() i32
    fn terminal_test_key(handle c::c_intptr_t) i32
    fn terminal_test_resize(handle c::c_intptr_t) i32
    fn terminal_test_restored(input c::c_intptr_t, output c::c_intptr_t) i32
}

fn main() i32 {
    if terminal_test_console() != 0 {
        return 1
    }
    input := terminal::open_controlling(true) @fail
    defer terminal::close(&input)
    output := terminal::open_controlling(false) @fail
    defer terminal::close(&output)
    raw := terminal::apply_mode(&input, terminal::raw_mode()) @fail
    defer terminal::restore(&raw)
    prepared := terminal::prepare_output(&output) @fail
    defer terminal::restore(&prepared)
    mode := terminal::get_mode(&input) @fail
    if mode.echo or mode.canonical or mode.signals {
        return 2
    }
    initial := terminal::wait(&input, terminal::WaitOptions{input: true, resize: true}) @fail
    if not initial.timed_out {
        return 3
    }
    if terminal_test_key(input.native) != 0 {
        return 4
    }
    incoming := terminal::wait(&input, terminal::WaitOptions{input: true, timeout_ms: 1000}) @fail
    if not incoming.input or incoming.timed_out {
        return 5
    }
    buffer := zero<u8[16]>()
    result := terminal::read(&input, buffer[:]) @fail
    if result.eof or result.count != 2 or buffer[0] != 195 or buffer[1] != 169 {
        return 6
    }
    before := terminal::size(&output) @fail
    if terminal_test_resize(output.native) != 0 {
        return 7
    }
    resized := terminal::wait(&input, terminal::WaitOptions{resize: true, timeout_ms: 1000}) @fail
    after := terminal::size(&input) @fail
    if not resized.resized or after.columns != before.columns - 1 or after.rows != before.rows - 1 {
        return 8
    }
    terminal::restore(&prepared) @fail
    terminal::restore(&raw) @fail
    if terminal_test_restored(input.native, output.native) != 0 {
        return 9
    }
    return 0
}
