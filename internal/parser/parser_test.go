package parser

import (
	"fmt"
	"testing"

	"github.com/HicaroD/Telia/internal/ast"
	"github.com/HicaroD/Telia/internal/diagnostics"
	"github.com/HicaroD/Telia/internal/lexer"
	"github.com/HicaroD/Telia/internal/lexer/token"
)

func TestFnDecl(t *testing.T) {
	filename := "test.tt"
	tests := []struct {
		input string
		check func(t *testing.T, node *ast.Node)
	}{
		{
			input: "fn do_nothing() {}",
			check: func(t *testing.T, node *ast.Node) {
				if node.Kind != ast.KIND_FN_DECL {
					t.Errorf("expected KIND_FN_DECL, got %v", node.Kind)
				}
				fnDecl := node.Node.(*ast.FnDecl)
				if string(fnDecl.Name.Lexeme) != "do_nothing" {
					t.Errorf("expected name 'do_nothing', got %s", fnDecl.Name.Lexeme)
				}
				if fnDecl.Params.Fields != nil {
					t.Errorf("expected no params, got %v", fnDecl.Params.Fields)
				}
				if !fnDecl.RetType.IsVoid() {
					t.Errorf("expected void return type, got %v", fnDecl.RetType)
				}
			},
		},
		{
			input: "fn do_nothing(a bool) {}",
			check: func(t *testing.T, node *ast.Node) {
				fnDecl := node.Node.(*ast.FnDecl)
				if string(fnDecl.Name.Lexeme) != "do_nothing" {
					t.Errorf("expected name 'do_nothing', got %s", fnDecl.Name.Lexeme)
				}
				if len(fnDecl.Params.Fields) != 1 {
					t.Errorf("expected 1 param, got %d", len(fnDecl.Params.Fields))
				}
				if !fnDecl.RetType.IsVoid() {
					t.Errorf("expected void return type")
				}
			},
		},
		{
			input: "fn do_nothing(a bool, b i32) {}",
			check: func(t *testing.T, node *ast.Node) {
				fnDecl := node.Node.(*ast.FnDecl)
				if len(fnDecl.Params.Fields) != 2 {
					t.Errorf("expected 2 params, got %d", len(fnDecl.Params.Fields))
				}
				if !fnDecl.RetType.IsVoid() {
					t.Errorf("expected void return type")
				}
			},
		},
		{
			input: "fn do_nothing(a bool, b i32) i8 {}",
			check: func(t *testing.T, node *ast.Node) {
				fnDecl := node.Node.(*ast.FnDecl)
				if fnDecl.RetType.T.(*ast.BasicType).Kind != token.I8_TYPE {
					t.Errorf("expected i8 return type, got %v", fnDecl.RetType.T)
				}
			},
		},
		{
			input: "fn do_nothing(a bool, b i32) u8 {}",
			check: func(t *testing.T, node *ast.Node) {
				fnDecl := node.Node.(*ast.FnDecl)
				if fnDecl.RetType.T.(*ast.BasicType).Kind != token.U8_TYPE {
					t.Errorf("expected u8 return type, got %v", fnDecl.RetType.T)
				}
			},
		},
		{
			input: "fn do_nothing(a bool, b i32) i16 {}",
			check: func(t *testing.T, node *ast.Node) {
				fnDecl := node.Node.(*ast.FnDecl)
				if fnDecl.RetType.T.(*ast.BasicType).Kind != token.I16_TYPE {
					t.Errorf("expected i16 return type")
				}
			},
		},
		{
			input: "fn do_nothing(a bool, b i32) u16 {}",
			check: func(t *testing.T, node *ast.Node) {
				fnDecl := node.Node.(*ast.FnDecl)
				if fnDecl.RetType.T.(*ast.BasicType).Kind != token.U16_TYPE {
					t.Errorf("expected u16 return type")
				}
			},
		},
		{
			input: "fn do_nothing(a bool, b i32) i32 {}",
			check: func(t *testing.T, node *ast.Node) {
				fnDecl := node.Node.(*ast.FnDecl)
				if fnDecl.RetType.T.(*ast.BasicType).Kind != token.I32_TYPE {
					t.Errorf("expected i32 return type")
				}
			},
		},
		{
			input: "fn do_nothing(a bool, b i32) u32 {}",
			check: func(t *testing.T, node *ast.Node) {
				fnDecl := node.Node.(*ast.FnDecl)
				if fnDecl.RetType.T.(*ast.BasicType).Kind != token.U32_TYPE {
					t.Errorf("expected u32 return type")
				}
			},
		},
		{
			input: "fn do_nothing(a bool, b i32) i64 {}",
			check: func(t *testing.T, node *ast.Node) {
				fnDecl := node.Node.(*ast.FnDecl)
				if fnDecl.RetType.T.(*ast.BasicType).Kind != token.I64_TYPE {
					t.Errorf("expected i64 return type")
				}
			},
		},
		{
			input: "fn do_nothing(a bool, b i32) u64 {}",
			check: func(t *testing.T, node *ast.Node) {
				fnDecl := node.Node.(*ast.FnDecl)
				if fnDecl.RetType.T.(*ast.BasicType).Kind != token.U64_TYPE {
					t.Errorf("expected u64 return type")
				}
			},
		},
		{
			input: "fn do_nothing(a bool, b i32) bool {}",
			check: func(t *testing.T, node *ast.Node) {
				fnDecl := node.Node.(*ast.FnDecl)
				if fnDecl.RetType.T.(*ast.BasicType).Kind != token.BOOL_TYPE {
					t.Errorf("expected bool return type")
				}
			},
		},
		{
			input: "fn do_nothing(a bool, b i32) *i8 {}",
			check: func(t *testing.T, node *ast.Node) {
				fnDecl := node.Node.(*ast.FnDecl)
				if !fnDecl.RetType.IsPointer() {
					t.Errorf("expected pointer return type")
				}
			},
		},
	}

	for _, test := range tests {
		t.Run(fmt.Sprintf("TestFnDecl('%s')", test.input), func(t *testing.T) {
			collector := diagnostics.New()
			lex := lexer.New(FakeLoc(filename), []byte(test.input), collector)
			p := NewForTest(lex, collector)

			node, err := p.ParseFnDecl(ast.Attributes{})
			if err != nil {
				t.Fatal(err)
			}

			test.check(t, node)
		})
	}
}

