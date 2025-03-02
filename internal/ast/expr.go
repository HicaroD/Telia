package ast

import (
	"fmt"

	"github.com/HicaroD/Telia/internal/lexer/token"
)

var LOGICAL map[token.Kind]bool = map[token.Kind]bool{
	token.AND: true,
	token.OR:  true,
}

var COMPARASION map[token.Kind]bool = map[token.Kind]bool{
	token.EQUAL_EQUAL: true,
	token.BANG_EQUAL:  true,
	token.GREATER:     true,
	token.GREATER_EQ:  true,
	token.LESS:        true,
	token.LESS_EQ:     true,
}

var TERM map[token.Kind]bool = map[token.Kind]bool{
	token.MINUS: true,
	token.PLUS:  true,
}

var FACTOR map[token.Kind]bool = map[token.Kind]bool{
	token.SLASH: true,
	token.STAR:  true,
}

var UNARY map[token.Kind]bool = map[token.Kind]bool{
	token.NOT:   true,
	token.MINUS: true,
}

type LiteralExpr struct {
	Type  *ExprType
	Value []byte
}

func (literal LiteralExpr) String() string {
	return fmt.Sprintf("%v %v\n", literal.Type, literal.Value)
}

type IdExpr struct {
	Name *token.Token
}

func (idExpr IdExpr) String() string {
	// Make it simpler to get lexeme
	return idExpr.Name.Name()
}

type NamespaceAccess struct {
	Left  *IdExpr
	Right *Node
}

func (namespaceAccess NamespaceAccess) String() string {
	return fmt.Sprintf("%v.%v\n", namespaceAccess.Left, namespaceAccess.Right)
}

type UnaryExpr struct {
	Op    token.Kind
	Value *Node
}

func (unary UnaryExpr) String() string {
	return fmt.Sprintf("%v %v\n", unary.Op, unary.Value)
}

type BinaryExpr struct {
	Left  *Node
	Op    token.Kind
	Right *Node
}

func (binExpr BinaryExpr) String() string {
	return fmt.Sprintf("(%v) %v (%v)\n", binExpr.Left, binExpr.Op, binExpr.Right)
}

type TupleExpr struct {
	Type  *TupleType
	Exprs []*Node
}

func (te TupleExpr) String() string {
	return fmt.Sprintf("%v\n", te.Exprs)
}
