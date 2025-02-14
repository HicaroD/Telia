package sema

import (
	"fmt"
	"log"
	"reflect"
	"strconv"

	"github.com/HicaroD/Telia/internal/diagnostics"
	"github.com/HicaroD/Telia/internal/frontend/ast"
	"github.com/HicaroD/Telia/internal/frontend/lexer"
	"github.com/HicaroD/Telia/internal/frontend/lexer/token"
	"github.com/HicaroD/Telia/internal/frontend/parser"
)

type sema struct {
	collector *diagnostics.Collector
}

func New(collector *diagnostics.Collector) *sema {
	return &sema{collector}
}

func (s *sema) Check(program *ast.Program) error {
	return s.checkModule(program.Root)
}

func (s *sema) checkModule(module *ast.Package) error {
	for _, file := range module.Files {
		err := s.checkFile(file)
		if err != nil {
			return err
		}
	}

	for _, innerPackage := range module.Packages {
		err := s.checkModule(innerPackage)
		if err != nil {
			return err
		}
	}
	return nil
}

func (s *sema) checkFile(file *ast.File) error {
	for _, node := range file.Body {
		switch n := node.(type) {
		case *ast.FunctionDecl:
			err := s.checkFnDecl(n)
			if err != nil {
				return err
			}
		case *ast.ExternDecl:
			err := s.checkExternDecl(n)
			if err != nil {
				return err
			}
		case *ast.PkgDecl:
			continue
		default:
			log.Fatalf("unimplemented ast node for sema: %s\n", reflect.TypeOf(n))
		}
	}
	return nil
}

func (sema *sema) checkFnDecl(function *ast.FunctionDecl) error {
	err := sema.checkBlock(function.Block, function.RetType, function.Scope)
	return err
}

func (sema *sema) checkExternAttributes(externAttributes *ast.ExternAttrs) {
	// TODO: make sure attributes values are valid
}

func (sema *sema) checkExternDecl(extern *ast.ExternDecl) error {
	// TODO: check for attributes

	for _, proto := range extern.Prototypes {
		symbols := make(map[string]bool, len(proto.Params.Fields))
		for _, param := range proto.Params.Fields {
			if _, found := symbols[param.Name.Name()]; found {
				// TODO(errors): add proper error here + tests
				return fmt.Errorf("redeclaration of '%s' parameter on '%s' prototype at extern declaration '%s'\n", param.Name.Name(), proto.Name.Name(), extern.Name.Name())
			}
			symbols[param.Name.Name()] = true
		}
	}

	return nil
}

func (sema *sema) checkBlock(
	block *ast.BlockStmt,
	returnTy ast.ExprType,
	scope *ast.Scope,
) error {
	for i := range block.Statements {
		err := sema.checkStmt(block.Statements[i], scope, returnTy)
		if err != nil {
			return err
		}
	}
	return nil
}

func (sema *sema) checkStmt(
	stmt ast.Stmt,
	scope *ast.Scope,
	returnTy ast.ExprType,
) error {
	switch statement := stmt.(type) {
	case *ast.FunctionCall:
		err := sema.checkFunctionCall(statement, scope)
		return err
	case *ast.MultiVarStmt, *ast.VarStmt:
		err := sema.checkVarDecl(statement, scope)
		return err
	case *ast.CondStmt:
		err := sema.checkCondStmt(statement, returnTy, scope)
		return err
	case *ast.ReturnStmt:
		_, err := sema.inferExprTypeWithContext(statement.Value, returnTy, scope)
		return err
	case *ast.FieldAccess:
		_, err := sema.checkFieldAccessExpr(statement, scope)
		return err
	case *ast.ForLoop:
		err := sema.checkForLoop(statement, scope, returnTy)
		return err
	case *ast.WhileLoop:
		err := sema.checkWhileLoop(statement, scope, returnTy)
		return err
	default:
		log.Fatalf("unimplemented statement on sema: %s", statement)
	}
	return nil
}

