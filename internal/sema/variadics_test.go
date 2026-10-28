package sema_test

import "testing"

func TestCVariadicMarkersAreNotOrdinaryValues(t *testing.T) {
	c := checkDiagnostics(
		t,
		"package main\nfn inspect(marker i32, args @c ...) {\nvalue := args\n}\nfn main() {}\n",
	)
	if len(c.Diags) != 1 || c.Diags[0].Title != "C variadic marker is not a value" ||
		!c.Diags[0].Span.Valid() ||
		c.Diags[0].Help == "" {
		t.Fatalf("want variadic marker diagnostic, got %v", c.Diags)
	}
}

func TestCVariadicDefinitionsRequireAFixedParameter(t *testing.T) {
	c := checkDiagnostics(t, "package main\nfn inspect(args @c ...) {}\nfn main() {}\n")
	if len(c.Diags) != 1 || c.Diags[0].Title != "C variadic fixed parameter" ||
		!c.Diags[0].Span.Valid() ||
		c.Diags[0].Help == "" {
		t.Fatalf("want fixed-parameter diagnostic, got %v", c.Diags)
	}
}

func TestCVariadicsRejectUnverifiedAggregateRepresentations(t *testing.T) {
	c := checkDiagnostics(t, `package main
@[default_cc="c"]
extern {
fn send(marker i32, args @c ...)
}
struct Record {
value i32
}
fn main() {
send(0, Record.{value: 42})
}
`)
	if len(c.Diags) != 1 || c.Diags[0].Title != "Unsupported C variadic argument" ||
		!c.Diags[0].Span.Valid() ||
		c.Diags[0].Help == "" ||
		len(c.Diags[0].Related) != 1 {
		t.Fatalf("want actionable variadic diagnostic, got %v", c.Diags)
	}
}
