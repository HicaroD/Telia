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

type VarStmt struct {
	Stmt
	Name           *token.Token
	Type           ExprType
	Value          Expr
	NeedsInference bool
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

type CondStmt struct {
	IfStmt    *IfCondStmt
	ElifStmts []*ElifCondStmt
	ElseStmt  *ElseCondStmt
}

func (condStmt CondStmt) String() string {
	return ""
}
func (cond CondStmt) stmtNode() {}

type IfCondStmt struct {
	If    *token.Position
	Expr  Expr
	Block *BlockStmt
}

// NOTE: it is the same as IfCondStmt, I might want to reuse the same struct
type ElifCondStmt struct {
	Elif  *token.Position
	Expr  Expr
	Block *BlockStmt
}

// NOTE: Maybe use an alias to *BlockStmt given that this struct only has a single field
type ElseCondStmt struct {
	Else  *token.Position
	Block *BlockStmt
}
