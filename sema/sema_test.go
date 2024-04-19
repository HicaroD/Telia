package sema

import (
	"fmt"
	"reflect"
	"testing"

	"github.com/HicaroD/Telia/ast"
	"github.com/HicaroD/Telia/lexer/token"
	"github.com/HicaroD/Telia/lexer/token/kind"
	"github.com/HicaroD/Telia/scope"
)

type varTest struct {
	input    string
	ty       ast.ExprType
	inferred bool
}

func TestVarDecl(t *testing.T) {
	filename := "test.tt"
	tests := []varTest{
		{
			input:    `name := "Hicaro";`,
			ty:       &ast.PointerType{Type: &ast.BasicType{Kind: kind.I8_TYPE}},
			inferred: true,
		},
		{
			input:    "age := 18;",
			ty:       &ast.BasicType{Kind: kind.INT_TYPE},
			inferred: true,
		},
		{
			input:    "score := -18;",
			ty:       &ast.BasicType{Kind: kind.INT_TYPE},
			inferred: true,
		},
		{
			input:    "age := 1 + 1;",
			ty:       &ast.BasicType{Kind: kind.INT_TYPE},
			inferred: true,
		},
		{
			input:    "age := 1 - 1;",
			ty:       &ast.BasicType{Kind: kind.INT_TYPE},
			inferred: true,
		},
		{
			input:    "can_vote := true;",
			ty:       &ast.BasicType{Kind: kind.BOOL_TYPE},
			inferred: true,
		},
		{
			input:    "can_vote := false;",
			ty:       &ast.BasicType{Kind: kind.BOOL_TYPE},
			inferred: true,
		},
		{
			input:    "can_vote := false;",
			ty:       &ast.BasicType{Kind: kind.BOOL_TYPE},
			inferred: true,
		},
		{
			input:    "is_greater := 2 > 1;",
			ty:       &ast.BasicType{Kind: kind.BOOL_TYPE},
			inferred: true,
		},
		{
			input:    "is_greater_or_eq := 2 >= 1;",
			ty:       &ast.BasicType{Kind: kind.BOOL_TYPE},
			inferred: true,
		},
		{
			input:    "is_lesser := 2 < 1;",
			ty:       &ast.BasicType{Kind: kind.BOOL_TYPE},
			inferred: true,
		},
		{
			input:    "is_lesser_or_eq := 2 <= 1;",
			ty:       &ast.BasicType{Kind: kind.BOOL_TYPE},
			inferred: true,
		},
		{
			input:    "is_eq := 2 == 1;",
			ty:       &ast.BasicType{Kind: kind.BOOL_TYPE},
			inferred: true,
		},
		{
			input:    "is_true := true and true;",
			ty:       &ast.BasicType{Kind: kind.BOOL_TYPE},
			inferred: true,
		},
		{
			input:    "is_true := true or true;",
			ty:       &ast.BasicType{Kind: kind.BOOL_TYPE},
			inferred: true,
		},
		{
			input:    "is_not_true := not true;",
			ty:       &ast.BasicType{Kind: kind.BOOL_TYPE},
			inferred: true,
		},
	}

	for _, test := range tests {
		t.Run(fmt.Sprintf("TestVarDecl('%s')", test.input), func(t *testing.T) {
			varDecl, err := analyzeVarDeclFrom(test.input, filename)
			if err != nil {
				t.Fatal(err)
			}
			if !(varDecl.NeedsInference && test.inferred) {
				t.Fatalf("inference error, expected %v, but got %v", test.inferred, varDecl.NeedsInference)
			}
			if !reflect.DeepEqual(varDecl.Type, test.ty) {
				t.Fatalf("type mismatch, expect %s, but got %s", test.ty, varDecl.Type)
			}
		})
	}
}