func TestFnDeclWithBody(t *testing.T) {
	tests := []struct {
		input string
		check func(t *testing.T, node *ast.Node)
	}{
		{
			input: "fn add(a i32, b i32) i32 {\nreturn a + b\n}",
			check: func(t *testing.T, node *ast.Node) {
				fnDecl := node.Node.(*ast.FnDecl)
				if fnDecl.Block == nil {
					t.Errorf("expected block to not be nil")
				}
				if len(fnDecl.Block.Statements) == 0 {
					t.Errorf("expected at least one statement")
				}
			},
		},
		{
			input: "fn empty() {}",
			check: func(t *testing.T, node *ast.Node) {
				fnDecl := node.Node.(*ast.FnDecl)
				if fnDecl.Block == nil {
					t.Errorf("expected block to not be nil")
				}
				if len(fnDecl.Block.Statements) != 0 {
					t.Errorf("expected empty block")
				}
			},
		},
	}

	for _, test := range tests {
		t.Run(fmt.Sprintf("TestFnDeclWithBody('%s')", test.input), func(t *testing.T) {
			collector := diagnostics.New()
			lex := lexer.New(FakeLoc("test.tt"), []byte(test.input), collector)
			p := NewForTest(lex, collector)

			node, err := p.ParseFnDecl(ast.Attributes{})
			if err != nil {
				t.Fatal(err)
			}

			test.check(t, node)
		})
	}
}

func TestForLoop(t *testing.T) {
	filename := "test.tt"
	tests := []struct {
		input string
		check func(t *testing.T, node *ast.Node)
	}{
		{
			input: "for i := 0; i < 10; i = i + 1 {}",
			check: func(t *testing.T, node *ast.Node) {
				if node.Kind != ast.KIND_FOR_LOOP_STMT {
					t.Errorf("expected KIND_FOR_LOOP_STMT, got %v", node.Kind)
				}
				forLoop := node.Node.(*ast.ForLoop)
				if forLoop.Init == nil {
					t.Errorf("expected init to not be nil")
				}
				if forLoop.Cond == nil {
					t.Errorf("expected condition to not be nil")
				}
				if forLoop.Update == nil {
					t.Errorf("expected update to not be nil")
				}
			},
		},
	}

	for _, test := range tests {
		t.Run(fmt.Sprintf("TestForLoop('%s')", test.input), func(t *testing.T) {
			forLoop, err := ParseForLoopFrom(test.input, filename)
			if err != nil {
				t.Fatal(err)
			}

			node := &ast.Node{
				Kind: ast.KIND_FOR_LOOP_STMT,
				Node: forLoop,
			}
			test.check(t, node)
		})
	}
}

func TestWhileLoop(t *testing.T) {
	filename := "test.tt"
	tests := []struct {
		input string
		check func(t *testing.T, node *ast.Node)
	}{
		{
			input: "while true {}",
			check: func(t *testing.T, node *ast.Node) {
				if node.Kind != ast.KIND_WHILE_LOOP_STMT {
					t.Errorf("expected KIND_WHILE_LOOP_STMT, got %v", node.Kind)
				}
				whileLoop := node.Node.(*ast.WhileLoop)
				if whileLoop.Cond == nil {
					t.Errorf("expected condition to not be nil")
				}
			},
		},
	}

	for _, test := range tests {
		t.Run(fmt.Sprintf("TestWhileLoop('%s')", test.input), func(t *testing.T) {
			whileLoop, err := ParseWhileLoopFrom(test.input, filename)
			if err != nil {
				t.Fatal(err)
			}

			node := &ast.Node{
				Kind: ast.KIND_WHILE_LOOP_STMT,
				Node: whileLoop,
			}
			test.check(t, node)
		})
	}
}

