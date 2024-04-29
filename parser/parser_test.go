package parser

import (
	"bufio"
	"fmt"
	"reflect"
	"strings"
	"testing"

	"github.com/HicaroD/Telia/ast"
	"github.com/HicaroD/Telia/collector"
	"github.com/HicaroD/Telia/lexer"
	"github.com/HicaroD/Telia/lexer/token"
	"github.com/HicaroD/Telia/lexer/token/kind"
)

type functionDeclTest struct {
	input string
	node  *ast.FunctionDecl
}

func TestFunctionDecl(t *testing.T) {
	filename := "test.tt"
	tests := []functionDeclTest{
		{
			input: "fn do_nothing() {}",
			node: &ast.FunctionDecl{
				Scope: nil,
				Name:  token.New("do_nothing", kind.ID, token.NewPosition(filename, 4, 1)),
				Params: &ast.FieldList{
					Open:   token.New("", kind.OPEN_PAREN, token.NewPosition(filename, 14, 1)),
					Fields: nil,
					Close: token.New(
						"",
						kind.CLOSE_PAREN,
						token.NewPosition(filename, 15, 1),
					),
					IsVariadic: false,
				},
				RetType: &ast.BasicType{Kind: kind.VOID_TYPE},
				Block: &ast.BlockStmt{
					OpenCurly:  token.NewPosition(filename, 17, 1),
					Statements: nil,
					CloseCurly: token.NewPosition(filename, 18, 1),
				},
			},
		},
		{
			input: "fn do_nothing(a bool) {}",
			node: &ast.FunctionDecl{
				Scope: nil,
				Name:  token.New("do_nothing", kind.ID, token.NewPosition(filename, 4, 1)),
				Params: &ast.FieldList{
					Open: token.New("", kind.OPEN_PAREN, token.NewPosition(filename, 14, 1)),
					Fields: []*ast.Field{
						{
							Name: token.New("a", kind.ID, token.NewPosition(filename, 15, 1)),
							Type: &ast.BasicType{Kind: kind.BOOL_TYPE},
						},
					},
					Close: token.New(
						"",
						kind.CLOSE_PAREN,
						token.NewPosition(filename, 21, 1),
					),
					IsVariadic: false,
				},
				RetType: &ast.BasicType{Kind: kind.VOID_TYPE},
				Block: &ast.BlockStmt{
					OpenCurly:  token.NewPosition(filename, 23, 1),
					Statements: nil,
					CloseCurly: token.NewPosition(filename, 24, 1),
				},
			},
		},
		{
			input: "fn do_nothing(a bool, b i32) {}",
			node: &ast.FunctionDecl{
				Scope: nil,
				Name:  token.New("do_nothing", kind.ID, token.NewPosition(filename, 4, 1)),
				Params: &ast.FieldList{
					Open: token.New("", kind.OPEN_PAREN, token.NewPosition(filename, 14, 1)),
					Fields: []*ast.Field{
						{
							Name: token.New("a", kind.ID, token.NewPosition(filename, 15, 1)),
							Type: &ast.BasicType{Kind: kind.BOOL_TYPE},
						},
						{
							Name: token.New("b", kind.ID, token.NewPosition(filename, 23, 1)),
							Type: &ast.BasicType{Kind: kind.I32_TYPE},
						},
					},
					Close: token.New(
						"",
						kind.CLOSE_PAREN,
						token.NewPosition(filename, 28, 1),
					),
					IsVariadic: false,
				},
				RetType: &ast.BasicType{Kind: kind.VOID_TYPE},
				Block: &ast.BlockStmt{
					OpenCurly:  token.NewPosition(filename, 30, 1),
					Statements: nil,
					CloseCurly: token.NewPosition(filename, 31, 1),
				},
			},
		},
		{
			input: "fn do_nothing(a bool, b i32) i8 {}",
			node: &ast.FunctionDecl{
				Scope: nil,
				Name:  token.New("do_nothing", kind.ID, token.NewPosition(filename, 4, 1)),
				Params: &ast.FieldList{
					Open: token.New("", kind.OPEN_PAREN, token.NewPosition(filename, 14, 1)),
					Fields: []*ast.Field{
						{
							Name: token.New("a", kind.ID, token.NewPosition(filename, 15, 1)),
							Type: &ast.BasicType{Kind: kind.BOOL_TYPE},
						},
						{
							Name: token.New("b", kind.ID, token.NewPosition(filename, 23, 1)),
							Type: &ast.BasicType{Kind: kind.I32_TYPE},
						},
					},
					Close: token.New(
						"",
						kind.CLOSE_PAREN,
						token.NewPosition(filename, 28, 1),
					),
					IsVariadic: false,
				},
				RetType: &ast.BasicType{Kind: kind.I8_TYPE},
				Block: &ast.BlockStmt{
					OpenCurly:  token.NewPosition(filename, 33, 1),
					Statements: nil,
					CloseCurly: token.NewPosition(filename, 34, 1),
				},
			},
		},
		{
			input: "fn do_nothing(a bool, b i32) u8 {}",
			node: &ast.FunctionDecl{
				Scope: nil,
				Name:  token.New("do_nothing", kind.ID, token.NewPosition(filename, 4, 1)),
				Params: &ast.FieldList{
					Open: token.New("", kind.OPEN_PAREN, token.NewPosition(filename, 14, 1)),
					Fields: []*ast.Field{
						{
							Name: token.New("a", kind.ID, token.NewPosition(filename, 15, 1)),
							Type: &ast.BasicType{Kind: kind.BOOL_TYPE},
						},
						{
							Name: token.New("b", kind.ID, token.NewPosition(filename, 23, 1)),
							Type: &ast.BasicType{Kind: kind.I32_TYPE},
						},
					},
					Close: token.New(
						"",
						kind.CLOSE_PAREN,
						token.NewPosition(filename, 28, 1),
					),
					IsVariadic: false,
				},
				RetType: &ast.BasicType{Kind: kind.U8_TYPE},
				Block: &ast.BlockStmt{
					OpenCurly:  token.NewPosition(filename, 33, 1),
					Statements: nil,
					CloseCurly: token.NewPosition(filename, 34, 1),
				},
			},
		},
		{
			input: "fn do_nothing(a bool, b i32) i16 {}",
			node: &ast.FunctionDecl{
				Scope: nil,
				Name:  token.New("do_nothing", kind.ID, token.NewPosition(filename, 4, 1)),
				Params: &ast.FieldList{
					Open: token.New("", kind.OPEN_PAREN, token.NewPosition(filename, 14, 1)),
					Fields: []*ast.Field{
						{
							Name: token.New("a", kind.ID, token.NewPosition(filename, 15, 1)),
							Type: &ast.BasicType{Kind: kind.BOOL_TYPE},
						},
						{
							Name: token.New("b", kind.ID, token.NewPosition(filename, 23, 1)),
							Type: &ast.BasicType{Kind: kind.I32_TYPE},
						},
					},
					Close: token.New(
						"",
						kind.CLOSE_PAREN,
						token.NewPosition(filename, 28, 1),
					),
					IsVariadic: false,
				},
				RetType: &ast.BasicType{Kind: kind.I16_TYPE},
				Block: &ast.BlockStmt{
					OpenCurly:  token.NewPosition(filename, 34, 1),
					Statements: nil,
					CloseCurly: token.NewPosition(filename, 35, 1),
				},
			},
		},
		{
			input: "fn do_nothing(a bool, b i32) u16 {}",
			node: &ast.FunctionDecl{
				Scope: nil,
				Name:  token.New("do_nothing", kind.ID, token.NewPosition(filename, 4, 1)),
				Params: &ast.FieldList{
					Open: token.New("", kind.OPEN_PAREN, token.NewPosition(filename, 14, 1)),
					Fields: []*ast.Field{
						{
							Name: token.New("a", kind.ID, token.NewPosition(filename, 15, 1)),
							Type: &ast.BasicType{Kind: kind.BOOL_TYPE},
						},
						{
							Name: token.New("b", kind.ID, token.NewPosition(filename, 23, 1)),
							Type: &ast.BasicType{Kind: kind.I32_TYPE},
						},
					},
					Close: token.New(
						"",
						kind.CLOSE_PAREN,
						token.NewPosition(filename, 28, 1),
					),
					IsVariadic: false,
				},
				RetType: &ast.BasicType{Kind: kind.U16_TYPE},
				Block: &ast.BlockStmt{
					OpenCurly:  token.NewPosition(filename, 34, 1),
					Statements: nil,
					CloseCurly: token.NewPosition(filename, 35, 1),
				},
			},
		},
		{
			input: "fn do_nothing(a bool, b i32) i32 {}",
			node: &ast.FunctionDecl{
				Scope: nil,
				Name:  token.New("do_nothing", kind.ID, token.NewPosition(filename, 4, 1)),
				Params: &ast.FieldList{
					Open: token.New("", kind.OPEN_PAREN, token.NewPosition(filename, 14, 1)),
					Fields: []*ast.Field{
						{
							Name: token.New("a", kind.ID, token.NewPosition(filename, 15, 1)),
							Type: &ast.BasicType{Kind: kind.BOOL_TYPE},
						},
						{
							Name: token.New("b", kind.ID, token.NewPosition(filename, 23, 1)),
							Type: &ast.BasicType{Kind: kind.I32_TYPE},
						},
					},
					Close: token.New(
						"",
						kind.CLOSE_PAREN,
						token.NewPosition(filename, 28, 1),
					),
					IsVariadic: false,
				},
				RetType: &ast.BasicType{Kind: kind.I32_TYPE},
				Block: &ast.BlockStmt{
					OpenCurly:  token.NewPosition(filename, 34, 1),
					Statements: nil,
					CloseCurly: token.NewPosition(filename, 35, 1),
				},
			},
		},
		{
			input: "fn do_nothing(a bool, b i32) u32 {}",
			node: &ast.FunctionDecl{
				Scope: nil,
				Name:  token.New("do_nothing", kind.ID, token.NewPosition(filename, 4, 1)),
				Params: &ast.FieldList{
					Open: token.New("", kind.OPEN_PAREN, token.NewPosition(filename, 14, 1)),
					Fields: []*ast.Field{
						{
							Name: token.New("a", kind.ID, token.NewPosition(filename, 15, 1)),
							Type: &ast.BasicType{Kind: kind.BOOL_TYPE},
						},
						{
							Name: token.New("b", kind.ID, token.NewPosition(filename, 23, 1)),
							Type: &ast.BasicType{Kind: kind.I32_TYPE},
						},
					},
					Close: token.New(
						"",
						kind.CLOSE_PAREN,
						token.NewPosition(filename, 28, 1),
					),
					IsVariadic: false,
				},
				RetType: &ast.BasicType{Kind: kind.U32_TYPE},
				Block: &ast.BlockStmt{
					OpenCurly:  token.NewPosition(filename, 34, 1),
					Statements: nil,
					CloseCurly: token.NewPosition(filename, 35, 1),
				},
			},
		},
		{
			input: "fn do_nothing(a bool, b i32) i64 {}",
			node: &ast.FunctionDecl{
				Scope: nil,
				Name:  token.New("do_nothing", kind.ID, token.NewPosition(filename, 4, 1)),
				Params: &ast.FieldList{
					Open: token.New("", kind.OPEN_PAREN, token.NewPosition(filename, 14, 1)),
					Fields: []*ast.Field{
						{
							Name: token.New("a", kind.ID, token.NewPosition(filename, 15, 1)),
							Type: &ast.BasicType{Kind: kind.BOOL_TYPE},
						},
						{
							Name: token.New("b", kind.ID, token.NewPosition(filename, 23, 1)),
							Type: &ast.BasicType{Kind: kind.I32_TYPE},
						},
					},
					Close: token.New(
						"",
						kind.CLOSE_PAREN,
						token.NewPosition(filename, 28, 1),
					),
					IsVariadic: false,
				},
				RetType: &ast.BasicType{Kind: kind.I64_TYPE},
				Block: &ast.BlockStmt{
					OpenCurly:  token.NewPosition(filename, 34, 1),
					Statements: nil,
					CloseCurly: token.NewPosition(filename, 35, 1),
				},
			},
		},
		{
			input: "fn do_nothing(a bool, b i32) u64 {}",
			node: &ast.FunctionDecl{
				Scope: nil,
				Name:  token.New("do_nothing", kind.ID, token.NewPosition(filename, 4, 1)),
				Params: &ast.FieldList{
					Open: token.New("", kind.OPEN_PAREN, token.NewPosition(filename, 14, 1)),
					Fields: []*ast.Field{
						{
							Name: token.New("a", kind.ID, token.NewPosition(filename, 15, 1)),
							Type: &ast.BasicType{Kind: kind.BOOL_TYPE},
						},
						{
							Name: token.New("b", kind.ID, token.NewPosition(filename, 23, 1)),
							Type: &ast.BasicType{Kind: kind.I32_TYPE},
						},
					},
					Close: token.New(
						"",
						kind.CLOSE_PAREN,
						token.NewPosition(filename, 28, 1),
					),
					IsVariadic: false,
				},
				RetType: &ast.BasicType{Kind: kind.U64_TYPE},
				Block: &ast.BlockStmt{
					OpenCurly:  token.NewPosition(filename, 34, 1),
					Statements: nil,
					CloseCurly: token.NewPosition(filename, 35, 1),
				},
			},
		},
		{
			input: "fn do_nothing(a bool, b i32) bool {}",
			node: &ast.FunctionDecl{
				Scope: nil,
				Name:  token.New("do_nothing", kind.ID, token.NewPosition(filename, 4, 1)),
				Params: &ast.FieldList{
					Open: token.New("", kind.OPEN_PAREN, token.NewPosition(filename, 14, 1)),
					Fields: []*ast.Field{
						{
							Name: token.New("a", kind.ID, token.NewPosition(filename, 15, 1)),
							Type: &ast.BasicType{Kind: kind.BOOL_TYPE},
						},
						{
							Name: token.New("b", kind.ID, token.NewPosition(filename, 23, 1)),
							Type: &ast.BasicType{Kind: kind.I32_TYPE},
						},
					},
					Close: token.New(
						"",
						kind.CLOSE_PAREN,
						token.NewPosition(filename, 28, 1),
					),
					IsVariadic: false,
				},
				RetType: &ast.BasicType{Kind: kind.BOOL_TYPE},
				Block: &ast.BlockStmt{
					OpenCurly:  token.NewPosition(filename, 35, 1),
					Statements: nil,
					CloseCurly: token.NewPosition(filename, 36, 1),
				},
			},
		},
		{
			input: "fn do_nothing(a bool, b i32) *i8 {}",
			node: &ast.FunctionDecl{
				Scope: nil,
				Name:  token.New("do_nothing", kind.ID, token.NewPosition(filename, 4, 1)),
				Params: &ast.FieldList{
					Open: token.New("", kind.OPEN_PAREN, token.NewPosition(filename, 14, 1)),
					Fields: []*ast.Field{
						{
							Name: token.New("a", kind.ID, token.NewPosition(filename, 15, 1)),
							Type: &ast.BasicType{Kind: kind.BOOL_TYPE},
						},
						{
							Name: token.New("b", kind.ID, token.NewPosition(filename, 23, 1)),
							Type: &ast.BasicType{Kind: kind.I32_TYPE},
						},
					},
					Close: token.New(
						"",
						kind.CLOSE_PAREN,
						token.NewPosition(filename, 28, 1),
					),
					IsVariadic: false,
				},
				RetType: &ast.PointerType{Type: &ast.BasicType{Kind: kind.I8_TYPE}},
				Block: &ast.BlockStmt{
					OpenCurly:  token.NewPosition(filename, 34, 1),
					Statements: nil,
					CloseCurly: token.NewPosition(filename, 35, 1),
				},
			},
		},
		// TODO(tests): test variadic arguments on functions
	}
	for _, test := range tests {
		t.Run(fmt.Sprintf("TestFunctionDecl('%s')", test.input), func(t *testing.T) {
			fnDecl, err := parseFnDeclFrom(filename, test.input)
			if err != nil {
				t.Fatal(err)
			}
			if !reflect.DeepEqual(fnDecl, test.node) {
				t.Fatalf("\n-------\nexpected:\n%s\n-------\ngot:\n%s\n-------", test.node, fnDecl)
			}
		})
	}
}

