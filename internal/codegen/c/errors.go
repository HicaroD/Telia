package c

import (
	"fmt"
	"path/filepath"
	"strings"

	"github.com/HicaroD/Telia/config"
	"github.com/HicaroD/Telia/internal/ast"
)

func errorKind(d *ast.EnumDecl, v *ast.EnumVariant) string {
	// The standard runtime and Telia declarations share these private symbols.
	// Check the canonical package, not a user package or enum with the same name.
	if config.ENVS != nil && d.Name.Name() == "Error" &&
		filepath.Clean(d.Owner.Loc.Path) == filepath.Join(config.ENVS.STD, "errors") {
		return "&_telia_runtime_kind_" + v.Name.Name()
	}
	return fmt.Sprintf("&_telia_kind_%s_%d", d.CName, v.Index)
}

func (c *CCodegen) emitErrorType() {
	c.buf.WriteString(
		"struct _Error { _String msg; bool present; const void *kind; union { unsigned char empty;\n",
	)
	for _, d := range c.program.Enums {
		if d.IsError && d.HasPayload {
			fmt.Fprintf(c.buf, "%s %s;\n", d.CName, d.CName)
		}
	}
	c.buf.WriteString("} payload; };\n")
	for _, name := range []string{"OutOfBounds", "InvalidRange", "SizeOverflow", "AllocationFailed", "InvalidAlignment", "EmbeddedZero", "EmptySeparator", "InvalidUTF8"} {
		fmt.Fprintf(c.buf, "static const char _telia_runtime_kind_%s;\n", name)
	}
	for _, d := range c.program.Enums {
		if d.IsError {
			for _, v := range d.Variants {
				fmt.Fprintf(c.buf, "static const char _telia_kind_%s_%d;\n", d.CName, v.Index)
			}
		}
	}
}

func (c *CCodegen) emitError(e *ast.ErrorExpr) string {
	// A statement expression preserves source order and the evaluation site,
	// including short-circuit operands and loop conditions.
	outer := c.buf
	c.buf = &strings.Builder{}
	defer func() { c.buf = outer }()
	identity, payload := "NULL", ""
	if e.Enum != nil {
		value := c.emitExpr(e.Args[0])
		tmp := c.nextTmp()
		ty := &ast.ExprType{Kind: ast.EXPR_TYPE_ENUM, T: &ast.EnumType{Decl: e.Enum}}
		fmt.Fprintf(c.buf, "%s %s = %s;\n", emitCType(ty), tmp, value)
		tag := tmp
		if e.Enum.HasPayload {
			tag += ".tag"
			payload = fmt.Sprintf(", .payload.%s = %s", e.Enum.CName, tmp)
		}
		for _, v := range e.Enum.Variants {
			identity = fmt.Sprintf(
				"(%s == %s ? %s : %s)",
				tag,
				c.emitConstant(v.Value),
				errorKind(e.Enum, v),
				identity,
			)
		}
	}
	message := c.emitExpr(e.Args[len(e.Args)-1])
	fmt.Fprintf(
		c.buf,
		"(_Error){.msg = %s, .present = true, .kind = %s%s};\n",
		message,
		identity,
		payload,
	)
	return "({\n" + c.buf.String() + "})"
}
