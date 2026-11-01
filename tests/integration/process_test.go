package integration

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"testing"
	"time"

	"github.com/HicaroD/Telia/tests/compiler"
)

func runProcessSource(t *testing.T, body string, args []string, flags ...string) {
	t.Helper()
	path := compileProcessSource(t, body, flags...)
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	output, err := exec.CommandContext(ctx, path, args...).CombinedOutput()
	if err != nil {
		t.Fatalf("execution: %v; output %q", err, output)
	}
}

func compileProcessSource(t *testing.T, body string, flags ...string) string {
	t.Helper()
	file := filepath.Join(t.TempDir(), "process.t")
	source := "package main\nuse \"std::process\"\nuse \"std::errors\"\n" + body
	if err := os.WriteFile(file, []byte(source), 0600); err != nil {
		t.Fatal(err)
	}
	path, diags := compiler.CompileOnlyWithCFlags(file, flags...)
	if diags.HasErrors() {
		t.Fatalf("compile: %v", diags.Diags)
	}
	t.Cleanup(func() { removeTestBinary(t, path) })
	return path
}

func TestProcessRejectsInvalidNativeUTF8(t *testing.T) {
	if runtime.GOOS == "windows" {
		t.Skip("POSIX native byte arguments only")
	}
	for _, opt := range []string{"-O0", "-O3"} {
		runProcessSource(t, `fn main() i32 {
snapshot, err := process::args()
if err.kind != errors::Error::InvalidUTF8 or process::count(&snapshot) != 0 {
return 1
}
value, access_err := process::get(&snapshot, 0)
if access_err.kind != errors::Error::OutOfBounds or value != "" {
return 2
}
process::free(&snapshot)
return 0
}`, []string{"valid", "bad\xff", "later"}, processControlFlags(t, opt, "allocation_control.h")...)
	}
}

func TestProcessArgumentsPreserveValuesAndOrder(t *testing.T) {
	args := []string{"", "two words", `a"b\c`, "世界-é-😀"}
	body := fmt.Sprintf(`fn main() i32 {
snapshot, err := process::args()
if err != nil or process::count(&snapshot) != %d {
return 1
}
`, len(args))
	for i, arg := range args {
		body += fmt.Sprintf(`value%d, err%d := process::get(&snapshot, %d)
if err%d != nil or value%d != %s {
return %d
}
`, i, i, i, i, i, utf8Literal(arg), i+2)
	}
	body += "process::free(&snapshot)\nreturn 0\n}\n"
	for _, opt := range []string{"-O0", "-O3"} {
		t.Run(opt, func(t *testing.T) { runProcessSource(t, body, args, opt) })
	}
}

func TestProcessArgumentBoundsAndCleanup(t *testing.T) {
	for _, args := range [][]string{nil, {"first", "", "last"}} {
		body := fmt.Sprintf(`fn main() i32 {
snapshot, err := process::args()
if err != nil or process::count(&snapshot) != %d {
return 1
}
value, access_err := process::get(&snapshot, %d)
if value != "" or access_err.kind != errors::Error::OutOfBounds {
return 2
}
huge, huge_err := process::get(&snapshot, 18446744073709551615)
if huge != "" or huge_err.kind != errors::Error::OutOfBounds {
return 3
}
other, other_err := process::args()
if other_err != nil {
return 4
}
process::free(&snapshot)
process::free(&snapshot)
if process::count(&snapshot) != 0 or process::count(&other) != %d {
return 5
}
empty, empty_err := process::get(&snapshot, 0)
if empty != "" or empty_err.kind != errors::Error::OutOfBounds {
return 6
}
process::free(&other)
return 0
}`, len(args), len(args), len(args))
		for _, opt := range []string{"-O0", "-O3"} {
			runProcessSource(t, body, args, opt)
		}
	}
}

const processAllocationControls = `@[default_cc="c"]
extern {
fn _test_fail_after(count i64)
fn _test_live_allocations() i64
}
`

