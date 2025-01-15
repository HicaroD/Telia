package parser

import (
	"fmt"
	"log"

	"github.com/HicaroD/Telia/diagnostics"
	"github.com/HicaroD/Telia/frontend/ast"
	"github.com/HicaroD/Telia/frontend/lexer"
	"github.com/HicaroD/Telia/frontend/lexer/token"
	"github.com/HicaroD/Telia/scope"
)

type Parser struct {
	lex       *lexer.Lexer
	collector *diagnostics.Collector

	moduleScope *scope.Scope[ast.Node] // current module scope being analyzed
}

func New(lex *lexer.Lexer, collector *diagnostics.Collector) *Parser {
	return &Parser{lex: lex, collector: collector}
}

func (p *Parser) Parse() (*ast.Program, error) {
	// NOTE: temporary change for handling the new architecture
	// I will have many modules (std included), many files
	// inside modules and more

	// Universe scope has a nil parent
	universe := scope.New[ast.Node](nil)

	// TODO: properly set module name
	// NOTE: not every module will have the universe as the
	// parent, for inner-modules we could have other modules as parents
	moduleScope := scope.New(universe)
	// NOTE: every time a new module is being parsed, this variable should be
	// set appropriately, deal with it carefuly
	p.moduleScope = moduleScope

	file, err := p.parseFile()
	if err != nil {
		return nil, err
	}
	file.Scope = scope.New(moduleScope)

	module := &ast.Module{Scope: moduleScope, Name: "main", Body: []*ast.File{file}}
	program := &ast.Program{Universe: universe, Body: []*ast.Module{module}}

	return program, nil
}

func (p *Parser) parseFile() (*ast.File, error) {
	var nodes []ast.Node
	for {
		node, eof, err := p.Next()
		if err != nil {
			return nil, err
		}
		if eof {
			break
		}
		nodes = append(nodes, node)
	}
	// TODO: properly set filename here
	file := &ast.File{Name: "my_program", Body: nodes}
	return file, nil
}

func (p *Parser) Next() (ast.Node, bool, error) {
	eof := false

	tok := p.lex.Peek()
	if tok.Kind == token.EOF {
		eof = true
		return nil, eof, nil
	}
	switch tok.Kind {
	case token.FN:
		fnDecl, err := p.parseFnDecl()
		return fnDecl, eof, err
	case token.EXTERN:
		externDecl, err := p.parseExternDecl()
		return externDecl, eof, err
	default:
		pos := tok.Pos
		unexpectedTokenOnGlobalScope := diagnostics.Diag{
			Message: fmt.Sprintf(
				"%s:%d:%d: unexpected non-declaration statement on global scope",
				pos.Filename,
				pos.Line,
				pos.Column,
			),
		}
		p.collector.ReportAndSave(unexpectedTokenOnGlobalScope)
		return nil, eof, diagnostics.COMPILER_ERROR_FOUND
	}
}

// Useful for testing
func ParseExprFrom(expr, filename string) (ast.Expr, error) {
	collector := diagnostics.New()

	src := []byte(expr)
	lex := lexer.New(filename, src, collector)
	parser := New(lex, collector)

	exprAst, err := parser.parseExpr()
	if err != nil {
		return nil, err
	}
	return exprAst, nil
}

