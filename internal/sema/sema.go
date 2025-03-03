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

	pkg  *ast.Package
	file *ast.File
}

func New(collector *diagnostics.Collector) *sema {
	s := new(sema)
	s.pkg = nil
	s.file = nil
	s.mainPackageFound = false
	s.collector = collector
	return s
}

func (s *sema) Check(program *ast.Program) error {
	return s.checkPackage(program.Root)
}

func (s *sema) checkPackage(pkg *ast.Package) error {
	prevPkg := s.pkg
	defer func() { s.pkg = prevPkg }()

	s.pkg = pkg

	err := s.checkPackageFiles(pkg)
	if err != nil {
		return err
	}

	for _, innerPackage := range pkg.Packages {
		err := s.checkPackage(innerPackage)
		if err != nil {
			return err
		}
	}

	return nil
}

func (s *sema) checkPackageFiles(pkg *ast.Package) error {
	if pkg.Loc.Name == "main" {
		return fmt.Errorf("package name is not allowed to be 'main'")
	}

	hasMainMethod := false
	requiresMainMethod := false

	for _, file := range pkg.Files {
		prevFile := s.file
		defer func() { s.file = prevFile }()

		s.file = file

		for _, imp := range file.Imports {
			err := s.checkPackage(imp)
			if err != nil {
				return err
			}
		}

		if file.PkgName == "main" {
			if s.mainPackageFound {
				// TODO(errors)
				return fmt.Errorf("error: main package already defined somewhere else")
			}
			s.mainPackageFound = true
			requiresMainMethod = true
		}

		if requiresMainMethod && file.PkgName != "main" {
			// TODO(errors)
			return fmt.Errorf("error: expected package name to be 'main'\n")
		}
		if !requiresMainMethod && file.PkgName != pkg.Loc.Name {
			return fmt.Errorf("error: expected package name to be '%s'\n", pkg.Loc.Name)
		}

		fileHasMain, err := s.checkFile(file)
		if err != nil {
			return err
		}
		hasMainMethod = hasMainMethod || fileHasMain
	}

	// TODO(errors)
	if requiresMainMethod && !hasMainMethod {
		return fmt.Errorf("main package requires a 'main' method as entrypoint\n")
	}

	// TODO(errors)
	if !requiresMainMethod && hasMainMethod {
		return fmt.Errorf("'main' method not allowed on non-main packages\n")
	}

	return nil
}

func (s *sema) checkFile(file *ast.File) (bool, error) {
	var foundMain bool

	for _, node := range file.Body {
		switch node.Kind {
		case ast.KIND_FN_DECL:
			fnDecl := node.Node.(*ast.FnDecl)
			if !foundMain {
				foundMain = fnDecl.Name.Name() == "main"
			}
			err := s.checkFnDecl(fnDecl, s.pkg.Scope)
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
			continue
		default:
			return false, fmt.Errorf("unimplemented ast node for sema: %s\n", reflect.TypeOf(node.Node))
		}
	}

	return foundMain, nil
}

func (sema *sema) checkFnDecl(function *ast.FnDecl, declScope *ast.Scope) error {
	if function.Attributes != nil {
		err := sema.checkFnAttributes(function.Attributes)
		if err != nil {
			return err
		}
	}
	err := sema.checkBlock(function.Block, function.RetType, function.Scope, declScope)
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
	referenceScope *ast.Scope,
	declScope *ast.Scope,
) error {
	for _, statement := range block.Statements {
		err := sema.checkStmt(statement, referenceScope, declScope, returnTy)
		if err != nil {
			return err
		}
	}
	return nil
}

