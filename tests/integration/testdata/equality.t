package main

use "std::strings"
use "std::memory"
use "std::ffi::c"

@[header="equality_native.h", default_cc="c"]
extern {
 @[link_name="equality_nan"]
 fn make_nan() f64
 @[link_name="equality_infinity"]
 fn infinity() f64
 @[link_name="equality_negative_zero"]
 fn negative_zero() f64
 @[link_name="equality_nan32"]
 fn make_nan32() f32
 @[link_name="equality_negative_zero32"]
 fn negative_zero32() f32
 @[link_name="equality_malformed"]
 fn malformed_text() c::ConstCString
}

struct Point {
 x i32
 y i32
}

enum Choice<T> {
 Empty
 Value {
  item T
 }
}

struct Box<T> {
 value T
}

fn same<T>(left T, right T) bool {
 return left == right
}

struct Empty {}

struct Node {
 next *Node
 unsupported i32[]
}

struct Padded {
 small u8
 big u64
}

enum Dirty {
 Empty
 Large {
  text string
  values f64[8]
 }
}

enum Failure: error {
 Detail {
  text string
 }
}

enum NonzeroPayload u8 {
 Only = 7 {
  text string
  value f64
 }
}

enum NonzeroScalar u8 {
 Only = 7
}

fn pair(n i32) (i32, string) {
 return n, "pair"
}

fn pair_text(text string) (i32, string) {
 return 1, text
}

fn next(counter *i32, expected i32) Point {
 *counter = *counter + 1
 return Point{x: *counter - expected, y: 7}
}

fn callback(value i32) i32 {
 return value
}

fn loop_next(counter *i32) Point {
 *counter = *counter + 1
 return Point{x: *counter / 3, y: 0}
}

struct Callbacks {
 call fn(value i32) i32
}

