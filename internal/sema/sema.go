package sema

import (
	"fmt"
	"log"
	"reflect"
	"strconv"

	"github.com/HicaroD/Telia/internal/ast"
	"github.com/HicaroD/Telia/internal/diagnostics"
	"github.com/HicaroD/Telia/internal/lexer"
	"github.com/HicaroD/Telia/internal/lexer/token"
	"github.com/HicaroD/Telia/internal/parser"
)

type sema struct {
	mainPackageFound bool
	collector        *diagnostics.Collector
}

func New(collector *diagnostics.Collector) *sema {
	return &sema{false, collector}
}

func (s *sema) Check(program *ast.Program) error {
	return s.checkPackage(program.Root)
}

func (s *sema) checkPackage(pkg *ast.Package) error {
	// TODO(erros)
	// TODO: maybe restricting "main" as name for a package is not the way to go
	if pkg.Name == "main" {
		return fmt.Errorf("package name is not allowed to be 'main'")
	}

	requiresMain := false
	hasMainMethod := false

	for _, file := range pkg.Files {
		if file.PkgName == "main" {
			if s.mainPackageFound {
				// TODO(errors)
				return fmt.Errorf("error: main package already defined somewhere else")
			}
			s.mainPackageFound = true
			requiresMain = true
		}

		if requiresMain && file.PkgName != "main" {
			// TODO(errors)
			return fmt.Errorf("error: expected package name to be 'main'\n")
		}
		if !requiresMain && file.PkgName != pkg.Name {
			return fmt.Errorf("error: expected package name to be '%s'\n", pkg.Name)
		}

		fileHasMain, err := s.checkFile(file)
		if err != nil {
			return err
		}
		hasMainMethod = hasMainMethod || fileHasMain
	}

	// TODO(errors)
	if requiresMain && !hasMainMethod {
		return fmt.Errorf("main package requires a 'main' method as entrypoint\n")
	}

	// TODO(errors)
	if !requiresMain && hasMainMethod {
		return fmt.Errorf("'main' method not allowed on non-main packages\n")
	}

	for _, innerPackage := range pkg.Packages {
		err := s.checkPackage(innerPackage)
		if err != nil {
			return err
		}
	}

	return nil
}

func (s *sema) checkFile(file *ast.File) (bool, error) {
	foundMain := false

	for _, node := range file.Body {
		switch node.Kind {
		case ast.KIND_FN_DECL:
			fnDecl := node.Node.(*ast.FnDecl)
			if !foundMain {
				foundMain = fnDecl.Name.Name() == "main"
			}
			err := s.checkFnDecl(fnDecl)
			if err != nil {
				return false, err
			}
		case ast.KIND_EXTERN_DECL:
			err := s.checkExternDecl(node.Node.(*ast.ExternDecl))
			if err != nil {
				return false, err
			}
		case ast.KIND_PKG_DECL:
			err := s.checkUseDecl()
			if err != nil {
				return false, err
			}
		case ast.KIND_USE_DECL:
			continue
		default:
			return false, fmt.Errorf("unimplemented ast node for sema: %s\n", reflect.TypeOf(node.Node))
		}
	}

	return foundMain, nil
}

func (sema *sema) checkUseDecl() error {
	return nil
}

func (sema *sema) checkFnDecl(function *ast.FnDecl) error {
	err := sema.checkBlock(function.Block, function.RetType, function.Scope)
	return err
}

var VALID_CALLING_CONVENTIONS []string = []string{
	"c", "fast", "cold",
}

func (sema *sema) checkExternAttributes(attributes *ast.ExternAttrs) error {
	ccFound := false
	for _, cc := range VALID_CALLING_CONVENTIONS {
		if cc == attributes.DefaultCallingConvention {
			ccFound = true
			break
		}
	}
	if !ccFound {
		return fmt.Errorf("invalid calling convention: '%s'\n", attributes.DefaultCallingConvention)
	}
	return nil
}

