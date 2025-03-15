package ast

import (
	"fmt"

	"github.com/HicaroD/Telia/internal/lexer/token"
)

type BlockStmt struct {
	OpenCurly   token.Pos
	DeferStack  []*DeferStmt
	Statements  []*Node
	FoundReturn bool
	CloseCurly  token.Pos
}

func (block BlockStmt) String() string {
	return fmt.Sprintf("\n'{' %s\n%s\n'}' %s", block.OpenCurly, block.Statements, block.CloseCurly)
}

type VarIdStmt struct {
	Name            *token.Token
	Type            *ExprType
	NeedsInference  bool
	PointerReceiver bool

	// Codegen
	BackendType any
	N           *Node
}

type VarStmt struct {
	IsDecl             bool
	HasFieldAccess     bool
	HasPointerReceiver bool
	Names              []*Node
	Expr               *Node
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
	Name     *token.Token
	Args     []*Node
	Variadic bool
	AtOp     *AtOperator

	// Codegen
	BackendType any
	Decl        *FnDecl
	Proto       *Proto
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
	Skip bool
}

func (d DeferStmt) String() string {
	return fmt.Sprintf("defer %s;", d.Stmt)
}
