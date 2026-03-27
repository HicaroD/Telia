package sema

import (
	"testing"

	"github.com/HicaroD/Telia/internal/ast"
	"github.com/HicaroD/Telia/internal/diagnostics"
	"github.com/HicaroD/Telia/internal/lexer"
	"github.com/HicaroD/Telia/internal/parser"
)

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
			collector := parseAndCheck(tt.src)
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
			collector := parseAndCheck(tt.src)
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
			collector := parseAndCheck(tt.src)
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
			collector := parseAndCheck(tt.src)
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
			collector := parseAndCheck(tt.src)
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
			collector := parseAndCheck(tt.src)
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
			collector := parseAndCheck(tt.src)
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
			collector := parseAndCheck(tt.src)
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
			collector := parseAndCheck(tt.src)
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
			collector := parseAndCheck(tt.src)
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
			collector := parseAndCheck(tt.src)
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

func TestParseNextDecl(t *testing.T) {
	src := "package main\n\nfn foo() {}\nfn main() {}"
	collector := diagnostics.New()

	loc := &ast.Loc{Name: "test.t"}
	lex := lexer.New(loc, []byte(src), collector)
	p := parser.NewForTest(lex, collector)

	file := &ast.File{
		PkgNameDefined: false,
		Imports:        make(map[string]*ast.UseDecl),
		IsFirstNode:    true,
	}
	pkg := &ast.Package{Scope: ast.NewScope(nil)}
	p.SetFileAndPkg(file, pkg)

	node, done, err := parseNextDecl(p, file)
	if err != nil {
		t.Fatalf("parseNextDecl error: %v", err)
	}
	if done {
		t.Errorf("expected not done after package decl")
	}
	if node == nil {
		t.Errorf("expected node, got nil")
	}
	if node.Kind != ast.KIND_PKG_DECL {
		t.Errorf("expected KIND_PKG_DECL, got %v", node.Kind)
	}

	node, done, err = parseNextDecl(p, file)
	if err != nil {
		t.Fatalf("parseNextDecl error: %v", err)
	}
	if done {
		t.Errorf("expected not done after first decl")
	}
	if node == nil {
		t.Errorf("expected node, got nil")
	}
	if node.Kind != ast.KIND_FN_DECL {
		t.Errorf("expected KIND_FN_DECL, got %v", node.Kind)
	}
	fnDecl := node.Node.(*ast.FnDecl)
	if string(fnDecl.Name.Lexeme) != "foo" {
		t.Errorf("expected 'foo', got %s", fnDecl.Name.Lexeme)
	}

	node, done, err = parseNextDecl(p, file)
	if err != nil {
		t.Fatalf("parseNextDecl error: %v", err)
	}
	if done {
		t.Errorf("expected not done after second decl")
	}
	if node == nil {
		t.Errorf("expected node, got nil")
	}
	if node.Kind != ast.KIND_FN_DECL {
		t.Errorf("expected KIND_FN_DECL, got %v", node.Kind)
	}
	fnDecl = node.Node.(*ast.FnDecl)
	if string(fnDecl.Name.Lexeme) != "main" {
		t.Errorf("expected 'main', got %s", fnDecl.Name.Lexeme)
	}

	node, done, err = parseNextDecl(p, file)
	if err != nil {
		t.Fatalf("parseNextDecl error: %v", err)
	}
	if !done {
		t.Errorf("expected done after all decls")
	}
	if node != nil {
		t.Errorf("expected nil node, got %v", node)
	}
}

