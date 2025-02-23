package ast

import (
	"fmt"

	"github.com/HicaroD/Telia/internal/lexer/token"
)

type ExprTypeKind int

const (
	EXPR_TYPE_BASIC ExprTypeKind = iota
	EXPR_TYPE_ID
	EXPR_TYPE_POINTER
	EXPR_TYPE_ALIAS
)

type ExprType struct {
	Kind ExprTypeKind
	T    any
}

func (ty *ExprType) IsBoolean() bool {
	if ty.Kind != EXPR_TYPE_BASIC {
		return false
	}
	basic := ty.T.(*BasicType)
	return basic.Kind == token.BOOL_TYPE
}

func (ty *ExprType) IsVoid() bool {
	if ty.Kind != EXPR_TYPE_BASIC {
		return false
	}
	basic := ty.T.(*BasicType)
	return basic.Kind == token.VOID_TYPE
}

func (ty *ExprType) IsNumeric() bool {
	if ty.Kind != EXPR_TYPE_BASIC {
		return false
	}
	basic := ty.T.(*BasicType)
	return basic.Kind > token.NUMERIC_TYPE_START && basic.Kind < token.NUMERIC_TYPE_END
}

type BasicType struct {
	Kind token.Kind
}

func NewBasicType(kind token.Kind) *ExprType {
	ty := new(ExprType)
	ty.Kind = EXPR_TYPE_BASIC
	ty.T = &BasicType{Kind: kind}
	return ty
}

func (basicType BasicType) String() string {
	return basicType.Kind.String()
}

func (bt *BasicType) IsAnyStringType() bool {
	return bt.Kind == token.STRING_TYPE || bt.Kind == token.CSTRING_TYPE
}

func (bt *BasicType) IsIntegerType() bool {
	return bt.Kind > token.INTEGER_TYPE_START && bt.Kind < token.INTEGER_TYPE_END
}

type IdType struct {
	Name *token.Token
}

func (idType IdType) String() string {
	return fmt.Sprintf("IdType: %s", idType.Name.Lexeme)
}

type PointerType struct {
	Type *ExprType
}

func (pointer PointerType) String() string {
	return fmt.Sprintf("*%s", pointer.Type)
}

type TypeAlias struct {
	Node
	Name *token.Token
	Type *ExprType
}

func (alias TypeAlias) String() string {
	return fmt.Sprintf("ALIAS: %s - TYPE: %s\n", alias.Name, alias.Type)
}
