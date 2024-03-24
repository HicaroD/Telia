package lexer

import (
	"bufio"
	"io"
	"log"
	"strconv"
	"unicode"

	"github.com/HicaroD/telia-lang/lexer/token"
	"github.com/HicaroD/telia-lang/lexer/token/kind"
)

var KEYWORDS map[string]kind.TokenKind = map[string]kind.TokenKind{
	"fn":     kind.FN,
	"return": kind.RETURN,
}

type Lexer struct {
	filename string
	cursor   *Cursor
}

func NewLexer(filename string, reader *bufio.Reader) *Lexer {
	return &Lexer{filename: filename, cursor: newCursor(filename, reader)}
}

func (lex *Lexer) Tokenize() []token.Token {
	var tokens []token.Token
	for {
		lex.cursor.SkipWhitespace()
		character, err := lex.cursor.Peek()
		if err != nil {
			if err == io.EOF {
				break
			}
		}

		token := lex.getToken(character)
		tokens = append(tokens, token)
	}

	tokens = append(tokens, lex.consumeToken(nil, kind.EOF))
	return tokens
}

func (lex *Lexer) getToken(character rune) token.Token {
	switch character {
	case '(':
		token := lex.consumeToken(nil, kind.OPEN_PAREN)
		lex.cursor.Skip()
		return token
	case ')':
		token := lex.consumeToken(nil, kind.CLOSE_PAREN)
		lex.cursor.Skip()
		return token
	case '{':
		token := lex.consumeToken(nil, kind.OPEN_CURLY)
		lex.cursor.Skip()
		return token
	case '}':
		token := lex.consumeToken(nil, kind.CLOSE_CURLY)
		lex.cursor.Skip()
		return token
	case '"':
		return lex.getStringLiteral()
	case ';':
		token := lex.consumeToken(nil, kind.SEMICOLON)
		lex.cursor.Skip()
		return token
	default:
		position := lex.cursor.Position()
		if unicode.IsLetter(character) || character == '_' {
			// position := lex.cursor.Position()
			identifier := lex.cursor.ReadWhile(func(chr rune) bool { return unicode.IsNumber(chr) || unicode.IsLetter(chr) || chr == '_' })
			token := lex.classifyIdentifier(identifier, position)
			return token
		} else if unicode.IsNumber(character) {
			// position := lex.cursor.Position()
			number := lex.cursor.ReadWhile(func(chr rune) bool { return unicode.IsNumber(chr) || chr == '_' })

			// TODO: deal with floating pointer numbers
			// if strings.Contains(number, ".") {}

			convertedNumber, err := strconv.Atoi(number)
			if err != nil {
				// TODO(errors): unable to convert string to integer
				log.Fatalf("unable to convert string to integer due to error: '%s'", err)
			}
			token := token.NewToken(convertedNumber, kind.INTEGER_LITERAL, position)
			return token
		} else {
			// TODO(errors): invalid character
			log.Fatalf("invalid character: '%c' at %s", character, lex.cursor.Position())
		}
	}

	// TODO(errors): invalid character
	return lex.consumeToken(nil, kind.INVALID)
}

func (lex *Lexer) getStringLiteral() token.Token {
	position := lex.cursor.Position()

	lex.cursor.Skip() // "
	strLiteral := lex.cursor.ReadWhile(func(character rune) bool { return character != '"' })

	currentCharacter, err := lex.cursor.Peek()
	if err != nil {
		if err == io.EOF {
			// TODO(errors): deal with unterminated string literal error
			log.Fatal("unterminated string literal error")
		}
	}
	if currentCharacter == '"' {
		lex.cursor.Skip()
	}

	return token.NewToken(strLiteral, kind.STRING_LITERAL, position)
}

func (lex *Lexer) classifyIdentifier(identifier string, position token.Position) token.Token {
	idKind, ok := KEYWORDS[identifier]
	if ok {
		return token.NewToken(identifier, idKind, position)
	}
	return token.NewToken(identifier, kind.ID, position)
}

func (lex *Lexer) consumeToken(lexeme any, kind kind.TokenKind) token.Token {
	return token.NewToken(lexeme, kind, lex.cursor.Position())
}
