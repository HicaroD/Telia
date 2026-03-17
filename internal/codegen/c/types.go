package c

import (
	"fmt"

	"github.com/HicaroD/Telia/internal/ast"
	"github.com/HicaroD/Telia/internal/lexer/token"
)

// CVariable represents a variable in the C backend.
// CType holds the C type string (e.g. "int32_t", "char *", "Point").
// Name holds the C variable name (e.g. "x", "_t0").
// This type is stored in VarIdStmt.BackendType and Param.BackendType,
// replacing the *Variable type used by the LLVM backend.
type CVariable struct {
	CType string
	Name  string
}

// emitCType converts a Telia ExprType into its C type string.
func emitCType(ty *ast.ExprType) string {
	switch ty.Kind {
	case ast.EXPR_TYPE_BASIC:
		b := ty.T.(*ast.BasicType)
		switch b.Kind {
		case token.BOOL_TYPE:
			return "_Bool"
		case token.INT_TYPE:
			return "int"
		case token.UINT_TYPE:
			return "unsigned int"
		case token.I8_TYPE:
			return "int8_t"
		case token.U8_TYPE:
			return "uint8_t"
		case token.I16_TYPE:
			return "int16_t"
		case token.U16_TYPE:
			return "uint16_t"
		case token.I32_TYPE:
			return "int32_t"
		case token.U32_TYPE:
			return "uint32_t"
		case token.I64_TYPE:
			return "int64_t"
		case token.U64_TYPE:
			return "uint64_t"
		case token.I128_TYPE:
			return "__int128"
		case token.U128_TYPE:
			return "unsigned __int128"
		case token.F32_TYPE, token.FLOAT_TYPE:
			return "float"
		case token.F64_TYPE:
			return "double"
		case token.VOID_TYPE:
			return "void"
		case token.STRING_TYPE, token.CSTRING_TYPE:
			return "char *"
		case token.RAWPTR_TYPE:
			return "void *"
		default:
			panic(fmt.Sprintf("emitCType: unhandled basic type kind: %v", b.Kind))
		}
	case ast.EXPR_TYPE_POINTER:
		ptr := ty.T.(*ast.PointerType)
		return emitCType(ptr.Type) + " *"
	case ast.EXPR_TYPE_STRUCT:
		st := ty.T.(*ast.StructType)
		return st.Decl.Name.Name()
	case ast.EXPR_TYPE_ALIAS:
		// Type aliases are fully resolved by sema before codegen runs.
		// This branch should never be reached.
		panic("emitCType: type alias should have been resolved by sema")
	case ast.EXPR_TYPE_TUPLE:
		// Tuple typedefs are generated separately by emitTupleTypedef (issue #73).
		panic("emitCType: tuple types must be handled by emitTupleTypedef")
	default:
		panic(fmt.Sprintf("emitCType: unhandled ExprType kind: %v", ty.Kind))
	}
}
