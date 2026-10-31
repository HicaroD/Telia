package parser

import (
	"fmt"
	"strings"

	"github.com/HicaroD/Telia/internal/ast"
	"github.com/HicaroD/Telia/internal/diagnostics"
	"github.com/HicaroD/Telia/internal/lexer/token"
)

// Test blocks share function checking and code generation, but are not callable
// by source name and therefore are not inserted in the package scope.
func (p *Parser) parseTestBlock() (*ast.Node, bool, error) {
	p.lex.Skip()
	if !strings.HasSuffix(p.file.Loc.Path, "_test.t") {
		return nil, false, fmt.Errorf("test blocks must be declared in *_test.t files")
	}
	description, ok := p.expect(token.STRING_TYPE)
	if !ok || !description.Literal {
		return nil, false, fmt.Errorf("expected a string description after 'test'")
	}
	if p.testDescriptions == nil {
		p.testDescriptions = make(map[string]*token.Token)
	}
	if previous := p.testDescriptions[description.Name()]; previous != nil {
		return nil, false, diagnostics.Diag{
			Title: "Duplicate test", Message: "duplicate test description " + description.Name(), Span: description.Span,
			Related: []diagnostics.Related{
				{Span: previous.Span, Message: "Previously declared here"},
			},
		}
	}
	p.testDescriptions[description.Name()] = description
	scope := ast.NewScope(p.pkg.Scope)
	block, err := p.ParseBlock(scope)
	if err != nil {
		return nil, false, err
	}
	text := description.Name()
	name := *description
	name.Kind = token.ID
	name.Lexeme = []byte(fmt.Sprintf("_telia_test_%d", p.testCount))
	p.testCount++
	fn := &ast.FnDecl{
		Name: &name, CName: name.Name(), TestDescription: &text,
		Scope: scope, Params: &ast.Params{}, RetType: ast.NewBasicType(token.VOID_TYPE), Block: block,
	}
	return &ast.Node{Kind: ast.KIND_FN_DECL, Node: fn}, false, nil
}