// TODO(tests)
// type externDeclTest struct {
// 	input string
// 	node  *ast.ExternDecl
// }
//
// func TestExternDecl(t *testing.T) {
// 	filename := "test.tt"
// 	tests := []externDeclTest{
// 		{
// 			input: "extern libc {}",
// 		},
// 	}
//
// 	for _, test := range tests {
// 		t.Run(fmt.Sprintf("TestExternDecl('%s')", test.input), func(t *testing.T) {
// 		})
// 	}
// }

type exprTest struct {
	input string
	node  ast.Expr
}

func TestLiteralExpr(t *testing.T) {
	filename := "test.tt"
	tests := []exprTest{
		{
			input: "1",
			node:  &ast.LiteralExpr{Value: "1", Type: &ast.BasicType{Kind: kind.INTEGER_LITERAL}},
		},
		{
			input: "true",
			node: &ast.LiteralExpr{
				Value: "true",
				Type:  &ast.BasicType{Kind: kind.TRUE_BOOL_LITERAL},
			},
		},
		{
			input: "false",
			node: &ast.LiteralExpr{
				Value: "false",
				Type:  &ast.BasicType{Kind: kind.FALSE_BOOL_LITERAL},
			},
		},
		{
			input: "\"Hello, world\"",
			node: &ast.LiteralExpr{
				Value: "Hello, world",
				Type:  &ast.BasicType{Kind: kind.STRING_LITERAL},
			},
		},
	}

	for _, test := range tests {
		t.Run(fmt.Sprintf("TestLiteralExpr('%s')", test.input), func(t *testing.T) {
			actualNode, err := ParseExprFrom(test.input, filename)
			if err != nil {
				t.Errorf("TestLiteralExpr('%s'): unexpected error '%v'", test.input, err)
			}
			if !reflect.DeepEqual(test.node, actualNode) {
				t.Errorf(
					"TestLiteralExpr('%s'): expression node differs\nexpected: '%v', but got '%v'\n",
					test.input,
					test.node,
					actualNode,
				)
			}
		})
	}
}

