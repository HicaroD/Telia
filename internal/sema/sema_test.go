package sema

import (
	"strings"
	"testing"

	"github.com/HicaroD/Telia/internal/ast"
	"github.com/HicaroD/Telia/internal/diagnostics"
	"github.com/HicaroD/Telia/internal/lexer"
	"github.com/HicaroD/Telia/internal/parser"
)

func parseAndCheck(t *testing.T, src string) *diagnostics.Collector {
	collector := diagnostics.New()

	loc := &ast.Loc{Name: "test.t"}
	lex := lexer.New(loc, []byte(src), collector)
	p := parser.NewWithLex(lex, collector)

	file, err := p.ParseFileForTest(loc)
	if err != nil {
		diag := diagnostics.Diag{Message: "parse error: " + err.Error()}
		collector.ReportAndSave(diag)
		return collector
	}

	pkg := p.GetPkg()
	pkg.Loc = loc
	pkg.Files = []*ast.File{file}

	prog := &ast.Program{Root: pkg}

	sema := New(collector)
	err = sema.Check(prog, nil)
	if err != nil {
		diag := diagnostics.Diag{Message: err.Error()}
		collector.ReportAndSave(diag)
	}

	return collector
}

func parseNextDecl(p *parser.Parser, file *ast.File) (*ast.Node, bool, error) {
	return p.NextForTest(file)
}

func containsDiag(diags []diagnostics.Diag, substr string) bool {
	for _, d := range diags {
		if strings.Contains(d.Message, substr) {
			return true
		}
	}
	return false
}

func TestTypeInference(t *testing.T) {
	tests := []struct {
		name     string
		src      string
		hasError bool
	}{
		{
			name:     "int literal inference",
			src:      "package main\n\nfn main() {\n  a := 1\n}",
			hasError: false,
		},
		{
			name:     "float literal inference",
			src:      "package main\n\nfn main() {\n  a := 1.5\n}",
			hasError: false,
		},
		{
			name:     "bool literal inference",
			src:      "package main\n\nfn main() {\n  a := true\n}",
			hasError: false,
		},
		{
			name:     "string literal inference",
			src:      "package main\n\nfn main() {\n  a := \"hello\"\n}",
			hasError: false,
		},
		{
			name:     "binary expr int inference",
			src:      "package main\n\nfn main() {\n  a := 1 + 2\n}",
			hasError: false,
		},
		{
			name:     "unary minus int inference",
			src:      "package main\n\nfn main() {\n  a := -5\n}",
			hasError: false,
		},
		{
			name:     "comparison inference",
			src:      "package main\n\nfn main() {\n  a := 1 > 2\n}",
			hasError: false,
		},
		{
			name:     "logical and inference",
			src:      "package main\n\nfn main() {\n  a := true and false\n}",
			hasError: false,
		},
		{
			name:     "logical or inference",
			src:      "package main\n\nfn main() {\n  a := true or false\n}",
			hasError: false,
		},
		{
			name:     "not inference",
			src:      "package main\n\nfn main() {\n  a := not true\n}",
			hasError: false,
		},
		{
			name:     "variable reference inference",
			src:      "package main\n\nfn main() {\n  a := 1\n  b := a\n}",
			hasError: false,
		},
		{
			name:     "function call inference",
			src:      "package main\n\nfn foo() int {\n  return 1\n}\nfn main() {\n  a := foo()\n}",
			hasError: false,
		},
		{
			name:     "explicit type annotation",
			src:      "package main\n\nfn main() {\n  a i32 := 1\n}",
			hasError: false,
		},
		{
			name:     "explicit type annotation mismatch",
			src:      "package main\n\nfn main() {\n  a i32 := 1.5\n}",
			hasError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			collector := parseAndCheck(t, tt.src)
			hasError := len(collector.Diags) > 0
			if hasError != tt.hasError {
				t.Errorf("expected hasError=%v, got=%v. Diags: %v", tt.hasError, hasError, collector.Diags)
			}
		})
	}
}

func TestTypeChecking(t *testing.T) {
	tests := []struct {
		name      string
		src       string
		errSubstr string
	}{
		{
			name:      "binary expr type mismatch",
			src:       "package main\n\nfn main() {\n  a := 1 + 1.5\n}",
			errSubstr: "invalid operands types",
		},
		{
			name:      "int plus string error",
			src:       "package main\n\nfn main() {\n  a := 1 + \"hello\"\n}",
			errSubstr: "invalid operands types",
		},
		{
			name:      "bool plus int error",
			src:       "package main\n\nfn main() {\n  a := true + 1\n}",
			errSubstr: "invalid operands types",
		},
		{
			name:      "division type mismatch",
			src:       "package main\n\nfn main() {\n  a := 1 / true\n}",
			errSubstr: "invalid operands types",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			collector := parseAndCheck(t, tt.src)
			if tt.errSubstr == "" {
				if len(collector.Diags) > 0 {
					t.Errorf("expected no error, got: %v", collector.Diags)
				}
			} else {
				if !containsDiag(collector.Diags, tt.errSubstr) {
					t.Errorf("expected error containing '%s', got: %v", tt.errSubstr, collector.Diags)
				}
			}
		})
	}
}

