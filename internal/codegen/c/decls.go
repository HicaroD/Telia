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
	"_telia_println": true,
	"printf":         true,
	"fprintf":        true,
	"sprintf":        true,
	"scanf":          true,
	"puts":           true,
	"fputs":          true,
	"fgets":          true,
	"fopen":          true,
	"fclose":         true,
	"fread":          true,
	"fwrite":         true,
	"malloc":         true,
	"calloc":         true,
	"realloc":        true,
	"free":           true,
	"exit":           true,
	"abort":          true,
	"memcpy":         true,
	"memmove":        true,
	"memset":         true,
	"strlen":         true,
	"strcmp":         true,
	"strcpy":         true,
	"strcat":         true,
	"strdup":         true,
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
		// Header-owned names can be function-like macros (SDL redirects several
		// libc functions this way). A prototype would expand the macro with types
		// as arguments. Let the header own those declarations and call expansions.
		if ext.Attributes.Header != "" {
			fmt.Fprintf(c.buf, "#ifndef %s\n", cName)
		}
		c.buf.WriteString(
			fmt.Sprintf(
				"%s %s%s(%s);\n",
				retType,
				callingConvention(proto.Attributes),
				cName,
				params,
			),
		)
		if ext.Attributes.Header != "" {
			c.buf.WriteString("#endif\n")
		}
	}
}

// emitProtoParams renders the C parameter list for a prototype, including
// variadic handling and pointer qualifiers.
func (c *CCodegen) emitProtoParams(proto *ast.Proto) string {
	var parts []string
	for _, param := range proto.Params.Fields {
		if param.Variadic {
			// @c variadic param — emit as C "..."
			parts = append(parts, "...")
			continue
		}
		cType := emitCType(param.Type)
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
	c.buf.WriteString(fmt.Sprintf("struct %s {\n", name))
	for i, elem := range tt.Types {
		c.buf.WriteString(fmt.Sprintf("    %s _%d;\n", emitCType(elem), i))
	}
	c.buf.WriteString("};\n\n")
}

// emitFnForwardDecl emits a C forward declaration for a Telia function.
func (c *CCodegen) emitFnForwardDecl(fn *ast.FnDecl) {
	retType := emitCType(fn.RetType)
	name := c.mangledName(fn)
	params := c.emitFnParams(fn)
	c.buf.WriteString(
		fmt.Sprintf("%s %s%s(%s);\n", retType, callingConvention(fn.Attributes), name, params),
	)
}

// emitFnParams renders the C parameter list for a function declaration.
func (c *CCodegen) emitFnParams(fn *ast.FnDecl) string {
	var parts []string
	for _, param := range fn.Params.Fields {
		if param.Variadic && param.Attributes != nil && param.Attributes.C {
			parts = append(parts, "...")
			continue
		}
		cType := emitCType(param.Type)
		cVar := &CVariable{CType: cType, Name: c.nextTmp() + "_" + param.Name.Name()}
		param.BackendType = cVar
		parts = append(parts, fmt.Sprintf("%s %s", cType, cVar.Name))
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
	if st.CName != "" {
		name = st.CName
	}
	c.buf.WriteString(fmt.Sprintf("struct %s {\n", name))
	for _, field := range st.Fields {
		cType := emitCType(field.Type)
		c.buf.WriteString(fmt.Sprintf("    %s %s;\n", cType, field.Name.Name()))
	}
	if len(st.Fields) == 0 {
		c.buf.WriteString("    unsigned char _empty;\n")
	}
	c.buf.WriteString("};\n\n")
}

// emitFnBody emits the full C function definition including body.
func (c *CCodegen) emitFnBody(fn *ast.FnDecl) {
	retType := emitCType(fn.RetType)
	name := c.mangledName(fn)
	params := c.emitFnParams(fn)
	c.buf.WriteString(
		fmt.Sprintf("%s %s%s(%s) {\n", retType, callingConvention(fn.Attributes), name, params),
	)
	prevRetTy := c.currentRetTy
	c.currentRetTy = fn.RetType
	c.emitBlock(fn.Block, "    ")
	c.currentRetTy = prevRetTy
	c.buf.WriteString("}\n\n")
}
