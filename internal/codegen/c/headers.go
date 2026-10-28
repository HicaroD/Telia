package c

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/HicaroD/Telia/internal/ast"
	"github.com/HicaroD/Telia/internal/diagnostics"
)

func foreignFieldMarker(d *ast.StructDecl, field *ast.StructField) string {
	return "foreign field mismatch: " + d.CName + "." + field.Name.Name()
}

func foreignFieldType(ty *ast.ExprType) string {
	if ty.Kind == ast.EXPR_TYPE_ARRAY {
		seq := ty.T.(*ast.SequenceType)
		return fmt.Sprintf("%s[%d]", foreignFieldType(seq.Element), seq.Length)
	}
	return emitCType(ty)
}

func (c *CCodegen) emitForeignFieldChecks() {
	for _, d := range c.program.Enums {
		if d.Foreign == nil || d.Type == nil {
			continue
		}
		name, repr := d.CName, emitCType(d.Type)
		fmt.Fprintf(
			c.buf,
			"_Static_assert(sizeof(%s) == sizeof(%s) && (((%s)-1 > 0) == ((%s)-1 > 0)), %q);\n",
			name,
			repr,
			name,
			repr,
			"foreign enum representation mismatch: "+name,
		)
	}
	for _, d := range c.program.Structs {
		if d.Foreign == nil {
			continue
		}
		for _, field := range d.Fields {
			fmt.Fprintf(
				c.buf,
				"_Static_assert(__builtin_types_compatible_p(__typeof__(((%s *)0)->%s), %s), %q);\n",
				d.CName,
				field.Name.Name(),
				foreignFieldType(field.Type),
				foreignFieldMarker(d, field),
			)
		}
	}
}

func (c *CCodegen) foreignFieldDiagnostic(output string) error {
	for _, d := range c.program.Enums {
		if d.Foreign != nil &&
			strings.Contains(output, "foreign enum representation mismatch: "+d.CName) {
			return diagnostics.Diag{
				Title:   "Foreign enum representation mismatch",
				Message: "the representation annotation does not match the header enum's size and signedness",
				Span:    d.Name.Span,
				Help:    "Remove the annotation to use the native C enum, or match the header's exact representation.",
			}
		}
	}
	for _, d := range c.program.Structs {
		if d.Foreign == nil {
			continue
		}
		for _, field := range d.Fields {
			if strings.Contains(output, foreignFieldMarker(d, field)) {
				return diagnostics.Diag{
					Title: "Foreign field type mismatch",
					Message: fmt.Sprintf(
						"field '%s' is not compatible with its declaration in the C header",
						field.Name.Name(),
					),
					Span: field.Name.Span,
					Help: "Match the header's exact field type, pointer qualifiers, and array dimensions; do not substitute a Telia wrapper array or stand-in record.",
					Related: []diagnostics.Related{
						{Span: d.Name.Span, Message: "header-owned record declared here"},
					},
				}
			}
		}
	}
	return nil
}

func (c *CCodegen) emitHeaderIncludes(
	pkg *ast.Package,
	packages map[*ast.Package]bool,
	headers map[string]bool,
) {
	if pkg == nil || packages[pkg] {
		return
	}
	packages[pkg] = true
	for _, file := range pkg.Files {
		for _, imp := range file.Imports {
			c.emitHeaderIncludes(imp.Package, packages, headers)
		}
		for _, node := range file.Body {
			if ext, ok := node.Node.(*ast.ExternDecl); ok && ext.Attributes.Header != "" &&
				!headers[ext.Attributes.Header] {
				headers[ext.Attributes.Header] = true
				fmt.Fprintf(c.buf, "#include %s\n", strconv.Quote(ext.Attributes.Header))
			}
		}
	}
}