func (sema *sema) checkExternDecl(extern *ast.ExternDecl) error {
	if extern.Attributes != nil {
		if err := sema.checkExternAttributes(extern.Attributes); err != nil {
			return err
		}
	}

	for _, proto := range extern.Prototypes {
		err := sema.checkExternPrototype(extern, proto)
		if err != nil {
			return err
		}
	}

	return nil
}

var VALID_FUNCTION_LINKAGES []string = []string{
	"external", "internal", "weak", "link_once",
}

func (sema *sema) checkExternPrototype(extern *ast.ExternDecl, proto *ast.Proto) error {
	if proto.Attributes != nil {
		linkageFound := false
		if proto.Attributes.Linkage != "" {
			for _, l := range VALID_FUNCTION_LINKAGES {
				if l == proto.Attributes.Linkage {
					linkageFound = true
					break
				}
			}
			// if not found, set a default linkage
			if !linkageFound {
				return fmt.Errorf("invalid linkage type: %s\n", proto.Attributes.Linkage)
			}
		}
	}

	symbols := make(map[string]bool, len(proto.Params.Fields))
	for _, param := range proto.Params.Fields {
		if _, found := symbols[param.Name.Name()]; found {
			// TODO(errors): add proper error here + tests
			return fmt.Errorf("redeclaration of '%s' parameter on '%s' prototype at extern declaration '%s'\n", param.Name.Name(), proto.Name.Name(), extern.Name.Name())
		}
		symbols[param.Name.Name()] = true
	}

	return nil
}

func (sema *sema) checkBlock(
	block *ast.BlockStmt,
	returnTy *ast.MyExprType,
	scope *ast.Scope,
) error {
	for _, statement := range block.Statements {
		err := sema.checkStmt(statement, scope, returnTy)
		if err != nil {
			return err
		}
	}
	return nil
}

func (sema *sema) checkStmt(
	stmt *ast.MyNode,
	scope *ast.Scope,
	returnTy *ast.MyExprType,
) error {
	switch stmt.Kind {
	case ast.KIND_FN_CALL:
		err := sema.checkFunctionCall(stmt.Node.(*ast.FnCall), scope)
		return err
	case ast.KIND_MULTI_VAR_STMT, ast.KIND_VAR_STMT:
		err := sema.checkVarDecl(stmt, scope)
		return err
	case ast.KIND_COND_STMT:
		err := sema.checkCondStmt(stmt.Node.(*ast.CondStmt), returnTy, scope)
		return err
	case ast.KIND_RETURN_STMT:
		returnStmt := stmt.Node.(*ast.ReturnStmt)
		_, err := sema.inferExprTypeWithContext(returnStmt.Value, returnTy, scope)
		return err
	case ast.KIND_FOR_LOOP_STMT:
		err := sema.checkForLoop(stmt.Node.(*ast.ForLoop), returnTy, scope)
		return err
	case ast.KIND_WHILE_LOOP_STMT:
		err := sema.checkWhileLoop(stmt.Node.(*ast.WhileLoop), scope, returnTy)
		return err
	default:
		return fmt.Errorf("error: unimplemented statement on sema: %s\n", stmt)
	}
}

func (sema *sema) checkVarDecl(
	variable *ast.MyNode,
	currentScope *ast.Scope,
) error {
	switch variable.Kind {
	case ast.KIND_MULTI_VAR_STMT:
		err := sema.checkMultiVar(variable.Node.(*ast.MultiVarStmt), currentScope)
		return err
	case ast.KIND_VAR_STMT:
		err := sema.checkVar(variable.Node.(*ast.VarStmt), currentScope)
		return err
	}
	return nil
}

