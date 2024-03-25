package ast

import (
	"github.com/HicaroD/telia-lang/lexer/token/kind"
)

type Expr interface {
	AstNode
	exprNode()
}

type LiteralExpr struct {
	Expr
	Kind  kind.TokenKind
	Value any
}

type BinaryExpr struct {
	Expr
	Left  Expr
	Op    kind.TokenKind
	Right Expr
}
