package integration

import "testing"

func TestNamedExternEnumMembersUseHeaderTags(t *testing.T) {
	runRecordFixture(t, `package main
@[header=%q, default_cc="c"]
extern ffi {
@[link_name="FixtureKind"]
enum Kind {
@[link_name="FIXTURE_KIND_READY"]
Ready
}
@[link_name="fixture_kind"]
fn kind(known bool) Kind
}
fn main() i32 {
if ffi::kind(true) != ffi::Kind::Ready {
return 1
}
return 0
}
`)
}

func TestForeignEnumUsesHeaderTagsAndMatchesUnknownValues(t *testing.T) {
	runRecordFixture(t, `package main
@[header=%q, default_cc="c"]
extern {
@[link_name="FixtureKind"]
enum Kind {
@[link_name="FIXTURE_KIND_READY"]
Ready
}
@[link_name="fixture_kind"]
fn kind(known bool) Kind
}
fn main() i32 {
if kind(true) != Kind::Ready {
return 1
}
result := match kind(false) {
Kind::Ready => 1
_ => 42
}
if result != 42 {
return 2
}
return 0
}
`)
}
