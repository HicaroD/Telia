package integration

import "testing"

func TestNativeCallbackArrayResultsKeepWrapperIdentity(t *testing.T) {
	runRecordFixture(t, `package main
type Make = fn() i32[2]
fn make() i32[2] {
return .{40, 42}
}
fn main() i32 {
callback Make := make
values := callback()
if values[1] != 42 {
return 1
}
return 0
}
`)
}
