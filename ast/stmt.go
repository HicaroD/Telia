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
	return "BLOCK"
}
func (block BlockStmt) stmtNode() {}

type VarDeclStmt struct {
	Stmt
	Name           *token.Token
	Type           ExprType
	Value          Expr
	NeedsInference bool
}

func (variable VarDeclStmt) String() string {
	return fmt.Sprintf("Variable: %s", variable.Name)
}
func (variable VarDeclStmt) stmtNode() {}

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

type CondStmt struct {
	Stmt
	IfStmt    *IfElifCond
	ElifStmts []*IfElifCond
	ElseStmt  *ElseCond
}

func (condStmt CondStmt) String() string {
	return "IF"
}
func (cond CondStmt) stmtNode() {}

type IfElifCond struct {
	If    *token.Position
	Expr  Expr
	Block *BlockStmt
}

type ElseCond struct {
	Else  *token.Position
	Block *BlockStmt
}
