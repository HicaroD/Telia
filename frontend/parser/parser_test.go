package parser

import (
	"fmt"
	"reflect"
	"testing"

	"github.com/HicaroD/Telia/diagnostics"
	"github.com/HicaroD/Telia/frontend/ast"
	"github.com/HicaroD/Telia/frontend/lexer"
	"github.com/HicaroD/Telia/frontend/lexer/token"
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
				Name:  token.New([]byte("do_nothing"), token.ID, token.NewPosition(filename, 4, 1)),
				Params: &ast.FieldList{
					Open:   token.New(nil, token.OPEN_PAREN, token.NewPosition(filename, 14, 1)),
					Fields: nil,
					Close: token.New(
						nil,
						token.CLOSE_PAREN,
						token.NewPosition(filename, 15, 1),
					),
					IsVariadic: false,
				},
				RetType: &ast.BasicType{Kind: token.VOID_TYPE},
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
				Name:  token.New([]byte("do_nothing"), token.ID, token.NewPosition(filename, 4, 1)),
				Params: &ast.FieldList{
					Open: token.New(nil, token.OPEN_PAREN, token.NewPosition(filename, 14, 1)),
					Fields: []*ast.Field{
						{
							Name: token.New([]byte("a"), token.ID, token.NewPosition(filename, 15, 1)),
							Type: &ast.BasicType{Kind: token.BOOL_TYPE},
						},
					},
					Close: token.New(
						nil,
						token.CLOSE_PAREN,
						token.NewPosition(filename, 21, 1),
					),
					IsVariadic: false,
				},
				RetType: &ast.BasicType{Kind: token.VOID_TYPE},
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
				Name:  token.New([]byte("do_nothing"), token.ID, token.NewPosition(filename, 4, 1)),
				Params: &ast.FieldList{
					Open: token.New(nil, token.OPEN_PAREN, token.NewPosition(filename, 14, 1)),
					Fields: []*ast.Field{
						{
							Name: token.New([]byte("a"), token.ID, token.NewPosition(filename, 15, 1)),
							Type: &ast.BasicType{Kind: token.BOOL_TYPE},
						},
						{
							Name: token.New([]byte("b"), token.ID, token.NewPosition(filename, 23, 1)),
							Type: &ast.BasicType{Kind: token.I32_TYPE},
						},
					},
					Close: token.New(
						nil,
						token.CLOSE_PAREN,
						token.NewPosition(filename, 28, 1),
					),
					IsVariadic: false,
				},
				RetType: &ast.BasicType{Kind: token.VOID_TYPE},
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
				Name:  token.New([]byte("do_nothing"), token.ID, token.NewPosition(filename, 4, 1)),
				Params: &ast.FieldList{
					Open: token.New(nil, token.OPEN_PAREN, token.NewPosition(filename, 14, 1)),
					Fields: []*ast.Field{
						{
							Name: token.New([]byte("a"), token.ID, token.NewPosition(filename, 15, 1)),
							Type: &ast.BasicType{Kind: token.BOOL_TYPE},
						},
						{
							Name: token.New([]byte("b"), token.ID, token.NewPosition(filename, 23, 1)),
							Type: &ast.BasicType{Kind: token.I32_TYPE},
						},
					},
					Close: token.New(
						nil,
						token.CLOSE_PAREN,
						token.NewPosition(filename, 28, 1),
					),
					IsVariadic: false,
				},
				RetType: &ast.BasicType{Kind: token.I8_TYPE},
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
				Name:  token.New([]byte("do_nothing"), token.ID, token.NewPosition(filename, 4, 1)),
				Params: &ast.FieldList{
					Open: token.New(nil, token.OPEN_PAREN, token.NewPosition(filename, 14, 1)),
					Fields: []*ast.Field{
						{
							Name: token.New([]byte("a"), token.ID, token.NewPosition(filename, 15, 1)),
							Type: &ast.BasicType{Kind: token.BOOL_TYPE},
						},
						{
							Name: token.New([]byte("b"), token.ID, token.NewPosition(filename, 23, 1)),
							Type: &ast.BasicType{Kind: token.I32_TYPE},
						},
					},
					Close: token.New(
						nil,
						token.CLOSE_PAREN,
						token.NewPosition(filename, 28, 1),
					),
					IsVariadic: false,
				},
				RetType: &ast.BasicType{Kind: token.U8_TYPE},
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
				Name:  token.New([]byte("do_nothing"), token.ID, token.NewPosition(filename, 4, 1)),
				Params: &ast.FieldList{
					Open: token.New(nil, token.OPEN_PAREN, token.NewPosition(filename, 14, 1)),
					Fields: []*ast.Field{
						{
							Name: token.New([]byte("a"), token.ID, token.NewPosition(filename, 15, 1)),
							Type: &ast.BasicType{Kind: token.BOOL_TYPE},
						},
						{
							Name: token.New([]byte("b"), token.ID, token.NewPosition(filename, 23, 1)),
							Type: &ast.BasicType{Kind: token.I32_TYPE},
						},
					},
					Close: token.New(
						nil,
						token.CLOSE_PAREN,
						token.NewPosition(filename, 28, 1),
					),
					IsVariadic: false,
				},
				RetType: &ast.BasicType{Kind: token.I16_TYPE},
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
				Name:  token.New([]byte("do_nothing"), token.ID, token.NewPosition(filename, 4, 1)),
				Params: &ast.FieldList{
					Open: token.New(nil, token.OPEN_PAREN, token.NewPosition(filename, 14, 1)),
					Fields: []*ast.Field{
						{
							Name: token.New([]byte("a"), token.ID, token.NewPosition(filename, 15, 1)),
							Type: &ast.BasicType{Kind: token.BOOL_TYPE},
						},
						{
							Name: token.New([]byte("b"), token.ID, token.NewPosition(filename, 23, 1)),
							Type: &ast.BasicType{Kind: token.I32_TYPE},
						},
					},
					Close: token.New(
						nil,
						token.CLOSE_PAREN,
						token.NewPosition(filename, 28, 1),
					),
					IsVariadic: false,
				},
				RetType: &ast.BasicType{Kind: token.U16_TYPE},
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
				Name:  token.New([]byte("do_nothing"), token.ID, token.NewPosition(filename, 4, 1)),
				Params: &ast.FieldList{
					Open: token.New(nil, token.OPEN_PAREN, token.NewPosition(filename, 14, 1)),
					Fields: []*ast.Field{
						{
							Name: token.New([]byte("a"), token.ID, token.NewPosition(filename, 15, 1)),
							Type: &ast.BasicType{Kind: token.BOOL_TYPE},
						},
						{
							Name: token.New([]byte("b"), token.ID, token.NewPosition(filename, 23, 1)),
							Type: &ast.BasicType{Kind: token.I32_TYPE},
						},
					},
					Close: token.New(
						nil,
						token.CLOSE_PAREN,
						token.NewPosition(filename, 28, 1),
					),
					IsVariadic: false,
				},
				RetType: &ast.BasicType{Kind: token.I32_TYPE},
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
				Name:  token.New([]byte("do_nothing"), token.ID, token.NewPosition(filename, 4, 1)),
				Params: &ast.FieldList{
					Open: token.New(nil, token.OPEN_PAREN, token.NewPosition(filename, 14, 1)),
					Fields: []*ast.Field{
						{
							Name: token.New([]byte("a"), token.ID, token.NewPosition(filename, 15, 1)),
							Type: &ast.BasicType{Kind: token.BOOL_TYPE},
						},
						{
							Name: token.New([]byte("b"), token.ID, token.NewPosition(filename, 23, 1)),
							Type: &ast.BasicType{Kind: token.I32_TYPE},
						},
					},
					Close: token.New(
						nil,
						token.CLOSE_PAREN,
						token.NewPosition(filename, 28, 1),
					),
					IsVariadic: false,
				},
				RetType: &ast.BasicType{Kind: token.U32_TYPE},
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
				Name:  token.New([]byte("do_nothing"), token.ID, token.NewPosition(filename, 4, 1)),
				Params: &ast.FieldList{
					Open: token.New(nil, token.OPEN_PAREN, token.NewPosition(filename, 14, 1)),
					Fields: []*ast.Field{
						{
							Name: token.New([]byte("a"), token.ID, token.NewPosition(filename, 15, 1)),
							Type: &ast.BasicType{Kind: token.BOOL_TYPE},
						},
						{
							Name: token.New([]byte("b"), token.ID, token.NewPosition(filename, 23, 1)),
							Type: &ast.BasicType{Kind: token.I32_TYPE},
						},
					},
					Close: token.New(
						nil,
						token.CLOSE_PAREN,
						token.NewPosition(filename, 28, 1),
					),
					IsVariadic: false,
				},
				RetType: &ast.BasicType{Kind: token.I64_TYPE},
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
				Name:  token.New([]byte("do_nothing"), token.ID, token.NewPosition(filename, 4, 1)),
				Params: &ast.FieldList{
					Open: token.New(nil, token.OPEN_PAREN, token.NewPosition(filename, 14, 1)),
					Fields: []*ast.Field{
						{
							Name: token.New([]byte("a"), token.ID, token.NewPosition(filename, 15, 1)),
							Type: &ast.BasicType{Kind: token.BOOL_TYPE},
						},
						{
							Name: token.New([]byte("b"), token.ID, token.NewPosition(filename, 23, 1)),
							Type: &ast.BasicType{Kind: token.I32_TYPE},
						},
					},
					Close: token.New(
						nil,
						token.CLOSE_PAREN,
						token.NewPosition(filename, 28, 1),
					),
					IsVariadic: false,
				},
				RetType: &ast.BasicType{Kind: token.U64_TYPE},
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
				Name:  token.New([]byte("do_nothing"), token.ID, token.NewPosition(filename, 4, 1)),
				Params: &ast.FieldList{
					Open: token.New(nil, token.OPEN_PAREN, token.NewPosition(filename, 14, 1)),
					Fields: []*ast.Field{
						{
							Name: token.New([]byte("a"), token.ID, token.NewPosition(filename, 15, 1)),
							Type: &ast.BasicType{Kind: token.BOOL_TYPE},
						},
						{
							Name: token.New([]byte("b"), token.ID, token.NewPosition(filename, 23, 1)),
							Type: &ast.BasicType{Kind: token.I32_TYPE},
						},
					},
					Close: token.New(
						nil,
						token.CLOSE_PAREN,
						token.NewPosition(filename, 28, 1),
					),
					IsVariadic: false,
				},
				RetType: &ast.BasicType{Kind: token.BOOL_TYPE},
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
				Name:  token.New([]byte("do_nothing"), token.ID, token.NewPosition(filename, 4, 1)),
				Params: &ast.FieldList{
					Open: token.New(nil, token.OPEN_PAREN, token.NewPosition(filename, 14, 1)),
					Fields: []*ast.Field{
						{
							Name: token.New([]byte("a"), token.ID, token.NewPosition(filename, 15, 1)),
							Type: &ast.BasicType{Kind: token.BOOL_TYPE},
						},
						{
							Name: token.New([]byte("b"), token.ID, token.NewPosition(filename, 23, 1)),
							Type: &ast.BasicType{Kind: token.I32_TYPE},
						},
					},
					Close: token.New(
						nil,
						token.CLOSE_PAREN,
						token.NewPosition(filename, 28, 1),
					),
					IsVariadic: false,
				},
				RetType: &ast.PointerType{Type: &ast.BasicType{Kind: token.I8_TYPE}},
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
			universeScope := ast.NewScope(nil)
			moduleScope := ast.NewScope(universeScope)
			fileScope := ast.NewScope(moduleScope)

			fnDecl, err := parseFnDeclFrom(filename, test.input, fileScope)
			if err != nil {
				t.Fatal(err)
			}
			test.node.Scope = ast.NewScope(fileScope)

			if !reflect.DeepEqual(fnDecl, test.node) {
				t.Fatal("Function declarations are not the same")
			}
		})
	}
}

