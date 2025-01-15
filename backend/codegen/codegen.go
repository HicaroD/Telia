package codegen

import (
	"bytes"
	"log"
	"os/exec"
	"path/filepath"
	"reflect"
	"strconv"
	"strings"

	"github.com/HicaroD/Telia/backend/codegen/values"
	"github.com/HicaroD/Telia/frontend/ast"
	"github.com/HicaroD/Telia/frontend/lexer/token"
	"tinygo.org/x/go-llvm"
)

type codegen struct {
	filename string

	context llvm.Context
	module  llvm.Module
	builder llvm.Builder

	// NOTE: temporary - find a better way of doing this ( preferebly don't do this :) )
	strLiterals map[string]llvm.Value
}

func New(filename string) *codegen {
	context := llvm.NewContext()
	// TODO: properly define the module name
	// The name of the module could be file name
	module := context.NewModule("tmpmod")
	builder := context.NewBuilder()

	return &codegen{
		filename: filename,

		context: context,
		module:  module,
		builder: builder,

		strLiterals: map[string]llvm.Value{},
	}
}

func (codegen *codegen) Generate(program *ast.Program) error {
	for _, module := range program.Body {
		for _, file := range module.Body {
			codegen.generateFile(file)
		}
	}
	err := codegen.generateExecutable()
	return err
}

func (codegen *codegen) generateFile(file *ast.File) {
	for _, node := range file.Body {
		switch n := node.(type) {
		case *ast.FunctionDecl:
			codegen.generateFnDecl(n)
		case *ast.ExternDecl:
			codegen.generateExternDecl(n)
		default:
			log.Fatalf("unimplemented: %s\n", reflect.TypeOf(node))
		}
	}
}

func (codegen *codegen) generateExecutable() error {
	module := codegen.module.String()

	filenameNoExt := strings.TrimSuffix(filepath.Base(codegen.filename), filepath.Ext(codegen.filename))
	cmd := exec.Command("clang", "-O3", "-Wall", "-x", "ir", "-", "-o", filenameNoExt)
	cmd.Stdin = bytes.NewReader([]byte(module))

	var stderr bytes.Buffer
	cmd.Stderr = &stderr
	err := cmd.Run()

	return err
}

func (codegen *codegen) generateFnDecl(functionDecl *ast.FunctionDecl) {
	returnType := codegen.getType(functionDecl.RetType)
	paramsTypes := codegen.getFieldListTypes(functionDecl.Params)
	functionType := llvm.FunctionType(returnType, paramsTypes, functionDecl.Params.IsVariadic)
	functionValue := llvm.AddFunction(codegen.module, functionDecl.Name.Name(), functionType)
	functionBlock := codegen.context.AddBasicBlock(functionValue, "entry")
	fnValue := values.NewFunctionValue(functionValue, functionType, &functionBlock)
	codegen.builder.SetInsertPointAtEnd(functionBlock)

	functionDecl.BackendType = fnValue

	codegen.generateParameters(fnValue, functionDecl, paramsTypes)

	_ = codegen.generateBlock(functionDecl.Block, functionDecl.Scope, functionDecl, fnValue)
}

func (codegen *codegen) generateBlock(
	block *ast.BlockStmt,
	parentScope *ast.Scope,
	functionDecl *ast.FunctionDecl,
	functionLlvm *values.Function,
) (stoppedOnReturn bool) {
	stoppedOnReturn = false

	for _, stmt := range block.Statements {
		codegen.generateStmt(stmt, parentScope, functionDecl, functionLlvm)
		if stmt.IsReturn() {
			stoppedOnReturn = true
			return
		}
	}
	return
}

