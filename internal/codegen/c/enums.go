package c

import (
	"fmt"
	"strings"

	"github.com/HicaroD/Telia/internal/ast"
)

type matchScope struct {
	result, label string
	deferDepth    int
}

func (c *CCodegen) emitEnumDecl(d *ast.EnumDecl) {
	if !d.HasPayload {
		return
	}
	fmt.Fprintf(c.buf, "struct %s { %s tag; union {\n", d.CName, emitCType(d.Type))
	for _, v := range d.Variants {
		if len(v.Fields) == 0 {
			continue
		}
		c.buf.WriteString("struct {\n")
		for _, field := range v.Fields {
			fmt.Fprintf(c.buf, "%s %s;\n", emitCType(field.Type), field.Name.Name())
		}
		fmt.Fprintf(c.buf, "} v%d;\n", v.Index)
	}
	// An enum can declare an empty payload block. C unions cannot be empty.
	c.buf.WriteString("unsigned char empty;\n} payload; };\n")
}

func (c *CCodegen) emitEnumConstruction(ns *ast.NamespaceAccess) string {
	d, v := ns.EnumDecl, ns.EnumVariant
	if !d.HasPayload {
		return c.emitConstant(v.Value)
	}
	var fields []string
	if ns.Right.Kind == ast.KIND_STRUCT_EXPR {
		for _, value := range ns.Right.Node.(*ast.StructLiteralExpr).Values {
			// Materialize in source order. C initializer evaluation order is unspecified.
			var ty *ast.ExprType
			for _, field := range v.Fields {
				if field.Name.Name() == value.Name.Name() {
					ty = field.Type
					break
				}
			}
			expr := c.emitExpr(value.Value)
			tmp := c.nextTmp()
			fmt.Fprintf(c.buf, "%s %s = %s;\n", emitCType(ty), tmp, expr)
			fields = append(fields, fmt.Sprintf(".%s = %s", value.Name.Name(), tmp))
		}
	}
	payload := ""
	if len(fields) != 0 {
		payload = fmt.Sprintf(", .payload.v%d = {%s}", v.Index, strings.Join(fields, ", "))
	}
	return fmt.Sprintf("(%s){.tag = %s%s}", d.CName, c.emitConstant(v.Value), payload)
}

func (c *CCodegen) emitMatch(m *ast.MatchExpr) string {
	// Keep the entire match at its expression evaluation site. Hoisting its
	// statements would break short-circuit operators and loop conditions.
	outer := c.buf
	c.buf = &strings.Builder{}
	defer func() { c.buf = outer }()
	value := c.emitExpr(m.Value)
	scrutinee, result := c.nextTmp(), c.nextTmp()
	label := c.nextTmp()
	fmt.Fprintf(c.buf, "%s %s = %s;\n%s %s;\n", emitCType(&ast.ExprType{Kind: ast.EXPR_TYPE_ENUM, T: &ast.EnumType{Decl: m.Enum}}), scrutinee, value, emitCType(m.Type), result)
	tag := scrutinee
	if m.Enum.HasPayload {
		tag += ".tag"
	}
	for i, arm := range m.Arms {
		if i != 0 {
			c.buf.WriteString("else ")
		}
		if !arm.Fallback {
			fmt.Fprintf(c.buf, "if (%s == %s) ", tag, c.emitConstant(arm.Variant.Value))
		}
		c.buf.WriteString("{\n")
		for _, field := range arm.Fields {
			name := c.nextTmp()
			field.Binding.BackendType = &CVariable{Name: name, CType: emitCType(field.Binding.Type)}
			fmt.Fprintf(c.buf, "%s %s = %s.payload.v%d.%s;\n", emitCType(field.Binding.Type), name, scrutinee, arm.Variant.Index, field.Name.Name())
		}
		if arm.Block != nil {
			c.matchScopes = append(c.matchScopes, matchScope{result: result, label: label, deferDepth: len(c.deferScopes)})
			c.emitBlock(arm.Block, "")
			c.matchScopes = c.matchScopes[:len(c.matchScopes)-1]
		} else {
			value := c.emitExpr(arm.Result)
			fmt.Fprintf(c.buf, "%s = %s;\n", result, value)
		}
		c.buf.WriteString("}\n")
	}
	fmt.Fprintf(c.buf, "%s: ;\n%s;\n", label, result)
	return "({\n" + c.buf.String() + "})"
}
