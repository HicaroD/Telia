package sema

import (
	"errors"
	"fmt"
	"strings"

	"github.com/HicaroD/Telia/internal/ast"
	"github.com/HicaroD/Telia/internal/diagnostics"
	"github.com/HicaroD/Telia/internal/lexer/token"
)

// Validate names and operations which do not depend on substitution without
// checking a template as a concrete function or modifying its syntax tree.
type templateValidator struct {
	matchDepth int
	s          *sema
	params     map[string]bool
	locals     map[string]*ast.ExprType
}

func (s *sema) validateTemplates(pkg *ast.Package, seen map[*ast.Package]bool) error {
	if seen[pkg] {
		return nil
	}
	seen[pkg] = true
	for _, file := range pkg.Files {
		for _, imp := range file.Imports {
			if err := s.validateTemplates(imp.Package, seen); err != nil {
				if !errors.Is(err, diagnostics.COMPILER_ERROR_FOUND) {
					return err
				}
			}
		}
		if err := s.inContext(pkg, file, nil, func() error {
			for _, node := range file.Body {
				var params []*token.Token
				switch d := node.Node.(type) {
				case *ast.EnumDecl:
					params = d.TypeParams
				case *ast.StructDecl:
					params = d.TypeParams
				case *ast.FnDecl:
					params = d.TypeParams
				}
				if len(params) == 0 {
					continue
				}
				err := func() error {
					v := templateValidator{s: s, params: map[string]bool{}, locals: map[string]*ast.ExprType{}}
					for _, p := range params {
						v.params[p.Name()] = true
					}
					switch d := node.Node.(type) {
					case *ast.EnumDecl:
						if err := v.typeExpr(d.Type, map[*ast.TypeAlias]bool{}); err != nil {
							return err
						}
						for _, variant := range d.Variants {
							for _, field := range variant.Fields {
								if err := v.typeExpr(field.Type, map[*ast.TypeAlias]bool{}); err != nil {
									return err
								}
							}
						}
					case *ast.StructDecl:
						if err := s.checkStructDecl(d); err != nil {
							return err
						}
						if d.Attributes.LinkName != "" || d.Attributes.Linkage != "" {
							return fmt.Errorf("%s: generic templates cannot have C linkage", d.Name.Pos)
						}
						for _, f := range d.Fields {
							if err := v.typeExpr(f.Type, map[*ast.TypeAlias]bool{}); err != nil {
								return err
							}
						}
					case *ast.FnDecl:
						if d.Name.Name() == "main" || token.IsAnyBuiltin(d.Name.Name()) {
							return fmt.Errorf("%s: '%s' cannot be generic", d.Name.Pos, d.Name.Name())
						}
						if d.Attributes.LinkName != "" || d.Attributes.Linkage != "" {
							return fmt.Errorf("%s: generic templates cannot have C linkage", d.Name.Pos)
						}
						if err := s.checkFnAttributes(d.Attributes); err != nil {
							return err
						}
						if err := v.typeExpr(d.RetType, map[*ast.TypeAlias]bool{}); err != nil {
							return err
						}
						for _, p := range d.Params.Fields {
							if err := v.typeExpr(p.Type, map[*ast.TypeAlias]bool{}); err != nil {
								return err
							}
							v.locals[p.Name.Name()] = p.Type
						}
						if err := v.block(d.Block); err != nil {
							return err
						}
						if !d.RetType.IsVoid() && !s.checkNonVoidFunctionReturns(d.Block) {
							return fmt.Errorf("%s: non-void function '%s' must always return a value", d.Name.Pos, d.Name.Name())
						}
					}
					return nil
				}()
				if err != nil {
					s.collector.ReportError(err, node.DeclarationSpan())
					node.Invalid = true
				}
			}
			return nil
		}); err != nil {
			return err
		}
	}
	if s.collector.HasErrors() {
		return diagnostics.COMPILER_ERROR_FOUND
	}
	return nil
}