func (sema *sema) checkMultiVar(
	multi *ast.MultiVarStmt,
	currentScope *ast.Scope,
) error {
	for _, variable := range multi.Variables {
		err := sema.checkVar(variable.Node.(*ast.VarStmt), currentScope)
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
func checkVarDeclFrom(input, filename string) (*ast.VarStmt, error) {
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

	return varStmt.Node.(*ast.VarStmt), nil
}

func (sema *sema) checkCondStmt(
	condStmt *ast.CondStmt,
	returnTy *ast.MyExprType,
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
	fnCall *ast.FnCall,
	currentScope *ast.Scope,
) error {
	symbol, err := currentScope.LookupAcrossScopes(fnCall.Name.Name())
	if err != nil {
		if err == ast.ERR_SYMBOL_NOT_FOUND_ON_SCOPE {
			pos := fnCall.Name.Pos
			functionNotDefined := diagnostics.Diag{
				Message: fmt.Sprintf(
					"%s:%d:%d: function '%s' not defined on scope",
					pos.Filename,
					pos.Line,
					pos.Column,
					fnCall.Name.Name(),
				),
			}
			sema.collector.ReportAndSave(functionNotDefined)
			return diagnostics.COMPILER_ERROR_FOUND
		}
		return err
	}

	if symbol.Node == ast.KIND_FN_DECL {
		pos := fnCall.Name.Pos
		notCallable := diagnostics.Diag{
			Message: fmt.Sprintf(
				"%s:%d:%d: '%s' is not callable",
				pos.Filename,
				pos.Line,
				pos.Column,
				fnCall.Name.Name(),
			),
		}
		sema.collector.ReportAndSave(notCallable)
		return diagnostics.COMPILER_ERROR_FOUND
	}

	fnDecl := symbol.Node.(*ast.FnDecl)

	if len(fnCall.Args) != len(fnDecl.Params.Fields) {
		pos := fnCall.Name.Pos
		// TODO(errors): show which arguments were passed and which types we
		// were expecting
		notEnoughArguments := diagnostics.Diag{
			Message: fmt.Sprintf(
				"%s:%d:%d: not enough arguments in call to '%s'",
				pos.Filename,
				pos.Line,
				pos.Column,
				fnCall.Name.Name(),
			),
		}
		sema.collector.ReportAndSave(notEnoughArguments)
		return diagnostics.COMPILER_ERROR_FOUND
	}

	for i, arg := range fnCall.Args {
		paramType := fnDecl.Params.Fields[i].Type
		argType, err := sema.inferExprTypeWithContext(arg, paramType, currentScope)
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

func (sema *sema) checkIfExpr(expr *ast.MyNode, scope *ast.Scope) error {
	inferedExprType, _, err := sema.inferExprTypeWithoutContext(expr, scope)
	// TODO(errors)
	if err != nil {
		return err
	}

	// TODO(erros)
	if !inferedExprType.IsBoolean() {
		return fmt.Errorf("invalid non-boolean condition on if statement: %s", inferedExprType)
	}

	return nil
}

func (s *sema) inferExprTypeWithContext(
	expr *ast.MyNode,
	expectedType *ast.MyExprType,
	scope *ast.Scope,
) (*ast.MyExprType, error) {
	switch expr.Kind {
	case ast.KIND_LITERAl_EXPR:
		return s.inferLiteralExprTypeWithContext(expr.Node.(*ast.LiteralExpr), expectedType, scope)
	case ast.KIND_ID_EXPR:
		return s.inferIdExprTypeWithContext(expr.Node.(*ast.IdExpr), expectedType, scope)
	case ast.KIND_BINARY_EXPR:
		return s.inferBinaryExprTypeWithContext(expr.Node.(*ast.BinaryExpr), expectedType, scope)
	case ast.KIND_UNARY_EXPR:
		return s.inferUnaryExprTypeWithContext(expr.Node.(*ast.UnaryExpr), expectedType, scope)
	case ast.KIND_FN_CALL:
		return s.inferFnCallExprTypeWithContext(expr.Node.(*ast.FnCall), expectedType, scope)
	case ast.KIND_VOID_EXPR:
		return s.inferVoidExprTypeWithContext(expr.Node.(*ast.VoidExpr), expectedType, scope)
	default:
		log.Fatalf("unimplemented expression: %s\n", expr.Kind)
		return nil, nil
	}
}

func (sema *sema) inferLiteralExprTypeWithContext(
	literal *ast.LiteralExpr,
	expectedType *ast.MyExprType,
	scope *ast.Scope,
) (*ast.MyExprType, error) {
	// TODO(errors)
	if literal.Type.Kind != ast.EXPR_TYPE_BASIC && expectedType.Kind != ast.EXPR_TYPE_BASIC {
		return nil, fmt.Errorf("error: type mismatch for literal expression")
	}

	expectedBasicType := expectedType.T.(*ast.BasicType)
	actualBasicType := literal.Type.T.(*ast.BasicType)

	_, err := sema.inferBasicExprTypeWithContext(actualBasicType, expectedBasicType)
	if err != nil {
		return nil, err
	}

	return expectedType, nil
}

func (sema *sema) inferIdExprTypeWithContext(
	id *ast.IdExpr,
	expectedType *ast.MyExprType,
	scope *ast.Scope,
) (*ast.MyExprType, error) {
	symbol, err := scope.LookupAcrossScopes(id.Name.Name())
	// TODO(errors)
	if err != nil {
		return nil, err
	}

	var ty *ast.MyExprType

	switch symbol.Kind {
	case ast.KIND_VAR_STMT:
		ty = symbol.Node.(*ast.VarStmt).Type
	case ast.KIND_FIELD:
		ty = symbol.Node.(*ast.Field).Type
	default:
		// TODO(errors)
		return nil, fmt.Errorf("expected to be a variable or parameter, but got %s", symbol.Kind)
	}

	switch ty.Kind {
	case ast.EXPR_TYPE_BASIC:
		actualBasicType := ty.T.(*ast.BasicType)
		// TODO(errors)
		if expectedType.Kind != ast.EXPR_TYPE_BASIC {
			return nil, fmt.Errorf("error: type mismatch for id expression")
		}
		expectedBasicType := expectedType.T.(*ast.BasicType)

		_, err := sema.inferBasicExprTypeWithContext(actualBasicType, expectedBasicType)
		if err != nil {
			return nil, err
		}

		return expectedType, nil
	default:
		// TODO(errors)
		panic("unimplemented infer id expr type with context")
	}
}

func (sema *sema) inferBinaryExprTypeWithContext(
	binary *ast.BinaryExpr,
	expectedType *ast.MyExprType,
	scope *ast.Scope,
) (*ast.MyExprType, error) {
	lhsType, err := sema.inferExprTypeWithContext(binary.Left, expectedType, scope)
	if err != nil {
		return nil, err
	}

	rhsType, err := sema.inferExprTypeWithContext(binary.Right, expectedType, scope)
	if err != nil {
		return nil, err
	}

	// NOTE: figure out a way to not use reflect.DeepEqual
	if !reflect.DeepEqual(lhsType, rhsType) {
		return nil, fmt.Errorf("mismatched types: %s %s %s", lhsType, binary.Op, rhsType)
	}
	return lhsType, nil
}

func (sema *sema) inferUnaryExprTypeWithContext(
	unary *ast.UnaryExpr,
	expectedType *ast.MyExprType,
	scope *ast.Scope,
) (*ast.MyExprType, error) {
	switch unary.Op {
	case token.MINUS:
		unaryExprType, err := sema.inferExprTypeWithContext(unary.Value, expectedType, scope)
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
		log.Fatalf("unimplemented unary expr operator: %s", unary.Op)
		return nil, nil
	}
}

func (sema *sema) inferFnCallExprTypeWithContext(
	fnCall *ast.FnCall,
	expectedType *ast.MyExprType,
	scope *ast.Scope,
) (*ast.MyExprType, error) {
	err := sema.checkFunctionCall(fnCall, scope)
	// TODO(errors)
	if err != nil {
		return nil, err
	}
	symbol, _ := scope.LookupAcrossScopes(fnCall.Name.Name())
	fnDecl := symbol.Node.(*ast.FnDecl)
	return fnDecl.RetType, nil
}

func (sema *sema) inferVoidExprTypeWithContext(
	void *ast.VoidExpr,
	expectedType *ast.MyExprType,
	scope *ast.Scope,
) (*ast.MyExprType, error) {
	// TODO(errors)
	if !expectedType.IsVoid() {
		return nil, fmt.Errorf("expected return type to be '%s'", expectedType)
	}
	return expectedType, nil
}

func (sema *sema) inferBasicExprTypeWithContext(
	actual *ast.BasicType,
	expected *ast.BasicType,
) (*ast.BasicType, error) {

	switch actual.Kind {
	case token.STRING_LITERAL, token.UNTYPED_STRING:
		if !expected.IsAnyStringType() {
			return nil, fmt.Errorf("error: expected %s type, got string\n", expected.Kind)
		}
		actual.Kind = expected.Kind
		return expected, nil
	case token.INTEGER_LITERAL, token.UNTYPED_INT:
		if !expected.IsIntegerType() {
			return nil, fmt.Errorf("error: expected %s type, got integer\n", expected.Kind)
		}
		actual.Kind = expected.Kind
		return expected, nil
	case token.TRUE_BOOL_LITERAL, token.FALSE_BOOL_LITERAL:
		if expected.Kind != token.BOOL_TYPE {
			return nil, fmt.Errorf("expected %s type, got boolean\n", expected.Kind)
		}
		actual.Kind = token.BOOL_TYPE
		return expected, nil
	default:
		// TODO(errors)
		return nil, fmt.Errorf("error: type mismatch for literal expression")
	}
}

// Useful for testing
func inferExprTypeWithoutContext(
	input, filename string,
	scope *ast.Scope,
) (*ast.MyNode, *ast.MyExprType, error) {
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
	ty *ast.MyExprType,
	scope *ast.Scope,
) (*ast.MyExprType, error) {
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
	expr *ast.MyNode,
	scope *ast.Scope,
) (*ast.MyExprType, bool, error) {
	switch expression := expr.(type) {
	case *ast.LiteralExpr:
		switch ty := expression.Type.(type) {
		case *ast.BasicType:
			switch ty.Kind {
			case token.STRING_LITERAL:
				finalTy := &ast.BasicType{Kind: token.UNTYPED_STRING}
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
	case *ast.FnCall:
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
			return nil, false, err
		}
		rhsType = rhsTypeWithContext
	}
	if !lhsFoundContext && rhsFoundContext {
		lhsTypeWithContext, err := sema.inferExprTypeWithContext(expression.Left, rhsType, scope)
		// TODO(errors)
		if err != nil {
			return nil, false, err
		}
		lhsType = lhsTypeWithContext
	}

	// TODO: get rid of reflect.DeepEqual for comparing types somehow
	// TODO(errors)
	if !reflect.DeepEqual(lhsType, rhsType) {
		return nil, false, fmt.Errorf("mismatched types: %s %s %s", lhsType, expression.Op, rhsType)
	}

	// TODO: it needs to be more flexible - automatically evaluate correct
	// operators
	switch expression.Op {
	case token.PLUS, token.MINUS, token.SLASH, token.STAR:
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

func (sema *sema) inferIntegerType(value []byte) (ast.ExprType, error) {
	integerType := token.UNTYPED_INT
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
		case *ast.FnCall:
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
	prototypeCall *ast.FnCall,
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
				if _, err := sema.inferExprTypeWithContext(prototypeCall.Args[i], paramType, callScope); err != nil {
					return err
				}
			}
		} else {
			if len(prototypeCall.Args) != len(proto.Params.Fields) {
				// TODO(errors)
				log.Fatalf("expected %d arguments, but got %d", len(proto.Params.Fields), len(prototypeCall.Args))
			}
			for i := range len(prototypeCall.Args) {
				paramType := proto.Params.Fields[i].Type
				if _, err := sema.inferExprTypeWithContext(prototypeCall.Args[i], paramType, callScope); err != nil {
					return err
				}
				fmt.Println(prototypeCall.Args[i])
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
	returnTy *ast.MyExprType,
	scope *ast.Scope,
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

func (sema *sema) checkWhileLoop(
	whileLoop *ast.WhileLoop,
	scope *ast.Scope,
	returnTy *ast.MyExprType,
) error {
	err := sema.checkIfExpr(whileLoop.Cond, scope)
	if err != nil {
		return err
	}
	err = sema.checkBlock(whileLoop.Block, returnTy, scope)
	return err
}
