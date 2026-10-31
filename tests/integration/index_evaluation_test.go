package integration

import (
	"context"
	"os/exec"
	"testing"
	"time"

	"github.com/HicaroD/Telia/tests/compiler"
)

func TestGuardedIndexEvaluatesCollectionAndIndexOncePerAccess(t *testing.T) {
	path := loopControlSource(t, `package main
use "std::io"
fn source(values u8[]) u8[] {
io::println("source")
return values
}
fn pick_index(i u64) u64 {
io::println("index")
return i
}
fn main() i32 {
values u8[2] := .{1, 2}
view := values[:]
i u64 := 0
while i < len(view) and source(view)[pick_index(i)] != 9 {
i = i + 1
continue
}
if i != 2 {
return 1
}
return 0
}
`)
	for _, opt := range []string{"-O0", "-O3"} {
		t.Run(opt, func(t *testing.T) {
			exe, diags := compiler.CompileOnlyWithCFlags(
				path,
				opt,
				"-fsanitize=undefined",
				"-fno-sanitize-recover=all",
			)
			if len(diags.Diags) != 0 {
				t.Fatalf("guarded index should compile: %v", diags.Diags)
			}
			defer removeTestBinary(t, exe)
			ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
			defer cancel()
			output, err := exec.CommandContext(ctx, exe).CombinedOutput()
			if want := "source\nindex\nsource\nindex\n"; err != nil || string(output) != want {
				t.Fatalf("want %q, got %q, %v", want, output, err)
			}
		})
	}
}
