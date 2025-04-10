package llvm

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"slices"
	"strconv"
	"strings"

	"github.com/HicaroD/Telia/config"
	"github.com/HicaroD/Telia/internal/ast"
	"github.com/HicaroD/Telia/internal/lexer/token"
	"tinygo.org/x/go-llvm"
)

var (
	context llvm.Context = llvm.NewContext()
	builder llvm.Builder = context.NewBuilder()
)

// B stands for (B)ackend
var (
	B_LIT_FALSE = llvm.ConstInt(B_BOOL_TYPE, 0, false)
	B_LIT_TRUE  = llvm.ConstInt(B_BOOL_TYPE, 1, false)

	B_VOID_TYPE   = context.VoidType()
	B_BOOL_TYPE   = context.Int1Type()
	B_INT8_TYPE   = context.Int8Type()
	B_INT32_TYPE  = context.Int32Type()
	B_F32_TYPE    = context.FloatType()
	B_F64_TYPE    = context.DoubleType()
	B_RAWPTR_TYPE = llvm.PointerType(B_INT8_TYPE, 0) // *u8
)

type codegen struct {
	module llvm.Module

	loc     *ast.Loc
	program *ast.Program
	runtime *ast.Package
	pkg     *ast.Package
}

func NewCG(loc *ast.Loc, program *ast.Program, runtime *ast.Package) *codegen {
	module := context.NewModule(loc.Dir)

	defaultTargetTriple := llvm.DefaultTargetTriple()
	module.SetTarget(defaultTargetTriple)

	return &codegen{
		loc:     loc,
		program: program,
		module:  module,
		runtime: runtime,
	}
}

func (c *codegen) Generate(buildType config.BuildType) error {
	c.generatePackage(c.runtime)
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
			c.emitFnBody(node.Node.(*ast.FnDecl))
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

	functionValue.SetFunctionCallConv(
		c.getCallingConvention(fnDecl.Attributes.DefaultCallingConvention),
	)
}

func (c *codegen) emitFnBody(fnDecl *ast.FnDecl) {
	fn := c.module.NamedFunction(fnDecl.Name.Name())
	if fn.IsNil() {
		panic("function is nil when generating fn body")
	}
	body := context.AddBasicBlock(fn, "entry")
	builder.SetInsertPointAtEnd(body)

	paramsTypes := fn.GlobalValueType().ParamTypes()
	for i, fnParam := range fn.Params() {
		argPtr := builder.CreateAlloca(paramsTypes[i], "")
		builder.CreateStore(fnParam, argPtr)
		variable := &Variable{
			Ty:  paramsTypes[i],
			Ptr: argPtr,
		}
		field := fnDecl.Params.Fields[i]
		field.BackendType = variable
	}

	stoppedOnReturn := c.emitBlock(fn, fnDecl.Block)
	if !stoppedOnReturn && fnDecl.RetType.IsVoid() {
		builder.CreateRetVoid()
	}
}

func (c *codegen) emitBlock(
	fn llvm.Value,
	block *ast.BlockStmt,
) bool {
	stoppedOnReturn := false

	for _, stmt := range block.Statements {
		c.emitStmt(fn, block, stmt)
		if stmt.IsReturn() {
			stoppedOnReturn = true
			break
		}
	}

	return stoppedOnReturn
}

func (c *codegen) emitStmt(
	function llvm.Value,
	block *ast.BlockStmt,
	stmt *ast.Node,
) {
	switch stmt.Kind {
	case ast.KIND_FN_CALL:
		c.emitFnCall(stmt.Node.(*ast.FnCall))
	case ast.KIND_RETURN_STMT:
		if len(block.DeferStack) > 0 {
			c.emitDefers(function, block, block.DeferStack)
		}
		c.emitReturn(stmt.Node.(*ast.ReturnStmt))
	case ast.KIND_VAR_STMT:
		c.emitVar(stmt.Node.(*ast.VarStmt))
	case ast.KIND_NAMESPACE_ACCESS:
		c.emitNamespaceAccess(stmt.Node.(*ast.NamespaceAccess))
	case ast.KIND_COND_STMT:
		c.emitCond(stmt.Node.(*ast.CondStmt), function)
	case ast.KIND_FOR_LOOP_STMT:
		c.emitForLoop(stmt.Node.(*ast.ForLoop), function)
	case ast.KIND_WHILE_LOOP_STMT:
		c.emitWhileLoop(stmt.Node.(*ast.WhileLoop), function)
	case ast.KIND_ASSIGNMENT_STMT:
		c.emitAssignment(stmt.Node.(*ast.AssignmentStmt))
	case ast.KIND_DEFER_STMT:
		break
	default:
		log.Fatalf("unimplemented block statement: %s\n", stmt)
	}
}