func (sema *sema) checkVarDecl(
	variable ast.Stmt,
	currentScope *ast.Scope,
) error {
	switch varStmt := variable.(type) {
	case *ast.MultiVarStmt:
		err := sema.checkMultiVar(varStmt, currentScope)
		return err
	case *ast.VarStmt:
		err := sema.checkVar(varStmt, currentScope)
		return err
	}
	return nil
}

func (sema *sema) checkMultiVar(
	multi *ast.MultiVarStmt,
	currentScope *ast.Scope,
) error {
	// TODO: refactor this code
	// It is pretty repetitive, but I need tests before refactoring

	// if multi.IsDecl {
	// 	allVariablesDefined := true
	// 	for i := range multi.Variables {
	// 		_, err := currentScope.LookupCurrentScope(multi.Variables[i].Name.Name())
	// 		if err != nil {
	// 			if err == ast.ERR_SYMBOL_NOT_FOUND_ON_SCOPE {
	// 				allVariablesDefined = false
	// 				break
	// 			}
	// 			return err
	// 		}
	// 		multi.Variables[i].Decl = false
	// 	}
	// 	if allVariablesDefined {
	// 		firstVariable := multi.Variables[0]
	// 		pos := firstVariable.Name.Pos
	// 		// TODO: give user a hint for consider using = instead of :=
	// 		noNewVariablesDeclared := diagnostics.Diag{
	// 			Message: fmt.Sprintf(
	// 				"%s:%d:%d: no new variables declared",
	// 				pos.Filename,
	// 				pos.Line,
	// 				pos.Column,
	// 			),
	// 		}
	// 		sema.collector.ReportAndSave(noNewVariablesDeclared)
	// 		return diagnostics.COMPILER_ERROR_FOUND
	// 	}
	// } else {
	// 	allVariablesDefined := true
	//
	// 	var undefinedVar *ast.VarStmt
	// 	for i := range multi.Variables {
	// 		_, err := currentScope.LookupAcrossScopes(multi.Variables[i].Name.Name())
	// 		if err != nil {
	// 			if err == ast.ERR_SYMBOL_NOT_FOUND_ON_SCOPE {
	// 				allVariablesDefined = false
	// 				undefinedVar = multi.Variables[i]
	// 				break
	// 			}
	// 			return err
	// 		}
	// 	}
	// 	if !allVariablesDefined {
	// 		if undefinedVar == nil {
	// 			log.Fatal("panic: variable at non decl is nil, but it should never be")
	// 		}
	// 		pos := undefinedVar.Name.Pos
	// 		notDeclared := diagnostics.Diag{
	// 			Message: fmt.Sprintf("%s:%d:%d: '%s' not declared", pos.Filename, pos.Line, pos.Column, undefinedVar.Name.Name()),
	// 		}
	// 		sema.collector.ReportAndSave(notDeclared)
	// 		return diagnostics.COMPILER_ERROR_FOUND
	// 	}
	// }

	for i := range multi.Variables {
		// if multi.Variables[i].Decl {
		// 	varName := multi.Variables[i].Name.Name()
		// 	err := currentScope.Insert(varName, multi.Variables[i])
		// 	if err != nil {
		// 		return err
		// 	}
		// }

		// err := sema.checkVariableType(multi.Variables[i], currentScope)
		// // TODO(errors)
		// if err != nil {
		// 	return err
		// }

		err := sema.checkVar(multi.Variables[i], currentScope)
		if err != nil {
			return err
		}
	}
	return nil
}

func (sema *sema) checkVar(variable *ast.VarStmt, currentScope *ast.Scope) error {
	err := sema.checkVariableType(variable, currentScope)
	return err
}

