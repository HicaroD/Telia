package parser

import "github.com/HicaroD/Telia/internal/lexer/token"

// synchronize stops at a safe sibling boundary, leaving the enclosing brace
// for its owner. Every iteration either consumes a token or returns.
func (p *Parser) synchronize(depth int, declaration bool) {
	for {
		tok := p.lex.Peek()
		if tok.Kind == token.EOF {
			return
		}
		if p.lex.BraceDepth <= depth {
			if !declaration && tok.Kind == token.CLOSE_CURLY {
				return
			}
			if tok.Kind == token.NEWLINE {
				p.lex.Skip()
				return
			}
			if declaration {
				switch tok.Kind {
				case token.FN, token.STRUCT, token.PACKAGE, token.USE, token.EXTERN:
					return
				}
			}
		}
		p.lex.Skip()
	}
}
