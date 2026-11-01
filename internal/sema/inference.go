package sema

import (
	"fmt"

	"github.com/HicaroD/Telia/internal/ast"
	"github.com/HicaroD/Telia/internal/lexer/token"
)

// typeInference collects concrete bindings before checking contextual expressions.
// Explicit bindings are fixed; untyped expressions are considered only after
// concrete argument types have contributed their constraints.
type typeInference struct {
	s        *sema
	params   []*token.Token
	bindings map[string]*ast.ExprType
	scope    *ast.Scope
	file     *ast.File
	dry      bool // read-only inference for advisory rewrites
}

// Context inferred from an expression's own explicit arguments is not independent
// evidence that those same arguments can be removed. Keep nested hints conservative.
func (s *sema) withInferenceDependency(dependent bool, check func() error) error {
	if dependent {
		s.inferenceDependentContext++
		defer func() { s.inferenceDependentContext-- }()
	}
	return check()
}

func (s *sema) assignmentContext(names []*ast.Node, scope *ast.Scope) (*ast.ExprType, error) {
	var types []*ast.ExprType
	for _, name := range names {
		var ty *ast.ExprType
		switch v := name.Node.(type) {
		case *ast.VarIdStmt:
			ty = v.Type
			if ty == nil && v.N != nil {
				switch original := v.N.Node.(type) {
				case *ast.VarIdStmt:
					ty = original.Type
				case *ast.Param:
					ty = original.Type
				}
			}
		case *ast.FieldAccess:
			_, field, err := s.getAccessedField(v, scope)
			if err != nil {
				return nil, err
			}
			ty = field.Type
		case *ast.IndexExpr:
			var err error
			ty, err = s.inferIndex(v, scope, scope, false)
			if err != nil {
				return nil, err
			}
		case *ast.DerefPointerExpr:
			var err error
			ty, _, err = s.inferDerefPtrExprTypeWithoutContext(v, scope, scope, false, false)
			if err != nil {
				return nil, err
			}
		}
		if ty == nil {
			return nil, nil
		}
		resolved, err := s.resolveType(ty, scope)
		if err != nil {
			return nil, err
		}
		types = append(types, resolved)
	}
	if len(types) == 1 {
		return types[0], nil
	}
	return &ast.ExprType{Kind: ast.EXPR_TYPE_TUPLE, T: &ast.TupleType{Types: types}}, nil
}

func (s *sema) newInference(
	params []*token.Token,
	args []*ast.ExprType,
	scope *ast.Scope,
) (*typeInference, error) {
	if len(args) != 0 && len(args) != len(params) {
		return nil, fmt.Errorf("requires %d type arguments, got %d", len(params), len(args))
	}
	i := &typeInference{s: s, params: params, bindings: map[string]*ast.ExprType{}, scope: scope}
	for n, arg := range args {
		if inferenceHole(arg) {
			continue
		}
		ty, err := s.resolveType(arg, scope)
		if err != nil {
			return nil, err
		}
		i.bindings[params[n].Name()] = ty
	}
	return i, nil
}

func inferenceHole(t *ast.ExprType) bool {
	if t == nil {
		return true
	}
	id, ok := t.T.(*ast.IdType)
	return ok && id.Namespace == nil && id.Name.Name() == "_"
}

func needsInference(args []*ast.ExprType) bool {
	if len(args) == 0 {
		return true
	}
	for _, arg := range args {
		if inferenceHole(arg) {
			return true
		}
	}
	return false
}

func (i *typeInference) parameter(t *ast.ExprType) string {
	id, ok := t.T.(*ast.IdType)
	if !ok || id.Namespace != nil {
		return ""
	}
	for _, p := range i.params {
		if p.Name() == id.Name.Name() {
			return p.Name()
		}
	}
	return ""
}

