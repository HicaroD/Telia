package main

import (
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
	"testing"
)

func writeTestPackage(t *testing.T, files map[string]string) string {
	t.Helper()
	dir := filepath.Join(t.TempDir(), "math")
	for name, text := range files {
		path := filepath.Join(dir, name)
		if err := os.MkdirAll(filepath.Dir(path), 0755); err != nil {
			t.Fatal(err)
		}
		if err := os.WriteFile(path, []byte(text), 0600); err != nil {
			t.Fatal(err)
		}
	}
	return dir
}

func TestTestCommandDiscoversNamedBlocksAndHelpers(t *testing.T) {
	dir := writeTestPackage(t, map[string]string{
		"math.t":      "package math\nfn add(a i32, b i32) i32 {\nreturn a + b\n}\n",
		"math_test.t": "package math\nfn helper() {\nvalue := add(2, 3)\n}\nfn test_not_discovered() {\n}\ntest \"uses package helpers\" {\nhelper()\n}\ntest \"second test\" {\n}\n",
	})
	out, err := runTestCompiler(t, dir, "test")
	if err != nil {
		t.Fatalf("test command: %v\n%s", err, out)
	}
	for _, want := range []string{"PASS uses package helpers", "PASS second test", "2 passed, 0 failed"} {
		if !strings.Contains(out, want) {
			t.Fatalf("missing %q:\n%s", want, out)
		}
	}
	if strings.Contains(out, "PASS test_not_discovered") {
		t.Fatalf("discovered helper:\n%s", out)
	}
}

func TestTestCommandOwnsEntryPoint(t *testing.T) {
	for _, main := range []string{"", "fn main() i32 {\nreturn 73\n}\n"} {
		body := ""
		if main != "" {
			body = "testing::assert_eq(main(), 73)\n"
		}
		dir := writeTestPackage(t, map[string]string{
			"main.t":      "package main\n" + main,
			"main_test.t": "package main\nuse \"std::testing\"\ntest \"runs instead of main\" {\n" + body + "}\n",
		})
		out, err := runTestCompiler(t, dir, "test")
		if err != nil || !strings.Contains(out, "PASS runs instead of main") {
			t.Fatalf("test entry point: %v\n%s", err, out)
		}
	}
}

func TestTestingAssertionsStopHelperAndContinueTests(t *testing.T) {
	dir := writeTestPackage(t, map[string]string{
		"math_test.t": `package math
use "std::testing"
use "std::io"
fn add(a i32, b i32) i32 {
return a + b
}
fn check_sum() {
testing::assert_eq(add(2, 3), 7)
io::println("UNREACHABLE helper")
}
test "passing assertions" {
testing::assert(true)
testing::assert_eq(add(2, 3), 5)
testing::assert_ne(2, 3)
}
test "bad sum" {
check_sum()
io::println("UNREACHABLE test")
}
test "still runs" {
testing::assert(true)
}
test "false condition" {
testing::assert(false)
}
test "equal values" {
testing::assert_ne(4, 4)
}
`,
	})
	out, err := runTestCompiler(t, dir, "test", "--color=never")
	if err == nil {
		t.Fatalf("failing assertions exited successfully:\n%s", out)
	}
	for _, want := range []string{"PASS passing assertions", "FAIL bad sum", "math_test.t:8:", "expected: 7", "actual: 5", "PASS still runs", "FAIL false condition", "assertion failed", "FAIL equal values", "unexpected: 4", "2 passed, 3 failed"} {
		if !strings.Contains(out, want) {
			t.Fatalf("missing %q:\n%s", want, out)
		}
	}
	if strings.Contains(out, "UNREACHABLE") {
		t.Fatalf("assertion did not stop test:\n%s", out)
	}
}