func TestUnaryExpr(t *testing.T) {
	filename := "test.tt"
	tests := []exprTest{
		{
			input: "-1",
			node: &ast.UnaryExpr{
				Op: kind.MINUS,
				Value: &ast.LiteralExpr{
					Value: "1",
					Type:  &ast.BasicType{Kind: kind.INTEGER_LITERAL},
				},
			},
		},
		{
			input: "not true",
			node: &ast.UnaryExpr{
				Op: kind.NOT,
				Value: &ast.LiteralExpr{
					Value: "true",
					Type:  &ast.BasicType{Kind: kind.TRUE_BOOL_LITERAL},
				},
			},
		},
	}
	for _, test := range tests {
		t.Run(fmt.Sprintf("TestUnaryExpr('%s')", test.input), func(t *testing.T) {
			actualNode, err := ParseExprFrom(test.input, filename)
			if err != nil {
				t.Errorf("TestUnaryExpr('%s'): unexpected error '%v'", test.input, err)
			}
			if !reflect.DeepEqual(test.node, actualNode) {
				t.Errorf(
					"TestUnaryExpr('%s'): expression node differs\nexpected: '%v', but got '%v'\n",
					test.input,
					test.node,
					actualNode,
				)
			}
		})
	}
}

func TestBinaryExpr(t *testing.T) {
	filename := "test.tt"
	tests := []exprTest{
		{
			input: "1 + 1",
			node: &ast.BinaryExpr{
				Left: &ast.LiteralExpr{
					Value: "1",
					Type:  &ast.BasicType{Kind: kind.INTEGER_LITERAL},
				},
				Op: kind.PLUS,
				Right: &ast.LiteralExpr{
					Value: "1",
					Type:  &ast.BasicType{Kind: kind.INTEGER_LITERAL},
				},
			},
		},
		{
			input: "2 - 1",
			node: &ast.BinaryExpr{
				Left: &ast.LiteralExpr{
					Value: "2",
					Type:  &ast.BasicType{Kind: kind.INTEGER_LITERAL},
				},
				Op: kind.MINUS,
				Right: &ast.LiteralExpr{
					Value: "1",
					Type:  &ast.BasicType{Kind: kind.INTEGER_LITERAL},
				},
			},
		},
		{
			input: "5 * 10",
			node: &ast.BinaryExpr{
				Left: &ast.LiteralExpr{
					Value: "5",
					Type:  &ast.BasicType{Kind: kind.INTEGER_LITERAL},
				},
				Op: kind.STAR,
				Right: &ast.LiteralExpr{
					Value: "10",
					Type:  &ast.BasicType{Kind: kind.INTEGER_LITERAL},
				},
			},
		},
		{
			input: "3 + 4 * 5",
			node: &ast.BinaryExpr{
				Left: &ast.LiteralExpr{
					Value: "3",
					Type:  &ast.BasicType{Kind: kind.INTEGER_LITERAL},
				},
				Op: kind.PLUS,
				Right: &ast.BinaryExpr{
					Left: &ast.LiteralExpr{
						Value: "4",
						Type:  &ast.BasicType{Kind: kind.INTEGER_LITERAL},
					},
					Op: kind.STAR,
					Right: &ast.LiteralExpr{
						Value: "5",
						Type:  &ast.BasicType{Kind: kind.INTEGER_LITERAL},
					},
				},
			},
		},
		{
			input: "3 + (4 * 5)",
			node: &ast.BinaryExpr{
				Left: &ast.LiteralExpr{
					Value: "3",
					Type:  &ast.BasicType{Kind: kind.INTEGER_LITERAL},
				},
				Op: kind.PLUS,
				Right: &ast.BinaryExpr{
					Left: &ast.LiteralExpr{
						Value: "4",
						Type:  &ast.BasicType{Kind: kind.INTEGER_LITERAL},
					},
					Op: kind.STAR,
					Right: &ast.LiteralExpr{
						Value: "5",
						Type:  &ast.BasicType{Kind: kind.INTEGER_LITERAL},
					},
				},
			},
		},
		{
			input: "10 / 1",
			node: &ast.BinaryExpr{
				Left: &ast.LiteralExpr{
					Value: "10",
					Type:  &ast.BasicType{Kind: kind.INTEGER_LITERAL},
				},
				Op: kind.SLASH,
				Right: &ast.LiteralExpr{
					Value: "1",
					Type:  &ast.BasicType{Kind: kind.INTEGER_LITERAL},
				},
			},
		},
		{
			input: "6 / 3 - 1",
			node: &ast.BinaryExpr{
				Left: &ast.BinaryExpr{
					Left: &ast.LiteralExpr{
						Value: "6",
						Type:  &ast.BasicType{Kind: kind.INTEGER_LITERAL},
					},
					Op: kind.SLASH,
					Right: &ast.LiteralExpr{
						Value: "3",
						Type:  &ast.BasicType{Kind: kind.INTEGER_LITERAL},
					},
				},
				Op: kind.MINUS,
				Right: &ast.LiteralExpr{
					Value: "1",
					Type:  &ast.BasicType{Kind: kind.INTEGER_LITERAL},
				},
			},
		},
		{
			input: "6 / (3 - 1)",
			node: &ast.BinaryExpr{
				Left: &ast.LiteralExpr{
					Value: "6",
					Type:  &ast.BasicType{Kind: kind.INTEGER_LITERAL},
				},
				Op: kind.SLASH,
				Right: &ast.BinaryExpr{
					Left: &ast.LiteralExpr{
						Value: "3",
						Type:  &ast.BasicType{Kind: kind.INTEGER_LITERAL},
					},
					Op: kind.MINUS,
					Right: &ast.LiteralExpr{
						Value: "1",
						Type:  &ast.BasicType{Kind: kind.INTEGER_LITERAL},
					},
				},
			},
		},
		{
			input: "1 / (1 + 1)",
			node: &ast.BinaryExpr{
				Left: &ast.LiteralExpr{
					Value: "1",
					Type:  &ast.BasicType{Kind: kind.INTEGER_LITERAL},
				},
				Op: kind.SLASH,
				Right: &ast.BinaryExpr{
					Left: &ast.LiteralExpr{
						Value: "1",
						Type:  &ast.BasicType{Kind: kind.INTEGER_LITERAL},
					},
					Op: kind.PLUS,
					Right: &ast.LiteralExpr{
						Value: "1",
						Type:  &ast.BasicType{Kind: kind.INTEGER_LITERAL},
					},
				},
			},
		},
		{
			input: "1 > 1",
			node: &ast.BinaryExpr{
				Left: &ast.LiteralExpr{
					Value: "1",
					Type:  &ast.BasicType{Kind: kind.INTEGER_LITERAL},
				},
				Op: kind.GREATER,
				Right: &ast.LiteralExpr{
					Value: "1",
					Type:  &ast.BasicType{Kind: kind.INTEGER_LITERAL},
				},
			},
		},
		{
			input: "1 >= 1",
			node: &ast.BinaryExpr{
				Left: &ast.LiteralExpr{
					Value: "1",
					Type:  &ast.BasicType{Kind: kind.INTEGER_LITERAL},
				},
				Op: kind.GREATER_EQ,
				Right: &ast.LiteralExpr{
					Value: "1",
					Type:  &ast.BasicType{Kind: kind.INTEGER_LITERAL},
				},
			},
		},
		{
			input: "1 < 1",
			node: &ast.BinaryExpr{
				Left: &ast.LiteralExpr{
					Value: "1",
					Type:  &ast.BasicType{Kind: kind.INTEGER_LITERAL},
				},
				Op: kind.LESS,
				Right: &ast.LiteralExpr{
					Value: "1",
					Type:  &ast.BasicType{Kind: kind.INTEGER_LITERAL},
				},
			},
		},
		{
			input: "1 <= 1",
			node: &ast.BinaryExpr{
				Left: &ast.LiteralExpr{
					Value: "1",
					Type:  &ast.BasicType{Kind: kind.INTEGER_LITERAL},
				},
				Op: kind.LESS_EQ,
				Right: &ast.LiteralExpr{
					Value: "1",
					Type:  &ast.BasicType{Kind: kind.INTEGER_LITERAL},
				},
			},
		},
		{
			// not 1 > 1 is invalid in Golang
			input: "not (1 > 1)",
			node: &ast.UnaryExpr{
				Op: kind.NOT,
				Value: &ast.BinaryExpr{
					Left: &ast.LiteralExpr{
						Value: "1",
						Type:  &ast.BasicType{Kind: kind.INTEGER_LITERAL},
					},
					Op: kind.GREATER,
					Right: &ast.LiteralExpr{
						Value: "1",
						Type:  &ast.BasicType{Kind: kind.INTEGER_LITERAL},
					},
				},
			},
		},
		{
			input: "1 > 1 and 1 > 1",
			node: &ast.BinaryExpr{
				Left: &ast.BinaryExpr{
					Left: &ast.LiteralExpr{
						Value: "1",
						Type:  &ast.BasicType{Kind: kind.INTEGER_LITERAL},
					},
					Op: kind.GREATER,
					Right: &ast.LiteralExpr{
						Value: "1",
						Type:  &ast.BasicType{Kind: kind.INTEGER_LITERAL},
					},
				},
				Op: kind.AND,
				Right: &ast.BinaryExpr{
					Left: &ast.LiteralExpr{
						Value: "1",
						Type:  &ast.BasicType{Kind: kind.INTEGER_LITERAL},
					},
					Op: kind.GREATER,
					Right: &ast.LiteralExpr{
						Value: "1",
						Type:  &ast.BasicType{Kind: kind.INTEGER_LITERAL},
					},
				},
			},
		},
		{
			input: "1 > 1 or 1 > 1",
			node: &ast.BinaryExpr{
				Left: &ast.BinaryExpr{
					Left: &ast.LiteralExpr{
						Value: "1",
						Type:  &ast.BasicType{Kind: kind.INTEGER_LITERAL},
					},
					Op: kind.GREATER,
					Right: &ast.LiteralExpr{
						Value: "1",
						Type:  &ast.BasicType{Kind: kind.INTEGER_LITERAL},
					},
				},
				Op: kind.OR,
				Right: &ast.BinaryExpr{
					Left: &ast.LiteralExpr{
						Value: "1",
						Type:  &ast.BasicType{Kind: kind.INTEGER_LITERAL},
					},
					Op: kind.GREATER,
					Right: &ast.LiteralExpr{
						Value: "1",
						Type:  &ast.BasicType{Kind: kind.INTEGER_LITERAL},
					},
				},
			},
		},
		{
			input: "celsius*9/5+32",
			node: &ast.BinaryExpr{
				Left: &ast.BinaryExpr{
					Left: &ast.BinaryExpr{
						Left: &ast.IdExpr{
							Name: token.New("celsius", kind.ID, token.NewPosition("test.tt", 1, 1)),
						},
						Op: kind.STAR,
						Right: &ast.LiteralExpr{
							Value: "9",
							Type:  &ast.BasicType{Kind: kind.INTEGER_LITERAL},
						},
					},
					Op: kind.SLASH,
					Right: &ast.LiteralExpr{
						Value: "5",
						Type:  &ast.BasicType{Kind: kind.INTEGER_LITERAL},
					},
				},
				Op: kind.PLUS,
				Right: &ast.LiteralExpr{
					Value: "32",
					Type:  &ast.BasicType{Kind: kind.INTEGER_LITERAL},
				},
			},
		},
		{
			input: "get_celsius()*9/5+32",
			node: &ast.BinaryExpr{
				Left: &ast.BinaryExpr{
					Left: &ast.BinaryExpr{
						Left: &ast.FunctionCall{
							Name: token.New(
								"get_celsius",
								kind.ID,
								token.NewPosition(filename, 1, 1),
							),
							Args: nil,
						},
						Op: kind.STAR,
						Right: &ast.LiteralExpr{
							Value: "9",
							Type:  &ast.BasicType{Kind: kind.INTEGER_LITERAL},
						},
					},
					Op: kind.SLASH,
					Right: &ast.LiteralExpr{
						Value: "5",
						Type:  &ast.BasicType{Kind: kind.INTEGER_LITERAL},
					},
				},
				Op: kind.PLUS,
				Right: &ast.LiteralExpr{
					Value: "32",
					Type:  &ast.BasicType{Kind: kind.INTEGER_LITERAL},
				},
			},
		},
		{
			input: "1 > 1 > 1",
			node: &ast.BinaryExpr{
				Left: &ast.BinaryExpr{
					Left: &ast.LiteralExpr{
						Value: "1",
						Type:  &ast.BasicType{Kind: kind.INTEGER_LITERAL},
					},
					Op: kind.GREATER,
					Right: &ast.LiteralExpr{
						Value: "1",
						Type:  &ast.BasicType{Kind: kind.INTEGER_LITERAL},
					},
				},
				Op: kind.GREATER,
				Right: &ast.LiteralExpr{
					Value: "1",
					Type:  &ast.BasicType{Kind: kind.INTEGER_LITERAL},
				},
			},
		},
		{
			input: "n == 1",
			node: &ast.BinaryExpr{
				Left: &ast.IdExpr{
					Name: token.New("n", kind.ID, token.NewPosition("test.tt", 1, 1)),
				},
				Op: kind.EQUAL_EQUAL,
				Right: &ast.LiteralExpr{
					Value: "1",
					Type:  &ast.BasicType{Kind: kind.INTEGER_LITERAL},
				},
			},
		},
		{
			input: "n == 1 or n == 2",
			node: &ast.BinaryExpr{
				Left: &ast.BinaryExpr{
					Left: &ast.IdExpr{
						Name: token.New("n", kind.ID, token.NewPosition("test.tt", 1, 1)),
					},
					Op: kind.EQUAL_EQUAL,
					Right: &ast.LiteralExpr{
						Type:  &ast.BasicType{Kind: kind.INTEGER_LITERAL},
						Value: "1",
					},
				},
				Op: kind.OR,
				Right: &ast.BinaryExpr{
					Left: &ast.IdExpr{
						Name: token.New("n", kind.ID, token.NewPosition("test.tt", 11, 1)),
					},
					Op: kind.EQUAL_EQUAL,
					Right: &ast.LiteralExpr{
						Type:  &ast.BasicType{Kind: kind.INTEGER_LITERAL},
						Value: "2",
					},
				},
			},
		},
		{
			input: "1 + 1 > 2 and 1 == 1",
			node: &ast.BinaryExpr{
				Left: &ast.BinaryExpr{
					Left: &ast.BinaryExpr{
						Left: &ast.LiteralExpr{
							Value: "1",
							Type:  &ast.BasicType{Kind: kind.INTEGER_LITERAL},
						},
						Op: kind.PLUS,
						Right: &ast.LiteralExpr{
							Value: "1",
							Type:  &ast.BasicType{Kind: kind.INTEGER_LITERAL},
						},
					},
					Op: kind.GREATER,
					Right: &ast.LiteralExpr{
						Value: "2",
						Type:  &ast.BasicType{Kind: kind.INTEGER_LITERAL},
					},
				},
				Op: kind.AND,
				Right: &ast.BinaryExpr{
					Left: &ast.LiteralExpr{
						Value: "1",
						Type:  &ast.BasicType{Kind: kind.INTEGER_LITERAL},
					},
					Op: kind.EQUAL_EQUAL,
					Right: &ast.LiteralExpr{
						Value: "1",
						Type:  &ast.BasicType{Kind: kind.INTEGER_LITERAL},
					},
				},
			},
		},
		{
			input: "true and true and true and true",
			node: &ast.BinaryExpr{
				Left: &ast.BinaryExpr{
					Left: &ast.BinaryExpr{
						Left: &ast.LiteralExpr{
							Value: "true",
							Type:  &ast.BasicType{Kind: kind.TRUE_BOOL_LITERAL},
						},
						Op: kind.AND,
						Right: &ast.LiteralExpr{
							Value: "true",
							Type:  &ast.BasicType{Kind: kind.TRUE_BOOL_LITERAL},
						},
					},
					Op: kind.AND,
					Right: &ast.LiteralExpr{
						Value: "true",
						Type:  &ast.BasicType{Kind: kind.TRUE_BOOL_LITERAL},
					},
				},
				Op: kind.AND,
				Right: &ast.LiteralExpr{
					Value: "true",
					Type:  &ast.BasicType{Kind: kind.TRUE_BOOL_LITERAL},
				},
			},
		},
		{
			input: "(((true and true) and true) and true)",
			node: &ast.BinaryExpr{
				Left: &ast.BinaryExpr{
					Left: &ast.BinaryExpr{
						Left: &ast.LiteralExpr{
							Value: "true",
							Type:  &ast.BasicType{Kind: kind.TRUE_BOOL_LITERAL},
						},
						Op: kind.AND,
						Right: &ast.LiteralExpr{
							Value: "true",
							Type:  &ast.BasicType{Kind: kind.TRUE_BOOL_LITERAL},
						},
					},
					Op: kind.AND,
					Right: &ast.LiteralExpr{
						Value: "true",
						Type:  &ast.BasicType{Kind: kind.TRUE_BOOL_LITERAL},
					},
				},
				Op: kind.AND,
				Right: &ast.LiteralExpr{
					Value: "true",
					Type:  &ast.BasicType{Kind: kind.TRUE_BOOL_LITERAL},
				},
			},
		},
		{
			input: "1 + multiply_by_2(10)",
			node: &ast.BinaryExpr{
				Left: &ast.LiteralExpr{
					Value: "1",
					Type:  &ast.BasicType{Kind: kind.INTEGER_LITERAL},
				},
				Op: kind.PLUS,
				Right: &ast.FunctionCall{
					Name: token.New("multiply_by_2", kind.ID, token.NewPosition(filename, 5, 1)),
					Args: []ast.Expr{
						&ast.LiteralExpr{
							Value: "10",
							Type:  &ast.BasicType{Kind: kind.INTEGER_LITERAL},
						},
					},
				},
			},
		},
	}

	for _, test := range tests {
		t.Run(fmt.Sprintf("TestBinaryExpr('%s')", test.input), func(t *testing.T) {
			actualNode, err := ParseExprFrom(test.input, filename)
			if err != nil {
				t.Errorf("unexpected error '%v'", err)
			}
			if !reflect.DeepEqual(test.node, actualNode) {
				t.Errorf(
					"expression node differs\nexpected: '%v' '%v'\ngot:      '%v' '%v'\n",
					test.node,
					reflect.TypeOf(test.node),
					actualNode,
					reflect.TypeOf(actualNode),
				)
			}
		})
	}
}

