package sema_test

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/HicaroD/Telia/internal/ast"
	"github.com/HicaroD/Telia/internal/diagnostics"
	"github.com/HicaroD/Telia/internal/parser"
	"github.com/HicaroD/Telia/internal/sema"
)

func checkDiagnostics(t *testing.T, text string) *diagnostics.Collector {
	t.Helper()
	path := filepath.Join(t.TempDir(), "main.t")
	if err := os.WriteFile(path, []byte(text), 0600); err != nil {
		t.Fatal(err)
	}
	loc, err := ast.LocFromPath(path)
	if err != nil {
		t.Fatal(err)
	}
	c := diagnostics.New()
	p, err := parser.New(c).ParseFileAsProgram(path, loc, c)
	if err != nil {
		t.Fatalf("parse: %v, %v", err, c.Diags)
	}
	if err := sema.New(c).Check(p); err == nil {
		t.Fatal("invalid program succeeded")
	}
	return c
}

func TestIndependentTypeErrorsAreCollected(t *testing.T) {
	c := checkDiagnostics(t, "package main\nfn main() {\n  age i32 := \"twenty\"\n  count i32 := \"many\"\n}\nfn other() {\n  size i32 := \"large\"\n}\n")
	if len(c.Diags) != 3 {
		t.Fatalf("want three independent errors, got %v", c.Diags)
	}
	for i, want := range []string{"\"twenty\"", "\"many\"", "\"large\""} {
		if c.Diags[i].Title != "Type mismatch" {
			t.Fatalf("missing concise error title: %+v", c.Diags[i])
		}
		s := c.Diags[i].Span
		if !s.Valid() || s.File.Text[s.Start:s.End] != want {
			t.Fatalf("want highlight %s, got %+v", want, s)
		}
	}
}

func TestFailedInferenceSuppressesDependentErrors(t *testing.T) {
	c := checkDiagnostics(t, "package main\nfn main() {\n  broken := missing\n  dependent := broken + 1\n  chained := dependent + 1\n  independent i32 := \"bad\"\n}\n")
	if len(c.Diags) != 2 {
		t.Fatalf("want only root errors, got %v", c.Diags)
	}
	if !strings.Contains(c.Diags[0].Message, "missing") {
		t.Fatalf("unknown name omitted: %s", c.Diags[0].Message)
	}
}

func TestDuplicateVariablePointsToBothDeclarations(t *testing.T) {
	c := checkDiagnostics(t, "package main\nfn main() {\n  value := 1\n  value := 2\n}\n")
	if len(c.Diags) != 1 {
		t.Fatalf("want one duplicate error, got %v", c.Diags)
	}
	d := c.Diags[0]
	if len(d.Related) != 1 || !d.Span.Valid() || !d.Related[0].Span.Valid() {
		t.Fatalf("missing declaration locations: %+v", d)
	}
	line, _ := d.Span.File.Position(d.Span.Start)
	previous, _ := d.Related[0].Span.File.Position(d.Related[0].Span.Start)
	if line != 4 || previous != 3 {
		t.Fatalf("want duplicate on 4, original on 3; got %d, %d", line, previous)
	}
}

func TestGenericErrorRetainsExpressionLocation(t *testing.T) {
	c := checkDiagnostics(t, "package main\nfn double<T>(value T) T {\n  return value + value\n}\nfn main() {\n  result := double<string>(\"bad\")\n}\n")
	if len(c.Diags) != 1 {
		t.Fatalf("want one generic error, got %v", c.Diags)
	}
	s := c.Diags[0].Span
	if !s.Valid() || s.File.Text[s.Start:s.End] != "value + value" {
		t.Fatalf("missing generic expression location: %+v", s)
	}
}

func TestInvalidValuesDoNotCrashOtherExpressionForms(t *testing.T) {
	for _, expr := range []string{"broken.field", "broken[0]", "*broken", "&broken", "broken()"} {
		t.Run(expr, func(t *testing.T) {
			c := checkDiagnostics(t, "package main\nfn main() {\n  broken := missing\n  next := "+expr+"\n  independent i32 := \"bad\"\n}\n")
			if len(c.Diags) != 2 {
				t.Fatalf("want root errors only, got %v", c.Diags)
			}
		})
	}
}

