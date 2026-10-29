package parser

import (
	"testing"

	"github.com/HicaroD/Telia/internal/ast"
)

func TestNamedLiteralWithoutDot(t *testing.T) {
	node, err := ParseExprFrom("Point{x: 42}", "literal.t")
	if err != nil {
		t.Fatal(err)
	}
	if node.Kind != ast.KIND_STRUCT_EXPR {
		t.Fatalf("expected struct literal, got %v", node.Kind)
	}
	lit := node.Node.(*ast.StructLiteralExpr)
	if lit.Name.Name() != "Point" || len(lit.Values) != 1 || lit.Values[0].Name.Name() != "x" {
		t.Fatalf("literal name and fields were not preserved: %+v", lit)
	}
}

func TestExistingDottedLiteralsAndArraysRemainValid(t *testing.T) {
	for _, input := range []string{
		"Point.{x: 42}",
		"Box<i32>.{value: 42}",
		"Option<i32>::Some.{value: 42}",
		"match event {\nOption<i32>::Some.{value} => value\n}",
		".{Point{x: 42}}",
		".{}",
	} {
		t.Run(input, func(t *testing.T) {
			_, err := ParseDeclFrom("fn test() {\nvalue := "+input+"\n}\n", "compatibility.t")
			if err != nil {
				t.Fatal(err)
			}
		})
	}
}

func TestDirectNamedLiteralsInHeadersNeedParentheses(t *testing.T) {
	for _, body := range []string{
		"if point == Point{x: 42} {}",
		"while point == Point{x: 42} {}",
		"value := match Option<i32>::Some{value: 42} {\nOption<i32>::Some{value} => value\n}",
	} {
		t.Run(body, func(t *testing.T) {
			if _, err := ParseDeclFrom("fn test() {\n"+body+"\n}\n", "header.t"); err == nil {
				t.Fatal("expected ambiguous literal to require parentheses")
			}
		})
	}
}

func TestControlFlowBracesAreNotLiterals(t *testing.T) {
	for _, body := range []string{
		"if ready {}",
		"if ready {} elif other {} else {}",
		"while ready {}",
		"for item in items {}",
		"for i := 0; i < 3; i = next {}",
		"value := match event {\nEvent::Quit => 0\n}",
		"value := match Event::Quit {\nEvent::Quit => 0\n}",
		"if a < b > c {}",
	} {
		t.Run(body, func(t *testing.T) {
			_, err := ParseDeclFrom("fn test() {\n"+body+"\n}\n", "header.t")
			if err != nil {
				t.Fatal(err)
			}
		})
	}
}

func TestGenericAndQualifiedLiteralsWithoutDot(t *testing.T) {
	for _, input := range []string{
		"Box<i32>{value: 42}",
		"Box<Box<i32>>{value: Box<i32>{value: 42}}",
		"pkg::Point{x: 42}",
		"Option<i32>::Some{value: 42}",
		"pkg::Option<i32>::Some{value: 42}",
	} {
		t.Run(input, func(t *testing.T) {
			_, err := ParseDeclFrom("fn test() {\nvalue := "+input+"\n}\n", "literal.t")
			if err != nil {
				t.Fatal(err)
			}
		})
	}
}

func TestPayloadPatternWithoutDot(t *testing.T) {
	_, err := ParseDeclFrom(`fn test() {
value := match event {
Option<i32>::Some{value: result, ..} => result
Option<i32>::None => 0
}
}
`, "pattern.t")
	if err != nil {
		t.Fatal(err)
	}
}

func TestDelimitedLiteralsInControlFlowHeaders(t *testing.T) {
	for _, body := range []string{
		"if point == (Point{x: 42}) {}",
		"if point == Point.{nested: Point{x: 42}} {}",
		"if point == (Box<i32>{value: 42}) {}",
		"if contains(Point{x: 42}) {}",
		"if contains<Box<i32>>(Box<i32>{value: 42}) {}",
		"if flags[Point{x: 42}] {}",
		"while contains(Point{x: 42}) {}",
		"if maybe() @catch err {\npoint := Point{x: 42}\nreturn\n} {}",
		"if match event {\nEvent::Quit => check(Point{x: 42})\n} {}",
		"if match event {\nEvent::Quit => {\npoint := Point{x: 42}\nyield true\n}\n} {}",
		"if match event {\nEvent::Quit => Point{x: 42}\n} == point {}",
		"for item in .{Point{x: 42}} {}",
		"for item in values(Point{x: 42}) {}",
		"value := match (Option<i32>::Some{value: 42}) {\nOption<i32>::Some{value} => value\n}",
	} {
		t.Run(body, func(t *testing.T) {
			_, err := ParseDeclFrom("fn test() {\n"+body+"\n}\n", "header.t")
			if err != nil {
				t.Fatal(err)
			}
		})
	}
}
