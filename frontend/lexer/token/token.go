package token

import (
	"github.com/HicaroD/Telia/frontend/lexer/token/kind"
)

type Token struct {
	Lexeme   string
	Kind     kind.TokenKind
	Position Position
}

func New(lexeme string, kind kind.TokenKind, position Position) *Token {
	return &Token{Lexeme: lexeme, Kind: kind, Position: position}
}

func (token *Token) Name() string {
	if token.Kind == kind.ID {
		return token.Lexeme
	}
	return token.Kind.String()
}
