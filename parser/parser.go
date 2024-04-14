package parser

import (
	"bufio"
	"fmt"
	"log"
	"strings"

	"github.com/HicaroD/telia-lang/ast"
	"github.com/HicaroD/telia-lang/lexer"
	"github.com/HicaroD/telia-lang/lexer/token"
	"github.com/HicaroD/telia-lang/lexer/token/kind"
)

type parser struct {
	cursor *cursor
}

func New(tokens []*token.Token) *parser {
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
			externDecl, err := parser.parseExternDecl()
			// TODO(errors)
			if err != nil {
				return nil, err
			}
			astNodes = append(astNodes, externDecl)
		default:
			// TODO(errors)
			return nil, fmt.Errorf("unimplemented on parser: %s", token.Lexeme)
		}
	}
	return astNodes, nil
}

// Useful for testing
func parseExprFrom(expr, filename string) (ast.Expr, error) {
	reader := bufio.NewReader(strings.NewReader(expr))
	lex := lexer.New(filename, reader)
	tokens, err := lex.Tokenize()
	if err != nil {
		return nil, err
	}

	parser := New(tokens)
	exprAst, err := parser.parseExpr()
	if err != nil {
		return nil, err
	}
	return exprAst, nil
}

func (parser *parser) parseExternDecl() (*ast.ExternDecl, error) {
	_, ok := parser.expect(kind.EXTERN)
	// TODO(errors)
	if !ok {
		return nil, fmt.Errorf("expected 'extern'")
	}

	externName, ok := parser.expect(kind.STRING_LITERAL)
	// TODO(errors)
	if !ok {
		return nil, fmt.Errorf("expected string literal")
	}

	_, ok = parser.expect(kind.OPEN_CURLY)
	// TODO(errors)
	if !ok {
		return nil, fmt.Errorf("expected '{'")
	}

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

	_, ok = parser.expect(kind.CLOSE_CURLY)
	// TODO(errors)
	if !ok {
		return nil, fmt.Errorf("expected '}'")
	}
	return &ast.ExternDecl{Scope: nil, Name: externName, Prototypes: prototypes}, nil
}

