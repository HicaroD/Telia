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
	for {
		token := parser.cursor.peek()
		if token == nil || token.Kind == kind.EOF {
			break
		}

		switch token.Kind {
		case kind.FN:
			fnDecl, err := parser.parseFnDecl()
			// TODO(errors)
			if err != nil {
				return nil, err
			}
			astNodes = append(astNodes, fnDecl)
		case kind.EXTERN:
			externDecl, err := parser.parseExternBlockDecl()
			// TODO(errors)
			if err != nil {
				return nil, err
			}
			astNodes = append(astNodes, externDecl)
		default:
			// TODO(errors)
			return nil, fmt.Errorf("unimplemented: %s %s", token.Kind, token.Lexeme)
		}
	}
	return astNodes, nil
}

func (parser *parser) parseExternBlockDecl() (*ast.ExternDecl, error) {
	var err error

	_, err = parser.expect(kind.EXTERN)
	// TODO(errors)
	if err != nil {
		return nil, err
	}

	externName, err := parser.expect(kind.STRING_LITERAL)
	// TODO(errors)
	if err != nil {
		return nil, err
	}

	_, err = parser.expect(kind.OPEN_CURLY)

	var prototypes []*ast.Proto
	for {
		proto, err := parser.parsePrototype()
		// TODO(errors)
		if err != nil {
			return nil, err
		}
		prototypes = append(prototypes, proto)

		token := parser.cursor.peek()
		// TODO(errors)
		if token == nil {
			return nil, fmt.Errorf("can't peek next token")
		}
		if token.Kind == kind.CLOSE_CURLY {
			break
		}
	}
	_, err = parser.expect(kind.CLOSE_CURLY)
	return &ast.ExternDecl{Name: externName, Prototypes: prototypes}, nil
}

func (parser *parser) parsePrototype() (*ast.Proto, error) {
	var err error

	_, err = parser.expect(kind.FN)
	// TODO(errors)
	if err != nil {
		return nil, err
	}

	name, err := parser.expect(kind.ID)
	// TODO(errors)
	if err != nil {
		return nil, err
	}

	params, err := parser.parseFunctionParams()
	// TODO(errors)
	if err != nil {
		return nil, err
	}

	returnType, err := parser.parseExprType()
	// TODO(errors)
	if err != nil {
		return nil, err
	}

	_, err = parser.expect(kind.SEMICOLON)
	// TODO(errors)
	if err != nil {
		return nil, err
	}

	return &ast.Proto{Name: name.Lexeme.(string), Params: params, RetType: returnType}, nil
}

func (parser *parser) parseFnDecl() (*ast.FunctionDecl, error) {
	var err error

	_, err = parser.expect(kind.FN)
	// TODO(errors)
	if err != nil {
		return nil, err
	}

	name, err := parser.expect(kind.ID)
	// TODO(errors)
	if err != nil {
		return nil, err
	}

	params, err := parser.parseFunctionParams()
	// TODO(errors)
	if err != nil {
		return nil, err
	}

	returnType, err := parser.parseFnReturnType()
	// TODO(errors)
	if err != nil {
		return nil, err
	}

	block, err := parser.parseBlock()
	// TODO(errors)
	if err != nil {
		return nil, err
	}

	fnDecl := ast.FunctionDecl{Name: name.Lexeme.(string), Params: params, Block: block, RetType: returnType}
	return &fnDecl, nil
}

func (parser *parser) parseFunctionParams() (*ast.FieldList, error) {
	var err error
	var params []*ast.Field
	isVariadic := false

	openParen, err := parser.expect(kind.OPEN_PAREN)
	// TODO(errors)
	if err != nil {
		return nil, err
	}

	for {
		if parser.cursor.nextIs(kind.DOT_DOT_DOT) {
			isVariadic = true
			parser.cursor.skip()
			// TODO(errors)
			if !parser.cursor.nextIs(kind.CLOSE_PAREN) {
				return nil, fmt.Errorf("... is not at the end of function parameters")
			}
			break
		}
		if parser.cursor.nextIs(kind.CLOSE_PAREN) {
			break
		}
		if parser.cursor.nextIs(kind.COMMA) {
			parser.cursor.skip() // ,
			continue
		}

		name, err := parser.expect(kind.ID)
		// TODO(errors): unable to parse identifier
		if err != nil {
			return nil, err
		}
		paramType, err := parser.parseExprType()
		if err != nil {
			return nil, err
		}
		params = append(params, &ast.Field{Name: name.Lexeme.(string), Type: paramType})

	}

	closeParen, err := parser.expect(kind.CLOSE_PAREN)
	// TODO(errors)
	if err != nil {
		return nil, err
	}

	return &ast.FieldList{Open: openParen, Fields: params, Close: closeParen, IsVariadic: isVariadic}, nil
}

