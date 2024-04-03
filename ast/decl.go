package ast

import (
	"fmt"

	"github.com/HicaroD/telia-lang/lexer/token"
)

type Decl interface {
	AstNode
	declNode()
}

type FunctionDecl struct {
	Decl
	Scope   *Scope
	Name    string
	Params  *FieldList
	RetType ExprType
	Block   *BlockStmt
}

func (fnDecl FunctionDecl) String() string {
	return fmt.Sprintf("FN: %s", fnDecl.Name)
}
func (fnDecl FunctionDecl) declNode() {}

type ExternDecl struct {
	Decl
	Scope      *Scope
	Name       *token.Token
	Prototypes []*Proto
}

func (extern ExternDecl) String() string {
	return fmt.Sprintf("EXTERN: %s", extern.Name)
}
func (extern ExternDecl) declNode() {}

// NOTE: Proto implementing AstNode is temporary
type Proto struct {
	AstNode
	Name    string
	Params  *FieldList
	RetType ExprType
}

func (proto Proto) String() string { return "" }
