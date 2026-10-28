package ast

import "github.com/HicaroD/Telia/internal/lexer/token"

type OpaqueDecl struct {
	Name       *token.Token
	Attributes Attributes
}

type OpaqueType struct{ Decl *OpaqueDecl }

func (o *OpaqueType) String() string { return o.Decl.Name.Name() }