func (sema *sema) checkStmt(
	stmt *ast.Node,
	referenceScope *ast.Scope,
	declScope *ast.Scope,
	returnTy *ast.ExprType,
) error {
	switch stmt.Kind {
	case ast.KIND_FN_CALL:
		_, err := sema.checkFnCall(stmt.Node.(*ast.FnCall), referenceScope, declScope)
		return err
	case ast.KIND_VAR_STMT:
		err := sema.checkVar(stmt.Node.(*ast.VarStmt), referenceScope, declScope)
		return err
	case ast.KIND_COND_STMT:
		err := sema.checkCondStmt(stmt.Node.(*ast.CondStmt), returnTy, referenceScope, declScope)
		return err
	case ast.KIND_RETURN_STMT:
		returnStmt := stmt.Node.(*ast.ReturnStmt)
		_, err := sema.inferExprTypeWithContext(returnStmt.Value, returnTy, referenceScope, declScope)
		return err
	case ast.KIND_NAMESPACE_ACCESS:
		_, err := sema.checkNamespaceAccess(stmt.Node.(*ast.NamespaceAccess), referenceScope, declScope)
		return err
	case ast.KIND_FOR_LOOP_STMT:
		err := sema.checkForLoop(stmt.Node.(*ast.ForLoop), returnTy, referenceScope, declScope)
		return err
	case ast.KIND_WHILE_LOOP_STMT:
		err := sema.checkWhileLoop(stmt.Node.(*ast.WhileLoop), referenceScope, declScope, returnTy)
		return err
	case ast.KIND_DEFER_STMT:
		deferStmt := stmt.Node.(*ast.DeferStmt)
		err := sema.checkStmt(deferStmt.Stmt, referenceScope, declScope, returnTy)
		return err
	default:
		return fmt.Errorf("error: unimplemented statement on sema: %s\n", stmt.Kind)
	}
}

func (sema *sema) checkVar(variable *ast.VarStmt, referenceScope *ast.Scope, declScope *ast.Scope) error {
	for _, currentVar := range variable.Names {
		if variable.IsDecl {
			_, err := referenceScope.LookupCurrentScope(currentVar.Name.Name())
			if err == nil {
				return fmt.Errorf("'%s' already declared in the current block", currentVar.Name.Name())
			}

			n := new(ast.Node)
			n.Kind = ast.KIND_VAR_ID_STMT
			n.Node = currentVar
			currentVar.N = n

			err = referenceScope.Insert(currentVar.Name.Name(), n)
			if err != nil {
				return err
			}
		} else {
			symbol, err := referenceScope.LookupAcrossScopes(currentVar.Name.Name())
			// TODO(errors)
			if err != nil {
				if err == ast.ErrSymbolNotFoundOnScope {
					name := currentVar.Name.Name()
					pos := currentVar.Name.Pos
					return fmt.Errorf("%s '%s' not declared in the current block", pos, name)
				}
				return err
			}
			currentVar.N = symbol
		}
	}

	switch variable.Expr.Kind {
	case ast.KIND_TUPLE_EXPR:
		tuple := variable.Expr.Node.(*ast.TupleExpr)
		err := sema.checkTupleExprAssignedToVariable(variable, tuple, referenceScope, declScope)
		return err
	case ast.KIND_FN_CALL:
		fnCall := variable.Expr.Node.(*ast.FnCall)
		fnDecl, err := sema.checkFnCall(fnCall, referenceScope, declScope)
		if err != nil {
			return err
		}

		if fnDecl.RetType.Kind == ast.EXPR_TYPE_TUPLE {
			err := sema.checkTupleTypeAssignedToVariable(variable.Names, fnDecl.RetType.T.(*ast.TupleType), referenceScope)
			return err
		} else {
			// TODO(errors)
			if len(variable.Names) != 1 {
				return fmt.Errorf("more variables than expressions\n")
			}
			err := sema.checkVarExpr(variable.Names[0], variable.Expr, referenceScope, declScope)
			if err != nil {
				return err
			}
		}
	default:
		// TODO(errors)
		if len(variable.Names) != 1 {
			return fmt.Errorf("more variables than expressions\n")
		}
		err := sema.checkVarExpr(variable.Names[0], variable.Expr, referenceScope, declScope)
		if err != nil {
			return err
		}
	}

	return nil
}

