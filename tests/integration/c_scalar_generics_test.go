package integration

import "testing"

func TestCScalarsSupportUnaryArithmeticAndGenericSpecializations(t *testing.T) {
	runRecordFixture(t, `package main
use "std::libc"
fn negate<T>(value libc::c_long, ignored T) libc::c_long {
return -value
}
fn invoke<T>(value T) T {
return value
}
fn main() i32 {
value libc::c_long := 42
negative := invoke<libc::c_long>(negate<bool>(value, true))
if negative != -42 {
return 1
}
return 0
}
`)
}
