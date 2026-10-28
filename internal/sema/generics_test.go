package sema

import (
	"fmt"
	"strings"
	"testing"

	"github.com/HicaroD/Telia/internal/ast"
	"github.com/HicaroD/Telia/internal/diagnostics"
	"github.com/HicaroD/Telia/internal/lexer"
	"github.com/HicaroD/Telia/internal/parser"
)

func TestGenericInstantiationDiagnostics(t *testing.T) {
	for _, tc := range []struct{ name, source, want string }{
		{"operation", `fn double<T>(value T) T {
return value + value
}
fn main() {
value := double<string>("bad")
}
`, "instantiated double<string>"},
		{"explicit arguments", `fn identity<T>(value T) T {
return value
}
fn main() {
value := identity(1)
}
`, "requires 1 explicit type arguments, got 0"},
		{"wrong struct arity", "struct Box<T> {\nvalue T\n}\nfn main() {\nx := Box<i32, string>.{}\n}\n", "requires 1 type arguments, got 2"},
		{"unapplied type", "struct Box<T> {\nvalue T\n}\nfn main() {\nx := Box<Box>.{}\n}\n", "requires 1 type arguments, got 0"},
		{"literal argument", "struct Box<T> {\nvalue T\n}\nfn main() {\nx := Box<42>.{}\n}\n", "invalid generic type argument"},
		{"nil argument", "struct Box<T> {\nvalue T\n}\nfn main() {\nx := Box<nil>.{}\n}\n", "invalid generic type argument"},
		{"by value recursion", "struct Node<T> {\nnext Node<T>\n}\nfn main() {\nx := Node<i32>.{}\n}\n", "infinitely sized by-value layout"},
		{"specialization growth", "fn grow<T>() {\ngrow<*T>()\n}\nfn main() {\ngrow<i32>()\n}\n", "instantiation depth limit (128)"},
		{"recursive alias argument", "struct Box<T> {}\ntype Loop = Box<Loop>\nfn main() {\nx := Box<Loop>.{}\n}\n", "cyclic type alias 'Loop'"},
	} {
		t.Run(tc.name, func(t *testing.T) {
			diags := parseAndCheck("package main\n" + tc.source)
			if !containsDiag(diags.Diags, tc.want) {
				t.Fatalf("want diagnostic %q, got %v", tc.want, diags.Diags)
			}
			for _, d := range diags.Diags {
				if strings.Contains(d.Message, "C compiler error") {
					t.Fatalf("expected a Telia diagnostic: %s", d.Message)
				}
			}
		})
	}
}

func TestGenericSpecializationBudget(t *testing.T) {
	var source strings.Builder
	source.WriteString("package main\nfn mark<T>() {}\n")
	for i := 0; i < 4097; i++ {
		fmt.Fprintf(&source, "struct Tag%d {}\n", i)
	}
	source.WriteString("fn main() {\n")
	for i := 0; i < 4097; i++ {
		fmt.Fprintf(&source, "mark<Tag%d>()\n", i)
	}
	source.WriteString("}\n")
	diags := parseAndCheck(source.String())
	if !containsDiag(diags.Diags, "specialization limit (4096)") {
		t.Fatalf("expected specialization budget diagnostic, got %v", diags.Diags)
	}
}

func TestGenericCallResultMustMatchItsContext(t *testing.T) {
	diags := parseAndCheck(`package main
fn identity<T>(value T) T {
return value
}
fn wrong() string {
return identity<i32>(42)
}
fn main() {}
`)
	if !containsDiag(diags.Diags, "type mismatch") {
		t.Fatalf("expected a Telia type mismatch, got %v", diags.Diags)
	}
}

func TestUnusedTemplatesRejectIndependentErrors(t *testing.T) {
	for _, tc := range []struct{ source, want string }{
		{"struct Box<T> {\nvalue Missing\n}\n", "'Missing' type not found"},
		{"struct Box<T> {\nvalue T\nvalue T\n}\n", "duplicate field"},
		{"fn bad<T>(value T) T {\nreturn missing\n}\n", "'missing' not defined"},
		{"fn bad<T>() {\nvalue := 1 + \"bad\"\n}\n", "invalid operand types"},
		{"fn bad<T>() {\nvalue := not 1\n}\n", "invalid operand type"},
		{"fn helper(value i32) {}\nfn bad<T>() {\nhelper()\n}\n", "requires 1 value arguments, got 0"},
		{"struct Box<T> {}\nfn bad<T>(value Box<T, i32>) {}\n", "requires 1 type arguments, got 2"},
	} {
		diags := parseAndCheck("package main\n" + tc.source + "fn main() {}\n")
		if !containsDiag(diags.Diags, tc.want) {
			t.Fatalf("want %q for %s, got %v", tc.want, tc.source, diags.Diags)
		}
	}
}

