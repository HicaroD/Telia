package integration

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/HicaroD/Telia/internal/ctarget"
	"github.com/HicaroD/Telia/tests/compiler"
)

func TestLibcWideCharacterAndPointerIntegerScalarsMatchC(t *testing.T) {
	dir := t.TempDir()
	header := filepath.Join(dir, "scalars.h")
	if err := os.WriteFile(header, []byte(`#include <stddef.h>
#include <stdint.h>
static inline int check_scalars(wchar_t *wide, intptr_t *signed_pointer, uintptr_t *unsigned_pointer) {
    return *wide == 65 && *signed_pointer == -42 && *unsigned_pointer == 42;
}
`), 0600); err != nil {
		t.Fatal(err)
	}
	path := filepath.Join(dir, "main.t")
	source := fmt.Sprintf(`package main
use "std::libc"
@[header=%q, default_cc="c"]
extern {
  fn check_scalars(wide *libc::c_wchar_t, signed_pointer *libc::c_intptr_t, unsigned_pointer *libc::c_uintptr_t) i32
}
fn main() i32 {
  wide libc::c_wchar_t := 65
  signed_pointer libc::c_intptr_t := -42
  unsigned_pointer libc::c_uintptr_t := 42
  return check_scalars(&wide, &signed_pointer, &unsigned_pointer) - 1
}
`, header)
	if err := os.WriteFile(path, []byte(source), 0600); err != nil {
		t.Fatal(err)
	}
	_, diags := compiler.CompileFile(path)
	if len(diags.Diags) != 0 {
		t.Fatalf("wide and pointer-sized scalars: %v", diags.Diags)
	}
}

func TestCScalarABIChangingFlagsCannotSilentlyChangeSemantics(t *testing.T) {
	scalars, err := ctarget.Scalars()
	if err != nil {
		t.Fatal(err)
	}
	flag := "-funsigned-char"
	if scalars["__c_char"].Unsigned {
		flag = "-fsigned-char"
	}
	path := filepath.Join(t.TempDir(), "main.t")
	if err := os.WriteFile(path, []byte("package main\nuse \"std::libc\"\nfn main() {\nvalue libc::c_char := 0\n}\n"), 0600); err != nil {
		t.Fatal(err)
	}
	exe, diags := compiler.CompileOnlyWithCFlags(path, flag)
	if exe != "" {
		defer removeTestBinary(t, exe)
	}
	for _, d := range diags.Diags {
		if strings.Contains(d.Message, "C scalar ABI mismatch: char") {
			return
		}
	}
	t.Fatalf("want target ABI mismatch, got %v", diags.Diags)
}

func TestLibcScalarsUseTheSelectedCCompilersTypes(t *testing.T) {
	runRecordFixture(t, `package main
use "std::libc"
@[header=%q, default_cc="c"]
extern {
@[link_name="fixture_long_size"]
fn long_size() u64
@[link_name="fixture_long_increment"]
fn increment(value *libc::c_long) libc::c_long
@[link_name="fixture_plain_char"]
fn plain_char() const *libc::c_char
@[link_name="fixture_char_check"]
fn check(value const *libc::c_char) i32
}
fn main() i32 {
value libc::c_long := 41
if sizeof<libc::c_long>() != long_size() or increment(&value) != 42 {
return 1
}
if check(plain_char()) != 1 {
return 2
}
return 0
}
`)
}