func TestTestResultColorsRespectFlagsAndEnvironment(t *testing.T) {
	dir := writeTestPackage(
		t,
		map[string]string{
			"math_test.t": "package math\nuse \"std::testing\"\ntest \"passes\" {\ntesting::assert(true)\n}\ntest \"fails\" {\ntesting::assert(false)\n}\n",
		},
	)
	for _, tc := range []struct {
		mode          string
		noColor, ansi bool
	}{
		{"auto", false, false}, {"auto", true, false}, {"always", false, true}, {"always", true, true}, {"never", false, false},
	} {
		t.Run(tc.mode+"/NO_COLOR="+strconv.FormatBool(tc.noColor), func(t *testing.T) {
			cmd := exec.Command(
				os.Args[0],
				"-test.run=^TestCompilerDiagnosticProcess$",
				"--",
				"test",
				"--color="+tc.mode,
			)
			cmd.Dir = dir
			for _, env := range testCompilerEnv(t) {
				if !strings.HasPrefix(env, "NO_COLOR=") {
					cmd.Env = append(cmd.Env, env)
				}
			}
			if tc.noColor {
				cmd.Env = append(cmd.Env, "NO_COLOR=1")
			}
			out, err := cmd.CombinedOutput()
			if err == nil {
				t.Fatalf("failing test succeeded:\n%s", out)
			}
			for _, pair := range []struct{ plain, colored string }{{"PASS passes", "\x1b[32mPASS\x1b[0m passes"}, {"FAIL fails", "\x1b[31mFAIL\x1b[0m fails"}} {
				want := pair.plain
				if tc.ansi {
					want = pair.colored
				}
				if !strings.Contains(string(out), want) {
					t.Fatalf("missing %q:\n%s", want, out)
				}
			}
			if strings.Contains(string(out), "\x1b[") != tc.ansi {
				t.Fatalf("unexpected color output:\n%s", out)
			}
		})
	}
}

func TestDuplicateTestDescriptionsAreRejected(t *testing.T) {
	dir := writeTestPackage(t, map[string]string{
		"a_test.t": "package math\ntest \"same description\" {\n}\n",
		"b_test.t": "package math\ntest \"same description\" {\n}\n",
	})
	out, err := runTestCompiler(t, dir, "test", "--color=never")
	if err == nil || !strings.Contains(out, "duplicate test description") ||
		!strings.Contains(out, "b_test.t:2:") ||
		!strings.Contains(out, "a_test.t:2:") {
		t.Fatalf("expected duplicate diagnostic with both locations: %v\n%s", err, out)
	}
}

func TestTestSyntaxAndCompilationDiagnostics(t *testing.T) {
	for _, tc := range []struct{ name, file, text, want string }{
		{"outside test file", "math.t", "test \"wrong place\" {\n}\n", "test blocks must be declared in *_test.t files"},
		{"missing description", "math_test.t", "test {\n}\n", "expected a string description"},
		{"description is type", "math_test.t", "test string {\n}\n", "expected a string description"},
		{"missing body", "math_test.t", "test \"missing body\"\n", "expected '{'"},
		{"invalid body", "math_test.t", "test \"bad body\" {\nunknown()\n}\n", "not defined"},
		{"wrong condition", "math_test.t", "use \"std::testing\"\ntest \"bad condition\" {\ntesting::assert(42)\n}\n", "Type mismatch"},
		{"wrong arity", "math_test.t", "use \"std::testing\"\ntest \"missing operand\" {\ntesting::assert_eq(42)\n}\n", "requires 2 value arguments"},
		{"mismatched values", "math_test.t", "use \"std::testing\"\ntest \"bad types\" {\ntesting::assert_eq(42, \"42\")\n}\n", "Type mismatch"},
		{"unsupported values", "math_test.t", "use \"std::testing\"\nstruct Point {\nx i32\n}\ntest \"bad type\" {\ntesting::assert_eq(Point{x: 1}, Point{x: 1})\n}\n", "equality assertions support"},
	} {
		t.Run(tc.name, func(t *testing.T) {
			dir := writeTestPackage(t, map[string]string{tc.file: "package math\n" + tc.text})
			out, err := runTestCompiler(t, dir, "test", "--color=never")
			if err == nil || !strings.Contains(out, tc.want) ||
				!strings.Contains(out, tc.file+":") {
				t.Fatalf("expected %q with source location: %v\n%s", tc.want, err, out)
			}
			if strings.Contains(out, "PASS ") || strings.Contains(out, "panic:") {
				t.Fatalf("invalid program ran or compiler panicked:\n%s", out)
			}
		})
	}
}

