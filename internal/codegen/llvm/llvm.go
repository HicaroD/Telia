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

var (
	context llvm.Context = llvm.NewContext()
	builder llvm.Builder = context.NewBuilder()
)

type codegen struct {
	module llvm.Module

	loc     *ast.Loc
	program *ast.Program
	pkg     *ast.Package
}

func NewCG(loc *ast.Loc, program *ast.Program) *codegen {
	module := context.NewModule(loc.Dir)

	defaultTargetTriple := llvm.DefaultTargetTriple()
	module.SetTarget(defaultTargetTriple)

	return &codegen{
		loc:     loc,
		program: program,
		module:  module,
	}
}

func (c *codegen) Generate(buildType config.BuildType) error {
	c.generatePackage(c.program.Root)
	err := c.generateExe(buildType)
	return err
}

func (c *codegen) generatePackage(pkg *ast.Package) {
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
		c.emitDeclarations(file)
	}

	for _, file := range pkg.Files {
		c.emitBodies(file)
	}

	pkg.Processed = true
}

func (c *codegen) emitDeclarations(file *ast.File) {
	for _, node := range file.Body {
		switch node.Kind {
		case ast.KIND_FN_DECL:
			c.emitFnSignature(node.Node.(*ast.FnDecl))
		case ast.KIND_EXTERN_DECL:
			c.emitExternDecl(node.Node.(*ast.ExternDecl))
		case ast.KIND_STRUCT_DECL:
			c.emitStructDecl(node.Node.(*ast.StructDecl))
		default:
			continue
		}
	}
}

func (c *codegen) emitBodies(file *ast.File) {
	for _, node := range file.Body {
		switch node.Kind {
		case ast.KIND_FN_DECL:
			c.generateFnBody(node.Node.(*ast.FnDecl))
		default:
			continue
		}
	}
}

func (c *codegen) emitFnSignature(fnDecl *ast.FnDecl) {
	returnType := c.emitType(fnDecl.RetType)
	paramsTypes := c.getFieldListTypes(fnDecl.Params)
	functionType := llvm.FunctionType(returnType, paramsTypes, fnDecl.Params.IsVariadic)
	functionValue := llvm.AddFunction(c.module, fnDecl.Name.Name(), functionType)

	if fnDecl.Attributes != nil {
		functionValue.SetFunctionCallConv(c.getCallingConvention(fnDecl.Attributes.DefaultCallingConvention))
	}
}

func (c *codegen) generateFnBody(fnDecl *ast.FnDecl) {
	fn := c.module.NamedFunction(fnDecl.Name.Name())
	if fn.IsNil() {
		panic("function is nil when generating fn body")
	}
	body := context.AddBasicBlock(fn, "entry")
	builder.SetInsertPointAtEnd(body)

	paramsTypes := fn.GlobalValueType().ParamTypes()
	for i, paramPtrValue := range fn.Params() {
		field := fnDecl.Params.Fields[i]
		c.generateParam(field, paramsTypes[i], paramPtrValue)
	}

	stoppedOnReturn := c.generateBlock(fn, fnDecl.Block)
	if !stoppedOnReturn && fnDecl.RetType.IsVoid() {
		builder.CreateRetVoid()
	}
}

func (c *codegen) generateParam(field *ast.Param, ty llvm.Type, val llvm.Value) {
	paramVal := builder.CreateAlloca(ty, ".param")
	builder.CreateStore(val, paramVal)
	variable := &Variable{
		Ty:  ty,
		Ptr: paramVal,
	}
	field.BackendType = variable
}

func (c *codegen) generateBlock(
	fn llvm.Value,
	block *ast.BlockStmt,
) bool {
	stoppedOnReturn := false

	for _, stmt := range block.Statements {
		c.generateStmt(fn, block, stmt)
		if stmt.IsReturn() {
			stoppedOnReturn = true
			break
		}
	}

	return stoppedOnReturn
}

