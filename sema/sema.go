package sema

import (
	"bufio"
	"fmt"
	"log"
	"reflect"
	"strings"

	"github.com/HicaroD/telia-lang/ast"
	"github.com/HicaroD/telia-lang/lexer"
	"github.com/HicaroD/telia-lang/lexer/token/kind"
	"github.com/HicaroD/telia-lang/parser"
	"github.com/HicaroD/telia-lang/scope"
)

type sema struct {
	universe *scope.Scope[ast.AstNode]
}

func New() *sema {
	// "universe" scope does not have any parent, it is the root of the tree of
	// scopes
	var nilScope *scope.Scope[ast.AstNode] = nil
	universe := scope.New(nilScope)
	return &sema{universe}
}

func (sema *sema) Analyze(astNodes []ast.AstNode) error {
	for i := range astNodes {
		switch astNode := astNodes[i].(type) {
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
			log.Fatalf("unimplemented ast node for sema: %s\n", reflect.TypeOf(astNode))
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
	for _, param := range function.Params.Fields {
		paramName := param.Name.Lexeme.(string)
		err := function.Scope.Insert(paramName, param)
		// TODO(errors)
		if err != nil {
			return err
		}
	}

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
			err := sema.analyzeFunctionCall(statement, currentScope)
			// TODO(errors)
			if err != nil {
				return err
			}
		case *ast.VarDeclStmt:
			err := sema.analyzeVarDecl(statement, currentScope)
			// TODO(errors)
			if err != nil {
				return err
			}
		case *ast.CondStmt:
			err := sema.analyzeCondStmt(statement, returnTy, currentScope)
			// TODO(errors)
			if err != nil {
				return err
			}
		case *ast.ReturnStmt:
			_, err := sema.getExprType(statement.Value, returnTy, currentScope)
			// TODO(errors)
			if err != nil {
				return err
			}
		default:
			log.Fatalf("unimplemented statement on sema: %s", statement)
		}
	}
	return nil
}

func (sema *sema) analyzeVarDecl(varDecl *ast.VarDeclStmt, scope *scope.Scope[ast.AstNode]) error {
	if varDecl.NeedsInference {
		// TODO(errors): need a test for it
		if varDecl.Type != nil {
			log.Fatalf("needs inference, but variable already has a type: %s", varDecl.Type)
		}
		exprType, err := sema.inferExprType(varDecl.Value, scope)
		// TODO(errors)
		if err != nil {
			return err
		}
		varDecl.Type = exprType
	} else {
		// TODO(errors)
		if varDecl.Type == nil {
			log.Fatalf("variable does not have a type and it said it does not need inference")
		}
		// TODO: what do I assert here in order to make it right?
		_, err := sema.getExprType(varDecl.Value, varDecl.Type, scope)
		// TODO(errors): Deal with type mismatch
		if err != nil {
			return err
		}
	}

	varName := varDecl.Name.Lexeme.(string)
	err := scope.Insert(varName, varDecl)
	if err != nil {
		return err
	}
	return nil
}

// Useful for testing
func analyzeVarDeclFrom(input, filename string) (*ast.VarDeclStmt, error) {
	reader := bufio.NewReader(strings.NewReader(input))

	lexer := lexer.New(filename, reader)
	tokens, err := lexer.Tokenize()
	if err != nil {
		return nil, err
	}

	par := parser.New(tokens)
	idStmt, err := par.ParseIdStmt()
	if err != nil {
		return nil, err
	}
	varDecl := idStmt.(*ast.VarDeclStmt)

	sema := New()

	parent := scope.New[ast.AstNode](nil)
	scope := scope.New(parent)

	err = sema.analyzeVarDecl(varDecl, scope)
	if err != nil {
		return nil, err
	}

	return varDecl, nil
}

func (sema *sema) analyzeCondStmt(condStmt *ast.CondStmt, returnTy ast.ExprType, outterScope *scope.Scope[ast.AstNode]) error {
	ifScope := scope.New(outterScope)

	err := sema.analyzeIfExpr(condStmt.IfStmt.Expr, outterScope)
	// TODO(errors)
	if err != nil {
		return err
	}

	err = sema.analyzeBlock(ifScope, condStmt.IfStmt.Block, returnTy)
	// TODO(errors)
	if err != nil {
		return err
	}

	for i := range condStmt.ElifStmts {
		elifScope := scope.New(outterScope)

		err := sema.analyzeIfExpr(condStmt.ElifStmts[i].Expr, elifScope)
		// TODO(errors)
		if err != nil {
			return err
		}
		err = sema.analyzeBlock(elifScope, condStmt.ElifStmts[i].Block, returnTy)
		// TODO(errors)
		if err != nil {
			return err
		}
	}

	if condStmt.ElseStmt != nil {
		elseScope := scope.New(outterScope)
		err = sema.analyzeBlock(elseScope, condStmt.ElseStmt.Block, returnTy)
		// TODO(errors)
		if err != nil {
			return err
		}
	}

	return nil
}

