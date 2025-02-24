package llvm

import (
	"errors"
	"fmt"
	"io/fs"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"reflect"
	"strconv"
	"strings"

	"github.com/HicaroD/Telia/internal/ast"
	"github.com/HicaroD/Telia/internal/config"
	"github.com/HicaroD/Telia/internal/lexer/token"
	"tinygo.org/x/go-llvm"
)

type llvmCodegen struct {
	path    string
	program *ast.Program

	context llvm.Context
	module  llvm.Module
	builder llvm.Builder
}

func NewCG(parentDirName, path string, program *ast.Program) *llvmCodegen {
	context := llvm.NewContext()
	module := context.NewModule(parentDirName)
	builder := context.NewBuilder()

	defaultTargetTriple := llvm.DefaultTargetTriple()
	module.SetTarget(defaultTargetTriple)

	return &llvmCodegen{
		path:    path,
		program: program,

		context: context,
		module:  module,
		builder: builder,
	}
}

func (c *llvmCodegen) Generate(buildType config.BuildType) error {
	c.generateModule(c.program.Root)
	err := c.generateExe(buildType)
	return err
}

func (c *llvmCodegen) generateModule(module *ast.Package) {
	for _, file := range module.Files {
		c.generateFile(file)
	}
	for _, module := range module.Packages {
		c.generateModule(module)
	}
}

func (c *llvmCodegen) generateFile(file *ast.File) {
	for _, node := range file.Body {
		switch node.Kind {
		case ast.KIND_FN_DECL:
			c.generateFnDecl(node.Node.(*ast.FnDecl))
		case ast.KIND_EXTERN_DECL:
			c.generateExternDecl(node.Node.(*ast.ExternDecl))
		case ast.KIND_PKG_DECL:
			continue
		case ast.KIND_USE_DECL:
			continue
		default:
			log.Fatalf("unimplemented: %s\n", reflect.TypeOf(node))
		}
	}
}

func (c *llvmCodegen) exists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}
	if errors.Is(err, fs.ErrNotExist) {
		return false, nil
	}
	return false, err

}

func (c *llvmCodegen) generateExe(buildType config.BuildType) error {
	filenameNoExt := strings.TrimSuffix(filepath.Base(c.path), filepath.Ext(c.path))

	dirName := "__build__"
	exists, err := c.exists(dirName)
	if err != nil {
		return err
	}
	if exists {
		err := os.RemoveAll(dirName)
		if err != nil {
			return err
		}
	}

	err = os.Mkdir(dirName, os.ModePerm)
	if err != nil {
		return err
	}

	irFileName := filepath.Join(dirName, c.program.Root.Name)
	irFilepath := irFileName + ".ll"
	optimizedIrFilepath := irFileName + "_optimized.ll"

	irFile, err := os.Create(irFilepath)
	// TODO(errors)
	if err != nil {
		return err
	}
	defer irFile.Close()

	module := c.module.String()
	_, err = irFile.WriteString(module)
	// TODO(errors)
	if err != nil {
		return err
	}

	optLevel := ""
	compilerFlags := ""

	switch buildType {
	case config.RELEASE:
		optLevel = "-O3"
		compilerFlags = "-Wl,-s"
	case config.DEBUG:
		optLevel = "-O0"
	default:
		panic("invalid build type: " + buildType.String())
	}

	cmd := exec.Command("opt", optLevel, "-o", optimizedIrFilepath, irFilepath)
	err = cmd.Run()
	// TODO(errors)
	if err != nil {
		return err
	}

	cmd = exec.Command("clang-18", compilerFlags, optLevel, "-o", filenameNoExt, optimizedIrFilepath)
	err = cmd.Run()
	// TODO(errors)
	if err != nil {
		return err
	}

	if config.DEBUG_MODE {
		fmt.Println("[DEBUG MODE] keeping __build__ directory")
	} else {
		err = os.RemoveAll(dirName)
		if err != nil {
			return err
		}
	}

	return nil
}