func (sema *sema) checkVariableType(
	varDecl *ast.VarStmt,
	currentScope *ast.Scope,
) error {
	if varDecl.NeedsInference {
		// TODO(errors): need a test for it
		if varDecl.Type != nil {
			return fmt.Errorf(
				"needs inference, but variable already has a type: %s",
				varDecl.Type,
			)
		}
		exprType, _, err := sema.inferExprTypeWithoutContext(varDecl.Value, currentScope)
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
		exprTy, err := sema.inferExprTypeWithContext(varDecl.Value, varDecl.Type, currentScope)
		// TODO(errors): Deal with type mismatch
		if err != nil {
			return err
		}
		if !reflect.DeepEqual(varDecl.Type, exprTy) {
			return fmt.Errorf("type mismatch on variable decl, expected %s, got %s", varDecl.Type, exprTy)
		}
	}
	return nil
}

// Useful for testing
func checkVarDeclFrom(input, filename string) (ast.Stmt, error) {
	collector := diagnostics.New()

	src := []byte(input)
	lexer := lexer.New(filename, src, collector)
	par := parser.NewWithLex(lexer, collector)

	tmpScope := ast.NewScope(nil)
	varStmt, err := par.ParseIdStmt(tmpScope)
	if err != nil {
		return nil, err
	}

	sema := New(collector)
	parent := ast.NewScope(nil)
	scope := ast.NewScope(parent)
	err = sema.checkVarDecl(varStmt, scope)
	if err != nil {
		return nil, err
	}

	return varStmt, nil
}

func (sema *sema) checkCondStmt(
	condStmt *ast.CondStmt,
	returnTy ast.ExprType,
	outterScope *ast.Scope,
) error {
	err := sema.checkIfExpr(condStmt.IfStmt.Expr, outterScope)
	// TODO(errors)
	if err != nil {
		return err
	}

	err = sema.checkBlock(condStmt.IfStmt.Block, returnTy, condStmt.IfStmt.Scope)
	// TODO(errors)
	if err != nil {
		return err
	}

	for i := range condStmt.ElifStmts {
		err := sema.checkIfExpr(condStmt.ElifStmts[i].Expr, condStmt.ElifStmts[i].Scope)
		// TODO(errors)
		if err != nil {
			return err
		}
		err = sema.checkBlock(condStmt.ElifStmts[i].Block, returnTy, condStmt.ElifStmts[i].Scope)
		// TODO(errors)
		if err != nil {
			return err
		}
	}

	if condStmt.ElseStmt != nil {
		err = sema.checkBlock(condStmt.ElseStmt.Block, returnTy, condStmt.ElseStmt.Scope)
		// TODO(errors)
		if err != nil {
			return err
		}
	}

	return nil
}

func (sema *sema) checkFunctionCall(
	functionCall *ast.FunctionCall,
	currentScope *ast.Scope,
) error {
	function, err := currentScope.LookupAcrossScopes(functionCall.Name.Name())
	if err != nil {
		if err == ast.ERR_SYMBOL_NOT_FOUND_ON_SCOPE {
			pos := functionCall.Name.Pos
			functionNotDefined := diagnostics.Diag{
				Message: fmt.Sprintf(
					"%s:%d:%d: function '%s' not defined on scope",
					pos.Filename,
					pos.Line,
					pos.Column,
					functionCall.Name.Name(),
				),
			}
			sema.collector.ReportAndSave(functionNotDefined)
			return diagnostics.COMPILER_ERROR_FOUND
		}
		return err
	}

	decl, ok := function.(*ast.FunctionDecl)
	if !ok {
		pos := functionCall.Name.Pos
		notCallable := diagnostics.Diag{
			Message: fmt.Sprintf(
				"%s:%d:%d: '%s' is not callable",
				pos.Filename,
				pos.Line,
				pos.Column,
				functionCall.Name.Name(),
			),
		}
		sema.collector.ReportAndSave(notCallable)
		return diagnostics.COMPILER_ERROR_FOUND
	}

	if len(functionCall.Args) != len(decl.Params.Fields) {
		pos := functionCall.Name.Pos
		// TODO(errors): show which arguments were passed and which types we
		// were expecting
		notEnoughArguments := diagnostics.Diag{
			Message: fmt.Sprintf(
				"%s:%d:%d: not enough arguments in call to '%s'",
				pos.Filename,
				pos.Line,
				pos.Column,
				functionCall.Name.Name(),
			),
		}
		sema.collector.ReportAndSave(notEnoughArguments)
		return diagnostics.COMPILER_ERROR_FOUND
	}

	for i := range len(functionCall.Args) {
		paramType := decl.Params.Fields[i].Type
		argType, err := sema.inferExprTypeWithContext(functionCall.Args[i], paramType, currentScope)
		if err != nil {
			return err
		}
		if !reflect.DeepEqual(argType, paramType) {
			mismatchedArgType := diagnostics.Diag{
				// TODO(errors): add position of the error
				Message: fmt.Sprintf("can't use %s on argument of type %s", argType, paramType),
			}
			sema.collector.ReportAndSave(mismatchedArgType)
			return diagnostics.COMPILER_ERROR_FOUND
		}
	}

	// TODO: deal with variadic arguments
	return nil
}