func (sema *sema) checkTupleExprAssignedToVariable(variable *ast.VarStmt, tuple *ast.TupleExpr, referenceScope *ast.Scope, declScope *ast.Scope) error {
	numExprs, err := sema.countExprsOnTuple(tuple, referenceScope, declScope)
	if err != nil {
		return err
	}

	// TODO(errors)
	if len(variable.Names) != numExprs {
		return fmt.Errorf("%d != %d", len(variable.Names), numExprs)
	}

	t := 0
	for _, expr := range tuple.Exprs {
		switch expr.Kind {
		case ast.KIND_TUPLE_EXPR:
			innerTupleExpr := expr.Node.(*ast.TupleExpr)
			for _, innerExpr := range innerTupleExpr.Exprs {
				err := sema.checkVarExpr(variable.Names[t], innerExpr, referenceScope, declScope)
				if err != nil {
					return err
				}
				t++
			}
		case ast.KIND_FN_CALL:
			fnCall := expr.Node.(*ast.FnCall)
			fnDecl, err := sema.checkFnCall(fnCall, referenceScope, declScope)
			if err != nil {
				return err
			}

			if fnDecl.RetType.Kind == ast.EXPR_TYPE_TUPLE {
				tupleType := fnDecl.RetType.T.(*ast.TupleType)
				affectedVariables := variable.Names[t : t+len(tupleType.Types)]
				sema.checkTupleTypeAssignedToVariable(affectedVariables, tupleType, referenceScope)
				t += len(affectedVariables)
			} else {
				err := sema.checkVarExpr(variable.Names[t], expr, referenceScope, declScope)
				if err != nil {
					return err
				}
				t++
			}
		default:
			err := sema.checkVarExpr(variable.Names[t], expr, referenceScope, declScope)
			if err != nil {
				return err
			}
			t++
		}
	}
	return nil
}

func (sema *sema) checkTupleTypeAssignedToVariable(variables []*ast.VarId, tupleTy *ast.TupleType, referenceScope *ast.Scope) error {
	if len(variables) != len(tupleTy.Types) {
		return fmt.Errorf("expected %d variables, but got %d", len(tupleTy.Types), len(variables))
	}
	for i, ty := range tupleTy.Types {
		variables[i].Type = ty
	}
	return nil
}

func (sema *sema) checkVarExpr(variable *ast.VarId, expr *ast.Node, referenceScope *ast.Scope, declScope *ast.Scope) error {
	if variable.NeedsInference {
		// TODO(errors): need a test for it
		if variable.Type != nil {
			return fmt.Errorf(
				"needs inference, but variable already has a type: %s",
				variable.Type,
			)
		}
		exprType, _, err := sema.inferExprTypeWithoutContext(expr, referenceScope, declScope)
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
		exprTy, err := sema.inferExprTypeWithContext(expr, variable.Type, referenceScope, declScope)
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

func (sema *sema) countExprsOnTuple(tuple *ast.TupleExpr, referenceScope *ast.Scope, declScope *ast.Scope) (int, error) {
	counter := 0
	for _, expr := range tuple.Exprs {
		switch expr.Kind {
		case ast.KIND_TUPLE_EXPR:
			varTuple := expr.Node.(*ast.TupleExpr)
			innerTupleExprs, err := sema.countExprsOnTuple(varTuple, referenceScope, declScope)
			if err != nil {
				return -1, err
			}
			counter += innerTupleExprs
		case ast.KIND_FN_CALL:
			fnCall := expr.Node.(*ast.FnCall)
			fnDecl, err := sema.checkFnCall(fnCall, referenceScope, declScope)
			if err != nil {
				return -1, err
			}
			if fnDecl.RetType.Kind == ast.EXPR_TYPE_TUPLE {
				fnTuple := fnDecl.RetType.T.(*ast.TupleType)
				counter += len(fnTuple.Types)
			} else {
				counter++
			}
		default:
			counter++
		}
	}

	return counter, nil
}

// Useful for testing
func checkVarDeclFrom(input, filename string) (*ast.VarId, error) {
	collector := diagnostics.New()

	src := []byte(input)
	loc := new(ast.Loc)
	loc.Name = filename
	lex := lexer.New(loc, src, collector)
	par := parser.NewWithLex(lex, collector)

	tmpScope := ast.NewScope(nil)
	varStmt, err := par.ParseIdStmt(tmpScope)
	if err != nil {
		return nil, err
	}

	sema := New(collector)
	parent := ast.NewScope(nil)
	referenceScope := ast.NewScope(parent)
	declScope := ast.NewScope(parent)
	err = sema.checkVar(varStmt.Node.(*ast.VarStmt), referenceScope, declScope)
	if err != nil {
		return nil, err
	}

	return varStmt.Node.(*ast.VarId), nil
}

func (sema *sema) checkCondStmt(
	condStmt *ast.CondStmt,
	returnTy *ast.ExprType,
	referenceScope *ast.Scope,
	declScope *ast.Scope,
) error {
	err := sema.checkIfExpr(condStmt.IfStmt.Expr, referenceScope, declScope)
	// TODO(errors)
	if err != nil {
		return err
	}

	err = sema.checkBlock(condStmt.IfStmt.Block, returnTy, condStmt.IfStmt.Scope, condStmt.IfStmt.Scope)
	// TODO(errors)
	if err != nil {
		return err
	}

	for i := range condStmt.ElifStmts {
		err := sema.checkIfExpr(condStmt.ElifStmts[i].Expr, condStmt.ElifStmts[i].Scope, declScope)
		// TODO(errors)
		if err != nil {
			return err
		}
		err = sema.checkBlock(condStmt.ElifStmts[i].Block, returnTy, condStmt.ElifStmts[i].Scope, declScope)
		// TODO(errors)
		if err != nil {
			return err
		}
	}

	if condStmt.ElseStmt != nil {
		err = sema.checkBlock(condStmt.ElseStmt.Block, returnTy, condStmt.ElseStmt.Scope, declScope)
		// TODO(errors)
		if err != nil {
			return err
		}
	}

	return nil
}

func (sema *sema) checkFnCall(
	fnCall *ast.FnCall,
	referenceScope *ast.Scope,
	declScope *ast.Scope,
) (*ast.FnDecl, error) {
	symbol, err := declScope.LookupAcrossScopes(fnCall.Name.Name())
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
			return nil, diagnostics.COMPILER_ERROR_FOUND
		}
		return nil, err
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
		return nil, diagnostics.COMPILER_ERROR_FOUND
	}

	fnDecl := symbol.Node.(*ast.FnDecl)
	fnCall.Decl = fnDecl

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
		return nil, diagnostics.COMPILER_ERROR_FOUND
	}

	for i, arg := range fnCall.Args {
		paramType := fnDecl.Params.Fields[i].Type
		argType, err := sema.inferExprTypeWithContext(arg, paramType, referenceScope, declScope)
		if err != nil {
			return nil, err
		}
		if !reflect.DeepEqual(argType, paramType) {
			mismatchedArgType := diagnostics.Diag{
				// TODO(errors): add position of the error
				Message: fmt.Sprintf("can't use %s on argument of type %s", argType, paramType),
			}
			sema.collector.ReportAndSave(mismatchedArgType)
			return nil, diagnostics.COMPILER_ERROR_FOUND
		}
	}

	// TODO: deal with variadic arguments
	return fnDecl, nil
}

