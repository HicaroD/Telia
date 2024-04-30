package lexer

import (
	"bufio"
	"fmt"
	"unicode"

	"github.com/HicaroD/Telia/collector"
	"github.com/HicaroD/Telia/lexer/token"
	"github.com/HicaroD/Telia/lexer/token/kind"
)

type lexer struct {
	filename  string
	collector *collector.DiagCollector
	cursor    *cursor
}

func New(filename string, reader *bufio.Reader, diagCollector *collector.DiagCollector) *lexer {
	return &lexer{filename: filename, cursor: new(filename, reader), collector: diagCollector}
}

func (lex *lexer) Tokenize() ([]*token.Token, error) {
	var tokens []*token.Token
	for {
		lex.cursor.skipWhitespace()
		character, ok := lex.cursor.peek()
		if !ok {
			break
		}
		token, err := lex.getToken(character)
		if err != nil {
			return nil, err
		}
		tokens = append(tokens, token)
	}
	tokens = append(tokens, lex.consumeToken("", kind.EOF))
	return tokens, nil
}

func (lex *lexer) getToken(character rune) (*token.Token, error) {
	switch character {
	case '(':
		token := lex.consumeToken("", kind.OPEN_PAREN)
		lex.cursor.skip()
		return token, nil
	case ')':
		token := lex.consumeToken("", kind.CLOSE_PAREN)
		lex.cursor.skip()
		return token, nil
	case '{':
		token := lex.consumeToken("", kind.OPEN_CURLY)
		lex.cursor.skip()
		return token, nil
	case '}':
		token := lex.consumeToken("", kind.CLOSE_CURLY)
		lex.cursor.skip()
		return token, nil
	case '"':
		stringLiteral, err := lex.getStringLiteral()
		if err != nil {
			return nil, err
		}
		return stringLiteral, nil
	case ',':
		token := lex.consumeToken("", kind.COMMA)
		lex.cursor.skip()
		return token, nil
	case ';':
		token := lex.consumeToken("", kind.SEMICOLON)
		lex.cursor.skip()
		return token, nil
	case '+':
		token := lex.consumeToken("", kind.PLUS)
		lex.cursor.skip()
		return token, nil
	case '-':
		token := lex.consumeToken("", kind.MINUS)
		lex.cursor.skip()
		return token, nil
	case '*':
		token := lex.consumeToken("", kind.STAR)
		lex.cursor.skip()
		return token, nil
	case '/':
		token := lex.consumeToken("", kind.SLASH)
		lex.cursor.skip()
		return token, nil
	case '!':
		tokenPosition := lex.cursor.Position()
		lex.cursor.skip()

		invalidCharacter := collector.Diag{
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
			return nil, collector.COMPILER_ERROR_FOUND
		}
		if next == '=' {
			lex.cursor.skip()
			return token.New("", kind.BANG_EQUAL, tokenPosition), nil
		}

		lex.collector.ReportAndSave(invalidCharacter)
		return nil, collector.COMPILER_ERROR_FOUND
	case '>':
		tokenKind := kind.GREATER
		tokenPosition := lex.cursor.Position()
		lex.cursor.skip() // >

		next, ok := lex.cursor.peek()
		if !ok || next != '=' {
			return token.New("", tokenKind, tokenPosition), nil
		}
		lex.cursor.skip() // =
		tokenKind = kind.GREATER_EQ
		return token.New("", tokenKind, tokenPosition), nil
	case '<':
		tokenKind := kind.LESS
		tokenPosition := lex.cursor.Position()
		lex.cursor.skip() // >

		next, ok := lex.cursor.peek()
		if !ok || next != '=' {
			return token.New("", tokenKind, tokenPosition), nil
		}
		lex.cursor.skip() // =
		tokenKind = kind.LESS_EQ
		return token.New("", tokenKind, tokenPosition), nil
	case '=':
		tokenKind := kind.EQUAL
		tokenPosition := lex.cursor.Position()
		lex.cursor.skip() // =

		next, ok := lex.cursor.peek()
		if !ok || next != '=' {
			return token.New("", tokenKind, tokenPosition), nil
		}
		lex.cursor.skip() // =
		tokenKind = kind.EQUAL_EQUAL
		return token.New("", tokenKind, tokenPosition), nil
	case '.':
		tokenKind := kind.DOT
		tokenPosition := lex.cursor.Position()
		lex.cursor.skip() // .

		next, ok := lex.cursor.peek()
		if !ok || next != '.' {
			return token.New("", tokenKind, tokenPosition), nil
		}
		lex.cursor.skip() // .
		tokenKind = kind.DOT_DOT

		next, ok = lex.cursor.peek()
		if !ok || next != '.' {
			return token.New("", tokenKind, tokenPosition), nil
		}
		lex.cursor.skip() // .
		tokenKind = kind.DOT_DOT_DOT
		return token.New("", tokenKind, tokenPosition), nil
	case ':':
		tokenPosition := lex.cursor.Position()
		lex.cursor.skip() // :

		invalidCharacter := collector.Diag{
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
			return nil, collector.COMPILER_ERROR_FOUND
		}
		if next == '=' {
			lex.cursor.skip() // =
			return token.New("", kind.COLON_EQUAL, tokenPosition), nil
		}

		lex.collector.ReportAndSave(invalidCharacter)
		return nil, collector.COMPILER_ERROR_FOUND
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
			invalidCharacter := collector.Diag{
				Message: fmt.Sprintf("%s:%d:%d: invalid character %c", tokenPosition.Filename, tokenPosition.Line, tokenPosition.Column, character),
			}
			lex.collector.ReportAndSave(invalidCharacter)
			return nil, collector.COMPILER_ERROR_FOUND
		}
	}
}

func (lex *lexer) getStringLiteral() (*token.Token, error) {
	position := lex.cursor.Position()

	lex.cursor.skip() // "
	strLiteral := lex.cursor.readWhile(func(character rune) bool { return character != '"' })

	currentCharacter, ok := lex.cursor.peek()
	if !ok {
		unterminatedStringLiteral := collector.Diag{
			Message: fmt.Sprintf(
				"%s:%d:%d: unterminated string literal",
				position.Filename,
				position.Line,
				position.Column,
			),
		}
		lex.collector.ReportAndSave(unterminatedStringLiteral)
		return nil, collector.COMPILER_ERROR_FOUND
	}
	if currentCharacter == '"' {
		lex.cursor.skip()
	}

	return token.New(strLiteral, kind.STRING_LITERAL, position), nil
}

func (lex *lexer) getNumberLiteral(position token.Position) *token.Token {
	number := lex.cursor.readWhile(
		func(chr rune) bool { return unicode.IsNumber(chr) || chr == '_' },
	)

	// TODO: deal with floating pointer numbers
	// if strings.Contains(number, ".") {}

	token := token.New(number, kind.INTEGER_LITERAL, position)
	return token
}

func (lex *lexer) classifyIdentifier(identifier string, position token.Position) *token.Token {
	idKind, ok := kind.KEYWORDS[identifier]
	if ok {
		return token.New(identifier, idKind, position)
	}
	return token.New(identifier, kind.ID, position)
}

func (lex *lexer) consumeToken(lexeme string, kind kind.TokenKind) *token.Token {
	return token.New(lexeme, kind, lex.cursor.Position())
}
