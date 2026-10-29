package parser

import "github.com/HicaroD/Telia/internal/ast"

func (p *Parser) parseHeaderExpr(scope *ast.Scope) (*ast.Node, error) {
	previous := p.noNamedLiteral
	p.noNamedLiteral = true
	defer func() { p.noNamedLiteral = previous }()
	return p.ParseSingleExpr(scope)
}

// Parentheses and brackets remove the ambiguity with a following body brace.
func (p *Parser) parseDelimitedExpr(scope *ast.Scope) (*ast.Node, error) {
	previous := p.noNamedLiteral
	p.noNamedLiteral = false
	defer func() { p.noNamedLiteral = previous }()
	return p.ParseSingleExpr(scope)
}