func (sema *sema) checkIfExpr(expr *ast.Node, referenceScope *ast.Scope, declScope *ast.Scope) error {
	inferedExprType, _, err := sema.inferExprTypeWithoutContext(expr, referenceScope, declScope)
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
	referenceScope *ast.Scope,
	declScope *ast.Scope,
) (*ast.ExprType, error) {
	switch expr.Kind {
	case ast.KIND_LITERAl_EXPR:
		return s.inferLiteralExprTypeWithContext(expr.Node.(*ast.LiteralExpr), expectedType)
	case ast.KIND_ID_EXPR:
		return s.inferIdExprTypeWithContext(expr.Node.(*ast.IdExpr), expectedType, referenceScope)
	case ast.KIND_BINARY_EXPR:
		return s.inferBinaryExprTypeWithContext(expr.Node.(*ast.BinaryExpr), expectedType, referenceScope, declScope)
	case ast.KIND_UNARY_EXPR:
		return s.inferUnaryExprTypeWithContext(expr.Node.(*ast.UnaryExpr), expectedType, referenceScope, declScope)
	case ast.KIND_FN_CALL:
		return s.inferFnCallExprTypeWithContext(expr.Node.(*ast.FnCall), expectedType, referenceScope, declScope)
	case ast.KIND_VOID_EXPR:
		return s.inferVoidExprTypeWithContext(expectedType)
	case ast.KIND_TUPLE_EXPR:
		return s.inferTupleExprTypeWithContext(expr.Node.(*ast.TupleExpr), expectedType, referenceScope, declScope)
	default:
		log.Fatalf("unimplemented expression: %s\n", expr.Kind)
		return nil, nil
	}
}

