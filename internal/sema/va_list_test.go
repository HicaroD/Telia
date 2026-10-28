package sema_test

import "testing"

func TestVaListIntrinsicsAreNotOrdinaryCallbackValues(t *testing.T) {
	c := checkDiagnostics(t, "package main\ntype VaList = __c_va_list_owner\nfn main() {\ncallback := va_end\n}\n")
	if len(c.Diags) != 1 || c.Diags[0].Title != "Intrinsic is not a function value" || !c.Diags[0].Span.Valid() || c.Diags[0].Help == "" {
		t.Fatalf("want direct-invocation diagnostic, got %v", c.Diags)
	}
}

func TestVaStartRejectsUnverifiedCallingConventionMacros(t *testing.T) {
	c := checkDiagnostics(t, "package main\ntype VaList = __c_va_list_owner\n@[cc_macro=\"FOREIGN_CC\"]\nfn inspect(last i32, args @c ...) {\nowner := va_start(args)\ndefer va_end(&owner)\n}\nfn main() {}\n")
	if len(c.Diags) != 1 || c.Diags[0].Title != "Unsupported va_start convention" || !c.Diags[0].Span.Valid() || c.Diags[0].Help == "" {
		t.Fatalf("want native-convention diagnostic, got %v", c.Diags)
	}
}

func TestVaStartRejectsLastFixedParametersChangedByPromotion(t *testing.T) {
	c := checkDiagnostics(t, "package main\ntype VaList = __c_va_list_owner\nfn inspect(last f32, args @c ...) {\nowner := va_start(args)\ndefer va_end(&owner)\n}\nfn main() {}\n")
	if len(c.Diags) != 1 || c.Diags[0].Title != "Invalid va_start parameter" || !c.Diags[0].Span.Valid() || c.Diags[0].Help == "" || len(c.Diags[0].Related) != 1 {
		t.Fatalf("want stable last-parameter diagnostic, got %v", c.Diags)
	}
}

func TestVaListOwnerAddressCannotEscapeIntrinsicOperations(t *testing.T) {
	c := checkDiagnostics(t, "package main\ntype VaList = __c_va_list_owner\ntype Args = __c_va_list\nfn inspect(args Args) {\nowner := va_copy(args)\ndefer va_end(&owner)\npointer := &owner\n}\nfn main() {}\n")
	if len(c.Diags) != 1 || c.Diags[0].Title != "va_list owner escape" || !c.Diags[0].Span.Valid() || c.Diags[0].Help == "" {
		t.Fatalf("want owner escape diagnostic, got %v", c.Diags)
	}
}

func TestVaListOwnersCannotBeCopiedByAssignment(t *testing.T) {
	for _, declaration := range []string{"copy := owner", "copy VaList := owner"} {
		t.Run(declaration, func(t *testing.T) {
			c := checkDiagnostics(t, "package main\ntype VaList = __c_va_list_owner\ntype Args = __c_va_list\nfn inspect(args Args) {\nowner := va_copy(args)\ndefer va_end(&owner)\n"+declaration+"\n}\nfn main() {}\n")
			if len(c.Diags) != 1 || c.Diags[0].Title != "Noncopyable va_list owner" || !c.Diags[0].Span.Valid() || c.Diags[0].Help == "" {
				t.Fatalf("want owner copy diagnostic, got %v", c.Diags)
			}
		})
	}
}

func TestVaArgRejectsTypesChangedByDefaultPromotions(t *testing.T) {
	for _, ty := range []string{"i8", "u16", "bool", "f32"} {
		t.Run(ty, func(t *testing.T) {
			c := checkDiagnostics(t, "package main\ntype VaList = __c_va_list_owner\ntype Args = __c_va_list\nfn inspect(args Args) {\nowner := va_copy(args)\ndefer va_end(&owner)\nvalue := va_arg<"+ty+">(&owner)\n}\nfn main() {}\n")
			if len(c.Diags) != 1 || c.Diags[0].Title != "Invalid va_arg type" || !c.Diags[0].Span.Valid() || c.Diags[0].Help == "" || len(c.Diags[0].Related) != 1 {
				t.Fatalf("want promoted-type diagnostic, got %v", c.Diags)
			}
		})
	}
}

func TestBorrowedVaListCannotEscapeThroughItsAddress(t *testing.T) {
	c := checkDiagnostics(t, "package main\ntype Args = __c_va_list\nfn inspect(args Args) {\npointer := &args\n}\nfn main() {}\n")
	if len(c.Diags) != 1 || c.Diags[0].Title != "Borrowed va_list escape" || !c.Diags[0].Span.Valid() || c.Diags[0].Help == "" {
		t.Fatalf("want borrowed va_list address diagnostic, got %v", c.Diags)
	}
}

func TestBorrowedVaListCannotBeReturned(t *testing.T) {
	c := checkDiagnostics(t, "package main\ntype Args = __c_va_list\nfn escape(args Args) Args {\nreturn args\n}\nfn main() {}\n")
	if len(c.Diags) != 1 || c.Diags[0].Title != "Borrowed va_list escape" || !c.Diags[0].Span.Valid() || c.Diags[0].Help == "" {
		t.Fatalf("want borrowed va_list return diagnostic, got %v", c.Diags)
	}
}

func TestBorrowedVaListCannotBeCopiedWithExplicitType(t *testing.T) {
	c := checkDiagnostics(t, "package main\ntype Args = __c_va_list\nfn inspect(args Args) {\ncopy Args := args\n}\nfn main() {}\n")
	if len(c.Diags) != 1 || c.Diags[0].Title != "Borrowed va_list copy" || !c.Diags[0].Span.Valid() || c.Diags[0].Help == "" {
		t.Fatalf("want borrowed va_list diagnostic, got %v", c.Diags)
	}
}

func TestBorrowedVaListCannotBeCopied(t *testing.T) {
	c := checkDiagnostics(t, "package main\ntype Args = __c_va_list\nfn inspect(args Args) {\ncopy := args\n}\nfn main() {}\n")
	if len(c.Diags) != 1 || c.Diags[0].Title != "Borrowed va_list copy" || !c.Diags[0].Span.Valid() || c.Diags[0].Help == "" {
		t.Fatalf("want borrowed va_list diagnostic, got %v", c.Diags)
	}
}
