package llvm

import (
	"fmt"
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
	context llvm.Context
	module  llvm.Module
	builder llvm.Builder

	loc *ast.Loc

	program *ast.Program
	pkg     *ast.Package
}

func NewCG(loc *ast.Loc, program *ast.Program) *llvmCodegen {
	context := llvm.NewContext()
	module := context.NewModule(loc.Dir)
	builder := context.NewBuilder()

	defaultTargetTriple := llvm.DefaultTargetTriple()
	module.SetTarget(defaultTargetTriple)

	return &llvmCodegen{
		loc:     loc,
		program: program,
		context: context,
		module:  module,
		builder: builder,
	}
}

func (c *llvmCodegen) Generate(buildType config.BuildType) error {
	c.generatePackage(c.program.Root)
	err := c.generateExe(buildType)
	return err
}

func (c *llvmCodegen) generatePackage(pkg *ast.Package) {
	if pkg.Processed {
		return
	}

	currentPkg := c.pkg
	defer func() { c.pkg = currentPkg }()
	c.pkg = pkg

	for _, file := range pkg.Files {
		for _, imp := range file.Imports {
			c.generatePackage(imp.Package)
		}
		c.generateDeclarations(file)
	}

	for _, file := range pkg.Files {
		c.generateBodies(file)
	}

	pkg.Processed = true
}

func (c *llvmCodegen) generateDeclarations(file *ast.File) {
	for _, node := range file.Body {
		switch node.Kind {
		case ast.KIND_FN_DECL:
			c.generateFnSignature(node.Node.(*ast.FnDecl))
		case ast.KIND_EXTERN_DECL:
			c.generateExternDecl(node.Node.(*ast.ExternDecl))
		default:
			continue
		}
	}
}

func (c *llvmCodegen) generateBodies(file *ast.File) {
	for _, node := range file.Body {
		switch node.Kind {
		case ast.KIND_FN_DECL:
			c.generateFnBody(node.Node.(*ast.FnDecl))
		default:
			continue
		}
	}
}

func (c *llvmCodegen) generateFnSignature(fnDecl *ast.FnDecl) {
	returnType := c.getType(fnDecl.RetType)
	paramsTypes := c.getFieldListTypes(fnDecl.Params)
	functionType := llvm.FunctionType(returnType, paramsTypes, fnDecl.Params.IsVariadic)
	functionValue := llvm.AddFunction(c.module, fnDecl.Name.Name(), functionType)

	if fnDecl.Attributes != nil {
		functionValue.SetFunctionCallConv(c.getCallingConvention(fnDecl.Attributes.DefaultCallingConvention))
	}

	fnValue := NewFunctionValue(functionValue, functionType)
	fnDecl.BackendType = fnValue
}

func (c *llvmCodegen) generateFnBody(fnDecl *ast.FnDecl) {
	fnValue := fnDecl.BackendType.(*Function)
	functionBlock := c.context.AddBasicBlock(fnValue.Fn, "entry")
	c.builder.SetInsertPointAtEnd(functionBlock)

	c.generateFnParams(fnValue, fnDecl)
	stoppedOnReturn := c.generateBlock(fnDecl.Block, fnValue)
	if !stoppedOnReturn && fnDecl.RetType.IsVoid() {
		c.builder.CreateRetVoid()
	}
}

func (c *llvmCodegen) generateFnParams(fnValue *Function, fnDecl *ast.FnDecl) {
	paramsTypes := fnValue.Ty.ParamTypes()
	fmt.Println(fnValue.Fn.Params())
	for i, paramPtrValue := range fnValue.Fn.Params() {
		field := fnDecl.Params.Fields[i]
		c.generateParam(field, paramsTypes[i], paramPtrValue)
	}
}

func (c *llvmCodegen) generateParam(field *ast.Field, ty llvm.Type, val llvm.Value) {
	paramVal := c.builder.CreateAlloca(ty, ".param")
	c.builder.CreateStore(val, paramVal)
	variable := &Variable{
		Ty:  ty,
		Ptr: paramVal,
	}
	field.BackendType = variable
}

func (c *llvmCodegen) generateBlock(
	block *ast.BlockStmt,
	function *Function,
) bool {
	stoppedOnReturn := false

	for _, stmt := range block.Statements {
		c.generateStmt(stmt, function, block)
		if stmt.IsReturn() {
			stoppedOnReturn = true
			break
		}
	}

	return stoppedOnReturn
}

