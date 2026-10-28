package integration

import (
	"fmt"
	"os"
	"path/filepath"
	"testing"

	"github.com/HicaroD/Telia/tests/compiler"
)

func TestHeaderOwnedFunctionMacroUsesHeaderDeclaration(t *testing.T) {
	dir := t.TempDir()
	header := filepath.Join(dir, "macro.h")
	text := "static inline int increment(int value) { return value + 1; }\n#define fixture_increment(value) increment((value) + 1)\n"
	if err := os.WriteFile(header, []byte(text), 0600); err != nil {
		t.Fatal(err)
	}
	source := fmt.Sprintf(`package main
@[header=%q, default_cc="c"]
extern {
  @[link_name="fixture_increment"]
  fn increment(value i32) i32
}
fn main() i32 {
  return increment(40) - 42
}
`, header)
	path := filepath.Join(dir, "main.t")
	if err := os.WriteFile(path, []byte(source), 0600); err != nil {
		t.Fatal(err)
	}
	_, diags := compiler.CompileFile(path)
	if len(diags.Diags) != 0 {
		t.Fatalf("header macro call: %v", diags.Diags)
	}
}
