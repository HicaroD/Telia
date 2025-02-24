package ast

import (
	"fmt"

	"github.com/HicaroD/Telia/internal/lexer/token"
)

type BlockStmt struct {
	OpenCurly  token.Pos
	Statements []*Node
	CloseCurly token.Pos
}

func (block BlockStmt) String() string {
	return fmt.Sprintf("\n'{' %s\n%s\n'}' %s", block.OpenCurly, block.Statements, block.CloseCurly)
}

type MultiVarStmt struct {
	IsDecl    bool
	Variables []*Node
}

func (multi MultiVarStmt) String() string {
	return fmt.Sprintf("Multi: %v %v", multi.IsDecl, multi.Variables)
}

type VarStmt struct {
	Decl           bool
	Name           *token.Token
	Type           *ExprType
	Value          *Node
	NeedsInference bool

	BackendType any // LLVM: *values.Variable
}

func (variable VarStmt) String() string {
	return fmt.Sprintf("Variable: %s %v %s", variable.Name, variable.NeedsInference, variable.Value)
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
