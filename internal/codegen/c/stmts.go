package c

import (
	"fmt"

	"github.com/HicaroD/Telia/internal/ast"
)

// emitBlock emits each statement in a block at the given indentation level.
// Deferred statements are emitted in LIFO order before each return.
func (c *CCodegen) emitBlock(block *ast.BlockStmt, indent string) {
	for _, stmt := range block.Statements {
		// Before a return, flush the defer stack in reverse order.
		if stmt.Kind == ast.KIND_RETURN_STMT {
			for i := len(block.DeferStack) - 1; i >= 0; i-- {
				d := block.DeferStack[i]
				if !d.Skip {
					c.buf.WriteString(indent)
					c.emitStmtNode(d.Stmt, indent)
					c.buf.WriteString(";\n")
				}
			}
		}
		c.buf.WriteString(indent)
		c.emitStmtNode(stmt, indent)
	}
}

// emitStmtNode dispatches to the correct statement emitter.
func (c *CCodegen) emitStmtNode(node *ast.Node, indent string) {
	switch node.Kind {
	case ast.KIND_RETURN_STMT:
		c.emitReturn(node.Node.(*ast.ReturnStmt), indent)
	case ast.KIND_VAR_STMT:
		c.emitVarStmt(node.Node.(*ast.VarStmt), indent)
	case ast.KIND_ASSIGNMENT_STMT:
		c.emitAssignment(node.Node.(*ast.AssignmentStmt), indent)
	case ast.KIND_COND_STMT:
		c.emitCond(node.Node.(*ast.CondStmt), indent)
	case ast.KIND_FOR_LOOP_STMT:
		c.emitForLoop(node.Node.(*ast.ForLoop), indent)
	case ast.KIND_WHILE_LOOP_STMT:
		c.emitWhileLoop(node.Node.(*ast.WhileLoop), indent)
	case ast.KIND_DEFER_STMT:
		// Defer statements are emitted before returns; skip here.
	case ast.KIND_FN_CALL:
		call := node.Node.(*ast.FnCall)
		c.buf.WriteString(c.emitFnCall(call))
		c.buf.WriteString(";\n")
	case ast.KIND_NAMESPACE_ACCESS:
		ns := node.Node.(*ast.NamespaceAccess)
		c.buf.WriteString(c.emitNamespaceAccess(ns))
		c.buf.WriteString(";\n")
	case ast.KIND_FIELD_ACCESS:
		c.buf.WriteString(c.emitExpr(node))
		c.buf.WriteString(";\n")
	default:
		panic(fmt.Sprintf("emitStmtNode: unhandled statement kind: %v", node.Kind))
	}
}

// emitReturn emits a return statement.
func (c *CCodegen) emitReturn(ret *ast.ReturnStmt, _ string) {
	if ret.Value == nil || ret.Value.Kind == ast.KIND_VOID_EXPR {
		c.buf.WriteString("return;\n")
		return
	}
	c.buf.WriteString(fmt.Sprintf("return %s;\n", c.emitExpr(ret.Value)))
}

// emitVarStmt emits a variable declaration or reassignment.
func (c *CCodegen) emitVarStmt(stmt *ast.VarStmt, indent string) {
	if stmt.IsDecl {
		// Declaration: may be a single var or a multi-assignment from a tuple.
		// For now handle the common single-name case.
		if len(stmt.Names) == 1 {
			varId := stmt.Names[0].Node.(*ast.VarIdStmt)
			cType := emitCType(varId.Type)
			cVar := &CVariable{CType: cType, Name: varId.Name.Name()}
			varId.BackendType = cVar
			val := c.emitExpr(stmt.Expr)
			c.buf.WriteString(fmt.Sprintf("%s %s = %s;\n", cType, varId.Name.Name(), val))
		} else {
			// Multi-variable declaration (tuple unpack) — deferred to #73.
			// Emit each name as a separate declaration for now.
			val := c.emitExpr(stmt.Expr)
			tmp := c.nextTmp()
			// We don't know the tuple type yet — emit a comment placeholder.
			c.buf.WriteString(fmt.Sprintf("/* tuple unpack: %s = %s */\n", tmp, val))
		}
	} else {
		// Reassignment
		if len(stmt.Names) == 1 {
			target := c.emitVarTarget(stmt.Names[0])
			val := c.emitExpr(stmt.Expr)
			c.buf.WriteString(fmt.Sprintf("%s = %s;\n", target, val))
		}
	}
}

// emitAssignment emits an assignment statement (targets = values).
func (c *CCodegen) emitAssignment(stmt *ast.AssignmentStmt, indent string) {
	if stmt.Decl {
		for i, target := range stmt.Targets {
			varId := target.Node.(*ast.VarIdStmt)
			cType := emitCType(varId.Type)
			cVar := &CVariable{CType: cType, Name: varId.Name.Name()}
			varId.BackendType = cVar
			val := c.emitExpr(stmt.Values[i])
			c.buf.WriteString(fmt.Sprintf("%s %s = %s;\n", cType, varId.Name.Name(), val))
			if i < len(stmt.Targets)-1 {
				c.buf.WriteString(indent)
			}
		}
	} else {
		for i, target := range stmt.Targets {
			t := c.emitExpr(target)
			val := c.emitExpr(stmt.Values[i])
			c.buf.WriteString(fmt.Sprintf("%s = %s;\n", t, val))
			if i < len(stmt.Targets)-1 {
				c.buf.WriteString(indent)
			}
		}
	}
}

