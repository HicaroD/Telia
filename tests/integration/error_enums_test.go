package integration

import (
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"testing"

	"github.com/HicaroD/Telia/tests/compiler"
)

func errorProgram(t *testing.T, source string) string {
	t.Helper()
	path := filepath.Join(t.TempDir(), "main.t")
	if err := os.WriteFile(path, []byte(source), 0600); err != nil {
		t.Fatal(err)
	}
	return path
}

func TestValueCatchRejectsFallthrough(t *testing.T) {
	_, diags := compiler.CompileOnly(errorProgram(t, `package main
fn read() (i32, error) {
return 42, nil
}

fn main() i32 {
value := read() @catch err {
if err != nil {
return 1
}
}
return value
}
`))
	for _, diag := range diags.Diags {
		if strings.Contains(diag.Message, "value-producing @catch") {
			return
		}
	}
	t.Fatalf("expected catch fallthrough diagnostic, got %v", diags.Diags)
}

func TestErrorConstructorAcceptsBorrowedMessage(t *testing.T) {
	output, diags := compiler.CompileFile(errorProgram(t, `package main
use "std::io"
fn failure(message string) error {
return error(message)
}
fn main() i32 {
err := failure("dynamic")
io::println(err.msg)
return 0
}
`))
	if len(diags.Diags) != 0 || output != "dynamic\n" {
		t.Fatalf("want borrowed message, got %q, %v", output, diags.Diags)
	}
}

func TestErrorEnumIdentitySurvivesReturn(t *testing.T) {
	_, diags := compiler.CompileFile(errorProgram(t, `package main
enum FileError: error {
NotFound
}
fn failure() error {
return error(FileError::NotFound, "missing")
}
fn main() i32 {
err := failure()
if err == nil {
return 1
}
copy := err
found := match copy {
FileError::NotFound => true
_ => false
}
if found == false {
return 2
}
return 0
}
`))
	if len(diags.Diags) != 0 {
		t.Fatalf("want classified error after return and assignment, got %v", diags.Diags)
	}
}

func TestErrorEnumPayloadSurvivesTupleAndCatch(t *testing.T) {
	output, diags := compiler.CompileFile(errorProgram(t, `package main
use "std::io"
enum FileError: error {
ReadFailed {
path string
code i32
}
}
enum NetworkError: error {
ReadFailed {
code i64
}
}
fn identity<T>(value T) T {
return value
}
fn failure() (i32, error) {
err := error(FileError::ReadFailed{path: "file.txt", code: 42}, "read failed")
return 0, identity<error>(err)
}
fn main() i32 {
value := failure() @catch err {
description := match err {
NetworkError::ReadFailed{code} => "wrong enum"
FileError::ReadFailed{path, code} => {
if code != 42 {
return 2
}
yield path
}
_ => "wrong fallback"
}
io::println(description)
io::println(err.msg)
return 0
}
return value
}
`))
	if len(diags.Diags) != 0 || output != "file.txt\nread failed\n" {
		t.Fatalf("want preserved payload through tuple and catch, got %q, %v", output, diags.Diags)
	}
}

func TestErrorKindIgnoresPayloadAndMessage(t *testing.T) {
	_, diags := compiler.CompileFile(errorProgram(t, `package main
enum FileError: error {
ReadFailed {
code i32
}
}
enum OtherError: error {
ReadFailed
}
fn main() i32 {
first := error(FileError::ReadFailed{code: 1}, "first")
second := error(FileError::ReadFailed{code: 2}, "second")
if first.kind != second.kind {
return 1
}
if first.kind != FileError::ReadFailed {
return 2
}
if first.kind == OtherError::ReadFailed {
return 3
}
plain := error("")
success error := nil
if plain.kind != nil or success.kind != nil {
return 4
}
if plain == nil {
return 5
}
return 0
}
`))
	if len(diags.Diags) != 0 {
		t.Fatalf("want payload-independent error identity, got %v", diags.Diags)
	}
}

func TestErrorKindCanBeSavedButNotChanged(t *testing.T) {
	_, diags := compiler.CompileFile(errorProgram(t, `package main
enum Failure: error {
Bad {
code i32
}
}
fn main() i32 {
err := error(Failure::Bad{code: 7}, "")
kind := err.kind
if kind != Failure::Bad {
return 1
}
return 0
}
`))
	if len(diags.Diags) != 0 {
		t.Fatalf("want saved identity, got %v", diags.Diags)
	}
	_, diags = compiler.CompileOnly(errorProgram(t, `package main
fn main() i32 {
err := error("failure")
err.kind = nil
return 0
}
`))
	if len(diags.Diags) == 0 {
		t.Fatal("error identity must be read-only")
	}
}

