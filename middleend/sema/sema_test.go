package sema

import (
	"fmt"
	"reflect"
	"testing"

	"github.com/HicaroD/Telia/diagnostics"
	"github.com/HicaroD/Telia/frontend/ast"
	"github.com/HicaroD/Telia/frontend/lexer"
	"github.com/HicaroD/Telia/frontend/lexer/token"
	"github.com/HicaroD/Telia/frontend/parser"
	"github.com/HicaroD/Telia/scope"
)

type varTest struct {
	input    string
	ty       ast.ExprType
	inferred bool
}

func TestVarDeclForInference(t *testing.T) {
	filename := "test.tt"
	tests := []varTest{
		{
			input:    `name := "Hicaro";`,
			ty:       &ast.PointerType{Type: &ast.BasicType{Kind: token.U8_TYPE}},
			inferred: true,
		},
		{
			input:    "age := 18;",
			ty:       &ast.BasicType{Kind: token.INT_TYPE},
			inferred: true,
		},
		{
			input:    "score := -18;",
			ty:       &ast.BasicType{Kind: token.INT_TYPE},
			inferred: true,
		},
		{
			input:    "age := 1 + 1;",
			ty:       &ast.BasicType{Kind: token.INT_TYPE},
			inferred: true,
		},
		{
			input:    "age := 1 - 1;",
			ty:       &ast.BasicType{Kind: token.INT_TYPE},
			inferred: true,
		},
		{
			input:    "can_vote := true;",
			ty:       &ast.BasicType{Kind: token.BOOL_TYPE},
			inferred: true,
		},
		{
			input:    "can_vote := false;",
			ty:       &ast.BasicType{Kind: token.BOOL_TYPE},
			inferred: true,
		},
		{
			input:    "can_vote := false;",
			ty:       &ast.BasicType{Kind: token.BOOL_TYPE},
			inferred: true,
		},
		{
			input:    "is_greater := 2 > 1;",
			ty:       &ast.BasicType{Kind: token.BOOL_TYPE},
			inferred: true,
		},
		{
			input:    "is_greater_or_eq := 2 >= 1;",
			ty:       &ast.BasicType{Kind: token.BOOL_TYPE},
			inferred: true,
		},
		{
			input:    "is_lesser := 2 < 1;",
			ty:       &ast.BasicType{Kind: token.BOOL_TYPE},
			inferred: true,
		},
		{
			input:    "is_lesser_or_eq := 2 <= 1;",
			ty:       &ast.BasicType{Kind: token.BOOL_TYPE},
			inferred: true,
		},
		{
			input:    "is_eq := 2 == 1;",
			ty:       &ast.BasicType{Kind: token.BOOL_TYPE},
			inferred: true,
		},
		{
			input:    "is_true := true and true;",
			ty:       &ast.BasicType{Kind: token.BOOL_TYPE},
			inferred: true,
		},
		{
			input:    "is_true := true or true;",
			ty:       &ast.BasicType{Kind: token.BOOL_TYPE},
			inferred: true,
		},
		{
			input:    "is_not_true := not true;",
			ty:       &ast.BasicType{Kind: token.BOOL_TYPE},
			inferred: true,
		},
		// TODO: test variable decl with explicit type annotation
	}

	for _, test := range tests {
		t.Run(fmt.Sprintf("TestVarDeclForInference('%s')", test.input), func(t *testing.T) {
			varStmt, err := analyzeVarDeclFrom(test.input, filename)
			if err != nil {
				t.Fatal(err)
			}

			variable, ok := varStmt.(*ast.VarStmt)
			if !ok {
				t.Fatalf("unable to cast %s to *ast.VarStmt", reflect.TypeOf(varStmt))
			}

			if !(variable.NeedsInference && test.inferred) {
				t.Fatalf(
					"inference error, expected %v, but got %v",
					test.inferred,
					variable.NeedsInference,
				)
			}
			if !reflect.DeepEqual(variable.Type, test.ty) {
				t.Fatalf("type mismatch, expect %s, but got %s", test.ty, variable.Type)
			}
		})
	}
}

type exprInferenceTest struct {
	scope *scope.Scope[ast.Node]
	tests []struct {
		input string
		ty    ast.ExprType
		value ast.Expr
	}
}