func (i *typeInference) match(pattern, actual *ast.ExprType) error {
	if pattern == nil || actual == nil || !inferenceValueType(actual) {
		return nil
	}
	if name := i.parameter(pattern); name != "" {
		if bound := i.bindings[name]; bound != nil {
			if !bound.Equals(actual) {
				return fmt.Errorf("conflicting types for '%s': %s and %s", name, bound.T, actual.T)
			}
		} else {
			i.bindings[name] = actual
		}
		return nil
	}
	switch p := pattern.T.(type) {
	case *ast.IdType:
		scope := i.scope
		if p.Namespace != nil && i.file != nil {
			if imp := i.file.Imports[p.Namespace.Name()]; imp != nil {
				scope = imp.Package.Scope
			}
		}
		sym, err := scope.LookupAcrossScopes(p.Name.Name())
		if err != nil {
			return err
		}
		var origin any
		var args []*ast.ExprType
		switch a := actual.T.(type) {
		case *ast.StructType:
			origin, args = a.Decl.Template, a.Decl.TypeArgs
		case *ast.EnumType:
			origin, args = a.Decl.Template, a.Decl.TypeArgs
		}
		if origin == sym.Node && len(args) == len(p.Args) {
			for n := range args {
				if err := i.match(p.Args[n], args[n]); err != nil {
					return err
				}
			}
		}
	case *ast.PointerType:
		if a, ok := actual.T.(*ast.PointerType); ok {
			return i.match(p.Type, a.Type)
		}
	case *ast.SequenceType:
		if a, ok := actual.T.(*ast.SequenceType); ok {
			return i.match(p.Element, a.Element)
		}
	case *ast.TupleType:
		if a, ok := actual.T.(*ast.TupleType); ok && len(p.Types) == len(a.Types) {
			for n := range p.Types {
				if err := i.match(p.Types[n], a.Types[n]); err != nil {
					return err
				}
			}
		}
	case *ast.FunctionType:
		if a, ok := actual.T.(*ast.FunctionType); ok && len(p.Params.Fields) == len(a.Params.Fields) {
			for n := range p.Params.Fields {
				if err := i.match(p.Params.Fields[n].Type, a.Params.Fields[n].Type); err != nil {
					return err
				}
			}
			return i.match(p.Result, a.Result)
		}
	}
	return nil
}

func inferenceValueType(t *ast.ExprType) bool {
	return valueType(t) && inferenceConcreteType(t)
}

// Pointees need a known identity, not a by-value layout: opaque types and void
// are valid behind pointers, while an untyped nil pointee supplies no evidence.
func inferenceConcreteType(t *ast.ExprType) bool {
	switch p := t.T.(type) {
	case *ast.IdType:
		return false
	case *ast.BasicType:
		return p.Kind <= token.UNTYPED_START || p.Kind >= token.UNTYPED_END
	case *ast.PointerType:
		return inferenceConcreteType(p.Type)
	case *ast.TupleType:
		for _, part := range p.Types {
			if !inferenceConcreteType(part) {
				return false
			}
		}
	case *ast.SequenceType:
		return inferenceConcreteType(p.Element)
	case *ast.FunctionType:
		for _, param := range p.Params.Fields {
			if !inferenceConcreteType(param.Type) {
				return false
			}
		}
		return inferenceConcreteType(p.Result)
	}
	return true
}

