package token

import (
	"github.com/HicaroD/telia-lang/lexer/token/kind"
)

// TODO: define lexeme, position and more related to token
type Token struct {
	kind kind.TokenKind
}
