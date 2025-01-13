package lexer

import (
	"bufio"
	"fmt"
	"unicode"

	"github.com/HicaroD/Telia/diagnostics"
	"github.com/HicaroD/Telia/frontend/lexer/token"
)

type Lexer struct {
	filename  string
	collector *diagnostics.Collector
	cursor    *cursor

	curr *token.Token
	next *token.Token
}

func New(filename string, reader *bufio.Reader, diagCollector *diagnostics.Collector) *Lexer {
	return &Lexer{filename: filename, cursor: new(filename, reader), collector: diagCollector}
}

func (lex *Lexer) Next() (*token.Token, error) {
	lex.cursor.skipWhitespace()
	character, ok := lex.cursor.peek()
	if !ok {
		return lex.consumeToken("", token.EOF), nil
	}
	token, err := lex.getToken(character)
	if err != nil {
		return nil, err
	}
	return token, nil
}

// TODO: don't return error
func (lex *Lexer) Peek() (*token.Token, error) {
	currentPos := lex.cursor.Position()
	token, err := lex.Next()
	if err != nil {
		return nil, err
	}
	lex.cursor.SetPosition(currentPos)
	return token, nil
}

// TODO: don't return error
func (lex *Lexer) Peek1() (*token.Token, error) {
	currentPos := lex.cursor.Position()
	var token *token.Token

	_, err := lex.Next()
	if err != nil {
		return nil, err
	}

	token, err = lex.Next()
	if err != nil {
		return nil, err
	}

	lex.cursor.SetPosition(currentPos)
	return token, nil
}

// TODO: do I need to check the error here?
func (lex *Lexer) Skip() {
	lex.Next()
}

func (lex *Lexer) NextIs(expectedKind token.Kind) bool {
	token, err := lex.Peek()
	// TODO: HANDLE THIS PROPERLY
	if err != nil {
		return false
	}
	return token.Kind == expectedKind
}

func (lex *Lexer) Tokenize() ([]*token.Token, error) {
	var tokens []*token.Token
	for {
		tok, err := lex.Next()
		if err != nil {
			return nil, err
		}
		tokens = append(tokens, tok)
		if tok.Kind == token.EOF {
			break
		}
	}
	return tokens, nil
}

