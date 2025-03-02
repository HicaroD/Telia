package ast

import (
	"fmt"

	"github.com/HicaroD/Telia/internal/lexer/token"
)

type BlockStmt struct {
	OpenCurly  token.Pos
	DeferStack []*Node
	Statements []*Node
	CloseCurly token.Pos
}

func (block BlockStmt) String() string {
	return fmt.Sprintf("\n'{' %s\n%s\n'}' %s", block.OpenCurly, block.Statements, block.CloseCurly)
}

type VarId struct {
	Name           *token.Token
	Type           *ExprType
	NeedsInference bool
	BackendType    any
}

type VarStmt struct {
	IsDecl bool
	Names  []*VarId
	Expr   *Node
}

func (v VarStmt) String() string {
	return fmt.Sprintf("VarStmt: %v %v %v", v.IsDecl, v.Names, v.Expr)
}

type ReturnStmt struct {
	Return *token.Token
	Value  *Node
}

func (ret ReturnStmt) String() string {
	return fmt.Sprintf("RETURN: %s", ret.Value)
}

type FnCall struct {
	Name *token.Token
	Args []*Node
	AtOp *AtOperator

	BackendType any
}

func (call FnCall) String() string {
	return fmt.Sprintf("CALL: %s - ARGS: %s", call.Name, call.Args)
}

type CondStmt struct {
	IfStmt    *IfElifCond
	ElifStmts []*IfElifCond
	ElseStmt  *ElseCond
}

func (condStmt CondStmt) String() string {
	return "IF"
}

type IfElifCond struct {
	If    *token.Pos
	Expr  *Node
	Block *BlockStmt
	Scope *Scope
}

type ElseCond struct {
	Else  *token.Pos
	Block *BlockStmt
	Scope *Scope
}

type ForLoop struct {
	Init   *Node
	Cond   *Node
	Update *Node
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

type WhileLoop struct {
	Cond  *Node
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

type DeferStmt struct {
	Stmt *Node
}

func (d DeferStmt) String() string {
	return fmt.Sprintf("defer %s;", d.Stmt)
}
