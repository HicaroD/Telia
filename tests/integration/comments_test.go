package integration

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/HicaroD/Telia/tests/compiler"
)

func TestCommentsCompile(t *testing.T) {
	path := filepath.Join(t.TempDir(), "main.t")
	src := `/* File header */
package main // package comment
fn /* declaration */ main() i32 {
value := 8 / /* operand */ 2 // trailing comment
/* Ignored syntax: ! " /* non-nested */
other := 2 /* A block newline separates statements.
*/ result := value / other
if result != 2 {
return 1
}
return 0
} // EOF comment`
	if err := os.WriteFile(path, []byte(src), 0600); err != nil {
		t.Fatal(err)
	}
	_, diags := compiler.CompileFile(path)
	if len(diags.Diags) != 0 {
		t.Fatalf("comments must be ignored during compilation: %v", diags.Diags)
	}
}
