package ast

import (
	"fmt"

	"github.com/HicaroD/telia-lang/lexer/token"
	"github.com/HicaroD/telia-lang/scope"
)

type Decl interface {
	AstNode
	declNode()
}

type FunctionDecl struct {
	Decl
	Scope   *scope.Scope[AstNode]
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
	Scope      *scope.Scope[AstNode]
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