func TestLiteralExpr(t *testing.T) {
	filename := "test.tt"
	tests := []struct {
		input string
		check func(t *testing.T, node *ast.Node)
	}{
		{
			input: "1",
			check: func(t *testing.T, node *ast.Node) {
				if node.Kind != ast.KIND_LITERAL_EXPR {
					t.Errorf("expected KIND_LITERAL_EXPR, got %v", node.Kind)
				}
				lit := node.Node.(*ast.LiteralExpr)
				if string(lit.Value) != "1" {
					t.Errorf("expected value '1', got %s", lit.Value)
				}
			},
		},
		{
			input: "true",
			check: func(t *testing.T, node *ast.Node) {
				if node.Kind != ast.KIND_LITERAL_EXPR {
					t.Errorf("expected KIND_LITERAL_EXPR")
				}
				lit := node.Node.(*ast.LiteralExpr)
				if string(lit.Value) != "true" {
					t.Errorf("expected value 'true', got %s", lit.Value)
				}
			},
		},
		{
			input: "false",
			check: func(t *testing.T, node *ast.Node) {
				lit := node.Node.(*ast.LiteralExpr)
				if string(lit.Value) != "false" {
					t.Errorf("expected value 'false', got %s", lit.Value)
				}
			},
		},
		{
			input: "\"Hello, world\"",
			check: func(t *testing.T, node *ast.Node) {
				lit := node.Node.(*ast.LiteralExpr)
				if string(lit.Value) != "Hello, world" {
					t.Errorf("expected value 'Hello, world', got %s", lit.Value)
				}
			},
		},
		{
			input: "3.14",
			check: func(t *testing.T, node *ast.Node) {
				lit := node.Node.(*ast.LiteralExpr)
				if string(lit.Value) != "3.14" {
					t.Errorf("expected value '3.14', got %s", lit.Value)
				}
			},
		},
	}

	for _, test := range tests {
		t.Run(fmt.Sprintf("TestLiteralExpr('%s')", test.input), func(t *testing.T) {
			actualNode, err := ParseExprFrom(test.input, filename)
			if err != nil {
				t.Errorf("unexpected error '%v'", err)
			}
			test.check(t, actualNode)
		})
	}
}

func TestUnaryExpr(t *testing.T) {
	filename := "test.tt"
	tests := []struct {
		input string
		check func(t *testing.T, node *ast.Node)
	}{
		{
			input: "-1",
			check: func(t *testing.T, node *ast.Node) {
				if node.Kind != ast.KIND_UNARY_EXPR {
					t.Errorf("expected KIND_UNARY_EXPR, got %v", node.Kind)
				}
				unary := node.Node.(*ast.UnaryExpr)
				if unary.Op != token.MINUS {
					t.Errorf("expected MINUS operator")
				}
				if unary.Value == nil {
					t.Errorf("expected value to not be nil")
				}
			},
		},
		{
			input: "not true",
			check: func(t *testing.T, node *ast.Node) {
				unary := node.Node.(*ast.UnaryExpr)
				if unary.Op != token.NOT {
					t.Errorf("expected NOT operator")
				}
			},
		},
		{
			input: "*ptr",
			check: func(t *testing.T, node *ast.Node) {
				if node.Kind != ast.KIND_DEREF_POINTER_EXPR {
					t.Errorf("expected KIND_DEREF_POINTER_EXPR, got %v", node.Kind)
				}
			},
		},
		{
			input: "&value",
			check: func(t *testing.T, node *ast.Node) {
				if node.Kind != ast.KIND_ADDRESS_OF_EXPR {
					t.Errorf("expected KIND_ADDRESS_OF_EXPR, got %v", node.Kind)
				}
			},
		},
	}

	for _, test := range tests {
		t.Run(fmt.Sprintf("TestUnaryExpr('%s')", test.input), func(t *testing.T) {
			actualNode, err := ParseExprFrom(test.input, filename)
			if err != nil {
				t.Errorf("unexpected error '%v'", err)
			}
			test.check(t, actualNode)
		})
	}
}

