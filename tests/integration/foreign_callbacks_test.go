package integration

import "testing"

func TestCallbacksCrossCBoundaryAndUseTypedOutputSlots(t *testing.T) {
	runRecordFixture(t, `package main
type Compare = fn(left i32, right i32) i32
@[header=%q, default_cc="c"]
extern {
@[link_name="fixture_callback_invoke"]
fn invoke(callback Compare) i32
@[link_name="fixture_callback_get"]
fn get(output *Compare)
}
fn compare(left i32, right i32) i32 {
return left - right
}
fn main() i32 {
if invoke(compare) != 42 {
return 1
}
callback Compare := nil
if callback != nil {
return 3
}
get(&callback)
if callback == nil {
return 4
}
if callback(40, 2) != 42 {
return 2
}
return 0
}
`)
}
