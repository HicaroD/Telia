package parser

import (
	"github.com/HicaroD/Telia/internal/ast"
	"github.com/HicaroD/Telia/internal/diagnostics"
	"github.com/HicaroD/Telia/internal/lexer/token"
)

func (p *Parser) matchExpect(kind token.Kind, help string) (*token.Token, error) {
	tk, ok := p.expect(kind)
	if !ok {
		return nil, diagnostics.Diag{
			Title:   "Invalid match syntax",
			Message: "expected " + kind.String(),
			Span:    tk.Span,
			Help:    help,
		}
	}
	return tk, nil
}

func (p *Parser) parseMatch(scope *ast.Scope) (*ast.Node, error) {
	start := p.lex.Next()
	m := &ast.MatchExpr{}
	var err error
	m.Value, err = p.parseHeaderExpr(scope)
	if err != nil {
		return nil, err
	}
	if _, err := p.matchExpect(token.OPEN_CURLY, "Write match value { Enum::Variant => expression }."); err != nil {
		return nil, err
	}
	previous := p.noNamedLiteral
	p.noNamedLiteral = false
	defer func() { p.noNamedLiteral = previous }()
	p.skipNewLines()
	for !p.lex.NextIs(token.CLOSE_CURLY) {
		arm := &ast.MatchArm{Scope: ast.NewScope(scope)}
		pattern, err := p.parseMatchPattern(arm.Scope)
		if err != nil {
			return nil, err
		}
		arm.MatchPattern = *pattern
		for p.lex.NextIs(token.PIPE) {
			p.lex.Skip()
			pattern, err := p.parseMatchPattern(arm.Scope)
			if err != nil {
				return nil, err
			}
			arm.Alternatives = append(arm.Alternatives, pattern)
		}
		if p.lex.NextIs(token.IF) {
			p.lex.Skip()
			arm.Guard, err = p.ParseSingleExpr(arm.Scope)
			if err != nil {
				return nil, err
			}
		}
		if _, err = p.matchExpect(token.FAT_ARROW, "Separate the pattern and result with =>."); err != nil {
			return nil, err
		}
		if p.lex.NextIs(token.OPEN_CURLY) {
			arm.Block, err = p.ParseBlock(arm.Scope)
		} else {
			arm.Result, err = p.ParseSingleExpr(arm.Scope)
		}
		if err != nil {
			return nil, err
		}
		m.Arms = append(m.Arms, arm)
		if _, err = p.matchExpect(token.NEWLINE, "Put each match arm on a new line."); err != nil {
			return nil, err
		}
		p.skipNewLines()
	}
	p.lex.Skip()
	return &ast.Node{Kind: ast.KIND_MATCH_EXPR, Node: m, Span: p.lex.SpanFrom(start.Span)}, nil
}

func (p *Parser) parseMatchPattern(scope *ast.Scope) (*ast.MatchPattern, error) {
	arm := &ast.MatchPattern{}
	next := p.lex.Peek()
	if next.Kind.IsInteger() || next.Kind == token.MINUS {
		arm.VariantName = next
		var err error
		arm.Literal, err = p.parseUnary(scope)
		return arm, err
	}
	name, err := p.matchExpect(token.ID, "Start an arm with Enum::Variant or _.")
	if err != nil {
		return nil, err
	}
	if name.Name() == "_" {
		arm.Fallback = true
		arm.VariantName = name
	} else if p.lex.NextIs(token.IF) || p.lex.NextIs(token.PIPE) || p.lex.NextIs(token.FAT_ARROW) {
		arm.VariantName = name
		arm.Binding = &ast.Param{Name: name, Attributes: &ast.ParamAttributes{}}
	} else {
		id := &ast.IdType{Name: name}
		id.Args, err = p.parseTypeArgs()
		if err != nil {
			return nil, err
		}
		if _, err = p.matchExpect(token.COLON_COLON, "Qualify the variant as Enum::Variant."); err != nil {
			return nil, err
		}
		next, err := p.matchExpect(token.ID, "Write a variant name after ::.")
		if err != nil {
			return nil, err
		}
		if p.lex.NextIs(token.COLON_COLON) || p.genericSuffix() == token.COLON_COLON {
			id.Namespace, id.Name = name, next
			id.Args, err = p.parseTypeArgs()
			if err != nil {
				return nil, err
			}
			if _, err = p.matchExpect(token.COLON_COLON, "Write package::Enum::Variant."); err != nil {
				return nil, err
			}
			next, err = p.matchExpect(token.ID, "Write a variant name after ::.")
			if err != nil {
				return nil, err
			}
		}
		arm.PatternType = &ast.ExprType{Kind: ast.EXPR_TYPE_ID, T: id}
		arm.VariantName = next
		if p.lex.NextIs(token.DOT) || p.lex.NextIs(token.OPEN_CURLY) {
			if p.lex.NextIs(token.DOT) {
				p.lex.Skip()
			}
			if _, err = p.matchExpect(token.OPEN_CURLY, "Write {field}, {field: binding}, or {field, ..}."); err != nil {
				return nil, err
			}
			for !p.lex.NextIs(token.CLOSE_CURLY) {
				if p.lex.NextIs(token.DOT_DOT) {
					p.lex.Skip()
					arm.OmitRest = true
					break
				}
				field, err := p.matchExpect(token.ID, "Use named fields; write .. to omit remaining fields.")
				if err != nil {
					return nil, err
				}
				binding := field
				if p.lex.NextIs(token.COLON) {
					p.lex.Skip()
					binding, err = p.matchExpect(token.ID, "Write a binding name after :.")
					if err != nil {
						return nil, err
					}
				}
				arm.Fields = append(arm.Fields, &ast.PatternField{Name: field, Binding: &ast.Param{Name: binding, Attributes: &ast.ParamAttributes{}}})
				if !p.lex.NextIs(token.COMMA) {
					break
				}
				p.lex.Skip()
			}
			if _, err = p.matchExpect(token.CLOSE_CURLY, "Close the named-field pattern with }."); err != nil {
				return nil, err
			}
		}
	}
	return arm, nil
}
