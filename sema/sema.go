package sema

import (
	"bufio"
	"fmt"
	"log"
	"reflect"
	"strconv"
	"strings"

	"github.com/HicaroD/Telia/ast"
	"github.com/HicaroD/Telia/lexer"
	"github.com/HicaroD/Telia/lexer/token/kind"
	"github.com/HicaroD/Telia/parser"
	"github.com/HicaroD/Telia/scope"
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
			_, err := sema.inferExprTypeWithContext(statement.Value, returnTy, currentScope)
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
		exprType, err := sema.inferExprTypeWithoutContext(varDecl.Value, scope)
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
		_, err := sema.inferExprTypeWithContext(varDecl.Value, varDecl.Type, scope)
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
				argType, err := sema.inferExprTypeWithContext(functionCall.Args[i], paramType, scope)
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
				argType, err := sema.inferExprTypeWithContext(functionCall.Args[i], paramType, scope)
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
	inferedExpr, err := sema.inferExprTypeWithoutContext(expr, scope)
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
func (sema *sema) inferExprTypeWithContext(exprNode ast.Expr, expectedType ast.ExprType, scope *scope.Scope[ast.AstNode]) (ast.ExprType, error) {
	switch expr := exprNode.(type) {
	case *ast.LiteralExpr:
		switch expr.Kind {
		case kind.STRING_LITERAL:
			return ast.PointerType{Type: ast.BasicType{Kind: kind.I8_TYPE}}, nil
		case kind.INTEGER_LITERAL:
			switch ty := expectedType.(type) {
			case ast.BasicType:
				value := expr.Value.(string)
				bitSize := ty.Kind.BitSize()
				if bitSize == -1 {
					return nil, fmt.Errorf("not a valid numeric type: %s", ty.Kind)
				}
				_, err := strconv.ParseInt(value, bitSize, 8)
				// TODO(errors)
				if err != nil {
					return nil, fmt.Errorf("%s integer overflow: %s", ty.Kind, value)
				}
				return ast.BasicType{Kind: ty.Kind}, nil
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
		lhs, err := sema.inferExprTypeWithContext(expr.Left, expectedType, scope)
		// TODO(errors)
		if err != nil {
			return nil, err
		}
		rhs, err := sema.inferExprTypeWithContext(expr.Right, expectedType, scope)
		// TODO(errors)
		if err != nil {
			return nil, err
		}
		if !reflect.DeepEqual(lhs, rhs) {
			log.Fatalf("mismatched types: %s %s %s", lhs, expr.Op, rhs)
		}
		return lhs, nil
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

// Useful for testing
func inferExprTypeWithoutContext(input, filename string, scope *scope.Scope[ast.AstNode]) (ast.Expr, ast.ExprType, error) {
	expr, err := parser.ParseExprFrom(input, filename)
	if err != nil {
		return nil, nil, err
	}

	analyzer := New()
	exprType, err := analyzer.inferExprTypeWithoutContext(expr, scope)
	if err != nil {
		return nil, nil, err
	}
	return expr, exprType, nil
}

func inferExprTypeWithContext(input, filename string, ty ast.ExprType, scope *scope.Scope[ast.AstNode]) (ast.ExprType, error) {
	expr, err := parser.ParseExprFrom(input, filename)
	if err != nil {
		return nil, err
	}

	exprType, err := New().inferExprTypeWithContext(expr, ty, scope)
	if err != nil {
		return nil, err
	}
	return exprType, nil
}

func (sema *sema) inferExprTypeWithoutContext(expr ast.Expr, scope *scope.Scope[ast.AstNode]) (ast.ExprType, error) {
	switch expression := expr.(type) {
	case *ast.LiteralExpr:
		switch expression.Kind {
		// TODO: what if I want a defined string type on my programming
		// language? I don't think *i8 is correct type for this
		case kind.STRING_LITERAL:
			return ast.PointerType{Type: ast.BasicType{Kind: kind.I8_TYPE}}, nil
		case kind.INTEGER_LITERAL:
			ty, err := sema.inferIntegerType(expression.Value.(string))
			if err != nil {
				return nil, err
			}
			return ty, nil
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
	case *ast.UnaryExpr:
		switch expression.Op {
		case kind.MINUS:
			switch unaryExpr := expression.Node.(type) {
			case *ast.LiteralExpr:
				if unaryExpr.Kind == kind.INTEGER_LITERAL {
					integerType, err := sema.inferIntegerType(unaryExpr.Value.(string))
					// TODO(errors)
					if err != nil {
						return nil, err
					}
					return integerType, nil
				}
			}
		case kind.NOT:
			unaryExpr, err := sema.inferExprTypeWithoutContext(expression.Expr, scope)
			// TODO(errors)
			if err != nil {
				return nil, err
			}
			if !unaryExpr.IsBoolean() {
				return nil, fmt.Errorf("expected boolean expression on not unary expression")
			}
			return unaryExpr, nil
		}
		// TODO(errors)
		log.Fatalf("unimplemented unary expression on sema")
	case *ast.BinaryExpr:
		ty, err := sema.inferBinaryExprType(expression, scope)
		// TODO(errors)
		if err != nil {
			return nil, err
		}
		return ty, nil
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
	default:
		log.Fatalf("unimplemented expression on sema: %s", reflect.TypeOf(expression))
	}
	// TODO(errors)
	// NOTE: this should be unreachable
	log.Fatalf("UNREACHABLE - inferExprType")
	return nil, nil
}

func (sema *sema) inferBinaryExprType(expression *ast.BinaryExpr, scope *scope.Scope[ast.AstNode]) (ast.ExprType, error) {
	lhs, err := sema.inferExprTypeWithoutContext(expression.Left, scope)
	// TODO(errors)
	if err != nil {
		return nil, err
	}

	rhs, err := sema.inferExprTypeWithoutContext(expression.Right, scope)
	// TODO(errors)
	if err != nil {
		return nil, err
	}

	if !reflect.DeepEqual(lhs, rhs) {
		log.Fatalf("mismatched types: %s %s %s", lhs, expression.Op, rhs)
	}

	switch expression.Op {
	case kind.PLUS, kind.MINUS:
		if lhs.IsNumeric() && rhs.IsNumeric() {
			return lhs, nil
		}
	default:
		if _, ok := kind.LOGICAL_OP[expression.Op]; ok {
			return ast.BasicType{Kind: kind.BOOL_TYPE}, nil
		}
	}
	// TODO
	return nil, nil
}

// TODO: deal with other types of integer
func (sema *sema) inferIntegerType(value string) (ast.ExprType, error) {
	integerType := kind.INT_TYPE
	base := 10

	_, err := strconv.ParseUint(value, base, 32)
	if err == nil {
		return ast.BasicType{Kind: integerType}, nil
	}

	_, err = strconv.ParseUint(value, base, 64)
	if err == nil {
		return ast.BasicType{Kind: integerType}, nil
	}

	return nil, fmt.Errorf("can't parse integer 32 bits: %s", value)
}

// func (sema *sema) inferSignedIntegerType(value string) (ast.ExprType, error) {
// 	integerType := kind.I32_TYPE
// 	base := 10
//
// 	_, err := strconv.ParseInt(value, base, 32)
// 	if err == nil {
// 		return ast.BasicType{Kind: integerType}, nil
// 	}
//
// 	_, err = strconv.ParseInt(value, base, 64)
// 	if err == nil {
// 		return ast.BasicType{Kind: integerType}, nil
// 	}
//
// 	return nil, fmt.Errorf("can't parse signed integer: %s", value)
// }

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
