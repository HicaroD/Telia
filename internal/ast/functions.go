package ast

import "fmt"

type FunctionType struct {
	Params     *Params
	Result     *ExprType
	CName      string
	Attributes Attributes
}

func (f *FunctionType) String() string { return fmt.Sprintf("fn(...) %s", f.Result.T) }

func (f *FunctionType) HasForeignABI() bool {
	return f.Attributes.CallingConvention != "" || f.Attributes.CallingConventionMacro != ""
}

func (f *FunctionType) Equals(g *FunctionType) bool {
	if f.Attributes.CallingConvention != g.Attributes.CallingConvention ||
		f.Attributes.CallingConventionMacro != g.Attributes.CallingConventionMacro {
		return false
	}
	if f.Params.IsVariadic != g.Params.IsVariadic || len(f.Params.Fields) != len(g.Params.Fields) ||
		!f.Result.Equals(g.Result) {
		return false
	}
	for i, p := range f.Params.Fields {
		if !p.Type.Equals(g.Params.Fields[i].Type) {
			return false
		}
	}
	return true
}
