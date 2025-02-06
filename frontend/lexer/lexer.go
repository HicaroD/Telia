package lexer

import (
	"fmt"
	"os"
	"unicode"

	"github.com/HicaroD/Telia/diagnostics"
	"github.com/HicaroD/Telia/frontend/lexer/token"
)

const eof = '\000'

type Lexer struct {
	ParentDirName string
	Path          string

	src    []byte
	offset int
	pos    token.Pos

	collector *diagnostics.Collector
}

func New(path string, src []byte, collector *diagnostics.Collector) *Lexer {
	lexer := new(Lexer)

	lexer.ParentDirName = ""
	lexer.Path = path
	lexer.collector = collector
	lexer.src = src
	lexer.offset = 0
	lexer.pos = token.NewPosition(path, 1, 1)

	return lexer
}

func NewFromFilePath(parentDirName, path string, collector *diagnostics.Collector) (*Lexer, error) {
	src, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	l := New(path, src, collector)
	l.ParentDirName = parentDirName
	return l, nil
}

func (lex *Lexer) Peek() *token.Token {
	prevPos := lex.pos
	prevOffset := lex.offset

	token := lex.next()

	lex.pos.SetPosition(prevPos)
	lex.offset = prevOffset
	return token
}

func (lex *Lexer) Peek1() *token.Token {
	prevPos := lex.pos
	prevOffset := lex.offset

	var token *token.Token

	_ = lex.next()
	token = lex.next()

	lex.pos.SetPosition(prevPos)
	lex.offset = prevOffset

	return token
}

func (lex *Lexer) Skip() {
	lex.next()
}

func (lex *Lexer) NextIs(expectedKind token.Kind) bool {
	token := lex.Peek()
	return token.Kind == expectedKind
}

func (lex *Lexer) next() *token.Token {
	lex.skipWhitespace()
	character := lex.peekChar()
	if character == eof {
		return lex.consumeToken(nil, token.EOF)
	}
	token := lex.getToken(character)
	return token
}

// Useful for testing
func (lex *Lexer) Tokenize() ([]*token.Token, error) {
	var tokens []*token.Token
	for {
		tok := lex.next()
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

func (lex *Lexer) getToken(ch byte) *token.Token {
	tok := &token.Token{}
	tok.Lexeme = nil
	tok.Kind = token.INVALID

	switch ch {
	case '(':
		tok = lex.consumeToken(nil, token.OPEN_PAREN)
		lex.nextChar()
	case ')':
		tok = lex.consumeToken(nil, token.CLOSE_PAREN)
		lex.nextChar()
	case '{':
		tok = lex.consumeToken(nil, token.OPEN_CURLY)
		lex.nextChar()
	case '}':
		tok = lex.consumeToken(nil, token.CLOSE_CURLY)
		lex.nextChar()
	case '"':
		tok = lex.getStringLiteral()
	case ',':
		tok = lex.consumeToken(nil, token.COMMA)
		lex.nextChar()
	case ';':
		tok = lex.consumeToken(nil, token.SEMICOLON)
		lex.nextChar()
	case '+':
		tok = lex.consumeToken(nil, token.PLUS)
		lex.nextChar()
	case '-':
		tok = lex.consumeToken(nil, token.MINUS)
		lex.nextChar()
	case '*':
		tok = lex.consumeToken(nil, token.STAR)
		lex.nextChar()
	case '/':
		tok = lex.consumeToken(nil, token.SLASH)
		lex.nextChar()
	case '#':
		tok = lex.consumeToken(nil, token.SHARP)
		lex.nextChar()
	case '[':
		tok = lex.consumeToken(nil, token.OPEN_BRACKET)
		lex.nextChar()
	case ']':
		tok = lex.consumeToken(nil, token.CLOSE_BRACKET)
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
			lex.collector.ReportAndSave(invalidCharacter)
			return tok
		}
		if next == '=' {
			lex.nextChar()
			tok.Kind = token.BANG_EQUAL
			return tok
		}

		lex.collector.ReportAndSave(invalidCharacter)
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
			lex.collector.ReportAndSave(invalidCharacter)
			return tok
		}
		if next == '=' {
			lex.nextChar() // =
			tok.Kind = token.COLON_EQUAL
			return tok
		}

		lex.collector.ReportAndSave(invalidCharacter)
	default:
		position := lex.pos

		if unicode.IsLetter(rune(ch)) || ch == '_' {
			identifier := lex.readWhile(
				func(chr byte) bool { return unicode.IsNumber(rune(chr)) || unicode.IsLetter(rune(chr)) || chr == '_' },
			)
			tok = lex.classifyIdentifier(identifier, position)
		} else if ch >= '0' && ch <= '9' {
			tok = lex.getNumberLiteral(position)
		} else {
			tokenPosition := lex.pos
			invalidCharacter := diagnostics.Diag{
				Message: fmt.Sprintf("%s:%d:%d: invalid character %c", tokenPosition.Filename, tokenPosition.Line, tokenPosition.Column, ch),
			}
			lex.collector.ReportAndSave(invalidCharacter)
		}
	}
	return tok
}

func (lex *Lexer) getStringLiteral() *token.Token {
	tok := &token.Token{}
	tok.Kind = token.INVALID
	tok.Pos = lex.pos

	lex.nextChar() // "
	str := lex.readWhile(func(ch byte) bool { return ch != '"' })

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
		lex.collector.ReportAndSave(unterminatedStringLiteral)
		return tok
	}
	if ch == '"' {
		lex.nextChar()
	}

	tok.Kind = token.STRING_LITERAL
	tok.Lexeme = str

	return tok
}

func (lex *Lexer) getNumberLiteral(position token.Pos) *token.Token {
	number := lex.readWhile(
		func(chr byte) bool { return (chr >= '0' && chr <= '9') || chr == '_' },
	)

	// TODO: deal with floating pointer numbers
	// if strings.Contains(number, ".") {}

	token := token.New(number, token.INTEGER_LITERAL, position)
	return token
}

func (lex *Lexer) classifyIdentifier(identifier []byte, position token.Pos) *token.Token {
	idKind, ok := token.KEYWORDS[string(identifier)]
	if ok {
		return token.New(identifier, idKind, position)
	}
	return token.New(identifier, token.ID, position)
}

func (lex *Lexer) consumeToken(lexeme []byte, kind token.Kind) *token.Token {
	return token.New(lexeme, kind, lex.pos)
}

func (lex *Lexer) skipWhitespace() {
	lex.readWhile(func(ch byte) bool { return ch == ' ' || ch == '\t' || ch == '\n' || ch == '\r' })
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
