package ast

import (
	"fmt"
	"github.com/HicaroD/Telia/frontend/lexer/token"
)

type Stmt interface {
	Node
	IsReturn() bool
	stmtNode()
}

type BlockStmt struct {
	Stmt
	OpenCurly  token.Pos
	Statements []Stmt
	CloseCurly token.Pos
}

func (block BlockStmt) String() string {
	return fmt.Sprintf("\n'{' %s\n%s\n'}' %s", block.OpenCurly, block.Statements, block.CloseCurly)
}
func (block BlockStmt) IsReturn() bool { return false }
func (block BlockStmt) astNode()       {}
func (block BlockStmt) stmtNode()      {}

type MultiVarStmt struct {
	Stmt
	IsDecl    bool
	Variables []*VarStmt
}

func (multi MultiVarStmt) String() string {
	return fmt.Sprintf("Multi: %v %v", multi.IsDecl, multi.Variables)
}
func (multi MultiVarStmt) IsReturn() bool { return false }
func (multi MultiVarStmt) astNode()       {}
func (multi MultiVarStmt) stmtNode()      {}

type VarStmt struct {
	Stmt
	Decl           bool
	Name           *token.Token
	Type           ExprType
	Value          Expr
	NeedsInference bool

	BackendType any // LLVM: *values.Variable
}

func (variable VarStmt) String() string {
	return fmt.Sprintf("Variable: %s %v %s", variable.Name, variable.NeedsInference, variable.Value)
}
func (variable VarStmt) IsReturn() bool { return false }
func (variable VarStmt) astNode()       {}
func (variable VarStmt) stmtNode()      {}

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
	Name *token.Token
	Args []Expr

	BackendType any
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
	If    *token.Pos
	Expr  Expr
	Block *BlockStmt
}

type ElseCond struct {
	Else  *token.Pos
	Block *BlockStmt
}

type ForLoop struct {
	Stmt
	Init   Stmt
	Cond   Expr
	Update Stmt
	Block  *BlockStmt
}

func (forLoop ForLoop) String() string {
	return fmt.Sprintf(
		"for(%s;%s;%s) %s",
		forLoop.Init,
		forLoop.Cond,
		forLoop.Update,
		forLoop.Block,
	)
}
func (forLoop ForLoop) IsReturn() bool { return false }
func (forLoop ForLoop) astNode()       {}
func (forLoop ForLoop) stmtNode()      {}

type WhileLoop struct {
	Stmt
	Cond  Expr
	Block *BlockStmt
}

func (whileLoop WhileLoop) String() string {
	return fmt.Sprintf(
		"while(%s) %s",
		whileLoop.Cond,
		whileLoop.Block,
	)
}
func (whileLoop WhileLoop) IsReturn() bool { return false }
func (whileLoop WhileLoop) astNode()       {}
func (whileLoop WhileLoop) stmtNode()      {}