func (p *Parser) parseExternDecl() (*ast.ExternDecl, error) {
	_, ok := p.expect(token.EXTERN)
	if !ok {
		return nil, fmt.Errorf("expected 'extern'")
	}

	name, ok := p.expect(token.ID)
	if !ok {
		pos := name.Pos
		expectedName := diagnostics.Diag{
			Message: fmt.Sprintf(
				"%s:%d:%d: expected name, not %s",
				pos.Filename,
				pos.Line,
				pos.Column,
				name.Kind,
			),
		}
		p.collector.ReportAndSave(expectedName)
		return nil, diagnostics.COMPILER_ERROR_FOUND
	}

	openCurly, ok := p.expect(token.OPEN_CURLY)
	if !ok {
		pos := openCurly.Pos
		expectedOpenCurly := diagnostics.Diag{
			Message: fmt.Sprintf(
				"%s:%d:%d: expected {, not %s",
				pos.Filename,
				pos.Line,
				pos.Column,
				openCurly.Kind,
			),
		}
		p.collector.ReportAndSave(expectedOpenCurly)
		return nil, diagnostics.COMPILER_ERROR_FOUND
	}

	var prototypes []*ast.Proto
	for {
		if p.lex.NextIs(token.CLOSE_CURLY) {
			break
		}

		proto, err := p.parsePrototype()
		if err != nil {
			return nil, err
		}
		prototypes = append(prototypes, proto)
	}

	closeCurly, ok := p.expect(token.CLOSE_CURLY)
	// QUESTION: will it ever false?
	if !ok {
		pos := closeCurly.Pos
		expectedCloseCurly := diagnostics.Diag{
			Message: fmt.Sprintf(
				"%s:%d:%d: expected }, not %s",
				pos.Filename,
				pos.Line,
				pos.Column,
				closeCurly.Kind,
			),
		}
		p.collector.ReportAndSave(expectedCloseCurly)
		return nil, diagnostics.COMPILER_ERROR_FOUND
	}

	return &ast.ExternDecl{Scope: nil, Name: name, Prototypes: prototypes}, nil
}

func (p *Parser) parsePrototype() (*ast.Proto, error) {
	fn, ok := p.expect(token.FN)
	if !ok {
		pos := fn.Pos
		expectedCloseCurly := diagnostics.Diag{
			Message: fmt.Sprintf(
				"%s:%d:%d: expected prototype or }, not %s",
				pos.Filename,
				pos.Line,
				pos.Column,
				fn.Kind,
			),
		}
		p.collector.ReportAndSave(expectedCloseCurly)
		return nil, diagnostics.COMPILER_ERROR_FOUND
	}

	name, ok := p.expect(token.ID)
	if !ok {
		pos := name.Pos
		expectedName := diagnostics.Diag{
			Message: fmt.Sprintf(
				"%s:%d:%d: expected name, not %s",
				pos.Filename,
				pos.Line,
				pos.Column,
				name.Kind,
			),
		}
		p.collector.ReportAndSave(expectedName)
		return nil, diagnostics.COMPILER_ERROR_FOUND
	}

	params, err := p.parseFunctionParams()
	if err != nil {
		return nil, err
	}

	returnType, err := p.parseReturnType( /*isPrototype=*/ true)
	if err != nil {
		return nil, err
	}

	semicolon, ok := p.expect(token.SEMICOLON)
	if !ok {
		pos := semicolon.Pos
		expectedSemicolon := diagnostics.Diag{
			Message: fmt.Sprintf(
				"%s:%d:%d: expected ; at the end of prototype, not %s",
				pos.Filename,
				pos.Line,
				pos.Column,
				semicolon.Kind,
			),
		}
		p.collector.ReportAndSave(expectedSemicolon)
		return nil, diagnostics.COMPILER_ERROR_FOUND
	}

	return &ast.Proto{Name: name, Params: params, RetType: returnType}, nil
}

func (p *Parser) parseFnDecl() (*ast.FunctionDecl, error) {
	var err error

	_, ok := p.expect(token.FN)
	// TODO(errors): this should never hit
	if !ok {
		return nil, fmt.Errorf("expected 'fn'")
	}

	name, ok := p.expect(token.ID)
	if !ok {
		pos := name.Pos
		expectedIdentifier := diagnostics.Diag{
			Message: fmt.Sprintf(
				"%s:%d:%d: expected name, not %s",
				pos.Filename,
				pos.Line,
				pos.Column,
				name.Kind,
			),
		}
		p.collector.ReportAndSave(expectedIdentifier)
		return nil, diagnostics.COMPILER_ERROR_FOUND
	}

	params, err := p.parseFunctionParams()
	if err != nil {
		return nil, err
	}

	returnType, err := p.parseReturnType( /*isPrototype=*/ false)
	if err != nil {
		return nil, err
	}

	block, err := p.parseBlock()
	if err != nil {
		return nil, err
	}

	fnDecl := &ast.FunctionDecl{
		Scope:   nil,
		Name:    name,
		Params:  params,
		Block:   block,
		RetType: returnType,
	}

	err = p.moduleScope.Insert(name.Name(), fnDecl)
	if err != nil {
		return nil, err
	}
	return fnDecl, nil
}

