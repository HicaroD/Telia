package token

import (
	"github.com/HicaroD/Telia/lexer/token/kind"
)

// TODO: define lexeme, position and more related to token
type Token struct {
	Lexeme   any
	Kind     kind.TokenKind
	Position Position
}

func New(lexeme any, kind kind.TokenKind, position Position) *Token {
	return &Token{Lexeme: lexeme, Kind: kind, Position: position}
}