func (codegen *codegen) generateStmt(
	stmt ast.Stmt,
	parentScope *ast.Scope,
	functionDecl *ast.FunctionDecl,
	functionLlvm *values.Function,
) {
	switch statement := stmt.(type) {
	case *ast.FunctionCall:
		codegen.generateFunctionCall(parentScope, statement)
	case *ast.ReturnStmt:
		codegen.generateReturnStmt(statement, parentScope)
	case *ast.CondStmt:
		codegen.generateCondStmt(parentScope, statement, functionDecl, functionLlvm)
	case *ast.VarStmt:
		codegen.generateVar(statement, parentScope)
	case *ast.MultiVarStmt:
		codegen.generateMultiVar(statement, parentScope)
	case *ast.FieldAccess:
		codegen.generateFieldAccessStmt(statement, parentScope)
	case *ast.ForLoop:
		codegen.generateForLoop(statement, functionDecl, functionLlvm, parentScope)
	case *ast.WhileLoop:
		codegen.generateWhileLoop(statement, functionDecl, functionLlvm, parentScope)
	default:
		log.Fatalf("unimplemented block statement: %s", statement)
	}
}

func (codegen *codegen) generateReturnStmt(
	ret *ast.ReturnStmt,
	scope *ast.Scope,
) {
	if ret.Value.IsVoid() {
		codegen.builder.CreateRetVoid()
		return
	}
	returnValue := codegen.getExpr(ret.Value, scope)
	codegen.builder.CreateRet(returnValue)
}

func (codegen *codegen) generateMultiVar(
	varDecl *ast.MultiVarStmt,
	scope *ast.Scope,
) {
	for _, variable := range varDecl.Variables {
		codegen.generateVar(variable, scope)
	}
}

func (codegen *codegen) generateVar(
	varStmt *ast.VarStmt,
	scope *ast.Scope,
) {
	if varStmt.Decl {
		codegen.generateVarDecl(varStmt, scope)
	} else {
		codegen.generateVarReassign(varStmt, scope)
	}
}

func (codegen *codegen) generateVarDecl(
	varDecl *ast.VarStmt,
	scope *ast.Scope,
) {
	varTy := codegen.getType(varDecl.Type)
	varPtr := codegen.builder.CreateAlloca(varTy, ".ptr")
	varExpr := codegen.getExpr(varDecl.Value, scope)
	codegen.builder.CreateStore(varExpr, varPtr)

	variableLlvm := &values.Variable{
		Ty:  varTy,
		Ptr: varPtr,
	}
	varDecl.BackendType = variableLlvm
}

func (codegen *codegen) generateVarReassign(
	varDecl *ast.VarStmt,
	scope *ast.Scope,
) {
	expr := codegen.getExpr(varDecl.Value, scope)
	symbol, _ := scope.LookupAcrossScopes(varDecl.Name.Name())

	var variable *values.Variable
	switch sy := symbol.(type) {
	case *ast.VarStmt:
		variable = sy.BackendType.(*values.Variable)
	case *ast.Field:
		variable = sy.BackendType.(*values.Variable)
	default:
		log.Fatalf("invalid symbol on generateVarReassign: %v\n", reflect.TypeOf(variable))
	}

	codegen.builder.CreateStore(expr, variable.Ptr)
}

func (codegen *codegen) generateParameters(
	fnValue *values.Function,
	functionNode *ast.FunctionDecl,
	paramsTypes []llvm.Type,
) {
	for i, paramPtrValue := range fnValue.Fn.Params() {
		paramType := paramsTypes[i]
		paramPtr := codegen.builder.CreateAlloca(paramType, ".param")
		codegen.builder.CreateStore(paramPtrValue, paramPtr)
		variable := values.Variable{
			Ty:  paramType,
			Ptr: paramPtr,
		}
		functionNode.Params.Fields[i].BackendType = &variable
	}
}

func (codegen *codegen) generateFunctionCall(
	functionScope *ast.Scope,
	functionCall *ast.FunctionCall,
) llvm.Value {
	symbol, _ := functionScope.LookupAcrossScopes(functionCall.Name.Name())

	calledFunction := symbol.(*ast.FunctionDecl)
	calledFunctionLlvm := calledFunction.BackendType.(*values.Function)
	args := codegen.getExprList(functionScope, functionCall.Args)

	return codegen.builder.CreateCall(calledFunctionLlvm.Ty, calledFunctionLlvm.Fn, args, "")
}

