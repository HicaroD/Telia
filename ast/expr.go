package ast

import (
	"fmt"

	"github.com/HicaroD/telia-lang/lexer/token"
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

func (literal LiteralExpr) String() string {
	return fmt.Sprintf("Literal: %s", literal.Kind)
}
func (literal LiteralExpr) exprNode() {}

type BinaryExpr struct {
	Expr
	Left  Expr
	Op    kind.TokenKind
	Right Expr
}

func (binExpr BinaryExpr) String() string {
	return fmt.Sprintf("BinOp: %s", binExpr.Op)
}
func (binExpr BinaryExpr) exprNode() {}

type IdExpr struct {
	Expr
	Name *token.Token
}

func (idExpr IdExpr) String() string {
	return fmt.Sprintf("IdExpr: %s", idExpr.Name)
}
func (idExpr IdExpr) exprNode() {}