fn main() i32 {
 if nil != nil {
  return 30
 }
 a := Point{x: 1, y: 2}
 b := Point{x: 1, y: 2}
 c := Point{x: 1, y: 3}
 if not (a == b) {
 return 1
 }
 if a != b {
 return 2
 }
 if a == c {
 return 3
 }
 if not (a != c) {
 return 4
 }
 x := Choice<Box<Point>>::Value{item: Box<Point>{value: a}}
 y := Choice<Box<Point>>::Value{item: Box<Point>{value: b}}
 z := Choice<Box<Point>>::Empty
 if not same(x, y) {
  return 5
 }
 if x == z {
  return 6
 }
 if z != Choice<Box<Point>>::Empty {
  return 7
 }
 different := Choice<Box<Point>>::Value{item: Box<Point>{value: c}}
 if x == different or not (x != different) {
  return 31
 }

 values Point[2] := .{a, b}
 copy := values
 if values != copy {
  return 8
 }
 copy[1] = c
 if values == copy {
  return 9
 }
 if pair(1) != pair(1) or pair(1) == pair(2) {
  return 10
 }
 if pair_text("a") == pair_text("b") or not (pair_text("a") != pair_text("b")) {
  return 43
 }
 if (zero<Empty>() != Empty{}) {
  return 11
 }
 empty Empty[0] := .{}
 if empty != empty {
  return 12
 }

 count i32 := 0
 if next(&count, 1) != next(&count, 2) or count != 2 {
  return 13
 }
 if false and next(&count, 3) == next(&count, 4) {
  return 14
 }
 if count != 2 {
  return 15
 }
 count = 0
 while loop_next(&count) == zero<Point>() {
  if count > 2 {
   return 32
  }
 }
 if count != 3 {
  return 33
 }

 owned := strings::clone("a\0b") @fail
 defer strings::free(&owned)
 other := strings::clone("a\0b") @fail
 defer strings::free(&other)
 if owned != other {
  return 16
 }
 text := Box<string>{value: strings::view(&owned)}
 if (text != Box<string>{value: "a\0b"}) {
  return 17
 }
 // Native strings also accept malformed UTF-8. from_c borrows these bytes.
 malformed := strings::from_c(malformed_text())
 if (Box<string>{value: malformed} != Box<string>{value: malformed}) {
  return 18
 }
 if (zero<Box<string>>() != Box<string>{value: ""}) {
  return 34
 }
 c_owner := strings::to_c("same") @fail
 defer strings::free_c(&c_owner)
 c_other := strings::to_c("same") @fail
 defer strings::free_c(&c_other)
 c_pointer := strings::c_view(&c_owner)
 if (Box<c::CString>{value: c_pointer} != Box<c::CString>{value: c_pointer}) {
  return 35
 }
 if (Box<c::CString>{value: c_pointer} == Box<c::CString>{value: strings::c_view(&c_other)}) {
  return 36
 }
 if (zero<Box<c::ConstCString>>() != Box<c::ConstCString>{value: nil}) {
  return 37
 }
 if (Box<c::c_long>{value: 17} != Box<c::c_long>{value: 17}) {
  return 38
 }

 n := make_nan()
 inf := infinity()
 zero_value f64 := 0.0
 negzero := negative_zero()
 if n == n or not (n != n) or inf != inf or inf == -inf {
  return 19
 }
 if (zero_value != negzero or Box<f64>{value: zero_value} != Box<f64>{value: negzero}) {
  return 20
 }
 floats f64[2] := .{inf, n}
 if (floats == floats or not (floats != floats) or Box<f64>{value: n} == Box<f64>{value: n}) {
  return 21
 }
 n32 := make_nan32()
 zero32 f32 := 0.0
 if (n32 == n32 or Box<f32>{value: n32} == Box<f32>{value: n32}) {
  return 39
 }
 if (zero32 != negative_zero32() or Box<f32>{value: zero32} != Box<f32>{value: negative_zero32()}) {
  return 40
 }
 nan_enum := Choice<f64>::Value{item: n}
 if nan_enum == nan_enum or not (nan_enum != nan_enum) {
  return 41
 }

 node := Node{}
 other_node := Node{}
 p := Box<*Node>{value: &node}
 if (p != Box<*Node>{value: &node} or p == Box<*Node>{value: &other_node}) {
  return 22
 }
 if (zero<Box<*Node>>() != Box<*Node>{value: nil}) {
  return 23
 }
 node_pointer := &node
 raw rawptr := node_pointer
 if (Box<rawptr>{value: raw} != Box<rawptr>{value: raw}) {
  return 24
 }
 calls := Callbacks{call: callback}
 if (calls != Callbacks{call: callback} or calls == zero<Callbacks>()) {
  return 25
 }
 if (Failure::Detail{text: "same"} != Failure::Detail{text: "same"}) {
  return 26
 }
 if (Failure::Detail{text: "same"} == Failure::Detail{text: "different"}) {
  return 27
 }
 // Different inactive bytes include invalid string pointers; never inspect them.
 dirty := zero<Dirty>()
 dirty_other := zero<Dirty>()
 dirty_pointer := &dirty
 other_dirty_pointer := &dirty_other
 bytes := slice_from_raw<u8>(dirty_pointer, sizeof<Dirty>())
 other_bytes := slice_from_raw<u8>(other_dirty_pointer, sizeof<Dirty>())
 for i u64 := sizeof<int>(); i < sizeof<Dirty>(); i = i + 1 {
  bytes[i] = 255
  other_bytes[i] = 254
 }
 if dirty != dirty_other {
  return 28
 }
 large := Dirty::Large{text: "large", values: zero<f64[8]>()}
 if dirty == large or large != large {
  return 42
 }
 // Padding differs; only the declared fields can contribute to equality.
 padded := zero<Padded>()
 padded_other := zero<Padded>()
 padded_pointer := &padded
 padding := slice_from_raw<u8>(padded_pointer, sizeof<Padded>())
 for i u64 := 0; i < sizeof<Padded>(); i = i + 1 {
  padding[i] = 255
 }
 padded.small = 0
 padded.big = 0
 if padded != padded_other {
  return 29
 }
 // zero<T>() retains zero bytes, not the first declared variant (tag 7).
 unknown := zero<NonzeroPayload>()
 unknown_other := zero<NonzeroPayload>()
 if unknown != unknown_other or not (unknown == unknown_other) {
  return 44
 }
 if zero<NonzeroScalar>() != zero<NonzeroScalar>() {
  return 45
 }
 if zero<NonzeroScalar>() == NonzeroScalar::Only {
  return 50
 }
 unknown_pointer := &unknown
 unknown_other_pointer := &unknown_other
 unknown_bytes := slice_from_raw<u8>(unknown_pointer, sizeof<NonzeroPayload>())
 unknown_other_bytes := slice_from_raw<u8>(unknown_other_pointer, sizeof<NonzeroPayload>())
 for i u64 := sizeof<u8>(); i < sizeof<NonzeroPayload>(); i = i + 1 {
  unknown_bytes[i] = 255
  unknown_other_bytes[i] = 254
 }
 if unknown != unknown_other or not (unknown == unknown_other) {
  return 46
 }
 if (Box<NonzeroPayload>{value: unknown} != Box<NonzeroPayload>{value: unknown_other}) {
  return 51
 }
 known := NonzeroPayload::Only{text: "known", value: 0.0}
 if unknown == known or not (unknown != known) {
  return 47
 }
 unknown_other_bytes[0] = 1
 if unknown == unknown_other or not (unknown != unknown_other) {
  return 48
 }
 // A second undeclared tag also ignores invalid payload bytes.
 unknown_bytes[0] = 1
 if unknown != unknown_other {
  return 49
 }
 return 0
}
