package ast

import (
	"github.com/HicaroD/telia-lang/lexer/token"
)

type Stmt interface {
	AstNode
	stmtNode()
}

type BlockStmt struct {
	Stmt
	OpenCurly  token.Position
	Statements []Stmt
	CloseCurly token.Position
}

type VarStmt struct {
	Stmt
	Name  string
	Type  ExprType
	Value Expr
}

type ReturnStmt struct {
	Stmt
	Return *token.Token
	Value  Expr
}
