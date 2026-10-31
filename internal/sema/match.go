package sema

import (
	"strings"

	"github.com/HicaroD/Telia/internal/ast"
	"github.com/HicaroD/Telia/internal/diagnostics"
)

type yieldSite struct {
	expr     *ast.Node
	scope    *ast.Scope
	ty       *ast.ExprType
	explicit bool
}
type yieldContext struct{ sites []yieldSite }

func (s *sema) checkMatch(
	m *ast.MatchExpr,
	expected *ast.ExprType,
	scope *ast.Scope,
) (*ast.ExprType, error) {
	if m.Type != nil {
		if expected != nil && !m.Type.Equals(expected) {
			return nil, diagnostics.Diag{
				Title:   "Match result type mismatch",
				Message: "match result does not have the expected type",
				Span:    m.Value.Span,
				Help:    "Use compatible result types in all match arms.",
			}
		}
		return m.Type, nil
	}
	ty, _, err := s.inferExprTypeWithoutContext(m.Value, scope, scope, false, false)
	if err != nil {
		return nil, err
	}
	if ty.Kind != ast.EXPR_TYPE_ENUM && !ty.IsError() && !ty.IsInteger() {
		return nil, diagnostics.Diag{
			Title:   "Invalid match value",
			Message: "match requires an enum, error, or integer value",
			Span:    m.Value.Span,
			Help:    "Match an enum value using its named variants.",
		}
	}
	m.ValueType = ty
	if ty.IsInteger() {
		if err := s.checkIntegerPatterns(m, ty, scope); err != nil {
			return nil, err
		}
	} else {
		var d *ast.EnumDecl
		m.IsError = ty.IsError()
		if !m.IsError {
			d = ty.T.(*ast.EnumType).Decl
		}
		m.Enum = d
		seen := map[*ast.EnumVariant]*ast.MatchArm{}
		var fallback *ast.MatchArm
		for _, group := range m.Arms {
			arm := group
			// Semantic scopes may differ from the parser's scopes (notably @catch).
			arm.Scope.Parent = scope
			if fallback != nil {
				return nil, diagnostics.Diag{
					Title:   "Unreachable match arm",
					Message: "an earlier fallback already matches every remaining value",
					Span:    arm.VariantName.Span,
					Help:    "Keep only one fallback and put it last.",
					Related: []diagnostics.Related{
						{Span: fallback.VariantName.Span, Message: "fallback arm"},
					},
				}
			}
			for index, pattern := range group.Patterns() {
				patternScope := group.Scope
				if index != 0 {
					patternScope = ast.NewScope(scope)
				}
				arm := pattern
				if arm.Fallback || arm.Binding != nil {
					if group.Guard == nil {
						fallback = group
					}
					if arm.Binding != nil {
						arm.Binding.Type = ty
						if err := patternScope.Insert(arm.Binding.Name.Name(), &ast.Node{Kind: ast.KIND_PARAM, Node: arm.Binding, Span: arm.Binding.Name.Span}); err != nil {
							return nil, err
						}
					}
					continue
				}
				if arm.PatternType == nil {
					return nil, diagnostics.Diag{Title: "Invalid enum pattern", Message: "expected an enum variant", Span: arm.VariantName.Span, Help: "Use a variant of the matched enum."}
				}
				patternType, err := s.resolveType(arm.PatternType, scope)
				if err != nil {
					return nil, err
				}
				if m.IsError {
					if patternType.Kind != ast.EXPR_TYPE_ENUM ||
						!patternType.T.(*ast.EnumType).Decl.IsError {
						return nil, diagnostics.Diag{
							Message: "error pattern requires an enum declared with : error",
							Span:    arm.VariantName.Span,
						}
					}
					d = patternType.T.(*ast.EnumType).Decl
				} else if !ty.Equals(patternType) {
					return nil, diagnostics.Diag{
						Title:   "Wrong enum in pattern",
						Message: "pattern belongs to a different enum",
						Span:    arm.VariantName.Span,
						Help:    "Use variants of the matched enum.",
						Related: []diagnostics.Related{
							{Span: d.Name.Span, Message: "matched enum declared here"},
						},
					}
				}
				arm.Enum = d
				for _, v := range d.Variants {
					if v.Name.Name() == arm.VariantName.Name() {
						arm.Variant = v
						break
					}
				}
				if arm.Variant == nil {
					return nil, diagnostics.Diag{
						Title:   "Unknown enum variant",
						Message: "variant is not declared by the matched enum",
						Span:    arm.VariantName.Span,
						Help:    "Use a variant declared in this enum.",
						Related: []diagnostics.Related{{Span: d.Name.Span, Message: "enum declared here"}},
					}
				}
				if previous := seen[arm.Variant]; previous != nil && previous != group {
					return nil, diagnostics.Diag{
						Title:   "Duplicate match arm",
						Message: "variant has already been matched",
						Span:    arm.VariantName.Span,
						Help:    "Match each variant once.",
						Related: []diagnostics.Related{
							{Span: previous.VariantName.Span, Message: "previous arm"},
						},
					}
				}
				if group.Guard == nil {
					seen[arm.Variant] = group
				}
				fields := map[string]*ast.PatternField{}
				for _, binding := range arm.Fields {
					if previous := fields[binding.Name.Name()]; previous != nil {
						return nil, diagnostics.Diag{
							Title:   "Duplicate pattern field",
							Message: "payload field is bound twice",
							Span:    binding.Name.Span,
							Help:    "Bind each payload field at most once.",
							Related: []diagnostics.Related{
								{Span: previous.Name.Span, Message: "previous binding"},
							},
						}
					}
					fields[binding.Name.Name()] = binding
					for _, f := range arm.Variant.Fields {
						if f.Name.Name() == binding.Name.Name() {
							binding.Binding.Type = f.Type
							break
						}
					}
					if binding.Binding.Type == nil {
						return nil, diagnostics.Diag{
							Title:   "Unknown payload field",
							Message: "pattern names an undeclared payload field",
							Span:    binding.Name.Span,
							Help:    "Use a field declared in this variant.",
							Related: []diagnostics.Related{
								{Span: arm.Variant.Name.Span, Message: "variant declared here"},
							},
						}
					}
					node := &ast.Node{
						Kind: ast.KIND_PARAM,
						Node: binding.Binding,
						Span: binding.Binding.Name.Span,
					}
					if previous := patternScope.Nodes[binding.Binding.Name.Name()]; previous != nil {
						return nil, diagnostics.Diag{
							Title:   "Duplicate pattern binding",
							Message: "two fields use the same binding name",
							Span:    binding.Binding.Name.Span,
							Help:    "Give the bindings distinct names.",
							Related: []diagnostics.Related{
								{Span: previous.DeclarationSpan(), Message: "previous binding"},
							},
						}
					}
					if err := patternScope.Insert(binding.Binding.Name.Name(), node); err != nil {
						return nil, err
					}
				}
				if !arm.OmitRest && len(fields) != len(arm.Variant.Fields) {
					return nil, diagnostics.Diag{
						Title:   "Incomplete payload pattern",
						Message: "pattern omits payload fields without ..",
						Span:    arm.VariantName.Span,
						Help:    "Bind every payload field or add .. to omit the remaining fields.",
						Related: []diagnostics.Related{
							{Span: arm.Variant.Name.Span, Message: "variant declared here"},
						},
					}
				}
			}
			if err := shareMatchBindings(group); err != nil {
				return nil, err
			}
		}
		if fallback == nil && m.IsError {
			return nil, diagnostics.Diag{
				Message: "matching error requires a final fallback arm",
				Span:    m.Value.Span,
			}
		}
		if fallback == nil && d.Foreign != nil {
			return nil, diagnostics.Diag{
				Title:   "Foreign enum match requires fallback",
				Message: "foreign enums may contain values not named in the declaration",
				Span:    m.Value.Span,
				Help:    "Add a final _ => result arm for unknown C enum values.",
				Related: []diagnostics.Related{
					{Span: d.Name.Span, Message: "foreign enum declared here"},
				},
			}
		}
		if fallback == nil && len(seen) != len(d.Variants) {
			var missing []string
			var related []diagnostics.Related
			for _, v := range d.Variants {
				if seen[v] == nil {
					missing = append(missing, v.Name.Name())
					related = append(
						related,
						diagnostics.Related{Span: v.Name.Span, Message: "unmatched variant"},
					)
				}
			}
			return nil, diagnostics.Diag{
				Title:   "Non-exhaustive match",
				Message: "missing variants: " + strings.Join(missing, ", "),
				Span:    m.Value.Span,
				Help:    "Add arms for the missing variants or a final _ fallback.",
				Related: related,
			}
		}
	}
	if len(m.Arms) == 0 {
		return nil, diagnostics.Diag{
			Title:   "Empty match",
			Message: "a value-producing match needs a result arm",
			Span:    m.Value.Span,
			Help:    "Add a variant arm or a fallback that produces a value.",
		}
	}
	var sites []yieldSite
	for _, arm := range m.Arms {
		if arm.Guard != nil {
			guardType, _, err := s.inferExprTypeWithoutContext(
				arm.Guard,
				arm.Scope,
				arm.Scope,
				false,
				false,
			)
			if err != nil {
				return nil, err
			}
			if !guardType.IsBoolean() {
				return nil, diagnostics.Diag{
					Title:   "Invalid match guard",
					Message: "match guard must be boolean",
					Span:    arm.Guard.Span,
					Help:    "Use a boolean condition after if.",
				}
			}
		}
		if arm.Block != nil {
			ctx := &yieldContext{}
			previous := s.currentYield
			s.currentYield = ctx
			err := s.checkBlockStmt(
				arm.Block,
				arm.Scope,
				arm.Scope,
				s.currentReturnType,
				false,
				false,
			)
			s.currentYield = previous
			if err != nil {
				return nil, err
			}
			if !blockYields(arm.Block) {
				return nil, diagnostics.Diag{
					Title:   "Missing match yield",
					Message: "match arm block can finish without producing a value",
					Span:    arm.VariantName.Span,
					Help:    "Use yield expression on every path that does not return or jump to an enclosing loop.",
				}
			}
			sites = append(sites, ctx.sites...)
		} else {
			actual, explicit, err := s.inferExprTypeWithoutContext(arm.Result, arm.Scope, arm.Scope, false, false)
			if err != nil {
				return nil, err
			}
			sites = append(sites, yieldSite{expr: arm.Result, scope: arm.Scope, ty: actual, explicit: explicit})
		}
	}
	result := expected
	if result == nil {
		for _, site := range sites {
			if result == nil || site.explicit {
				result = site.ty
				if site.explicit {
					break
				}
			}
		}
	}
	if result == nil {
		return nil, diagnostics.Diag{
			Title:   "Missing match result",
			Message: "match has no value-producing arm",
			Span:    m.Value.Span,
			Help:    "Produce a value with an expression arm or yield.",
		}
	}
	if !valueType(result) {
		return nil, diagnostics.Diag{
			Title:   "Invalid match result",
			Message: "match arms must produce a value",
			Span:    m.Value.Span,
			Help:    "Return compatible non-void values from every arm.",
		}
	}
	for _, site := range sites {
		if _, err := s.inferExprTypeWithContext(site.expr, result, site.scope, site.scope, false, false); err != nil {
			return nil, err
		}
	}
	m.Type = result
	return result, nil
}