func (codegen *codegen) generateExternDecl(external *ast.ExternDecl) {
	for i := range external.Prototypes {
		codegen.generatePrototype(external.Prototypes[i])
	}
}

func (codegen *codegen) generatePrototype(prototype *ast.Proto) {
	returnTy := codegen.getType(prototype.RetType)
	paramsTypes := codegen.getFieldListTypes(prototype.Params)
	ty := llvm.FunctionType(returnTy, paramsTypes, prototype.Params.IsVariadic)
	protoValue := llvm.AddFunction(codegen.module, prototype.Name.Name(), ty)
	proto := values.NewFunctionValue(protoValue, ty, nil)

	prototype.BackendType = proto
}

func (codegen *codegen) getType(ty ast.ExprType) llvm.Type {
	switch exprTy := ty.(type) {
	case *ast.BasicType:
		switch exprTy.Kind {
		case token.BOOL_TYPE:
			return codegen.context.Int1Type()
		case token.INT_TYPE, token.UINT_TYPE:
			// 32 bits or 64 bits
			bitSize := exprTy.Kind.BitSize()
			return codegen.context.IntType(bitSize)
		case token.I8_TYPE, token.U8_TYPE:
			return codegen.context.Int8Type()
		case token.I16_TYPE, token.U16_TYPE:
			return codegen.context.Int16Type()
		case token.I32_TYPE, token.U32_TYPE:
			return codegen.context.Int32Type()
		case token.I64_TYPE, token.U64_TYPE:
			return codegen.context.Int64Type()
		case token.VOID_TYPE:
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
	parentScope *ast.Scope,
	expressions []ast.Expr,
) []llvm.Value {
	values := make([]llvm.Value, len(expressions))
	for i, expr := range expressions {
		values[i] = codegen.getExpr(expr, parentScope)
	}
	return values
}

func (codegen *codegen) getExpr(
	expr ast.Expr,
	scope *ast.Scope,
) llvm.Value {
	switch currentExpr := expr.(type) {
	case *ast.LiteralExpr:
		switch ty := currentExpr.Type.(type) {
		case *ast.BasicType:
			integerValue, bitSize := codegen.getIntegerValue(currentExpr, ty)
			return llvm.ConstInt(codegen.context.IntType(bitSize), integerValue, false)
		case *ast.PointerType:
			switch ptrTy := ty.Type.(type) {
			case *ast.BasicType:
				switch ptrTy.Kind {
				case token.U8_TYPE:
					str := string(currentExpr.Value)
					// NOTE: huge string literals can affect performance because it
					// creates a new entry on the map
					globalStrLiteral, ok := codegen.strLiterals[str]
					if ok {
						return globalStrLiteral
					}
					globalStrPtr := codegen.builder.CreateGlobalStringPtr(str, ".str")
					codegen.strLiterals[str] = globalStrPtr
					return globalStrPtr
				default:
					log.Fatalf("unimplemented ptr basic type: %s", ptrTy.Kind)
				}
			}
		}
	case *ast.IdExpr:
		// REFACTOR: make it simpler to get the lexeme
		varName := currentExpr.Name.Name()
		symbol, _ := scope.LookupAcrossScopes(varName)

		var localVar *values.Variable

		switch symbol.(type) {
		case *ast.VarStmt:
			variable := symbol.(*ast.VarStmt)
			localVar = variable.BackendType.(*values.Variable)
		case *ast.Field:
			variable := symbol.(*ast.Field)
			localVar = variable.BackendType.(*values.Variable)
		}

		loadedVariable := codegen.builder.CreateLoad(localVar.Ty, localVar.Ptr, ".load")
		return loadedVariable
	case *ast.BinaryExpr:
		lhs := codegen.getExpr(currentExpr.Left, scope)
		rhs := codegen.getExpr(currentExpr.Right, scope)

		switch currentExpr.Op {
		// TODO: deal with signed or unsigned operations
		// I'm assuming all unsigned for now
		case token.EQUAL_EQUAL:
			// TODO: there a list of IntPredicate, I could map token kind to these
			// for code reability
			// See https://github.com/tinygo-org/go-llvm/blob/master/ir.go#L302
			return codegen.builder.CreateICmp(llvm.IntEQ, lhs, rhs, ".cmpeq")
		case token.STAR:
			return codegen.builder.CreateMul(lhs, rhs, ".mul")
		case token.MINUS:
			return codegen.builder.CreateSub(lhs, rhs, ".sub")
		case token.PLUS:
			return codegen.builder.CreateAdd(lhs, rhs, ".add")
		case token.LESS:
			return codegen.builder.CreateICmp(llvm.IntULT, lhs, rhs, ".cmplt")
		case token.LESS_EQ:
			return codegen.builder.CreateICmp(llvm.IntULE, lhs, rhs, ".cmple")
		case token.GREATER:
			return codegen.builder.CreateICmp(llvm.IntUGT, lhs, rhs, ".cmpgt")
		case token.GREATER_EQ:
			return codegen.builder.CreateICmp(llvm.IntUGE, lhs, rhs, ".cmpge")
		default:
			log.Fatalf("unimplemented binary operator: %s", currentExpr.Op)
		}
	case *ast.FunctionCall:
		call := codegen.generateFunctionCall(scope, currentExpr)
		return call
	case *ast.UnaryExpr:
		switch currentExpr.Op {
		case token.MINUS:
			expr := codegen.getExpr(currentExpr.Value, scope)
			return codegen.builder.CreateNeg(expr, ".neg")
		default:
			log.Fatalf("unimplemented unary operator: %s", currentExpr.Op)
		}
	default:
		log.Fatalf("unimplemented expr: %s", expr)
	}

	// NOTE: this line should be unreachable
	log.Fatalf("REACHING AN UNREACHABLE LINE AT getExpr")
	return llvm.Value{}
}

func (codegen *codegen) getIntegerValue(
	expr *ast.LiteralExpr,
	ty *ast.BasicType,
) (uint64, int) {
	bitSize := ty.Kind.BitSize()
	intLit := string(expr.Value)
	integerValue, err := strconv.ParseUint(intLit, 10, bitSize)
	// TODO: make sure to validate this during semantic analysis stage!
	if err != nil {
		log.Fatal(err)
	}
	return integerValue, bitSize
}

func (codegen *codegen) generateCondStmt(
	parentScope *ast.Scope,
	condStmt *ast.CondStmt,
	functionDecl *ast.FunctionDecl,
	functionLlvm *values.Function,
) {
	ifBlock := llvm.AddBasicBlock(functionLlvm.Fn, ".if")
	elseBlock := llvm.AddBasicBlock(functionLlvm.Fn, ".else")
	endBlock := llvm.AddBasicBlock(functionLlvm.Fn, ".end")

	ifExpr := codegen.getExpr(condStmt.IfStmt.Expr, parentScope)
	codegen.builder.CreateCondBr(ifExpr, ifBlock, elseBlock)
	codegen.builder.SetInsertPointAtEnd(ifBlock)

	stoppedOnReturn := codegen.generateBlock(condStmt.IfStmt.Block, parentScope, functionDecl, functionLlvm)
	if !stoppedOnReturn {
		codegen.builder.CreateBr(endBlock)
	}

	// TODO: implement elif statements (basically an if inside the else)

	codegen.builder.SetInsertPointAtEnd(elseBlock)
	if condStmt.ElseStmt != nil {
		elseScope := ast.NewScope(parentScope)
		elseStoppedOnReturn := codegen.generateBlock(
			condStmt.ElseStmt.Block,
			elseScope,
			functionDecl,
			functionLlvm,
		)
		if !elseStoppedOnReturn {
			codegen.builder.CreateBr(endBlock)
		}
	} else {
		codegen.builder.CreateBr(endBlock)
	}
	codegen.builder.SetInsertPointAtEnd(endBlock)
}

func (codegen *codegen) generateFieldAccessStmt(
	fieldAccess *ast.FieldAccess,
	scope *ast.Scope,
) {
	idExpr := fieldAccess.Left.(*ast.IdExpr)
	id := idExpr.Name.Name()

	symbol, _ := scope.LookupAcrossScopes(id)

	switch left := symbol.(type) {
	case *ast.ExternDecl:
		switch right := fieldAccess.Right.(type) {
		case *ast.FunctionCall:
			codegen.generatePrototypeCall(left, right, scope)
		default:
			// TODO(errors)
			log.Fatalf("unimplemented %s on field access statement", right)
		}
	default:
		// TODO(errors)
		log.Fatalf("unimplemented %s on extern", left)
	}
}

func (codegen *codegen) generatePrototypeCall(
	extern *ast.ExternDecl,
	call *ast.FunctionCall,
	callScope *ast.Scope,
) llvm.Value {
	prototype, _ := extern.Scope.LookupCurrentScope(call.Name.Name())
	proto := prototype.(*ast.Proto)
	protoLlvm := proto.BackendType.(*values.Function)
	args := codegen.getExprList(callScope, call.Args)

	return codegen.builder.CreateCall(protoLlvm.Ty, protoLlvm.Fn, args, "")
}

func (codegen *codegen) generateForLoop(
	forLoop *ast.ForLoop,
	functionDecl *ast.FunctionDecl,
	functionLlvm *values.Function,
	parentScope *ast.Scope,
) {
	forScope := ast.NewScope(parentScope)

	forPrepBlock := llvm.AddBasicBlock(functionLlvm.Fn, ".forprep")
	forInitBlock := llvm.AddBasicBlock(functionLlvm.Fn, ".forinit")
	forBodyBlock := llvm.AddBasicBlock(functionLlvm.Fn, ".forbody")
	forUpdateBlock := llvm.AddBasicBlock(functionLlvm.Fn, ".forupdate")
	endBlock := llvm.AddBasicBlock(functionLlvm.Fn, ".forend")

	codegen.builder.CreateBr(forPrepBlock)
	codegen.builder.SetInsertPointAtEnd(forPrepBlock)
	codegen.generateStmt(forLoop.Init, forScope, functionDecl, functionLlvm)

	codegen.builder.CreateBr(forInitBlock)
	codegen.builder.SetInsertPointAtEnd(forInitBlock)

	expr := codegen.getExpr(forLoop.Cond, forScope)

	codegen.builder.CreateCondBr(expr, forBodyBlock, endBlock)

	codegen.builder.SetInsertPointAtEnd(forBodyBlock)
	codegen.generateBlock(forLoop.Block, forScope, functionDecl, functionLlvm)

	codegen.builder.CreateBr(forUpdateBlock)
	codegen.builder.SetInsertPointAtEnd(forUpdateBlock)
	codegen.generateStmt(forLoop.Update, forScope, functionDecl, functionLlvm)

	codegen.builder.CreateBr(forInitBlock)

	codegen.builder.SetInsertPointAtEnd(endBlock)
}

func (codegen *codegen) generateWhileLoop(
	whileLoop *ast.WhileLoop,
	functionDecl *ast.FunctionDecl,
	functionLlvm *values.Function,
	parentScope *ast.Scope,
) {
	whileScope := ast.NewScope(parentScope)

	whileInitBlock := llvm.AddBasicBlock(functionLlvm.Fn, ".whileinit")
	whileBodyBlock := llvm.AddBasicBlock(functionLlvm.Fn, ".whilebody")
	endBlock := llvm.AddBasicBlock(functionLlvm.Fn, ".whileend")

	codegen.builder.CreateBr(whileInitBlock)
	codegen.builder.SetInsertPointAtEnd(whileInitBlock)
	expr := codegen.getExpr(whileLoop.Cond, whileScope)
	codegen.builder.CreateCondBr(expr, whileBodyBlock, endBlock)

	codegen.builder.SetInsertPointAtEnd(whileBodyBlock)
	codegen.generateBlock(whileLoop.Block, whileScope, functionDecl, functionLlvm)
	codegen.builder.CreateBr(whileInitBlock)

	codegen.builder.SetInsertPointAtEnd(endBlock)
}
