package sema

import (
	"github.com/HicaroD/Telia/internal/ast"
	"github.com/HicaroD/Telia/internal/diagnostics"
)

func (s *sema) functionValue(
	id *ast.IdExpr,
	symbol *ast.Node,
	context ...*ast.ExprType,
) (*ast.ExprType, error) {
	var params *ast.Params
	var result *ast.ExprType
	var attrs ast.Attributes
	switch d := symbol.Node.(type) {
	case *ast.FnDecl:
		if len(d.TypeParams) == 0 && len(id.TypeArgs) != 0 {
			return nil, diagnostics.Diag{Message: "function does not take type arguments", Span: id.Name.Span}
		}
		if len(d.TypeParams) != 0 {
			template := d
			i, err := s.newInference(d.TypeParams, id.TypeArgs, s.pkg.Scope)
			if err != nil {
				return nil, err
			}
			i.scope, i.file = d.Owner.Scope, d.File
			if len(context) != 0 {
				pattern := &ast.ExprType{Kind: ast.EXPR_TYPE_FUNCTION, T: &ast.FunctionType{Params: d.Params, Result: d.RetType, Attributes: d.Attributes}}
				if err := i.match(pattern, context[0]); err != nil {
					return nil, err
				}
			}
			args, err := i.result()
			if err != nil {
				return nil, diagnostics.Diag{Title: "Generic function value", Message: err.Error(), Span: id.Name.Span, Help: "Supply concrete type arguments or an expected callback signature."}
			}
			d, err = s.instantiateFunction(d, args, s.pkg.Scope, id.Name)
			if err != nil {
				return nil, err
			}
			id.N = &ast.Node{Kind: ast.KIND_FN_DECL, Node: d, Span: id.Name.Span}
			s.suggestTypeArguments(id.Name, id.TypeArgs, args, template.TypeParams, template.Owner, template.File, s.pkg.Scope, func(i *typeInference) error {
				if len(context) == 0 || context[0] == nil {
					return nil
				}
				return i.match(&ast.ExprType{Kind: ast.EXPR_TYPE_FUNCTION, T: &ast.FunctionType{Params: template.Params, Result: template.RetType, Attributes: template.Attributes}}, context[0])
			})
		}
		if err := s.resolveSignature(d); err != nil {
			return nil, err
		}
		params, result = d.Params, d.RetType
		attrs = d.Attributes
	case *ast.Proto:
		if len(id.TypeArgs) != 0 {
			return nil, diagnostics.Diag{Message: "foreign function does not take type arguments", Span: id.Name.Span}
		}
		params, result = d.Params, d.RetType
		attrs = d.Attributes
		// Prototypes without an explicit convention still belong to the foreign
		// boundary, unlike ordinary Telia function declarations.
		if attrs.CallingConvention == "" && attrs.CallingConventionMacro == "" {
			attrs.CallingConvention = "c"
		}
	}
	ty := &ast.ExprType{
		Kind:     ast.EXPR_TYPE_FUNCTION,
		Explicit: true,
		T:        &ast.FunctionType{Params: params, Result: result, Attributes: attrs},
	}
	return s.resolveType(ty, s.pkg.Scope)
}