func (i *typeInference) arguments(
	patterns []*ast.ExprType,
	values []*ast.Node,
	scope *ast.Scope,
) error {
	for pass := 0; pass < 3; pass++ {
		for n, value := range values {
			if n >= len(patterns) {
				break
			}
			if value == nil {
				continue
			} // Some lowered builtin calls carry omitted operands.
			literal := value.Kind == ast.KIND_LITERAL_EXPR || value.Kind == ast.KIND_UNARY_EXPR ||
				value.Kind == ast.KIND_BINARY_EXPR
			stage := i.expressionStage(value, scope)
			if stage != pass {
				continue
			}
			if literal, ok := value.Node.(*ast.ArrayLiteralExpr); ok {
				if shape, ok := patterns[n].T.(*ast.SequenceType); ok &&
					patterns[n].Kind == ast.EXPR_TYPE_ARRAY {
					elements := make([]*ast.ExprType, len(literal.Values))
					for n := range elements {
						elements[n] = shape.Element
					}
					if err := i.arguments(elements, literal.Values, scope); err != nil {
						return err
					}
				}
			}
			// A bound literal must be checked against its context, not used to
			// replace an already-known width with the default integer width.
			if literal && i.parameter(patterns[n]) != "" &&
				i.bindings[i.parameter(patterns[n])] != nil {
				if lit, ok := value.Node.(*ast.LiteralExpr); !ok || lit.Type.IsNumeric() {
					continue
				}
			}
			var ty *ast.ExprType
			var err error
			if i.dry {
				ty = i.knownExpression(value, i.concrete(patterns[n]), scope)
			} else if context := i.concrete(patterns[n]); context != nil && stage >= 1 {
				ty, err = i.s.probeInferenceExpression(value, context, scope)
			} else {
				ty, err = i.s.probeInferenceExpression(value, nil, scope)
			}
			if err != nil {
				continue
			} // Context-dependent expressions are checked once bindings are complete.
			if err := i.match(patterns[n], ty); err != nil {
				return err
			}
		}
		// A generic callback may learn its input from a sibling value and in
		// turn determine the outer call's result. Revisit until bindings stop growing.
		for round := 0; round < len(i.params); round++ {
			before := len(i.bindings)
			for n, value := range values {
				if n >= len(patterns) {
					break
				}
				if evidence := i.callbackEvidence(value, patterns[n], scope); evidence != nil {
					if err := i.match(patterns[n], evidence); err != nil {
						return err
					}
				}
			}
			if len(i.bindings) == before {
				break
			}
		}
	}
	return nil
}

func (s *sema) probeInferenceExpression(
	value *ast.Node,
	context *ast.ExprType,
	scope *ast.Scope,
) (*ast.ExprType, error) {
	s.inferenceProbeDepth++
	defer func() { s.inferenceProbeDepth-- }()
	probe := (syntaxClone{}).node(value, scope)
	if context != nil {
		return s.inferExprTypeWithContext(probe, context, scope, scope, false, true)
	}
	ty, _, err := s.inferExprTypeWithoutContext(probe, scope, scope, false, true)
	return ty, err
}

func (i *typeInference) callbackEvidence(
	value *ast.Node,
	pattern *ast.ExprType,
	scope *ast.Scope,
) *ast.ExprType {
	if value == nil {
		return nil
	}
	want, ok := pattern.T.(*ast.FunctionType)
	if !ok {
		return nil
	}
	id, ok := value.Node.(*ast.IdExpr)
	lookupScope := scope
	if !ok {
		access, namespace := value.Node.(*ast.NamespaceAccess)
		if !namespace || !access.IsImport || i.s.file.Imports[access.Left.Name.Name()] == nil {
			return nil
		}
		id, ok = access.Right.Node.(*ast.IdExpr)
		if !ok {
			return nil
		}
		lookupScope = i.s.file.Imports[access.Left.Name.Name()].Package.Scope
	}
	symbol, err := lookupScope.LookupAcrossScopes(id.Name.Name())
	if err != nil {
		return nil
	}
	fn, ok := symbol.Node.(*ast.FnDecl)
	if !ok || len(fn.TypeParams) == 0 || len(fn.Params.Fields) != len(want.Params.Fields) {
		return nil
	}
	inner := &typeInference{
		s:        i.s,
		params:   fn.TypeParams,
		bindings: map[string]*ast.ExprType{},
		scope:    fn.Owner.Scope,
		file:     fn.File,
		dry:      i.dry,
	}
	if len(id.TypeArgs) != 0 && len(id.TypeArgs) != len(fn.TypeParams) {
		return nil
	}
	for n, arg := range id.TypeArgs {
		if inferenceHole(arg) {
			continue
		}
		var ty *ast.ExprType
		if i.dry {
			ty = i.knownType(arg, scope, i.s.file)
		} else {
			ty, err = i.s.resolveType(arg, scope)
		}
		if err != nil || ty == nil {
			return nil
		}
		inner.bindings[fn.TypeParams[n].Name()] = ty
	}
	for n, param := range fn.Params.Fields {
		if inner.match(param.Type, i.concrete(want.Params.Fields[n].Type)) != nil {
			return nil
		}
	}
	if inner.match(fn.RetType, i.concrete(want.Result)) != nil {
		return nil
	}
	if _, err := inner.result(); err != nil {
		return nil
	}
	signature := cloneType(
		&ast.ExprType{
			Kind: ast.EXPR_TYPE_FUNCTION,
			T: &ast.FunctionType{
				Params:     fn.Params,
				Result:     fn.RetType,
				Attributes: fn.Attributes,
			},
		},
		inner.bindings,
	)
	if i.dry {
		return inner.knownType(signature, fn.Owner.Scope, fn.File)
	}
	var resolved *ast.ExprType
	err = i.s.inContext(fn.Owner, fn.File, inner.bindings, func() error {
		var err error
		resolved, err = i.s.resolveType(signature, fn.Owner.Scope)
		return err
	})
	if err != nil {
		return nil
	}
	return resolved
}

