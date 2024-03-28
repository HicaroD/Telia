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

func (block BlockStmt) String() string {
	return fmt.Sprintf("BLOCK")
}
func (block BlockStmt) stmtNode() {}

type VarStmt struct {
	Stmt
	Name  string
	Type  ExprType
	Value Expr
}

func (variable VarStmt) String() string {
	return fmt.Sprintf("Variable: %s", variable.Name)
}
func (variable VarStmt) stmtNode() {}

type ReturnStmt struct {
	Stmt
	Return *token.Token
	Value  Expr
}

func (ret ReturnStmt) String() string {
	return fmt.Sprintf("RETURN: %s", ret.Value)
}
func (ret ReturnStmt) stmtNode() {}

type FunctionCallStmt struct {
	Stmt
	Name string
	Args []Expr
}

func (call FunctionCallStmt) String() string {
	return fmt.Sprintf("CALL: %s - ARGS: %s", call.Name, call.Args)
}
func (call FunctionCallStmt) stmtNode() {}
