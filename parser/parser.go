package parser

import (
	"bufio"
	"fmt"
	"log"
	"strings"

	"github.com/HicaroD/Telia/ast"
	"github.com/HicaroD/Telia/collector"
	"github.com/HicaroD/Telia/lexer"
	"github.com/HicaroD/Telia/lexer/token"
	"github.com/HicaroD/Telia/lexer/token/kind"
)

type parser struct {
	Collector *collector.DiagCollector
	cursor    *cursor
}

func New(tokens []*token.Token, diagCollector *collector.DiagCollector) *parser {
	return &parser{cursor: newCursor(tokens), Collector: diagCollector}
}

func (parser *parser) Parse() ([]ast.Node, error) {
	var astNodes []ast.Node
	for {
		token := parser.cursor.peek()
		if token.Kind == kind.EOF {
			break
		}
		switch token.Kind {
		case kind.FN:
			fnDecl, err := parser.parseFnDecl()
			if err != nil {
				return nil, err
			}
			astNodes = append(astNodes, fnDecl)
		case kind.EXTERN:
			externDecl, err := parser.parseExternDecl()
			if err != nil {
				return nil, err
			}
			astNodes = append(astNodes, externDecl)
		default:
			pos := token.Position
			unexpectedTokenOnGlobalScope := collector.Diag{
				Message: fmt.Sprintf(
					"%s:%d:%d: unexpected non-declaration statement on global scope",
					pos.Filename,
					pos.Line,
					pos.Column,
				),
			}
			parser.Collector.ReportAndSave(unexpectedTokenOnGlobalScope)
			return nil, collector.COMPILER_ERROR_FOUND
		}
	}
	return astNodes, nil
}

// Useful for testing
func ParseExprFrom(expr, filename string) (ast.Expr, error) {
	diagCollector := collector.New()

	reader := bufio.NewReader(strings.NewReader(expr))
	lex := lexer.New(filename, reader, diagCollector)

	tokens, err := lex.Tokenize()
	if err != nil {
		return nil, err
	}

	parser := New(tokens, diagCollector)
	exprAst, err := parser.parseExpr()
	if err != nil {
		return nil, err
	}
	return exprAst, nil
}

func (parser *parser) parseExternDecl() (*ast.ExternDecl, error) {
	_, ok := parser.expect(kind.EXTERN)
	// TODO(errors): should never hit
	if !ok {
		return nil, fmt.Errorf("expected 'extern'")
	}

	name, ok := parser.expect(kind.ID)
	if !ok {
		pos := name.Position
		expectedName := collector.Diag{
			Message: fmt.Sprintf(
				"%s:%d:%d: expected name, not %s",
				pos.Filename,
				pos.Line,
				pos.Column,
				name.Kind,
			),
		}
		parser.Collector.ReportAndSave(expectedName)
		return nil, collector.COMPILER_ERROR_FOUND
	}

	openCurly, ok := parser.expect(kind.OPEN_CURLY)
	if !ok {
		pos := openCurly.Position
		expectedOpenCurly := collector.Diag{
			Message: fmt.Sprintf(
				"%s:%d:%d: expected {, not %s",
				pos.Filename,
				pos.Line,
				pos.Column,
				openCurly.Kind,
			),
		}
		parser.Collector.ReportAndSave(expectedOpenCurly)
		return nil, collector.COMPILER_ERROR_FOUND
	}

	var prototypes []*ast.Proto
	for {
		if parser.cursor.nextIs(kind.CLOSE_CURLY) {
			break
		}

		proto, err := parser.parsePrototype()
		if err != nil {
			return nil, err
		}
		prototypes = append(prototypes, proto)
	}

	closeCurly, ok := parser.expect(kind.CLOSE_CURLY)
	// QUESTION: will it ever false?
	if !ok {
		pos := closeCurly.Position
		expectedCloseCurly := collector.Diag{
			Message: fmt.Sprintf(
				"%s:%d:%d: expected }, not %s",
				pos.Filename,
				pos.Line,
				pos.Column,
				closeCurly.Kind,
			),
		}
		parser.Collector.ReportAndSave(expectedCloseCurly)
		return nil, collector.COMPILER_ERROR_FOUND
	}

	return &ast.ExternDecl{Scope: nil, Name: name, Prototypes: prototypes}, nil
}