func TestFunctionCalls(t *testing.T) {
	tests := []struct {
		name      string
		src       string
		errSubstr string
	}{
		{
			name:      "valid function call",
			src:       "package main\n\nfn foo() {}\nfn main() {\n  foo()\n}",
			errSubstr: "",
		},
		{
			name:      "function call with args",
			src:       "package main\n\nfn foo(a i32) {}\nfn main() {\n  foo(1)\n}",
			errSubstr: "",
		},
		{
			name:      "too few arguments",
			src:       "package main\n\nfn foo(a i32, b i32) {}\nfn main() {\n  foo(1)\n}",
			errSubstr: "not enough arguments",
		},
		{
			name:      "too many arguments",
			src:       "package main\n\nfn foo(a i32) {}\nfn main() {\n  foo(1, 2)\n}",
			errSubstr: "not enough arguments",
		},
		{
			name:      "wrong argument type",
			src:       "package main\n\nfn foo(a i32) {}\nfn main() {\n  foo(\"hello\")\n}",
			errSubstr: "cannot use",
		},
		{
			name:      "call non-function",
			src:       "package main\n\nfn main() {\n  a := 1\n  a()\n}",
			errSubstr: "not callable",
		},
		{
			name:      "undefined function",
			src:       "package main\n\nfn main() {\n  foo()\n}",
			errSubstr: "not defined on scope",
		},
		{
			name:      "function with return value",
			src:       "package main\n\nfn foo() int {\n  return 1\n}\nfn main() {\n  a := foo()\n}",
			errSubstr: "",
		},
		{
			name:      "nested function calls",
			src:       "package main\n\nfn foo() int {\n  return 1\n}\nfn bar() int {\n  return foo()\n}\nfn main() {\n  a := bar()\n}",
			errSubstr: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			collector := parseAndCheck(t, tt.src)
			if tt.errSubstr == "" {
				if len(collector.Diags) > 0 {
					t.Errorf("expected no error, got: %v", collector.Diags)
				}
			} else {
				if !containsDiag(collector.Diags, tt.errSubstr) {
					t.Errorf("expected error containing '%s', got: %v", tt.errSubstr, collector.Diags)
				}
			}
		})
	}
}

func TestScopeResolution(t *testing.T) {
	tests := []struct {
		name      string
		src       string
		errSubstr string
	}{
		{
			name:      "variable in scope",
			src:       "package main\n\nfn main() {\n  a := 1\n  b := a\n}",
			errSubstr: "",
		},
		{
			name:      "undefined variable",
			src:       "package main\n\nfn main() {\n  a := b\n}",
			errSubstr: "symbol not found",
		},
		{
			name:      "shadowing allowed",
			src:       "package main\n\nfn main() {\n  a := 1\n  {\n    a := 2\n  }\n}",
			errSubstr: "",
		},
		{
			name:      "function shadows variable",
			src:       "package main\n\nfn foo() int {\n  return 1\n}\nfn main() {\n  a := foo\n}",
			errSubstr: "not a variable",
		},
		{
			name:      "nested scope variable",
			src:       "package main\n\nfn main() {\n  {\n    a := 1\n  }\n  b := a\n}",
			errSubstr: "symbol not found",
		},
		{
			name:      "use outer scope variable",
			src:       "package main\n\nfn main() {\n  a := 1\n  {\n    b := a\n  }\n}",
			errSubstr: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			collector := parseAndCheck(t, tt.src)
			if tt.errSubstr == "" {
				if len(collector.Diags) > 0 {
					t.Errorf("expected no error, got: %v", collector.Diags)
				}
			} else {
				if !containsDiag(collector.Diags, tt.errSubstr) {
					t.Errorf("expected error containing '%s', got: %v", tt.errSubstr, collector.Diags)
				}
			}
		})
	}
}

