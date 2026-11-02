package sema_test

import (
	"strings"
	"testing"
)

func TestEqualityRejectsUnsupportedNestedComponents(t *testing.T) {
	for _, tc := range []struct{ name, decl, body, path string }{
		{"slice", "", "a := zero<i32[]>()\nsame := a == a", "value"},
		{"nested alias", "type Values = i32[]\nstruct Inner {\nvalues Values\n}\nstruct Outer<T> {\ninner T\n}", "a := zero<Outer<Inner>>()\nsame := a != a", "value.inner.values"},
		{"array", "", "a := zero<i32[][0]>()\nsame := a == a", "value[element]"},
		{"tuple", "", "same := zero<(i32, i32[])>() == zero<(i32, i32[])>()", "value[1]"},
		{"error", "", "a := error(\"failure\")\nsame := a == a", "value"},
		{"error field", "struct Box {\nerr error\n}", "a := Box{}\nsame := a == a", "value.err"},
		{"foreign record", "@[header=\"unused.h\"]\nextern {\nstruct Foreign {\nx i32\n}\n}", "a := Foreign{}\nsame := a == a", "value"},
		{"foreign union", "@[header=\"unused.h\"]\nextern {\nunion Foreign {\nx i32\n}\n}", "a := Foreign{x: 1}\nsame := a == a", "value"},
		{"foreign inline array", "@[header=\"unused.h\"]\nextern {\nstruct Foreign {\nx i32[2]\n}\n}", "a := Foreign{}\nsame := a.x == a.x", "value"},
	} {
		t.Run(tc.name, func(t *testing.T) {
			c := checkDiagnostics(t, "package main\n"+tc.decl+"\nfn main() {\n"+tc.body+"\n}\n")
			if len(c.Diags) != 1 || c.Diags[0].Title != "Unsupported equality" ||
				!strings.Contains(c.Diags[0].Message, tc.path) ||
				!c.Diags[0].Span.Valid() ||
				c.Diags[0].Help == "" {
				t.Fatalf("want component equality diagnostic at %s, got %v", tc.path, c.Diags)
			}
		})
	}
}

func TestGenericTemplateRejectsIndependentWholeErrorEquality(t *testing.T) {
	c := checkDiagnostics(
		t,
		"package main\nfn unused<T>(a error, b error) bool {\nreturn a == b\n}\nfn main() {}\n",
	)
	if len(c.Diags) != 1 || !strings.Contains(c.Diags[0].Message, "whole-value equality") {
		t.Fatalf("want unsupported error equality, got %v", c.Diags)
	}
}
