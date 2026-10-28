// Package ctarget discovers scalar properties from the same native compiler
// used by the C backend. Generated assertions also check ABI-changing flags.
package ctarget

import (
	"fmt"
	"os/exec"
	"strconv"
	"strings"
	"sync"
)

func Compiler() (string, error) {
	for _, name := range []string{"cc", "gcc", "clang"} {
		if path, err := exec.LookPath(name); err == nil {
			return path, nil
		}
	}
	return "", fmt.Errorf("no C compiler found on PATH; install gcc or clang")
}

type Scalar struct {
	Name     string
	Bytes    int
	Unsigned bool
	Float    bool
}

var cache sync.Map

func Scalars() (map[string]Scalar, error) {
	compiler, err := Compiler()
	if err != nil {
		return nil, err
	}
	if cached, ok := cache.Load(compiler); ok {
		return cached.(map[string]Scalar), nil
	}
	cmd := exec.Command(compiler, "-dM", "-E", "-x", "c", "-")
	cmd.Stdin = strings.NewReader("")
	output, err := cmd.CombinedOutput()
	if err != nil {
		return nil, fmt.Errorf("query C target scalar properties: %w: %s", err, output)
	}
	macros := map[string]string{}
	for _, line := range strings.Split(string(output), "\n") {
		parts := strings.Fields(line)
		if len(parts) >= 3 && parts[0] == "#define" {
			macros[parts[1]] = parts[2]
		}
	}
	if macros["__CHAR_BIT__"] != "8" {
		return nil, fmt.Errorf("unsupported C target: requires 8-bit bytes")
	}
	result := map[string]Scalar{}
	for _, spec := range []struct {
		key, name, size    string
		unsigned, floating bool
	}{
		{"char", "char", "1", macros["__CHAR_UNSIGNED__"] != "", false},
		{"schar", "signed char", "1", false, false}, {"uchar", "unsigned char", "1", true, false},
		{"short", "short", "__SIZEOF_SHORT__", false, false}, {"ushort", "unsigned short", "__SIZEOF_SHORT__", true, false},
		{"int", "int", "__SIZEOF_INT__", false, false}, {"uint", "unsigned int", "__SIZEOF_INT__", true, false},
		{"long", "long", "__SIZEOF_LONG__", false, false}, {"ulong", "unsigned long", "__SIZEOF_LONG__", true, false},
		{"longlong", "long long", "__SIZEOF_LONG_LONG__", false, false}, {"ulonglong", "unsigned long long", "__SIZEOF_LONG_LONG__", true, false},
		{"size_t", "size_t", "__SIZEOF_SIZE_T__", true, false}, {"ptrdiff_t", "ptrdiff_t", "__SIZEOF_PTRDIFF_T__", false, false},
		{"float", "float", "__SIZEOF_FLOAT__", false, true}, {"double", "double", "__SIZEOF_DOUBLE__", false, true},
	} {
		size := spec.size
		if size != "1" {
			size = macros[size]
		}
		bytes, err := strconv.Atoi(size)
		if err != nil || bytes != 1 && bytes != 2 && bytes != 4 && bytes != 8 && bytes != 16 {
			return nil, fmt.Errorf("unsupported C target scalar %s size %q", spec.name, size)
		}
		if spec.floating && (spec.key == "float" && bytes != 4 || spec.key == "double" && bytes != 8) {
			return nil, fmt.Errorf("unsupported C floating-point representation: %s", spec.name)
		}
		result["__c_"+spec.key] = Scalar{Name: spec.name, Bytes: bytes, Unsigned: spec.unsigned, Float: spec.floating}
	}
	cache.Store(compiler, result)
	return result, nil
}
