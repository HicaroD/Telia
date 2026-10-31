package sema_test

import "testing"

func TestInvalidMatchesHaveActionableDiagnostics(t *testing.T) {
	for _, tc := range []struct{ name, body, title string }{
		{"missing variant", "value := match Event::Quit {\nEvent::Quit => 0\n}", "Non-exhaustive match"},
		{"duplicate arm", "value := match Event::Quit {\nEvent::Quit => 0\nEvent::Quit => 1\nEvent::Resize{..} => 2\n}", "Duplicate match arm"},
		{"missing fields", "value := match Event::Quit {\nEvent::Quit => 0\nEvent::Resize{width} => width\n}", "Incomplete payload pattern"},
		{"unknown field", "value := match Event::Quit {\nEvent::Quit => 0\nEvent::Resize{depth, ..} => depth\n}", "Unknown payload field"},
		{"missing yield", "value := match Event::Quit {\nEvent::Quit => 0\nEvent::Resize{..} => {\nlocal := 1\n}\n}", "Missing match yield"},
		{"yield outside match", "yield 1", "Yield outside match"},
		{"unreachable arm", "value := match Event::Quit {\n_ => 0\nEvent::Quit => 1\n}", "Unreachable match arm"},
	} {
		t.Run(tc.name, func(t *testing.T) {
			c := checkDiagnostics(
				t,
				"package main\nenum Event {\nQuit\nResize {\nwidth i32\nheight i32\n}\n}\nfn main() {\n"+tc.body+"\n}\n",
			)
			if len(c.Diags) != 1 || c.Diags[0].Title != tc.title || !c.Diags[0].Span.Valid() ||
				c.Diags[0].Help == "" {
				t.Fatalf("want structured %s diagnostic, got %v", tc.title, c.Diags)
			}
		})
	}
}

func TestUnusedGenericFunctionsRejectYieldOutsideMatch(t *testing.T) {
	c := checkDiagnostics(
		t,
		"package main\nfn invalid<T>(value T) {\nyield value\n}\nfn main() {}\n",
	)
	if len(c.Diags) != 1 || c.Diags[0].Title != "Yield outside match" || !c.Diags[0].Span.Valid() ||
		c.Diags[0].Help == "" {
		t.Fatalf("want structured yield-context diagnostic, got %v", c.Diags)
	}
}

func TestDeferredCodeCannotYieldFromTheEnclosingMatch(t *testing.T) {
	c := checkDiagnostics(
		t,
		"package main\nenum Event {\nQuit\n}\nfn main() {\nvalue := match Event::Quit {\nEvent::Quit => {\ndefer yield 1\nyield 0\n}\n}\n}\n",
	)
	if len(c.Diags) != 1 || c.Diags[0].Title != "Yield outside match" || !c.Diags[0].Span.Valid() ||
		c.Diags[0].Help == "" {
		t.Fatalf("want a yield-context diagnostic instead of recursive cleanup, got %v", c.Diags)
	}
}

func TestConditionalMatchesRejectInvalidCoverageAndBindings(t *testing.T) {
	for _, tc := range []struct{ name, body, title string }{
		{"guard is boolean", "match 1 {\n_ if 1 => 0\n_ => 1\n}", "Invalid match guard"},
		{"guarded variant is not coverage", "match Event::Quit {\nEvent::Quit if true => 0\nEvent::Resize{..} | Event::Wide{..} => 1\n}", "Non-exhaustive match"},
		{"guarded fallback is not coverage", "match 1 {\n_ if true => 0\n}", "Non-exhaustive match"},
		{"integer requires fallback", "match 1 {\n1 | 2 => 0\n}", "Non-exhaustive match"},
		{"alternative missing binding", "match Event::Quit {\nEvent::Quit | Event::Resize{width: n} => 0\n}", "Incompatible alternative bindings"},
		{"alternative binding type", "match Event::Quit {\nEvent::Resize{width: n} | Event::Wide{width: n} => n\n_ => 0\n}", "Incompatible alternative bindings"},
		{"integer alternative missing binding", "match 1 {\n1 | n => 0\n}", "Incompatible alternative bindings"},
		{"integer duplicate after coverage", "match 1 {\n1 | 2 => 0\n2 if true => 1\n_ => 2\n}", "Duplicate match arm"},
		{"integer pattern is not an expression", "match 1 {\n-(1 + 2) => 0\n_ => 1\n}", "Invalid match pattern"},
	} {
		t.Run(tc.name, func(t *testing.T) {
			c := checkDiagnostics(
				t,
				"package main\nenum Event {\nQuit\nResize {\nwidth i32\n}\nWide {\nwidth u64\n}\n}\nfn main() {\nvalue := "+tc.body+"\n}\n",
			)
			if len(c.Diags) != 1 || c.Diags[0].Title != tc.title || !c.Diags[0].Span.Valid() ||
				c.Diags[0].Help == "" {
				t.Fatalf("want structured %s diagnostic, got %v", tc.title, c.Diags)
			}
		})
	}
}
