package lexer_test

import (
	"testing"

	"github.com/HicaroD/Telia/internal/ast"
	"github.com/HicaroD/Telia/internal/diagnostics"
	"github.com/HicaroD/Telia/internal/lexer"
	"github.com/HicaroD/Telia/internal/lexer/token"
)

func TestTokenSpansCoverOriginalSpelling(t *testing.T) {
	c := diagnostics.New()
	l := lexer.New(&ast.Loc{Name: "main.t", Path: "pkg/main.t"}, []byte("  \"a\\n\" != 12\n"), c)
	for _, want := range []struct {
		kind       token.Kind
		start, end int
	}{{token.STRING_TYPE, 2, 7}, {token.BANG_EQUAL, 8, 10}, {token.INT_TYPE, 11, 13}, {token.NEWLINE, 13, 14}, {token.EOF, 14, 14}} {
		tok := l.Next()
		if tok.Kind != want.kind || !tok.Span.Valid() || tok.Span.Start != want.start ||
			tok.Span.End != want.end {
			t.Fatalf(
				"token %v: span %+v; want %v [%d,%d)",
				tok,
				tok.Span,
				want.kind,
				want.start,
				want.end,
			)
		}
		if tok.Span.File.Path != "pkg/main.t" {
			t.Fatal(tok.Span.File.Path)
		}
	}
}

func TestLookaheadDoesNotReportAndInvalidTokensMakeProgress(t *testing.T) {
	c := diagnostics.New()
	l := lexer.New(&ast.Loc{Name: "main.t"}, []byte("! $ valid"), c)
	l.Peek()
	l.Peek1()
	l.PeekN(2)
	if len(c.Diags) != 0 {
		t.Fatalf("lookahead emitted %d diagnostics", len(c.Diags))
	}
	for _, want := range []token.Kind{token.INVALID, token.INVALID, token.ID, token.EOF} {
		if got := l.Next().Kind; got != want {
			t.Fatalf("invalid tokens did not advance: got %v, want %v", got, want)
		}
	}
	if len(c.Diags) != 2 {
		t.Fatalf("want two errors, got %v", c.Diags)
	}
	for i, start := range []int{0, 2} {
		if d := c.Diags[i]; !d.Span.Valid() || d.Span.Start != start || d.Span.End != start+1 {
			t.Fatalf("bad diagnostic span: %+v", d)
		}
	}
}

func TestTokenizeCollectsIndependentLexicalErrors(t *testing.T) {
	c := diagnostics.New()
	l := lexer.New(&ast.Loc{Name: "main.t"}, []byte("$ \"bad\\q\" 1.2.3 valid"), c)
	_, err := l.Tokenize()
	if err == nil || len(c.Diags) != 3 {
		t.Fatalf("want three lexical errors, got %v, %v", err, c.Diags)
	}
	if l.Next().Kind != token.EOF {
		t.Fatal("tokenization stopped before EOF")
	}
}
