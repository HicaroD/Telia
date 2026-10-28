package sema

import (
	"github.com/HicaroD/Telia/internal/ast"
	"github.com/HicaroD/Telia/internal/diagnostics"
	"github.com/HicaroD/Telia/internal/source"
)

func checkForeignPointerDeclarator(ty *ast.ExprType, span source.Span) error {
	if ty.Kind != ast.EXPR_TYPE_POINTER {
		return nil
	}
	inner := ty.T.(*ast.PointerType).Type
	if inner.Kind == ast.EXPR_TYPE_ARRAY || inner.Kind == ast.EXPR_TYPE_SLICE {
		return diagnostics.Diag{
			Title:   "Unsupported foreign declarator",
			Message: "a pointer to a Telia wrapper array or slice is not a pointer to a C array",
			Span:    span,
			Help:    "Use an explicit element pointer, such as *u8, or a header-compiled C adapter for pointer-to-array declarators.",
		}
	}
	return checkForeignPointerDeclarator(inner, span)
}