func (sema *sema) inferLiteralExprTypeWithContext(
	literal *ast.LiteralExpr,
	expectedType *ast.ExprType,
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
	referenceScope *ast.Scope,
) (*ast.ExprType, error) {
	symbol, err := referenceScope.LookupAcrossScopes(id.Name.Name())
	// TODO(errors)
	if err != nil {
		return nil, err
	}
	id.N = symbol

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
	referenceScope *ast.Scope,
	declScope *ast.Scope,
) (*ast.ExprType, error) {
	lhsType, err := sema.inferExprTypeWithContext(binary.Left, expectedType, referenceScope, declScope)
	if err != nil {
		return nil, err
	}

	rhsType, err := sema.inferExprTypeWithContext(binary.Right, expectedType, referenceScope, declScope)
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
	referenceScope *ast.Scope,
	declScope *ast.Scope,
) (*ast.ExprType, error) {
	switch unary.Op {
	case token.MINUS:
		unaryExprType, err := sema.inferExprTypeWithContext(unary.Value, expectedType, referenceScope, declScope)
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
	referenceScope *ast.Scope,
	declScope *ast.Scope,
) (*ast.ExprType, error) {
	fnDecl, err := sema.checkFnCall(fnCall, referenceScope, declScope)
	if err != nil {
		return nil, err
	}
	return fnDecl.RetType, nil
}

func (sema *sema) inferVoidExprTypeWithContext(expectedType *ast.ExprType) (*ast.ExprType, error) {
	// TODO(errors)
	if !expectedType.IsVoid() {
		return nil, fmt.Errorf("expected return type to be '%s'", expectedType)
	}
	return expectedType, nil
}

func (sema *sema) inferTupleExprTypeWithContext(
	tuple *ast.TupleExpr,
	expectedType *ast.ExprType,
	referenceScope *ast.Scope,
	declScope *ast.Scope,
) (*ast.ExprType, error) {
	// TODO(errors)
	if expectedType.Kind != ast.EXPR_TYPE_TUPLE {
		return nil, fmt.Errorf("expected to be %s, got tuple\n", expectedType)
	}
	expectedTuple := expectedType.T.(*ast.TupleType)

	// TODO(errors)
	if len(expectedTuple.Types) != len(tuple.Exprs) {
		return nil, fmt.Errorf("expected %d expressions, but got %d\n", len(expectedTuple.Types), len(tuple.Exprs))
	}

	tupleTy := new(ast.ExprType)
	tupleTy.Kind = ast.EXPR_TYPE_TUPLE

	types := make([]*ast.ExprType, 0)
	for i, expr := range tuple.Exprs {
		ty, err := sema.inferExprTypeWithContext(expr, expectedTuple.Types[i], referenceScope, declScope)
		if err != nil {
			return nil, err
		}
		types = append(types, ty)
	}

	ty := &ast.TupleType{Types: types}
	tupleTy.T = ty
	tuple.Type = ty
	return tupleTy, nil
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
	referenceScope *ast.Scope,
	declScope *ast.Scope,
) (*ast.Node, *ast.ExprType, error) {
	collector := diagnostics.New()

	expr, err := parser.ParseExprFrom(input, filename)
	if err != nil {
		return nil, nil, err
	}

	analyzer := New(collector)
	exprType, _, err := analyzer.inferExprTypeWithoutContext(expr, referenceScope, declScope)
	if err != nil {
		return nil, nil, err
	}
	return expr, exprType, nil
}

// Useful for testing
func inferExprTypeWithContext(
	input, filename string,
	ty *ast.ExprType,
	referenceScope *ast.Scope,
	declScope *ast.Scope,
) (*ast.ExprType, error) {
	collector := diagnostics.New()

	expr, err := parser.ParseExprFrom(input, filename)
	if err != nil {
		return nil, err
	}

	analyzer := New(collector)
	exprType, err := analyzer.inferExprTypeWithContext(expr, ty, referenceScope, declScope)
	if err != nil {
		return nil, err
	}
	return exprType, nil
}

func (sema *sema) inferExprTypeWithoutContext(
	expr *ast.Node,
	referenceScope *ast.Scope,
	declScope *ast.Scope,
) (*ast.ExprType, bool, error) {
	switch expr.Kind {
	case ast.KIND_LITERAl_EXPR:
		return sema.inferLiteralExprTypeWithoutContext(expr.Node.(*ast.LiteralExpr), referenceScope)
	case ast.KIND_ID_EXPR:
		return sema.inferIdExprTypeWithoutContext(expr.Node.(*ast.IdExpr), referenceScope)
	case ast.KIND_UNARY_EXPR:
		return sema.inferUnaryExprTypeWithoutContext(expr.Node.(*ast.UnaryExpr), referenceScope, declScope)
	case ast.KIND_BINARY_EXPR:
		return sema.inferBinaryExprTypeWithoutContext(expr.Node.(*ast.BinaryExpr), referenceScope, declScope)
	case ast.KIND_FN_CALL:
		return sema.inferFnCallExprTypeWithoutContext(expr.Node.(*ast.FnCall), referenceScope, declScope)
	case ast.KIND_NAMESPACE_ACCESS:
		return sema.inferNamespaceAccessExprTypeWithoutContext(expr.Node.(*ast.NamespaceAccess), referenceScope, declScope)
	case ast.KIND_TUPLE_EXPR:
		return sema.inferTupleExprTypeWithoutContext(expr.Node.(*ast.TupleExpr), referenceScope, declScope)
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

	id.N = variable

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
	referenceScope *ast.Scope,
	declScope *ast.Scope,
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
		unaryExpr, foundContext, err := sema.inferExprTypeWithoutContext(unary.Value, referenceScope, declScope)
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
	referenceScope *ast.Scope,
	declScope *ast.Scope,
) (*ast.ExprType, bool, error) {
	lhsType, lhsFoundContext, err := sema.inferExprTypeWithoutContext(binary.Left, referenceScope, declScope)
	// TODO(errors)
	if err != nil {
		return nil, false, err
	}

	rhsType, rhsFoundContext, err := sema.inferExprTypeWithoutContext(binary.Right, referenceScope, declScope)
	// TODO(errors)
	if err != nil {
		return nil, false, err
	}

	if lhsFoundContext && !rhsFoundContext {
		rhsTypeWithContext, err := sema.inferExprTypeWithContext(binary.Right, lhsType, referenceScope, declScope)
		// TODO(errors)
		if err != nil {
			return nil, false, err
		}
		rhsType = rhsTypeWithContext
	}
	if !lhsFoundContext && rhsFoundContext {
		lhsTypeWithContext, err := sema.inferExprTypeWithContext(binary.Left, rhsType, referenceScope, declScope)
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
	referenceScope *ast.Scope,
	declScope *ast.Scope,
) (*ast.ExprType, bool, error) {
	fnDecl, err := sema.checkFnCall(fnCall, referenceScope, declScope)
	return fnDecl.RetType, true, err
}

func (sema *sema) inferNamespaceAccessExprTypeWithoutContext(
	namespaceAccess *ast.NamespaceAccess,
	referenceScope *ast.Scope,
	declScope *ast.Scope,
) (*ast.ExprType, bool, error) {
	ty, err := sema.checkNamespaceAccess(namespaceAccess, referenceScope, declScope)
	return ty, true, err
}

func (sema *sema) inferTupleExprTypeWithoutContext(
	tuple *ast.TupleExpr,
	referenceScope *ast.Scope,
	declScope *ast.Scope,
) (*ast.ExprType, bool, error) {
	tupleTy := new(ast.ExprType)
	tupleTy.Kind = ast.EXPR_TYPE_TUPLE

	types := make([]*ast.ExprType, 0)
	for _, expr := range tuple.Exprs {
		innerTy, _, err := sema.inferExprTypeWithoutContext(expr, referenceScope, declScope)
		if err != nil {
			return nil, false, err
		}
		types = append(types, innerTy)
	}

	tupleTy.T = &ast.TupleType{Types: types}
	return tupleTy, false, nil
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

func (sema *sema) checkNamespaceAccess(
	namespaceAccess *ast.NamespaceAccess,
	referenceScope *ast.Scope,
	declScope *ast.Scope,
) (*ast.ExprType, error) {
	if namespaceAccess.IsImport {
		imp := namespaceAccess.Left.Name.Name()
		pkg, ok := sema.file.Imports[imp]
		if !ok {
			panic("package not found")
		}
		return sema.checkImportAccess(namespaceAccess.Right, referenceScope, pkg.Scope)
	}

	left, err := referenceScope.LookupAcrossScopes(namespaceAccess.Left.Name.Name())
	if err != nil {
		return nil, err
	}
	namespaceAccess.Left.N = left

	// TODO(errors)
	if left.Kind != ast.KIND_EXTERN_DECL {
		return nil, fmt.Errorf("invalid access")
	}

	externDecl := left.Node.(*ast.ExternDecl)
	if namespaceAccess.Right.Kind != ast.KIND_FN_CALL {
		return nil, fmt.Errorf("expected prototype call")
	}
	protoCall, err := sema.checkPrototypeCall(externDecl, namespaceAccess.Right.Node.(*ast.FnCall), referenceScope, declScope)
	if err != nil {
		return nil, err
	}
	return protoCall.RetType, nil
}

func (sema *sema) checkImportAccess(node *ast.Node, referenceScope *ast.Scope, declScope *ast.Scope) (*ast.ExprType, error) {
	switch node.Kind {
	case ast.KIND_FN_CALL:
		fnCall := node.Node.(*ast.FnCall)
		fnDecl, err := sema.checkFnCall(fnCall, referenceScope, declScope)
		return fnDecl.RetType, err
	case ast.KIND_NAMESPACE_ACCESS:
		namespaceAccess := node.Node.(*ast.NamespaceAccess)
		ty, err := sema.checkNamespaceAccess(namespaceAccess, referenceScope, declScope)
		return ty, err
	case ast.KIND_EXTERN_DECL:
		return nil, fmt.Errorf("nothing do with a extern declaration, try accessing a prototype")
	default:
		return nil, fmt.Errorf("unimplemented symbol to import: %s\n", node.Kind)
	}
}

func (sema *sema) checkPrototypeCall(
	extern *ast.ExternDecl,
	prototypeCall *ast.FnCall,
	referenceScope *ast.Scope,
	declScope *ast.Scope,
) (*ast.Proto, error) {
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
		return nil, diagnostics.COMPILER_ERROR_FOUND
	}
	if symbol.Kind != ast.KIND_PROTO {
		return nil, fmt.Errorf("expected prototype function, not %s\n", symbol.Kind)
	}

	prototype := symbol.Node.(*ast.Proto)
	if prototype.Params.IsVariadic && len(prototypeCall.Args) < len(prototype.Params.Fields) {
		return nil, fmt.Errorf("expected at least %d arguments, got %s\n", len(prototype.Params.Fields), len(prototypeCall.Args))
	}
	prototypeCall.Proto = prototype

	argIndex := 0
	for i, param := range prototype.Params.Fields {
		arg := prototypeCall.Args[i]
		if _, err := sema.inferExprTypeWithContext(arg, param.Type, referenceScope, declScope); err != nil {
			return nil, err
		}
		argIndex++
	}

	if prototype.Params.IsVariadic {
		for i := argIndex; i < len(prototypeCall.Args); i++ {
			// TODO: deal with variadic argument
			fmt.Printf("variadic arg: %s\n", prototypeCall.Args[i])
		}
	}

	return prototype, nil
}

func (sema *sema) checkForLoop(
	forLoop *ast.ForLoop,
	returnTy *ast.ExprType,
	referenceScope *ast.Scope,
	declScope *ast.Scope,
) error {
	err := sema.checkStmt(forLoop.Init, referenceScope, declScope, returnTy)
	if err != nil {
		return err
	}

	err = sema.checkIfExpr(forLoop.Cond, referenceScope, declScope)
	if err != nil {
		return err
	}

	err = sema.checkStmt(forLoop.Update, referenceScope, declScope, returnTy)
	if err != nil {
		return err
	}

	err = sema.checkBlock(forLoop.Block, returnTy, referenceScope, declScope)
	return err
}

func (sema *sema) checkWhileLoop(
	whileLoop *ast.WhileLoop,
	referenceScope *ast.Scope,
	declScope *ast.Scope,
	returnTy *ast.ExprType,
) error {
	err := sema.checkIfExpr(whileLoop.Cond, referenceScope, declScope)
	if err != nil {
		return err
	}
	err = sema.checkBlock(whileLoop.Block, returnTy, referenceScope, declScope)
	return err
}
