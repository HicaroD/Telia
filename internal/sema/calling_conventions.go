package sema

import (
	"github.com/HicaroD/Telia/internal/diagnostics"
)

func validateCallingConvention(cc, macro string) error {
	if cc != "" && macro != "" {
		return diagnostics.Diag{Title: "Conflicting calling conventions", Message: "choose a calling convention or a header macro, not both", Help: "Use cc/default_cc or cc_macro/default_cc_macro."}
	}
	if cc != "" && cc != "c" {
		return diagnostics.Diag{Title: "Unsupported calling convention", Message: "this backend currently supports only the C calling convention and header convention macros", Help: "Use cc=\"c\" or the library's cc_macro."}
	}
	for i, c := range macro {
		if !(c == '_' || c >= 'a' && c <= 'z' || c >= 'A' && c <= 'Z' || i > 0 && c >= '0' && c <= '9') {
			return diagnostics.Diag{Title: "Invalid calling convention macro", Message: "a calling convention macro must be one C identifier", Help: "Use the header's macro name, not a C expression or declaration."}
		}
	}
	return nil
}
