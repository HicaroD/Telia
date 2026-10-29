package parser

import (
	"testing"

	"github.com/HicaroD/Telia/internal/ast"
)

func TestParseGenericStructDeclaration(t *testing.T) {
	node, err := ParseDeclFrom("struct Pair<A, B> {\nfirst A\nsecond B\n}\n", "generic.t")
	if err != nil {
		t.Fatalf("generic declaration should parse: %v", err)
	}
	decl := node.Node.(*ast.StructDecl)
	if len(decl.TypeParams) != 2 || decl.TypeParams[0].Name() != "A" ||
		decl.TypeParams[1].Name() != "B" {
		t.Fatalf("type parameters were not preserved: %v", decl.TypeParams)
	}
	if len(decl.Fields) != 2 || decl.Fields[1].Name.Name() != "second" {
		t.Fatalf("generic fields were not preserved: %v", decl.Fields)
	}
}

func TestGenericParsingPreservesComparisons(t *testing.T) {
	for _, input := range []string{"a < b", "a > b", "a < b && b > c", "a < b > c"} {
		node, err := ParseExprFrom(input, "comparison.t")
		if err != nil || node.Kind != ast.KIND_BINARY_EXPR {
			t.Fatalf("%s should remain a comparison: node %v, error %v", input, node, err)
		}
	}
}

func TestInvalidGenericDeclarations(t *testing.T) {
	for _, input := range []string{
		"struct Pair<T, T> {}\n", "struct Pair<> {}\n",
		"struct Pair<i32> {}\n", "fn f<T,>() {}\n",
	} {
		if _, err := ParseDeclFrom(input, "invalid.t"); err == nil {
			t.Fatalf("expected invalid generic declaration: %s", input)
		}
	}
}

func TestNestedGenericLiteralsPreserveMultilineParsing(t *testing.T) {
	_, err := ParseExprFrom(`Box<Box<i32>>{
    value: Box<i32>{
        value: 42
    }
}`, "nested.t")
	if err != nil {
		t.Fatalf("nested multiline literal should parse: %v", err)
	}
}

func TestParseGenericComposition(t *testing.T) {
	_, err := ParseDeclFrom(`fn wrap<T>(value T) pkg::Pair<T, *pkg::Box<(i32, string)>> {
    copy := pkg::identity<T>(value)
    return pkg::Pair<T, *pkg::Box<(i32, string)>>{ first: copy }
}
`, "generic.t")
	if err != nil {
		t.Fatalf("nested generic types, calls and literals should parse: %v", err)
	}
}
