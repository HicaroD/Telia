package c

import (
	"fmt"
	"strings"

	"github.com/HicaroD/Telia/internal/ast"
	"github.com/HicaroD/Telia/internal/lexer/token"
)

// emitExpr returns the C expression string for the given node.
// For complex sub-expressions that must be hoisted to a temp variable,
// it writes the declaration to c.buf and returns the temp name.
func (c *CCodegen) emitExpr(node *ast.Node) string {
	switch node.Kind {
	case ast.KIND_ERROR_EXPR:
		return c.emitError(node.Node.(*ast.ErrorExpr))
	case ast.KIND_MATCH_EXPR:
		return c.emitMatch(node.Node.(*ast.MatchExpr))
	case ast.KIND_ARRAY_LITERAL_EXPR:
		lit := node.Node.(*ast.ArrayLiteralExpr)
		if lit.Type.T.(*ast.SequenceType).Foreign {
			var values []string
			for _, value := range lit.Values {
				values = append(values, c.emitExpr(value))
			}
			if len(values) == 0 {
				return "{0}"
			}
			return "{" + strings.Join(values, ", ") + "}"
		}
		if len(lit.Values) == 0 {
			return fmt.Sprintf("(%s){0}", emitCType(lit.Type))
		}
		if hasVaEffects(node) {
			return c.vaExpression(func() string {
				return fmt.Sprintf(
					"(%s){.data = %s}",
					emitCType(lit.Type),
					c.emitVaArrayInitializer(lit),
				)
			})
		}
		var values []string
		for _, value := range lit.Values {
			values = append(values, c.emitExpr(value))
		}
		return fmt.Sprintf("(%s){.data = {%s}}", emitCType(lit.Type), strings.Join(values, ", "))
	case ast.KIND_INDEX_EXPR:
		index := node.Node.(*ast.IndexExpr)
		// Keep collection capture, index evaluation and bounds checks at the
		// access site, including inside short-circuit operands and loop tests.
		// Return a pointer for addressable elements to preserve assignment targets.
		if node.IsAddressable() {
			return "(*" + c.vaExpression(
				func() string { return "&(" + c.emitIndex(index) + ")" },
			) + ")"
		}
		return c.vaExpression(func() string { return c.emitIndex(index) })
	case ast.KIND_LITERAL_EXPR:
		return c.emitLiteral(node.Node.(*ast.LiteralExpr))

	case ast.KIND_ID_EXPR:
		return c.emitIdExpr(node.Node.(*ast.IdExpr))

	case ast.KIND_BINARY_EXPR:
		return c.emitBinExpr(node.Node.(*ast.BinExpr))

	case ast.KIND_UNARY_EXPR:
		return c.emitUnaryExpr(node.Node.(*ast.UnaryExpr))

	case ast.KIND_FN_CALL:
		call := node.Node.(*ast.FnCall)
		if call.HandledValue {
			return c.emitPropagatedValue(call)
		}
		return c.emitFnCall(call)

	case ast.KIND_NAMESPACE_ACCESS:
		return c.emitNamespaceAccess(node.Node.(*ast.NamespaceAccess))

	case ast.KIND_NULLPTR_EXPR:
		if ty := node.Node.(*ast.NullPtrExpr).Type; ty != nil && ty.IsError() {
			return "(_Error){0}"
		}
		return "NULL"

	case ast.KIND_VOID_EXPR:
		return ""

	case ast.KIND_ADDRESS_OF_EXPR:
		inner := c.emitExpr(node.Node.(*ast.AddressOfExpr).Expr)
		return fmt.Sprintf("(&(%s))", inner)

	case ast.KIND_DEREF_POINTER_EXPR:
		if hasVaEffects(node.Node.(*ast.DerefPointerExpr).Expr) {
			return "(*" + c.vaExpression(func() string {
				inner := c.emitExpr(node.Node.(*ast.DerefPointerExpr).Expr)
				pointer := c.nextTmp()
				fmt.Fprintf(
					c.buf,
					"__auto_type %s = (%s);\n_check_nil_pointer_deref((void *)(%s));\n",
					pointer,
					inner,
					pointer,
				)
				return pointer
			}) + ")"
		}
		inner := c.emitExpr(node.Node.(*ast.DerefPointerExpr).Expr)
		c.buf.WriteString(fmt.Sprintf("_check_nil_pointer_deref((void *)(%s));\n", inner))
		return fmt.Sprintf("(*(%s))", inner)

	case ast.KIND_FIELD_ACCESS:
		return c.emitFieldAccess(node.Node.(*ast.FieldAccess))

	case ast.KIND_STRUCT_EXPR:
		if hasVaEffects(node) {
			return c.vaExpression(
				func() string { return c.emitStructLiteral(node.Node.(*ast.StructLiteralExpr)) },
			)
		}
		return c.emitStructLiteral(node.Node.(*ast.StructLiteralExpr))

	case ast.KIND_TUPLE_LITERAL_EXPR:
		if hasVaEffects(node) {
			return c.vaExpression(
				func() string { return c.emitTupleLiteral(node.Node.(*ast.TupleExpr)) },
			)
		}
		return c.emitTupleLiteral(node.Node.(*ast.TupleExpr))

	case ast.KIND_VARG_EXPR:
		// VarArgsExpr is only encountered inside emitFnCall args handling.
		panic("emitExpr: VARG_EXPR should be unwrapped by emitFnCall")

	default:
		panic(fmt.Sprintf("emitExpr: unhandled node kind: %v", node.Kind))
	}
}

