package sema

import (
	"bufio"
	"fmt"
	"log"
	"reflect"
	"strconv"
	"strings"

	"github.com/HicaroD/Telia/ast"
	"github.com/HicaroD/Telia/collector"
	"github.com/HicaroD/Telia/lexer"
	"github.com/HicaroD/Telia/lexer/token/kind"
	"github.com/HicaroD/Telia/parser"
	"github.com/HicaroD/Telia/scope"
)

type sema struct {
	collector *collector.DiagCollector
	universe  *scope.Scope[ast.Node]
}

func New(collector *collector.DiagCollector) *sema {
	// "universe" scope does not have any parent, it is the
	// root of the tree of scopes
	var nilScope *scope.Scope[ast.Node] = nil
	universe := scope.New(nilScope)
	return &sema{collector, universe}
}

func (sema *sema) Analyze(nodes []ast.Node) error {
	for i := range nodes {
		switch astNode := nodes[i].(type) {
		case *ast.FunctionDecl:
			err := sema.analyzeFnDecl(astNode)
			if err != nil {
				return err
			}
		case *ast.ExternDecl:
			err := sema.analyzeExtern(astNode)
			if err != nil {
				return err
			}
		default:
			log.Fatalf("unimplemented ast node for sema: %s\n", reflect.TypeOf(astNode))
		}
	}
	return nil
}

func (sema *sema) analyzeExtern(extern *ast.ExternDecl) error {
	externScope := scope.New(sema.universe)
	for i := range extern.Prototypes {
		prototypeName := extern.Prototypes[i].Name.Name()
		err := externScope.Insert(prototypeName, extern.Prototypes[i])
		if err != nil {
			if err == scope.ERR_SYMBOL_ALREADY_DEFINED_ON_SCOPE {
				pos := extern.Prototypes[i].Name.Position
				prototypeRedeclaration := collector.Diag{
					Message: fmt.Sprintf(
						"%s:%d:%d: prototype '%s' already declared on extern '%s'",
						pos.Filename,
						pos.Line,
						pos.Column,
						prototypeName,
						extern.Name.Name(),
					),
				}
				sema.collector.ReportAndSave(prototypeRedeclaration)
				return collector.COMPILER_ERROR_FOUND
			}
			return err
		}
	}
	extern.Scope = externScope
	err := sema.universe.Insert(extern.Name.Name(), extern)
	if err != nil {
		if err == scope.ERR_SYMBOL_ALREADY_DEFINED_ON_SCOPE {
			pos := extern.Name.Position
			prototypeRedeclaration := collector.Diag{
				Message: fmt.Sprintf(
					"%s:%d:%d: extern '%s' already declared on scope",
					pos.Filename,
					pos.Line,
					pos.Column,
					extern.Name.Name(),
				),
			}
			sema.collector.ReportAndSave(prototypeRedeclaration)
			return collector.COMPILER_ERROR_FOUND
		}
		return err
	}
	return nil
}

func (sema *sema) analyzeFnDecl(function *ast.FunctionDecl) error {
	// TODO: Is it really correct to insert in the universe scope?
	// In the future, I'll isolate these functions into their modules
	err := sema.universe.Insert(function.Name.Name(), function)
	if err != nil {
		// TODO: show the position of the first declaration
		// for helping the program
		if err == scope.ERR_SYMBOL_ALREADY_DEFINED_ON_SCOPE {
			pos := function.Name.Position
			functionRedeclaration := collector.Diag{
				Message: fmt.Sprintf(
					"%s:%d:%d: function '%s' already declared on scope",
					pos.Filename,
					pos.Line,
					pos.Column,
					function.Name.Name(),
				),
			}
			sema.collector.ReportAndSave(functionRedeclaration)
			return collector.COMPILER_ERROR_FOUND
		}
		return err
	}

	function.Scope = scope.New(sema.universe)
	err = sema.addParametersToScope(function.Params, function.Name.Name(), function.Scope)
	if err != nil {
		return err
	}

	err = sema.analyzeBlock(function.Scope, function.Block, function.RetType)
	if err != nil {
		return err
	}
	return nil
}

func (sema *sema) addParametersToScope(
	params *ast.FieldList,
	functionName string,
	functionScope *scope.Scope[ast.Node],
) error {
	for _, param := range params.Fields {
		paramName := param.Name.Name()
		err := functionScope.Insert(paramName, param)
		if err != nil {
			if err == scope.ERR_SYMBOL_ALREADY_DEFINED_ON_SCOPE {
				pos := param.Name.Position
				parameterRedeclaration := collector.Diag{
					Message: fmt.Sprintf(
						"%s:%d:%d: parameter '%s' already declared on function '%s'",
						pos.Filename,
						pos.Line,
						pos.Column,
						paramName,
						functionName,
					),
				}
				sema.collector.ReportAndSave(parameterRedeclaration)
				return collector.COMPILER_ERROR_FOUND
			}
			return err
		}
	}
	return nil
}