func (c *llvmCodegen) generateFnDecl(functionDecl *ast.FnDecl) {
	returnType := c.getType(functionDecl.RetType)
	paramsTypes := c.getFieldListTypes(functionDecl.Params)
	functionType := llvm.FunctionType(returnType, paramsTypes, functionDecl.Params.IsVariadic)
	functionValue := llvm.AddFunction(c.module, functionDecl.Name.Name(), functionType)
	functionBlock := c.context.AddBasicBlock(functionValue, "entry")
	c.builder.SetInsertPointAtEnd(functionBlock)

	fnValue := NewFunctionValue(functionValue, functionType, &functionBlock)
	functionDecl.BackendType = fnValue

	c.generateParameters(fnValue, functionDecl, paramsTypes)
	_ = c.generateBlock(functionDecl.Block, fnValue, functionDecl.Scope)
}

func (c *llvmCodegen) generateBlock(
	block *ast.BlockStmt,
	function *Function,
	parentScope *ast.Scope,
) (stoppedOnReturn bool) {
	stoppedOnReturn = false

	for _, stmt := range block.Statements {
		c.generateStmt(stmt, function, parentScope)
		if stmt.IsReturn() {
			stoppedOnReturn = true
			return
		}
	}
	return
}

func (c *llvmCodegen) generateStmt(
	stmt *ast.Node,
	function *Function,
	parentScope *ast.Scope,
) {
	switch stmt.Kind {
	case ast.KIND_FN_CALL:
		c.generateFunctionCall(stmt.Node.(*ast.FnCall), parentScope)
	case ast.KIND_RETURN_STMT:
		c.generateReturnStmt(stmt.Node.(*ast.ReturnStmt), parentScope)
	case ast.KIND_VAR_STMT:
		c.generateVar(stmt.Node.(*ast.VarStmt), parentScope)
	case ast.KIND_MULTI_VAR_STMT:
		c.generateMultiVar(stmt.Node.(*ast.MultiVarStmt), parentScope)
	case ast.KIND_FIELD_ACCESS:
		c.generateFieldAccessStmt(stmt.Node.(*ast.FieldAccess), parentScope)
	case ast.KIND_COND_STMT:
		c.generateCondStmt(stmt.Node.(*ast.CondStmt), function)
	case ast.KIND_FOR_LOOP_STMT:
		c.generateForLoop(stmt.Node.(*ast.ForLoop), function)
	case ast.KIND_WHILE_LOOP_STMT:
		c.generateWhileLoop(stmt.Node.(*ast.WhileLoop), function)
	default:
		log.Fatalf("unimplemented block statement: %s\n", stmt)
	}
}

func (c *llvmCodegen) generateReturnStmt(
	ret *ast.ReturnStmt,
	scope *ast.Scope,
) {
	if ret.Value.IsVoid() {
		c.builder.CreateRetVoid()
		return
	}
	returnValue := c.getExpr(ret.Value, scope)
	// TODO: learn more about noundef for return values
	c.builder.CreateRet(returnValue)
}

func (c *llvmCodegen) generateMultiVar(
	varDecl *ast.MultiVarStmt,
	scope *ast.Scope,
) {
	for _, variable := range varDecl.Variables {
		c.generateVar(variable.Node.(*ast.VarStmt), scope)
	}
}

func (c *llvmCodegen) generateVar(
	varStmt *ast.VarStmt,
	scope *ast.Scope,
) {
	if varStmt.Decl {
		c.generateVarDecl(varStmt, scope)
	} else {
		c.generateVarReassign(varStmt, scope)
	}
}

func (c *llvmCodegen) generateVarDecl(
	varDecl *ast.VarStmt,
	scope *ast.Scope,
) {
	ty := c.getType(varDecl.Type)
	ptr := c.builder.CreateAlloca(ty, ".ptr")
	expr := c.getExpr(varDecl.Value, scope)
	c.builder.CreateStore(expr, ptr)

	variableLlvm := &Variable{
		Ty:  ty,
		Ptr: ptr,
	}
	varDecl.BackendType = variableLlvm
}

