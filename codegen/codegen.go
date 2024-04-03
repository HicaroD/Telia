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
	cache    *cache
}

func New(astNodes []ast.AstNode) *codegen {
	context := llvm.NewContext()
	// TODO: properly define the module name
	module := context.NewModule("tmpmod")
	builder := context.NewBuilder()

	return &codegen{
		context:  context,
		module:   module,
		builder:  builder,
		astNodes: astNodes,
		cache:    newCache(),
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

	err := codegen.generateBitcodeFile()
	return err
}

func (codegen *codegen) generateBitcodeFile() error {
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

	fn := codegen.cache.InsertFunction(function.Name, &functionValue, &functionType)

	functionBody := codegen.context.AddBasicBlock(functionValue, "entry")
	codegen.builder.SetInsertPointAtEnd(functionBody)
	codegen.generateBlock(fn, function.Block)
	return nil
}

func (codegen *codegen) generateBlock(fn *function, stmts *ast.BlockStmt) {
	for i := range stmts.Statements {
		switch statement := stmts.Statements[i].(type) {
		case *ast.FunctionCallStmt:
			function := codegen.cache.GetFunction(statement.Name)
			// TODO(errors): function not found on cache
			if function == nil {
				log.Fatalf("FUNCTION %s not found on module", statement.Name)
			}
			args := codegen.getExprList(fn, statement.Args)
			codegen.builder.CreateCall(*function.ty, *function.fn, args, "call")
		case *ast.ReturnStmt:
			returnValue := codegen.getExpr(statement.Value, fn)
			codegen.builder.CreateRet(returnValue)
		case *ast.CondStmt:
			codegen.generateCondStmt(fn, statement)
		case *ast.VarDeclStmt:
			valueType := codegen.getType(statement.Type)
			allocaInst := codegen.builder.CreateAlloca(valueType, ".ptr")
			varExpr := codegen.getExpr(statement.Value, fn)
			codegen.builder.CreateStore(varExpr, allocaInst)
			err := fn.InsertLocal(statement.Name.Lexeme.(string), valueType, allocaInst)
			if err != nil {
				log.Fatal(err)
			}
		default:
			log.Fatalf("unimplemented block statement: %s", statement)
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
	codegen.cache.InsertFunction(prototype.Name, &functionValue, &functionType)
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
		// case kind.I128_TYPE:
		// 	return codegen.context.IntType(128)
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

func (codegen *codegen) getExprList(fn *function, expressions []ast.Expr) []llvm.Value {
	values := make([]llvm.Value, len(expressions))
	for i := range expressions {
		values[i] = codegen.getExpr(expressions[i], fn)
	}
	return values
}

func (codegen *codegen) getExpr(expr ast.Expr, fn *function) llvm.Value {
	switch currentExpr := expr.(type) {
	case *ast.LiteralExpr:
		switch currentExpr.Kind {
		case kind.INTEGER_LITERAL:
			integerLiteral := uint64(currentExpr.Value.(int))
			return llvm.ConstInt(codegen.context.Int32Type(), integerLiteral, false)
		case kind.STRING_LITERAL:
			stringLiteral := currentExpr.Value.(string)
			globalStrLiteral := codegen.cache.GetGlobal(stringLiteral)
			if globalStrLiteral != nil {
				return *globalStrLiteral
			}
			globalStrPtr := codegen.builder.CreateGlobalStringPtr(stringLiteral, ".str")
			codegen.cache.InsertGlobal(stringLiteral, &globalStrPtr)
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
	case *ast.IdExpr:
		// TODO: globals?
		varName := currentExpr.Name.Lexeme.(string)
		localVar := fn.GetLocal(varName)
		if localVar == nil {
			log.Fatalf("local not defined: %s", varName)
		}
		loadedVariable := codegen.builder.CreateLoad(localVar.ty, localVar.ptr, ".load")
		return loadedVariable
	default:
		log.Fatalf("unimplemented expr: %s", expr)
	}
	// NOTE: this line should be unreachable
	fmt.Println("REACHING AN UNREACHABLE LINE AT getExpr AT getExpr")
	return llvm.Value{}
}

func (codegen *codegen) generateCondStmt(fn *function, condStmt *ast.CondStmt) {
	ifBlock := llvm.AddBasicBlock(*fn.fn, ".if")
	elseBlock := llvm.AddBasicBlock(*fn.fn, ".else")
	// TODO(errors): if there are no instructions left to process, the end
	// block will be empty, and this will cause an error, I need to deal with
	// it
	endBlock := llvm.AddBasicBlock(*fn.fn, ".end")

	ifExpr := codegen.getExpr(condStmt.IfStmt.Expr, fn)
	codegen.builder.CreateCondBr(ifExpr, ifBlock, elseBlock)

	codegen.builder.SetInsertPointAtEnd(ifBlock)
	codegen.generateBlock(fn, condStmt.IfStmt.Block)
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
		codegen.generateBlock(fn, condStmt.ElseStmt.Block)
	}
	codegen.builder.CreateBr(endBlock)

	codegen.builder.SetInsertPointAtEnd(endBlock)
}
