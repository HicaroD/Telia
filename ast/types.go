package ast

import (
	"fmt"

	"github.com/HicaroD/telia-lang/lexer/token/kind"
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
	exprTypeNode()
}

// void, bool, i8, i16, i32, i64, u8, u16, u32, u64
type BasicType struct {
	ExprType
	Kind kind.TokenKind
}

func (basicType BasicType) IsVoid() bool  { return basicType.Kind == kind.VOID_TYPE }
func (basicType BasicType) exprTypeNode() {}
func (basicType BasicType) String() string {
	return basicType.Kind.String()
}

type PointerType struct {
	ExprType
	Type ExprType
}

func (pointer PointerType) IsVoid() bool  { return false }
func (pointer PointerType) exprTypeNode() {}
func (pointer PointerType) String() string {
	return fmt.Sprintf("*%s", pointer.Type)
}
