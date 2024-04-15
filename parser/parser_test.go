package parser

import (
	"fmt"
	"reflect"
	"testing"

	"github.com/HicaroD/Telia/ast"
	"github.com/HicaroD/Telia/lexer/token"
	"github.com/HicaroD/Telia/lexer/token/kind"
)

type functionDeclTest struct {
	input string
	node  *ast.FunctionDecl
}

func TestFunctionDecl(t *testing.T) {
	filename := "test.tt"

	// parent := scope.New[ast.AstNode](nil)
	// scope := scope.New(parent)

	tests := []functionDeclTest{
		{
			input: "fn do_nothing() {}",
			node: &ast.FunctionDecl{
				Scope: nil,
				Name:  "do_nothing",
				Params: &ast.FieldList{
					Open:       token.New(nil, kind.OPEN_PAREN, token.NewPosition(filename, 14, 1)),
					Fields:     nil,
					Close:      token.New(nil, kind.CLOSE_PAREN, token.NewPosition(filename, 15, 1)),
					IsVariadic: false,
				},
				RetType: ast.BasicType{Kind: kind.VOID_TYPE},
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
				Name:  "do_nothing",
				Params: &ast.FieldList{
					Open: token.New(nil, kind.OPEN_PAREN, token.NewPosition(filename, 14, 1)),
					Fields: []*ast.Field{
						{
							Name: token.New("a", kind.ID, token.NewPosition(filename, 15, 1)),
							Type: ast.BasicType{Kind: kind.BOOL_TYPE},
						},
					},
					Close:      token.New(nil, kind.CLOSE_PAREN, token.NewPosition(filename, 21, 1)),
					IsVariadic: false,
				},
				RetType: ast.BasicType{Kind: kind.VOID_TYPE},
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
				Name:  "do_nothing",
				Params: &ast.FieldList{
					Open: token.New(nil, kind.OPEN_PAREN, token.NewPosition(filename, 14, 1)),
					Fields: []*ast.Field{
						{
							Name: token.New("a", kind.ID, token.NewPosition(filename, 15, 1)),
							Type: ast.BasicType{Kind: kind.BOOL_TYPE},
						},
						{
							Name: token.New("b", kind.ID, token.NewPosition(filename, 23, 1)),
							Type: ast.BasicType{Kind: kind.I32_TYPE},
						},
					},
					Close:      token.New(nil, kind.CLOSE_PAREN, token.NewPosition(filename, 28, 1)),
					IsVariadic: false,
				},
				RetType: ast.BasicType{Kind: kind.VOID_TYPE},
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
				Name:  "do_nothing",
				Params: &ast.FieldList{
					Open: token.New(nil, kind.OPEN_PAREN, token.NewPosition(filename, 14, 1)),
					Fields: []*ast.Field{
						{
							Name: token.New("a", kind.ID, token.NewPosition(filename, 15, 1)),
							Type: ast.BasicType{Kind: kind.BOOL_TYPE},
						},
						{
							Name: token.New("b", kind.ID, token.NewPosition(filename, 23, 1)),
							Type: ast.BasicType{Kind: kind.I32_TYPE},
						},
					},
					Close:      token.New(nil, kind.CLOSE_PAREN, token.NewPosition(filename, 28, 1)),
					IsVariadic: false,
				},
				RetType: ast.BasicType{Kind: kind.I8_TYPE},
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
				Name:  "do_nothing",
				Params: &ast.FieldList{
					Open: token.New(nil, kind.OPEN_PAREN, token.NewPosition(filename, 14, 1)),
					Fields: []*ast.Field{
						{
							Name: token.New("a", kind.ID, token.NewPosition(filename, 15, 1)),
							Type: ast.BasicType{Kind: kind.BOOL_TYPE},
						},
						{
							Name: token.New("b", kind.ID, token.NewPosition(filename, 23, 1)),
							Type: ast.BasicType{Kind: kind.I32_TYPE},
						},
					},
					Close:      token.New(nil, kind.CLOSE_PAREN, token.NewPosition(filename, 28, 1)),
					IsVariadic: false,
				},
				RetType: ast.BasicType{Kind: kind.U8_TYPE},
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
				Name:  "do_nothing",
				Params: &ast.FieldList{
					Open: token.New(nil, kind.OPEN_PAREN, token.NewPosition(filename, 14, 1)),
					Fields: []*ast.Field{
						{
							Name: token.New("a", kind.ID, token.NewPosition(filename, 15, 1)),
							Type: ast.BasicType{Kind: kind.BOOL_TYPE},
						},
						{
							Name: token.New("b", kind.ID, token.NewPosition(filename, 23, 1)),
							Type: ast.BasicType{Kind: kind.I32_TYPE},
						},
					},
					Close:      token.New(nil, kind.CLOSE_PAREN, token.NewPosition(filename, 28, 1)),
					IsVariadic: false,
				},
				RetType: ast.BasicType{Kind: kind.I16_TYPE},
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
				Name:  "do_nothing",
				Params: &ast.FieldList{
					Open: token.New(nil, kind.OPEN_PAREN, token.NewPosition(filename, 14, 1)),
					Fields: []*ast.Field{
						{
							Name: token.New("a", kind.ID, token.NewPosition(filename, 15, 1)),
							Type: ast.BasicType{Kind: kind.BOOL_TYPE},
						},
						{
							Name: token.New("b", kind.ID, token.NewPosition(filename, 23, 1)),
							Type: ast.BasicType{Kind: kind.I32_TYPE},
						},
					},
					Close:      token.New(nil, kind.CLOSE_PAREN, token.NewPosition(filename, 28, 1)),
					IsVariadic: false,
				},
				RetType: ast.BasicType{Kind: kind.U16_TYPE},
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
				Name:  "do_nothing",
				Params: &ast.FieldList{
					Open: token.New(nil, kind.OPEN_PAREN, token.NewPosition(filename, 14, 1)),
					Fields: []*ast.Field{
						{
							Name: token.New("a", kind.ID, token.NewPosition(filename, 15, 1)),
							Type: ast.BasicType{Kind: kind.BOOL_TYPE},
						},
						{
							Name: token.New("b", kind.ID, token.NewPosition(filename, 23, 1)),
							Type: ast.BasicType{Kind: kind.I32_TYPE},
						},
					},
					Close:      token.New(nil, kind.CLOSE_PAREN, token.NewPosition(filename, 28, 1)),
					IsVariadic: false,
				},
				RetType: ast.BasicType{Kind: kind.I32_TYPE},
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
				Name:  "do_nothing",
				Params: &ast.FieldList{
					Open: token.New(nil, kind.OPEN_PAREN, token.NewPosition(filename, 14, 1)),
					Fields: []*ast.Field{
						{
							Name: token.New("a", kind.ID, token.NewPosition(filename, 15, 1)),
							Type: ast.BasicType{Kind: kind.BOOL_TYPE},
						},
						{
							Name: token.New("b", kind.ID, token.NewPosition(filename, 23, 1)),
							Type: ast.BasicType{Kind: kind.I32_TYPE},
						},
					},
					Close:      token.New(nil, kind.CLOSE_PAREN, token.NewPosition(filename, 28, 1)),
					IsVariadic: false,
				},
				RetType: ast.BasicType{Kind: kind.U32_TYPE},
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
				Name:  "do_nothing",
				Params: &ast.FieldList{
					Open: token.New(nil, kind.OPEN_PAREN, token.NewPosition(filename, 14, 1)),
					Fields: []*ast.Field{
						{
							Name: token.New("a", kind.ID, token.NewPosition(filename, 15, 1)),
							Type: ast.BasicType{Kind: kind.BOOL_TYPE},
						},
						{
							Name: token.New("b", kind.ID, token.NewPosition(filename, 23, 1)),
							Type: ast.BasicType{Kind: kind.I32_TYPE},
						},
					},
					Close:      token.New(nil, kind.CLOSE_PAREN, token.NewPosition(filename, 28, 1)),
					IsVariadic: false,
				},
				RetType: ast.BasicType{Kind: kind.I64_TYPE},
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
				Name:  "do_nothing",
				Params: &ast.FieldList{
					Open: token.New(nil, kind.OPEN_PAREN, token.NewPosition(filename, 14, 1)),
					Fields: []*ast.Field{
						{
							Name: token.New("a", kind.ID, token.NewPosition(filename, 15, 1)),
							Type: ast.BasicType{Kind: kind.BOOL_TYPE},
						},
						{
							Name: token.New("b", kind.ID, token.NewPosition(filename, 23, 1)),
							Type: ast.BasicType{Kind: kind.I32_TYPE},
						},
					},
					Close:      token.New(nil, kind.CLOSE_PAREN, token.NewPosition(filename, 28, 1)),
					IsVariadic: false,
				},
				RetType: ast.BasicType{Kind: kind.U64_TYPE},
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
				Name:  "do_nothing",
				Params: &ast.FieldList{
					Open: token.New(nil, kind.OPEN_PAREN, token.NewPosition(filename, 14, 1)),
					Fields: []*ast.Field{
						{
							Name: token.New("a", kind.ID, token.NewPosition(filename, 15, 1)),
							Type: ast.BasicType{Kind: kind.BOOL_TYPE},
						},
						{
							Name: token.New("b", kind.ID, token.NewPosition(filename, 23, 1)),
							Type: ast.BasicType{Kind: kind.I32_TYPE},
						},
					},
					Close:      token.New(nil, kind.CLOSE_PAREN, token.NewPosition(filename, 28, 1)),
					IsVariadic: false,
				},
				RetType: ast.BasicType{Kind: kind.BOOL_TYPE},
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
				Name:  "do_nothing",
				Params: &ast.FieldList{
					Open: token.New(nil, kind.OPEN_PAREN, token.NewPosition(filename, 14, 1)),
					Fields: []*ast.Field{
						{
							Name: token.New("a", kind.ID, token.NewPosition(filename, 15, 1)),
							Type: ast.BasicType{Kind: kind.BOOL_TYPE},
						},
						{
							Name: token.New("b", kind.ID, token.NewPosition(filename, 23, 1)),
							Type: ast.BasicType{Kind: kind.I32_TYPE},
						},
					},
					Close:      token.New(nil, kind.CLOSE_PAREN, token.NewPosition(filename, 28, 1)),
					IsVariadic: false,
				},
				RetType: ast.PointerType{Type: ast.BasicType{Kind: kind.I8_TYPE}},
				Block: &ast.BlockStmt{
					OpenCurly:  token.NewPosition(filename, 34, 1),
					Statements: nil,
					CloseCurly: token.NewPosition(filename, 35, 1),
				},
			},
		},
		// TODO: test variadic arguments on functions
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

func TestExternDecl(t *testing.T) {}

type exprTest struct {
	input string
	node  ast.Expr
}

func TestLiteralExpr(t *testing.T) {
	filename := "test.tt"
	tests := []exprTest{
		{
			input: "1",
			node:  &ast.LiteralExpr{Value: 1, Kind: kind.INTEGER_LITERAL},
		},
		{
			input: "true",
			node:  &ast.LiteralExpr{Value: "true", Kind: kind.TRUE_BOOL_LITERAL},
		},
		{
			input: "false",
			node:  &ast.LiteralExpr{Value: "false", Kind: kind.FALSE_BOOL_LITERAL},
		},
		{
			input: "\"Hello, world\"",
			node:  &ast.LiteralExpr{Value: "Hello, world", Kind: kind.STRING_LITERAL},
		},
	}

	for _, test := range tests {
		t.Run(fmt.Sprintf("TestLiteralExpr('%s')", test.input), func(t *testing.T) {
			actualNode, err := ParseExprFrom(test.input, filename)
			if err != nil {
				t.Errorf("TestLiteralExpr('%s'): unexpected error '%v'", test.input, err)
			}
			if !reflect.DeepEqual(test.node, actualNode) {
				t.Errorf("TestLiteralExpr('%s'): expression node differs\nexpected: '%v', but got '%v'\n", test.input, test.node, actualNode)
			}
		})
	}
}

func TestUnaryExpr(t *testing.T) {
	filename := "test.tt"
	unaryExprs := []exprTest{
		{
			input: "-1",
			node: &ast.UnaryExpr{
				Op:   kind.MINUS,
				Node: &ast.LiteralExpr{Value: 1, Kind: kind.INTEGER_LITERAL},
			},
		},
		{
			input: "not true",
			node: &ast.UnaryExpr{
				Op:   kind.NOT,
				Node: &ast.LiteralExpr{Value: "true", Kind: kind.TRUE_BOOL_LITERAL},
			},
		},
	}
	for _, test := range unaryExprs {
		t.Run(fmt.Sprintf("TestBinaryExpr('%s')", test.input), func(t *testing.T) {
			actualNode, err := ParseExprFrom(test.input, filename)
			if err != nil {
				t.Errorf("TestBinaryExpr('%s'): unexpected error '%v'", test.input, err)
			}
			if !reflect.DeepEqual(test.node, actualNode) {
				t.Errorf("TestBinaryExpr('%s'): expression node differs\nexpected: '%v', but got '%v'\n", test.input, test.node, actualNode)
			}
		})
	}
}

func TestBinaryExpr(t *testing.T) {
	filename := "test.tt"
	binExprs := []exprTest{
		{
			input: "1 + 1",
			node: &ast.BinaryExpr{
				Left:  &ast.LiteralExpr{Value: 1, Kind: kind.INTEGER_LITERAL},
				Op:    kind.PLUS,
				Right: &ast.LiteralExpr{Value: 1, Kind: kind.INTEGER_LITERAL},
			},
		},
		{
			input: "2 - 1",
			node: &ast.BinaryExpr{
				Left:  &ast.LiteralExpr{Value: 2, Kind: kind.INTEGER_LITERAL},
				Op:    kind.MINUS,
				Right: &ast.LiteralExpr{Value: 1, Kind: kind.INTEGER_LITERAL},
			},
		},
		{
			input: "5 * 10",
			node: &ast.BinaryExpr{
				Left:  &ast.LiteralExpr{Value: 5, Kind: kind.INTEGER_LITERAL},
				Op:    kind.STAR,
				Right: &ast.LiteralExpr{Value: 10, Kind: kind.INTEGER_LITERAL},
			},
		},
		{
			input: "3 + 4 * 5",
			node: &ast.BinaryExpr{
				Left: &ast.LiteralExpr{Value: 3, Kind: kind.INTEGER_LITERAL},
				Op:   kind.PLUS,
				Right: &ast.BinaryExpr{
					Left:  &ast.LiteralExpr{Value: 4, Kind: kind.INTEGER_LITERAL},
					Op:    kind.STAR,
					Right: &ast.LiteralExpr{Value: 5, Kind: kind.INTEGER_LITERAL},
				},
			},
		},
		{
			input: "3 + (4 * 5)",
			node: &ast.BinaryExpr{
				Left: &ast.LiteralExpr{Value: 3, Kind: kind.INTEGER_LITERAL},
				Op:   kind.PLUS,
				Right: &ast.BinaryExpr{
					Left:  &ast.LiteralExpr{Value: 4, Kind: kind.INTEGER_LITERAL},
					Op:    kind.STAR,
					Right: &ast.LiteralExpr{Value: 5, Kind: kind.INTEGER_LITERAL},
				},
			},
		},
		{
			input: "10 / 1",
			node: &ast.BinaryExpr{
				Left:  &ast.LiteralExpr{Value: 10, Kind: kind.INTEGER_LITERAL},
				Op:    kind.SLASH,
				Right: &ast.LiteralExpr{Value: 1, Kind: kind.INTEGER_LITERAL},
			},
		},
		{
			input: "6 / 3 - 1",
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
			input: "6 / (3 - 1)",
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
			input: "1 / (1 + 1)",
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
			input: "1 > 1",
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
			input: "1 >= 1",
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
			input: "1 < 1",
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
			input: "1 <= 1",
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
			input: "not (1 > 1)",
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
			input: "1 > 1 and 1 > 1",
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
			input: "1 > 1 or 1 > 1",
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
			input: "celsius*9/5+32",
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
		{
			input: "get_celsius()*9/5+32",
			node: &ast.BinaryExpr{
				Left: &ast.BinaryExpr{
					Left: &ast.BinaryExpr{
						Left: &ast.FunctionCall{
							Name: "get_celsius",
							Args: nil,
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
		{
			input: "1 > 1 > 1",
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
			input: "n == 1",
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
			input: "n == 1 or n == 2",
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
			input: "1 + 1 > 2 and 1 == 1",
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
		{
			input: "true and true and true and true",
			node: &ast.BinaryExpr{
				Left: &ast.BinaryExpr{
					Left: &ast.BinaryExpr{
						Left: &ast.LiteralExpr{
							Value: "true",
							Kind:  kind.TRUE_BOOL_LITERAL,
						},
						Op: kind.AND,
						Right: &ast.LiteralExpr{
							Value: "true",
							Kind:  kind.TRUE_BOOL_LITERAL,
						},
					},
					Op: kind.AND,
					Right: &ast.LiteralExpr{
						Value: "true",
						Kind:  kind.TRUE_BOOL_LITERAL,
					},
				},
				Op: kind.AND,
				Right: &ast.LiteralExpr{
					Value: "true",
					Kind:  kind.TRUE_BOOL_LITERAL,
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
							Kind:  kind.TRUE_BOOL_LITERAL,
						},
						Op: kind.AND,
						Right: &ast.LiteralExpr{
							Value: "true",
							Kind:  kind.TRUE_BOOL_LITERAL,
						},
					},
					Op: kind.AND,
					Right: &ast.LiteralExpr{
						Value: "true",
						Kind:  kind.TRUE_BOOL_LITERAL,
					},
				},
				Op: kind.AND,
				Right: &ast.LiteralExpr{
					Value: "true",
					Kind:  kind.TRUE_BOOL_LITERAL,
				},
			},
		},
		{
			input: "1 + multiply_by_2(10)",
			node: &ast.BinaryExpr{
				Left: &ast.LiteralExpr{
					Value: 1,
					Kind:  kind.INTEGER_LITERAL,
				},
				Op: kind.PLUS,
				Right: &ast.FunctionCall{
					Name: "multiply_by_2",
					Args: []ast.Expr{
						&ast.LiteralExpr{
							Value: 10,
							Kind:  kind.INTEGER_LITERAL,
						},
					},
				},
			},
		},
	}

	for _, test := range binExprs {
		t.Run(fmt.Sprintf("TestBinaryExpr('%s')", test.input), func(t *testing.T) {
			actualNode, err := ParseExprFrom(test.input, filename)
			if err != nil {
				t.Errorf("unexpected error '%v'", err)
			}
			if !reflect.DeepEqual(test.node, actualNode) {
				t.Errorf("expression node differs\nexpected: '%v' '%v'\ngot:      '%v' '%v'\n", test.node, reflect.TypeOf(test.node), actualNode, reflect.TypeOf(actualNode))
			}
		})
	}
}

func TestFuncCallStmt(t *testing.T) {}

func TestIfStmt(t *testing.T) {}
