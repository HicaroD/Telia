package integration

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/HicaroD/Telia/tests/compiler"
)

func TestBitwiseOrCombinesIntegerFlags(t *testing.T) {
	path := filepath.Join(t.TempDir(), "main.t")
	if err := os.WriteFile(path, []byte("package main\nfn main() i32 {\na u64 := 8\nb u64 := 32\nif a | b != 40 {\nreturn 1\n}\nreturn 0\n}\n"), 0600); err != nil {
		t.Fatal(err)
	}
	_, diags := compiler.CompileFile(path)
	if len(diags.Diags) != 0 {
		t.Fatalf("bitwise OR must combine flags before comparison: %v", diags.Diags)
	}
}

func TestBitwiseAndAndXorHaveDistinctPrecedence(t *testing.T) {
	path := filepath.Join(t.TempDir(), "main.t")
	if err := os.WriteFile(path, []byte("package main\nfn main() i32 {\na u32 := 1\nb u32 := 2\nc u32 := 3\nif a | b ^ c & a != 3 {\nreturn 1\n}\nreturn 0\n}\n"), 0600); err != nil {
		t.Fatal(err)
	}
	_, diags := compiler.CompileFile(path)
	if len(diags.Diags) != 0 {
		t.Fatalf("expected AND before XOR before OR before comparison: %v", diags.Diags)
	}
}

func TestBitwiseNotPreservesNarrowUnsignedWidth(t *testing.T) {
	path := filepath.Join(t.TempDir(), "main.t")
	if err := os.WriteFile(path, []byte("package main\nfn main() i32 {\na u8 := 0\nif ~a != 255 {\nreturn 1\n}\nreturn 0\n}\n"), 0600); err != nil {
		t.Fatal(err)
	}
	_, diags := compiler.CompileFile(path)
	if len(diags.Diags) != 0 {
		t.Fatalf("NOT must complement only the operand's bits: %v", diags.Diags)
	}
}

func TestLogicalAndBindsMoreTightlyThanOr(t *testing.T) {
	path := filepath.Join(t.TempDir(), "main.t")
	if err := os.WriteFile(path, []byte("package main\nfn main() i32 {\nif true or false and false {\nreturn 0\n}\nreturn 1\n}\n"), 0600); err != nil {
		t.Fatal(err)
	}
	_, diags := compiler.CompileFile(path)
	if len(diags.Diags) != 0 {
		t.Fatalf("logical precedence must follow the agreed operator table: %v", diags.Diags)
	}
}
