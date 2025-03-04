package ast

import (
	"fmt"

	"github.com/HicaroD/Telia/internal/lexer/token"
)

type Field struct {
	Name        *token.Token
	Type        *ExprType
	Variadic    bool
	BackendType any // LLVM: *values.Variable
}

func (field Field) String() string {
	return fmt.Sprintf("Name: %v\nType: %v", field.Name, field.Type)
}

// Field list for function parameters
type FieldList struct {
	Open       *token.Token
	Fields     []*Field
	Len        int
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
