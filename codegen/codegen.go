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

type moduleCache struct {
	functions map[string]*cachedFunction
}

type cachedFunction struct {
	fn *llvm.Value
	ty *llvm.Type
}

func newModuleCache() *moduleCache {
	return &moduleCache{
		functions: map[string]*cachedFunction{},
	}
}

func (cache *moduleCache) InsertFunction(name string, fn *llvm.Value, ty *llvm.Type) {
	// TODO(errors)
	if _, ok := cache.functions[name]; ok {
		log.Fatalf("FUNCTION %s ALREADY DECLARED", name)
	}
	function := cachedFunction{
		fn: fn,
		ty: ty,
	}
	cache.functions[name] = &function
}

func (cache *moduleCache) GetFunction(name string) *cachedFunction {
	if fn, ok := cache.functions[name]; ok {
		return fn
	}
	return nil
}

type codegen struct {
	context     llvm.Context
	module      llvm.Module
	builder     llvm.Builder
	astNodes    []ast.AstNode
	moduleCache *moduleCache
}

func NewCodegen(astNodes []ast.AstNode) *codegen {
	context := llvm.NewContext()
	module := context.NewModule("tmpmod")
	builder := context.NewBuilder()

	return &codegen{
		context:     context,
		module:      module,
		builder:     builder,
		astNodes:    astNodes,
		moduleCache: newModuleCache(),
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

	filename := "telia.ll"
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

	codegen.moduleCache.InsertFunction(function.Name, &functionValue, &functionType)

	functionBody := codegen.context.AddBasicBlock(functionValue, "entry")
	codegen.builder.SetInsertPointAtEnd(functionBody)
	codegen.generateBlock(function.Block)
	return nil
}

func (codegen *codegen) generateBlock(block *ast.BlockStmt) {
	for i := range block.Statements {
		switch block.Statements[i].(type) {
		case *ast.FunctionCallStmt:
			calledFn := block.Statements[i].(*ast.FunctionCallStmt)
			function := codegen.moduleCache.GetFunction(calledFn.Name)
			// TODO(errors): function not found on cache
			if function == nil {
				log.Fatalf("FUNCTION %s not found on module", calledFn.Name)
			}
			args := codegen.getExprList(calledFn.Args)
			call := codegen.builder.CreateCall(*function.ty, *function.fn, args, "call")
			fmt.Printf("BUILDING CALL: %s\n", call)
		case *ast.ReturnStmt:
			fmt.Println("RETURN STATEMENT")
			ret := block.Statements[i].(*ast.ReturnStmt)
			returnValue := codegen.getExpr(ret.Value)
			generatedRet := codegen.builder.CreateRet(returnValue)
			fmt.Printf("BUILDING RETURN: %s\n", generatedRet)
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
		functionValue := llvm.AddFunction(codegen.module, external.Prototypes[i].Name, functionType)
		codegen.moduleCache.InsertFunction(external.Prototypes[i].Name, &functionValue, &functionType)
	}
	return nil
}

func (codegen *codegen) getType(ty ast.ExprType) llvm.Type {
	switch ty.(type) {
	case *ast.BasicType:
		basicType := ty.(*ast.BasicType)
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
	case *ast.PointerType:
		pointerType := ty.(*ast.PointerType)
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

func (codegen *codegen) getFieldListTypes(fields *ast.FieldList) []llvm.Type {
	types := make([]llvm.Type, len(fields.Fields))
	for i := range fields.Fields {
		types[i] = codegen.getType(fields.Fields[i].Type)
	}
	return types
}

func (codegen *codegen) getExprList(expressions []ast.Expr) []llvm.Value {
	values := make([]llvm.Value, len(expressions))
	for i := range expressions {
		values[i] = codegen.getExpr(expressions[i])
	}
	return values
}

func (codegen *codegen) getExpr(expr ast.Expr) llvm.Value {
	switch expr.(type) {
	case *ast.LiteralExpr:
		literalExpr := expr.(*ast.LiteralExpr)
		switch literalExpr.Kind {
		case kind.INTEGER_LITERAL:
			integerLiteral := uint64(literalExpr.Value.(int))
			return llvm.ConstInt(codegen.context.Int32Type(), integerLiteral, false)
		case kind.STRING_LITERAL:
			stringLiteral := literalExpr.Value.(string)
			globalStrPtr := codegen.builder.CreateGlobalStringPtr(stringLiteral, ".str")
			return globalStrPtr
		default:
			log.Fatalf("unimplemented literal expr: %s", expr)
		}
	default:
		log.Fatalf("unimplemented expr: %s", expr)
	}
	// NOTE: this line should be unreachable
	fmt.Println("REACHING AN UNREACHABLE LINE AT getExpr AT getExpr")
	return llvm.Value{}
}
