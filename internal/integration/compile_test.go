package integration

import (
	"strings"
	"testing"

	"github.com/HicaroD/Telia/internal/itest"
)

func TestCompileHelloWorld(t *testing.T) {
	output, diags := itest.CompileFile("testdata/hello_world.t")
	if len(diags.Diags) > 0 {
		t.Fatalf("unexpected errors: %v", diags.Diags)
	}
	expected := "Hello, world!\n"
	if output != expected {
		t.Errorf("expected %q, got %q", expected, output)
	}
}

func TestCompileFibonacci(t *testing.T) {
	output, diags := itest.CompileFile("testdata/fib.t")
	if len(diags.Diags) > 0 {
		t.Fatalf("unexpected errors: %v", diags.Diags)
	}
	expected := "55\n"
	if output != expected {
		t.Errorf("expected %q, got %q", expected, output)
	}
}

func TestCompileCalculator(t *testing.T) {
	output, diags := itest.CompileFile("testdata/calculator.t")
	if len(diags.Diags) > 0 {
		t.Fatalf("unexpected errors: %v", diags.Diags)
	}
	expected := "8\n6\n42\n"
	if output != expected {
		t.Errorf("expected %q, got %q", expected, output)
	}
}

func TestUndefinedVariable(t *testing.T) {
	_, diags := itest.CompileFile("testdata/errors/undefined_var.t")
	if len(diags.Diags) == 0 {
		t.Fatalf("expected errors, got none")
	}
	found := false
	for _, diag := range diags.Diags {
		if strings.Contains(diag.Message, "y") {
			found = true
			break
		}
	}
	if !found {
		t.Errorf("expected error about undefined variable 'y', got: %v", diags.Diags)
	}
}

func TestTypeMismatch(t *testing.T) {
	_, diags := itest.CompileFile("testdata/errors/type_mismatch.t")
	if len(diags.Diags) == 0 {
		t.Fatalf("expected errors, got none")
	}
	found := false
	for _, diag := range diags.Diags {
		if strings.Contains(diag.Message, "string") || strings.Contains(diag.Message, "int") {
			found = true
			break
		}
	}
	if !found {
		t.Errorf("expected error about type mismatch, got: %v", diags.Diags)
	}
}

func TestBadSyntax(t *testing.T) {
	_, diags := itest.CompileFile("testdata/errors/bad_syntax.t")
	if len(diags.Diags) == 0 {
		t.Fatalf("expected errors, got none")
	}
}
