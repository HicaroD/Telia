package codegen

import (
	"fmt"
	"log"
	"os"
	"reflect"
	"strconv"

	"github.com/HicaroD/Telia/ast"
	"github.com/HicaroD/Telia/codegen/values"
	"github.com/HicaroD/Telia/lexer/token/kind"
	"github.com/HicaroD/Telia/scope"
	"tinygo.org/x/go-llvm"
)

type codegen struct {
	universe    *scope.Scope[values.LLVMValue]
	strLiterals map[string]llvm.Value
	context     llvm.Context
	module      llvm.Module
	builder     llvm.Builder
	astNodes    []ast.AstNode
}

func New(astNodes []ast.AstNode) *codegen {
	// parent of universe scope is nil
	var nilScope *scope.Scope[values.LLVMValue] = nil
	universe := scope.New(nilScope)

	context := llvm.NewContext()
	// TODO: properly define the module name
	// The name of the module could be file name
	module := context.NewModule("tmpmod")
	builder := context.NewBuilder()

	return &codegen{
		universe:    universe,
		context:     context,
		module:      module,
		builder:     builder,
		astNodes:    astNodes,
		strLiterals: map[string]llvm.Value{},
	}
}

func (codegen *codegen) Generate() error {
	for i := range codegen.astNodes {
		switch astNode := codegen.astNodes[i].(type) {
		case *ast.FunctionDecl:
			err := codegen.generateFnDecl(astNode)
			// TODO(errors)
			if err != nil {
				return err
			}
		case *ast.ExternDecl:
			err := codegen.generateExternDecl(astNode)
			// TODO(errors)
			if err != nil {
				return err
			}
		default:
			log.Fatalf("unimplemented: %s\n", reflect.TypeOf(codegen.astNodes[i]))
		}
	}

	err := codegen.generateBitcodeFile()
	return err
}

func (codegen *codegen) generateBitcodeFile() error {
	filename := "telia.ll"
	file, err := os.Create(filename)
	// TODO(errors)
	if err != nil {
		return err
	}

	module := codegen.module.String()
	_, err = file.Write([]byte(module))
	// TODO(errors)
	if err != nil {
		return err
	}

	fmt.Printf("'%s' file generated successfuly\n", filename)
	return nil
}

func (codegen *codegen) generateFnDecl(function *ast.FunctionDecl) error {
	returnType := codegen.getType(function.RetType)
	paramsTypes := codegen.getFieldListTypes(function.Params)
	functionType := llvm.FunctionType(returnType, paramsTypes, function.Params.IsVariadic)
	functionValue := llvm.AddFunction(codegen.module, function.Name, functionType)
	functionBlock := codegen.context.AddBasicBlock(functionValue, "entry")
	fnValue := values.NewFunctionValue(functionValue, functionType, &functionBlock)

	err := codegen.universe.Insert(function.Name, fnValue)
	// TODO(errors)
	if err != nil {
		return err
	}

	fnScope := scope.New(codegen.universe)

	codegen.builder.SetInsertPointAtEnd(functionBlock)
	err = codegen.generateParameters(fnValue, function, fnScope, paramsTypes)
	// TODO(errors)
	if err != nil {
		return err
	}

	_, err = codegen.generateBlock(function.Block, fnScope, fnValue)
	// TODO(errors)
	if err != nil {
		return err
	}

	return err
}

func (codegen *codegen) generateBlock(stmts *ast.BlockStmt, scope *scope.Scope[values.LLVMValue], function *values.Function) (bool, error) {
	for i := range stmts.Statements {
		stmt := stmts.Statements[i]
		err := codegen.generateStmt(stmt, scope, function)
		if stmt.IsReturn() {
			return true, nil
		}
		if err != nil {
			return false, err
		}
	}
	return false, nil
}

func (codegen *codegen) generateStmt(stmt ast.Stmt, scope *scope.Scope[values.LLVMValue], function *values.Function) error {
	switch statement := stmt.(type) {
	case *ast.FunctionCall:
		_, err := codegen.generateFunctionCall(scope, statement)
		// TODO(errors)
		if err != nil {
			return err
		}
	case *ast.ReturnStmt:
		err := codegen.generateReturnStmt(statement, scope)
		// TODO(errors)
		if err != nil {
			return nil
		}
	case *ast.CondStmt:
		err := codegen.generateCondStmt(scope, function, statement)
		// TODO(errors)
		if err != nil {
			return err
		}
	case *ast.VarDeclStmt:
		err := codegen.generateVariableDecl(statement, scope)
		// TODO(errors)
		if err != nil {
			return err
		}
	default:
		log.Fatalf("unimplemented block statement: %s", statement)
	}
	return nil
}