// emitLiteral emits a C literal for any basic type.
func (c *CCodegen) emitLiteral(lit *ast.LiteralExpr) string {
	val := string(lit.Value)
	if lit.Type.Kind == ast.EXPR_TYPE_BASIC {
		b := lit.Type.T.(*ast.BasicType)
		switch b.Kind {
		case token.STRING_TYPE:
			return emitStringLiteral(lit.Value)
		case token.CSTRING_TYPE:
			return emitCBytes(lit.Value)
		case token.ERROR_TYPE:
			return fmt.Sprintf("(_Error){.msg = %s, .present = true}", emitStringLiteral(lit.Value))
		case token.BOOL_TYPE:
			if val == "true" {
				return "1"
			}
			return "0"
		case token.F32_TYPE, token.FLOAT_TYPE:
			// Ensure float suffix for f32.
			if !strings.ContainsAny(val, ".eE") {
				val += ".0"
			}
			return val + "f"
		case token.F64_TYPE:
			if !strings.ContainsAny(val, ".eE") {
				val += ".0"
			}
			return val
		default:
			// Integer types: emit as-is.
			return val
		}
	}
	return val
}

// Fixed-width octal escapes cannot consume a following hex digit, unlike \x.
func emitCBytes(data []byte) string {
	var b strings.Builder
	b.WriteByte('"')
	for _, v := range data {
		fmt.Fprintf(&b, "\\%03o", v)
	}
	b.WriteByte('"')
	return b.String()
}

func emitStringLiteral(data []byte) string {
	return fmt.Sprintf("(_String){(const unsigned char *)%s, %d}", emitCBytes(data), len(data))
}

// emitIdExpr resolves an identifier to its C name.
func (c *CCodegen) emitIdExpr(id *ast.IdExpr) string {
	if id.Constant != nil {
		return c.emitConstant(id.Constant)
	}
	if id.N == nil {
		// Unresolved — just use the name directly.
		return id.Name.Name()
	}
	switch id.N.Kind {
	case ast.KIND_FN_DECL:
		return c.mangledFnDeclName(id.N.Node.(*ast.FnDecl))
	case ast.KIND_PROTO:
		d := id.N.Node.(*ast.Proto)
		if d.Attributes.LinkName != "" {
			return d.Attributes.LinkName
		}
		return d.Name.Name()
	case ast.KIND_CONST_DECL:
		d := id.N.Node.(*ast.ConstDecl)
		if d.Foreign != nil {
			name := d.Attributes.LinkName
			if name == "" {
				name = d.Name.Name()
			}
			return fmt.Sprintf("((%s)(%s))", emitCType(d.Type), name)
		}
		return c.emitConstant(d.Value)
	case ast.KIND_VAR_ID_STMT:
		varId := id.N.Node.(*ast.VarIdStmt)
		if varId.BackendType != nil {
			return varId.BackendType.(*CVariable).Name
		}
		return varId.Name.Name()
	case ast.KIND_PARAM:
		param := id.N.Node.(*ast.Param)
		if param.BackendType != nil {
			return param.BackendType.(*CVariable).Name
		}
		return param.Name.Name()
	default:
		return id.Name.Name()
	}
}

