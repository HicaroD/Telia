package integration

import (
	"bytes"
	"context"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"testing"
	"time"

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

func TestStringOutputPreservesBytes(t *testing.T) {
	output, diags := compiler.CompileFile("testdata/string_bytes.t")
	if len(diags.Diags) > 0 {
		t.Fatalf("unexpected errors: %v", diags.Diags)
	}
	if want := "A\x00B — Olá\n\n"; output != want {
		t.Fatalf("want %q, got %q", want, output)
	}
}

func TestStringEqualityComparesBytes(t *testing.T) {
	output, diags := compiler.CompileFile("testdata/string_equality.t")
	if len(diags.Diags) > 0 {
		t.Fatalf("unexpected errors: %v", diags.Diags)
	}
	if want := "equal\ndifferent\nempty\nlength\n"; output != want {
		t.Fatalf("want %q, got %q", want, output)
	}
}

func TestStringSlicesUseCheckedByteRanges(t *testing.T) {
	output, diags := compiler.CompileFile("testdata/string_slice.t")
	if len(diags.Diags) > 0 {
		t.Fatalf("unexpected errors: %v", diags.Diags)
	}
	if want := "\x00B\n\xc3\nempty\nbytes\nreversed\npast end\n"; output != want {
		t.Fatalf("want %q, got %q", want, output)
	}
}

func TestOwnedStringsCloneIndependentlyAndResetOnFree(t *testing.T) {
	output, diags := compiler.CompileFile("testdata/string_owned.t")
	if len(diags.Diags) > 0 {
		t.Fatalf("unexpected errors: %v", diags.Diags)
	}
	if want := "Hello\x00é\nreset\nHello\x00é!\n"; output != want {
		t.Fatalf("want %q, got %q", want, output)
	}
}

func TestDeferPreservesReturnValueAndUnwindsEnclosingScopes(t *testing.T) {
	output, diags := compiler.CompileFile("testdata/defer_return_value.t")
	if len(diags.Diags) > 0 {
		t.Fatalf("unexpected errors: %v", diags.Diags)
	}
	if want := "7\ninner\nouter\n1\n"; output != want {
		t.Fatalf("want %q, got %q", want, output)
	}
}

func TestDeferredCleanupUsesItsDeclaringScope(t *testing.T) {
	output, diags := compiler.CompileFile("testdata/defer_shadow.t")
	if len(diags.Diags) > 0 {
		t.Fatalf("unexpected errors: %v", diags.Diags)
	}
	if want := "2\n1\n"; output != want {
		t.Fatalf("want %q, got %q", want, output)
	}
}

func TestCStringConversionOwnsTerminatedStorageAndRejectsZeros(t *testing.T) {
	output, diags := compiler.CompileFile("testdata/string_c_interop.t")
	if len(diags.Diags) > 0 {
		t.Fatalf("unexpected errors: %v", diags.Diags)
	}
	want := "Olá\nOlá\n[]\nrejected\nstatic\n"
	// These lines use C stdio, whose Windows text mode translates LF to CRLF.
	// Keep captured output untouched so byte-preservation tests remain strict.
	if runtime.GOOS == "windows" {
		want = strings.ReplaceAll(want, "\n", "\r\n")
	}
	if output != want {
		t.Fatalf("want %q, got %q", want, output)
	}
}

func TestStringAllocationFailurePreservesExistingOwners(t *testing.T) {
	header, err := filepath.Abs("testdata/allocation_control.h")
	if err != nil {
		t.Fatal(err)
	}
	path, diags := compiler.CompileOnlyWithCFlags("testdata/string_allocation_failure.t", "-include", header)
	if len(diags.Diags) > 0 {
		t.Fatalf("unexpected errors: %v", diags.Diags)
	}
	output, err := compiler.RunBinary(path)
	if err != nil {
		t.Fatalf("execution: %v: %s", err, output)
	}
	if want := "zero allocation\nheap failed\nclone failed\nconcat failed\nconversion failed\npreserved\nreleased\n"; output != want {
		t.Fatalf("want %q, got %q", want, output)
	}
}

func TestArenaStringsSurviveGrowthAndSupportResetAndDestroy(t *testing.T) {
	output, diags := compiler.CompileFile("testdata/string_arena.t")
	if len(diags.Diags) > 0 {
		t.Fatalf("unexpected errors: %v", diags.Diags)
	}
	if want := "stable\x00é\nstable\x00é!\nC text\nnew text\nagain\n"; output != want {
		t.Fatalf("want %q, got %q", want, output)
	}
}

func TestArenaFailurePreservesStorageAndResetReusesIt(t *testing.T) {
	header, err := filepath.Abs("testdata/allocation_control.h")
	if err != nil {
		t.Fatal(err)
	}
	path, diags := compiler.CompileOnlyWithCFlags("testdata/arena_failure.t", "-include", header)
	if len(diags.Diags) > 0 {
		t.Fatalf("unexpected errors: %v", diags.Diags)
	}
	output, err := compiler.RunBinary(path)
	if err != nil {
		t.Fatalf("execution: %v: %s", err, output)
	}
	if want := "growth failed\npreserved\nreused\noverflow rejected\nalignment rejected\nreleased\n"; output != want {
		t.Fatalf("want %q, got %q", want, output)
	}
}

func TestQualifiedFailPanicsOnStringRangeError(t *testing.T) {
	path, diags := compiler.CompileOnly("testdata/string_slice_fail.t")
	if len(diags.Diags) > 0 {
		t.Fatalf("unexpected errors: %v", diags.Diags)
	}
	output, err := compiler.RunBinary(path)
	if err == nil || output != "runtime panic: invalid string range\n" {
		t.Fatalf("want range panic, got output %q, error %v", output, err)
	}
}

func TestStringMemoryWithNativeSanitizers(t *testing.T) {
	sanitizers := os.Getenv("TELIA_TEST_SANITIZERS")
	if sanitizers == "" {
		sanitizers = "undefined"
	}
	for _, fixture := range []string{"string_owned.t", "string_c_interop.t", "string_arena.t", "string_allocation_failure.t", "arena_failure.t"} {
		t.Run(fixture, func(t *testing.T) {
			header, err := filepath.Abs("testdata/allocation_control.h")
			if err != nil {
				t.Fatal(err)
			}
			path, diags := compiler.CompileOnlyWithCFlags("testdata/"+fixture,
				"-fsanitize="+sanitizers, "-fno-sanitize-recover=all", "-fno-omit-frame-pointer", "-include", header)
			if len(diags.Diags) > 0 {
				t.Fatalf("sanitized compilation: %v", diags.Diags)
			}
			ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
			defer cancel()
			cmd := exec.CommandContext(ctx, path)
			output, err := cmd.CombinedOutput()
			if err != nil {
				t.Fatalf("sanitized execution: %v\n%s", err, output)
			}
			if bytes.Contains(output, []byte("runtime error:")) || bytes.Contains(output, []byte("Sanitizer")) {
				t.Fatalf("sanitizer finding: %s", output)
			}
		})
	}
}

func TestStringValuesAndEmptyErrors(t *testing.T) {
	for _, tc := range []struct{ file, want string }{
		{"string_values.t", "\nA\x00é\n\n"},
		{"string_empty_error.t", "present\n\ncaught\n"},
	} {
		t.Run(tc.file, func(t *testing.T) {
			output, diags := compiler.CompileFile("testdata/" + tc.file)
			if len(diags.Diags) > 0 {
				t.Fatalf("unexpected errors: %v", diags.Diags)
			}
			if output != tc.want {
				t.Fatalf("want %q, got %q", tc.want, output)
			}
		})
	}
}

func TestPanicPreservesMessageBytes(t *testing.T) {
	path, diags := compiler.CompileOnly("testdata/string_error_bytes.t")
	if len(diags.Diags) > 0 {
		t.Fatalf("unexpected errors: %v", diags.Diags)
	}
	output, err := compiler.RunBinary(path)
	if err == nil || output != "runtime panic: A\x00Bé\n" {
		t.Fatalf("want byte-preserving panic, got %q, %v", output, err)
	}
}

func TestDocumentedStringsExample(t *testing.T) {
	header, err := filepath.Abs("testdata/allocation_control.h")
	if err != nil {
		t.Fatal(err)
	}
	path, diags := compiler.CompileOnlyWithCFlags("../../examples/strings.t",
		"-fsanitize=undefined", "-fno-sanitize-recover=all", "-include", header)
	if len(diags.Diags) > 0 {
		t.Fatalf("unexpected errors: %v", diags.Diags)
	}
	output, err := compiler.RunBinary(path)
	if err != nil {
		t.Fatalf("example execution: %v: %s", err, output)
	}
	if want := "Hello\nHello, Telia!\nHello, Telia!\ninvalid string range\n"; output != want {
		t.Fatalf("want %q, got %q", want, output)
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

func TestErrorStructField(t *testing.T) {
	output, diags := compiler.CompileFile("testdata/error_struct_field.t")
	if len(diags.Diags) > 0 {
		t.Fatalf("unexpected errors: %v", diags.Diags)
	}
	expected := "field broke\n"
	if output != expected {
		t.Errorf("expected %q, got %q", expected, output)
	}
}

func TestErrorParam(t *testing.T) {
	output, diags := compiler.CompileFile("testdata/error_param.t")
	if len(diags.Diags) > 0 {
		t.Fatalf("unexpected errors: %v", diags.Diags)
	}
	expected := "param broke\n"
	if output != expected {
		t.Errorf("expected %q, got %q", expected, output)
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