func parseFnDeclFrom(filename, input string) (*ast.FunctionDecl, error) {
	universeScope := scope.New[ast.Node](nil)
	moduleScope := scope.New(universeScope)

	collector := diagnostics.New()

	src := []byte(input)
	lexer := lexer.New(filename, src, collector)
	parser := New(lexer, collector)
	parser.moduleScope = moduleScope

	fnDecl, err := parser.parseFnDecl()
	if err != nil {
		return nil, err
	}

	return fnDecl, nil
}

func (p *Parser) parseFunctionParams() (*ast.FieldList, error) {
	var params []*ast.Field
	isVariadic := false

	openParen, ok := p.expect(token.OPEN_PAREN)
	if !ok {
		pos := openParen.Pos
		expectedOpenParen := diagnostics.Diag{
			Message: fmt.Sprintf(
				"%s:%d:%d: expected (, not %s",
				pos.Filename,
				pos.Line,
				pos.Column,
				openParen.Kind,
			),
		}
		p.collector.ReportAndSave(expectedOpenParen)
		return nil, diagnostics.COMPILER_ERROR_FOUND
	}

	for {
		if p.lex.NextIs(token.CLOSE_PAREN) {
			break
		}
		if p.lex.NextIs(token.DOT_DOT_DOT) {
			isVariadic = true
			tok := p.lex.Peek()
			pos := tok.Pos
			p.lex.Skip()

			if !p.lex.NextIs(token.CLOSE_PAREN) {
				// TODO(errors):
				// "fn name(a int, ...,) {}" because of the comma
				unexpectedDotDotDot := diagnostics.Diag{
					Message: fmt.Sprintf(
						"%s:%d:%d: ... is only allowed at the end of parameter list",
						pos.Filename,
						pos.Line,
						pos.Column,
					),
				}
				p.collector.ReportAndSave(unexpectedDotDotDot)
				return nil, diagnostics.COMPILER_ERROR_FOUND
			}
			break
		}

		name, ok := p.expect(token.ID)
		if !ok {
			pos := name.Pos
			expectedCloseParenOrId := diagnostics.Diag{
				Message: fmt.Sprintf(
					"%s:%d:%d: expected parameter or ), not %s",
					pos.Filename,
					pos.Line,
					pos.Column,
					name.Kind,
				),
			}
			p.collector.ReportAndSave(expectedCloseParenOrId)
			return nil, diagnostics.COMPILER_ERROR_FOUND
		}
		paramType, err := p.parseExprType()
		if err != nil {
			tok := p.lex.Peek()
			pos := tok.Pos
			expectedParamType := diagnostics.Diag{
				Message: fmt.Sprintf(
					"%s:%d:%d: expected parameter type for '%s', not %s",
					pos.Filename,
					pos.Line,
					pos.Column,
					name.Lexeme,
					tok.Kind,
				),
			}
			p.collector.ReportAndSave(expectedParamType)
			return nil, diagnostics.COMPILER_ERROR_FOUND
		}
		params = append(params, &ast.Field{Name: name, Type: paramType})

		if p.lex.NextIs(token.COMMA) {
			p.lex.Skip() // ,
			continue
		}
	}

	closeParen, ok := p.expect(token.CLOSE_PAREN)
	if !ok {
		pos := closeParen.Pos
		expectedCloseParen := diagnostics.Diag{
			Message: fmt.Sprintf(
				"%s:%d:%d: expected ), not %s",
				pos.Filename,
				pos.Line,
				pos.Column,
				closeParen.Kind,
			),
		}
		p.collector.ReportAndSave(expectedCloseParen)
		return nil, diagnostics.COMPILER_ERROR_FOUND
	}

	return &ast.FieldList{
		Open:       openParen,
		Fields:     params,
		Close:      closeParen,
		IsVariadic: isVariadic,
	}, nil
}

