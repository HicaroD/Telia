package ast

import (
	"github.com/HicaroD/telia-lang/lexer/token/kind"
)

type ExprType interface {
	exprType()
}

type LiteralType struct {
	ExprType
	Kind kind.TokenKind
}
