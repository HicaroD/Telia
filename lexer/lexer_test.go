package lexer

import (
	"bufio"
	"strings"
	"testing"

	"github.com/HicaroD/telia-lang/lexer/token/kind"
)

type tokenTest struct {
	lexeme string
	kind   kind.TokenKind
}

var tokenKinds []*tokenTest = []*tokenTest{
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

func TestTokenPos(t *testing.T) {}

func TestIsIdentifier(t *testing.T) {}

func TestIsLiteral(t *testing.T) {}
