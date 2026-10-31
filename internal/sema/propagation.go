package sema

import (
	"fmt"

	"github.com/HicaroD/Telia/internal/ast"
)

// Type-check calls against their full result first, then consume only the
// successful value when a propagating call is used as a value expression.
func (s *sema) propagatedValue(node *ast.Node, result *ast.ExprType) (*ast.ExprType, error) {
	for node.Kind == ast.KIND_NAMESPACE_ACCESS {
		node = node.Node.(*ast.NamespaceAccess).Right
	}
	call, ok := node.Node.(*ast.FnCall)
	if !ok || call.AtOp == nil || !call.AtOp.Propagates {
		return result, nil
	}
	if result.IsError() {
		return nil, fmt.Errorf(
			"@try on an error-only call has no success value; use it as a statement",
		)
	}
	call.HandledValue = true
	call.ResultType = result
	success := unwrapErrorFromTuple(result)
	if success.Kind == ast.EXPR_TYPE_TUPLE {
		s.registerTupleType(success)
	}
	return success, nil
}
