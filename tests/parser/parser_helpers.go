package parser

import (
	"github.com/HicaroD/Telia/internal/ast"
	"github.com/HicaroD/Telia/internal/diagnostics"
	"github.com/HicaroD/Telia/internal/lexer"
	"github.com/HicaroD/Telia/internal/parser"
)

func ParseExprFrom(expr, filename string) (*ast.Node, error) {
	collector := diagnostics.New()

	src := []byte(expr)
	loc := new(ast.Loc)
	loc.Name = filename
	lex := lexer.New(loc, src, collector)
	p := parser.NewWithLex(lex, collector)

	exprAst, err := p.ParseSingleExpr(nil)
	if err != nil {
		return nil, err
	}
	return exprAst, nil
}

func ParseDeclFrom(src, filename string) (*ast.Node, error) {
	collector := diagnostics.New()
	loc := &ast.Loc{Name: filename}
	lex := lexer.New(loc, []byte(src), collector)
	p := parser.NewWithLex(lex, collector)

	file := &ast.File{
		PkgNameDefined: false,
		Imports:        make(map[string]*ast.UseDecl),
		IsFirstNode:    true,
	}
	pkg := &ast.Package{Scope: ast.NewScope(nil)}
	p.SetFileAndPkg(file, pkg)

	node, _, err := p.NextForTest(file)
	return node, err
}

func ParseForLoopFrom(input, filename string) (*ast.ForLoop, error) {
	collector := diagnostics.New()

	src := []byte(input)
	loc := new(ast.Loc)
	loc.Name = filename
	lex := lexer.New(loc, src, collector)
	p := parser.NewWithLex(lex, collector)

	tempScope := ast.NewScope(nil)
	forLoop, err := p.ParseForLoop(tempScope)
	return forLoop.Node.(*ast.ForLoop), err
}

func ParseWhileLoopFrom(input, filename string) (*ast.WhileLoop, error) {
	collector := diagnostics.New()

	src := []byte(input)
	loc := new(ast.Loc)
	loc.Name = filename
	lex := lexer.New(loc, src, collector)
	p := parser.NewWithLex(lex, collector)

	tempScope := ast.NewScope(nil)
	whileLoop, err := p.ParseWhileLoop(tempScope)
	return whileLoop.Node.(*ast.WhileLoop), err
}

func parseFnDeclFrom(filename, input string, scope *ast.Scope) (*ast.FnDecl, error) {
	collector := diagnostics.New()

	src := []byte(input)
	loc := new(ast.Loc)
	loc.Name = filename
	lex := lexer.New(loc, src, collector)
	p := parser.NewWithLex(lex, collector)
	p.GetPkg().Scope = scope

	fnDecl, err := p.ParseFnDecl(ast.Attributes{})
	if err != nil {
		return nil, err
	}

	return fnDecl.Node.(*ast.FnDecl), nil
}

func parseVarFrom(filename, input string) (*ast.VarIdStmt, error) {
	collector := diagnostics.New()

	src := []byte(input)
	loc := new(ast.Loc)
	loc.Name = filename
	lex := lexer.New(loc, src, collector)
	p := parser.NewWithLex(lex, collector)

	tmpScope := ast.NewScope(nil)
	stmt, err := p.ParseIdStmt(tmpScope)
	if err != nil {
		return nil, err
	}
	return stmt.Node.(*ast.VarIdStmt), nil
}
