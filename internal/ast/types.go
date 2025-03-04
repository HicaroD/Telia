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
	EXPR_TYPE_TUPLE
)

type ExprType struct {
	Kind ExprTypeKind
	T    any
}

func (t *ExprType) Equals(other *ExprType) bool {
	if t.Kind != other.Kind {
		return false
	}

	switch t.Kind {
	case EXPR_TYPE_BASIC:
		return t.T.(*BasicType).Equal(other.T.(*BasicType))
	case EXPR_TYPE_POINTER:
		return t.T.(*PointerType).Equal(other.T.(*PointerType))
	case EXPR_TYPE_TUPLE:
		return t.T.(*TupleType).Equal(other.T.(*TupleType))
	// Add other type cases
	default:
		return false
	}
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
	return basic.Kind > token.NUMERIC_TYPE_START && basic.Kind < token.NUMERIC_TYPE_END || basic.Kind == token.UNTYPED_INT
}

func (ty *ExprType) IsUntyped() bool {
	if ty.Kind != EXPR_TYPE_BASIC {
		return false
	}
	basic := ty.T.(*BasicType)
	return basic.IsUntyped()
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

func (left *BasicType) Equal(right *BasicType) bool {
	if left.Kind.IsUntyped() || right.Kind.IsUntyped() {
		return left.IsCompatibleWith(right)
	}
	return left.Kind == right.Kind
}

func (left *BasicType) IsCompatibleWith(right *BasicType) bool {
	if left.Kind.IsNumeric() && right.Kind.IsNumeric() {
		return true
	}

	if left.Kind.IsStringLiteral() && right.Kind.IsStringType() {
		return true
	}

	return false
}

func (b *BasicType) IsUntyped() bool {
	return b.Kind.IsUntyped()
}

func (basicType BasicType) String() string {
	return basicType.Kind.String()
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

func (p *PointerType) Equal(other *PointerType) bool {
	return p.Type.Equals(other.Type)
}

func (pointer PointerType) String() string {
	return fmt.Sprintf("*%v", pointer.Type)
}

type TypeAlias struct {
	Name *token.Token
	Type *ExprType
}

func (alias TypeAlias) String() string {
	return fmt.Sprintf("ALIAS: %v - TYPE: %v\n", alias.Name, alias.Type)
}

type TupleType struct {
	Types []*ExprType
}

func (t *TupleType) Equal(other *TupleType) bool {
	if len(t.Types) != len(other.Types) {
		return false
	}

	for i := range t.Types {
		if !t.Types[i].Equals(other.Types[i]) {
			return false
		}
	}
	return true
}

func (tt TupleType) String() string {
	return fmt.Sprintf("TUPLE: %v\n", tt.Types)
}

// Operator validation

type OperatorValidation struct {
	ValidTypes []*ExprType
	ResultType *ExprType
	Handler    func(operands []*ExprType) (*ExprType, error)
}

type OperatorTable map[token.Kind]OperatorValidation

var UnaryOperators = OperatorTable{
	token.MINUS: {
		ValidTypes: []*ExprType{
			NewBasicType(token.I8_TYPE),
			NewBasicType(token.I16_TYPE),
			NewBasicType(token.I32_TYPE),
			NewBasicType(token.I64_TYPE),
		},
		ResultType: nil,
		Handler:    handleNumericUnary,
	},
	token.NOT: {
		ValidTypes: []*ExprType{NewBasicType(token.BOOL_TYPE)},
		ResultType: NewBasicType(token.BOOL_TYPE),
	},
}

var BinaryOperators = OperatorTable{
	token.PLUS: {
		ValidTypes: []*ExprType{
			NewBasicType(token.I8_TYPE),
			NewBasicType(token.I16_TYPE),
			NewBasicType(token.I32_TYPE),
			NewBasicType(token.I64_TYPE),
		},
		Handler: handleNumericType,
	},
	token.MINUS: {
		ValidTypes: []*ExprType{
			NewBasicType(token.I8_TYPE),
			NewBasicType(token.I16_TYPE),
			NewBasicType(token.I32_TYPE),
			NewBasicType(token.I64_TYPE),
		},
		Handler: handleNumericType,
	},
	token.STAR: {
		ValidTypes: []*ExprType{
			NewBasicType(token.I8_TYPE),
			NewBasicType(token.I16_TYPE),
			NewBasicType(token.I32_TYPE),
			NewBasicType(token.I64_TYPE),
		},
		Handler: handleNumericType,
	},
	token.SLASH: {
		ValidTypes: []*ExprType{
			NewBasicType(token.I8_TYPE),
			NewBasicType(token.I16_TYPE),
			NewBasicType(token.I32_TYPE),
			NewBasicType(token.I64_TYPE),
		},
		Handler: handleNumericType,
	},
	token.OR: {
		ValidTypes: []*ExprType{
			NewBasicType(token.BOOL_TYPE),
		},
		ResultType: NewBasicType(token.BOOL_TYPE),
	},
	token.AND: {
		ValidTypes: []*ExprType{
			NewBasicType(token.BOOL_TYPE),
		},
		ResultType: NewBasicType(token.BOOL_TYPE),
	},
	token.LESS: {
		ValidTypes: []*ExprType{
			NewBasicType(token.I8_TYPE),
			NewBasicType(token.I16_TYPE),
			NewBasicType(token.I32_TYPE),
			NewBasicType(token.I64_TYPE),
		},
		ResultType: NewBasicType(token.BOOL_TYPE),
	},
	token.LESS_EQ: {
		ValidTypes: []*ExprType{
			NewBasicType(token.I8_TYPE),
			NewBasicType(token.I16_TYPE),
			NewBasicType(token.I32_TYPE),
			NewBasicType(token.I64_TYPE),
		},
		ResultType: NewBasicType(token.BOOL_TYPE),
	},
	token.GREATER: {
		ValidTypes: []*ExprType{
			NewBasicType(token.I8_TYPE),
			NewBasicType(token.I16_TYPE),
			NewBasicType(token.I32_TYPE),
			NewBasicType(token.I64_TYPE),
		},
		ResultType: NewBasicType(token.BOOL_TYPE),
	},
	token.EQUAL_EQUAL: {
		ValidTypes: []*ExprType{
			NewBasicType(token.I8_TYPE),
			NewBasicType(token.I16_TYPE),
			NewBasicType(token.I32_TYPE),
			NewBasicType(token.I64_TYPE),
		},
		ResultType: NewBasicType(token.BOOL_TYPE),
	},
	token.BANG_EQUAL: {
		ValidTypes: []*ExprType{
			NewBasicType(token.I8_TYPE),
			NewBasicType(token.I16_TYPE),
			NewBasicType(token.I32_TYPE),
			NewBasicType(token.I64_TYPE),
		},
		ResultType: NewBasicType(token.BOOL_TYPE),
	},
}

func handleNumericUnary(operands []*ExprType) (*ExprType, error) {
	return operands[0], nil
}

func handleNumericType(operands []*ExprType) (*ExprType, error) {
	leftType := operands[0]
	rightType := operands[1]

	if !leftType.Equals(rightType) {
		return nil, fmt.Errorf("type mismatch")
	}
	return leftType, nil
}