func TestBinaryExpr(t *testing.T) {
	filename := "test.tt"
	tests := []struct {
		input string
		check func(t *testing.T, node *ast.Node)
	}{
		{
			input: "1 + 1",
			check: func(t *testing.T, node *ast.Node) {
				if node.Kind != ast.KIND_BINARY_EXPR {
					t.Errorf("expected KIND_BINARY_EXPR")
				}
				bin := node.Node.(*ast.BinExpr)
				if bin.Op != token.PLUS {
					t.Errorf("expected PLUS operator")
				}
			},
		},
		{
			input: "2 - 1",
			check: func(t *testing.T, node *ast.Node) {
				bin := node.Node.(*ast.BinExpr)
				if bin.Op != token.MINUS {
					t.Errorf("expected MINUS operator")
				}
			},
		},
		{
			input: "5 * 10",
			check: func(t *testing.T, node *ast.Node) {
				bin := node.Node.(*ast.BinExpr)
				if bin.Op != token.STAR {
					t.Errorf("expected STAR operator")
				}
			},
		},
		{
			input: "10 / 1",
			check: func(t *testing.T, node *ast.Node) {
				bin := node.Node.(*ast.BinExpr)
				if bin.Op != token.SLASH {
					t.Errorf("expected SLASH operator")
				}
			},
		},
		{
			input: "3 + 4 * 5",
			check: func(t *testing.T, node *ast.Node) {
				bin := node.Node.(*ast.BinExpr)
				if bin.Op != token.PLUS {
					t.Errorf("expected PLUS operator at root")
				}
				rightBin := bin.Right.Node.(*ast.BinExpr)
				if rightBin.Op != token.STAR {
					t.Errorf("expected STAR operator on right")
				}
			},
		},
		{
			input: "1 > 1",
			check: func(t *testing.T, node *ast.Node) {
				bin := node.Node.(*ast.BinExpr)
				if bin.Op != token.GREATER {
					t.Errorf("expected GREATER operator")
				}
			},
		},
		{
			input: "1 >= 1",
			check: func(t *testing.T, node *ast.Node) {
				bin := node.Node.(*ast.BinExpr)
				if bin.Op != token.GREATER_EQ {
					t.Errorf("expected GREATER_EQ operator")
				}
			},
		},
		{
			input: "1 < 1",
			check: func(t *testing.T, node *ast.Node) {
				bin := node.Node.(*ast.BinExpr)
				if bin.Op != token.LESS {
					t.Errorf("expected LESS operator")
				}
			},
		},
		{
			input: "1 <= 1",
			check: func(t *testing.T, node *ast.Node) {
				bin := node.Node.(*ast.BinExpr)
				if bin.Op != token.LESS_EQ {
					t.Errorf("expected LESS_EQ operator")
				}
			},
		},
		{
			input: "not (1 > 1)",
			check: func(t *testing.T, node *ast.Node) {
				if node.Kind != ast.KIND_UNARY_EXPR {
					t.Errorf("expected KIND_UNARY_EXPR for not")
				}
			},
		},
		{
			input: "1 > 1 and 1 > 1",
			check: func(t *testing.T, node *ast.Node) {
				bin := node.Node.(*ast.BinExpr)
				if bin.Op != token.AND {
					t.Errorf("expected AND operator")
				}
			},
		},
		{
			input: "1 > 1 or 1 > 1",
			check: func(t *testing.T, node *ast.Node) {
				bin := node.Node.(*ast.BinExpr)
				if bin.Op != token.OR {
					t.Errorf("expected OR operator")
				}
			},
		},
		{
			input: "n == 1",
			check: func(t *testing.T, node *ast.Node) {
				bin := node.Node.(*ast.BinExpr)
				if bin.Op != token.EQUAL_EQUAL {
					t.Errorf("expected EQUAL_EQUAL operator")
				}
			},
		},
		{
			input: "n == 1 or n == 2",
			check: func(t *testing.T, node *ast.Node) {
				bin := node.Node.(*ast.BinExpr)
				if bin.Op != token.OR {
					t.Errorf("expected OR operator at root")
				}
			},
		},
		{
			input: "1 + 1 > 2 and 1 == 1",
			check: func(t *testing.T, node *ast.Node) {
				bin := node.Node.(*ast.BinExpr)
				if bin.Op != token.AND {
					t.Errorf("expected AND operator at root")
				}
			},
		},
		{
			input: "true and true and true and true",
			check: func(t *testing.T, node *ast.Node) {
				bin := node.Node.(*ast.BinExpr)
				if bin.Op != token.AND {
					t.Errorf("expected AND operator")
				}
			},
		},
		{
			input: "1 + multiply_by_2(10)",
			check: func(t *testing.T, node *ast.Node) {
				bin := node.Node.(*ast.BinExpr)
				if bin.Op != token.PLUS {
					t.Errorf("expected PLUS operator")
				}
				if bin.Right.Kind != ast.KIND_FN_CALL {
					t.Errorf("expected right side to be function call")
				}
			},
		},
	}

	for _, test := range tests {
		t.Run(fmt.Sprintf("TestBinaryExpr('%s')", test.input), func(t *testing.T) {
			actualNode, err := ParseExprFrom(test.input, filename)
			if err != nil {
				t.Errorf("unexpected error '%v'", err)
			}
			test.check(t, actualNode)
		})
	}
}

func TestFieldAccessExpr(t *testing.T) {
	filename := "test.tt"
	tests := []struct {
		input string
		check func(t *testing.T, node *ast.Node)
	}{
		{
			input: "first.second",
			check: func(t *testing.T, node *ast.Node) {
				if node.Kind != ast.KIND_FIELD_ACCESS {
					t.Errorf("expected KIND_FIELD_ACCESS, got %v", node.Kind)
				}
				access := node.Node.(*ast.FieldAccess)
				if access.Left == nil {
					t.Errorf("expected left to not be nil")
				}
				if access.Right == nil {
					t.Errorf("expected right to not be nil")
				}
			},
		},
		{
			input: "first.second.third",
			check: func(t *testing.T, node *ast.Node) {
				access := node.Node.(*ast.FieldAccess)
				rightAccess := access.Right.Node.(*ast.FieldAccess)
				if rightAccess.Right.Kind != ast.KIND_ID_EXPR {
					t.Errorf("expected far right to be ID_EXPR")
				}
			},
		},
	}

	for _, test := range tests {
		t.Run(fmt.Sprintf("TestFieldAccessExpr('%s')", test.input), func(t *testing.T) {
			actualNode, err := ParseExprFrom(test.input, filename)
			if err != nil {
				t.Errorf("unexpected error '%v'", err)
			}
			test.check(t, actualNode)
		})
	}
}

