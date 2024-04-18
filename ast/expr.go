package ast

import (
	"fmt"

	"github.com/HicaroD/Telia/lexer/token"
	"github.com/HicaroD/Telia/lexer/token/kind"
)

var LOGICAL map[kind.TokenKind]bool = map[kind.TokenKind]bool{
	kind.AND: true,
	kind.OR:  true,
}

var COMPARASION map[kind.TokenKind]bool = map[kind.TokenKind]bool{
	kind.EQUAL_EQUAL: true,
	kind.BANG_EQUAL:  true,
	kind.GREATER:     true,
	kind.GREATER_EQ:  true,
	kind.LESS:        true,
	kind.LESS_EQ:     true,
}

var TERM map[kind.TokenKind]bool = map[kind.TokenKind]bool{
	kind.MINUS: true,
	kind.PLUS:  true,
}

var FACTOR map[kind.TokenKind]bool = map[kind.TokenKind]bool{
	kind.SLASH: true,
	kind.STAR:  true,
}

var UNARY map[kind.TokenKind]bool = map[kind.TokenKind]bool{
	kind.NOT:   true,
	kind.MINUS: true,
}

type Expr interface {
	AstNode
	IsVoid() bool
	exprNode()
}

// Used on empty return
type VoidExpr struct {
	Expr
}

func (void VoidExpr) String() string {
	return "void"
}
func (void VoidExpr) IsVoid() bool { return true }
func (void VoidExpr) exprNode()    {}

type LiteralExpr struct {
	Expr
	Kind  kind.TokenKind
	Value any
}

func (literal LiteralExpr) String() string {
	return fmt.Sprintf("%s %s", literal.Kind, literal.Value)
}
func (literal LiteralExpr) IsVoid() bool { return false }
func (literal LiteralExpr) exprNode()    {}

type IdExpr struct {
	Expr
	Name *token.Token
}

func (idExpr IdExpr) String() string {
	return fmt.Sprintf("%s", idExpr.Name)
}
func (idExpr IdExpr) IsVoid() bool { return false }
func (idExpr IdExpr) exprNode()    {}

type UnaryExpr struct {
	Expr
	Op    kind.TokenKind
	Value Expr
}

func (unary UnaryExpr) String() string {
	return fmt.Sprintf("%s %s", unary.Op, unary.Value)
}
func (unary UnaryExpr) IsVoid() bool { return false }
func (unary UnaryExpr) exprNode()    {}

type BinaryExpr struct {
	Expr
	Left  Expr
	Op    kind.TokenKind
	Right Expr
}

func (binExpr BinaryExpr) String() string {
	return fmt.Sprintf("(%s) %s (%s)", binExpr.Left, binExpr.Op, binExpr.Right)
}
func (binExpr BinaryExpr) IsVoid() bool { return false }
func (binExpr BinaryExpr) exprNode()    {}
