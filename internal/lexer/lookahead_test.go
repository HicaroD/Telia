package lexer_test

import (
	"reflect"
	"testing"

	"github.com/HicaroD/Telia/internal/ast"
	"github.com/HicaroD/Telia/internal/diagnostics"
	"github.com/HicaroD/Telia/internal/lexer"
	"github.com/HicaroD/Telia/internal/lexer/token"
)

func lookaheadLexer(src string) (*lexer.Lexer, *diagnostics.Collector) {
	collector := diagnostics.New()
	return lexer.New(&ast.Loc{Name: "lookahead.t"}, []byte(src), collector), collector
}

func TestLookaheadReusesTokensWithoutConsuming(t *testing.T) {
	lex, _ := lookaheadLexer("first second third")
	first := lex.Peek()
	third := lex.PeekN(2)
	if lex.Peek() != first || lex.Peek1() != lex.PeekN(1) || lex.PeekN(2) != third {
		t.Fatal("repeated lookahead should reuse tokens")
	}
	if span := lex.SpanFrom(first.Span); span.End != 0 {
		t.Fatalf("lookahead advanced consumed cursor: %+v", span)
	}
	if lex.Next() != first || lex.Next().Name() != "second" || lex.Next() != third {
		t.Fatal("lookahead changed consumption order")
	}
	if span := lex.SpanFrom(first.Span); span.End != len("first second third") {
		t.Fatalf("consumed cursor not restored: %+v", span)
	}
	for i := 0; i < 10; i++ {
		if lex.PeekN(10000).Kind != token.EOF || lex.Next().Kind != token.EOF {
			t.Fatal("expected stable EOF")
		}
	}
	if first.Name() != "first" || third.Name() != "third" {
		t.Fatal("buffer reuse mutated a previously returned token")
	}
}

func TestLookaheadBraceDepthChangesOnlyOnConsumption(t *testing.T) {
	lex, _ := lookaheadLexer("{{}}")
	lex.PeekN(3)
	if lex.BraceDepth != 0 {
		t.Fatal("peeking changed brace depth")
	}
	for _, expected := range []int{1, 2, 1, 0} {
		lex.Next()
		if lex.BraceDepth != expected {
			t.Fatalf("brace depth: got %d, want %d", lex.BraceDepth, expected)
		}
	}
}

func TestLookaheadNewlineModeInvalidation(t *testing.T) {
	lex, _ := lookaheadLexer("a\n b\n c")
	lex.PeekN(4)
	lex.Next() // a
	lex.StrictNewline = false
	if lex.Peek().Name() != "b" || lex.Peek1().Name() != "c" {
		t.Fatal("stale strict-newline lookahead")
	}
	lex.Next() // b
	lex.StrictNewline = true
	if lex.Next().Kind != token.NEWLINE || lex.Next().Name() != "c" {
		t.Fatal("changing back to strict mode lost the unconsumed newline")
	}
}

func TestLookaheadDefersDiagnosticsExactlyOnce(t *testing.T) {
	for _, src := range []string{"!", "\"bad\\q\\z\"", "1.2.3", "\"unterminated"} {
		t.Run(src, func(t *testing.T) {
			direct, want := lookaheadLexer(src)
			direct.Next()
			lex, got := lookaheadLexer(src)
			lex.PeekN(5)
			lex.Peek()
			if got.HasErrors() {
				t.Fatal("lookahead published diagnostics")
			}
			lex.Next()
			if len(got.Diags) == 0 || len(got.Diags) != len(want.Diags) {
				t.Fatalf("got %d diagnostics, want %d", len(got.Diags), len(want.Diags))
			}
			for i, diag := range got.Diags {
				expected := want.Diags[i]
				if diag.Message != expected.Message || diag.Span.Start != expected.Span.Start ||
					diag.Span.End != expected.Span.End {
					t.Fatalf("diagnostic changed: %+v", diag)
				}
			}
			count := len(got.Diags)
			lex.Next()
			if len(got.Diags) != count {
				t.Fatal("diagnostic published twice")
			}
		})
	}
}

func TestLookaheadMatchesDirectTokenization(t *testing.T) {
	for _, src := range []string{
		"fn test[T](x T) T {\nreturn x\n}\n",
		"a::b(12, 3.5, \"escaped\\ntext\")\n{}[] + - != ==",
		"\"bad\\q\" ! 1.2.3\n",
	} {
		direct, _ := lookaheadLexer(src)
		buffered, _ := lookaheadLexer(src)
		for i := 0; ; i++ {
			buffered.PeekN(i % 7)
			want, got := direct.Next(), buffered.Next()
			if want.Kind != got.Kind ||
				want.Pos != got.Pos ||
				!reflect.DeepEqual(want.Lexeme, got.Lexeme) ||
				want.Literal != got.Literal ||
				want.Span.Start != got.Span.Start ||
				want.Span.End != got.Span.End {
				t.Fatalf("token %d differs: got %+v, want %+v", i, got, want)
			}
			if got.Kind == token.EOF {
				break
			}
		}
	}
}

func TestLookaheadInvalidationDoesNotPublishDiscardedErrors(t *testing.T) {
	lex, collector := lookaheadLexer("\n!")
	lex.PeekN(1)
	lex.StrictNewline = false
	lex.Peek()
	lex.StrictNewline = true
	lex.Next() // newline, not the speculative invalid token
	if collector.HasErrors() {
		t.Fatal("discarded speculative error was published")
	}
	lex.Next()
	if len(collector.Diags) != 1 {
		t.Fatalf("got %d errors after consuming invalid token", len(collector.Diags))
	}
}
