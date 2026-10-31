package integration

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"testing"
)

func runStrconvSource(t *testing.T, source string, flags ...string) {
	t.Helper()
	file := filepath.Join(t.TempDir(), "strconv.t")
	if err := os.WriteFile(file, []byte(source), 0600); err != nil {
		t.Fatal(err)
	}
	runStandardLibraryFixture(t, file, flags...)
	if runtime.GOOS != "windows" {
		t.Run("sanitized", func(t *testing.T) {
			sanitizers := os.Getenv("TELIA_TEST_SANITIZERS")
			if sanitizers == "" {
				sanitizers = "undefined"
			}
			cflags := append(append([]string{}, flags...), "-fsanitize="+sanitizers,
				"-fno-sanitize-recover=all", "-fno-omit-frame-pointer")
			runStandardLibraryFixture(t, file, cflags...)
		})
	}
}

func TestStrconvOwnershipAndAllocationFailures(t *testing.T) {
	header, err := filepath.Abs("testdata/allocation_control.h")
	if err != nil {
		t.Fatal(err)
	}
	source, err := os.ReadFile("testdata/strconv_ownership.t")
	if err != nil {
		t.Fatal(err)
	}
	runStrconvSource(t, string(source), "-include", header)
}

func TestStrconvExampleCleanup(t *testing.T) {
	header, err := filepath.Abs("testdata/allocation_control.h")
	if err != nil {
		t.Fatal(err)
	}
	example, err := filepath.Abs(filepath.Join("..", "..", "examples", "strconv.t"))
	if err != nil {
		t.Fatal(err)
	}
	runStandardLibraryFixture(t, example, "-include", header)
}

func TestStrconvIntegerBoundaries(t *testing.T) {
	cases := []struct {
		kind, text, want, errorKind string
	}{
		{"u64", "0", "0", ""}, {"u64", "00042", "42", ""},
		{"u64", "18446744073709551615", "18446744073709551615", ""},
		{"u64", "18446744073709551616", "0", "ConversionOutOfRange"},
		{"i64", "42", "42", ""}, {"i64", "+00042", "42", ""},
		{"i64", "-42", "-42", ""}, {"i64", "-0", "0", ""},
		{"i64", "9223372036854775807", "9223372036854775807", ""},
		{"i64", "-9223372036854775808", "(-9223372036854775807 - 1)", ""},
		{"i64", "9223372036854775808", "0", "ConversionOutOfRange"},
		{"i64", "-9223372036854775809", "0", "ConversionOutOfRange"},
	}
	for _, kind := range []string{"i64", "u64"} {
		for _, text := range []string{"", "+", "-", " 1", "1 ", "1\t", "\n1", "1x", "1.0", "1_0", "0x10", "１２", "١", "1\x00", "\x00", "\xff", "--1", "+-1", "\u00a01"} {
			cases = append(
				cases,
				struct{ kind, text, want, errorKind string }{kind, text, "0", "InvalidNumber"},
			)
		}
		cases = append(
			cases,
			struct{ kind, text, want, errorKind string }{
				kind,
				strings.Repeat("9", 2000),
				"0",
				"ConversionOutOfRange",
			},
			struct{ kind, text, want, errorKind string }{
				kind,
				strings.Repeat("0", 2000) + "7",
				"7",
				"",
			},
		)
	}
	for _, text := range []string{"+1", "-1", "+0", "-0"} {
		cases = append(
			cases,
			struct{ kind, text, want, errorKind string }{"u64", text, "0", "InvalidNumber"},
		)
	}
	var source strings.Builder
	source.WriteString("package main\nuse \"std::strconv\"\nuse \"std::errors\"\nfn main() i32 {\n")
	for i, tc := range cases {
		fmt.Fprintf(
			&source,
			"v%d, e%d := strconv::parse_%s(%s)\nif v%d != %s {\nreturn %d\n}\n",
			i,
			i,
			tc.kind,
			utf8Literal(tc.text),
			i,
			tc.want,
			i+1,
		)
		if tc.errorKind == "" {
			fmt.Fprintf(&source, "if e%d != nil {\nreturn %d\n}\n", i, i+1)
		} else {
			fmt.Fprintf(&source, "if e%d.kind != errors::Error::%s {\nreturn %d\n}\n", i, tc.errorKind, i+1)
		}
	}
	source.WriteString("return 0\n}\n")
	runStrconvSource(t, source.String())
}