func (c *llvmCodegen) generateStmt(
	stmt *ast.Node,
	function *Function,
	block *ast.BlockStmt,
) {
	switch stmt.Kind {
	case ast.KIND_FN_CALL:
		c.generateFnCall(stmt.Node.(*ast.FnCall))
	case ast.KIND_RETURN_STMT:
		if len(block.DeferStack) > 0 {
			for i := len(block.DeferStack) - 1; i >= 0; i-- {
				if block.DeferStack[i].Skip {
					continue
				}
				c.generateStmt(block.DeferStack[i].Stmt, function, block)
			}
		}
		c.generateReturnStmt(stmt.Node.(*ast.ReturnStmt))
	case ast.KIND_VAR_STMT:
		c.generateVar(stmt.Node.(*ast.VarStmt))
	case ast.KIND_NAMESPACE_ACCESS:
		c.generateNamespaceAccess(stmt.Node.(*ast.NamespaceAccess))
	case ast.KIND_COND_STMT:
		c.generateCondStmt(stmt.Node.(*ast.CondStmt), function)
	case ast.KIND_FOR_LOOP_STMT:
		c.generateForLoop(stmt.Node.(*ast.ForLoop), function)
	case ast.KIND_WHILE_LOOP_STMT:
		c.generateWhileLoop(stmt.Node.(*ast.WhileLoop), function)
	case ast.KIND_DEFER_STMT:
		// NOTE: ignored because defer statements are handled during block
		// generation
		goto Ignore
	default:
		log.Fatalf("unimplemented block statement: %s\n", stmt)
	}
Ignore:
}

func (c *llvmCodegen) generateReturnStmt(
	ret *ast.ReturnStmt,
) {
	if ret.Value.IsVoid() {
		c.builder.CreateRetVoid()
		return
	}
	returnValue := c.getExpr(ret.Value)
	// TODO: learn more about noundef for return values
	c.builder.CreateRet(returnValue)
}

func (c *llvmCodegen) generateVar(variable *ast.VarStmt) {
	switch variable.Expr.Kind {
	case ast.KIND_TUPLE_EXPR:
		tuple := variable.Expr.Node.(*ast.TupleExpr)

		t := 0
		for _, expr := range tuple.Exprs {
			switch expr.Kind {
			case ast.KIND_TUPLE_EXPR:
				innerTupleExpr := expr.Node.(*ast.TupleExpr)
				for _, innerExpr := range innerTupleExpr.Exprs {
					c.generateVariable(variable.Names[t], innerExpr, variable.IsDecl)
					t++
				}
			case ast.KIND_FN_CALL:
				fnCall := expr.Node.(*ast.FnCall)
				if fnCall.Decl.RetType.Kind == ast.EXPR_TYPE_TUPLE {
					tupleType := fnCall.Decl.RetType.T.(*ast.TupleType)
					affectedVariables := variable.Names[t : t+len(tupleType.Types)]
					c.generateFnCallForTuple(affectedVariables, fnCall, variable.IsDecl)
					t += len(affectedVariables)
				} else {
					c.generateVariable(variable.Names[t], expr, variable.IsDecl)
					t++
				}
			default:
				c.generateVariable(variable.Names[t], expr, variable.IsDecl)
				t++
			}
		}
	case ast.KIND_FN_CALL:
		fnCall := variable.Expr.Node.(*ast.FnCall)
		c.generateFnCallForTuple(variable.Names, fnCall, variable.IsDecl)
	default:
		c.generateVariable(variable.Names[0], variable.Expr, variable.IsDecl)
	}
}

func (c *llvmCodegen) generateFnCallForTuple(variables []*ast.VarId, fnCall *ast.FnCall, isDecl bool) {
	genFnCall := c.generateFnCall(fnCall)

	if len(variables) == 1 {
		c.generateVariableWithValue(variables[0], genFnCall, isDecl)
	} else {
		for i, currentVar := range variables {
			value := c.builder.CreateExtractValue(genFnCall, i, ".arg")
			c.generateVariableWithValue(currentVar, value, isDecl)
		}
	}
}