func TestErrorConstructorsAndPatternsHaveDiagnostics(t *testing.T) {
	for _, tc := range []struct{ name, body, want string }{
		{"no arguments", `err := error()`, "error constructor requires"},
		{"extra arguments", `err := error("a", "b", "c")`, "error constructor requires"},
		{"missing comma", `err := error("a" "b")`, "expected , or )"},
		{"numeric message", `err := error(42)`, ""},
		{"nil message", `err := error(nil)`, ""},
		{"ordinary enum", `err := error(Ordinary::Bad, "bad")`, "declared with : error"},
		{"implicit conversion", `err error := Failure::Bad{code: 1}`, ""},
		{"missing field", `err := error(Failure::Bad{}, "bad")`, "missing initializer"},
		{"wrong field type", `err := error(Failure::Bad{code: "bad"}, "bad")`, ""},
		{"missing fallback", "err := error(\"bad\")\nresult := match err {\nFailure::Bad{code} => code\n}", "requires a final fallback"},
		{"ordinary pattern", "err := error(\"bad\")\nresult := match err {\nOrdinary::Bad => 1\n_ => 0\n}", "declared with : error"},
		{"missing pattern field", "err := error(\"bad\")\nresult := match err {\nFailure::Bad => 1\n_ => 0\n}", "omits payload fields"},
		{"unknown pattern field", "err := error(\"bad\")\nresult := match err {\nFailure::Bad{missing} => 1\n_ => 0\n}", "undeclared payload field"},
		{"duplicate pattern", "err := error(\"bad\")\nresult := match err {\nFailure::Bad{code} => code\nFailure::Bad{code} => code\n_ => 0\n}", "already been matched"},
	} {
		t.Run(tc.name, func(t *testing.T) {
			source := "package main\nenum Ordinary {\nBad\n}\nenum Failure: error {\nBad {\ncode i32\n}\n}\nfn main() i32 {\n" + tc.body + "\nreturn 0\n}\n"
			_, diags := compiler.CompileOnly(errorProgram(t, source))
			if len(diags.Diags) == 0 {
				t.Fatal("expected compiler diagnostic")
			}
			if tc.want != "" {
				for _, diag := range diags.Diags {
					if strings.Contains(diag.Message, tc.want) {
						return
					}
				}
				t.Fatalf("want diagnostic containing %q, got %v", tc.want, diags.Diags)
			}
		})
	}
}

func TestErrorInputsEvaluateOnceAtTheirExpressionSite(t *testing.T) {
	_, diags := compiler.CompileFile(errorProgram(t, `package main
enum Failure: error {
Bad {
first i32
second i32
}
}
fn next(count *i32) i32 {
*count = *count + 1
return *count
}
fn message(count *i32) string {
*count = *count + 1
return "message"
}
fn counted(err error, count *i32) error {
*count = *count + 1
return err
}
fn main() i32 {
count i32 := 0
skipped := false and (error(Failure::Bad{first: next(&count), second: next(&count)}, message(&count)) != nil)
if count != 0 {
return 1
}
err := error(Failure::Bad{first: next(&count), second: next(&count)}, message(&count))
if count != 3 {
return 2
}
result := match counted(err, &count) {
Failure::Bad{first, second} => first * 10 + second
_ => 0
}
if count != 4 or result != 12 {
return 3
}
empty string := ""
present := error(empty)
if present == nil {
return 4
}
return 0
}
`))
	if len(diags.Diags) != 0 {
		t.Fatalf("want single evaluation in source order, got %v", diags.Diags)
	}
}

func TestErrorInlinePayloadRejectsRecursiveStorage(t *testing.T) {
	_, diags := compiler.CompileOnly(errorProgram(t, `package main
struct Context {
cause error
}
enum Failure: error {
Bad {
context Context
}
}
fn main() i32 {
return 0
}
`))
	for _, diag := range diags.Diags {
		if strings.Contains(diag.Message, "infinitely sized") {
			return
		}
	}
	t.Fatalf("want recursive error storage diagnostic, got %v", diags.Diags)
}

func TestStandardLibraryErrorsHaveStableIdentities(t *testing.T) {
	_, diags := compiler.CompileFile(errorProgram(t, `package main
use "std::errors"
use "std::strings"
use "std::memory"
use "std::slices"
use "std::unicode::utf8"
fn main() i32 {
text, rangeErr := strings::slice("abc", 3, 4)
if rangeErr.kind != errors::Error::InvalidRange {
return 1
}
ctext, zeroErr := strings::to_c("a\0b")
if zeroErr.kind != errors::Error::EmbeddedZero {
return 2
}
arena := memory::arena_create() @fail
defer memory::arena_destroy(&arena)
ptr, alignErr := memory::arena_alloc(&arena, 1, 3)
if alignErr.kind != errors::Error::InvalidAlignment {
return 3
}
parts, separatorErr := strings::split("abc", "")
if separatorErr.kind != errors::Error::EmptySeparator {
return 4
}
values i32[1] := .{7}
item, indexErr := slices::get<i32>(values[:], 1)
if indexErr.kind != errors::Error::OutOfBounds {
return 5
}
bad := strings::slice("é", 0, 1) @fail
point, width, utfErr := utf8::decode(bad, 0)
if utfErr.kind != errors::Error::InvalidUTF8 {
return 6
}
return 0
}
`))
	if len(diags.Diags) != 0 {
		t.Fatalf("want documented standard-library identities, got %v", diags.Diags)
	}
}