// Evidence with a fixed type wins over expressions which still need a numeric
// default. This also applies to literals nested inside calls and constructors.
func (i *typeInference) expressionStage(n *ast.Node, scope *ast.Scope) int {
	if n == nil {
		return 0
	}
	switch e := n.Node.(type) {
	case *ast.LiteralExpr:
		if e.Type.IsNumeric() {
			return 2
		}
	case *ast.UnaryExpr:
		return i.expressionStage(e.Value, scope)
	case *ast.BinExpr:
		return max(i.expressionStage(e.Left, scope), i.expressionStage(e.Right, scope))
	case *ast.FnCall:
		sym, err := scope.LookupAcrossScopes(e.Name.Name())
		if err == nil {
			if fn, ok := sym.Node.(*ast.FnDecl); ok && (len(fn.TypeParams) == 0 || !needsInference(e.TypeArgs)) {
				return 0
			}
		}
		stage := 1
		for _, a := range e.Args {
			stage = max(stage, i.expressionStage(a, scope))
		}
		return stage
	case *ast.StructLiteralExpr:
		if !needsInference(e.TypeArgs) {
			return 0
		}
		stage := 1
		for _, f := range e.Values {
			stage = max(stage, i.expressionStage(f.Value, scope))
		}
		return stage
	case *ast.ArrayLiteralExpr:
		stage := 1
		for _, value := range e.Values {
			stage = max(stage, i.expressionStage(value, scope))
		}
		return stage
	case *ast.NamespaceAccess:
		if e.IsImport && i.s.file != nil {
			if imp := i.s.file.Imports[e.Left.Name.Name()]; imp != nil {
				return i.expressionStage(e.Right, imp.Package.Scope)
			}
		}
		if !needsInference(e.TypeArgs) {
			return 0
		}
		return max(1, i.expressionStage(e.Right, scope))
	}
	return 0
}

func (i *typeInference) concrete(pattern *ast.ExprType) *ast.ExprType {
	if pattern == nil {
		return nil
	}
	if name := i.parameter(pattern); name != "" {
		return cloneType(i.bindings[name], nil)
	}
	// Do not resolve an incomplete pattern (resolution is not an inference operation).
	var complete func(*ast.ExprType) bool
	complete = func(t *ast.ExprType) bool {
		if t == nil {
			return true
		}
		if name := i.parameter(t); name != "" {
			return i.bindings[name] != nil
		}
		switch p := t.T.(type) {
		case *ast.IdType:
			for _, a := range p.Args {
				if !complete(a) {
					return false
				}
			}
		case *ast.PointerType:
			return complete(p.Type)
		case *ast.SequenceType:
			return complete(p.Element)
		case *ast.TupleType:
			for _, a := range p.Types {
				if !complete(a) {
					return false
				}
			}
		case *ast.FunctionType:
			for _, p := range p.Params.Fields {
				if !complete(p.Type) {
					return false
				}
			}
			return complete(p.Result)
		}
		return true
	}
	if !complete(pattern) {
		return nil
	}
	if i.dry {
		return i.knownType(cloneType(pattern, i.bindings), i.scope, i.file)
	}
	var ty *ast.ExprType
	err := i.s.inContext(i.s.pkg, i.file, i.s.generics.bindings, func() error {
		var err error
		ty, err = i.s.resolveType(cloneType(pattern, i.bindings), i.scope)
		return err
	})
	if err != nil {
		return nil
	}
	return ty
}

