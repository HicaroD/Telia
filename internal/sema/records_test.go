package sema_test

import "testing"

func TestForeignMultidimensionalArrayFieldsDiagnoseUnsupportedDeclarators(t *testing.T) {
	c := checkDiagnostics(
		t,
		"package main\n@[header=\"fixture.h\", default_cc=\"c\"]\nextern {\nstruct Matrix {\ndata u8[3][2]\n}\n}\nfn main() {}\n",
	)
	if len(c.Diags) != 1 || c.Diags[0].Title != "Unsupported foreign declarator" ||
		!c.Diags[0].Span.Valid() ||
		c.Diags[0].Help == "" {
		t.Fatalf("want multidimensional array diagnostic, got %v", c.Diags)
	}
}

func TestForeignArrayAddressesRequireExplicitElementPointers(t *testing.T) {
	c := checkDiagnostics(
		t,
		"package main\n@[header=\"fixture.h\", default_cc=\"c\"]\nextern {\nstruct Bytes {\ndata u8[4]\n}\n}\nfn main() {\nvalue := Bytes{}\npointer := &value.data\n}\n",
	)
	if len(c.Diags) != 1 || c.Diags[0].Title != "Unsupported foreign declarator" ||
		!c.Diags[0].Span.Valid() ||
		c.Diags[0].Help == "" {
		t.Fatalf("want explicit element pointer diagnostic, got %v", c.Diags)
	}
}

func TestExternPointersCannotPretendWrapperArraysAreCArrayPointers(t *testing.T) {
	c := checkDiagnostics(
		t,
		"package main\n@[header=\"fixture.h\", default_cc=\"c\"]\nextern {\nfn send(values *u8[4])\n}\nfn main() {}\n",
	)
	if len(c.Diags) != 1 || c.Diags[0].Title != "Unsupported foreign declarator" ||
		!c.Diags[0].Span.Valid() ||
		c.Diags[0].Help == "" {
		t.Fatalf("want array-pointer ABI diagnostic, got %v", c.Diags)
	}
}

func TestForeignArrayFieldsCannotBeAssignedAsWholeArrays(t *testing.T) {
	c := checkDiagnostics(
		t,
		"package main\n@[header=\"fixture.h\", default_cc=\"c\"]\nextern {\nstruct Bytes {\ndata u8[4]\n}\n}\nfn main() {\nvalue := Bytes{}\nvalue.data = .{1, 2, 3, 4}\n}\n",
	)
	if len(c.Diags) != 1 || c.Diags[0].Title != "Foreign array copy" || !c.Diags[0].Span.Valid() ||
		c.Diags[0].Help == "" {
		t.Fatalf("want whole-array assignment diagnostic, got %v", c.Diags)
	}
}

func TestExternArrayParametersRequireExplicitElementPointers(t *testing.T) {
	c := checkDiagnostics(
		t,
		"package main\n@[header=\"fixture.h\", default_cc=\"c\"]\nextern {\nfn send(values u8[4])\n}\nfn main() {}\n",
	)
	if len(c.Diags) != 1 || c.Diags[0].Title != "Foreign array parameter" ||
		!c.Diags[0].Span.Valid() ||
		c.Diags[0].Help == "" {
		t.Fatalf("want explicit pointer diagnostic, got %v", c.Diags)
	}
}

func TestForeignArrayFieldsCannotBeCopiedAsWrapperArrays(t *testing.T) {
	c := checkDiagnostics(
		t,
		"package main\n@[header=\"fixture.h\", default_cc=\"c\"]\nextern {\nstruct Bytes {\ndata u8[4]\n}\n}\nfn main() {\nvalue := Bytes{}\ncopy := value.data\n}\n",
	)
	if len(c.Diags) != 1 || c.Diags[0].Title != "Foreign array copy" || !c.Diags[0].Span.Valid() ||
		c.Diags[0].Help == "" {
		t.Fatalf("want foreign array storage diagnostic, got %v", c.Diags)
	}
}

func TestForeignUnionLiteralRequiresExactlyOneMember(t *testing.T) {
	for _, literal := range []string{"Value{}", "Value{integer: 42, real: 1.0}"} {
		t.Run(literal, func(t *testing.T) {
			c := checkDiagnostics(
				t,
				"package main\n@[header=\"fixture.h\", default_cc=\"c\"]\nextern {\nunion Value {\ninteger i32\nreal f64\n}\n}\nfn main() {\nvalue := "+literal+"\n}\n",
			)
			if len(c.Diags) != 1 || c.Diags[0].Title != "Union member selection" ||
				!c.Diags[0].Span.Valid() ||
				c.Diags[0].Help == "" ||
				len(c.Diags[0].Related) != 1 {
				t.Fatalf("want actionable union diagnostic, got %v", c.Diags)
			}
		})
	}
}
