package diagnostics_test

import (
	"bytes"
	"strings"
	"testing"

	"github.com/HicaroD/Telia/internal/diagnostics"
	"github.com/HicaroD/Telia/internal/source"
)

func TestRenderSourceHelpRelatedAndCounts(t *testing.T) {
	f := source.New("main.telia", []byte("age: i32 = \"twenty\"\n"))
	c := diagnostics.New()
	c.ReportAndSave(
		diagnostics.Diag{
			Title:   "Type mismatch",
			Message: "Expected i32, but this value is string.",
			Span:    source.Span{File: f, Start: 11, End: 19},
			Help:    "Try an integer value, such as 20.",
			Related: []diagnostics.Related{
				{Span: source.Span{File: f, Start: 0, End: 3}, Message: "Declared here"},
			},
		},
	)
	c.ReportAndSave(diagnostics.Diag{Severity: diagnostics.Warning, Message: "Unused value"})
	var out bytes.Buffer
	if err := (diagnostics.Renderer{}).Render(&out, c); err != nil {
		t.Fatal(err)
	}
	want := "main.telia:1:12 · Type mismatch\n\n    1 │ age: i32 = \"twenty\"\n      │            ~~~~~~~~\n\n    Expected i32, but this value is string.\n    Try an integer value, such as 20.\n\n    main.telia:1:1 · Declared here\n\n    1 │ age: i32 = \"twenty\"\n      │ ~~~\n\nWarning · Unused value\n\nCompilation failed · 1 error · 1 warning\n"
	if out.String() != want {
		t.Fatalf("want:\n%s\ngot:\n%s", want, out.String())
	}
	if strings.Contains(out.String(), "\x1b") {
		t.Fatal("plain output contains ANSI")
	}
}

func TestColorPolicy(t *testing.T) {
	for _, tc := range []struct {
		mode                    string
		terminal, noColor, want bool
	}{
		{"auto", true, false, true}, {"auto", false, false, false}, {"auto", true, true, false},
		{"always", false, true, true}, {"never", true, false, false},
	} {
		if got := diagnostics.ColorEnabled(tc.mode, tc.terminal, tc.noColor); got != tc.want {
			t.Errorf("%+v: got %v", tc, got)
		}
	}
}

func TestColorDoesNotChangeDiagnosticText(t *testing.T) {
	c := diagnostics.New()
	f := source.New("main.t", []byte("bad\n"))
	c.ReportAndSave(
		diagnostics.Diag{
			Title:   "Problem",
			Message: "Explanation",
			Span:    source.Span{File: f, Start: 0, End: 3},
			Help:    "Fix it",
			Related: []diagnostics.Related{
				{Span: source.Span{File: f, Start: 0, End: 1}, Message: "Origin"},
			},
		},
	)
	c.ReportAndSave(diagnostics.Diag{Severity: diagnostics.Warning, Message: "Warning"})
	var plain, colored bytes.Buffer
	if err := (diagnostics.Renderer{}).Render(&plain, c); err != nil {
		t.Fatal(err)
	}
	if err := (diagnostics.Renderer{Color: true}).Render(&colored, c); err != nil {
		t.Fatal(err)
	}
	text := colored.String()
	for _, code := range []string{"31", "33", "36", "90"} {
		if !strings.Contains(text, "\x1b["+code+"m") {
			t.Errorf("missing ANSI color %s", code)
		}
		text = strings.ReplaceAll(text, "\x1b["+code+"m", "")
	}
	text = strings.ReplaceAll(text, "\x1b[0m", "")
	if text != plain.String() {
		t.Fatal("color changed diagnostic content")
	}
}

func TestRenderTabsUTF8MultilineAndEOF(t *testing.T) {
	f := source.New("unicode.t", []byte("\té = 1\r\nnext\n"))
	c := diagnostics.New()
	c.ReportAndSave(
		diagnostics.Diag{Message: "Across lines", Span: source.Span{File: f, Start: 1, End: 13}},
	)
	c.ReportAndSave(
		diagnostics.Diag{
			Message: "At EOF",
			Span:    source.Span{File: f, Start: len(f.Text), End: len(f.Text)},
		},
	)
	var out bytes.Buffer
	if err := (diagnostics.Renderer{}).Render(&out, c); err != nil {
		t.Fatal(err)
	}
	for _, want := range []string{"unicode.t:1:5", "    1 │     é = 1\n      │     ~~~~~", "    2 │ next\n      │ ~~~~", "unicode.t:3:1", "    3 │ \n      │ ~"} {
		if !strings.Contains(out.String(), want) {
			t.Fatalf("missing %q in:\n%s", want, out.String())
		}
	}
	if strings.Contains(out.String(), "\r") {
		t.Fatal("CRLF leaked into output")
	}
}

func TestHighlightsUseTerminalCellWidths(t *testing.T) {
	for _, tc := range []struct{ text, location, marker string }{
		{"界 bad", "wide.t:1:4", "      │    ~~~"},
		{"e\u0301 bad", "wide.t:1:3", "      │   ~~~"},
	} {
		c := diagnostics.New()
		f := source.New("wide.t", []byte(tc.text))
		start := strings.Index(tc.text, "bad")
		c.ReportAndSave(
			diagnostics.Diag{
				Message: "Bad",
				Span:    source.Span{File: f, Start: start, End: start + 3},
			},
		)
		var out bytes.Buffer
		if err := (diagnostics.Renderer{}).Render(&out, c); err != nil {
			t.Fatal(err)
		}
		if !strings.Contains(out.String(), tc.location) ||
			!strings.Contains(out.String(), tc.marker) {
			t.Fatalf("misaligned highlight:\n%s", out.String())
		}
	}
}
