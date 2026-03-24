package integration

import (
	"bytes"
	"os/exec"
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

func TestDefer(t *testing.T) {
	output, diags := compiler.CompileFile("testdata/defer.t")
	if len(diags.Diags) > 0 {
		t.Fatalf("unexpected errors: %v", diags.Diags)
	}
	// lifo(): LIFO order — third/second/first
	// no_defer(): no output
	// nested(1): inner defer runs at if-block fall-through; outer before return
	expected := "third\nsecond\nfirst\ninner\nouter\n"
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

func TestPointerArithmetic(t *testing.T) {
	output, diags := compiler.CompileFile("testdata/pointer_arith.t")
	if len(diags.Diags) > 0 {
		t.Fatalf("unexpected errors: %v", diags.Diags)
	}
	// *p = 100 mutates x; both *p and x print 100
	expected := "100\n100\n"
	if output != expected {
		t.Errorf("expected %q, got %q", expected, output)
	}
}

func TestCompilePackage(t *testing.T) {
	output, diags := compiler.CompilePackage("testdata/pkg_smoke")
	if len(diags.Diags) > 0 {
		t.Fatalf("unexpected errors: %v", diags.Diags)
	}
	// greet::hello() prints first, then main's io::println
	expected := "Hello from greet package!\nHello from main package!\n"
	if output != expected {
		t.Errorf("expected %q, got %q", expected, output)
	}
}

func TestErrorReturn(t *testing.T) {
	output, diags := compiler.CompileFile("testdata/error_return.t")
	if len(diags.Diags) > 0 {
		t.Fatalf("unexpected errors: %v", diags.Diags)
	}
	expected := "5\nok\n0\ndivision by zero\n"
	if output != expected {
		t.Errorf("expected %q, got %q", expected, output)
	}
}

func TestErrorBasic(t *testing.T) {
	output, diags := compiler.CompileFile("testdata/error_basic.t")
	if len(diags.Diags) > 0 {
		t.Fatalf("unexpected errors: %v", diags.Diags)
	}
	expected := "oops\n"
	if output != expected {
		t.Errorf("expected %q, got %q", expected, output)
	}
}

func TestErrorFailTuple(t *testing.T) {
	output, diags := compiler.CompileFile("testdata/error_fail_tuple.t")
	if len(diags.Diags) > 0 {
		t.Fatalf("unexpected errors: %v", diags.Diags)
	}
	expected := "99\n"
	if output != expected {
		t.Errorf("expected %q, got %q", expected, output)
	}
}

func TestErrorFailPanic(t *testing.T) {
	exePath, diags := compiler.CompileOnly("testdata/error_fail_panic.t")
	if len(diags.Diags) > 0 {
		t.Fatalf("unexpected compile errors: %v", diags.Diags)
	}

	stderr, err := compiler.RunBinary(exePath)
	if err == nil {
		t.Fatal("expected non-zero exit code, got success")
	}
	exitErr, ok := err.(*exec.ExitError)
	if !ok {
		t.Fatalf("expected *exec.ExitError, got: %v", err)
	}
	if exitErr.ExitCode() != 1 {
		t.Errorf("expected exit code 1, got %d", exitErr.ExitCode())
	}
	if !strings.Contains(stderr, "it failed") {
		t.Errorf("expected stderr to contain 'it failed', got: %q", stderr)
	}
}

func TestErrorFailVoidPanic(t *testing.T) {
	exePath, diags := compiler.CompileOnly("testdata/error_fail_void.t")
	if len(diags.Diags) > 0 {
		t.Fatalf("unexpected compile errors: %v", diags.Diags)
	}

	stderr, err := compiler.RunBinary(exePath)
	if err == nil {
		t.Fatal("expected non-zero exit code, got success")
	}
	exitErr, ok := err.(*exec.ExitError)
	if !ok {
		t.Fatalf("expected *exec.ExitError, got: %v", err)
	}
	if exitErr.ExitCode() != 1 {
		t.Errorf("expected exit code 1, got %d", exitErr.ExitCode())
	}
	if !strings.Contains(stderr, "boom") {
		t.Errorf("expected stderr to contain 'boom', got: %q", stderr)
	}
}

func TestErrorFailVoidOK(t *testing.T) {
	output, diags := compiler.CompileFile("testdata/error_fail_void_ok.t")
	if len(diags.Diags) > 0 {
		t.Fatalf("unexpected errors: %v", diags.Diags)
	}
	expected := "ok\n"
	if output != expected {
		t.Errorf("expected %q, got %q", expected, output)
	}
}

func TestErrorCatchVoid(t *testing.T) {
	exePath, diags := compiler.CompileOnly("testdata/error_catch_void.t")
	if len(diags.Diags) > 0 {
		t.Fatalf("unexpected compile errors: %v", diags.Diags)
	}

	cmd := exec.Command(exePath)
	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	err := cmd.Run()
	if err == nil {
		t.Fatal("expected non-zero exit code, got success")
	}
	exitErr, ok := err.(*exec.ExitError)
	if !ok {
		t.Fatalf("expected *exec.ExitError, got: %v", err)
	}
	if exitErr.ExitCode() != 1 {
		t.Errorf("expected exit code 1, got %d", exitErr.ExitCode())
	}
	output := stdout.String() + stderr.String()
	if !strings.Contains(output, "boom") {
		t.Errorf("expected output to contain 'boom', got stdout=%q stderr=%q", stdout.String(), stderr.String())
	}
}

func TestErrorCatchTuple(t *testing.T) {
	output, diags := compiler.CompileFile("testdata/error_catch_tuple.t")
	if len(diags.Diags) > 0 {
		t.Fatalf("unexpected errors: %v", diags.Diags)
	}
	expected := "val: 42\n"
	if output != expected {
		t.Errorf("expected %q, got %q", expected, output)
	}
}

func TestErrorCatchTupleFail(t *testing.T) {
	exePath, diags := compiler.CompileOnly("testdata/error_catch_tuple_fail.t")
	if len(diags.Diags) > 0 {
		t.Fatalf("unexpected compile errors: %v", diags.Diags)
	}

	cmd := exec.Command(exePath)
	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	err := cmd.Run()
	if err == nil {
		t.Fatal("expected non-zero exit code, got success")
	}
	exitErr, ok := err.(*exec.ExitError)
	if !ok {
		t.Fatalf("expected *exec.ExitError, got: %v", err)
	}
	if exitErr.ExitCode() != 1 {
		t.Errorf("expected exit code 1, got %d", exitErr.ExitCode())
	}
	output := stdout.String() + stderr.String()
	if !strings.Contains(output, "connection refused") {
		t.Errorf("expected output to contain 'connection refused', got stdout=%q stderr=%q", stdout.String(), stderr.String())
	}
}

func TestNilPointerPanic(t *testing.T) {
	exePath, diags := compiler.CompileOnly("testdata/nil_ptr_panic.t")
	if len(diags.Diags) > 0 {
		t.Fatalf("unexpected compile errors: %v", diags.Diags)
	}

	stderr, err := compiler.RunBinary(exePath)
	if err == nil {
		t.Fatal("expected non-zero exit code, got success")
	}
	exitErr, ok := err.(*exec.ExitError)
	if !ok {
		t.Fatalf("expected *exec.ExitError, got: %v", err)
	}
	if exitErr.ExitCode() != 1 {
		t.Errorf("expected exit code 1, got %d", exitErr.ExitCode())
	}
	if !strings.Contains(stderr, "null pointer deref") {
		t.Errorf("expected stderr to contain 'null pointer deref', got: %q", stderr)
	}
}
