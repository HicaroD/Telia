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

func (s *sema) checkMatch(m *ast.MatchExpr, expected *ast.ExprType, scope *ast.Scope) (*ast.ExprType, error) {
	if m.Type != nil {
		if expected != nil && !m.Type.Equals(expected) {
			return nil, diagnostics.Diag{Title: "Match result type mismatch", Message: "match result does not have the expected type", Span: m.Value.Span, Help: "Use compatible result types in all match arms."}
		}
		return m.Type, nil
	}
	ty, _, err := s.inferExprTypeWithoutContext(m.Value, scope, scope, false, false)
	if err != nil {
		return nil, err
	}
	if ty.Kind != ast.EXPR_TYPE_ENUM {
		return nil, diagnostics.Diag{Title: "Invalid match value", Message: "match requires an enum value", Span: m.Value.Span, Help: "Match an enum value using its named variants."}
	}
	d := ty.T.(*ast.EnumType).Decl
	m.Enum = d
	seen := map[*ast.EnumVariant]*ast.MatchArm{}
	var fallback *ast.MatchArm
	for _, arm := range m.Arms {
		if fallback != nil {
			return nil, diagnostics.Diag{Title: "Unreachable match arm", Message: "an earlier fallback already matches every remaining value", Span: arm.VariantName.Span, Help: "Keep only one fallback and put it last.", Related: []diagnostics.Related{{Span: fallback.VariantName.Span, Message: "fallback arm"}}}
		}
		if arm.Fallback {
			fallback = arm
			continue
		}
		patternType, err := s.resolveType(arm.PatternType, scope)
		if err != nil {
			return nil, err
		}
		if !ty.Equals(patternType) {
			return nil, diagnostics.Diag{Title: "Wrong enum in pattern", Message: "pattern belongs to a different enum", Span: arm.VariantName.Span, Help: "Use variants of the matched enum.", Related: []diagnostics.Related{{Span: d.Name.Span, Message: "matched enum declared here"}}}
		}
		for _, v := range d.Variants {
			if v.Name.Name() == arm.VariantName.Name() {
				arm.Variant = v
				break
			}
		}
		if arm.Variant == nil {
			return nil, diagnostics.Diag{Title: "Unknown enum variant", Message: "variant is not declared by the matched enum", Span: arm.VariantName.Span, Help: "Use a variant declared in this enum.", Related: []diagnostics.Related{{Span: d.Name.Span, Message: "enum declared here"}}}
		}
		if previous := seen[arm.Variant]; previous != nil {
			return nil, diagnostics.Diag{Title: "Duplicate match arm", Message: "variant has already been matched", Span: arm.VariantName.Span, Help: "Match each variant once.", Related: []diagnostics.Related{{Span: previous.VariantName.Span, Message: "previous arm"}}}
		}
		seen[arm.Variant] = arm
		fields := map[string]*ast.PatternField{}
		for _, binding := range arm.Fields {
			if previous := fields[binding.Name.Name()]; previous != nil {
				return nil, diagnostics.Diag{Title: "Duplicate pattern field", Message: "payload field is bound twice", Span: binding.Name.Span, Help: "Bind each payload field at most once.", Related: []diagnostics.Related{{Span: previous.Name.Span, Message: "previous binding"}}}
			}
			fields[binding.Name.Name()] = binding
			for _, f := range arm.Variant.Fields {
				if f.Name.Name() == binding.Name.Name() {
					binding.Binding.Type = f.Type
					break
				}
			}
			if binding.Binding.Type == nil {
				return nil, diagnostics.Diag{Title: "Unknown payload field", Message: "pattern names an undeclared payload field", Span: binding.Name.Span, Help: "Use a field declared in this variant.", Related: []diagnostics.Related{{Span: arm.Variant.Name.Span, Message: "variant declared here"}}}
			}
			node := &ast.Node{Kind: ast.KIND_PARAM, Node: binding.Binding, Span: binding.Binding.Name.Span}
			if previous := arm.Scope.Nodes[binding.Binding.Name.Name()]; previous != nil {
				return nil, diagnostics.Diag{Title: "Duplicate pattern binding", Message: "two fields use the same binding name", Span: binding.Binding.Name.Span, Help: "Give the bindings distinct names.", Related: []diagnostics.Related{{Span: previous.DeclarationSpan(), Message: "previous binding"}}}
			}
			if err := arm.Scope.Insert(binding.Binding.Name.Name(), node); err != nil {
				return nil, err
			}
		}
		if !arm.OmitRest && len(fields) != len(arm.Variant.Fields) {
			return nil, diagnostics.Diag{Title: "Incomplete payload pattern", Message: "pattern omits payload fields without ..", Span: arm.VariantName.Span, Help: "Bind every payload field or add .. to omit the remaining fields.", Related: []diagnostics.Related{{Span: arm.Variant.Name.Span, Message: "variant declared here"}}}
		}
	}
	if fallback == nil && len(seen) != len(d.Variants) {
		var missing []string
		var related []diagnostics.Related
		for _, v := range d.Variants {
			if seen[v] == nil {
				missing = append(missing, v.Name.Name())
				related = append(related, diagnostics.Related{Span: v.Name.Span, Message: "unmatched variant"})
			}
		}
		return nil, diagnostics.Diag{Title: "Non-exhaustive match", Message: "missing variants: " + strings.Join(missing, ", "), Span: m.Value.Span, Help: "Add arms for the missing variants or a final _ fallback.", Related: related}
	}
	if len(m.Arms) == 0 {
		return nil, diagnostics.Diag{Title: "Empty match", Message: "a value-producing match needs a result arm", Span: m.Value.Span, Help: "Add a variant arm or a fallback that produces a value."}
	}
	var sites []yieldSite
	for _, arm := range m.Arms {
		if arm.Block != nil {
			ctx := &yieldContext{}
			previous := s.currentYield
			s.currentYield = ctx
			err := s.checkBlockStmt(arm.Block, arm.Scope, arm.Scope, s.currentReturnType, false, false)
			s.currentYield = previous
			if err != nil {
				return nil, err
			}
			if !blockYields(arm.Block) {
				return nil, diagnostics.Diag{Title: "Missing match yield", Message: "match arm block can finish without producing a value", Span: arm.VariantName.Span, Help: "Use yield expression on every path that does not return from the function."}
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
		return nil, diagnostics.Diag{Title: "Missing match result", Message: "match has no value-producing arm", Span: m.Value.Span, Help: "Produce a value with an expression arm or yield."}
	}
	if !valueType(result) {
		return nil, diagnostics.Diag{Title: "Invalid match result", Message: "match arms must produce a value", Span: m.Value.Span, Help: "Return compatible non-void values from every arm."}
	}
	for _, site := range sites {
		if _, err := s.inferExprTypeWithContext(site.expr, result, site.scope, site.scope, false, false); err != nil {
			return nil, err
		}
	}
	m.Type = result
	return result, nil
}

func blockYields(block *ast.BlockStmt) bool {
	for _, stmt := range block.Statements {
		switch stmt.Kind {
		case ast.KIND_YIELD_STMT, ast.KIND_RETURN_STMT:
			return true
		case ast.KIND_BLOCK_STMT:
			if blockYields(stmt.Node.(*ast.BlockStmt)) {
				return true
			}
		case ast.KIND_COND_STMT:
			cond := stmt.Node.(*ast.CondStmt)
			if cond.ElseStmt == nil || !blockYields(cond.IfStmt.Block) || !blockYields(cond.ElseStmt.Block) {
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
