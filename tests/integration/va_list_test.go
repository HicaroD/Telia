package integration

import "testing"

func TestBorrowedVaListUsesNativeCParameterABI(t *testing.T) {
	runRecordFixture(t, `package main
use "std::libc"
type Inspect = fn(args libc::c_va_list) i32
@[header=%q, default_cc="c"]
extern {
@[link_name="fixture_va_drive"]
fn drive(callback Inspect) i32
@[link_name="fixture_va_check"]
fn check(args libc::c_va_list) i32
}
fn inspect(args libc::c_va_list) i32 {
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