func TestFnCall(t *testing.T) {
	filename := "test.tt"
	tests := []struct {
		input string
		check func(t *testing.T, node *ast.Node)
	}{
		{
			input: "foo()",
			check: func(t *testing.T, node *ast.Node) {
				if node.Kind != ast.KIND_FN_CALL {
					t.Errorf("expected KIND_FN_CALL, got %v", node.Kind)
				}
				fnCall := node.Node.(*ast.FnCall)
				if string(fnCall.Name.Lexeme) != "foo" {
					t.Errorf("expected name 'foo', got %s", fnCall.Name.Lexeme)
				}
				if fnCall.Args != nil && len(fnCall.Args) != 0 {
					t.Errorf("expected no args")
				}
			},
		},
		{
			input: "foo(1)",
			check: func(t *testing.T, node *ast.Node) {
				fnCall := node.Node.(*ast.FnCall)
				if len(fnCall.Args) != 1 {
					t.Errorf("expected 1 arg, got %d", len(fnCall.Args))
				}
			},
		},
		{
			input: "foo(1, 2, 3)",
			check: func(t *testing.T, node *ast.Node) {
				fnCall := node.Node.(*ast.FnCall)
				if len(fnCall.Args) != 3 {
					t.Errorf("expected 3 args, got %d", len(fnCall.Args))
				}
			},
		},
		{
			input: "foo(1 + 2)",
			check: func(t *testing.T, node *ast.Node) {
				fnCall := node.Node.(*ast.FnCall)
				if len(fnCall.Args) != 1 {
					t.Errorf("expected 1 arg")
				}
				if fnCall.Args[0].Kind != ast.KIND_BINARY_EXPR {
					t.Errorf("expected arg to be binary expr")
				}
			},
		},
	}

	for _, test := range tests {
		t.Run(fmt.Sprintf("TestFnCall('%s')", test.input), func(t *testing.T) {
			actualNode, err := ParseExprFrom(test.input, filename)
			if err != nil {
				t.Errorf("unexpected error '%v'", err)
			}
			test.check(t, actualNode)
		})
	}
}

func TestVar(t *testing.T) {
	filename := "test.tt"
	tests := []struct {
		input string
		check func(t *testing.T, node *ast.Node)
	}{
		{
			input: "age := 10",
			check: func(t *testing.T, node *ast.Node) {
				if node.Kind != ast.KIND_VAR_STMT {
					t.Errorf("expected KIND_VAR_STMT, got %v", node.Kind)
				}
				varStmt := node.Node.(*ast.VarStmt)
				if len(varStmt.Names) == 0 {
					t.Errorf("expected at least one name")
				}
				firstName := varStmt.Names[0].Node.(*ast.VarIdStmt)
				if string(firstName.Name.Lexeme) != "age" {
					t.Errorf("expected name 'age', got %s", firstName.Name.Lexeme)
				}
				if !varStmt.IsDecl {
					t.Errorf("expected declaration")
				}
			},
		},
		{
			input: "age = 10",
			check: func(t *testing.T, node *ast.Node) {
				varStmt := node.Node.(*ast.VarStmt)
				if varStmt.IsDecl {
					t.Errorf("expected assignment, not declaration")
				}
			},
		},
		{
			input: "age u8 := 10",
			check: func(t *testing.T, node *ast.Node) {
				varStmt := node.Node.(*ast.VarStmt)
				if len(varStmt.Names) == 0 {
					t.Errorf("expected at least one name")
				}
				firstName := varStmt.Names[0].Node.(*ast.VarIdStmt)
				if firstName.Type == nil {
					t.Errorf("expected type to be set")
				}
			},
		},
		{
			input: "a, b := 10, 20",
			check: func(t *testing.T, node *ast.Node) {
				if node.Kind != ast.KIND_VAR_STMT {
					t.Errorf("expected KIND_VAR_STMT, got %v", node.Kind)
				}
			},
		},
	}

	for _, test := range tests {
		t.Run(fmt.Sprintf("TestVar('%s')", test.input), func(t *testing.T) {
			collector := diagnostics.New()
			lex := lexer.New(FakeLoc(filename), []byte(test.input), collector)
			p := NewForTest(lex, collector)

			scope := ast.NewScope(nil)
			node, err := p.ParseVar(scope, false)
			if err != nil {
				t.Fatal(err)
			}

			test.check(t, node)
		})
	}
}

