package token

import (
	"github.com/HicaroD/telia-lang/lexer/token/kind"
)

// TODO: define lexeme, position and more related to token
type Token struct {
	Kind kind.TokenKind
}

func NewToken(kind kind.TokenKind) Token {
	return Token{Kind: kind}
}
