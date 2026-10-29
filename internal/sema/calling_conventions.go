package sema

import (
	"github.com/HicaroD/Telia/internal/ast"
	"github.com/HicaroD/Telia/internal/diagnostics"
	"github.com/HicaroD/Telia/internal/source"
)

// Native function values are not foreign callbacks, even when the backend uses
// the same machine calling convention for both.
func checkForeignCallbacks(ty *ast.ExprType, span source.Span) error {
	switch ty.Kind {
	case ast.EXPR_TYPE_FUNCTION:
		f := ty.T.(*ast.FunctionType)
		if !f.HasForeignABI() {
			return diagnostics.Diag{
				Title:   "Foreign callback requires an ABI",
				Message: "a native function type cannot be used as a foreign callback",
				Span:    span,
				Help:    "Annotate the callback alias and its implementation with @[cc=\"c\"] or the matching cc_macro.",
			}
		}
		for _, param := range f.Params.Fields {
			if err := checkForeignCallbacks(param.Type, span); err != nil {
				return err
			}
		}
		return checkForeignCallbacks(f.Result, span)
	case ast.EXPR_TYPE_POINTER:
		return checkForeignCallbacks(ty.T.(*ast.PointerType).Type, span)
	case ast.EXPR_TYPE_ARRAY, ast.EXPR_TYPE_SLICE:
		return checkForeignCallbacks(ty.T.(*ast.SequenceType).Element, span)
	}
	return nil
}

func validateCallingConvention(cc, macro string) error {
	if cc != "" && macro != "" {
		return diagnostics.Diag{
			Title:   "Conflicting calling conventions",
			Message: "choose a calling convention or a header macro, not both",
			Help:    "Use cc/default_cc or cc_macro/default_cc_macro.",
		}
	}
	if cc != "" && cc != "c" {
		return diagnostics.Diag{
			Title:   "Unsupported calling convention",
			Message: "this backend currently supports only the C calling convention and header convention macros",
			Help:    "Use cc=\"c\" or the library's cc_macro.",
		}
	}
	for i, c := range macro {
		if !(c == '_' || c >= 'a' && c <= 'z' || c >= 'A' && c <= 'Z' || i > 0 && c >= '0' && c <= '9') {
			return diagnostics.Diag{
				Title:   "Invalid calling convention macro",
				Message: "a calling convention macro must be one C identifier",
				Help:    "Use the header's macro name, not a C expression or declaration.",
			}
		}
	}
	return nil
}
