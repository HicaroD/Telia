package integration

import "testing"

func TestForeignArrayFieldsInitializeIndexAndPassExplicitElementPointers(t *testing.T) {
	runRecordFixture(t, `package main
@[header=%q, default_cc="c"]
extern {
@[link_name="FixtureBytes"]
struct Bytes {
data u8[4]
}
@[link_name="fixture_bytes_check"]
fn check(value Bytes, first *u8) i32
}
fn main() i32 {
value := Bytes{data: .{10, 20, 30, 40}}
value.data[1] = 42
if value.data[1] != 42 or check(value, &value.data[0]) != 1 {
return 1
}
return 0
}
`)
}