func TestNoTestsAndImportedTestExclusion(t *testing.T) {
	dir := writeTestPackage(t, map[string]string{
		"math.t":            "package math\nuse \"pkg::lib\"\nfn value() i32 {\nreturn lib::value()\n}\n",
		"lib/lib.t":         "package lib\nfn value() i32 {\nreturn 1\n}\n",
		"lib/broken_test.t": "this must not be parsed",
	})
	out, err := runTestCompiler(t, dir, "test")
	if err != nil || !strings.Contains(out, "no tests found") ||
		!strings.Contains(out, "0 passed, 0 failed") {
		t.Fatalf("no tests: %v\n%s", err, out)
	}
}

func TestAssertionValuesAndSingleEvaluation(t *testing.T) {
	dir := writeTestPackage(t, map[string]string{"math_test.t": `package math
use "std::testing"
use "std::strings"
fn next(counter *i32) i32 {
*counter = *counter + 1
return *counter
}
test "evaluates operands once in order" {
counter i32 := 0
testing::assert_eq(next(&counter), next(&counter) - 1)
testing::assert_eq(counter, 2)
}
test "string equality uses contents" {
first := strings::clone("hello") @fail
second := strings::clone("hello") @fail
defer strings::free(&first)
defer strings::free(&second)
testing::assert_eq(strings::view(&first), strings::view(&second))
testing::assert_ne("hello", "world")
testing::assert_eq("a\0b", "a\0b")
testing::assert_ne("a\0b", "a\0c")
}
test "strings show values" {
testing::assert_eq("actual text", "expected text")
}
test "bools show values" {
testing::assert_eq(false, true)
}
test "floats show values" {
actual f64 := 1.25
testing::assert_eq(actual, 2.5)
}
test "signed integers show values" {
actual i64 := -9223372036854775808
testing::assert_eq(actual, 3)
}
test "unsigned integers show values" {
actual u64 := 18446744073709551615
testing::assert_eq(actual, 2)
}
test "wide signed integers show values" {
actual i128 := -1
actual = actual << 127
testing::assert_eq(actual, 1)
}
test "wide unsigned integers show values" {
actual u128 := 0
actual = ~actual
testing::assert_eq(actual, 1)
}
`})
	out, err := runTestCompiler(t, dir, "test", "--color=never")
	if err == nil {
		t.Fatalf("failed assertions succeeded:\n%s", out)
	}
	for _, want := range []string{"PASS evaluates operands once in order", "PASS string equality uses contents", "expected: \"expected text\"", "actual: \"actual text\"", "expected: true", "actual: false", "expected: 2.5", "actual: 1.25", "actual: -9223372036854775808", "actual: 18446744073709551615", "actual: -170141183460469231731687303715884105728", "actual: 340282366920938463463374607431768211455", "2 passed, 7 failed"} {
		if !strings.Contains(out, want) {
			t.Fatalf("missing %q:\n%s", want, out)
		}
	}
}

func TestBuildRejectsExplicitTestFile(t *testing.T) {
	dir := writeTestPackage(
		t,
		map[string]string{"main_test.t": "package main\nfn main() i32 {\nreturn 0\n}\n"},
	)
	out, err := runTestCompiler(t, dir, "build", "main_test.t")
	if err == nil || !strings.Contains(out, "telia test") {
		t.Fatalf("explicit test file was not rejected: %v\n%s", err, out)
	}
}