func (c *codegen) generateStmt(
	function llvm.Value,
	block *ast.BlockStmt,
	stmt *ast.Node,
) {
	switch stmt.Kind {
	case ast.KIND_FN_CALL:
		c.emitFnCall(stmt.Node.(*ast.FnCall))
	case ast.KIND_RETURN_STMT:
		if len(block.DeferStack) > 0 {
			for i := len(block.DeferStack) - 1; i >= 0; i-- {
				if block.DeferStack[i].Skip {
					continue
				}
				c.generateStmt(function, block, block.DeferStack[i].Stmt)
			}
		}
		c.generateReturnStmt(stmt.Node.(*ast.ReturnStmt))
	case ast.KIND_VAR_STMT:
		c.generateVar(stmt.Node.(*ast.VarStmt))
	case ast.KIND_NAMESPACE_ACCESS:
		c.emitNamespaceAccess(stmt.Node.(*ast.NamespaceAccess))
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

func (c *codegen) generateReturnStmt(
	ret *ast.ReturnStmt,
) {
	if ret.Value.IsVoid() {
		builder.CreateRetVoid()
		return
	}
	returnValue, _, _ := c.emitExpr(ret.Value)
	// TODO: learn more about noundef for return values
	builder.CreateRet(returnValue)
}

func (c *codegen) generateVar(variable *ast.VarStmt) {
	switch variable.Expr.Kind {
	case ast.KIND_TUPLE_LITERAL_EXPR:
		tuple := variable.Expr.Node.(*ast.TupleExpr)

		t := 0
		for _, expr := range tuple.Exprs {
			switch expr.Kind {
			case ast.KIND_TUPLE_LITERAL_EXPR:
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

func (c *codegen) generateFnCallForTuple(variables []*ast.Node, fnCall *ast.FnCall, isDecl bool) {
	genFnCall := c.emitFnCall(fnCall)

	if len(variables) == 1 {
		c.generateVariableWithValue(variables[0], genFnCall, isDecl)
	} else {
		for i, currentVar := range variables {
			value := builder.CreateExtractValue(genFnCall, i, ".arg")
			c.generateVariableWithValue(currentVar, value, isDecl)
		}
	}
}

func (c *codegen) generateVariable(name *ast.Node, expr *ast.Node, isDecl bool) {
	if isDecl {
		c.generateVarDecl(name, expr)
	} else {
		c.generateVarReassign(name, expr)
	}
}

func (c *codegen) generateVariableWithValue(name *ast.Node, value llvm.Value, isDecl bool) {
	if isDecl {
		c.generateVarDeclWithValue(name, value)
	} else {
		c.generateVarReassignWithValue(name, value)
	}
}

func (c *codegen) generateVarDecl(
	name *ast.Node,
	expr *ast.Node,
) {
	var ty llvm.Type
	var ptr llvm.Value

	variable := name.Node.(*ast.VarIdStmt)
	ty = c.emitType(variable.Type)

	switch expr.Kind {
	case ast.KIND_STRUCT_LITERAL_EXPR:
		ptr = c.emitStructLiteral(expr.Node.(*ast.StructLiteralExpr))
	case ast.KIND_NAMESPACE_ACCESS:
		ptr = c.emitNamespaceAccess(expr.Node.(*ast.NamespaceAccess))
	default:
		ptr = builder.CreateAlloca(ty, ".ptr")
		generatedExpr, _, _ := c.emitExpr(expr)
		builder.CreateStore(generatedExpr, ptr)
	}

	variableLlvm := &Variable{
		Ty:  ty,
		Ptr: ptr,
	}
	variable.BackendType = variableLlvm
}

func (c *codegen) generateVarReassign(
	name *ast.Node,
	expr *ast.Node,
) {
	generatedExpr, _, _ := c.emitExpr(expr)
	var varPtr llvm.Value

	switch name.Kind {
	case ast.KIND_VAR_ID_STMT:
		varId := name.Node.(*ast.VarIdStmt)
		switch varId.N.Kind {
		case ast.KIND_VAR_ID_STMT:
			variable := varId.N.Node.(*ast.VarIdStmt)
			varPtr = variable.BackendType.(*Variable).Ptr
		case ast.KIND_FIELD:
			param := varId.N.Node.(*ast.Param)
			varPtr = param.BackendType.(*Variable).Ptr
		default:
			panic(fmt.Sprintf("unimplemented kind of name expression: %v\n", varId.N))
		}
	case ast.KIND_FIELD_ACCESS:
		node := name.Node.(*ast.FieldAccess)
		varPtr = c.getStructFieldPtr(node)
	default:
		log.Fatalf("invalid symbol on generateVarReassign: %v\n", name.Kind)
	}

	builder.CreateStore(generatedExpr, varPtr)
}

func (c *codegen) getStructFieldPtr(fieldAccess *ast.FieldAccess) llvm.Value {
	st := c.module.GetTypeByName(fieldAccess.Decl.Name.Name())
	if st.IsNil() {
		panic("struct is nil")
	}

	switch fieldAccess.Right.Kind {
	case ast.KIND_ID_EXPR:
		obj := fieldAccess.StructVar.BackendType.(*Variable)
		return builder.CreateStructGEP(st, obj.Ptr, fieldAccess.AccessedField.Index, ".field")
	default:
		panic("unimplemented other type of field access")
	}
}

func (c *codegen) generateVarDeclWithValue(
	name *ast.Node,
	value llvm.Value,
) {
	variable := name.Node.(*ast.VarIdStmt)

	ty := c.emitType(variable.Type)
	ptr := builder.CreateAlloca(ty, ".ptr")
	builder.CreateStore(value, ptr)

	variableLlvm := &Variable{
		Ty:  ty,
		Ptr: ptr,
	}
	variable.BackendType = variableLlvm
}

func (c *codegen) generateVarReassignWithValue(
	name *ast.Node,
	value llvm.Value,
) {
	var variable *Variable

	switch name.Kind {
	case ast.KIND_VAR_ID_STMT:
		node := name.Node.(*ast.VarIdStmt)
		switch node.N.Kind {
		case ast.KIND_VAR_ID_STMT:
			varId := node.N.Node.(*ast.VarIdStmt)
			variable = varId.BackendType.(*Variable)
		case ast.KIND_FIELD:
			param := node.N.Node.(*ast.Param)
			variable = param.BackendType.(*Variable)
		}
	// TODO: check if the case below is reachable
	case ast.KIND_FIELD:
		node := name.Node.(*ast.Param)
		variable = node.BackendType.(*Variable)
	default:
		log.Fatalf("invalid symbol on generateVarReassign: %v\n", reflect.TypeOf(variable))
	}

	builder.CreateStore(value, variable.Ptr)
}

func (c *codegen) emitFnCall(call *ast.FnCall) llvm.Value {
	args := c.getCallArgs(call)
	fn := c.module.NamedFunction(call.Name.Name())
	if fn.IsNil() {
		panic("function is nil when generating function call")
	}
	return builder.CreateCall(fn.GlobalValueType(), fn, args, "")
}

func (c *codegen) emitTupleExpr(tuple *ast.TupleExpr) llvm.Value {
	types := c.emitTypes(tuple.Type.Types)
	tupleTy := context.StructType(types, false)
	tupleVal := llvm.ConstNull(tupleTy)

	for i, expr := range tuple.Exprs {
		generatedExpr, _, _ := c.emitExpr(expr)
		tupleVal = builder.CreateInsertValue(tupleVal, generatedExpr, i, "")
	}

	return tupleVal
}

func (c *codegen) emitStructLiteral(lit *ast.StructLiteralExpr) llvm.Value {
	st := c.module.GetTypeByName(lit.Name.Name())
	if st.IsNil() {
		panic("struct is nil")
	}

	obj := builder.CreateAlloca(st, ".obj")
	for _, field := range lit.Values {
		expr, _, _ := c.emitExpr(field.Value)
		allocatedField := builder.CreateStructGEP(st, obj, field.Index, ".field")
		builder.CreateStore(expr, allocatedField)
	}

	return obj
}

func (c *codegen) emitFieldAccessExpr(fieldAccess *ast.FieldAccess) llvm.Value {
	ptr := c.getStructFieldPtr(fieldAccess)
	ty := c.emitType(fieldAccess.AccessedField.Type)
	loadedPtr := builder.CreateLoad(ty, ptr, ".access")
	return loadedPtr
}

func (c *codegen) emitExternDecl(external *ast.ExternDecl) {
	for _, proto := range external.Prototypes {
		c.generatePrototype(proto.Attributes, proto)
	}
}

func (c *codegen) emitStructDecl(st *ast.StructDecl) {
	fields := c.getStructFieldList(st.Fields)
	structTy := context.StructCreateNamed(st.Name.Name())
	// TODO: set packed properly
	packed := false
	structTy.StructSetBody(fields, packed)
}

func (c *codegen) generatePrototype(externAttributes *ast.Attributes, prototype *ast.Proto) {
	returnTy := c.emitType(prototype.RetType)
	paramsTypes := c.getFieldListTypes(prototype.Params)
	ty := llvm.FunctionType(returnTy, paramsTypes, prototype.Params.IsVariadic)
	protoValue := llvm.AddFunction(c.module, prototype.Name.Name(), ty)

	if externAttributes != nil {
		protoValue.SetFunctionCallConv(c.getCallingConvention(externAttributes.DefaultCallingConvention))
	}
	if prototype.Attributes == nil {
		return
	}

	if prototype.Attributes.DefaultCallingConvention != "" {
		protoValue.SetLinkage(c.getFunctionLinkage(prototype.Attributes.Linkage))
	}
	if prototype.Attributes.LinkName != "" {
		protoValue.SetName(externAttributes.LinkName)
	}
}

func (c *codegen) getFunctionLinkage(linkage string) llvm.Linkage {
	switch linkage {
	case "internal":
		return llvm.InternalLinkage
	// NOTE: there are several types of weak linkages
	// Make sure to choose the perfect match
	case "weak":
		return llvm.WeakAnyLinkage
	// NOTE: there are several types of link_once linkage
	// Makesure to choose the perfect match
	case "link_once":
		return llvm.LinkOnceAnyLinkage
	case "external":
		fallthrough
	default:
		return llvm.ExternalLinkage
	}
}

func (c *codegen) getCallingConvention(callingConvention string) llvm.CallConv {
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

func (c *codegen) emitType(ty *ast.ExprType) llvm.Type {
	switch ty.Kind {
	case ast.EXPR_TYPE_BASIC:
		b := ty.T.(*ast.BasicType)
		switch b.Kind {
		case token.BOOL_TYPE, token.UNTYPED_BOOL, token.INT_TYPE,
			token.UINT_TYPE, token.UNTYPED_INT, token.I8_TYPE, token.U8_TYPE,
			token.I16_TYPE, token.U16_TYPE, token.I32_TYPE, token.U32_TYPE,
			token.I64_TYPE, token.U64_TYPE, token.I128_TYPE, token.U128_TYPE:

			bitSize := b.Kind.BitSize()
			return context.IntType(bitSize)
		case token.F32_TYPE, token.FLOAT_TYPE, token.UNTYPED_FLOAT:
			return context.FloatType()
		case token.F64_TYPE:
			return context.DoubleType()
		case token.VOID_TYPE:
			return context.VoidType()
		case token.STRING_TYPE, token.CSTRING_TYPE:
			u8 := ast.NewBasicType(token.U8_TYPE)
			u8Type := c.emitType(u8)
			return c.emitPtrType(u8Type)
		case token.RAWPTR_TYPE:
			u8 := ast.NewBasicType(token.U8_TYPE)
			u8Type := c.emitType(u8)
			return c.emitPtrType(u8Type)
		}
	case ast.EXPR_TYPE_TUPLE:
		tuple := ty.T.(*ast.TupleType)
		innerTypes := c.emitTypes(tuple.Types)
		tupleTy := context.StructType(innerTypes, false)
		return tupleTy
	case ast.EXPR_TYPE_STRUCT:
		ty := ty.T.(*ast.StructType)
		return c.module.GetTypeByName(ty.Decl.Name.Name())
	default:
		log.Fatalf("unimplemented type: %v\n", ty.Kind)
	}
	return context.VoidType()
}

func (c *codegen) emitPtrType(ty llvm.Type) llvm.Type {
	return llvm.PointerType(ty, 0)
}

func (c *codegen) emitTypes(types []*ast.ExprType) []llvm.Type {
	tys := make([]llvm.Type, len(types))
	for i, ty := range types {
		tys[i] = c.emitType(ty)
	}
	return tys
}

func (c *codegen) getStructFieldList(fields []*ast.StructField) []llvm.Type {
	types := make([]llvm.Type, len(fields))
	for i, field := range fields {
		types[i] = c.emitType(field.Type)
	}
	return types
}

func (c *codegen) getFieldListTypes(params *ast.Params) []llvm.Type {
	types := make([]llvm.Type, params.Len)
	for i := range params.Len {
		types[i] = c.emitType(params.Fields[i].Type)
	}
	return types
}

func (c *codegen) getCallArgs(call *ast.FnCall) []llvm.Value {
	nExprs := len(call.Args)
	if call.Variadic {
		nExprs-- // desconsider the variadic argument place for now
	}

	values := make([]llvm.Value, nExprs)
	for i := range nExprs {
		value, _, _ := c.emitExpr(call.Args[i])
		values[i] = value
	}

	if call.Variadic {
		variadic := call.Args[nExprs].Node.(*ast.VarArgsExpr)
		for _, arg := range variadic.Args {
			varArg, _, _ := c.emitExpr(arg)
			values = append(values, varArg)
		}
	}
	return values
}

func (c *codegen) emitExpr(expr *ast.Node) (llvm.Value, int, bool) {
	hasFloat := false
	bitSize := 32

	switch expr.Kind {
	case ast.KIND_LITERAL_EXPR:
		lit, bitSize, hasFloat := c.emitLiteralExpr(expr.Node.(*ast.LiteralExpr))
		return lit, bitSize, hasFloat
	case ast.KIND_ID_EXPR:
		id, bitSize, hasFloat := c.emitIdExpr(expr.Node.(*ast.IdExpr))
		return id, bitSize, hasFloat
	case ast.KIND_BINARY_EXPR:
		binExpr, bitSize, hasFloat := c.emitBinExpr(expr.Node.(*ast.BinExpr))
		return binExpr, bitSize, hasFloat
	case ast.KIND_FN_CALL:
		fnCall := expr.Node.(*ast.FnCall)
		call := c.emitFnCall(fnCall)
		return call, bitSize, hasFloat
	case ast.KIND_UNARY_EXPR:
		unary, bitSize, hasFloat := c.emitUnaryExpr(expr.Node.(*ast.UnaryExpr))
		return unary, bitSize, hasFloat
	case ast.KIND_NAMESPACE_ACCESS:
		value := c.emitNamespaceAccess(expr.Node.(*ast.NamespaceAccess))
		return value, bitSize, hasFloat
	case ast.KIND_TUPLE_LITERAL_EXPR:
		return c.emitTupleExpr(expr.Node.(*ast.TupleExpr)), bitSize, hasFloat
	case ast.KIND_STRUCT_LITERAL_EXPR:
		return c.emitStructLiteral(expr.Node.(*ast.StructLiteralExpr)), bitSize, hasFloat
	case ast.KIND_FIELD_ACCESS:
		return c.emitFieldAccessExpr(expr.Node.(*ast.FieldAccess)), bitSize, hasFloat
	case ast.KIND_VARG_EXPR:
		panic("unimplemented var args")
	default:
		panic(fmt.Sprintf("unimplemented expr: %v", expr))
	}
}

func (c *codegen) emitLiteralExpr(lit *ast.LiteralExpr) (llvm.Value, int, bool) {
	hasFloat := false
	bitSize := 32

	switch lit.Type.Kind {
	case ast.EXPR_TYPE_BASIC:
		basic := lit.Type.T.(*ast.BasicType)

		if basic.Kind.IsFloat() {
			hasFloat = true
			floatValue, bitSize := c.getFloatValue(lit, basic)
			if bitSize == 32 {
				return llvm.ConstFloat(context.FloatType(), floatValue), bitSize, hasFloat
			} else {
				return llvm.ConstFloat(context.DoubleType(), floatValue), bitSize, hasFloat
			}
		}

		if basic.Kind.IsInteger() {
			integerValue, bitSize := c.getIntegerValue(lit, basic)
			return llvm.ConstInt(context.IntType(bitSize), integerValue, false), bitSize, hasFloat
		}

		switch basic.Kind {
		case token.STRING_TYPE, token.CSTRING_TYPE:
			strlen := len(lit.Value) + 1
			arrTy := llvm.ArrayType(context.Int8Type(), strlen)
			arr := llvm.ConstArray(arrTy, c.llvmConstInt8s(lit.Value, strlen))

			globalVal := llvm.AddGlobal(c.module, arrTy, ".str")
			globalVal.SetInitializer(arr)
			globalVal.SetLinkage(llvm.PrivateLinkage)
			globalVal.SetGlobalConstant(true)
			globalVal.SetAlignment(1)
			globalVal.SetUnnamedAddr(true)

			zero := llvm.ConstInt(context.Int32Type(), 0, false)
			indices := []llvm.Value{zero, zero}
			ptr := llvm.ConstInBoundsGEP(arrTy, globalVal, indices)

			return ptr, bitSize, hasFloat
		case token.BOOL_TYPE, token.UNTYPED_BOOL:
			ty := c.emitType(lit.Type)
			boolVal := string(lit.Value)
			var val uint64
			if boolVal == "true" {
				val = 1
			} else {
				val = 0
			}
			return llvm.ConstInt(ty, val, false), bitSize, hasFloat
		default:
			log.Fatalf("unimplemented basic type: %d %s\n", basic.Kind, string(lit.Value))
			return llvm.Value{}, bitSize, hasFloat
		}
	case ast.EXPR_TYPE_POINTER:
		panic("unimplemented pointer type")
	default:
		panic("unimplemented literal type")
	}
}

func (c *codegen) emitIdExpr(id *ast.IdExpr) (llvm.Value, int, bool) {
	var localVar *Variable
	var varTy *ast.ExprType

	hasFloat := false
	bitSize := 32

	switch id.N.Kind {
	case ast.KIND_VAR_ID_STMT:
		variable := id.N.Node.(*ast.VarIdStmt)
		varTy = variable.Type
		localVar = variable.BackendType.(*Variable)
	case ast.KIND_FIELD:
		variable := id.N.Node.(*ast.Param)
		varTy = variable.Type
		localVar = variable.BackendType.(*Variable)
	}

	if varTy.IsFloat() {
		bt := varTy.T.(*ast.BasicType)
		bitSize = bt.Kind.BitSize()
		hasFloat = true
	}

	loadedVariable := builder.CreateLoad(localVar.Ty, localVar.Ptr, ".load")
	return loadedVariable, bitSize, hasFloat
}

func (c *codegen) emitBinExpr(bin *ast.BinExpr) (llvm.Value, int, bool) {
	lhs, lhsBitSize, lhsHasFloat := c.emitExpr(bin.Left)
	rhs, rhsBitSize, rhsHasFloat := c.emitExpr(bin.Right)

	if lhsBitSize != rhsBitSize {
		panic("expected bit size to be the same")
	}

	hasFloat := lhsHasFloat || rhsHasFloat
	bitSize := lhsBitSize // since they are the same

	var binExpr llvm.Value
	if hasFloat {
		binExpr = c.emitFloatBinExpr(lhs, bin.Op, rhs)
	} else {
		binExpr = c.emitIntBinExpr(lhs, bin.Op, rhs)
	}
	return binExpr, bitSize, hasFloat
}

func (c *codegen) emitUnaryExpr(unary *ast.UnaryExpr) (llvm.Value, int, bool) {
	expr, bitSize, hasFloat := c.emitExpr(unary.Value)

	if hasFloat {
		switch unary.Op {
		case token.MINUS:
			return builder.CreateFNeg(expr, ".fneg"), bitSize, hasFloat
		case token.NOT:
			var ty llvm.Type
			if bitSize == 32 {
				ty = context.FloatType()
			} else {
				ty = context.DoubleType()
			}
			falseVal := llvm.ConstFloat(ty, 0)
			return builder.CreateICmp(llvm.IntEQ, falseVal, expr, ".not"), bitSize, hasFloat
		default:
			panic(fmt.Sprintf("unimplemented unary operator: %s", unary.Op))
		}
	} else {
		switch unary.Op {
		case token.MINUS:
			return builder.CreateNeg(expr, ".neg"), bitSize, hasFloat
		case token.NOT:
			falseVal := llvm.ConstInt(context.Int1Type(), 0, false)
			return builder.CreateICmp(llvm.IntEQ, falseVal, expr, ".not"), bitSize, hasFloat
		default:
			panic(fmt.Sprintf("unimplemented unary operator: %s", unary.Op))
		}
	}
}

func (c *codegen) emitIntBinExpr(lhs llvm.Value, binOp token.Kind, rhs llvm.Value) llvm.Value {
	if binOp.IsCmpOp() {
		var predicate llvm.IntPredicate

		switch binOp {
		case token.EQUAL_EQUAL:
			predicate = llvm.IntEQ
		case token.BANG_EQUAL:
			predicate = llvm.IntNE
		case token.LESS:
			predicate = llvm.IntULT
		case token.LESS_EQ:
			predicate = llvm.IntULE
		case token.GREATER:
			predicate = llvm.IntUGT
		case token.GREATER_EQ:
			predicate = llvm.IntUGE
		default:
			panic(fmt.Sprintf("unsupported comparison operator: %v", binOp))
		}

		cmp := builder.CreateICmp(predicate, lhs, rhs, "")
		return cmp
	}

	switch binOp {
	case token.STAR:
		return builder.CreateMul(lhs, rhs, "")
	case token.MINUS:
		return builder.CreateSub(lhs, rhs, "")
	case token.PLUS:
		return builder.CreateAdd(lhs, rhs, "")
	case token.SLASH:
		return builder.CreateExactSDiv(lhs, rhs, "")
	default:
		panic(fmt.Sprintf("unsupported binary operator: %v", binOp))
	}
}

func (c *codegen) emitFloatBinExpr(lhs llvm.Value, binOp token.Kind, rhs llvm.Value) llvm.Value {
	if binOp.IsCmpOp() {
		var predicate llvm.FloatPredicate

		switch binOp {
		case token.EQUAL_EQUAL:
			predicate = llvm.FloatOEQ // ordered and equal
		case token.BANG_EQUAL:
			predicate = llvm.FloatONE // ordered and not equal
		case token.LESS:
			predicate = llvm.FloatOLT // ordered and less than
		case token.LESS_EQ:
			predicate = llvm.FloatOLE // ordered and less equal
		case token.GREATER:
			predicate = llvm.FloatOGT // ordered and greater than
		case token.GREATER_EQ:
			predicate = llvm.FloatOGE // ordered and greater equal
		default:
			panic(fmt.Sprintf("unsupported comparison operator: %v", binOp))
		}

		return builder.CreateFCmp(predicate, lhs, rhs, "")
	}

	switch binOp {
	case token.STAR:
		return builder.CreateFMul(lhs, rhs, "")
	case token.MINUS:
		return builder.CreateFSub(lhs, rhs, "")
	case token.PLUS:
		return builder.CreateFAdd(lhs, rhs, "")
	case token.SLASH:
		return builder.CreateFDiv(lhs, rhs, "")
	default:
		panic(fmt.Sprintf("unsupported binary operator: %v", binOp))
	}
}

func (c *codegen) llvmConstInt8s(data []byte, length int) []llvm.Value {
	out := make([]llvm.Value, length)
	for i, b := range data {
		out[i] = llvm.ConstInt(context.Int8Type(), uint64(b), false)
	}
	// c-strings are null terminated
	out[length-1] = llvm.ConstInt(context.Int8Type(), uint64(0), false)
	return out
}

func (c *codegen) getIntegerValue(
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

func (c *codegen) getFloatValue(
	expr *ast.LiteralExpr,
	ty *ast.BasicType,
) (float64, int) {
	bitSize := ty.Kind.BitSize()
	floatLit := string(expr.Value)
	val, _ := strconv.ParseFloat(floatLit, bitSize)
	return val, bitSize
}

func (c *codegen) generateCondStmt(
	condStmt *ast.CondStmt,
	function llvm.Value,
) {
	ifBlock := llvm.AddBasicBlock(function, ".if")
	elseBlock := llvm.AddBasicBlock(function, ".else")
	endBlock := llvm.AddBasicBlock(function, ".end")

	// If
	ifExpr, _, _ := c.emitExpr(condStmt.IfStmt.Expr)
	builder.CreateCondBr(ifExpr, ifBlock, elseBlock)
	builder.SetInsertPointAtEnd(ifBlock)
	stoppedOnReturn := c.generateBlock(function, condStmt.IfStmt.Block)
	if !stoppedOnReturn {
		builder.CreateBr(endBlock)
	}

	// TODO: implement elif statements (at the end, they are just ifs and elses)

	// Else
	builder.SetInsertPointAtEnd(elseBlock)
	if condStmt.ElseStmt != nil {
		elseStoppedOnReturn := c.generateBlock(
			function,
			condStmt.ElseStmt.Block,
		)
		if !elseStoppedOnReturn {
			builder.CreateBr(endBlock)
		}
	} else {
		builder.CreateBr(endBlock)
	}
	builder.SetInsertPointAtEnd(endBlock)
}

func (c *codegen) emitNamespaceAccess(
	namespaceAccess *ast.NamespaceAccess,
) llvm.Value {
	if namespaceAccess.IsImport {
		return c.generateImportAccess(namespaceAccess.Right)
	}

	switch namespaceAccess.Left.N.Kind {
	case ast.KIND_EXTERN_DECL:
		fnCall := namespaceAccess.Right.Node.(*ast.FnCall)
		return c.emitFnCall(fnCall)
	default:
		panic("unreachable")
	}
}

func (c *codegen) generateImportAccess(right *ast.Node) llvm.Value {
	switch right.Kind {
	case ast.KIND_FN_CALL:
		fnCall := right.Node.(*ast.FnCall)
		return c.emitFnCall(fnCall)
	case ast.KIND_NAMESPACE_ACCESS:
		namespaceAccess := right.Node.(*ast.NamespaceAccess)
		return c.emitNamespaceAccess(namespaceAccess)
	case ast.KIND_STRUCT_LITERAL_EXPR:
		stLit := right.Node.(*ast.StructLiteralExpr)
		return c.emitStructLiteral(stLit)
	default:
		panic(fmt.Sprintf("unimplemented import access: %v\n", right.Kind))
	}
}

func (c *codegen) generateForLoop(
	forLoop *ast.ForLoop,
	function llvm.Value,
) {
	forPrepBlock := llvm.AddBasicBlock(function, ".forprep")
	forInitBlock := llvm.AddBasicBlock(function, ".forinit")
	forBodyBlock := llvm.AddBasicBlock(function, ".forbody")
	forUpdateBlock := llvm.AddBasicBlock(function, ".forupdate")
	endBlock := llvm.AddBasicBlock(function, ".forend")

	builder.CreateBr(forPrepBlock)
	builder.SetInsertPointAtEnd(forPrepBlock)
	c.generateStmt(function, forLoop.Block, forLoop.Init)

	builder.CreateBr(forInitBlock)
	builder.SetInsertPointAtEnd(forInitBlock)

	expr, _, _ := c.emitExpr(forLoop.Cond)

	builder.CreateCondBr(expr, forBodyBlock, endBlock)

	builder.SetInsertPointAtEnd(forBodyBlock)
	c.generateBlock(function, forLoop.Block)

	builder.CreateBr(forUpdateBlock)
	builder.SetInsertPointAtEnd(forUpdateBlock)
	c.generateStmt(function, forLoop.Block, forLoop.Update)

	builder.CreateBr(forInitBlock)

	builder.SetInsertPointAtEnd(endBlock)
}

func (c *codegen) generateWhileLoop(
	whileLoop *ast.WhileLoop,
	function llvm.Value,
) {
	whileInitBlock := llvm.AddBasicBlock(function, ".whileinit")
	whileBodyBlock := llvm.AddBasicBlock(function, ".whilebody")
	endBlock := llvm.AddBasicBlock(function, ".whileend")

	builder.CreateBr(whileInitBlock)
	builder.SetInsertPointAtEnd(whileInitBlock)
	expr, _, _ := c.emitExpr(whileLoop.Cond)
	builder.CreateCondBr(expr, whileBodyBlock, endBlock)

	builder.SetInsertPointAtEnd(whileBodyBlock)
	c.generateBlock(function, whileLoop.Block)
	builder.CreateBr(whileInitBlock)

	builder.SetInsertPointAtEnd(endBlock)
}

func (c *codegen) generateExe(buildType config.BuildType) error {
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

	if err := llvm.VerifyModule(c.module, llvm.AbortProcessAction); err != nil {
		fmt.Println("error: module is not valid")
		return err
	}

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
		compilerFlags += "-Wl,-s"
	case config.DEBUG:
		optLevel = "-O0"
	default:
		panic("invalid build type: " + buildType.String())
	}

	optCommand := []string{
		optLevel,
		"-o",
		optimizedIrFilepath,
		irFilepath,
	}
	cmd := exec.Command("opt", optCommand...)
	err = cmd.Run()
	// TODO(errors)
	if err != nil {
		if config.DEBUG_MODE {
			fmt.Printf("[DEBUG MODE] OPT COMMAND: %s\n", cmd)
		}
		return err
	}

	clangCommand := []string{
		compilerFlags,
		optLevel,
		"-o",
		filenameNoExt,
		optimizedIrFilepath,
		"-lm", // math library
	}
	cmd = exec.Command("clang-18", clangCommand...)
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
