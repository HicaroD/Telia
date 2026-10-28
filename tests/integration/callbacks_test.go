package integration

import "testing"

func TestCallbackValuesCallBareFunctions(t *testing.T) {
	runRecordFixture(t, `package main
type Compare = fn(left i32, right i32) i32
fn compare(left i32, right i32) i32 {
return left - right
}
fn invoke(callback Compare) i32 {
return callback(44, 2)
}
fn main() i32 {
callback Compare := compare
if callback(44, 2) != 42 or invoke(compare) != 42 {
return 1
}
return 0
}
`)
}
