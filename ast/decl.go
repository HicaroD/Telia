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
	Name    string
	Params  *FieldList
	RetType ExprType
	Block   *BlockStmt
}

func (fnDecl FunctionDecl) String() string {
	return fmt.Sprintf("FN: %s", fnDecl.Name)
}

/*
Extern blocks contains a list of function prototypes.

extern "C" {
  fn printf(format *i8, ...) i32;
}
*/

type ExternDecl struct {
	Decl
	Name       *token.Token
	Prototypes []*Proto
}

func (extern ExternDecl) String() string {
	return fmt.Sprintf("EXTERN: %s", extern.Name)
}

type Proto struct {
	Name    string
	Params  *FieldList
	RetType ExprType
}
