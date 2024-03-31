package lexer

import (
	"bufio"
	"log"
	"strconv"
	"unicode"

	"github.com/HicaroD/telia-lang/lexer/token"
	"github.com/HicaroD/telia-lang/lexer/token/kind"
)

type lexer struct {
	filename string
	cursor   *cursor
}

func NewLexer(filename string, reader *bufio.Reader) *lexer {
	return &lexer{filename: filename, cursor: newCursor(filename, reader)}
}

func (lex *lexer) Tokenize() []*token.Token {
	var tokens []*token.Token
	for {
		lex.cursor.SkipWhitespace()
		character, ok := lex.cursor.Peek()
		if !ok {
			break
		}
		token := lex.getToken(character)
		tokens = append(tokens, token)
	}
	tokens = append(tokens, lex.consumeToken(nil, kind.EOF))
	return tokens
}

func (lex *lexer) getToken(character rune) *token.Token {
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
	case ',':
		token := lex.consumeToken(nil, kind.COMMA)
		lex.cursor.Skip()
		return token
	case ';':
		token := lex.consumeToken(nil, kind.SEMICOLON)
		lex.cursor.Skip()
		return token
	case '*':
		token := lex.consumeToken(nil, kind.STAR)
		lex.cursor.Skip()
		return token
	case '=':
		tokenKind := kind.EQUAL
		tokenPosition := lex.cursor.Position()
		lex.cursor.Skip() // =

		next, ok := lex.cursor.Peek()
		if !ok || next != '=' {
			return token.New(nil, tokenKind, tokenPosition)
		}
		lex.cursor.Skip() // =
		tokenKind = kind.EQUAL_EQUAL
		return token.New(nil, tokenKind, tokenPosition)
	case '.':
		tokenKind := kind.DOT
		tokenPosition := lex.cursor.Position()
		lex.cursor.Skip() // .

		next, ok := lex.cursor.Peek()
		if !ok || next != '.' {
			return token.New(nil, tokenKind, tokenPosition)
		}
		lex.cursor.Skip() // .
		tokenKind = kind.DOT_DOT

		next, ok = lex.cursor.Peek()
		if !ok || next != '.' {
			return token.New(nil, tokenKind, tokenPosition)
		}
		lex.cursor.Skip() // .
		tokenKind = kind.DOT_DOT_DOT
		return token.New(nil, tokenKind, tokenPosition)
	case '-':
		token := lex.consumeToken(nil, kind.MINUS)
		lex.cursor.Skip()
		return token
	case ':':
		tokenPosition := lex.cursor.Position()
		lex.cursor.Skip() // :

		next, ok := lex.cursor.Peek()
		if !ok {
			// TODO(errors)
			log.Fatal("invalid ':' character") // EOF
		}
		if next == '=' {
			lex.cursor.Skip() // =
			return token.New(nil, kind.COLON_EQUAL, tokenPosition)
		}
		// TODO(errors)
		log.Fatal("invalid ':' character")
	default:
		position := lex.cursor.Position()
		if unicode.IsLetter(character) || character == '_' {
			identifier := lex.cursor.ReadWhile(func(chr rune) bool { return unicode.IsNumber(chr) || unicode.IsLetter(chr) || chr == '_' })
			token := lex.classifyIdentifier(identifier, position)
			return token
		} else if unicode.IsNumber(character) {
			number := lex.cursor.ReadWhile(func(chr rune) bool { return unicode.IsNumber(chr) || chr == '_' })

			// TODO: deal with floating pointer numbers
			// if strings.Contains(number, ".") {}

			convertedNumber, err := strconv.Atoi(number)
			if err != nil {
				// TODO(errors): unable to convert string to integer
				log.Fatalf("unable to convert string to integer due to error: '%s'", err)
			}
			token := token.New(convertedNumber, kind.INTEGER_LITERAL, position)
			return token
		} else {
			// TODO(errors): invalid character
			log.Fatalf("invalid character: '%c' at %s", character, lex.cursor.Position())
		}
	}

	// TODO(errors): invalid character
	return lex.consumeToken(nil, kind.INVALID)
}

func (lex *lexer) getStringLiteral() *token.Token {
	position := lex.cursor.Position()

	lex.cursor.Skip() // "
	strLiteral := lex.cursor.ReadWhile(func(character rune) bool { return character != '"' })

	currentCharacter, ok := lex.cursor.Peek()
	if !ok {
		// TODO(errors): deal with unterminated string literal error
		log.Fatal("unterminated string literal error")
	}
	if currentCharacter == '"' {
		lex.cursor.Skip()
	}

	return token.New(strLiteral, kind.STRING_LITERAL, position)
}

func (lex *lexer) classifyIdentifier(identifier string, position token.Position) *token.Token {
	idKind, ok := kind.KEYWORDS[identifier]
	if ok {
		return token.New(identifier, idKind, position)
	}
	return token.New(identifier, kind.ID, position)
}

func (lex *lexer) consumeToken(lexeme any, kind kind.TokenKind) *token.Token {
	return token.New(lexeme, kind, lex.cursor.Position())
}