func TestStrconvUnsignedParsing(t *testing.T) {
	runStrconvSource(t, `package main
use "std::strconv"
fn main() i32 {
    value, err := strconv::parse_u64("42")
    if err != nil or value != 42 {
        return 1
    }
    return 0
}
`)
}

func TestStrconvBooleanParsing(t *testing.T) {
	var source strings.Builder
	source.WriteString("package main\nuse \"std::strconv\"\nuse \"std::errors\"\nfn main() i32 {\n")
	for i, text := range []string{"true", "false", "", "True", "FALSE", "1", "0", " true", "false ", "true\x00", "truex", "\xff"} {
		fmt.Fprintf(
			&source,
			"v%d, e%d := strconv::parse_bool(%s)\nif v%d != %t {\nreturn %d\n}\n",
			i,
			i,
			utf8Literal(text),
			i,
			text == "true",
			i+1,
		)
		if i < 2 {
			fmt.Fprintf(&source, "if e%d != nil {\nreturn %d\n}\n", i, i+1)
		} else {
			fmt.Fprintf(&source, "if e%d.kind != errors::Error::InvalidBoolean {\nreturn %d\n}\n", i, i+1)
		}
	}
	source.WriteString("return 0\n}\n")
	runStrconvSource(t, source.String())
}

func TestStrconvUnsignedFormatting(t *testing.T) {
	var source strings.Builder
	source.WriteString(
		"package main\nuse \"std::strconv\"\nuse \"std::strings\"\nfn main() i32 {\n",
	)
	for i, value := range []string{"0", "1", "9", "10", "42", "100", "4294967296", "9223372036854775808", "18446744073709551615"} {
		fmt.Fprintf(
			&source,
			"owner%d := strconv::format_u64(%s) @fail\nif strings::view(&owner%d) != %s {\nreturn %d\n}\nparsed%d := strconv::parse_u64(strings::view(&owner%d)) @fail\nif parsed%d != %s {\nreturn %d\n}\nstrings::free(&owner%d)\n",
			i,
			value,
			i,
			utf8Literal(value),
			i+1,
			i,
			i,
			i,
			value,
			i+1,
			i,
		)
	}
	source.WriteString("return 0\n}\n")
	runStrconvSource(t, source.String())
}

func TestStrconvSignedAndBooleanFormatting(t *testing.T) {
	var source strings.Builder
	source.WriteString(
		"package main\nuse \"std::strconv\"\nuse \"std::strings\"\nfn main() i32 {\n",
	)
	for i, text := range []string{"0", "1", "-1", "9", "-10", "42", "-100", "9223372036854775807", "-9223372036854775808"} {
		value := text
		if text == "-9223372036854775808" {
			value = "(-9223372036854775807 - 1)"
		}
		fmt.Fprintf(
			&source,
			"owner%d := strconv::format_i64(%s) @fail\nif strings::view(&owner%d) != %s {\nreturn %d\n}\nparsed%d := strconv::parse_i64(strings::view(&owner%d)) @fail\nif parsed%d != %s {\nreturn %d\n}\nstrings::free(&owner%d)\n",
			i,
			value,
			i,
			utf8Literal(text),
			i+1,
			i,
			i,
			i,
			value,
			i+1,
			i,
		)
	}
	for i, value := range []string{"true", "false"} {
		fmt.Fprintf(
			&source,
			"b%d := strconv::format_bool(%s) @fail\nif strings::view(&b%d) != %s {\nreturn 20\n}\nv%d := strconv::parse_bool(strings::view(&b%d)) @fail\nif v%d != %s {\nreturn 21\n}\nstrings::free(&b%d)\n",
			i,
			value,
			i,
			utf8Literal(value),
			i,
			i,
			i,
			value,
			i,
		)
	}
	source.WriteString("return 0\n}\n")
	runStrconvSource(t, source.String())
}
