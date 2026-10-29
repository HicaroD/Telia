package integration

import "testing"

func TestCScalarsSupportUnaryArithmeticAndGenericSpecializations(t *testing.T) {
	runRecordFixture(t, `package main
use "std::ffi::c"
fn negate<T>(value c::c_long, ignored T) c::c_long {
return -value
}
fn invoke<T>(value T) T {
return value
}
fn main() i32 {
value c::c_long := 42
negative := invoke<c::c_long>(negate<bool>(value, true))
if negative != -42 {
return 1
}
return 0
}
`)
}
