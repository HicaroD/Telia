package c

import (
	"fmt"
	"strconv"

	"github.com/HicaroD/Telia/internal/ast"
)

func (c *CCodegen) emitHeaderIncludes(pkg *ast.Package, packages map[*ast.Package]bool, headers map[string]bool) {
	if pkg == nil || packages[pkg] {
		return
	}
	packages[pkg] = true
	for _, file := range pkg.Files {
		for _, imp := range file.Imports {
			c.emitHeaderIncludes(imp.Package, packages, headers)
		}
		for _, node := range file.Body {
			if ext, ok := node.Node.(*ast.ExternDecl); ok && ext.Attributes.Header != "" && !headers[ext.Attributes.Header] {
				headers[ext.Attributes.Header] = true
				fmt.Fprintf(c.buf, "#include %s\n", strconv.Quote(ext.Attributes.Header))
			}
		}
	}
}
