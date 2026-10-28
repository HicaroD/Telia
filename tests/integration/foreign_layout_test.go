package integration

import "testing"

func TestHeaderOwnedPackingSizeAlignmentAndOffsetsMatchC(t *testing.T) {
	runRecordFixture(t, `package main
@[header=%q, default_cc="c"]
extern {
@[link_name="FixturePacked"]
struct Packed {
tag u8
value u64
}
@[link_name="fixture_packed_check"]
fn check(record *Packed, field *u64, size u64, alignment u64) i32
@[link_name="fixture_packed_roundtrip"]
fn roundtrip(record Packed) Packed
}
fn main() i32 {
record := Packed.{tag: 7, value: 40}
if check(&record, &record.value, sizeof<Packed>(), alignof<Packed>()) != 1 {
return 1
}
returned := roundtrip(record)
if returned.tag != 7 or returned.value != 42 {
return 2
}
return 0
}
`)
}
