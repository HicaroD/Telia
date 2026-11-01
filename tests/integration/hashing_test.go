package integration

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"testing"

	"github.com/HicaroD/Telia/tests/compiler"
)

func TestHashRejectsNestedStandardOwner(t *testing.T) {
	path := filepath.Join(t.TempDir(), "main.t")
	source := `package main
use "std::strings"
struct Box<T> {
 value T
}
fn main() {
 h := __key_hash<Box<strings::Owned>>(zero<Box<strings::Owned>>(), 0, 0)
}
`
	if err := os.WriteFile(path, []byte(source), 0600); err != nil {
		t.Fatal(err)
	}
	_, diags := compiler.CompileOnly(path)
	for _, d := range diags.Diags {
		if strings.Contains(d.Message, "value.value") &&
			strings.Contains(d.Message, "owning handle") {
			return
		}
	}
	t.Fatalf("want nested nominal owner diagnostic, got %v", diags.Diags)
}

func TestHashEntropyInjection(t *testing.T) {
	include, err := filepath.Abs("testdata")
	if err != nil {
		t.Fatal(err)
	}
	for _, fail := range []int{0, 1} {
		for _, opt := range []string{"-O0", "-O3"} {
			t.Run(fmt.Sprintf("%s/fail=%d", opt, fail), func(t *testing.T) {
				path, diags := compiler.CompileOnlyWithCFlags(
					"testdata/hash_entropy.t",
					opt,
					"-I",
					include,
					"-include",
					filepath.Join(include, "hash_entropy.h"),
					fmt.Sprintf("-DTELIA_ENTROPY_FAIL=%d", fail),
				)
				if diags.HasErrors() {
					t.Fatalf("entropy compilation: %v", diags.Diags)
				}
				defer removeTestBinary(t, path)
				if out, err := compiler.RunBinary(path); err != nil {
					t.Fatalf("entropy execution: %v: %s", err, out)
				}
			})
		}
	}
}

func TestSemanticHashing(t *testing.T) {
	include, err := filepath.Abs("testdata")
	if err != nil {
		t.Fatal(err)
	}
	modes := [][]string{{"-O0"}, {"-O3"}}
	if runtime.GOOS != "windows" {
		modes = append(modes, []string{"-O3", "-fsanitize=undefined", "-fno-sanitize-recover=all"})
	}
	for _, flags := range modes {
		t.Run(flags[0]+flags[len(flags)-1], func(t *testing.T) {
			flags = append(flags, "-I", include)
			path, diags := compiler.CompileOnlyWithCFlags("testdata/hashing.t", flags...)
			if diags.HasErrors() {
				t.Fatalf("hash compilation: %v", diags.Diags)
			}
			defer removeTestBinary(t, path)
			if out, err := compiler.RunBinary(path); err != nil {
				t.Fatalf("hash execution: %v: %s", err, out)
			}
		})
	}
}
