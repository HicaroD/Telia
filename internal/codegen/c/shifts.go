package c

import (
	"fmt"
	"math/big"
	"strings"

	"github.com/HicaroD/Telia/internal/ast"
	"github.com/HicaroD/Telia/internal/lexer/token"
)

// Shift unsigned bits, then reconstruct signed results without out-of-range
// unsigned-to-signed conversions or C's undefined signed-left-shift behavior.
func (c *CCodegen) emitShift(b *ast.BinExpr) string {
	outer := c.buf
	c.buf = &strings.Builder{}
	defer func() { c.buf = outer }()
	kind := b.OperandType.T.(*ast.BasicType).Kind
	width := kind.BitSize()
	unsigned := kind >= token.UINT_TYPE && kind <= token.U128_TYPE
	unsignedType := ast.NewBasicType(token.U8_TYPE)
	switch width {
	case 16:
		unsignedType = ast.NewBasicType(token.U16_TYPE)
	case 32:
		unsignedType = ast.NewBasicType(token.U32_TYPE)
	case 64:
		unsignedType = ast.NewBasicType(token.U64_TYPE)
	case 128:
		unsignedType = ast.NewBasicType(token.U128_TYPE)
	}
	left, count, bits := c.nextTmp(), c.nextTmp(), c.nextTmp()
	value := c.emitExpr(b.Left)
	fmt.Fprintf(c.buf, "%s %s = %s;\n", emitCType(b.OperandType), left, value)
	value = c.emitExpr(b.Right)
	fmt.Fprintf(c.buf, "%s %s = %s;\n", emitCType(b.RightType), count, value)
	fmt.Fprintf(c.buf, "if (%s < 0 || (unsigned __int128)%s >= %d) _panic(%s);\n", count, count, width, emitStringLiteral([]byte("shift count out of range")))
	op := "<<"
	if b.Op == token.SHIFT_RIGHT {
		op = ">>"
	}
	fmt.Fprintf(c.buf, "%s %s = (%s)((unsigned __int128)(%s)%s %s (unsigned)%s);\n", emitCType(unsignedType), bits, emitCType(unsignedType), emitCType(unsignedType), left, op, count)
	if !unsigned && b.Op == token.SHIFT_RIGHT {
		fmt.Fprintf(c.buf, "if (%s < 0 && %s != 0) %s |= (%s)(~(unsigned __int128)0 << (%d - (unsigned)%s));\n", left, count, bits, emitCType(unsignedType), width, count)
	}
	result := bits
	if !unsigned {
		high := new(big.Int).Lsh(big.NewInt(1), uint(width-1))
		mask := new(big.Int).Sub(new(big.Int).Set(high), big.NewInt(1))
		minimum := new(big.Int).Neg(new(big.Int).Set(high))
		emit := func(value *big.Int, ty *ast.ExprType) string {
			return c.emitConstant(&ast.LiteralExpr{Type: ty, Value: []byte(value.String())})
		}
		result = fmt.Sprintf("((%s)(%s & %s) + ((%s & %s) ? %s : 0))", emitCType(b.OperandType), bits, emit(mask, unsignedType), bits, emit(high, unsignedType), emit(minimum, b.OperandType))
	}
	fmt.Fprintf(c.buf, "(%s)(%s);\n", emitCType(b.OperandType), result)
	return "({\n" + c.buf.String() + "})"
}