func TestIfStmt(t *testing.T) {
	filename := "test.tt"
	tests := []struct {
		input string
		check func(t *testing.T, node *ast.Node)
	}{
		{
			input: "if true {}",
			check: func(t *testing.T, node *ast.Node) {
				if node.Kind != ast.KIND_COND_STMT {
					t.Errorf("expected KIND_COND_STMT, got %v", node.Kind)
				}
			},
		},
		{
			input: "if true {} else {}\n",
			check: func(t *testing.T, node *ast.Node) {
				ifStmt := node.Node.(*ast.CondStmt)
				if ifStmt.ElseStmt == nil {
					t.Errorf("expected else block")
				}
			},
		},
		{
			input: "if true {} elif false {} else {}\n",
			check: func(t *testing.T, node *ast.Node) {
				ifStmt := node.Node.(*ast.CondStmt)
				if len(ifStmt.ElifStmts) != 1 {
					t.Errorf("expected 1 elif")
				}
			},
		},
	}

	for _, test := range tests {
		t.Run(fmt.Sprintf("TestIfStmt('%s')", test.input), func(t *testing.T) {
			collector := diagnostics.New()
			lex := lexer.New(FakeLoc(filename), []byte(test.input), collector)
			p := NewForTest(lex, collector)

			scope := ast.NewScope(nil)
			node, err := p.ParseCondStmt(scope)
			if err != nil {
				t.Fatal(err)
			}

			test.check(t, node)
		})
	}
}

func TestReturnStmt(t *testing.T) {
	filename := "test.tt"
	tests := []struct {
		input string
		check func(t *testing.T, node *ast.Node)
	}{
		{
			input: "return\n",
			check: func(t *testing.T, node *ast.Node) {
				if node.Kind != ast.KIND_RETURN_STMT {
					t.Errorf("expected KIND_RETURN_STMT")
				}
				retStmt := node.Node.(*ast.ReturnStmt)
				if retStmt.Value.Kind != ast.KIND_VOID_EXPR {
					t.Errorf("expected void return value")
				}
			},
		},
		{
			input: "return 1\n",
			check: func(t *testing.T, node *ast.Node) {
				retStmt := node.Node.(*ast.ReturnStmt)
				if retStmt.Value == nil {
					t.Errorf("expected return value")
				}
			},
		},
		{
			input: "return 1 + 2\n",
			check: func(t *testing.T, node *ast.Node) {
				retStmt := node.Node.(*ast.ReturnStmt)
				if retStmt.Value.Kind != ast.KIND_BINARY_EXPR {
					t.Errorf("expected binary expr as return value")
				}
			},
		},
	}

	for _, test := range tests {
		t.Run(fmt.Sprintf("TestReturnStmt('%s')", test.input), func(t *testing.T) {
			collector := diagnostics.New()
			lex := lexer.New(FakeLoc(filename), []byte(test.input), collector)
			p := NewForTest(lex, collector)

			scope := ast.NewScope(nil)
			block := &ast.BlockStmt{Statements: []*ast.Node{}}
			node, err := p.ParseStmt(block, scope, false)
			if err != nil {
				t.Fatal(err)
			}

			test.check(t, node)
		})
	}
}

func TestStructDecl(t *testing.T) {
	filename := "test.tt"
	tests := []struct {
		input string
		check func(t *testing.T, node *ast.Node)
	}{
		{
			input: "struct Point {\n x i32\n y i32\n}\n",
			check: func(t *testing.T, node *ast.Node) {
				if node.Kind != ast.KIND_STRUCT_DECL {
					t.Errorf("expected KIND_STRUCT_DECL, got %v", node.Kind)
				}
				structDecl := node.Node.(*ast.StructDecl)
				if string(structDecl.Name.Lexeme) != "Point" {
					t.Errorf("expected name 'Point', got %s", structDecl.Name.Lexeme)
				}
				if len(structDecl.Fields) != 2 {
					t.Errorf("expected 2 fields, got %d", len(structDecl.Fields))
				}
			},
		},
		{
			input: "struct Empty {}\n",
			check: func(t *testing.T, node *ast.Node) {
				structDecl := node.Node.(*ast.StructDecl)
				if len(structDecl.Fields) != 0 {
					t.Errorf("expected 0 fields")
				}
			},
		},
	}

	for _, test := range tests {
		t.Run(fmt.Sprintf("TestStructDecl('%s')", test.input), func(t *testing.T) {
			collector := diagnostics.New()
			lex := lexer.New(FakeLoc(filename), []byte(test.input), collector)
			p := NewForTest(lex, collector)

			node, err := p.ParseStruct(ast.Attributes{})
			if err != nil {
				t.Fatal(err)
			}

			test.check(t, node)
		})
	}
}

