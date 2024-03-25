package ast

import (
	"github.com/HicaroD/telia-lang/lexer/token/kind"
)

type Decl interface {
	AstNode
	declNode()
}

type FunctionDecl struct {
	Stmt
	Name    string
	Params  *FieldList
	RetType kind.TokenKind
	Block   *BlockStmt
}
