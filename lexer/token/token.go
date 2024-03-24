package token

import (
	"github.com/HicaroD/telia-lang/lexer/token/kind"
)

// TODO: define lexeme, position and more related to token
type Token struct {
	Lexeme   any
	Kind     kind.TokenKind
	Position Position
}

func NewToken(lexeme any, kind kind.TokenKind, position Position) Token {
	return Token{Lexeme: lexeme, Kind: kind, Position: position}
}