func TestDuplicateDetection(t *testing.T) {
	tests := []struct {
		name      string
		src       string
		errSubstr string
	}{
		{
			name:      "duplicate parameters",
			src:       "package main\n\nfn foo(a i32, a i32) {}\nfn main() {}",
			errSubstr: "already declared",
		},
		{
			name:      "no duplicate parameters",
			src:       "package main\n\nfn foo(a i32, b i32) {}\nfn main() {}",
			errSubstr: "",
		},
		{
			name:      "duplicate extern prototypes",
			src:       "package main\nextern libc {\n  fn puts()\n}\nextern libc {\n  fn puts()\n}\nfn main() {}",
			errSubstr: "already declared",
		},
		{
			name:      "duplicate extern declarations",
			src:       "package main\nextern libc {}\nextern libc {}",
			errSubstr: "already declared",
		},
		{
			name:      "duplicate function",
			src:       "package main\n\nfn foo() {}\nfn foo() {}",
			errSubstr: "already declared",
		},
		{
			name:      "duplicate variable in block",
			src:       "package main\n\nfn main() {\n  a := 1\n  a := 2\n}",
			errSubstr: "already declared",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			collector := parseAndCheck(t, tt.src)
			if tt.errSubstr == "" {
				if len(collector.Diags) > 0 {
					t.Errorf("expected no error, got: %v", collector.Diags)
				}
			} else {
				if !containsDiag(collector.Diags, tt.errSubstr) {
					t.Errorf("expected error containing '%s', got: %v", tt.errSubstr, collector.Diags)
				}
			}
		})
	}
}

func TestSemanticErrors(t *testing.T) {
	tests := []struct {
		name      string
		src       string
		errSubstr string
	}{
		{
			name:      "extern function not found",
			src:       "package main\nextern libc {}\nfn main() {\n  libc::puts()\n}",
			errSubstr: "invalid calling convention",
		},
		{
			name:      "extern not defined",
			src:       "package main\n\nfn main() {\n  libc::printf()\n}",
			errSubstr: "symbol not found",
		},
		{
			name:      "multiple variable declaration",
			src:       "package main\n\nfn main() {\n  a, b := 10, 20\n}",
			errSubstr: "",
		},
		{
			name:      "multiple assignment",
			src:       "package main\n\nfn main() {\n  a := 1\n  b := 2\n  a, b = 10, 10\n}",
			errSubstr: "",
		},
		{
			name:      "multiple assignment undeclared left",
			src:       "package main\n\nfn main() {\n  a := 1\n  a, b = 10, 10\n}",
			errSubstr: "not declared",
		},
		{
			name:      "multiple assignment no new vars",
			src:       "package main\n\nfn main() {\n  a := 1\n  b := 2\n  a, b := 10, 10\n}\n",
			errSubstr: "already declared",
		},
		{
			name:      "struct field access",
			src:       "package main\n\nstruct Point {\n  x i32\n  y i32\n}\nfn main() {\n  p := Point.{x: 1, y: 2}\n  q := p.x\n}\n",
			errSubstr: "",
		},
		{
			name:      "struct field access undefined",
			src:       "package main\n\nstruct Point {\n  x i32\n}\nfn main() {\n  p := Point.{x: 1}\n  q := p.z\n}\n",
			errSubstr: "not found",
		},
		{
			name:      "return type mismatch",
			src:       "package main\n\nfn foo() i32 {\n  return 1.5\n}\n",
			errSubstr: "cannot use",
		},
		{
			name:      "return from void function",
			src:       "package main\n\nfn foo() {\n  return 1\n}\n",
			errSubstr: "cannot use",
		},
		{
			name:      "missing return value",
			src:       "package main\n\nfn foo() i32 {}\n",
			errSubstr: "must always return",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			collector := parseAndCheck(t, tt.src)
			if tt.errSubstr == "" {
				if len(collector.Diags) > 0 {
					t.Errorf("expected no error, got: %v", collector.Diags)
				}
			} else {
				if !containsDiag(collector.Diags, tt.errSubstr) {
					t.Errorf("expected error containing '%s', got: %v", tt.errSubstr, collector.Diags)
				}
			}
		})
	}
}