func (parser *parser) parsePrototype() (*ast.Proto, error) {
	fn, ok := parser.expect(kind.FN)
	if !ok {
		pos := fn.Position
		expectedCloseCurly := collector.Diag{
			Message: fmt.Sprintf(
				"%s:%d:%d: expected prototype or }, not %s",
				pos.Filename,
				pos.Line,
				pos.Column,
				fn.Kind,
			),
		}
		parser.Collector.ReportAndSave(expectedCloseCurly)
		return nil, collector.COMPILER_ERROR_FOUND
	}

	name, ok := parser.expect(kind.ID)
	if !ok {
		pos := name.Position
		expectedName := collector.Diag{
			Message: fmt.Sprintf(
				"%s:%d:%d: expected name, not %s",
				pos.Filename,
				pos.Line,
				pos.Column,
				name.Kind,
			),
		}
		parser.Collector.ReportAndSave(expectedName)
		return nil, collector.COMPILER_ERROR_FOUND
	}

	params, err := parser.parseFunctionParams()
	if err != nil {
		return nil, err
	}

	returnType, err := parser.parseReturnType( /*isPrototype=*/ true)
	if err != nil {
		return nil, err
	}

	semicolon, ok := parser.expect(kind.SEMICOLON)
	if !ok {
		pos := semicolon.Position
		expectedSemicolon := collector.Diag{
			Message: fmt.Sprintf(
				"%s:%d:%d: expected ; at the end of prototype, not %s",
				pos.Filename,
				pos.Line,
				pos.Column,
				semicolon.Kind,
			),
		}
		parser.Collector.ReportAndSave(expectedSemicolon)
		return nil, collector.COMPILER_ERROR_FOUND
	}

	return &ast.Proto{Name: name, Params: params, RetType: returnType}, nil
}

func (parser *parser) parseFnDecl() (*ast.FunctionDecl, error) {
	var err error

	_, ok := parser.expect(kind.FN)
	// TODO(errors): this should never hit
	if !ok {
		return nil, fmt.Errorf("expected 'fn'")
	}

	name, ok := parser.expect(kind.ID)
	if !ok {
		pos := name.Position
		expectedIdentifier := collector.Diag{
			Message: fmt.Sprintf(
				"%s:%d:%d: expected name, not %s",
				pos.Filename,
				pos.Line,
				pos.Column,
				name.Kind,
			),
		}
		parser.Collector.ReportAndSave(expectedIdentifier)
		return nil, collector.COMPILER_ERROR_FOUND
	}

	params, err := parser.parseFunctionParams()
	if err != nil {
		return nil, err
	}

	returnType, err := parser.parseReturnType( /*isPrototype=*/ false)
	if err != nil {
		return nil, err
	}

	block, err := parser.parseBlock()
	if err != nil {
		return nil, err
	}

	fnDecl := ast.FunctionDecl{
		Scope:   nil,
		Name:    name,
		Params:  params,
		Block:   block,
		RetType: returnType,
	}
	return &fnDecl, nil
}

func parseFnDeclFrom(filename, input string) (*ast.FunctionDecl, error) {
	diagCollector := collector.New()

	reader := bufio.NewReader(strings.NewReader(input))

	lexer := lexer.New(filename, reader, diagCollector)
	tokens, err := lexer.Tokenize()
	if err != nil {
		return nil, err
	}

	parser := New(tokens, diagCollector)
	fnDecl, err := parser.parseFnDecl()
	if err != nil {
		return nil, err
	}

	return fnDecl, nil
}

