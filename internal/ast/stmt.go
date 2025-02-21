package ast

import (
	"fmt"

	"github.com/HicaroD/Telia/internal/lexer/token"
)

type Stmt interface {
	Node
	IsReturn() bool
	stmtNode()
}

type BlockStmt struct {
	Stmt
	OpenCurly  token.Pos
	Statements []*MyNode
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
	Variables []*MyNode
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
	Type           *MyExprType
	Value          *MyNode
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
	Value  *MyNode
}

func (ret ReturnStmt) String() string {
	return fmt.Sprintf("RETURN: %s", ret.Value)
}
func (ret ReturnStmt) IsReturn() bool { return true }
func (ret ReturnStmt) astNode()       {}
func (ret ReturnStmt) stmtNode()      {}

type FnCall struct {
	Stmt
	Expr
	Name *token.Token
	Args []*MyNode

	BackendType any
}

func (call FnCall) String() string {
	return fmt.Sprintf("CALL: %s - ARGS: %s", call.Name, call.Args)
}
func (call FnCall) IsReturn() bool { return false }
func (call FnCall) astNode()       {}
func (call FnCall) stmtNode()      {}
func (call FnCall) exprNode()      {}

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
	Expr  *MyNode
	Block *BlockStmt
	Scope *Scope
}

type ElseCond struct {
	Else  *token.Pos
	Block *BlockStmt
	Scope *Scope
}

type ForLoop struct {
	Stmt
	Init   *MyNode
	Cond   *MyNode
	Update *MyNode
	Block  *BlockStmt
	Scope  *Scope
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
	Cond  *MyNode
	Block *BlockStmt
	Scope *Scope
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