func (c *llvmCodegen) generateVariable(name *ast.VarId, expr *ast.Node, isDecl bool) {
	if isDecl {
		c.generateVarDecl(name, expr)
	} else {
		c.generateVarReassign(name, expr)
	}
}

func (c *llvmCodegen) generateVariableWithValue(name *ast.VarId, value llvm.Value, isDecl bool) {
	if isDecl {
		c.generateVarDeclWithValue(name, value)
	} else {
		c.generateVarReassignWithValue(name, value)
	}
}

func (c *llvmCodegen) generateVarDecl(
	name *ast.VarId,
	expr *ast.Node,
) {
	ty := c.getType(name.Type)
	ptr := c.builder.CreateAlloca(ty, ".ptr")
	generatedExpr := c.getExpr(expr)
	c.builder.CreateStore(generatedExpr, ptr)

	variableLlvm := &Variable{
		Ty:  ty,
		Ptr: ptr,
	}
	name.BackendType = variableLlvm
}

func (c *llvmCodegen) generateVarReassign(
	name *ast.VarId,
	expr *ast.Node,
) {
	generatedExpr := c.getExpr(expr)

	var variable *Variable

	switch name.N.Kind {
	case ast.KIND_VAR_ID_STMT:
		node := name.N.Node.(*ast.VarId)
		variable = node.BackendType.(*Variable)
	case ast.KIND_FIELD:
		node := name.N.Node.(*ast.Field)
		variable = node.BackendType.(*Variable)
	default:
		log.Fatalf("invalid symbol on generateVarReassign: %v\n", reflect.TypeOf(variable))
	}

	c.builder.CreateStore(generatedExpr, variable.Ptr)
}

func (c *llvmCodegen) generateVarDeclWithValue(
	name *ast.VarId,
	value llvm.Value,
) {
	ty := c.getType(name.Type)
	ptr := c.builder.CreateAlloca(ty, ".ptr")
	c.builder.CreateStore(value, ptr)

	variableLlvm := &Variable{
		Ty:  ty,
		Ptr: ptr,
	}
	name.BackendType = variableLlvm
}

func (c *llvmCodegen) generateVarReassignWithValue(
	name *ast.VarId,
	value llvm.Value,
) {
	var variable *Variable

	switch name.N.Kind {
	case ast.KIND_VAR_ID_STMT:
		node := name.N.Node.(*ast.VarId)
		variable = node.BackendType.(*Variable)
	case ast.KIND_FIELD:
		node := name.N.Node.(*ast.Field)
		variable = node.BackendType.(*Variable)
	default:
		log.Fatalf("invalid symbol on generateVarReassign: %v\n", reflect.TypeOf(variable))
	}

	c.builder.CreateStore(value, variable.Ptr)
}

func (c *llvmCodegen) generateFnCall(
	fnCall *ast.FnCall,
) llvm.Value {
	callLlvm := fnCall.Decl.BackendType.(*Function)
	args := c.getExprList(fnCall.Args)
	return c.builder.CreateCall(callLlvm.Ty, callLlvm.Fn, args, "")
}

func (c *llvmCodegen) generateTupleExpr(tuple *ast.TupleExpr) llvm.Value {
	types := c.getTypes(tuple.Type.Types)
	tupleTy := c.context.StructType(types, false)
	tupleVal := llvm.ConstNull(tupleTy)

	for i, expr := range tuple.Exprs {
		generatedExpr := c.getExpr(expr)
		tupleVal = c.builder.CreateInsertValue(tupleVal, generatedExpr, i, "")
	}

	return tupleVal
}

func (c *llvmCodegen) generateExternDecl(external *ast.ExternDecl) {
	for i := range external.Prototypes {
		c.generatePrototype(external.Attributes, external.Prototypes[i])
	}
}

func (c *llvmCodegen) generatePrototype(attributes *ast.Attributes, prototype *ast.Proto) {
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

	proto := NewFunctionValue(protoValue, ty)
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
		case token.BOOL_TYPE, token.UNTYPED_END:
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
	case ast.EXPR_TYPE_TUPLE:
		tuple := ty.T.(*ast.TupleType)
		innerTypes := c.getTypes(tuple.Types)
		tupleTy := c.context.StructType(innerTypes, false)
		return tupleTy
	default:
		log.Fatalf("unimplemented type: %v\n", ty.Kind)
	}
	return c.context.VoidType()
}

