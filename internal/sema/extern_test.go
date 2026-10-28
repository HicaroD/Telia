package sema

import "testing"

func TestExternRejectsUndefinedNestedParameterType(t *testing.T) {
	diags := parseAndCheck("package main\n@[default_cc=\"c\"]\nextern {\nfn foreign(value **Missing)\n}\nfn main() {}\n")
	if !containsDiag(diags.Diags, "'Missing' type not found") {
		t.Fatalf("expected undefined nested type diagnostic, got %v", diags.Diags)
	}
}

func TestExternRejectsUndefinedNestedReturnType(t *testing.T) {
	diags := parseAndCheck("package main\n@[default_cc=\"c\"]\nextern {\nfn foreign() **Missing\n}\nfn main() {}\n")
	if !containsDiag(diags.Diags, "'Missing' type not found") {
		t.Fatalf("expected undefined nested return type diagnostic, got %v", diags.Diags)
	}
}

func TestExternNestedTypesPreserveAliasCycleDiagnostics(t *testing.T) {
	for _, signature := range []string{"fn foreign(value **Loop)", "fn foreign() **Loop"} {
		diags := parseAndCheck("package main\ntype Loop = *Loop\n@[default_cc=\"c\"]\nextern {\n" + signature + "\n}\nfn main() {}\n")
		if !containsDiag(diags.Diags, "cyclic type alias 'Loop'") {
			t.Fatalf("expected alias cycle diagnostic for %s, got %v", signature, diags.Diags)
		}
	}
}