func (parser *parser) parseFunctionParams() (*ast.FieldList, error) {
	var params []*ast.Field
	isVariadic := false

	openParen, ok := parser.expect(kind.OPEN_PAREN)
	if !ok {
		pos := openParen.Position
		expectedOpenParen := collector.Diag{
			Message: fmt.Sprintf(
				"%s:%d:%d: expected (, not %s",
				pos.Filename,
				pos.Line,
				pos.Column,
				openParen.Kind,
			),
		}
		parser.Collector.ReportAndSave(expectedOpenParen)
		return nil, collector.COMPILER_ERROR_FOUND
	}

	for {
		if parser.cursor.nextIs(kind.CLOSE_PAREN) {
			break
		}
		if parser.cursor.nextIs(kind.DOT_DOT_DOT) {
			isVariadic = true
			tok := parser.cursor.peek()
			pos := tok.Position
			parser.cursor.skip()

			if !parser.cursor.nextIs(kind.CLOSE_PAREN) {
				// TODO(errors):
				// "fn name(a int, ...,) {}" because of the comma
				unexpectedDotDotDot := collector.Diag{
					Message: fmt.Sprintf(
						"%s:%d:%d: ... is only allowed at the end of parameter list",
						pos.Filename,
						pos.Line,
						pos.Column,
					),
				}
				parser.Collector.ReportAndSave(unexpectedDotDotDot)
				return nil, collector.COMPILER_ERROR_FOUND
			}
			break
		}

		name, ok := parser.expect(kind.ID)
		if !ok {
			pos := name.Position
			expectedCloseParenOrId := collector.Diag{
				Message: fmt.Sprintf(
					"%s:%d:%d: expected parameter or ), not %s",
					pos.Filename,
					pos.Line,
					pos.Column,
					name.Kind,
				),
			}
			parser.Collector.ReportAndSave(expectedCloseParenOrId)
			return nil, collector.COMPILER_ERROR_FOUND
		}
		paramType, err := parser.parseExprType()
		if err != nil {
			tok := parser.cursor.peek()
			pos := tok.Position
			expectedParamType := collector.Diag{
				Message: fmt.Sprintf(
					"%s:%d:%d: expected parameter type for '%s', not %s",
					pos.Filename,
					pos.Line,
					pos.Column,
					name.Lexeme,
					tok.Kind,
				),
			}
			parser.Collector.ReportAndSave(expectedParamType)
			return nil, collector.COMPILER_ERROR_FOUND
		}
		params = append(params, &ast.Field{Name: name, Type: paramType})

		if parser.cursor.nextIs(kind.COMMA) {
			parser.cursor.skip() // ,
			continue
		}
	}

	closeParen, ok := parser.expect(kind.CLOSE_PAREN)
	if !ok {
		pos := closeParen.Position
		expectedCloseParen := collector.Diag{
			Message: fmt.Sprintf(
				"%s:%d:%d: expected ), not %s",
				pos.Filename,
				pos.Line,
				pos.Column,
				closeParen.Kind,
			),
		}
		parser.Collector.ReportAndSave(expectedCloseParen)
		return nil, collector.COMPILER_ERROR_FOUND
	}

	return &ast.FieldList{
		Open:       openParen,
		Fields:     params,
		Close:      closeParen,
		IsVariadic: isVariadic,
	}, nil
}

func (parser *parser) parseReturnType(isPrototype bool) (ast.ExprType, error) {
	if (isPrototype && parser.cursor.nextIs(kind.SEMICOLON)) ||
		parser.cursor.nextIs(kind.OPEN_CURLY) {
		return &ast.BasicType{Kind: kind.VOID_TYPE}, nil
	}

	returnType, err := parser.parseExprType()
	if err != nil {
		tok := parser.cursor.peek()
		pos := tok.Position
		expectedReturnTy := collector.Diag{
			Message: fmt.Sprintf(
				"%s:%d:%d: expected type or {, not %s",
				pos.Filename,
				pos.Line,
				pos.Column,
				tok.Kind,
			),
		}
		parser.Collector.ReportAndSave(expectedReturnTy)
		return nil, collector.COMPILER_ERROR_FOUND
	}
	return returnType, nil
}

func (parser *parser) expect(expectedKind kind.TokenKind) (*token.Token, bool) {
	token := parser.cursor.peek()
	if token.Kind != expectedKind {
		return token, false
	}
	parser.cursor.skip()
	return token, true
}