func (parser *parser) parsePrototype() (*ast.Proto, error) {
	_, ok := parser.expect(kind.FN)
	// TODO(errors)
	if !ok {
		return nil, fmt.Errorf("expected 'fn'")
	}

	name, ok := parser.expect(kind.ID)
	// TODO(errors)
	if !ok {
		return nil, fmt.Errorf("expected an identifier")
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

	// TODO(errors)
	_, ok = parser.expect(kind.SEMICOLON)
	if !ok {
		return nil, fmt.Errorf("expected ';'")
	}

	return &ast.Proto{Name: name.Lexeme.(string), Params: params, RetType: returnType}, nil
}

func (parser *parser) parseFnDecl() (*ast.FunctionDecl, error) {
	var err error

	_, ok := parser.expect(kind.FN)
	// TODO(errors)
	if !ok {
		return nil, fmt.Errorf("expected 'fn'")
	}

	name, ok := parser.expect(kind.ID)
	// TODO(errors)
	if !ok {
		return nil, fmt.Errorf("expected an identifier")
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

	fnDecl := ast.FunctionDecl{Scope: nil, Name: name.Lexeme.(string), Params: params, Block: block, RetType: returnType}
	return &fnDecl, nil
}

func parseFnDeclFrom(filename, input string) (*ast.FunctionDecl, error) {
	reader := bufio.NewReader(strings.NewReader(input))

	lexer := lexer.New(filename, reader)
	tokens, err := lexer.Tokenize()
	if err != nil {
		return nil, err
	}

	par := New(tokens)

	fnDecl, err := par.parseFnDecl()
	if err != nil {
		return nil, err
	}

	return fnDecl, nil
}

func (parser *parser) parseFunctionParams() (*ast.FieldList, error) {
	var params []*ast.Field
	isVariadic := false

	openParen, ok := parser.expect(kind.OPEN_PAREN)
	// TODO(errors)
	if !ok {
		return nil, fmt.Errorf("expected '('")
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

		name, ok := parser.expect(kind.ID)
		// TODO(errors): unable to parse identifier
		if !ok {
			return nil, fmt.Errorf("expected an identifier")
		}
		paramType, err := parser.parseExprType()
		// TODO(errors)
		if err != nil {
			return nil, err
		}
		params = append(params, &ast.Field{Name: name, Type: paramType})
	}

	closeParen, ok := parser.expect(kind.CLOSE_PAREN)
	// TODO(errors)
	if !ok {
		return nil, fmt.Errorf("expected ')'")
	}
	return &ast.FieldList{Open: openParen, Fields: params, Close: closeParen, IsVariadic: isVariadic}, nil
}

func (parser *parser) parseFnReturnType() (ast.ExprType, error) {
	if parser.cursor.nextIs(kind.OPEN_CURLY) {
		return ast.BasicType{Kind: kind.VOID_TYPE}, nil
	}

	returnType, err := parser.parseExprType()
	// TODO(errors)
	if err != nil {
		return nil, err
	}
	return returnType, nil
}

// TODO: maybe return a boolean for saying if matches or not
func (parser *parser) expect(expectedKind kind.TokenKind) (*token.Token, bool) {
	token := parser.cursor.peek()
	// TODO(errors)
	if token == nil {
		return nil, false
	}

	// TODO(errors)
	if token.Kind != expectedKind {
		return token, false
	}

	parser.cursor.skip()
	return token, true
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
		// TODO(errors)
		if err != nil {
			return nil, err
		}
		return ast.PointerType{Type: ty}, nil
	default:
		if _, ok := kind.BASIC_TYPES[token.Kind]; ok {
			parser.cursor.skip()
			return ast.BasicType{Kind: token.Kind}, nil
		}
		// TODO(errors)
		return nil, fmt.Errorf("token %s %s is not a proper type", token.Kind, token.Lexeme)
	}
}

func (parser *parser) parseBlock() (*ast.BlockStmt, error) {
	openCurly, ok := parser.expect(kind.OPEN_CURLY)
	// TODO(errors)
	if !ok {
		return nil, fmt.Errorf("expected '{', but got %s", openCurly)
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
			if parser.cursor.nextIs(kind.SEMICOLON) {
				statements = append(statements, &ast.ReturnStmt{Return: token, Value: &ast.VoidExpr{}})
				parser.cursor.skip()
				break
			}

			returnValue, err := parser.parseExpr()
			// TODO(errors)
			if err != nil {
				return nil, err
			}
			_, ok = parser.expect(kind.SEMICOLON)
			// TODO(errors)
			if !ok {
				return nil, fmt.Errorf("expected ';'")
			}
			statements = append(statements, &ast.ReturnStmt{Return: token, Value: returnValue})
		case kind.ID:
			idNode, err := parser.ParseIdStmt()
			// TODO(errors)
			if err != nil {
				return nil, err
			}
			_, ok = parser.expect(kind.SEMICOLON)
			// TODO(errors)
			if !ok {
				return nil, fmt.Errorf("expected ';'")
			}
			statements = append(statements, idNode)
		case kind.IF:
			condStmt, err := parser.parseCondStmt()
			if err != nil {
				return nil, err
			}
			statements = append(statements, condStmt)
		default:
			return nil, fmt.Errorf("invalid token for statement parsing: %s %s", token.Kind, token.Lexeme)
		}
	}

	closeCurly, ok := parser.expect(kind.CLOSE_CURLY)
	// TODO(errors)
	if !ok {
		return nil, fmt.Errorf("expected '}'")
	}

	return &ast.BlockStmt{OpenCurly: openCurly.Position, Statements: statements, CloseCurly: closeCurly.Position}, nil
}