func (c *llvmCodegen) generateVarReassign(
	varDecl *ast.VarStmt,
	scope *ast.Scope,
) {
	expr := c.getExpr(varDecl.Value, scope)
	symbol, _ := scope.LookupAcrossScopes(varDecl.Name.Name())

	var variable *Variable

	switch symbol.Kind {
	case ast.KIND_VAR_STMT:
		node := symbol.Node.(*ast.VarStmt)
		variable = node.BackendType.(*Variable)
	case ast.KIND_FIELD:
		node := symbol.Node.(*ast.Field)
		variable = node.BackendType.(*Variable)
	default:
		log.Fatalf("invalid symbol on generateVarReassign: %v\n", reflect.TypeOf(variable))
	}

	c.builder.CreateStore(expr, variable.Ptr)
}

func (c *llvmCodegen) generateParameters(
	fnValue *Function,
	functionNode *ast.FnDecl,
	paramsTypes []llvm.Type,
) {
	// TODO: learn more about noundef parameter attribute
	for i, paramPtrValue := range fnValue.Fn.Params() {
		paramType := paramsTypes[i]
		paramPtr := c.builder.CreateAlloca(paramType, ".param")
		c.builder.CreateStore(paramPtrValue, paramPtr)
		variable := Variable{
			Ty:  paramType,
			Ptr: paramPtr,
		}
		functionNode.Params.Fields[i].BackendType = &variable
	}
}

func (c *llvmCodegen) generateFunctionCall(
	functionCall *ast.FnCall,
	functionScope *ast.Scope,
) llvm.Value {
	symbol, _ := functionScope.LookupAcrossScopes(functionCall.Name.Name())

	calledFunction := symbol.Node.(*ast.FnDecl)
	calledFunctionLlvm := calledFunction.BackendType.(*Function)
	args := c.getExprList(functionScope, functionCall.Args)

	return c.builder.CreateCall(calledFunctionLlvm.Ty, calledFunctionLlvm.Fn, args, "")
}

func (c *llvmCodegen) generateExternDecl(external *ast.ExternDecl) {
	for i := range external.Prototypes {
		c.generatePrototype(external.Attributes, external.Prototypes[i])
	}
}

func (c *llvmCodegen) generatePrototype(attributes *ast.ExternAttrs, prototype *ast.Proto) {
	returnTy := c.getType(prototype.RetType)
	paramsTypes := c.getFieldListTypes(prototype.Params)
	ty := llvm.FunctionType(returnTy, paramsTypes, prototype.Params.IsVariadic)
	protoValue := llvm.AddFunction(c.module, prototype.Name.Name(), ty)

	if attributes != nil {
		protoValue.SetFunctionCallConv(c.getCallingConvention(attributes.DefaultCallingConvention))
	}

	if prototype.Attributes != nil {
		protoValue.SetLinkage(c.getFunctionLinkage(prototype.Attributes.Linkage))
	}

	proto := NewFunctionValue(protoValue, ty, nil)
	prototype.BackendType = proto
}

func (c *llvmCodegen) getFunctionLinkage(linkage string) llvm.Linkage {
	switch linkage {
	case "internal":
		return llvm.InternalLinkage
	// NOTE: there are several types of weak linkages
	// Make sure to choose the perfect match
	case "weak":
		return llvm.WeakAnyLinkage
	// NOTE: there are several types of link once linkage
	// Makesure to choose the perfect match
	case "link_once":
		return llvm.LinkOnceAnyLinkage
	case "external":
		fallthrough
	default:
		return llvm.ExternalLinkage
	}
}

func (c *llvmCodegen) getCallingConvention(callingConvention string) llvm.CallConv {
	// TODO: define other types of calling conventions
	switch callingConvention {
	case "fast":
		return llvm.FastCallConv
	case "cold":
		return llvm.ColdCallConv
	case "c":
		fallthrough
	default:
		return llvm.CCallConv
	}
}