// emitBinExpr emits a binary expression with parentheses for safety.
func (c *CCodegen) emitBinExpr(bin *ast.BinExpr) string {
	if bin.Op == token.SHIFT_LEFT || bin.Op == token.SHIFT_RIGHT {
		return c.emitShift(bin)
	}
	// Special case: error != nil — C can't compare structs by value.
	// Inspect presence instead of comparing C structs.
	if isErrorNilComparison(bin) {
		return c.emitErrorNilComparison(bin)
	}
	if bin.Op != token.AND && bin.Op != token.OR &&
		(hasVaEffects(bin.Left) || hasVaEffects(bin.Right)) {
		return c.emitSequencedVaBinary(bin)
	}
	left := c.emitExpr(bin.Left)
	right := c.emitExpr(bin.Right)
	return emitBinaryOperation(bin, left, right)
}

func emitBinaryOperation(bin *ast.BinExpr, left, right string) string {
	if bin.OperandType != nil && bin.OperandType.Equals(ast.NewBasicType(token.STRING_TYPE)) {
		comparison := fmt.Sprintf("_telia_string_equal(%s, %s)", left, right)
		if bin.Op == token.BANG_EQUAL {
			return "(!" + comparison + ")"
		}
		return comparison
	}
	op := binOpToC(bin.Op)
	return fmt.Sprintf("(%s %s %s)", left, op, right)
}

// isErrorNilComparison checks if a binary expression compares an error value with nil.
func isErrorNilComparison(bin *ast.BinExpr) bool {
	if bin.Op != token.EQUAL_EQUAL && bin.Op != token.BANG_EQUAL {
		return false
	}
	leftNil := bin.Left.Kind == ast.KIND_NULLPTR_EXPR
	rightNil := bin.Right.Kind == ast.KIND_NULLPTR_EXPR
	if !leftNil && !rightNil {
		return false
	}
	if bin.OperandType != nil && bin.OperandType.IsError() {
		return true
	}
	// Check if the non-nil side is error type (via field access on error, or variable)
	errorSide := bin.Right
	if rightNil {
		errorSide = bin.Left
	}
	return isErrorExpr(errorSide)
}

// isErrorExpr returns true if the expression produces an error type.
func isErrorExpr(node *ast.Node) bool {
	switch node.Kind {
	case ast.KIND_ID_EXPR:
		id := node.Node.(*ast.IdExpr)
		if id.N != nil {
			if id.N.Kind == ast.KIND_VAR_ID_STMT {
				varId := id.N.Node.(*ast.VarIdStmt)
				return varId.Type != nil && varId.Type.IsError()
			}
			if id.N.Kind == ast.KIND_PARAM {
				param := id.N.Node.(*ast.Param)
				return param.Type != nil && param.Type.IsError()
			}
		}
	case ast.KIND_LITERAL_EXPR:
		lit := node.Node.(*ast.LiteralExpr)
		return lit.Type != nil && lit.Type.IsError()
	case ast.KIND_ERROR_EXPR:
		return true
	}
	return false
}

// emitErrorNilComparison emits C code for error != nil / error == nil.
// Presence is independent of the borrowed message's storage.
func (c *CCodegen) emitErrorNilComparison(bin *ast.BinExpr) string {
	errorExpr := bin.Right
	if bin.Left.Kind != ast.KIND_NULLPTR_EXPR {
		errorExpr = bin.Left
	}
	val := c.emitExpr(errorExpr)
	op := binOpToC(bin.Op)
	return fmt.Sprintf("((%s).present %s false)", val, op)
}

