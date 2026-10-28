package ast

import "github.com/HicaroD/Telia/internal/lexer/token"

type ConstDecl struct {
	Foreign    *ExternDecl
	Attributes Attributes
	Owner      *Package
	File       *File
	Scope      *Scope
	Name       *token.Token
	Type       *ExprType
	Expr       *Node
	Value      *LiteralExpr
	Evaluating bool
}
