package integration

import "testing"

func TestCVariadicsApplyDefaultIntegerAndFloatPromotions(t *testing.T) {
	runRecordFixture(t, `package main
@[header=%q, default_cc="c"]
extern {
@[link_name="fixture_variadic_check"]
fn check(marker i32, args @c ...) i32
}
fn main() i32 {
small i8 := -7
unsigned u16 := 60000
fraction f32 := 1.5
result i32 := check(42, true, small, unsigned, fraction)
if result != 1 {
return 1
}
return 0
}
`)
}
