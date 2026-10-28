package sema

import (
	"fmt"
	"go/constant"
	gotoken "go/token"

	"github.com/HicaroD/Telia/internal/ast"
	"github.com/HicaroD/Telia/internal/diagnostics"
	"github.com/HicaroD/Telia/internal/lexer/token"
)

func (s *sema) inferShift(b *ast.BinExpr, expected *ast.ExprType, reference, declaration *ast.Scope, imported bool) (*ast.ExprType, bool, error) {
	left, explicit, err := s.inferExprTypeWithoutContext(b.Left, reference, declaration, imported, false)
	if err != nil {
		return nil, false, err
	}
	if expected != nil && !explicit && expected.IsInteger() {
		left, err = s.inferExprTypeWithContext(b.Left, expected, reference, declaration, imported, false)
		if err != nil {
			return nil, false, err
		}
	}
	right, _, err := s.inferExprTypeWithoutContext(b.Right, reference, declaration, imported, false)
	if err != nil {
		return nil, false, err
	}
	if !left.IsInteger() || !right.IsInteger() {
		return nil, false, diagnostics.Diag{Title: "Invalid shift operands", Message: "shifts require integer operands", Span: b.Left.Span, Help: "Use an integer value and an integer count; the count may have a different integer type."}
	}
	if expected != nil && !left.Equals(expected) {
		return nil, false, diagnostics.Diag{Title: "Shift result type mismatch", Message: fmt.Sprintf("shift has type %s, expected %s", left.T, expected.T), Span: b.Left.Span, Help: "The left operand determines the result type; convert it to the required integer type."}
	}
	b.OperandType, b.RightType = left, right
	if value, err := evaluateConstant(b.Right); err == nil {
		if err := validateShiftCount(value, b); err != nil {
			return nil, false, err
		}
	}
	return left, explicit, nil
}

func validateShiftCount(value constant.Value, b *ast.BinExpr) error {
	width := b.OperandType.T.(*ast.BasicType).Kind.BitSize()
	if value.Kind() != constant.Int || constant.Sign(value) < 0 || !constant.Compare(value, gotoken.LSS, constant.MakeInt64(int64(width))) {
		return diagnostics.Diag{Title: "Invalid shift count", Message: fmt.Sprintf("shift count must be between 0 and %d", width-1), Span: b.Right.Span, Help: "Use a nonnegative count smaller than the left operand's bit width.", Related: []diagnostics.Related{{Span: b.Left.Span, Message: fmt.Sprintf("left operand is %d bits wide", width)}}}
	}
	return nil
}

func evaluateShift(left, right constant.Value, b *ast.BinExpr) (constant.Value, error) {
	if err := validateShiftCount(right, b); err != nil {
		return nil, err
	}
	count, _ := constant.Uint64Val(right)
	if b.Op == token.SHIFT_RIGHT {
		return constant.Shift(left, gotoken.SHR, uint(count)), nil
	}
	result := constant.Shift(left, gotoken.SHL, uint(count))
	kind := b.OperandType.T.(*ast.BasicType).Kind
	modulus := constant.Shift(constant.MakeInt64(1), gotoken.SHL, uint(kind.BitSize()))
	mask := constant.BinaryOp(modulus, gotoken.SUB, constant.MakeInt64(1))
	result = constant.BinaryOp(result, gotoken.AND, mask)
	unsigned := kind >= token.UINT_TYPE && kind <= token.U128_TYPE
	if !unsigned && !constant.Compare(result, gotoken.LSS, constant.Shift(constant.MakeInt64(1), gotoken.SHL, uint(kind.BitSize()-1))) {
		result = constant.BinaryOp(result, gotoken.SUB, modulus)
	}
	return result, nil
}
