package c

import (
	"fmt"

	"github.com/HicaroD/Telia/internal/ast"
)

func (c *CCodegen) emitIndex(index *ast.IndexExpr) string {
	data, length := c.captureSequence(index.Value, index.Type)
	i := c.emitExpr(index.Index)
	tmp := c.nextTmp()
	fmt.Fprintf(c.buf, "uint64_t %s = %s;\n", tmp, i)
	fmt.Fprintf(
		c.buf,
		"if (%s >= %s) _panic(%s);\n",
		tmp,
		length,
		emitStringLiteral([]byte("index out of bounds")),
	)
	return fmt.Sprintf("(%s)[%s]", data, tmp)
}

func (c *CCodegen) emitSliceCall(call *ast.FnCall) string {
	data, length := c.captureSequence(call.Args[0], call.BuiltinType)
	if call.Builtin == "view" {
		return fmt.Sprintf("(%s){%s, %s}", emitCType(call.ResultType), data, length)
	}
	start, end := "0", length
	if call.Args[1] != nil {
		start = c.emitExpr(call.Args[1])
	}
	lo := c.nextTmp()
	c.buf.WriteString(fmt.Sprintf("uint64_t %s = %s;\n", lo, start))
	if call.Args[2] != nil {
		end = c.emitExpr(call.Args[2])
	}
	hi := c.nextTmp()
	c.buf.WriteString(fmt.Sprintf("uint64_t %s = %s;\n", hi, end))
	result := c.nextTmp()
	c.buf.WriteString(fmt.Sprintf("%s %s = {0};\n", emitCType(call.ResultType), result))
	c.buf.WriteString(
		fmt.Sprintf(
			"if (%s > %s || %s > %s) { %s._1 = _TELIA_ERROR(InvalidRange, \"slice range out of bounds\"); } else {\n",
			lo,
			hi,
			hi,
			length,
			result,
		),
	)
	c.buf.WriteString(
		fmt.Sprintf(
			"%s._0.data = %s == %s ? NULL : (%s + %s);\n%s._0.len = %s - %s;\n}\n",
			result,
			lo,
			hi,
			data,
			lo,
			result,
			hi,
			lo,
		),
	)
	return result
}

func (c *CCodegen) emitRangeLoop(stmt *ast.ForLoop, indent string) {
	if stmt.RangeIter != nil {
		c.emitContainerLoop(stmt, indent)
		return
	}
	data, length := c.captureSequence(stmt.Range, stmt.RangeType)
	index := c.nextTmp()
	seq := stmt.RangeType.T.(*ast.SequenceType)
	element := fmt.Sprintf("(%s)[%s]", data, index)
	c.buf.WriteString(
		fmt.Sprintf("for (uint64_t %s = 0; %s < %s; ++%s) {\n", index, index, length, index),
	)
	variable := &CVariable{
		CType: emitCType(seq.Element),
		Name:  c.nextTmp() + "_" + stmt.Bindings[0].Name.Name(),
	}
	stmt.Bindings[0].BackendType = variable
	c.buf.WriteString(fmt.Sprintf("%s %s = %s;\n", variable.CType, variable.Name, element))
	c.emitLoopBody(stmt.Block, indent+"    ")
	c.buf.WriteString(indent + "}\n")
}

func (c *CCodegen) emitContainerLoop(stmt *ast.ForLoop, indent string) {
	// Capture an address, never a second shallow cleanup owner. In particular,
	// indexed/indirect source expressions must run only once.
	value := c.emitExpr(stmt.Range)
	source, iterator, result := c.nextTmp(), c.nextTmp(), c.nextTmp()
	fmt.Fprintf(c.buf, "{\n%s *%s = &(%s);\n", emitCType(stmt.RangeType), source, value)
	fmt.Fprintf(
		c.buf,
		"%s %s = %s(%s);\n",
		emitCType(stmt.RangeIter.RetType),
		iterator,
		c.mangledFnDeclName(stmt.RangeIter),
		source,
	)
	c.buf.WriteString("while (true) {\n")
	fmt.Fprintf(
		c.buf,
		"%s %s = %s(&%s);\n",
		emitCType(stmt.RangeNext.RetType),
		result,
		c.mangledFnDeclName(stmt.RangeNext),
		iterator,
	)
	fmt.Fprintf(c.buf, "if (!%s._%d) break;\n", result, len(stmt.Bindings))
	for i, binding := range stmt.Bindings {
		variable := &CVariable{
			CType: emitCType(binding.Type),
			Name:  c.nextTmp() + "_" + binding.Name.Name(),
		}
		binding.BackendType = variable
		fmt.Fprintf(c.buf, "%s %s = %s._%d;\n", variable.CType, variable.Name, result, i)
	}
	c.emitLoopBody(stmt.Block, indent+"    ")
	c.buf.WriteString(indent + "}\n}\n")
}

// Preserve source evaluation order without copying an addressable array's
// storage. Slice handles and temporary array values can safely be copied.
func (c *CCodegen) captureSequence(node *ast.Node, ty *ast.ExprType) (data, length string) {
	value := c.emitExpr(node)
	seq := ty.T.(*ast.SequenceType)
	readOnly := false
	if deref, ok := node.Node.(*ast.DerefPointerExpr); ok {
		readOnly = deref.ReadOnly
	}
	if field, ok := node.Node.(*ast.FieldAccess); ok {
		for field.Right.Kind == ast.KIND_FIELD_ACCESS {
			field = field.Right.Node.(*ast.FieldAccess)
		}
		readOnly = field.ReadOnly
	}
	if seq.Foreign {
		base := c.nextTmp()
		element := emitCType(seq.Element)
		if readOnly {
			element = constCType(seq.Element)
		}
		fmt.Fprintf(c.buf, "%s *%s = %s;\n", element, base, value)
		return base, fmt.Sprint(seq.Length)
	}
	base := c.nextTmp()
	access := base + "."
	if ty.Kind == ast.EXPR_TYPE_ARRAY && node.IsAddressable() {
		ctype := emitCType(ty)
		if readOnly {
			ctype = constCType(ty)
		}
		c.buf.WriteString(fmt.Sprintf("%s *%s = &(%s);\n", ctype, base, value))
		access = base + "->"
	} else {
		c.buf.WriteString(fmt.Sprintf("%s %s = %s;\n", emitCType(ty), base, value))
	}
	if ty.Kind == ast.EXPR_TYPE_ARRAY {
		if seq.Length == 0 {
			return fmt.Sprintf("((%s *)NULL)", emitCType(seq.Element)), "0"
		}
		return access + "data", fmt.Sprint(seq.Length)
	}
	return access + "data", access + "len"
}