// binOpToC maps a Telia binary operator token to its C operator string.
func binOpToC(op token.Kind) string {
	switch op {
	case token.PIPE:
		return "|"
	case token.CARET:
		return "^"
	case token.AMPERSAND:
		return "&"
	case token.PLUS:
		return "+"
	case token.MINUS:
		return "-"
	case token.STAR:
		return "*"
	case token.SLASH:
		return "/"
	case token.EQUAL_EQUAL:
		return "=="
	case token.BANG_EQUAL:
		return "!="
	case token.LESS:
		return "<"
	case token.LESS_EQ:
		return "<="
	case token.GREATER:
		return ">"
	case token.GREATER_EQ:
		return ">="
	case token.AND:
		return "&&"
	case token.OR:
		return "||"
	default:
		panic(fmt.Sprintf("binOpToC: unhandled operator: %v", op))
	}
}

// emitUnaryExpr emits a unary expression.
func (c *CCodegen) emitUnaryExpr(u *ast.UnaryExpr) string {
	val := c.emitExpr(u.Value)
	switch u.Op {
	case token.TILDE:
		return fmt.Sprintf("((%s)(~(%s)))", emitCType(u.Type), val)
	case token.MINUS:
		return fmt.Sprintf("(-(%s))", val)
	case token.NOT:
		return fmt.Sprintf("(!(%s))", val)
	default:
		panic(fmt.Sprintf("emitUnaryExpr: unhandled operator: %v", u.Op))
	}
}

// emitFnCall emits a function call expression, handling variadic args.
func (c *CCodegen) emitFnCall(call *ast.FnCall) string {
	if c.program.TestMode && call.Builtin == "" {
		return c.emitCheckedTestCall(call)
	}
	return c.emitFnCallExpr(call)
}

func (c *CCodegen) emitFnCallExpr(call *ast.FnCall) string {
	if c.program.TestMode && call.Proto != nil &&
		call.Proto.Attributes.LinkName == "_telia_testing_fail" {
		return c.vaExpression(func() string {
			c.emitTestFailureHeader(call.Name)
			c.emitTestFailureExit()
			return "((void)0)"
		})
	}
	if strings.HasPrefix(call.Builtin, "testing_") {
		return c.emitTestingAssertion(call)
	}
	switch call.Builtin {
	case "widen", "ptr_cast":
		return fmt.Sprintf("((%s)(%s))", emitCType(call.BuiltinType), c.emitExpr(call.Args[0]))
	case "narrow":
		return c.emitNarrow(call)
	case "va_end", "va_borrow", "va_arg":
		return c.emitVaOperation(call)
	case "sizeof":
		return fmt.Sprintf("((uint64_t)sizeof(%s))", emitCType(call.BuiltinType))
	case "alignof":
		return fmt.Sprintf("((uint64_t)_Alignof(%s))", emitCType(call.BuiltinType))
	case "zero":
		return fmt.Sprintf("(%s){0}", emitCType(call.BuiltinType))
	case "slice_from_raw":
		data, count := c.emitExpr(call.Args[0]), c.emitExpr(call.Args[1])
		return fmt.Sprintf(
			"(%s){(%s *)%s, %s}",
			emitCType(call.ResultType),
			emitCType(call.BuiltinType),
			data,
			count,
		)
	case "slice_data":
		return fmt.Sprintf("((void *)(%s).data)", c.emitExpr(call.Args[0]))
	}
	if call.Builtin == "view" || call.Builtin == "slice" {
		return c.emitSliceCall(call)
	}
	if call.Builtin == "len" {
		value := c.emitExpr(call.Args[0])
		if call.BuiltinType.Kind == ast.EXPR_TYPE_ARRAY {
			c.buf.WriteString(fmt.Sprintf("(void)(%s);\n", value))
			return fmt.Sprintf("UINT64_C(%d)", call.BuiltinType.T.(*ast.SequenceType).Length)
		}
		return fmt.Sprintf("(%s).len", value)
	}
	for _, argument := range call.Args {
		if hasVaEffects(argument) {
			return c.emitSequencedVaCall(call)
		}
	}
	name := c.getFnCallName(call)
	args := c.emitCallArgs(call.Args)
	return fmt.Sprintf("%s(%s)", name, args)
}

