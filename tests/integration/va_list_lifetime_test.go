package integration

import (
	"strings"
	"testing"

	"github.com/HicaroD/Telia/tests/compiler"
)

func TestVaListOwnersRequireCleanupOnEveryExecutedScopeExit(t *testing.T) {
	runVaLifetimeFailure(
		t,
		"owner := c::va_copy(args)\nreturn 1",
		"va_list owner left scope without va_end",
	)
}

func TestVaEndRejectsEndingTheSameOwnerTwice(t *testing.T) {
	runVaLifetimeFailure(
		t,
		"owner := c::va_copy(args)\nc::va_end(&owner)\nc::va_end(&owner)\nreturn 0",
		"va_list has already ended",
	)
}

func TestVaArgRejectsReadingAnEndedOwner(t *testing.T) {
	runVaLifetimeFailure(
		t,
		"owner := c::va_copy(args)\nc::va_end(&owner)\nvalue := c::va_arg<c::c_int>(&owner)\nreturn 0",
		"va_list is ended or consumed",
	)
}

func TestVaArgRejectsReusingAnOwnerConsumedByForeignCode(t *testing.T) {
	runVaLifetimeFailure(
		t,
		"owner := c::va_copy(args)\ndefer c::va_end(&owner)\ncheck(c::va_borrow(&owner))\nvalue := c::va_arg<c::c_int>(&owner)\nreturn 0",
		"va_list is ended or consumed",
	)
}

func runVaLifetimeFailure(t *testing.T, body, want string) {
	t.Helper()
	source := `package main
use "std::ffi::c"
@[cc="c"]
type Inspect = fn(args c::c_va_list) i32
@[header=%q, default_cc="c"]
extern {
@[link_name="fixture_va_drive"]
fn drive(callback Inspect) i32
@[link_name="fixture_va_check"]
fn check(args c::c_va_list) i32
}
@[cc="c"]
fn inspect(args c::c_va_list) i32 {
BODY
}
fn main() i32 {
return drive(inspect)
}
`
	exe, diags := compileRecordFixture(t, strings.ReplaceAll(source, "BODY", body))
	if len(diags.Diags) != 0 {
		t.Fatalf("compile lifetime fixture: %v", diags.Diags)
	}
	defer removeTestBinary(t, exe)
	output, err := compiler.RunBinary(exe)
	if err == nil || !strings.Contains(output, want) {
		t.Fatalf("want panic %q, got %v: %s", want, err, output)
	}
}

func TestVaBorrowCannotBeStoredOutsideItsCallExpression(t *testing.T) {
	exe, diags := compileRecordFixture(t, `package main
use "std::ffi::c"
fn inspect(args c::c_va_list) {
owner := c::va_copy(args)
defer c::va_end(&owner)
borrowed := c::va_borrow(&owner)
}
fn main() {}
`)
	if exe != "" {
		defer removeTestBinary(t, exe)
	}
	if len(diags.Diags) != 1 || diags.Diags[0].Title != "Borrowed va_list copy" ||
		!diags.Diags[0].Span.Valid() ||
		diags.Diags[0].Help == "" {
		t.Fatalf("want call-scoped borrow diagnostic, got %v", diags.Diags)
	}
}