func processControlFlags(t *testing.T, opt, headerName string) []string {
	t.Helper()
	header, err := filepath.Abs("testdata/" + headerName)
	if err != nil {
		t.Fatal(err)
	}
	flags := []string{opt, "-include", header}
	if runtime.GOOS != "windows" {
		sanitizers := os.Getenv("TELIA_TEST_SANITIZERS")
		if sanitizers == "" {
			sanitizers = "undefined"
		}
		flags = append(flags, "-fsanitize="+sanitizers, "-fno-sanitize-recover=all")
	}
	return flags
}

func TestProcessArgumentAllocationFailureCleanup(t *testing.T) {
	for _, opt := range []string{"-O0", "-O3"} {
		runProcessSource(t, processAllocationControls+`fn main() i32 {
limit i64 := 0
while limit < 4 {
_test_fail_after(limit)
snapshot, err := process::args()
if err.kind != errors::Error::AllocationFailed or process::count(&snapshot) != 0 {
return 1
}
process::free(&snapshot)
if _test_live_allocations() != 0 {
return 2
}
limit = limit + 1
}
_test_fail_after(4)
snapshot, err := process::args()
if err != nil or process::count(&snapshot) != 4 {
return 3
}
process::free(&snapshot)
if _test_live_allocations() != 0 {
return 4
}
return 0
}`, []string{"first", "", "世界", "last"}, processControlFlags(t, opt, "allocation_control.h")...)
	}
}

func TestProcessWindowsRejectsInvalidWideArguments(t *testing.T) {
	if runtime.GOOS != "windows" {
		t.Skip("native Windows UTF-16 boundary")
	}
	for _, opt := range []string{"-O0", "-O3"} {
		runProcessSource(t, processAllocationControls+`@[default_cc="c"]
extern {
fn _test_wide_mode(mode i32)
}
fn main() i32 {
mode i32 := 1
while mode <= 3 {
_test_wide_mode(mode)
snapshot, err := process::args()
if err.kind != errors::Error::InvalidUTF8 or process::count(&snapshot) != 0 {
return 1
}
process::free(&snapshot)
if _test_live_allocations() != 0 {
return 2
}
mode = mode + 1
}
return 0
}`, nil, processControlFlags(t, opt, "process_wide_control.h")...)
	}
}

func TestProcessEntryPointPreservesExitAndDeferredCleanup(t *testing.T) {
	for _, opt := range []string{"-O0", "-O3"} {
		for _, signature := range []string{"fn main() i32", "fn main()"} {
			result, exitCode := "return 37", 37
			if signature == "fn main()" {
				result, exitCode = "return", 0
			}
			body := "use \"std::io\"\n" + processAllocationControls + fmt.Sprintf(`%s {
snapshot, err := process::args()
if err != nil {
%s
}
defer process::free(&snapshot)
defer io::println("cleanup")
%s
}
`, signature, result, result)
			path := compileProcessSource(
				t,
				body,
				processControlFlags(t, opt, "allocation_control.h")...)
			ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
			cmd := exec.CommandContext(ctx, path, "arg")
			output, err := cmd.CombinedOutput()
			cancel()
			if cmd.ProcessState == nil || cmd.ProcessState.ExitCode() != exitCode ||
				string(output) != "cleanup\n" {
				t.Fatalf("entry point: exit %v, output %q, error %v", cmd.ProcessState, output, err)
			}
		}
	}
}

func TestProcessArgumentExample(t *testing.T) {
	path, diags := compiler.CompileOnly(filepath.Join("..", "..", "examples", "process_args.t"))
	if diags.HasErrors() {
		t.Fatalf("compile example: %v", diags.Diags)
	}
	defer removeTestBinary(t, path)
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	output, err := exec.CommandContext(ctx, path, "", "two words", "世界", "😀").CombinedOutput()
	if err != nil || string(output) != "\ntwo words\n世界\n😀\n" {
		t.Fatalf("example: %v, output %q", err, output)
	}
}