func (c *codegen) emitDefers(fn llvm.Value, block *ast.BlockStmt, defers []*ast.DeferStmt) {
	for i := len(defers) - 1; i >= 0; i-- {
		if defers[i].Skip {
			continue
		}
		c.emitStmt(fn, block, defers[i].Stmt)
	}
}

func (c *codegen) emitReturn(
	ret *ast.ReturnStmt,
) {
	if ret.Value.IsVoid() {
		builder.CreateRetVoid()
		return
	}
	_, v, _ := c.emitExprWithLoadIfNeeded(ret.Value)
	// TODO: learn more about noundef for return values
	builder.CreateRet(v)
}

func (c *codegen) emitVar(variable *ast.VarStmt) {
	switch variable.Expr.Kind {
	case ast.KIND_TUPLE_LITERAL_EXPR:
		tuple := variable.Expr.Node.(*ast.TupleExpr)
		c.emitTupleLiteralAsVarValue(variable, tuple)
	case ast.KIND_FN_CALL:
		fnCall := variable.Expr.Node.(*ast.FnCall)
		c.emitFnCallForTuple(variable.Names, fnCall, variable.IsDecl)
	default:
		c.emitVariable(
			variable.Names[0],
			variable.Expr,
			variable.IsDecl,
		)
	}
}

func (c *codegen) emitAssignment(variable *ast.AssignmentStmt) {
	t := 0
	for _, value := range variable.Values {
		switch value.Kind {
		case ast.KIND_FN_CALL:
			fnCall := value.Node.(*ast.FnCall)
			if fnCall.Decl.RetType.Kind == ast.EXPR_TYPE_TUPLE {
				tupleType := fnCall.Decl.RetType.T.(*ast.TupleType)
				affectedVariables := variable.Targets[t : t+len(tupleType.Types)]
				c.emitFnCallForTuple(affectedVariables, fnCall, variable.Decl)
				t += len(affectedVariables)
			} else {
				c.emitVariable(variable.Targets[t], value, variable.Decl)
				t++
			}
		default:
			c.emitVariable(variable.Targets[t], value, variable.Decl)
			t++
		}
	}
}

func (c *codegen) emitTupleLiteralAsVarValue(variable *ast.VarStmt, tuple *ast.TupleExpr) {
	t := 0

	for _, expr := range tuple.Exprs {
		switch expr.Kind {
		case ast.KIND_TUPLE_LITERAL_EXPR:
			innerTupleExpr := expr.Node.(*ast.TupleExpr)
			for _, innerExpr := range innerTupleExpr.Exprs {
				c.emitVariable(
					variable.Names[t],
					innerExpr,
					variable.IsDecl,
				)
				t++
			}
		case ast.KIND_FN_CALL:
			fnCall := expr.Node.(*ast.FnCall)
			if fnCall.Decl.RetType.Kind == ast.EXPR_TYPE_TUPLE {
				tupleType := fnCall.Decl.RetType.T.(*ast.TupleType)
				affectedVariables := variable.Names[t : t+len(tupleType.Types)]
				c.emitFnCallForTuple(affectedVariables, fnCall, variable.IsDecl)
				t += len(affectedVariables)
			} else {
				c.emitVariable(variable.Names[t], expr, variable.IsDecl)
				t++
			}
		default:
			c.emitVariable(variable.Names[t], expr, variable.IsDecl)
			t++
		}
	}
}

func (c *codegen) emitFnCallForTuple(variables []*ast.Node, fnCall *ast.FnCall, isDecl bool) {
	_, genFnCall, _ := c.emitFnCall(fnCall)

	if len(variables) == 1 {
		c.emitVariableWithValue(variables[0], genFnCall, isDecl)
	} else {
		for i, currentVar := range variables {
			value := builder.CreateExtractValue(genFnCall, i, "")
			c.emitVariableWithValue(currentVar, value, isDecl)
		}
	}
}

func (c *codegen) emitVariable(name *ast.Node, expr *ast.Node, isDecl bool) {
	if isDecl {
		c.emitVarDecl(name, expr)
	} else {
		c.emitVarReassign(name, expr)
	}
}

func (c *codegen) emitVariableWithValue(name *ast.Node, value llvm.Value, isDecl bool) {
	if isDecl {
		c.emitVarDeclWithValue(name, value)
	} else {
		c.emitVarReassignWithValue(name, value)
	}
}

