package sema

import (
	"fmt"

	"github.com/HicaroD/Telia/internal/ast"
	"github.com/HicaroD/Telia/internal/lexer/token"
)

func (s *sema) checkError(e *ast.ErrorExpr, scope *ast.Scope) (*ast.ExprType, error) {
	if len(e.Args) != 1 && len(e.Args) != 2 {
		return nil, fmt.Errorf(
			"error constructor requires a message or an error-enum value and a message",
		)
	}
	if len(e.Args) == 2 {
		ty, _, err := s.inferExprTypeWithoutContext(e.Args[0], scope, scope, false, false)
		if err != nil {
			return nil, err
		}
		if ty.Kind != ast.EXPR_TYPE_ENUM || !ty.T.(*ast.EnumType).Decl.IsError {
			return nil, fmt.Errorf("error identity requires an enum declared with : error")
		}
		e.Enum = ty.T.(*ast.EnumType).Decl
	}
	if _, err := s.inferExprTypeWithContext(e.Args[len(e.Args)-1], ast.NewBasicType(token.STRING_TYPE), scope, scope, false, false); err != nil {
		return nil, err
	}
	return ast.NewBasicType(token.ERROR_TYPE), nil
}

// A bare variant in a kind comparison names its identity, not a payload value.
func (s *sema) markErrorIdentityComparison(b *ast.BinExpr, scope *ast.Scope) {
	for _, pair := range [][2]*ast.Node{{b.Left, b.Right}, {b.Right, b.Left}} {
		if pair[0].Kind != ast.KIND_FIELD_ACCESS && pair[0].Kind != ast.KIND_ID_EXPR {
			continue
		}
		if pair[1].Kind != ast.KIND_NAMESPACE_ACCESS {
			continue
		}
		ty, _, err := s.inferExprTypeWithoutContext(pair[0], scope, scope, false, false)
		if err != nil || ty.Kind != ast.EXPR_TYPE_ERROR_KIND {
			continue
		}
		var mark func(*ast.Node)
		mark = func(n *ast.Node) {
			if ns, ok := n.Node.(*ast.NamespaceAccess); ok {
				ns.ErrorIdentity = true
				mark(ns.Right)
			}
		}
		mark(pair[1])
	}
}
