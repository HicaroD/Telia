package integration

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/HicaroD/Telia/tests/compiler"
)

func TestBooleanResultsCanBeComparedForEquality(t *testing.T) {
	path := filepath.Join(t.TempDir(), "main.t")
	source := `package main
fn enabled() bool {
  return true
}
fn main() i32 {
  if enabled() == false {
    return 1
  }
  if enabled() != true {
    return 2
  }
  if enabled() == true {
    if enabled() != false {
      return 0
    }
  }
  return 3
}
`
	if err := os.WriteFile(path, []byte(source), 0644); err != nil {
		t.Fatal(err)
	}
	_, diags := compiler.CompileFile(path)
	if len(diags.Diags) != 0 {
		t.Fatalf("boolean equality: %v", diags.Diags)
	}
}
