package integration

import "testing"

func TestForeignArrayFieldAliasesStillUseInlineCStorage(t *testing.T) {
	runRecordFixture(t, `package main
type Data = u8[4]
@[header=%q, default_cc="c"]
extern {
@[link_name="FixtureBytes"]
struct Bytes {
data Data
}
@[link_name="fixture_bytes_check"]
fn check(value Bytes, first *u8) i32
}
fn main() i32 {
value := Bytes{data: .{10, 42, 30, 40}}
if check(value, &value.data[0]) != 1 {
return 1
}
return 0
}
`)
}
