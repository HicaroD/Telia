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
	context llvm.Context
	module  llvm.Module
	builder llvm.Builder

	astNodes    []ast.Node
	universe    *scope.Scope[values.LLVMValue]
	strLiterals map[string]llvm.Value
}

func New(astNodes []ast.Node) *codegen {
	// parent of universe scope is nil
	var nilScope *scope.Scope[values.LLVMValue] = nil
	universe := scope.New(nilScope)

	context := llvm.NewContext()
	// TODO: properly define the module name
	// The name of the module could be file name
	module := context.NewModule("tmpmod")
	builder := context.NewBuilder()

	return &codegen{
		context: context,
		module:  module,
		builder: builder,

		astNodes:    astNodes,
		universe:    universe,
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
	// TODO: properly set the name of generated file
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
	functionValue := llvm.AddFunction(codegen.module, function.Name.Name(), functionType)
	functionBlock := codegen.context.AddBasicBlock(functionValue, "entry")
	fnValue := values.NewFunctionValue(functionValue, functionType, &functionBlock)
	codegen.builder.SetInsertPointAtEnd(functionBlock)

	err := codegen.universe.Insert(function.Name.Name(), fnValue)
	// TODO(errors)
	if err != nil {
		return err
	}

	fnScope := scope.New(codegen.universe)
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

func (codegen *codegen) generateBlock(
	stmts *ast.BlockStmt,
	scope *scope.Scope[values.LLVMValue],
	function *values.Function,
) (bool, error) {
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

func (codegen *codegen) generateStmt(
	stmt ast.Stmt,
	scope *scope.Scope[values.LLVMValue],
	function *values.Function,
) error {
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
	case *ast.MultiVarStmt:
		err := codegen.generateVar(statement, scope)
		// TODO(errors)
		if err != nil {
			return err
		}
	case *ast.FieldAccess:
		err := codegen.generateFieldAccessStmt(statement, scope)
		// TODO(errors)
		if err != nil {
			return err
		}
	case *ast.ForLoop:
		err := codegen.generateForLoop(statement, scope)
		if err != nil {
			return err
		}
	default:
		log.Fatalf("unimplemented block statement: %s", statement)
	}
	return nil
}

func (codegen *codegen) generateReturnStmt(
	ret *ast.ReturnStmt,
	scope *scope.Scope[values.LLVMValue],
) error {
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

func (codegen *codegen) generateVar(
	varDecl *ast.MultiVarStmt,
	scope *scope.Scope[values.LLVMValue],
) error {
	if varDecl.IsDecl {
		for i := range varDecl.Variables {
			err := codegen.generateVarDecl(varDecl.Variables[i], scope)
			if err != nil {
				return err
			}
		}
	} else {
		log.Fatal("unimplemented var redeclaration on codegen")
	}
	return nil
}

func (codegen *codegen) generateVarDecl(
	varDecl *ast.VarStmt,
	scope *scope.Scope[values.LLVMValue],
) error {
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
	// REFACTOR: make it simpler to get the lexeme
	err = scope.Insert(varDecl.Name.Name(), &variable)

	// TODO(errors)
	if err != nil {
		return err
	}
	return nil
}

func (codegen *codegen) generateParameters(
	fnValue *values.Function,
	functionNode *ast.FunctionDecl,
	fnScope *scope.Scope[values.LLVMValue],
	paramsTypes []llvm.Type,
) error {
	for i, paramPtrValue := range fnValue.Fn.Params() {
		// REFACTOR: make it simpler to get the lexeme
		paramName := functionNode.Params.Fields[i].Name.Name()
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

func (codegen *codegen) generateFunctionCall(
	scope *scope.Scope[values.LLVMValue],
	functionCall *ast.FunctionCall,
) (llvm.Value, error) {
	symbol, err := scope.LookupAcrossScopes(functionCall.Name.Name())
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

	return codegen.builder.CreateCall(function.Ty, function.Fn, args, ""), nil
}

func (codegen *codegen) generateExternDecl(external *ast.ExternDecl) error {
	var err error

	scope := scope.New(codegen.universe)

	for i := range external.Prototypes {
		prototype, err := codegen.generatePrototype(external.Prototypes[i])
		// TODO(errors)
		if err != nil {
			return err
		}

		err = scope.Insert(external.Prototypes[i].Name.Name(), prototype)
		// TODO(errors)
		if err != nil {
			return err
		}
	}

	externDecl := values.NewExtern(scope)
	// REFACTOR: make it simpler to get the lexeme
	err = codegen.universe.Insert(external.Name.Name(), externDecl)
	return err
}

func (codegen *codegen) generatePrototype(prototype *ast.Proto) (*values.Function, error) {
	returnTy := codegen.getType(prototype.RetType)
	paramsTypes := codegen.getFieldListTypes(prototype.Params)
	ty := llvm.FunctionType(returnTy, paramsTypes, prototype.Params.IsVariadic)
	protoValue := llvm.AddFunction(codegen.module, prototype.Name.Name(), ty)
	proto := values.NewFunctionValue(protoValue, ty, nil)
	return proto, nil
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

func (codegen *codegen) getExprList(
	parentScope *scope.Scope[values.LLVMValue],
	expressions []ast.Expr,
) ([]llvm.Value, error) {
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

func (codegen *codegen) getExpr(
	expr ast.Expr,
	scope *scope.Scope[values.LLVMValue],
) (llvm.Value, error) {
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
		case *ast.PointerType:
			switch ptrTy := ty.Type.(type) {
			case *ast.BasicType:
				switch ptrTy.Kind {
				case kind.U8_TYPE:
					stringLiteral := currentExpr.Value
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
		// REFACTOR: make it simpler to get the lexeme
		varName := currentExpr.Name.Name()
		symbol, err := scope.LookupAcrossScopes(varName)
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
		// TODO: deal with signed or unsigned operations
		// I'm assuming all unsigned for now
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
		case kind.LESS:
			return codegen.builder.CreateICmp(llvm.IntULT, lhs, rhs, ".cmplt"), nil
		case kind.LESS_EQ:
			return codegen.builder.CreateICmp(llvm.IntULE, lhs, rhs, ".cmple"), nil
		case kind.GREATER:
			return codegen.builder.CreateICmp(llvm.IntUGT, lhs, rhs, ".cmpgt"), nil
		case kind.GREATER_EQ:
			return codegen.builder.CreateICmp(llvm.IntUGE, lhs, rhs, ".cmpge"), nil
		default:
			log.Fatalf("unimplemented binary operator: %s", currentExpr.Op)
		}
	case *ast.FunctionCall:
		symbol, err := scope.LookupAcrossScopes(currentExpr.Name.Name())
		// TODO(errors)
		if err != nil {
			log.Fatalf("at this point of code generation, every symbol should be located")
		}
		switch sym := symbol.(type) {
		case *values.Function:
			fnCall, err := codegen.generateFunctionCall(scope, currentExpr)
			// TODO(errors)
			if err != nil {
				return llvm.Value{}, err
			}
			return fnCall, nil
		default:
			log.Fatalf("unimplemented value: %s %s", expr, reflect.TypeOf(sym))
		}
	case *ast.UnaryExpr:
		switch currentExpr.Op {
		case kind.MINUS:
			expr, err := codegen.getExpr(currentExpr.Value, scope)
			// TODO(errors)
			if err != nil {
				return llvm.Value{}, err
			}
			return codegen.builder.CreateNeg(expr, ".neg"), nil
		default:
			log.Fatalf("unimplemented unary operator: %s", currentExpr.Op)
		}
	default:
		log.Fatalf("unimplemented expr: %s", expr)
	}
	// NOTE: this line should be unreachable
	log.Fatalf("REACHING AN UNREACHABLE LINE AT getExpr")
	return llvm.Value{}, nil
}

func (codegen *codegen) getIntegerValue(
	expr *ast.LiteralExpr,
	ty *ast.BasicType,
) (uint64, int, error) {
	bitSize := ty.Kind.BitSize()
	if bitSize == -1 {
		return 0, bitSize, nil
	}
	integerLiteral := expr.Value
	integerValue, err := strconv.ParseUint(integerLiteral, 10, bitSize)
	return integerValue, bitSize, err
}

func (codegen *codegen) generateCondStmt(
	parentScope *scope.Scope[values.LLVMValue],
	function *values.Function,
	condStmt *ast.CondStmt,
) error {
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
		elseStoppedOnReturn, err := codegen.generateBlock(
			condStmt.ElseStmt.Block,
			elseScope,
			function,
		)
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

func (codegen *codegen) generateFieldAccessStmt(
	fieldAccess *ast.FieldAccess,
	scope *scope.Scope[values.LLVMValue],
) error {
	// REFACTOR: make it simpler to get the lexeme
	idExpr := fieldAccess.Left.(*ast.IdExpr)
	id := idExpr.Name.Name()

	symbol, err := scope.LookupAcrossScopes(id)
	// TODO(errors): should not be a problem, sema needs to caught this errors
	if err != nil {
		log.Fatal("panic: ", err)
	}

	switch left := symbol.(type) {
	case *values.Extern:
		switch right := fieldAccess.Right.(type) {
		case *ast.FunctionCall:
			_, err := codegen.generatePrototypeCall(left, right, scope)
			// TODO(errors)
			if err != nil {
				return err
			}
		default:
			// TODO(errors)
			return fmt.Errorf("unimplemented %s on field access statement", right)
		}
	default:
		// TODO(errors)
		return fmt.Errorf("unimplemented %s on extern", left)
	}
	return nil
}

func (codegen *codegen) generatePrototypeCall(
	extern *values.Extern,
	call *ast.FunctionCall,
	callScope *scope.Scope[values.LLVMValue],
) (llvm.Value, error) {
	prototype, err := extern.Scope.LookupCurrentScope(call.Name.Name())
	// TODO(errors)
	if err != nil {
		return llvm.Value{}, err
	}

	proto := prototype.(*values.Function)
	args, err := codegen.getExprList(callScope, call.Args)
	// TODO(errors)
	if err != nil {
		return llvm.Value{}, err
	}

	return codegen.builder.CreateCall(proto.Ty, proto.Fn, args, ""), nil
}

func (codegen *codegen) generateForLoop(
	forLoop *ast.ForLoop,
	scope *scope.Scope[values.LLVMValue],
) error {
	log.Fatal("unimplemented for loop on code generation")
	return nil
}