func (c *codegen) emitVarDecl(
	name *ast.Node,
	expr *ast.Node,
) {
	var ty llvm.Type
	var ptr llvm.Value

	variable := name.Node.(*ast.VarIdStmt)
	ty = c.emitType(variable.Type)

	switch expr.Kind {
	case ast.KIND_STRUCT_EXPR:
		_, ptr, _ = c.emitStructLiteral(expr.Node.(*ast.StructLiteralExpr))
	case ast.KIND_NAMESPACE_ACCESS:
		_, ptr, _ = c.emitNamespaceAccess(expr.Node.(*ast.NamespaceAccess))
	default:
		ptr = builder.CreateAlloca(ty, "")
		_, val, _ := c.emitExprWithLoadIfNeeded(expr)
		builder.CreateStore(val, ptr)
	}

	variableLlvm := &Variable{
		Ty:  ty,
		Ptr: ptr,
	}
	variable.BackendType = variableLlvm
}

func (c *codegen) emitVarReassign(
	name *ast.Node,
	expr *ast.Node,
) {
	var p llvm.Value

	switch name.Kind {
	case ast.KIND_VAR_ID_STMT:
		varId := name.Node.(*ast.VarIdStmt)

		var v *Variable

		switch varId.N.Kind {
		case ast.KIND_VAR_ID_STMT:
			variable := varId.N.Node.(*ast.VarIdStmt)
			v = variable.BackendType.(*Variable)
		case ast.KIND_PARAM:
			param := varId.N.Node.(*ast.Param)
			v = param.BackendType.(*Variable)
		default:
			panic(fmt.Sprintf("unimplemented kind of name expression: %v\n", varId.N))
		}

		p = v.Ptr
	case ast.KIND_FIELD_ACCESS:
		f := name.Node.(*ast.FieldAccess)
		_, p = c.getStructFieldPtr(f)
	case ast.KIND_DEREF_POINTER_EXPR:
		_, p, _ = c.emitExpr(name)
	default:
		log.Fatalf("invalid symbol on generateVarReassign: %v\n", name)
	}

	if p.IsNil() {
		panic("variable ptr is nil")
	}

	_, e, _ := c.emitExprWithLoadIfNeeded(expr)
	builder.CreateStore(e, p)
}

func (c *codegen) getStructFieldPtr(fieldAccess *ast.FieldAccess) (llvm.Type, llvm.Value) {
	st := c.module.GetTypeByName(fieldAccess.Decl.Name.Name())
	if st.IsNil() {
		panic("struct is nil")
	}
	fTy := fieldAccess.AccessedField.Type

	switch fieldAccess.Right.Kind {
	case ast.KIND_ID_EXPR:
		var objTy *ast.ExprType
		var obj *Variable

		if fieldAccess.Param {
			objTy = fieldAccess.StructParam.Type
			obj = fieldAccess.StructParam.BackendType.(*Variable)
		} else {
			objTy = fieldAccess.StructVar.Type
			obj = fieldAccess.StructVar.BackendType.(*Variable)
		}

		// NOTE: it only allows a single pointer deref
		ptr := obj.Ptr
		if objTy.IsPointer() {
			ty := c.emitType(objTy)
			ptr = builder.CreateLoad(ty, ptr, "")
		}
		gep := builder.CreateStructGEP(
			st,
			ptr,
			fieldAccess.AccessedField.Index,
			"",
		)
		return c.emitType(fTy), gep
	default:
		panic("unimplemented other type of field access")
	}
}

func (c *codegen) emitVarDeclWithValue(
	name *ast.Node,
	value llvm.Value,
) {
	variable := name.Node.(*ast.VarIdStmt)

	ty := c.emitType(variable.Type)
	ptr := builder.CreateAlloca(ty, "")
	builder.CreateStore(value, ptr)

	variableLlvm := &Variable{
		Ty:  ty,
		Ptr: ptr,
	}

	variable.BackendType = variableLlvm
}

func (c *codegen) emitVarReassignWithValue(
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
		case ast.KIND_PARAM:
			param := node.N.Node.(*ast.Param)
			variable = param.BackendType.(*Variable)
		}
	default:
		panic(fmt.Sprintf("unimplemented symbol on generateVarReassign: %v\n", name.Node))
	}
	builder.CreateStore(value, variable.Ptr)
}

func (c *codegen) emitFnCall(call *ast.FnCall) (llvm.Type, llvm.Value, bool) {
	args := c.emitCallArgs(call)

	name := c.getFnCallName(call)
	fn := c.module.NamedFunction(name)
	if fn.IsNil() {
		panic("function is nil when generating function call")
	}

	var hasFloat bool
	if call.Decl != nil {
		hasFloat, _ = c.checkFloatTypeForBitSize(call.Decl.RetType)
	} else {
		hasFloat, _ = c.checkFloatTypeForBitSize(call.Proto.RetType)
	}

	ty := fn.GlobalValueType()
	return ty, builder.CreateCall(ty, fn, args, ""), hasFloat
}

