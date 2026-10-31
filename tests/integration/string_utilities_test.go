package integration

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"testing"
)

func runStringSource(t *testing.T, source string, flags ...string) {
	t.Helper()
	file := filepath.Join(t.TempDir(), "strings.t")
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

func TestStringFindExactBytes(t *testing.T) {
	cases := []struct {
		text, needle string
		offset       uint64
		found        bool
	}{
		{"abcabc", "bc", 1, true}, {"abc", "z", 0, false},
		{"", "", 0, true}, {"abc", "", 0, true}, {"", "x", 0, false},
		{"a", "ab", 0, false}, {"é世界😀", "😀", 8, true},
		{"a\x00b\xff", "\x00b", 1, true}, {"a\x00b\xff", "\xff", 3, true},
		{"é", "e\u0301", 0, false}, {"ABC", "abc", 0, false},
		{"aaaaab", "aaab", 2, true}, {"é", "\xa9", 1, true},
	}
	var source strings.Builder
	source.WriteString("package main\nuse \"std::strings\"\nfn main() i32 {\n")
	for i, tc := range cases {
		fmt.Fprintf(
			&source,
			"offset%d, found%d := strings::find(%s, %s)\nif offset%d != %d or found%d != %t {\nreturn %d\n}\n",
			i,
			i,
			utf8Literal(tc.text),
			utf8Literal(tc.needle),
			i,
			tc.offset,
			i,
			tc.found,
			i+1,
		)
	}
	source.WriteString("return 0\n}\n")
	runStringSource(t, source.String())
}

func TestStringExactPredicates(t *testing.T) {
	cases := []struct {
		text, needle             string
		contains, prefix, suffix bool
	}{
		{"", "", true, true, true}, {"abc", "", true, true, true},
		{"", "a", false, false, false}, {"ab", "abc", false, false, false},
		{"abc", "abc", true, true, true}, {"abc", "a", true, true, false},
		{"abc", "b", true, false, false}, {"abc", "c", true, false, true},
		{"é世界😀", "世界", true, false, false}, {"é", "e\u0301", false, false, false},
		{"ABC", "abc", false, false, false}, {"\xff\x00a\x00", "\xff\x00", true, true, false},
		{"\xff\x00a\x00", "a\x00", true, false, true},
	}
	var source strings.Builder
	source.WriteString("package main\nuse \"std::strings\"\nfn main() i32 {\n")
	for i, tc := range cases {
		fmt.Fprintf(
			&source,
			"if strings::contains(%s, %s) != %t or strings::starts_with(%s, %s) != %t or strings::ends_with(%s, %s) != %t {\nreturn %d\n}\n",
			utf8Literal(tc.text),
			utf8Literal(tc.needle),
			tc.contains,
			utf8Literal(tc.text),
			utf8Literal(tc.needle),
			tc.prefix,
			utf8Literal(tc.text),
			utf8Literal(tc.needle),
			tc.suffix,
			i+1,
		)
	}
	source.WriteString("return 0\n}\n")
	runStringSource(t, source.String())
}

func TestStringTrimUnicodeWhitespace(t *testing.T) {
	// Unicode 17.0.0 White_Space, explicitly enumerated independently of the implementation.
	spaces := []rune{0x09, 0x0a, 0x0b, 0x0c, 0x0d, 0x20, 0x85, 0xa0, 0x1680,
		0x2000, 0x2001, 0x2002, 0x2003, 0x2004, 0x2005, 0x2006, 0x2007, 0x2008, 0x2009, 0x200a,
		0x2028, 0x2029, 0x202f, 0x205f, 0x3000}
	cases := []struct{ text, want string }{
		{"", ""}, {"世界😀", "世界😀"}, {"  a b  ", "a b"},
		{" \u2003\t\u00a0", ""}, {"\u3000 \té\u00a0\n", "é"},
		{" \x00 ", "\x00"}, {"\t é\u2003世 😀 \n", "é\u2003世 😀"},
	}
	for _, point := range spaces {
		s := string(point)
		cases = append(cases, struct{ text, want string }{s + "é世界😀" + s, "é世界😀"},
			struct{ text, want string }{s + s, ""})
	}
	// Adjacent nonmembers and historically mistaken whitespace (notably U+180E and BOM).
	for _, point := range []rune{0, 8, 14, 0x1c, 0x1f, 0x21, 0x84, 0x86, 0x9f, 0xa1,
		0x167f, 0x1681, 0x180e, 0x1fff, 0x200b, 0x2027, 0x202a, 0x202e, 0x2030,
		0x205e, 0x2060, 0x2fff, 0x3001, 0xfeff, 0x10ffff} {
		s := string(point)
		cases = append(cases, struct{ text, want string }{" " + s + " ", s})
	}
	var source strings.Builder
	source.WriteString("package main\nuse \"std::strings\"\nfn main() i32 {\n")
	for i, tc := range cases {
		fmt.Fprintf(
			&source,
			"text%d, err%d := strings::trim_space(%s)\nif err%d != nil or text%d != %s {\nreturn %d\n}\n",
			i,
			i,
			utf8Literal(tc.text),
			i,
			i,
			utf8Literal(tc.want),
			i+1,
		)
	}
	source.WriteString("return 0\n}\n")
	runStringSource(t, source.String())
}

func TestStringTrimRejectsMalformedInterior(t *testing.T) {
	var source strings.Builder
	source.WriteString("package main\nuse \"std::strings\"\nfn main() i32 {\n")
	for i, text := range []string{
		"a\xffb", " \xff ", "valid\x00\xff", "\xc0\x80", "\xed\xa0\x80",
		"\xf4\x90\x80\x80", "\xe2\x80", "\x80", "a\xe1A\x80b",
	} {
		fmt.Fprintf(
			&source,
			"text%d, err%d := strings::trim_space(%s)\nif err%d == nil or text%d != \"\" {\nreturn %d\n}\n",
			i,
			i,
			utf8Literal(text),
			i,
			i,
			i+1,
		)
	}
	source.WriteString("return 0\n}\n")
	runStringSource(t, source.String())
}

func TestStringSplitPreservesFields(t *testing.T) {
	cases := []struct {
		text, separator string
		want            []string
	}{
		{"a,b,c", ",", []string{"a", "b", "c"}},
		{"", ",", []string{""}}, {",", ",", []string{"", ""}},
		{",a,,b,", ",", []string{"", "a", "", "b", ""}},
		{"abc", "x", []string{"abc"}}, {"a", "abc", []string{"a"}},
		{"aaaaa", "aa", []string{"", "", "a"}},
		{"ababa", "aba", []string{"", "ba"}},
		{"é世界😀世界", "世界", []string{"é", "😀", ""}},
		{"\xff\x00a\xff\x00", "\xff\x00", []string{"", "a", ""}},
		{"a\x00b\x00c", "\x00", []string{"a", "b", "c"}},
		{"é", "e\u0301", []string{"é"}},
	}
	var source strings.Builder
	source.WriteString("package main\nuse \"std::strings\"\nuse \"std::lists\"\nfn main() i32 {\n")
	for i, tc := range cases {
		fmt.Fprintf(
			&source,
			"fields%d, err%d := strings::split(%s, %s)\nif err%d != nil or lists::count<string>(&fields%d) != %d {\nreturn %d\n}\n",
			i,
			i,
			utf8Literal(tc.text),
			utf8Literal(tc.separator),
			i,
			i,
			len(tc.want),
			i+1,
		)
		for j, want := range tc.want {
			fmt.Fprintf(
				&source,
				"value%d_%d := lists::get<string>(&fields%d, %d) @fail\nif value%d_%d != %s {\nreturn %d\n}\n",
				i,
				j,
				i,
				j,
				i,
				j,
				utf8Literal(want),
				i+1,
			)
		}
		fmt.Fprintf(&source, "lists::free<string>(&fields%d)\n", i)
	}
	source.WriteString("return 0\n}\n")
	runStringSource(t, source.String())
}

func TestStringUtilitiesOwnershipAndAllocationFailures(t *testing.T) {
	header, err := filepath.Abs("testdata/allocation_control.h")
	if err != nil {
		t.Fatal(err)
	}
	runStandardLibraryFixture(t, "string_utilities_ownership.t", "-include", header)
}

func TestStringUtilitiesWithNativeSanitizers(t *testing.T) {
	if runtime.GOOS == "windows" {
		t.Skip("MinGW-w64 does not provide the required sanitizer runtime")
	}
	sanitizers := os.Getenv("TELIA_TEST_SANITIZERS")
	if sanitizers == "" {
		sanitizers = "undefined"
	}
	header, err := filepath.Abs("testdata/allocation_control.h")
	if err != nil {
		t.Fatal(err)
	}
	runStandardLibraryFixture(t, "string_utilities_ownership.t", "-include", header,
		"-fsanitize="+sanitizers, "-fno-sanitize-recover=all", "-fno-omit-frame-pointer")
}
