package sema_test

import "testing"

func TestOpaqueTypesCannotBePassedByValue(t *testing.T) {
	c := checkDiagnostics(
		t,
		"package main\n@[header=\"fixture.h\", default_cc=\"c\"]\nextern {\nopaque Window\nfn consume(value Window)\n}\nfn main() {}\n",
	)
	if len(c.Diags) != 1 || c.Diags[0].Title != "Incomplete opaque value" ||
		!c.Diags[0].Span.Valid() ||
		c.Diags[0].Help == "" ||
		len(c.Diags[0].Related) != 1 {
		t.Fatalf("want actionable incomplete-type diagnostic, got %v", c.Diags)
	}
}

func TestOpaqueTypesCannotBeSized(t *testing.T) {
	c := checkDiagnostics(
		t,
		"package main\n@[header=\"fixture.h\", default_cc=\"c\"]\nextern {\nopaque Window\n}\nfn main() {\nsize := sizeof<Window>()\n}\n",
	)
	if len(c.Diags) != 1 || c.Diags[0].Title != "Incomplete opaque value" ||
		!c.Diags[0].Span.Valid() ||
		c.Diags[0].Help == "" ||
		len(c.Diags[0].Related) != 1 {
		t.Fatalf("want actionable incomplete-type diagnostic, got %v", c.Diags)
	}
}

func TestOpaquePointersCannotBeDereferencedIntoValues(t *testing.T) {
	c := checkDiagnostics(
		t,
		"package main\n@[header=\"fixture.h\", default_cc=\"c\"]\nextern {\nopaque Window\n}\nfn main() {\nwindow *Window := nil\nvalue := *window\n}\n",
	)
	if len(c.Diags) != 1 || c.Diags[0].Title != "Incomplete opaque value" ||
		!c.Diags[0].Span.Valid() ||
		c.Diags[0].Help == "" ||
		len(c.Diags[0].Related) != 1 {
		t.Fatalf("want actionable incomplete-type diagnostic, got %v", c.Diags)
	}
}

func TestOpaquePointersDoNotExposeHeaderFields(t *testing.T) {
	c := checkDiagnostics(
		t,
		"package main\n@[header=\"fixture.h\", default_cc=\"c\"]\nextern {\nopaque Window\n}\nfn main() {\nwindow *Window := nil\nvalue := window.width\n}\n",
	)
	if len(c.Diags) != 1 || c.Diags[0].Title != "Opaque field access" || !c.Diags[0].Span.Valid() ||
		c.Diags[0].Help == "" ||
		len(c.Diags[0].Related) != 1 {
		t.Fatalf("want actionable opaque-field diagnostic, got %v", c.Diags)
	}
}
