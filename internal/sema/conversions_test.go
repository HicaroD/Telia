package sema_test

import (
	"strings"
	"testing"
)

func TestConversionsRejectOverflowBeforeCheckingTheDestination(t *testing.T) {
	for _, operation := range []string{"narrow<u8>", "widen<i64>"} {
		c := checkDiagnostics(
			t,
			"package main\nfn main() i32 {\n"+operation+"(4294967297)\nreturn 0\n}\n",
		)
		if len(c.Diags) != 1 || !strings.Contains(c.Diags[0].Message, "outside the range") {
			t.Fatalf("want source constant overflow diagnostic, got %v", c.Diags)
		}
	}
}
