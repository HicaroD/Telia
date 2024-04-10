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
	return fmt.Sprintf("%s %s", literal.Kind, literal.Value)
}
func (literal LiteralExpr) exprNode() {}

type BinaryExpr struct {
	Expr
	Left  Expr
	Op    kind.TokenKind
	Right Expr
}

func (binExpr BinaryExpr) String() string {
	return fmt.Sprintf("%s %s %s", binExpr.Left, binExpr.Op, binExpr.Right)
}
func (binExpr BinaryExpr) exprNode() {}

type IdExpr struct {
	Expr
	Name *token.Token
}

func (idExpr IdExpr) String() string {
	return fmt.Sprintf("%s", idExpr.Name)
}
func (idExpr IdExpr) exprNode() {}