func TestTemplateErrorsDoNotHideIndependentDeclarations(t *testing.T) {
	c := checkDiagnostics(t, "package main\nfn first<T>(value Missing) {}\nfn second<T>(value Unknown) {}\nfn main() {\n  independent i32 := \"bad\"\n}\n")
	if len(c.Diags) != 3 {
		t.Fatalf("want two template errors and one ordinary error, got %v", c.Diags)
	}
}

func TestDuplicateStructFieldHasRelatedLocation(t *testing.T) {
	c := checkDiagnostics(t, "package main\nstruct Point {\n  x i32\n  x i32\n}\nfn main() {}\n")
	if len(c.Diags) != 1 || len(c.Diags[0].Related) != 1 {
		t.Fatalf("want duplicate with related field, got %+v", c.Diags)
	}
}

func TestConstantReassignmentPointsToUseAndDeclaration(t *testing.T) {
	c := checkDiagnostics(t, "package main\nconst Answer = 42\nfn main() {\n  Answer = 43\n}\n")
	if len(c.Diags) != 1 {
		t.Fatalf("want one reassignment diagnostic, got %+v", c.Diags)
	}
	d := c.Diags[0]
	if d.Title != "Constant reassignment" || d.Help == "" || len(d.Related) != 1 {
		t.Fatalf("missing structured constant diagnostic: %+v", d)
	}
	if !d.Span.Valid() || d.Span.File.Text[d.Span.Start:d.Span.End] != "Answer" {
		t.Fatalf("want constant name highlighted: %+v", d.Span)
	}
	line, _ := d.Span.File.Position(d.Span.Start)
	previous, _ := d.Related[0].Span.File.Position(d.Related[0].Span.Start)
	if line != 4 || previous != 2 {
		t.Fatalf("want reassignment on 4 and declaration on 2, got %d and %d", line, previous)
	}
}

func TestConstantOverflowHasAValueHighlight(t *testing.T) {
	c := checkDiagnostics(t, "package main\nconst TooLarge u8 = 256\nfn main() {}\n")
	if len(c.Diags) != 1 {
		t.Fatalf("want one overflow diagnostic, got %+v", c.Diags)
	}
	d := c.Diags[0]
	if d.Title != "Constant overflow" || d.Help == "" || !d.Span.Valid() || d.Span.File.Text[d.Span.Start:d.Span.End] != "256" {
		t.Fatalf("want structured overflow diagnostic highlighting 256, got %+v", d)
	}
}

func TestConstantUseChecksTheDestinationRange(t *testing.T) {
	c := checkDiagnostics(t, "package main\nconst Count = 256\nfn main() {\n  count u8 := Count\n}\n")
	if len(c.Diags) != 1 {
		t.Fatalf("want one constant overflow, got %+v", c.Diags)
	}
	d := c.Diags[0]
	if d.Title != "Constant overflow" || !d.Span.Valid() || d.Span.File.Text[d.Span.Start:d.Span.End] != "Count" {
		t.Fatalf("want overflow at the constant use, got %+v", d)
	}
}

func TestRuntimeConstantInitializerExplainsCompileTimeRequirement(t *testing.T) {
	c := checkDiagnostics(t, "package main\nfn get() i32 {\nreturn 42\n}\nconst Answer = get()\nfn main() {}\n")
	if len(c.Diags) != 1 {
		t.Fatalf("want one nonconstant initializer diagnostic, got %+v", c.Diags)
	}
	d := c.Diags[0]
	if d.Title != "Nonconstant initializer" || d.Help == "" || !d.Span.Valid() || d.Span.File.Text[d.Span.Start:d.Span.End] != "get()" {
		t.Fatalf("want structured initializer diagnostic highlighting get(), got %+v", d)
	}
}