func TestErrorHandling(t *testing.T) {
	tests := []struct {
		name     string
		src      string
		hasError bool
		errMsg   string
	}{
		{
			name: "error constructor and field access",
			src: `package main

fn main() {
  e := error("oops")
  m := e.msg
}`,
			hasError: false,
		},
		{
			name: "error nil comparison !=",
			src: `package main

fn main() {
  e := error("oops")
  if e != nil {
    m := e.msg
  }
}`,
			hasError: false,
		},
		{
			name: "error nil comparison ==",
			src: `package main

fn main() {
  e := error("oops")
  if e == nil {
    m := e.msg
  }
}`,
			hasError: false,
		},
		{
			name: "error nil comparison reversed (nil != err)",
			src: `package main

fn main() {
  e := error("oops")
  if nil != e {
    m := e.msg
  }
}`,
			hasError: false,
		},
		{
			name: "error nil comparison reversed (nil == err)",
			src: `package main

fn main() {
  e := error("oops")
  if nil == e {
    m := e.msg
  }
}`,
			hasError: false,
		},
		{
			name: "@fail on non-error function",
			src: `package main

fn greet() {
}

fn main() {
  greet() @fail
}`,
			hasError: true,
			errMsg:   "@fail requires error-returning function",
		},
		{
			name: "@fail on error-only function",
			src: `package main

fn failOnly() error {
  return error("boom")
}

fn main() {
  failOnly() @fail
}`,
			hasError: false,
		},
		{
			name: "@fail on (i32, error) function",
			src: `package main

fn connect() (i32, error) {
  return 42, nil
}

fn main() {
  connect() @fail
}`,
			hasError: false,
		},
		{
			name: "@fail on error-only function with variable assignment",
			src: `package main

fn something() error {
  return error("something")
}

fn main() {
  a := something() @fail
}`,
			hasError: true,
			errMsg:   "@fail on error-only function",
		},
		{
			name: "@catch on non-error function",
			src: `package main

fn greet() {
}

fn main() {
  greet() @catch err {
  }
}`,
			hasError: true,
			errMsg:   "@catch requires error-returning function",
		},
		{
			name: "@catch on error-only function",
			src: `package main

fn failOnly() error {
  return error("boom")
}

fn main() {
  failOnly() @catch err {
    m := err.msg
  }
}`,
			hasError: false,
		},
		{
			name: "@catch handler block is type-checked",
			src: `package main

fn failOnly() error {
  return error("boom")
}

fn main() {
  failOnly() @catch err {
    x := undefinedVar
  }
}`,
			hasError: true,
			errMsg:   "symbol not found on scope",
		},
		{
			name: "@catch on (i32, error) tuple function",
			src: `package main

fn connect() (i32, error) {
  return 42, nil
}

fn main() {
  connect() @catch err {
    m := err.msg
  }
}`,
			hasError: false,
		},
		{
			name: "@catch handler return type must match enclosing function",
			src: `package main

fn failOnly() error {
  return error("boom")
}

fn main() i32 {
  failOnly() @catch err {
    return "bad"
  }
  return 0
}`,
			hasError: true,
			errMsg:   "cannot use string as i32",
		},
		{
			name: "@catch handler with correct return type",
			src: `package main

fn failOnly() error {
  return error("boom")
}

fn main() i32 {
  failOnly() @catch err {
    return 1
  }
  return 0
}`,
			hasError: false,
		},
		{
			name: "@catch on error-only function with variable assignment",
			src: `package main

fn something() error {
  return error("something")
}

fn main() {
  a := something() @catch err {
  }
}`,
			hasError: true,
			errMsg:   "@catch on error-only function",
		},
		{
			name: "error struct field exposes nested msg access",
			src: `package main

struct Result {
  err error
}

fn main() {
  result := Result.{err: error("oops")}
  msg := result.err.msg
}`,
			hasError: false,
		},
		{
			name: "error parameter accepts error values",
			src: `package main

fn print_error(err error) {
  msg := err.msg
}

fn make_error() error {
  return error("oops")
}

fn main() {
  print_error(make_error())
}`,
			hasError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			diags := parseAndCheck(tt.src)
			if tt.hasError {
				if len(diags.Diags) == 0 {
					t.Fatal("expected errors, got none")
				}
				if tt.errMsg != "" && !containsDiag(diags.Diags, tt.errMsg) {
					t.Errorf("expected error containing %q, got %v", tt.errMsg, diags.Diags)
				}
			} else {
				if len(diags.Diags) > 0 {
					t.Errorf("unexpected errors: %v", diags.Diags)
				}
			}
		})
	}
}
