package integration

import "testing"

func TestImportedBareFunctionNamesAreCallbackValues(t *testing.T) {
	runRecordFixture(t, `package main
use "std::libc"
type Allocate = fn(size libc::c_size_t) rawptr
fn main() i32 {
allocate Allocate := libc::malloc
pointer := allocate(1)
if pointer == nil {
return 1
}
libc::free(pointer)
return 0
}
`)
}
