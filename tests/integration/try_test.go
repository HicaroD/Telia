package integration

import (
	"strings"
	"testing"

	"github.com/HicaroD/Telia/tests/compiler"
)

func TestTrySuccessContinuesForErrorOnlyAndTupleCalls(t *testing.T) {
	path := loopControlSource(t, `package main
use "std::io"
fn okay() error {
return nil
}
fn value<T>(item T) (T, error) {
return item, nil
}
fn run() error {
defer io::println("cleanup")
okay() @try
text := value<string>("continued") @try
io::println(text)
return nil
}
fn main() i32 {
run() @fail
return 0
}
`)
	output, diags := compiler.CompileFile(path)
	if len(diags.Diags) != 0 {
		t.Fatalf("successful @try should compile: %v", diags.Diags)
	}
	if want := "continued\ncleanup\n"; output != want {
		t.Fatalf("want %q, got %q", want, output)
	}
}

func TestTryRejectsCallsOrEnclosingFunctionsWithoutErrors(t *testing.T) {
	for _, body := range []string{
		"fn value() i32 {\nreturn 1\n}\nfn run() error {\nx := value() @try\nreturn nil\n}\n",
		"fn value() error {\nreturn nil\n}\nfn run() i32 {\nvalue() @try\nreturn 0\n}\n",
	} {
		path := loopControlSource(t, "package main\n"+body+"fn main() i32 {\nreturn 0\n}\n")
		_, diags := compiler.CompileOnly(path)
		if len(diags.Diags) == 0 || !strings.Contains(diags.Diags[0].Message, "@try requires") {
			t.Fatalf("want propagation contract diagnostic, got %v", diags.Diags)
		}
	}
}

func TestTryPropagatesErrorsThroughDeferredCleanup(t *testing.T) {
	path := loopControlSource(t, `package main
use "std::io"
fn zero() i32 {
return 99
}
fn failure() (i32, error) {
return 0, error("failure")
}
fn worker() (string, error) {
defer io::println("worker first")
defer io::println("worker second")
{
defer io::println("inner")
value := failure() @try
}
return "unexpected", nil
}
fn run() error {
defer io::println("caller")
value := worker() @try
io::println(value)
return nil
}
fn main() i32 {
run() @catch err {
io::println(err.msg)
return 0
}
return 1
}
`)
	output, diags := compiler.CompileFile(path)
	if len(diags.Diags) != 0 {
		t.Fatalf("@try should compile: %v", diags.Diags)
	}
	if want := "inner\nworker second\nworker first\ncaller\nfailure\n"; output != want {
		t.Fatalf("want %q, got %q", want, output)
	}
}

func TestTryComposesInsideArgumentsArraysAndShortCircuitExpressions(t *testing.T) {
	path := loopControlSource(t, `package main
use "std::io"
fn value(fail bool) (i32, error) {
io::println("called")
if fail {
return 0, error("failure")
}
return 42, nil
}
fn same(value i32) i32 {
return value
}
fn check() error {
defer io::println("cleanup")
values i32[1] := .{same(value(false) @try)}
assigned i32 := 0
assigned = value(false) @try
values[0] = value(false) @try
if assigned != 42 {
return error("assignment lost")
}
if false and (value(true) @try) == 0 {
return error("short circuit failed")
}
if values[0] != 42 {
return error("array value lost")
}
if true and (value(true) @try) == 0 {
return error("propagation failed")
}
return nil
}
fn main() i32 {
check() @catch err {
io::println(err.msg)
return 0
}
return 1
}
`)
	output, diags := compiler.CompileFile(path)
	if len(diags.Diags) != 0 {
		t.Fatalf("composed @try should compile: %v", diags.Diags)
	}
	if want := "called\ncalled\ncalled\ncalled\ncleanup\nfailure\n"; output != want {
		t.Fatalf("want %q, got %q", want, output)
	}
}

func TestTryUnwrapsMultipleSuccessValuesAndPreservesErrors(t *testing.T) {
	path := loopControlSource(t, `package main
use "std::io"
fn pair(fail bool) (i32, string, error) {
if fail {
return 0, "", error("pair failed")
}
return 42, "pair", nil
}
fn run() error {
defer io::println("cleanup")
number, text := pair(false) @try
if number != 42 {
return error("number lost")
}
io::println(text)
callback := pair
callback(true) @try
io::println("unreachable")
return nil
}
fn main() i32 {
run() @catch err {
io::println(err.msg)
return 0
}
return 1
}
`)
	output, diags := compiler.CompileFile(path)
	if len(diags.Diags) != 0 {
		t.Fatalf("multi-result @try should compile: %v", diags.Diags)
	}
	if want := "pair\ncleanup\npair failed\n"; output != want {
		t.Fatalf("want %q, got %q", want, output)
	}
}
