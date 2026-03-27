package c

import (
	"fmt"
	"strings"

	"github.com/HicaroD/Telia/internal/ast"
)

// isLLVMIntrinsic reports whether a C name is an LLVM intrinsic (contains a dot).
// These are not valid C identifiers and must be silently dropped by the C backend.
func isLLVMIntrinsic(name string) bool {
	return strings.Contains(name, ".")
}

// knownLibcSymbols are symbols already declared by the standard #include headers
// in the preamble, or LLVM intrinsics that must be silently dropped for the C backend.
// Re-declaring known libc symbols causes duplicate-declaration compile errors;
// LLVM intrinsics (llvm.*) are not valid C identifiers and must not be emitted.
var knownLibcSymbols = map[string]bool{
	"printf":  true,
	"fprintf": true,
	"sprintf": true,
	"scanf":   true,
	"puts":    true,
	"fputs":   true,
	"fgets":   true,
	"fopen":   true,
	"fclose":  true,
	"fread":   true,
	"fwrite":  true,
	"malloc":  true,
	"calloc":  true,
	"realloc": true,
	"free":    true,
	"exit":    true,
	"abort":   true,
	"memcpy":  true,
	"memmove": true,
	"memset":  true,
	"strlen":  true,
	"strcmp":  true,
	"strcpy":  true,
	"strcat":  true,
	"strdup":  true,
}

// emitExternDecl emits C forward declarations for all prototypes in an extern block.
func (c *CCodegen) emitExternDecl(ext *ast.ExternDecl) {
	for _, proto := range ext.Prototypes {
		cName := proto.Name.Name()
		if proto.Attributes.LinkName != "" {
			cName = proto.Attributes.LinkName
		}

		// Skip LLVM intrinsics (contain dots — not valid C identifiers).
		if isLLVMIntrinsic(cName) {
			continue
		}
		// Skip symbols already provided by the preamble headers.
		if knownLibcSymbols[cName] {
			continue
		}

		retType := emitCType(proto.RetType)
		params := c.emitProtoParams(proto)
		c.buf.WriteString(fmt.Sprintf("%s %s(%s);\n", retType, cName, params))
	}
}

// emitProtoParams renders the C parameter list for a prototype, including
// variadic handling and @const qualifiers.
func (c *CCodegen) emitProtoParams(proto *ast.Proto) string {
	var parts []string
	for _, param := range proto.Params.Fields {
		if param.Variadic {
			// @c variadic param — emit as C "..."
			parts = append(parts, "...")
			continue
		}
		cType := emitCType(param.Type)
		if param.Attributes != nil && param.Attributes.Const {
			cType = "const " + cType
		}
		parts = append(parts, cType)
	}
	if proto.Params.IsVariadic && len(parts) > 0 {
		// Ensure "..." is last if it wasn't already added via @c.
		if parts[len(parts)-1] != "..." {
			parts = append(parts, "...")
		}
	}
	return strings.Join(parts, ", ")
}

// emitTupleTypedef emits the C typedef struct for one unique tuple shape.
//
//	typedef struct { int32_t _0; int64_t _1; } _Tuple_int32_t_int64_t;
//
// Callers are responsible for deduplication; this function always emits.
func (c *CCodegen) emitTupleTypedef(ty *ast.ExprType) {
	name := tupleTypedefName(ty)
	tt := ty.T.(*ast.TupleType)
	c.buf.WriteString("typedef struct {\n")
	for i, elem := range tt.Types {
		c.buf.WriteString(fmt.Sprintf("    %s _%d;\n", emitCType(elem), i))
	}
	c.buf.WriteString(fmt.Sprintf("} %s;\n\n", name))
}

// emitFnForwardDecl emits a C forward declaration for a Telia function.
func (c *CCodegen) emitFnForwardDecl(fn *ast.FnDecl) {
	retType := emitCType(fn.RetType)
	name := c.mangledName(fn)
	params := c.emitFnParams(fn)
	c.buf.WriteString(fmt.Sprintf("%s %s(%s);\n", retType, name, params))
}

// emitFnParams renders the C parameter list for a function declaration.
func (c *CCodegen) emitFnParams(fn *ast.FnDecl) string {
	var parts []string
	for _, param := range fn.Params.Fields {
		cType := emitCType(param.Type)
		cVar := &CVariable{CType: cType, Name: param.Name.Name()}
		param.BackendType = cVar
		parts = append(parts, fmt.Sprintf("%s %s", cType, param.Name.Name()))
	}
	return strings.Join(parts, ", ")
}

// emitStructDecl emits a C typedef struct definition for a Telia struct declaration.
// It must be emitted before any function forward declarations to avoid forward
// reference errors in the generated C file.
//
//	typedef struct { field_type field_name; ... } StructName;
func (c *CCodegen) emitStructDecl(st *ast.StructDecl) {
	name := st.Name.Name()
	c.buf.WriteString(fmt.Sprintf("typedef struct %s {\n", name))
	for _, field := range st.Fields {
		cType := emitCType(field.Type)
		c.buf.WriteString(fmt.Sprintf("    %s %s;\n", cType, field.Name.Name()))
	}
	c.buf.WriteString(fmt.Sprintf("} %s;\n\n", name))
}

// emitFnBody emits the full C function definition including body.
func (c *CCodegen) emitFnBody(fn *ast.FnDecl) {
	retType := emitCType(fn.RetType)
	name := c.mangledName(fn)
	params := c.emitFnParams(fn)
	c.buf.WriteString(fmt.Sprintf("%s %s(%s) {\n", retType, name, params))
	prevRetTy := c.currentRetTy
	c.currentRetTy = fn.RetType
	c.emitBlock(fn.Block, "    ")
	c.currentRetTy = prevRetTy
	c.buf.WriteString("}\n\n")
}