func (c *llvmCodegen) getType(ty *ast.ExprType) llvm.Type {
	switch ty.Kind {
	case ast.EXPR_TYPE_BASIC:
		b := ty.T.(*ast.BasicType)
		switch b.Kind {
		case token.BOOL_TYPE:
			return c.context.Int1Type()
		case token.UINT_TYPE, token.INT_TYPE, token.UNTYPED_INT:
			// 32 bits or 64 bits depends on the architecture
			bitSize := b.Kind.BitSize()
			return c.context.IntType(bitSize)
		case token.I8_TYPE, token.U8_TYPE:
			return c.context.Int8Type()
		case token.I16_TYPE, token.U16_TYPE:
			return c.context.Int16Type()
		case token.I32_TYPE, token.U32_TYPE:
			return c.context.Int32Type()
		case token.I64_TYPE, token.U64_TYPE:
			return c.context.Int64Type()
		case token.VOID_TYPE:
			return c.context.VoidType()
		// NOTE: token.STRING_TYPE in the same place as token.CSTRING_TYPE is a placeholder
		// NOTE: string type is actually more complex than that
		case token.STRING_TYPE, token.CSTRING_TYPE:
			u8 := ast.NewBasicType(token.U8_TYPE)
			u8Type := c.getType(u8)
			return c.getPtrType(u8Type)
		default:
			log.Fatalf("invalid basic type token: '%s'", b.Kind)
		}
	}
	return c.context.VoidType()
}

func (c *llvmCodegen) getPtrType(ty llvm.Type) llvm.Type {
	return llvm.PointerType(ty, 0)
}

func (c *llvmCodegen) getFieldListTypes(fields *ast.FieldList) []llvm.Type {
	types := make([]llvm.Type, len(fields.Fields))
	for i := range fields.Fields {
		types[i] = c.getType(fields.Fields[i].Type)
	}
	return types
}

func (c *llvmCodegen) getExprList(
	parentScope *ast.Scope,
	expressions []*ast.Node,
) []llvm.Value {
	values := make([]llvm.Value, len(expressions))
	for i, expr := range expressions {
		values[i] = c.getExpr(expr, parentScope)
	}
	return values
}

