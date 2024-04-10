package parser

import (
	"fmt"
	"reflect"
	"testing"

	"github.com/HicaroD/telia-lang/ast"
	"github.com/HicaroD/telia-lang/lexer/token/kind"
)

type exprTest struct {
	expr string
	node ast.Expr
}

func TestFunctionDecl(t *testing.T) {
}
func TestExternDecl(t *testing.T) {}

// TODO: build a utility method for parsing expressions for
// testing
func TestLiteralExpr(t *testing.T) {
	filename := "test.tt"
	literals := []exprTest{
		{
			expr: "1",
			node: &ast.LiteralExpr{Value: 1, Kind: kind.INTEGER_LITERAL},
		},
		{
			expr: "-1",
			node: &ast.LiteralExpr{Value: -1, Kind: kind.NEGATIVE_INTEGER_LITERAL},
		},
		{
			expr: "true",
			node: &ast.LiteralExpr{Value: "true", Kind: kind.TRUE_BOOL_LITERAL},
		},
		{
			expr: "false",
			node: &ast.LiteralExpr{Value: "false", Kind: kind.FALSE_BOOL_LITERAL},
		},
	}

	for _, test := range literals {
		t.Run(fmt.Sprintf("TestLiteralExpr('%s')", test.expr), func(t *testing.T) {
			actualNode, err := ParseExprFrom(test.expr, filename)
			if err != nil {
				t.Errorf("TestLiteralExpr('%s'): unexpected error '%v'", test.expr, err)
			}
			if !reflect.DeepEqual(test.node, actualNode) {
				t.Errorf("TestLiteralExpr('%s'): expression node differs\nexpected: '%v', but got '%v'\n", test.expr, test.node, actualNode)
			}
		})
	}
}
func TestFuncCallExpr(t *testing.T) {}
func TestBinaryExpr(t *testing.T) {
	filename := "test.tt"
	binExprs := []exprTest{
		{
			expr: "1 + 1",
			node: &ast.BinaryExpr{
				Left:  &ast.LiteralExpr{Value: 1, Kind: kind.INTEGER_LITERAL},
				Op:    kind.PLUS,
				Right: &ast.LiteralExpr{Value: 1, Kind: kind.INTEGER_LITERAL},
			},
		},
	}

	for _, test := range binExprs {
		t.Run(fmt.Sprintf("TestBinaryExpr('%s')", test.expr), func(t *testing.T) {
			actualNode, err := ParseExprFrom(test.expr, filename)
			if err != nil {
				t.Errorf("TestBinaryExpr('%s'): unexpected error '%v'", test.expr, err)
			}
			if !reflect.DeepEqual(test.node, actualNode) {
				t.Errorf("TestBinaryExpr('%s'): expression node differs\nexpected: '%v', but got '%v'\n", test.expr, test.node, actualNode)
			}
		})
	}
}

func TestFuncCallStmt(t *testing.T) {}
func TestIfStmt(t *testing.T)       {}