func TestExprInferenceWithoutContext(t *testing.T) {
	filename := "test.tt"
	tests := []exprInferenceTest{
		{
			scope: &scope.Scope[ast.Node]{
				Parent: nil,
				Nodes: map[string]ast.Node{
					"a": &ast.VarStmt{
						Name: token.New([]byte("a"), token.ID, token.NewPosition(filename, 1, 1)),
						Type: &ast.BasicType{Kind: token.I8_TYPE},
						Value: ast.LiteralExpr{
							Type:  &ast.BasicType{Kind: token.I8_TYPE},
							Value: []byte("1"),
						},
						NeedsInference: false,
					},
				},
			},
			tests: []struct {
				input string
				ty    ast.ExprType
				value ast.Expr
			}{
				{
					input: "true",
					ty:    &ast.BasicType{Kind: token.BOOL_TYPE},
					value: &ast.LiteralExpr{
						Value: []byte("1"),
						Type:  &ast.BasicType{Kind: token.BOOL_TYPE},
					},
				},
				{
					input: "false",
					ty:    &ast.BasicType{Kind: token.BOOL_TYPE},
					value: &ast.LiteralExpr{
						Value: []byte("0"),
						Type:  &ast.BasicType{Kind: token.BOOL_TYPE},
					},
				},
				{
					input: "1",
					ty:    &ast.BasicType{Kind: token.INT_TYPE},
					value: &ast.LiteralExpr{
						Value: []byte("1"),
						Type:  &ast.BasicType{Kind: token.INT_TYPE},
					},
				},
				{
					input: "1 + 1",
					ty:    &ast.BasicType{Kind: token.INT_TYPE},
					value: &ast.BinaryExpr{
						Left: &ast.LiteralExpr{
							Value: []byte("1"),
							Type:  &ast.BasicType{Kind: token.INT_TYPE},
						},
						Op: token.PLUS,
						Right: &ast.LiteralExpr{
							Value: []byte("1"),
							Type:  &ast.BasicType{Kind: token.INT_TYPE},
						},
					},
				},
				{
					input: "-1",
					ty:    &ast.BasicType{Kind: token.INT_TYPE},
					value: &ast.UnaryExpr{
						Op: token.MINUS,
						Value: &ast.LiteralExpr{
							Type:  &ast.BasicType{Kind: token.INT_TYPE},
							Value: []byte("1"),
						},
					},
				},
				{
					input: "-1 + 1",
					ty:    &ast.BasicType{Kind: token.INT_TYPE},
					value: &ast.BinaryExpr{
						Left: &ast.UnaryExpr{
							Op: token.MINUS,
							Value: &ast.LiteralExpr{
								Type:  &ast.BasicType{Kind: token.INT_TYPE},
								Value: []byte("1"),
							},
						},
						Op: token.PLUS,
						Right: &ast.LiteralExpr{
							Value: []byte("1"),
							Type:  &ast.BasicType{Kind: token.INT_TYPE},
						},
					},
				},
				{
					input: "a + 1",
					ty:    &ast.BasicType{Kind: token.I8_TYPE},
					value: &ast.BinaryExpr{
						Left: &ast.IdExpr{
							Name: token.New([]byte("a"), token.ID, token.NewPosition(filename, 1, 1)),
						},
						Op: token.PLUS,
						Right: &ast.LiteralExpr{
							Value: []byte("1"),
							Type:  &ast.BasicType{Kind: token.I8_TYPE},
						},
					},
				},
				{
					input: "1 + a",
					ty:    &ast.BasicType{Kind: token.I8_TYPE},
					value: &ast.BinaryExpr{
						Left: &ast.LiteralExpr{
							Value: []byte("1"),
							Type:  &ast.BasicType{Kind: token.I8_TYPE},
						},
						Op: token.PLUS,
						Right: &ast.IdExpr{
							Name: token.New([]byte("a"), token.ID, token.NewPosition(filename, 5, 1)),
						},
					},
				},
				{
					input: "1 + 2 + a",
					ty:    &ast.BasicType{Kind: token.I8_TYPE},
					value: &ast.BinaryExpr{
						Left: &ast.BinaryExpr{
							Left: &ast.LiteralExpr{
								Value: []byte("1"),
								Type:  &ast.BasicType{Kind: token.I8_TYPE},
							},
							Op: token.PLUS,
							Right: &ast.LiteralExpr{
								Value: []byte("2"),
								Type:  &ast.BasicType{Kind: token.I8_TYPE},
							},
						},
						Op: token.PLUS,
						Right: &ast.IdExpr{
							Name: token.New([]byte("a"), token.ID, token.NewPosition(filename, 9, 1)),
						},
					},
				},
				{
					input: "1 + a + 3",
					ty:    &ast.BasicType{Kind: token.I8_TYPE},
					value: &ast.BinaryExpr{
						Left: &ast.BinaryExpr{
							Left: &ast.LiteralExpr{
								Value: []byte("1"),
								Type:  &ast.BasicType{Kind: token.I8_TYPE},
							},
							Op: token.PLUS,
							Right: &ast.IdExpr{
								Name: token.New([]byte("a"), token.ID, token.NewPosition(filename, 5, 1)),
							},
						},
						Op: token.PLUS,
						Right: &ast.LiteralExpr{
							Value: []byte("3"),
							Type:  &ast.BasicType{Kind: token.I8_TYPE},
						},
					},
				},
			},
		},
	}
	for _, test := range tests {
		for _, unit := range test.tests {
			t.Run(
				fmt.Sprintf("TestExprInferenceWithoutContext('%s')", unit.input),
				func(t *testing.T) {
					actualExpr, actualExprTy, err := inferExprTypeWithoutContext(
						unit.input,
						filename,
						test.scope,
					)
					if err != nil {
						t.Fatal(err)
					}
					if !reflect.DeepEqual(actualExpr, unit.value) {
						t.Fatalf("\nexpected expr: %s\ngot expr: %s\n", unit.value, actualExpr)
					}
					if !reflect.DeepEqual(actualExprTy, unit.ty) {
						t.Fatalf("\nexpected ty: %s\ngot ty: %s\n", unit.ty, actualExprTy)
					}
				},
			)
		}
	}
}

