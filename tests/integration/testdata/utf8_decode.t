package main

use "std::unicode::utf8"
use "std::strings"

fn main() i32 {
    text := "A\0é世😀"
    points u32[5] := .{65, 0, 233, 19990, 128512}
    widths u64[5] := .{1, 1, 2, 3, 4}
    offset u64 := 0
    index u64 := 0
    while index < 5 {
        decoded_point, decoded_width, decoded_err := utf8::decode(text, offset)
        if decoded_err != nil or decoded_point != points[index] or decoded_width != widths[index] {
            return 1
        }
        offset = offset + decoded_width
        index = index + 1
    }
    if strings::len(text) != 11 {
        return 2
    }
    count, err := utf8::count(text)
    if err != nil or count != 5 {
        return 3
    }
    empty, empty_err := utf8::count("")
    if empty_err != nil or empty != 0 {
        return 4
    }
    point, width, decode_err := utf8::decode(text, 3)
    if decode_err == nil or point != 0 or width != 0 {
        return 5
    }
    point, width, decode_err = utf8::decode(text, 11)
    if decode_err == nil or point != 0 or width != 0 {
        return 6
    }
    point, width, decode_err = utf8::decode(text, 18446744073709551615)
    if decode_err == nil or point != 0 or width != 0 {
        return 7
    }
    broken := strings::slice(text, 0, 10) @fail
    failed, failed_err := utf8::count(broken)
    if failed_err == nil or failed != 0 {
        return 8
    }
    local_point, local_width, local_err := utf8::decode(broken, 0)
    if local_err != nil or local_point != 65 or local_width != 1 {
        return 9
    }
    zero_point, zero_width, zero_err := utf8::decode("", 0)
    if zero_err == nil or zero_point != 0 or zero_width != 0 {
        return 10
    }
    return 0
}
