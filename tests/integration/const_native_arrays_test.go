package integration

import "testing"

func TestConstNativeArrayPointersCanBeIndexedForReading(t *testing.T) {
	runRecordFixture(t, `package main
fn read_fixture_values(values const *i32[2]) i32 {
return (*values)[1]
}
fn main() i32 {
values i32[2] := .{1, 42}
if read_fixture_values(&values) != 42 {
return 1
}
return 0
}
`)
}
