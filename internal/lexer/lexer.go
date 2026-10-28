package lexer

import (
	"fmt"
	"os"
	"unicode"

	"github.com/HicaroD/Telia/internal/ast"
	"github.com/HicaroD/Telia/internal/diagnostics"
	"github.com/HicaroD/Telia/internal/lexer/token"
	"github.com/HicaroD/Telia/internal/source"
)

const eof = '\000'

type Lexer struct {
	Loc           *ast.Loc
	Collector     *diagnostics.Collector
	StrictNewline bool
	BraceDepth    int

	src         []byte
	offset      int
	pos         token.Pos
	file        *source.File
	speculative int
}

func New(loc *ast.Loc, src []byte, collector *diagnostics.Collector) *Lexer {
	lexer := new(Lexer)

	lexer.Loc = loc
	lexer.Collector = collector
	lexer.StrictNewline = true
	lexer.pos = token.NewPosition(loc.Name, 1, 1)
	lexer.src = src
	path := loc.Path
	if path == "" {
		path = loc.Name
	}
	lexer.file = source.New(path, src)
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

// SpanFrom covers consumed source starting at a previously observed token.
func (lex *Lexer) SpanFrom(start source.Span) source.Span {
	return source.Span{File: lex.file, Start: start.Start, End: lex.offset}
}

func (lex *Lexer) Peek() *token.Token {
	lex.speculative++
	defer func() { lex.speculative-- }()
	prevPos := lex.pos
	prevOffset := lex.offset

	token := lex.Next()

	lex.pos.SetPosition(prevPos)
	lex.offset = prevOffset
	return token
}

func (lex *Lexer) Peek1() *token.Token {
	lex.speculative++
	defer func() { lex.speculative-- }()
	prevPos := lex.pos
	prevOffset := lex.offset

	var token *token.Token

	_ = lex.Next()
	token = lex.Next()

	lex.pos.SetPosition(prevPos)
	lex.offset = prevOffset

	return token
}

func (lex *Lexer) PeekN(n int) *token.Token {
	lex.speculative++
	defer func() { lex.speculative-- }()
	prevPos := lex.pos
	prevOffset := lex.offset

	var token *token.Token

	for n != 0 {
		lex.Next()
		n--
	}
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
	start := lex.offset
	character := lex.peekChar()

	tok := &token.Token{}
	defer func() {
		tok.Span = source.Span{File: lex.file, Start: start, End: lex.offset}
		lex.Collector.RegisterSpan(tok.Span)
		if lex.speculative == 0 {
			if tok.Kind == token.OPEN_CURLY {
				lex.BraceDepth++
			}
			if tok.Kind == token.CLOSE_CURLY {
				lex.BraceDepth--
			}
		}
	}()
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
	failed := false
	for {
		tok := lex.Next()
		if tok.Kind == token.INVALID {
			failed = true
			continue
		}
		tokens = append(tokens, tok)
		if tok.Kind == token.EOF {
			break
		}
	}
	if failed {
		return tokens, diagnostics.COMPILER_ERROR_FOUND
	}
	return tokens, nil
}

func (lex *Lexer) getToken(tok *token.Token, ch byte) *token.Token {
	tok.Pos = lex.pos

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
		tok.Literal = true
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
	case '&':
		lex.consumeTokenNoLex(tok, token.AMPERSAND)
		lex.nextChar()
	case '|':
		lex.consumeTokenNoLex(tok, token.PIPE)
		lex.nextChar()
	case '^':
		lex.consumeTokenNoLex(tok, token.CARET)
		lex.nextChar()
	case '~':
		lex.consumeTokenNoLex(tok, token.TILDE)
		lex.nextChar()
	case '!':
		lex.nextChar()

		next := lex.peekChar()
		if next == eof {
			lex.report(
				lex.offset-1,
				"Invalid character",
				"The character ! is not a valid operator; use != for inequality.",
			)
			return tok
		}
		if next == '=' {
			lex.nextChar()
			tok.Kind = token.BANG_EQUAL
			return tok
		}

		lex.report(
			lex.offset-1,
			"Invalid character",
			"The character ! is not a valid operator; use != for inequality.",
		)
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
		if next == '>' {
			lex.nextChar()
			tok.Kind = token.FAT_ARROW
			return tok
		}
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
		tok.Kind = token.COLON
		lex.nextChar() // :

		next := lex.peekChar()
		if next == '=' {
			lex.nextChar() // =
			tok.Kind = token.COLON_EQUAL
			break
		} else if next == ':' {
			lex.nextChar() // :
			tok.Kind = token.COLON_COLON
			break
		}
	default:
		if unicode.IsLetter(rune(ch)) || ch == '_' {
			lex.getIdOrKeyword(tok)
			word := string(tok.Lexeme)
			tok.Literal = word == "true" || word == "false" || word == "nil"
		} else if ch >= '0' && ch <= '9' {
			tok.Literal = true
			lex.getNumberLit(tok)
		} else {
			start := lex.offset
			lex.nextChar()
			lex.report(start, "Invalid character", fmt.Sprintf("The character %q is not valid here.", ch))
		}
	}
	return tok
}

func (lex *Lexer) getStringLit(tok *token.Token, isRaw bool) *token.Token {
	start := lex.offset
	invalid := false
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
			case '0':
				escape = 0
			default:
				lex.nextChar()
				lex.report(
					lex.offset-2,
					"Invalid escape sequence",
					fmt.Sprintf("Unknown escape \\%c.", escapeSym),
				)
				invalid = true
				continue
			}
			str = append(str, escape)
		} else {
			str = append(str, ch)
		}

		lex.nextChar()
	}

	ch := lex.peekChar()
	if ch != '"' {
		lex.report(start, "Unterminated string literal", "Expected a closing double quote.")
		return tok
	}
	if ch == '"' {
		lex.nextChar()
	}

	tok.Kind = token.STRING_TYPE
	if invalid {
		tok.Kind = token.INVALID
	}
	tok.Lexeme = str
	return tok
}

func (lex *Lexer) getNumberLit(tok *token.Token) {
	start := lex.offset
	var dotFound, dotRepeated bool
	numberType := token.INT_TYPE

	number := lex.readWhile(
		func(chr byte) bool {
			if chr == '.' {
				if dotFound {
					dotRepeated = true
					return false
				}
				numberType = token.FLOAT_TYPE
				dotFound = true
				return true
			}
			return (chr >= '0' && chr <= '9') || chr == '_'
		},
	)

	// TODO(errors)
	if dotRepeated {
		tok.Kind = token.INVALID
		lex.readWhile(func(ch byte) bool { return ch == '.' || ch >= '0' && ch <= '9' })
		lex.report(start, "Invalid number", "A number cannot contain more than one decimal point.")
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
	if lex.offset >= len(lex.src) {
		return eof
	}
	character := lex.src[lex.offset]
	lex.pos.Move(character)
	lex.offset++
	return character
}

func (lex *Lexer) report(start int, title, message string) {
	if lex.speculative != 0 {
		return
	}
	lex.Collector.ReportAndSave(
		diagnostics.Diag{
			Title:   title,
			Message: message,
			Span:    source.Span{File: lex.file, Start: start, End: lex.offset},
		},
	)
}

func (lex *Lexer) peekChar() byte {
	if lex.offset >= len(lex.src) {
		return eof
	}
	character := lex.src[lex.offset]
	return character
}
