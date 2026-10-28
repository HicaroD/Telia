package integration

import "testing"

func TestConstCStringCanBeExplicitlyViewedAsTeliaText(t *testing.T) {
	runRecordFixture(t, `package main
use "std::strings"
fn main() i32 {
text const cstring := "hello"
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
@[header=%q, default_cc="c"]
extern {
@[link_name="fixture_text_check"]
fn check(text const cstring) i32
}
fn main() i32 {
if check("hello") != 1 {
return 1
}
text cstring := "hello"
if check(text) != 1 {
return 2
}
return 0
}
`)
}
