package ast

import (
	"fmt"

	"github.com/HicaroD/Telia/internal/lexer/token"
)

var (
	RAWPTR_TYPE = NewBasicType(token.RAWPTR_TYPE)
)

type ExprTypeKind int

const (
	EXPR_TYPE_BASIC ExprTypeKind = iota
	EXPR_TYPE_ID
	EXPR_TYPE_STRUCT
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
		return t.T.(*BasicType).Equals(other.T.(*BasicType))
	case EXPR_TYPE_POINTER:
		return t.T.(*PointerType).Equals(other.T.(*PointerType))
	case EXPR_TYPE_TUPLE:
		return t.T.(*TupleType).Equals(other.T.(*TupleType))
	case EXPR_TYPE_ID:
		leftId := t.T.(*IdType)
		rightId := other.T.(*IdType)
		return leftId.Name.Name() == rightId.Name.Name()
	case EXPR_TYPE_STRUCT:
		return t.T.(*StructType).Equals(other.T.(*StructType))
	default:
		return false
	}
}

func (ty *ExprType) Promote() error {
	switch ty.Kind {
	case EXPR_TYPE_BASIC:
		basicType := ty.T.(*BasicType)
		if basicType.Kind.IsInteger() {
			ty.T = &BasicType{Kind: token.INT_TYPE}
		} else if basicType.Kind.IsFloat() {
			ty.T = &BasicType{Kind: token.FLOAT_TYPE}
		}
	case EXPR_TYPE_POINTER:
		ptrType := ty.T.(*PointerType)
		if err := ptrType.Type.Promote(); err != nil {
			return err
		}
	case EXPR_TYPE_TUPLE:
		tupleType := ty.T.(*TupleType)
		for _, t := range tupleType.Types {
			if err := t.Promote(); err != nil {
				return err
			}
		}
		// Add more cases for other types, such as structs or aliases, if necessary.
	}
	return nil
}

