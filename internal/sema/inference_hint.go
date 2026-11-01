package sema

import (
	"fmt"
	"strings"

	"github.com/HicaroD/Telia/internal/ast"
	"github.com/HicaroD/Telia/internal/diagnostics"
	"github.com/HicaroD/Telia/internal/lexer"
	"github.com/HicaroD/Telia/internal/lexer/token"
	"github.com/HicaroD/Telia/internal/source"
)

// knownType resolves only existing types. Unlike resolveType it cannot create
// specializations, check bodies, register layouts, or emit diagnostics.
func (i *typeInference) knownType(t *ast.ExprType, scope *ast.Scope, file *ast.File) *ast.ExprType {
	if t == nil {
		return nil
	}
	switch p := t.T.(type) {
	case *ast.IdType:
		if p.Namespace != nil {
			if file == nil || file.Imports[p.Namespace.Name()] == nil {
				return nil
			}
			scope = file.Imports[p.Namespace.Name()].Package.Scope
		}
		n, err := scope.LookupAcrossScopes(p.Name.Name())
		if err != nil {
			return nil
		}
		switch d := n.Node.(type) {
		case *ast.TypeAlias:
			return i.knownType(d.Type, d.Owner.Scope, d.File)
		case *ast.StructDecl:
			if len(d.TypeParams) != len(p.Args) {
				return nil
			}
			if len(p.Args) == 0 {
				return &ast.ExprType{Kind: ast.EXPR_TYPE_STRUCT, T: &ast.StructType{Decl: d}}
			}
			args := make([]*ast.ExprType, len(p.Args))
			for n, arg := range p.Args {
				args[n] = i.knownType(arg, i.scope, file)
				if args[n] == nil {
					return nil
				}
			}
			for _, instance := range i.s.generics.structs {
				if instance.Template == d && sameTypeArguments(args, instance.TypeArgs) {
					return &ast.ExprType{Kind: ast.EXPR_TYPE_STRUCT, T: &ast.StructType{Decl: instance}}
				}
			}
		case *ast.EnumDecl:
			if len(d.TypeParams) != len(p.Args) {
				return nil
			}
			if len(p.Args) == 0 {
				return &ast.ExprType{Kind: ast.EXPR_TYPE_ENUM, T: &ast.EnumType{Decl: d}}
			}
			args := make([]*ast.ExprType, len(p.Args))
			for n, arg := range p.Args {
				args[n] = i.knownType(arg, i.scope, file)
				if args[n] == nil {
					return nil
				}
			}
			for _, instance := range i.s.generics.enums {
				if instance.Template == d && sameTypeArguments(args, instance.TypeArgs) {
					return &ast.ExprType{Kind: ast.EXPR_TYPE_ENUM, T: &ast.EnumType{Decl: instance}}
				}
			}
		}
		return nil
	case *ast.PointerType:
		inner := i.knownType(p.Type, scope, file)
		if inner == nil {
			return nil
		}
		return &ast.ExprType{Kind: t.Kind, T: &ast.PointerType{Type: inner, Const: p.Const}}
	case *ast.SequenceType:
		inner := i.knownType(p.Element, scope, file)
		if inner == nil {
			return nil
		}
		return &ast.ExprType{Kind: t.Kind, T: &ast.SequenceType{Element: inner, Length: p.Length, Foreign: p.Foreign}}
	case *ast.TupleType:
		parts := make([]*ast.ExprType, len(p.Types))
		for n, part := range p.Types {
			parts[n] = i.knownType(part, scope, file)
			if parts[n] == nil {
				return nil
			}
		}
		return &ast.ExprType{Kind: t.Kind, T: &ast.TupleType{Types: parts}}
	case *ast.FunctionType:
		f := cloneType(t, nil)
		fn := f.T.(*ast.FunctionType)
		for _, p := range fn.Params.Fields {
			p.Type = i.knownType(p.Type, scope, file)
			if p.Type == nil {
				return nil
			}
		}
		fn.Result = i.knownType(fn.Result, scope, file)
		if fn.Result == nil {
			return nil
		}
		return f
	default:
		return t
	}
}

func sameTypeArguments(a, b []*ast.ExprType) bool {
	if len(a) != len(b) {
		return false
	}
	for n := range a {
		if a[n] == nil || b[n] == nil || typeKey(a[n]) != typeKey(b[n]) {
			return false
		}
	}
	return true
}