func (parser *parser) ParseIdStmt() (ast.Stmt, error) {
	identifier, ok := parser.expect(kind.ID)
	// TODO(errors)
	if !ok {
		return nil, fmt.Errorf("expected identifier")
	}

	next := parser.cursor.peek()
	// TODO(errors)
	if next == nil {
		return nil, fmt.Errorf("invalid id statement")
	}

	switch next.Kind {
	case kind.OPEN_PAREN:
		fnCall, err := parser.parseFnCall(identifier.Lexeme.(string))
		// TODO(errors)
		if err != nil {
			return nil, err
		}
		return fnCall, nil
	case kind.COLON_EQUAL:
		varDecl, err := parser.parseVarDecl(identifier)
		if err != nil {
			return nil, err
		}
		return varDecl, nil
	// TODO: variable reassignment
	case kind.EQUAL:
		log.Fatalf("unimplemented var reassigment")
	// TODO(errors)
	default:
		return nil, fmt.Errorf("unable to parse id statement: %s", next)
	}
	// TODO(errors)
	log.Fatalf("should be unreachable - parseIdStmt")
	return nil, nil
}

func (parser *parser) parseVarDecl(identifier *token.Token) (*ast.VarDeclStmt, error) {
	// TODO: parse variable declaration with type annotation
	// Right now I am only parsing variables with type inference
	_, ok := parser.expect(kind.COLON_EQUAL)
	if !ok {
		return nil, fmt.Errorf("expected ':=' at parseVarDecl")
	}

	varExpr, err := parser.parseExpr()
	if err != nil {
		return nil, err
	}
	return &ast.VarDeclStmt{Name: identifier, Type: nil, Value: varExpr, NeedsInference: true}, nil
}

func (parser *parser) parseCondStmt() (*ast.CondStmt, error) {
	ifCond, err := parser.parseIfCond()
	// TODO(errors)
	if err != nil {
		return nil, err
	}

	elifConds, err := parser.parseElifConds()
	// TODO(errors)
	if err != nil {
		return nil, err
	}

	elseCond, err := parser.parseElseCond()
	// TODO(errors)
	if err != nil {
		return nil, err
	}

	return &ast.CondStmt{IfStmt: ifCond, ElifStmts: elifConds, ElseStmt: elseCond}, nil
}

func (parser *parser) parseIfCond() (*ast.IfElifCond, error) {
	ifToken, ok := parser.expect(kind.IF)
	// TODO(errors)
	if !ok {
		return nil, fmt.Errorf("expected 'if'")
	}

	// TODO(errors)
	ifExpr, err := parser.parseExpr()
	if err != nil {
		return nil, err
	}

	ifBlock, err := parser.parseBlock()
	// TODO(errors)
	if err != nil {
		return nil, err
	}
	return &ast.IfElifCond{If: &ifToken.Position, Expr: ifExpr, Block: ifBlock}, nil
}

func (parser *parser) parseElifConds() ([]*ast.IfElifCond, error) {
	var elifConds []*ast.IfElifCond
	for {
		elifToken, ok := parser.expect(kind.ELIF)
		if !ok {
			break
		}
		elifExpr, err := parser.parseExpr()
		// TODO(errors)
		if err != nil {
			return nil, err
		}
		elifBlock, err := parser.parseBlock()
		// TODO(errors)
		if err != nil {
			return nil, err
		}
		elifConds = append(elifConds, &ast.IfElifCond{If: &elifToken.Position, Expr: elifExpr, Block: elifBlock})
	}
	return elifConds, nil
}

func (parser *parser) parseElseCond() (*ast.ElseCond, error) {
	elseToken, ok := parser.expect(kind.ELSE)
	if !ok {
		return nil, nil
	}

	elseBlock, err := parser.parseBlock()
	if err != nil {
		return nil, err
	}
	return &ast.ElseCond{Else: &elseToken.Position, Block: elseBlock}, nil
}

func (parser *parser) parseExpr() (ast.Expr, error) {
	return parser.parseLogical()
}

func (parser *parser) parseLogical() (ast.Expr, error) {
	lhs, err := parser.parseComparasion()
	// TODO(errors)
	if err != nil {
		return nil, err
	}

	// TODO: refactor this code, it seems there is a better way of writing this
	for {
		next := parser.cursor.peek()
		if next == nil {
			break
		}
		if _, ok := ast.LOGICAL[next.Kind]; ok {
			parser.cursor.skip()
			rhs, err := parser.parseComparasion()
			// TODO(errors)
			if err != nil {
				return nil, err
			}
			lhs = &ast.BinaryExpr{Left: lhs, Op: next.Kind, Right: rhs}
		} else {
			break
		}
	}
	return lhs, nil
}

