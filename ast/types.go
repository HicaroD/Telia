package ast

import (
	"github.com/HicaroD/telia-lang/lexer/token/kind"
)

type ExprType interface {
	exprTypeNode()
}

// bool, i8, i16, i32, i64, i128
type BasicType struct {
	ExprType
	Kind kind.TokenKind
}

func (basicType BasicType) exprTypeNode() {}

type PointerType struct {
	ExprType
	Type ExprType
}

func (pointer PointerType) exprTypeNode() {}