func (c *codegen) getFnCallName(call *ast.FnCall) string {
	var attrs ast.Attributes
	if call.Proto != nil {
		attrs = call.Proto.Attributes
	} else {
		attrs = call.Decl.Attributes
	}

	if attrs.LinkName != "" {
		return attrs.LinkName
	}
	return call.Name.Name()
}

func (c *codegen) emitTupleExpr(tuple *ast.TupleExpr) (llvm.Type, llvm.Value, bool) {
	types := c.emitTypes(tuple.Type.Types)
	tupleTy := context.StructType(types, false)
	tupleVal := llvm.ConstNull(tupleTy)

	for i, expr := range tuple.Exprs {
		_, e, _ := c.emitExprWithLoadIfNeeded(expr)
		tupleVal = builder.CreateInsertValue(tupleVal, e, i, "")
	}

	return tupleTy, tupleVal, false
}

func (c *codegen) emitStructDecl(st *ast.StructDecl) {
	fields := c.getStructFieldList(st.Fields)
	structTy := context.StructCreateNamed(st.Name.Name())
	// TODO: set packed properly
	packed := false
	structTy.StructSetBody(fields, packed)
}

func (c *codegen) emitStructLiteral(lit *ast.StructLiteralExpr) (llvm.Type, llvm.Value, bool) {
	st := c.module.GetTypeByName(lit.Name.Name())
	if st.IsNil() {
		panic("struct is nil")
	}

	values := make([]llvm.Value, len(lit.Values))
	for _, field := range lit.Values {
		_, val, _ := c.emitExpr(field.Value)
		values[field.Index] = val
	}

	stConst := llvm.ConstNamedStruct(st, values)

	globalConst := llvm.AddGlobal(c.module, st, "")
	globalConst.SetInitializer(stConst)
	globalConst.SetLinkage(llvm.PrivateLinkage)
	globalConst.SetGlobalConstant(true)
	globalConst.SetUnnamedAddr(true)

	structAlloca := builder.CreateAlloca(st, "")
	sizeOf := c.emitSizeOfType(st)
	c.emitMemcpy(structAlloca, globalConst, sizeOf)
	return st, structAlloca, false
}

func (c *codegen) emitSizeOfType(ty llvm.Type) llvm.Value {
	nullPtr := llvm.ConstNull(llvm.PointerType(ty, 0))
	gep := llvm.ConstGEP(ty, nullPtr, []llvm.Value{
		llvm.ConstInt(context.Int32Type(), 1, false),
	})
	size := builder.CreatePtrToInt(gep, context.Int64Type(), "")
	return size
}

func (c *codegen) emitMemcpy(dest, src, size llvm.Value) {
	c.emitRuntimeCall("llvm.memcpy.p0.p0.i64",
		[]llvm.Value{
			dest, // destination (ptr)
			src,  // source (ptr)
			size, // size in bytes
			llvm.ConstInt(context.Int1Type(), 0, false), // isVolatile
		},
	)
}

func (c *codegen) emitFieldAccessExpr(fieldAccess *ast.FieldAccess) (llvm.Type, llvm.Value, bool) {
	ty, ptr := c.getStructFieldPtr(fieldAccess)
	hasFloat, _ := c.checkFloatTypeForBitSize(fieldAccess.AccessedField.Type)
	return ty, ptr, hasFloat
}

func (c *codegen) emitNullPtrExpr(nullPtr *ast.NullPtrExpr) (llvm.Type, llvm.Value, bool) {
	ty := c.emitType(nullPtr.Type)
	return ty, llvm.ConstNull(ty), false
}

func (c *codegen) emitPtrExpr(ptr *ast.AddressOfExpr) (llvm.Type, llvm.Value, bool) {
	var t llvm.Type
	var p llvm.Value

	switch ptr.Expr.Kind {
	case ast.KIND_ID_EXPR:
		id := ptr.Expr.Node.(*ast.IdExpr)

		switch id.N.Kind {
		case ast.KIND_VAR_ID_STMT:
			variable := id.N.Node.(*ast.VarIdStmt)
			v := variable.BackendType.(*Variable)
			t = v.Ty
			p = v.Ptr
		case ast.KIND_PARAM:
			variable := id.N.Node.(*ast.Param)
			v := variable.BackendType.(*Variable)
			t = v.Ty
			p = v.Ptr
		}
	case ast.KIND_FIELD_ACCESS:
		fieldAccess := ptr.Expr.Node.(*ast.FieldAccess)
		_, p = c.getStructFieldPtr(fieldAccess)
	}

	return t, p, false

}

func (c *codegen) emitExternDecl(external *ast.ExternDecl) {
	for _, proto := range external.Prototypes {
		c.emitPrototype(external.Attributes, proto)
	}
}

