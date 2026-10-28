package integration

import (
	"testing"

	"github.com/HicaroD/Telia/tests/compiler"
)

func TestGenericStructsHaveIndependentConcreteLayouts(t *testing.T) {
	output, diags := compiler.CompileFile("testdata/generic_structs.t")
	if len(diags.Diags) != 0 {
		t.Fatalf("generic structs should compile: %v", diags.Diags)
	}
	if want := "42\nhello\n7\n"; output != want {
		t.Fatalf("want %q, got %q", want, output)
	}
}

func TestGenericTupleResultsKeepConcreteElementIdentity(t *testing.T) {
	output, diags := compiler.CompileFile("testdata/generic_tuples.t")
	if len(diags.Diags) != 0 {
		t.Fatalf("generic tuple results should compile: %v", diags.Diags)
	}
	if want := "42\nhello\n"; output != want {
		t.Fatalf("want %q, got %q", want, output)
	}
}

func TestGenericBodiesPreservePointersControlFlowAndCleanup(t *testing.T) {
	output, diags := compiler.CompileFile("testdata/generic_control.t")
	if len(diags.Diags) != 0 {
		t.Fatalf("generic body cloning should preserve behavior: %v", diags.Diags)
	}
	if want := "cleanup\ncleanup\n42\nhello\n"; output != want {
		t.Fatalf("want %q, got %q", want, output)
	}
}

func TestGenericTypeParametersCanConstructConcreteStructs(t *testing.T) {
	output, diags := compiler.CompileFile("testdata/generic_zero.t")
	if len(diags.Diags) != 0 {
		t.Fatalf("type-parameter literal should specialize: %v", diags.Diags)
	}
	if output != "0\n" {
		t.Fatalf("want zero-initialized struct, got %q", output)
	}
}

func TestGenericBodiesWithNativeSanitizers(t *testing.T) {
	path, diags := compiler.CompileOnlyWithCFlags(
		"testdata/generic_control.t",
		"-fsanitize=undefined",
		"-fno-sanitize-recover=all",
	)
	if len(diags.Diags) != 0 {
		t.Fatalf("sanitized generic compilation: %v", diags.Diags)
	}
	defer removeTestBinary(t, path)
	output, err := compiler.RunBinary(path)
	if err != nil {
		t.Fatalf("sanitized execution: %v: %s", err, output)
	}
	if output != "cleanup\ncleanup\n42\nhello\n" {
		t.Fatalf("unexpected sanitized output: %q", output)
	}
}

func TestDocumentedGenericsExample(t *testing.T) {
	output, diags := compiler.CompileFile("../../examples/generics.t")
	if len(diags.Diags) != 0 {
		t.Fatalf("documented example should compile: %v", diags.Diags)
	}
	if output != "count\n42\n" {
		t.Fatalf("unexpected example output: %q", output)
	}
}

func TestImportedGenericsUseDeclarationContextAndCallerTypeArguments(t *testing.T) {
	output, diags := compiler.CompilePackage("testdata/generic_imports")
	if len(diags.Diags) != 0 {
		t.Fatalf("imported generics should compile: %v", diags.Diags)
	}
	if want := "caller\nlibrary\nother\n9\n9\n"; output != want {
		t.Fatalf("want %q, got %q", want, output)
	}
}

func TestGenericFunctionsComposeAndReuseRecursiveInstances(t *testing.T) {
	output, diags := compiler.CompileFile("testdata/generic_composition.t")
	if len(diags.Diags) != 0 {
		t.Fatalf("generic composition should compile: %v", diags.Diags)
	}
	if want := "42\nhello\n5\n"; output != want {
		t.Fatalf("want %q, got %q", want, output)
	}
}

func TestGenericFunctionsSpecializeWithoutSharingState(t *testing.T) {
	output, diags := compiler.CompileFile("testdata/generic_functions.t")
	if len(diags.Diags) != 0 {
		t.Fatalf("generic functions should compile: %v", diags.Diags)
	}
	if want := "42\nhello\n7\n"; output != want {
		t.Fatalf("want %q, got %q", want, output)
	}
}

func TestGenericStructsComposeAndAllowPointerRecursion(t *testing.T) {
	output, diags := compiler.CompileFile("testdata/generic_layouts.t")
	if len(diags.Diags) != 0 {
		t.Fatalf("nested and recursive generic structs should compile: %v", diags.Diags)
	}
	if want := "nested\n9\n"; output != want {
		t.Fatalf("want %q, got %q", want, output)
	}
}