func (i *typeInference) knownExpression(
	n *ast.Node,
	expected *ast.ExprType,
	scope *ast.Scope,
) *ast.ExprType {
	if n == nil {
		return nil
	}
	switch e := n.Node.(type) {
	case *ast.LiteralExpr:
		if expected != nil {
			return expected
		}
		if e.Type.IsInteger() {
			return ast.NewBasicType(token.INT_TYPE)
		}
		if e.Type.IsFloat() {
			return ast.NewBasicType(token.FLOAT_TYPE)
		}
		return e.Type
	case *ast.IdExpr:
		sym, err := scope.LookupAcrossScopes(e.Name.Name())
		if err != nil {
			return nil
		}
		switch d := sym.Node.(type) {
		case *ast.VarIdStmt:
			return d.Type
		case *ast.Param:
			return d.Type
		case *ast.ConstDecl:
			return i.knownType(d.Type, d.Owner.Scope, d.File)
		case *ast.FnDecl:
			if len(d.TypeParams) == 0 {
				return i.knownType(&ast.ExprType{Kind: ast.EXPR_TYPE_FUNCTION, T: &ast.FunctionType{Params: d.Params, Result: d.RetType, Attributes: d.Attributes}}, d.Owner.Scope, d.File)
			}
		}
	case *ast.AddressOfExpr:
		ty := i.knownExpression(e.Expr, nil, scope)
		if ty != nil {
			return ast.NewPointerType(ty)
		}
	case *ast.DerefPointerExpr:
		if ty := i.knownExpression(e.Expr, nil, scope); ty != nil {
			if p, ok := ty.T.(*ast.PointerType); ok {
				return p.Type
			}
		}
	case *ast.FieldAccess:
		if e.AccessedField != nil {
			return e.AccessedField.Type
		}
	case *ast.UnaryExpr:
		return i.knownExpression(e.Value, expected, scope)
	case *ast.FnCall:
		sym, err := scope.LookupAcrossScopes(e.Name.Name())
		if err != nil {
			return nil
		}
		fn, ok := sym.Node.(*ast.FnDecl)
		if !ok {
			return nil
		}
		if len(fn.TypeParams) == 0 {
			return i.knownType(fn.RetType, fn.Owner.Scope, fn.File)
		}
		bindings := map[string]*ast.ExprType{}
		if len(e.TypeArgs) != 0 && len(e.TypeArgs) != len(fn.TypeParams) {
			return nil
		}
		for n, arg := range e.TypeArgs {
			if inferenceHole(arg) {
				continue
			}
			ty := i.knownType(arg, scope, i.s.file)
			if ty == nil {
				return nil
			}
			bindings[fn.TypeParams[n].Name()] = ty
		}
		child := &typeInference{s: i.s, params: fn.TypeParams, bindings: bindings, scope: fn.Owner.Scope, file: fn.File, dry: true}
		if expected != nil {
			if child.match(inferenceCallResult(fn.RetType, e, expected), expected) != nil {
				return nil
			}
		}
		var patterns []*ast.ExprType
		for _, p := range fn.Params.Fields {
			patterns = append(patterns, p.Type)
		}
		if child.arguments(patterns, e.Args, scope) != nil {
			return nil
		}
		if _, err := child.result(); err != nil {
			return nil
		}
		return child.knownType(cloneType(fn.RetType, child.bindings), child.scope, child.file)
	}
	return nil // Unproven rewrites stay silent.
}

