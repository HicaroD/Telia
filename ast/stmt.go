package ast

import (
	"fmt"

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

func (variable BlockStmt) String() string {
	return fmt.Sprintf("BLOCK")
}

type VarStmt struct {
	Stmt
	Name  string
	Type  ExprType
	Value Expr
}

func (variable VarStmt) String() string {
	return fmt.Sprintf("Variable: %s", variable.Name)
}

type ReturnStmt struct {
	Stmt
	Return *token.Token
	Value  Expr
}

func (ret ReturnStmt) String() string {
	return fmt.Sprintf("RETURN: %s", ret.Value)
}
