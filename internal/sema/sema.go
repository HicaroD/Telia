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
			externDecl := node.Node.(*ast.ExternDecl)
			err := s.checkExternDecl(externDecl)
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
	if function.Attributes != nil {
		err := sema.checkFnAttributes(function.Attributes)
		if err != nil {
			return err
		}
	}
	err := sema.checkBlock(function.Block, function.RetType, function.Scope)
	return err
}

var VALID_CALLING_CONVENTIONS []string = []string{
	"c", "fast", "cold",
}

func (sema *sema) checkExternAttributes(attributes *ast.Attributes) error {
	if attributes.Linkage != "" {
		return fmt.Errorf("'linkage' is not a valid attribute for extern declaration\n")
	}

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

func (sema *sema) checkFnAttributes(attributes *ast.Attributes) error {
	linkageFound := false

	if attributes.Linkage != "" {
		for _, l := range VALID_FUNCTION_LINKAGES {
			if l == attributes.Linkage {
				linkageFound = true
				break
			}
		}

		if !linkageFound {
			return fmt.Errorf("invalid linkage type: %s\n", attributes.Linkage)
		}
	}

	return nil
}

func (sema *sema) checkExternPrototype(extern *ast.ExternDecl, proto *ast.Proto) error {
	if proto.Attributes != nil {
		err := sema.checkFnAttributes(proto.Attributes)
		if err != nil {
			return err
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
	returnTy *ast.ExprType,
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
	stmt *ast.Node,
	scope *ast.Scope,
	returnTy *ast.ExprType,
) error {
	switch stmt.Kind {
	case ast.KIND_FN_CALL:
		err := sema.checkFunctionCall(stmt.Node.(*ast.FnCall), scope)
		return err
	case ast.KIND_VAR_STMT:
		err := sema.checkVar(stmt.Node.(*ast.Var), scope)
		return err
	case ast.KIND_COND_STMT:
		err := sema.checkCondStmt(stmt.Node.(*ast.CondStmt), returnTy, scope)
		return err
	case ast.KIND_RETURN_STMT:
		returnStmt := stmt.Node.(*ast.ReturnStmt)
		_, err := sema.inferExprTypeWithContext(returnStmt.Value, returnTy, scope)
		return err
	case ast.KIND_FIELD_ACCESS:
		_, err := sema.checkFieldAccessExpr(stmt.Node.(*ast.FieldAccess), scope)
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

func (sema *sema) checkVar(variable *ast.Var, currentScope *ast.Scope) error {
	for _, currentVar := range variable.Names {
		if variable.IsDecl {
			_, err := currentScope.LookupCurrentScope(currentVar.Name.Name())
			if err == nil {
				return fmt.Errorf("'%s' already declared in the current block", currentVar.Name.Name())
			}

			n := new(ast.Node)
			n.Kind = ast.KIND_VAR_ID_STMT
			n.Node = currentVar

			err = currentScope.Insert(currentVar.Name.Name(), n)
			if err != nil {
				return err
			}
		} else {
			_, err := currentScope.LookupAcrossScopes(currentVar.Name.Name())
			// TODO(errors)
			if err != nil {
				if err == ast.ErrSymbolNotFoundOnScope {
					name := currentVar.Name.Name()
					pos := currentVar.Name.Pos
					return fmt.Errorf("%s '%s' not declared in the current block", pos, name)
				}
				return err
			}
		}
	}

	switch variable.Expr.Kind {
	case ast.KIND_TUPLE_EXPR:
		tuple := variable.Expr.Node.(*ast.TupleExpr)
		numExprs := sema.countExprsOnTuple(tuple)

		// TODO(errors)
		if len(variable.Names) != numExprs {
			return fmt.Errorf("%d != %d", len(tuple.Exprs), numExprs)
		}

		t := 0
		for _, expr := range tuple.Exprs {
			switch expr.Kind {
			case ast.KIND_TUPLE_EXPR:
				innerTupleExpr := expr.Node.(*ast.TupleExpr)
				for _, innerExpr := range innerTupleExpr.Exprs {
					err := sema.checkVarExpr(variable.Names[t], innerExpr, currentScope)
					if err != nil {
						return err
					}
					t++
				}
			default:
				err := sema.checkVarExpr(variable.Names[t], expr, currentScope)
				if err != nil {
					return err
				}
				t++
			}
		}
	default:
		if len(variable.Names) != 1 {
			return fmt.Errorf("more variables than expressions\n")
		}
		err := sema.checkVarExpr(variable.Names[0], variable.Expr, currentScope)
		if err != nil {
			return err
		}
	}

	return nil
}

func (sema *sema) checkVarExpr(variable *ast.VarId, expr *ast.Node, scope *ast.Scope) error {
	if variable.NeedsInference {
		// TODO(errors): need a test for it
		if variable.Type != nil {
			return fmt.Errorf(
				"needs inference, but variable already has a type: %s",
				variable.Type,
			)
		}
		exprType, _, err := sema.inferExprTypeWithoutContext(expr, scope)
		// TODO(errors)
		if err != nil {
			return err
		}
		variable.Type = exprType
	} else {
		// TODO(errors)
		if variable.Type == nil {
			log.Fatalf("variable does not have a type and it said it does not need inference")
		}
		exprTy, err := sema.inferExprTypeWithContext(expr, variable.Type, scope)
		// TODO(errors): Deal with type mismatch
		if err != nil {
			return err
		}
		if !reflect.DeepEqual(variable.Type, exprTy) {
			return fmt.Errorf("type mismatch on variable decl, expected %s, got %s", variable.Type, exprTy)
		}
	}
	return nil
}

func (sema *sema) countExprsOnTuple(tuple *ast.TupleExpr) int {
	counter := 0
	for _, expr := range tuple.Exprs {
		switch expr.Kind {
		case ast.KIND_TUPLE_EXPR:
			varTuple := expr.Node.(*ast.TupleExpr)
			counter += sema.countExprsOnTuple(varTuple)
		default:
			counter++
		}
	}
	return counter
}

// Useful for testing
func checkVarDeclFrom(input, filename string) (*ast.VarId, error) {
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
	err = sema.checkVar(varStmt.Node.(*ast.Var), scope)
	if err != nil {
		return nil, err
	}

	return varStmt.Node.(*ast.VarId), nil
}

func (sema *sema) checkCondStmt(
	condStmt *ast.CondStmt,
	returnTy *ast.ExprType,
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
		if err == ast.ErrSymbolNotFoundOnScope {
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

	if symbol.Kind != ast.KIND_FN_DECL {
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

func (sema *sema) checkIfExpr(expr *ast.Node, scope *ast.Scope) error {
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
	expr *ast.Node,
	expectedType *ast.ExprType,
	scope *ast.Scope,
) (*ast.ExprType, error) {
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
		return s.inferVoidExprTypeWithContext(expectedType, scope)
	// TODO: infer tuple expr with context
	default:
		log.Fatalf("unimplemented expression: %s\n", expr.Kind)
		return nil, nil
	}
}

func (sema *sema) inferLiteralExprTypeWithContext(
	literal *ast.LiteralExpr,
	expectedType *ast.ExprType,
	scope *ast.Scope,
) (*ast.ExprType, error) {
	// TODO(errors)
	if literal.Type.Kind != ast.EXPR_TYPE_BASIC || expectedType.Kind != ast.EXPR_TYPE_BASIC {
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
	expectedType *ast.ExprType,
	scope *ast.Scope,
) (*ast.ExprType, error) {
	symbol, err := scope.LookupAcrossScopes(id.Name.Name())
	// TODO(errors)
	if err != nil {
		return nil, err
	}

	var ty *ast.ExprType

	switch symbol.Kind {
	case ast.KIND_VAR_ID_STMT:
		ty = symbol.Node.(*ast.VarId).Type
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
	expectedType *ast.ExprType,
	scope *ast.Scope,
) (*ast.ExprType, error) {
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
	expectedType *ast.ExprType,
	scope *ast.Scope,
) (*ast.ExprType, error) {
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
	expectedType *ast.ExprType,
	scope *ast.Scope,
) (*ast.ExprType, error) {
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
	expectedType *ast.ExprType,
	scope *ast.Scope,
) (*ast.ExprType, error) {
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
	var untyped token.Kind

	switch {
	case expected.Kind == token.BOOL_TYPE:
		untyped = token.UNTYPED_BOOL
	case expected.IsAnyStringType():
		untyped = token.UNTYPED_STRING
	case expected.IsIntegerType():
		untyped = token.UNTYPED_INT
	default:
		return nil, fmt.Errorf("unimplemented type: %s %s", actual.String(), expected.String())
	}

	if actual.Kind == untyped {
		actual.Kind = expected.Kind
	}

	// TODO: in the case of integer, i need to make sure that the value really fits the expected size in bits
	if actual.Kind != expected.Kind {
		return nil, fmt.Errorf("type mismatch - expected %s, got %s", expected.String(), actual.String())
	}

	return expected, nil
}

// Useful for testing
func inferExprTypeWithoutContext(
	input, filename string,
	scope *ast.Scope,
) (*ast.Node, *ast.ExprType, error) {
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
	ty *ast.ExprType,
	scope *ast.Scope,
) (*ast.ExprType, error) {
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
	expr *ast.Node,
	scope *ast.Scope,
) (*ast.ExprType, bool, error) {
	switch expr.Kind {
	case ast.KIND_LITERAl_EXPR:
		return sema.inferLiteralExprTypeWithoutContext(expr.Node.(*ast.LiteralExpr), scope)
	case ast.KIND_ID_EXPR:
		return sema.inferIdExprTypeWithoutContext(expr.Node.(*ast.IdExpr), scope)
	case ast.KIND_UNARY_EXPR:
		return sema.inferUnaryExprTypeWithoutContext(expr.Node.(*ast.UnaryExpr), scope)
	case ast.KIND_BINARY_EXPR:
		return sema.inferBinaryExprTypeWithoutContext(expr.Node.(*ast.BinaryExpr), scope)
	case ast.KIND_FN_CALL:
		return sema.inferFnCallExprTypeWithoutContext(expr.Node.(*ast.FnCall), scope)
	case ast.KIND_FIELD_ACCESS:
		return sema.inferFieldAccessExprTypeWithoutContext(expr.Node.(*ast.FieldAccess), scope)
	// TODO: infer tuple expr without context
	default:
		log.Fatalf("unimplemented expression: %s\n", expr.Kind)
		return nil, false, nil
	}
}

func (sema *sema) inferLiteralExprTypeWithoutContext(
	literal *ast.LiteralExpr,
	scope *ast.Scope,
) (*ast.ExprType, bool, error) {
	// TODO(errors)
	if literal.Type.Kind != ast.EXPR_TYPE_BASIC {
		return nil, false, fmt.Errorf("type mismatch for literal expression without context")
	}

	actualBasicType := literal.Type.T.(*ast.BasicType)
	switch actualBasicType.Kind {
	case token.UNTYPED_STRING:
		actualBasicType.Kind = token.UNTYPED_STRING
		return literal.Type, false, nil
	case token.UNTYPED_INT:
		ty, err := sema.inferIntegerType(literal.Value)
		if err != nil {
			return nil, false, err
		}
		literal.Type = ty
		return literal.Type, false, nil
	case token.UNTYPED_BOOL:
		actualBasicType.Kind = token.BOOL_TYPE
		return literal.Type, true, nil
	default:
		// TODO(errors)
		log.Fatalf("unimplemented literal expression: %s\n", actualBasicType.Kind)
		return nil, false, nil
	}
}

func (sema *sema) inferIdExprTypeWithoutContext(
	id *ast.IdExpr,
	scope *ast.Scope,
) (*ast.ExprType, bool, error) {
	variableName := id.Name.Name()
	variable, err := scope.LookupAcrossScopes(variableName)
	// TODO(errors)
	if err != nil {
		return nil, false, err
	}

	switch variable.Kind {
	case ast.KIND_VAR_ID_STMT:
		ty := variable.Node.(*ast.VarId).Type
		return ty, !ty.IsUntyped(), nil
	case ast.KIND_FIELD:
		return variable.Node.(*ast.Field).Type, true, nil
	default:
		// TODO(errors)
		return nil, false, fmt.Errorf("'%s' is not a variable or parameter", variableName)
	}
}

func (sema *sema) inferUnaryExprTypeWithoutContext(
	unary *ast.UnaryExpr,
	scope *ast.Scope,
) (*ast.ExprType, bool, error) {
	switch unary.Op {
	case token.MINUS:
		switch unary.Value.Kind {
		case ast.KIND_LITERAl_EXPR:
			lit := unary.Value.Node.(*ast.LiteralExpr)
			if !lit.Type.IsNumeric() {
				// TODO(errors)
				return nil, false, fmt.Errorf("error: expected numeric type")
			}
			// TODO: this method should infer any numeric type
			numericTy, err := sema.inferIntegerType(lit.Value)
			if err != nil {
				return nil, false, fmt.Errorf("error: unable to infer number type")
			}

			lit.Type = numericTy
			return lit.Type, false, nil
		default:
			return nil, false, fmt.Errorf("invalid unary expression: %s\n", reflect.TypeOf(unary.Value))
		}
	case token.NOT:
		unaryExpr, foundContext, err := sema.inferExprTypeWithoutContext(unary.Value, scope)
		if err != nil {
			return nil, false, err
		}
		if !unaryExpr.IsBoolean() {
			return nil, false, fmt.Errorf("expected boolean expression on not unary expression")
		}
		return unaryExpr, foundContext, nil
	default:
		log.Fatalf("unimplemented unary operator: %s\n", unary.Op)
		return nil, false, nil
	}
}

func (sema *sema) inferBinaryExprTypeWithoutContext(
	binary *ast.BinaryExpr,
	scope *ast.Scope,
) (*ast.ExprType, bool, error) {
	lhsType, lhsFoundContext, err := sema.inferExprTypeWithoutContext(binary.Left, scope)
	// TODO(errors)
	if err != nil {
		return nil, false, err
	}

	rhsType, rhsFoundContext, err := sema.inferExprTypeWithoutContext(binary.Right, scope)
	// TODO(errors)
	if err != nil {
		return nil, false, err
	}

	if lhsFoundContext && !rhsFoundContext {
		rhsTypeWithContext, err := sema.inferExprTypeWithContext(binary.Right, lhsType, scope)
		// TODO(errors)
		if err != nil {
			return nil, false, err
		}
		rhsType = rhsTypeWithContext
	}
	if !lhsFoundContext && rhsFoundContext {
		lhsTypeWithContext, err := sema.inferExprTypeWithContext(binary.Left, rhsType, scope)
		// TODO(errors)
		if err != nil {
			return nil, false, err
		}
		lhsType = lhsTypeWithContext
	}

	// TODO: get rid of reflect.DeepEqual for comparing types somehow
	// TODO(errors)
	if !reflect.DeepEqual(lhsType, rhsType) {
		return nil, false, fmt.Errorf("mismatched types: %s %s %s", lhsType, binary.Op, rhsType)
	}

	// TODO: it needs to be more flexible - easily evaluate correct operators
	switch binary.Op {
	case token.PLUS, token.MINUS, token.SLASH, token.STAR:
		if lhsType.IsNumeric() && rhsType.IsNumeric() {
			return lhsType, lhsFoundContext || rhsFoundContext, nil
		}
	default:
		if binary.Op.IsLogicalOp() {
			t := ast.NewBasicType(token.BOOL_TYPE)
			return t, lhsFoundContext || rhsFoundContext, nil
		}
	}
	// TODO(errors)
	log.Fatalf("UNREACHABLE - inferBinaryExprType")
	return nil, false, nil
}

func (sema *sema) inferFnCallExprTypeWithoutContext(
	fnCall *ast.FnCall,
	scope *ast.Scope,
) (*ast.ExprType, bool, error) {
	err := sema.checkFunctionCall(fnCall, scope)
	// TODO(errors)
	if err != nil {
		return nil, false, err
	}
	// At this point, function should exists!
	function, _ := scope.LookupAcrossScopes(fnCall.Name.Name())
	functionDecl := function.Node.(*ast.FnDecl)
	return functionDecl.RetType, true, nil
}

func (sema *sema) inferFieldAccessExprTypeWithoutContext(
	fieldAccess *ast.FieldAccess,
	scope *ast.Scope,
) (*ast.ExprType, bool, error) {
	proto, err := sema.checkFieldAccessExpr(fieldAccess, scope)
	return proto.RetType, true, err
}

// TODO: properly infer the type if possible
func (sema *sema) inferIntegerType(value []byte) (*ast.ExprType, error) {
	integerType := token.UNTYPED_INT
	base := 10
	intSize := strconv.IntSize

	_, err := strconv.ParseUint(string(value), base, intSize)
	if err != nil {
		return nil, fmt.Errorf("can't parse integer literal: %s", value)
	}

	t := ast.NewBasicType(integerType)
	return t, nil
}

func (sema *sema) checkFieldAccessExpr(
	fieldAccess *ast.FieldAccess,
	currentScope *ast.Scope,
) (*ast.Proto, error) {
	if !fieldAccess.Left.IsId() {
		return nil, fmt.Errorf("invalid expression on field accessing: %s", fieldAccess.Left)
	}

	idExpr := fieldAccess.Left.Node.(*ast.IdExpr)
	id := idExpr.Name.Name()

	symbol, err := currentScope.LookupAcrossScopes(id)
	if err != nil {
		if err == ast.ErrSymbolNotFoundOnScope {
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

	switch symbol.Kind {
	case ast.KIND_EXTERN_DECL:
		extern := symbol.Node.(*ast.ExternDecl)
		switch fieldAccess.Right.Kind {
		case ast.KIND_FN_CALL:
			fnCall := fieldAccess.Right.Node.(*ast.FnCall)
			err := sema.checkPrototypeCall(fnCall, currentScope, extern)
			if err != nil {
				return nil, err
			}

			proto, err := extern.Scope.LookupCurrentScope(fnCall.Name.Name())
			if err != nil {
				return nil, err
			}

			return proto.Node.(*ast.Proto), nil
		default:
			return nil, fmt.Errorf("expected prototype function in the right side, not %s\n", fieldAccess.Right.Kind)
		}
	default:
		// TODO(errors)
		return nil, fmt.Errorf("invalid expression %s when accessing field", fieldAccess.Right.Kind)
	}
}

func (sema *sema) checkPrototypeCall(
	prototypeCall *ast.FnCall,
	callScope *ast.Scope,
	extern *ast.ExternDecl,
) error {
	symbol, err := extern.Scope.LookupCurrentScope(prototypeCall.Name.Name())
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

	if symbol.Kind != ast.KIND_PROTO {
		return fmt.Errorf("expected prototype function, not %s\n", symbol.Kind)
	}

	prototype := symbol.Node.(*ast.Proto)
	if prototype.Params.IsVariadic && len(prototypeCall.Args) < len(prototype.Params.Fields) {
		return fmt.Errorf("expected at least %d arguments, got %s\n", len(prototype.Params.Fields), len(prototypeCall.Args))
	}

	argIndex := 0
	for i, param := range prototype.Params.Fields {
		arg := prototypeCall.Args[i]
		if _, err := sema.inferExprTypeWithContext(arg, param.Type, callScope); err != nil {
			return err
		}
		argIndex++
	}

	if prototype.Params.IsVariadic {
		for i := argIndex; i < len(prototypeCall.Args); i++ {
			// TODO: deal with variadic argument
			fmt.Printf("variadic arg: %s\n", prototypeCall.Args[i])
		}
	}

	return nil
}

func (sema *sema) checkForLoop(
	forLoop *ast.ForLoop,
	returnTy *ast.ExprType,
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
	returnTy *ast.ExprType,
) error {
	err := sema.checkIfExpr(whileLoop.Cond, scope)
	if err != nil {
		return err
	}
	err = sema.checkBlock(whileLoop.Block, returnTy, scope)
	return err
}