func (parser *parser) parseExprType() (ast.ExprType, error) {
	token := parser.cursor.peek()
	switch token.Kind {
	case kind.STAR:
		parser.cursor.skip() // *
		ty, err := parser.parseExprType()
		if err != nil {
			return nil, err
		}
		return &ast.PointerType{Type: ty}, nil
	case kind.ID:
		parser.cursor.skip()
		return &ast.IdType{Name: token}, nil
	default:
		if token.Kind.IsBasicType() {
			parser.cursor.skip()
			return &ast.BasicType{Kind: token.Kind}, nil
		}
		return nil, collector.COMPILER_ERROR_FOUND
	}
}

func (parser *parser) parseStmt() (ast.Stmt, error) {
	token := parser.cursor.peek()
	switch token.Kind {
	case kind.RETURN:
		parser.cursor.skip()
		returnStmt := &ast.ReturnStmt{Return: token, Value: &ast.VoidExpr{}}
		if parser.cursor.nextIs(kind.SEMICOLON) {
			parser.cursor.skip()
			return returnStmt, nil
		}
		returnValue, err := parser.parseExpr()
		if err != nil {
			tok := parser.cursor.peek()
			pos := tok.Position
			expectedSemicolon := collector.Diag{
				Message: fmt.Sprintf(
					"%s:%d:%d: expected expression or ;, not %s",
					pos.Filename,
					pos.Line,
					pos.Column,
					tok.Kind,
				),
			}
			parser.Collector.ReportAndSave(expectedSemicolon)
			return nil, collector.COMPILER_ERROR_FOUND
		}
		returnStmt.Value = returnValue

		_, ok := parser.expect(kind.SEMICOLON)
		if !ok {
			tok := parser.cursor.peek()
			pos := tok.Position
			expectedSemicolon := collector.Diag{
				Message: fmt.Sprintf(
					"%s:%d:%d: expected ; at the end of statement, not %s",
					pos.Filename,
					pos.Line,
					pos.Column,
					tok.Kind,
				),
			}
			parser.Collector.ReportAndSave(expectedSemicolon)
			return nil, collector.COMPILER_ERROR_FOUND
		}

		return returnStmt, nil
	case kind.ID:
		idStmt, err := parser.ParseIdStmt()
		if err != nil {
			return nil, err
		}
		semicolon, ok := parser.expect(kind.SEMICOLON)
		if !ok {
			pos := semicolon.Position
			expectedSemicolon := collector.Diag{
				Message: fmt.Sprintf(
					"%s:%d:%d: expected ; at the end of statement, not %s",
					pos.Filename,
					pos.Line,
					pos.Column,
					semicolon.Kind,
				),
			}
			parser.Collector.ReportAndSave(expectedSemicolon)
			return nil, collector.COMPILER_ERROR_FOUND
		}
		// TODO(errors)
		return idStmt, err
	case kind.IF:
		condStmt, err := parser.parseCondStmt()
		// TODO(errors)
		return condStmt, err
	case kind.FOR:
		forLoop, err := parser.parseForLoop()
		// TODO(errors)
		return forLoop, err
	default:
		return nil, nil
	}
}

func (parser *parser) parseBlock() (*ast.BlockStmt, error) {
	openCurly, ok := parser.expect(kind.OPEN_CURLY)
	// TODO(errors): should never hit
	if !ok {
		return nil, fmt.Errorf("expected '{', but got %s", openCurly)
	}

	var statements []ast.Stmt

	for {
		token := parser.cursor.peek()
		if token.Kind == kind.CLOSE_CURLY {
			break
		}

		stmt, err := parser.parseStmt()
		// TODO(errors)
		if err != nil {
			return nil, err
		}

		if stmt != nil {
			statements = append(statements, stmt)
		} else {
			break
		}
	}

	closeCurly, ok := parser.expect(kind.CLOSE_CURLY)
	if !ok {
		pos := closeCurly.Position
		expectedStatementOrCloseCurly := collector.Diag{
			Message: fmt.Sprintf(
				"%s:%d:%d: expected statement or }, not %s",
				pos.Filename,
				pos.Line,
				pos.Column,
				closeCurly.Kind,
			),
		}
		parser.Collector.ReportAndSave(expectedStatementOrCloseCurly)
		return nil, collector.COMPILER_ERROR_FOUND
	}

	return &ast.BlockStmt{
		OpenCurly:  openCurly.Position,
		Statements: statements,
		CloseCurly: closeCurly.Position,
	}, nil
}

