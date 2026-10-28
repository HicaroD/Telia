package integration

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/HicaroD/Telia/internal/ctarget"
	"github.com/HicaroD/Telia/tests/compiler"
)

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
		defer os.RemoveAll(filepath.Dir(exe))
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
