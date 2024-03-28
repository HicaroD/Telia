package codegen

import (
	"fmt"
	"log"
	"os"
	"reflect"

	"github.com/HicaroD/telia-lang/ast"
	"github.com/HicaroD/telia-lang/lexer/token/kind"
	"tinygo.org/x/go-llvm"
)

type codegen struct {
	context  llvm.Context
	module   llvm.Module
	builder  llvm.Builder
	astNodes []ast.AstNode
}

func NewCodegen(astNodes []ast.AstNode) *codegen {
	context := llvm.NewContext()
	module := context.NewModule("tmpmod")
	builder := context.NewBuilder()

	return &codegen{
		context:  context,
		module:   module,
		builder:  builder,
		astNodes: astNodes,
	}
}

func (codegen *codegen) Generate() error {
	for i := range codegen.astNodes {
		switch codegen.astNodes[i].(type) {
		case *ast.FunctionDecl:
			fnDecl := codegen.astNodes[i].(*ast.FunctionDecl)
			err := codegen.generateFnDecl(fnDecl)
			if err != nil {
				return err
			}
		case *ast.ExternDecl:
			externDecl := codegen.astNodes[i].(*ast.ExternDecl)
			err := codegen.generateExternDecl(externDecl)
			if err != nil {
				return err
			}
		default:
			fmt.Printf("Generating...: %s\n", reflect.TypeOf(codegen.astNodes[i]))
		}
	}

	filename := "llvm.ll"
	file, err := os.Create(filename)
	if err != nil {
		return err
	}

	module := codegen.module.String()
	_, err = file.Write([]byte(module))
	if err != nil {
		return err
	}
	fmt.Printf("'%s' file generated successfuly\n", filename)
	return nil
}

func (codegen *codegen) generateFnDecl(function *ast.FunctionDecl) error {
	returnType := codegen.getType(function.RetType)
	paramsTypes := codegen.getFieldListTypes(function.Params)
	functionType := llvm.FunctionType(returnType, paramsTypes, function.Params.IsVariadic)

	functionValue := llvm.AddFunction(codegen.module, function.Name, functionType)
	functionBody := codegen.context.AddBasicBlock(functionValue, "entry")

	codegen.builder.SetInsertPointAtEnd(functionBody)

	codegen.generateBlock(function.Block)
	return nil
}

func (codegen *codegen) generateBlock(block *ast.BlockStmt) {
	for i := range block.Statements {
		switch block.Statements[i].(type) {
		// TODO: build function call stmt
		case ast.FunctionCallStmt:
			fmt.Println("FUNCTION CALL STATEMENT")
		// TODO: build return stmt
		case ast.ReturnStmt:
			fmt.Println("RETURN STATEMENT")
		default:
			log.Fatalf("unimplemented block statement: %s", reflect.TypeOf(block.Statements[i]))
		}
	}
}

// TODO: generate extern blocks
func (codegen *codegen) generateExternDecl(external *ast.ExternDecl) error {
	for i := range external.Prototypes {
		returnType := codegen.getType(external.Prototypes[i].RetType)
		paramsTypes := codegen.getFieldListTypes(external.Prototypes[i].Params)
		functionType := llvm.FunctionType(returnType, paramsTypes, external.Prototypes[i].Params.IsVariadic)
		llvm.AddFunction(codegen.module, external.Prototypes[i].Name, functionType)
	}
	return nil
}

func (codegen *codegen) getType(ty ast.ExprType) llvm.Type {
	switch ty.(type) {
	case ast.BasicType:
		basicType := ty.(ast.BasicType)
		switch basicType.Kind {
		case kind.BOOL_TYPE:
			return codegen.context.Int1Type()
		case kind.I8_TYPE:
			return codegen.context.Int8Type()
		case kind.I16_TYPE:
			return codegen.context.Int16Type()
		case kind.I32_TYPE:
			return codegen.context.Int32Type()
		case kind.I64_TYPE:
			return codegen.context.Int64Type()
		case kind.I128_TYPE:
			return codegen.context.IntType(128)
		default:
			log.Fatalf("invalid basic type token: '%s'", basicType.Kind)
		}
	case ast.PointerType:
		pointerType := ty.(ast.PointerType)
		underlyingExprType := codegen.getType(pointerType.Type)
		// TODO: learn about how to properly define a pointer address space
		return llvm.PointerType(underlyingExprType, 0)
	case nil:
		return codegen.context.VoidType()
	default:
		log.Fatalf("invalid type: %s", reflect.TypeOf(ty))
	}
	// NOTE: this line should be unreachable
	return codegen.context.VoidType()
}

func (codegen *codegen) getFieldListTypes(fieldList *ast.FieldList) []llvm.Type {
	types := make([]llvm.Type, len(fieldList.Fields))
	for i := range fieldList.Fields {
		types[i] = codegen.getType(fieldList.Fields[i].Type)
	}
	return types
}
