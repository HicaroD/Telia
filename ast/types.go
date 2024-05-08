package ast

import (
	"fmt"

	"github.com/HicaroD/Telia/lexer/token"
	"github.com/HicaroD/Telia/lexer/token/kind"
)

type ExprType interface {
	IsVoid() bool
	IsBoolean() bool
	IsNumeric() bool
	exprTypeNode()
}

// void, bool, int, i8, i16, i32, i64, uint, u8, u16, u32, u64
type BasicType struct {
	ExprType
	Kind kind.TokenKind
}

func (basicType BasicType) IsNumeric() bool {
	_, ok := kind.NUMERIC_TYPES[basicType.Kind]
	return ok
}
func (basicType BasicType) IsBoolean() bool { return basicType.Kind == kind.BOOL_TYPE }
func (basicType BasicType) IsVoid() bool    { return basicType.Kind == kind.VOID_TYPE }
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
	Type ExprType
}

func (pointer PointerType) IsNumeric() bool { return false }
func (pointer PointerType) IsBoolean() bool { return false }
func (pointer PointerType) IsVoid() bool    { return false }
func (pointer PointerType) exprTypeNode()   {}
func (pointer PointerType) String() string {
	return fmt.Sprintf("*%s", pointer.Type)
}

type TupleType struct {
	ExprType
	OpenParen  token.Position
	Types      []ExprType
	CloseParen token.Position
}

func (tuple TupleType) IsNumeric() bool { return false }
func (tuple TupleType) IsBoolean() bool { return false }
func (tuple TupleType) IsVoid() bool    { return false }
func (tuple TupleType) exprTypeNode()   {}
func (tuple TupleType) String() string {
	return fmt.Sprintf("(%s)", tuple.Types)
}