func (c *codegen) emitPrototype(externAttributes ast.Attributes, prototype *ast.Proto) {
	returnTy := c.emitType(prototype.RetType)
	paramsTypes := c.getFieldListTypes(prototype.Params)
	ty := llvm.FunctionType(returnTy, paramsTypes, prototype.Params.IsVariadic)
	protoValue := llvm.AddFunction(c.module, prototype.Name.Name(), ty)

	defaultCc := c.getCallingConvention(externAttributes.DefaultCallingConvention)
	protoValue.SetFunctionCallConv(defaultCc)

	if prototype.Attributes.DefaultCallingConvention != "" {
		protoValue.SetLinkage(c.getFunctionLinkage(prototype.Attributes.Linkage))
	}
	if prototype.Attributes.LinkName != "" {
		protoValue.SetName(prototype.Attributes.LinkName)
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
		case token.BOOL_TYPE, token.INT_TYPE,
			token.UINT_TYPE, token.I8_TYPE, token.U8_TYPE,
			token.I16_TYPE, token.U16_TYPE, token.I32_TYPE, token.U32_TYPE,
			token.I64_TYPE, token.U64_TYPE, token.I128_TYPE, token.U128_TYPE:

			bitSize := b.Kind.BitSize()
			return context.IntType(bitSize)
		case token.F32_TYPE, token.FLOAT_TYPE:
			return B_F32_TYPE
		case token.F64_TYPE:
			return B_F64_TYPE
		case token.VOID_TYPE:
			return B_VOID_TYPE
		case token.STRING_TYPE, token.CSTRING_TYPE:
			u8 := ast.NewBasicType(token.U8_TYPE)
			u8Type := c.emitType(u8)
			return c.emitPtrType(u8Type)
		case token.RAWPTR_TYPE:
			return B_RAWPTR_TYPE
		// case token.UNTYPED_BOOL:
		// 	panic("unimplemented untyped bool")
		// case token.UNTYPED_INT:
		// 	panic("unimplemented untyped int")
		// case token.UNTYPED_FLOAT:
		// 	panic("unimplemented untyped float")
		// case token.UNTYPED_STRING:
		// 	panic("unimplemented untyped string")
		case token.UNTYPED_NULLPTR:
			panic("unimplemented untyped nullptr")
		default:
			panic(fmt.Sprintf("unimplemented basic type: %v\n", ty.Kind))
		}
	case ast.EXPR_TYPE_TUPLE:
		tuple := ty.T.(*ast.TupleType)
		innerTypes := c.emitTypes(tuple.Types)
		tupleTy := context.StructType(innerTypes, false)
		return tupleTy
	case ast.EXPR_TYPE_STRUCT:
		ty := ty.T.(*ast.StructType)
		return c.module.GetTypeByName(ty.Decl.Name.Name())
	case ast.EXPR_TYPE_POINTER:
		ptr := ty.T.(*ast.PointerType)
		ptrTy := c.emitType(ptr.Type)
		return c.emitPtrType(ptrTy)
	default:
		panic(fmt.Sprintf("unimplemented type: %v\n", ty.T))
	}
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

func (c *codegen) emitCallArgs(call *ast.FnCall) []llvm.Value {
	nExprs := len(call.Args)
	if call.Variadic {
		nExprs-- // desconsider the variadic argument place for now
	}

	values := make([]llvm.Value, nExprs)
	for i := range nExprs {
		_, v, _ := c.emitExprWithLoadIfNeeded(call.Args[i])
		values[i] = v
	}

	if call.Variadic {
		variadic := call.Args[nExprs].Node.(*ast.VarArgsExpr)
		for _, arg := range variadic.Args {
			_, v, _ := c.emitExprWithLoadIfNeeded(arg)
			values = append(values, v)
		}
	}

	return values
}

var NO_PRELOAD = []ast.NodeKind{
	ast.KIND_ADDRESS_OF_EXPR,
	ast.KIND_LITERAL_EXPR,
	ast.KIND_FN_CALL,
	ast.KIND_BINARY_EXPR,
	ast.KIND_UNARY_EXPR,
	ast.KIND_TUPLE_LITERAL_EXPR,
	ast.KIND_NULLPTR_EXPR,
}

func (c *codegen) emitExprWithLoadIfNeeded(expr *ast.Node) (llvm.Type, llvm.Value, bool) {
	ty, val, hasFloat := c.emitExpr(expr)

	noPreload := slices.Contains(NO_PRELOAD, expr.Kind)
	if noPreload {
		return ty, val, hasFloat
	}

	l := builder.CreateLoad(ty, val, "")
	return ty, l, hasFloat
}