func (c *llvmCodegen) getExpr(
	expr *ast.Node,
	scope *ast.Scope,
) llvm.Value {
	switch expr.Kind {
	case ast.KIND_LITERAl_EXPR:
		lit := expr.Node.(*ast.LiteralExpr)
		switch lit.Type.Kind {
		case ast.EXPR_TYPE_BASIC:
			basic := lit.Type.T.(*ast.BasicType)

			switch {
			case basic.IsIntegerType():
				integerValue, bitSize := c.getIntegerValue(lit, basic)
				return llvm.ConstInt(c.context.IntType(bitSize), integerValue, false)
			case basic.Kind == token.CSTRING_TYPE:
				strlen := len(lit.Value) + 1
				arrTy := llvm.ArrayType(c.context.Int8Type(), strlen)
				arr := llvm.ConstArray(arrTy, c.llvmConstInt8s(lit.Value, strlen))

				globalVal := llvm.AddGlobal(c.module, arrTy, ".str")
				globalVal.SetInitializer(arr)
				globalVal.SetLinkage(llvm.PrivateLinkage)
				globalVal.SetGlobalConstant(true)
				globalVal.SetAlignment(1)
				globalVal.SetUnnamedAddr(true)

				zero := llvm.ConstInt(c.context.Int32Type(), 0, false)
				indices := []llvm.Value{zero, zero}
				ptr := llvm.ConstInBoundsGEP(arrTy, globalVal, indices)

				return ptr
			default:
				log.Fatalf("unimplemented basic type: %s\n", lit)
				return llvm.Value{}
			}
		case ast.EXPR_TYPE_POINTER:
			panic("unimplemented pointer type")
		default:
			panic("unimplemented literal type")
		}
	case ast.KIND_ID_EXPR:
		id := expr.Node.(*ast.IdExpr)

		varName := id.Name.Name()
		symbol, _ := scope.LookupAcrossScopes(varName)

		var localVar *Variable

		switch symbol.Kind {
		case ast.KIND_VAR_STMT:
			variable := symbol.Node.(*ast.VarStmt)
			localVar = variable.BackendType.(*Variable)
		case ast.KIND_FIELD:
			variable := symbol.Node.(*ast.Field)
			localVar = variable.BackendType.(*Variable)
		}

		loadedVariable := c.builder.CreateLoad(localVar.Ty, localVar.Ptr, ".load")
		return loadedVariable
	case ast.KIND_BINARY_EXPR:
		binary := expr.Node.(*ast.BinaryExpr)

		lhs := c.getExpr(binary.Left, scope)
		rhs := c.getExpr(binary.Right, scope)

		switch binary.Op {
		// TODO: deal with signed or unsigned operations
		// I'm assuming all unsigned for now
		case token.EQUAL_EQUAL:
			// TODO: there a list of IntPredicate, I could map token kind to these
			// for code reability
			// See https://github.com/tinygo-org/go-llvm/blob/master/ir.go#L302
			return c.builder.CreateICmp(llvm.IntEQ, lhs, rhs, ".cmpeq")
		case token.BANG_EQUAL:
			return c.builder.CreateICmp(llvm.IntNE, lhs, rhs, ".cmpneq")
		case token.STAR:
			return c.builder.CreateMul(lhs, rhs, ".mul")
		case token.MINUS:
			return c.builder.CreateSub(lhs, rhs, ".sub")
		case token.PLUS:
			return c.builder.CreateAdd(lhs, rhs, ".add")
		case token.LESS:
			return c.builder.CreateICmp(llvm.IntULT, lhs, rhs, ".cmplt")
		case token.LESS_EQ:
			return c.builder.CreateICmp(llvm.IntULE, lhs, rhs, ".cmple")
		case token.GREATER:
			return c.builder.CreateICmp(llvm.IntUGT, lhs, rhs, ".cmpgt")
		case token.GREATER_EQ:
			return c.builder.CreateICmp(llvm.IntUGE, lhs, rhs, ".cmpge")
		case token.SLASH:
			return c.builder.CreateExactSDiv(lhs, rhs, ".exactdiv")
		default:
			log.Fatalf("unimplemented binary operator: %s", binary.Op)
			return llvm.Value{}
		}
	case ast.KIND_FN_CALL:
		fnCall := expr.Node.(*ast.FnCall)
		call := c.generateFunctionCall(fnCall, scope)
		return call
	case ast.KIND_UNARY_EXPR:
		unary := expr.Node.(*ast.UnaryExpr)

		switch unary.Op {
		case token.MINUS:
			expr := c.getExpr(unary.Value, scope)
			return c.builder.CreateNeg(expr, ".neg")
		default:
			log.Fatalf("unimplemented unary operator: %s", unary.Op)
			return llvm.Value{}
		}
	case ast.KIND_FIELD_ACCESS:
		fieldAccess := expr.Node.(*ast.FieldAccess)
		value := c.generateFieldAccessStmt(fieldAccess, scope)
		return value
	default:
		log.Fatalf("unimplemented expr: %s", reflect.TypeOf(expr))
		return llvm.Value{}
	}
}

func (c *llvmCodegen) llvmConstInt8s(data []byte, length int) []llvm.Value {
	out := make([]llvm.Value, length)
	for i, b := range data {
		out[i] = llvm.ConstInt(c.context.Int8Type(), uint64(b), false)
	}
	// c-string null terminated string
	out[length-1] = llvm.ConstInt(c.context.Int8Type(), uint64(0), false)
	return out
}

func (c *llvmCodegen) getIntegerValue(
	expr *ast.LiteralExpr,
	ty *ast.BasicType,
) (uint64, int) {
	bitSize := ty.Kind.BitSize()
	intLit := string(expr.Value)
	switch intLit {
	case "true":
		intLit = "1"
	case "false":
		intLit = "0"
	}
	val, _ := strconv.ParseUint(intLit, 10, bitSize)
	return val, bitSize
}

func (c *llvmCodegen) generateCondStmt(
	condStmt *ast.CondStmt,
	function *Function,
) {
	ifBlock := llvm.AddBasicBlock(function.Fn, ".if")
	elseBlock := llvm.AddBasicBlock(function.Fn, ".else")
	endBlock := llvm.AddBasicBlock(function.Fn, ".end")

	// If
	ifExpr := c.getExpr(condStmt.IfStmt.Expr, condStmt.IfStmt.Scope)
	c.builder.CreateCondBr(ifExpr, ifBlock, elseBlock)
	c.builder.SetInsertPointAtEnd(ifBlock)
	stoppedOnReturn := c.generateBlock(condStmt.IfStmt.Block, function, condStmt.IfStmt.Scope)
	if !stoppedOnReturn {
		c.builder.CreateBr(endBlock)
	}

	// TODO: implement elif statements (at the end, they are just ifs and elses)

	// Else
	c.builder.SetInsertPointAtEnd(elseBlock)
	if condStmt.ElseStmt != nil {
		elseStoppedOnReturn := c.generateBlock(
			condStmt.ElseStmt.Block,
			function,
			condStmt.ElseStmt.Scope,
		)
		if !elseStoppedOnReturn {
			c.builder.CreateBr(endBlock)
		}
	} else {
		c.builder.CreateBr(endBlock)
	}
	c.builder.SetInsertPointAtEnd(endBlock)
}

