package ast

import (
	"fmt"

	"github.com/HicaroD/Telia/lexer/token/kind"
)

var LOGICAL_OP map[kind.TokenKind]bool = map[kind.TokenKind]bool{
	kind.AND:         true,
	kind.OR:          true,
	kind.BANG_EQUAL:  true,
	kind.EQUAL_EQUAL: true,
	kind.GREATER:     true,
	kind.GREATER_EQ:  true,
	kind.LESS:        true,
	kind.LESS_EQ:     true,
}

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