func (lex *Lexer) getToken(character rune) (*token.Token, error) {
	switch character {
	case '(':
		token := lex.consumeToken("", token.OPEN_PAREN)
		lex.cursor.skip()
		return token, nil
	case ')':
		token := lex.consumeToken("", token.CLOSE_PAREN)
		lex.cursor.skip()
		return token, nil
	case '{':
		token := lex.consumeToken("", token.OPEN_CURLY)
		lex.cursor.skip()
		return token, nil
	case '}':
		token := lex.consumeToken("", token.CLOSE_CURLY)
		lex.cursor.skip()
		return token, nil
	case '"':
		stringLiteral, err := lex.getStringLiteral()
		if err != nil {
			return nil, err
		}
		return stringLiteral, nil
	case ',':
		token := lex.consumeToken("", token.COMMA)
		lex.cursor.skip()
		return token, nil
	case ';':
		token := lex.consumeToken("", token.SEMICOLON)
		lex.cursor.skip()
		return token, nil
	case '+':
		token := lex.consumeToken("", token.PLUS)
		lex.cursor.skip()
		return token, nil
	case '-':
		token := lex.consumeToken("", token.MINUS)
		lex.cursor.skip()
		return token, nil
	case '*':
		token := lex.consumeToken("", token.STAR)
		lex.cursor.skip()
		return token, nil
	case '/':
		token := lex.consumeToken("", token.SLASH)
		lex.cursor.skip()
		return token, nil
	case '!':
		tokenPosition := lex.cursor.Position()
		lex.cursor.skip()

		invalidCharacter := diagnostics.Diag{
			Message: fmt.Sprintf(
				"%s:%d:%d: invalid character !",
				tokenPosition.Filename,
				tokenPosition.Line,
				tokenPosition.Column,
			),
		}

		next, ok := lex.cursor.peek()
		if !ok {
			lex.collector.ReportAndSave(invalidCharacter)
			return nil, diagnostics.COMPILER_ERROR_FOUND
		}
		if next == '=' {
			lex.cursor.skip()
			return token.New("", token.BANG_EQUAL, tokenPosition), nil
		}

		lex.collector.ReportAndSave(invalidCharacter)
		return nil, diagnostics.COMPILER_ERROR_FOUND
	case '>':
		tokenKind := token.GREATER
		tokenPosition := lex.cursor.Position()
		lex.cursor.skip() // >

		next, ok := lex.cursor.peek()
		if !ok || next != '=' {
			return token.New("", tokenKind, tokenPosition), nil
		}
		lex.cursor.skip() // =
		tokenKind = token.GREATER_EQ
		return token.New("", tokenKind, tokenPosition), nil
	case '<':
		tokenKind := token.LESS
		tokenPosition := lex.cursor.Position()
		lex.cursor.skip() // >

		next, ok := lex.cursor.peek()
		if !ok || next != '=' {
			return token.New("", tokenKind, tokenPosition), nil
		}
		lex.cursor.skip() // =
		tokenKind = token.LESS_EQ
		return token.New("", tokenKind, tokenPosition), nil
	case '=':
		tokenKind := token.EQUAL
		tokenPosition := lex.cursor.Position()
		lex.cursor.skip() // =

		next, ok := lex.cursor.peek()
		if !ok || next != '=' {
			return token.New("", tokenKind, tokenPosition), nil
		}
		lex.cursor.skip() // =
		tokenKind = token.EQUAL_EQUAL
		return token.New("", tokenKind, tokenPosition), nil
	case '.':
		tokenKind := token.DOT
		tokenPosition := lex.cursor.Position()
		lex.cursor.skip() // .

		next, ok := lex.cursor.peek()
		if !ok || next != '.' {
			return token.New("", tokenKind, tokenPosition), nil
		}
		lex.cursor.skip() // .
		tokenKind = token.DOT_DOT

		next, ok = lex.cursor.peek()
		if !ok || next != '.' {
			return token.New("", tokenKind, tokenPosition), nil
		}
		lex.cursor.skip() // .
		tokenKind = token.DOT_DOT_DOT
		return token.New("", tokenKind, tokenPosition), nil
	case ':':
		tokenPosition := lex.cursor.Position()
		lex.cursor.skip() // :

		invalidCharacter := diagnostics.Diag{
			Message: fmt.Sprintf(
				"%s:%d:%d: invalid character :",
				tokenPosition.Filename,
				tokenPosition.Line,
				tokenPosition.Column,
			),
		}

		next, ok := lex.cursor.peek()
		if !ok {
			lex.collector.ReportAndSave(invalidCharacter)
			return nil, diagnostics.COMPILER_ERROR_FOUND
		}
		if next == '=' {
			lex.cursor.skip() // =
			return token.New("", token.COLON_EQUAL, tokenPosition), nil
		}

		lex.collector.ReportAndSave(invalidCharacter)
		return nil, diagnostics.COMPILER_ERROR_FOUND
	default:
		position := lex.cursor.Position()
		if unicode.IsLetter(character) || character == '_' {
			identifier := lex.cursor.readWhile(
				func(chr rune) bool { return unicode.IsNumber(chr) || unicode.IsLetter(chr) || chr == '_' },
			)
			token := lex.classifyIdentifier(identifier, position)
			return token, nil
		} else if unicode.IsNumber(character) {
			return lex.getNumberLiteral(position), nil
		} else {
			tokenPosition := lex.cursor.Position()
			invalidCharacter := diagnostics.Diag{
				Message: fmt.Sprintf("%s:%d:%d: invalid character %c", tokenPosition.Filename, tokenPosition.Line, tokenPosition.Column, character),
			}
			lex.collector.ReportAndSave(invalidCharacter)
			return nil, diagnostics.COMPILER_ERROR_FOUND
		}
	}
}

func (lex *Lexer) getStringLiteral() (*token.Token, error) {
	position := lex.cursor.Position()

	lex.cursor.skip() // "
	strLiteral := lex.cursor.readWhile(func(character rune) bool { return character != '"' })

	currentCharacter, ok := lex.cursor.peek()
	if !ok {
		unterminatedStringLiteral := diagnostics.Diag{
			Message: fmt.Sprintf(
				"%s:%d:%d: unterminated string literal",
				position.Filename,
				position.Line,
				position.Column,
			),
		}
		lex.collector.ReportAndSave(unterminatedStringLiteral)
		return nil, diagnostics.COMPILER_ERROR_FOUND
	}
	if currentCharacter == '"' {
		lex.cursor.skip()
	}

	return token.New(strLiteral, token.STRING_LITERAL, position), nil
}

func (lex *Lexer) getNumberLiteral(position token.Position) *token.Token {
	number := lex.cursor.readWhile(
		func(chr rune) bool { return unicode.IsNumber(chr) || chr == '_' },
	)

	// TODO: deal with floating pointer numbers
	// if strings.Contains(number, ".") {}

	token := token.New(number, token.INTEGER_LITERAL, position)
	return token
}

func (lex *Lexer) classifyIdentifier(identifier string, position token.Position) *token.Token {
	idKind, ok := token.KEYWORDS[identifier]
	if ok {
		return token.New(identifier, idKind, position)
	}
	return token.New(identifier, token.ID, position)
}

func (lex *Lexer) consumeToken(lexeme string, kind token.Kind) *token.Token {
	return token.New(lexeme, kind, lex.cursor.Position())
}
