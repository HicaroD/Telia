package lexer

import (
	"fmt"
	"log"
	"os"
	"unicode"

	"github.com/HicaroD/Telia/internal/ast"
	"github.com/HicaroD/Telia/internal/diagnostics"
	"github.com/HicaroD/Telia/internal/lexer/token"
)

const eof = '\000'

type Lexer struct {
	Loc           *ast.Loc
	Collector     *diagnostics.Collector
	StrictNewline bool

	src    []byte
	offset int
	pos    token.Pos
}

func New(loc *ast.Loc, src []byte, collector *diagnostics.Collector) *Lexer {
	lexer := new(Lexer)

	lexer.Loc = loc
	lexer.Collector = collector
	lexer.StrictNewline = true
	lexer.pos = token.NewPosition(loc.Name, 1, 1)
	lexer.src = src
	lexer.offset = 0

	return lexer
}

func NewFromFilePath(loc *ast.Loc, collector *diagnostics.Collector) (*Lexer, error) {
	src, err := os.ReadFile(loc.Path)
	if err != nil {
		return nil, err
	}
	l := New(loc, src, collector)
	return l, nil
}

func (lex *Lexer) Filename() string { return lex.pos.Filename }

func (lex *Lexer) Peek() *token.Token {
	prevPos := lex.pos
	prevOffset := lex.offset

	token := lex.Next()

	lex.pos.SetPosition(prevPos)
	lex.offset = prevOffset
	return token
}

func (lex *Lexer) Peek1() *token.Token {
	prevPos := lex.pos
	prevOffset := lex.offset

	var token *token.Token

	_ = lex.Next()
	token = lex.Next()

	lex.pos.SetPosition(prevPos)
	lex.offset = prevOffset

	return token
}

func (lex *Lexer) Skip() {
	lex.Next()
}

func (lex *Lexer) NextIs(expectedKind token.Kind) bool {
	token := lex.Peek()
	return token.Kind == expectedKind
}

func (lex *Lexer) Next() *token.Token {
	lex.skipWhitespace()
	character := lex.peekChar()

	tok := &token.Token{}
	tok.Kind = token.INVALID

	if character == eof {
		lex.consumeTokenNoLex(tok, token.EOF)
		return tok
	}

	token := lex.getToken(tok, character)
	return token
}

// Useful for testing
func (lex *Lexer) Tokenize() ([]*token.Token, error) {
	var tokens []*token.Token
	for {
		tok := lex.Next()
		if tok.Kind == token.INVALID {
			return nil, diagnostics.COMPILER_ERROR_FOUND
		}
		tokens = append(tokens, tok)
		if tok.Kind == token.EOF {
			break
		}
	}
	return tokens, nil
}

func (lex *Lexer) getToken(tok *token.Token, ch byte) *token.Token {
	switch ch {
	case '\n':
		lex.consumeTokenNoLex(tok, token.NEWLINE)
		lex.nextChar()
	case '(':
		lex.consumeTokenNoLex(tok, token.OPEN_PAREN)
		lex.nextChar()
	case ')':
		lex.consumeTokenNoLex(tok, token.CLOSE_PAREN)
		lex.nextChar()
	case '{':
		lex.consumeTokenNoLex(tok, token.OPEN_CURLY)
		lex.nextChar()
	case '}':
		lex.consumeTokenNoLex(tok, token.CLOSE_CURLY)
		lex.nextChar()
	case '"':
		// TODO: implement raw strings
		isRaw := false
		lex.getStringLit(tok, isRaw)
	case ',':
		lex.consumeTokenNoLex(tok, token.COMMA)
		lex.nextChar()
	case ';':
		lex.consumeTokenNoLex(tok, token.SEMICOLON)
		lex.nextChar()
	case '+':
		lex.consumeTokenNoLex(tok, token.PLUS)
		lex.nextChar()
	case '-':
		lex.consumeTokenNoLex(tok, token.MINUS)
		lex.nextChar()
	case '*':
		lex.consumeTokenNoLex(tok, token.STAR)
		lex.nextChar()
	case '/':
		lex.consumeTokenNoLex(tok, token.SLASH)
		lex.nextChar()
	case '#':
		lex.consumeTokenNoLex(tok, token.SHARP)
		lex.nextChar()
	case '@':
		lex.consumeTokenNoLex(tok, token.AT)
		lex.nextChar()
	case '[':
		lex.consumeTokenNoLex(tok, token.OPEN_BRACKET)
		lex.nextChar()
	case ']':
		lex.consumeTokenNoLex(tok, token.CLOSE_BRACKET)
		lex.nextChar()
	case '!':
		tok.Pos = lex.pos
		lex.nextChar()

		invalidCharacter := diagnostics.Diag{
			Message: fmt.Sprintf(
				"%s:%d:%d: invalid character !",
				tok.Pos.Filename,
				tok.Pos.Line,
				tok.Pos.Column,
			),
		}

		next := lex.peekChar()
		if next == eof {
			lex.Collector.ReportAndSave(invalidCharacter)
			return tok
		}
		if next == '=' {
			lex.nextChar()
			tok.Kind = token.BANG_EQUAL
			return tok
		}

		lex.Collector.ReportAndSave(invalidCharacter)
		return tok
	case '>':
		tok.Pos = lex.pos

		tok.Kind = token.GREATER
		lex.nextChar() // >

		next := lex.peekChar()
		if next != '=' {
			return tok
		}
		lex.nextChar() // =
		tok.Kind = token.GREATER_EQ
	case '<':
		tok.Kind = token.LESS
		tok.Pos = lex.pos
		lex.nextChar() // >

		next := lex.peekChar()
		if next != '=' {
			return tok
		}
		lex.nextChar() // =
		tok.Kind = token.LESS_EQ
	case '=':
		tok.Kind = token.EQUAL
		tok.Pos = lex.pos
		lex.nextChar() // =

		next := lex.peekChar()
		if next != '=' {
			return tok
		}
		lex.nextChar() // =
		tok.Kind = token.EQUAL_EQUAL
	case '.':
		tok.Kind = token.DOT
		tok.Pos = lex.pos
		lex.nextChar() // .

		next := lex.peekChar()
		if next != '.' {
			return tok
		}
		lex.nextChar() // .
		tok.Kind = token.DOT_DOT

		next = lex.peekChar()
		if next != '.' {
			return tok
		}
		lex.nextChar() // .
		tok.Kind = token.DOT_DOT_DOT
	case ':':
		tok.Pos = lex.pos
		lex.nextChar() // :

		invalidCharacter := diagnostics.Diag{
			Message: fmt.Sprintf(
				"%s:%d:%d: invalid character :",
				tok.Pos.Filename,
				tok.Pos.Line,
				tok.Pos.Column,
			),
		}

		next := lex.peekChar()
		if next == eof {
			lex.Collector.ReportAndSave(invalidCharacter)
			return tok
		}

		if next == '=' {
			lex.nextChar() // =
			tok.Kind = token.COLON_EQUAL
			break
		} else if next == ':' {
			lex.nextChar() // :
			tok.Kind = token.COLON_COLON
			break
		}

		lex.Collector.ReportAndSave(invalidCharacter)
	default:
		if unicode.IsLetter(rune(ch)) || ch == '_' {
			lex.getIdOrKeyword(tok)
		} else if ch >= '0' && ch <= '9' {
			lex.getNumberLit(tok)
		} else {
			tokenPosition := lex.pos
			invalidCharacter := diagnostics.Diag{
				Message: fmt.Sprintf("%s:%d:%d: invalid character %c", tokenPosition.Filename, tokenPosition.Line, tokenPosition.Column, ch),
			}
			lex.Collector.ReportAndSave(invalidCharacter)
		}
	}
	return tok
}