// emitCond emits an if/elif/else chain.
func (c *CCodegen) emitCond(stmt *ast.CondStmt, indent string) {
	inner := indent + "    "

	cond := c.emitExpr(stmt.IfStmt.Expr)
	c.buf.WriteString(fmt.Sprintf("if (%s) {\n", cond))
	c.emitBlock(stmt.IfStmt.Block, inner)
	c.buf.WriteString(indent + "}")

	for _, elif := range stmt.ElifStmts {
		cond = c.emitExpr(elif.Expr)
		c.buf.WriteString(fmt.Sprintf(" else if (%s) {\n", cond))
		c.emitBlock(elif.Block, inner)
		c.buf.WriteString(indent + "}")
	}

	if stmt.ElseStmt != nil {
		c.buf.WriteString(" else {\n")
		c.emitBlock(stmt.ElseStmt.Block, inner)
		c.buf.WriteString(indent + "}")
	}

	c.buf.WriteString("\n")
}

// emitForLoop emits a for loop.
func (c *CCodegen) emitForLoop(stmt *ast.ForLoop, indent string) {
	inner := indent + "    "

	// Emit init inline without trailing newline/semicolon so it can go into the for(;;)
	init := c.emitForInit(stmt.Init)
	cond := c.emitExpr(stmt.Cond)
	update := c.emitForUpdate(stmt.Update)

	c.buf.WriteString(fmt.Sprintf("for (%s; %s; %s) {\n", init, cond, update))
	c.emitBlock(stmt.Block, inner)
	c.buf.WriteString(indent + "}\n")
}

// emitForInit extracts the init expression/declaration for a for loop header.
func (c *CCodegen) emitForInit(node *ast.Node) string {
	if node == nil {
		return ""
	}
	switch node.Kind {
	case ast.KIND_VAR_STMT:
		stmt := node.Node.(*ast.VarStmt)
		if stmt.IsDecl && len(stmt.Names) == 1 {
			varId := stmt.Names[0].Node.(*ast.VarIdStmt)
			cType := emitCType(varId.Type)
			cVar := &CVariable{CType: cType, Name: varId.Name.Name()}
			varId.BackendType = cVar
			val := c.emitExpr(stmt.Expr)
			return fmt.Sprintf("%s %s = %s", cType, varId.Name.Name(), val)
		}
	case ast.KIND_ASSIGNMENT_STMT:
		stmt := node.Node.(*ast.AssignmentStmt)
		if len(stmt.Targets) == 1 {
			t := c.emitExpr(stmt.Targets[0])
			val := c.emitExpr(stmt.Values[0])
			return fmt.Sprintf("%s = %s", t, val)
		}
	}
	return c.emitExpr(node)
}

// emitForUpdate extracts the update expression for a for loop header.
func (c *CCodegen) emitForUpdate(node *ast.Node) string {
	if node == nil {
		return ""
	}
	switch node.Kind {
	case ast.KIND_ASSIGNMENT_STMT:
		stmt := node.Node.(*ast.AssignmentStmt)
		if len(stmt.Targets) == 1 {
			t := c.emitExpr(stmt.Targets[0])
			val := c.emitExpr(stmt.Values[0])
			return fmt.Sprintf("%s = %s", t, val)
		}
	case ast.KIND_VAR_STMT:
		// The parser emits the for-loop update as a VarStmt (reassignment).
		stmt := node.Node.(*ast.VarStmt)
		if !stmt.IsDecl && len(stmt.Names) == 1 {
			target := c.emitVarTarget(stmt.Names[0])
			val := c.emitExpr(stmt.Expr)
			return fmt.Sprintf("%s = %s", target, val)
		}
	}
	return c.emitExpr(node)
}

// emitVarTarget resolves the C name for the left-hand side of a VarStmt reassignment.
// VarStmt.Names entries are KIND_VAR_ID_STMT nodes, not expressions.
func (c *CCodegen) emitVarTarget(node *ast.Node) string {
	switch node.Kind {
	case ast.KIND_VAR_ID_STMT:
		varId := node.Node.(*ast.VarIdStmt)
		if varId.BackendType != nil {
			return varId.BackendType.(*CVariable).Name
		}
		return varId.Name.Name()
	default:
		// For field access or pointer deref targets, fall through to emitExpr.
		return c.emitExpr(node)
	}
}

// emitWhileLoop emits a while loop.
func (c *CCodegen) emitWhileLoop(stmt *ast.WhileLoop, indent string) {
	inner := indent + "    "
	cond := c.emitExpr(stmt.Cond)
	c.buf.WriteString(fmt.Sprintf("while (%s) {\n", cond))
	c.emitBlock(stmt.Block, inner)
	c.buf.WriteString(indent + "}\n")
}