func (c *codegen) emitExpr(expr *ast.Node) (llvm.Type, llvm.Value, bool) {
	switch expr.Kind {
	case ast.KIND_LITERAL_EXPR:
		return c.emitLiteralExpr(expr.Node.(*ast.LiteralExpr))
	case ast.KIND_ID_EXPR:
		return c.emitIdExpr(expr.Node.(*ast.IdExpr))
	case ast.KIND_BINARY_EXPR:
		return c.emitBinExpr(expr.Node.(*ast.BinExpr))
	case ast.KIND_FN_CALL:
		return c.emitFnCall(expr.Node.(*ast.FnCall))
	case ast.KIND_UNARY_EXPR:
		return c.emitUnaryExpr(expr.Node.(*ast.UnaryExpr))
	case ast.KIND_NAMESPACE_ACCESS:
		return c.emitNamespaceAccess(expr.Node.(*ast.NamespaceAccess))
	case ast.KIND_TUPLE_LITERAL_EXPR:
		return c.emitTupleExpr(expr.Node.(*ast.TupleExpr))
	case ast.KIND_STRUCT_EXPR:
		return c.emitStructLiteral(expr.Node.(*ast.StructLiteralExpr))
	case ast.KIND_FIELD_ACCESS:
		return c.emitFieldAccessExpr(expr.Node.(*ast.FieldAccess))
	case ast.KIND_DEREF_POINTER_EXPR:
		return c.emitDerefPtrExpr(expr.Node.(*ast.DerefPointerExpr))
	case ast.KIND_ADDRESS_OF_EXPR:
		return c.emitPtrExpr(expr.Node.(*ast.AddressOfExpr))
	case ast.KIND_NULLPTR_EXPR:
		return c.emitNullPtrExpr(expr.Node.(*ast.NullPtrExpr))
	case ast.KIND_VARG_EXPR:
		panic("unimplemented var args")
	default:
		panic(fmt.Sprintf("unimplemented expr: %v", expr))
	}
}

func (c *codegen) emitLiteralExpr(lit *ast.LiteralExpr) (llvm.Type, llvm.Value, bool) {
	switch lit.Type.Kind {
	case ast.EXPR_TYPE_BASIC:
		basic := lit.Type.T.(*ast.BasicType)

		if basic.Kind.IsFloat() {
			return c.emitFloat(lit, basic)
		}

		if basic.Kind.IsInteger() {
			integerValue, bitSize := c.getIntegerValue(lit, basic)
			ty := context.IntType(bitSize)
			return ty, llvm.ConstInt(ty, integerValue, false), false
		}

		switch basic.Kind {
		case token.STRING_TYPE, token.CSTRING_TYPE:
			ty, val := c.emitCstringLit(lit.Value)
			return ty, val, false
		case token.BOOL_TYPE:
			ty, val := c.emitBool(lit)
			return ty, val, false
		default:
			panic(
				fmt.Sprintf(
					"unimplemented basic type: %s %s\n",
					basic.Kind.String(),
					string(lit.Value),
				),
			)
		}
	case ast.EXPR_TYPE_POINTER:
		panic("unimplemented pointer type")
	default:
		panic("unimplemented literal type")
	}
}

func (c *codegen) emitFloat(lit *ast.LiteralExpr, ty *ast.BasicType) (llvm.Type, llvm.Value, bool) {
	hasFloat := true
	floatValue, bitSize := c.getFloatValue(lit, ty)
	if bitSize == 32 {
		return B_F32_TYPE, llvm.ConstFloat(B_F32_TYPE, floatValue), hasFloat
	} else {
		return B_F64_TYPE, llvm.ConstFloat(B_F64_TYPE, floatValue), hasFloat
	}
}

func (c *codegen) emitCstringLit(str []byte) (llvm.Type, llvm.Value) {
	strlen := len(str) + 1
	arrTy := llvm.ArrayType(B_INT8_TYPE, strlen)
	arr := llvm.ConstArray(arrTy, c.emitConstInt8sForCstring(str, strlen))

	globalVal := llvm.AddGlobal(c.module, arrTy, ".str")
	globalVal.SetInitializer(arr)
	globalVal.SetLinkage(llvm.PrivateLinkage)
	globalVal.SetGlobalConstant(true)
	globalVal.SetAlignment(1)
	globalVal.SetUnnamedAddr(true)

	zero := llvm.ConstInt(B_INT32_TYPE, 0, false)
	indices := []llvm.Value{zero, zero}
	ptr := llvm.ConstInBoundsGEP(arrTy, globalVal, indices)
	return arrTy, ptr
}

func (c *codegen) emitBool(lit *ast.LiteralExpr) (llvm.Type, llvm.Value) {
	boolVal := string(lit.Value)

	var val uint64
	if boolVal == "true" {
		val = 1
	} else {
		val = 0
	}

	return B_BOOL_TYPE, llvm.ConstInt(B_BOOL_TYPE, val, false)
}