func TestFieldAccessExpr(t *testing.T) {
	filename := "test.tt"
	tests := []exprTest{
		{
			input: "first.second.third",
			node: &ast.FieldAccess{
				Left: &ast.IdExpr{
					Name: token.New("first", kind.ID, token.NewPosition(filename, 1, 1)),
				},
				Right: &ast.FieldAccess{
					Left: &ast.IdExpr{
						Name: token.New("second", kind.ID, token.NewPosition(filename, 7, 1)),
					},
					Right: &ast.IdExpr{
						Name: token.New("third", kind.ID, token.NewPosition(filename, 14, 1)),
					},
				},
			},
		},
		{
			input: "first.second",
			node: &ast.FieldAccess{
				Left: &ast.IdExpr{
					Name: token.New("first", kind.ID, token.NewPosition(filename, 1, 1)),
				},
				Right: &ast.IdExpr{
					Name: token.New("second", kind.ID, token.NewPosition(filename, 7, 1)),
				},
			},
		},
	}

	for _, test := range tests {
		t.Run(fmt.Sprintf("TestBinaryExpr('%s')", test.input), func(t *testing.T) {
			actualNode, err := ParseExprFrom(test.input, filename)
			if err != nil {
				t.Errorf("unexpected error '%v'", err)
			}
			if !reflect.DeepEqual(test.node, actualNode) {
				t.Errorf(
					"expression node differs\nexpected: '%v' '%v'\ngot:      '%v' '%v'\n",
					test.node,
					reflect.TypeOf(test.node),
					actualNode,
					reflect.TypeOf(actualNode),
				)
			}
		})
	}
}

