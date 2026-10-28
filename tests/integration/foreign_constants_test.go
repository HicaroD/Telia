package integration

import "testing"

func TestNamedExternConstantUsesItsHeaderValue(t *testing.T) {
	runRecordFixture(t, `package main
@[header=%q, default_cc="c"]
extern ffi {
@[link_name="FIXTURE_RECORD_FLAG"]
const Flag u64
}
fn main() i32 {
if ffi::Flag != 1099511627776 {
return 1
}
return 0
}
`)
}

func TestForeignConstantUsesHeaderValue(t *testing.T) {
	runRecordFixture(t, `package main
@[header=%q, default_cc="c"]
extern {
@[link_name="FIXTURE_RECORD_FLAG"]
const Flag u64
}
fn main() i32 {
flags := Flag | 1
if flags != 1099511627777 {
return 1
}
return 0
}
`)
}
