package sema_test

import "testing"

func TestForeignEnumMatchRequiresFallbackEvenWhenAllNamedMembersAreCovered(t *testing.T) {
	c := checkDiagnostics(t, `package main
@[header="fixture.h", default_cc="c"]
extern {
enum Kind {
Ready
}
}
fn main() {
value := match Kind::Ready {
Kind::Ready => 1
}
}
`)
	if len(c.Diags) != 1 || c.Diags[0].Title != "Foreign enum match requires fallback" || !c.Diags[0].Span.Valid() || c.Diags[0].Help == "" || len(c.Diags[0].Related) != 1 {
		t.Fatalf("want actionable fallback diagnostic, got %v", c.Diags)
	}
}