func (codegen *codegen) generateReturnStmt(ret *ast.ReturnStmt, scope *scope.Scope[values.LLVMValue]) error {
	if ret.Value.IsVoid() {
		codegen.builder.CreateRetVoid()
		return nil
	}
	returnValue, err := codegen.getExpr(ret.Value, scope)
	// TODO(errors)
	if err != nil {
		return err
	}
	codegen.builder.CreateRet(returnValue)
	return nil
}

func (codegen *codegen) generateVariableDecl(varDecl *ast.VarDeclStmt, scope *scope.Scope[values.LLVMValue]) error {
	varTy := codegen.getType(varDecl.Type)
	varPtr := codegen.builder.CreateAlloca(varTy, ".ptr")
	varExpr, err := codegen.getExpr(varDecl.Value, scope)

	// TODO(errors)
	if err != nil {
		return err
	}
	codegen.builder.CreateStore(varExpr, varPtr)

	variable := values.Variable{
		Ty:  varTy,
		Ptr: varPtr,
	}
	err = scope.Insert(varDecl.Name.Lexeme.(string), &variable)

	// TODO(errors)
	if err != nil {
		return err
	}
	return nil
}

func (codegen *codegen) generateParameters(fnValue *values.Function, functionNode *ast.FunctionDecl, fnScope *scope.Scope[values.LLVMValue], paramsTypes []llvm.Type) error {
	for i, paramPtrValue := range fnValue.Fn.Params() {
		paramName := functionNode.Params.Fields[i].Name.Lexeme.(string)
		paramType := paramsTypes[i]
		paramPtr := codegen.builder.CreateAlloca(paramType, ".param")
		codegen.builder.CreateStore(paramPtrValue, paramPtr)
		variable := values.Variable{
			Ty:  paramType,
			Ptr: paramPtr,
		}
		err := fnScope.Insert(paramName, &variable)
		// TODO(errors)
		if err != nil {
			return err
		}
	}
	return nil
}

func (codegen *codegen) generateFunctionCall(scope *scope.Scope[values.LLVMValue], functionCall *ast.FunctionCall) (llvm.Value, error) {
	symbol, err := scope.Lookup(functionCall.Name)
	// TODO(errors)
	if err != nil {
		return llvm.Value{}, err
	}

	function := symbol.(*values.Function)
	args, err := codegen.getExprList(scope, functionCall.Args)
	// TODO(errors)
	if err != nil {
		return llvm.Value{}, err
	}

	// NOTE: do I really need to define the name?
	callName := ""
	return codegen.builder.CreateCall(function.Ty, function.Fn, args, callName), nil
}

func (codegen *codegen) generateExternDecl(external *ast.ExternDecl) error {
	for i := range external.Prototypes {
		err := codegen.generatePrototype(external.Prototypes[i])
		// TODO(errors)
		if err != nil {
			return err
		}
	}
	return nil
}

func (codegen *codegen) generatePrototype(prototype *ast.Proto) error {
	returnType := codegen.getType(prototype.RetType)
	paramsTypes := codegen.getFieldListTypes(prototype.Params)
	functionType := llvm.FunctionType(returnType, paramsTypes, prototype.Params.IsVariadic)
	functionValue := llvm.AddFunction(codegen.module, prototype.Name, functionType)

	function := values.NewFunctionValue(functionValue, functionType, nil)
	err := codegen.universe.Insert(prototype.Name, function)
	return err
}

func (codegen *codegen) getType(ty ast.ExprType) llvm.Type {
	switch exprTy := ty.(type) {
	case *ast.BasicType:
		switch exprTy.Kind {
		case kind.BOOL_TYPE:
			return codegen.context.Int1Type()
		case kind.INT_TYPE, kind.UINT_TYPE:
			// 32 bits or 64 bits
			bitSize := exprTy.Kind.BitSize()
			return codegen.context.IntType(bitSize)
		case kind.I8_TYPE, kind.U8_TYPE:
			return codegen.context.Int8Type()
		case kind.I16_TYPE, kind.U16_TYPE:
			return codegen.context.Int16Type()
		case kind.I32_TYPE, kind.U32_TYPE:
			return codegen.context.Int32Type()
		case kind.I64_TYPE, kind.U64_TYPE:
			return codegen.context.Int64Type()
		case kind.VOID_TYPE:
			return codegen.context.VoidType()
		default:
			log.Fatalf("invalid basic type token: '%s'", exprTy.Kind)
		}
	case *ast.PointerType:
		underlyingExprType := codegen.getType(exprTy.Type)
		// TODO: learn about how to properly define a pointer address space
		return llvm.PointerType(underlyingExprType, 0)
	default:
		log.Fatalf("invalid type: %s", reflect.TypeOf(exprTy))
	}
	// NOTE: this line should be unreachable
	return codegen.context.VoidType()
}