func (c *llvmCodegen) generateFieldAccessStmt(
	fieldAccess *ast.FieldAccess,
	scope *ast.Scope,
) llvm.Value {
	idExpr := fieldAccess.Left.Node.(*ast.IdExpr)
	id := idExpr.Name.Name()

	symbol, _ := scope.LookupAcrossScopes(id)

	left := symbol.Node.(*ast.ExternDecl)
	right := fieldAccess.Right.Node.(*ast.FnCall)
	return c.generatePrototypeCall(left, right, scope)
}

func (c *llvmCodegen) generatePrototypeCall(
	extern *ast.ExternDecl,
	call *ast.FnCall,
	callScope *ast.Scope,
) llvm.Value {
	prototype, _ := extern.Scope.LookupCurrentScope(call.Name.Name())
	proto := prototype.Node.(*ast.Proto)
	protoLlvm := proto.BackendType.(*Function)
	args := c.getExprList(callScope, call.Args)

	return c.builder.CreateCall(protoLlvm.Ty, protoLlvm.Fn, args, "")
}

func (c *llvmCodegen) generateForLoop(
	forLoop *ast.ForLoop,
	function *Function,
) {
	forPrepBlock := llvm.AddBasicBlock(function.Fn, ".forprep")
	forInitBlock := llvm.AddBasicBlock(function.Fn, ".forinit")
	forBodyBlock := llvm.AddBasicBlock(function.Fn, ".forbody")
	forUpdateBlock := llvm.AddBasicBlock(function.Fn, ".forupdate")
	endBlock := llvm.AddBasicBlock(function.Fn, ".forend")

	c.builder.CreateBr(forPrepBlock)
	c.builder.SetInsertPointAtEnd(forPrepBlock)
	c.generateStmt(forLoop.Init, function, forLoop.Scope)

	c.builder.CreateBr(forInitBlock)
	c.builder.SetInsertPointAtEnd(forInitBlock)

	expr := c.getExpr(forLoop.Cond, forLoop.Scope)

	c.builder.CreateCondBr(expr, forBodyBlock, endBlock)

	c.builder.SetInsertPointAtEnd(forBodyBlock)
	c.generateBlock(forLoop.Block, function, forLoop.Scope)

	c.builder.CreateBr(forUpdateBlock)
	c.builder.SetInsertPointAtEnd(forUpdateBlock)
	c.generateStmt(forLoop.Update, function, forLoop.Scope)

	c.builder.CreateBr(forInitBlock)

	c.builder.SetInsertPointAtEnd(endBlock)
}

func (c *llvmCodegen) generateWhileLoop(
	whileLoop *ast.WhileLoop,
	function *Function,
) {
	whileInitBlock := llvm.AddBasicBlock(function.Fn, ".whileinit")
	whileBodyBlock := llvm.AddBasicBlock(function.Fn, ".whilebody")
	endBlock := llvm.AddBasicBlock(function.Fn, ".whileend")

	c.builder.CreateBr(whileInitBlock)
	c.builder.SetInsertPointAtEnd(whileInitBlock)
	expr := c.getExpr(whileLoop.Cond, whileLoop.Scope)
	c.builder.CreateCondBr(expr, whileBodyBlock, endBlock)

	c.builder.SetInsertPointAtEnd(whileBodyBlock)
	c.generateBlock(whileLoop.Block, function, whileLoop.Scope)
	c.builder.CreateBr(whileInitBlock)

	c.builder.SetInsertPointAtEnd(endBlock)
}