func (c *llvmCodegen) getPtrType(ty llvm.Type) llvm.Type {
	return llvm.PointerType(ty, 0)
}

func (c *llvmCodegen) getTypes(types []*ast.ExprType) []llvm.Type {
	tys := make([]llvm.Type, len(types))
	for i, ty := range types {
		tys[i] = c.getType(ty)
	}
	return tys
}

func (c *llvmCodegen) getFieldListTypes(params *ast.FieldList) []llvm.Type {
	types := make([]llvm.Type, params.Len)
	for i := range params.Len {
		types[i] = c.getType(params.Fields[i].Type)
	}
	return types
}

func (c *llvmCodegen) getExprList(expressions []*ast.Node) []llvm.Value {
	values := make([]llvm.Value, len(expressions))
	for i, expr := range expressions {
		values[i] = c.getExpr(expr)
	}
	return values
}

func (c *llvmCodegen) getExpr(expr *ast.Node) llvm.Value {
	switch expr.Kind {
	case ast.KIND_LITERAl_EXPR:
		lit := expr.Node.(*ast.LiteralExpr)
		switch lit.Type.Kind {
		case ast.EXPR_TYPE_BASIC:
			basic := lit.Type.T.(*ast.BasicType)

			if basic.Kind.IsNumeric() {
				// TODO: generate numeric type properly
				// Right now it only accepts integers
				integerValue, bitSize := c.getIntegerValue(lit, basic)
				return llvm.ConstInt(c.context.IntType(bitSize), integerValue, false)
			}

			switch basic.Kind {
			// case basic.IsIntegerType():
			// 	integerValue, bitSize := c.getIntegerValue(lit, basic)
			// 	return llvm.ConstInt(c.context.IntType(bitSize), integerValue, false)
			case token.STRING_TYPE, token.CSTRING_TYPE:
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
			case token.BOOL_TYPE, token.UNTYPED_BOOL:
				ty := c.getType(lit.Type)
				boolVal := string(lit.Value)

				var val uint64
				if boolVal == "true" {
					val = 1
				} else {
					val = 0
				}

				return llvm.ConstInt(ty, val, false)
			default:
				log.Fatalf("unimplemented basic type: %d %s\n", basic.Kind, string(lit.Value))
				return llvm.Value{}
			}
		case ast.EXPR_TYPE_POINTER:
			panic("unimplemented pointer type")
		default:
			panic("unimplemented literal type")
		}
	case ast.KIND_ID_EXPR:
		id := expr.Node.(*ast.IdExpr)

		var localVar *Variable
		switch id.N.Kind {
		case ast.KIND_VAR_ID_STMT:
			variable := id.N.Node.(*ast.VarId)
			localVar = variable.BackendType.(*Variable)
		case ast.KIND_FIELD:
			variable := id.N.Node.(*ast.Field)
			if variable.Variadic {
				ty := c.getType(variable.Type)
				val := c.builder.CreateAlloca(ty, ".ptr")
				c.generateParam(variable, ty, val)
			}
			localVar = variable.BackendType.(*Variable)
		}

		loadedVariable := c.builder.CreateLoad(localVar.Ty, localVar.Ptr, ".load")
		return loadedVariable
	case ast.KIND_BINARY_EXPR:
		binary := expr.Node.(*ast.BinaryExpr)

		lhs := c.getExpr(binary.Left)
		rhs := c.getExpr(binary.Right)

		switch binary.Op {
		// TODO: deal with signed or unsigned operations
		// I'm assuming all unsigned for now
		case token.EQUAL_EQUAL:
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
		call := c.generateFnCall(fnCall)
		return call
	case ast.KIND_UNARY_EXPR:
		unary := expr.Node.(*ast.UnaryExpr)
		expr := c.getExpr(unary.Value)

		switch unary.Op {
		case token.MINUS:
			return c.builder.CreateNeg(expr, ".neg")
		case token.NOT:
			falseVal := llvm.ConstInt(c.context.Int1Type(), 0, false)
			return c.builder.CreateICmp(llvm.IntEQ, falseVal, expr, ".not")
		default:
			log.Fatalf("unimplemented unary operator: %s", unary.Op)
			return llvm.Value{}
		}
	case ast.KIND_NAMESPACE_ACCESS:
		namespaceAccess := expr.Node.(*ast.NamespaceAccess)
		value := c.generateNamespaceAccess(namespaceAccess)
		return value
	case ast.KIND_TUPLE_EXPR:
		return c.generateTupleExpr(expr.Node.(*ast.TupleExpr))
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
	ifExpr := c.getExpr(condStmt.IfStmt.Expr)
	c.builder.CreateCondBr(ifExpr, ifBlock, elseBlock)
	c.builder.SetInsertPointAtEnd(ifBlock)
	stoppedOnReturn := c.generateBlock(condStmt.IfStmt.Block, function)
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
		)
		if !elseStoppedOnReturn {
			c.builder.CreateBr(endBlock)
		}
	} else {
		c.builder.CreateBr(endBlock)
	}
	c.builder.SetInsertPointAtEnd(endBlock)
}