func TestErrorPayloadLifetimesWithNativeSanitizers(t *testing.T) {
	if runtime.GOOS == "windows" {
		t.Skip("MinGW-w64 does not provide the required sanitizer runtime")
	}
	sanitizers := os.Getenv("TELIA_TEST_SANITIZERS")
	if sanitizers == "" {
		sanitizers = "undefined"
	}
	runStandardLibraryFixture(
		t,
		"error_enum_lifetimes.t",
		"-fsanitize="+sanitizers,
		"-fno-sanitize-recover=all",
		"-fno-omit-frame-pointer",
	)
}

func TestValueCatchAcceptsReturnInNestedBlock(t *testing.T) {
	_, diags := compiler.CompileFile(errorProgram(t, `package main
fn failure() (i32, error) {
return 0, error("bad")
}
fn main() i32 {
value := failure() @catch err {
{
return 0
}
}
return 1
}
`))
	if len(diags.Diags) != 0 {
		t.Fatalf("nested return must exit catch, got %v", diags.Diags)
	}
}

func TestErrorConstructionDoesNotAllocate(t *testing.T) {
	header, err := filepath.Abs("testdata/allocation_control.h")
	if err != nil {
		t.Fatal(err)
	}
	runStandardLibraryFixture(t, "error_enum_allocation_failure.t", "-include", header)
}

func TestImportedErrorEnumsWithSameNamesRemainDistinct(t *testing.T) {
	dir := t.TempDir()
	for _, name := range []string{"first", "second"} {
		if err := os.Mkdir(filepath.Join(dir, name), 0700); err != nil {
			t.Fatal(err)
		}
		source := "package " + name + "\nenum Failure: error {\nBad {\ncode i32\n}\n}\n"
		if err := os.WriteFile(filepath.Join(dir, name, "errors.t"), []byte(source), 0600); err != nil {
			t.Fatal(err)
		}
	}
	source := `package main
use "pkg::first"
use "pkg::second"
enum GenericFailure<T>: error {
Bad {
value T
}
}
fn main() i32 {
one := error(first::Failure::Bad{code: 1}, "same")
two := error(second::Failure::Bad{code: 1}, "same")
if one.kind == two.kind or one.kind != first::Failure::Bad {
return 1
}
value := match two {
first::Failure::Bad{code} => 0
second::Failure::Bad{code} => code
_ => 0
}
if value != 1 {
return 2
}
generic := error(GenericFailure<i32>::Bad{value: 42}, "generic")
result := match generic {
GenericFailure<i64>::Bad{value} => 0
GenericFailure<i32>::Bad{value} => value
_ => 0
}
if result != 42 {
return 3
}
return 0
}
`
	if err := os.WriteFile(filepath.Join(dir, "main.t"), []byte(source), 0600); err != nil {
		t.Fatal(err)
	}
	_, diags := compiler.CompilePackage(dir)
	if len(diags.Diags) != 0 {
		t.Fatalf("want nominal identities across packages and specializations, got %v", diags.Diags)
	}
}

func TestErrorConstructorAllowsMultilineArguments(t *testing.T) {
	_, diags := compiler.CompileFile(errorProgram(t, `package main
enum Failure: error {
Bad {
code i32
}
}
fn main() i32 {
err := error(
Failure::Bad{code: 42},
"message",
)
if err.kind != Failure::Bad {
return 1
}
return 0
}
`))
	if len(diags.Diags) != 0 {
		t.Fatalf("want multiline error construction, got %v", diags.Diags)
	}
}

func TestErrorFallbackHandlesUnclassifiedErrorsAndSuccess(t *testing.T) {
	_, diags := compiler.CompileFile(errorProgram(t, `package main
enum Failure: error {
Bad
}
fn fallback(err error) bool {
return match err {
Failure::Bad => false
_ => true
}
}
fn empty_failure(message string) error {
return error(message)
}
fn main() i32 {
if fallback(nil) == false or fallback(error("plain")) == false {
return 1
}
empty string := ""
empty_failure(empty) @catch err {
if err == nil or err.msg != "" {
return 2
}
return 0
}
return 3
}
`))
	if len(diags.Diags) != 0 {
		t.Fatalf("want fallback and empty-message catch behavior, got %v", diags.Diags)
	}
}
