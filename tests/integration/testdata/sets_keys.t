package main

use "std::sets"
use "std::testing"
use "std::strings"
use "std::memory"
use "std::ffi::c"
use "std::errors"

@[default_cc="c"]
extern {
    fn _map_nan() f64
    fn _map_negative_zero() f64
    fn _map_malformed() c::ConstCString
}

struct Empty {}
struct Box<T> {
    value T
}
enum Choice<T> {
    Empty
    Value {
        item T
    }
}
type Nested = Box<Box<string>>

fn roundtrip<T>(first T, equal T) {
    set := sets::create<T>()
    defer sets::free(&set)
    added := sets::insert(&set, first) @fail
    duplicate := sets::insert(&set, equal) @fail
    testing::assert(added and not duplicate and sets::count(&set) == 1)
    sets::reserve(&set, 50) @fail
    testing::assert(sets::contains(&set, equal))
    stored, found := sets::remove(&set, equal)
    testing::assert(found and stored == first)
    return
}

struct Nonreflexive {
    nan f64
    id u64
}

fn verify_nans(set *sets::Set<Nonreflexive>) {
    testing::assert(sets::count(set) == 100)
    seen bool[100] := .{}
    cursor := sets::iter(set)
    total u64 := 0
    while true {
        element, found := sets::next(&cursor)
        if not found {
            break
        }
        testing::assert(element.nan != element.nan)
        testing::assert(element.id < 100 and not seen[element.id])
        seen[element.id] = true
        total = total + 1
    }
    testing::assert(total == 100)
    return
}

fn nans(arena *memory::Arena) {
    source := sets::create_in<Nonreflexive>(arena)
    defer sets::free(&source)
    index u64 := 0
    while index < 100 {
        element := Nonreflexive{nan: _map_nan(), id: index}
        added := sets::insert(&source, element) @fail
        testing::assert(added and not sets::contains(&source, element))
        missing, found := sets::remove(&source, element)
        testing::assert(not found)
        index = index + 1
    }
    sets::reserve(&source, 1000) @fail
    heap := sets::clone(&source) @fail
    defer sets::free(&heap)
    selected := sets::clone_in(&source, arena) @fail
    defer sets::free(&selected)
    verify_nans(&source)
    sets::free(&source)
    sets::reserve(&heap, 2000) @fail
    sets::reserve(&selected, 2000) @fail
    verify_nans(&heap)
    verify_nans(&selected)
    // Identical NaNs also remain distinct (id is not what separates entries).
    scalar := sets::create<f64>()
    index = 0
    while index < 100 {
        added := sets::insert(&scalar, _map_nan()) @fail
        testing::assert(added)
        index = index + 1
    }
    copy := sets::clone(&scalar) @fail
    testing::assert(sets::count(&scalar) == 100 and sets::count(&copy) == 100)
    sets::free(&copy)
    sets::free(&scalar)
    return
}

fn tuple() (i32, string) {
    return 7, "tuple"
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
    roundtrip(maximum, maximum)
    roundtrip(-signed_maximum - 1, -signed_maximum - 1)
    roundtrip<int>(-7, -7)
    roundtrip<uint>(7, 7)
    roundtrip<c::c_long>(-7, -7)
    roundtrip<c::c_ulong>(7, 7)
    roundtrip<f32>(0.0, -0.0)
    roundtrip<float>(1.5, 1.5)
    roundtrip<f64>(0.0, _map_negative_zero())
    zeros := sets::create<f64>()
    sets::insert(&zeros, _map_negative_zero()) @fail
    duplicate := sets::insert(&zeros, 0.0) @fail
    stored_zero, found := sets::remove(&zeros, 0.0)
    testing::assert(not duplicate and found and 1.0 / stored_zero < 0.0)
    sets::free(&zeros)
    roundtrip<string>("", zero<string>())
    roundtrip<string>("a\0b", "a\0b")
    invalid := strings::from_c(_map_malformed())
    copy := strings::clone(invalid) @fail
    defer strings::free(&copy)
    roundtrip(invalid, strings::view(&copy))
    roundtrip(Empty{}, Empty{})
    roundtrip(Box<Empty>{value: Empty{}}, Box<Empty>{value: Empty{}})
    roundtrip(zero<Empty[0]>(), zero<Empty[0]>())
    nested := Nested{value: Box<string>{value: "nested"}}
    roundtrip(nested, nested)
    roundtrip(Choice<Box<string>>::Value{item: Box<string>{value: "payload"}}, Choice<Box<string>>::Value{item: Box<string>{value: "payload"}})
    roundtrip(Choice<Empty>::Empty, Choice<Empty>::Empty)
    roundtrip(Choice<Empty>::Value{item: Empty{}}, Choice<Empty>::Value{item: Empty{}})
    array Box<string>[2] := .{Box<string>{value: "a"}, Box<string>{value: "b"}}
    roundtrip(array, array)
    roundtrip(tuple(), tuple())
    roundtrip<*Box<string>>(&array[0], &array[0])
    roundtrip<*i32>(nil, nil)
    roundtrip<rawptr>(nil, nil)
    roundtrip<c::ConstCString>("native", "native")
    roundtrip<c::CString>(nil, nil)
    classified := error(errors::Error::OutOfBounds, "one")
    roundtrip(classified.kind, classified.kind)
    nans(nil)
    arena := memory::arena_create() @fail
    defer memory::arena_destroy(&arena)
    nans(&arena)
    return 0
}