func (v *templateValidator) typeExpr(ty *ast.ExprType, aliases map[*ast.TypeAlias]bool) error {
	if ty == nil {
		return nil
	}
	switch t := ty.T.(type) {
	case *ast.FunctionType:
		for _, p := range t.Params.Fields {
			if err := v.typeExpr(p.Type, aliases); err != nil {
				return err
			}
		}
		return v.typeExpr(t.Result, aliases)
	case *ast.SequenceType:
		return v.typeExpr(t.Element, aliases)
	case *ast.IdType:
		if t.Namespace == nil && v.params[t.Name.Name()] {
			if len(t.Args) != 0 {
				return fmt.Errorf("%s: type parameter '%s' cannot take type arguments", t.Name.Pos, t.Name.Name())
			}
			return nil
		}
		scope := v.s.pkg.Scope
		if t.Namespace != nil {
			imp, ok := v.s.file.Imports[t.Namespace.Name()]
			if !ok {
				return fmt.Errorf("%s: unknown type namespace %s", t.Name.Pos, t.Namespace.Name())
			}
			scope = imp.Package.Scope
		}
		sym, err := scope.LookupAcrossScopes(t.Name.Name())
		if err != nil {
			if t.Namespace == nil && strings.HasPrefix(t.Name.Name(), "__c_") {
				_, err := v.s.resolveType(ty, scope)
				return err
			}
			return fmt.Errorf("%s: '%s' type not found on scope", t.Name.Pos, t.Name.Name())
		}
		want := 0
		switch d := sym.Node.(type) {
		case *ast.StructDecl:
			want = len(d.TypeParams)
		case *ast.EnumDecl:
			want = len(d.TypeParams)
		case *ast.OpaqueDecl:
			want = 0
		case *ast.TypeAlias:
			if aliases[d] {
				return fmt.Errorf("%s: cyclic type alias '%s'", t.Name.Pos, t.Name.Name())
			}
			aliases[d] = true
			defer delete(aliases, d)
			if err := v.s.inContext(d.Owner, d.File, nil, func() error {
				aliasValidator := templateValidator{s: v.s, params: map[string]bool{}}
				return aliasValidator.typeExpr(d.Type, aliases)
			}); err != nil {
				return err
			}
		default:
			return fmt.Errorf("%s: '%s' is not a type", t.Name.Pos, t.Name.Name())
		}
		if len(t.Args) != want {
			return fmt.Errorf("%s: '%s' requires %d type arguments, got %d", t.Name.Pos, t.Name.Name(), want, len(t.Args))
		}
		for _, arg := range t.Args {
			if err := v.typeExpr(arg, aliases); err != nil {
				return err
			}
		}
	case *ast.PointerType:
		return v.typeExpr(t.Type, aliases)
	case *ast.TupleType:
		for _, t := range t.Types {
			if err := v.typeExpr(t, aliases); err != nil {
				return err
			}
		}
	}
	return nil
}

func (v *templateValidator) block(block *ast.BlockStmt) error {
	if block == nil {
		return nil
	}
	outer := v.locals
	v.locals = map[string]*ast.ExprType{}
	for name, ty := range outer {
		v.locals[name] = ty
	}
	defer func() { v.locals = outer }()
	for _, n := range block.Statements {
		if _, err := v.node(n); err != nil {
			return err
		}
	}
	return nil
}