func (parser *parser) ParseIdStmt() (ast.Stmt, error) {
	aheadId := parser.cursor.peekN(1)
	switch aheadId.Kind {
	case kind.OPEN_PAREN:
		fnCall, err := parser.parseFnCall()
		return fnCall, err
	case kind.DOT:
		fieldAccessing := parser.parseFieldAccess()
		return fieldAccessing, nil
	default:
		return parser.parseVar()
	}
}

func (parser *parser) parseVar() (ast.Stmt, error) {
	variables := make([]*ast.VarStmt, 0)
	isDecl := false

VarDecl:
	for {
		name, ok := parser.expect(kind.ID)
		// TODO(errors)
		if !ok {
			return nil, fmt.Errorf("expected ID")
		}

		variable := &ast.VarStmt{
			Name:           name,
			Type:           nil,
			Value:          nil,
			Decl:           true,
			NeedsInference: true,
		}
		variables = append(variables, variable)

		next := parser.cursor.peek()
		switch next.Kind {
		case kind.COLON_EQUAL, kind.EQUAL:
			parser.cursor.skip() // := or =
			isDecl = next.Kind == kind.COLON_EQUAL
			break VarDecl
		case kind.COMMA:
			parser.cursor.skip()
			continue
		}

		ty, err := parser.parseExprType()
		// TODO(errors)
		if err != nil {
			return nil, err
		}
		variable.Type = ty
		variable.NeedsInference = false

		next = parser.cursor.peek()
		switch next.Kind {
		case kind.COLON_EQUAL, kind.EQUAL:
			parser.cursor.skip() // := or =
			isDecl = next.Kind == kind.COLON_EQUAL
			break VarDecl
		case kind.COMMA:
			parser.cursor.skip()
			continue
		}
	}

	exprs, err := parser.parseExprList([]kind.TokenKind{kind.SEMICOLON, kind.CLOSE_PAREN})
	// TODO(errors)
	if err != nil {
		return nil, err
	}

	// TODO(errors)
	if len(variables) != len(exprs) {
		return nil, fmt.Errorf("%d != %d", len(variables), len(exprs))
	}
	for i := range variables {
		variables[i].Value = exprs[i]
		variables[i].Decl = isDecl
	}

	if len(variables) == 1 {
		return variables[0], nil
	}
	return &ast.MultiVarStmt{IsDecl: isDecl, Variables: variables}, nil
}

// Useful for testing
func parseVarFrom(filename, input string) (ast.Stmt, error) {
	diagCollector := collector.New()

	reader := bufio.NewReader(strings.NewReader(input))
	lex := lexer.New(filename, reader, diagCollector)

	tokens, err := lex.Tokenize()
	if err != nil {
		return nil, err
	}

	parser := New(tokens, diagCollector)
	stmt, err := parser.ParseIdStmt()
	if err != nil {
		return nil, err
	}
	return stmt, nil
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
		elifConds = append(
			elifConds,
			&ast.IfElifCond{If: &elifToken.Position, Expr: elifExpr, Block: elifBlock},
		)
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
	if _, ok := ast.UNARY[next.Kind]; ok {
		parser.cursor.skip()
		rhs, err := parser.parseUnary()
		// TODO(errors)
		if err != nil {
			return nil, err
		}
		return &ast.UnaryExpr{Op: next.Kind, Value: rhs}, nil
	}

	return parser.parsePrimary()
}

func (parser *parser) parsePrimary() (ast.Expr, error) {
	token := parser.cursor.peek()
	switch token.Kind {
	case kind.ID:
		idExpr := &ast.IdExpr{Name: token}

		next := parser.cursor.peekN(1)
		switch next.Kind {
		case kind.OPEN_PAREN:
			return parser.parseFnCall()
		case kind.DOT:
			fieldAccess := parser.parseFieldAccess()
			return fieldAccess, nil
		}

		parser.cursor.skip()
		return idExpr, nil
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
			return &ast.LiteralExpr{
				Type:  &ast.BasicType{Kind: token.Kind},
				Value: token.Lexeme,
			}, nil
		}
		return nil, fmt.Errorf(
			"invalid token for expression parsing: %s %s %s",
			token.Kind,
			token.Lexeme,
			token.Position,
		)
	}
}

