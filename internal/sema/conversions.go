package sema

import (
	"fmt"

	"github.com/HicaroD/Telia/internal/ast"
	"github.com/HicaroD/Telia/internal/lexer/token"
)

func isConversionBuiltin(name string) bool {
	return name == "widen" || name == "narrow" || name == "ptr_cast"
}

func unsignedInteger(kind token.Kind) bool {
	switch kind {
	case token.U8_TYPE,
		token.U16_TYPE,
		token.U32_TYPE,
		token.U64_TYPE,
		token.U128_TYPE,
		token.UINT_TYPE:
		return true
	}
	return false
}

func (s *sema) checkConversion(
	call *ast.FnCall,
	scope, decl *ast.Scope,
	imported bool,
) (*ast.ExprType, error) {
	name := call.Name.Name()
	if len(call.TypeArgs) != 1 || len(call.Args) != 1 {
		return nil, fmt.Errorf("%s requires one type argument and one value argument", name)
	}
	target, err := s.resolveType(call.TypeArgs[0], scope)
	if err != nil {
		return nil, err
	}
	source, _, err := s.inferExprTypeWithoutContext(call.Args[0], scope, decl, imported, false)
	if err != nil {
		return nil, err
	}
	if name == "ptr_cast" {
		if !(source.IsPointer() || source.Equals(ast.RAWPTR_TYPE)) ||
			!(target.IsPointer() || target.Equals(ast.RAWPTR_TYPE)) {
			return nil, fmt.Errorf(
				"ptr_cast requires object-pointer or rawptr source and destination types",
			)
		}
		if source.IsPointer() {
			pointer := source.T.(*ast.PointerType)
			if pointer.Const && (!target.IsPointer() || !target.T.(*ast.PointerType).Const) {
				return nil, fmt.Errorf("ptr_cast cannot discard const protection")
			}
			if pointer.Type.IsPointer() && !pointerAssignable(source, target) {
				return nil, fmt.Errorf("ptr_cast cannot reinterpret nested pointer qualifiers")
			}
		}
		call.Builtin, call.BuiltinType, call.ResultType = name, target, target
		return target, nil
	}
	if !source.IsInteger() || !target.IsInteger() {
		return nil, fmt.Errorf("%s requires integer source and destination types", name)
	}
	if value, err := evaluateConstant(call.Args[0]); err == nil {
		if err := checkConstantRange(value, source, call.Args[0].Span); err != nil {
			return nil, err
		}
	}
	a, b := source.T.(*ast.BasicType).Kind, target.T.(*ast.BasicType).Kind
	preservesRange := unsignedInteger(a) == unsignedInteger(b) && b.BitSize() >= a.BitSize() ||
		unsignedInteger(a) && !unsignedInteger(b) && b.BitSize() > a.BitSize()
	if name == "widen" && !preservesRange {
		return nil, fmt.Errorf(
			"widen must preserve every source value; use narrow for a checked conversion",
		)
	}
	call.Builtin, call.BuiltinType, call.ResultType = name, target, target
	call.ConversionSource = source
	if name == "narrow" {
		call.ResultType = &ast.ExprType{
			Kind: ast.EXPR_TYPE_TUPLE,
			T: &ast.TupleType{
				Types: []*ast.ExprType{target, ast.NewBasicType(token.ERROR_TYPE)},
			},
		}
		s.registerTupleType(call.ResultType)
	}
	return call.ResultType, nil
}
