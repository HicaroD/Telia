package sema

import (
	"github.com/HicaroD/Telia/internal/ast"
	"github.com/HicaroD/Telia/internal/ctarget"
	"github.com/HicaroD/Telia/internal/diagnostics"
	"github.com/HicaroD/Telia/internal/lexer/token"
	"github.com/HicaroD/Telia/internal/source"
)

func terminalVaCall(node *ast.Node) *ast.FnCall {
	if node == nil {
		return nil
	}
	for node.Kind == ast.KIND_NAMESPACE_ACCESS {
		node = node.Node.(*ast.NamespaceAccess).Right
	}
	call, _ := node.Node.(*ast.FnCall)
	return call
}

func ownerCopyDiagnostic(span source.Span) error {
	return diagnostics.Diag{Title: "Noncopyable va_list owner", Message: "owning va_list storage cannot be copied or reassigned", Span: span, Help: "Create an independent owner with va_copy(va_borrow(&owner)), and end each owner separately."}
}

func vaIntrinsicValueDiagnostic(span source.Span) error {
	return diagnostics.Diag{Title: "Intrinsic is not a function value", Message: "va_list operations require direct compiler-recognized invocation", Span: span, Help: "Call va_start, va_copy, va_arg, va_borrow, or va_end directly; these operations do not have ordinary callback addresses."}
}

func (s *sema) prepareVaIntrinsics(pkg *ast.Package) {
	// The owner's compiler type exports these operations into its declaring
	// package. They have no fake C prototypes or addressable function bodies.
	for _, node := range pkg.Scope.Nodes {
		d, ok := node.Node.(*ast.TypeAlias)
		if !ok {
			continue
		}
		id, ok := d.Type.T.(*ast.IdType)
		if !ok || id.Name.Name() != "__c_va_list_owner" {
			continue
		}
		for _, name := range []string{"va_copy", "va_end", "va_borrow", "va_arg", "va_start"} {
			tok := *d.Name
			tok.Lexeme = []byte(name)
			if pkg.Scope.Nodes[name] == nil {
				pkg.Scope.Nodes[name] = &ast.Node{Kind: ast.KIND_VA_INTRINSIC, Node: &ast.VaIntrinsic{Name: &tok}, Span: tok.Span}
			}
		}
		return
	}
}

func (s *sema) checkVaIntrinsic(call *ast.FnCall, scope, decl *ast.Scope, imported bool) (*ast.ExprType, error) {
	name := call.Name.Name()
	if (name == "va_start" || name == "va_copy") && s.vaInitializer != call {
		return nil, diagnostics.Diag{Title: "va_list owner initialization", Message: "a va_list owner must be initialized directly in a new local variable", Span: call.Name.Span, Help: "Write owner := libc::" + name + "(...), then arrange va_end(&owner) on every exit path."}
	}
	typeCount := 0
	if name == "va_arg" {
		typeCount = 1
	}
	if len(call.TypeArgs) != typeCount || len(call.Args) != 1 || call.AtOp != nil {
		return nil, diagnostics.Diag{Title: "Invalid va_list operation", Message: name + " requires one value argument and the operation's exact type arguments", Span: call.Name.Span, Help: "Use va_copy(borrowed), va_end(&owner), va_borrow(&owner), or va_arg<PromotedType>(&owner), without an error operator."}
	}
	result := ast.NewBasicType(token.VOID_TYPE)
	if name == "va_start" {
		fn := s.currentFunction
		id, named := call.Args[0].Node.(*ast.IdExpr)
		if fn == nil || !fn.Params.IsVariadic || fn.Params.Len == 0 || !named {
			return nil, diagnostics.Diag{Title: "Invalid va_start context", Message: "va_start requires this function's C variadic marker and at least one fixed parameter", Span: call.Name.Span, Help: "Use va_start(args) directly inside fn name(fixed Type, args @c ...)."}
		}
		if fn.Attributes.CallingConventionMacro != "" {
			return nil, diagnostics.Diag{Title: "Unsupported va_start convention", Message: "native va_start storage has not been validated for this calling-convention macro", Span: call.Name.Span, Help: "Use a native cc=\"c\" variadic wrapper or a header-compiled C adapter for convention-specific va_list types.", Related: []diagnostics.Related{{Span: fn.Name.Span, Message: "variadic function declared here"}}}
		}
		marker := fn.Params.Fields[fn.Params.Len]
		symbol, err := scope.LookupAcrossScopes(id.Name.Name())
		if err != nil || symbol.Node != marker || marker.Attributes == nil || !marker.Attributes.C {
			return nil, diagnostics.Diag{Title: "Invalid va_start marker", Message: "the argument is not this function's C variadic marker", Span: id.Name.Span, Help: "Pass the args @c ... parameter belonging to this function.", Related: []diagnostics.Related{{Span: marker.Name.Span, Message: "variadic parameter declared here"}}}
		}
		id.N = symbol
		call.VaLast = fn.Params.Fields[fn.Params.Len-1]
		valid, err := vaPromotedType(call.VaLast.Type, call.Name.Span)
		if err != nil {
			return nil, err
		}
		if !valid {
			return nil, diagnostics.Diag{Title: "Invalid va_start parameter", Message: "the last fixed parameter must have a supported type unchanged by C default promotions", Span: call.Name.Span, Help: "Use a final libc::c_int, libc::c_double, or pointer parameter before args @c ....", Related: []diagnostics.Related{{Span: call.VaLast.Name.Span, Message: "last fixed parameter declared here"}}}
		}
		result = &ast.ExprType{Kind: ast.EXPR_TYPE_VA_LIST_OWNER, Explicit: true, T: &ast.VaListOwnerType{}}
	} else if name == "va_copy" {
		ty, _, err := s.inferExprTypeWithoutContext(call.Args[0], scope, decl, imported, true)
		if err != nil {
			return nil, err
		}
		if ty.Kind != ast.EXPR_TYPE_C_VA_LIST {
			return nil, diagnostics.Diag{Title: "Invalid va_copy source", Message: "va_copy requires a borrowed native va_list", Span: call.Args[0].Span, Help: "Copy a c_va_list parameter, or use va_borrow(&owner)."}
		}
		result = &ast.ExprType{Kind: ast.EXPR_TYPE_VA_LIST_OWNER, Explicit: true, T: &ast.VaListOwnerType{}}
	} else {
		address, ok := call.Args[0].Node.(*ast.AddressOfExpr)
		if !ok {
			return nil, diagnostics.Diag{Title: "Invalid va_list owner", Message: "this operation requires the address of a local va_list owner", Span: call.Args[0].Span, Help: "Pass &owner, where owner was created by va_copy or va_start."}
		}
		id, ok := address.Expr.Node.(*ast.IdExpr)
		if !ok {
			return nil, diagnostics.Diag{Title: "Invalid va_list owner", Message: "va_list owners must be named local storage", Span: address.Expr.Span, Help: "Use the address of the local owner variable."}
		}
		ty, _, err := s.inferIdExprTypeWithoutContext(id, scope)
		if err != nil {
			return nil, err
		}
		owner, local := id.N.Node.(*ast.VarIdStmt)
		if ty.Kind != ast.EXPR_TYPE_VA_LIST_OWNER || !local {
			return nil, diagnostics.Diag{Title: "Invalid va_list owner", Message: "the argument is not local owning va_list storage", Span: id.Name.Span, Help: "Create a local owner with va_copy or va_start first."}
		}
		call.VaOwner = owner
		if name == "va_borrow" {
			result = &ast.ExprType{Kind: ast.EXPR_TYPE_C_VA_LIST, Explicit: true, T: &ast.CVaListType{}}
		}
		if name == "va_arg" {
			ty, err := s.resolveType(call.TypeArgs[0], scope)
			if err != nil {
				return nil, err
			}
			valid, err := vaPromotedType(ty, call.Name.Span)
			if err != nil {
				return nil, err
			}
			if !valid {
				return nil, diagnostics.Diag{Title: "Invalid va_arg type", Message: "va_arg must name a supported actual type after C default promotions", Span: call.Name.Span, Help: "Read narrow integers and bool as libc::c_int, and f32 as libc::c_double; otherwise use the actual promoted scalar or pointer type passed by the caller.", Related: []diagnostics.Related{{Span: owner.Name.Span, Message: "owning cursor declared here"}}}
			}
			call.BuiltinType, result = ty, ty
		}
	}
	call.Builtin, call.ResultType = name, result
	return result, nil
}