type exprInferenceTest struct {
	scope *scope.Scope[ast.AstNode]
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
			scope: &scope.Scope[ast.AstNode]{
				Parent: nil,
				Nodes: map[string]ast.AstNode{
					"a": &ast.VarDeclStmt{
						Name: token.New("a", kind.ID, token.NewPosition(filename, 1, 1)),
						Type: &ast.BasicType{Kind: kind.I8_TYPE},
						Value: ast.LiteralExpr{
							Type:  &ast.BasicType{Kind: kind.I8_TYPE},
							Value: 1,
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
					ty:    &ast.BasicType{Kind: kind.BOOL_TYPE},
					value: &ast.LiteralExpr{
						Value: "true",
						Type:  &ast.BasicType{Kind: kind.BOOL_TYPE},
					},
				},
				{
					input: "false",
					ty:    &ast.BasicType{Kind: kind.BOOL_TYPE},
					value: &ast.LiteralExpr{
						Value: "false",
						Type:  &ast.BasicType{Kind: kind.BOOL_TYPE},
					},
				},
				{
					input: "1",
					ty:    &ast.BasicType{Kind: kind.INT_TYPE},
					value: &ast.LiteralExpr{
						Value: "1",
						Type:  &ast.BasicType{Kind: kind.INT_TYPE},
					},
				},
				{
					input: "1 + 1",
					ty:    &ast.BasicType{Kind: kind.INT_TYPE},
					value: &ast.BinaryExpr{
						Left: &ast.LiteralExpr{
							Value: "1",
							Type:  &ast.BasicType{Kind: kind.INT_TYPE},
						},
						Op: kind.PLUS,
						Right: &ast.LiteralExpr{
							Value: "1",
							Type:  &ast.BasicType{Kind: kind.INT_TYPE},
						},
					},
				},
				{
					input: "-1",
					ty:    &ast.BasicType{Kind: kind.INT_TYPE},
					value: &ast.UnaryExpr{
						Op: kind.MINUS,
						Value: &ast.LiteralExpr{
							Type:  &ast.BasicType{Kind: kind.INT_TYPE},
							Value: "1",
						},
					},
				},
				{
					input: "-1 + 1",
					ty:    &ast.BasicType{Kind: kind.INT_TYPE},
					value: &ast.BinaryExpr{
						Left: &ast.UnaryExpr{
							Op: kind.MINUS,
							Value: &ast.LiteralExpr{
								Type:  &ast.BasicType{Kind: kind.INT_TYPE},
								Value: "1",
							},
						},
						Op: kind.PLUS,
						Right: &ast.LiteralExpr{
							Value: "1",
							Type:  &ast.BasicType{Kind: kind.INT_TYPE},
						},
					},
				},
				{
					input: "a + 1",
					ty:    &ast.BasicType{Kind: kind.I8_TYPE},
					value: &ast.BinaryExpr{
						Left: &ast.IdExpr{
							Name: token.New("a", kind.ID, token.NewPosition(filename, 1, 1)),
						},
						Op: kind.PLUS,
						Right: &ast.LiteralExpr{
							Value: "1",
							Type:  &ast.BasicType{Kind: kind.I8_TYPE},
						},
					},
				},
				{
					input: "1 + a",
					ty:    &ast.BasicType{Kind: kind.I8_TYPE},
					value: &ast.BinaryExpr{
						Left: &ast.LiteralExpr{
							Value: "1",
							Type:  &ast.BasicType{Kind: kind.I8_TYPE},
						},
						Op: kind.PLUS,
						Right: &ast.IdExpr{
							Name: token.New("a", kind.ID, token.NewPosition(filename, 5, 1)),
						},
					},
				},
				{
					input: "1 + 2 + a",
					ty:    &ast.BasicType{Kind: kind.I8_TYPE},
					value: &ast.BinaryExpr{
						Left: &ast.BinaryExpr{
							Left: &ast.LiteralExpr{
								Value: "1",
								Type:  &ast.BasicType{Kind: kind.I8_TYPE},
							},
							Op: kind.PLUS,
							Right: &ast.LiteralExpr{
								Value: "2",
								Type:  &ast.BasicType{Kind: kind.I8_TYPE},
							},
						},
						Op: kind.PLUS,
						Right: &ast.IdExpr{
							Name: token.New("a", kind.ID, token.NewPosition(filename, 9, 1)),
						},
					},
				},
				{
					input: "1 + a + 3",
					ty:    &ast.BasicType{Kind: kind.I8_TYPE},
					value: &ast.BinaryExpr{
						Left: &ast.BinaryExpr{
							Left: &ast.LiteralExpr{
								Value: "1",
								Type:  &ast.BasicType{Kind: kind.I8_TYPE},
							},
							Op: kind.PLUS,
							Right: &ast.IdExpr{
								Name: token.New("a", kind.ID, token.NewPosition(filename, 5, 1)),
							},
						},
						Op: kind.PLUS,
						Right: &ast.LiteralExpr{
							Value: "3",
							Type:  &ast.BasicType{Kind: kind.I8_TYPE},
						},
					},
				},
			},
		},
	}
	for _, test := range tests {
		for _, unit := range test.tests {
			t.Run(fmt.Sprintf("TestExprInferenceWithoutContext('%s')", unit.input), func(t *testing.T) {
				actualExpr, actualExprTy, err := inferExprTypeWithoutContext(unit.input, filename, test.scope)
				if err != nil {
					t.Fatal(err)
				}
				if !reflect.DeepEqual(actualExpr, unit.value) {
					t.Fatalf("\nexpected expr: %s\ngot expr: %s\n", unit.value, actualExpr)
				}
				if !reflect.DeepEqual(actualExprTy, unit.ty) {
					t.Fatalf("\nexpected ty: %s\ngot ty: %s\n", unit.ty, actualExprTy)
				}
			})
		}
	}
}

