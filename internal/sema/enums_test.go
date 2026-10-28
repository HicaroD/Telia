package sema_test

import "testing"

func TestEnumRejectsRepeatedTagsWithBothLocations(t *testing.T) {
	c := checkDiagnostics(
		t,
		"package main\nenum Event u8 {\nQuit = 7\nResize = 7\n}\nfn main() {}\n",
	)
	if len(c.Diags) != 1 {
		t.Fatalf("want one repeated-tag diagnostic, got %v", c.Diags)
	}
	d := c.Diags[0]
	if d.Title != "Duplicate enum tag" || !d.Span.Valid() || d.Help == "" || len(d.Related) != 1 ||
		!d.Related[0].Span.Valid() {
		t.Fatalf("want structured duplicate-tag diagnostic, got %+v", d)
	}
}

func TestEnumRejectsTypeArgumentsOnNonGenericDeclaration(t *testing.T) {
	c := checkDiagnostics(
		t,
		"package main\nenum Event {\nQuit\n}\nfn main() {\nvalue Event<i32> := Event::Quit\n}\n",
	)
	if len(c.Diags) != 1 || c.Diags[0].Title != "Invalid enum type arguments" ||
		!c.Diags[0].Span.Valid() ||
		c.Diags[0].Help == "" ||
		len(c.Diags[0].Related) != 1 {
		t.Fatalf("want structured enum arity diagnostic, got %v", c.Diags)
	}
}

func TestEnumRejectsCyclicRepresentation(t *testing.T) {
	c := checkDiagnostics(t, "package main\nenum Event Event {\nQuit\n}\nfn main() {}\n")
	if len(c.Diags) != 1 {
		t.Fatalf("want one cycle diagnostic, got %v", c.Diags)
	}
	d := c.Diags[0]
	if d.Title != "Cyclic enum declaration" || !d.Span.Valid() || d.Help == "" ||
		len(d.Related) != 1 {
		t.Fatalf("want structured cycle diagnostic, got %+v", d)
	}
}

func TestEnumTypesRemainDistinctFromTheirRepresentation(t *testing.T) {
	c := checkDiagnostics(
		t,
		"package main\nenum Event u8 {\nQuit\n}\nfn main() {\nvalue u8 := Event::Quit\n}\n",
	)
	if len(c.Diags) != 1 {
		t.Fatalf("want one enum type mismatch, got %v", c.Diags)
	}
	d := c.Diags[0]
	if d.Title != "Enum type mismatch" || !d.Span.Valid() || d.Help == "" || len(d.Related) != 1 {
		t.Fatalf("want structured nominal-type diagnostic, got %+v", d)
	}
}

func TestEnumTagValidationRegressions(t *testing.T) {
	for _, tc := range []struct{ name, declaration, title string }{
		{"explicit overflow", "enum Event u8 {\nQuit = 256\n}", "Constant overflow"},
		{"implicit overflow", "enum Event u8 {\nQuit = 255\nResize\n}", "Constant overflow"},
		{"negative unsigned", "enum Event u8 {\nQuit = -1\n}", "Constant overflow"},
		{"noninteger representation", "enum Event bool {\nQuit\n}", "Invalid enum representation"},
		{"repeated name", "enum Event {\nQuit\nQuit\n}", "Duplicate enum variant"},
	} {
		t.Run(tc.name, func(t *testing.T) {
			c := checkDiagnostics(t, "package main\n"+tc.declaration+"\nfn main() {}\n")
			if len(c.Diags) != 1 || c.Diags[0].Title != tc.title || !c.Diags[0].Span.Valid() ||
				c.Diags[0].Help == "" {
				t.Fatalf("want structured %s diagnostic, got %v", tc.title, c.Diags)
			}
		})
	}
}

func TestPayloadEnumsRequireMatchingInsteadOfTagOnlyEquality(t *testing.T) {
	c := checkDiagnostics(
		t,
		"package main\nenum Event {\nResize {\nwidth i32\n}\n}\nfn main() {\nsame := Event::Resize.{width: 1} == Event::Resize.{width: 2}\n}\n",
	)
	if len(c.Diags) != 1 || c.Diags[0].Title != "Payload enum comparison" ||
		!c.Diags[0].Span.Valid() ||
		c.Diags[0].Help == "" ||
		len(c.Diags[0].Related) != 1 {
		t.Fatalf("want structured payload-comparison diagnostic, got %v", c.Diags)
	}
}

func TestInvalidPayloadConstructionHasActionableDiagnostics(t *testing.T) {
	for _, tc := range []struct{ expr, title string }{
		{"Event::Resize", "Missing payload field"},
		{"Event::Resize.{width: 1}", "Missing payload field"},
		{"Event::Resize.{width: 1, height: 2, depth: 3}", "Unknown payload field"},
		{"Event::Resize.{width: 1, width: 2, height: 3}", "Duplicate payload initializer"},
	} {
		t.Run(tc.title+tc.expr, func(t *testing.T) {
			c := checkDiagnostics(
				t,
				"package main\nenum Event {\nResize {\nwidth i32\nheight i32\n}\n}\nfn main() {\nvalue := "+tc.expr+"\n}\n",
			)
			if len(c.Diags) != 1 || c.Diags[0].Title != tc.title || !c.Diags[0].Span.Valid() ||
				c.Diags[0].Help == "" ||
				len(c.Diags[0].Related) != 1 {
				t.Fatalf("want structured %s diagnostic, got %v", tc.title, c.Diags)
			}
		})
	}
}
