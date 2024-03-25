package parser

import (
	"fmt"

	"github.com/HicaroD/telia-lang/ast"
	"github.com/HicaroD/telia-lang/lexer/token"
	"github.com/HicaroD/telia-lang/lexer/token/kind"
)

type parser struct {
	cursor *cursor
}

func NewParser(tokens []token.Token) *parser {
	return &parser{cursor: newCursor(tokens)}
}

func (parser *parser) Parse() ([]ast.AstNode, error) {
	var astNodes []ast.AstNode
	for parser.cursor.peek().Kind != kind.EOF {
		token := parser.cursor.peek()
		switch token.Kind {
		case kind.FN:
			fnDecl, err := parser.parseFnDecl()
			if err != nil {
				return nil, err
			}
			astNodes = append(astNodes, fnDecl)
		}
	}
	return astNodes, nil
}

func (parser *parser) parseFnDecl() (ast.AstNode, error) {
	var err error

	_, err = parser.parseExpectedToken(kind.FN)
	// TODO(errors)
	if err != nil {
		return nil, err
	}

	name, err := parser.parseIdentifier()
	// TODO(errors)
	if err != nil {
		return nil, err
	}

	params, err := parser.parseFunctionParams()
	// TODO(errors)
	if err != nil {
		return nil, err
	}

	block, err := parser.parseBlock()
	// TODO(errors)
	if err != nil {
		return nil, err
	}

	fnDecl := ast.FunctionDecl{Name: name, Params: params, Block: block}
	return fnDecl, nil
}

func (parser *parser) parseExpectedToken(expectedKind kind.TokenKind) (*token.Token, error) {
	token := parser.cursor.peek()
	// TODO(errors)
	if token == nil {
		return nil, fmt.Errorf("can't peek next token because it is null")
	}

	// TODO(errors)
	if token.Kind != expectedKind {
		return nil, fmt.Errorf("expected %s, but got %s", expectedKind, token.Kind)
	}

	parser.cursor.skip()
	return token, nil
}

func (parser *parser) parseFunctionParams() (*ast.FieldList, error) {
	var err error
	var params []*ast.Field

	openParen, err := parser.parseExpectedToken(kind.OPEN_PAREN)
	// TODO(errors)
	if err != nil {
		return nil, err
	}

	for {
		if parser.cursor.nextIs(kind.CLOSE_PAREN) {
			break
		}

		if parser.cursor.nextIs(kind.COMMA) {
			parser.cursor.skip() // ,
			continue
		}

		name, err := parser.parseIdentifier()
		// TODO(errors): unable to parse identifier
		if err != nil {
			return nil, err
		}
		paramType, err := parser.parseExprType()

		params = append(params, &ast.Field{Name: name, Type: paramType})
	}

	closeParen, err := parser.parseExpectedToken(kind.CLOSE_PAREN)
	// TODO(errors)
	if err != nil {
		return nil, err
	}

	return &ast.FieldList{Open: openParen, Fields: params, Close: closeParen}, nil
}

func (parser *parser) parseIdentifier() (string, error) {
	token := parser.cursor.peek()
	if token == nil {
		// TODO(errors)
		return "", fmt.Errorf("can't peek next token because it is nil")
	}

	if token.Kind != kind.ID {
		// TODO(errors)
		return "", fmt.Errorf("expected identifier, but got %s", token.Kind)
	}
	parser.cursor.skip()

	name := token.Lexeme.(string)
	return name, nil
}

func (parser *parser) parseExprType() (ast.ExprType, error) {
	token := parser.cursor.peek()
	// TODO(errors)
	if token == nil {
		return nil, fmt.Errorf("can't peek next token because it is nil")
	}

	switch token.Kind {
	case kind.BOOL_TYPE:
		return ast.LiteralType{Kind: token.Kind}, nil
	default:
		// TODO(errors)
		return nil, fmt.Errorf("Token %s is not a proper type", token.Kind)
	}
}

func (parser *parser) parseBlock() (*ast.BlockStmt, error) {
	openCurly, err := parser.parseExpectedToken(kind.OPEN_CURLY)

	// TODO(errors)
	if err != nil {
		return nil, err
	}
	var statements []ast.Stmt

	for {
		token := parser.cursor.peek()
		// TODO(errors)
		if token == nil {
			return nil, fmt.Errorf("can't peek next token because it is null")
		}
		if token.Kind == kind.CLOSE_CURLY {
			break
		}

		switch token.Kind {
		case kind.RETURN:
			parser.cursor.skip()
			returnValue, err := parser.parseExpr()
			// TODO(errors)
			if err != nil {
				return nil, err
			}

			_, err = parser.parseExpectedToken(kind.SEMICOLON)
			// TODO(errors)
			if err != nil {
				return nil, err
			}
			statements = append(statements, ast.ReturnStmt{Return: token, Value: returnValue})
		default:
			return nil, fmt.Errorf("invalid token for statement parsing: %s", token.Kind)
		}
	}

	closeCurly, err := parser.parseExpectedToken(kind.CLOSE_CURLY)
	// TODO(errors)
	if err != nil {
		return nil, err
	}

	return &ast.BlockStmt{OpenCurly: openCurly.Position, Statements: statements, CloseCurly: closeCurly.Position}, nil
}

func (parser *parser) parseExpr() (ast.Expr, error) {
	token := parser.cursor.peek()
	if token == nil {
		return nil, fmt.Errorf("can't peek next token because it is null")
	}
	switch token.Kind {
	case kind.INTEGER_LITERAL, kind.STRING_LITERAL:
		parser.cursor.skip()
		return ast.LiteralExpr{Kind: token.Kind, Value: token.Lexeme}, nil
	default:
		return nil, fmt.Errorf("invalid token for expression parsing: %s", token.Kind)
	}
}
