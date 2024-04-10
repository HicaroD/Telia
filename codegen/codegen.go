package codegen

import (
	"fmt"
	"log"
	"os"
	"reflect"

	"github.com/HicaroD/telia-lang/ast"
	"github.com/HicaroD/telia-lang/codegen/values"
	"github.com/HicaroD/telia-lang/lexer/token/kind"
	"github.com/HicaroD/telia-lang/scope"
	"tinygo.org/x/go-llvm"
)

type codegen struct {
	universe          *scope.Scope[values.LLVMValue]
	globalStrLiterals map[string]llvm.Value
	context           llvm.Context
	module            llvm.Module
	builder           llvm.Builder
	astNodes          []ast.AstNode
}

func New(astNodes []ast.AstNode) *codegen {
	// parent of universe scope is nil
	var nilScope *scope.Scope[values.LLVMValue] = nil
	universe := scope.New(nilScope)

	context := llvm.NewContext()
	// TODO: properly define the module name
	module := context.NewModule("tmpmod")
	builder := context.NewBuilder()

	return &codegen{
		universe:          universe,
		context:           context,
		module:            module,
		builder:           builder,
		astNodes:          astNodes,
		globalStrLiterals: map[string]llvm.Value{},
	}
}

func (codegen *codegen) Generate() error {
	for i := range codegen.astNodes {
		switch astNode := codegen.astNodes[i].(type) {
		case *ast.FunctionDecl:
			err := codegen.generateFnDecl(astNode)
			// TODO(errors)
			if err != nil {
				return err
			}
		case *ast.ExternDecl:
			err := codegen.generateExternDecl(astNode)
			// TODO(errors)
			if err != nil {
				return err
			}
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
	// TODO(errors)
	if err != nil {
		return err
	}

	module := codegen.module.String()
	_, err = file.Write([]byte(module))
	// TODO(errors)
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
	functionBlock := codegen.context.AddBasicBlock(functionValue, "entry")
	fn := values.NewFunctionValue(functionValue, functionType, &functionBlock)

	err := codegen.universe.Insert(function.Name, fn)
	// TODO(errors)
	if err != nil {
		return err
	}

	codegen.builder.SetInsertPointAtEnd(functionBlock)
	err = codegen.generateBlock(codegen.universe, fn, function.Block)
	// TODO(errors)
	if err != nil {
		return err
	}

	return err
}

func (codegen *codegen) generateBlock(parentScope *scope.Scope[values.LLVMValue], function values.Function, stmts *ast.BlockStmt) error {
	currentBlockScope := scope.New(parentScope)

	for i := range stmts.Statements {
		switch statement := stmts.Statements[i].(type) {
		case *ast.FunctionCall:
			symbol, err := parentScope.Lookup(statement.Name)
			// TODO(errors)
			if err != nil {
				return err
			}

			function := symbol.(values.Function)
			args, err := codegen.getExprList(currentBlockScope, statement.Args)
			// TODO(errors)
			if err != nil {
				return err
			}
			codegen.builder.CreateCall(function.Ty, function.Fn, args, "call")
		case *ast.ReturnStmt:
			returnValue, err := codegen.getExpr(currentBlockScope, statement.Value)
			// TODO(errors)
			if err != nil {
				return err
			}
			codegen.builder.CreateRet(returnValue)
		case *ast.CondStmt:
			err := codegen.generateCondStmt(currentBlockScope, function, statement)
			if err != nil {
				return err
			}
		case *ast.VarDeclStmt:
			varTy := codegen.getType(statement.Type)
			varPtr := codegen.builder.CreateAlloca(varTy, ".ptr")
			varExpr, err := codegen.getExpr(currentBlockScope, statement.Value)
			// TODO(errors)
			if err != nil {
				return nil
			}
			codegen.builder.CreateStore(varExpr, varPtr)

			variable := values.Variable{
				Ty:  varTy,
				Ptr: varPtr,
			}
			err = currentBlockScope.Insert(statement.Name.Lexeme.(string), variable)
			// TODO(errors)
			if err != nil {
				return err
			}
		default:
			log.Fatalf("unimplemented block statement: %s", statement)
		}
	}
	return nil
}

func (codegen *codegen) generateExternDecl(external *ast.ExternDecl) error {
	for i := range external.Prototypes {
		err := codegen.generatePrototype(external.Prototypes[i])
		// TODO(errors)
		if err != nil {
			return err
		}
	}
	return nil
}

func (codegen *codegen) generatePrototype(prototype *ast.Proto) error {
	returnType := codegen.getType(prototype.RetType)
	paramsTypes := codegen.getFieldListTypes(prototype.Params)
	functionType := llvm.FunctionType(returnType, paramsTypes, prototype.Params.IsVariadic)
	functionValue := llvm.AddFunction(codegen.module, prototype.Name, functionType)

	function := values.NewFunctionValue(functionValue, functionType, nil)
	err := codegen.universe.Insert(prototype.Name, function)
	return err
}

func (codegen *codegen) getType(ty ast.ExprType) llvm.Type {
	switch exprTy := ty.(type) {
	case ast.BasicType:
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
		default:
			log.Fatalf("invalid basic type token: '%s'", exprTy.Kind)
		}
	case ast.PointerType:
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

func (codegen *codegen) getExprList(parentScope *scope.Scope[values.LLVMValue], expressions []ast.Expr) ([]llvm.Value, error) {
	values := make([]llvm.Value, len(expressions))
	for i := range expressions {
		expr, err := codegen.getExpr(parentScope, expressions[i])
		// TODO(errors)
		if err != nil {
			return nil, err
		}
		values[i] = expr
	}
	return values, nil
}

func (codegen *codegen) getExpr(parentScope *scope.Scope[values.LLVMValue], expr ast.Expr) (llvm.Value, error) {
	switch currentExpr := expr.(type) {
	case *ast.LiteralExpr:
		switch currentExpr.Kind {
		case kind.INTEGER_LITERAL:
			integerLiteral := uint64(currentExpr.Value.(int))
			return llvm.ConstInt(codegen.context.Int32Type(), integerLiteral, false), nil
		case kind.NEGATIVE_INTEGER_LITERAL:
			negativeIntegerLiteral := int64(currentExpr.Value.(int))
			return llvm.ConstInt(codegen.context.Int32Type(), uint64(negativeIntegerLiteral), false), nil
		case kind.STRING_LITERAL:
			stringLiteral := currentExpr.Value.(string)
			globalStrLiteral, ok := codegen.globalStrLiterals[stringLiteral]
			if ok {
				return globalStrLiteral, nil
			}
			globalStrPtr := codegen.builder.CreateGlobalStringPtr(stringLiteral, ".str")
			codegen.globalStrLiterals[stringLiteral] = globalStrPtr
			return globalStrPtr, nil
		case kind.TRUE_BOOL_LITERAL:
			trueBoolLiteral := llvm.ConstInt(codegen.context.Int1Type(), 1, false)
			return trueBoolLiteral, nil
		case kind.FALSE_BOOL_LITERAL:
			falseBoolLiteral := llvm.ConstInt(codegen.context.Int1Type(), 0, false)
			return falseBoolLiteral, nil
		default:
			log.Fatalf("unimplemented literal expr: %s", expr)
		}
	case *ast.IdExpr:
		varName := currentExpr.Name.Lexeme.(string)
		symbol, err := parentScope.Lookup(varName)
		// TODO(errors)
		if err != nil {
			return llvm.Value{}, err
		}
		// TODO(errors)
		if symbol == nil {
			log.Fatalf("local not defined: %s", varName)
		}
		localVar := symbol.(values.Variable)
		loadedVariable := codegen.builder.CreateLoad(localVar.Ty, localVar.Ptr, ".load")
		return loadedVariable, nil
	default:
		log.Fatalf("unimplemented expr: %s", expr)
	}
	// NOTE: this line should be unreachable
	log.Fatalf("REACHING AN UNREACHABLE LINE AT getExpr AT getExpr")
	return llvm.Value{}, nil
}

func (codegen *codegen) generateCondStmt(parentScope *scope.Scope[values.LLVMValue], function values.Function, condStmt *ast.CondStmt) error {
	ifBlock := llvm.AddBasicBlock(function.Fn, ".if")
	elseBlock := llvm.AddBasicBlock(function.Fn, ".else")
	endBlock := llvm.AddBasicBlock(function.Fn, ".end")

	ifScope := scope.New(parentScope)
	ifExpr, err := codegen.getExpr(ifScope, condStmt.IfStmt.Expr)
	if err != nil {
		return err
	}

	codegen.builder.CreateCondBr(ifExpr, ifBlock, elseBlock)

	codegen.builder.SetInsertPointAtEnd(ifBlock)
	err = codegen.generateBlock(ifScope, function, condStmt.IfStmt.Block)
	// TODO(errors)
	if err != nil {
		return nil
	}
	codegen.builder.CreateBr(endBlock)

	// TODO: implement elif statements (basically an if inside the else)

	codegen.builder.SetInsertPointAtEnd(elseBlock)
	if condStmt.ElseStmt != nil {
		elseScope := scope.New(parentScope)
		err := codegen.generateBlock(elseScope, function, condStmt.ElseStmt.Block)
		// TODO(errors)
		if err != nil {
			return err
		}
	}
	codegen.builder.CreateBr(endBlock)
	codegen.builder.SetInsertPointAtEnd(endBlock)
	return nil
}
