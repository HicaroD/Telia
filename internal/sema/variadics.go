package sema

import (
	"github.com/HicaroD/Telia/internal/ast"
	"github.com/HicaroD/Telia/internal/diagnostics"
	"github.com/HicaroD/Telia/internal/lexer/token"
	"github.com/HicaroD/Telia/internal/source"
)

func rejectCVariadicValue(param *ast.Param, span source.Span) error {
	if param.Variadic && param.Attributes != nil && param.Attributes.C {
		return diagnostics.Diag{
			Title:   "C variadic marker is not a value",
			Message: "args @c ... names a native variadic argument pack, not ordinary storage",
			Span:    span,
			Help:    "Use c::va_start(args) from std::ffi::c to initialize local owning cursor storage.",
			Related: []diagnostics.Related{
				{Span: param.Name.Span, Message: "C variadic marker declared here"},
			},
		}
	}
	return nil
}

func checkCVariadicParams(params *ast.Params) error {
	if params.IsVariadic && params.Len == 0 && len(params.Fields) != 0 &&
		params.Fields[0].Attributes != nil &&
		params.Fields[0].Attributes.C {
		return diagnostics.Diag{
			Title:   "C variadic fixed parameter",
			Message: "C variadic signatures require at least one fixed parameter",
			Span:    params.Fields[0].Name.Span,
			Help:    "Declare a fixed marker or format parameter before args @c ....",
		}
	}
	return nil
}

func cVariadicType(ty *ast.ExprType) bool {
	switch ty.Kind {
	case ast.EXPR_TYPE_FUNCTION:
		return ty.T.(*ast.FunctionType).HasForeignABI()
	case ast.EXPR_TYPE_POINTER:
		return checkForeignCallbacks(ty, source.Span{}) == nil
	case ast.EXPR_TYPE_ENUM:
		return !ty.T.(*ast.EnumType).Decl.HasPayload
	case ast.EXPR_TYPE_BASIC:
		kind := ty.T.(*ast.BasicType).Kind
		return kind.IsNumeric() || kind == token.BOOL_TYPE || kind == token.CSTRING_TYPE ||
			kind == token.RAWPTR_TYPE
	}
	return false
}
