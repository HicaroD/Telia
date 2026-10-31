package main

use "std::strconv"
use "std::strings"
use "std::io"

fn main() i32 {
    text, trim_err := strings::trim_space("　-0042 ")
    if trim_err != nil {
        return 1
    }
    signed, signed_err := strconv::parse_i64(text)
    if signed_err != nil {
        return 2
    }
    unsigned, unsigned_err := strconv::parse_u64("18446744073709551615")
    if unsigned_err != nil {
        return 3
    }
    enabled, bool_err := strconv::parse_bool("true")
    if bool_err != nil {
        return 4
    }
    signed_text, signed_format_err := strconv::format_i64(signed)
    if signed_format_err != nil {
        return 5
    }
    defer strings::free(&signed_text)
    unsigned_text, unsigned_format_err := strconv::format_u64(unsigned)
    if unsigned_format_err != nil {
        return 6
    }
    defer strings::free(&unsigned_text)
    bool_text, bool_format_err := strconv::format_bool(enabled)
    if bool_format_err != nil {
        return 7
    }
    defer strings::free(&bool_text)
    io::println(strings::view(&signed_text))
    io::println(strings::view(&unsigned_text))
    io::println(strings::view(&bool_text))
    return 0
}