func (lex *Lexer) getStringLit(tok *token.Token, isRaw bool) *token.Token {
	tok.Pos = lex.pos
	lex.nextChar() // "

	var str []byte
	for {
		ch := lex.peekChar()
		if ch == eof || ch == '"' {
			break
		}

		if ch == '\\' && !isRaw {
			lex.nextChar()
			escapeSym := lex.peekChar()

			var escape byte

			switch escapeSym {
			case 'n':
				escape = '\n'
			case 't':
				escape = '\t'
			case '\\':
				escape = '\\'
			case '"':
				escape = '"'
			default:
				// TODO(errors)
				log.Fatalf("invalid escape sequence: %b", escape)
				return tok
			}
			str = append(str, escape)
		} else {
			str = append(str, ch)
		}

		lex.nextChar()
	}

	ch := lex.peekChar()
	if ch != '"' {
		unterminatedStringLiteral := diagnostics.Diag{
			Message: fmt.Sprintf(
				"%s:%d:%d: unterminated string literal",
				tok.Pos.Filename,
				tok.Pos.Line,
				tok.Pos.Column,
			),
		}
		lex.Collector.ReportAndSave(unterminatedStringLiteral)
		return tok
	}
	if ch == '"' {
		lex.nextChar()
	}

	tok.Kind = token.UNTYPED_STRING
	tok.Lexeme = str
	return tok
}

func (lex *Lexer) getNumberLit(tok *token.Token) {
	var dotFound, dotRepeated bool
	numberType := token.UNTYPED_INT

	number := lex.readWhile(
		func(chr byte) bool {
			if chr == '.' {
				if dotFound {
					dotRepeated = true
					return false
				}
				numberType = token.UNTYPED_FLOAT
				dotFound = true
				return true
			}
			return (chr >= '0' && chr <= '9') || chr == '_'
		},
	)

	// TODO(errors)
	if dotRepeated {
		tok.Kind = token.INVALID
		log.Fatal("error: invalid float format")
		return
	}

	tok.Kind = numberType
	tok.Lexeme = number
}

func (lex *Lexer) getIdOrKeyword(tok *token.Token) {
	identifier := lex.readWhile(
		func(chr byte) bool { return unicode.IsNumber(rune(chr)) || unicode.IsLetter(rune(chr)) || chr == '_' },
	)
	tok.Kind = token.ID
	tok.Lexeme = identifier
	keyword, ok := token.KEYWORDS[string(identifier)]
	if ok {
		tok.Kind = keyword
	}
}

func (lex *Lexer) consumeTokenNoLex(tok *token.Token, kind token.Kind) {
	tok.Lexeme = nil
	tok.Kind = kind
	tok.Pos = lex.pos
}

func (lex *Lexer) skipWhitespace() {
	lex.readWhile(func(ch byte) bool {
		return ch == ' ' || ch == '\t' || ch == '\r' || (ch == '\n' && !lex.StrictNewline)
	})
}

func (lex *Lexer) readWhile(isValid func(byte) bool) []byte {
	var start, end int
	start = lex.offset

	for {
		character := lex.peekChar()
		if character == eof {
			break
		}

		if isValid(character) {
			lex.nextChar()
		} else {
			break
		}
	}

	end = lex.offset

	return lex.src[start:end]
}

func (lex *Lexer) nextChar() byte {
	if lex.offset > len(lex.src) {
		return eof
	}
	character := lex.src[lex.offset]
	lex.pos.Move(character)
	lex.offset++
	return character
}

func (lex *Lexer) peekChar() byte {
	if lex.offset >= len(lex.src) {
		return eof
	}
	character := lex.src[lex.offset]
	return character
}