func TestGenericHelpersCanUseTestingAssertions(t *testing.T) {
	dir := writeTestPackage(
		t,
		map[string]string{
			"math_test.t": "package math\nuse \"std::testing\"\nfn check<T>(actual T, expected T) {\ntesting::assert_eq(actual, expected)\n}\ntest \"generic helper\" {\ncheck<i32>(3, 3)\n}\n",
		},
	)
	out, err := runTestCompiler(t, dir, "test")
	if err != nil || !strings.Contains(out, "PASS generic helper") {
		t.Fatalf("generic helper: %v\n%s", err, out)
	}
}

func TestTestDescriptionsPreserveStringBytes(t *testing.T) {
	dir := writeTestPackage(
		t,
		map[string]string{
			"math_test.t": "package math\nuse \"std::testing\"\ntest \"π 100% \\\"quoted\\\" \\0a\" {\ntesting::assert(false)\n}\n",
		},
	)
	out, err := runTestCompiler(t, dir, "test", "--color=never")
	if err == nil || !strings.Contains(out, "FAIL π 100% \"quoted\" \x00a") ||
		!strings.Contains(out, "assertion failed in π 100% \"quoted\" \x00a") {
		t.Fatalf("description bytes lost: %v\n%q", err, out)
	}
}

func TestDocumentedTestingExample(t *testing.T) {
	dir, err := filepath.Abs("../../examples/testing/math")
	if err != nil {
		t.Fatal(err)
	}
	out, err := runTestCompiler(t, dir, "test", "--color=never")
	if err != nil || !strings.Contains(out, "2 passed, 0 failed") {
		t.Fatalf("documented example: %v\n%s", err, out)
	}
}

func TestAssertionFailureUnwindsDeferredCleanupAcrossHelpers(t *testing.T) {
	dir := writeTestPackage(t, map[string]string{"math_test.t": `package math
use "std::testing"
use "std::io"
fn inner() i32 {
defer io::println("cleanup inner first")
defer io::println("cleanup inner second")
testing::assert(false)
io::println("UNREACHABLE inner")
return 1
}
fn outer() {
defer io::println("cleanup outer")
value := inner() + 1
io::println("UNREACHABLE outer")
}
test "unwinds helpers" {
defer io::println("cleanup test")
outer()
io::println("UNREACHABLE test")
}
test "continues after cleanup" {
testing::assert(true)
}
`})
	out, err := runTestCompiler(t, dir, "test", "--color=never")
	if err == nil {
		t.Fatalf("failed assertion succeeded:\n%s", out)
	}
	want := "cleanup inner second\ncleanup inner first\ncleanup outer\ncleanup test\nFAIL unwinds helpers\nPASS continues after cleanup\n1 passed, 1 failed"
	if !strings.Contains(out, want) || strings.Contains(out, "UNREACHABLE") {
		t.Fatalf("cleanup did not unwind in order:\n%s", out)
	}
}