func (sema *sema) analyzeBlock(
	scope *scope.Scope[ast.Node],
	block *ast.BlockStmt,
	returnTy ast.ExprType,
) error {
	for i := range block.Statements {
		switch statement := block.Statements[i].(type) {
		case *ast.FunctionCall:
			err := sema.analyzeFunctionCall(statement, scope)
			if err != nil {
				return err
			}
		case *ast.MultiVarStmt:
			err := sema.analyzeVarDecl(statement, scope)
			// TODO(errors)
			if err != nil {
				return err
			}
		case *ast.CondStmt:
			err := sema.analyzeCondStmt(statement, returnTy, scope)
			// TODO(errors)
			if err != nil {
				return err
			}
		case *ast.ReturnStmt:
			_, err := sema.inferExprTypeWithContext(statement.Value, returnTy, scope)
			// TODO(errors)
			if err != nil {
				return err
			}
		case *ast.FieldAccess:
			err := sema.analyzeFieldAccessExpr(statement, scope)
			// TODO(errors)
			if err != nil {
				return err
			}
		case *ast.ForLoop:
			err := sema.analyzeForLoop(statement, scope)
			if err != nil {
				return err
			}
		default:
			log.Fatalf("unimplemented statement on sema: %s", statement)
		}
	}
	return nil
}

func (sema *sema) analyzeVarDecl(
	multi *ast.MultiVarStmt,
	currentScope *scope.Scope[ast.Node],
) error {
	for i := range multi.Variables {
		varDecl := multi.Variables[i]
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
			// TODO: what do I assert here in order to make it right?
			exprTy, err := sema.inferExprTypeWithContext(varDecl.Value, varDecl.Type, currentScope)
			// TODO(errors): Deal with type mismatch
			if err != nil {
				return err
			}
			if !reflect.DeepEqual(varDecl.Type, exprTy) {
				return fmt.Errorf("type mismatch on variable decl, expected %s, got %s", varDecl.Type, exprTy)
			}
		}

		varName := varDecl.Name.Name()
		err := currentScope.Insert(varName, varDecl)
		// TODO(errors): symbol already defined on scope
		if err != nil {
			return err
		}
	}
	return nil
}

// Useful for testing
func analyzeVarDeclFrom(input, filename string) (*ast.MultiVarStmt, error) {
	diagCollector := collector.New()

	reader := bufio.NewReader(strings.NewReader(input))
	lexer := lexer.New(filename, reader, diagCollector)

	tokens, err := lexer.Tokenize()
	if err != nil {
		return nil, err
	}

	par := parser.New(tokens, diagCollector)
	idStmt, err := par.ParseIdStmt()
	if err != nil {
		return nil, err
	}

	varDecl, ok := idStmt.(*ast.MultiVarStmt)
	if !ok {
		return nil, fmt.Errorf("unable to cast %s to *ast.VarDeclStmt", reflect.TypeOf(varDecl))
	}

	sema := New(diagCollector)

	parent := scope.New[ast.Node](nil)
	scope := scope.New(parent)

	err = sema.analyzeVarDecl(varDecl, scope)
	if err != nil {
		return nil, err
	}

	return varDecl, nil
}

