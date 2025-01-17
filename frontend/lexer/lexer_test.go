package lexer

import (
	"fmt"
	"reflect"
	"testing"

	"github.com/HicaroD/Telia/diagnostics"
	"github.com/HicaroD/Telia/frontend/lexer/token"
)

type tokenKindTest struct {
	lexeme string
	kind   token.Kind
}

func TestTokenKinds(t *testing.T) {
	filename := "test.tt"

	tests := []*tokenKindTest{
		// Keywords
		{"fn", token.FN},
		{"for", token.FOR},
		{"while", token.WHILE},
		{"return", token.RETURN},
		{"extern", token.EXTERN},
		{"if", token.IF},
		{"elif", token.ELIF},
		{"else", token.ELSE},
		{"not", token.NOT},

		// Types
		{"bool", token.BOOL_TYPE},

		{"int", token.INT_TYPE},
		{"i8", token.I8_TYPE},
		{"i16", token.I16_TYPE},
		{"i32", token.I32_TYPE},
		{"i64", token.I64_TYPE},

		{"uint", token.UINT_TYPE},
		{"u8", token.U8_TYPE},
		{"u16", token.U16_TYPE},
		{"u32", token.U32_TYPE},
		{"u64", token.U64_TYPE},

		// Other tokens
		{"(", token.OPEN_PAREN},
		{")", token.CLOSE_PAREN},
		{"{", token.OPEN_CURLY},
		{"}", token.CLOSE_CURLY},
		{",", token.COMMA},
		{";", token.SEMICOLON},
		{".", token.DOT},
		{"..", token.DOT_DOT},
		{"...", token.DOT_DOT_DOT},
		{"=", token.EQUAL},
		{":=", token.COLON_EQUAL},
		{"!=", token.BANG_EQUAL},
		{"==", token.EQUAL_EQUAL},
		{">", token.GREATER},
		{">=", token.GREATER_EQ},
		{"<", token.LESS},
		{"<=", token.LESS_EQ},
		{"+", token.PLUS},
		{"-", token.MINUS},
		{"*", token.STAR},
		{"/", token.SLASH},
	}

	for _, test := range tests {
		t.Run(fmt.Sprintf("TestTokenKind('%q')", test.lexeme), func(t *testing.T) {
			collector := diagnostics.New()

			src := []byte(test.lexeme)
			lexer := New(filename, src, collector)

			tokenResult, err := lexer.Tokenize()
			if err != nil {
				t.Errorf("unexpected error '%v'", err)
			}

			if len(tokenResult) != 2 {
				t.Errorf("expected len(tokenResult) == 2, but got %q", len(tokenResult))
			}
			if tokenResult[1].Kind != token.EOF {
				t.Errorf("expected last token to be EOF, but got %q", tokenResult[1].Kind)
			}
			if tokenResult[0].Kind != test.kind {
				t.Errorf("expected token to be %q, but got %q", test.kind, tokenResult[0].Kind)
			}
		})
	}
}

type tokenPosTest struct {
	input     string
	positions []token.Pos
}

func TestTokenPos(t *testing.T) {
	filename := "test.tt"

	tests := []*tokenPosTest{
		{";", []token.Pos{
			{Filename: "test.tt", Line: 1, Column: 1},
			{Filename: "test.tt", Line: 1, Column: 2}},
		},
		{";\n;", []token.Pos{
			{Filename: "test.tt", Line: 1, Column: 1},
			{Filename: "test.tt", Line: 2, Column: 1},
			{Filename: "test.tt", Line: 2, Column: 2}},
		},
		{"fn\nhello world\n;", []token.Pos{
			{Filename: "test.tt", Line: 1, Column: 1},
			{Filename: "test.tt", Line: 2, Column: 1},
			{Filename: "test.tt", Line: 2, Column: 7},
			{Filename: "test.tt", Line: 3, Column: 1},
			{Filename: "test.tt", Line: 3, Column: 2}},
		},
	}

	for _, test := range tests {
		t.Run(fmt.Sprintf("TestTokenPos(%q)", test.input), func(t *testing.T) {
			collector := diagnostics.New()

			src := []byte(test.input)
			lexer := New(filename, src, collector)

			tokenResult, err := lexer.Tokenize()
			if err != nil {
				t.Errorf("unexpected error '%v'", err)
			}

			if len(tokenResult) == 1 && tokenResult[0].Kind == token.EOF {
				t.Errorf("expected at least one token, but only got EOF")
			}

			if len(tokenResult) != len(test.positions) {
				t.Errorf(
					"expected len(tokenResult) == len(expectedPos.positions), expected %d, but got %d",
					len(tokenResult),
					len(test.positions),
				)
			}

			for i, expectedPos := range test.positions {
				actualPos := tokenResult[i].Pos
				if expectedPos != actualPos {
					t.Errorf(
						"expected token position to be the same, expected %q, but got %q",
						expectedPos,
						actualPos,
					)
				}
			}
		})
	}
}

type tokenIdentTest struct {
	lexeme string
	isId   bool
}

