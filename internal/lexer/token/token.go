package token

import (
	"fmt"
	"github.com/HicaroD/Telia/internal/source"
)

type Token struct {
	Literal bool // distinguishes literal spellings from type keywords
	Lexeme  []byte
	Kind    Kind
	Pos     Pos
	Span    source.Span
}

func New(lexeme []byte, kind Kind, position Pos) *Token {
	return &Token{Lexeme: lexeme, Kind: kind, Pos: position}
}

func (token *Token) Name() string {
	if token.Kind == ID || token.Kind == STRING_TYPE {
		return string(token.Lexeme)
	}
	return token.Kind.String()
}

func (token *Token) String() string {
	return fmt.Sprintf("%s | %s | %s", string(token.Lexeme), token.Kind, token.Pos)
}
