package integration

import "testing"

func TestConstCStringAliasesComposeThroughPointerOutputSlots(t *testing.T) {
	runRecordFixture(t, `package main
type Text = const cstring
@[header=%q, default_cc="c"]
extern {
@[link_name="fixture_text_replace"]
fn replace(output *Text)
@[link_name="fixture_text_inspect"]
fn inspect(texts const *Text) i32
}
fn main() i32 {
text Text := "before"
replace(&text)
if inspect(&text) != 1 {
return 1
}
return 0
}
`)
}

func TestConstPointerAliasesProtectOnlyTheImmediatePointee(t *testing.T) {
	runRecordFixture(t, `package main
type Mutable = *i32
type Inspect = const *Mutable
fn inspect(slot Inspect) {
pointer := *slot
*pointer = 42
}
fn main() i32 {
value i32 := 0
pointer := &value
inspect(&pointer)
if value != 42 {
return 1
}
return 0
}
`)
}