// TODO: test for mismatched types errors

func TestExprInferenceWithContext(t *testing.T) {
	filename := "test.tt"
	tests := []exprInferenceTest{
		{
			scope: &scope.Scope[ast.AstNode]{
				Parent: nil,
				Nodes: map[string]ast.AstNode{
					"a": &ast.VarDeclStmt{
						Name: token.New("a", kind.ID, token.NewPosition(filename, 1, 1)),
						Type: &ast.BasicType{Kind: kind.I8_TYPE},
						Value: ast.LiteralExpr{
							Type:  &ast.BasicType{Kind: kind.I8_TYPE},
							Value: 1,
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
					ty:    &ast.BasicType{Kind: kind.I8_TYPE},
					value: &ast.BinaryExpr{
						Left: &ast.IdExpr{
							Name: token.New("a", kind.ID, token.NewPosition(filename, 1, 1)),
						},
						Op: kind.PLUS,
						Right: &ast.LiteralExpr{
							Value: 1,
							Type:  &ast.BasicType{Kind: kind.INT_TYPE},
						},
					},
				},
				{
					input: "1 + a",
					ty:    &ast.BasicType{Kind: kind.I8_TYPE},
					value: &ast.BinaryExpr{
						Left: &ast.LiteralExpr{
							Value: 1,
							Type:  &ast.BasicType{Kind: kind.INT_TYPE},
						},
						Op: kind.PLUS,
						Right: &ast.IdExpr{
							Name: token.New("a", kind.ID, token.NewPosition(filename, 1, 1)),
						},
					},
				},
			},
		},
	}
	for _, test := range tests {
		for _, unit := range test.tests {
			t.Run(fmt.Sprintf("TestExprInferenceWithContext('%s')", unit.input), func(t *testing.T) {
				actualExprTy, err := inferExprTypeWithContext(unit.input, filename, unit.ty, test.scope)
				if err != nil {
					t.Fatal(err)
				}
				if !reflect.DeepEqual(actualExprTy, unit.ty) {
					t.Fatalf("\nexpected: %s\ngot: %s\n", unit.ty, actualExprTy)
				}
			})
		}
	}
}