func TestFailuresDuringCleanupDoNotRepeatOrSkipRemainingActions(t *testing.T) {
	dir := writeTestPackage(t, map[string]string{"math_test.t": `package math
use "std::testing"
use "std::io"
fn bad_cleanup() {
defer io::println("cleanup nested")
testing::assert_eq(1, 2)
io::println("UNREACHABLE cleanup")
}
fn fallthrough() {
defer io::println("cleanup helper outer")
{
defer bad_cleanup()
}
io::println("UNREACHABLE fallthrough")
}
test "primary and cleanup failure" {
defer io::println("cleanup last")
defer bad_cleanup()
defer io::println("cleanup first")
testing::assert_eq(3, 4)
defer io::println("UNREACHABLE registration")
}
test "normal scope cleanup failure" {
defer io::println("cleanup test outer")
fallthrough()
io::println("UNREACHABLE caller")
}
test "deferred block failure" {
defer io::println("cleanup after block")
defer {
defer io::println("cleanup inside block")
testing::assert(false)
io::println("UNREACHABLE block")
}
}
test "still executes" {
testing::assert(true)
}
`})
	out, err := runTestCompiler(t, dir, "test", "--color=never")
	if err == nil {
		t.Fatalf("cleanup failures succeeded:\n%s", out)
	}
	var events []string
	for _, line := range strings.Split(out, "\n") {
		if strings.HasPrefix(line, "cleanup ") {
			events = append(events, line)
		}
	}
	want := "cleanup first\ncleanup nested\ncleanup last\ncleanup nested\ncleanup helper outer\ncleanup test outer\ncleanup inside block\ncleanup after block"
	if strings.Join(events, "\n") != want || strings.Contains(out, "UNREACHABLE") ||
		!strings.Contains(out, "1 passed, 3 failed") ||
		!strings.Contains(out, "expected: 4") ||
		!strings.Contains(out, "expected: 2") {
		t.Fatalf("cleanup failures were mishandled:\n%s", out)
	}
}

func TestAssertionCleanupReleasesResourcesAndReadsCurrentLocals(t *testing.T) {
	dir := writeTestPackage(t, map[string]string{"math_test.t": `package math
use "std::testing"
use "std::strings"
use "std::io"
fn fail_and_free(owner *strings::Owned) {
defer strings::free(owner)
testing::assert(false)
}
fn verify_freed(owner *strings::Owned) {
testing::assert_eq(strings::view(owner), "")
io::println("cleanup released owner")
}
test "resource cleanup" {
owner := strings::clone("owned allocation") @fail
defer verify_freed(&owner)
fail_and_free(&owner)
}
test "cleanup observes latest local value" {
value i32 := 1
defer testing::assert_eq(value, 2)
value = 2
testing::assert(false)
}
`})
	out, err := runTestCompiler(t, dir, "test", "--color=never")
	if err == nil || !strings.Contains(out, "cleanup released owner") ||
		!strings.Contains(out, "0 passed, 2 failed") ||
		strings.Count(out, "assertion failed") != 2 {
		t.Fatalf("cleanup lost resources or current locals: %v\n%s", err, out)
	}
}

func TestCleanupFailuresAbortReturnBreakAndContinue(t *testing.T) {
	dir := writeTestPackage(t, map[string]string{"math_test.t": `package math
use "std::testing"
use "std::io"
fn returning() i32 {
defer testing::assert(false)
return 42
}
fn looping(kind i32) {
while true {
defer testing::assert(false)
if kind == 0 {
break
}
continue
}
io::println("UNREACHABLE after loop")
}
test "return cleanup fails" {
value := returning()
io::println("UNREACHABLE after return")
}
test "break cleanup fails" {
looping(0)
io::println("UNREACHABLE after break")
}
test "continue cleanup fails" {
looping(1)
io::println("UNREACHABLE after continue")
}
`})
	out, err := runTestCompiler(t, dir, "test", "--color=never")
	if err == nil || strings.Contains(out, "UNREACHABLE") ||
		!strings.Contains(out, "0 passed, 3 failed") {
		t.Fatalf("cleanup failure did not abort control flow: %v\n%s", err, out)
	}
}

func TestSuccessfulTestsKeepNormalCleanupControlFlow(t *testing.T) {
	dir := writeTestPackage(t, map[string]string{"math_test.t": `package math
use "std::testing"
enum Choice {
Value
}
fn increment(counter *i32) {
*counter = *counter + 1
}
fn returning(counter *i32) i32 {
defer increment(counter)
return *counter
}
test "normal cleanup runs once" {
counter i32 := 0
testing::assert_eq(returning(&counter), 0)
testing::assert_eq(counter, 1)
index i32 := 0
while index < 3 {
defer increment(&counter)
index = index + 1
if index < 3 {
continue
}
break
}
testing::assert_eq(counter, 4)
result := match Choice::Value {
Choice::Value => {
defer increment(&counter)
yield 42
}
}
testing::assert_eq(result, 42)
testing::assert_eq(counter, 5)
}
`})
	out, err := runTestCompiler(t, dir, "test", "--color=never")
	if err != nil || !strings.Contains(out, "1 passed, 0 failed") {
		t.Fatalf("normal cleanup semantics changed: %v\n%s", err, out)
	}
}

