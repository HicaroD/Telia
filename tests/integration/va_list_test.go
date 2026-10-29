package integration

import "testing"

func TestBorrowedVaListUsesNativeCParameterABI(t *testing.T) {
	runRecordFixture(t, `package main
use "std::ffi::c"
@[cc="c"]
type Inspect = fn(args c::c_va_list) i32
@[header=%q, default_cc="c"]
extern {
@[link_name="fixture_va_drive"]
fn drive(callback Inspect) i32
@[link_name="fixture_va_check"]
fn check(args c::c_va_list) i32
}
@[cc="c"]
fn inspect(args c::c_va_list) i32 {
return check(args)
}
fn main() i32 {
if drive(inspect) != 1 {
return 1
}
return 0
}
`)
}
