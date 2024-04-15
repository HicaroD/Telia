package sema

import (
	"fmt"
	"reflect"
	"testing"

	"github.com/HicaroD/Telia/ast"
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
			ty:       ast.PointerType{Type: ast.BasicType{Kind: kind.I8_TYPE}},
			inferred: true,
		},
		{
			input:    "age := 18;",
			ty:       ast.BasicType{Kind: kind.I32_TYPE},
			inferred: true,
		},
		// TODO: analyze more types of operators
		// {
		// 	input:    "age := 1 + 1;",
		// 	ty:       ast.BasicType{Kind: kind.I32_TYPE},
		// 	inferred: true,
		// },
		// TODO: test more types of integer inference before refactoring
		{
			input:    "can_vote := true;",
			ty:       ast.BasicType{Kind: kind.BOOL_TYPE},
			inferred: true,
		},
		{
			input:    "can_vote := false;",
			ty:       ast.BasicType{Kind: kind.BOOL_TYPE},
			inferred: true,
		},
		{
			input:    "can_vote := false;",
			ty:       ast.BasicType{Kind: kind.BOOL_TYPE},
			inferred: true,
		},
		{
			input:    "is_greater := 2 > 1;",
			ty:       ast.BasicType{Kind: kind.BOOL_TYPE},
			inferred: true,
		},
		{
			input:    "is_greater_or_eq := 2 >= 1;",
			ty:       ast.BasicType{Kind: kind.BOOL_TYPE},
			inferred: true,
		},
		{
			input:    "is_lesser := 2 < 1;",
			ty:       ast.BasicType{Kind: kind.BOOL_TYPE},
			inferred: true,
		},
		{
			input:    "is_lesser_or_eq := 2 <= 1;",
			ty:       ast.BasicType{Kind: kind.BOOL_TYPE},
			inferred: true,
		},
		{
			input:    "is_eq := 2 == 1;",
			ty:       ast.BasicType{Kind: kind.BOOL_TYPE},
			inferred: true,
		},
		{
			input:    "is_true := true and true;",
			ty:       ast.BasicType{Kind: kind.BOOL_TYPE},
			inferred: true,
		},
		{
			input:    "is_true := true or true;",
			ty:       ast.BasicType{Kind: kind.BOOL_TYPE},
			inferred: true,
		},
		// TODO: analyze unary expressions
		// {
		// 	input:    "is_not_true := not true;",
		// 	ty:       ast.BasicType{Kind: kind.BOOL_TYPE},
		// 	inferred: true,
		// },
	}

	for _, test := range tests {
		t.Run(fmt.Sprintf("TestBinaryExpr('%s')", test.input), func(t *testing.T) {
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
	}
}

func TestExprInference(t *testing.T) {
	filename := "test.tt"
	tests := []exprInferenceTest{
		{
			scope: nil,
			tests: []struct {
				input string
				ty    ast.ExprType
			}{
				{
					input: "true",
					ty:    ast.BasicType{Kind: kind.BOOL_TYPE},
				},
				{
					input: "false",
					ty:    ast.BasicType{Kind: kind.BOOL_TYPE},
				},
			},
		},
		{
			scope: &scope.Scope[ast.AstNode]{
				Parent: nil,
				Nodes: map[string]ast.AstNode{
					"a": ast.VarDeclStmt{
						Name: "a",
					},
				},
			},
			tests: []struct {
				input string
				ty    ast.ExprType
			}{
				{
					input: "true",
					ty:    ast.BasicType{Kind: kind.BOOL_TYPE},
				},
			},
		},
	}
	for _, test := range tests {
		for _, unit := range test.tests {
			t.Run(fmt.Sprintf("TestExprInference('%s')", unit.input), func(t *testing.T) {
				actualExprTy, err := analyzeExprType(unit.input, filename, test.scope)
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