func (parser *parser) parseFnReturnType() (ast.ExprType, error) {
	if parser.cursor.nextIs(kind.OPEN_CURLY) {
		return nil, nil
	}

	// TODO(errors)
	if !parser.nextIsPossibleType() {
		return nil, fmt.Errorf("Not a valid function return type annotation")
	}

	returnType, err := parser.parseExprType()
	// TODO(errors)
	if err != nil {
		return nil, err
	}
	return returnType, nil
}

// TODO: maybe return a boolean for saying if matches or not
func (parser *parser) expect(expectedKind kind.TokenKind) (*token.Token, error) {
	token := parser.cursor.peek()
	// TODO(errors)
	if token == nil {
		return nil, fmt.Errorf("can't peek next token because it is null")
	}

	// TODO(errors)
	if token.Kind != expectedKind {
		return nil, fmt.Errorf("expected '%s', but got '%s' '%s'", expectedKind, token.Kind, token.Lexeme)
	}

	parser.cursor.skip()
	return token, nil
}

func (parser *parser) nextIsPossibleType() bool {
	token := parser.cursor.peek()
	if token == nil {
		return false
	}
	switch token.Kind {
	case kind.ID:
		return true
	default:
		_, isBasicType := kind.BASIC_TYPES[token.Kind]
		return isBasicType
	}
}

func (parser *parser) parseExprType() (ast.ExprType, error) {
	token := parser.cursor.peek()
	// TODO(errors)
	if token == nil {
		return nil, fmt.Errorf("can't peek next token because it is nil")
	}

	switch token.Kind {
	case kind.STAR:
		parser.cursor.skip() // *
		ty, err := parser.parseExprType()
		if err != nil {
			return nil, err
		}
		return &ast.PointerType{Type: ty}, nil
	default:
		if _, ok := kind.BASIC_TYPES[token.Kind]; ok {
			parser.cursor.skip()
			return &ast.BasicType{Kind: token.Kind}, nil
		}
		// TODO(errors)
		return nil, fmt.Errorf("Token %s %s is not a proper type", token.Kind, token.Lexeme)
	}
}

func (parser *parser) parseBlock() (*ast.BlockStmt, error) {
	openCurly, err := parser.expect(kind.OPEN_CURLY)

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
			_, err = parser.expect(kind.SEMICOLON)
			// TODO(errors)
			if err != nil {
				return nil, err
			}
			statements = append(statements, ast.ReturnStmt{Return: token, Value: returnValue})
		case kind.ID:
			idNode, err := parser.parseIdStmt()
			// TODO(errors)
			if err != nil {
				return nil, err
			}
			_, err = parser.expect(kind.SEMICOLON)
			// TODO(errors)
			if err != nil {
				return nil, err
			}
			statements = append(statements, idNode)
		default:
			return nil, fmt.Errorf("invalid token for statement parsing: %s", token.Kind)
		}
	}

	closeCurly, err := parser.expect(kind.CLOSE_CURLY)
	// TODO(errors)
	if err != nil {
		return nil, err
	}

	return &ast.BlockStmt{OpenCurly: openCurly.Position, Statements: statements, CloseCurly: closeCurly.Position}, nil
}

func (parser *parser) parseIdStmt() (ast.Stmt, error) {
	identifier, err := parser.expect(kind.ID)
	// TODO(errors)
	if err != nil {
		return nil, err
	}

	next := parser.cursor.peek()
	// TODO(errors)
	if next == nil {
		return nil, fmt.Errorf("Invalid id statement")
	}

	switch next.Kind {
	case kind.OPEN_PAREN:
		fnCall, err := parser.parseFnCall(identifier.Lexeme.(string))
		// TODO(errors)
		if err != nil {
			return nil, err
		}
		return fnCall, nil
	// TODO: deal with variable decl, variable reassignment
	default:
		// TODO(errors)
		return nil, fmt.Errorf("unable to parse id statement")
	}
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

func (parser *parser) parseFnCall(fnName string) (*ast.FuncCallStmt, error) {
	_, err := parser.expect(kind.OPEN_PAREN)
	// TODO(errors)
	if err != nil {
		return nil, err
	}

	var callArgs []ast.Expr
	for {
		if parser.cursor.nextIs(kind.CLOSE_PAREN) {
			parser.cursor.skip()
			break
		}
		expr, err := parser.parseExpr()
		// TODO(errors)
		if err != nil {
			return nil, err
		}
		callArgs = append(callArgs, expr)

		if parser.cursor.nextIs(kind.COMMA) {
			parser.cursor.skip()
			continue
		}
	}

	return &ast.FuncCallStmt{Name: fnName, Args: callArgs}, nil
}
