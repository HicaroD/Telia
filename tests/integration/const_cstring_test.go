package integration

import "testing"

func TestConstCStringCanBeCheckedForNullWithoutDiscardingConst(t *testing.T) {
	runRecordFixture(t, `package main
use "std::ffi::c"
fn main() i32 {
  missing c::ConstCString := nil
  text c::ConstCString := "present"
  if missing != nil or text == nil {
    return 1
  }
  return 0
}
`)
}

func TestConstCStringCanBeExplicitlyViewedAsTeliaText(t *testing.T) {
	runRecordFixture(t, `package main
use "std::strings"
use "std::ffi::c"
fn main() i32 {
text c::ConstCString := "hello"
view := strings::from_c(text)
if strings::len(view) != 5 {
return 1
}
return 0
}
`)
}

func TestConstCStringPreservesExplicitStringConversions(t *testing.T) {
	runRecordFixture(t, `package main
use "std::ffi::c"
@[header=%q, default_cc="c"]
extern {
@[link_name="fixture_text_check"]
fn check(text c::ConstCString) i32
}
fn main() i32 {
if check("hello") != 1 {
return 1
}
text c::CString := "hello"
if check(text) != 1 {
return 2
}
return 0
}
`)
}
