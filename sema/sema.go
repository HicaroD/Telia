package sema

import (
	"fmt"
	"log"
	"reflect"

	"github.com/HicaroD/telia-lang/ast"
	"github.com/HicaroD/telia-lang/lexer/token/kind"
)

type sema struct {
	astNodes []ast.AstNode
	universe *ast.Scope
}

func New(astNodes []ast.AstNode) *sema {
	// "universe" scope does not have any parent, it is the root of the tree of
	// scopes
	universe := ast.NewScope(nil)
	return &sema{astNodes, universe}
}

func (sema *sema) Analyze() error {
	for i := range sema.astNodes {
		switch astNode := sema.astNodes[i].(type) {
		case *ast.FunctionDecl:
			err := sema.analyzeFnDecl(astNode)
			// TODO(errors)
			if err != nil {
				if err == ast.SYMBOL_ALREADY_DEFINED_ON_SCOPE {
					log.Fatalf("Function '%s' already defined on scope", astNode.Name)
				}
			}
		case *ast.ExternDecl:
			err := sema.analyzeExternDecl(astNode)
			// TODO(errors)
			if err != nil {
				if err == ast.SYMBOL_ALREADY_DEFINED_ON_SCOPE {
					log.Fatalf("Extern block '%s' already defined on scope", astNode.Name)
				}
			}
		default:
			// fmt.Printf("OTHER AST NODE: %s\n", sema.astNodes[i].String())
			fmt.Printf("OTHER AST NODE: %s\n", reflect.TypeOf(astNode))
		}
	}
	return nil
}

func (sema *sema) analyzeFnDecl(function *ast.FunctionDecl) error {
	err := sema.universe.Insert(function.Name, function)
	if err != nil {
		return err
	}
	function.Scope = ast.NewScope(sema.universe)

	for i := range function.Block.Statements {
		switch statement := function.Block.Statements[i].(type) {
		case *ast.FunctionCallStmt:
			function := function.Scope.Lookup(statement.Name)
			if function == nil {
				return ast.SYMBOL_NOT_FOUND_ON_SCOPE
			}
			// TODO(errors)
			if _, ok := function.(*ast.FunctionDecl); !ok {
				log.Fatalf("expected symbol to be a function, not %s", reflect.TypeOf(function))
			}

			functionDecl := function.(*ast.FunctionDecl)
			if functionDecl.Params.IsVariadic {
				minimumNumberOfArgs := len(functionDecl.Params.Fields)
				// TODO(errors)
				if len(statement.Args) < minimumNumberOfArgs {
					log.Fatalf("the number of args is lesser than the minimum number of parameters")
				}

				for i := range minimumNumberOfArgs {
					paramType := functionDecl.Params.Fields[i].Type
					argType, err := sema.getExprType(statement.Args[i], paramType)
					// TODO(errors)
					if err != nil {
						return err
					}
					// TODO(errors)
					if argType != paramType {
						log.Fatalf("mismatched argument type, expected %s, but got %s", paramType, argType)
					}
				}
			}
		// TODO: add function scope, but don't forget to check for redeclarations
		case *ast.VarStmt:
			if statement.NeedsInference {
				// TODO(errors): need a test for it
				if statement.Type != nil {
					log.Fatalf("needs inference, but variable already has a type: %s", statement.Type)
				}
				exprType, err := sema.inferExprType(statement.Value)
				// TODO(errors)
				if err != nil {
					return err
				}
				statement.Type = exprType
			} else {
				// TODO(errors)
				if statement.Type == nil {
					log.Fatalf("variable does not have a type and it said it does not need inference")
				}
				// TODO: what do I assert here in order to make it right?
				_, err := sema.getExprType(statement.Value, statement.Type)
				// TODO(errors): Deal with type mismatch
				if err != nil {
					return err
				}
			}

			return nil
		default:
			log.Fatalf("unimplemented statement: %s", statement)
		}
	}

	return nil
}