func (sema *sema) checkIfExpr(expr ast.Expr, scope *ast.Scope) error {
	inferedExprType, _, err := sema.inferExprTypeWithoutContext(expr, scope)
	// TODO(errors)
	if err != nil {
		return err
	}

	if !inferedExprType.IsBoolean() {
		// TODO(errors)
		log.Fatalf("invalid non-boolean condition on if statement: %s", inferedExprType)
	}

	return nil
}

// TODO: think about the way I'm inferring or getting the expr type correctly
func (sema *sema) inferExprTypeWithContext(
	exprNode ast.Expr,
	expectedType ast.ExprType,
	scope *ast.Scope,
) (ast.ExprType, error) {
	switch expression := exprNode.(type) {
	case *ast.LiteralExpr:
		switch ty := expression.Type.(type) {
		case *ast.BasicType:
			switch ty.Kind {
			case token.STRING_LITERAL:
				finalTy := &ast.PointerType{Type: &ast.BasicType{Kind: token.U8_TYPE}}
				expression.Type = finalTy
				return finalTy, nil
			case token.INTEGER_LITERAL, token.INT_TYPE:
				switch ty := expectedType.(type) {
				case *ast.BasicType:
					value := string(expression.Value)
					bitSize := ty.Kind.BitSize()
					// TODO(errors)
					if bitSize == -1 {
						return nil, fmt.Errorf("not a valid numeric type: %s", ty.Kind)
					}
					_, err := strconv.ParseUint(value, 10, bitSize)
					// TODO(errors)
					if err != nil {
						return nil, fmt.Errorf("kind: %s bitSize: %d - integer overflow %s - error: %s", ty.Kind, bitSize, value, err)
					}
					finalTy := &ast.BasicType{Kind: ty.Kind}
					expression.Type = finalTy
					return finalTy, nil
				default:
					// TODO(errors)
					log.Fatalf("unimplemented type at integer literal: %s %s", ty, reflect.TypeOf(ty))
				}
			case token.TRUE_BOOL_LITERAL:
				finalTy := &ast.BasicType{Kind: token.BOOL_TYPE}
				expression.Type = finalTy
				expression.Value = []byte("1")
				return finalTy, nil
			case token.FALSE_BOOL_LITERAL:
				finalTy := &ast.BasicType{Kind: token.BOOL_TYPE}
				expression.Type = finalTy
				expression.Value = []byte("0")
				return finalTy, nil
			default:
				// TODO(errors)
				log.Fatalf("unimplemented basic type kind: %s", ty.Kind)
			}
		default:
			// TODO(errors)
			log.Fatalf("unimplemented literal expr type on inferExprTypeWithContext: %s", reflect.TypeOf(ty))
		}
	case *ast.IdExpr:
		symbol, err := scope.LookupAcrossScopes(expression.Name.Name())
		// TODO(errors)
		if err != nil {
			return nil, err
		}
		switch symTy := symbol.(type) {
		case *ast.VarStmt:
			return symTy.Type, nil
		case *ast.Field:
			return symTy.Type, nil
		default:
			// TODO(errors)
			log.Fatalf("expected to be a variable or parameter, but got %s", reflect.TypeOf(symTy))
		}
	case *ast.BinaryExpr:
		ty, err := sema.inferBinaryExprTypeWithContext(expression, expectedType, scope)
		// TODO(errors)
		if err != nil {
			return nil, err
		}
		return ty, nil
	case *ast.FunctionCall:
		err := sema.checkFunctionCall(expression, scope)
		// TODO(errors)
		if err != nil {
			return nil, err
		}
		function, _ := scope.LookupAcrossScopes(expression.Name.Name())
		functionDecl := function.(*ast.FunctionDecl)
		return functionDecl.RetType, nil
	case *ast.VoidExpr:
		// TODO(errors)
		if !expectedType.IsVoid() {
			return nil, fmt.Errorf("expected return type to be '%s'", expectedType)
		}
		return expectedType, nil
	case *ast.UnaryExpr:
		switch expression.Op {
		case token.MINUS:
			unaryExprType, err := sema.inferExprTypeWithContext(expression.Value, expectedType, scope)
			// TODO(errors)
			if err != nil {
				return nil, err
			}
			if !unaryExprType.IsNumeric() {
				return nil, fmt.Errorf("can't use - operator on a non-numeric value")
			}
			return unaryExprType, nil
		default:
			// TODO(errors)
			log.Fatalf("unimplemented unary expr operator: %s", reflect.TypeOf(expression.Op))
		}
	default:
		// TODO(errors)
		log.Fatalf("unimplemented on getExprType: %s", reflect.TypeOf(expression))
	}

	// NOTE: this line should be unreachable
	log.Fatalf("unreachable line at getExprTy")
	return nil, nil
}

