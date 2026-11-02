package sema_test

import (
	"strings"
	"testing"
)

func TestHashKeyDiagnostics(t *testing.T) {
	for _, tc := range []struct{ name, decl, ty, path string }{
		{"slice", "", "i32[]", "value"},
		{"empty array of slices", "", "i32[][0]", "value[element]"},
		{"tuple", "", "(i32, i32[])", "value[1]"},
		{"function", "", "fn() i32", "function values"},
		{"error", "", "error", "general errors"},
		{"nested generic", "struct Box<T> {\nvalue T\n}", "Box<i32[]>", "value.value"},
		{"inactive enum", "enum Choice {\nEmpty\nBad {\nitems i32[]\n}\n}", "Choice", "value::Bad.items"},
		{"foreign record", "@[header=\"unused.h\"]\nextern {\nstruct Foreign {\nx i32\n}\n}", "Foreign", "foreign records"},
	} {
		t.Run(tc.name, func(t *testing.T) {
			diags := checkDiagnostics(
				t,
				"package main\n"+tc.decl+"\nfn hash<T>(value T) u64 {\nreturn __key_hash<T>(value, 0, 0)\n}\nfn main() {\nh := hash(zero<"+tc.ty+">())\n}\n",
			)
			found := false
			for _, d := range diags.Diags {
				if strings.Contains(d.Message, "ineligible hash key") &&
					strings.Contains(d.Message, tc.path) {
					found = true
				}
			}
			if !found {
				t.Fatalf("want hashing diagnostic at %s, got %v", tc.path, diags.Diags)
			}
		})
	}
}