func (v *templateValidator) node(n *ast.Node) (*ast.ExprType, error) {
	if n == nil {
		return nil, nil
	}
	switch t := n.Node.(type) {
	case *ast.ErrorExpr:
		if len(t.Args) != 1 && len(t.Args) != 2 {
			return nil, fmt.Errorf("error constructor requires a message or an error-enum value and a message")
		}
		for _, arg := range t.Args {
			if _, err := v.node(arg); err != nil {
				return nil, err
			}
		}
		return ast.NewBasicType(token.ERROR_TYPE), nil
	case *ast.MatchExpr:
		v.matchDepth++
		defer func() { v.matchDepth-- }()
		if _, err := v.node(t.Value); err != nil {
			return nil, err
		}
		for _, arm := range t.Arms {
			for _, pattern := range arm.Patterns() {
				if err := v.typeExpr(pattern.PatternType, map[*ast.TypeAlias]bool{}); err != nil {
					return nil, err
				}
				if _, err := v.node(pattern.Literal); err != nil {
					return nil, err
				}
			}
			outer := v.locals
			v.locals = map[string]*ast.ExprType{}
			for name, ty := range outer {
				v.locals[name] = ty
			}
			for _, f := range arm.Fields {
				v.locals[f.Binding.Name.Name()] = nil
			}
			if arm.Binding != nil {
				v.locals[arm.Binding.Name.Name()] = nil
			}
			if _, err := v.node(arm.Guard); err != nil {
				return nil, err
			}
			var err error
			if arm.Block != nil {
				err = v.block(arm.Block)
			} else {
				_, err = v.node(arm.Result)
			}
			v.locals = outer
			if err != nil {
				return nil, err
			}
		}
		return nil, nil
	case *ast.YieldStmt:
		if v.matchDepth == 0 {
			return nil, diagnostics.Diag{Title: "Yield outside match", Message: "yield is only valid inside a match arm block", Span: n.Span, Help: "Use return to return from a function, or yield inside a match arm."}
		}
		return v.node(t.Value)
	case *ast.ConstDecl:
		if err := v.typeExpr(t.Type, map[*ast.TypeAlias]bool{}); err != nil {
			return nil, err
		}
		ty, err := v.node(t.Expr)
		if err != nil {
			return nil, err
		}
		if t.Type != nil {
			ty = t.Type
		}
		v.locals[t.Name.Name()] = ty
		return nil, nil
	case *ast.ArrayLiteralExpr:
		for _, value := range t.Values {
			if _, err := v.node(value); err != nil {
				return nil, err
			}
		}
		return nil, nil
	case *ast.IndexExpr:
		if _, err := v.node(t.Value); err != nil {
			return nil, err
		}
		_, err := v.node(t.Index)
		return nil, err
	case *ast.LiteralExpr:
		return t.Type, nil
	case *ast.NullPtrExpr:
		return t.Type, nil
	case *ast.IdExpr:
		if ty, ok := v.locals[t.Name.Name()]; ok {
			return ty, nil
		}
		if symbol, err := v.s.pkg.Scope.LookupAcrossScopes(t.Name.Name()); err == nil && symbol.Kind == ast.KIND_CONST_DECL {
			d := symbol.Node.(*ast.ConstDecl)
			if err := v.s.checkConstDecl(d); err != nil {
				return nil, err
			}
			return d.Type, nil
		}
		return nil, fmt.Errorf("%s: '%s' not defined", t.Name.Pos, t.Name.Name())
	case *ast.BinExpr:
		left, err := v.node(t.Left)
		if err != nil {
			return nil, err
		}
		right, err := v.node(t.Right)
		if err != nil {
			return nil, err
		}
		if left == nil || right == nil || left.Kind != ast.EXPR_TYPE_BASIC || right.Kind != ast.EXPR_TYPE_BASIC {
			return nil, nil
		}
		if t.Op == token.SHIFT_LEFT || t.Op == token.SHIFT_RIGHT {
			if !left.IsInteger() || !right.IsInteger() {
				return nil, diagnostics.Diag{Title: "Invalid shift operands", Message: "shifts require integer operands", Span: n.Span, Help: "Use an integer left operand and an integer count."}
			}
			return left, nil
		}
		if (t.Op == token.EQUAL_EQUAL || t.Op == token.BANG_EQUAL) &&
			(left.IsError() && right.T.(*ast.BasicType).Kind == token.UNTYPED_NULLPTR || right.IsError() && left.T.(*ast.BasicType).Kind == token.UNTYPED_NULLPTR) {
			return ast.NewBasicType(token.BOOL_TYPE), nil
		}
		op := ast.BinaryOperators[t.Op]
		// Literal operands adopt the other numeric operand's type, just as
		// concrete expression checking does with an expected type.
		if left.IsInteger() && right.IsInteger() || left.IsFloat() && right.IsFloat() {
			if _, literal := t.Right.Node.(*ast.LiteralExpr); literal {
				right = left
			}
			if _, literal := t.Left.Node.(*ast.LiteralExpr); literal {
				left = right
			}
		}
		valid := false
		for _, ty := range op.ValidTypes {
			if operatorTypeMatches(ty, left) && operatorTypeMatches(ty, right) {
				valid = true
				break
			}
		}
		if !valid {
			return nil, fmt.Errorf("invalid operand types for %s: %s and %s", t.Op, left.T, right.T)
		}
		if op.Handler != nil {
			return op.Handler([]*ast.ExprType{left, right})
		}
		return op.ResultType, nil
	case *ast.UnaryExpr:
		operand, err := v.node(t.Value)
		if err != nil || operand == nil || operand.Kind != ast.EXPR_TYPE_BASIC {
			return nil, err
		}
		op := ast.UnaryOperators[t.Op]
		for _, ty := range op.ValidTypes {
			if operatorTypeMatches(ty, operand) {
				if op.Handler != nil {
					return op.Handler([]*ast.ExprType{operand})
				}
				return op.ResultType, nil
			}
		}
		return nil, fmt.Errorf("invalid operand type for %s: %s", t.Op, operand.T)
	case *ast.AddressOfExpr:
		ty, err := v.node(t.Expr)
		if err != nil || ty == nil {
			return nil, err
		}
		return ast.NewPointerType(ty), nil
	case *ast.DerefPointerExpr:
		ty, err := v.node(t.Expr)
		if err != nil || ty == nil {
			return nil, err
		}
		if ty.Kind == ast.EXPR_TYPE_POINTER {
			return ty.T.(*ast.PointerType).Type, nil
		}
		return nil, nil
	case *ast.FieldAccess:
		if t.Receiver != nil {
			_, err := v.node(t.Receiver)
			return nil, err
		}
		_, err := v.node(&ast.Node{Kind: ast.KIND_ID_EXPR, Node: t.Left})
		return nil, err
	case *ast.StructLiteralExpr:
		ty := &ast.ExprType{Kind: ast.EXPR_TYPE_ID, T: &ast.IdType{Name: t.Name, Namespace: t.Namespace, Args: t.TypeArgs}}
		if err := v.typeExpr(ty, map[*ast.TypeAlias]bool{}); err != nil {
			return nil, err
		}
		for _, f := range t.Values {
			if _, err := v.node(f.Value); err != nil {
				return nil, err
			}
		}
		return ty, nil
	case *ast.FnCall:
		return v.call(t, v.s.pkg.Scope)
	case *ast.NamespaceAccess:
		scope := v.s.pkg.Scope
		if imp, ok := v.s.file.Imports[t.Left.Name.Name()]; ok {
			scope = imp.Package.Scope
		} else {
			sym, err := scope.LookupAcrossScopes(t.Left.Name.Name())
			if err != nil {
				return nil, err
			}
			if ext, ok := sym.Node.(*ast.ExternDecl); ok {
				scope = ext.Scope
			} else if d, ok := sym.Node.(*ast.EnumDecl); ok {
				ty := &ast.ExprType{Kind: ast.EXPR_TYPE_ID, T: &ast.IdType{Name: t.Left.Name, Args: t.TypeArgs}}
				if err := v.typeExpr(ty, map[*ast.TypeAlias]bool{}); err != nil {
					return nil, err
				}
				return v.enumVariant(t, d, ty)
			} else {
				return nil, fmt.Errorf("%s: invalid namespace", t.Left.Name.Pos)
			}
		}
		if call, ok := t.Right.Node.(*ast.FnCall); ok {
			return v.call(call, scope)
		}
		if nested, ok := t.Right.Node.(*ast.NamespaceAccess); ok {
			symbol, err := scope.LookupAcrossScopes(nested.Left.Name.Name())
			if err != nil {
				return nil, err
			}
			if d, ok := symbol.Node.(*ast.EnumDecl); ok {
				ty := &ast.ExprType{Kind: ast.EXPR_TYPE_ID, T: &ast.IdType{Name: nested.Left.Name, Namespace: t.Left.Name, Args: nested.TypeArgs}}
				if err := v.typeExpr(ty, map[*ast.TypeAlias]bool{}); err != nil {
					return nil, err
				}
				return v.enumVariant(nested, d, ty)
			}
		}
		return nil, fmt.Errorf("unsupported nested template namespace")
	case *ast.TupleExpr:
		for _, n := range t.Exprs {
			_, err := v.node(n)
			if err != nil {
				return nil, err
			}
		}
		return nil, nil
	case *ast.VarArgsExpr:
		for _, n := range t.Args {
			if _, err := v.node(n); err != nil {
				return nil, err
			}
		}
	case *ast.VarStmt:
		inferred, err := v.node(t.Expr)
		if err != nil {
			return nil, err
		}
		for _, name := range t.Names {
			if variable, ok := name.Node.(*ast.VarIdStmt); ok {
				if err := v.typeExpr(variable.Type, map[*ast.TypeAlias]bool{}); err != nil {
					return nil, err
				}
				if t.IsDecl {
					ty := variable.Type
					if ty == nil {
						ty = inferred
					}
					v.locals[variable.Name.Name()] = ty
				} else if _, ok := v.locals[variable.Name.Name()]; !ok {
					return nil, fmt.Errorf("%s: '%s' not defined", variable.Name.Pos, variable.Name.Name())
				}
			} else if _, err := v.node(name); err != nil {
				return nil, err
			}
		}
	case *ast.ReturnStmt:
		return v.node(t.Value)
	case *ast.BlockStmt:
		return nil, v.block(t)
	case *ast.DeferStmt:
		previousDepth := v.matchDepth
		v.matchDepth = 0
		defer func() { v.matchDepth = previousDepth }()
		return v.node(t.Stmt)
	case *ast.CondStmt:
		conds := append([]*ast.IfElifCond{t.IfStmt}, t.ElifStmts...)
		for _, cond := range conds {
			if _, err := v.node(cond.Expr); err != nil {
				return nil, err
			}
			if err := v.block(cond.Block); err != nil {
				return nil, err
			}
		}
		if t.ElseStmt != nil {
			return nil, v.block(t.ElseStmt.Block)
		}
	case *ast.ForLoop:
		outer := v.locals
		v.locals = map[string]*ast.ExprType{}
		for k, val := range outer {
			v.locals[k] = val
		}
		defer func() { v.locals = outer }()
		if t.Range != nil {
			if _, err := v.node(t.Range); err != nil {
				return nil, err
			}
			v.locals[t.Binding.Name.Name()] = nil
			return nil, v.block(t.Block)
		}
		for _, n := range []*ast.Node{t.Init, t.Cond, t.Update} {
			if _, err := v.node(n); err != nil {
				return nil, err
			}
		}
		return nil, v.block(t.Block)
	case *ast.WhileLoop:
		if _, err := v.node(t.Cond); err != nil {
			return nil, err
		}
		return nil, v.block(t.Block)
	case *ast.AssignmentStmt:
		for _, n := range t.Values {
			if _, err := v.node(n); err != nil {
				return nil, err
			}
		}
		for _, n := range t.Targets {
			if t.Decl {
				if id, ok := n.Node.(*ast.IdExpr); ok {
					v.locals[id.Name.Name()] = nil
					continue
				}
			}
			if _, err := v.node(n); err != nil {
				return nil, err
			}
		}
	}
	return nil, nil
}

