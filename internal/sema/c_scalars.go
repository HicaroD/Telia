package sema

import (
	"github.com/HicaroD/Telia/internal/ast"
	"github.com/HicaroD/Telia/internal/ctarget"
	"github.com/HicaroD/Telia/internal/diagnostics"
	"github.com/HicaroD/Telia/internal/lexer/token"
)

func operatorTypeMatches(valid, actual *ast.ExprType) bool {
	if valid.Equals(actual) {
		return true
	}
	return valid.Kind == ast.EXPR_TYPE_BASIC && actual.Kind == ast.EXPR_TYPE_BASIC &&
		(actual.IsNumeric() || actual.T.(*ast.BasicType).Kind == token.CSTRING_TYPE) &&
		valid.T.(*ast.BasicType).Kind == actual.T.(*ast.BasicType).Kind
}

func (s *sema) cScalar(id *ast.IdType) (*ast.ExprType, error) {
	scalars, err := ctarget.Scalars()
	if err != nil {
		return nil, diagnostics.Diag{
			Title:   "C target unavailable",
			Message: err.Error(),
			Span:    id.Name.Span,
			Help:    "Install a supported native GCC or Clang toolchain and ensure cc is on PATH.",
		}
	}
	spec, ok := scalars[id.Name.Name()]
	if !ok {
		return nil, diagnostics.Diag{
			Title:   "Unknown C scalar",
			Message: "unknown compiler C scalar intrinsic",
			Span:    id.Name.Span,
			Help:    "Use the public C scalar aliases in std::ffi::c.",
		}
	}
	kinds := map[int]token.Kind{
		1:  token.I8_TYPE,
		2:  token.I16_TYPE,
		4:  token.I32_TYPE,
		8:  token.I64_TYPE,
		16: token.I128_TYPE,
	}
	if spec.Unsigned {
		kinds = map[int]token.Kind{
			1:  token.U8_TYPE,
			2:  token.U16_TYPE,
			4:  token.U32_TYPE,
			8:  token.U64_TYPE,
			16: token.U128_TYPE,
		}
	}
	kind := kinds[spec.Bytes]
	if spec.Float {
		if spec.Bytes == 4 {
			kind = token.F32_TYPE
		} else {
			kind = token.F64_TYPE
		}
	}
	ty := concreteBasicType(kind)
	ty.T.(*ast.BasicType).CName = spec.Name
	s.program.CScalars = append(s.program.CScalars, ty)
	return ty, nil
}
