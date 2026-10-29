package integration

import (
	"path/filepath"
	"testing"

	"github.com/HicaroD/Telia/tests/compiler"
)

func TestExamplesUsePublishedStandardLibraryAPIs(t *testing.T) {
	for _, test := range []struct {
		file string
		want string
	}{
		{"floats.t", "3.000000\n"},
		{"hello2.t", "Hello, world\n"},
	} {
		t.Run(test.file, func(t *testing.T) {
			output, diags := compiler.CompileFile(filepath.Join("..", "..", "examples", test.file))
			if len(diags.Diags) != 0 {
				t.Fatalf("compile example: %v", diags.Diags)
			}
			if output != test.want {
				t.Fatalf("want %q, got %q", test.want, output)
			}
		})
	}
}
