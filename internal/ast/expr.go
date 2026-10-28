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
	token.TILDE: true,
	token.NOT:   true,
	token.MINUS: true,
}

type LiteralExpr struct {
	Type  *ExprType
	Value []byte
}

func (literal *LiteralExpr) String() string {
	return fmt.Sprintf("%v %v\n", literal.Type, literal.Value)
}

type IdExpr struct {
	Name     *token.Token
	Constant *LiteralExpr

	// Codegen
	N *Node
}

func (idExpr IdExpr) String() string {
	// Make it simpler to get lexeme
	return idExpr.Name.Name()
}

type NamespaceAccess struct {
	EnumVariant *EnumVariant
	EnumDecl    *EnumDecl
	TypeArgs    []*ExprType
	IsImport    bool
	Left        *IdExpr
	Right       *Node
}

func (n *NamespaceAccess) String() string {
	return fmt.Sprintf("%v::%v\n", n.Left, n.Right)
}

type FieldAccess struct {
	Left  *IdExpr
	Right *Node

	// Codegen
	AccessedField *StructField
	Param         bool
	StructVar     *VarIdStmt
	StructParam   *Param
	Decl          *StructDecl
	BackendType   any
}

func (f *FieldAccess) String() string {
	return fmt.Sprintf("%v.%v\n", f.Left, f.Right)
}

type UnaryExpr struct {
	Op    token.Kind
	Value *Node
	Type  *ExprType
}

func (unary *UnaryExpr) String() string {
	return fmt.Sprintf("%v %v\n", unary.Op, unary.Value)
}

type BinExpr struct {
	Left        *Node
	Op          token.Kind
	Right       *Node
	OperandType *ExprType
	RightType   *ExprType
}

func (binExpr *BinExpr) String() string {
	return fmt.Sprintf("(%v) %v (%v)\n", binExpr.Left, binExpr.Op, binExpr.Right)
}

type TupleExpr struct {
	Type  *TupleType
	Exprs []*Node
}

func (te *TupleExpr) String() string {
	return fmt.Sprintf("%v\n", te.Exprs)
}

type VarArgsExpr struct {
	Args []*Node
}

func (v *VarArgsExpr) String() string {
	return fmt.Sprintf("VARGS: %v\n", v.Args)
}

type StructLiteralExpr struct {
	TypeArgs  []*ExprType
	Namespace *token.Token
	Decl      *StructDecl
	Name      *token.Token
	Values    []*StructFieldValue
}

type StructFieldValue struct {
	Name  *token.Token
	Index int
	Value *Node
}

type StructFieldAccessExpr struct {
	Access *Node
}

type AddressOfExpr struct {
	Expr *Node
}

type DerefPointerExpr struct {
	Expr *Node
	Type *ExprType
}

type NullPtrExpr struct {
	Type *ExprType
}

type ArrayLiteralExpr struct {
	Values []*Node
	Type   *ExprType
}

type IndexExpr struct {
	Value *Node
	Index *Node
	Type  *ExprType // type of the indexed sequence
}

func (n *Node) IsAddressable() bool {
	switch n.Kind {
	case KIND_ID_EXPR:
		id := n.Node.(*IdExpr)
		return id.N == nil || id.N.Kind != KIND_CONST_DECL
	case KIND_FIELD_ACCESS, KIND_DEREF_POINTER_EXPR:
		return true
	case KIND_INDEX_EXPR:
		i := n.Node.(*IndexExpr)
		return i.Type != nil && (i.Type.Kind == EXPR_TYPE_SLICE || i.Value.IsAddressable())
	}
	return false
}
