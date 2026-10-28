package ast

import "github.com/HicaroD/Telia/internal/lexer/token"

type ConstDecl struct {
	Owner      *Package
	File       *File
	Scope      *Scope
	Name       *token.Token
	Type       *ExprType
	Expr       *Node
	Value      *LiteralExpr
	Evaluating bool
}
