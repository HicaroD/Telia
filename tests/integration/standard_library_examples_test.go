package integration

import (
	"context"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"testing"
	"time"

	"github.com/HicaroD/Telia/tests/compiler"
)

func runStandardLibraryFixture(t *testing.T, fixture string, flags ...string) {
	t.Helper()
	if !filepath.IsAbs(fixture) {
		fixture = filepath.Join("testdata", fixture)
	}
	for _, opt := range []string{"-O0", "-O3"} {
		t.Run(opt, func(t *testing.T) {
			cflags := append([]string{opt}, flags...)
			path, diags := compiler.CompileOnlyWithCFlags(fixture, cflags...)
			if len(diags.Diags) != 0 {
				t.Fatalf("compile: %v", diags.Diags)
			}
			defer removeTestBinary(t, path)
			ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
			defer cancel()
			output, err := exec.CommandContext(ctx, path).CombinedOutput()
			if err != nil {
				t.Fatalf("execution: %v: %s", err, output)
			}
		})
	}
}

func TestExamplesUsePublishedStandardLibraryAPIs(t *testing.T) {
	for _, test := range []struct {
		file      string
		want      string
		libcLines int // Leading lines printed by libc, which uses CRLF on Windows.
	}{
		{"floats.t", "3.000000\n", 1},
		{"error_enums.t", "settings.txt\nfile does not exist\n", 0},
		{"hello2.t", "Hello, world\n", 0},
		{"string_utilities.t", "世界 at byte 3\né\n世界\n😀\n", 1},
		{"utf8.t", "14 bytes, 7 code points\n72 at byte 0\n105 at byte 1\n44 at byte 2\n32 at byte 3\n19990 at byte 4\n30028 at byte 7\n128512 at byte 10\n", 8},
	} {
		t.Run(test.file, func(t *testing.T) {
			output, diags := compiler.CompileFile(filepath.Join("..", "..", "examples", test.file))
			if len(diags.Diags) != 0 {
				t.Fatalf("compile example: %v", diags.Diags)
			}
			want := test.want
			if runtime.GOOS == "windows" {
				want = strings.Replace(want, "\n", "\r\n", test.libcLines)
			}
			if output != want {
				t.Fatalf("want %q, got %q", want, output)
			}
		})
	}
}
