package ctarget

import (
	"fmt"
	"os/exec"
	"strings"
	"testing"
)

func TestWcharSignednessFromCompilerMacros(t *testing.T) {
	const sizes = `#define __CHAR_BIT__ 8
#define __SIZEOF_SHORT__ 2
#define __SIZEOF_INT__ 4
#define __SIZEOF_LONG__ 4
#define __SIZEOF_LONG_LONG__ 8
#define __SIZEOF_SIZE_T__ 8
#define __SIZEOF_PTRDIFF_T__ 8
#define __SIZEOF_POINTER__ 8
#define __SIZEOF_FLOAT__ 4
#define __SIZEOF_DOUBLE__ 8
`
	for _, tc := range []struct {
		name, macros string
		bytes        int
		unsigned     bool
	}{
		{"mingw-gcc", "#define __WCHAR_TYPE__ short unsigned int\n#define __SIZEOF_WCHAR_T__ 2", 2, true},
		{"unsigned-first", "#define __WCHAR_TYPE__ unsigned short\n#define __SIZEOF_WCHAR_T__ 2", 2, true},
		{"clang", "#define __WCHAR_UNSIGNED__ 1\n#define __SIZEOF_WCHAR_T__ 2", 2, true},
		{"signed", "#define __WCHAR_TYPE__ int\n#define __SIZEOF_WCHAR_T__ 4", 4, false},
	} {
		t.Run(tc.name, func(t *testing.T) {
			scalars, err := parseScalars(sizes + tc.macros)
			if err != nil {
				t.Fatal(err)
			}
			got := scalars["__c_wchar_t"]
			if got.Bytes != tc.bytes || got.Unsigned != tc.unsigned {
				t.Fatalf("wchar_t = %+v, want bytes=%d unsigned=%t", got, tc.bytes, tc.unsigned)
			}
		})
	}
}

func TestScalarsMatchNativeCompiler(t *testing.T) {
	scalars, err := Scalars()
	if err != nil {
		t.Fatal(err)
	}
	compiler, err := Compiler()
	if err != nil {
		t.Fatal(err)
	}
	var source strings.Builder
	source.WriteString("#include <stddef.h>\n#include <stdint.h>\n")
	for _, scalar := range scalars {
		fmt.Fprintf(
			&source,
			"_Static_assert(sizeof(%s) == %d, \"%s size\");\n",
			scalar.Name,
			scalar.Bytes,
			scalar.Name,
		)
		if !scalar.Float {
			fmt.Fprintf(
				&source,
				"_Static_assert((((%s)-1) > 0) == %d, \"%s signedness\");\n",
				scalar.Name,
				boolInt(scalar.Unsigned),
				scalar.Name,
			)
		}
	}
	cmd := exec.Command(compiler, "-std=c11", "-fsyntax-only", "-x", "c", "-")
	cmd.Stdin = strings.NewReader(source.String())
	if output, err := cmd.CombinedOutput(); err != nil {
		t.Fatalf("scalar ABI mismatch: %v\n%s", err, output)
	}
}

func boolInt(value bool) int {
	if value {
		return 1
	}
	return 0
}