func vaPromotedType(ty *ast.ExprType, span source.Span) (bool, error) {
	if ty.Kind == ast.EXPR_TYPE_POINTER || ty.Kind == ast.EXPR_TYPE_FUNCTION {
		return true, nil
	}
	if ty.Kind != ast.EXPR_TYPE_BASIC {
		return false, nil
	}
	basic := ty.T.(*ast.BasicType)
	if basic.Kind == token.F64_TYPE || basic.Kind == token.RAWPTR_TYPE || basic.Kind == token.CSTRING_TYPE {
		return true, nil
	}
	if !basic.Kind.IsInteger() {
		return false, nil
	}
	scalars, err := ctarget.Scalars()
	if err != nil {
		return false, diagnostics.Diag{Title: "C target unavailable", Message: err.Error(), Span: span, Help: "Install the native C toolchain used by this build."}
	}
	switch basic.CName {
	case "char", "signed char", "unsigned char", "short", "unsigned short":
		return false, nil
	}
	return basic.Kind.BitSize() >= scalars["__c_int"].Bytes*8, nil
}

func rejectVaListEscape(ty *ast.ExprType, span source.Span) error {
	switch ty.Kind {
	case ast.EXPR_TYPE_VA_LIST_OWNER:
		return diagnostics.Diag{Title: "va_list owner escape", Message: "owning va_list storage must remain local to the function that initialized it", Span: span, Help: "Use va_arg, va_end, va_borrow, or va_copy explicitly; do not pass, return, retain, or take ordinary addresses of the owner."}
	case ast.EXPR_TYPE_C_VA_LIST:
		return diagnostics.Diag{Title: "Borrowed va_list escape", Message: "a borrowed va_list may only be forwarded during the current call", Span: span, Help: "Create local owning storage with libc::va_copy if you need an independent cursor; do not return or retain the borrowed parameter."}
	case ast.EXPR_TYPE_POINTER:
		return rejectVaListEscape(ty.T.(*ast.PointerType).Type, span)
	case ast.EXPR_TYPE_ARRAY, ast.EXPR_TYPE_SLICE:
		return rejectVaListEscape(ty.T.(*ast.SequenceType).Element, span)
	case ast.EXPR_TYPE_TUPLE:
		for _, item := range ty.T.(*ast.TupleType).Types {
			if err := rejectVaListEscape(item, span); err != nil {
				return err
			}
		}
	}
	return nil
}
