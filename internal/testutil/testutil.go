package testutil

import (
	"github.com/HicaroD/Telia/internal/ast"
	"github.com/HicaroD/Telia/internal/diagnostics"
	"github.com/HicaroD/Telia/internal/lexer"
	"github.com/HicaroD/Telia/internal/lexer/token"
	"github.com/HicaroD/Telia/internal/sema"
)

const DefaultFilename = "test.t"

func NewLexer(src []byte, filename string) *lexer.Lexer {
	if filename == "" {
		filename = DefaultFilename
	}
	loc := &ast.Loc{Name: filename}
	collector := diagnostics.New()
	return lexer.New(loc, src, collector)
}

func NewLexerWithCollector(src []byte, filename string) (*lexer.Lexer, *diagnostics.Collector) {
	if filename == "" {
		filename = DefaultFilename
	}
	loc := &ast.Loc{Name: filename}
	collector := diagnostics.New()
	return lexer.New(loc, src, collector), collector
}

func NewSemaChecker(collector *diagnostics.Collector) {
	sema.New(collector)
}

func CheckProgram(program *ast.Program, runtime *ast.Package, collector *diagnostics.Collector) error {
	checker := sema.New(collector)
	return checker.Check(program, runtime)
}

func NewExprType(kind token.Kind) *ast.ExprType {
	return ast.NewBasicType(kind)
}

func NewBasicType(kind token.Kind) *ast.ExprType {
	return ast.NewBasicType(kind)
}

func NewPointerType(baseType *ast.ExprType) *ast.ExprType {
	return ast.NewPointerType(baseType)
}

func NewIdExpr(name string, filename string) *ast.IdExpr {
	if filename == "" {
		filename = DefaultFilename
	}
	return &ast.IdExpr{
		Name: token.New([]byte(name), token.ID, token.NewPosition(filename, 1, 1)),
	}
}

func NewLiteralExpr(value string, kind token.Kind) *ast.LiteralExpr {
	return &ast.LiteralExpr{
		Value: []byte(value),
		Type:  ast.NewBasicType(kind),
	}
}

func NewBinExpr(left *ast.Node, op token.Kind, right *ast.Node) *ast.BinExpr {
	return &ast.BinExpr{
		Left:  left,
		Op:    op,
		Right: right,
	}
}

func NewUnaryExpr(op token.Kind, value *ast.Node) *ast.UnaryExpr {
	return &ast.UnaryExpr{
		Op:    op,
		Value: value,
	}
}

func NewNode(kind ast.NodeKind, node any) *ast.Node {
	return &ast.Node{Kind: kind, Node: node}
}

func NodeFromExpr(expr any, kind ast.NodeKind) *ast.Node {
	return &ast.Node{Kind: kind, Node: expr}
}

func FakeLoc(filename string) *ast.Loc {
	if filename == "" {
		filename = DefaultFilename
	}
	return &ast.Loc{Name: filename}
}
