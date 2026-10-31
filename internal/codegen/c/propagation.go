package c

import (
	"fmt"
	"strings"

	"github.com/HicaroD/Telia/internal/ast"
)

func (c *CCodegen) emitPropagatedValue(call *ast.FnCall) string {
	// Keep the call and propagation at the expression's evaluation site. In
	// particular, a short-circuited operand must not call or run cleanup.
	return c.vaExpression(func() string {
		value := c.emitFnCall(call)
		tmp := c.nextTmp()
		fmt.Fprintf(c.buf, "%s %s = %s;\n", emitCType(call.ResultType), tmp, value)
		types := call.ResultType.T.(*ast.TupleType).Types
		last := len(types) - 1
		op := call.AtOp.Op.(*ast.CatchAtOperator)
		fmt.Fprintf(
			c.buf,
			"if (%s._%d.present) {\n_Error %s = %s._%d;\n",
			tmp,
			last,
			op.ErrVarName.Name(),
			tmp,
			last,
		)
		c.emitBlock(op.Block, "")
		c.buf.WriteString("}\n")
		if last == 1 {
			return tmp + "._0"
		}
		var values []string
		for i := 0; i < last; i++ {
			values = append(values, fmt.Sprintf("%s._%d", tmp, i))
		}
		success := &ast.ExprType{Kind: ast.EXPR_TYPE_TUPLE, T: &ast.TupleType{Types: types[:last]}}
		return fmt.Sprintf("(%s){%s}", emitCType(success), strings.Join(values, ", "))
	})
}