func (p *Parser) parseReturnType(isPrototype bool) (ast.ExprType, error) {
	if (isPrototype && p.lex.NextIs(token.SEMICOLON)) ||
		p.lex.NextIs(token.OPEN_CURLY) {
		return &ast.BasicType{Kind: token.VOID_TYPE}, nil
	}

	returnType, err := p.parseExprType()
	if err != nil {
		tok := p.lex.Peek()
		pos := tok.Pos
		expectedReturnTy := diagnostics.Diag{
			Message: fmt.Sprintf(
				"%s:%d:%d: expected type or {, not %s",
				pos.Filename,
				pos.Line,
				pos.Column,
				tok.Kind,
			),
		}
		p.collector.ReportAndSave(expectedReturnTy)
		return nil, diagnostics.COMPILER_ERROR_FOUND
	}
	return returnType, nil
}

func (p *Parser) expect(expectedKind token.Kind) (*token.Token, bool) {
	tok := p.lex.Peek()
	if tok.Kind != expectedKind {
		return tok, false
	}
	p.lex.Skip()
	return tok, true
}

func (p *Parser) parseExprType() (ast.ExprType, error) {
	tok := p.lex.Peek()
	switch tok.Kind {
	case token.STAR:
		p.lex.Skip() // *
		ty, err := p.parseExprType()
		if err != nil {
			return nil, err
		}
		return &ast.PointerType{Type: ty}, nil
	case token.ID:
		p.lex.Skip()
		return &ast.IdType{Name: tok}, nil
	default:
		if tok.Kind.IsBasicType() {
			p.lex.Skip()
			return &ast.BasicType{Kind: tok.Kind}, nil
		}
		return nil, diagnostics.COMPILER_ERROR_FOUND
	}
}

func (p *Parser) parseStmt() (ast.Stmt, error) {
	tok := p.lex.Peek()
	switch tok.Kind {
	case token.RETURN:
		p.lex.Skip()
		returnStmt := &ast.ReturnStmt{Return: tok, Value: &ast.VoidExpr{}}
		if p.lex.NextIs(token.SEMICOLON) {
			p.lex.Skip()
			return returnStmt, nil
		}
		returnValue, err := p.parseExpr()
		if err != nil {
			tok := p.lex.Peek()
			pos := tok.Pos
			expectedSemicolon := diagnostics.Diag{
				Message: fmt.Sprintf(
					"%s:%d:%d: expected expression or ;, not %s",
					pos.Filename,
					pos.Line,
					pos.Column,
					tok.Kind,
				),
			}
			p.collector.ReportAndSave(expectedSemicolon)
			return nil, diagnostics.COMPILER_ERROR_FOUND
		}
		returnStmt.Value = returnValue

		_, ok := p.expect(token.SEMICOLON)
		if !ok {
			tok := p.lex.Peek()
			pos := tok.Pos
			expectedSemicolon := diagnostics.Diag{
				Message: fmt.Sprintf(
					"%s:%d:%d: expected ; at the end of statement, not %s",
					pos.Filename,
					pos.Line,
					pos.Column,
					tok.Kind,
				),
			}
			p.collector.ReportAndSave(expectedSemicolon)
			return nil, diagnostics.COMPILER_ERROR_FOUND
		}

		return returnStmt, nil
	case token.ID:
		idStmt, err := p.ParseIdStmt()
		if err != nil {
			return nil, err
		}
		semicolon, ok := p.expect(token.SEMICOLON)
		if !ok {
			pos := semicolon.Pos
			expectedSemicolon := diagnostics.Diag{
				Message: fmt.Sprintf(
					"%s:%d:%d: expected ; at the end of statement, not %s",
					pos.Filename,
					pos.Line,
					pos.Column,
					semicolon.Kind,
				),
			}
			p.collector.ReportAndSave(expectedSemicolon)
			return nil, diagnostics.COMPILER_ERROR_FOUND
		}
		// TODO(errors)
		return idStmt, err
	case token.IF:
		condStmt, err := p.parseCondStmt()
		// TODO(errors)
		return condStmt, err
	case token.FOR:
		forLoop, err := p.parseForLoop()
		return forLoop, err
	case token.WHILE:
		whileLoop, err := p.parseWhileLoop()
		return whileLoop, err
	default:
		return nil, nil
	}
}

