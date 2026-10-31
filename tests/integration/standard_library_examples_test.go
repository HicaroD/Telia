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
		{"utf8.t", "14 bytes, 7 code points\n72 at byte 0\n105 at byte 1\n44 at byte 2\n32 at byte 3\n19990 at byte 4\n30028 at byte 7\n128512 at byte 10\n"},
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