func (ty *ExprType) IsNumeric() bool {
	if ty.Kind != EXPR_TYPE_BASIC {
		return false
	}
	basic := ty.T.(*BasicType)
	return basic.Kind.IsNumeric()
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

func (ty *ExprType) IsInteger() bool {
	if ty.Kind != EXPR_TYPE_BASIC {
		return false
	}
	basic := ty.T.(*BasicType)
	return basic.Kind > token.INTEGER_TYPE_START && basic.Kind < token.NUMERIC_TYPE_END ||
		basic.Kind == token.INTEGER_TYPE_END
}

func (ty *ExprType) IsFloat() bool {
	if ty.Kind != EXPR_TYPE_BASIC {
		return false
	}
	basic := ty.T.(*BasicType)
	return basic.Kind > token.FLOAT_TYPE_START && basic.Kind < token.NUMERIC_TYPE_END ||
		basic.Kind == token.FLOAT_TYPE_END
}

func (ty *ExprType) IsPointer() bool {
	return ty.Kind == EXPR_TYPE_POINTER
}

func (ty *ExprType) PointerTo(pointeeTy ExprTypeKind) bool {
	if !ty.IsPointer() {
		return false
	}
	ptr := ty.T.(*PointerType)
	return ptr.Type.Kind == pointeeTy
}

func (ty *ExprType) IsUntyped() bool {
	if ty.Kind == EXPR_TYPE_BASIC {
		basic := ty.T.(*BasicType)
		return basic.IsUntyped()
	}
	if ty.Kind == EXPR_TYPE_POINTER {
		ptr := ty.T.(*PointerType)
		return ptr.Type.IsUntyped()
	}
	return false
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

func NewPointerType(ty *ExprType) *ExprType {
	t := new(ExprType)
	t.Kind = EXPR_TYPE_POINTER
	t.T = &PointerType{Type: ty}
	return t
}

func (left *BasicType) Equals(right *BasicType) bool {
	if left.Kind.IsUntyped() || right.Kind.IsUntyped() {
		return left.IsCompatibleWith(right)
	}
	return left.Kind == right.Kind
}

func (left *BasicType) IsCompatibleWith(right *BasicType) bool {
	if left.Kind.IsInteger() && right.Kind.IsInteger() {
		return true
	}

	if left.Kind.IsFloat() && right.Kind.IsFloat() {
		return true
	}

	if left.Kind.IsStringLiteral() && right.Kind.IsStringType() {
		return true
	}

	if left.Kind == token.UNTYPED_BOOL && right.Kind == token.BOOL_TYPE {
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

func (p *PointerType) Equals(other *PointerType) bool {
	return p.Type.Equals(other.Type)
}

func (pointer PointerType) String() string {
	return fmt.Sprintf("*%s", pointer.Type.T)
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

func (t *TupleType) Equals(other *TupleType) bool {
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

type StructType struct {
	Decl *StructDecl
}

func (st *StructType) Equals(other *StructType) bool {
	return st.Decl.Name == other.Decl.Name
}

// Operator validation

type OperatorValidation struct {
	ValidTypes []*ExprType
	ResultType *ExprType
	Handler    func(operands []*ExprType) (*ExprType, error)
}

type OperatorTable map[token.Kind]OperatorValidation

// Useful for dealing with recursive pointer type during semantic analysis
var POINTER_TYPE = &ExprType{Kind: EXPR_TYPE_POINTER, T: new(PointerType)}

var UnaryOperators = OperatorTable{
	token.MINUS: {
		ValidTypes: []*ExprType{
			// integers
			NewBasicType(token.UINT_TYPE),
			NewBasicType(token.INT_TYPE),
			NewBasicType(token.I8_TYPE),
			NewBasicType(token.I16_TYPE),
			NewBasicType(token.I32_TYPE),
			NewBasicType(token.I64_TYPE),
			NewBasicType(token.I128_TYPE),
			NewBasicType(token.U8_TYPE),
			NewBasicType(token.U16_TYPE),
			NewBasicType(token.U32_TYPE),
			NewBasicType(token.U64_TYPE),
			NewBasicType(token.U128_TYPE),
			// floats
			NewBasicType(token.FLOAT_TYPE),
			NewBasicType(token.F32_TYPE),
			NewBasicType(token.F64_TYPE),
		},
		ResultType: nil,
		Handler:    handleNumericUnary,
	},
	token.NOT: {
		ValidTypes: []*ExprType{
			NewBasicType(token.BOOL_TYPE),
			NewBasicType(token.UNTYPED_BOOL),
		},
		ResultType: NewBasicType(token.BOOL_TYPE),
	},
}

var BinaryOperators = OperatorTable{
	token.PLUS: {
		ValidTypes: []*ExprType{
			// integers
			NewBasicType(token.UINT_TYPE),
			NewBasicType(token.INT_TYPE),
			NewBasicType(token.I8_TYPE),
			NewBasicType(token.I16_TYPE),
			NewBasicType(token.I32_TYPE),
			NewBasicType(token.I64_TYPE),
			NewBasicType(token.I128_TYPE),
			NewBasicType(token.U8_TYPE),
			NewBasicType(token.U16_TYPE),
			NewBasicType(token.U32_TYPE),
			NewBasicType(token.U64_TYPE),
			NewBasicType(token.U128_TYPE),
			// floats
			NewBasicType(token.FLOAT_TYPE),
			NewBasicType(token.F32_TYPE),
			NewBasicType(token.F64_TYPE),
		},
		Handler: handleNumericType,
	},
	token.MINUS: {
		ValidTypes: []*ExprType{
			// integers
			NewBasicType(token.UINT_TYPE),
			NewBasicType(token.INT_TYPE),
			NewBasicType(token.I8_TYPE),
			NewBasicType(token.I16_TYPE),
			NewBasicType(token.I32_TYPE),
			NewBasicType(token.I64_TYPE),
			NewBasicType(token.I128_TYPE),
			NewBasicType(token.U8_TYPE),
			NewBasicType(token.U16_TYPE),
			NewBasicType(token.U32_TYPE),
			NewBasicType(token.U64_TYPE),
			NewBasicType(token.U128_TYPE),
			// floats
			NewBasicType(token.FLOAT_TYPE),
			NewBasicType(token.F32_TYPE),
			NewBasicType(token.F64_TYPE),
		},
		Handler: handleNumericType,
	},
	token.STAR: {
		ValidTypes: []*ExprType{
			// integers
			NewBasicType(token.UINT_TYPE),
			NewBasicType(token.INT_TYPE),
			NewBasicType(token.I8_TYPE),
			NewBasicType(token.I16_TYPE),
			NewBasicType(token.I32_TYPE),
			NewBasicType(token.I64_TYPE),
			NewBasicType(token.I128_TYPE),
			NewBasicType(token.U8_TYPE),
			NewBasicType(token.U16_TYPE),
			NewBasicType(token.U32_TYPE),
			NewBasicType(token.U64_TYPE),
			NewBasicType(token.U128_TYPE),
			// floats
			NewBasicType(token.FLOAT_TYPE),
			NewBasicType(token.F32_TYPE),
			NewBasicType(token.F64_TYPE),
		},
		Handler: handleNumericType,
	},
	token.SLASH: {
		ValidTypes: []*ExprType{
			// integers
			NewBasicType(token.UINT_TYPE),
			NewBasicType(token.INT_TYPE),
			NewBasicType(token.I8_TYPE),
			NewBasicType(token.I16_TYPE),
			NewBasicType(token.I32_TYPE),
			NewBasicType(token.I64_TYPE),
			NewBasicType(token.I128_TYPE),
			NewBasicType(token.U8_TYPE),
			NewBasicType(token.U16_TYPE),
			NewBasicType(token.U32_TYPE),
			NewBasicType(token.U64_TYPE),
			NewBasicType(token.U128_TYPE),
			// floats
			NewBasicType(token.FLOAT_TYPE),
			NewBasicType(token.F32_TYPE),
			NewBasicType(token.F64_TYPE),
		},
		Handler: handleNumericType,
	},
	token.OR: {
		ValidTypes: []*ExprType{NewBasicType(token.BOOL_TYPE)},
		ResultType: NewBasicType(token.BOOL_TYPE),
	},
	token.AND: {
		ValidTypes: []*ExprType{NewBasicType(token.BOOL_TYPE)},
		ResultType: NewBasicType(token.BOOL_TYPE),
	},
	token.LESS: {
		ValidTypes: []*ExprType{
			// integers
			NewBasicType(token.UINT_TYPE),
			NewBasicType(token.INT_TYPE),
			NewBasicType(token.I8_TYPE),
			NewBasicType(token.I16_TYPE),
			NewBasicType(token.I32_TYPE),
			NewBasicType(token.I64_TYPE),
			NewBasicType(token.I128_TYPE),
			NewBasicType(token.U8_TYPE),
			NewBasicType(token.U16_TYPE),
			NewBasicType(token.U32_TYPE),
			NewBasicType(token.U64_TYPE),
			NewBasicType(token.U128_TYPE),
			// floats
			NewBasicType(token.FLOAT_TYPE),
			NewBasicType(token.F32_TYPE),
			NewBasicType(token.F64_TYPE),
		},
		ResultType: NewBasicType(token.BOOL_TYPE),
	},
	token.LESS_EQ: {
		ValidTypes: []*ExprType{
			// integers
			NewBasicType(token.UINT_TYPE),
			NewBasicType(token.INT_TYPE),
			NewBasicType(token.I8_TYPE),
			NewBasicType(token.I16_TYPE),
			NewBasicType(token.I32_TYPE),
			NewBasicType(token.I64_TYPE),
			NewBasicType(token.I128_TYPE),
			NewBasicType(token.U8_TYPE),
			NewBasicType(token.U16_TYPE),
			NewBasicType(token.U32_TYPE),
			NewBasicType(token.U64_TYPE),
			NewBasicType(token.U128_TYPE),
			// floats
			NewBasicType(token.FLOAT_TYPE),
			NewBasicType(token.F32_TYPE),
			NewBasicType(token.F64_TYPE),
		},
		ResultType: NewBasicType(token.BOOL_TYPE),
	},
	token.GREATER: {
		ValidTypes: []*ExprType{
			// integers
			NewBasicType(token.UINT_TYPE),
			NewBasicType(token.INT_TYPE),
			NewBasicType(token.I8_TYPE),
			NewBasicType(token.I16_TYPE),
			NewBasicType(token.I32_TYPE),
			NewBasicType(token.I64_TYPE),
			NewBasicType(token.I128_TYPE),
			NewBasicType(token.U8_TYPE),
			NewBasicType(token.U16_TYPE),
			NewBasicType(token.U32_TYPE),
			NewBasicType(token.U64_TYPE),
			NewBasicType(token.U128_TYPE),
			// floats
			NewBasicType(token.FLOAT_TYPE),
			NewBasicType(token.F32_TYPE),
			NewBasicType(token.F64_TYPE),
		},
		ResultType: NewBasicType(token.BOOL_TYPE),
	},
	token.GREATER_EQ: {
		ValidTypes: []*ExprType{
			// integers
			NewBasicType(token.UINT_TYPE),
			NewBasicType(token.INT_TYPE),
			NewBasicType(token.I8_TYPE),
			NewBasicType(token.I16_TYPE),
			NewBasicType(token.I32_TYPE),
			NewBasicType(token.I64_TYPE),
			NewBasicType(token.I128_TYPE),
			NewBasicType(token.U8_TYPE),
			NewBasicType(token.U16_TYPE),
			NewBasicType(token.U32_TYPE),
			NewBasicType(token.U64_TYPE),
			NewBasicType(token.U128_TYPE),
			// floats
			NewBasicType(token.FLOAT_TYPE),
			NewBasicType(token.F32_TYPE),
			NewBasicType(token.F64_TYPE),
		},
		ResultType: NewBasicType(token.BOOL_TYPE),
	},
	token.EQUAL_EQUAL: {
		ValidTypes: []*ExprType{
			// integers
			NewBasicType(token.UINT_TYPE),
			NewBasicType(token.INT_TYPE),
			NewBasicType(token.I8_TYPE),
			NewBasicType(token.I16_TYPE),
			NewBasicType(token.I32_TYPE),
			NewBasicType(token.I64_TYPE),
			NewBasicType(token.I128_TYPE),
			NewBasicType(token.U8_TYPE),
			NewBasicType(token.U16_TYPE),
			NewBasicType(token.U32_TYPE),
			NewBasicType(token.U64_TYPE),
			NewBasicType(token.U128_TYPE),
			// floats
			NewBasicType(token.FLOAT_TYPE),
			NewBasicType(token.F32_TYPE),
			NewBasicType(token.F64_TYPE),
			NewBasicType(token.F64_TYPE),
			// pointer
			POINTER_TYPE,
			NewBasicType(token.RAWPTR_TYPE),
		},
		Handler: handleEqualityComparison,
	},
	token.BANG_EQUAL: {
		ValidTypes: []*ExprType{
			NewBasicType(token.UINT_TYPE),
			NewBasicType(token.INT_TYPE),
			NewBasicType(token.I8_TYPE),
			NewBasicType(token.I16_TYPE),
			NewBasicType(token.I32_TYPE),
			NewBasicType(token.I64_TYPE),
			NewBasicType(token.I128_TYPE),
			NewBasicType(token.U8_TYPE),
			NewBasicType(token.U16_TYPE),
			NewBasicType(token.U32_TYPE),
			NewBasicType(token.U64_TYPE),
			NewBasicType(token.U128_TYPE),
			// floats
			NewBasicType(token.FLOAT_TYPE),
			NewBasicType(token.F32_TYPE),
			NewBasicType(token.F64_TYPE),
			// pointers
			POINTER_TYPE,
			NewBasicType(token.RAWPTR_TYPE),
		},
		Handler: handleEqualityComparison,
	},
}

func handleEqualityComparison(operands []*ExprType) (*ExprType, error) {
	leftType := operands[0]
	rightType := operands[1]

	boolTy := NewBasicType(token.BOOL_TYPE)

	if leftType.IsPointer() && rightType.IsPointer() {
		if leftType.Equals(rightType) {
			return boolTy, nil
		}
		return nil, fmt.Errorf("type mismatch: cannot compare %s with %s", leftType.T, rightType.T)
	}

	if !leftType.IsPointer() && !rightType.IsPointer() {
		return boolTy, nil
	}

	return nil, fmt.Errorf("type mismatch: cannot compare pointer and non-pointer type")
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