// Useful for testing
func inferExprTypeWithoutContext(
	input, filename string,
	scope *ast.Scope,
) (ast.Expr, ast.ExprType, error) {
	collector := diagnostics.New()

	expr, err := parser.ParseExprFrom(input, filename)
	if err != nil {
		return nil, nil, err
	}

	analyzer := New(collector)
	exprType, _, err := analyzer.inferExprTypeWithoutContext(expr, scope)
	if err != nil {
		return nil, nil, err
	}
	return expr, exprType, nil
}

// Useful for testing
func inferExprTypeWithContext(
	input, filename string,
	ty ast.ExprType,
	scope *ast.Scope,
) (ast.ExprType, error) {
	collector := diagnostics.New()

	expr, err := parser.ParseExprFrom(input, filename)
	if err != nil {
		return nil, err
	}

	analyzer := New(collector)
	exprType, err := analyzer.inferExprTypeWithContext(expr, ty, scope)
	if err != nil {
		return nil, err
	}
	return exprType, nil
}

func (sema *sema) inferExprTypeWithoutContext(
	expr ast.Expr,
	scope *ast.Scope,
) (ast.ExprType, bool, error) {
	switch expression := expr.(type) {
	case *ast.LiteralExpr:
		switch ty := expression.Type.(type) {
		case *ast.BasicType:
			// TODO: what if I want a defined string type on my programming
			// language? I don't think *i8 is correct type for this
			switch ty.Kind {
			case token.STRING_LITERAL:
				finalTy := &ast.PointerType{Type: &ast.BasicType{Kind: token.U8_TYPE}}
				expression.Type = finalTy
				return finalTy, false, nil
			case token.INTEGER_LITERAL:
				ty, err := sema.inferIntegerType(expression.Value)
				if err != nil {
					return nil, false, err
				}
				expression.Type = ty
				return ty, false, nil
			case token.TRUE_BOOL_LITERAL:
				finalTy := &ast.BasicType{Kind: token.BOOL_TYPE}
				expression.Type = finalTy
				expression.Value = []byte("1")
				return finalTy, false, nil
			case token.FALSE_BOOL_LITERAL:
				finalTy := &ast.BasicType{Kind: token.BOOL_TYPE}
				expression.Type = finalTy
				expression.Value = []byte("0")
				return finalTy, false, nil
			default:
				// TODO(errors)
				log.Fatalf("unimplemented literal expr: %s", expression)
			}
		default:
			// TODO(errors)
			log.Fatalf("unimplemented literal expr type: %s", reflect.TypeOf(ty))
		}
	case *ast.IdExpr:
		variableName := expression.Name.Name()
		variable, err := scope.LookupAcrossScopes(variableName)
		// TODO(errors)
		if err != nil {
			return nil, false, err
		}

		switch node := variable.(type) {
		case *ast.VarStmt:
			return node.Type, true, nil
		case *ast.Field:
			return node.Type, true, nil
		default:
			return nil, false, fmt.Errorf("symbol '%s' is not a variable", node)
		}
	case *ast.UnaryExpr:
		switch expression.Op {
		case token.MINUS:
			switch unaryExpr := expression.Value.(type) {
			case *ast.LiteralExpr:
				switch unaryTy := unaryExpr.Type.(type) {
				case *ast.BasicType:
					if unaryTy.Kind == token.INTEGER_LITERAL {
						integerType, err := sema.inferIntegerType(unaryExpr.Value)
						// TODO(errors)
						if err != nil {
							return nil, false, err
						}
						unaryExpr.Type = integerType
						return integerType, false, nil
					}
				default:
					// TODO(errors)
					log.Fatalf("unimplemented unary expr type: %s", unaryExpr)
				}
			}
		case token.NOT:
			unaryExpr, foundContext, err := sema.inferExprTypeWithoutContext(expression.Value, scope)
			// TODO(errors)
			if err != nil {
				return nil, false, err
			}
			if !unaryExpr.IsBoolean() {
				return nil, false, fmt.Errorf("expected boolean expression on not unary expression")
			}
			return unaryExpr, foundContext, nil
		}
		// TODO(errors)
		log.Fatalf("unimplemented unary expression on sema")
	case *ast.BinaryExpr:
		ty, foundContext, err := sema.inferBinaryExprTypeWithoutContext(expression, scope)
		// TODO(errors)
		if err != nil {
			return nil, false, err
		}
		return ty, foundContext, nil
	case *ast.FunctionCall:
		err := sema.checkFunctionCall(expression, scope)
		// TODO(errors)
		if err != nil {
			return nil, false, err
		}
		// At this point, function should exists!
		function, err := scope.LookupAcrossScopes(expression.Name.Name())
		if err != nil {
			// TODO(errors)
			log.Fatalf("panic: at this point, function '%s' should exists in current block", expression.Name)
		}
		functionDecl := function.(*ast.FunctionDecl)
		return functionDecl.RetType, true, nil
	case *ast.FieldAccess:
		proto, err := sema.checkFieldAccessExpr(expression, scope)
		return proto.RetType, true, err
	default:
		// TODO(errors)
		log.Fatalf("unimplemented expression on sema: %s", reflect.TypeOf(expression))
	}
	// TODO(errors)
	// NOTE: this should be unreachable
	log.Fatalf("UNREACHABLE - inferExprType: %s", reflect.TypeOf(expr))
	return nil, false, nil
}

