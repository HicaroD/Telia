package sema

import (
	"fmt"
	"path/filepath"

	"github.com/HicaroD/Telia/config"
	"github.com/HicaroD/Telia/internal/ast"
	"github.com/HicaroD/Telia/internal/diagnostics"
	"github.com/HicaroD/Telia/internal/lexer/token"
)

func testingAssertion(fn *ast.FnDecl) bool {
	if config.ENVS == nil || fn.Owner == nil || fn.Owner.Loc == nil ||
		filepath.Clean(fn.Owner.Loc.Path) != filepath.Join(config.ENVS.STD, "testing") {
		return false
	}
	switch fn.Name.Name() {
	case "assert", "assert_eq", "assert_ne":
		return true
	}
	return false
}

func (s *sema) checkTestingAssertion(call *ast.FnCall, scope *ast.Scope) (*ast.ExprType, error) {
	name := call.Name.Name()
	count := 2
	if name == "assert" {
		count = 1
	}
	if len(call.Args) != count || len(call.TypeArgs) != 0 || call.AtOp != nil {
		return nil, diagnostics.Diag{
			Title: "Invalid assertion",
			Message: fmt.Sprintf(
				"%s requires %d value arguments and no type arguments or error operator",
				name,
				count,
			),
			Span: call.Name.Span,
		}
	}
	var ty *ast.ExprType
	var err error
	if name == "assert" {
		ty = ast.NewBasicType(token.BOOL_TYPE)
		_, err = s.inferExprTypeWithContext(call.Args[0], ty, scope, scope, false, true)
	} else {
		ty, _, err = s.inferExprTypeWithoutContext(call.Args[0], scope, scope, false, true)
	}
	if err != nil {
		return nil, err
	}
	if ty.Kind != ast.EXPR_TYPE_BASIC ||
		!(ty.IsNumeric() || ty.Equals(ast.NewBasicType(token.BOOL_TYPE)) || ty.Equals(ast.NewBasicType(token.STRING_TYPE))) {
		return nil, diagnostics.Diag{
			Title:   "Unsupported assertion type",
			Message: "equality assertions support booleans, integers, floats, and strings",
			Span:    call.Args[0].Span,
		}
	}
	if count == 2 {
		if _, err := s.inferExprTypeWithContext(call.Args[1], ty, scope, scope, false, true); err != nil {
			return nil, err
		}
	}
	call.Builtin, call.BuiltinType = "testing_"+name, ty
	call.ResultType = ast.NewBasicType(token.VOID_TYPE)
	return call.ResultType, nil
}
