package c

import (
	"fmt"
	"strings"

	"github.com/HicaroD/Telia/internal/ast"
	"github.com/HicaroD/Telia/internal/lexer/token"
)

// Keep captures at the expression site (including loop conditions and short
// circuit operands). C argument evaluation order is not source order.
func (c *CCodegen) emitKeyHash(call *ast.FnCall) string {
	outer := c.buf
	c.buf = &strings.Builder{}
	defer func() { c.buf = outer }()
	args := make([]string, 3)
	for i, arg := range call.Args {
		value := c.emitExpr(arg)
		ty := "uint64_t"
		if i == 0 {
			ty = emitCType(call.BuiltinType)
		}
		args[i] = c.nextTmp()
		fmt.Fprintf(c.buf, "%s %s = %s;\n", ty, args[i], value)
	}
	state := c.nextTmp()
	fmt.Fprintf(c.buf, "_TeliaSip %s = _telia_sip_init(%s,%s);\n", state, args[1], args[2])
	c.hashValue(call.BuiltinType, args[0], state)
	fmt.Fprintf(c.buf, "_telia_sip_finish(%s);\n", state)
	return "({\n" + c.buf.String() + "})"
}

func (c *CCodegen) hashValue(ty *ast.ExprType, value, state string) {
	switch ty.Kind {
	case ast.EXPR_TYPE_STRUCT:
		for _, f := range ty.T.(*ast.StructType).Decl.Fields {
			c.hashValue(f.Type, value+"."+f.Name.Name(), state)
		}
	case ast.EXPR_TYPE_TUPLE:
		for i, t := range ty.T.(*ast.TupleType).Types {
			c.hashValue(t, fmt.Sprintf("%s._%d", value, i), state)
		}
	case ast.EXPR_TYPE_ARRAY:
		seq := ty.T.(*ast.SequenceType)
		i := c.nextTmp()
		fmt.Fprintf(c.buf, "for(uint64_t %s=0; %s<%d; ++%s) {\n", i, i, seq.Length, i)
		// Zero-length arrays use a dummy byte, not an element array.
		if seq.Length > 0 {
			c.hashValue(seq.Element, value+".data["+i+"]", state)
		}
		c.buf.WriteString("}\n")
	case ast.EXPR_TYPE_ENUM:
		d := ty.T.(*ast.EnumType).Decl
		if !d.HasPayload {
			c.hashInteger(value, "sizeof("+value+")", state)
			return
		}
		c.hashValue(d.Type, value+".tag", state)
		fmt.Fprintf(c.buf, "switch (%s.tag) {\n", value)
		for _, v := range d.Variants {
			fmt.Fprintf(c.buf, "case %s: {\n", c.emitConstant(v.Value))
			for _, f := range v.Fields {
				c.hashValue(
					f.Type,
					fmt.Sprintf("%s.payload.v%d.%s", value, v.Index, f.Name.Name()),
					state,
				)
			}
			c.buf.WriteString("break; }\n")
		}
		// Undeclared tags have no active payload (the equality contract).
		c.buf.WriteString("default: break; }\n")
	case ast.EXPR_TYPE_BASIC:
		b := ty.T.(*ast.BasicType)
		switch b.Kind {
		case token.STRING_TYPE:
			c.hashInteger(value+".len", "8", state)
			fmt.Fprintf(c.buf, "_telia_sip_write(&%s,%s.data,%s.len);\n", state, value, value)
		case token.RAWPTR_TYPE, token.CSTRING_TYPE:
			c.hashInteger("(uintptr_t)"+value, "sizeof(uintptr_t)", state)
		case token.F32_TYPE, token.FLOAT_TYPE, token.F64_TYPE:
			tmp, bits := c.nextTmp(), c.nextTmp()
			integer := "uint32_t"
			if b.Kind == token.F64_TYPE {
				integer = "uint64_t"
			}
			fmt.Fprintf(
				c.buf,
				"%s %s = %s == 0 ? 0 : %s; %s %s; memcpy(&%s,&%s,sizeof(%s));\n",
				emitCType(ty),
				tmp,
				value,
				value,
				integer,
				bits,
				bits,
				tmp,
				bits,
			)
			c.hashInteger(bits, "sizeof("+bits+")", state)
		default:
			c.hashInteger(value, "sizeof("+value+")", state)
		}
	case ast.EXPR_TYPE_POINTER, ast.EXPR_TYPE_ERROR_KIND:
		c.hashInteger("(uintptr_t)"+value, "sizeof(uintptr_t)", state)
	default:
		panic("hash eligibility must be checked before code generation")
	}
}

// Explicit unsigned conversion, shifts and byte order: no signed overflow,
// padding, alignment requirement, or host-endian object representation.
func (c *CCodegen) hashInteger(value, width, state string) {
	tmp, i, b := c.nextTmp(), c.nextTmp(), c.nextTmp()
	fmt.Fprintf(
		c.buf,
		"{ unsigned __int128 %s = (unsigned __int128)(%s); for(unsigned %s=0; %s<%s; ++%s) { unsigned char %s = (unsigned char)(%s >> (8*%s)); _telia_sip_write(&%s,&%s,1); } }\n",
		tmp,
		value,
		i,
		i,
		width,
		i,
		b,
		tmp,
		i,
		state,
		b,
	)
}