func (c *codegen) emitIdExpr(id *ast.IdExpr) (llvm.Type, llvm.Value, bool) {
	var localVar *Variable
	var varTy *ast.ExprType

	hasFloat := false

	switch id.N.Kind {
	case ast.KIND_VAR_ID_STMT:
		variable := id.N.Node.(*ast.VarIdStmt)
		varTy = variable.Type
		localVar = variable.BackendType.(*Variable)
	case ast.KIND_PARAM:
		variable := id.N.Node.(*ast.Param)
		varTy = variable.Type
		localVar = variable.BackendType.(*Variable)
	}

	if varTy.IsNumeric() {
		bt := varTy.T.(*ast.BasicType)
		hasFloat = bt.Kind.IsFloat()
	}

	return localVar.Ty, localVar.Ptr, hasFloat
}

func (c *codegen) emitDerefPtrExpr(deref *ast.DerefPointerExpr) (llvm.Type, llvm.Value, bool) {
	ty, v, hasFloat := c.emitExprWithLoadIfNeeded(deref.Expr)
	c.emitRuntimeCall("_check_nil_pointer_deref", []llvm.Value{v})
	return ty, v, hasFloat
}

func (c *codegen) emitBinExpr(bin *ast.BinExpr) (llvm.Type, llvm.Value, bool) {
	tLhs, vLhs, lhsHasFloat := c.emitExprWithLoadIfNeeded(bin.Left)
	_, vRhs, rhsHasFloat := c.emitExprWithLoadIfNeeded(bin.Right)
	hasFloat := lhsHasFloat || rhsHasFloat

	var binExpr llvm.Value
	if hasFloat {
		binExpr = c.emitFloatBinExpr(vLhs, bin.Op, vRhs)
	} else {
		binExpr = c.emitIntBinExpr(vLhs, bin.Op, vRhs)
	}
	binTy := tLhs
	return binTy, binExpr, hasFloat
}

func (c *codegen) emitUnaryExpr(unary *ast.UnaryExpr) (llvm.Type, llvm.Value, bool) {
	ty, e, hasFloat := c.emitExprWithLoadIfNeeded(unary.Value)
	if hasFloat && unary.Op == token.MINUS {
		return ty, builder.CreateFNeg(e, ""), hasFloat
	}

	switch unary.Op {
	case token.MINUS:
		return ty, builder.CreateNeg(e, ""), hasFloat
	case token.NOT:
		return ty, builder.CreateICmp(llvm.IntEQ, B_LIT_FALSE, e, ""), hasFloat
	default:
		panic(fmt.Sprintf("unimplemented unary operator: %s", unary.Op))
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
	case token.AND:
		return builder.CreateAnd(lhs, rhs, "")
	case token.OR:
		return builder.CreateOr(lhs, rhs, "")
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
	case token.AND:
		return builder.CreateAnd(lhs, rhs, "")
	case token.OR:
		return builder.CreateOr(lhs, rhs, "")
	default:
		panic(fmt.Sprintf("unsupported binary operator: %v", binOp))
	}
}