func (parser *parser) parseExprList(possibleEnds []kind.TokenKind) ([]ast.Expr, error) {
	var exprs []ast.Expr
Var:
	for {
		for _, end := range possibleEnds {
			if parser.cursor.nextIs(end) {
				break Var
			}
		}

		expr, err := parser.parseExpr()
		if err != nil {
			return nil, err
		}
		exprs = append(exprs, expr)

		if parser.cursor.nextIs(kind.COMMA) {
			parser.cursor.skip()
			continue
		}
	}
	return exprs, nil
}

func (parser *parser) parseFnCall() (*ast.FunctionCall, error) {
	name, ok := parser.expect(kind.ID)
	// TODO(errors): should never hit
	if !ok {
		return nil, fmt.Errorf("expected 'id'")
	}

	_, ok = parser.expect(kind.OPEN_PAREN)
	// TODO(errors)
	if !ok {
		return nil, fmt.Errorf("expected '('")
	}

	args, err := parser.parseExprList([]kind.TokenKind{kind.CLOSE_PAREN})
	if err != nil {
		return nil, err
	}

	_, ok = parser.expect(kind.CLOSE_PAREN)
	// TODO(errors)
	if !ok {
		return nil, fmt.Errorf("expected ')'")
	}

	return &ast.FunctionCall{Name: name, Args: args}, nil
}

func (parser *parser) parseFieldAccess() *ast.FieldAccess {
	id, ok := parser.expect(kind.ID)
	// TODO(errors): should never hit
	if !ok {
		log.Fatal("expected ID")
	}
	left := &ast.IdExpr{Name: id}

	_, ok = parser.expect(kind.DOT)
	// TODO(errors): should never hit
	if !ok {
		log.Fatal("expect a dot")
	}

	right, err := parser.parsePrimary()
	// TODO(errors)
	if err != nil {
		log.Fatal(err)
	}
	return &ast.FieldAccess{Left: left, Right: right}
}

func (parser *parser) parseForLoop() (*ast.ForLoop, error) {
	_, ok := parser.expect(kind.FOR)
	// TODO(errors): should never hit
	if !ok {
		return nil, fmt.Errorf("expected 'for'")
	}

	_, ok = parser.expect(kind.OPEN_PAREN)
	// TODO(errors): should never hit
	if !ok {
		return nil, fmt.Errorf("expected '('")
	}

	init, err := parser.parseVar()
	// TODO(errors)
	if err != nil {
		return nil, err
	}

	_, ok = parser.expect(kind.SEMICOLON)
	// TODO(errors)
	if !ok {
		return nil, fmt.Errorf("expected ';'")
	}

	cond, err := parser.parseExpr()
	// TODO(errors)
	if err != nil {
		return nil, err
	}

	_, ok = parser.expect(kind.SEMICOLON)
	// TODO(errors)
	if !ok {
		return nil, fmt.Errorf("expected ';'")
	}

	update, err := parser.parseVar()
	// TODO(errors)
	if err != nil {
		return nil, err
	}

	_, ok = parser.expect(kind.CLOSE_PAREN)
	// TODO(errors): should never hit
	if !ok {
		return nil, fmt.Errorf("expected ')'")
	}

	block, err := parser.parseBlock()
	if err != nil {
		return nil, err
	}
	return &ast.ForLoop{Init: init, Cond: cond, Update: update, Block: block}, nil
}

// Useful for testing
func ParseForLoopFrom(input, filename string) (*ast.ForLoop, error) {
	diagCollector := collector.New()

	reader := bufio.NewReader(strings.NewReader(input))
	lex := lexer.New(filename, reader, diagCollector)

	tokens, err := lex.Tokenize()
	if err != nil {
		return nil, err
	}

	parser := New(tokens, diagCollector)
	forLoop, err := parser.parseForLoop()
	return forLoop, err
}
