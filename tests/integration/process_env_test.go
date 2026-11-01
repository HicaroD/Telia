package integration

import (
	"context"
	"fmt"
	"os/exec"
	"path/filepath"
	"runtime"
	"testing"
	"time"

	"github.com/HicaroD/Telia/tests/compiler"
)

func runEnvironmentSource(t *testing.T, body string, env []string, flags ...string) {
	t.Helper()
	path := compileProcessSource(t, "use \"std::strings\"\n"+body, flags...)
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	cmd := exec.CommandContext(ctx, path)
	cmd.Env = env
	output, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("execution: %v; output %q", err, output)
	}
}

func TestProcessEnvironmentRejectsInvalidText(t *testing.T) {
	for _, opt := range []string{"-O0", "-O3"} {
		for _, value := range []struct{ name, kind string }{
			{"", "InvalidEnvironmentName"},
			{"A=B", "InvalidEnvironmentName"},
			{"A\x00B", "EmbeddedZero"},
			{"A\xff", "InvalidUTF8"},
		} {
			body := fmt.Sprintf(`fn main() i32 {
value, found, err := process::lookup_env(%s)
if err.kind != errors::Error::%s or found or strings::view(&value) != "" {
return 1
}
strings::free(&value)
return 0
}`, utf8Literal(value.name), value.kind)
			runEnvironmentSource(t, body, []string{}, opt)
		}
		if runtime.GOOS != "windows" {
			runEnvironmentSource(t, `fn main() i32 {
value, found, err := process::lookup_env("TELIA_ENV_TEST")
if err.kind != errors::Error::InvalidUTF8 or found or strings::view(&value) != "" {
return 1
}
strings::free(&value)
return 0
}`, []string{"TELIA_ENV_TEST=bad\xff"}, processControlFlags(t, opt, "allocation_control.h")...)
		}
	}
}

func TestProcessEnvironmentMissingEmptyAndPresent(t *testing.T) {
	for _, opt := range []string{"-O0", "-O3"} {
		for _, value := range []struct {
			env   []string
			text  string
			found bool
		}{
			{[]string{}, "", false},
			{[]string{"TELIA_ENV_TEST="}, "", true},
			{[]string{"TELIA_ENV_TEST=two words 世界 😀"}, "two words 世界 😀", true},
		} {
			body := fmt.Sprintf(`fn main() i32 {
value, found, err := process::lookup_env("TELIA_ENV_TEST")
if err != nil or found != %t or strings::view(&value) != %s {
return 1
}
strings::free(&value)
strings::free(&value)
return 0
}`, value.found, utf8Literal(value.text))
			runEnvironmentSource(t, body, value.env, opt)
		}
	}
}

func TestProcessEnvironmentUnicodeNamesAndNativeComparison(t *testing.T) {
	for _, opt := range []string{"-O0", "-O3"} {
		body := fmt.Sprintf(`fn main() i32 {
value, found, err := process::lookup_env(%s)
if err != nil or not found or strings::view(&value) != %s {
return 1
}
strings::free(&value)
other, other_found, other_err := process::lookup_env("telia_case_test")
if other_err != nil or other_found != %t {
return 2
}
strings::free(&other)
return 0
}`, utf8Literal("TELIA_世界_😀"), utf8Literal("é 世界 😀"), runtime.GOOS == "windows")
		runEnvironmentSource(t, body, []string{"TELIA_世界_😀=é 世界 😀", "TELIA_CASE_TEST=value"}, opt)
	}
}

func TestProcessEnvironmentOwnershipAndAllocationFailures(t *testing.T) {
	for _, opt := range []string{"-O0", "-O3"} {
		runEnvironmentSource(t, processAllocationControls+`@[default_cc="c"]
extern {
fn _test_change_environment()
}

fn main() i32 {
limit i64 := 0
succeeded := false
while limit < 8 and not succeeded {
_test_fail_after(limit)
value, found, err := process::lookup_env("TELIA_ENV_TEST")
if err != nil {
if err.kind != errors::Error::AllocationFailed or found or strings::view(&value) != "" {
return 1
}
} else {
if not found or strings::view(&value) != "original" or limit == 0 {
return 2
}
_test_change_environment()
if strings::view(&value) != "original" {
return 3
}
succeeded = true
}
strings::free(&value)
if _test_live_allocations() != 0 {
return 4
}
limit = limit + 1
}
if not succeeded {
return 5
}
return 0
}`, []string{"TELIA_ENV_TEST=original"}, processControlFlags(t, opt, "process_env_control.h")...)
	}
}

func TestProcessEnvironmentWindowsFailures(t *testing.T) {
	if runtime.GOOS != "windows" {
		t.Skip("native Windows UTF-16 boundary")
	}
	for _, opt := range []string{"-O0", "-O3"} {
		body := processAllocationControls + `@[default_cc="c"]
extern {
fn _test_env_mode(mode i32)
}
fn main() i32 {
mode i32 := 1
while mode <= 7 {
_test_env_mode(mode)
value, found, err := process::lookup_env("TELIA_ENV_TEST")
if found or strings::view(&value) != "" {
return 1
}
if mode == 4 or mode == 5 {
if err.kind != errors::Error::EnvironmentLookupFailed {
return 2
}
} else {
if err.kind != errors::Error::InvalidUTF8 {
return 3
}
}
strings::free(&value)
if _test_live_allocations() != 0 {
return 4
}
mode = mode + 1
}
return 0
}`
		runEnvironmentSource(
			t,
			body,
			[]string{},
			processControlFlags(t, opt, "process_env_wide_control.h")...)
	}
}

func TestProcessEnvironmentExample(t *testing.T) {
	for _, opt := range []string{"-O0", "-O3"} {
		path, diags := compiler.CompileOnlyWithCFlags(
			filepath.Join("..", "..", "examples", "process_env.t"),
			opt,
		)
		if diags.HasErrors() {
			t.Fatalf("compile example: %v", diags.Diags)
		}
		t.Cleanup(func() { removeTestBinary(t, path) })
		cases := []struct {
			env    []string
			output string
			exit   int
		}{
			{[]string{}, "Hello from Telia", 0},
			{[]string{"TELIA_GREETING="}, "", 0},
			{[]string{"TELIA_GREETING=Hello 世界 😀"}, "Hello 世界 😀", 0},
		}
		if runtime.GOOS != "windows" {
			cases = append(cases, struct {
				env    []string
				output string
				exit   int
			}{
				[]string{"TELIA_GREETING=bad\xff"}, "could not read TELIA_GREETING\n", 1,
			})
		}
		for _, value := range cases {
			ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
			cmd := exec.CommandContext(ctx, path)
			cmd.Env = value.env
			output, err := cmd.CombinedOutput()
			cancel()
			if cmd.ProcessState == nil || cmd.ProcessState.ExitCode() != value.exit ||
				string(output) != value.output {
				t.Fatalf("example: %v; output %q, want %q", err, output, value.output)
			}
		}
	}
}