type varDeclTest struct {
	input   string
	varDecl *ast.MultiVarStmt
}

func TestVar(t *testing.T) {
	filename := "test.tt"
	tests := []varDeclTest{
		{
			input: "age := 10;",
			varDecl: &ast.MultiVarStmt{
				IsDecl: true,
				Variables: []*ast.VarDeclStmt{
					{
						Name: token.New(
							"age",
							kind.ID,
							token.NewPosition(filename, 1, 1),
						),
						Type:           nil,
						NeedsInference: true,
						Value: &ast.LiteralExpr{
							Type:  &ast.BasicType{Kind: kind.INTEGER_LITERAL},
							Value: "10",
						},
					},
				},
			},
		},
		{
			input: "score u8 := 10;",
			varDecl: &ast.MultiVarStmt{
				IsDecl: true,
				Variables: []*ast.VarDeclStmt{
					{
						Name: token.New(
							"score",
							kind.ID,
							token.NewPosition(filename, 1, 1),
						),
						Type:           &ast.BasicType{Kind: kind.U8_TYPE},
						NeedsInference: false,
						Value: &ast.LiteralExpr{
							Type:  &ast.BasicType{Kind: kind.INTEGER_LITERAL},
							Value: "10",
						},
					},
				},
			},
		},
		{
			input: "age int := 10;",
			varDecl: &ast.MultiVarStmt{
				IsDecl: true,
				Variables: []*ast.VarDeclStmt{
					{
						Name: token.New(
							"age",
							kind.ID,
							token.NewPosition(filename, 1, 1),
						),
						Type:           &ast.BasicType{Kind: kind.INT_TYPE},
						NeedsInference: false,
						Value: &ast.LiteralExpr{
							Type:  &ast.BasicType{Kind: kind.INTEGER_LITERAL},
							Value: "10",
						},
					},
				},
			},
		},
		// This code is not valid semantically (depends!), but
		// the parser needs to be able to analyze it.
		{
			input: "score SomeType := 10;",
			varDecl: &ast.MultiVarStmt{
				IsDecl: true,
				Variables: []*ast.VarDeclStmt{
					{
						Name: token.New(
							"score",
							kind.ID,
							token.NewPosition(filename, 1, 1),
						),
						Type: &ast.IdType{
							Name: token.New("SomeType", kind.ID, token.NewPosition(filename, 7, 1)),
						},
						NeedsInference: false,
						Value: &ast.LiteralExpr{
							Type:  &ast.BasicType{Kind: kind.INTEGER_LITERAL},
							Value: "10",
						},
					},
				},
			},
		},
		{
			input: "a, b := 10, 10;",
			varDecl: &ast.MultiVarStmt{
				IsDecl: true,
				Variables: []*ast.VarDeclStmt{
					{
						Name:           token.New("a", kind.ID, token.NewPosition(filename, 1, 1)),
						Type:           nil,
						NeedsInference: true,
						Value: &ast.LiteralExpr{
							Type:  &ast.BasicType{Kind: kind.INTEGER_LITERAL},
							Value: "10",
						},
					},
					{
						Name:           token.New("b", kind.ID, token.NewPosition(filename, 4, 1)),
						Type:           nil,
						NeedsInference: true,
						Value: &ast.LiteralExpr{
							Type:  &ast.BasicType{Kind: kind.INTEGER_LITERAL},
							Value: "10",
						},
					},
				},
			},
		},
	}

	for _, test := range tests {
		t.Run(fmt.Sprintf("TestVar('%s')", test.input), func(t *testing.T) {
			varDecl, err := parseVarDecl(filename, test.input)
			if err != nil {
				t.Fatal(err)
			}
			if !reflect.DeepEqual(varDecl, test.varDecl) {
				t.Fatalf("\nexp: %s\ngot: %s\n", test.varDecl, varDecl)
			}
		})
	}
}

