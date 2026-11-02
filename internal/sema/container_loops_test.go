package sema_test

import (
	"strings"
	"testing"
)

func TestRangeBindingsDiagnoseCountsAndLookalikes(t *testing.T) {
	for _, tc := range []struct{ decl, body, want string }{
		{"", "values i32[1] := .{1}\nfor a, b in values {}", "requires 1 binding(s), got 2"},
		{"", "values := zero<i32[]>()\nfor a, b in values {}", "requires 1 binding(s), got 2"},
		{"struct Map<K, V> {\nkey K\nvalue V\n}", "values := Map<i32, i32>{}\nfor key, value in values {}", "iteration requires an array"},
		{"struct Set<T> {\nvalue T\n}", "values := Set<i32>{}\nfor value in values {}", "iteration requires an array"},
	} {
		t.Run(tc.body, func(t *testing.T) {
			c := checkDiagnostics(t, "package main\n"+tc.decl+"\nfn main() {\n"+tc.body+"\n}\n")
			if len(c.Diags) != 1 || !strings.Contains(c.Diags[0].Message, tc.want) ||
				!c.Diags[0].Span.Valid() {
				t.Fatalf("want located %q, got %v", tc.want, c.Diags)
			}
		})
	}
}
