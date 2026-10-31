package sema_test

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/HicaroD/Telia/internal/ast"
	"github.com/HicaroD/Telia/internal/diagnostics"
	"github.com/HicaroD/Telia/internal/parser"
	"github.com/HicaroD/Telia/internal/sema"
)

func checkHintProgram(t *testing.T, text string) *diagnostics.Collector {
	t.Helper()
	path := filepath.Join(t.TempDir(), "main.t")
	if err := os.WriteFile(path, []byte(text), 0600); err != nil {
		t.Fatal(err)
	}
	loc, err := ast.LocFromPath(path)
	if err != nil {
		t.Fatal(err)
	}
	c := diagnostics.New()
	program, err := parser.New(c).ParseFileAsProgram(path, loc, c)
	if err != nil {
		t.Fatalf("parse: %v, %v", err, c.Diags)
	}
	if err := sema.New(c).Check(program); err != nil {
		t.Fatalf("check: %v, %v", err, c.Diags)
	}
	return c
}

func TestIfElifElseSuggestsOneMatchRewrite(t *testing.T) {
	c := checkHintProgram(t, `package main
fn classify(code i32) i32 {
if code == 1 or code == 2 {
return 3
} elif code == 5 {
return 4
} elif code == 8 {
return 6
} else {
return 0
}
}
fn main() {}
`)
	if c.HasErrors() || len(c.Diags) != 1 {
		t.Fatalf("want one non-blocking hint, got %v", c.Diags)
	}
	d := c.Diags[0]
	if d.Severity != diagnostics.Hint || d.Title != "Consider using match" || !d.Span.Valid() {
		t.Fatalf("want a located match hint, got %+v", d)
	}
	want := "return match code {\n    1 | 2 => 3\n    5 => 4\n    8 => 6\n    _ => 0\n}"
	if !strings.Contains(d.Help, want) {
		t.Fatalf("want rewrite:\n%s\ngot:\n%s", want, d.Help)
	}
	checkHintProgram(t, "package main\nfn classify(code i32) i32 {\n"+want+"\n}\nfn main() {}\n")
}

func TestMatchHintIncludesRangeGuardsAndElse(t *testing.T) {
	c := checkHintProgram(t, `package main
fn classify(code i32) i32 {
if code == 1 or code == 2 {
return 3
} elif code == 5 {
return 6
} elif code > 10 and code < 20 {
return 4
} else {
return 0
}
}
fn main() {}
`)
	if len(c.Diags) != 1 || c.Diags[0].Severity != diagnostics.Hint {
		t.Fatalf("want one hint, got %v", c.Diags)
	}
	want := "return match code {\n    1 | 2 => 3\n    5 => 6\n    _ if code > 10 and code < 20 => 4\n    _ => 0\n}"
	if !strings.Contains(c.Diags[0].Help, want) {
		t.Fatalf("want guarded rewrite:\n%s\ngot:\n%s", want, c.Diags[0].Help)
	}
	checkHintProgram(t, "package main\nfn classify(code i32) i32 {\n"+want+"\n}\nfn main() {}\n")
}

