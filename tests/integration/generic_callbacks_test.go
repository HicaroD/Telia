package integration

import "testing"

func TestGenericCallbackSignaturesSpecializeIndependently(t *testing.T) {
	runRecordFixture(t, `package main
fn identity32(value i32) i32 {
return value
}
fn identity64(value i64) i64 {
return value
}
fn apply<T>(callback fn(value T) T, value T) T {
return callback(value)
}
fn main() i32 {
if apply<i32>(identity32, 42) != 42 {
return 1
}
if apply<i64>(identity64, 99) != 99 {
return 2
}
return 0
}
`)
}
