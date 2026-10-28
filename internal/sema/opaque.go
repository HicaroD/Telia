package sema

import (
	"github.com/HicaroD/Telia/internal/ast"
	"github.com/HicaroD/Telia/internal/diagnostics"
	"github.com/HicaroD/Telia/internal/source"
)

func rejectOpaqueValue(ty *ast.ExprType, span source.Span) error {
	if ty.Kind == ast.EXPR_TYPE_OPAQUE {
		d := ty.T.(*ast.OpaqueType).Decl
		return diagnostics.Diag{Title: "Incomplete opaque value", Message: "opaque types cannot be constructed, sized, or used by value", Span: span, Help: "Use a typed pointer to the opaque type; its definition belongs to the C header.", Related: []diagnostics.Related{{Span: d.Name.Span, Message: "opaque type declared here"}}}
	}
	if ty.Kind == ast.EXPR_TYPE_TUPLE {
		for _, t := range ty.T.(*ast.TupleType).Types {
			if err := rejectOpaqueValue(t, span); err != nil {
				return err
			}
		}
	}
	return nil
}
