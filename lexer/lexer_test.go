package lexer

import (
	"bufio"
	"fmt"
	"strings"
	"testing"

	"github.com/HicaroD/Telia/lexer/token"
	"github.com/HicaroD/Telia/lexer/token/kind"
)

type tokenKindTest struct {
	lexeme string
	kind   kind.TokenKind
}

func TestTokenKinds(t *testing.T) {
	filename := "test.tt"

	tests := []*tokenKindTest{
		// Keywords
		{"fn", kind.FN},
		{"return", kind.RETURN},
		{"extern", kind.EXTERN},
		{"if", kind.IF},
		{"elif", kind.ELIF},
		{"else", kind.ELSE},
		{"not", kind.NOT},

		// Types
		{"bool", kind.BOOL_TYPE},
		{"i8", kind.I8_TYPE},
		{"i16", kind.I16_TYPE},
		{"i32", kind.I32_TYPE},
		{"i64", kind.I64_TYPE},
		{"u8", kind.U8_TYPE},
		{"u16", kind.U16_TYPE},
		{"u32", kind.U32_TYPE},
		{"u64", kind.U64_TYPE},

		// Other tokens
		{"(", kind.OPEN_PAREN},
		{")", kind.CLOSE_PAREN},
		{"{", kind.OPEN_CURLY},
		{"}", kind.CLOSE_CURLY},
		{",", kind.COMMA},
		{";", kind.SEMICOLON},
		{".", kind.DOT},
		{"..", kind.DOT_DOT},
		{"...", kind.DOT_DOT_DOT},
		{"=", kind.EQUAL},
		{":=", kind.COLON_EQUAL},
		{"!=", kind.BANG_EQUAL},
		{"==", kind.EQUAL_EQUAL},
		{">", kind.GREATER},
		{">=", kind.GREATER_EQ},
		{"<", kind.LESS},
		{"<=", kind.LESS_EQ},
		{"+", kind.PLUS},
		{"-", kind.MINUS},
		{"*", kind.STAR},
		{"/", kind.SLASH},
	}

	for _, test := range tests {
		t.Run(fmt.Sprintf("TestTokenKind('%q')", test.lexeme), func(t *testing.T) {
			reader := bufio.NewReader(strings.NewReader(test.lexeme))
			lexer := New(filename, reader)
			tokenResult, err := lexer.Tokenize()
			if err != nil {
				t.Errorf("unexpected error '%v'", err)
			}

			if len(tokenResult) != 2 {
				t.Errorf("expected len(tokenResult) == 2, but got %q", len(tokenResult))
			}
			if tokenResult[1].Kind != kind.EOF {
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
	positions []token.Position
}

func TestTokenPos(t *testing.T) {
	filename := "test.tt"

	tests := []*tokenPosTest{
		{";", []token.Position{
			{Filename: "test.tt", Line: 1, Column: 1},
			{Filename: "test.tt", Line: 1, Column: 2}},
		},
		{";\n;", []token.Position{
			{Filename: "test.tt", Line: 1, Column: 1},
			{Filename: "test.tt", Line: 2, Column: 1},
			{Filename: "test.tt", Line: 2, Column: 2}},
		},
		{"fn\nhello world\n;", []token.Position{
			{Filename: "test.tt", Line: 1, Column: 1},
			{Filename: "test.tt", Line: 2, Column: 1},
			{Filename: "test.tt", Line: 2, Column: 7},
			{Filename: "test.tt", Line: 3, Column: 1},
			{Filename: "test.tt", Line: 3, Column: 2}},
		},
	}

	for _, test := range tests {
		t.Run(fmt.Sprintf("TestTokenPos(%q)", test.input), func(t *testing.T) {
			reader := bufio.NewReader(strings.NewReader(test.input))
			lexer := New(filename, reader)
			tokenResult, err := lexer.Tokenize()
			if err != nil {
				t.Errorf("unexpected error '%v'", err)
			}

			if len(tokenResult) == 1 && tokenResult[0].Kind == kind.EOF {
				t.Errorf("expected at least one token, but only got EOF")
			}

			if len(tokenResult) != len(test.positions) {
				t.Errorf("expected len(tokenResult) == len(expectedPos.positions), expected %d, but got %d", len(tokenResult), len(test.positions))
			}

			for i, expectedPos := range test.positions {
				actualPos := tokenResult[i].Position
				if expectedPos != actualPos {
					t.Errorf("expected token position to be the same, expected %q, but got %q", expectedPos, actualPos)
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
		{"foo६४", true},
		{"a۰۱۸", true},
		{"bar９８７６", true},
		{"ŝ", true},
		{"ŝfoo", true},
		{"a123456789", true}, // NOTE: starts with "a"
		{"123456789", false},
		// TODO: add float here

		{"true", false},
		{"false", false},
		{"fn", false},
		{"return", false},
		{"if", false},
		{"elif", false},
		{"else", false},

		{"bool", false},
		{"i8", false},
		{"i16", false},
		{"i32", false},
		{"i64", false},
		{"u8", false},
		{"u16", false},
		{"u32", false},
		{"u64", false},
	}

	for _, test := range tests {
		t.Run(fmt.Sprintf("TestIsIdentifier('%q')", test.lexeme), func(t *testing.T) {
			reader := bufio.NewReader(strings.NewReader(test.lexeme))
			lexer := New(filename, reader)
			tokenResult, err := lexer.Tokenize()
			if err != nil {
				t.Errorf("unexpected error '%v'", err)
			}
			if len(tokenResult) != 2 {
				t.Errorf("expected a single token, but got %d", len(tokenResult))
			}
			if tokenResult[1].Kind != kind.EOF {
				t.Errorf("expected last token to be EOF, but got %q", tokenResult[1].Kind)
			}
			if tokenResult[0].Kind != kind.ID && test.isId {
				t.Errorf("expecting to be an identifier, but got %q", tokenResult[0].Kind)
			}
		})
	}
}

type tokenLiteralTest struct {
	lexeme      string
	literalKind kind.TokenKind
}

func TestIsLiteral(t *testing.T) {
	filename := "test.tt"

	tests := []*tokenLiteralTest{
		{"1", kind.INTEGER_LITERAL},
		{"2", kind.INTEGER_LITERAL},
		{"3", kind.INTEGER_LITERAL},
		{"4", kind.INTEGER_LITERAL},
		{"5", kind.INTEGER_LITERAL},
		{"6", kind.INTEGER_LITERAL},
		{"7", kind.INTEGER_LITERAL},
		{"8", kind.INTEGER_LITERAL},
		{"9", kind.INTEGER_LITERAL},
		{"123456789", kind.INTEGER_LITERAL},
		// TODO: add float here
		{"\"Hello world\"", kind.STRING_LITERAL},
		{"true", kind.TRUE_BOOL_LITERAL},
		{"false", kind.FALSE_BOOL_LITERAL},
	}

	for _, test := range tests {
		t.Run(fmt.Sprintf("TestIsLiteral('%q')", test.lexeme), func(t *testing.T) {
			reader := bufio.NewReader(strings.NewReader(test.lexeme))
			lexer := New(filename, reader)
			tokenResult, err := lexer.Tokenize()
			if err != nil {
				t.Errorf("unexpected error '%v'", err)
			}

			if len(tokenResult) != 2 {
				t.Errorf("expected a single token, but got %d", len(tokenResult))
			}
			if tokenResult[1].Kind != kind.EOF {
				t.Errorf("expected last token to be EOF, but got %q", tokenResult[1].Kind)
			}
			if tokenResult[0].Kind != test.literalKind {
				t.Errorf("expected to be a %q, but got %q", test.literalKind, tokenResult[0].Kind)
			}
		})
	}
}
