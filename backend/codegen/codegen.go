package codegen

import (
	"bytes"
	"fmt"
	"log"
	"os/exec"
	"path/filepath"
	"reflect"
	"strconv"
	"strings"

	"github.com/HicaroD/Telia/backend/codegen/values"
	"github.com/HicaroD/Telia/frontend/ast"
	"github.com/HicaroD/Telia/frontend/lexer/token"
	"github.com/HicaroD/Telia/scope"
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
			err := codegen.generateFile(file)
			if err != nil {
				return err
			}
		}
	}
	err := codegen.generateExecutable()
	return err
}

func (codegen *codegen) generateFile(file *ast.File) error {
	for _, node := range file.Body {
		switch n := node.(type) {
		case *ast.FunctionDecl:
			err := codegen.generateFnDecl(n)
			// TODO(errors)
			if err != nil {
				return err
			}
		case *ast.ExternDecl:
			err := codegen.generateExternDecl(n)
			// TODO(errors)
			if err != nil {
				return err
			}
		default:
			log.Fatalf("unimplemented: %s\n", reflect.TypeOf(node))
		}
	}

	return nil
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

func (codegen *codegen) generateFnDecl(functionDecl *ast.FunctionDecl) error {
	returnType := codegen.getType(functionDecl.RetType)
	paramsTypes := codegen.getFieldListTypes(functionDecl.Params)
	functionType := llvm.FunctionType(returnType, paramsTypes, functionDecl.Params.IsVariadic)
	functionValue := llvm.AddFunction(codegen.module, functionDecl.Name.Name(), functionType)
	functionBlock := codegen.context.AddBasicBlock(functionValue, "entry")
	fnValue := values.NewFunctionValue(functionValue, functionType, &functionBlock)
	codegen.builder.SetInsertPointAtEnd(functionBlock)

	functionDecl.BackendType = fnValue

	err := codegen.generateParameters(fnValue, functionDecl, paramsTypes)
	// TODO(errors)
	if err != nil {
		return err
	}

	_, err = codegen.generateBlock(functionDecl.Block, functionDecl.Scope, functionDecl, fnValue)
	// TODO(errors)
	if err != nil {
		return err
	}

	return err
}

func (codegen *codegen) generateBlock(
	block *ast.BlockStmt,
	parentScope *scope.Scope[ast.Node],
	functionDecl *ast.FunctionDecl,
	functionLlvm *values.Function,
) (stoppedOnReturn bool, err error) {
	err = nil
	stoppedOnReturn = false

	for _, stmt := range block.Statements {
		err = codegen.generateStmt(stmt, parentScope, functionDecl, functionLlvm)
		if stmt.IsReturn() {
			stoppedOnReturn = true
			return
		}
		if err != nil {
			return
		}
	}
	return
}

func (codegen *codegen) generateStmt(
	stmt ast.Stmt,
	parentScope *scope.Scope[ast.Node],
	functionDecl *ast.FunctionDecl,
	functionLlvm *values.Function,
) error {
	switch statement := stmt.(type) {
	case *ast.FunctionCall:
		_, err := codegen.generateFunctionCall(parentScope, statement)
		return err
	case *ast.ReturnStmt:
		err := codegen.generateReturnStmt(statement, parentScope)
		return err
	case *ast.CondStmt:
		err := codegen.generateCondStmt(parentScope, statement, functionDecl, functionLlvm)
		return err
	case *ast.VarStmt:
		err := codegen.generateVar(statement, parentScope)
		return err
	case *ast.MultiVarStmt:
		err := codegen.generateMultiVar(statement, parentScope)
		return err
	case *ast.FieldAccess:
		err := codegen.generateFieldAccessStmt(statement, parentScope)
		return err
	case *ast.ForLoop:
		err := codegen.generateForLoop(statement, functionDecl, functionLlvm, parentScope)
		return err
	case *ast.WhileLoop:
		err := codegen.generateWhileLoop(statement, functionDecl, functionLlvm, parentScope)
		return err
	default:
		log.Fatalf("unimplemented block statement: %s", statement)
	}
	return nil
}

func (codegen *codegen) generateReturnStmt(
	ret *ast.ReturnStmt,
	scope *scope.Scope[ast.Node],
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

func (codegen *codegen) generateMultiVar(
	varDecl *ast.MultiVarStmt,
	scope *scope.Scope[ast.Node],
) error {
	for _, variable := range varDecl.Variables {
		err := codegen.generateVar(variable, scope)
		if err != nil {
			return err
		}
	}
	return nil
}

func (codegen *codegen) generateVar(
	varStmt *ast.VarStmt,
	scope *scope.Scope[ast.Node],
) error {
	if varStmt.Decl {
		return codegen.generateVarDecl(varStmt, scope)
	} else {
		return codegen.generateVarReassign(varStmt, scope)
	}
}

func (codegen *codegen) generateVarDecl(
	varDecl *ast.VarStmt,
	scope *scope.Scope[ast.Node],
) error {
	varTy := codegen.getType(varDecl.Type)
	varPtr := codegen.builder.CreateAlloca(varTy, ".ptr")
	varExpr, err := codegen.getExpr(varDecl.Value, scope)
	// TODO(errors)
	if err != nil {
		return err
	}
	codegen.builder.CreateStore(varExpr, varPtr)

	variableLlvm := &values.Variable{
		Ty:  varTy,
		Ptr: varPtr,
	}
	varDecl.BackendType = variableLlvm
	return err
}

func (codegen *codegen) generateVarReassign(
	varDecl *ast.VarStmt,
	scope *scope.Scope[ast.Node],
) error {
	expr, err := codegen.getExpr(varDecl.Value, scope)
	if err != nil {
		return err
	}
	symbol, err := scope.LookupAcrossScopes(varDecl.Name.Name())
	if err != nil {
		return err
	}

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
	return nil
}

func (codegen *codegen) generateParameters(
	fnValue *values.Function,
	functionNode *ast.FunctionDecl,
	paramsTypes []llvm.Type,
) error {
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
	return nil
}

func (codegen *codegen) generateFunctionCall(
	functionScope *scope.Scope[ast.Node],
	functionCall *ast.FunctionCall,
) (llvm.Value, error) {
	symbol, err := functionScope.LookupAcrossScopes(functionCall.Name.Name())
	// TODO(errors)
	if err != nil {
		return llvm.Value{}, err
	}

	calledFunction := symbol.(*ast.FunctionDecl)
	calledFunctionLlvm := calledFunction.BackendType.(*values.Function)
	args, err := codegen.getExprList(functionScope, functionCall.Args)
	// TODO(errors)
	if err != nil {
		return llvm.Value{}, err
	}

	return codegen.builder.CreateCall(calledFunctionLlvm.Ty, calledFunctionLlvm.Fn, args, ""), nil
}

func (codegen *codegen) generateExternDecl(external *ast.ExternDecl) error {
	var err error

	for i := range external.Prototypes {
		codegen.generatePrototype(external.Prototypes[i])
	}

	return err
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
	parentScope *scope.Scope[ast.Node],
	expressions []ast.Expr,
) ([]llvm.Value, error) {
	values := make([]llvm.Value, len(expressions))
	for i, expr := range expressions {
		expr, err := codegen.getExpr(expr, parentScope)
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
	scope *scope.Scope[ast.Node],
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
				case token.U8_TYPE:
					str := string(currentExpr.Value)
					// NOTE: huge string literals can affect performance because it
					// creates a new entry on the map
					globalStrLiteral, ok := codegen.strLiterals[str]
					if ok {
						return globalStrLiteral, nil
					}
					globalStrPtr := codegen.builder.CreateGlobalStringPtr(str, ".str")
					codegen.strLiterals[str] = globalStrPtr
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
		case token.EQUAL_EQUAL:
			// TODO: there a list of IntPredicate, I could map token kind to these
			// for code reability
			// See https://github.com/tinygo-org/go-llvm/blob/master/ir.go#L302
			return codegen.builder.CreateICmp(llvm.IntEQ, lhs, rhs, ".cmpeq"), nil
		case token.STAR:
			return codegen.builder.CreateMul(lhs, rhs, ".mul"), nil
		case token.MINUS:
			return codegen.builder.CreateSub(lhs, rhs, ".sub"), nil
		case token.PLUS:
			return codegen.builder.CreateAdd(lhs, rhs, ".add"), nil
		case token.LESS:
			return codegen.builder.CreateICmp(llvm.IntULT, lhs, rhs, ".cmplt"), nil
		case token.LESS_EQ:
			return codegen.builder.CreateICmp(llvm.IntULE, lhs, rhs, ".cmple"), nil
		case token.GREATER:
			return codegen.builder.CreateICmp(llvm.IntUGT, lhs, rhs, ".cmpgt"), nil
		case token.GREATER_EQ:
			return codegen.builder.CreateICmp(llvm.IntUGE, lhs, rhs, ".cmpge"), nil
		default:
			log.Fatalf("unimplemented binary operator: %s", currentExpr.Op)
		}
	case *ast.FunctionCall:
		call, err := codegen.generateFunctionCall(scope, currentExpr)
		// TODO(errors)
		if err != nil {
			return llvm.Value{}, err
		}
		return call, nil
	case *ast.UnaryExpr:
		switch currentExpr.Op {
		case token.MINUS:
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
	intLit := string(expr.Value)
	integerValue, err := strconv.ParseUint(intLit, 10, bitSize)
	return integerValue, bitSize, err
}

func (codegen *codegen) generateCondStmt(
	parentScope *scope.Scope[ast.Node],
	condStmt *ast.CondStmt,
	functionDecl *ast.FunctionDecl,
	functionLlvm *values.Function,
) error {
	ifBlock := llvm.AddBasicBlock(functionLlvm.Fn, ".if")
	elseBlock := llvm.AddBasicBlock(functionLlvm.Fn, ".else")
	endBlock := llvm.AddBasicBlock(functionLlvm.Fn, ".end")

	ifExpr, err := codegen.getExpr(condStmt.IfStmt.Expr, parentScope)
	if err != nil {
		return err
	}

	stoppedOnReturn := false

	codegen.builder.CreateCondBr(ifExpr, ifBlock, elseBlock)
	codegen.builder.SetInsertPointAtEnd(ifBlock)
	stoppedOnReturn, err = codegen.generateBlock(condStmt.IfStmt.Block, parentScope, functionDecl, functionLlvm)
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
			functionDecl,
			functionLlvm,
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
	scope *scope.Scope[ast.Node],
) error {
	// REFACTOR: make it simpler to get the lexeme
	idExpr := fieldAccess.Left.(*ast.IdExpr)
	id := idExpr.Name.Name()

	symbol, err := scope.LookupAcrossScopes(id)
	if err != nil {
		return err
	}

	switch left := symbol.(type) {
	case *ast.ExternDecl:
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
	extern *ast.ExternDecl,
	call *ast.FunctionCall,
	callScope *scope.Scope[ast.Node],
) (llvm.Value, error) {
	prototype, err := extern.Scope.LookupCurrentScope(call.Name.Name())
	// TODO(errors)
	if err != nil {
		return llvm.Value{}, err
	}

	proto := prototype.(*ast.Proto)
	protoLlvm := proto.BackendType.(*values.Function)
	args, err := codegen.getExprList(callScope, call.Args)
	// TODO(errors)
	if err != nil {
		return llvm.Value{}, err
	}

	return codegen.builder.CreateCall(protoLlvm.Ty, protoLlvm.Fn, args, ""), nil
}

func (codegen *codegen) generateForLoop(
	forLoop *ast.ForLoop,
	functionDecl *ast.FunctionDecl,
	functionLlvm *values.Function,
	parentScope *scope.Scope[ast.Node],
) error {
	forScope := scope.New(parentScope)

	forPrepBlock := llvm.AddBasicBlock(functionLlvm.Fn, ".forprep")
	forInitBlock := llvm.AddBasicBlock(functionLlvm.Fn, ".forinit")
	forBodyBlock := llvm.AddBasicBlock(functionLlvm.Fn, ".forbody")
	forUpdateBlock := llvm.AddBasicBlock(functionLlvm.Fn, ".forupdate")
	endBlock := llvm.AddBasicBlock(functionLlvm.Fn, ".forend")

	codegen.builder.CreateBr(forPrepBlock)
	codegen.builder.SetInsertPointAtEnd(forPrepBlock)
	err := codegen.generateStmt(forLoop.Init, forScope, functionDecl, functionLlvm)
	if err != nil {
		return err
	}

	codegen.builder.CreateBr(forInitBlock)
	codegen.builder.SetInsertPointAtEnd(forInitBlock)

	expr, err := codegen.getExpr(forLoop.Cond, forScope)
	if err != nil {
		return err
	}

	codegen.builder.CreateCondBr(expr, forBodyBlock, endBlock)

	codegen.builder.SetInsertPointAtEnd(forBodyBlock)
	_, err = codegen.generateBlock(forLoop.Block, forScope, functionDecl, functionLlvm)
	if err != nil {
		return err
	}

	codegen.builder.CreateBr(forUpdateBlock)
	codegen.builder.SetInsertPointAtEnd(forUpdateBlock)
	err = codegen.generateStmt(forLoop.Update, forScope, functionDecl, functionLlvm)
	if err != nil {
		return err
	}

	codegen.builder.CreateBr(forInitBlock)

	codegen.builder.SetInsertPointAtEnd(endBlock)
	return nil
}

func (codegen *codegen) generateWhileLoop(
	whileLoop *ast.WhileLoop,
	functionDecl *ast.FunctionDecl,
	functionLlvm *values.Function,
	parentScope *scope.Scope[ast.Node],
) error {
	whileScope := scope.New(parentScope)

	whileInitBlock := llvm.AddBasicBlock(functionLlvm.Fn, ".whileinit")
	whileBodyBlock := llvm.AddBasicBlock(functionLlvm.Fn, ".whilebody")
	endBlock := llvm.AddBasicBlock(functionLlvm.Fn, ".whileend")

	codegen.builder.CreateBr(whileInitBlock)
	codegen.builder.SetInsertPointAtEnd(whileInitBlock)
	expr, err := codegen.getExpr(whileLoop.Cond, whileScope)
	if err != nil {
		return err
	}
	codegen.builder.CreateCondBr(expr, whileBodyBlock, endBlock)

	codegen.builder.SetInsertPointAtEnd(whileBodyBlock)
	_, err = codegen.generateBlock(whileLoop.Block, whileScope, functionDecl, functionLlvm)
	if err != nil {
		return err
	}
	codegen.builder.CreateBr(whileInitBlock)

	codegen.builder.SetInsertPointAtEnd(endBlock)
	return nil
}