// emitCallArgs renders the argument list for a call, unwrapping VarArgsExpr.
func (c *CCodegen) emitCallArgs(args []*ast.Node) string {
	var parts []string
	for _, arg := range args {
		if arg.Kind == ast.KIND_VARG_EXPR {
			// Flatten variadic args.
			for _, varg := range arg.Node.(*ast.VarArgsExpr).Args {
				parts = append(parts, c.emitExpr(varg))
			}
		} else {
			parts = append(parts, c.emitExpr(arg))
		}
	}
	return strings.Join(parts, ", ")
}

// getFnCallName resolves the C function name for a call.
// Handles link_name, mangled names, and direct proto names.
func (c *CCodegen) getFnCallName(call *ast.FnCall) string {
	if call.CalleeExpr != nil {
		return "(" + c.emitExpr(call.CalleeExpr) + ")"
	}
	if call.Callee != nil {
		return c.emitIdExpr(call.Callee)
	}
	if call.Proto != nil {
		if call.Proto.Attributes.LinkName != "" {
			return call.Proto.Attributes.LinkName
		}
		return call.Proto.Name.Name()
	}
	if call.Decl != nil {
		return c.mangledFnDeclName(call.Decl)
	}
	// Fallback: bare name (builtin or unresolved).
	return call.Name.Name()
}

// mangledFnDeclName returns the C name for a FnDecl, using the package it
// belongs to. We walk the program to find which package owns the decl.
func (c *CCodegen) mangledFnDeclName(fn *ast.FnDecl) string {
	if fn.CName != "" {
		return fn.CName
	}
	// Check if this fn belongs to the root package.
	if c.fnBelongsToPackage(fn, c.program.Root) {
		return fn.Name.Name()
	}
	// Search imported packages.
	pkg := c.findPackageForFn(fn, c.program.Root)
	if pkg != nil {
		return pkg.Loc.Name + "__" + fn.Name.Name()
	}
	// Fallback.
	return fn.Name.Name()
}

// fnBelongsToPackage reports whether fn is declared in pkg.
func (c *CCodegen) fnBelongsToPackage(fn *ast.FnDecl, pkg *ast.Package) bool {
	if pkg == nil {
		return false
	}
	for _, file := range pkg.Files {
		for _, node := range file.Body {
			if node.Kind == ast.KIND_FN_DECL && node.Node.(*ast.FnDecl) == fn {
				return true
			}
		}
	}
	return false
}

// findPackageForFn recursively searches pkg and its imports for the package
// that owns fn. Returns nil if not found.
func (c *CCodegen) findPackageForFn(fn *ast.FnDecl, pkg *ast.Package) *ast.Package {
	if pkg == nil {
		return nil
	}
	for _, file := range pkg.Files {
		for _, node := range file.Body {
			if node.Kind == ast.KIND_FN_DECL && node.Node.(*ast.FnDecl) == fn {
				return pkg
			}
		}
		for _, imp := range file.Imports {
			if found := c.findPackageForFn(fn, imp.Package); found != nil {
				return found
			}
		}
	}
	return nil
}

// emitNamespaceAccess emits a namespace-qualified call (e.g. io::println or libc::printf).
func (c *CCodegen) emitNamespaceAccess(ns *ast.NamespaceAccess) string {
	if ns.EnumVariant != nil {
		if hasVaEffects(ns.Right) {
			return c.vaExpression(func() string { return c.emitEnumConstruction(ns) })
		}
		return c.emitEnumConstruction(ns)
	}
	if ns.Right.Kind == ast.KIND_FN_CALL {
		return c.emitExpr(ns.Right)
	}
	// Nested namespace access (e.g. pkg::sub::fn) — recurse.
	if ns.Right.Kind == ast.KIND_NAMESPACE_ACCESS {
		return c.emitNamespaceAccess(ns.Right.Node.(*ast.NamespaceAccess))
	}
	return c.emitExpr(ns.Right)
}

