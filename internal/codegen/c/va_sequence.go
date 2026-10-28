package c

import (
	"fmt"
	"strings"

	"github.com/HicaroD/Telia/internal/ast"
)

func (c *CCodegen) emitVaArrayInitializer(lit *ast.ArrayLiteralExpr) string {
	if len(lit.Values) == 0 {
		return "{0}"
	}
	var values []string
	for _, element := range lit.Values {
		value, tmp := c.emitExpr(element), c.nextTmp()
		fmt.Fprintf(c.buf, "%s %s = (%s);\n", vaTemporaryType(element), tmp, value)
		values = append(values, tmp)
	}
	return "{" + strings.Join(values, ", ") + "}"
}

// Native va_arg mutates shared cursor storage. C does not specify ordinary
// operand/argument order; sequence these expressions at their original site.
func hasVaEffects(node *ast.Node) bool {
	if node == nil {
		return false
	}
	switch n := node.Node.(type) {
	case *ast.FnCall:
		if n.Builtin == "va_arg" || n.Builtin == "va_end" || n.Builtin == "va_borrow" {
			return true
		}
		if hasVaEffects(n.CalleeExpr) {
			return true
		}
		for _, arg := range n.Args {
			if hasVaEffects(arg) {
				return true
			}
		}
	case *ast.NamespaceAccess:
		return hasVaEffects(n.Right)
	case *ast.BinExpr:
		return hasVaEffects(n.Left) || hasVaEffects(n.Right)
	case *ast.UnaryExpr:
		return hasVaEffects(n.Value)
	case *ast.AddressOfExpr:
		return hasVaEffects(n.Expr)
	case *ast.DerefPointerExpr:
		return hasVaEffects(n.Expr)
	case *ast.IndexExpr:
		return hasVaEffects(n.Value) || hasVaEffects(n.Index)
	case *ast.VarArgsExpr:
		for _, arg := range n.Args {
			if hasVaEffects(arg) {
				return true
			}
		}
	case *ast.ArrayLiteralExpr:
		for _, value := range n.Values {
			if hasVaEffects(value) {
				return true
			}
		}
	case *ast.StructLiteralExpr:
		for _, value := range n.Values {
			if hasVaEffects(value.Value) {
				return true
			}
		}
	case *ast.TupleExpr:
		for _, value := range n.Exprs {
			if hasVaEffects(value) {
				return true
			}
		}
	case *ast.MatchExpr:
		if hasVaEffects(n.Value) {
			return true
		}
		for _, arm := range n.Arms {
			if hasVaEffects(arm.Result) {
				return true
			}
			if arm.Block != nil {
				for _, stmt := range arm.Block.Statements {
					if hasVaEffects(stmt) {
						return true
					}
				}
			}
		}
	case *ast.VarStmt:
		return hasVaEffects(n.Expr)
	case *ast.YieldStmt:
		return hasVaEffects(n.Value)
	case *ast.ReturnStmt:
		return hasVaEffects(n.Value)
	case *ast.DeferStmt:
		return hasVaEffects(n.Stmt)
	case *ast.AssignmentStmt:
		for _, expr := range n.Values {
			if hasVaEffects(expr) {
				return true
			}
		}
	case *ast.CondStmt:
		if hasVaEffects(n.IfStmt.Expr) || hasVaEffects(&ast.Node{Node: n.IfStmt.Block}) {
			return true
		}
		for _, branch := range n.ElifStmts {
			if hasVaEffects(branch.Expr) || hasVaEffects(&ast.Node{Node: branch.Block}) {
				return true
			}
		}
		return n.ElseStmt != nil && hasVaEffects(&ast.Node{Node: n.ElseStmt.Block})
	case *ast.WhileLoop:
		return hasVaEffects(n.Cond) || hasVaEffects(&ast.Node{Node: n.Block})
	case *ast.ForLoop:
		return hasVaEffects(n.Init) || hasVaEffects(n.Cond) || hasVaEffects(n.Update) || hasVaEffects(n.Range) || hasVaEffects(&ast.Node{Node: n.Block})
	case *ast.BlockStmt:
		for _, stmt := range n.Statements {
			if hasVaEffects(stmt) {
				return true
			}
		}
	}
	return false
}

func (c *CCodegen) vaExpression(emit func() string) string {
	outer := c.buf
	c.buf = &strings.Builder{}
	defer func() { c.buf = outer }()
	result := emit()
	return "({\n" + c.buf.String() + result + ";\n})"
}

func (c *CCodegen) emitSequencedVaBinary(bin *ast.BinExpr) string {
	return c.vaExpression(func() string {
		left, right := c.nextTmp(), c.nextTmp()
		value := c.emitExpr(bin.Left)
		fmt.Fprintf(c.buf, "%s %s = (%s);\n", vaTemporaryType(bin.Left), left, value)
		value = c.emitExpr(bin.Right)
		fmt.Fprintf(c.buf, "%s %s = (%s);\n", vaTemporaryType(bin.Right), right, value)
		return emitBinaryOperation(bin, left, right)
	})
}

func (c *CCodegen) emitSequencedVaCall(call *ast.FnCall) string {
	return c.vaExpression(func() string {
		name := c.getFnCallName(call)
		if call.Callee != nil || call.CalleeExpr != nil {
			tmp := c.nextTmp()
			fmt.Fprintf(c.buf, "__auto_type %s = (%s);\n", tmp, name)
			name = tmp
		}
		var arguments []*ast.Node
		for _, arg := range call.Args {
			if pack, ok := arg.Node.(*ast.VarArgsExpr); ok {
				arguments = append(arguments, pack.Args...)
			} else {
				arguments = append(arguments, arg)
			}
		}
		var values []string
		for _, arg := range arguments {
			if borrow := vaCall(arg); borrow != nil && borrow.Builtin == "va_borrow" {
				owner := borrow.VaOwner.BackendType.(*CVariable).Name
				fmt.Fprintf(c.buf, "if (%s.state != 1) _panic(%s);\n%s.state = 2;\n", owner, vaPanicMessage(borrow, "va_list is ended or consumed"), owner)
				values = append(values, owner+".args")
				continue
			}
			if id, ok := arg.Node.(*ast.IdExpr); ok && id.N != nil && id.N.Kind == ast.KIND_PARAM && id.N.Node.(*ast.Param).Type.Kind == ast.EXPR_TYPE_C_VA_LIST {
				values = append(values, c.emitExpr(arg)) // Native array adjustment, no cursor assignment.
				continue
			}
			value, tmp := c.emitExpr(arg), c.nextTmp()
			fmt.Fprintf(c.buf, "%s %s = (%s);\n", vaTemporaryType(arg), tmp, value)
			values = append(values, tmp)
		}
		return fmt.Sprintf("%s(%s)", name, strings.Join(values, ", "))
	})
}

func vaTemporaryType(node *ast.Node) string {
	if value, ok := node.Node.(*ast.NullPtrExpr); ok && value.Type != nil {
		return emitCType(value.Type)
	}
	return "__auto_type"
}
