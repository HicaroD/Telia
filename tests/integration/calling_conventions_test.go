package integration

import "testing"

func TestCallbackCallingConventionMacroComesFromHeader(t *testing.T) {
	runRecordFixture(t, `package main
@[cc_macro="FIXTURE_CC"]
type Compare = fn(left i32, right i32) i32
@[header=%q, default_cc_macro="FIXTURE_CC"]
extern {
@[link_name="fixture_convention_invoke"]
fn invoke(callback Compare) i32
}
@[cc_macro="FIXTURE_CC"]
fn compare(left i32, right i32) i32 {
return left - right
}
fn main() i32 {
if invoke(compare) != 42 {
return 1
}
return 0
}
`)
}
