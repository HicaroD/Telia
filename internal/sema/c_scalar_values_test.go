package sema_test

import (
	"strings"
	"testing"

	"github.com/HicaroD/Telia/internal/diagnostics"
)

func TestCScalarValueCompatibilityDoesNotErasePointerOrCallbackIdentity(t *testing.T) {
	for _, test := range []struct {
		name, source string
	}{
		{"pointer variable", "fn accept(value *CFloat) {}\nfn main() {\nvalue f32 := 1.0\npointer := &value\naccept(pointer)\n}\n"},
		{"address expression", "fn accept(value *CFloat) {}\nfn main() {\nvalue f32 := 1.0\naccept(&value)\n}\n"},
		{"callback", "type Callback = fn(value CFloat) CFloat\nfn native(value f32) f32 {\nreturn value\n}\nfn main() {\ncallback Callback := native\n}\n"},
	} {
		t.Run(test.name, func(t *testing.T) {
			c := checkDiagnostics(t, "package main\ntype CFloat = __c_float\n"+test.source)
			if len(c.Diags) != 1 || !c.Diags[0].Span.Valid() {
				t.Fatalf("want one located identity mismatch, got %+v", c.Diags)
			}
		})
	}
}

func TestCFloatArithmeticNarrowingHasAnActionableDiagnostic(t *testing.T) {
	c := checkDiagnostics(t, "package main\ntype CFloat = __c_float\nfn narrow(value f64) CFloat {\nreturn value / 2.0\n}\nfn main() {}\n")
	if len(c.Diags) != 1 {
		t.Fatalf("want one narrowing error, got %+v", c.Diags)
	}
	d := c.Diags[0]
	if d.Title != "Type mismatch" || !d.Span.Valid() || d.Help == "" || !strings.Contains(d.Message, "f64") || !strings.Contains(d.Message, "float") {
		t.Fatalf("want structured scalar mismatch, got %+v", d)
	}
	if got := d.Span.File.Text[d.Span.Start:d.Span.End]; got != "value / 2.0" {
		t.Fatalf("want complete arithmetic expression highlighted, got %q", got)
	}
	var output strings.Builder
	if err := (diagnostics.Renderer{}).Render(&output, c); err != nil {
		t.Fatal(err)
	}
	for _, text := range []string{"Type mismatch", "return value / 2.0", "matching widths", "Compilation failed · 1 error"} {
		if !strings.Contains(output.String(), text) {
			t.Fatalf("rendered diagnostic lacks %q:\n%s", text, output.String())
		}
	}
}
