package ast

import "github.com/HicaroD/telia-lang/lexer/token"

type AstNode interface {
	String() string
}

// Field list for function parameters
type FieldList struct {
	Open       *token.Token
	Fields     []*Field
	IsVariadic bool
	Close      *token.Token
}

type Field struct {
	Name string
	Type ExprType
}
