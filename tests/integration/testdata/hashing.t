package main

use "std::strings"
use "std::ffi::c"
use "std::errors"

@[header="equality_native.h", default_cc="c"]
extern {
 @[link_name="equality_nan"]
 fn nan() f64
 @[link_name="equality_infinity"]
 fn infinity() f64
 @[link_name="equality_negative_zero"]
 fn negzero() f64
 @[link_name="equality_nan32"]
 fn nan32() f32
 @[link_name="equality_negative_zero32"]
 fn negzero32() f32
 @[link_name="equality_malformed"]
 fn malformed() c::ConstCString
}
@[header="records.h", default_cc="c"]
extern {
 @[link_name="FixtureKind"]
 enum NativeKind {
  @[link_name="FIXTURE_KIND_READY"]
  Ready
 }
}

struct Pair {
 a string
 b string
}
struct Box<T> {
 value T
}
enum Choice<T> {
 Empty
 Value {
 item T
 }
}
struct Empty {}
struct Padded {
 small u8
 large u64
}
enum Unknown u8 {
 Only = 7 {
  text string
  value f64
 }
}
enum Scalar u8 {
 Only = 7
}
enum Failure: error {
 Detail {
  text string
 }
}
fn tuple(text string) (i32, string) {
 return 7, text
}
fn consistent<T>(left T, right T) bool {
 return left == right and hash(left) == hash(right)
}
fn next(counter *u64) u64 {
 *counter = *counter + 1
 return *counter
}
fn hash<T>(value T) u64 {
 return __key_hash<T>(value, 506097522914230528, 1084818905618843912)
}
fn main() i32 {
 seed := __key_seed() @fail
 if __key_hash<u64>(7, seed[0], seed[1]) != __key_hash<u64>(7, seed[0], seed[1]) {
 return 3
 }
 if hash(Pair{a: "a", b: "bc"}) == hash(Pair{a: "ab", b: "c"}) {
 return 1
 }
 a := Choice<Box<Pair>>::Value{item: Box<Pair>{value: Pair{a: "a", b: "b"}}}
 b := a
 if a != b or hash(a) != hash(b) {
 return 2
 }
 // Independent raw-byte vectors also cross the semantic compiler boundary.
 if hash<u8>(0) != 14552427766023178899 or hash<u16>(256) != 9424797079400076365 or hash<u32>(50462976) != 14948950610250269480 or hash<u64>(506097522914230528) != 3931806377309739662 {
  return 32
 }
 sequence u128 := 1084818905618843912
 sequence = (sequence << 64) + 506097522914230528
 if hash(sequence) != 14722228812742953830 {
  return 33
 }
 // All widths, sign extension, floats, C scalars, and pointer identities.
 if not consistent<bool>(true, true) or not consistent<bool>(false, false) {
  return 4
 }
 if not consistent<i8>(-128, -128) or not consistent<u8>(255, 255) or not consistent<i16>(-32768, -32768) or not consistent<u16>(65535, 65535) {
  return 5
 }
 if not consistent<i32>(-2147483648, -2147483648) or not consistent<u32>(4294967295, 4294967295) or not consistent<i64>(-9223372036854775808, -9223372036854775808) or not consistent<u64>(18446744073709551615, 18446744073709551615) {
  return 6
 }
 one u128 := 1
 high := one << 127
 max_u128 := (high - 1) + high
 signed_one i128 := 1
 quarter := signed_one << 126
 max_i128 := (quarter - 1) + quarter
 min_i128 := -max_i128 - 1
 if not consistent(min_i128, min_i128) or not consistent(max_u128, max_u128) {
  return 7
 }
 low := (one << 64) - 1
 if hash<u128>(0) == hash(one << 64) or hash<i128>(-1) == hash(low) {
  return 8
 }
 if not consistent<int>(-7, -7) or not consistent<uint>(7, 7) or not consistent<c::c_long>(-7, -7) or not consistent<c::c_ulong>(7, 7) {
  return 9
 }
 if not consistent<f64>(0.0, negzero()) or not consistent<f32>(0.0, negzero32()) or not consistent<float>(0.0, negzero32()) {
  return 10
 }
 if not consistent(infinity(), infinity()) or hash(infinity()) == hash(-infinity()) {
  return 11
 }
 n := nan()
 n32 := nan32()
 if n == n or n32 == n32 or hash(n) != hash(n) or hash(n32) != hash(n32) {
  return 12
 }
 // NaNs have no equal-value promise, but must hash safely, also in aggregates.
 if hash(Choice<f64>::Value{item: n}) != hash(Choice<f64>::Value{item: n}) {
  return 13
 }
 text := strings::clone("a\0b") @fail
 defer strings::free(&text)
 other := strings::clone("a\0b") @fail
 defer strings::free(&other)
 if not consistent(strings::view(&text), strings::view(&other)) or not consistent<string>("", zero<string>()) {
  return 14
 }
 invalid := strings::from_c(malformed())
 copy := strings::clone(invalid) @fail
 defer strings::free(&copy)
 if not consistent(invalid, strings::view(&copy)) or hash("a\0b") == hash("a") {
  return 15
 }
 if not consistent(tuple("a"), tuple("a")) or hash(tuple("a")) == hash(tuple("b")) {
  return 16
 }
 values Pair[2] := .{Pair{a: "a", b: "b"}, Pair{a: "c", b: "d"}}
 if not consistent(values, values) or not consistent(zero<Empty[0]>(), zero<Empty[0]>()) or not consistent(Empty{}, zero<Empty>()) {
  return 17
 }
 if not consistent<*Pair>(nil, nil) or not consistent<*Pair>(&values[0], &values[0]) or hash(&values[0]) == hash(&values[1]) {
  return 18
 }
 pointer := &values
 raw rawptr := pointer
 if not consistent(raw, raw) or not consistent<rawptr>(nil, nil) {
  return 19
 }
 if not consistent<c::ConstCString>(malformed(), malformed()) or not consistent<c::CString>(nil, nil) {
  return 20
 }
 if not consistent(NativeKind::Ready, NativeKind::Ready) or not consistent(Scalar::Only, Scalar::Only) or not consistent(zero<Scalar>(), zero<Scalar>()) or not consistent(Failure::Detail{text: "same"}, Failure::Detail{text: "same"}) {
  return 21
 }
 err := error(errors::Error::InvalidState, "one")
 other_err := error(errors::Error::InvalidState, "two")
 if not consistent(err.kind, other_err.kind) {
  return 22
 }
 // Deliberately invalid payload bytes must never be read for unknown tags.
 unknown := zero<Unknown>()
 unknown_other := zero<Unknown>()
 up := &unknown
 op := &unknown_other
 bytes := slice_from_raw<u8>(up, sizeof<Unknown>())
 other_bytes := slice_from_raw<u8>(op, sizeof<Unknown>())
 for i u64 := 1; i < sizeof<Unknown>(); i = i + 1 {
  bytes[i] = 255
  other_bytes[i] = 254
 }
 if not consistent(unknown, unknown_other) {
  return 23
 }
 bytes[0] = 1
 other_bytes[0] = 1
 if not consistent(unknown, unknown_other) {
  return 24
 }
 // Known payload-free variants also ignore dirty inactive storage.
 empty := zero<Choice<string>>()
 empty_other := empty
 ep := &empty
 eb := slice_from_raw<u8>(ep, sizeof<Choice<string>>())
 for i u64 := sizeof<int>(); i < sizeof<Choice<string>>(); i = i + 1 {
  eb[i] = 255
 }
 if not consistent(empty, empty_other) {
  return 25
 }
 padded := zero<Padded>()
 pp := &padded
 pb := slice_from_raw<u8>(pp, sizeof<Padded>())
 for i u64 := 0; i < sizeof<Padded>(); i = i + 1 {
  pb[i] = 255
 }
 padded.small = 0
 padded.large = 0
 if not consistent(padded, zero<Padded>()) {
  return 26
 }
 counter u64 := 0
 // Source order is value, k0, k1; every operand evaluated exactly once.
 ordered := __key_hash<u64>(next(&counter), next(&counter), next(&counter))
 if counter != 3 or ordered != __key_hash<u64>(1, 2, 3) {
  return 27
 }
 if false and __key_hash<u64>(next(&counter), 0, 0) == 0 {
  return 28
 }
 if counter != 3 {
  return 29
 }
 while hash(next(&counter)) != hash<u64>(6) {
  if counter > 6 {
   return 30
  }
 }
 if counter != 6 {
  return 31
 }
 return 0
}