func (p *Parser) parseBlock() (*ast.BlockStmt, error) {
	openCurly, ok := p.expect(token.OPEN_CURLY)
	// TODO(errors): should never hit
	if !ok {
		return nil, fmt.Errorf("expected '{', but got %s", openCurly)
	}

	var statements []ast.Stmt

	for {
		tok := p.lex.Peek()
		if tok.Kind == token.CLOSE_CURLY {
			break
		}

		stmt, err := p.parseStmt()
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

	closeCurly, ok := p.expect(token.CLOSE_CURLY)
	if !ok {
		pos := closeCurly.Pos
		expectedStatementOrCloseCurly := diagnostics.Diag{
			Message: fmt.Sprintf(
				"%s:%d:%d: expected statement or }, not %s",
				pos.Filename,
				pos.Line,
				pos.Column,
				closeCurly.Kind,
			),
		}
		p.collector.ReportAndSave(expectedStatementOrCloseCurly)
		return nil, diagnostics.COMPILER_ERROR_FOUND
	}

	return &ast.BlockStmt{
		OpenCurly:  openCurly.Pos,
		Statements: statements,
		CloseCurly: closeCurly.Pos,
	}, nil
}

func (p *Parser) ParseIdStmt() (ast.Stmt, error) {
	aheadId := p.lex.Peek1()
	switch aheadId.Kind {
	case token.OPEN_PAREN:
		fnCall, err := p.parseFnCall()
		return fnCall, err
	case token.DOT:
		fieldAccessing := p.parseFieldAccess()
		return fieldAccessing, nil
	default:
		return p.parseVar()
	}
}

func (p *Parser) parseVar() (ast.Stmt, error) {
	variables := make([]*ast.VarStmt, 0)
	isDecl := false

VarDecl:
	for {
		name, ok := p.expect(token.ID)
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

		next := p.lex.Peek()
		switch next.Kind {
		case token.COLON_EQUAL, token.EQUAL:
			p.lex.Skip() // := or =
			isDecl = next.Kind == token.COLON_EQUAL
			break VarDecl
		case token.COMMA:
			p.lex.Skip()
			continue
		}

		ty, err := p.parseExprType()
		// TODO(errors)
		if err != nil {
			return nil, err
		}
		variable.Type = ty
		variable.NeedsInference = false

		next = p.lex.Peek()
		switch next.Kind {
		case token.COLON_EQUAL, token.EQUAL:
			p.lex.Skip() // := or =
			isDecl = next.Kind == token.COLON_EQUAL
			break VarDecl
		case token.COMMA:
			p.lex.Skip()
			continue
		}
	}

	exprs, err := p.parseExprList([]token.Kind{token.SEMICOLON, token.CLOSE_PAREN})
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
	collector := diagnostics.New()

	src := []byte(input)
	lex := lexer.New(filename, src, collector)
	parser := New(lex, collector)

	stmt, err := parser.ParseIdStmt()
	if err != nil {
		return nil, err
	}
	return stmt, nil
}

func (p *Parser) parseCondStmt() (*ast.CondStmt, error) {
	ifCond, err := p.parseIfCond()
	// TODO(errors)
	if err != nil {
		return nil, err
	}

	elifConds, err := p.parseElifConds()
	// TODO(errors)
	if err != nil {
		return nil, err
	}

	elseCond, err := p.parseElseCond()
	// TODO(errors)
	if err != nil {
		return nil, err
	}

	return &ast.CondStmt{IfStmt: ifCond, ElifStmts: elifConds, ElseStmt: elseCond}, nil
}

func (p *Parser) parseIfCond() (*ast.IfElifCond, error) {
	ifToken, ok := p.expect(token.IF)
	// TODO(errors)
	if !ok {
		return nil, fmt.Errorf("expected 'if'")
	}

	// TODO(errors)
	ifExpr, err := p.parseExpr()
	if err != nil {
		return nil, err
	}

	ifBlock, err := p.parseBlock()
	// TODO(errors)
	if err != nil {
		return nil, err
	}
	return &ast.IfElifCond{If: &ifToken.Pos, Expr: ifExpr, Block: ifBlock}, nil
}

func (p *Parser) parseElifConds() ([]*ast.IfElifCond, error) {
	var elifConds []*ast.IfElifCond
	for {
		elifToken, ok := p.expect(token.ELIF)
		if !ok {
			break
		}
		elifExpr, err := p.parseExpr()
		// TODO(errors)
		if err != nil {
			return nil, err
		}
		elifBlock, err := p.parseBlock()
		// TODO(errors)
		if err != nil {
			return nil, err
		}
		elifConds = append(
			elifConds,
			&ast.IfElifCond{If: &elifToken.Pos, Expr: elifExpr, Block: elifBlock},
		)
	}
	return elifConds, nil
}

func (p *Parser) parseElseCond() (*ast.ElseCond, error) {
	elseToken, ok := p.expect(token.ELSE)
	if !ok {
		return nil, nil
	}

	elseBlock, err := p.parseBlock()
	if err != nil {
		return nil, err
	}
	return &ast.ElseCond{Else: &elseToken.Pos, Block: elseBlock}, nil
}

func (p *Parser) parseExpr() (ast.Expr, error) {
	return p.parseLogical()
}

func (p *Parser) parseLogical() (ast.Expr, error) {
	lhs, err := p.parseComparasion()
	// TODO(errors)
	if err != nil {
		return nil, err
	}

	// TODO: refactor this code, it seems there is a better way of writing this
	for {
		next := p.lex.Peek()
		if _, ok := ast.LOGICAL[next.Kind]; ok {
			p.lex.Skip()
			rhs, err := p.parseComparasion()
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

func (p *Parser) parseComparasion() (ast.Expr, error) {
	lhs, err := p.parseTerm()
	// TODO(errors)
	if err != nil {
		return nil, err
	}

	// TODO: refactor this code, it seems there is a better way of writing this
	for {
		next := p.lex.Peek()
		if _, ok := ast.COMPARASION[next.Kind]; ok {
			p.lex.Skip()
			rhs, err := p.parseTerm()
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

func (p *Parser) parseTerm() (ast.Expr, error) {
	lhs, err := p.parseFactor()
	// TODO(errors)
	if err != nil {
		return nil, err
	}

	// TODO: refactor this code, it seems there is a better way of writing this
	for {
		next := p.lex.Peek()
		if _, ok := ast.TERM[next.Kind]; ok {
			p.lex.Skip()
			rhs, err := p.parseFactor()
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

func (p *Parser) parseFactor() (ast.Expr, error) {
	lhs, err := p.parseUnary()
	// TODO(errors)
	if err != nil {
		return nil, err
	}

	// TODO: refactor this code, it seems there is a better way of writing this
	for {
		next := p.lex.Peek()
		if _, ok := ast.FACTOR[next.Kind]; ok {
			p.lex.Skip()
			rhs, err := p.parseUnary()
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

func (p *Parser) parseUnary() (ast.Expr, error) {
	next := p.lex.Peek()
	if _, ok := ast.UNARY[next.Kind]; ok {
		p.lex.Skip()
		rhs, err := p.parseUnary()
		// TODO(errors)
		if err != nil {
			return nil, err
		}
		return &ast.UnaryExpr{Op: next.Kind, Value: rhs}, nil
	}

	return p.parsePrimary()
}

func (p *Parser) parsePrimary() (ast.Expr, error) {
	tok := p.lex.Peek()
	switch tok.Kind {
	case token.ID:
		idExpr := &ast.IdExpr{Name: tok}

		next := p.lex.Peek1()
		switch next.Kind {
		case token.OPEN_PAREN:
			return p.parseFnCall()
		case token.DOT:
			fieldAccess := p.parseFieldAccess()
			return fieldAccess, nil
		}

		p.lex.Skip()
		return idExpr, nil
	case token.OPEN_PAREN:
		p.lex.Skip() // (
		expr, err := p.parseExpr()
		// TODO(errors)
		if err != nil {
			return nil, err
		}
		_, ok := p.expect(token.CLOSE_PAREN)
		if !ok {
			return nil, fmt.Errorf("expected closing parenthesis")
		}
		return expr, nil
	default:
		if _, ok := token.LITERAL_KIND[tok.Kind]; ok {
			p.lex.Skip()
			return &ast.LiteralExpr{
				Type:  &ast.BasicType{Kind: tok.Kind},
				Value: tok.Lexeme,
			}, nil
		}
		return nil, fmt.Errorf(
			"invalid token for expression parsing: %s %s %s",
			tok.Kind,
			tok.Lexeme,
			tok.Pos,
		)
	}
}

func (p *Parser) parseExprList(possibleEnds []token.Kind) ([]ast.Expr, error) {
	var exprs []ast.Expr
Var:
	for {
		for _, end := range possibleEnds {
			if p.lex.NextIs(end) {
				break Var
			}
		}

		expr, err := p.parseExpr()
		if err != nil {
			return nil, err
		}
		exprs = append(exprs, expr)

		if p.lex.NextIs(token.COMMA) {
			p.lex.Skip()
			continue
		}
	}
	return exprs, nil
}

func (parser *Parser) parseFnCall() (*ast.FunctionCall, error) {
	name, ok := parser.expect(token.ID)
	// TODO(errors): should never hit
	if !ok {
		return nil, fmt.Errorf("expected 'id'")
	}

	_, ok = parser.expect(token.OPEN_PAREN)
	// TODO(errors)
	if !ok {
		return nil, fmt.Errorf("expected '('")
	}

	args, err := parser.parseExprList([]token.Kind{token.CLOSE_PAREN})
	if err != nil {
		return nil, err
	}

	_, ok = parser.expect(token.CLOSE_PAREN)
	// TODO(errors)
	if !ok {
		return nil, fmt.Errorf("expected ')'")
	}

	return &ast.FunctionCall{Name: name, Args: args}, nil
}

func (parser *Parser) parseFieldAccess() *ast.FieldAccess {
	id, ok := parser.expect(token.ID)
	// TODO(errors): should never hit
	if !ok {
		log.Fatal("expected ID")
	}
	left := &ast.IdExpr{Name: id}

	_, ok = parser.expect(token.DOT)
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

func (parser *Parser) parseForLoop() (*ast.ForLoop, error) {
	_, ok := parser.expect(token.FOR)
	// TODO(errors): should never hit
	if !ok {
		return nil, fmt.Errorf("expected 'for'")
	}

	_, ok = parser.expect(token.OPEN_PAREN)
	// TODO(errors): should never hit
	if !ok {
		return nil, fmt.Errorf("expected '('")
	}

	init, err := parser.parseVar()
	// TODO(errors)
	if err != nil {
		return nil, err
	}

	_, ok = parser.expect(token.SEMICOLON)
	// TODO(errors)
	if !ok {
		return nil, fmt.Errorf("expected ';'")
	}

	cond, err := parser.parseExpr()
	// TODO(errors)
	if err != nil {
		return nil, err
	}

	_, ok = parser.expect(token.SEMICOLON)
	// TODO(errors)
	if !ok {
		return nil, fmt.Errorf("expected ';'")
	}

	update, err := parser.parseVar()
	// TODO(errors)
	if err != nil {
		return nil, err
	}

	_, ok = parser.expect(token.CLOSE_PAREN)
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
	collector := diagnostics.New()

	src := []byte(input)
	lex := lexer.New(filename, src, collector)
	parser := New(lex, collector)

	forLoop, err := parser.parseForLoop()
	return forLoop, err
}

func ParseWhileLoopFrom(input, filename string) (*ast.WhileLoop, error) {
	collector := diagnostics.New()

	src := []byte(input)
	lex := lexer.New(filename, src, collector)
	parser := New(lex, collector)

	whileLoop, err := parser.parseWhileLoop()
	return whileLoop, err
}

func (p *Parser) parseWhileLoop() (*ast.WhileLoop, error) {
	_, ok := p.expect(token.WHILE)
	// TODO(errors): should never hit
	if !ok {
		return nil, fmt.Errorf("expected 'while'")
	}

	expr, err := p.parseExpr()
	if err != nil {
		return nil, err
	}

	block, err := p.parseBlock()
	if err != nil {
		return nil, err
	}
	return &ast.WhileLoop{Cond: expr, Block: block}, nil
}
