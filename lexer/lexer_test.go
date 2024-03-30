package lexer

import (
	"bufio"
	"strings"
	"testing"

	"github.com/HicaroD/telia-lang/lexer/token"
	"github.com/HicaroD/telia-lang/lexer/token/kind"
)

type tokenKindTest struct {
	lexeme string
	kind   kind.TokenKind
}

var tokenKinds []*tokenKindTest = []*tokenKindTest{
	// Keywords
	{"fn", kind.FN},
	{"return", kind.RETURN},
	{"extern", kind.EXTERN},

	// Types
	{"bool", kind.BOOL_TYPE},
	{"i8", kind.I8_TYPE},
	{"i16", kind.I16_TYPE},
	{"i32", kind.I32_TYPE},
	{"i64", kind.I64_TYPE},
	{"i128", kind.I128_TYPE},

	{"(", kind.OPEN_PAREN},
	{")", kind.CLOSE_PAREN},
	{"{", kind.OPEN_CURLY},
	{"}", kind.CLOSE_CURLY},
	{",", kind.COMMA},
	{";", kind.SEMICOLON},
	{".", kind.DOT},
	{"..", kind.DOT_DOT},
	{"...", kind.DOT_DOT_DOT},
	{"*", kind.STAR},
}

func TestTokenKinds(t *testing.T) {
	testFilename := "test.tt"

	for _, expectedToken := range tokenKinds {
		reader := bufio.NewReader(strings.NewReader(expectedToken.lexeme))
		lexer := NewLexer(testFilename, reader)
		tokenResult := lexer.Tokenize()

		if len(tokenResult) != 2 {
			t.Errorf("TestTokenKind(%q): expected len(tokenResult) == 2, but got %q", expectedToken.lexeme, len(tokenResult))
		}
		if tokenResult[1].Kind != kind.EOF {
			t.Errorf("TestTokenKind(%q): expected last token to be EOF, but got %q", expectedToken.lexeme, tokenResult[1].Kind)
		}
		if tokenResult[0].Kind != expectedToken.kind {
			t.Errorf("TestTokenKind(%q): expected token to be %q, but got %q", expectedToken.lexeme, expectedToken.kind, tokenResult[0].Kind)
		}
	}
}

type tokenPosTest struct {
	input     string
	positions []token.Position
}

var tokenPos []*tokenPosTest = []*tokenPosTest{
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

func TestTokenPos(t *testing.T) {
	testFilename := "test.tt"

	for _, expectedPos := range tokenPos {
		reader := bufio.NewReader(strings.NewReader(expectedPos.input))
		lexer := NewLexer(testFilename, reader)
		tokenResult := lexer.Tokenize()

		if len(tokenResult) == 1 && tokenResult[0].Kind == kind.EOF {
			t.Errorf("TestTokenPos(%q): expected at least one token, but only got EOF", expectedPos.input)
		}

		if len(tokenResult) != len(expectedPos.positions) {
			t.Errorf("TestTokenPos(%q): expected len(tokenResult) == len(expectedPos.positions), expected %d, but got %d", expectedPos.input, len(tokenResult), len(expectedPos.positions))
		}

		for i, expectedPos := range expectedPos.positions {
			lexeme := tokenResult[i].Lexeme
			actualPos := tokenResult[i].Position
			if expectedPos != actualPos {
				t.Errorf("TestTokenPos(%q): expected token position to be the same, expected %q, but got %q", lexeme, expectedPos, actualPos)
			}
		}
	}
}

type tokenIdentTest struct {
	lexeme      string
	expectingId bool
}

var tokenIdent []*tokenIdentTest = []*tokenIdentTest{
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
	{"fn", false},
	{"return", false},
}

func TestIsIdentifier(t *testing.T) {
	testFilename := "test.tt"

	for _, expectedTokenIdent := range tokenIdent {
		reader := bufio.NewReader(strings.NewReader(expectedTokenIdent.lexeme))
		lexer := NewLexer(testFilename, reader)
		tokenResult := lexer.Tokenize()

		if len(tokenResult) != 2 {
			t.Errorf("TestIsIdentifier(%q): expected a single token, but got %d", expectedTokenIdent.lexeme, len(tokenResult))
		}
		if tokenResult[1].Kind != kind.EOF {
			t.Errorf("TestIsIdentifier(%q): expected last token to be EOF, but got %q", expectedTokenIdent.lexeme, tokenResult[1].Kind)
		}
		if tokenResult[0].Kind != kind.ID && expectedTokenIdent.expectingId {
			t.Errorf("TestIsIdentifier(%q): expected to be an identifier, but got %q", expectedTokenIdent.lexeme, tokenResult[0].Kind)
		}
	}
}

type tokenLiteralTest struct {
	lexeme      string
	literalKind kind.TokenKind
}

var tokenLiterals []*tokenLiteralTest = []*tokenLiteralTest{
	// TODO: add float here
	// TODO: add bool here
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
	{"\"Hello world\"", kind.STRING_LITERAL},
}

func TestIsLiteral(t *testing.T) {
	testFilename := "test.tt"
	for _, expectedTokenLiteral := range tokenLiterals {
		reader := bufio.NewReader(strings.NewReader(expectedTokenLiteral.lexeme))
		lexer := NewLexer(testFilename, reader)
		tokenResult := lexer.Tokenize()

		if len(tokenResult) != 2 {
			t.Errorf("TestIsIdentifier(%q): expected a single token, but got %d", expectedTokenLiteral.lexeme, len(tokenResult))
		}
		if tokenResult[1].Kind != kind.EOF {
			t.Errorf("TestIsIdentifier(%q): expected last token to be EOF, but got %q", expectedTokenLiteral.lexeme, tokenResult[1].Kind)
		}
		if tokenResult[0].Kind != expectedTokenLiteral.literalKind {
			t.Errorf("TestIsIdentifier(%q): expected to be a %q, but got %q", expectedTokenLiteral.lexeme, expectedTokenLiteral.literalKind, tokenResult[0].Kind)
		}
	}
}