// TODO: refactor this algorithm
func (sema *sema) inferBinaryExprTypeWithoutContext(
	expression *ast.BinaryExpr,
	scope *ast.Scope,
) (ast.ExprType, bool, error) {
	lhsType, lhsFoundContext, err := sema.inferExprTypeWithoutContext(expression.Left, scope)
	// TODO(errors)
	if err != nil {
		return nil, false, err
	}

	rhsType, rhsFoundContext, err := sema.inferExprTypeWithoutContext(expression.Right, scope)
	// TODO(errors)
	if err != nil {
		return nil, false, err
	}

	if lhsFoundContext && !rhsFoundContext {
		rhsTypeWithContext, err := sema.inferExprTypeWithContext(expression.Right, lhsType, scope)
		// TODO(errors)
		if err != nil {
			log.Fatal(err)
		}
		rhsType = rhsTypeWithContext
	}
	if !lhsFoundContext && rhsFoundContext {
		lhsTypeWithContext, err := sema.inferExprTypeWithContext(expression.Left, rhsType, scope)
		// TODO(errors)
		if err != nil {
			log.Fatal(err)
		}
		lhsType = lhsTypeWithContext
	}

	// TODO(errors)
	if !reflect.DeepEqual(lhsType, rhsType) {
		log.Fatalf("mismatched types: %s %s %s", lhsType, expression.Op, rhsType)
	}

	switch expression.Op {
	case token.PLUS, token.MINUS:
		if lhsType.IsNumeric() && rhsType.IsNumeric() {
			return lhsType, lhsFoundContext || rhsFoundContext, nil
		}
	default:
		if _, ok := token.LOGICAL_OP[expression.Op]; ok {
			return &ast.BasicType{Kind: token.BOOL_TYPE}, lhsFoundContext || rhsFoundContext, nil
		}
	}
	// TODO(errors)
	log.Fatalf("UNREACHABLE - inferBinaryExprType")
	return nil, false, nil
}