func TestAssertionFunctionValuesPreserveCleanupAndContinuation(t *testing.T) {
	dir := writeTestPackage(t, map[string]string{"math_test.t": `package math
use "std::testing"
use "std::io"
test "assertion callback" {
defer io::println("cleanup callback")
check := testing::assert
check(false)
io::println("UNREACHABLE callback")
}
test "continues" {
testing::assert(true)
}
`})
	out, err := runTestCompiler(t, dir, "test", "--color=never")
	if err == nil || !strings.Contains(out, "cleanup callback") ||
		!strings.Contains(out, "1 passed, 1 failed") ||
		strings.Contains(out, "UNREACHABLE") {
		t.Fatalf("assertion callback bypassed runner cleanup: %v\n%s", err, out)
	}
}

func TestTestCommandDoesNotDependOnUserLibraryConfiguration(t *testing.T) {
	t.Setenv("XDG_CONFIG_HOME", t.TempDir())
	t.Setenv("T_STD_PATH", "")
	dir := writeTestPackage(t, map[string]string{
		"math_test.t": "package math\nuse \"std::testing\"\ntest \"fresh configuration\" {\ntesting::assert(true)\n}\n",
	})
	out, err := runTestCompiler(t, dir, "test", "--color=never")
	if err != nil || !strings.Contains(out, "PASS fresh configuration") {
		t.Fatalf("test command depended on user library configuration: %v\n%s", err, out)
	}
}

func testCompilerEnv(t *testing.T) []string {
	t.Helper()
	std, err := filepath.Abs("../../base/std")
	if err != nil {
		t.Fatal(err)
	}
	var env []string
	for _, entry := range os.Environ() {
		if !strings.HasPrefix(entry, "T_STD_PATH=") &&
			!strings.HasPrefix(entry, "XDG_CONFIG_HOME=") &&
			!strings.HasPrefix(entry, "TELIA_DIAGNOSTIC_TEST=") {
			env = append(env, entry)
		}
	}
	return append(env, "T_STD_PATH="+std, "XDG_CONFIG_HOME="+t.TempDir(), "TELIA_DIAGNOSTIC_TEST=1")
}

func runTestCompiler(t *testing.T, dir string, args ...string) (string, error) {
	t.Helper()
	cmd := exec.Command(
		os.Args[0],
		append([]string{"-test.run=^TestCompilerDiagnosticProcess$", "--"}, args...)...)
	cmd.Dir = dir
	cmd.Env = testCompilerEnv(t)
	out, err := cmd.CombinedOutput()
	return string(out), err
}

func TestBuildExcludesTestFilesIncludingImports(t *testing.T) {
	dir := writeTestPackage(t, map[string]string{
		"main.t":            "package main\nuse \"pkg::lib\"\nfn main() i32 {\nreturn lib::value()\n}\n",
		"broken_test.t":     "this is deliberately invalid Telia",
		"lib/lib.t":         "package lib\nfn value() i32 {\nreturn 0\n}\n",
		"lib/broken_test.t": "this is deliberately invalid Telia",
	})
	if out, err := runTestCompiler(t, dir, "build"); err != nil {
		t.Fatalf("normal build included tests: %v\n%s", err, out)
	}
	if out, err := exec.Command(filepath.Join(dir, "math")).CombinedOutput(); err != nil {
		t.Fatalf("run: %v\n%s", err, out)
	}
}
