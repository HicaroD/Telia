package integration

import (
	"strings"
	"testing"

	"github.com/HicaroD/Telia/tests/compiler"
)

func TestAssignmentArityMismatchReportsDiagnostic(t *testing.T) {
	for _, assignment := range []string{
		"record.value = read()",
		"record.value = 1, 2",
		"record.value, record.other = 1",
	} {
		t.Run(assignment, func(t *testing.T) {
			path := loopControlSource(t, `package main
struct Record {
value i32
other i32
}
fn read() (i32, error) {
return 42, nil
}
fn main() i32 {
record := Record{value: 0, other: 0}
`+assignment+`
return 0
}
`)
			_, diags := compiler.CompileOnly(path)
			if len(diags.Diags) != 1 ||
				!strings.Contains(diags.Diags[0].Message, "assignment") ||
				!diags.Diags[0].Span.Valid() {
				t.Fatalf("want source-located assignment diagnostic, got %v", diags.Diags)
			}
		})
	}
}
