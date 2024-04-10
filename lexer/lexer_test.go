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
	{"if", kind.IF},
	{"elif", kind.ELIF},
	{"else", kind.ELSE},

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
	{"*", kind.STAR},
	{"=", kind.EQUAL},
	{"==", kind.EQUAL_EQUAL},
	{"-", kind.MINUS},
	{":=", kind.COLON_EQUAL},
}

func TestTokenKinds(t *testing.T) {
	testFilename := "test.tt"

	for _, expectedToken := range tokenKinds {
		reader := bufio.NewReader(strings.NewReader(expectedToken.lexeme))
		lexer := New(testFilename, reader)
		tokenResult, err := lexer.Tokenize()
		if err != nil {
			t.Errorf("TestTokenKind(%q): unexpected error '%v'", expectedToken.lexeme, err)
		}

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
		lexer := New(testFilename, reader)
		tokenResult, err := lexer.Tokenize()
		if err != nil {
			t.Errorf("TestTokenPos(%q): unexpected error '%v'", expectedPos.input, err)
		}

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
	lexeme string
	isId   bool
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

func TestIsIdentifier(t *testing.T) {
	testFilename := "test.tt"

	for _, expectedTokenIdent := range tokenIdent {
		reader := bufio.NewReader(strings.NewReader(expectedTokenIdent.lexeme))
		lexer := New(testFilename, reader)
		tokenResult, err := lexer.Tokenize()
		if err != nil {
			t.Errorf("TestIsIdentifier(%q): unexpected error '%v'", expectedTokenIdent.lexeme, err)
		}
		if len(tokenResult) != 2 {
			t.Errorf("TestIsIdentifier(%q): expected a single token, but got %d", expectedTokenIdent.lexeme, len(tokenResult))
		}
		if tokenResult[1].Kind != kind.EOF {
			t.Errorf("TestIsIdentifier(%q): expected last token to be EOF, but got %q", expectedTokenIdent.lexeme, tokenResult[1].Kind)
		}
		if tokenResult[0].Kind != kind.ID && expectedTokenIdent.isId {
			t.Errorf("TestIsIdentifier(%q): expecting to be an identifier, but got %q", expectedTokenIdent.lexeme, tokenResult[0].Kind)
		}
	}
}

type tokenLiteralTest struct {
	lexeme      string
	literalKind kind.TokenKind
}

var tokenLiterals []*tokenLiteralTest = []*tokenLiteralTest{
	// TODO: add float here
	{"1", kind.INTEGER_LITERAL},
	{"2", kind.INTEGER_LITERAL},
	{"3", kind.INTEGER_LITERAL},
	{"4", kind.INTEGER_LITERAL},
	{"5", kind.INTEGER_LITERAL},
	{"6", kind.INTEGER_LITERAL},
	{"7", kind.INTEGER_LITERAL},
	{"8", kind.INTEGER_LITERAL},
	{"9", kind.INTEGER_LITERAL},
	{"-1", kind.NEGATIVE_INTEGER_LITERAL},
	{"-2", kind.NEGATIVE_INTEGER_LITERAL},
	{"-3", kind.NEGATIVE_INTEGER_LITERAL},
	{"-4", kind.NEGATIVE_INTEGER_LITERAL},
	{"-5", kind.NEGATIVE_INTEGER_LITERAL},
	{"-6", kind.NEGATIVE_INTEGER_LITERAL},
	{"-7", kind.NEGATIVE_INTEGER_LITERAL},
	{"-8", kind.NEGATIVE_INTEGER_LITERAL},
	{"-9", kind.NEGATIVE_INTEGER_LITERAL},
	{"123456789", kind.INTEGER_LITERAL},
	{"-123456789", kind.NEGATIVE_INTEGER_LITERAL},
	{"\"Hello world\"", kind.STRING_LITERAL},
	{"true", kind.TRUE_BOOL_LITERAL},
	{"false", kind.FALSE_BOOL_LITERAL},
}

func TestIsLiteral(t *testing.T) {
	testFilename := "test.tt"
	for _, expectedTokenLiteral := range tokenLiterals {
		reader := bufio.NewReader(strings.NewReader(expectedTokenLiteral.lexeme))
		lexer := New(testFilename, reader)
		tokenResult, err := lexer.Tokenize()
		if err != nil {
			t.Errorf("TestIsIdentifier(%q): unexpected error '%v'", expectedTokenLiteral.lexeme, err)
		}

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