func (c *llvmCodegen) generateNamespaceAccess(
	namespaceAccess *ast.NamespaceAccess,
) llvm.Value {
	if namespaceAccess.IsImport {
		return c.generateImportAccess(namespaceAccess.Right)
	}

	switch namespaceAccess.Left.N.Kind {
	case ast.KIND_EXTERN_DECL:
		fnCall := namespaceAccess.Right.Node.(*ast.FnCall)
		return c.generatePrototypeCall(fnCall)
	default:
		panic("unreachable")
	}
}

func (c *llvmCodegen) generateImportAccess(right *ast.Node) llvm.Value {
	switch right.Kind {
	case ast.KIND_FN_CALL:
		fnCall := right.Node.(*ast.FnCall)
		return c.generateFnCall(fnCall)
	case ast.KIND_NAMESPACE_ACCESS:
		namespaceAccess := right.Node.(*ast.NamespaceAccess)
		return c.generateNamespaceAccess(namespaceAccess)
	default:
		panic("unimplemented import access")
	}
}

func (c *llvmCodegen) generatePrototypeCall(call *ast.FnCall) llvm.Value {
	protoLlvm := call.Proto.BackendType.(*Function)
	args := c.getExprList(call.Args)
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
	c.generateStmt(forLoop.Init, function, forLoop.Block)

	c.builder.CreateBr(forInitBlock)
	c.builder.SetInsertPointAtEnd(forInitBlock)

	expr := c.getExpr(forLoop.Cond)

	c.builder.CreateCondBr(expr, forBodyBlock, endBlock)

	c.builder.SetInsertPointAtEnd(forBodyBlock)
	c.generateBlock(forLoop.Block, function)

	c.builder.CreateBr(forUpdateBlock)
	c.builder.SetInsertPointAtEnd(forUpdateBlock)
	c.generateStmt(forLoop.Update, function, forLoop.Block)

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
	expr := c.getExpr(whileLoop.Cond)
	c.builder.CreateCondBr(expr, whileBodyBlock, endBlock)

	c.builder.SetInsertPointAtEnd(whileBodyBlock)
	c.generateBlock(whileLoop.Block, function)
	c.builder.CreateBr(whileInitBlock)

	c.builder.SetInsertPointAtEnd(endBlock)
}

func (c *llvmCodegen) generateExe(buildType config.BuildType) error {
	dir, err := os.MkdirTemp("", "build")
	if err != nil {
		return err
	}

	filenameNoExt := strings.TrimSuffix(filepath.Base(c.loc.Name), filepath.Ext(c.loc.Name))
	irFileName := filepath.Join(dir, filenameNoExt)
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
		if config.DEBUG_MODE {
			fmt.Printf("[DEBUG MODE] OPT COMMAND: %s\n", cmd)
		}
		return err
	}

	cmd = exec.Command("clang-18", compilerFlags, optLevel, "-o", filenameNoExt, optimizedIrFilepath)
	err = cmd.Run()
	// TODO(errors)
	if err != nil {
		if config.DEBUG_MODE {
			fmt.Printf("[DEBUG MODE] CLANG COMMAND: %s\n", cmd)
		}
		return err
	}

	if config.DEBUG_MODE {
		fmt.Printf("[DEBUG MODE] keeping '%s' build directory\n", dir)
	} else {
		err = os.RemoveAll(dir)
		if err != nil {
			return err
		}
	}

	return nil
}
