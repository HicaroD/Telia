package c

import (
	"fmt"

	"github.com/HicaroD/Telia/internal/ast"
)

// emitBlock emits each statement in a block at the given indentation level.
//
// Defer is block-scoped in Telia: deferred statements run before the return
// of the block they appear in, not necessarily before the function return.
// The DeferStack is flushed in two places:
//  1. Immediately before any explicit return statement in this block.
//  2. At the end of the block when there is no explicit return (fall-through),
//     e.g. void functions and nested if-blocks. This is guarded by
//     !block.FoundReturn to avoid double-emission.
func (c *CCodegen) emitBlock(block *ast.BlockStmt, indent string) {
	for _, stmt := range block.Statements {
		// Before a return, flush the defer stack in reverse order.
		if stmt.Kind == ast.KIND_RETURN_STMT {
			c.flushDeferStack(block, indent)
		}
		c.buf.WriteString(indent)
		c.emitStmtNode(stmt, indent)
	}
	// Flush for fall-through blocks (no explicit return): void functions and
	// nested blocks such as if-bodies that end without a return statement.
	if !block.FoundReturn {
		c.flushDeferStack(block, indent)
	}
}

// flushDeferStack emits all non-skipped deferred statements from block in LIFO
// order at the given indentation level.
func (c *CCodegen) flushDeferStack(block *ast.BlockStmt, indent string) {
	for i := len(block.DeferStack) - 1; i >= 0; i-- {
		d := block.DeferStack[i]
		if !d.Skip {
			c.buf.WriteString(indent)
			c.emitStmtNode(d.Stmt, indent)
			c.buf.WriteString(";\n")
		}
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
		if call.AtOp != nil && call.AtOp.Kind == ast.AT_OPERATOR_FAIL {
			c.emitAtFailBareCall(call)
		} else if call.AtOp != nil && call.AtOp.Kind == ast.AT_OPERATOR_CATCH {
			c.emitAtCatchBareCall(call, indent)
		} else {
			c.buf.WriteString(c.emitFnCall(call))
			c.buf.WriteString(";\n")
		}
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
	// Special case: return nil in error function -> return (_Error){NULL};
	if ret.Value.Kind == ast.KIND_NULLPTR_EXPR && c.currentRetTy != nil && c.currentRetTy.IsError() {
		c.buf.WriteString("return (_Error){NULL};\n")
		return
	}
	c.buf.WriteString(fmt.Sprintf("return %s;\n", c.emitExpr(ret.Value)))
}

// emitAtFailBareCall emits a bare @fail function call (statement position).
// For error-only functions: _Error _t0 = fn(); if (_t0.msg != NULL) { _panic(_t0.msg); }
// For tuple functions:      _Tuple... _t0 = fn(); if (_t0._N.msg != NULL) { _panic(...); }
func (c *CCodegen) emitAtFailBareCall(call *ast.FnCall) {
	fullRetType := getFullFnRetType(call)
	fnCallExpr := c.emitFnCall(call)
	tmp := c.nextTmp()

	if fullRetType == nil {
		// Error-only: Decl/Proto not available or not a tuple.
		c.buf.WriteString(fmt.Sprintf("_Error %s = %s;\n", tmp, fnCallExpr))
		c.buf.WriteString(fmt.Sprintf("if (%s.msg != NULL) { _panic(%s.msg); }\n", tmp, tmp))
	} else {
		// Tuple: store full result, check error field (last element).
		tupleName := tupleTypedefName(fullRetType)
		c.buf.WriteString(fmt.Sprintf("%s %s = %s;\n", tupleName, tmp, fnCallExpr))
		tt := fullRetType.T.(*ast.TupleType)
		errIdx := len(tt.Types) - 1
		c.buf.WriteString(fmt.Sprintf("if (%s._%d.msg != NULL) { _panic(%s._%d.msg); }\n", tmp, errIdx, tmp, errIdx))
	}
}

// emitAtCatchBareCall emits a bare @catch function call (statement position).
// For error-only functions:
//
//	_Error _t0 = fn();
//	if (_t0.msg != NULL) {
//	    _Error err = _t0;
//	    <handler block>
//	}
//
// For tuple functions:
//
//	_Tuple... _t0 = fn();
//	if (_t0._N.msg != NULL) {
//	    _Error err = _t0._N;
//	    <handler block>
//	}
func (c *CCodegen) emitAtCatchBareCall(call *ast.FnCall, indent string) {
	catchOp := call.AtOp.Op.(*ast.CatchAtOperator)
	errVarName := catchOp.ErrVarName.Name()
	fullRetType := getFullFnRetType(call)
	fnCallExpr := c.emitFnCall(call)
	tmp := c.nextTmp()
	inner := indent + "    "

	if fullRetType == nil {
		// Error-only function
		c.buf.WriteString(fmt.Sprintf("_Error %s = %s;\n", tmp, fnCallExpr))
		c.buf.WriteString(fmt.Sprintf("if (%s.msg != NULL) {\n", tmp))
		c.buf.WriteString(fmt.Sprintf("%s_Error %s = %s;\n", inner, errVarName, tmp))
		c.emitBlock(catchOp.Block, inner)
		c.buf.WriteString(indent + "}\n")
	} else {
		// Tuple function
		tupleName := tupleTypedefName(fullRetType)
		c.buf.WriteString(fmt.Sprintf("%s %s = %s;\n", tupleName, tmp, fnCallExpr))
		tt := fullRetType.T.(*ast.TupleType)
		errIdx := len(tt.Types) - 1
		c.buf.WriteString(fmt.Sprintf("if (%s._%d.msg != NULL) {\n", tmp, errIdx))
		c.buf.WriteString(fmt.Sprintf("%s_Error %s = %s._%d;\n", inner, errVarName, tmp, errIdx))
		c.emitBlock(catchOp.Block, inner)
		c.buf.WriteString(indent + "}\n")
	}
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

			// @fail on a tuple-returning function: store full tuple in temp,
			// check error field, panic if non-nil, then extract non-error fields.
			if isAtFailTupleCall(stmt) {
				fnCall := stmt.Expr.Node.(*ast.FnCall)
				fullRetType := getFullFnRetType(fnCall)
				tupleName := tupleTypedefName(fullRetType)
				tmp := c.nextTmp()
				fnCallExpr := c.emitExpr(stmt.Expr)
				c.buf.WriteString(fmt.Sprintf("%s %s = %s;\n", tupleName, tmp, fnCallExpr))

				tt := fullRetType.T.(*ast.TupleType)
				errIdx := len(tt.Types) - 1
				c.buf.WriteString(fmt.Sprintf("if (%s._%d.msg != NULL) { _panic(%s._%d.msg); }\n", tmp, errIdx, tmp, errIdx))

				if len(tt.Types) > 2 {
					for i := 0; i < errIdx; i++ {
						elemType := emitCType(tt.Types[i])
						c.buf.WriteString(fmt.Sprintf("%s %s_%d = %s._%d;\n", elemType, tmp, i, tmp, i))
					}
				} else {
					c.buf.WriteString(fmt.Sprintf("%s %s = %s._0;\n", cType, varId.Name.Name(), tmp))
				}
			} else if isAtCatchTupleCall(stmt) {
				// @catch on a tuple-returning function: store full tuple in temp,
				// check error field, branch to handler or extract non-error value.
				// Declare the variable before the if/else so it's available after.
				fnCall := stmt.Expr.Node.(*ast.FnCall)
				catchOp := fnCall.AtOp.Op.(*ast.CatchAtOperator)
				errVarName := catchOp.ErrVarName.Name()
				fullRetType := getFullFnRetType(fnCall)
				tupleName := tupleTypedefName(fullRetType)
				tmp := c.nextTmp()
				fnCallExpr := c.emitExpr(stmt.Expr)
				inner := indent + "    "

				// Declare variable upfront
				c.buf.WriteString(fmt.Sprintf("%s %s;\n", cType, varId.Name.Name()))

				c.buf.WriteString(fmt.Sprintf("%s %s = %s;\n", tupleName, tmp, fnCallExpr))

				tt := fullRetType.T.(*ast.TupleType)
				errIdx := len(tt.Types) - 1
				c.buf.WriteString(fmt.Sprintf("if (%s._%d.msg != NULL) {\n", tmp, errIdx))
				c.buf.WriteString(fmt.Sprintf("%s_Error %s = %s._%d;\n", inner, errVarName, tmp, errIdx))
				c.emitBlock(catchOp.Block, inner)
				c.buf.WriteString(indent + "} else {\n")

				if len(tt.Types) > 2 {
					for i := 0; i < errIdx; i++ {
						elemType := emitCType(tt.Types[i])
						c.buf.WriteString(fmt.Sprintf("%s%s %s_%d = %s._%d;\n", inner, elemType, tmp, i, tmp, i))
					}
				} else {
					c.buf.WriteString(fmt.Sprintf("%s%s = %s._0;\n", inner, varId.Name.Name(), tmp))
				}
				c.buf.WriteString(indent + "}\n")
			} else {
				val := c.emitExpr(stmt.Expr)
				c.buf.WriteString(fmt.Sprintf("%s %s = %s;\n", cType, varId.Name.Name(), val))
			}
		} else if stmt.Expr.Kind == ast.KIND_TUPLE_LITERAL_EXPR {
			// Multi-variable declaration from a literal tuple: `a, b := 1, 2`.
			// Sema split the expressions and set each VarIdStmt.Type directly;
			// TupleExpr.Type is not set. Emit each variable from its paired expr.
			te := stmt.Expr.Node.(*ast.TupleExpr)
			for i, nameNode := range stmt.Names {
				varId := nameNode.Node.(*ast.VarIdStmt)
				cType := emitCType(varId.Type)
				varId.BackendType = &CVariable{CType: cType, Name: varId.Name.Name()}
				val := c.emitExpr(te.Exprs[i])
				c.buf.WriteString(fmt.Sprintf("%s %s = %s;\n", cType, varId.Name.Name(), val))
				if i < len(stmt.Names)-1 {
					c.buf.WriteString(indent)
				}
			}
		} else {
			// Multi-variable declaration from a tuple-returning function call:
			// `a, b := fnCall()`. Emit:
			//   _Tuple_X _t0 = fnCall();
			//   T0 a = _t0._0;
			//   T1 b = _t0._1;
			elemTypes := make([]*ast.ExprType, len(stmt.Names))
			for i, nameNode := range stmt.Names {
				elemTypes[i] = nameNode.Node.(*ast.VarIdStmt).Type
			}
			tupleType := &ast.ExprType{
				Kind: ast.EXPR_TYPE_TUPLE,
				T:    &ast.TupleType{Types: elemTypes},
			}
			tupleName := tupleTypedefName(tupleType)
			tmp := c.nextTmp()
			val := c.emitExpr(stmt.Expr)
			c.buf.WriteString(fmt.Sprintf("%s %s = %s;\n", tupleName, tmp, val))
			for i, nameNode := range stmt.Names {
				varId := nameNode.Node.(*ast.VarIdStmt)
				cType := emitCType(varId.Type)
				varId.BackendType = &CVariable{CType: cType, Name: varId.Name.Name()}
				c.buf.WriteString(indent)
				c.buf.WriteString(fmt.Sprintf("%s %s = %s._%d;\n", cType, varId.Name.Name(), tmp, i))
			}
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

// isAtFailTupleCall reports whether stmt is a single-var declaration from a
// @fail call on a tuple-returning function.
func isAtFailTupleCall(stmt *ast.VarStmt) bool {
	if stmt.Expr == nil || stmt.Expr.Kind != ast.KIND_FN_CALL {
		return false
	}
	call := stmt.Expr.Node.(*ast.FnCall)
	return call.AtOp != nil && call.AtOp.Kind == ast.AT_OPERATOR_FAIL &&
		getFullFnRetType(call) != nil
}

// isAtCatchTupleCall reports whether stmt is a single-var declaration from a
// @catch call on a tuple-returning function.
func isAtCatchTupleCall(stmt *ast.VarStmt) bool {
	if stmt.Expr == nil || stmt.Expr.Kind != ast.KIND_FN_CALL {
		return false
	}
	call := stmt.Expr.Node.(*ast.FnCall)
	return call.AtOp != nil && call.AtOp.Kind == ast.AT_OPERATOR_CATCH &&
		getFullFnRetType(call) != nil
}

// getFullFnRetType returns the function's declared return type (before sema
// unwrapping).  It checks Decl first, then Proto. Returns nil if unavailable
// or not a tuple type.
func getFullFnRetType(call *ast.FnCall) *ast.ExprType {
	var retType *ast.ExprType
	if call.Decl != nil {
		retType = call.Decl.RetType
	} else if call.Proto != nil {
		retType = call.Proto.RetType
	}
	if retType == nil || retType.Kind != ast.EXPR_TYPE_TUPLE {
		return nil
	}
	return retType
}
