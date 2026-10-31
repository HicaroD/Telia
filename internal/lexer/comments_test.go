package lexer_test

import (
	"testing"

	"github.com/HicaroD/Telia/internal/lexer/token"
)

func TestComments(t *testing.T) {
	for _, test := range []struct {
		name  string
		src   string
		kinds []token.Kind
	}{
		{"line at EOF", "// ignored ! /*", []token.Kind{token.EOF}},
		{"empty line", "//", []token.Kind{token.EOF}},
		{"line newline", "a// ignored\r\nb", []token.Kind{token.ID, token.NEWLINE, token.ID, token.EOF}},
		{"empty block", "/**/", []token.Kind{token.EOF}},
		{"block at EOF", "a/* ignored */", []token.Kind{token.ID, token.EOF}},
		{"adjacent", "/*a*//*b*///c", []token.Kind{token.EOF}},
		{"separate identifiers", "a/*ignored*/b", []token.Kind{token.ID, token.ID, token.EOF}},
		{"multiline block", "a/*\n\n*/b", []token.Kind{token.ID, token.NEWLINE, token.NEWLINE, token.ID, token.EOF}},
		{"non nested", "/* outer /* inner */tail", []token.Kind{token.ID, token.EOF}},
		{"quotes and line markers in block", "/* \" // */a", []token.Kind{token.ID, token.EOF}},
		{"division at EOF", "/", []token.Kind{token.SLASH, token.EOF}},
		{"division", "8/2 / 1", []token.Kind{token.INT_TYPE, token.SLASH, token.INT_TYPE, token.SLASH, token.INT_TYPE, token.EOF}},
		{"division by commented operand", "8/ /*comment*/2", []token.Kind{token.INT_TYPE, token.SLASH, token.INT_TYPE, token.EOF}},
	} {
		t.Run(test.name, func(t *testing.T) {
			lex, collector := lookaheadLexer(test.src)
			lex.PeekN(100)
			for _, want := range test.kinds {
				if got := lex.Next().Kind; got != want {
					t.Fatalf("got %v, want %v", got, want)
				}
			}
			if collector.HasErrors() {
				t.Fatalf("unexpected diagnostics: %v", collector.Diags)
			}
		})
	}
}

func TestCommentMarkersInStrings(t *testing.T) {
	for _, text := range []string{"// line", "/* block */", "/* unterminated", "https://example.com"} {
		lex, collector := lookaheadLexer("\"" + text + "\"")
		tok := lex.Next()
		if tok.Kind != token.STRING_TYPE || string(tok.Lexeme) != text || collector.HasErrors() {
			t.Fatalf("comment markers changed string: %+v, %v", tok, collector.Diags)
		}
	}
}

func TestCommentPositionsAndSpans(t *testing.T) {
	lex, _ := lookaheadLexer("//x\r\na/*x\n y\n*/b")
	lex.PeekN(20)
	for _, want := range []struct {
		kind                     token.Kind
		line, column, start, end int
	}{
		{token.NEWLINE, 1, 5, 4, 5},
		{token.ID, 2, 1, 5, 6},
		{token.NEWLINE, 2, 5, 9, 10},
		{token.NEWLINE, 3, 3, 12, 13},
		{token.ID, 4, 3, 15, 16},
		{token.EOF, 4, 4, 16, 16},
	} {
		got := lex.Next()
		if got.Kind != want.kind || got.Pos.Line != want.line || got.Pos.Column != want.column ||
			got.Span.Start != want.start || got.Span.End != want.end {
			t.Fatalf("got %+v span %+v, want %+v", got, got.Span, want)
		}
	}
}

func TestUnterminatedBlockComments(t *testing.T) {
	for _, src := range []string{"/*", "/*/", "/* text*", "/*\ntext", "a /*\n/* nested"} {
		t.Run(src, func(t *testing.T) {
			lex, collector := lookaheadLexer(src)
			lex.PeekN(100)
			if collector.HasErrors() {
				t.Fatal("lookahead published diagnostics")
			}
			_, err := lex.Tokenize()
			if err == nil || len(collector.Diags) != 1 {
				t.Fatalf("expected one error, got %v, %v", err, collector.Diags)
			}
			start := 0
			if src[0] == 'a' {
				start = 2
			}
			diag := collector.Diags[0]
			if diag.Title != "Unterminated block comment" || diag.Span.Start != start ||
				diag.Span.End != len(src) {
				t.Fatalf("unexpected diagnostic: %+v", diag)
			}
			if lex.Next().Kind != token.EOF || len(collector.Diags) != 1 {
				t.Fatal("unterminated comment did not reach stable EOF")
			}
		})
	}
}

func TestCommentNewlineModeChanges(t *testing.T) {
	lex, _ := lookaheadLexer("a/*\nx\ny\n*/b//line\nc")
	lex.PeekN(20)
	lex.Next() // a
	lex.Next() // first newline inside block
	lex.StrictNewline = false
	if lex.Peek().Name() != "b" {
		t.Fatal("non-strict lookahead did not skip block newlines")
	}
	lex.StrictNewline = true
	for _, want := range []token.Kind{token.NEWLINE, token.NEWLINE, token.ID, token.NEWLINE, token.ID, token.EOF} {
		if got := lex.Next().Kind; got != want {
			t.Fatalf("mode change lost comment state: got %v, want %v", got, want)
		}
	}

	lex, _ = lookaheadLexer("/*\n*/a//line\nb")
	lex.StrictNewline = false
	for _, want := range []token.Kind{token.ID, token.ID, token.EOF} {
		if got := lex.Next().Kind; got != want {
			t.Fatalf("non-strict comment scan: got %v, want %v", got, want)
		}
	}
}

func TestCommentDiagnosticAfterNewlineModeChanges(t *testing.T) {
	lex, collector := lookaheadLexer("/*\nunterminated")
	lex.PeekN(20)
	lex.StrictNewline = false
	lex.PeekN(20)
	lex.StrictNewline = true
	if lex.Next().Kind != token.NEWLINE || collector.HasErrors() {
		t.Fatal("discarded comment diagnostic was published")
	}
	lex.StrictNewline = false
	if lex.Next().Kind != token.INVALID || len(collector.Diags) != 1 {
		t.Fatal("unterminated comment must report once on consumption")
	}
	if lex.Next().Kind != token.EOF || len(collector.Diags) != 1 {
		t.Fatal("unterminated comment diagnostic was repeated")
	}
}
