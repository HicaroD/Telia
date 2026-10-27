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
