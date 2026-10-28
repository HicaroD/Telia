package integration

import (
	"os"
	"path/filepath"
	"testing"
)

func TestHeaderOwnedRecordFieldTypesAreCheckedAgainstHeader(t *testing.T) {
	exe, diags := compileRecordFixture(t, `package main
@[header=%q, default_cc="c"]
extern {
@[link_name="FixtureRecord"]
struct Record {
value f64
}
}
fn main() {
record := Record.{value: 42.0}
}
`)
	if exe != "" {
		defer os.RemoveAll(filepath.Dir(exe))
	}
	if len(diags.Diags) != 1 || diags.Diags[0].Title != "Foreign field type mismatch" || !diags.Diags[0].Span.Valid() || diags.Diags[0].Help == "" || len(diags.Diags[0].Related) != 1 {
		t.Fatalf("want header field type diagnostic, got %v", diags.Diags)
	}
	span := diags.Diags[0].Span
	if span.File.Text[span.Start:span.End] != "value" {
		t.Fatalf("wrong field span: %+v", span)
	}
}