// emitFieldAccess emits a struct field access expression.
func (c *CCodegen) emitFieldAccess(fa *ast.FieldAccess) string {
	if fa.Receiver != nil {
		recv := c.emitExpr(fa.Receiver)
		separator := "."
		if fa.ReceiverType.IsPointer() {
			separator = "->"
		}
		return fieldAccessTail("("+recv+")"+separator+fa.AccessedField.Name.Name(), fa)
	}
	var recv string
	if fa.StructVar != nil {
		if fa.StructVar.BackendType != nil {
			recv = fa.StructVar.BackendType.(*CVariable).Name
		} else {
			recv = fa.StructVar.Name.Name()
		}
	} else if fa.StructParam != nil {
		if fa.StructParam.BackendType != nil {
			recv = fa.StructParam.BackendType.(*CVariable).Name
		} else {
			recv = fa.StructParam.Name.Name()
		}
	} else {
		recv = fa.Left.Name.Name()
	}

	fieldName := fa.AccessedField.Name.Name()

	// Pointer receiver uses ->, value receiver uses .
	if fa.StructVar != nil {
		varId := fa.StructVar
		isPointerVar := varId.Pointer ||
			varId.NumberOfPointerReceivers > 0 ||
			(varId.Type != nil && varId.Type.Kind == ast.EXPR_TYPE_POINTER)
		if isPointerVar {
			return fieldAccessTail(fmt.Sprintf("%s->%s", recv, fieldName), fa)
		}
	} else if fa.StructParam != nil {
		if fa.StructParam.Type != nil && fa.StructParam.Type.Kind == ast.EXPR_TYPE_POINTER {
			return fieldAccessTail(fmt.Sprintf("%s->%s", recv, fieldName), fa)
		}
	}
	return fieldAccessTail(fmt.Sprintf("%s.%s", recv, fieldName), fa)
}

func fieldAccessTail(expr string, fa *ast.FieldAccess) string {
	if fa.Right.Kind != ast.KIND_FIELD_ACCESS {
		return expr
	}
	next := fa.Right.Node.(*ast.FieldAccess)
	separator := "."
	if fa.AccessedField.Type.IsPointer() {
		separator = "->"
	}
	return fieldAccessTail(expr+separator+next.AccessedField.Name.Name(), next)
}

// emitStructLiteral emits a C compound literal for a struct.
func (c *CCodegen) emitStructLiteral(sl *ast.StructLiteralExpr) string {
	name := sl.Name.Name()
	if sl.Decl != nil && sl.Decl.CName != "" {
		name = sl.Decl.CName
	}
	if len(sl.Values) == 0 {
		return fmt.Sprintf("(%s){0}", name)
	}
	var fields []string
	sequence := false
	for _, value := range sl.Values {
		sequence = sequence || hasVaEffects(value.Value)
	}
	for _, fv := range sl.Values {
		var val string
		if lit, ok := fv.Value.Node.(*ast.ArrayLiteralExpr); sequence && ok &&
			lit.Type.T.(*ast.SequenceType).Foreign {
			val = c.emitVaArrayInitializer(lit)
		} else {
			val = c.emitExpr(fv.Value)
			if sequence {
				tmp := c.nextTmp()
				fmt.Fprintf(c.buf, "%s %s = (%s);\n", vaTemporaryType(fv.Value), tmp, val)
				val = tmp
			}
		}
		fields = append(fields, fmt.Sprintf(".%s = %s", fv.Name.Name(), val))
	}
	return fmt.Sprintf("(%s){%s}", name, strings.Join(fields, ", "))
}

// emitTupleLiteral emits a C compound literal for a tuple expression, e.g.:
//
//	(_Tuple_int32_t_int32_t){._0 = 1, ._1 = (2 + value)}
func (c *CCodegen) emitTupleLiteral(te *ast.TupleExpr) string {
	tupleType := &ast.ExprType{
		Kind: ast.EXPR_TYPE_TUPLE,
		T:    te.Type,
	}
	name := tupleTypedefName(tupleType)
	parts := make([]string, len(te.Exprs))
	sequence := false
	for _, expr := range te.Exprs {
		sequence = sequence || hasVaEffects(expr)
	}
	for i, expr := range te.Exprs {
		value := c.emitExpr(expr)
		if sequence {
			tmp := c.nextTmp()
			fmt.Fprintf(c.buf, "%s %s = (%s);\n", vaTemporaryType(expr), tmp, value)
			value = tmp
		}
		parts[i] = fmt.Sprintf("._%d = %s", i, value)
	}
	return fmt.Sprintf("(%s){%s}", name, strings.Join(parts, ", "))
}