// TODO: test for mismatched types errors

func TestExprInferenceWithContext(t *testing.T) {
	filename := "test.tt"
	tests := []exprInferenceTest{
		{
			scope: &scope.Scope[ast.Node]{
				Parent: nil,
				Nodes: map[string]ast.Node{
					"a": &ast.VarStmt{
						Name: token.New([]byte("a"), token.ID, token.NewPosition(filename, 1, 1)),
						Type: &ast.BasicType{Kind: token.I8_TYPE},
						Value: ast.LiteralExpr{
							Type:  &ast.BasicType{Kind: token.I8_TYPE},
							Value: []byte("1"),
						},
						NeedsInference: false,
					},
				},
			},
			tests: []struct {
				input string
				ty    ast.ExprType
				value ast.Expr
			}{
				{
					input: "a + 1",
					ty:    &ast.BasicType{Kind: token.I8_TYPE},
					value: &ast.BinaryExpr{
						Left: &ast.IdExpr{
							Name: token.New([]byte("a"), token.ID, token.NewPosition(filename, 1, 1)),
						},
						Op: token.PLUS,
						Right: &ast.LiteralExpr{
							Value: []byte("1"),
							Type:  &ast.BasicType{Kind: token.INT_TYPE},
						},
					},
				},
				{
					input: "1 + a",
					ty:    &ast.BasicType{Kind: token.I8_TYPE},
					value: &ast.BinaryExpr{
						Left: &ast.LiteralExpr{
							Value: []byte("1"),
							Type:  &ast.BasicType{Kind: token.INT_TYPE},
						},
						Op: token.PLUS,
						Right: &ast.IdExpr{
							Name: token.New([]byte("a"), token.ID, token.NewPosition(filename, 1, 1)),
						},
					},
				},
			},
		},
	}
	for _, test := range tests {
		for _, unit := range test.tests {
			t.Run(
				fmt.Sprintf("TestExprInferenceWithContext('%s')", unit.input),
				func(t *testing.T) {
					actualExprTy, err := inferExprTypeWithContext(
						unit.input,
						filename,
						unit.ty,
						test.scope,
					)
					if err != nil {
						t.Fatal(err)
					}
					if !reflect.DeepEqual(actualExprTy, unit.ty) {
						t.Fatalf("\nexpected: %s\ngot: %s\n", unit.ty, actualExprTy)
					}
				},
			)
		}
	}
}

