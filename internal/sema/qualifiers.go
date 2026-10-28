package sema

import (
	"github.com/HicaroD/Telia/internal/ast"
	"github.com/HicaroD/Telia/internal/diagnostics"
	"github.com/HicaroD/Telia/internal/lexer/token"
	"github.com/HicaroD/Telia/internal/source"
)

func constMutation(span source.Span) error {
	return diagnostics.Diag{Title: "Const pointee mutation", Message: "cannot modify storage protected by a const pointer", Span: span, Help: "Use a mutable pointer only when the API permits writing to this storage."}
}

func (s *sema) readOnlyStorage(node *ast.Node, scope *ast.Scope) bool {
	switch n := node.Node.(type) {
	case *ast.FieldAccess:
		for n.Right.Kind == ast.KIND_FIELD_ACCESS {
			n = n.Right.Node.(*ast.FieldAccess)
		}
		return n.ReadOnly
	case *ast.IndexExpr:
		return s.readOnlyStorage(n.Value, scope)
	case *ast.DerefPointerExpr:
		ty, _, err := s.inferExprTypeWithoutContext(n.Expr, scope, scope, false, false)
		return err == nil && ty.IsPointer() && ty.T.(*ast.PointerType).Const
	}
	return false
}

// Only the immediate pointee may gain const. Changing nested qualifiers would
// permit writing an incompatible pointer through a mutable output slot.
func pointerAssignable(actual, expected *ast.ExprType) bool {
	if !actual.IsPointer() || !expected.IsPointer() {
		return false
	}
	a, e := actual.T.(*ast.PointerType), expected.T.(*ast.PointerType)
	return (!a.Const || e.Const) && a.Type.Equals(e.Type)
}

func typeAssignable(actual, expected *ast.ExprType) bool {
	if actual.Equals(expected) || pointerAssignable(actual, expected) {
		return true
	}
	if actual.Kind == ast.EXPR_TYPE_BASIC && expected.Kind == ast.EXPR_TYPE_BASIC {
		a, e := actual.T.(*ast.BasicType), expected.T.(*ast.BasicType)
		if actual.IsNumeric() && a.Kind == e.Kind {
			return true
		}
		return a.Kind == token.CSTRING_TYPE && e.Kind == token.CSTRING_TYPE && !a.Const && e.Const
	}
	return false
}
