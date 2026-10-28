package c

import (
	"fmt"
	"math/big"

	"github.com/HicaroD/Telia/internal/ast"
)

func (c *CCodegen) emitConstant(lit *ast.LiteralExpr) string {
	if !lit.Type.IsInteger() {
		return c.emitLiteral(lit)
	}
	value, ok := new(big.Int).SetString(string(lit.Value), 10)
	if !ok {
		panic("integer constant must have been evaluated before code generation")
	}
	negative := value.Sign() < 0
	if negative {
		// Emit -(magnitude - 1) - 1 so even the signed minimum never
		// requires an out-of-range positive signed intermediate.
		value.Neg(value)
		value.Sub(value, big.NewInt(1))
	}
	text := fmt.Sprintf("UINT64_C(%d)", value.Uint64())
	if value.BitLen() > 64 {
		high := new(big.Int).Rsh(new(big.Int).Set(value), 64)
		text = fmt.Sprintf("(((unsigned __int128)UINT64_C(%d) << 64) | %s)", high.Uint64(), text)
	}
	ctype := emitCType(lit.Type)
	if negative {
		return fmt.Sprintf("((%s)(-((%s)(%s)) - 1))", ctype, ctype, text)
	}
	return fmt.Sprintf("((%s)(%s))", ctype, text)
}