func TestExternDecl(t *testing.T) {
	filename := "test.tt"
	tests := []struct {
		input string
		check func(t *testing.T, node *ast.Node)
	}{
		{
			input: "extern libc {}\n",
			check: func(t *testing.T, node *ast.Node) {
				if node.Kind != ast.KIND_EXTERN_DECL {
					t.Errorf("expected KIND_EXTERN_DECL, got %v", node.Kind)
				}
				externDecl := node.Node.(*ast.ExternDecl)
				if string(externDecl.Name.Lexeme) != "libc" {
					t.Errorf("expected name 'libc', got %s", externDecl.Name.Lexeme)
				}
			},
		},
		{
			input: "extern libc {\n fn printf()\n}\n",
			check: func(t *testing.T, node *ast.Node) {
				externDecl := node.Node.(*ast.ExternDecl)
				if len(externDecl.Prototypes) != 1 {
					t.Errorf("expected 1 prototype")
				}
			},
		},
		{
			input: "extern libc {\n fn printf(fmt *u8)\n}\n",
			check: func(t *testing.T, node *ast.Node) {
				externDecl := node.Node.(*ast.ExternDecl)
				proto := externDecl.Prototypes[0]
				if len(proto.Params.Fields) != 1 {
					t.Errorf("expected 1 param")
				}
			},
		},
	}

	for _, test := range tests {
		t.Run(fmt.Sprintf("TestExternDecl('%s')", test.input), func(t *testing.T) {
			collector := diagnostics.New()
			lex := lexer.New(FakeLoc(filename), []byte(test.input), collector)
			p := NewForTest(lex, collector)

			node, err := p.ParseExternDecl(ast.Attributes{})
			if err != nil {
				t.Fatal(err)
			}

			test.check(t, node)
		})
	}
}

func TestSyntaxErrors(t *testing.T) {
	filename := "test.tt"
	tests := []struct {
		input       string
		expectError bool
	}{
		{
			input:       "{",
			expectError: true,
		},
		{
			input:       "if",
			expectError: true,
		},
		{
			input:       "elif",
			expectError: true,
		},
		{
			input:       "else",
			expectError: true,
		},
		{
			input:       "fn name(){}",
			expectError: false,
		},
		{
			input:       "fn (){}",
			expectError: true,
		},
		{
			input:       "fn",
			expectError: true,
		},
		{
			input:       "fn name){}",
			expectError: true,
		},
		{
			input:       "fn name",
			expectError: true,
		},
		{
			input:       "fn name({}",
			expectError: true,
		},
		{
			input:       "fn name(",
			expectError: true,
		},
		{
			input:       "fn name(a, b int){}",
			expectError: true,
		},
		{
			input:       "fn name(a",
			expectError: true,
		},
		{
			input:       "fn name(a int, ..., b int){}",
			expectError: true,
		},
		{
			input:       "fn name() }",
			expectError: true,
		},
		{
			input:       "fn name()",
			expectError: true,
		},
		{
			input:       "fn name() {",
			expectError: true,
		},
		{
			input:       "extern libc {}",
			expectError: false,
		},
		{
			input:       "extern {}",
			expectError: false,
		},
		{
			input:       "extern libc }",
			expectError: true,
		},
		{
			input:       "extern libc {",
			expectError: true,
		},
		{
			input:       "extern libc",
			expectError: true,
		},
	}

	for _, test := range tests {
		t.Run(fmt.Sprintf("TestSyntaxErrors('%s')", test.input), func(t *testing.T) {
			collector := diagnostics.New()

			node, err := ParseDeclFrom(test.input, filename)

			if test.expectError && len(collector.Diags) == 0 && err == nil {
				t.Fatalf("expected error but got none")
			}
			if !test.expectError && (len(collector.Diags) > 0 || err != nil) {
				t.Fatalf("unexpected error: %v, diags: %v", err, collector.Diags)
			}
			_ = node
		})
	}
}

func TestSyntaxErrorsOnBlock(t *testing.T) {
	filename := "test.tt"

	tests := []struct {
		input       string
		expectError bool
	}{
		{
			input:       "{",
			expectError: true,
		},
		{
			input:       "{ return }",
			expectError: true,
		},
		{
			input:       "{ return }",
			expectError: true,
		},
	}

	for _, test := range tests {
		t.Run(fmt.Sprintf("TestSyntaxErrorsOnBlock('%s')", test.input), func(t *testing.T) {
			collector := diagnostics.New()

			src := []byte(test.input)
			lex := lexer.New(FakeLoc(filename), src, collector)
			p := NewForTest(lex, collector)

			tmpScope := ast.NewScope(nil)
			_, err := p.ParseBlock(tmpScope)

			if test.expectError && len(collector.Diags) == 0 && err == nil {
				t.Fatalf("expected error but got none")
			}
			if !test.expectError && (len(collector.Diags) > 0 || err != nil) {
				t.Fatalf("unexpected error: %v", err)
			}
		})
	}
}

