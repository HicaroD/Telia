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

func New(astNodes []ast.AstNode) *codegen {
	context := llvm.NewContext()
	// TODO: properly define the module name
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
		switch astNode := codegen.astNodes[i].(type) {
		case *ast.FunctionDecl:
			err := codegen.generateFnDecl(astNode)
			if err != nil {
				return err
			}
		case *ast.ExternDecl:
			codegen.generateExternDecl(astNode)
		default:
			log.Fatalf("unimplemented: %s\n", reflect.TypeOf(codegen.astNodes[i]))
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
	codegen.generateBlock(functionValue, function.Block)
	return nil
}

func (codegen *codegen) generateBlock(function llvm.Value, stmts *ast.BlockStmt) {
	for i := range stmts.Statements {
		switch statement := stmts.Statements[i].(type) {
		case *ast.FunctionCallStmt:
			function := codegen.moduleCache.GetFunction(statement.Name)
			// TODO(errors): function not found on cache
			if function == nil {
				log.Fatalf("FUNCTION %s not found on module", statement.Name)
			}
			args := codegen.getExprList(statement.Args)
			codegen.builder.CreateCall(*function.ty, *function.fn, args, "call")
		case *ast.ReturnStmt:
			returnValue := codegen.getExpr(statement.Value)
			codegen.builder.CreateRet(returnValue)
		case *ast.CondStmt:
			codegen.generateCondStmt(function, statement)
		default:
			log.Fatalf("unimplemented block statement: %s", reflect.TypeOf(stmts.Statements[i]))
		}
	}
}

func (codegen *codegen) generateExternDecl(external *ast.ExternDecl) {
	for i := range external.Prototypes {
		codegen.generatePrototype(external.Prototypes[i])
	}
}

func (codegen *codegen) generatePrototype(prototype *ast.Proto) {
	returnType := codegen.getType(prototype.RetType)
	paramsTypes := codegen.getFieldListTypes(prototype.Params)
	functionType := llvm.FunctionType(returnType, paramsTypes, prototype.Params.IsVariadic)
	functionValue := llvm.AddFunction(codegen.module, prototype.Name, functionType)
	codegen.moduleCache.InsertFunction(prototype.Name, &functionValue, &functionType)
}

func (codegen *codegen) getType(ty ast.ExprType) llvm.Type {
	switch exprTy := ty.(type) {
	case *ast.BasicType:
		switch exprTy.Kind {
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
			log.Fatalf("invalid basic type token: '%s'", exprTy.Kind)
		}
	case *ast.PointerType:
		underlyingExprType := codegen.getType(exprTy.Type)
		// TODO: learn about how to properly define a pointer address space
		return llvm.PointerType(underlyingExprType, 0)
	case nil:
		return codegen.context.VoidType()
	default:
		log.Fatalf("invalid type: %s", reflect.TypeOf(exprTy))
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
	switch currentExpr := expr.(type) {
	case *ast.LiteralExpr:
		switch currentExpr.Kind {
		case kind.INTEGER_LITERAL:
			integerLiteral := uint64(currentExpr.Value.(int))
			return llvm.ConstInt(codegen.context.Int32Type(), integerLiteral, false)
		case kind.STRING_LITERAL:
			stringLiteral := currentExpr.Value.(string)
			globalStrPtr := codegen.builder.CreateGlobalStringPtr(stringLiteral, ".str")
			return globalStrPtr
		case kind.TRUE_BOOL_LITERAL:
			trueBoolLiteral := llvm.ConstInt(codegen.context.Int1Type(), 1, false)
			return trueBoolLiteral
		case kind.FALSE_BOOL_LITERAL:
			falseBoolLiteral := llvm.ConstInt(codegen.context.Int1Type(), 0, false)
			return falseBoolLiteral
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

func (codegen *codegen) generateCondStmt(function llvm.Value, condStmt *ast.CondStmt) {
	ifBlock := llvm.AddBasicBlock(function, ".if")
	elseBlock := llvm.AddBasicBlock(function, ".else")
	endBlock := llvm.AddBasicBlock(function, ".end")

	ifExpr := codegen.getExpr(condStmt.IfStmt.Expr)
	codegen.builder.CreateCondBr(ifExpr, ifBlock, elseBlock)

	codegen.builder.SetInsertPointAtEnd(ifBlock)
	codegen.generateBlock(function, condStmt.IfStmt.Block)
	codegen.builder.CreateBr(endBlock)

	// TODO: implement elif
	/*
			Lembre-se de que elif são traduzidos para um else e dentro terá if:

			if condition {
		    	// A
			}
			elif condition {
		    	// B
			}
			else {
		    	// C
			}

			Será traduzido para isso:

			if condition {
		    	// A
			}
			else {
		    	if condition {
		        	// B
		    	}
		    	else {
		        	// C
		    	}
			}

			Na pasta de "Prototypes", em "learn_c", tem um exemplo de código em C
			e o código LLVM. Lá eu posso me inspirar para aprender como else if
			funcionam de fato, mas eu já tenho uma ideia.
	*/

	codegen.builder.SetInsertPointAtEnd(elseBlock)
	if condStmt.ElseStmt != nil {
		codegen.generateBlock(function, condStmt.ElseStmt.Block)
	}
	codegen.builder.CreateBr(endBlock)

	codegen.builder.SetInsertPointAtEnd(endBlock)
}
