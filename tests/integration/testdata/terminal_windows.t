package main

use "std::terminal"
use "std::ffi::c"

@[header="terminal_console.h", default_cc="c"]
extern {
    fn terminal_test_console() i32
    fn terminal_test_key(handle c::c_intptr_t) i32
    fn terminal_test_scalar(handle c::c_intptr_t, scalar u32) i32
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
    if terminal_test_scalar(input.native, 128578) != 0 {
        return 10
    }
    expected u8[4] := .{240, 159, 153, 130}
    single := zero<u8[1]>()
    index u64 := 0
    while index < 4 {
        ready := terminal::wait(&input, terminal::WaitOptions{input: true, timeout_ms: 1000}) @fail
        if not ready.input or ready.timed_out {
            return 11
        }
        byte := terminal::read(&input, single[:]) @fail
        if byte.eof or byte.count != 1 or single[0] != expected[index] {
            return 12
        }
        empty_buffer := single[0:0] @fail
        no_bytes := terminal::read(&input, empty_buffer) @fail
        if no_bytes.eof or no_bytes.count != 0 {
            return 14
        }
        index = index + 1
    }
    empty := terminal::wait(&input, terminal::WaitOptions{input: true}) @fail
    if empty.input or not empty.timed_out {
        return 13
    }
    if terminal_test_scalar(input.native, 20013) != 0 {
        return 15
    }
    bmp := terminal::read(&input, buffer[:]) @fail
    if bmp.eof or bmp.count != 3 or buffer[0] != 228 or buffer[1] != 184 or buffer[2] != 173 {
        return 16
    }
    if terminal_test_scalar(input.native, 55296) != 0 or terminal_test_scalar(input.native, 65) != 0 {
        return 17
    }
    replacement := terminal::read(&input, buffer[:]) @fail
    if replacement.eof or replacement.count != 3 or buffer[0] != 239 or buffer[1] != 191 or buffer[2] != 189 {
        return 18
    }
    retained := terminal::wait(&input, terminal::WaitOptions{input: true, timeout_ms: 1000}) @fail
    if not retained.input or retained.timed_out {
        return 19
    }
    ascii := terminal::read(&input, buffer[:]) @fail
    if ascii.eof or ascii.count != 1 or buffer[0] != 65 {
        return 20
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
