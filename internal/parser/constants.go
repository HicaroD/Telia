package parser

import (
	"fmt"

	"github.com/HicaroD/Telia/internal/ast"
	"github.com/HicaroD/Telia/internal/diagnostics"
	"github.com/HicaroD/Telia/internal/lexer/token"
)

func (p *Parser) parseConst(scope *ast.Scope) (*ast.Node, error) {
	start := p.lex.Next()
	name, ok := p.expect(token.ID)
	if !ok {
		return nil, diagnostics.Diag{Title: "Expected constant name", Message: fmt.Sprintf("expected constant name, got %s", name.Kind), Span: name.Span, Help: "Write const Name [Type] = expression."}
	}
	d := &ast.ConstDecl{Owner: p.pkg, File: p.file, Scope: scope, Name: name}
	if p.lex.NextIs(token.COLON_EQUAL) {
		return nil, diagnostics.Diag{Title: "Expected =", Message: "constant declarations use =, not :=", Span: p.lex.Peek().Span, Help: "Replace := with = in this constant declaration."}
	}
	if !p.lex.NextIs(token.EQUAL) {
		var err error
		d.Type, err = p.ParseExprType()
		if err != nil {
			return nil, err
		}
	}
	if found, ok := p.expect(token.EQUAL); !ok {
		return nil, diagnostics.Diag{Title: "Expected =", Message: "expected = in constant declaration", Span: found.Span, Help: "Write const Name [Type] = expression; do not use :=."}
	}
	var err error
	d.Expr, err = p.ParseSingleExpr(scope)
	if err != nil {
		return nil, err
	}
	if found, ok := p.expect(token.NEWLINE); !ok {
		return nil, diagnostics.Diag{Title: "Expected new line", Message: "expected new line after constant declaration", Span: found.Span, Help: "Place the next declaration or statement on a new line."}
	}
	node := &ast.Node{Kind: ast.KIND_CONST_DECL, Node: d, Span: p.lex.SpanFrom(start.Span)}
	if scope == p.pkg.Scope {
		if err := scope.Insert(name.Name(), node); err != nil {
			return nil, diagnostics.Diag{Title: "Duplicate declaration", Message: fmt.Sprintf("constant '%s' already declared on scope", name.Name()), Span: name.Span, Related: []diagnostics.Related{{Span: scope.Nodes[name.Name()].DeclarationSpan(), Message: "previous declaration"}}}
		}
	}
	return node, nil
}
