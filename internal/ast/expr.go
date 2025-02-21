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

type Expr interface {
	Node
	IsId() bool
	IsVoid() bool
	IsFieldAccess() bool
	exprNode()
}

// Used on empty return
type VoidExpr struct {
	Expr
}

func (void VoidExpr) String() string {
	return "void"
}
func (void VoidExpr) IsId() bool          { return false }
func (void VoidExpr) IsVoid() bool        { return true }
func (void VoidExpr) IsFieldAccess() bool { return false }
func (void VoidExpr) exprNode()           {}

type LiteralExpr struct {
	Expr
	Type  ExprType
	Value []byte
}

func (literal LiteralExpr) String() string {
	return fmt.Sprintf("%s %s", literal.Type, literal.Value)
}
func (literal LiteralExpr) IsId() bool          { return false }
func (literal LiteralExpr) IsVoid() bool        { return false }
func (literal LiteralExpr) IsFieldAccess() bool { return false }
func (literal LiteralExpr) exprNode()           {}

type IdExpr struct {
	Expr
	Name *token.Token
}

func (idExpr IdExpr) String() string {
	// Make it simpler to get lexeme
	return idExpr.Name.Name()
}
func (idExpr IdExpr) IsId() bool          { return true }
func (idExpr IdExpr) IsVoid() bool        { return false }
func (idExpr IdExpr) IsFieldAccess() bool { return false }
func (idExpr IdExpr) exprNode()           {}

type FieldAccess struct {
	Stmt
	Expr
	Left  *MyNode
	Right *MyNode
}

func (fieldAccess FieldAccess) String() string {
	return fmt.Sprintf("%s.%s", fieldAccess.Left, fieldAccess.Right)
}
func (fieldAccess FieldAccess) IsId() bool          { return false }
func (fieldAccess FieldAccess) IsReturn() bool      { return false }
func (fieldAccess FieldAccess) IsFieldAccess() bool { return true }
func (fieldAccess FieldAccess) astNode()            {}
func (fieldAccess FieldAccess) stmtNode()           {}
func (fieldAccess FieldAccess) exprNode()           {}

type UnaryExpr struct {
	Expr
	Op    token.Kind
	Value *MyNode
}

func (unary UnaryExpr) String() string {
	return fmt.Sprintf("%s %s", unary.Op, unary.Value)
}
func (unary UnaryExpr) IsId() bool          { return false }
func (unary UnaryExpr) IsVoid() bool        { return false }
func (unary UnaryExpr) IsFieldAccess() bool { return false }
func (unary UnaryExpr) exprNode()           {}

type BinaryExpr struct {
	Expr
	Left  *MyNode
	Op    token.Kind
	Right *MyNode
}

func (binExpr BinaryExpr) String() string {
	return fmt.Sprintf("(%s) %s (%s)", binExpr.Left, binExpr.Op, binExpr.Right)
}
func (binExpr BinaryExpr) IsId() bool          { return false }
func (binExpr BinaryExpr) IsVoid() bool        { return false }
func (binExpr BinaryExpr) IsFieldAccess() bool { return false }
func (binExpr BinaryExpr) exprNode()           {}
