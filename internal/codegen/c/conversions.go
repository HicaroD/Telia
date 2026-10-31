package c

import (
	"fmt"

	"github.com/HicaroD/Telia/internal/ast"
	"github.com/HicaroD/Telia/internal/lexer/token"
)

func unsignedInteger(kind token.Kind) bool {
	switch kind {
	case token.U8_TYPE,
		token.U16_TYPE,
		token.U32_TYPE,
		token.U64_TYPE,
		token.U128_TYPE,
		token.UINT_TYPE:
		return true
	}
	return false
}

func (c *CCodegen) emitNarrow(call *ast.FnCall) string {
	value := c.emitExpr(call.Args[0])
	source, result := c.nextTmp(), c.nextTmp()
	fmt.Fprintf(c.buf, "%s %s = %s;\n", emitCType(call.ConversionSource), source, value)
	fmt.Fprintf(c.buf, "%s %s = {0};\n", emitCType(call.ResultType), result)
	a := call.ConversionSource.T.(*ast.BasicType).Kind
	b := call.BuiltinType.T.(*ast.BasicType).Kind
	bits := b.BitSize()
	max := "(~(unsigned __int128)0)"
	if !unsignedInteger(b) {
		bits--
	}
	if bits < 128 {
		max = fmt.Sprintf("(((unsigned __int128)1 << %d) - 1)", bits)
	}
	overflow := fmt.Sprintf("(unsigned __int128)%s > %s", source, max)
	if !unsignedInteger(a) {
		if unsignedInteger(b) {
			overflow = fmt.Sprintf("%s < 0 || %s", source, overflow)
		} else {
			overflow = fmt.Sprintf("(__int128)%s < (-(__int128)%s - 1) || (%s >= 0 && (%s))", source, max, source, overflow)
		}
	}
	fmt.Fprintf(
		c.buf,
		"if (%s) { %s._1 = _TELIA_ERROR(ConversionOutOfRange, \"integer conversion out of range\"); } else { %s._0 = (%s)%s; }\n",
		overflow,
		result,
		result,
		emitCType(call.BuiltinType),
		source,
	)
	return result
}