func (c *codegen) emitConstInt8sForCstring(data []byte, length int) []llvm.Value {
	out := make([]llvm.Value, length)
	for i, b := range data {
		out[i] = llvm.ConstInt(B_INT8_TYPE, uint64(b), false)
	}
	// c-strings are null terminated
	out[length-1] = llvm.ConstInt(B_INT8_TYPE, uint64(0), false)
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

func (c *codegen) emitCond(
	condStmt *ast.CondStmt,
	function llvm.Value,
) {
	ifBlock := llvm.AddBasicBlock(function, ".if")
	elseBlock := llvm.AddBasicBlock(function, ".else")
	endBlock := llvm.AddBasicBlock(function, ".end")

	// If
	_, ifExpr, _ := c.emitExprWithLoadIfNeeded(condStmt.IfStmt.Expr)
	builder.CreateCondBr(ifExpr, ifBlock, elseBlock)
	builder.SetInsertPointAtEnd(ifBlock)
	stoppedOnReturn := c.emitBlock(function, condStmt.IfStmt.Block)
	if !stoppedOnReturn {
		builder.CreateBr(endBlock)
	}

	// TODO: implement elif statements (at the end, they are just ifs and elses)

	// Else
	builder.SetInsertPointAtEnd(elseBlock)
	if condStmt.ElseStmt != nil {
		elseStoppedOnReturn := c.emitBlock(
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
) (llvm.Type, llvm.Value, bool) {
	if namespaceAccess.IsImport {
		return c.emitImportAccess(namespaceAccess.Right)
	}

	if namespaceAccess.Left.N.Kind != ast.KIND_EXTERN_DECL {
		panic("expected extern declaration for left side of namespace access")
	}

	fnCall := namespaceAccess.Right.Node.(*ast.FnCall)
	return c.emitFnCall(fnCall)
}

func (c *codegen) emitImportAccess(right *ast.Node) (llvm.Type, llvm.Value, bool) {
	switch right.Kind {
	case ast.KIND_FN_CALL:
		fnCall := right.Node.(*ast.FnCall)
		return c.emitFnCall(fnCall)
	case ast.KIND_NAMESPACE_ACCESS:
		namespaceAccess := right.Node.(*ast.NamespaceAccess)
		return c.emitNamespaceAccess(namespaceAccess)
	case ast.KIND_STRUCT_EXPR:
		stLit := right.Node.(*ast.StructLiteralExpr)
		return c.emitStructLiteral(stLit)
	default:
		panic(fmt.Sprintf("unimplemented import access: %v\n", right.Kind))
	}
}

func (c *codegen) emitForLoop(
	forLoop *ast.ForLoop,
	function llvm.Value,
) {
	forPrepBlock := llvm.AddBasicBlock(function, ".forprep")
	forInitBlock := llvm.AddBasicBlock(function, ".forinit")
	forBodyBlock := llvm.AddBasicBlock(function, ".forbody")
	forUpdateBlock := llvm.AddBasicBlock(function, ".forupdate")
	endBlock := llvm.AddBasicBlock(function, ".forend")

	// INIT
	builder.CreateBr(forPrepBlock)
	builder.SetInsertPointAtEnd(forPrepBlock)
	c.emitStmt(function, forLoop.Block, forLoop.Init)
	builder.CreateBr(forInitBlock)
	builder.SetInsertPointAtEnd(forInitBlock)

	// COND
	_, expr, _ := c.emitExprWithLoadIfNeeded(forLoop.Cond)
	builder.CreateCondBr(expr, forBodyBlock, endBlock)

	// FOR LOOP BLOCK
	builder.SetInsertPointAtEnd(forBodyBlock)
	c.emitBlock(function, forLoop.Block)

	// UPDATE
	builder.CreateBr(forUpdateBlock)
	builder.SetInsertPointAtEnd(forUpdateBlock)
	c.emitStmt(function, forLoop.Block, forLoop.Update)
	builder.CreateBr(forInitBlock)

	// END
	builder.SetInsertPointAtEnd(endBlock)
}

func (c *codegen) emitWhileLoop(
	whileLoop *ast.WhileLoop,
	function llvm.Value,
) {
	whileInitBlock := llvm.AddBasicBlock(function, ".whileinit")
	whileBodyBlock := llvm.AddBasicBlock(function, ".whilebody")
	endBlock := llvm.AddBasicBlock(function, ".whileend")

	builder.CreateBr(whileInitBlock)
	builder.SetInsertPointAtEnd(whileInitBlock)
	_, expr, _ := c.emitExprWithLoadIfNeeded(whileLoop.Cond)
	builder.CreateCondBr(expr, whileBodyBlock, endBlock)

	builder.SetInsertPointAtEnd(whileBodyBlock)
	c.emitBlock(function, whileLoop.Block)
	builder.CreateBr(whileInitBlock)

	builder.SetInsertPointAtEnd(endBlock)
}

func (c *codegen) emitRuntimeCall(name string, args []llvm.Value) {
	fn := c.module.NamedFunction(name)
	if fn.IsNil() {
		panic("runtime function is nil")
	}

	ty := fn.GlobalValueType()
	builder.CreateCall(ty, fn, args, "")
}

func (c *codegen) checkFloatTypeForBitSize(ty *ast.ExprType) (bool, int) {
	if ty.IsFloat() {
		b := ty.T.(*ast.BasicType)
		return true, b.Kind.BitSize()
	}
	return false, -1
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

	module := c.module.String()
	_, err = irFile.WriteString(module)
	// TODO(errors)
	if err != nil {
		return err
	}

	if err := llvm.VerifyModule(c.module, llvm.ReturnStatusAction); err != nil {
		fmt.Printf("error: module is not valid, generated at '%s'\n", dir)
		return err
	}

	optLevel := ""
	compilerFlags := ""
	switch buildType {
	case config.RELEASE:
		optLevel = "-O3"
		compilerFlags += "-Wl,-s" // no debug symbols
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
		if config.DEV {
			fmt.Printf("[DEV] OPT COMMAND: %s\n", cmd)
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
		if config.DEV {
			fmt.Printf("[DEV] CLANG COMMAND: %s\n", cmd)
		}
		return err
	}

	if config.DEV {
		fmt.Printf("[DEV] keeping '%s' build directory\n", dir)
	} else {
		err = os.RemoveAll(dir)
		if err != nil {
			return err
		}
	}

	return nil
}
