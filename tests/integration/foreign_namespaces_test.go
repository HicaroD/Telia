package integration

import "testing"

func TestNamedExternRecordsResolveTheirDeclaringScope(t *testing.T) {
	runRecordFixture(t, `package main
@[header=%q, default_cc="c"]
extern ffi {
@[link_name="FixtureRecord"]
struct Record {
value i32
omitted i32
}
@[link_name="fixture_record_check"]
fn check(record *Record) i32
}
fn main() i32 {
record := ffi::Record{value: 42}
if ffi::check(&record) != 1 {
return 1
}
return 0
}
`)
}
