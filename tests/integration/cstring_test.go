package integration

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/HicaroD/Telia/tests/compiler"
)

func TestCStringCanRepresentAndCheckNullPointers(t *testing.T) {
	path := filepath.Join(t.TempDir(), "main.t")
	source := `package main
use "std::ffi::c"
fn main() i32 {
  value c::CString := nil
  if value != nil {
    return 1
  }
  if value == nil {
    empty c::CString := ""
    if empty != nil {
      return 0
    }
  }
  return 2
}
`
	if err := os.WriteFile(path, []byte(source), 0644); err != nil {
		t.Fatal(err)
	}
	_, diags := compiler.CompileFile(path)
	if len(diags.Diags) != 0 {
		t.Fatalf("nullable C strings: %v", diags.Diags)
	}
}