func (sema *sema) inferBinaryExprTypeWithContext(
	expression *ast.BinaryExpr,
	expectedType ast.ExprType,
	scope *ast.Scope,
) (ast.ExprType, error) {
	lhsType, err := sema.inferExprTypeWithContext(expression.Left, expectedType, scope)
	if err != nil {
		return nil, err
	}

	rhsType, err := sema.inferExprTypeWithContext(expression.Right, expectedType, scope)
	if err != nil {
		return nil, err
	}

	if !reflect.DeepEqual(lhsType, rhsType) {
		return nil, fmt.Errorf("mismatched types: %s %s %s", lhsType, expression.Op, rhsType)
	}
	return lhsType, nil
}

func (sema *sema) inferIntegerType(value []byte) (ast.ExprType, error) {
	integerType := token.INT_TYPE
	base := 10
	intSize := strconv.IntSize

	_, err := strconv.ParseUint(string(value), base, intSize)
	if err == nil {
		return &ast.BasicType{Kind: integerType}, nil
	}

	return nil, fmt.Errorf("can't parse integer literal: %s", value)
}

func (sema *sema) checkFieldAccessExpr(
	fieldAccess *ast.FieldAccess,
	currentScope *ast.Scope,
) (*ast.Proto, error) {
	if !fieldAccess.Left.IsId() {
		return nil, fmt.Errorf("invalid expression on field accessing: %s", fieldAccess.Left)
	}

	idExpr := fieldAccess.Left.(*ast.IdExpr)
	id := idExpr.Name.Name()

	symbol, err := currentScope.LookupAcrossScopes(id)
	if err != nil {
		if err == ast.ERR_SYMBOL_NOT_FOUND_ON_SCOPE {
			pos := idExpr.Name.Pos
			symbolNotDefined := diagnostics.Diag{
				Message: fmt.Sprintf(
					"%s:%d:%d: '%s' not defined on scope",
					pos.Filename,
					pos.Line,
					pos.Column,
					id,
				),
			}
			sema.collector.ReportAndSave(symbolNotDefined)
			return nil, diagnostics.COMPILER_ERROR_FOUND
		}
		return nil, err
	}

	switch sym := symbol.(type) {
	case *ast.ExternDecl:
		switch right := fieldAccess.Right.(type) {
		case *ast.FunctionCall:
			err := sema.checkPrototypeCall(right, currentScope, sym)
			if err != nil {
				return nil, err
			}

			proto, err := sym.Scope.LookupCurrentScope(right.Name.Name())
			if err != nil {
				return nil, err
			}
			return proto.(*ast.Proto), nil
		default:
			// TODO(errors)
			return nil, fmt.Errorf("invalid expression %s when accessing field", right)
		}
	}
	return nil, nil
}