// suggestTypeArguments proves rewrites using read-only constraints. It never
// runs the ordinary checker speculatively or reads a specialization's inferred
// types as evidence for that specialization's own omitted arguments.
func (s *sema) suggestTypeArguments(
	site *token.Token,
	explicit, actual []*ast.ExprType,
	params []*token.Token,
	owner *ast.Package,
	file *ast.File,
	scope *ast.Scope,
	constrain func(*typeInference) error,
) {
	if len(explicit) == 0 || len(s.generics.chain) != 0 || s.inferenceProbeDepth != 0 ||
		s.inferenceDependentContext != 0 ||
		!site.Span.Valid() {
		return
	}
	span, original := typeApplicationSource(site)
	if !span.Valid() || len(original) != len(params) {
		return
	}
	if s.inferenceHints[span] {
		return
	}
	tail := applicationTail(span)
	rangeSpan := source.Span{File: span.File, Start: site.Span.Start, End: span.End + len(tail)}
	for _, outer := range s.inferenceHintRanges {
		if spanContains(outer, rangeSpan) {
			return
		}
	}
	prove := func(keep []bool) bool {
		bindings := map[string]*ast.ExprType{}
		for n, fixed := range keep {
			if fixed {
				bindings[params[n].Name()] = actual[n]
			}
		}
		i := &typeInference{
			s:        s,
			params:   params,
			bindings: bindings,
			scope:    owner.Scope,
			file:     file,
			dry:      true,
		}
		if constrain(i) != nil {
			return false
		}
		args, err := i.result()
		return err == nil && sameTypeArguments(args, actual)
	}
	keep := make([]bool, len(params))
	replacement := ""
	if !prove(keep) {
		for n, arg := range explicit {
			keep[n] = !inferenceHole(arg)
		}
		changed := false
		for n := range keep {
			if !keep[n] {
				continue
			}
			keep[n] = false
			if prove(keep) {
				changed = true
			} else {
				keep[n] = true
			}
		}
		if !changed {
			return
		}
		parts := make([]string, len(params))
		for n := range parts {
			parts[n] = "_"
			if keep[n] {
				parts[n] = original[n]
			}
		}
		replacement = "<" + strings.Join(parts, ", ") + ">"
	}
	if s.inferenceHints == nil {
		s.inferenceHints = map[source.Span]bool{}
	}
	s.inferenceHints[span] = true
	// Prefer an enclosing rewrite to overlapping child suggestions. Two removals
	// can each be safe against the original context without being safe together.
	kept := s.collector.Diags[:0]
	for _, diag := range s.collector.Diags {
		if diag.Title != "Type arguments can be inferred" || !spanContains(rangeSpan, diag.Span) {
			kept = append(kept, diag)
		}
	}
	s.collector.Diags = kept
	s.inferenceHintRanges = append(s.inferenceHintRanges, rangeSpan)
	s.collector.ReportAndSave(
		diagnostics.Diag{
			Severity: diagnostics.Hint,
			Title:    "Type arguments can be inferred",
			Message:  "Explicit type arguments can be shortened without changing the concrete types.",
			Span:     span,
			Help:     fmt.Sprintf("Write %s%s%s", site.Name(), replacement, tail),
		},
	)
}

func spanContains(outer, inner source.Span) bool {
	return outer.File == inner.File && outer.Start <= inner.Start && inner.End <= outer.End
}

func applicationTail(span source.Span) string {
	text := span.File.Text[span.End:]
	l := lexer.New(&ast.Loc{Name: span.File.Path}, []byte(text), diagnostics.New())
	end := 0
	if l.Peek().Kind == token.COLON_COLON {
		l.Skip()
		if l.Peek().Kind != token.ID {
			return ""
		}
		end = l.Peek().Span.End
		l.Skip()
	}
	if l.Peek().Kind == token.DOT && l.Peek1().Kind == token.OPEN_CURLY {
		l.Skip() // Preserve the legacy Type.{...} suffix in the suggested spelling.
	}
	if l.Peek().Kind != token.OPEN_PAREN && l.Peek().Kind != token.OPEN_CURLY {
		return text[:end]
	}
	depth := 0
	for {
		t := l.Peek()
		switch t.Kind {
		case token.EOF:
			return text[:end]
		case token.OPEN_PAREN, token.OPEN_CURLY, token.OPEN_BRACKET:
			depth++
		case token.CLOSE_PAREN, token.CLOSE_CURLY, token.CLOSE_BRACKET:
			depth--
			if depth == 0 {
				return text[:t.Span.End]
			}
		}
		l.Skip()
	}
}

// Lexing the list retains precise byte spans and handles comments and nested
// applications without mistaking their commas for top-level argument separators.
func typeApplicationSource(site *token.Token) (source.Span, []string) {
	span := site.Span
	if !span.Valid() {
		return source.Span{}, nil
	}
	text := span.File.Text
	start := span.End
	l := lexer.New(&ast.Loc{Name: span.File.Path}, []byte(text[start:]), diagnostics.New())
	if l.Peek().Kind != token.LESS {
		return source.Span{}, nil
	}
	open := l.Peek().Span.Start
	l.Skip()
	depth, part, nested := 1, open+1, 0
	var args []string
	for {
		t := l.Peek()
		if t.Kind == token.EOF {
			return source.Span{}, nil
		}
		switch t.Kind {
		case token.OPEN_PAREN, token.OPEN_BRACKET:
			nested++
		case token.CLOSE_PAREN, token.CLOSE_BRACKET:
			nested--
		case token.LESS:
			depth++
		case token.GREATER:
			depth--
			if depth == 0 {
				args = append(args, strings.TrimSpace(text[start+part:start+t.Span.Start]))
				return source.Span{
					File:  span.File,
					Start: start + open,
					End:   start + t.Span.End,
				}, args
			}
		case token.COMMA:
			if depth == 1 && nested == 0 {
				args = append(args, strings.TrimSpace(text[start+part:start+t.Span.Start]))
				part = t.Span.End
			}
		}
		l.Skip()
	}
}
