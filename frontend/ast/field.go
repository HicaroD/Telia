package ast

import (
	"fmt"

	"github.com/HicaroD/Telia/frontend/lexer/token"
)

type Field struct {
	Node
	Name *token.Token
	Type ExprType
}

func (field Field) String() string {
	return fmt.Sprintf("Name: %s\nType: %s", field.Name, field.Type)
}
func (field Field) astNode() {}

// Field list for function parameters
type FieldList struct {
	Open       *token.Token
	Fields     []*Field
	IsVariadic bool
	Close      *token.Token
}

func (fieldList FieldList) String() string {
	return fmt.Sprintf(
		"\n'%s' %s\n%s\nIsVariadic: %t\n'%s' %s\n",
		fieldList.Open.Kind,
		fieldList.Open.Pos,
		fieldList.Fields,
		fieldList.IsVariadic,
		fieldList.Close.Kind,
		fieldList.Close.Pos,
	)
}
