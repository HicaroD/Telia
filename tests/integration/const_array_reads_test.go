package integration

import "testing"

func TestConstForeignArrayFieldsRemainReadable(t *testing.T) {
	runRecordFixture(t, `package main
@[header=%q, default_cc="c"]
extern {
@[link_name="FixtureBytes"]
struct Bytes {
data u8[4]
}
}
fn read(record const *Bytes) i32 {
if record.data[1] != 42 {
return 1
}
return 0
}
fn main() i32 {
value := Bytes{data: .{10, 42, 30, 40}}
return read(&value)
}
`)
}
