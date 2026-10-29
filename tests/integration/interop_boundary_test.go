package integration

import (
	"strings"
	"testing"
)

func TestLibcDoesNotExportCInteropConstructs(t *testing.T) {
	for _, test := range []struct {
		name string
		body string
		want string
	}{
		{"scalar", "fn inspect(value libc::c_int) {}", "'c_int' type not found"},
		{"text", "fn inspect(value libc::CString) {}", "'CString' type not found"},
		{"const text", "fn inspect(value libc::ConstCString) {}", "'ConstCString' type not found"},
		{"borrowed cursor", "fn inspect(value libc::c_va_list) {}", "'c_va_list' type not found"},
		{"owner", "fn inspect(value libc::VaList) {}", "'VaList' type not found"},
		{"operation", "fn inspect(value c::c_va_list) {\nowner := libc::va_copy(value)\ndefer c::va_end(&owner)\n}", "function 'va_copy' not defined"},
	} {
		t.Run(test.name, func(t *testing.T) {
			exe, diags := compileRecordFixture(
				t,
				"package main\nuse \"std::libc\"\nuse \"std::ffi::c\"\n"+
					test.body+"\nfn main() {}\n",
			)
			if exe != "" {
				defer removeTestBinary(t, exe)
			}
			for _, diag := range diags.Diags {
				if strings.Contains(diag.Message, test.want) {
					return
				}
			}
			t.Fatalf("want libc to reject old interop export: %q; got %v", test.want, diags.Diags)
		})
	}
}

func TestCInteropPackageSuppliesTypesWithoutLibcImports(t *testing.T) {
	runRecordFixture(t, `package main
use "std::ffi::c"
@[header=%q]
extern {
@[link_name="fixture_text_check"]
fn check(text c::ConstCString) i32
}
fn main() i32 {
text c::CString := "hello"
size c::c_size_t := 5
if check(text) != 1 or size != 5 {
return 1
}
return 0
}
`)
}

func TestLibcProvidesEnvironmentAndIntegerConversionWithoutLocalExterns(t *testing.T) {
	t.Setenv("TELIA_TEST_LIBC_NUMBER", "42")
	runRecordFixture(t, `package main
use "std::libc"
use "std::ffi::c"
fn main() i32 {
configured c::CString := libc::getenv("TELIA_TEST_LIBC_NUMBER")
if configured == nil {
return 1
}
number c::c_int := libc::atoi(configured)
if number != 42 or libc::atoi(" -7") != -7 {
return 2
}
return 0
}
`)
}

func TestLibcExposesStandardPutsAlongsidePrintf(t *testing.T) {
	runRecordFixture(t, `package main
use "std::libc"
fn main() i32 {
if libc::puts("hello") < 0 {
return 1
}
if libc::printf("%d", 42) != 2 {
return 2
}
return 0
}
`)
}

func TestImportedCStringsPreserveForeignCallsAndNativeViews(t *testing.T) {
	runRecordFixture(t, `package main
use "std::ffi::c"
use "std::strings"
type Text = c::ConstCString
@[header=%q]
extern {
@[link_name="fixture_text_check"]
fn check(text Text) i32
@[link_name="fixture_text_replace"]
fn replace(output *Text)
}
fn main() i32 {
mutable c::CString := "hello"
text Text := mutable
if check(text) != 1 or strings::len(strings::from_c(text)) != 5 {
return 1
}
replace(&text)
if strings::from_c(text) != "after" {
return 2
}
empty c::CString := nil
if empty != nil {
return 3
}
return 0
}
`)
}

func TestCInteropPackageOwnsNativeVariadicLifecycle(t *testing.T) {
	runRecordFixture(t, `package main
use "std::ffi::c"
@[cc="c"]
type Consume = fn(marker i32, args @c ...) i32
@[header=%q]
extern {
@[link_name="fixture_va_owned_drive"]
fn drive(callback Consume) i32
}
@[cc="c"]
fn consume(marker i32, args @c ...) i32 {
owner c::VaList := c::va_start(args)
defer c::va_end(&owner)
copy := c::va_copy(c::va_borrow(&owner))
defer c::va_end(&copy)
if c::va_arg<c::c_int>(&copy) != 42 or c::va_arg<c::c_double>(&copy) != 1.5 {
return 0
}
if c::va_arg<c::c_int>(&owner) != 42 {
return 0
}
return 1
}
fn main() i32 {
if drive(consume) != 1 {
return 1
}
return 0
}
`)
}

func TestNativeAndForeignCallbacksHaveIndependentGenericSpecializations(t *testing.T) {
	runRecordFixture(t, `package main
type Native = fn(value i32) i32
@[cc="c"]
type Foreign = fn(value i32) i32
fn native(value i32) i32 {
return value + 1
}
@[cc="c"]
fn foreign(value i32) i32 {
return value + 2
}
fn identity<T>(value T) T {
return value
}
fn main() i32 {
first Native := identity<Native>(native)
second Foreign := identity<Foreign>(foreign)
if first(40) != 41 or second(40) != 42 {
return 1
}
return 0
}
`)
}
