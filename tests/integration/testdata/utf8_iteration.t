package main

use "std::unicode::utf8"
use "std::strings"

fn main() i32 {
    owner := strings::clone("A\0é世😀") @fail
    defer strings::free(&owner)
    iter := utf8::iterator(strings::view(&owner)) @fail
    copied := iter
    points u32[5] := .{65, 0, 233, 19990, 128512}
    offsets u64[5] := .{0, 1, 2, 4, 7}
    index u64 := 0
    while index < 5 {
        point, offset, found, err := utf8::next(&iter)
        if err != nil or not found or point != points[index] or offset != offsets[index] {
            return 1
        }
        index = index + 1
    }
    end_point, end_offset, end_found, end_err := utf8::next(&iter)
    if end_err != nil or end_found or end_point != 0 or end_offset != 0 {
        return 2
    }
    end_point, end_offset, end_found, end_err = utf8::next(&iter)
    if end_err != nil or end_found {
        return 3
    }
    empty := utf8::iterator("") @fail
    empty_point, empty_offset, empty_found, empty_err := utf8::next(&empty)
    if empty_err != nil or empty_found or empty_point != 0 or empty_offset != 0 {
        return 4
    }
    broken := strings::slice("valid😀", 0, 8) @fail
    bad, bad_err := utf8::iterator(broken)
    if bad_err == nil {
        return 5
    }
    bad_point, bad_offset, bad_found, next_err := utf8::next(&bad)
    if next_err != nil or bad_found or bad_point != 0 or bad_offset != 0 {
        return 6
    }
    copy_point, copy_offset, copy_found, copy_err := utf8::next(&copied)
    if copy_err != nil or not copy_found or copy_point != 65 or copy_offset != 0 {
        return 7
    }
    initial := utf8::Iterator{}
    initial_point, initial_offset, initial_found, initial_err := utf8::next(&initial)
    if initial_err != nil or initial_found or initial_point != 0 or initial_offset != 0 {
        return 8
    }
    return 0
}