func (codegen *codegen) getFieldListTypes(fields *ast.FieldList) []llvm.Type {
	types := make([]llvm.Type, len(fields.Fields))
	for i := range fields.Fields {
		types[i] = codegen.getType(fields.Fields[i].Type)
	}
	return types
}

func (codegen *codegen) getExprList(parentScope *scope.Scope[values.LLVMValue], expressions []ast.Expr) ([]llvm.Value, error) {
	values := make([]llvm.Value, len(expressions))
	for i := range expressions {
		expr, err := codegen.getExpr(expressions[i], parentScope)
		// TODO(errors)
		if err != nil {
			return nil, err
		}
		values[i] = expr
	}
	return values, nil
}

func (codegen *codegen) getExpr(expr ast.Expr, scope *scope.Scope[values.LLVMValue]) (llvm.Value, error) {
	switch currentExpr := expr.(type) {
	case *ast.LiteralExpr:
		switch ty := currentExpr.Type.(type) {
		case *ast.BasicType:
			// TODO: test this
			integerValue, bitSize, err := codegen.getIntegerValue(currentExpr, ty)
			if bitSize == -1 {
				return llvm.Value{}, fmt.Errorf("%s is not a valid literal basic type", ty.Kind)
			}
			if err != nil {
				return llvm.Value{}, err
			}
			return llvm.ConstInt(codegen.context.IntType(bitSize), integerValue, false), nil

			// switch ty.Kind {
			// case kind.INT_TYPE, kind.UINT_TYPE:
			// 	integerValue, bitSize, err := codegen.getIntegerValue(currentExpr, ty)
			// 	if err != nil {
			// 		return llvm.Value{}, err
			// 	}
			// 	return llvm.ConstInt(codegen.context.IntType(bitSize), integerValue, false), nil
			// case kind.U8_TYPE, kind.I8_TYPE:
			// 	integerValue, _, err := codegen.getIntegerValue(currentExpr, ty)
			// 	if err != nil {
			// 		return llvm.Value{}, err
			// 	}
			// 	return llvm.ConstInt(codegen.context.Int8Type(), integerValue, false), nil
			// case kind.U16_TYPE, kind.I16_TYPE:
			// 	integerValue, _, err := codegen.getIntegerValue(currentExpr, ty)
			// 	if err != nil {
			// 		return llvm.Value{}, err
			// 	}
			// 	return llvm.ConstInt(codegen.context.Int16Type(), integerValue, false), nil
			// case kind.U32_TYPE, kind.I32_TYPE:
			// 	integerValue, _, err := codegen.getIntegerValue(currentExpr, ty)
			// 	if err != nil {
			// 		return llvm.Value{}, err
			// 	}
			// 	return llvm.ConstInt(codegen.context.Int32Type(), integerValue, false), nil
			// case kind.U64_TYPE, kind.I64_TYPE:
			// 	integerValue, _, err := codegen.getIntegerValue(currentExpr, ty)
			// 	if err != nil {
			// 		return llvm.Value{}, err
			// 	}
			// 	return llvm.ConstInt(codegen.context.Int64Type(), integerValue, false), nil
			// case kind.BOOL_TYPE:
			// 	boolLiteral := currentExpr.Value.(string)
			// 	value := 0
			// 	if boolLiteral == "true" {
			// 		value = 1
			// 	}
			// 	return llvm.ConstInt(codegen.context.Int1Type(), uint64(value), false), nil
			// default:
			// 	log.Fatalf("unimplemented literal expr: %s %s", ty.Kind, reflect.TypeOf(expr))
			// }
		case *ast.PointerType:
			switch ptrTy := ty.Type.(type) {
			case *ast.BasicType:
				switch ptrTy.Kind {
				case kind.U8_TYPE:
					stringLiteral := currentExpr.Value.(string)
					// NOTE: huge string literals can affect performance because it
					// creates a new entry on the map
					globalStrLiteral, ok := codegen.strLiterals[stringLiteral]
					if ok {
						return globalStrLiteral, nil
					}
					globalStrPtr := codegen.builder.CreateGlobalStringPtr(stringLiteral, ".str")
					codegen.strLiterals[stringLiteral] = globalStrPtr
					return globalStrPtr, nil
				default:
					log.Fatalf("unimplemented ptr basic type: %s", ptrTy.Kind)
				}
			}
		}
	case *ast.IdExpr:
		varName := currentExpr.Name.Lexeme.(string)
		symbol, err := scope.Lookup(varName)
		// TODO(errors)
		if err != nil {
			return llvm.Value{}, err
		}
		// TODO(errors)
		if symbol == nil {
			log.Fatalf("local not defined: %s", varName)
		}
		localVar := symbol.(*values.Variable)
		loadedVariable := codegen.builder.CreateLoad(localVar.Ty, localVar.Ptr, ".load")
		return loadedVariable, nil
	case *ast.BinaryExpr:
		lhs, err := codegen.getExpr(currentExpr.Left, scope)
		// TODO(errors)
		if err != nil {
			log.Fatalf("can't generate lhs expr: %s", err)
		}
		rhs, err := codegen.getExpr(currentExpr.Right, scope)
		// TODO(errors)
		if err != nil {
			log.Fatalf("can't generate rhs expr: %s", err)
		}
		switch currentExpr.Op {
		case kind.EQUAL_EQUAL:
			// TODO: there a list of IntPredicate, I could map token kind to these
			// for code reability
			// See https://github.com/tinygo-org/go-llvm/blob/master/ir.go#L302
			return codegen.builder.CreateICmp(llvm.IntEQ, lhs, rhs, ".cmpeq"), nil
		case kind.STAR:
			return codegen.builder.CreateMul(lhs, rhs, ".mul"), nil
		case kind.MINUS:
			return codegen.builder.CreateSub(lhs, rhs, ".sub"), nil
		case kind.PLUS:
			return codegen.builder.CreateAdd(lhs, rhs, ".add"), nil
		default:
			log.Fatalf("unimplemented binary operator: %s", currentExpr.Op)
		}
	case *ast.FunctionCall:
		symbol, err := scope.Lookup(currentExpr.Name)
		// TODO(errors)
		if err != nil {
			log.Fatalf("at this point of code generation, every symbol should be located")
		}
		switch sym := symbol.(type) {
		case *values.Function:
			fnCall, err := codegen.generateFunctionCall(scope, currentExpr)
			// TODO(errors)
			if err != nil {
				return llvm.Value{}, nil
			}
			return fnCall, nil
		default:
			log.Fatalf("unimplemented value: %s %s", expr, reflect.TypeOf(sym))
		}
	default:
		log.Fatalf("unimplemented expr: %s", expr)
	}
	// NOTE: this line should be unreachable
	log.Fatalf("REACHING AN UNREACHABLE LINE AT getExpr")
	return llvm.Value{}, nil
}

