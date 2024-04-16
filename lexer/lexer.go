package lexer

import (
	"bufio"
	"log"
	"unicode"

	"github.com/HicaroD/Telia/lexer/token"
	"github.com/HicaroD/Telia/lexer/token/kind"
)

type lexer struct {
	filename string
	cursor   *cursor
}

func New(filename string, reader *bufio.Reader) *lexer {
	return &lexer{filename: filename, cursor: new(filename, reader)}
}

func (lex *lexer) Tokenize() ([]*token.Token, error) {
	var tokens []*token.Token
	for {
		lex.cursor.skipWhitespace()
		character, ok := lex.cursor.peek()
		if !ok {
			break
		}
		token := lex.getToken(character)
		tokens = append(tokens, token)
	}
	tokens = append(tokens, lex.consumeToken(nil, kind.EOF))
	return tokens, nil
}

func (lex *lexer) getToken(character rune) *token.Token {
	switch character {
	case '(':
		token := lex.consumeToken(nil, kind.OPEN_PAREN)
		lex.cursor.skip()
		return token
	case ')':
		token := lex.consumeToken(nil, kind.CLOSE_PAREN)
		lex.cursor.skip()
		return token
	case '{':
		token := lex.consumeToken(nil, kind.OPEN_CURLY)
		lex.cursor.skip()
		return token
	case '}':
		token := lex.consumeToken(nil, kind.CLOSE_CURLY)
		lex.cursor.skip()
		return token
	case '"':
		return lex.getStringLiteral()
	case ',':
		token := lex.consumeToken(nil, kind.COMMA)
		lex.cursor.skip()
		return token
	case ';':
		token := lex.consumeToken(nil, kind.SEMICOLON)
		lex.cursor.skip()
		return token
	case '+':
		token := lex.consumeToken(nil, kind.PLUS)
		lex.cursor.skip()
		return token
	case '-':
		token := lex.consumeToken(nil, kind.MINUS)
		lex.cursor.skip()
		return token
	case '*':
		token := lex.consumeToken(nil, kind.STAR)
		lex.cursor.skip()
		return token
	case '/':
		token := lex.consumeToken(nil, kind.SLASH)
		lex.cursor.skip()
		return token
	case '!':
		tokenPosition := lex.cursor.Position()
		lex.cursor.skip()

		next, ok := lex.cursor.peek()
		// TODO(errors)
		if !ok {
			log.Fatal("invalid '!' token")
		}
		if next == '=' {
			lex.cursor.skip()
			return token.New(nil, kind.BANG_EQUAL, tokenPosition)
		}
		// TODO(errors)
		log.Fatal("invalid '!' token")
	case '>':
		tokenKind := kind.GREATER
		tokenPosition := lex.cursor.Position()
		lex.cursor.skip() // >

		next, ok := lex.cursor.peek()
		if !ok || next != '=' {
			return token.New(nil, tokenKind, tokenPosition)
		}
		lex.cursor.skip() // =
		tokenKind = kind.GREATER_EQ
		return token.New(nil, tokenKind, tokenPosition)
	case '<':
		tokenKind := kind.LESS
		tokenPosition := lex.cursor.Position()
		lex.cursor.skip() // >

		next, ok := lex.cursor.peek()
		if !ok || next != '=' {
			return token.New(nil, tokenKind, tokenPosition)
		}
		lex.cursor.skip() // =
		tokenKind = kind.LESS_EQ
		return token.New(nil, tokenKind, tokenPosition)
	case '=':
		tokenKind := kind.EQUAL
		tokenPosition := lex.cursor.Position()
		lex.cursor.skip() // =

		next, ok := lex.cursor.peek()
		if !ok || next != '=' {
			return token.New(nil, tokenKind, tokenPosition)
		}
		lex.cursor.skip() // =
		tokenKind = kind.EQUAL_EQUAL
		return token.New(nil, tokenKind, tokenPosition)
	case '.':
		tokenKind := kind.DOT
		tokenPosition := lex.cursor.Position()
		lex.cursor.skip() // .

		next, ok := lex.cursor.peek()
		if !ok || next != '.' {
			return token.New(nil, tokenKind, tokenPosition)
		}
		lex.cursor.skip() // .
		tokenKind = kind.DOT_DOT

		next, ok = lex.cursor.peek()
		if !ok || next != '.' {
			return token.New(nil, tokenKind, tokenPosition)
		}
		lex.cursor.skip() // .
		tokenKind = kind.DOT_DOT_DOT
		return token.New(nil, tokenKind, tokenPosition)
	case ':':
		tokenPosition := lex.cursor.Position()
		lex.cursor.skip() // :
		next, ok := lex.cursor.peek()
		if !ok {
			// TODO(errors)
			log.Fatal("invalid ':' character") // EOF
		}
		if next == '=' {
			lex.cursor.skip() // =
			return token.New(nil, kind.COLON_EQUAL, tokenPosition)
		}
		// TODO(errors)
		log.Fatal("invalid ':' character")
	default:
		position := lex.cursor.Position()
		if unicode.IsLetter(character) || character == '_' {
			identifier := lex.cursor.readWhile(func(chr rune) bool { return unicode.IsNumber(chr) || unicode.IsLetter(chr) || chr == '_' })
			token := lex.classifyIdentifier(identifier, position)
			return token
		} else if unicode.IsNumber(character) {
			return lex.getNumberLiteral(position)
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

	lex.cursor.skip() // "
	strLiteral := lex.cursor.readWhile(func(character rune) bool { return character != '"' })

	currentCharacter, ok := lex.cursor.peek()
	if !ok {
		// TODO(errors): deal with unterminated string literal error
		log.Fatal("unterminated string literal error")
	}
	if currentCharacter == '"' {
		lex.cursor.skip()
	}

	return token.New(strLiteral, kind.STRING_LITERAL, position)
}

func (lex *lexer) getNumberLiteral(position token.Position) *token.Token {
	number := lex.cursor.readWhile(func(chr rune) bool { return unicode.IsNumber(chr) || chr == '_' })

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

func (lex *lexer) consumeToken(lexeme any, kind kind.TokenKind) *token.Token {
	return token.New(lexeme, kind, lex.cursor.Position())
}