func TestBitwiseOperandsExplainTheIntegerRequirement(t *testing.T) {
	c := checkDiagnostics(t, "package main\nfn main() {\n  value := true | false\n}\n")
	if len(c.Diags) != 1 {
		t.Fatalf("want one invalid operand diagnostic, got %+v", c.Diags)
	}
	d := c.Diags[0]
	if d.Title != "Invalid bitwise operands" || d.Help == "" || !d.Span.Valid() || d.Span.File.Text[d.Span.Start:d.Span.End] != "true | false" {
		t.Fatalf("want structured bitwise diagnostic, got %+v", d)
	}
}

func TestOtherInvalidBitwiseOperandsHaveStructuredDiagnostics(t *testing.T) {
	for _, body := range []string{"value := ~1.5", "a u32 := 1\nb u64 := 2\nvalue := a & b"} {
		c := checkDiagnostics(t, "package main\nfn main() {\n"+body+"\n}\n")
		if len(c.Diags) != 1 || c.Diags[0].Title != "Invalid bitwise operands" || c.Diags[0].Help == "" || !c.Diags[0].Span.Valid() {
			t.Fatalf("want structured bitwise diagnostic for %s, got %+v", body, c.Diags)
		}
	}
}

func TestUndefinedNestedExternTypeHighlightsTheTypeName(t *testing.T) {
	c := checkDiagnostics(t, "package main\n@[default_cc=\"c\"]\nextern {\nfn foreign(value **Missing)\n}\nfn main() {}\n")
	if len(c.Diags) != 1 {
		t.Fatalf("want one undefined type diagnostic, got %+v", c.Diags)
	}
	d := c.Diags[0]
	if d.Title != "Unknown type" || d.Help == "" || !d.Span.Valid() || d.Span.File.Text[d.Span.Start:d.Span.End] != "Missing" {
		t.Fatalf("want a structured diagnostic at the nested type name, got %+v", d)
	}
}

func TestConstantsDoNotExposeMutableStorage(t *testing.T) {
	c := checkDiagnostics(t, "package main\nconst Answer = 42\nfn main() {\n  pointer := &Answer\n}\n")
	if len(c.Diags) != 1 || c.Diags[0].Title != "Not addressable" || !c.Diags[0].Span.Valid() || c.Diags[0].Span.File.Text[c.Diags[0].Span.Start:c.Diags[0].Span.End] != "Answer" {
		t.Fatalf("want a structured addressability diagnostic on Answer, got %+v", c.Diags)
	}
}

func TestFloatingConstantOverflowIsDiagnosedBeforeCodeGeneration(t *testing.T) {
	c := checkDiagnostics(t, "package main\nconst Huge f32 = 9999999999999999999999999999999999999999.0\nfn main() {}\n")
	if len(c.Diags) != 1 || c.Diags[0].Title != "Constant overflow" || !c.Diags[0].Span.Valid() {
		t.Fatalf("want a structured floating-point overflow diagnostic, got %+v", c.Diags)
	}
}

func TestTypedConstantIntermediateOverflowIsNotHidden(t *testing.T) {
	c := checkDiagnostics(t, "package main\nconst Limit u8 = 255\nconst Value u8 = (Limit + 1) - 1\nfn main() {}\n")
	if len(c.Diags) != 1 || c.Diags[0].Title != "Constant overflow" || !c.Diags[0].Span.Valid() {
		t.Fatalf("want intermediate overflow diagnostic, got %+v", c.Diags)
	}
}

func TestUnnamedExternDuplicateParameterHasRelatedLocation(t *testing.T) {
	c := checkDiagnostics(t, "package main\n@[default_cc=\"c\"]\nextern {\nfn foreign(value i32, value i32)\n}\nfn main() {}\n")
	if len(c.Diags) != 1 || c.Diags[0].Title != "Duplicate parameter" || len(c.Diags[0].Related) != 1 || !c.Diags[0].Span.Valid() || !c.Diags[0].Related[0].Span.Valid() {
		t.Fatalf("want structured duplicate parameter diagnostic, got %+v", c.Diags)
	}
}
