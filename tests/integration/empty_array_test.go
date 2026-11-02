package integration

import (
	"context"
	"os/exec"
	"runtime"
	"testing"
	"time"

	"github.com/HicaroD/Telia/tests/compiler"
)

func TestStructuralEqualityEmptyArraysAndKeys(t *testing.T) {
	for _, opt := range []string{"-O0", "-O3"} {
		for _, sanitized := range []bool{false, true} {
			name := opt
			flags := []string{opt}
			if sanitized {
				if runtime.GOOS == "windows" {
					continue
				}
				name += "/UBSan"
				flags = append(flags, "-fsanitize=undefined", "-fno-sanitize-recover=all")
			}
			t.Run(name, func(t *testing.T) {
				path, diags := compiler.CompileOnlyWithCFlags("testdata/empty_array.t", flags...)
				if diags.HasErrors() {
					t.Fatalf("empty-array compilation: %v", diags.Diags)
				}
				defer removeTestBinary(t, path)
				ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
				defer cancel()
				if output, err := exec.CommandContext(ctx, path).CombinedOutput(); err != nil {
					t.Fatalf("empty-array execution: %v: %s", err, output)
				}
			})
		}
	}
}