func TestGenericLoweringReusesAliasesAndPreservesTemplates(t *testing.T) {
	source := `package main
type Number = i32
struct Box<T> {
value T
}
fn box<T>(value T) Box<T> {
copy := value
return Box<T>.{ value: copy }
}
fn main() {
a := box<i32>(1)
b := box<Number>(2)
c := box<string>("hello")
}
`
	collector := diagnostics.New()
	loc := &ast.Loc{Name: "test.t"}
	p := parser.NewForTest(lexer.New(loc, []byte(source), collector), collector)
	file, err := p.ParseFileForTest(loc)
	if err != nil {
		t.Fatal(err)
	}
	pkg := p.GetPkg()
	pkg.Loc, pkg.Files = loc, []*ast.File{file}
	program := &ast.Program{Root: pkg}
	if err := New(collector).Check(program); err != nil {
		t.Fatal(err)
	}
	if len(program.Functions) != 2 || len(program.Structs) != 2 {
		t.Fatalf(
			"want two concrete functions and layouts, got %d and %d",
			len(program.Functions),
			len(program.Structs),
		)
	}
	templateNode, _ := pkg.Scope.LookupCurrentScope("box")
	template := templateNode.Node.(*ast.FnDecl)
	if template.Params.Fields[0].Type.Kind != ast.EXPR_TYPE_ID ||
		template.RetType.Kind != ast.EXPR_TYPE_ID {
		t.Fatal("specializing changed the template signature")
	}
	literal := template.Block.Statements[1].Node.(*ast.ReturnStmt).Value.Node.(*ast.StructLiteralExpr)
	if literal.Decl != nil || len(template.Scope.Nodes) != 1 {
		t.Fatal("specializing populated template semantic state")
	}
	mainNode, _ := pkg.Scope.LookupCurrentScope("main")
	main := mainNode.Node.(*ast.FnDecl)
	first := main.Block.Statements[0].Node.(*ast.VarStmt).Expr.Node.(*ast.FnCall).Decl
	second := main.Block.Statements[1].Node.(*ast.VarStmt).Expr.Node.(*ast.FnCall).Decl
	third := main.Block.Statements[2].Node.(*ast.VarStmt).Expr.Node.(*ast.FnCall).Decl
	if first != second || first == third {
		t.Fatal(
			"canonical aliases must reuse an instance while different arguments stay independent",
		)
	}
	if first.Scope == third.Scope || first.Block == third.Block {
		t.Fatal("concrete functions share mutable semantic state")
	}
}

func TestExternFunctionsRejectGenericArguments(t *testing.T) {
	for _, call := range []string{"foreign<i32>(1)", "ffi::foreign<i32>(1)"} {
		extern := "extern {\nfn foreign(value i32) i32\n}\n"
		if strings.HasPrefix(call, "ffi") {
			extern = "extern ffi {\nfn foreign(value i32) i32\n}\n"
		}
		diags := parseAndCheck(
			"package main\n@[default_cc=\"c\"]\n" + extern + "fn main() {\nvalue := " + call + "\n}\n",
		)
		if !containsDiag(diags.Diags, "does not take type arguments") {
			t.Fatalf("expected concrete C ABI diagnostic, got %v", diags.Diags)
		}
	}
}

func TestGenericDiagnosticsIncludeTemplateLocationAndInstantiationChain(t *testing.T) {
	diags := parseAndCheck(`package main
fn inner<T>(value T) T {
return value + value
}
fn outer<T>(value T) T {
return inner<T>(value)
}
fn main() {
value := outer<string>("bad")
}
`)
	var message strings.Builder
	for _, d := range diags.Diags {
		message.WriteString(d.Message)
	}
	for _, want := range []string{"[:test.t:2:4]", "instantiated inner<string>", "instantiated outer<string>"} {
		if !strings.Contains(message.String(), want) {
			t.Fatalf("want %q in diagnostic, got %s", want, message.String())
		}
	}
}

func TestGenericLayoutFailureIncludesInstantiation(t *testing.T) {
	diags := parseAndCheck(
		"package main\nstruct Node<T> {\nnext Node<T>\n}\nfn main() {\nx := Node<i32>.{}\n}\n",
	)
	if !containsDiag(diags.Diags, "instantiated Node<i32>") {
		t.Fatalf("want layout instantiation context, got %v", diags.Diags)
	}
}
