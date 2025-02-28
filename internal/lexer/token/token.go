package token

import "fmt"

type Token struct {
	Lexeme []byte
	Kind   Kind
	Pos    Pos
}

func New(lexeme []byte, kind Kind, position Pos) *Token {
	return &Token{Lexeme: lexeme, Kind: kind, Pos: position}
}

func (token *Token) Name() string {
	if token.Kind == ID || token.Kind == UNTYPED_STRING {
		return string(token.Lexeme)
	}
	return token.Kind.String()
}

func (token *Token) String() string {
	return fmt.Sprintf("%s | %s | %s", string(token.Lexeme), token.Kind, token.Pos)
}
