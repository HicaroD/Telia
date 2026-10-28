package main

use "std::libc"

type Number = libc::c_int
type NumberPointer = *Number

@[default_cc="c"]
extern {
    @[link_name="telia_test_pointer_round_trip"]
    fn round_trip(slot *NumberPointer) *Number
}

fn main() i32 {
    value Number := 40
    slot := &value
    result := round_trip(&slot)
    if *result != 42 {
        return 1
    }
    if value != 42 {
        return 2
    }
    return 0
}
