package integration

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"testing"
	"time"

	"github.com/HicaroD/Telia/tests/compiler"
)

func TestUTF8Validation(t *testing.T) {
	runUTF8Fixture(t, "utf8_validation.t")
}

func TestUTF8DecodingAndCounting(t *testing.T) {
	runUTF8Fixture(t, "utf8_decode.t")
}

func TestUTF8BorrowedIteration(t *testing.T) {
	runUTF8Fixture(t, "utf8_iteration.t")
}

func TestUTF8WithNativeSanitizers(t *testing.T) {
	if runtime.GOOS == "windows" {
		t.Skip("MinGW-w64 does not provide the required sanitizer runtime")
	}
	sanitizers := os.Getenv("TELIA_TEST_SANITIZERS")
	if sanitizers == "" {
		sanitizers = "undefined"
	}
	for _, fixture := range []string{"utf8_validation.t", "utf8_decode.t", "utf8_iteration.t"} {
		t.Run(fixture, func(t *testing.T) {
			runUTF8Fixture(
				t,
				fixture,
				"-fsanitize="+sanitizers,
				"-fno-sanitize-recover=all",
				"-fno-omit-frame-pointer",
			)
		})
	}
}

// Keep raw bytes in literals: Telia strings are byte views, not C strings.
func utf8Literal(text string) string {
	r := strings.NewReplacer("\\", "\\\\", "\"", "\\\"", "\x00", "\\0", "\n", "\\n", "\t", "\\t")
	return "\"" + r.Replace(text) + "\""
}

func TestUTF8StrictValidation(t *testing.T) {
	cases := []struct {
		text  string
		valid bool
	}{
		{
			"é世界😀",
			true,
		}, {"\x7f\xc2\x80\xdf\xbf\xe0\xa0\x80\xed\x9f\xbf\xee\x80\x80\xef\xbf\xbf\xf0\x90\x80\x80\xf4\x8f\xbf\xbf", true},
		{"\x80", false}, {"\xbf", false}, {"\xc0\x80", false}, {"\xc1\xbf", false},
		{"\xc2", false}, {"\xc2A", false}, {"\xe0\x80\x80", false},
		{"\xed\xa0\x80", false}, {"\xed\xbf\xbf", false}, {"\xe1\x80", false},
		{"\xe1A\x80", false}, {"\xe1\x80A", false}, {"\xf0\x80\x80\x80", false},
		{"\xf4\x90\x80\x80", false}, {"\xf5\x80\x80\x80", false}, {"\xff", false},
		{"\xf0\x90\x80", false}, {"\xf0A\x80\x80", false}, {"\xf0\x90A\x80", false},
		{"\xf0\x90\x80A", false}, {"valid\x00\xff", false},
		{"\xe0", false}, {"\xf0", false}, {"\xf0\x90", false}, {"\xfe", false},
	}
	var source strings.Builder
	source.WriteString("package main\nuse \"std::unicode::utf8\"\nfn main() i32 {\n")
	for i, tc := range cases {
		op := "=="
		if tc.valid {
			op = "!="
		}
		fmt.Fprintf(
			&source,
			"if utf8::validate(%s) %s nil {\nreturn %d\n}\n",
			utf8Literal(tc.text),
			op,
			i+1,
		)
		if !tc.valid {
			fmt.Fprintf(
				&source,
				"count%d, err%d := utf8::count(%s)\nif count%d != 0 or err%d == nil {\nreturn %d\n}\n",
				i,
				i,
				utf8Literal(tc.text),
				i,
				i,
				i+1,
			)
			fmt.Fprintf(
				&source,
				"iter%d, iter_err%d := utf8::iterator(%s)\nif iter_err%d == nil {\nreturn %d\n}\n",
				i,
				i,
				utf8Literal(tc.text),
				i,
				i+1,
			)
			if !strings.HasPrefix(tc.text, "valid") {
				fmt.Fprintf(
					&source,
					"point%d, width%d, decode_err%d := utf8::decode(%s, 0)\nif point%d != 0 or width%d != 0 or decode_err%d == nil {\nreturn %d\n}\n",
					i,
					i,
					i,
					utf8Literal(tc.text),
					i,
					i,
					i,
					i+1,
				)
			}
		}
	}
	source.WriteString("return 0\n}\n")
	file := filepath.Join(t.TempDir(), "validation.t")
	if err := os.WriteFile(file, []byte(source.String()), 0600); err != nil {
		t.Fatal(err)
	}
	runUTF8Fixture(t, file)
	if runtime.GOOS != "windows" {
		t.Run("undefined", func(t *testing.T) {
			runUTF8Fixture(t, file, "-fsanitize=undefined", "-fno-sanitize-recover=all")
		})
	}
}

func TestUTF8ScalarBoundaries(t *testing.T) {
	var source strings.Builder
	source.WriteString("package main\nuse \"std::unicode::utf8\"\nfn main() i32 {\n")
	for i, tc := range []struct {
		text  string
		point uint32
		width int
	}{
		{"\x00", 0, 1}, {"\x7f", 127, 1}, {"\xc2\x80", 128, 2}, {"\xdf\xbf", 2047, 2},
		{"\xe0\xa0\x80", 2048, 3}, {"\xed\x9f\xbf", 55295, 3}, {"\xee\x80\x80", 57344, 3},
		{"\xef\xbf\xbf", 65535, 3}, {"\xf0\x90\x80\x80", 65536, 4}, {"\xf4\x8f\xbf\xbf", 1114111, 4},
	} {
		fmt.Fprintf(
			&source,
			"point%d, width%d, err%d := utf8::decode(%s, 0)\nif point%d != %d or width%d != %d or err%d != nil {\nreturn %d\n}\n",
			i,
			i,
			i,
			utf8Literal(tc.text),
			i,
			tc.point,
			i,
			tc.width,
			i,
			i+1,
		)
	}
	source.WriteString("return 0\n}\n")
	file := filepath.Join(t.TempDir(), "boundaries.t")
	if err := os.WriteFile(file, []byte(source.String()), 0600); err != nil {
		t.Fatal(err)
	}
	runUTF8Fixture(t, file)
}

func runUTF8Fixture(t *testing.T, fixture string, flags ...string) {
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
