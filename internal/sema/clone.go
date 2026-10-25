package sema

import "github.com/HicaroD/Telia/internal/ast"

// Clone syntax only. Resolved symbols, inferred types, backend values and
// previously populated scopes must never cross specialization boundaries.
func cloneType(t *ast.ExprType, bindings map[string]*ast.ExprType) *ast.ExprType {
	if t == nil {
		return nil
	}
	copy := *t
	switch v := t.T.(type) {
	case *ast.IdType:
		if v.Namespace == nil && len(v.Args) == 0 && bindings[v.Name.Name()] != nil {
			return cloneType(bindings[v.Name.Name()], nil)
		}
		copy.T = &ast.IdType{Name: v.Name, Namespace: v.Namespace, Args: cloneTypes(v.Args, bindings)}
	case *ast.PointerType:
		copy.T = &ast.PointerType{Type: cloneType(v.Type, bindings), Explicit: v.Explicit}
	case *ast.TupleType:
		copy.T = &ast.TupleType{Types: cloneTypes(v.Types, bindings)}
	case *ast.BasicType:
		b := *v
		copy.T = &b
	case *ast.StructType:
		copy.T = &ast.StructType{Decl: v.Decl}
	}
	return &copy
}

func cloneTypes(types []*ast.ExprType, bindings map[string]*ast.ExprType) []*ast.ExprType {
	var result []*ast.ExprType
	for _, t := range types {
		result = append(result, cloneType(t, bindings))
	}
	return result
}

type syntaxClone struct{ bindings map[string]*ast.ExprType }

func (c syntaxClone) block(b *ast.BlockStmt, scope *ast.Scope) *ast.BlockStmt {
	if b == nil {
		return nil
	}
	copy := &ast.BlockStmt{OpenCurly: b.OpenCurly, CloseCurly: b.CloseCurly, FoundReturn: b.FoundReturn}
	copy.Statements = c.nodes(b.Statements, scope)
	for _, n := range copy.Statements {
		if n.Kind == ast.KIND_DEFER_STMT {
			copy.DeferStack = append(copy.DeferStack, n.Node.(*ast.DeferStmt))
		}
	}
	return copy
}

func (c syntaxClone) nodes(nodes []*ast.Node, scope *ast.Scope) []*ast.Node {
	var copies []*ast.Node
	for _, n := range nodes {
		copies = append(copies, c.node(n, scope))
	}
	return copies
}

