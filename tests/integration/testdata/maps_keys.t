package main

use "std::maps"
use "std::strings"
use "std::testing"
use "std::ffi::c"
use "std::errors"
use "std::io"

@[default_cc="c"]
extern {
    fn _map_nan() f64
    fn _map_infinity() f64
    fn _map_negative_zero() f64
    fn _map_malformed() c::ConstCString
}

struct Box<T> {
    value T
}
struct Empty {}
enum Choice<T> {
    Empty
    Value {
        item T
    }
}
enum Unknown u8 {
    Only = 7 {
        text string
    }
}

fn roundtrip<T>(first T, equal T) {
    map := maps::create<T, u64>()
    defer maps::free(&map)
    old, replaced := maps::insert(&map, first, 1) @fail
    testing::assert(old == 0 and not replaced)
    previous, updated := maps::insert(&map, equal, 2) @fail
    testing::assert(previous == 1 and updated)
    testing::assert(maps::count(&map) == 1)
    maps::reserve(&map, 50) @fail
    value, found := maps::get(&map, first)
    testing::assert(found and value == 2)
    stored, removed, has := maps::remove(&map, equal)
    testing::assert(has and stored == first and removed == 2)
    return
}

fn tuple() (i32, string) {
    return 7, "tuple"
}

fn nan_growth() {
    map := maps::create<Box<f64>, u64>()
    defer maps::free(&map)
    key := Box<f64>{value: _map_nan()}
    index u64 := 0
    while index < 100 {
        old, replaced := maps::insert(&map, key, index) @fail
        testing::assert(old == 0 and not replaced)
        index = index + 1
    }
    testing::assert(maps::count(&map) == 100)
    testing::assert(not maps::contains(&map, key))
    stored, removed, found := maps::remove(&map, key)
    testing::assert(not found)
    maps::reserve(&map, 1000) @fail
    seen bool[100] := .{}
    cursor := maps::iter(&map)
    total u64 := 0
    while true {
        nan_key, value, more := maps::next(&cursor)
        if not more {
            break
        }
        testing::assert(nan_key.value != nan_key.value)
        testing::assert(value < 100 and not seen[value])
        seen[value] = true
        total = total + 1
    }
    testing::assert(total == 100)
    maps::clear(&map)
    testing::assert(maps::count(&map) == 0)
    return
}

fn main() i32 {
    roundtrip<bool>(false, false)
    roundtrip<i8>(-128, -128)
    roundtrip<u8>(255, 255)
    roundtrip<i16>(-32768, -32768)
    roundtrip<u16>(65535, 65535)
    roundtrip<i32>(-2147483648, -2147483648)
    roundtrip<u32>(4294967295, 4294967295)
    roundtrip<i64>(-9223372036854775808, -9223372036854775808)
    roundtrip<u64>(18446744073709551615, 18446744073709551615)
    one u128 := 1
    high := one << 127
    maximum := high + (high - 1)
    signed_one i128 := 1
    quarter := signed_one << 126
    signed_maximum := quarter + (quarter - 1)
    signed_minimum := -signed_maximum - 1
    roundtrip(maximum, maximum)
    roundtrip(signed_minimum, signed_minimum)
    roundtrip<int>(-7, -7)
    roundtrip<uint>(7, 7)
    roundtrip<c::c_long>(-7, -7)
    roundtrip<c::c_ulong>(7, 7)
    roundtrip<f64>(0.0, _map_negative_zero())
    roundtrip<f64>(_map_infinity(), _map_infinity())
    roundtrip<f64>(-_map_infinity(), -_map_infinity())
    roundtrip<f32>(0.0, -0.0)
    roundtrip<float>(1.5, 1.5)
    roundtrip<string>("", zero<string>())
    first := strings::clone("a\0b") @fail
    defer strings::free(&first)
    equal := strings::clone("a\0b") @fail
    defer strings::free(&equal)
    roundtrip(strings::view(&first), strings::view(&equal))
    invalid := strings::from_c(_map_malformed())
    invalid_copy := strings::clone(invalid) @fail
    defer strings::free(&invalid_copy)
    roundtrip(invalid, strings::view(&invalid_copy))
    roundtrip(Box<string>{value: "nested"}, Box<string>{value: "nested"})
    roundtrip(Choice<Box<string>>::Value{item: Box<string>{value: "payload"}}, Choice<Box<string>>::Value{item: Box<string>{value: "payload"}})
    roundtrip(Choice<string>::Empty, Choice<string>::Empty)
    roundtrip(Empty{}, Empty{})
    roundtrip(zero<Empty[0]>(), zero<Empty[0]>())
    array Box<string>[2] := .{Box<string>{value: "a"}, Box<string>{value: "b"}}
    roundtrip(array, array)
    roundtrip(tuple(), tuple())
    roundtrip<*Box<string>>(&array[0], &array[0])
    roundtrip<*i32>(nil, nil)
    roundtrip<rawptr>(nil, nil)
    pointer := &array
    raw rawptr := pointer
    roundtrip(raw, raw)
    roundtrip<c::ConstCString>("native", "native")
    roundtrip<c::CString>(nil, nil)
    classified := error(errors::Error::OutOfBounds, "one")
    roundtrip(classified.kind, classified.kind)
    // Unknown tags hash and compare only tags, never poisoned inactive data.
    unknown := zero<Unknown>()
    other := zero<Unknown>()
    unknown_pointer := &unknown
    other_pointer := &other
    bytes := slice_from_raw<u8>(unknown_pointer, sizeof<Unknown>())
    other_bytes := slice_from_raw<u8>(other_pointer, sizeof<Unknown>())
    index u64 := 1
    while index < sizeof<Unknown>() {
        bytes[index] = 255
        other_bytes[index] = 254
        index = index + 1
    }
    roundtrip(unknown, other)
    nan_growth()
    io::println("keys ok")
    return 0
}