type forLoopTest struct {
	input string
	node  *ast.ForLoop
}

func TestForLoop(t *testing.T) {
	filename := "test.tt"
	tests := []forLoopTest{
		{
			input: "for (i := 0; i < 10; i = i + 1) {}",
			node: &ast.ForLoop{
				Init: &ast.VarStmt{
					Decl: true,
					Name: token.New(
						[]byte("i"),
						token.ID,
						token.NewPosition(filename, 6, 1),
					),
					Type:           nil,
					NeedsInference: true,
					Value: &ast.LiteralExpr{
						Type:  &ast.BasicType{Kind: token.INTEGER_LITERAL},
						Value: []byte("0"),
					},
				},
				Cond: &ast.BinaryExpr{
					Left: &ast.IdExpr{
						Name: token.New([]byte("i"), token.ID, token.NewPosition(filename, 14, 1)),
					},
					Op: token.LESS,
					Right: &ast.LiteralExpr{
						Type:  &ast.BasicType{Kind: token.INTEGER_LITERAL},
						Value: []byte("10"),
					},
				},
				Update: &ast.VarStmt{
					Decl: false,
					Name: token.New(
						[]byte("i"),
						token.ID,
						token.NewPosition(filename, 22, 1),
					),
					Type:           nil,
					NeedsInference: true,
					Value: &ast.BinaryExpr{
						Left: &ast.IdExpr{
							Name: token.New(
								[]byte("i"),
								token.ID,
								token.NewPosition(filename, 26, 1),
							),
						},
						Op: token.PLUS,
						Right: &ast.LiteralExpr{
							Type:  &ast.BasicType{Kind: token.INTEGER_LITERAL},
							Value: []byte("1"),
						},
					},
				},
				Block: &ast.BlockStmt{
					OpenCurly:  token.NewPosition(filename, 33, 1),
					Statements: nil,
					CloseCurly: token.NewPosition(filename, 34, 1),
				},
			},
		},
	}

	for _, test := range tests {
		t.Run(fmt.Sprintf("TestForLoop('%s')", test.input), func(t *testing.T) {
			forLoop, err := ParseForLoopFrom(test.input, filename)
			if err != nil {
				t.Fatal(err)
			}

			if !reflect.DeepEqual(forLoop, test.node) {
				t.Fatalf("\nexp: %s\ngot: %s\n", test.node, forLoop)
			}
		})
	}
}

