package sema

import (
	"fmt"
	"log"
	"reflect"

	"github.com/HicaroD/telia-lang/ast"
	"github.com/HicaroD/telia-lang/lexer/token/kind"
	"github.com/HicaroD/telia-lang/scope"
)

type sema struct {
	astNodes []ast.AstNode
	universe *scope.Scope[ast.AstNode]
}

func New(astNodes []ast.AstNode) *sema {
	// "universe" scope does not have any parent, it is the root of the tree of
	// scopes
	var nilScope *scope.Scope[ast.AstNode] = nil
	universe := scope.New(nilScope)
	return &sema{astNodes, universe}
}

func (sema *sema) Analyze() error {
	for i := range sema.astNodes {
		switch astNode := sema.astNodes[i].(type) {
		case *ast.FunctionDecl:
			err := sema.analyzeFnDecl(astNode)
			// TODO(errors)
			if err != nil {
				return err
			}
		case *ast.ExternDecl:
			err := sema.analyzeExternDecl(astNode)
			// TODO(errors)
			if err != nil {
				return err
			}
		default:
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

	function.Scope = scope.New(sema.universe)
	err = sema.analyzeBlock(function.Scope, function.Block, function.RetType)
	// TODO(errors)
	if err != nil {
		return err
	}
	return nil
}

func (sema *sema) analyzeBlock(currentScope *scope.Scope[ast.AstNode], block *ast.BlockStmt, returnTy ast.ExprType) error {
	for i := range block.Statements {
		switch statement := block.Statements[i].(type) {
		case *ast.FunctionCall:
			function, err := currentScope.Lookup(statement.Name)
			// TODO(errors)
			if err != nil {
				return err
			}
			// REFACTOR: basically the same code for function decl and
			// prototypes
			switch decl := function.(type) {
			case *ast.FunctionDecl:
				if decl.Params.IsVariadic {
					minimumNumberOfArgs := len(decl.Params.Fields)
					// TODO(errors)
					if len(statement.Args) < minimumNumberOfArgs {
						log.Fatalf("the number of args is lesser than the minimum number of parameters")
					}

					for i := range minimumNumberOfArgs {
						paramType := decl.Params.Fields[i].Type
						argType, err := sema.getExprType(statement.Args[i], paramType, currentScope)
						// TODO(errors)
						if err != nil {
							return err
						}
						// TODO(errors)
						if argType != paramType {
							log.Fatalf("mismatched argument type on function '%s', expected %s, but got %s", decl.Name, paramType, argType)
						}
					}
				}
				// TODO: deal with non variadic functions
			case *ast.Proto:
				if decl.Params.IsVariadic {
					minimumNumberOfArgs := len(decl.Params.Fields)
					// TODO(errors)
					if len(statement.Args) < minimumNumberOfArgs {
						log.Fatalf("the number of args is lesser than the minimum number of parameters")
					}

					for i := range minimumNumberOfArgs {
						paramType := decl.Params.Fields[i].Type
						argType, err := sema.getExprType(statement.Args[i], paramType, currentScope)
						// TODO(errors)
						if err != nil {
							return err
						}
						// TODO(errors)
						if !reflect.DeepEqual(argType, paramType) {
							log.Fatalf("mismatched argument type on prototype '%s', expected %s, but got %s", decl.Name, paramType, argType)
						}
					}
				}
				// TODO: deal with non variadic prototypes
			default:
				// TODO(errors)
				log.Fatalf("expected symbol to be a function or proto, not %s", reflect.TypeOf(function))
			}
		// TODO: add function scope, but don't forget to check for redeclarations
		case *ast.VarDeclStmt:
			if statement.NeedsInference {
				// TODO(errors): need a test for it
				if statement.Type != nil {
					log.Fatalf("needs inference, but variable already has a type: %s", statement.Type)
				}
				exprType, err := sema.inferExprType(statement.Value, currentScope)
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
				_, err := sema.getExprType(statement.Value, statement.Type, currentScope)
				// TODO(errors): Deal with type mismatch
				if err != nil {
					return err
				}
			}

			varName := statement.Name.Lexeme.(string)
			err := currentScope.Insert(varName, statement)
			if err != nil {
				return err
			}
		case *ast.CondStmt:
			err := sema.analyzeIfExpr(statement.IfStmt.Expr, currentScope)
			// TODO(errors)
			if err != nil {
				return err
			}

			ifScope := scope.New(currentScope)
			err = sema.analyzeBlock(ifScope, statement.IfStmt.Block, returnTy)
			// TODO(errors)
			if err != nil {
				return err
			}

			for i := range statement.ElifStmts {
				err := sema.analyzeIfExpr(statement.ElifStmts[i].Expr, currentScope)
				// TODO(errors)
				if err != nil {
					return err
				}

				elifScope := scope.New(currentScope)
				err = sema.analyzeBlock(elifScope, statement.ElseStmt.Block, returnTy)
				// TODO(errors)
				if err != nil {
					return err
				}
			}
		case *ast.ReturnStmt:
			_, err := sema.getExprType(statement.Value, returnTy, currentScope)
			// TODO(errors)
			if err != nil {
				return err
			}
		default:
			log.Fatalf("unimplemented statement: %s", statement)
		}
	}
	return nil
}

func (sema *sema) analyzeIfExpr(expr ast.Expr, scope *scope.Scope[ast.AstNode]) error {
	inferedExpr, err := sema.inferExprType(expr, scope)
	// TODO(errors)
	if err != nil {
		return err
	}
	// TODO(errors)
	if !sema.isValidExprToBeOnIf(inferedExpr) {
		// TODO(errors)
		log.Fatalf("invalid non-boolean condition on if statement: %s", inferedExpr)
	}
	return nil
}

func (sema *sema) isValidExprToBeOnIf(exprType ast.ExprType) bool {
	switch ty := exprType.(type) {
	case ast.BasicType:
		switch ty.Kind {
		case kind.BOOL_TYPE:
			return true
		// TODO: deal with binary expressions
		default:
			return false
		}
	}
	return false
}

// TODO: think about the way I'm inferring or getting the expr type correctly

func (sema *sema) getExprType(exprNode ast.Expr, expectedType ast.ExprType, scope *scope.Scope[ast.AstNode]) (ast.ExprType, error) {
	switch expr := exprNode.(type) {
	case *ast.LiteralExpr:
		switch expr.Kind {
		case kind.STRING_LITERAL:
			return ast.PointerType{Type: ast.BasicType{Kind: kind.I8_TYPE}}, nil
		case kind.INTEGER_LITERAL:
			switch ty := expectedType.(type) {
			case *ast.BasicType:
				// TODO: refactor this
				// This switch seems unnecessary because it is quite repetitive
				switch ty.Kind {
				case kind.I8_TYPE:
					value := expr.Value.(int)
					// TODO(errors)
					if !(value >= -128 && value <= 127) {
						log.Fatalf("i8 integer overflow: %d", value)
					}
					return ast.BasicType{Kind: kind.I8_TYPE}, nil
				case kind.I16_TYPE:
					value := expr.Value.(int)
					// TODO(errors)
					if !(value >= -32768 && value <= 32767) {
						log.Fatalf("i16 integer overflow: %d", value)
					}
					return ast.BasicType{Kind: kind.I16_TYPE}, nil
				case kind.I32_TYPE:
					value := expr.Value.(int)
					// TODO(errors)
					if !(value >= -2147483648 && value <= 2147483647) {
						log.Fatalf("i32 integer overflow: %d", value)
					}
					return ast.BasicType{Kind: kind.I32_TYPE}, nil
				case kind.I64_TYPE:
					value := expr.Value.(int)
					// TODO(errors)
					if !(value >= -9223372036854775808 && value <= 9223372036854775807) {
						log.Fatalf("i64 integer overflow: %d", value)
					}
					return ast.BasicType{Kind: kind.I64_TYPE}, nil
				default:
					// TODO(errors)
					log.Fatalf("expected to be an %s, but got %s", ty, expr.Kind)
				}
			}
		case kind.TRUE_BOOL_LITERAL, kind.FALSE_BOOL_LITERAL:
			return ast.BasicType{Kind: kind.BOOL_TYPE}, nil
		}
	case *ast.IdExpr:
		symbol, err := scope.Lookup(expr.Name.Lexeme.(string))
		if err != nil {
			return nil, err
		}
		switch symTy := symbol.(type) {
		case *ast.VarDeclStmt:
			return symTy.Type, nil
		// TODO(errors)
		default:
			log.Fatalf("expected to be a variable, but got %s", reflect.TypeOf(symTy))
		}
	default:
		// TODO(errors)
		log.Fatalf("unimplemented: %s", expr)
	}
	// NOTE: this line should be unreachable
	return nil, nil
}

func (sema *sema) inferExprType(expr ast.Expr, scope *scope.Scope[ast.AstNode]) (ast.ExprType, error) {
	switch expression := expr.(type) {
	case *ast.LiteralExpr:
		switch expression.Kind {
		case kind.STRING_LITERAL:
			return ast.PointerType{Type: ast.BasicType{Kind: kind.I8_TYPE}}, nil
		case kind.INTEGER_LITERAL, kind.NEGATIVE_INTEGER_LITERAL:
			return sema.inferIntegerType(expression.Value.(int)), nil
		// TODO: is this correct?
		case kind.TRUE_BOOL_LITERAL, kind.FALSE_BOOL_LITERAL:
			return ast.BasicType{Kind: kind.BOOL_TYPE}, nil
		default:
			log.Fatalf("unimplemented literal expr: %s", expression)
		}
	case *ast.IdExpr:
		variableName := expression.Name.Lexeme.(string)
		variable, err := scope.Lookup(variableName)
		// TODO(errors)
		if err != nil {
			return nil, err
		}

		switch node := variable.(type) {
		case *ast.VarDeclStmt:
			return node.Type, nil
		default:
			return nil, fmt.Errorf("symbol '%s' is not a variable", node)
		}
	default:
		log.Fatalf("unimplemented expression: %s", expression)
	}
	// TODO(errors)
	// NOTE: this should be unreachable
	log.Fatalf("UNREACHABLE - inferExprType")
	return nil, nil
}

func (sema *sema) inferIntegerType(value int) ast.ExprType {
	integerType := kind.I32_TYPE
	if value >= 2147483647 { // Max i32 size
		integerType = kind.I64_TYPE
	}
	// TODO: deal with i64 overflow
	return ast.BasicType{Kind: integerType}
}

func (sema *sema) analyzeExternDecl(extern *ast.ExternDecl) error {
	err := sema.universe.Insert(extern.Name.Lexeme.(string), extern)
	if err != nil {
		return err
	}
	extern.Scope = scope.New(sema.universe)
	// NOTE: this move is temporary, the idea is to access prototypes
	// methods using something like "libc.printf()" where libc is the name
	// of extern block and "printf" is one of the functions defined inside
	// the "libc" extern block
	// I'm only adding to Universe for testing purposes
	for i := range extern.Prototypes {
		err := sema.universe.Insert(extern.Prototypes[i].Name, extern.Prototypes[i])
		// TODO(errors)
		if err != nil {
			return err
		}
	}
	return nil
}