// TODO(tests)
func TestFuncCallStmt(t *testing.T) {}

// TODO(tests)
func TestIfStmt(t *testing.T) {}

type syntaxErrorTest struct {
	input string
	diags []collector.Diag
}

func TestSyntaxErrors(t *testing.T) {
	filename := "test.tt"
	tests := []syntaxErrorTest{
		{
			input: "{",
			diags: []collector.Diag{
				{
					Message: "test.tt:1:1: unexpected non-declaration statement on global scope",
				},
			},
		},
		{
			input: "if",
			diags: []collector.Diag{
				{
					Message: "test.tt:1:1: unexpected non-declaration statement on global scope",
				},
			},
		},
		{
			input: "elif",
			diags: []collector.Diag{
				{
					Message: "test.tt:1:1: unexpected non-declaration statement on global scope",
				},
			},
		},
		{
			input: "else",
			diags: []collector.Diag{
				{
					Message: "test.tt:1:1: unexpected non-declaration statement on global scope",
				},
			},
		},
		// Function declaration
		{
			input: "fn name(){}",
			diags: nil, // no errors,
		},
		{
			input: "fn (){}",
			diags: []collector.Diag{
				{
					Message: "test.tt:1:4: expected name, not (",
				},
			},
		},
		{
			input: "fn",
			diags: []collector.Diag{
				{
					Message: "test.tt:1:3: expected name, not end of file",
				},
			},
		},
		{
			input: "fn name){}",
			diags: []collector.Diag{
				{
					Message: "test.tt:1:8: expected (, not )",
				},
			},
		},
		{
			input: "fn name",
			diags: []collector.Diag{
				{
					Message: "test.tt:1:8: expected (, not end of file",
				},
			},
		},
		{
			input: "fn name({}",
			diags: []collector.Diag{
				{
					Message: "test.tt:1:9: expected parameter or ), not {",
				},
			},
		},
		{
			input: "fn name(",
			diags: []collector.Diag{
				{
					Message: "test.tt:1:9: expected parameter or ), not end of file",
				},
			},
		},
		{
			input: "fn name(a, b int){}",
			diags: []collector.Diag{
				{
					Message: "test.tt:1:10: expected parameter type for 'a', not ,",
				},
			},
		},
		{
			input: "fn name(a",
			diags: []collector.Diag{
				{
					Message: "test.tt:1:10: expected parameter type for 'a', not end of file",
				},
			},
		},
		{
			input: "fn name(a int, ..., b int){}",
			diags: []collector.Diag{
				{
					Message: "test.tt:1:16: ... is only allowed at the end of parameter list",
				},
			},
		},
		{
			input: "fn name() }",
			diags: []collector.Diag{
				{
					Message: "test.tt:1:11: expected type or {, not }",
				},
			},
		},
		{
			input: "fn name()",
			diags: []collector.Diag{
				{
					Message: "test.tt:1:10: expected type or {, not end of file",
				},
			},
		},
		{
			input: "fn name() {",
			diags: []collector.Diag{
				{
					Message: "test.tt:1:12: expected statement or }, not end of file",
				},
			},
		},
		// External declarations
		{
			input: "extern libc {}",
			diags: nil, // no errors
		},
		{
			input: // no formatting
			`extern libc {
				fn method();
			}`,
			diags: nil, // no errors
		},
		{
			input: // no formatting
			`extern libc {
				fn method() i8;
			}`,
			diags: nil, // no errors
		},
		{
			input: // no formatting
			`extern libc {
			fn method() {}
			}`,
			diags: []collector.Diag{
				{
					Message: "test.tt:2:16: expected ; at the end of prototype, not {",
				},
			},
		},
		{
			input: "extern {}",
			diags: []collector.Diag{
				{
					Message: "test.tt:1:8: expected name, not {",
				},
			},
		},
		{
			input: "extern libc }",
			diags: []collector.Diag{
				{
					Message: "test.tt:1:13: expected {, not }",
				},
			},
		},
		{
			input: "extern libc {",
			diags: []collector.Diag{
				{
					Message: "test.tt:1:14: expected prototype or }, not end of file",
				},
			},
		},
		{
			input: "extern libc",
			diags: []collector.Diag{
				{
					Message: "test.tt:1:12: expected {, not end of file",
				},
			},
		},
		{
			input: // no formatting
			`extern libc {
			fn method() {}
			}`,
			diags: []collector.Diag{
				{
					Message: "test.tt:2:16: expected ; at the end of prototype, not {",
				},
			},
		},
		{
			input: // no formatting
			`extern libc {
			fn ();
			}`,
			diags: []collector.Diag{
				{
					Message: "test.tt:2:7: expected name, not (",
				},
			},
		},
		{
			input: // no formatting
			`extern libc {
			fn name();`,
			diags: []collector.Diag{
				{
					Message: "test.tt:2:14: expected prototype or }, not end of file",
				},
			},
		},
		{
			input: // no formatting
			`extern libc {
			fn name);
			}`,
			diags: []collector.Diag{
				{
					Message: "test.tt:2:11: expected (, not )",
				},
			},
		},
		{
			input: // no formatting
			`extern libc {
			fn name(;
			}`,
			diags: []collector.Diag{
				{
					Message: "test.tt:2:12: expected parameter or ), not ;",
				},
			},
		},
		{
			input: // no formatting
			`extern libc {
			fn name(a);
			}`,
			diags: []collector.Diag{
				{
					Message: "test.tt:2:13: expected parameter type for 'a', not )",
				},
			},
		},
		{
			input: // no formatting
			`extern libc {
			fn name(a;
			}`,
			diags: []collector.Diag{
				{
					Message: "test.tt:2:13: expected parameter type for 'a', not ;",
				},
			},
		},
		{
			input: // no formatting
			`extern libc {
			fn name(a int;
			}`,
			diags: []collector.Diag{
				{
					Message: "test.tt:2:17: expected parameter or ), not ;",
				},
			},
		},
	}

	for _, test := range tests {
		t.Run(fmt.Sprintf("TestSyntaxErrors('%s')", test.input), func(t *testing.T) {
			diagCollector := collector.New()
			reader := bufio.NewReader(strings.NewReader(test.input))

			lex := lexer.New(filename, reader, diagCollector)
			tokens, err := lex.Tokenize()
			if err != nil {
				t.Fatalf("unexpected error: %s", err)
			}

			parser := New(tokens, diagCollector)
			_, err = parser.Parse()

			if err != nil && len(parser.Collector.Diags) == 0 {
				t.Fatalf(
					"error detected, but diagnostic collector is empty.\nError: %s",
					err,
				)
			}

			if len(test.diags) != len(parser.Collector.Diags) {
				t.Fatalf(
					"expected to have %d diag(s), but got %d\n\ngot: %s\nexp: %s\n",
					len(test.diags),
					len(parser.Collector.Diags),
					parser.Collector.Diags,
					test.diags,
				)
			}
			if !reflect.DeepEqual(test.diags, parser.Collector.Diags) {
				t.Fatalf("\nexpected diags: %v\ngot diags: %v\n", test.diags, parser.Collector)
			}
		})
	}
}