type whileLoopTest struct {
	input string
	node  *ast.WhileLoop
}

func TestWhileLoop(t *testing.T) {
	filename := "test.tt"
	tests := []*whileLoopTest{
		{
			input: "while true {}",
			node: &ast.WhileLoop{
				Cond: &ast.LiteralExpr{
					Type:  &ast.BasicType{Kind: token.TRUE_BOOL_LITERAL},
					Value: []byte("true"),
				},
				Block: &ast.BlockStmt{
					OpenCurly:  token.NewPosition(filename, 12, 1),
					Statements: nil,
					CloseCurly: token.NewPosition(filename, 13, 1),
				},
			},
		},
	}

	for _, test := range tests {
		t.Run(fmt.Sprintf("TestWhileLoop('%s')", test.input), func(t *testing.T) {
			whileLoop, err := ParseWhileLoopFrom(test.input, filename)
			if err != nil {
				t.Fatal(err)
			}

			if !reflect.DeepEqual(whileLoop, test.node) {
				t.Fatalf("\nexp: %s\ngot: %s\n", test.node, whileLoop)
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
			node:  &ast.LiteralExpr{Value: []byte("1"), Type: &ast.BasicType{Kind: token.INTEGER_LITERAL}},
		},
		{
			input: "true",
			node: &ast.LiteralExpr{
				Value: []byte("true"),
				Type:  &ast.BasicType{Kind: token.TRUE_BOOL_LITERAL},
			},
		},
		{
			input: "false",
			node: &ast.LiteralExpr{
				Value: []byte("false"),
				Type:  &ast.BasicType{Kind: token.FALSE_BOOL_LITERAL},
			},
		},
		{
			input: "\"Hello, world\"",
			node: &ast.LiteralExpr{
				Value: []byte("Hello, world"),
				Type:  &ast.BasicType{Kind: token.STRING_LITERAL},
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
				Op: token.MINUS,
				Value: &ast.LiteralExpr{
					Value: []byte("1"),
					Type:  &ast.BasicType{Kind: token.INTEGER_LITERAL},
				},
			},
		},
		{
			input: "not true",
			node: &ast.UnaryExpr{
				Op: token.NOT,
				Value: &ast.LiteralExpr{
					Value: []byte("true"),
					Type:  &ast.BasicType{Kind: token.TRUE_BOOL_LITERAL},
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
					Value: []byte("1"),
					Type:  &ast.BasicType{Kind: token.INTEGER_LITERAL},
				},
				Op: token.PLUS,
				Right: &ast.LiteralExpr{
					Value: []byte("1"),
					Type:  &ast.BasicType{Kind: token.INTEGER_LITERAL},
				},
			},
		},
		{
			input: "2 - 1",
			node: &ast.BinaryExpr{
				Left: &ast.LiteralExpr{
					Value: []byte("2"),
					Type:  &ast.BasicType{Kind: token.INTEGER_LITERAL},
				},
				Op: token.MINUS,
				Right: &ast.LiteralExpr{
					Value: []byte("1"),
					Type:  &ast.BasicType{Kind: token.INTEGER_LITERAL},
				},
			},
		},
		{
			input: "5 * 10",
			node: &ast.BinaryExpr{
				Left: &ast.LiteralExpr{
					Value: []byte("5"),
					Type:  &ast.BasicType{Kind: token.INTEGER_LITERAL},
				},
				Op: token.STAR,
				Right: &ast.LiteralExpr{
					Value: []byte("10"),
					Type:  &ast.BasicType{Kind: token.INTEGER_LITERAL},
				},
			},
		},
		{
			input: "3 + 4 * 5",
			node: &ast.BinaryExpr{
				Left: &ast.LiteralExpr{
					Value: []byte("3"),
					Type:  &ast.BasicType{Kind: token.INTEGER_LITERAL},
				},
				Op: token.PLUS,
				Right: &ast.BinaryExpr{
					Left: &ast.LiteralExpr{
						Value: []byte("4"),
						Type:  &ast.BasicType{Kind: token.INTEGER_LITERAL},
					},
					Op: token.STAR,
					Right: &ast.LiteralExpr{
						Value: []byte("5"),
						Type:  &ast.BasicType{Kind: token.INTEGER_LITERAL},
					},
				},
			},
		},
		{
			input: "3 + (4 * 5)",
			node: &ast.BinaryExpr{
				Left: &ast.LiteralExpr{
					Value: []byte("3"),
					Type:  &ast.BasicType{Kind: token.INTEGER_LITERAL},
				},
				Op: token.PLUS,
				Right: &ast.BinaryExpr{
					Left: &ast.LiteralExpr{
						Value: []byte("4"),
						Type:  &ast.BasicType{Kind: token.INTEGER_LITERAL},
					},
					Op: token.STAR,
					Right: &ast.LiteralExpr{
						Value: []byte("5"),
						Type:  &ast.BasicType{Kind: token.INTEGER_LITERAL},
					},
				},
			},
		},
		{
			input: "10 / 1",
			node: &ast.BinaryExpr{
				Left: &ast.LiteralExpr{
					Value: []byte("10"),
					Type:  &ast.BasicType{Kind: token.INTEGER_LITERAL},
				},
				Op: token.SLASH,
				Right: &ast.LiteralExpr{
					Value: []byte("1"),
					Type:  &ast.BasicType{Kind: token.INTEGER_LITERAL},
				},
			},
		},
		{
			input: "6 / 3 - 1",
			node: &ast.BinaryExpr{
				Left: &ast.BinaryExpr{
					Left: &ast.LiteralExpr{
						Value: []byte("6"),
						Type:  &ast.BasicType{Kind: token.INTEGER_LITERAL},
					},
					Op: token.SLASH,
					Right: &ast.LiteralExpr{
						Value: []byte("3"),
						Type:  &ast.BasicType{Kind: token.INTEGER_LITERAL},
					},
				},
				Op: token.MINUS,
				Right: &ast.LiteralExpr{
					Value: []byte("1"),
					Type:  &ast.BasicType{Kind: token.INTEGER_LITERAL},
				},
			},
		},
		{
			input: "6 / (3 - 1)",
			node: &ast.BinaryExpr{
				Left: &ast.LiteralExpr{
					Value: []byte("6"),
					Type:  &ast.BasicType{Kind: token.INTEGER_LITERAL},
				},
				Op: token.SLASH,
				Right: &ast.BinaryExpr{
					Left: &ast.LiteralExpr{
						Value: []byte("3"),
						Type:  &ast.BasicType{Kind: token.INTEGER_LITERAL},
					},
					Op: token.MINUS,
					Right: &ast.LiteralExpr{
						Value: []byte("1"),
						Type:  &ast.BasicType{Kind: token.INTEGER_LITERAL},
					},
				},
			},
		},
		{
			input: "1 / (1 + 1)",
			node: &ast.BinaryExpr{
				Left: &ast.LiteralExpr{
					Value: []byte("1"),
					Type:  &ast.BasicType{Kind: token.INTEGER_LITERAL},
				},
				Op: token.SLASH,
				Right: &ast.BinaryExpr{
					Left: &ast.LiteralExpr{
						Value: []byte("1"),
						Type:  &ast.BasicType{Kind: token.INTEGER_LITERAL},
					},
					Op: token.PLUS,
					Right: &ast.LiteralExpr{
						Value: []byte("1"),
						Type:  &ast.BasicType{Kind: token.INTEGER_LITERAL},
					},
				},
			},
		},
		{
			input: "1 > 1",
			node: &ast.BinaryExpr{
				Left: &ast.LiteralExpr{
					Value: []byte("1"),
					Type:  &ast.BasicType{Kind: token.INTEGER_LITERAL},
				},
				Op: token.GREATER,
				Right: &ast.LiteralExpr{
					Value: []byte("1"),
					Type:  &ast.BasicType{Kind: token.INTEGER_LITERAL},
				},
			},
		},
		{
			input: "1 >= 1",
			node: &ast.BinaryExpr{
				Left: &ast.LiteralExpr{
					Value: []byte("1"),
					Type:  &ast.BasicType{Kind: token.INTEGER_LITERAL},
				},
				Op: token.GREATER_EQ,
				Right: &ast.LiteralExpr{
					Value: []byte("1"),
					Type:  &ast.BasicType{Kind: token.INTEGER_LITERAL},
				},
			},
		},
		{
			input: "1 < 1",
			node: &ast.BinaryExpr{
				Left: &ast.LiteralExpr{
					Value: []byte("1"),
					Type:  &ast.BasicType{Kind: token.INTEGER_LITERAL},
				},
				Op: token.LESS,
				Right: &ast.LiteralExpr{
					Value: []byte("1"),
					Type:  &ast.BasicType{Kind: token.INTEGER_LITERAL},
				},
			},
		},
		{
			input: "1 <= 1",
			node: &ast.BinaryExpr{
				Left: &ast.LiteralExpr{
					Value: []byte("1"),
					Type:  &ast.BasicType{Kind: token.INTEGER_LITERAL},
				},
				Op: token.LESS_EQ,
				Right: &ast.LiteralExpr{
					Value: []byte("1"),
					Type:  &ast.BasicType{Kind: token.INTEGER_LITERAL},
				},
			},
		},
		{
			// not 1 > 1 is invalid in Golang
			input: "not (1 > 1)",
			node: &ast.UnaryExpr{
				Op: token.NOT,
				Value: &ast.BinaryExpr{
					Left: &ast.LiteralExpr{
						Value: []byte("1"),
						Type:  &ast.BasicType{Kind: token.INTEGER_LITERAL},
					},
					Op: token.GREATER,
					Right: &ast.LiteralExpr{
						Value: []byte("1"),
						Type:  &ast.BasicType{Kind: token.INTEGER_LITERAL},
					},
				},
			},
		},
		{
			input: "1 > 1 and 1 > 1",
			node: &ast.BinaryExpr{
				Left: &ast.BinaryExpr{
					Left: &ast.LiteralExpr{
						Value: []byte("1"),
						Type:  &ast.BasicType{Kind: token.INTEGER_LITERAL},
					},
					Op: token.GREATER,
					Right: &ast.LiteralExpr{
						Value: []byte("1"),
						Type:  &ast.BasicType{Kind: token.INTEGER_LITERAL},
					},
				},
				Op: token.AND,
				Right: &ast.BinaryExpr{
					Left: &ast.LiteralExpr{
						Value: []byte("1"),
						Type:  &ast.BasicType{Kind: token.INTEGER_LITERAL},
					},
					Op: token.GREATER,
					Right: &ast.LiteralExpr{
						Value: []byte("1"),
						Type:  &ast.BasicType{Kind: token.INTEGER_LITERAL},
					},
				},
			},
		},
		{
			input: "1 > 1 or 1 > 1",
			node: &ast.BinaryExpr{
				Left: &ast.BinaryExpr{
					Left: &ast.LiteralExpr{
						Value: []byte("1"),
						Type:  &ast.BasicType{Kind: token.INTEGER_LITERAL},
					},
					Op: token.GREATER,
					Right: &ast.LiteralExpr{
						Value: []byte("1"),
						Type:  &ast.BasicType{Kind: token.INTEGER_LITERAL},
					},
				},
				Op: token.OR,
				Right: &ast.BinaryExpr{
					Left: &ast.LiteralExpr{
						Value: []byte("1"),
						Type:  &ast.BasicType{Kind: token.INTEGER_LITERAL},
					},
					Op: token.GREATER,
					Right: &ast.LiteralExpr{
						Value: []byte("1"),
						Type:  &ast.BasicType{Kind: token.INTEGER_LITERAL},
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
							Name: token.New([]byte("celsius"), token.ID, token.NewPosition("test.tt", 1, 1)),
						},
						Op: token.STAR,
						Right: &ast.LiteralExpr{
							Value: []byte("9"),
							Type:  &ast.BasicType{Kind: token.INTEGER_LITERAL},
						},
					},
					Op: token.SLASH,
					Right: &ast.LiteralExpr{
						Value: []byte("5"),
						Type:  &ast.BasicType{Kind: token.INTEGER_LITERAL},
					},
				},
				Op: token.PLUS,
				Right: &ast.LiteralExpr{
					Value: []byte("32"),
					Type:  &ast.BasicType{Kind: token.INTEGER_LITERAL},
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
								[]byte("get_celsius"),
								token.ID,
								token.NewPosition(filename, 1, 1),
							),
							Args: nil,
						},
						Op: token.STAR,
						Right: &ast.LiteralExpr{
							Value: []byte("9"),
							Type:  &ast.BasicType{Kind: token.INTEGER_LITERAL},
						},
					},
					Op: token.SLASH,
					Right: &ast.LiteralExpr{
						Value: []byte("5"),
						Type:  &ast.BasicType{Kind: token.INTEGER_LITERAL},
					},
				},
				Op: token.PLUS,
				Right: &ast.LiteralExpr{
					Value: []byte("32"),
					Type:  &ast.BasicType{Kind: token.INTEGER_LITERAL},
				},
			},
		},
		{
			input: "1 > 1 > 1",
			node: &ast.BinaryExpr{
				Left: &ast.BinaryExpr{
					Left: &ast.LiteralExpr{
						Value: []byte("1"),
						Type:  &ast.BasicType{Kind: token.INTEGER_LITERAL},
					},
					Op: token.GREATER,
					Right: &ast.LiteralExpr{
						Value: []byte("1"),
						Type:  &ast.BasicType{Kind: token.INTEGER_LITERAL},
					},
				},
				Op: token.GREATER,
				Right: &ast.LiteralExpr{
					Value: []byte("1"),
					Type:  &ast.BasicType{Kind: token.INTEGER_LITERAL},
				},
			},
		},
		{
			input: "n == 1",
			node: &ast.BinaryExpr{
				Left: &ast.IdExpr{
					Name: token.New([]byte("n"), token.ID, token.NewPosition("test.tt", 1, 1)),
				},
				Op: token.EQUAL_EQUAL,
				Right: &ast.LiteralExpr{
					Value: []byte("1"),
					Type:  &ast.BasicType{Kind: token.INTEGER_LITERAL},
				},
			},
		},
		{
			input: "n == 1 or n == 2",
			node: &ast.BinaryExpr{
				Left: &ast.BinaryExpr{
					Left: &ast.IdExpr{
						Name: token.New([]byte("n"), token.ID, token.NewPosition("test.tt", 1, 1)),
					},
					Op: token.EQUAL_EQUAL,
					Right: &ast.LiteralExpr{
						Type:  &ast.BasicType{Kind: token.INTEGER_LITERAL},
						Value: []byte("1"),
					},
				},
				Op: token.OR,
				Right: &ast.BinaryExpr{
					Left: &ast.IdExpr{
						Name: token.New([]byte("n"), token.ID, token.NewPosition("test.tt", 11, 1)),
					},
					Op: token.EQUAL_EQUAL,
					Right: &ast.LiteralExpr{
						Type:  &ast.BasicType{Kind: token.INTEGER_LITERAL},
						Value: []byte("2"),
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
							Value: []byte("1"),
							Type:  &ast.BasicType{Kind: token.INTEGER_LITERAL},
						},
						Op: token.PLUS,
						Right: &ast.LiteralExpr{
							Value: []byte("1"),
							Type:  &ast.BasicType{Kind: token.INTEGER_LITERAL},
						},
					},
					Op: token.GREATER,
					Right: &ast.LiteralExpr{
						Value: []byte("2"),
						Type:  &ast.BasicType{Kind: token.INTEGER_LITERAL},
					},
				},
				Op: token.AND,
				Right: &ast.BinaryExpr{
					Left: &ast.LiteralExpr{
						Value: []byte("1"),
						Type:  &ast.BasicType{Kind: token.INTEGER_LITERAL},
					},
					Op: token.EQUAL_EQUAL,
					Right: &ast.LiteralExpr{
						Value: []byte("1"),
						Type:  &ast.BasicType{Kind: token.INTEGER_LITERAL},
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
							Value: []byte("true"),
							Type:  &ast.BasicType{Kind: token.TRUE_BOOL_LITERAL},
						},
						Op: token.AND,
						Right: &ast.LiteralExpr{
							Value: []byte("true"),
							Type:  &ast.BasicType{Kind: token.TRUE_BOOL_LITERAL},
						},
					},
					Op: token.AND,
					Right: &ast.LiteralExpr{
						Value: []byte("true"),
						Type:  &ast.BasicType{Kind: token.TRUE_BOOL_LITERAL},
					},
				},
				Op: token.AND,
				Right: &ast.LiteralExpr{
					Value: []byte("true"),
					Type:  &ast.BasicType{Kind: token.TRUE_BOOL_LITERAL},
				},
			},
		},
		{
			input: "(((true and true) and true) and true)",
			node: &ast.BinaryExpr{
				Left: &ast.BinaryExpr{
					Left: &ast.BinaryExpr{
						Left: &ast.LiteralExpr{
							Value: []byte("true"),
							Type:  &ast.BasicType{Kind: token.TRUE_BOOL_LITERAL},
						},
						Op: token.AND,
						Right: &ast.LiteralExpr{
							Value: []byte("true"),
							Type:  &ast.BasicType{Kind: token.TRUE_BOOL_LITERAL},
						},
					},
					Op: token.AND,
					Right: &ast.LiteralExpr{
						Value: []byte("true"),
						Type:  &ast.BasicType{Kind: token.TRUE_BOOL_LITERAL},
					},
				},
				Op: token.AND,
				Right: &ast.LiteralExpr{
					Value: []byte("true"),
					Type:  &ast.BasicType{Kind: token.TRUE_BOOL_LITERAL},
				},
			},
		},
		{
			input: "1 + multiply_by_2(10)",
			node: &ast.BinaryExpr{
				Left: &ast.LiteralExpr{
					Value: []byte("1"),
					Type:  &ast.BasicType{Kind: token.INTEGER_LITERAL},
				},
				Op: token.PLUS,
				Right: &ast.FunctionCall{
					Name: token.New([]byte("multiply_by_2"), token.ID, token.NewPosition(filename, 5, 1)),
					Args: []ast.Expr{
						&ast.LiteralExpr{
							Value: []byte("10"),
							Type:  &ast.BasicType{Kind: token.INTEGER_LITERAL},
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
					Name: token.New([]byte("first"), token.ID, token.NewPosition(filename, 1, 1)),
				},
				Right: &ast.FieldAccess{
					Left: &ast.IdExpr{
						Name: token.New([]byte("second"), token.ID, token.NewPosition(filename, 7, 1)),
					},
					Right: &ast.IdExpr{
						Name: token.New([]byte("third"), token.ID, token.NewPosition(filename, 14, 1)),
					},
				},
			},
		},
		{
			input: "first.second",
			node: &ast.FieldAccess{
				Left: &ast.IdExpr{
					Name: token.New([]byte("first"), token.ID, token.NewPosition(filename, 1, 1)),
				},
				Right: &ast.IdExpr{
					Name: token.New([]byte("second"), token.ID, token.NewPosition(filename, 7, 1)),
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
	varDecl ast.Stmt
}

func TestVar(t *testing.T) {
	filename := "test.tt"
	tests := []varDeclTest{
		{
			input: "age := 10;",
			varDecl: &ast.VarStmt{
				Decl: true,
				Name: token.New(
					[]byte("age"),
					token.ID,
					token.NewPosition(filename, 1, 1),
				),
				Type:           nil,
				NeedsInference: true,
				Value: &ast.LiteralExpr{
					Type:  &ast.BasicType{Kind: token.INTEGER_LITERAL},
					Value: []byte("10"),
				},
			},
		},
		{
			input: "score u8 := 10;",
			varDecl: &ast.VarStmt{
				Decl: true,
				Name: token.New(
					[]byte("score"),
					token.ID,
					token.NewPosition(filename, 1, 1),
				),
				Type:           &ast.BasicType{Kind: token.U8_TYPE},
				NeedsInference: false,
				Value: &ast.LiteralExpr{
					Type:  &ast.BasicType{Kind: token.INTEGER_LITERAL},
					Value: []byte("10"),
				},
			},
		},
		{
			input: "age int := 10;",
			varDecl: &ast.VarStmt{
				Decl: true,
				Name: token.New(
					[]byte("age"),
					token.ID,
					token.NewPosition(filename, 1, 1),
				),
				Type:           &ast.BasicType{Kind: token.INT_TYPE},
				NeedsInference: false,
				Value: &ast.LiteralExpr{
					Type:  &ast.BasicType{Kind: token.INTEGER_LITERAL},
					Value: []byte("10"),
				},
			},
		},
		// This code is not valid semantically (depends!), but
		// the parser needs to be able to analyze it.
		{
			input: "score SomeType := 10;",
			varDecl: &ast.VarStmt{
				Decl: true,
				Name: token.New(
					[]byte("score"),
					token.ID,
					token.NewPosition(filename, 1, 1),
				),
				Type: &ast.IdType{
					Name: token.New([]byte("SomeType"), token.ID, token.NewPosition(filename, 7, 1)),
				},
				NeedsInference: false,
				Value: &ast.LiteralExpr{
					Type:  &ast.BasicType{Kind: token.INTEGER_LITERAL},
					Value: []byte("10"),
				},
			},
		},
		{
			input: "a, b := 10, 10;",
			varDecl: &ast.MultiVarStmt{
				IsDecl: true,
				Variables: []*ast.VarStmt{
					{
						Decl:           true,
						Name:           token.New([]byte("a"), token.ID, token.NewPosition(filename, 1, 1)),
						Type:           nil,
						NeedsInference: true,
						Value: &ast.LiteralExpr{
							Type:  &ast.BasicType{Kind: token.INTEGER_LITERAL},
							Value: []byte("10"),
						},
					},
					{
						Decl:           true,
						Name:           token.New([]byte("b"), token.ID, token.NewPosition(filename, 4, 1)),
						Type:           nil,
						NeedsInference: true,
						Value: &ast.LiteralExpr{
							Type:  &ast.BasicType{Kind: token.INTEGER_LITERAL},
							Value: []byte("10"),
						},
					},
				},
			},
		},
		{
			input: "a, b = 10, 10;",
			varDecl: &ast.MultiVarStmt{
				IsDecl: false,
				Variables: []*ast.VarStmt{
					{
						Name:           token.New([]byte("a"), token.ID, token.NewPosition(filename, 1, 1)),
						Type:           nil,
						NeedsInference: true,
						Value: &ast.LiteralExpr{
							Type:  &ast.BasicType{Kind: token.INTEGER_LITERAL},
							Value: []byte("10"),
						},
					},
					{
						Name:           token.New([]byte("b"), token.ID, token.NewPosition(filename, 4, 1)),
						Type:           nil,
						NeedsInference: true,
						Value: &ast.LiteralExpr{
							Type:  &ast.BasicType{Kind: token.INTEGER_LITERAL},
							Value: []byte("10"),
						},
					},
				},
			},
		},
		{
			input: "a = 10;",
			varDecl: &ast.VarStmt{
				Decl:           false,
				Name:           token.New([]byte("a"), token.ID, token.NewPosition(filename, 1, 1)),
				Type:           nil,
				NeedsInference: true,
				Value: &ast.LiteralExpr{
					Type:  &ast.BasicType{Kind: token.INTEGER_LITERAL},
					Value: []byte("10"),
				},
			},
		},
		{
			input: "a, b u8 = 10, 10;",
			varDecl: &ast.MultiVarStmt{
				IsDecl: false,
				Variables: []*ast.VarStmt{
					{
						Name:           token.New([]byte("a"), token.ID, token.NewPosition(filename, 1, 1)),
						Type:           nil,
						NeedsInference: true,
						Value: &ast.LiteralExpr{
							Type:  &ast.BasicType{Kind: token.INTEGER_LITERAL},
							Value: []byte("10"),
						},
					},
					{
						Name:           token.New([]byte("b"), token.ID, token.NewPosition(filename, 4, 1)),
						Type:           &ast.BasicType{Kind: token.U8_TYPE},
						NeedsInference: false,
						Value: &ast.LiteralExpr{
							Type:  &ast.BasicType{Kind: token.INTEGER_LITERAL},
							Value: []byte("10"),
						},
					},
				},
			},
		},
	}

	for _, test := range tests {
		t.Run(fmt.Sprintf("TestVar('%s')", test.input), func(t *testing.T) {
			varDecl, err := parseVarFrom(filename, test.input)
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

type diagErrorTest struct {
	input string
	diags []diagnostics.Diag
}

func TestSyntaxErrors(t *testing.T) {
	filename := "test.tt"
	tests := []diagErrorTest{
		{
			input: "{",
			diags: []diagnostics.Diag{
				{
					Message: "test.tt:1:1: unexpected non-declaration statement on global scope",
				},
			},
		},
		{
			input: "if",
			diags: []diagnostics.Diag{
				{
					Message: "test.tt:1:1: unexpected non-declaration statement on global scope",
				},
			},
		},
		{
			input: "elif",
			diags: []diagnostics.Diag{
				{
					Message: "test.tt:1:1: unexpected non-declaration statement on global scope",
				},
			},
		},
		{
			input: "else",
			diags: []diagnostics.Diag{
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
			diags: []diagnostics.Diag{
				{
					Message: "test.tt:1:4: expected name, not (",
				},
			},
		},
		{
			input: "fn",
			diags: []diagnostics.Diag{
				{
					Message: "test.tt:1:3: expected name, not end of file",
				},
			},
		},
		{
			input: "fn name){}",
			diags: []diagnostics.Diag{
				{
					Message: "test.tt:1:8: expected (, not )",
				},
			},
		},
		{
			input: "fn name",
			diags: []diagnostics.Diag{
				{
					Message: "test.tt:1:8: expected (, not end of file",
				},
			},
		},
		{
			input: "fn name({}",
			diags: []diagnostics.Diag{
				{
					Message: "test.tt:1:9: expected parameter or ), not {",
				},
			},
		},
		{
			input: "fn name(",
			diags: []diagnostics.Diag{
				{
					Message: "test.tt:1:9: expected parameter or ), not end of file",
				},
			},
		},
		{
			input: "fn name(a, b int){}",
			diags: []diagnostics.Diag{
				{
					Message: "test.tt:1:10: expected parameter type for 'a', not ,",
				},
			},
		},
		{
			input: "fn name(a",
			diags: []diagnostics.Diag{
				{
					Message: "test.tt:1:10: expected parameter type for 'a', not end of file",
				},
			},
		},
		{
			input: "fn name(a int, ..., b int){}",
			diags: []diagnostics.Diag{
				{
					Message: "test.tt:1:16: ... is only allowed at the end of parameter list",
				},
			},
		},
		{
			input: "fn name() }",
			diags: []diagnostics.Diag{
				{
					Message: "test.tt:1:11: expected type or {, not }",
				},
			},
		},
		{
			input: "fn name()",
			diags: []diagnostics.Diag{
				{
					Message: "test.tt:1:10: expected type or {, not end of file",
				},
			},
		},
		{
			input: "fn name() {",
			diags: []diagnostics.Diag{
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
			diags: []diagnostics.Diag{
				{
					Message: "test.tt:2:16: expected ; at the end of prototype, not {",
				},
			},
		},
		{
			input: "extern {}",
			diags: []diagnostics.Diag{
				{
					Message: "test.tt:1:8: expected name, not {",
				},
			},
		},
		{
			input: "extern libc }",
			diags: []diagnostics.Diag{
				{
					Message: "test.tt:1:13: expected {, not }",
				},
			},
		},
		{
			input: "extern libc {",
			diags: []diagnostics.Diag{
				{
					Message: "test.tt:1:14: expected prototype or }, not end of file",
				},
			},
		},
		{
			input: "extern libc",
			diags: []diagnostics.Diag{
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
			diags: []diagnostics.Diag{
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
			diags: []diagnostics.Diag{
				{
					Message: "test.tt:2:7: expected name, not (",
				},
			},
		},
		{
			input: // no formatting
			`extern libc {
			fn name();`,
			diags: []diagnostics.Diag{
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
			diags: []diagnostics.Diag{
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
			diags: []diagnostics.Diag{
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
			diags: []diagnostics.Diag{
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
			diags: []diagnostics.Diag{
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
			diags: []diagnostics.Diag{
				{
					Message: "test.tt:2:17: expected parameter or ), not ;",
				},
			},
		},
	}

	for _, test := range tests {
		t.Run(fmt.Sprintf("TestSyntaxErrors('%s')", test.input), func(t *testing.T) {
			collector := diagnostics.New()

			src := []byte(test.input)

			lex := lexer.New(filename, src, collector)
			fmt.Println(lex)
			parser := New(collector)
			_, err := parser.ParseFileAsProgram(lex)

			if err != nil && len(parser.collector.Diags) == 0 {
				t.Fatalf(
					"error detected, but diagnostic collector is empty.\nError: %s",
					err,
				)
			}

			if len(test.diags) != len(parser.collector.Diags) {
				t.Fatalf(
					"expected to have %d diag(s), but got %d\n\ngot: %s\nexp: %s\n",
					len(test.diags),
					len(parser.collector.Diags),
					parser.collector.Diags,
					test.diags,
				)
			}
			if !reflect.DeepEqual(test.diags, parser.collector.Diags) {
				t.Fatalf("\nexpected diags: %v\ngot diags: %v\n", test.diags, parser.collector)
			}
		})
	}
}

func TestSyntaxErrorsOnBlock(t *testing.T) {
	filename := "test.tt"

	tests := []diagErrorTest{
		{
			input: "{",
			diags: []diagnostics.Diag{
				{
					Message: "test.tt:1:2: expected statement or }, not end of file",
				},
			},
		},
		{
			input: `{
			return
			}`,
			diags: []diagnostics.Diag{
				{
					Message: "test.tt:3:4: expected expression or ;, not }",
				},
			},
		},
		{
			input: `{
			return
			`,
			diags: []diagnostics.Diag{
				{
					Message: "test.tt:3:4: expected expression or ;, not end of file",
				},
			},
		},
		{
			input: `{
			return 10
			}`,
			diags: []diagnostics.Diag{
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
			collector := diagnostics.New()

			src := []byte(test.input)
			lex := lexer.New(filename, src, collector)
			parser := NewWithLex(lex, collector)

			_, err := parser.parseBlock()
			if err == nil {
				t.Fatal("expected to have syntax errors, but got nothing")
			}

			if len(test.diags) != len(parser.collector.Diags) {
				t.Fatalf(
					"expected to have %d diag(s), but got %d\n\ngot: %s\nexp: %s\n",
					len(test.diags),
					len(parser.collector.Diags),
					parser.collector.Diags,
					test.diags,
				)
			}
			if !reflect.DeepEqual(test.diags, parser.collector.Diags) {
				t.Fatalf("\nexpected diags: %v\ngot diags: %v\n", test.diags, parser.collector)
			}
		})
	}
}

func TestFunctionRedeclarationErrors(t *testing.T) {
	filename := "test.t"

	tests := []diagErrorTest{
		{
			input: "fn do_nothing() {}\nfn do_nothing() {}",
			diags: []diagnostics.Diag{
				{
					// TODO: show the first declaration position and the other
					Message: "test.t:2:4: function 'do_nothing' already declared on scope",
				},
			},
		},
	}

	for _, test := range tests {
		t.Run(fmt.Sprintf("TestFunctionRedeclarationErrors('%s')", test.input), func(t *testing.T) {
			collector := diagnostics.New()

			src := []byte(test.input)
			lex := lexer.New(filename, src, collector)
			parser := New(collector)

			_, err := parser.ParseFileAsProgram(lex)
			if err == nil {
				t.Fatal("expected to have syntax errors, but got nothing")
			}

			if len(test.diags) != len(parser.collector.Diags) {
				t.Fatalf(
					"expected to have %d diag(s), but got %d\n\ngot: %s\nexp: %s\n",
					len(test.diags),
					len(parser.collector.Diags),
					parser.collector.Diags,
					test.diags,
				)
			}
			if !reflect.DeepEqual(test.diags, parser.collector.Diags) {
				t.Fatalf("\nexpected diags: %v\ngot diags: %v\n", test.diags, parser.collector)
			}
		})
	}
}
