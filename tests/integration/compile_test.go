package integration

import (
	"strings"
	"testing"

	"github.com/HicaroD/Telia/tests/compiler"
)

func TestCompileHelloWorld(t *testing.T) {
	output, diags := compiler.CompileFile("testdata/hello_world.t")
	if len(diags.Diags) > 0 {
		t.Fatalf("unexpected errors: %v", diags.Diags)
	}
	expected := "Hello, world!\n"
	if output != expected {
		t.Errorf("expected %q, got %q", expected, output)
	}
}

func TestCompileFibonacci(t *testing.T) {
	output, diags := compiler.CompileFile("testdata/fib.t")
	if len(diags.Diags) > 0 {
		t.Fatalf("unexpected errors: %v", diags.Diags)
	}
	expected := "55\n"
	if output != expected {
		t.Errorf("expected %q, got %q", expected, output)
	}
}

func TestCompileCalculator(t *testing.T) {
	output, diags := compiler.CompileFile("testdata/calculator.t")
	if len(diags.Diags) > 0 {
		t.Fatalf("unexpected errors: %v", diags.Diags)
	}
	expected := "8\n6\n42\n"
	if output != expected {
		t.Errorf("expected %q, got %q", expected, output)
	}
}

func TestUndefinedVariable(t *testing.T) {
	_, diags := compiler.CompileFile("testdata/errors/undefined_var.t")
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
	_, diags := compiler.CompileFile("testdata/errors/type_mismatch.t")
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
	_, diags := compiler.CompileFile("testdata/errors/bad_syntax.t")
	if len(diags.Diags) == 0 {
		t.Fatalf("expected errors, got none")
	}
}

func TestForLoop(t *testing.T) {
	output, diags := compiler.CompileFile("testdata/for_loop.t")
	if len(diags.Diags) > 0 {
		t.Fatalf("unexpected errors: %v", diags.Diags)
	}
	expected := "10\n"
	if output != expected {
		t.Errorf("expected %q, got %q", expected, output)
	}
}

func TestWhileLoop(t *testing.T) {
	output, diags := compiler.CompileFile("testdata/while_loop.t")
	if len(diags.Diags) > 0 {
		t.Fatalf("unexpected errors: %v", diags.Diags)
	}
	expected := "1\n2\n3\n"
	if output != expected {
		t.Errorf("expected %q, got %q", expected, output)
	}
}

func TestCondStatement(t *testing.T) {
	output, diags := compiler.CompileFile("testdata/cond.t")
	if len(diags.Diags) > 0 {
		t.Fatalf("unexpected errors: %v", diags.Diags)
	}
	expected := "-1\n0\n1\n"
	if output != expected {
		t.Errorf("expected %q, got %q", expected, output)
	}
}

func TestStructLiteralAndFieldAccess(t *testing.T) {
	output, diags := compiler.CompileFile("testdata/struct_field_access.t")
	if len(diags.Diags) > 0 {
		t.Fatalf("unexpected errors: %v", diags.Diags)
	}
	// pt.x=3, pt.y=7, sum_fields(pt)=10, scale_x(&pt, 2)=6
	expected := "3\n7\n10\n6\n"
	if output != expected {
		t.Errorf("expected %q, got %q", expected, output)
	}
}

func TestMultipleReturnValues(t *testing.T) {
	output, diags := compiler.CompileFile("testdata/multi_ret.t")
	if len(diags.Diags) > 0 {
		t.Fatalf("unexpected errors: %v", diags.Diags)
	}
	// get(3) → (1, 5); literal tuple (10, 20) with heterogeneous types (i32, i64)
	expected := "1\n5\n10\n20\n"
	if output != expected {
		t.Errorf("expected %q, got %q", expected, output)
	}
}

func TestFloatArithmetic(t *testing.T) {
	output, diags := compiler.CompileFile("testdata/floats.t")
	if len(diags.Diags) > 0 {
		t.Fatalf("unexpected errors: %v", diags.Diags)
	}
	// f32/f64 literals, arithmetic (+, -, *, /), and all 6 comparison operators
	expected := "2.0\n1.5\n3.0\n5.0\n2.25\nlt\nle\ngt\nge\neq\nne\n"
	if output != expected {
		t.Errorf("expected %q, got %q", expected, output)
	}
}
