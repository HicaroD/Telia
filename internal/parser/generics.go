package parser

import (
	"fmt"

	"github.com/HicaroD/Telia/internal/ast"
	"github.com/HicaroD/Telia/internal/lexer/token"
)

func (p *Parser) parseTypeParams() ([]*token.Token, error) {
	if !p.lex.NextIs(token.LESS) {
		return nil, nil
	}
	p.lex.Skip()
	var params []*token.Token
	seen := map[string]bool{}
	for {
		name, ok := p.expect(token.ID)
		if !ok {
			return nil, fmt.Errorf("%s: expected type parameter name, not %s", name.Pos, name.Name())
		}
		if seen[name.Name()] {
			return nil, fmt.Errorf("%s: duplicate type parameter '%s'", name.Pos, name.Name())
		}
		seen[name.Name()] = true
		params = append(params, name)
		if p.lex.NextIs(token.GREATER) {
			p.lex.Skip()
			return params, nil
		}
		if _, ok := p.expect(token.COMMA); !ok {
			return nil, fmt.Errorf("%s: expected ',' or '>' in type parameters", p.lex.Peek().Pos)
		}
	}
}

func (p *Parser) parseTypeArgs() ([]*ast.ExprType, error) {
	if !p.lex.NextIs(token.LESS) {
		return nil, nil
	}
	p.lex.Skip()
	var args []*ast.ExprType
	for {
		if tok := p.lex.Peek(); tok.Literal {
			return nil, fmt.Errorf("%s: invalid generic type argument: expected a type, got a literal", tok.Pos)
		}
		ty, err := p.ParseExprType()
		if err != nil {
			return nil, err
		}
		args = append(args, ty)
		if p.lex.NextIs(token.GREATER) {
			p.lex.Skip()
			return args, nil
		}
		if _, ok := p.expect(token.COMMA); !ok {
			return nil, fmt.Errorf("%s: expected ',' or '>' in type arguments", p.lex.Peek().Pos)
		}
	}
}

// Only a balanced application followed by '(' or '.{' is an expression
// application. A bare '<' remains an ordinary comparison operator.
func (p *Parser) genericSuffix() token.Kind {
	if p.lex.Peek1().Kind != token.LESS {
		return token.INVALID
	}
	depth := 0
	for i := 1; ; i++ {
		t := p.lex.PeekN(i)
		switch t.Kind {
		case token.LESS:
			depth++
		case token.GREATER:
			depth--
			if depth == 0 {
				next := p.lex.PeekN(i + 1).Kind
				if next == token.OPEN_PAREN || next == token.DOT && p.lex.PeekN(i+2).Kind == token.OPEN_CURLY {
					return next
				}
				return token.INVALID
			}
		case token.NEWLINE, token.EOF, token.EQUAL, token.CLOSE_CURLY:
			return token.INVALID
		}
	}
}
