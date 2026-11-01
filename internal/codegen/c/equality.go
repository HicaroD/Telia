package c

import (
	"fmt"
	"strings"

	"github.com/HicaroD/Telia/internal/ast"
	"github.com/HicaroD/Telia/internal/lexer/token"
)

// Capture both operands at the expression site, in source order. Recursive
// comparisons see only these pure local values, never reevaluated expressions.
func (c *CCodegen) emitEquality(bin *ast.BinExpr) string {
	outer := c.buf
	c.buf = &strings.Builder{}
	defer func() { c.buf = outer }()
	left := c.emitExpr(bin.Left)
	a, b := c.nextTmp(), c.nextTmp()
	fmt.Fprintf(c.buf, "%s %s = %s;\n", emitCType(bin.OperandType), a, left)
	right := c.emitExpr(bin.Right)
	fmt.Fprintf(c.buf, "%s %s = %s;\n", emitCType(bin.OperandType), b, right)
	equal := c.equalValues(bin.OperandType, a, b)
	if bin.Op == token.BANG_EQUAL {
		equal = "!(" + equal + ")"
	}
	fmt.Fprintf(c.buf, "%s;\n", equal)
	return "({\n" + c.buf.String() + "})"
}

func (c *CCodegen) equalValues(ty *ast.ExprType, a, b string) string {
	var parts []string
	switch ty.Kind {
	case ast.EXPR_TYPE_BASIC:
		if ty.T.(*ast.BasicType).Kind == token.STRING_TYPE {
			return fmt.Sprintf("_telia_string_equal(%s, %s)", a, b)
		}
	case ast.EXPR_TYPE_STRUCT:
		for _, f := range ty.T.(*ast.StructType).Decl.Fields {
			parts = append(parts, c.equalValues(f.Type, a+"."+f.Name.Name(), b+"."+f.Name.Name()))
		}
		return conjunction(parts)
	case ast.EXPR_TYPE_TUPLE:
		for i, t := range ty.T.(*ast.TupleType).Types {
			field := fmt.Sprintf("._%d", i)
			parts = append(parts, c.equalValues(t, a+field, b+field))
		}
		return conjunction(parts)
	case ast.EXPR_TYPE_ARRAY:
		seq := ty.T.(*ast.SequenceType)
		i, result := c.nextTmp(), c.nextTmp()
		field := ".data[" + i + "]"
		return fmt.Sprintf(
			"({ _Bool %s = 1; for (uint64_t %s = 0; %s < %d && %s; ++%s) { %s = %s; } %s; })",
			result,
			i,
			i,
			seq.Length,
			result,
			i,
			result,
			c.equalValues(seq.Element, a+field, b+field),
			result,
		)
	case ast.EXPR_TYPE_ENUM:
		d := ty.T.(*ast.EnumType).Decl
		if d.HasPayload {
			// An undeclared tag has no active payload. In particular, zero<T>()
			// can produce tag zero when every declared variant tag is nonzero.
			payloadEqual := "1"
			for i := len(d.Variants) - 1; i >= 0; i-- {
				v := d.Variants[i]
				var fields []string
				for _, f := range v.Fields {
					field := fmt.Sprintf(".payload.v%d.%s", v.Index, f.Name.Name())
					fields = append(fields, c.equalValues(f.Type, a+field, b+field))
				}
				payloadEqual = fmt.Sprintf(
					"(%s.tag == %s ? %s : %s)",
					a,
					c.emitConstant(v.Value),
					conjunction(fields),
					payloadEqual,
				)
			}
			return fmt.Sprintf("((%s.tag == %s.tag) && %s)", a, b, payloadEqual)
		}
	}
	return fmt.Sprintf("(%s == %s)", a, b)
}

func conjunction(parts []string) string {
	if len(parts) == 0 {
		return "1"
	}
	return "(" + strings.Join(parts, " && ") + ")"
}
