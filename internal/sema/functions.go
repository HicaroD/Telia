package sema

import (
	"github.com/HicaroD/Telia/internal/ast"
	"github.com/HicaroD/Telia/internal/diagnostics"
)

func (s *sema) functionValue(id *ast.IdExpr, symbol *ast.Node) (*ast.ExprType, error) {
	var params *ast.Params
	var result *ast.ExprType
	var attrs ast.Attributes
	switch d := symbol.Node.(type) {
	case *ast.FnDecl:
		if len(d.TypeParams) != 0 {
			return nil, diagnostics.Diag{Title: "Generic function value", Message: "a generic function needs concrete type arguments", Span: id.Name.Span, Help: "Use a non-generic wrapper with the required callback signature."}
		}
		if err := s.resolveSignature(d); err != nil {
			return nil, err
		}
		params, result = d.Params, d.RetType
		attrs = d.Attributes
	case *ast.Proto:
		params, result = d.Params, d.RetType
		attrs = d.Attributes
	}
	ty := &ast.ExprType{Kind: ast.EXPR_TYPE_FUNCTION, Explicit: true, T: &ast.FunctionType{Params: params, Result: result, Attributes: attrs}}
	return s.resolveType(ty, s.pkg.Scope)
}
