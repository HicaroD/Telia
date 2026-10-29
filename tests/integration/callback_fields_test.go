package integration

import "testing"

func TestForeignCallbackFieldsAndReturnedCallbacksAreCallable(t *testing.T) {
	runRecordFixture(t, `package main
@[cc="c"]
type Compare = fn(left i32, right i32) i32
@[header=%q, default_cc="c"]
extern {
@[link_name="FixtureCallbacks"]
struct Callbacks {
compare Compare
}
@[link_name="fixture_callback_return"]
fn get() Compare
}
fn main() i32 {
callback := get()
holder := Callbacks.{compare: callback}
if holder.compare(40, 2) != 42 {
return 1
}
return 0
}
`)
}