func (v *templateValidator) enumVariant(
	access *ast.NamespaceAccess,
	d *ast.EnumDecl,
	ty *ast.ExprType,
) (*ast.ExprType, error) {
	var name *token.Token
	switch right := access.Right.Node.(type) {
	case *ast.IdExpr:
		name = right.Name
	case *ast.StructLiteralExpr:
		name = right.Name
		for _, field := range right.Values {
			if _, err := v.node(field.Value); err != nil {
				return nil, err
			}
		}
	default:
		return nil, diagnostics.Diag{Title: "Invalid enum construction", Message: "expected a named enum variant", Span: access.Left.Name.Span, Help: "Use Enum::Variant or Enum::Variant{field: value}."}
	}
	for _, variant := range d.Variants {
		if variant.Name.Name() == name.Name() {
			return ty, nil
		}
	}
	return nil, diagnostics.Diag{
		Title:   "Unknown enum variant",
		Message: "enum does not declare this variant",
		Span:    name.Span,
		Help:    "Use a variant declared by this enum.",
		Related: []diagnostics.Related{{Span: d.Name.Span, Message: "enum declared here"}},
	}
}

func (v *templateValidator) call(call *ast.FnCall, scope *ast.Scope) (*ast.ExprType, error) {
	if call.CalleeExpr != nil {
		if _, err := v.node(call.CalleeExpr); err != nil {
			return nil, err
		}
		for _, arg := range call.Args {
			if _, err := v.node(arg); err != nil {
				return nil, err
			}
		}
		return nil, nil
	}
	if ty, local := v.locals[call.Name.Name()]; local {
		for _, arg := range call.Args {
			if _, err := v.node(arg); err != nil {
				return nil, err
			}
		}
		if ty != nil && ty.Kind == ast.EXPR_TYPE_FUNCTION {
			return ty.T.(*ast.FunctionType).Result, nil
		}
		return nil, nil // Concrete specialization checks dependent callback types.
	}
	_, lookupErr := scope.LookupAcrossScopes(call.Name.Name())
	if isConversionBuiltin(call.Name.Name()) && lookupErr != nil {
		if len(call.TypeArgs) != 1 || len(call.Args) != 1 {
			return nil, fmt.Errorf(
				"%s requires one type argument and one value argument",
				call.Name.Name(),
			)
		}
		if err := v.typeExpr(call.TypeArgs[0], map[*ast.TypeAlias]bool{}); err != nil {
			return nil, err
		}
		_, err := v.node(call.Args[0])
		return nil, err
	}
	if isMemoryBuiltin(call.Name.Name()) && lookupErr != nil {
		count := 0
		if call.Name.Name() == "slice_from_raw" {
			count = 2
		}
		if call.Name.Name() == "slice_data" {
			count = 1
		}
		if len(call.TypeArgs) != 1 || len(call.Args) != count {
			return nil, fmt.Errorf(
				"%s requires one type argument and %d value arguments",
				call.Name.Name(),
				count,
			)
		}
		if err := v.typeExpr(call.TypeArgs[0], map[*ast.TypeAlias]bool{}); err != nil {
			return nil, err
		}
		for _, arg := range call.Args {
			if _, err := v.node(arg); err != nil {
				return nil, err
			}
		}
		return nil, nil
	}
	if call.Builtin == "view" || call.Builtin == "slice" || call.Name.Name() == "len" {
		for _, arg := range call.Args {
			if _, err := v.node(arg); err != nil {
				return nil, err
			}
		}
		return nil, nil
	}
	sym, err := scope.LookupAcrossScopes(call.Name.Name())
	if err != nil {
		return nil, fmt.Errorf("%s: '%s' not defined", call.Name.Pos, call.Name.Name())
	}
	if sym.Kind == ast.KIND_VA_INTRINSIC {
		count := 0
		if call.Name.Name() == "va_arg" {
			count = 1
		}
		if len(call.Args) != 1 || len(call.TypeArgs) != count || call.AtOp != nil {
			return nil, diagnostics.Diag{
				Title:   "Invalid va_list operation",
				Message: "invalid argument count for va_list intrinsic",
				Span:    call.Name.Span,
				Help:    "Use va_start(args), va_copy(borrowed), va_arg<T>(&owner), va_borrow(&owner), or va_end(&owner).",
			}
		}
		for _, ty := range call.TypeArgs {
			if err := v.typeExpr(ty, map[*ast.TypeAlias]bool{}); err != nil {
				return nil, err
			}
		}
		for _, arg := range call.Args {
			if _, err := v.node(arg); err != nil {
				return nil, err
			}
		}
		switch call.Name.Name() {
		case "va_start", "va_copy":
			return &ast.ExprType{
				Kind:     ast.EXPR_TYPE_VA_LIST_OWNER,
				T:        &ast.VaListOwnerType{},
				Explicit: true,
			}, nil
		case "va_borrow":
			return &ast.ExprType{
				Kind:     ast.EXPR_TYPE_C_VA_LIST,
				T:        &ast.CVaListType{},
				Explicit: true,
			}, nil
		case "va_arg":
			return call.TypeArgs[0], nil
		default:
			return ast.NewBasicType(token.VOID_TYPE), nil
		}
	}
	var params []*token.Token
	var valueParams *ast.Params
	var result *ast.ExprType
	switch fn := sym.Node.(type) {
	case *ast.FnDecl:
		params, result = fn.TypeParams, fn.RetType
		valueParams = fn.Params
	case *ast.Proto:
		result = fn.RetType
		valueParams = fn.Params
	default:
		return nil, fmt.Errorf("%s: '%s' is not callable", call.Name.Pos, call.Name.Name())
	}
	if len(params) != len(call.TypeArgs) {
		return nil, fmt.Errorf(
			"%s: '%s' requires %d explicit type arguments, got %d",
			call.Name.Pos,
			call.Name.Name(),
			len(params),
			len(call.TypeArgs),
		)
	}
	if !valueParams.IsVariadic && len(call.Args) != len(valueParams.Fields) ||
		valueParams.IsVariadic && len(call.Args) < len(valueParams.Fields)-1 {
		return nil, fmt.Errorf(
			"%s: '%s' requires %d value arguments, got %d",
			call.Name.Pos,
			call.Name.Name(),
			len(valueParams.Fields),
			len(call.Args),
		)
	}
	bindings := map[string]*ast.ExprType{}
	for i, arg := range call.TypeArgs {
		if err := v.typeExpr(arg, map[*ast.TypeAlias]bool{}); err != nil {
			return nil, err
		}
		bindings[params[i].Name()] = arg
	}
	for _, arg := range call.Args {
		if _, err := v.node(arg); err != nil {
			return nil, err
		}
	}
	if call.AtOp != nil {
		if op, ok := call.AtOp.Op.(*ast.CatchAtOperator); ok {
			outer := v.locals
			v.locals = map[string]*ast.ExprType{}
			for k, ty := range outer {
				v.locals[k] = ty
			}
			v.locals[op.ErrVarName.Name()] = ast.NewBasicType(token.ERROR_TYPE)
			err := v.block(op.Block)
			v.locals = outer
			if err != nil {
				return nil, err
			}
		}
	}
	return cloneType(result, bindings), nil
}
