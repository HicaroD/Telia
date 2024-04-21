package ast

import (
	"fmt"

	"github.com/HicaroD/Telia/lexer/token"
	"github.com/HicaroD/Telia/scope"
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
	return fmt.Sprintf(
		"Scope: %s\nName: %s\nParams: %s\nRetType: %s\nBlock: %s\n",
		fnDecl.Scope,
		fnDecl.Name,
		fnDecl.Params,
		fnDecl.RetType,
		fnDecl.Block,
	)
}
func (fnDecl FunctionDecl) astNode()  {}
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
func (extern ExternDecl) astNode()  {}
func (extern ExternDecl) declNode() {}

// NOTE: Proto implementing AstNode is temporary
type Proto struct {
	AstNode
	Name    string
	Params  *FieldList
	RetType ExprType
}

func (proto Proto) String() string { return fmt.Sprintf("PROTO: %s", proto.Name) }
func (proto Proto) astNode()       {}
