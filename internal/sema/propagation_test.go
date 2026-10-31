package sema_test

import (
	"strings"
	"testing"
)

func TestTryCannotReturnFromDeferredCleanup(t *testing.T) {
	c := checkDiagnostics(t, `package main
fn cleanup() error {
return nil
}
fn run() error {
defer cleanup() @try
return nil
}
fn main() i32 {
return 0
}
`)
	if len(c.Diags) != 1 ||
		!strings.Contains(c.Diags[0].Message, "@try cannot propagate from deferred cleanup") {
		t.Fatalf("want deferred propagation diagnostic, got %v", c.Diags)
	}
}
