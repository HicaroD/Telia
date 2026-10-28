package integration

import (
	"runtime"
	"testing"

	"github.com/HicaroD/Telia/tests/compiler"
)

func TestVaListOwnerNativeABIUnderUBSan(t *testing.T) {
	if runtime.GOOS == "windows" {
		t.Skip(
			"Windows sanitizer runtime is not configured; ordinary native lifecycle fixtures still run",
		)
	}
	exe, diags := compileRecordFixture(t, `package main
use "std::libc"
type Consume = fn(marker i32, args @c ...) i32
@[header=%q, default_cc="c"]
extern {
@[link_name="fixture_va_owned_drive"]
fn drive(callback Consume) i32
}
fn consume(marker i32, args @c ...) i32 {
owner libc::VaList := libc::va_start(args)
defer libc::va_end(&owner)
copy := libc::va_copy(libc::va_borrow(&owner))
defer libc::va_end(&copy)
if libc::va_arg<libc::c_int>(&copy) != 42 or libc::va_arg<libc::c_double>(&copy) != 1.5 {
return 0
}
if libc::va_arg<libc::c_int>(&owner) != 42 or libc::va_arg<libc::c_double>(&owner) != 1.5 {
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
`, "-fsanitize=undefined", "-fno-sanitize-recover=undefined", "-O2")
	if len(diags.Diags) != 0 {
		t.Fatalf("compile sanitized va_list fixture: %v", diags.Diags)
	}
	defer removeTestBinary(t, exe)
	if output, err := compiler.RunBinary(exe); err != nil {
		t.Fatalf("sanitized va_list ABI: %v: %s", err, output)
	}
}