func (sema *sema) checkPrototypeCall(
	prototypeCall *ast.FunctionCall,
	callScope *ast.Scope,
	extern *ast.ExternDecl,
) error {
	prototype, err := extern.Scope.LookupCurrentScope(prototypeCall.Name.Name())
	if err != nil {
		pos := prototypeCall.Name.Pos
		prototypeNotFound := diagnostics.Diag{
			Message: fmt.Sprintf(
				"%s:%d:%d: function '%s' not declared on extern '%s'",
				pos.Filename,
				pos.Line,
				pos.Column,
				prototypeCall.Name.Name(),
				extern.Name.Name(),
			),
		}
		sema.collector.ReportAndSave(prototypeNotFound)
		return diagnostics.COMPILER_ERROR_FOUND
	}

	if proto, ok := prototype.(*ast.Proto); ok {
		if proto.Params.IsVariadic {
			minimumNumberOfArgs := len(proto.Params.Fields)
			// TODO(errors)
			if len(prototypeCall.Args) < minimumNumberOfArgs {
				log.Fatalf("the number of args is lesser than the minimum number of parameters")
			}
			for i := range minimumNumberOfArgs {
				paramType := proto.Params.Fields[i].Type
				argType, err := sema.inferExprTypeWithContext(
					prototypeCall.Args[i],
					paramType,
					callScope,
				)
				// TODO(errors)
				if err != nil {
					return err
				}
				// TODO(errors)
				if !reflect.DeepEqual(argType, paramType) {
					log.Fatalf(
						"mismatched argument type on prototype '%s', expected %s, but got %s",
						proto.Name,
						paramType,
						argType,
					)
				}
			}
		} else {
			if len(prototypeCall.Args) != len(proto.Params.Fields) {
				// TODO(errors)
				log.Fatalf("expected %d arguments, but got %d", len(proto.Params.Fields), len(prototypeCall.Args))
			}
			for i := range len(prototypeCall.Args) {
				paramType := proto.Params.Fields[i].Type
				argType, err := sema.inferExprTypeWithContext(prototypeCall.Args[i], paramType, callScope)
				// TODO(errors)
				if err != nil {
					return err
				}
				// TODO(errors)
				if !reflect.DeepEqual(argType, paramType) {
					log.Fatalf("mismatched argument type on function '%s', expected %s, but got %s", proto.Name, paramType, argType)
				}
			}
		}
	} else {
		// TODO(errors)
		log.Fatalf("%s needs to be a prototype", proto)
	}
	return nil
}

func (sema *sema) checkForLoop(
	forLoop *ast.ForLoop,
	scope *ast.Scope,
	returnTy ast.ExprType,
) error {
	err := sema.checkStmt(forLoop.Init, scope, returnTy)
	if err != nil {
		return err
	}

	err = sema.checkIfExpr(forLoop.Cond, scope)
	if err != nil {
		return err
	}

	err = sema.checkStmt(forLoop.Update, scope, returnTy)
	if err != nil {
		return err
	}

	err = sema.checkBlock(forLoop.Block, returnTy, scope)
	return err
}

// TODO: need tests for it
func (sema *sema) checkWhileLoop(
	whileLoop *ast.WhileLoop,
	scope *ast.Scope,
	returnTy ast.ExprType,
) error {
	err := sema.checkIfExpr(whileLoop.Cond, scope)
	if err != nil {
		return err
	}
	err = sema.checkBlock(whileLoop.Block, returnTy, scope)
	return err
}