func (i *typeInference) result() ([]*ast.ExprType, error) {
	args := make([]*ast.ExprType, len(i.params))
	for n, p := range i.params {
		args[n] = i.bindings[p.Name()]
		if args[n] == nil {
			return nil, fmt.Errorf(
				"cannot infer type parameter '%s'; supply an explicit type argument",
				p.Name(),
			)
		}
	}
	return args, nil
}

func (s *sema) typeSymbol(name, namespace *token.Token, scope *ast.Scope) (*ast.Node, error) {
	if namespace != nil {
		if imp := s.file.Imports[namespace.Name()]; imp != nil {
			scope = imp.Package.Scope
		} else {
			n, err := scope.LookupAcrossScopes(namespace.Name())
			if err != nil {
				return nil, err
			}
			if ext, ok := n.Node.(*ast.ExternDecl); ok {
				scope = ext.Scope
			}
		}
	}
	return scope.LookupAcrossScopes(name.Name())
}

func (i *typeInference) fields(
	fields []*ast.StructField,
	values []*ast.StructFieldValue,
	scope *ast.Scope,
) error {
	var patterns []*ast.ExprType
	var exprs []*ast.Node
	for _, value := range values {
		for _, field := range fields {
			if field.Name.Name() == value.Name.Name() {
				patterns = append(patterns, field.Type)
				exprs = append(exprs, value.Value)
				break
			}
		}
	}
	return i.arguments(patterns, exprs, scope)
}

func (i *typeInference) expectedInstance(template any, expected *ast.ExprType) error {
	if expected == nil {
		return nil
	}
	var origin any
	var args []*ast.ExprType
	switch e := expected.T.(type) {
	case *ast.StructType:
		origin, args = e.Decl.Template, e.Decl.TypeArgs
	case *ast.EnumType:
		origin, args = e.Decl.Template, e.Decl.TypeArgs
	}
	if origin != template {
		return nil
	}
	for n, arg := range args {
		name := i.params[n].Name()
		if bound := i.bindings[name]; bound != nil && !bound.Equals(arg) {
			return fmt.Errorf("conflicting types for '%s': %s and %s", name, bound.T, arg.T)
		}
		i.bindings[name] = arg
	}
	return nil
}

func (s *sema) inferCallArguments(
	fn *ast.FnDecl,
	call *ast.FnCall,
	expected *ast.ExprType,
	scope *ast.Scope,
) ([]*ast.ExprType, error) {
	i, err := s.newInference(fn.TypeParams, call.TypeArgs, scope)
	if err != nil {
		return nil, err
	}
	i.scope, i.file = fn.Owner.Scope, fn.File
	if expected != nil {
		if err := i.match(inferenceCallResult(fn.RetType, call, expected), expected); err != nil {
			return nil, err
		}
	}
	var patterns []*ast.ExprType
	for _, p := range fn.Params.Fields {
		patterns = append(patterns, p.Type)
	}
	if err := i.arguments(patterns, call.Args, scope); err != nil {
		return nil, err
	}
	return i.result()
}

func inferenceCallResult(
	result *ast.ExprType,
	call *ast.FnCall,
	expected *ast.ExprType,
) *ast.ExprType {
	if call.AtOp == nil || result.Kind != ast.EXPR_TYPE_TUPLE {
		return result
	}
	parts := result.T.(*ast.TupleType).Types
	if len(parts) < 2 || !parts[len(parts)-1].IsError() {
		return result
	}
	if expected != nil && expected.Kind == ast.EXPR_TYPE_TUPLE &&
		len(expected.T.(*ast.TupleType).Types) == len(parts) {
		return result
	}
	if len(parts) == 2 {
		return parts[0]
	}
	return &ast.ExprType{Kind: ast.EXPR_TYPE_TUPLE, T: &ast.TupleType{Types: parts[:len(parts)-1]}}
}
