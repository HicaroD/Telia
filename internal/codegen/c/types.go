package c

import (
	"crypto/sha256"
	"fmt"
	"strings"

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
	case ast.EXPR_TYPE_ERROR_KIND:
		return "const void *"
	case ast.EXPR_TYPE_VA_LIST_OWNER:
		return "_TeliaVaList"
	case ast.EXPR_TYPE_C_VA_LIST:
		return "va_list"
	case ast.EXPR_TYPE_FUNCTION:
		return ty.T.(*ast.FunctionType).CName
	case ast.EXPR_TYPE_OPAQUE:
		d := ty.T.(*ast.OpaqueType).Decl
		if d.Attributes.LinkName != "" {
			return d.Attributes.LinkName
		}
		return d.Name.Name()
	case ast.EXPR_TYPE_ENUM:
		d := ty.T.(*ast.EnumType).Decl
		if d.Foreign != nil {
			return d.CName
		}
		if d.HasPayload {
			return d.CName
		}
		return emitCType(d.Type)
	case ast.EXPR_TYPE_ARRAY, ast.EXPR_TYPE_SLICE:
		t := ty.T.(*ast.SequenceType)
		return fmt.Sprintf(
			"telia_seq_%x",
			sha256.Sum256([]byte(fmt.Sprintf("%d:%d:%s", ty.Kind, t.Length, emitCType(t.Element)))),
		)
	case ast.EXPR_TYPE_BASIC:
		b := ty.T.(*ast.BasicType)
		if b.CName != "" {
			return b.CName
		}
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
		case token.STRING_TYPE:
			return "_String"
		case token.CSTRING_TYPE:
			if b.Const {
				return "const char *"
			}
			return "char *"
		case token.RAWPTR_TYPE:
			return "void *"
		case token.ERROR_TYPE:
			return "_Error"
		default:
			panic(fmt.Sprintf("emitCType: unhandled basic type kind: %v", b.Kind))
		}
	case ast.EXPR_TYPE_POINTER:
		ptr := ty.T.(*ast.PointerType)
		inner := emitCType(ptr.Type)
		if ptr.Const {
			inner = constCType(ptr.Type)
		}
		return inner + " *"
	case ast.EXPR_TYPE_STRUCT:
		st := ty.T.(*ast.StructType)
		if st.Decl.CName != "" {
			return st.Decl.CName
		}
		return st.Decl.Name.Name()
	case ast.EXPR_TYPE_ALIAS:
		// Type aliases are fully resolved by sema before codegen runs.
		// This branch should never be reached.
		panic("emitCType: type alias should have been resolved by sema")
	case ast.EXPR_TYPE_TUPLE:
		return tupleTypedefName(ty)
	default:
		panic(fmt.Sprintf("emitCType: unhandled ExprType kind: %v", ty.Kind))
	}
}

func constCType(ty *ast.ExprType) string {
	name := emitCType(ty)
	pointerLike := ty.IsPointer()
	if ty.Kind == ast.EXPR_TYPE_BASIC {
		kind := ty.T.(*ast.BasicType).Kind
		pointerLike = kind == token.CSTRING_TYPE || kind == token.RAWPTR_TYPE
	}
	if pointerLike {
		return name + " const"
	}
	return "const " + name
}

// tupleTypedefName returns the C typedef name for a tuple type.
// The name is derived from the C types of each element, sanitized for use
// as a C identifier, e.g. (i32, i64) → "_Tuple_int32_t_int64_t".
// Must match the key used by sema.tupleKey for consistent deduplication.
func tupleTypedefName(ty *ast.ExprType) string {
	tt := ty.T.(*ast.TupleType)
	parts := make([]string, len(tt.Types))
	for i, elem := range tt.Types {
		ctype := emitCType(elem)
		// Sanitize: spaces and * are not valid in C identifiers.
		safe := strings.NewReplacer(" ", "_", "*", "ptr").Replace(ctype)
		parts[i] = safe
	}
	return "_Tuple_" + strings.Join(parts, "_")
}