func TestIsIdentifier(t *testing.T) {
	filename := "test.tt"

	tests := []*tokenIdentTest{
		{"hello", true},
		{"world", true},
		{"foobar", true},
		{"hello_world_", true},

		// TODO: add support to Unicode
		// {"foo६४", true},
		// {"a۰۱۸", true},
		// {"bar９８７６", true},
		// {"ŝ", true},
		// {"ŝfoo", true},

		{"a123456789", true},
		{"123456789", false},
		// TODO: add float here

		{"true", false},
		{"false", false},
		{"fn", false},
		{"for", false},
		{"while", false},
		{"return", false},
		{"if", false},
		{"elif", false},
		{"else", false},

		{"bool", false},
		{"int", false},
		{"i8", false},
		{"i16", false},
		{"i32", false},
		{"i64", false},
		{"uint", false},
		{"u8", false},
		{"u16", false},
		{"u32", false},
		{"u64", false},
	}

	for _, test := range tests {
		t.Run(fmt.Sprintf("TestIsIdentifier('%q')", test.lexeme), func(t *testing.T) {
			collector := diagnostics.New()

			src := []byte(test.lexeme)
			lexer := New(filename, src, collector)

			tokenResult, err := lexer.Tokenize()
			if err != nil {
				t.Errorf("unexpected error '%v'", err)
			}
			if len(tokenResult) != 2 {
				t.Errorf("expected a single token, but got %d", len(tokenResult))
			}
			if tokenResult[1].Kind != token.EOF {
				t.Errorf("expected last token to be EOF, but got %q", tokenResult[1].Kind)
			}
			if tokenResult[0].Kind != token.ID && test.isId {
				t.Errorf("expecting to be an identifier, but got %q", tokenResult[0].Kind)
			}
		})
	}
}

type tokenLiteralTest struct {
	lexeme      string
	literalKind token.Kind
}

func TestIsLiteral(t *testing.T) {
	filename := "test.tt"

	tests := []*tokenLiteralTest{
		{"1", token.INTEGER_LITERAL},
		{"2", token.INTEGER_LITERAL},
		{"3", token.INTEGER_LITERAL},
		{"4", token.INTEGER_LITERAL},
		{"5", token.INTEGER_LITERAL},
		{"6", token.INTEGER_LITERAL},
		{"7", token.INTEGER_LITERAL},
		{"8", token.INTEGER_LITERAL},
		{"9", token.INTEGER_LITERAL},
		{"123456789", token.INTEGER_LITERAL},
		// TODO: add float here
		{"\"Hello world\"", token.STRING_LITERAL},
		{"true", token.TRUE_BOOL_LITERAL},
		{"false", token.FALSE_BOOL_LITERAL},
	}

	for _, test := range tests {
		t.Run(fmt.Sprintf("TestIsLiteral('%q')", test.lexeme), func(t *testing.T) {
			collector := diagnostics.New()

			src := []byte(test.lexeme)
			lexer := New(filename, src, collector)
			tokenResult, err := lexer.Tokenize()
			if err != nil {
				t.Errorf("unexpected error '%v'", err)
			}

			if len(tokenResult) != 2 {
				t.Errorf("expected a single token, but got %d", len(tokenResult))
			}
			if tokenResult[1].Kind != token.EOF {
				t.Errorf("expected last token to be EOF, but got %q", tokenResult[1].Kind)
			}
			if tokenResult[0].Kind != test.literalKind {
				t.Errorf("expected to be a %q, but got %q", test.literalKind, tokenResult[0].Kind)
			}
		})
	}
}

type lexicalErrorTest struct {
	input string
	diags []diagnostics.Diag
}

func TestLexicalErrors(t *testing.T) {
	filename := "test.tt"

	tests := []lexicalErrorTest{
		{
			input: "!",
			diags: []diagnostics.Diag{
				{
					Message: "test.tt:1:1: invalid character !",
				},
			},
		},
		{
			input: "!!",
			diags: []diagnostics.Diag{
				{
					Message: "test.tt:1:1: invalid character !",
				},
			},
		},
		{
			input: ":",
			diags: []diagnostics.Diag{
				{
					Message: "test.tt:1:1: invalid character :",
				},
			},
		},
		{
			input: "::",
			diags: []diagnostics.Diag{
				{
					Message: "test.tt:1:1: invalid character :",
				},
			},
		},
		{
			input: "?",
			diags: []diagnostics.Diag{
				{
					Message: "test.tt:1:1: invalid character ?",
				},
			},
		},
		{
			input: "\"Unterminated string literal here",
			diags: []diagnostics.Diag{
				{
					Message: "test.tt:1:1: unterminated string literal",
				},
			},
		},
		{
			input: "\"",
			diags: []diagnostics.Diag{
				{
					Message: "test.tt:1:1: unterminated string literal",
				},
			},
		},
	}

	for _, test := range tests {
		t.Run(fmt.Sprintf("TestLexicalErrors('%s')", test.input), func(t *testing.T) {
			collector := diagnostics.New()

			src := []byte(test.input)
			lex := New(filename, src, collector)
			_, err := lex.Tokenize()
			if err == nil {
				t.Fatal("expected to have lexical errors, but got nothing")
			}

			if len(test.diags) != len(lex.collector.Diags) {
				t.Fatalf(
					"expected to have %d diag(s), but got %d",
					len(test.diags),
					len(lex.collector.Diags),
				)
			}

			if !reflect.DeepEqual(test.diags, lex.collector.Diags) {
				t.Fatalf("\nexpected diags: %v\ngot diags: %v\n", test.diags, lex.collector)
			}
		})
	}
}
