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
	case ast.KIND_LITERAL_EXPR:
		return c.emitLiteral(node.Node.(*ast.LiteralExpr))

	case ast.KIND_ID_EXPR:
		return c.emitIdExpr(node.Node.(*ast.IdExpr))

	case ast.KIND_BINARY_EXPR:
		return c.emitBinExpr(node.Node.(*ast.BinExpr))

	case ast.KIND_UNARY_EXPR:
		return c.emitUnaryExpr(node.Node.(*ast.UnaryExpr))

	case ast.KIND_FN_CALL:
		return c.emitFnCall(node.Node.(*ast.FnCall))

	case ast.KIND_NAMESPACE_ACCESS:
		return c.emitNamespaceAccess(node.Node.(*ast.NamespaceAccess))

	case ast.KIND_NULLPTR_EXPR:
		return "NULL"

	case ast.KIND_VOID_EXPR:
		return ""

	case ast.KIND_ADDRESS_OF_EXPR:
		inner := c.emitExpr(node.Node.(*ast.AddressOfExpr).Expr)
		return fmt.Sprintf("(&(%s))", inner)

	case ast.KIND_DEREF_POINTER_EXPR:
		inner := c.emitExpr(node.Node.(*ast.DerefPointerExpr).Expr)
		c.buf.WriteString(fmt.Sprintf("_check_nil_pointer_deref((void *)(%s));\n", inner))
		return fmt.Sprintf("(*(%s))", inner)

	case ast.KIND_FIELD_ACCESS:
		return c.emitFieldAccess(node.Node.(*ast.FieldAccess))

	case ast.KIND_STRUCT_EXPR:
		return c.emitStructLiteral(node.Node.(*ast.StructLiteralExpr))

	case ast.KIND_TUPLE_LITERAL_EXPR:
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
		case token.STRING_TYPE, token.CSTRING_TYPE:
			// The lexer stores the literal without surrounding quotes.
			return fmt.Sprintf("%q", val)
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

// emitIdExpr resolves an identifier to its C name.
func (c *CCodegen) emitIdExpr(id *ast.IdExpr) string {
	if id.N == nil {
		// Unresolved — just use the name directly.
		return id.Name.Name()
	}
	switch id.N.Kind {
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
	left := c.emitExpr(bin.Left)
	right := c.emitExpr(bin.Right)
	op := binOpToC(bin.Op)
	return fmt.Sprintf("(%s %s %s)", left, op, right)
}

// binOpToC maps a Telia binary operator token to its C operator string.
func binOpToC(op token.Kind) string {
	switch op {
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
	if ns.Right.Kind == ast.KIND_FN_CALL {
		return c.emitFnCall(ns.Right.Node.(*ast.FnCall))
	}
	// Nested namespace access (e.g. pkg::sub::fn) — recurse.
	if ns.Right.Kind == ast.KIND_NAMESPACE_ACCESS {
		return c.emitNamespaceAccess(ns.Right.Node.(*ast.NamespaceAccess))
	}
	return c.emitExpr(ns.Right)
}

// emitFieldAccess emits a struct field access expression.
func (c *CCodegen) emitFieldAccess(fa *ast.FieldAccess) string {
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
			return fmt.Sprintf("%s->%s", recv, fieldName)
		}
	} else if fa.StructParam != nil {
		if fa.StructParam.Type != nil && fa.StructParam.Type.Kind == ast.EXPR_TYPE_POINTER {
			return fmt.Sprintf("%s->%s", recv, fieldName)
		}
	}
	return fmt.Sprintf("%s.%s", recv, fieldName)
}

// emitStructLiteral emits a C compound literal for a struct.
func (c *CCodegen) emitStructLiteral(sl *ast.StructLiteralExpr) string {
	var fields []string
	for _, fv := range sl.Values {
		val := c.emitExpr(fv.Value)
		fields = append(fields, fmt.Sprintf(".%s = %s", fv.Name.Name(), val))
	}
	return fmt.Sprintf("(%s){%s}", sl.Name.Name(), strings.Join(fields, ", "))
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
	for i, expr := range te.Exprs {
		parts[i] = fmt.Sprintf("._%d = %s", i, c.emitExpr(expr))
	}
	return fmt.Sprintf("(%s){%s}", name, strings.Join(parts, ", "))
}
