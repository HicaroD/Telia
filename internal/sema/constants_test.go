package sema

import "testing"

func TestConstantsCannotBeReassigned(t *testing.T) {
	diags := parseAndCheck("package main\nconst Answer = 42\nfn main() {\nAnswer = 43\n}\n")
	if !containsDiag(diags.Diags, "cannot reassign constant 'Answer'") {
		t.Fatalf("expected constant reassignment diagnostic, got %v", diags.Diags)
	}
}
