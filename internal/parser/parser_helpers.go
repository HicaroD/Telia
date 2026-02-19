package parser

import (
	"github.com/HicaroD/Telia/internal/ast"
	"github.com/HicaroD/Telia/internal/diagnostics"
	"github.com/HicaroD/Telia/internal/lexer"
)

const defaultFilename = "test.t"

func FakeLoc(filename string) *ast.Loc {
	if filename == "" {
		filename = defaultFilename
	}
	return &ast.Loc{Name: filename}
}

func NewForTest(lex *lexer.Lexer, collector *diagnostics.Collector) *Parser {
	universe := ast.NewScope(nil)
	pkgScope := ast.NewScope(universe)
	pkg := &ast.Package{
		Scope:  pkgScope,
		IsRoot: true,
	}
	return &Parser{
		lex:       lex,
		collector: collector,
		pkg:       pkg,
	}
}

func (p *Parser) NextForTest(file *ast.File) (*ast.Node, bool, error) {
	return p.next(file)
}

func (p *Parser) ParseFileForTest(loc *ast.Loc) (*ast.File, error) {
	file := &ast.File{
		Loc:              loc,
		PkgNameDefined:   false,
		Imports:          make(map[string]*ast.UseDecl),
		IsFirstNode:      true,
		AnyDeclNodeFound: false,
	}

	p.file = file

	err := p.parseFileDecls(file)
	if err != nil {
		return nil, err
	}

	return file, nil
}

func (p *Parser) GetPkg() *ast.Package {
	return p.pkg
}

func (p *Parser) SetFileAndPkg(file *ast.File, pkg *ast.Package) {
	p.file = file
	p.pkg = pkg
}

func ParseExprFrom(expr, filename string) (*ast.Node, error) {
	collector := diagnostics.New()

	src := []byte(expr)
	loc := new(ast.Loc)
	loc.Name = filename
	lex := lexer.New(loc, src, collector)
	p := NewForTest(lex, collector)

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
	p := NewForTest(lex, collector)

	file := &ast.File{
		PkgNameDefined: false,
		Imports:        make(map[string]*ast.UseDecl),
		IsFirstNode:    true,
	}
	pkg := &ast.Package{Scope: ast.NewScope(nil)}
	p.file = file
	p.pkg = pkg

	node, _, err := p.next(file)
	return node, err
}

func ParseForLoopFrom(input, filename string) (*ast.ForLoop, error) {
	collector := diagnostics.New()

	src := []byte(input)
	loc := new(ast.Loc)
	loc.Name = filename
	lex := lexer.New(loc, src, collector)
	p := NewForTest(lex, collector)

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
	p := NewForTest(lex, collector)

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
	p := NewForTest(lex, collector)
	p.pkg.Scope = scope

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
	p := NewForTest(lex, collector)

	tmpScope := ast.NewScope(nil)
	stmt, err := p.ParseIdStmt(tmpScope)
	if err != nil {
		return nil, err
	}
	return stmt.Node.(*ast.VarIdStmt), nil
}
