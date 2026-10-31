package integration

import (
	"context"
	"os/exec"
	"strings"
	"testing"
	"time"

	"github.com/HicaroD/Telia/tests/compiler"
)

func TestIndexedConditionsShortCircuitAndReevaluate(t *testing.T) {
	path := loopControlSource(t, `package main
fn main() i32 {
bytes u8[3] := .{1, 2, 3}
i u64 := 0
while i < len(bytes) and bytes[i] != 3 {
i = i + 1
}
if i != 2 {
return 1
}
i = 0
while i < len(bytes) and bytes[i] != 9 {
i = i + 1
}
if i != 3 {
return 2
}
empty u8[0] := .{}
if len(empty) > 0 and empty[0] == 1 {
return 3
}
if len(empty) == 0 or empty[0] == 1 {
return 0
}
return 4
}
`)
	for _, opt := range []string{"-O0", "-O3"} {
		t.Run(opt, func(t *testing.T) {
			exe, diags := compiler.CompileOnlyWithCFlags(path, opt)
			if len(diags.Diags) != 0 {
				t.Fatalf("indexed conditions should compile: %v", diags.Diags)
			}
			defer removeTestBinary(t, exe)
			ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
			defer cancel()
			if output, err := exec.CommandContext(ctx, exe).CombinedOutput(); err != nil {
				t.Fatalf("indexed conditions: %v: %s", err, output)
			}
		})
	}
}

func TestErrorHandledReassignmentReportsDiagnostic(t *testing.T) {
	for _, assignment := range []string{
		"record.value = read() @catch err {\nreturn 1\n}",
		"record.value = read() @fail",
		"value = read() @catch err {\nreturn 1\n}",
	} {
		t.Run(assignment, func(t *testing.T) {
			path := loopControlSource(t, `package main
struct Record {
value i32
}
fn read() (i32, error) {
return 42, nil
}
fn main() i32 {
record := Record{value: 0}
value := 0
`+assignment+`
return 0
}
`)
			_, diags := compiler.CompileOnly(path)
			if len(diags.Diags) != 1 ||
				!strings.Contains(diags.Diags[0].Message, "unsupported assignment") ||
				!diags.Diags[0].Span.Valid() {
				t.Fatalf(
					"want source-located unsupported assignment diagnostic, got %v",
					diags.Diags,
				)
			}
		})
	}
}
