package ast

import (
	"fmt"

	"github.com/HicaroD/Telia/internal/lexer/token"
)

type ParamAttributes struct {
	Const bool // @const
	ForC  bool // @for_c - only allowed on prototypes
}

func (p *ParamAttributes) String() string {
	return fmt.Sprintf("ATTRIBUTES - @for_c=%v @const=%v\n", p.ForC, p.Const)
}

type Param struct {
	Attributes  *ParamAttributes
	Name        *token.Token
	Type        *ExprType
	Variadic    bool
	BackendType any // LLVM: *values.Variable
}

func (param Param) String() string {
	return fmt.Sprintf("Name: %v\nType: %v", param.Name, param.Type)
}

// Field list for function parameters
type Params struct {
	Open       *token.Token
	Fields     []*Param
	Len        int
	IsVariadic bool
	Close      *token.Token
}

func (fieldList Params) String() string {
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