func TestBlockParsing(t *testing.T) {
	filename := "test.tt"
	tests := []struct {
		input string
		check func(t *testing.T, block *ast.BlockStmt)
	}{
		{
			input: "{}\n",
			check: func(t *testing.T, block *ast.BlockStmt) {
				if len(block.Statements) != 0 {
					t.Errorf("expected empty block")
				}
			},
		},
		{
			input: "{ x := 1\n}\n",
			check: func(t *testing.T, block *ast.BlockStmt) {
				if len(block.Statements) != 1 {
					t.Errorf("expected 1 statement")
				}
			},
		},
		{
			input: "{ x := 1\n y := 2\n}\n",
			check: func(t *testing.T, block *ast.BlockStmt) {
				if len(block.Statements) != 2 {
					t.Errorf("expected 2 statements, got %d", len(block.Statements))
				}
			},
		},
		{
			input: "{ x := 1\n y := 2\n return x + y\n}\n",
			check: func(t *testing.T, block *ast.BlockStmt) {
				if len(block.Statements) != 3 {
					t.Errorf("expected 3 statements")
				}
				if block.Statements[2].Kind != ast.KIND_RETURN_STMT {
					t.Errorf("expected last statement to be return")
				}
			},
		},
	}

	for _, test := range tests {
		t.Run(fmt.Sprintf("TestBlockParsing('%s')", test.input), func(t *testing.T) {
			collector := diagnostics.New()
			lex := lexer.New(FakeLoc(filename), []byte(test.input), collector)
			p := NewForTest(lex, collector)

			scope := ast.NewScope(nil)
			block, err := p.ParseBlock(scope)
			if err != nil {
				t.Fatal(err)
			}

			test.check(t, block)
		})
	}
}

func TestTypeParsing(t *testing.T) {
	tests := []struct {
		input string
		check func(t *testing.T, exprType *ast.ExprType)
	}{
		{
			input: "i32",
			check: func(t *testing.T, exprType *ast.ExprType) {
				if exprType.T.(*ast.BasicType).Kind != token.I32_TYPE {
					t.Errorf("expected i32 type")
				}
			},
		},
		{
			input: "u64",
			check: func(t *testing.T, exprType *ast.ExprType) {
				if exprType.T.(*ast.BasicType).Kind != token.U64_TYPE {
					t.Errorf("expected u64 type")
				}
			},
		},
		{
			input: "bool",
			check: func(t *testing.T, exprType *ast.ExprType) {
				if exprType.T.(*ast.BasicType).Kind != token.BOOL_TYPE {
					t.Errorf("expected bool type")
				}
			},
		},
		{
			input: "*i8",
			check: func(t *testing.T, exprType *ast.ExprType) {
				if !exprType.IsPointer() {
					t.Errorf("expected pointer type")
				}
			},
		},
		{
			input: "**i32",
			check: func(t *testing.T, exprType *ast.ExprType) {
				if !exprType.IsPointer() {
					t.Errorf("expected pointer type")
				}
				if !exprType.T.(*ast.PointerType).Type.IsPointer() {
					t.Errorf("expected pointer to pointer")
				}
			},
		},
		{
			input: "rawptr",
			check: func(t *testing.T, exprType *ast.ExprType) {
				if exprType.T.(*ast.BasicType).Kind != token.RAWPTR_TYPE {
					t.Errorf("expected rawptr type")
				}
			},
		},
		{
			input: "string",
			check: func(t *testing.T, exprType *ast.ExprType) {
				if exprType.T.(*ast.BasicType).Kind != token.STRING_TYPE {
					t.Errorf("expected string type")
				}
			},
		},
	}

	for _, test := range tests {
		t.Run(fmt.Sprintf("TestTypeParsing('%s')", test.input), func(t *testing.T) {
			collector := diagnostics.New()
			lex := lexer.New(FakeLoc("test.tt"), []byte(test.input), collector)
			p := NewForTest(lex, collector)

			exprType, err := p.ParseExprType()
			if err != nil {
				t.Fatal(err)
			}

			test.check(t, exprType)
		})
	}
}

func TestPkgUseDecl(t *testing.T) {
	filename := "test.tt"

	tests := []struct {
		input string
		check func(t *testing.T, node *ast.Node)
	}{
		{
			input: "package foo\n",
			check: func(t *testing.T, node *ast.Node) {
				if node.Kind != ast.KIND_PKG_DECL {
					t.Errorf("expected KIND_PKG_DECL")
				}
				pkgDecl := node.Node.(*ast.PkgDecl)
				if string(pkgDecl.Name.Lexeme) != "foo" {
					t.Errorf("expected package name 'foo'")
				}
			},
		},
		{
			input: "use \"std::foo\"\n",
			check: func(t *testing.T, node *ast.Node) {
				t.Skip("parseUse requires config setup - skipping integration test")
			},
		},
	}

	for _, test := range tests {
		t.Run(fmt.Sprintf("TestPkgUseDecl('%s')", test.input), func(t *testing.T) {
			collector := diagnostics.New()
			lex := lexer.New(FakeLoc(filename), []byte(test.input), collector)
			p := NewForTest(lex, collector)

			var node *ast.Node
			var err error

			if len(test.input) > 7 && test.input[:7] == "package" {
				_, node, err = p.ParsePkgDecl()
			} else {
				t.Skip("parseUse requires config setup - skipping integration test")
			}

			if err != nil {
				t.Fatal(err)
			}

			test.check(t, node)
		})
	}
}
