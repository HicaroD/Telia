package integration

import (
	"testing"
)

func TestForeignEnumRepresentationCannotOverrideHeaderABI(t *testing.T) {
	exe, diags := compileRecordFixture(t, `package main
@[header=%q, default_cc="c"]
extern {
@[link_name="FixtureKind"]
enum Kind u8 {
@[link_name="FIXTURE_KIND_READY"]
Ready
}
}
fn main() {}
`)
	if exe != "" {
		defer removeTestBinary(t, exe)
	}
	if len(diags.Diags) != 1 || diags.Diags[0].Title != "Foreign enum representation mismatch" ||
		!diags.Diags[0].Span.Valid() ||
		diags.Diags[0].Help == "" {
		t.Fatalf("want enum ABI diagnostic, got %v", diags.Diags)
	}
}