func (sema *sema) analyzeCondStmt(
	condStmt *ast.CondStmt,
	returnTy ast.ExprType,
	outterScope *scope.Scope[ast.Node],
) error {
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

func (sema *sema) analyzeFunctionCall(
	functionCall *ast.FunctionCall,
	currentScope *scope.Scope[ast.Node],
) error {
	function, err := currentScope.LookupAcrossScopes(functionCall.Name.Name())
	if err != nil {
		if err == scope.ERR_SYMBOL_NOT_FOUND_ON_SCOPE {
			pos := functionCall.Name.Position
			functionNotDefined := collector.Diag{
				Message: fmt.Sprintf(
					"%s:%d:%d: function '%s' not defined on scope",
					pos.Filename,
					pos.Line,
					pos.Column,
					functionCall.Name.Name(),
				),
			}
			sema.collector.ReportAndSave(functionNotDefined)
			return collector.COMPILER_ERROR_FOUND
		}
		return err
	}

	decl, ok := function.(*ast.FunctionDecl)
	if !ok {
		pos := functionCall.Name.Position
		notCallable := collector.Diag{
			Message: fmt.Sprintf(
				"%s:%d:%d: '%s' is not callable",
				pos.Filename,
				pos.Line,
				pos.Column,
				functionCall.Name.Name(),
			),
		}
		sema.collector.ReportAndSave(notCallable)
		return collector.COMPILER_ERROR_FOUND
	}

	if len(functionCall.Args) != len(decl.Params.Fields) {
		pos := functionCall.Name.Position
		// TODO(errors): show which arguments were passed and which types we
		// were expecting
		notEnoughArguments := collector.Diag{
			Message: fmt.Sprintf(
				"%s:%d:%d: not enough arguments in call to '%s'",
				pos.Filename,
				pos.Line,
				pos.Column,
				functionCall.Name.Name(),
			),
		}
		sema.collector.ReportAndSave(notEnoughArguments)
		return collector.COMPILER_ERROR_FOUND
	}

	for i := range len(functionCall.Args) {
		paramType := decl.Params.Fields[i].Type
		argType, err := sema.inferExprTypeWithContext(functionCall.Args[i], paramType, currentScope)
		if err != nil {
			return err
		}
		if !reflect.DeepEqual(argType, paramType) {
			mismatchedArgType := collector.Diag{
				// TODO(errors): add position of the error
				Message: fmt.Sprintf("can't use %s on argument of type %s", argType, paramType),
			}
			sema.collector.ReportAndSave(mismatchedArgType)
			return collector.COMPILER_ERROR_FOUND
		}
	}

	// TODO: deal with variadic arguments
	return nil
}

func (sema *sema) analyzeIfExpr(expr ast.Expr, scope *scope.Scope[ast.Node]) error {
	inferedExprType, _, err := sema.inferExprTypeWithoutContext(expr, scope)
	// TODO(errors)
	if err != nil {
		return err
	}

	if !inferedExprType.IsBoolean() {
		log.Fatalf("invalid non-boolean condition on if statement: %s", inferedExprType)
	}

	return nil
}

// TODO: think about the way I'm inferring or getting the expr type correctly
func (sema *sema) inferExprTypeWithContext(
	exprNode ast.Expr,
	expectedType ast.ExprType,
	scope *scope.Scope[ast.Node],
) (ast.ExprType, error) {
	switch expression := exprNode.(type) {
	case *ast.LiteralExpr:
		switch ty := expression.Type.(type) {
		case *ast.BasicType:
			switch ty.Kind {
			case kind.STRING_LITERAL:
				finalTy := &ast.PointerType{Type: &ast.BasicType{Kind: kind.U8_TYPE}}
				expression.Type = finalTy
				return finalTy, nil
			case kind.INTEGER_LITERAL, kind.INT_TYPE:
				switch ty := expectedType.(type) {
				case *ast.BasicType:
					value := expression.Value
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
					log.Fatalf("unimplemented type at integer literal: %s %s", ty, reflect.TypeOf(ty))
				}
			case kind.TRUE_BOOL_LITERAL:
				finalTy := &ast.BasicType{Kind: kind.BOOL_TYPE}
				expression.Type = finalTy
				expression.Value = "1"
				return finalTy, nil
			case kind.FALSE_BOOL_LITERAL:
				finalTy := &ast.BasicType{Kind: kind.BOOL_TYPE}
				expression.Type = finalTy
				expression.Value = "0"
				return finalTy, nil
			default:
				log.Fatalf("unimplemented basic type kind: %s", ty.Kind)
			}
		default:
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
		// TODO(errors)
		default:
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
		err := sema.analyzeFunctionCall(expression, scope)
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
		case kind.MINUS:
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
	scope *scope.Scope[ast.Node],
) (ast.Expr, ast.ExprType, error) {
	collector := collector.New()

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
	scope *scope.Scope[ast.Node],
) (ast.ExprType, error) {
	collector := collector.New()

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
	scope *scope.Scope[ast.Node],
) (ast.ExprType, bool, error) {
	switch expression := expr.(type) {
	case *ast.LiteralExpr:
		switch ty := expression.Type.(type) {
		case *ast.BasicType:
			// TODO: what if I want a defined string type on my programming
			// language? I don't think *i8 is correct type for this
			switch ty.Kind {
			case kind.STRING_LITERAL:
				finalTy := &ast.PointerType{Type: &ast.BasicType{Kind: kind.U8_TYPE}}
				expression.Type = finalTy
				return finalTy, false, nil
			case kind.INTEGER_LITERAL:
				ty, err := sema.inferIntegerType(expression.Value)
				if err != nil {
					return nil, false, err
				}
				expression.Type = ty
				return ty, false, nil
			case kind.TRUE_BOOL_LITERAL:
				finalTy := &ast.BasicType{Kind: kind.BOOL_TYPE}
				expression.Type = finalTy
				expression.Value = "1"
				return finalTy, false, nil
			case kind.FALSE_BOOL_LITERAL:
				finalTy := &ast.BasicType{Kind: kind.BOOL_TYPE}
				expression.Type = finalTy
				expression.Value = "0"
				return finalTy, false, nil
			default:
				log.Fatalf("unimplemented literal expr: %s", expression)
			}
		default:
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
		case kind.MINUS:
			switch unaryExpr := expression.Value.(type) {
			case *ast.LiteralExpr:
				switch unaryTy := unaryExpr.Type.(type) {
				case *ast.BasicType:
					if unaryTy.Kind == kind.INTEGER_LITERAL {
						integerType, err := sema.inferIntegerType(unaryExpr.Value)
						// TODO(errors)
						if err != nil {
							return nil, false, err
						}
						unaryExpr.Type = integerType
						return integerType, false, nil
					}
				default:
					log.Fatalf("unimplemented unary expr type: %s", unaryExpr)
				}
			}
		case kind.NOT:
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
		err := sema.analyzeFunctionCall(expression, scope)
		// TODO(errors)
		if err != nil {
			return nil, false, err
		}
		// At this point, function should exists!
		function, err := scope.LookupAcrossScopes(expression.Name.Name())
		if err != nil {
			log.Fatalf("panic: at this point, function '%s' should exists in current block", expression.Name)
		}
		functionDecl := function.(*ast.FunctionDecl)
		return functionDecl.RetType, true, nil
	default:
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
	scope *scope.Scope[ast.Node],
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

	if !reflect.DeepEqual(lhsType, rhsType) {
		log.Fatalf("mismatched types: %s %s %s", lhsType, expression.Op, rhsType)
	}

	switch expression.Op {
	case kind.PLUS, kind.MINUS:
		if lhsType.IsNumeric() && rhsType.IsNumeric() {
			return lhsType, lhsFoundContext || rhsFoundContext, nil
		}
	default:
		if _, ok := kind.LOGICAL_OP[expression.Op]; ok {
			return &ast.BasicType{Kind: kind.BOOL_TYPE}, lhsFoundContext || rhsFoundContext, nil
		}
	}
	// TODO
	log.Fatalf("UNREACHABLE - inferBinaryExprType")
	return nil, false, nil
}

func (sema *sema) inferBinaryExprTypeWithContext(
	expression *ast.BinaryExpr,
	expectedType ast.ExprType,
	scope *scope.Scope[ast.Node],
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

func (sema *sema) inferIntegerType(value string) (ast.ExprType, error) {
	integerType := kind.INT_TYPE
	base := 10
	intSize := strconv.IntSize

	_, err := strconv.ParseUint(value, base, intSize)
	if err == nil {
		return &ast.BasicType{Kind: integerType}, nil
	}

	return nil, fmt.Errorf("can't parse integer literal: %s", value)
}

func (sema *sema) analyzeFieldAccessExpr(
	fieldAccess *ast.FieldAccess,
	currentScope *scope.Scope[ast.Node],
) error {
	if !fieldAccess.Left.IsId() {
		return fmt.Errorf("invalid expression on field accessing: %s", fieldAccess.Left)
	}

	idExpr := fieldAccess.Left.(*ast.IdExpr)
	id := idExpr.Name.Name()

	symbol, err := currentScope.LookupAcrossScopes(id)
	if err != nil {
		if err == scope.ERR_SYMBOL_NOT_FOUND_ON_SCOPE {
			pos := idExpr.Name.Position
			symbolNotDefined := collector.Diag{
				Message: fmt.Sprintf(
					"%s:%d:%d: '%s' not defined on scope",
					pos.Filename,
					pos.Line,
					pos.Column,
					id,
				),
			}
			sema.collector.ReportAndSave(symbolNotDefined)
			return collector.COMPILER_ERROR_FOUND
		}
		return err
	}

	switch sym := symbol.(type) {
	case *ast.ExternDecl:
		switch right := fieldAccess.Right.(type) {
		case *ast.FunctionCall:
			err := sema.analyzePrototypeCall(right, currentScope, sym)
			if err != nil {
				return err
			}
		default:
			// TODO(errors)
			return fmt.Errorf("invalid expression %s when accessing field", right)
		}
	}
	return nil
}

func (sema *sema) analyzePrototypeCall(
	prototypeCall *ast.FunctionCall,
	callScope *scope.Scope[ast.Node],
	extern *ast.ExternDecl,
) error {
	prototype, err := extern.Scope.LookupCurrentScope(prototypeCall.Name.Name())
	// TODO(errors)
	if err != nil {
		return err
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

func (sema *sema) analyzeForLoop(forLoop *ast.ForLoop, scope *scope.Scope[ast.Node]) error {
	return nil
}