func (parser *parser) parseComparasion() (ast.Expr, error) {
	lhs, err := parser.parseTerm()
	// TODO(errors)
	if err != nil {
		return nil, err
	}

	// TODO: refactor this code, it seems there is a better way of writing this
	for {
		next := parser.cursor.peek()
		if next == nil {
			break
		}
		if _, ok := ast.COMPARASION[next.Kind]; ok {
			parser.cursor.skip()
			rhs, err := parser.parseTerm()
			// TODO(errors)
			if err != nil {
				return nil, err
			}
			lhs = &ast.BinaryExpr{Left: lhs, Op: next.Kind, Right: rhs}
		} else {
			break
		}
	}
	return lhs, nil
}

func (parser *parser) parseTerm() (ast.Expr, error) {
	lhs, err := parser.parseFactor()
	// TODO(errors)
	if err != nil {
		return nil, err
	}

	// TODO: refactor this code, it seems there is a better way of writing this
	for {
		next := parser.cursor.peek()
		if next == nil {
			break
		}
		if _, ok := ast.TERM[next.Kind]; ok {
			parser.cursor.skip()
			rhs, err := parser.parseFactor()
			// TODO(errors)
			if err != nil {
				return nil, err
			}
			lhs = &ast.BinaryExpr{Left: lhs, Op: next.Kind, Right: rhs}
		} else {
			break
		}
	}
	return lhs, nil
}

func (parser *parser) parseFactor() (ast.Expr, error) {
	lhs, err := parser.parseUnary()
	// TODO(errors)
	if err != nil {
		return nil, err
	}

	// TODO: refactor this code, it seems there is a better way of writing this
	for {
		next := parser.cursor.peek()
		if next == nil {
			break
		}
		if _, ok := ast.FACTOR[next.Kind]; ok {
			parser.cursor.skip()
			rhs, err := parser.parseUnary()
			// TODO(errors)
			if err != nil {
				return nil, err
			}
			lhs = &ast.BinaryExpr{Left: lhs, Op: next.Kind, Right: rhs}
		} else {
			break
		}
	}
	return lhs, nil

}

func (parser *parser) parseUnary() (ast.Expr, error) {
	next := parser.cursor.peek()
	// TODO(errors)
	if next == nil {
		return nil, fmt.Errorf("unable to peek next on parseUnary")
	}

	if _, ok := ast.UNARY[next.Kind]; ok {
		parser.cursor.skip()
		rhs, err := parser.parseUnary()
		// TODO(errors)
		if err != nil {
			return nil, err
		}
		return &ast.UnaryExpr{Op: next.Kind, Node: rhs}, nil
	}

	return parser.parsePrimary()
}

func (parser *parser) parsePrimary() (ast.Expr, error) {
	token := parser.cursor.peek()
	if token == nil {
		return nil, fmt.Errorf("can't peek next token because it is null")
	}
	switch token.Kind {
	case kind.ID:
		parser.cursor.skip()
		if parser.cursor.nextIs(kind.OPEN_PAREN) {
			return parser.parseFnCall(token.Lexeme.(string))
		}
		return &ast.IdExpr{Name: token}, nil
	case kind.OPEN_PAREN:
		parser.cursor.skip() // (
		expr, err := parser.parseExpr()
		// TODO(errors)
		if err != nil {
			return nil, err
		}
		_, ok := parser.expect(kind.CLOSE_PAREN)
		if !ok {
			return nil, fmt.Errorf("expected closing parenthesis")
		}
		return expr, nil
	default:
		if _, ok := kind.LITERAL_KIND[token.Kind]; ok {
			parser.cursor.skip()
			return &ast.LiteralExpr{Kind: token.Kind, Value: token.Lexeme}, nil
		}
		return nil, fmt.Errorf("invalid token for expression parsing: %s %s %s", token.Kind, token.Lexeme, token.Position)
	}
}

func (parser *parser) parseFnCall(fnName string) (*ast.FunctionCall, error) {
	_, ok := parser.expect(kind.OPEN_PAREN)
	// TODO(errors)
	if !ok {
		return nil, fmt.Errorf("expected '('")
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
	return &ast.FunctionCall{Name: fnName, Args: callArgs}, nil
}