func (c syntaxClone) node(n *ast.Node, scope *ast.Scope) *ast.Node {
	if n == nil {
		return nil
	}
	copy := &ast.Node{Kind: n.Kind}
	switch v := n.Node.(type) {
	case nil:
	case *ast.LiteralExpr:
		copy.Node = &ast.LiteralExpr{Type: cloneType(v.Type, c.bindings), Value: v.Value}
	case *ast.IdExpr:
		copy.Node = &ast.IdExpr{Name: v.Name}
	case *ast.NullPtrExpr:
		copy.Node = &ast.NullPtrExpr{Type: cloneType(v.Type, c.bindings)}
	case *ast.UnaryExpr:
		copy.Node = &ast.UnaryExpr{Op: v.Op, Value: c.node(v.Value, scope)}
	case *ast.BinExpr:
		copy.Node = &ast.BinExpr{Op: v.Op, Left: c.node(v.Left, scope), Right: c.node(v.Right, scope)}
	case *ast.AddressOfExpr:
		copy.Node = &ast.AddressOfExpr{Expr: c.node(v.Expr, scope)}
	case *ast.DerefPointerExpr:
		copy.Node = &ast.DerefPointerExpr{Expr: c.node(v.Expr, scope)}
	case *ast.TupleExpr:
		copy.Node = &ast.TupleExpr{Exprs: c.nodes(v.Exprs, scope)}
	case *ast.VarArgsExpr:
		copy.Node = &ast.VarArgsExpr{Args: c.nodes(v.Args, scope)}
	case *ast.StructLiteralExpr:
		lit := &ast.StructLiteralExpr{Name: v.Name, Namespace: v.Namespace, TypeArgs: cloneTypes(v.TypeArgs, c.bindings)}
		for _, f := range v.Values {
			lit.Values = append(lit.Values, &ast.StructFieldValue{Name: f.Name, Value: c.node(f.Value, scope)})
		}
		copy.Node = lit
	case *ast.FieldAccess:
		copy.Node = &ast.FieldAccess{Left: &ast.IdExpr{Name: v.Left.Name}, Right: c.node(v.Right, scope)}
	case *ast.NamespaceAccess:
		copy.Node = &ast.NamespaceAccess{Left: &ast.IdExpr{Name: v.Left.Name}, Right: c.node(v.Right, scope), IsImport: v.IsImport}
	case *ast.FnCall:
		call := &ast.FnCall{Name: v.Name, TypeArgs: cloneTypes(v.TypeArgs, c.bindings), Args: c.nodes(v.Args, scope), Variadic: v.Variadic}
		if v.AtOp != nil {
			call.AtOp = &ast.AtOperator{Kind: v.AtOp.Kind}
			if op, ok := v.AtOp.Op.(*ast.CatchAtOperator); ok {
				call.AtOp.Op = &ast.CatchAtOperator{ErrVarName: op.ErrVarName, Scope: ast.NewScope(scope), Block: c.block(op.Block, scope)}
			}
		}
		copy.Node = call
	case *ast.VarIdStmt:
		copy.Node = &ast.VarIdStmt{Name: v.Name, Type: cloneType(v.Type, c.bindings), NeedsInference: v.NeedsInference, Pointer: v.Pointer, NumberOfPointerReceivers: v.NumberOfPointerReceivers}
	case *ast.VarStmt:
		copy.Node = &ast.VarStmt{IsDecl: v.IsDecl, HasFieldAccess: v.HasFieldAccess, Pointer: v.Pointer, Names: c.nodes(v.Names, scope), Expr: c.node(v.Expr, scope)}
	case *ast.AssignmentStmt:
		copy.Node = &ast.AssignmentStmt{Decl: v.Decl, Targets: c.nodes(v.Targets, scope), Values: c.nodes(v.Values, scope)}
	case *ast.ReturnStmt:
		copy.Node = &ast.ReturnStmt{Return: v.Return, Value: c.node(v.Value, scope)}
	case *ast.DeferStmt:
		copy.Node = &ast.DeferStmt{Stmt: c.node(v.Stmt, scope)}
	case *ast.BlockStmt:
		copy.Node = c.block(v, ast.NewScope(scope))
	case *ast.CondStmt:
		cond := &ast.CondStmt{IfStmt: c.ifCond(v.IfStmt, scope)}
		for _, elif := range v.ElifStmts {
			cond.ElifStmts = append(cond.ElifStmts, c.ifCond(elif, scope))
		}
		if v.ElseStmt != nil {
			inner := ast.NewScope(scope)
			cond.ElseStmt = &ast.ElseCond{Else: v.ElseStmt.Else, Scope: inner, Block: c.block(v.ElseStmt.Block, inner)}
		}
		copy.Node = cond
	case *ast.ForLoop:
		inner := ast.NewScope(scope)
		copy.Node = &ast.ForLoop{Scope: inner, Init: c.node(v.Init, inner), Cond: c.node(v.Cond, inner), Update: c.node(v.Update, inner), Block: c.block(v.Block, inner)}
	case *ast.WhileLoop:
		inner := ast.NewScope(scope)
		copy.Node = &ast.WhileLoop{Scope: inner, Cond: c.node(v.Cond, inner), Block: c.block(v.Block, inner)}
	default:
		panic("unhandled syntax node in generic clone: " + n.String())
	}
	return copy
}

func (c syntaxClone) ifCond(v *ast.IfElifCond, scope *ast.Scope) *ast.IfElifCond {
	inner := ast.NewScope(scope)
	return &ast.IfElifCond{If: v.If, Scope: inner, Expr: c.node(v.Expr, inner), Block: c.block(v.Block, inner)}
}

func cloneFunction(template *ast.FnDecl, bindings map[string]*ast.ExprType) *ast.FnDecl {
	fn := &ast.FnDecl{Name: template.Name, Owner: template.Owner, File: template.File, Attributes: template.Attributes, Scope: ast.NewScope(template.Owner.Scope), RetType: cloneType(template.RetType, bindings)}
	params := *template.Params
	params.Fields = nil
	fn.Params = &params
	for _, p := range template.Params.Fields {
		param := &ast.Param{Name: p.Name, Attributes: p.Attributes, Variadic: p.Variadic, Type: cloneType(p.Type, bindings)}
		fn.Params.Fields = append(fn.Params.Fields, param)
		fn.Scope.Nodes[param.Name.Name()] = &ast.Node{Kind: ast.KIND_PARAM, Node: param}
	}
	fn.Block = (syntaxClone{bindings}).block(template.Block, fn.Scope)
	return fn
}
