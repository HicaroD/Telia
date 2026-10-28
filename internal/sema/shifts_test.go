package sema_test

import "testing"

func TestInvalidConstantShiftCountsHaveWidthAndCountLocations(t *testing.T) {
	for _, count := range []string{"-1", "8", "128"} {
		t.Run(count, func(t *testing.T) {
			c := checkDiagnostics(
				t,
				"package main\nconst Invalid u8 = 1 << "+count+"\nfn main() {}\n",
			)
			if len(c.Diags) != 1 || c.Diags[0].Title != "Invalid shift count" ||
				!c.Diags[0].Span.Valid() ||
				c.Diags[0].Help == "" ||
				len(c.Diags[0].Related) != 1 ||
				!c.Diags[0].Related[0].Span.Valid() {
				t.Fatalf("want structured shift-count diagnostic, got %v", c.Diags)
			}
		})
	}
}
