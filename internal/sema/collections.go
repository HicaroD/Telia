package sema

import (
	"fmt"
	"strconv"

	"github.com/HicaroD/Telia/internal/ast"
	"github.com/HicaroD/Telia/internal/lexer/token"
)

func isMemoryBuiltin(name string) bool {
	switch name {
	case "sizeof", "alignof", "zero", "slice_from_raw", "slice_data":
		return true
	}
	return false
}

func (s *sema) inferIndex(
	index *ast.IndexExpr,
	scope, decl *ast.Scope,
	imported bool,
) (*ast.ExprType, error) {
	ty, _, err := s.inferExprTypeWithoutContext(index.Value, scope, decl, imported, false)
	if err != nil {
		return nil, err
	}
	if ty.Kind != ast.EXPR_TYPE_ARRAY && ty.Kind != ast.EXPR_TYPE_SLICE {
		return nil, fmt.Errorf("indexing requires an array or slice")
	}
	if err := s.checkCollectionIndex(index.Index, scope, decl, imported); err != nil {
		return nil, err
	}
	index.Type = ty
	return ty.T.(*ast.SequenceType).Element, nil
}

func (s *sema) checkCollectionIndex(index *ast.Node, scope, decl *ast.Scope, imported bool) error {
	if unary, ok := index.Node.(*ast.UnaryExpr); ok && unary.Op == token.MINUS {
		return fmt.Errorf("collection indices must be nonnegative")
	}
	if literal, ok := index.Node.(*ast.LiteralExpr); ok && literal.Type.IsInteger() {
		if _, err := strconv.ParseUint(string(literal.Value), 10, 64); err != nil {
			return fmt.Errorf("collection index is outside u64 range")
		}
		_, err := s.inferExprTypeWithContext(
			index,
			concreteBasicType(token.U64_TYPE),
			scope,
			decl,
			imported,
			false,
		)
		return err
	}
	ty, _, err := s.inferExprTypeWithoutContext(index, scope, decl, imported, false)
	if err != nil {
		return err
	}
	if !ty.Equals(ast.NewBasicType(token.U64_TYPE)) {
		return fmt.Errorf(
			"collection index requires u64; other integer types need explicit conversion",
		)
	}
	return nil
}

func (s *sema) checkMemoryBuiltin(
	call *ast.FnCall,
	scope, decl *ast.Scope,
	imported bool,
) (*ast.ExprType, error) {
	name := call.Name.Name()
	count := 0
	if name == "slice_from_raw" {
		count = 2
	}
	if name == "slice_data" {
		count = 1
	}
	if len(call.TypeArgs) != 1 || len(call.Args) != count {
		return nil, fmt.Errorf("%s requires one type argument and %d value arguments", name, count)
	}
	ty, err := s.resolveType(call.TypeArgs[0], scope)
	if err != nil {
		return nil, err
	}
	if err := rejectOpaqueValue(ty, call.Name.Span); err != nil {
		return nil, err
	}
	if !valueType(ty) {
		return nil, fmt.Errorf("%s requires a concrete value type", name)
	}
	call.Builtin, call.BuiltinType = name, ty
	result := concreteBasicType(token.U64_TYPE)
	if name == "zero" {
		result = ty
	}
	if name == "slice_from_raw" || name == "slice_data" {
		view := &ast.ExprType{Kind: ast.EXPR_TYPE_SLICE, T: &ast.SequenceType{Element: ty}}
		s.program.SequenceTypes = append(s.program.SequenceTypes, view)
		if name == "slice_from_raw" {
			if _, err := s.inferExprTypeWithContext(call.Args[0], ast.NewBasicType(token.RAWPTR_TYPE), scope, decl, imported, false); err != nil {
				return nil, err
			}
			if _, err := s.inferExprTypeWithContext(call.Args[1], ast.NewBasicType(token.U64_TYPE), scope, decl, imported, false); err != nil {
				return nil, err
			}
			result = view
		} else {
			if _, err := s.inferExprTypeWithContext(call.Args[0], view, scope, decl, imported, false); err != nil {
				return nil, err
			}
			result = ast.NewBasicType(token.RAWPTR_TYPE)
		}
	}
	call.ResultType = result
	return result, nil
}

func concreteBasicType(kind token.Kind) *ast.ExprType {
	return &ast.ExprType{
		Kind:     ast.EXPR_TYPE_BASIC,
		Explicit: true,
		T:        &ast.BasicType{Kind: kind, Explicit: true},
	}
}

func (s *sema) checkSliceCall(
	call *ast.FnCall,
	scope, decl *ast.Scope,
	imported bool,
) (*ast.ExprType, error) {
	ty, _, err := s.inferExprTypeWithoutContext(call.Args[0], scope, decl, imported, false)
	if err != nil {
		return nil, err
	}
	if s.readOnlyStorage(call.Args[0], scope) {
		return nil, constMutation(call.Args[0].Span)
	}
	if ty.Kind != ast.EXPR_TYPE_ARRAY && ty.Kind != ast.EXPR_TYPE_SLICE {
		return nil, fmt.Errorf("slicing requires an array or slice")
	}
	if ty.Kind == ast.EXPR_TYPE_ARRAY && !call.Args[0].IsAddressable() {
		return nil, fmt.Errorf("borrowing requires addressable array storage")
	}
	for _, bound := range call.Args[1:] {
		if bound != nil {
			if err := s.checkCollectionIndex(bound, scope, decl, imported); err != nil {
				return nil, err
			}
		}
	}
	call.BuiltinType = ty
	view := &ast.ExprType{
		Kind: ast.EXPR_TYPE_SLICE,
		T:    &ast.SequenceType{Element: ty.T.(*ast.SequenceType).Element},
	}
	s.program.SequenceTypes = append(s.program.SequenceTypes, view)
	call.ResultType = view
	if call.Builtin == "slice" {
		call.ResultType = &ast.ExprType{
			Kind: ast.EXPR_TYPE_TUPLE,
			T:    &ast.TupleType{Types: []*ast.ExprType{view, ast.NewBasicType(token.ERROR_TYPE)}},
		}
		s.registerTupleType(call.ResultType)
	}
	return call.ResultType, nil
}
