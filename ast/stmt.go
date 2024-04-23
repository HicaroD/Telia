package ast

import (
	"fmt"

	"github.com/HicaroD/Telia/lexer/token"
)

type Stmt interface {
	Node
	IsReturn() bool
	stmtNode()
}

type BlockStmt struct {
	Stmt
	OpenCurly  token.Position
	Statements []Stmt
	CloseCurly token.Position
}

func (block BlockStmt) String() string {
	return fmt.Sprintf("\n'{' %s\n%s\n'}' %s", block.OpenCurly, block.Statements, block.CloseCurly)
}
func (block BlockStmt) IsReturn() bool { return false }
func (block BlockStmt) astNode()       {}
func (block BlockStmt) stmtNode()      {}

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
func (variable VarDeclStmt) IsReturn() bool { return false }
func (variable VarDeclStmt) astNode()       {}
func (variable VarDeclStmt) stmtNode()      {}

type ReturnStmt struct {
	Stmt
	Return *token.Token
	Value  Expr
}

func (ret ReturnStmt) String() string {
	return fmt.Sprintf("RETURN: %s", ret.Value)
}
func (ret ReturnStmt) IsReturn() bool { return true }
func (ret ReturnStmt) astNode()       {}
func (ret ReturnStmt) stmtNode()      {}

type FunctionCall struct {
	Stmt
	Expr
	Name string
	Args []Expr
}

func (call FunctionCall) String() string {
	return fmt.Sprintf("CALL: %s - ARGS: %s", call.Name, call.Args)
}
func (call FunctionCall) IsReturn() bool { return false }
func (call FunctionCall) astNode()       {}
func (call FunctionCall) stmtNode()      {}
func (call FunctionCall) exprNode()      {}

type CondStmt struct {
	Stmt
	IfStmt    *IfElifCond
	ElifStmts []*IfElifCond
	ElseStmt  *ElseCond
}

func (condStmt CondStmt) String() string {
	return "IF"
}
func (cond CondStmt) IsReturn() bool { return false }
func (cond CondStmt) astNode()       {}
func (cond CondStmt) stmtNode()      {}

type IfElifCond struct {
	If    *token.Position
	Expr  Expr
	Block *BlockStmt
}

type ElseCond struct {
	Else  *token.Position
	Block *BlockStmt
}
