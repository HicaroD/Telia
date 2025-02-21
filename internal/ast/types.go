package ast

import (
	"fmt"

	"github.com/HicaroD/Telia/internal/lexer/token"
)

type ExprTypeEnum int

const (
	EXPR_TYPE_BASIC ExprTypeEnum = iota
	EXPR_TYPE_ID
	EXPR_TYPE_POINTER
	EXPR_TYPE_ALIAS
)

type MyExprType struct {
	Kind ExprTypeEnum
	T    any
}

type ExprType interface {
	IsVoid() bool
	IsBoolean() bool
	IsNumeric() bool
	exprTypeNode()
}

// void, bool, int, i8, i16, i32, i64, uint,
// u8, u16, u32, u64, string, cstring
type BasicType struct {
	ExprType
	Kind token.Kind
}

func (basicType BasicType) IsNumeric() bool {
	_, ok := token.NUMERIC_TYPES[basicType.Kind]
	return ok
}
func (basicType BasicType) IsBoolean() bool { return basicType.Kind == token.BOOL_TYPE }
func (basicType BasicType) IsVoid() bool    { return basicType.Kind == token.VOID_TYPE }
func (basicType BasicType) exprTypeNode()   {}
func (basicType BasicType) String() string {
	return basicType.Kind.String()
}

type IdType struct {
	ExprType
	Name *token.Token
}

func (idType IdType) IsNumeric() bool { return false }
func (idType IdType) IsBoolean() bool { return false }
func (idType IdType) IsVoid() bool    { return false }
func (idType IdType) exprTypeNode()   {}
func (idType IdType) String() string {
	return fmt.Sprintf("IdType: %s", idType.Name.Lexeme)
}

type PointerType struct {
	ExprType
	Type *MyExprType
}

func (pointer PointerType) IsNumeric() bool { return false }
func (pointer PointerType) IsBoolean() bool { return false }
func (pointer PointerType) IsVoid() bool    { return false }
func (pointer PointerType) exprTypeNode()   {}
func (pointer PointerType) String() string {
	return fmt.Sprintf("*%s", pointer.Type)
}

type TypeAlias struct {
	Node
	ExprType
	Name *token.Token
	Type *MyExprType
}

func (alias TypeAlias) astNode()        {}
func (alias TypeAlias) IsNumeric() bool { return false }
func (alias TypeAlias) IsBoolean() bool { return false }
func (alias TypeAlias) IsVoid() bool    { return false }
func (alias TypeAlias) exprTypeNode()   {}
func (alias TypeAlias) String() string {
	return fmt.Sprintf("ALIAS: %s - TYPE: %s\n", alias.Name, alias.Type)
}
