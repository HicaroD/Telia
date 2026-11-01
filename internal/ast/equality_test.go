package ast_test

import (
	"strings"
	"testing"

	"github.com/HicaroD/Telia/internal/ast"
	"github.com/HicaroD/Telia/internal/lexer/token"
)

func TestEqualityAndKeyEligibilityTypeFamilies(t *testing.T) {
	for _, kind := range []token.Kind{
		token.BOOL_TYPE, token.INT_TYPE, token.UINT_TYPE,
		token.I8_TYPE, token.I16_TYPE, token.I32_TYPE, token.I64_TYPE, token.I128_TYPE,
		token.U8_TYPE, token.U16_TYPE, token.U32_TYPE, token.U64_TYPE, token.U128_TYPE,
		token.FLOAT_TYPE, token.F32_TYPE, token.F64_TYPE, token.STRING_TYPE, token.CSTRING_TYPE, token.RAWPTR_TYPE,
	} {
		for _, key := range []bool{false, true} {
			if err := ast.EqualityEligibility(ast.NewBasicType(kind), key); err != nil {
				t.Errorf("%s (key=%t): %v", kind, key, err)
			}
		}
	}
	callback := &ast.ExprType{Kind: ast.EXPR_TYPE_FUNCTION, T: &ast.FunctionType{}}
	for _, tc := range []struct {
		name          string
		ty            *ast.ExprType
		equality, key bool
	}{
		{"function", callback, true, false},
		{"C scalar", &ast.ExprType{Kind: ast.EXPR_TYPE_BASIC, T: &ast.BasicType{Kind: token.U64_TYPE, CName: "size_t"}}, true, true},
		{"const C string", &ast.ExprType{Kind: ast.EXPR_TYPE_BASIC, T: &ast.BasicType{Kind: token.CSTRING_TYPE, Const: true}}, true, true},
		{"error kind", ast.ErrorKind(), true, true},
		{"error", ast.NewBasicType(token.ERROR_TYPE), false, false},
		{"void", ast.NewBasicType(token.VOID_TYPE), false, false},
		{"opaque", &ast.ExprType{Kind: ast.EXPR_TYPE_OPAQUE, T: &ast.OpaqueType{Decl: &ast.OpaqueDecl{Name: name("Opaque")}}}, false, false},
		{"borrowed va_list", &ast.ExprType{Kind: ast.EXPR_TYPE_C_VA_LIST, T: &ast.CVaListType{}}, false, false},
		{"owning va_list", &ast.ExprType{Kind: ast.EXPR_TYPE_VA_LIST_OWNER, T: &ast.VaListOwnerType{}}, false, false},
		{"slice", &ast.ExprType{Kind: ast.EXPR_TYPE_SLICE, T: &ast.SequenceType{Element: ast.NewBasicType(token.I32_TYPE)}}, false, false},
		{"pointer to error", ast.NewPointerType(ast.NewBasicType(token.ERROR_TYPE)), true, true},
		{"bare nil", ast.NewPointerType(ast.NewBasicType(token.UNTYPED_NULLPTR)), false, false},
		{"alias", &ast.ExprType{Kind: ast.EXPR_TYPE_ALIAS, T: &ast.TypeAlias{Type: ast.NewBasicType(token.STRING_TYPE)}}, true, true},
		{"tuple callback", &ast.ExprType{Kind: ast.EXPR_TYPE_TUPLE, T: &ast.TupleType{Types: []*ast.ExprType{ast.NewBasicType(token.I32_TYPE), callback}}}, true, false},
		{"empty enum", &ast.ExprType{Kind: ast.EXPR_TYPE_ENUM, T: &ast.EnumType{Decl: &ast.EnumDecl{Name: name("Empty")}}}, true, true},
	} {
		t.Run(tc.name, func(t *testing.T) {
			if got := ast.EqualityEligibility(tc.ty, false); (got == nil) != tc.equality {
				t.Errorf("equality eligibility: %v", got)
			}
			if got := ast.EqualityEligibility(tc.ty, true); (got == nil) != tc.key {
				t.Errorf("key eligibility: %v", got)
			}
		})
	}
}

func TestKeyEligibilityReportsNestedCallbackAndStopsAtPointers(t *testing.T) {
	callback := &ast.ExprType{Kind: ast.EXPR_TYPE_FUNCTION, T: &ast.FunctionType{}}
	record := &ast.StructDecl{
		Name:   name("Record"),
		Fields: []*ast.StructField{{Name: name("callback"), Type: callback}},
	}
	ty := &ast.ExprType{
		Kind: ast.EXPR_TYPE_ARRAY,
		T: &ast.SequenceType{
			Length:  0,
			Element: &ast.ExprType{Kind: ast.EXPR_TYPE_STRUCT, T: &ast.StructType{Decl: record}},
		},
	}
	if err := ast.EqualityEligibility(ty, true); err == nil ||
		!strings.Contains(err.Error(), "value[element].callback") {
		t.Fatalf("want callback component path even for zero length array, got %v", err)
	}
	if err := ast.EqualityEligibility(ast.NewPointerType(ty), true); err != nil {
		t.Fatal(err)
	}
}

func name(s string) *token.Token { return &token.Token{Kind: token.ID, Lexeme: []byte(s)} }