func TestControlFlow(t *testing.T) {
	tests := []struct {
		name      string
		src       string
		errSubstr string
	}{
		{
			name:      "if statement",
			src:       "package main\n\nfn main() {\n  if true {}\n}",
			errSubstr: "",
		},
		{
			name:      "if else",
			src:       "package main\n\nfn main() {\n  if true {} else {}\n}",
			errSubstr: "",
		},
		{
			name:      "if elif else",
			src:       "package main\n\nfn main() {\n  if true {} elif false {} else {}\n}",
			errSubstr: "",
		},
		{
			name:      "if condition type error",
			src:       "package main\n\nfn main() {\n  if 1 {}\n}",
			errSubstr: "type",
		},
		{
			name:      "for loop",
			src:       "package main\n\nfn main() {\n  for i := 0; i < 10; i = i + 1 {}\n}",
			errSubstr: "",
		},
		{
			name:      "while loop",
			src:       "package main\n\nfn main() {\n  while true {}\n}",
			errSubstr: "",
		},
		{
			name:      "for loop condition type error",
			src:       "package main\n\nfn main() {\n  for i := 0; i; i = i + 1 {}\n}",
			errSubstr: "type",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			collector := parseAndCheck(t, tt.src)
			if tt.errSubstr == "" {
				if len(collector.Diags) > 0 {
					t.Errorf("expected no error, got: %v", collector.Diags)
				}
			} else {
				if !containsDiag(collector.Diags, tt.errSubstr) {
					t.Errorf("expected error containing '%s', got: %v", tt.errSubstr, collector.Diags)
				}
			}
		})
	}
}

func TestStructDecl(t *testing.T) {
	tests := []struct {
		name      string
		src       string
		errSubstr string
	}{
		{
			name:      "simple struct",
			src:       "package main\n\nstruct Point {\n  x i32\n  y i32\n}\nfn main() {}",
			errSubstr: "",
		},
		{
			name:      "duplicate struct fields",
			src:       "package main\n\nstruct Point {\n  x i32\n  x i32\n}\nfn main() {}",
			errSubstr: "duplicate",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			collector := parseAndCheck(t, tt.src)
			if tt.errSubstr == "" {
				if len(collector.Diags) > 0 {
					t.Errorf("expected no error, got: %v", collector.Diags)
				}
			} else {
				if !containsDiag(collector.Diags, tt.errSubstr) {
					t.Errorf("expected error containing '%s', got: %v", tt.errSubstr, collector.Diags)
				}
			}
		})
	}
}

func TestMainFunction(t *testing.T) {
	tests := []struct {
		name      string
		src       string
		errSubstr string
	}{
		{
			name:      "has main function",
			src:       "package main\n\nfn main() {}",
			errSubstr: "",
		},
		{
			name:      "no main function",
			src:       "package main\n\nfn foo() {}",
			errSubstr: "main",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			collector := parseAndCheck(t, tt.src)
			if tt.errSubstr == "" {
				if len(collector.Diags) > 0 {
					t.Errorf("expected no error, got: %v", collector.Diags)
				}
			} else {
				if !containsDiag(collector.Diags, tt.errSubstr) {
					t.Errorf("expected error containing '%s', got: %v", tt.errSubstr, collector.Diags)
				}
			}
		})
	}
}

func TestVariadicFunctions(t *testing.T) {
	tests := []struct {
		name      string
		src       string
		errSubstr string
	}{
		{
			name:      "variadic function call",
			src:       "package main\n\nfn foo(args ...i32) {}\nfn main() {\n  foo(1, 2, 3)\n}",
			errSubstr: "",
		},
		{
			name:      "variadic function no args",
			src:       "package main\n\nfn foo(args ...i32) {}\nfn main() {\n  foo()\n}",
			errSubstr: "",
		},
		{
			name:      "variadic wrong type",
			src:       "package main\n\nfn foo(args ...i32) {}\nfn main() {\n  foo(\"hello\")\n}\n",
			errSubstr: "cannot use",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			collector := parseAndCheck(t, tt.src)
			if tt.errSubstr == "" {
				if len(collector.Diags) > 0 {
					t.Errorf("expected no error, got: %v", collector.Diags)
				}
			} else {
				if !containsDiag(collector.Diags, tt.errSubstr) {
					t.Errorf("expected error containing '%s', got: %v", tt.errSubstr, collector.Diags)
				}
			}
		})
	}
}

func TestPointerOperations(t *testing.T) {
	tests := []struct {
		name      string
		src       string
		errSubstr string
	}{
		{
			name:      "pointer dereference",
			src:       "package main\n\nfn main() {\n  a := 1\n  p := &a\n  b := *p\n}\n",
			errSubstr: "",
		},
		{
			name:      "pointer type mismatch",
			src:       "package main\n\nfn main() {\n  a := 1\n  p := &a\n  b := *p + 1.5\n}\n",
			errSubstr: "cannot use",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			collector := parseAndCheck(t, tt.src)
			if tt.errSubstr == "" {
				if len(collector.Diags) > 0 {
					t.Errorf("expected no error, got: %v", collector.Diags)
				}
			} else {
				if !containsDiag(collector.Diags, tt.errSubstr) {
					t.Errorf("expected error containing '%s', got: %v", tt.errSubstr, collector.Diags)
				}
			}
		})
	}
}
