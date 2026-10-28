package parser_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/HicaroD/Telia/internal/ast"
	"github.com/HicaroD/Telia/internal/diagnostics"
	"github.com/HicaroD/Telia/internal/parser"
)

func TestIndependentSyntaxErrorsAreCollected(t *testing.T) {
	path := filepath.Join(t.TempDir(), "main.t")
	text := "package main\nfn main() {\n  first: = 1\n  second: = 2\n}\n"
	if err := os.WriteFile(path, []byte(text), 0600); err != nil {
		t.Fatal(err)
	}
	loc, err := ast.LocFromPath(path)
	if err != nil {
		t.Fatal(err)
	}
	c := diagnostics.New()
	_, err = parser.New(c).ParseFileAsProgram(path, loc, c)
	if err == nil {
		t.Fatal("invalid program succeeded")
	}
	if len(c.Diags) != 2 {
		t.Fatalf("want two independent errors, got %v", c.Diags)
	}
	for i, line := range []int{3, 4} {
		s := c.Diags[i].Span
		if !s.Valid() {
			t.Fatalf("missing span: %+v", c.Diags[i])
		}
		got, _ := s.File.Position(s.Start)
		if got != line {
			t.Fatalf("want line %d, got %d", line, got)
		}
	}
}

func TestDuplicateFunctionHasRelatedDeclaration(t *testing.T) {
	path := filepath.Join(t.TempDir(), "main.t")
	if err := os.WriteFile(path, []byte("package main\nfn same() {}\nfn same() {}\n"), 0600); err != nil {
		t.Fatal(err)
	}
	loc, err := ast.LocFromPath(path)
	if err != nil {
		t.Fatal(err)
	}
	c := diagnostics.New()
	_, err = parser.New(c).ParseFileAsProgram(path, loc, c)
	if err == nil || len(c.Diags) != 1 || len(c.Diags[0].Related) != 1 {
		t.Fatalf("want duplicate with related location, got %v, %+v", err, c.Diags)
	}
	d := c.Diags[0]
	if !d.Span.Valid() || !d.Related[0].Span.Valid() {
		t.Fatal("missing spans")
	}
	line, _ := d.Related[0].Span.File.Position(d.Related[0].Span.Start)
	if line != 2 {
		t.Fatalf("original declaration on line %d, want 2", line)
	}
}

func TestConstantDeclarationExplainsIncorrectAssignmentOperator(t *testing.T) {
	path := filepath.Join(t.TempDir(), "main.t")
	if err := os.WriteFile(path, []byte("package main\nconst Answer i32 := 42\nfn main() {}\n"), 0600); err != nil {
		t.Fatal(err)
	}
	loc, err := ast.LocFromPath(path)
	if err != nil {
		t.Fatal(err)
	}
	c := diagnostics.New()
	_, err = parser.New(c).ParseFileAsProgram(path, loc, c)
	if err == nil || len(c.Diags) != 1 {
		t.Fatalf("want one invalid constant diagnostic, got %v, %+v", err, c.Diags)
	}
	d := c.Diags[0]
	if d.Title != "Expected =" || d.Help == "" || !d.Span.Valid() ||
		d.Span.File.Text[d.Span.Start:d.Span.End] != ":=" {
		t.Fatalf("want a diagnostic highlighting := with a correction, got %+v", d)
	}
}

func TestPackageCollectsErrorsAcrossFiles(t *testing.T) {
	dir := t.TempDir()
	for _, name := range []string{"first.t", "second.t"} {
		if err := os.WriteFile(filepath.Join(dir, name), []byte("package main\nfn broken() {\n  value: = 1\n}\n"), 0600); err != nil {
			t.Fatal(err)
		}
	}
	loc, err := ast.LocFromPath(dir)
	if err != nil {
		t.Fatal(err)
	}
	c := diagnostics.New()
	_, err = parser.New(c).ParsePackageAsProgram(dir, loc)
	if err == nil || len(c.Diags) != 2 {
		t.Fatalf("want two file errors and failure, got %v, %v", err, c.Diags)
	}
	if c.Diags[0].Span.File == c.Diags[1].Span.File {
		t.Fatal("errors attributed to same file")
	}
}

func TestRecoveryHandlesEOFNestedBlocksAndLexicalErrors(t *testing.T) {
	for _, tc := range []struct {
		name, text string
		count      int
	}{
		{"EOF", "package main\nfn main() {\n", 1},
		{"nested", "package main\nfn main() {\n if true {\n  bad: = 1\n }\n other: = 2\n}\n", 2},
		{"lexical", "package main\nfn main() {\n $\n !\n}\n", 2},
		{"return lexical", "package main\nfn main() {\n return $\n}\n", 1},
		{"parameter lexical", "package main\nfn main(value $) {}\n", 1},
	} {
		t.Run(tc.name, func(t *testing.T) {
			path := filepath.Join(t.TempDir(), "main.t")
			if err := os.WriteFile(path, []byte(tc.text), 0600); err != nil {
				t.Fatal(err)
			}
			loc, err := ast.LocFromPath(path)
			if err != nil {
				t.Fatal(err)
			}
			c := diagnostics.New()
			_, err = parser.New(c).ParseFileAsProgram(path, loc, c)
			if err == nil || len(c.Diags) != tc.count {
				t.Fatalf("want %d errors, got %v, %v", tc.count, err, c.Diags)
			}
		})
	}
}

func TestNestedExpressionsRetainPreciseSpans(t *testing.T) {
	n, err := parser.ParseExprFrom("1 + 2 * 3", "expr.t")
	if err != nil {
		t.Fatal(err)
	}
	binary := n.Node.(*ast.BinExpr)
	for _, tc := range []struct {
		node     *ast.Node
		spelling string
	}{{binary.Left, "1"}, {binary.Right, "2 * 3"}, {binary.Right.Node.(*ast.BinExpr).Right, "3"}} {
		s := tc.node.Span
		if !s.Valid() || s.File.Text[s.Start:s.End] != tc.spelling {
			t.Fatalf("want %q span, got %+v", tc.spelling, s)
		}
	}
}