func (codegen *codegen) getIntegerValue(expr *ast.LiteralExpr, ty *ast.BasicType) (uint64, int, error) {
	bitSize := ty.Kind.BitSize()
	if bitSize == -1 {
		return 0, bitSize, nil
	}
	integerLiteral := expr.Value.(string)
	integerValue, err := strconv.ParseUint(integerLiteral, 10, bitSize)
	return integerValue, bitSize, err
}

func (codegen *codegen) generateCondStmt(parentScope *scope.Scope[values.LLVMValue], function *values.Function, condStmt *ast.CondStmt) error {
	ifBlock := llvm.AddBasicBlock(function.Fn, ".if")
	elseBlock := llvm.AddBasicBlock(function.Fn, ".else")
	endBlock := llvm.AddBasicBlock(function.Fn, ".end")

	ifScope := scope.New(parentScope)
	ifExpr, err := codegen.getExpr(condStmt.IfStmt.Expr, ifScope)
	if err != nil {
		return err
	}

	stoppedOnReturn := false

	codegen.builder.CreateCondBr(ifExpr, ifBlock, elseBlock)
	codegen.builder.SetInsertPointAtEnd(ifBlock)
	stoppedOnReturn, err = codegen.generateBlock(condStmt.IfStmt.Block, ifScope, function)
	// TODO(errors)
	if err != nil {
		return err
	}
	if !stoppedOnReturn {
		codegen.builder.CreateBr(endBlock)
	}

	// TODO: implement elif statements (basically an if inside the else)

	codegen.builder.SetInsertPointAtEnd(elseBlock)
	if condStmt.ElseStmt != nil {
		elseScope := scope.New(parentScope)
		elseStoppedOnReturn, err := codegen.generateBlock(condStmt.ElseStmt.Block, elseScope, function)
		// TODO(errors)
		if err != nil {
			return err
		}
		if !elseStoppedOnReturn {
			codegen.builder.CreateBr(endBlock)
		}
	} else {
		codegen.builder.CreateBr(endBlock)
	}
	codegen.builder.SetInsertPointAtEnd(endBlock)
	return nil
}
