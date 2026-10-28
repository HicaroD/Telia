package c

import "github.com/HicaroD/Telia/internal/ast"

func callingConvention(attrs ast.Attributes) string {
	if attrs.CallingConventionMacro != "" {
		return attrs.CallingConventionMacro + " "
	}
	return ""
}