func (sema *sema) analyzeFunctionCall(functionCall *ast.FunctionCall, scope *scope.Scope[ast.AstNode]) error {
	function, err := scope.Lookup(functionCall.Name)
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
			if len(functionCall.Args) < minimumNumberOfArgs {
				log.Fatalf("the number of args is lesser than the minimum number of parameters")
			}

			for i := range minimumNumberOfArgs {
				paramType := decl.Params.Fields[i].Type
				argType, err := sema.getExprType(functionCall.Args[i], paramType, scope)
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
			if len(functionCall.Args) < minimumNumberOfArgs {
				log.Fatalf("the number of args is lesser than the minimum number of parameters")
			}

			for i := range minimumNumberOfArgs {
				paramType := decl.Params.Fields[i].Type
				argType, err := sema.getExprType(functionCall.Args[i], paramType, scope)
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
			case ast.BasicType:
				// TODO: refactor this
				// This switch seems unnecessary because it is quite repetitive
				// There are probably better ways to deal with integer sizes
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
				case kind.VOID_TYPE:
					// TODO(errors)
					return nil, fmt.Errorf("mismatched type, expected to be a %s, but got %s", expr, expectedType)
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
		case *ast.Field:
			return symTy.Type, nil
		// TODO(errors)
		default:
			log.Fatalf("expected to be a variable or parameter, but got %s", reflect.TypeOf(symTy))
		}
	case *ast.BinaryExpr:
		lhs, err := sema.getExprType(expr.Left, expectedType, scope)
		// TODO(errors)
		if err != nil {
			return nil, err
		}
		rhs, err := sema.getExprType(expr.Right, expectedType, scope)
		// TODO(errors)
		if err != nil {
			return nil, err
		}
		if !reflect.DeepEqual(lhs, rhs) {
			log.Fatalf("mismatched types: %s %s %s", lhs, expr.Op, rhs)
		}
	case *ast.FunctionCall:
		err := sema.analyzeFunctionCall(expr, scope)
		// TODO(errors)
		if err != nil {
			return nil, err
		}
		// At this point, function should exists!
		function, err := scope.Lookup(expr.Name)
		if err != nil {
			log.Fatalf("panic: at this point, function '%s' should exists in current block", expr.Name)
		}
		functionDecl := function.(*ast.FunctionDecl)
		return functionDecl.RetType, nil
	case *ast.VoidExpr:
		// TODO(errors)
		if !expectedType.IsVoid() {
			return nil, fmt.Errorf("expected return type to be '%s'", expectedType)
		}
		return expectedType, nil
	default:
		// TODO(errors)
		log.Fatalf("unimplemented on getExprType: %s", expr)
	}
	// NOTE: this line should be unreachable
	log.Fatalf("unreachable line at getExprTy")
	return nil, nil
}

func (sema *sema) inferExprType(expr ast.Expr, scope *scope.Scope[ast.AstNode]) (ast.ExprType, error) {
	switch expression := expr.(type) {
	case *ast.LiteralExpr:
		switch expression.Kind {
		case kind.STRING_LITERAL:
			return ast.PointerType{Type: ast.BasicType{Kind: kind.I8_TYPE}}, nil
		case kind.INTEGER_LITERAL:
			return sema.inferIntegerType(expression.Value.(int)), nil // TODO: is this correct?
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
		case *ast.Field:
			return node.Type, nil
		default:
			return nil, fmt.Errorf("symbol '%s' is not a variable", node)
		}
	case *ast.BinaryExpr:
		lhs, err := sema.inferExprType(expression.Left, scope)
		// TODO(errors)
		if err != nil {
			return nil, err
		}
		rhs, err := sema.inferExprType(expression.Right, scope)
		// TODO(errors)
		if err != nil {
			return nil, err
		}
		if !reflect.DeepEqual(lhs, rhs) {
			log.Fatalf("mismatched types: %s %s %s", lhs, expression.Op, rhs)
		}
		// TODO: make sure to check not only if the operands are of the same type,
		// but also if they can be put together with the operator
		switch expression.Op {
		default:
			if _, ok := ast.LOGICAL_OP[expression.Op]; ok {
				return ast.BasicType{Kind: kind.BOOL_TYPE}, nil
			}
			log.Fatalf("invalid operator: %s", expression.Op)
		}
	case *ast.FunctionCall:
		err := sema.analyzeFunctionCall(expression, scope)
		// TODO(errors)
		if err != nil {
			return nil, err
		}
		// At this point, function should exists!
		function, err := scope.Lookup(expression.Name)
		if err != nil {
			log.Fatalf("panic: at this point, function '%s' should exists in current block", expression.Name)
		}
		functionDecl := function.(*ast.FunctionDecl)
		return functionDecl.RetType, nil
	// TODO: deal with unary expressions
	default:
		log.Fatalf("unimplemented expression on sema: %s", expression)
	}
	// TODO(errors)
	// NOTE: this should be unreachable
	log.Fatalf("UNREACHABLE - inferExprType")
	return nil, nil
}

// TODO: deal with other types of integer
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