// All alternatives expose the same locals to their shared guard and result.
func shareMatchBindings(arm *ast.MatchArm) error {
	bindings := func(p *ast.MatchPattern) map[string]*ast.Param {
		result := map[string]*ast.Param{}
		if p.Binding != nil {
			result[p.Binding.Name.Name()] = p.Binding
		}
		for _, field := range p.Fields {
			result[field.Binding.Name.Name()] = field.Binding
		}
		return result
	}
	first := bindings(&arm.MatchPattern)
	for _, pattern := range arm.Alternatives {
		next := bindings(pattern)
		valid := len(first) == len(next)
		for name, param := range first {
			other := next[name]
			if other == nil || !param.Type.Equals(other.Type) {
				valid = false
			}
		}
		if !valid {
			return diagnostics.Diag{
				Title:   "Incompatible alternative bindings",
				Message: "alternatives must bind the same names with the same types",
				Span:    pattern.VariantName.Span,
				Help:    "Use matching binding names and types in every alternative, or split the arms.",
			}
		}
		if pattern.Binding != nil {
			pattern.Binding = first[pattern.Binding.Name.Name()]
		}
		for _, field := range pattern.Fields {
			field.Binding = first[field.Binding.Name.Name()]
		}
	}
	return nil
}

func (s *sema) checkIntegerPatterns(m *ast.MatchExpr, ty *ast.ExprType, scope *ast.Scope) error {
	fallback := false
	seen := map[string]*ast.MatchArm{}
	for _, arm := range m.Arms {
		arm.Scope.Parent = scope
		if fallback {
			return diagnostics.Diag{
				Title:   "Unreachable match arm",
				Message: "an earlier fallback already matches every remaining value",
				Span:    arm.VariantName.Span,
				Help:    "Put the fallback last.",
			}
		}
		for index, pattern := range arm.Patterns() {
			patternScope := arm.Scope
			if index != 0 {
				patternScope = ast.NewScope(scope)
			}
			if pattern.Fallback || pattern.Binding != nil {
				fallback = arm.Guard == nil
				if pattern.Binding != nil {
					pattern.Binding.Type = ty
					if err := patternScope.Insert(pattern.Binding.Name.Name(), &ast.Node{Kind: ast.KIND_PARAM, Node: pattern.Binding, Span: pattern.Binding.Name.Span}); err != nil {
						return err
					}
				}
				continue
			}
			if pattern.Literal == nil {
				return diagnostics.Diag{
					Title:   "Invalid integer pattern",
					Message: "expected an integer literal or _",
					Span:    pattern.VariantName.Span,
					Help:    "Use integer literals in an integer match.",
				}
			}
			literal := pattern.Literal
			if literal.Kind == ast.KIND_UNARY_EXPR {
				literal = literal.Node.(*ast.UnaryExpr).Value
			}
			if literal.Kind != ast.KIND_LITERAL_EXPR {
				return diagnostics.Diag{
					Title:   "Invalid match pattern",
					Message: "expected an integer literal",
					Span:    pattern.Literal.Span,
					Help:    "Use an integer literal; put comparisons and other conditions in an if guard.",
				}
			}
			if _, err := s.inferExprTypeWithContext(pattern.Literal, ty, arm.Scope, arm.Scope, false, false); err != nil {
				return err
			}
			value, err := evaluateConstant(pattern.Literal)
			if err != nil {
				return err
			}
			if err := checkConstantRange(value, ty, pattern.Literal.Span); err != nil {
				return err
			}
			key := value.ExactString()
			if previous := seen[key]; previous != nil && previous != arm {
				return diagnostics.Diag{
					Title:   "Duplicate match arm",
					Message: "integer value has already been matched",
					Span:    pattern.VariantName.Span,
					Help:    "Remove the covered pattern or put guarded cases before the unguarded case.",
				}
			}
			if arm.Guard == nil {
				seen[key] = arm
			}
		}
		if err := shareMatchBindings(arm); err != nil {
			return err
		}
	}
	if !fallback {
		return diagnostics.Diag{
			Title:   "Non-exhaustive match",
			Message: "integer match requires a fallback",
			Span:    m.Value.Span,
			Help:    "Add a final _ => result arm.",
		}
	}
	return nil
}

func blockYields(block *ast.BlockStmt) bool {
	for _, stmt := range block.Statements {
		switch stmt.Kind {
		case ast.KIND_YIELD_STMT, ast.KIND_RETURN_STMT, ast.KIND_BREAK_STMT, ast.KIND_CONTINUE_STMT:
			return true
		case ast.KIND_BLOCK_STMT:
			if blockYields(stmt.Node.(*ast.BlockStmt)) {
				return true
			}
		case ast.KIND_COND_STMT:
			cond := stmt.Node.(*ast.CondStmt)
			if cond.ElseStmt == nil || !blockYields(cond.IfStmt.Block) ||
				!blockYields(cond.ElseStmt.Block) {
				continue
			}
			all := true
			for _, branch := range cond.ElifStmts {
				all = all && blockYields(branch.Block)
			}
			if all {
				return true
			}
		}
	}
	return false
}