// TODO: think about the way I'm inferring or getting the expr type correctly

func (sema *sema) getExprType(exprNode ast.Expr, expectedType ast.ExprType) (ast.ExprType, error) {
	switch expr := exprNode.(type) {
	case *ast.LiteralExpr:
		switch expr.Kind {
		// TODO: is this right?
		case kind.STRING_LITERAL:
			return &ast.PointerType{Type: ast.BasicType{Kind: kind.I8_TYPE}}, nil
		// TODO: the idea is to deal with different kinds of integer literals
		// REFACTOR: this code could be simpler
		// TODO(errors): be careful with number overflow when trying to convert
		// the any to a number
		case kind.INTEGER_LITERAL:
			switch ty := expectedType.(type) {
			case *ast.BasicType:
				switch ty.Kind {
				case kind.I8_TYPE:
					value := expr.Value.(int64)
					// TODO(errors)
					if !(value >= -128 && value <= 127) {
						log.Fatalf("i8 integer overflow: %d", value)
					}
					return &ast.BasicType{Kind: kind.I8_TYPE}, nil
				case kind.I16_TYPE:
					value := expr.Value.(int64)
					// TODO(errors)
					if !(value >= -32768 && value <= 32767) {
						log.Fatalf("i16 integer overflow: %d", value)
					}
					return &ast.BasicType{Kind: kind.I16_TYPE}, nil
				case kind.I32_TYPE:
					value := expr.Value.(int64)
					// TODO(errors)
					if !(value >= -2147483648 && value <= 2147483647) {
						log.Fatalf("i32 integer overflow: %d", value)
					}
					return &ast.BasicType{Kind: kind.I32_TYPE}, nil
				case kind.I64_TYPE:
					value := expr.Value.(int64)
					// TODO(errors)
					if !(value >= -9223372036854775808 && value <= 9223372036854775807) {
						log.Fatalf("i64 integer overflow: %d", value)
					}
					return &ast.BasicType{Kind: kind.I64_TYPE}, nil
				default:
					// TODO(errors)
					log.Fatalf("expected to be an %s, but got %s", ty, expr.Kind)
				}
			}
		case kind.TRUE_BOOL_LITERAL, kind.FALSE_BOOL_LITERAL:
			return &ast.BasicType{Kind: kind.BOOL_TYPE}, nil
		}
	default:
		// TODO(errors)
		log.Fatalf("unimplemented: %s", expr)
	}
	// NOTE: this line should be unreachable
	return nil, nil
}

func (sema *sema) inferExprType(expr ast.Expr) (ast.ExprType, error) {
	switch expression := expr.(type) {
	case *ast.LiteralExpr:
		switch expression.Kind {
		case kind.INTEGER_LITERAL:
			return sema.inferIntegerType(expression.Value.(int64)), nil
		// TODO: is this correct?
		case kind.STRING_LITERAL:
			return &ast.PointerType{Type: ast.BasicType{Kind: kind.I8_TYPE}}, nil
		case kind.TRUE_BOOL_LITERAL, kind.FALSE_BOOL_LITERAL:
			return &ast.BasicType{Kind: kind.BOOL_TYPE}, nil
		}
	}
	// TODO(errors)
	// NOTE: this should be unreachable
	log.Fatalf("UNREACHABLE - inferExprType")
	return nil, nil
}

func (sema *sema) inferIntegerType(value int64) ast.ExprType {
	integerType := kind.I32_TYPE
	if value >= 2147483647 { // Max i32 size
		integerType = kind.I64_TYPE
	}
	// TODO: deal with i64 overflow
	return &ast.BasicType{Kind: integerType}
}

func (sema *sema) analyzeExternDecl(extern *ast.ExternDecl) error {
	err := sema.universe.Insert(extern.Name.Lexeme.(string), extern)
	if err != nil {
		return err
	}
	extern.Scope = ast.NewScope(sema.universe)

	return nil
}