func TestMatchHintsAvoidUnsafeOrNoisySuggestions(t *testing.T) {
	for _, tc := range []struct{ name, body string }{
		{"two branches", "if code == 1 {\nreturn 1\n} else {\nreturn 0\n}"},
		{"if elif else is too small", "if code == 1 {\nreturn 1\n} elif code == 2 {\nreturn 2\n} else {\nreturn 0\n}"},
		{"alternatives do not inflate branch count", "if code == 1 or code == 2 {\nreturn 1\n} elif code == 3 {\nreturn 2\n} else {\nreturn 0\n}"},
		{"mostly range conditions", "if code == 1 {\nreturn 1\n} elif code > 10 and code < 20 {\nreturn 2\n} elif code > 30 {\nreturn 3\n} else {\nreturn 0\n}"},
		{"no complete fallback", "if code == 1 {\nreturn 1\n} elif code == 2 {\nreturn 2\n} elif code == 3 {\nreturn 3\n}\nreturn 0"},
		{"non-returning branch", "if code == 1 {\ncode = 4\n} elif code == 2 {\nreturn 2\n} elif code == 3 {\nreturn 3\n} else {\nreturn 0\n}\nreturn code"},
		{"different values", "if code == 1 {\nreturn 1\n} elif other == 2 {\nreturn 2\n} elif code == 3 {\nreturn 3\n} else {\nreturn 0\n}"},
		{"repeated calls", "if get(code) == 1 {\nreturn 1\n} elif get(code) == 2 {\nreturn 2\n} elif get(code) == 3 {\nreturn 3\n} else {\nreturn 0\n}"},
		{"condition calls", "if code == 1 {\nreturn 1\n} elif code > get(other) {\nreturn 2\n} elif code == 3 {\nreturn 3\n} else {\nreturn 0\n}"},
		{"unrelated condition", "if code == 1 {\nreturn 1\n} elif other > 10 and code < 20 {\nreturn 2\n} elif code == 3 {\nreturn 3\n} else {\nreturn 0\n}"},
		{"duplicate patterns", "if code == 1 {\nreturn 1\n} elif code == 1 {\nreturn 2\n} elif code == 3 {\nreturn 3\n} else {\nreturn 0\n}"},
		{"separate if statements", "if code == 1 {\nreturn 1\n}\nif code == 2 {\nreturn 2\n}\nreturn 0"},
	} {
		t.Run(tc.name, func(t *testing.T) {
			c := checkHintProgram(
				t,
				"package main\nfn get(n i32) i32 {\nreturn n\n}\nfn classify(code i32, other i32) i32 {\n"+tc.body+"\n}\nfn main() {}\n",
			)
			if len(c.Diags) != 0 {
				t.Fatalf("want no suggestion, got %v", c.Diags)
			}
		})
	}
}

func TestGenericMatchHintIsReportedOncePerSourceChain(t *testing.T) {
	c := checkHintProgram(t, `package main
fn classify<T>(code i32, unused T) i32 {
if code == 1 {
return 1
} elif code == 2 {
return 2
} elif code == 3 {
return 3
} else {
return 0
}
}
fn main() {
a := classify<i32>(1, 0)
b := classify<u64>(2, 0)
}
`)
	if len(c.Diags) != 1 || c.Diags[0].Severity != diagnostics.Hint {
		t.Fatalf("want one hint across specializations, got %v", c.Diags)
	}
}

func TestMatchHintSupportsEnumCases(t *testing.T) {
	const header = `package main
enum Event {
Quit
Idle
Resize
Other
}
fn classify(event Event) i32 {
`
	const chain = `if event == Event::Quit {
return 1
} elif event == Event::Idle {
return 2
} elif event == Event::Resize {
return 3
} else {
return 0
}`
	c := checkHintProgram(t, header+chain+"\n}\nfn main() {}\n")
	if len(c.Diags) != 1 || c.Diags[0].Severity != diagnostics.Hint {
		t.Fatalf("want an enum match hint, got %v", c.Diags)
	}
	rewrite := strings.TrimPrefix(c.Diags[0].Help, "Suggested rewrite:\n\n")
	for _, pattern := range []string{"Event::Quit => 1", "Event::Idle => 2", "Event::Resize => 3", "_ => 0"} {
		if !strings.Contains(rewrite, pattern) {
			t.Fatalf("missing %s in rewrite:\n%s", pattern, rewrite)
		}
	}
	checkHintProgram(t, header+rewrite+"\n}\nfn main() {}\n")
}

func TestMatchHintProducesValidPatternsFromParenthesizedLiterals(t *testing.T) {
	c := checkHintProgram(t, `package main
fn classify(code i32) i32 {
if code == (1) {
return 1
} elif code == (2) {
return 2
} elif code == (3) {
return 3
} else {
return 0
}
}
fn main() {}
`)
	if len(c.Diags) != 1 {
		t.Fatalf("want one hint, got %v", c.Diags)
	}
	rewrite := strings.TrimPrefix(c.Diags[0].Help, "Suggested rewrite:\n\n")
	checkHintProgram(t, "package main\nfn classify(code i32) i32 {\n"+rewrite+"\n}\nfn main() {}\n")
}
