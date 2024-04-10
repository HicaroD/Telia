package parser

import (
	"fmt"
	"reflect"
	"testing"

	"github.com/HicaroD/telia-lang/ast"
	"github.com/HicaroD/telia-lang/lexer/token"
	"github.com/HicaroD/telia-lang/lexer/token/kind"
)

type exprTest struct {
	expr string
	node ast.Expr
}

func TestFunctionDecl(t *testing.T) {}

func TestExternDecl(t *testing.T) {}

func TestLiteralExpr(t *testing.T) {
	filename := "test.tt"
	literals := []exprTest{
		{
			expr: "1",
			node: &ast.LiteralExpr{Value: 1, Kind: kind.INTEGER_LITERAL},
		},
		{
			expr: "true",
			node: &ast.LiteralExpr{Value: "true", Kind: kind.TRUE_BOOL_LITERAL},
		},
		{
			expr: "false",
			node: &ast.LiteralExpr{Value: "false", Kind: kind.FALSE_BOOL_LITERAL},
		},
		{
			expr: "\"Hello, world\"",
			node: &ast.LiteralExpr{Value: "Hello, world", Kind: kind.STRING_LITERAL},
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

func TestUnaryExpr(t *testing.T) {
	filename := "test.tt"
	unaryExprs := []exprTest{
		{
			expr: "-1",
			node: &ast.UnaryExpr{
				Op:   kind.MINUS,
				Node: &ast.LiteralExpr{Value: 1, Kind: kind.INTEGER_LITERAL},
			},
		},
		{
			expr: "not true",
			node: &ast.UnaryExpr{
				Op:   kind.NOT,
				Node: &ast.LiteralExpr{Value: "true", Kind: kind.TRUE_BOOL_LITERAL},
			},
		},
	}
	for _, test := range unaryExprs {
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
		{
			expr: "2 - 1",
			node: &ast.BinaryExpr{
				Left:  &ast.LiteralExpr{Value: 2, Kind: kind.INTEGER_LITERAL},
				Op:    kind.MINUS,
				Right: &ast.LiteralExpr{Value: 1, Kind: kind.INTEGER_LITERAL},
			},
		},
		{
			expr: "5 * 10",
			node: &ast.BinaryExpr{
				Left:  &ast.LiteralExpr{Value: 5, Kind: kind.INTEGER_LITERAL},
				Op:    kind.STAR,
				Right: &ast.LiteralExpr{Value: 10, Kind: kind.INTEGER_LITERAL},
			},
		},
		{
			expr: "10 / 1",
			node: &ast.BinaryExpr{
				Left:  &ast.LiteralExpr{Value: 10, Kind: kind.INTEGER_LITERAL},
				Op:    kind.SLASH,
				Right: &ast.LiteralExpr{Value: 1, Kind: kind.INTEGER_LITERAL},
			},
		},
		{
			expr: "6 / 3 - 1",
			node: &ast.BinaryExpr{
				Left: &ast.BinaryExpr{
					Left: &ast.LiteralExpr{
						Value: 6,
						Kind:  kind.INTEGER_LITERAL,
					},
					Op: kind.SLASH,
					Right: &ast.LiteralExpr{
						Value: 3,
						Kind:  kind.INTEGER_LITERAL,
					},
				},
				Op: kind.MINUS,
				Right: &ast.LiteralExpr{
					Value: 1,
					Kind:  kind.INTEGER_LITERAL,
				},
			},
		},
		{
			expr: "6 / (3 - 1)",
			node: &ast.BinaryExpr{
				Left: &ast.LiteralExpr{
					Value: 6,
					Kind:  kind.INTEGER_LITERAL,
				},
				Op: kind.SLASH,
				Right: &ast.BinaryExpr{
					Left: &ast.LiteralExpr{
						Value: 3,
						Kind:  kind.INTEGER_LITERAL,
					},
					Op: kind.MINUS,
					Right: &ast.LiteralExpr{
						Value: 1,
						Kind:  kind.INTEGER_LITERAL,
					},
				},
			},
		},
		{
			expr: "1 / (1 + 1)",
			node: &ast.BinaryExpr{
				Left: &ast.LiteralExpr{Value: 1, Kind: kind.INTEGER_LITERAL},
				Op:   kind.SLASH,
				Right: &ast.BinaryExpr{
					Left: &ast.LiteralExpr{
						Value: 1,
						Kind:  kind.INTEGER_LITERAL,
					},
					Op: kind.PLUS,
					Right: &ast.LiteralExpr{
						Value: 1,
						Kind:  kind.INTEGER_LITERAL,
					},
				},
			},
		},
		{
			expr: "1 > 1",
			node: &ast.BinaryExpr{
				Left: &ast.LiteralExpr{
					Value: 1,
					Kind:  kind.INTEGER_LITERAL,
				},
				Op: kind.GREATER,
				Right: &ast.LiteralExpr{
					Value: 1,
					Kind:  kind.INTEGER_LITERAL,
				},
			},
		},
		{
			expr: "1 >= 1",
			node: &ast.BinaryExpr{
				Left: &ast.LiteralExpr{
					Value: 1,
					Kind:  kind.INTEGER_LITERAL,
				},
				Op: kind.GREATER_EQ,
				Right: &ast.LiteralExpr{
					Value: 1,
					Kind:  kind.INTEGER_LITERAL,
				},
			},
		},
		{
			expr: "1 < 1",
			node: &ast.BinaryExpr{
				Left: &ast.LiteralExpr{
					Value: 1,
					Kind:  kind.INTEGER_LITERAL,
				},
				Op: kind.LESS,
				Right: &ast.LiteralExpr{
					Value: 1,
					Kind:  kind.INTEGER_LITERAL,
				},
			},
		},
		{
			expr: "1 <= 1",
			node: &ast.BinaryExpr{
				Left: &ast.LiteralExpr{
					Value: 1,
					Kind:  kind.INTEGER_LITERAL,
				},
				Op: kind.LESS_EQ,
				Right: &ast.LiteralExpr{
					Value: 1,
					Kind:  kind.INTEGER_LITERAL,
				},
			},
		},
		{
			// not 1 > 1 is invalid in Golang
			expr: "not (1 > 1)",
			node: &ast.UnaryExpr{
				Op: kind.NOT,
				Node: &ast.BinaryExpr{
					Left: &ast.LiteralExpr{
						Value: 1,
						Kind:  kind.INTEGER_LITERAL,
					},
					Op: kind.GREATER,
					Right: &ast.LiteralExpr{
						Value: 1,
						Kind:  kind.INTEGER_LITERAL,
					},
				},
			},
		},
		{
			expr: "1 > 1 and 1 > 1",
			node: &ast.BinaryExpr{
				Left: &ast.BinaryExpr{
					Left: &ast.LiteralExpr{
						Value: 1,
						Kind:  kind.INTEGER_LITERAL,
					},
					Op: kind.GREATER,
					Right: &ast.LiteralExpr{
						Value: 1,
						Kind:  kind.INTEGER_LITERAL,
					},
				},
				Op: kind.AND,
				Right: &ast.BinaryExpr{
					Left: &ast.LiteralExpr{
						Value: 1,
						Kind:  kind.INTEGER_LITERAL,
					},
					Op: kind.GREATER,
					Right: &ast.LiteralExpr{
						Value: 1,
						Kind:  kind.INTEGER_LITERAL,
					},
				},
			},
		},
		{
			expr: "1 > 1 or 1 > 1",
			node: &ast.BinaryExpr{
				Left: &ast.BinaryExpr{
					Left: &ast.LiteralExpr{
						Value: 1,
						Kind:  kind.INTEGER_LITERAL,
					},
					Op: kind.GREATER,
					Right: &ast.LiteralExpr{
						Value: 1,
						Kind:  kind.INTEGER_LITERAL,
					},
				},
				Op: kind.OR,
				Right: &ast.BinaryExpr{
					Left: &ast.LiteralExpr{
						Value: 1,
						Kind:  kind.INTEGER_LITERAL,
					},
					Op: kind.GREATER,
					Right: &ast.LiteralExpr{
						Value: 1,
						Kind:  kind.INTEGER_LITERAL,
					},
				},
			},
		},
		{
			expr: "celsius*9/5+32",
			node: &ast.BinaryExpr{
				Left: &ast.BinaryExpr{
					Left: &ast.BinaryExpr{
						Left: &ast.IdExpr{
							Name: token.New("celsius", kind.ID, token.NewPosition("test.tt", 1, 1)),
						},
						Op: kind.STAR,
						Right: &ast.LiteralExpr{
							Value: 9,
							Kind:  kind.INTEGER_LITERAL,
						},
					},
					Op: kind.SLASH,
					Right: &ast.LiteralExpr{
						Value: 5,
						Kind:  kind.INTEGER_LITERAL,
					},
				},
				Op: kind.PLUS,
				Right: &ast.LiteralExpr{
					Value: 32,
					Kind:  kind.INTEGER_LITERAL,
				},
			},
		},
		// Semantically, the input below is invalid
		// "bool > integer"
		{
			expr: "1 > 1 > 1",
			node: &ast.BinaryExpr{
				Left: &ast.BinaryExpr{
					Left: &ast.LiteralExpr{
						Value: 1,
						Kind:  kind.INTEGER_LITERAL,
					},
					Op: kind.GREATER,
					Right: &ast.LiteralExpr{
						Value: 1,
						Kind:  kind.INTEGER_LITERAL,
					},
				},
				Op: kind.GREATER,
				Right: &ast.LiteralExpr{
					Value: 1,
					Kind:  kind.INTEGER_LITERAL,
				},
			},
		},
		{
			expr: "n == 1",
			node: &ast.BinaryExpr{
				Left: &ast.IdExpr{
					Name: token.New("n", kind.ID, token.NewPosition("test.tt", 1, 1)),
				},
				Op: kind.EQUAL_EQUAL,
				Right: &ast.LiteralExpr{
					Value: 1,
					Kind:  kind.INTEGER_LITERAL,
				},
			},
		},
		{
			expr: "n == 1 or n == 2",
			node: &ast.BinaryExpr{
				Left: &ast.BinaryExpr{
					Left: &ast.IdExpr{
						Name: token.New("n", kind.ID, token.NewPosition("test.tt", 1, 1)),
					},
					Op: kind.EQUAL_EQUAL,
					Right: &ast.LiteralExpr{
						Kind:  kind.INTEGER_LITERAL,
						Value: 1,
					},
				},
				Op: kind.OR,
				Right: &ast.BinaryExpr{
					Left: &ast.IdExpr{
						Name: token.New("n", kind.ID, token.NewPosition("test.tt", 11, 1)),
					},
					Op: kind.EQUAL_EQUAL,
					Right: &ast.LiteralExpr{
						Kind:  kind.INTEGER_LITERAL,
						Value: 2,
					},
				},
			},
		},
		{
			expr: "1 + 1 > 2 and 1 == 1",
			node: &ast.BinaryExpr{
				Left: &ast.BinaryExpr{
					Left: &ast.BinaryExpr{
						Left: &ast.LiteralExpr{
							Value: 1,
							Kind:  kind.INTEGER_LITERAL,
						},
						Op: kind.PLUS,
						Right: &ast.LiteralExpr{
							Value: 1,
							Kind:  kind.INTEGER_LITERAL,
						},
					},
					Op: kind.GREATER,
					Right: &ast.LiteralExpr{
						Value: 2,
						Kind:  kind.INTEGER_LITERAL,
					},
				},
				Op: kind.AND,
				Right: &ast.BinaryExpr{
					Left: &ast.LiteralExpr{
						Value: 1,
						Kind:  kind.INTEGER_LITERAL,
					},
					Op: kind.EQUAL_EQUAL,
					Right: &ast.LiteralExpr{
						Value: 1,
						Kind:  kind.INTEGER_LITERAL,
					},
				},
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
				t.Errorf("TestBinaryExpr('%s'): expression node differs\nexpected: '%v' '%v'\ngot:      '%v' '%v'\n", test.expr, test.node, reflect.TypeOf(test.node), actualNode, reflect.TypeOf(actualNode))
			}
		})
	}
}

func TestFuncCallExpr(t *testing.T) {}

func TestFuncCallStmt(t *testing.T) {}

func TestIfStmt(t *testing.T) {}