func TestSyntaxErrorsOnBlock(t *testing.T) {
	filename := "test.tt"

	tests := []syntaxErrorTest{
		{
			input: "{",
			diags: []collector.Diag{
				{
					Message: "test.tt:1:2: expected statement or }, not end of file",
				},
			},
		},
		{
			input: `{
			return
			}`,
			diags: []collector.Diag{
				{
					Message: "test.tt:3:4: expected expression or ;, not }",
				},
			},
		},
		{
			input: `{
			return
			`,
			diags: []collector.Diag{
				{
					Message: "test.tt:3:4: expected expression or ;, not end of file",
				},
			},
		},
		{
			input: `{
			return 10
			}`,
			diags: []collector.Diag{
				{
					Message: "test.tt:3:4: expected ; at the end of statement, not }",
				},
			},
		},
		// TODO(tests): deal with id statement, such as function calls and variable
		// declarations
	}

	for _, test := range tests {
		t.Run(fmt.Sprintf("TestSyntaxErrorsOnBlock('%s')", test.input), func(t *testing.T) {
			collector := collector.New()

			reader := bufio.NewReader(strings.NewReader(test.input))
			lex := lexer.New(filename, reader, collector)
			tokens, err := lex.Tokenize()
			if err != nil {
				t.Fatalf("unexpected error: %s", err)
			}

			parser := New(tokens, collector)
			_, err = parser.parseBlock()
			if err == nil {
				t.Fatal("expected to have syntax errors, but got nothing")
			}

			if len(test.diags) != len(parser.Collector.Diags) {
				t.Fatalf(
					"expected to have %d diag(s), but got %d\n\ngot: %s\nexp: %s\n",
					len(test.diags),
					len(parser.Collector.Diags),
					parser.Collector.Diags,
					test.diags,
				)
			}
			if !reflect.DeepEqual(test.diags, parser.Collector.Diags) {
				t.Fatalf("\nexpected diags: %v\ngot diags: %v\n", test.diags, parser.Collector)
			}
		})
	}
}
