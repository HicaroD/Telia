package parser

import (
	"github.com/HicaroD/Telia/internal/ast"
	"github.com/HicaroD/Telia/internal/diagnostics"
	"github.com/HicaroD/Telia/internal/lexer/token"
)

func (p *Parser) parseEnum() (*ast.Node, error) {
	return p.parseEnumDeclaration(nil, ast.Attributes{})
}

func (p *Parser) parseEnumDeclaration(foreign *ast.ExternDecl, attrs ast.Attributes) (*ast.Node, error) {
	start := p.lex.Next()
	name, ok := p.expect(token.ID)
	if !ok {
		return nil, diagnostics.Diag{Title: "Expected enum name", Message: "expected a name after enum", Span: name.Span, Help: "Write enum Name { followed by named variants."}
	}
	d := &ast.EnumDecl{Owner: p.pkg, File: p.file, Name: name, Type: ast.NewBasicType(token.INT_TYPE), Foreign: foreign, Attributes: attrs}
	if foreign != nil {
		d.Type = nil
	}
	var err error
	d.TypeParams, err = p.parseTypeParams()
	if err != nil {
		return nil, err
	}
	if !p.lex.NextIs(token.OPEN_CURLY) {
		var err error
		d.Type, err = p.ParseExprType()
		if err != nil {
			return nil, err
		}
	}
	if found, ok := p.expect(token.OPEN_CURLY); !ok {
		return nil, diagnostics.Diag{Title: "Expected enum body", Message: "expected { after enum name", Span: found.Span, Help: "Put each variant on its own line inside braces."}
	}
	for p.lex.NextIs(token.NEWLINE) {
		p.lex.Skip()
	}
	for !p.lex.NextIs(token.CLOSE_CURLY) {
		attrs, err := p.parseAttributes()
		if err != nil {
			return nil, err
		}
		variant, ok := p.expect(token.ID)
		if !ok {
			return nil, diagnostics.Diag{Title: "Expected enum variant", Message: "expected a named enum variant", Span: variant.Span, Help: "Write a variant name or close the enum with }."}
		}
		v := &ast.EnumVariant{Name: variant, Attributes: attrs}
		if p.lex.NextIs(token.EQUAL) {
			p.lex.Skip()
			var err error
			v.Expr, err = p.ParseSingleExpr(p.pkg.Scope)
			if err != nil {
				return nil, err
			}
		}
		d.Variants = append(d.Variants, v)
		if p.lex.NextIs(token.OPEN_CURLY) {
			var err error
			v.Fields, err = p.parseStructFields()
			if err != nil {
				return nil, err
			}
			d.HasPayload = true
		} else if found, ok := p.expect(token.NEWLINE); !ok {
			return nil, diagnostics.Diag{Title: "Expected new line", Message: "expected a new line after enum variant", Span: found.Span, Help: "Put each enum variant on its own line."}
		}
		for p.lex.NextIs(token.NEWLINE) {
			p.lex.Skip()
		}
	}
	p.lex.Skip()
	node := &ast.Node{Kind: ast.KIND_ENUM_DECL, Node: d, Span: p.lex.SpanFrom(start.Span)}
	if foreign != nil {
		return node, nil
	}
	if previous, ok := p.pkg.Scope.Nodes[name.Name()]; ok {
		return nil, diagnostics.Diag{Title: "Duplicate declaration", Message: "enum name is already declared", Span: name.Span, Help: "Choose a distinct enum name.", Related: []diagnostics.Related{{Span: previous.DeclarationSpan(), Message: "previous declaration"}}}
	}
	if err := p.pkg.Scope.Insert(name.Name(), node); err != nil {
		return nil, err
	}
	return node, nil
}
