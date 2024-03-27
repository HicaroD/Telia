package ast

import (
	"github.com/HicaroD/telia-lang/lexer/token/kind"
)

type ExprType interface {
	exprType()
}

// bool, i8, i16, i32, i64, i128
type BasicType struct {
	ExprType
	Kind kind.TokenKind
}

type PointerType struct {
	ExprType
	Type ExprType
}
