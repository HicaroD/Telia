package sema_test

import "testing"

func TestNativeFunctionCannotInitializeCCallback(t *testing.T) {
	c := checkDiagnostics(t, `package main
@[cc="c"]
type Callback = fn(value i32) i32
fn native(value i32) i32 {
return value
}
fn main() {
callback Callback := native
}
`)
	if len(c.Diags) != 1 || c.Diags[0].Title != "Callback signature mismatch" ||
		!c.Diags[0].Span.Valid() {
		t.Fatalf("want a located native/foreign callback mismatch, got %+v", c.Diags)
	}
}

func TestForeignSignatureRequiresExplicitCallbackABI(t *testing.T) {
	c := checkDiagnostics(t, `package main
type Callback = fn(value i32) i32
extern {
fn register(callback *Callback)
}
fn main() {}
`)
	if len(c.Diags) != 1 || c.Diags[0].Title != "Foreign callback requires an ABI" ||
		!c.Diags[0].Span.Valid() || c.Diags[0].Help == "" {
		t.Fatalf("want a located explicit-ABI diagnostic with migration help, got %+v", c.Diags)
	}
}

func TestExternFunctionValueIsNotANativeFunctionValue(t *testing.T) {
	c := checkDiagnostics(t, `package main
type Native = fn(value i32) i32
extern {
fn foreign(value i32) i32
}
fn main() {
callback Native := foreign
}
`)
	if len(c.Diags) != 1 || c.Diags[0].Title != "Callback signature mismatch" {
		t.Fatalf("want extern function to retain its foreign ABI, got %+v", c.Diags)
	}
}

func TestForeignRecordRequiresExplicitCallbackABI(t *testing.T) {
	c := checkDiagnostics(t, `package main
type Callback = fn(value i32) i32
@[header="library.h"]
extern {
struct Callbacks {
callback Callback
}
}
fn main() {}
`)
	if len(c.Diags) != 1 || c.Diags[0].Title != "Foreign callback requires an ABI" {
		t.Fatalf("want native callback field rejected at the foreign boundary, got %+v", c.Diags)
	}
}

func TestCVariadicsRejectNativeFunctionValues(t *testing.T) {
	c := checkDiagnostics(t, `package main
extern {
fn inspect(marker i32, args @c ...)
}
fn native(value i32) i32 {
return value
}
fn main() {
inspect(0, native)
}
`)
	if len(c.Diags) != 1 || c.Diags[0].Title != "Unsupported C variadic argument" {
		t.Fatalf("want native callback rejected in C variadics, got %+v", c.Diags)
	}
}

func TestCVariadicsRejectNativeCallbackOutputSlots(t *testing.T) {
	c := checkDiagnostics(t, `package main
extern {
fn inspect(marker i32, args @c ...)
}
fn native(value i32) i32 {
return value
}
fn main() {
callback := native
inspect(0, &callback)
}
`)
	if len(c.Diags) != 1 || c.Diags[0].Title != "Unsupported C variadic argument" {
		t.Fatalf("want native callback storage rejected in C variadics, got %+v", c.Diags)
	}
}

func TestForeignCallableSignaturesRejectNestedNativeCallbacks(t *testing.T) {
	for _, declaration := range []string{
		"@[cc=\"c\"]\ntype Register = fn(callback Native)\nfn main() {\ncallback Register := nil\n}\n",
		"@[cc=\"c\"]\nfn register(callback Native) {}\nfn main() {}\n",
	} {
		t.Run(declaration, func(t *testing.T) {
			c := checkDiagnostics(t, "package main\ntype Native = fn(value i32) i32\n"+
				declaration)
			if len(c.Diags) != 1 || c.Diags[0].Title != "Foreign callback requires an ABI" {
				t.Fatalf("want nested native callback rejected, got %+v", c.Diags)
			}
		})
	}
}

func TestVaArgCannotConstructANativeFunctionValue(t *testing.T) {
	c := checkDiagnostics(t, `package main
type VaList = __c_va_list_owner
type Native = fn(value i32) i32
fn inspect(marker i32, args @c ...) {
owner := va_start(args)
defer va_end(&owner)
callback := va_arg<Native>(&owner)
}
fn main() {}
`)
	if len(c.Diags) != 1 || c.Diags[0].Title != "Invalid va_arg type" {
		t.Fatalf("want C storage unable to produce a native callback, got %+v", c.Diags)
	}
}