type semanticErrorTest struct {
	input string
	diags []diagnostics.Diag
}

func TestSemanticErrors(t *testing.T) {
	filename := "test.tt"

	tests := []semanticErrorTest{
		{
			input: "extern libc { fn puts(); fn puts(); }",
			diags: []diagnostics.Diag{
				{
					// TODO: show the first declaration and the other
					Message: "test.tt:1:29: prototype 'puts' already declared on extern 'libc'",
				},
			},
		},
		{
			input: "extern libc { }\nextern libc { }",
			diags: []diagnostics.Diag{
				{
					// TODO: show the first declaration and the other
					Message: "test.tt:2:8: extern 'libc' already declared on scope",
				},
			},
		},
		{
			input: "extern libc {}\nfn main() { libc.puts(); }",
			diags: []diagnostics.Diag{
				{
					// TODO: show the first declaration and the other
					Message: "test.tt:2:18: function 'puts' not declared on extern 'libc'",
				},
			},
		},
		{
			input: "fn do_nothing() {}\nfn do_nothing() {}",
			diags: []diagnostics.Diag{
				{
					// TODO: show the first declaration and the other
					Message: "test.tt:2:4: function 'do_nothing' already declared on scope",
				},
			},
		},
		{
			input: "fn do_nothing(a int, a int) {}",
			diags: []diagnostics.Diag{
				{
					Message: "test.tt:1:22: parameter 'a' already declared on function 'do_nothing'",
				},
			},
		},
		{
			input: "fn foo(a int, b int) {}\nfn main() { foo(); }",
			diags: []diagnostics.Diag{
				{
					Message: "test.tt:2:13: not enough arguments in call to 'foo'",
				},
			},
		},
		{
			input: "fn foo(a int) {}\nfn main() { foo(\"hello\"); }",
			diags: []diagnostics.Diag{
				{
					Message: "can't use *u8 on argument of type int",
				},
			},
		},
		{
			input: "fn main() { foo(); }",
			diags: []diagnostics.Diag{
				{
					Message: "test.tt:1:13: function 'foo' not defined on scope",
				},
			},
		},
		{
			input: "fn main() { a := 1; a(); }",
			diags: []diagnostics.Diag{
				{
					Message: "test.tt:1:21: 'a' is not callable",
				},
			},
		},
		{
			input: "fn main() { libc.printf(); }",
			diags: []diagnostics.Diag{
				{
					Message: "test.tt:1:13: 'libc' not defined on scope",
				},
			},
		},
		// Multiple variables
		{
			input: "fn main() { a, b := 10, 10; }",
			diags: nil, // no errors
		},
		{
			input: "fn main() { a := 1; b := 2; a, b = 10, 10; }",
			diags: nil, // no errors
		},
		{
			input: "fn main() { a := 1; a, b = 10, 10; }",
			diags: []diagnostics.Diag{
				{
					Message: "test.tt:1:24: 'b' not declared",
				},
			},
		},
		{
			input: "fn main() { b := 1; a, b = 10, 10; }",
			diags: []diagnostics.Diag{
				{
					Message: "test.tt:1:21: 'a' not declared",
				},
			},
		},
		{
			input: "fn main() { a := 1; b := 2; a, b := 10, 10; }",
			diags: []diagnostics.Diag{
				{
					Message: "test.tt:1:29: no new variables declared",
				},
			},
		},
	}

	for _, test := range tests {
		t.Run(fmt.Sprintf("TestSemanticErrors('%s')", test.input), func(t *testing.T) {
			collector := diagnostics.New()

			src := []byte(test.input)
			lex := lexer.New(filename, src, collector)
			parser := parser.New(lex, collector)

			program, err := parser.Parse()
			if err != nil {
				t.Fatal(err)
			}

			sema := New(collector)
			_ = sema.Check(program)

			if len(collector.Diags) != len(test.diags) {
				t.Fatalf(
					"expected to have %d diag(s), but got %d\n\ngot: %s\nexp: %s\n",
					len(test.diags),
					len(sema.collector.Diags),
					sema.collector.Diags,
					test.diags,
				)
			}

			if !reflect.DeepEqual(collector.Diags, test.diags) {
				t.Fatalf("\nexp: %v\ngot: %v\n", test.diags, collector.Diags)
			}
		})
	}
}
