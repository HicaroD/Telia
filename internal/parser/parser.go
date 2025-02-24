package parser

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/HicaroD/Telia/internal/ast"
	"github.com/HicaroD/Telia/internal/diagnostics"
	"github.com/HicaroD/Telia/internal/lexer"
	"github.com/HicaroD/Telia/internal/lexer/token"
)

type Parser struct {
	lex       *lexer.Lexer
	collector *diagnostics.Collector

	moduleScope *ast.Scope
}

func New(collector *diagnostics.Collector) *Parser {
	parser := new(Parser)
	parser.lex = nil
	parser.moduleScope = nil
	parser.collector = collector
	return parser
}

// Useful for testing
func NewWithLex(lex *lexer.Lexer, collector *diagnostics.Collector) *Parser {
	return &Parser{lex: lex, collector: collector}
}

func (p *Parser) ParseModuleDir(dirName, path string) (*ast.Program, error) {
	// Universe scope has a nil parent
	universe := ast.NewScope(nil)

	root := new(ast.Package)
	root.Name = dirName
	root.Scope = universe
	root.IsRoot = true

	err := p.buildModuleTree(path, root)
	if err != nil {
		return nil, err
	}

	return &ast.Program{Root: root}, nil
}

func (p *Parser) ParseFileAsProgram(lex *lexer.Lexer) (*ast.Program, error) {
	// Universe scope has a nil parent
	universe := ast.NewScope(nil)
	moduleScope := ast.NewScope(universe)

	file, err := p.parseFile(lex, moduleScope)
	if err != nil {
		return nil, err
	}

	module := &ast.Package{
		Name:   lex.ParentDirName,
		Files:  []*ast.File{file},
		Scope:  moduleScope,
		IsRoot: true,
	}
	program := &ast.Program{Root: module}

	return program, nil
}

func (p *Parser) parseFile(lex *lexer.Lexer, moduleScope *ast.Scope) (*ast.File, error) {
	file := &ast.File{
		Dir:            lex.ParentDirName,
		Path:           lex.Path,
		PkgNameDefined: false,
	}

	p.lex = lex
	p.moduleScope = moduleScope

	err := p.parseFileDecls(file)
	if err != nil {
		return nil, err
	}

	return file, nil
}

func (p *Parser) parseFileDecls(file *ast.File) error {
	var decls []*ast.Node
	firstNode := true

	for {
		node, eof, err := p.next()
		if err != nil {
			return err
		}

		if eof {
			break
		}

		if node == nil {
			continue
		}

		if node.Kind == ast.KIND_PKG_DECL {
			pkg := node.Node.(*ast.PkgDecl)
			if file.PkgNameDefined {
				// TODO(errors)
				return fmt.Errorf("redeclaration of package name\n")
			}

			file.PkgName = pkg.Name.Name()
			file.PkgNameDefined = true
		} else if firstNode {
			// TODO(errors)
			return fmt.Errorf("expected package declaration as first node")
		}

		firstNode = false
		decls = append(decls, node)
	}

	if !file.PkgNameDefined {
		// TODO(errors)
		return fmt.Errorf("Package name not defined\n")
	}

	file.Body = decls
	return nil
}

func (p *Parser) buildModuleTree(path string, module *ast.Package) error {
	return p.processModuleEntries(path, func(entry os.DirEntry, fullPath string) error {
		switch {
		case entry.IsDir():
			childModule := new(ast.Package)
			childScope := ast.NewScope(module.Scope)
			childModule.Name = filepath.Base(fullPath)
			childModule.Scope = childScope
			childModule.IsRoot = false

			module.Packages = append(module.Packages, childModule)

			return p.buildModuleTree(fullPath, childModule)
		case filepath.Ext(entry.Name()) == ".t":
			fileDirName := filepath.Base(filepath.Dir(fullPath))

			lex, err := lexer.NewFromFilePath(fileDirName, fullPath, p.collector)
			if err != nil {
				return err
			}

			file, err := p.parseFile(lex, module.Scope)
			if err != nil {
				return err
			}

			module.Files = append(module.Files, file)
		}
		return nil
	})
}

func (p *Parser) processModuleEntries(path string, handler func(entry os.DirEntry, fullPath string) error) error {
	entries, err := os.ReadDir(path)
	if err != nil {
		return fmt.Errorf("failed to read directory %q: %w", path, err)
	}

	for _, entry := range entries {
		fullPath := filepath.Join(path, entry.Name())
		if err := handler(entry, fullPath); err != nil {
			return err
		}
	}
	return nil
}

func (p *Parser) next() (*ast.Node, bool, error) {
	var err error
	var attributes *ast.Attributes
	attributesFound := false
	eof := false

peekAgain:
	tok := p.lex.Peek()
	if tok.Kind == token.EOF {
		eof = true
		return nil, eof, nil
	}
	if tok.Kind == token.SHARP {
		// TODO(errors)
		if attributesFound {
			return nil, eof, fmt.Errorf("attributes already defined\n")
		}
		attributes, err = p.parseAttributes()
		if err != nil {
			return nil, eof, err
		}
		attributesFound = true
		goto peekAgain
	}

	// TODO: try to parse the attribute before the actual declaration
	// TODO: single struct attribute to make it easier to parse

	switch tok.Kind {
	case token.PACKAGE:
		pkgDecl, err := p.parsePkgDecl()
		return pkgDecl, eof, err
	case token.USE:
		imp, err := p.parseUse()
		return imp, eof, err
	case token.TYPE:
		_, err := p.parseTypeAlias()
		return nil, eof, err
	case token.EXTERN:
		externDecl, err := p.parseExternDecl(attributes)
		return externDecl, eof, err
	case token.FN:
		fnDecl, err := p.parseFnDecl(attributes)
		return fnDecl, eof, err
	default:
		pos := tok.Pos
		unexpectedTokenOnGlobalScope := diagnostics.Diag{
			Message: fmt.Sprintf(
				"%s:%d:%d: unexpected non-declaration statement on global scope: %s\n",
				pos.Filename,
				pos.Line,
				pos.Column,
				tok.Kind.String(),
			),
		}

		p.collector.ReportAndSave(unexpectedTokenOnGlobalScope)
		return nil, eof, diagnostics.COMPILER_ERROR_FOUND
	}
}

// Useful for testing
func ParseExprFrom(expr, filename string) (*ast.Node, error) {
	collector := diagnostics.New()

	src := []byte(expr)
	lex := lexer.New(filename, src, collector)
	parser := NewWithLex(lex, collector)

	exprAst, err := parser.parseExpr(nil)
	if err != nil {
		return nil, err
	}
	return exprAst, nil
}

func (p *Parser) parseAttributes() (*ast.Attributes, error) {
	attributes := new(ast.Attributes)

	_, ok := p.expect(token.SHARP)
	if !ok {
		return nil, fmt.Errorf("expected '#'")
	}

	_, ok = p.expect(token.OPEN_BRACKET)
	if !ok {
		return nil, fmt.Errorf("expected '['")
	}

	for {
		attribute, ok := p.expect(token.ID)
		if !ok {
			return nil, fmt.Errorf("expected identifier")
		}

		_, ok = p.expect(token.EQUAL)
		if !ok {
			return nil, fmt.Errorf("expected '='")
		}

		attributeValue, ok := p.expect(token.UNTYPED_STRING)
		if !ok {
			return nil, fmt.Errorf("expected string literal")
		}

		switch attribute.Name() {
		case "default_cc":
			attributes.DefaultCallingConvention = attributeValue.Name()
		case "link_prefix":
			attributes.LinkPrefix = attributeValue.Name()
		case "link_name":
			attributes.LinkName = attributeValue.Name()
		case "linkage":
			attributes.Linkage = attributeValue.Name()
		default:
			return nil, fmt.Errorf("invalid attribute for extern declaration: %s\n", attribute.Name())
		}

		if p.lex.NextIs(token.CLOSE_BRACKET) {
			break
		}

		if !p.lex.NextIs(token.COMMA) {
			return nil, fmt.Errorf("expected either comma or closing bracket")
		}
	}

	_, ok = p.expect(token.CLOSE_BRACKET)
	if !ok {
		return nil, fmt.Errorf("expected ']'")
	}

	return attributes, nil
}

func (p *Parser) parseExternDecl(attributes *ast.Attributes) (*ast.Node, error) {
	externDecl := new(ast.ExternDecl)

	if attributes != nil {
		externDecl.Attributes = attributes
		attributes = nil
	}

	ext, ok := p.expect(token.EXTERN)
	if !ok {
		return nil, fmt.Errorf("expected 'extern', not %s\n", ext.Kind.String())
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
	externDecl.Name = name

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

	var err error
	var prototypes []*ast.Proto
	for {
		if p.lex.NextIs(token.CLOSE_CURLY) {
			break
		}

		var attributes *ast.Attributes
		if p.lex.NextIs(token.SHARP) {
			attributes, err = p.parseAttributes()
			if err != nil {
				return nil, err
			}
		}

		proto, err := p.parsePrototype(attributes)
		if err != nil {
			return nil, err
		}
		prototypes = append(prototypes, proto)
	}
	externDecl.Prototypes = prototypes

	closeCurly, ok := p.expect(token.CLOSE_CURLY)
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

	externScope := ast.NewScope(p.moduleScope)
	externDecl.Scope = externScope

	for i, prototype := range prototypes {
		n := new(ast.Node)
		n.Kind = ast.KIND_PROTO
		n.Node = prototype

		err := externScope.Insert(prototype.Name.Name(), n)
		if err != nil {
			if err == ast.ErrSymbolAlreadyDefinedOnScope {
				pos := prototypes[i].Name.Pos
				prototypeRedeclaration := diagnostics.Diag{
					Message: fmt.Sprintf(
						"%s:%d:%d: prototype '%s' already declared on extern '%s'",
						pos.Filename,
						pos.Line,
						pos.Column,
						prototype.Name.Name(),
						name.Name(),
					),
				}
				p.collector.ReportAndSave(prototypeRedeclaration)
				return nil, diagnostics.COMPILER_ERROR_FOUND
			}
			return nil, err
		}
	}

	n := new(ast.Node)
	n.Kind = ast.KIND_EXTERN_DECL
	n.Node = externDecl

	err = p.moduleScope.Insert(name.Name(), n)
	if err != nil {
		if err == ast.ErrSymbolAlreadyDefinedOnScope {
			pos := name.Pos
			prototypeRedeclaration := diagnostics.Diag{
				Message: fmt.Sprintf(
					"%s:%d:%d: extern '%s' already declared on scope",
					pos.Filename,
					pos.Line,
					pos.Column,
					name.Name(),
				),
			}
			p.collector.ReportAndSave(prototypeRedeclaration)
			return nil, diagnostics.COMPILER_ERROR_FOUND
		}
		return nil, err
	}

	return n, nil
}

func (p *Parser) parsePkgDecl() (*ast.Node, error) {
	pkg, ok := p.expect(token.PACKAGE)
	// TODO(errors)
	if !ok {
		return nil, fmt.Errorf("expected 'pkg' keyword, not %s\n", pkg.Kind.String())
	}

	name, ok := p.expect(token.ID)
	// TODO(errors)
	if !ok {
		return nil, fmt.Errorf("expected package name, not %s\n", name.Kind.String())
	}

	semi, ok := p.expect(token.SEMICOLON)
	// TODO(errors)
	if !ok {
		return nil, fmt.Errorf("expected semicolon, not %s\n", semi.Kind.String())
	}

	node := new(ast.Node)
	node.Kind = ast.KIND_PKG_DECL
	node.Node = &ast.PkgDecl{Name: name}
	return node, nil
}

func (p *Parser) parseUse() (*ast.Node, error) {
	imp := new(ast.UseDecl)

	use, ok := p.expect(token.USE)
	// TODO(errors)
	if !ok {
		return nil, fmt.Errorf("expected 'use' keyword, not %s\n", use.Kind.String())
	}

	useStr, ok := p.expect(token.UNTYPED_STRING)
	// TODO(errors)
	if !ok {
		return nil, fmt.Errorf("error: expected import string, not %s\n", useStr.Kind.String())
	}

	parts := strings.Split(useStr.Name(), "::")
	// TODO(errors)
	if len(parts) < 2 {
		return nil, fmt.Errorf("error: bad formatted import string")
	}

	switch parts[0] {
	case "std":
		imp.Std = true
	case "package":
		imp.Package = true
	default:
		// TODO(errors)
		return nil, fmt.Errorf("error: invalid use string prefix")
	}
	imp.Path = parts[1:]

	semi, ok := p.expect(token.SEMICOLON)
	// TODO(errors)
	if !ok {
		return nil, fmt.Errorf("expected semicolon, not %s\n", semi.Kind.String())
	}

	node := new(ast.Node)
	node.Kind = ast.KIND_USE_DECL
	node.Node = imp
	return node, nil
}

func (p *Parser) parseTypeAlias() (*ast.Node, error) {
	pkg, ok := p.expect(token.TYPE)
	// TODO(errors)
	if !ok {
		return nil, fmt.Errorf("expected 'type' keyword, not %s\n", pkg.Kind.String())
	}

	name, ok := p.expect(token.ID)
	// TODO(errors)
	if !ok {
		return nil, fmt.Errorf("expected type alias name, not %s\n", name.Kind.String())
	}

	_, ok = p.expect(token.EQUAL)
	// TODO(errors)
	if !ok {
		return nil, fmt.Errorf("expected semicolon, not %s\n", name.Kind.String())
	}

	ty, err := p.parseExprType()
	// TODO(errors)
	if err != nil {
		return nil, fmt.Errorf("expected valid alias type for '%s'\n", name.Name())
	}

	semi, ok := p.expect(token.SEMICOLON)
	// TODO(errors)
	if !ok {
		return nil, fmt.Errorf("expected semicolon, not %s\n", semi.Kind.String())
	}

	node := new(ast.Node)
	node.Kind = ast.KIND_TYPE_ALIAS_DECL

	alias := new(ast.TypeAlias)
	alias.Name = name
	alias.Type = ty

	node.Node = alias

	err = p.moduleScope.Insert(name.Name(), node)
	if err != nil {
		if err == ast.ErrSymbolAlreadyDefinedOnScope {
			return nil, fmt.Errorf("symbol '%s' already declared on scope\n", name.Name())
		}
		return nil, err
	}

	return node, nil
}

func (p *Parser) parsePrototype(attributes *ast.Attributes) (*ast.Proto, error) {
	prototype := new(ast.Proto)
	if attributes != nil {
		prototype.Attributes = attributes
		attributes = nil
	}

	fn, ok := p.expect(token.FN)
	if !ok {
		pos := fn.Pos
		expectedCloseCurly := diagnostics.Diag{
			Message: fmt.Sprintf(
				"%s:%d:%d: expected prototype declaration or }, not %s",
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
	prototype.Name = name

	params, err := p.parseFunctionParams(name, nil, true)
	if err != nil {
		return nil, err
	}
	prototype.Params = params

	returnType, err := p.parseReturnType( /*isPrototype=*/ true)
	if err != nil {
		return nil, err
	}
	prototype.RetType = returnType

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

	return prototype, nil
}

func (p *Parser) parseFnDecl(attributes *ast.Attributes) (*ast.Node, error) {
	var err error
	fnDecl := new(ast.FnDecl)

	if attributes != nil {
		fnDecl.Attributes = attributes
		attributes = nil
	}

	_, ok := p.expect(token.FN)
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
	fnDecl.Name = name

	fnScope := ast.NewScope(p.moduleScope)
	fnDecl.Scope = fnScope

	params, err := p.parseFunctionParams(name, fnScope, false)
	if err != nil {
		return nil, err
	}
	fnDecl.Params = params

	returnType, err := p.parseReturnType( /*isPrototype=*/ false)
	if err != nil {
		return nil, err
	}
	fnDecl.RetType = returnType

	block, err := p.parseBlock(fnScope)
	if err != nil {
		return nil, err
	}
	fnDecl.Block = block

	n := new(ast.Node)
	n.Kind = ast.KIND_FN_DECL
	n.Node = fnDecl

	err = p.moduleScope.Insert(name.Name(), n)
	if err != nil {
		if err == ast.ErrSymbolAlreadyDefinedOnScope {
			functionRedeclaration := diagnostics.Diag{
				Message: fmt.Sprintf(
					"%s:%d:%d: function '%s' already declared on scope",
					name.Pos.Filename,
					name.Pos.Line,
					name.Pos.Column,
					name.Name(),
				),
			}
			p.collector.ReportAndSave(functionRedeclaration)
		}
		return nil, err
	}

	return n, nil
}

// Useful for testing
func parseFnDeclFrom(filename, input string, scope *ast.Scope) (*ast.FnDecl, error) {
	collector := diagnostics.New()

	src := []byte(input)
	lexer := lexer.New(filename, src, collector)
	parser := NewWithLex(lexer, collector)
	parser.moduleScope = scope

	fnDecl, err := parser.parseFnDecl(nil)
	if err != nil {
		return nil, err
	}

	return fnDecl.Node.(*ast.FnDecl), nil
}

func (p *Parser) parseFunctionParams(functionName *token.Token, scope *ast.Scope, isPrototype bool) (*ast.FieldList, error) {
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

		param := new(ast.Field)

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
		param.Name = name

		ty, err := p.parseExprType()
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
		param.Type = ty

		// NOTE: prototypes parameters are validated at the semantic analyzer
		// stage
		if !isPrototype {
			if scope == nil {
				// TODO(errors): add proper error
				return nil, fmt.Errorf("error: scope should not be null when validating function parameters")
			}

			n := new(ast.Node)
			n.Kind = ast.KIND_FIELD
			n.Node = param

			err = scope.Insert(param.Name.Name(), n)
			if err != nil {
				if err == ast.ErrSymbolAlreadyDefinedOnScope {
					pos := param.Name.Pos
					parameterRedeclaration := diagnostics.Diag{
						Message: fmt.Sprintf(
							"%s:%d:%d: parameter '%s' already declared on function '%s'",
							pos.Filename,
							pos.Line,
							pos.Column,
							param.Name.Name(),
							functionName,
						),
					}
					p.collector.ReportAndSave(parameterRedeclaration)
					return nil, diagnostics.COMPILER_ERROR_FOUND
				}
				return nil, err
			}
		}

		params = append(params, param)
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

func (p *Parser) parseReturnType(isPrototype bool) (*ast.ExprType, error) {
	ty := new(ast.ExprType)

	if (isPrototype && p.lex.NextIs(token.SEMICOLON)) ||
		p.lex.NextIs(token.OPEN_CURLY) {
		ty.Kind = ast.EXPR_TYPE_BASIC
		ty.T = &ast.BasicType{Kind: token.VOID_TYPE}
		return ty, nil
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

func (p *Parser) parseExprType() (*ast.ExprType, error) {
	t := new(ast.ExprType)

	tok := p.lex.Peek()
	switch tok.Kind {
	case token.STAR:
		p.lex.Skip() // *
		ty, err := p.parseExprType()
		if err != nil {
			return nil, err
		}

		t.Kind = ast.EXPR_TYPE_POINTER
		t.T = &ast.PointerType{Type: ty}
	case token.ID:
		p.lex.Skip()
		symbol, err := p.moduleScope.LookupAcrossScopes(tok.Name())
		// TODO(errors)
		if err != nil {
			if err == ast.ErrSymbolNotFoundOnScope {
				return nil, fmt.Errorf("id type '%s' not found on scope\n", tok.Name())
			}
		}
		// TODO: add more id types, such as struct

		if symbol.Kind == ast.KIND_TYPE_ALIAS_DECL {
			alias := symbol.Node.(*ast.TypeAlias)
			return alias.Type, nil
		}

		// NOTE: I think ast.IdType won't be necessary anymore
		t.Kind = ast.EXPR_TYPE_ID
		t.T = &ast.IdType{Name: tok}
	default:
		if tok.Kind.IsBasicType() {
			p.lex.Skip()
			t.Kind = ast.EXPR_TYPE_BASIC
			t.T = &ast.BasicType{Kind: tok.Kind}
			return t, nil
		}
		return nil, diagnostics.COMPILER_ERROR_FOUND
	}
	return t, nil
}

func (p *Parser) parseStmt(parentScope *ast.Scope) (*ast.Node, error) {
	n := new(ast.Node)

	tok := p.lex.Peek()
	switch tok.Kind {
	case token.RETURN:
		p.lex.Skip()

		n.Kind = ast.KIND_RETURN_STMT

		returnStmt := new(ast.ReturnStmt)
		returnStmt.Return = tok
		returnStmt.Value = &ast.Node{Kind: ast.KIND_VOID_EXPR, Node: nil}

		n.Node = returnStmt

		if p.lex.NextIs(token.SEMICOLON) {
			p.lex.Skip()
			return n, nil
		}
		returnValue, err := p.parseExpr(parentScope)
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

		return n, nil
	case token.ID:
		idStmt, err := p.ParseIdStmt(parentScope)
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
		return idStmt, err
	case token.IF:
		condStmt, err := p.parseCondStmt(parentScope)
		return condStmt, err
	case token.FOR:
		forLoop, err := p.parseForLoop(parentScope)
		return forLoop, err
	case token.WHILE:
		whileLoop, err := p.parseWhileLoop(parentScope)
		return whileLoop, err
	default:
		return nil, nil
	}
}

func (p *Parser) parseBlock(parentScope *ast.Scope) (*ast.BlockStmt, error) {
	openCurly, ok := p.expect(token.OPEN_CURLY)
	if !ok {
		return nil, fmt.Errorf("expected '{', but got %s", openCurly)
	}

	var statements []*ast.Node

	for {
		tok := p.lex.Peek()
		if tok.Kind == token.CLOSE_CURLY {
			break
		}

		stmt, err := p.parseStmt(parentScope)
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

func (p *Parser) ParseIdStmt(parentScope *ast.Scope) (*ast.Node, error) {
	aheadId := p.lex.Peek1()
	switch aheadId.Kind {
	case token.OPEN_PAREN:
		fnCall, err := p.parseFnCall(parentScope)
		return fnCall, err
	case token.DOT:
		fieldAccessing, err := p.parseFieldAccess(parentScope)
		return fieldAccessing, err
	default:
		return p.parseVar(parentScope)
	}
}

func (p *Parser) parseVar(parentScope *ast.Scope) (*ast.Node, error) {
	variables := make([]*ast.Node, 0)
	isDecl := false

VarDecl:
	for {
		name, ok := p.expect(token.ID)
		// TODO(errors): add proper error
		if !ok {
			return nil, fmt.Errorf("expected ID")
		}

		n := new(ast.Node)
		n.Kind = ast.KIND_VAR_STMT
		variable := &ast.VarStmt{
			Name:           name,
			Type:           nil,
			Value:          nil,
			Decl:           true,
			NeedsInference: true,
			BackendType:    nil,
		}
		n.Node = variable
		variables = append(variables, n)

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

	exprs, err := p.parseExprList([]token.Kind{token.SEMICOLON, token.OPEN_CURLY}, parentScope)
	if err != nil {
		return nil, err
	}

	// TODO(errors): add proper error
	if len(variables) != len(exprs) {
		return nil, fmt.Errorf("%d != %d", len(variables), len(exprs))
	}
	for i, n := range variables {
		variable := n.Node.(*ast.VarStmt)
		variable.Value = exprs[i]
		variable.Decl = isDecl

		if variable.Decl {
			_, err := parentScope.LookupCurrentScope(variable.Name.Name())
			// TODO(errors)
			if err != nil {
				if err != ast.ErrSymbolNotFoundOnScope {
					return nil, fmt.Errorf("'%s' already declared in the current block", variable.Name.Name())
				}
			}
			err = parentScope.Insert(variable.Name.Name(), n)
			if err != nil {
				return nil, err
			}
		} else {
			_, err := parentScope.LookupAcrossScopes(variable.Name.Name())
			// TODO(errors)
			if err != nil {
				if err == ast.ErrSymbolNotFoundOnScope {
					name := variable.Name.Name()
					pos := variable.Name.Pos
					return nil, fmt.Errorf("%s '%s' not declared in the current block", pos, name)
				}
			}
		}
	}

	if len(variables) == 1 {
		return variables[0], nil
	}

	n := new(ast.Node)
	n.Kind = ast.KIND_MULTI_VAR_STMT
	n.Node = &ast.MultiVarStmt{IsDecl: isDecl, Variables: variables}
	return n, nil
}

// Useful for testing
func parseVarFrom(filename, input string) (*ast.VarStmt, error) {
	collector := diagnostics.New()

	src := []byte(input)
	lex := lexer.New(filename, src, collector)
	parser := NewWithLex(lex, collector)

	tmpScope := ast.NewScope(nil)
	stmt, err := parser.ParseIdStmt(tmpScope)
	if err != nil {
		return nil, err
	}
	return stmt.Node.(*ast.VarStmt), nil
}

func (p *Parser) parseCondStmt(parentScope *ast.Scope) (*ast.Node, error) {
	ifCond, err := p.parseIfCond(parentScope)
	if err != nil {
		return nil, err
	}

	elifConds, err := p.parseElifConds(parentScope)
	if err != nil {
		return nil, err
	}

	elseCond, err := p.parseElseCond(parentScope)
	if err != nil {
		return nil, err
	}

	n := new(ast.Node)
	n.Kind = ast.KIND_COND_STMT
	n.Node = &ast.CondStmt{IfStmt: ifCond, ElifStmts: elifConds, ElseStmt: elseCond}
	return n, nil
}

func (p *Parser) parseIfCond(parentScope *ast.Scope) (*ast.IfElifCond, error) {
	ifToken, ok := p.expect(token.IF)
	// TODO(errors): add proper error
	if !ok {
		return nil, fmt.Errorf("expected 'if'")
	}

	ifExpr, err := p.parseExpr(parentScope)
	if err != nil {
		return nil, err
	}

	ifScope := ast.NewScope(parentScope)
	ifBlock, err := p.parseBlock(ifScope)
	if err != nil {
		return nil, err
	}
	return &ast.IfElifCond{If: &ifToken.Pos, Expr: ifExpr, Block: ifBlock, Scope: ifScope}, nil
}

func (p *Parser) parseElifConds(parentScope *ast.Scope) ([]*ast.IfElifCond, error) {
	var elifConds []*ast.IfElifCond
	for {
		elifToken, ok := p.expect(token.ELIF)
		if !ok {
			break
		}
		elifExpr, err := p.parseExpr(parentScope)
		if err != nil {
			return nil, err
		}
		elifScope := ast.NewScope(parentScope)
		elifBlock, err := p.parseBlock(elifScope)
		if err != nil {
			return nil, err
		}
		elifConds = append(
			elifConds,
			&ast.IfElifCond{If: &elifToken.Pos, Expr: elifExpr, Block: elifBlock, Scope: elifScope},
		)
	}
	return elifConds, nil
}

func (p *Parser) parseElseCond(parentScope *ast.Scope) (*ast.ElseCond, error) {
	elseToken, ok := p.expect(token.ELSE)
	if !ok {
		return nil, nil
	}

	elseScope := ast.NewScope(parentScope)
	elseBlock, err := p.parseBlock(elseScope)
	if err != nil {
		return nil, err
	}
	return &ast.ElseCond{Else: &elseToken.Pos, Block: elseBlock, Scope: elseScope}, nil
}

func (p *Parser) parseExpr(parentScope *ast.Scope) (*ast.Node, error) {
	expr, err := p.parseLogical(parentScope)
	if err != nil {
		return nil, err
	}

	if p.lex.NextIs(token.AT) {
		atOperator, err := p.parseAtOperator(parentScope)
		if err != nil {
			return nil, err
		}
		fmt.Print(atOperator)
	}

	return expr, nil
}

func (p *Parser) parseAtOperator(parentScope *ast.Scope) (*ast.AtOperator, error) {
	_, ok := p.expect(token.AT)
	// TODO(errors)
	if !ok {
		return nil, fmt.Errorf("expected '@'\n")
	}

	name, ok := p.expect(token.ID)
	// TODO(errors)
	if !ok {
		return nil, fmt.Errorf("expected operator name, not %s\n", name.Name())
	}

	op := new(ast.AtOperator)

	switch name.Name() {
	case "prop":
		op.Kind = ast.AT_OPERATOR_PROP
		op.Op = nil
	case "fail":
		op.Kind = ast.AT_OPERATOR_FAIL
		op.Op = nil
	case "catch":
		op.Kind = ast.AT_OPERATOR_CATCH
		catchOp, err := p.parseCatchOperator(parentScope)
		if err != nil {
			return nil, err
		}
		op.Op = catchOp
	default:
		return nil, fmt.Errorf("invalid at operator '%s'\n", name.Name())
	}

	return op, nil
}

func (p *Parser) parseCatchOperator(parentScope *ast.Scope) (*ast.CatchAtOperator, error) {
	catchOp := new(ast.CatchAtOperator)

	scope := ast.NewScope(parentScope)
	catchOp.Scope = scope

	varName, ok := p.expect(token.ID)
	if !ok {
		return nil, fmt.Errorf("expected an identifier, not %s\n", varName.Name())
	}
	catchOp.ErrVarName = varName

	block, err := p.parseBlock(parentScope)
	if err != nil {
		return nil, err
	}
	catchOp.Block = block

	return catchOp, nil
}

func (p *Parser) parseLogical(parentScope *ast.Scope) (*ast.Node, error) {
	lhs, err := p.parseComparasion(parentScope)
	if err != nil {
		return nil, err
	}

	for {
		next := p.lex.Peek()
		if _, ok := ast.LOGICAL[next.Kind]; ok {
			p.lex.Skip()
			rhs, err := p.parseComparasion(parentScope)
			// TODO(errors): add proper error
			if err != nil {
				return nil, err
			}

			l := new(ast.Node)
			l.Kind = ast.KIND_BINARY_EXPR
			l.Node = &ast.BinaryExpr{Left: lhs, Op: next.Kind, Right: rhs}
			lhs = l
		} else {
			break
		}
	}

	return lhs, nil
}

func (p *Parser) parseComparasion(parentScope *ast.Scope) (*ast.Node, error) {
	lhs, err := p.parseTerm(parentScope)
	if err != nil {
		return nil, err
	}

	for {
		next := p.lex.Peek()
		if _, ok := ast.COMPARASION[next.Kind]; ok {
			p.lex.Skip()
			rhs, err := p.parseTerm(parentScope)
			if err != nil {
				return nil, err
			}

			l := new(ast.Node)
			l.Kind = ast.KIND_BINARY_EXPR
			l.Node = &ast.BinaryExpr{Left: lhs, Op: next.Kind, Right: rhs}
			lhs = l
		} else {
			break
		}
	}
	return lhs, nil
}

func (p *Parser) parseTerm(parentScope *ast.Scope) (*ast.Node, error) {
	lhs, err := p.parseFactor(parentScope)
	if err != nil {
		return nil, err
	}

	for {
		next := p.lex.Peek()
		if _, ok := ast.TERM[next.Kind]; ok {
			p.lex.Skip()
			rhs, err := p.parseFactor(parentScope)
			if err != nil {
				return nil, err
			}
			l := new(ast.Node)
			l.Kind = ast.KIND_BINARY_EXPR
			l.Node = &ast.BinaryExpr{Left: lhs, Op: next.Kind, Right: rhs}
			lhs = l
		} else {
			break
		}
	}
	return lhs, nil
}

func (p *Parser) parseFactor(parentScope *ast.Scope) (*ast.Node, error) {
	lhs, err := p.parseUnary(parentScope)
	if err != nil {
		return nil, err
	}

	for {
		next := p.lex.Peek()
		if _, ok := ast.FACTOR[next.Kind]; ok {
			p.lex.Skip()
			rhs, err := p.parseUnary(parentScope)
			if err != nil {
				return nil, err
			}
			l := new(ast.Node)
			l.Kind = ast.KIND_BINARY_EXPR
			l.Node = &ast.BinaryExpr{Left: lhs, Op: next.Kind, Right: rhs}
			lhs = l
		} else {
			break
		}
	}
	return lhs, nil

}

func (p *Parser) parseUnary(parentScope *ast.Scope) (*ast.Node, error) {
	next := p.lex.Peek()
	if _, ok := ast.UNARY[next.Kind]; ok {
		p.lex.Skip()
		rhs, err := p.parseUnary(parentScope)
		if err != nil {
			return nil, err
		}

		unary := new(ast.Node)
		unary.Kind = ast.KIND_UNARY_EXPR
		unary.Node = &ast.UnaryExpr{Op: next.Kind, Value: rhs}
		return unary, nil
	}

	return p.parsePrimary(parentScope)
}

func (p *Parser) parsePrimary(parentScope *ast.Scope) (*ast.Node, error) {
	n := new(ast.Node)

	tok := p.lex.Peek()
	switch tok.Kind {
	case token.ID:
		idExpr := &ast.IdExpr{Name: tok}

		next := p.lex.Peek1()
		switch next.Kind {
		case token.OPEN_PAREN:
			return p.parseFnCall(parentScope)
		case token.DOT:
			fieldAccess, err := p.parseFieldAccess(parentScope)
			return fieldAccess, err
		}

		p.lex.Skip()

		n.Kind = ast.KIND_ID_EXPR
		n.Node = idExpr
		return n, nil
	case token.OPEN_PAREN:
		p.lex.Skip() // (
		expr, err := p.parseExpr(parentScope)
		if err != nil {
			return nil, err
		}
		_, ok := p.expect(token.CLOSE_PAREN)
		if !ok {
			return nil, fmt.Errorf("expected closing parenthesis")
		}
		return expr, nil
	default:
		if tok.Kind.IsLiteral() {
			p.lex.Skip()

			n.Kind = ast.KIND_LITERAl_EXPR
			n.Node = &ast.LiteralExpr{
				Type:  ast.NewBasicType(tok.Kind),
				Value: tok.Lexeme,
			}

			return n, nil
		}
		return nil, fmt.Errorf(
			"invalid token for expression parsing: %s %s %s",
			tok.Kind,
			tok.Lexeme,
			tok.Pos,
		)
	}
}

func (p *Parser) parseExprList(possibleEnds []token.Kind, parentScope *ast.Scope) ([]*ast.Node, error) {
	var exprs []*ast.Node
Var:
	for {
		for _, end := range possibleEnds {
			if p.lex.NextIs(end) {
				break Var
			}
		}

		expr, err := p.parseExpr(parentScope)
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

func (parser *Parser) parseFnCall(parentScope *ast.Scope) (*ast.Node, error) {
	name, ok := parser.expect(token.ID)
	if !ok {
		return nil, fmt.Errorf("expected 'id'")
	}

	_, ok = parser.expect(token.OPEN_PAREN)
	// TODO(errors): add proper error
	if !ok {
		return nil, fmt.Errorf("expected '('")
	}

	args, err := parser.parseExprList([]token.Kind{token.CLOSE_PAREN}, parentScope)
	if err != nil {
		return nil, err
	}

	_, ok = parser.expect(token.CLOSE_PAREN)
	// TODO(errors): add proper error
	if !ok {
		return nil, fmt.Errorf("expected ')'")
	}

	n := new(ast.Node)
	n.Kind = ast.KIND_FN_CALL
	n.Node = &ast.FnCall{Name: name, Args: args}
	return n, nil
}

func (parser *Parser) parseFieldAccess(parentScope *ast.Scope) (*ast.Node, error) {
	id, ok := parser.expect(token.ID)
	// TODO(errors): add proper error
	if !ok {
		return nil, fmt.Errorf("error: expected ID")
	}

	left := new(ast.Node)
	left.Kind = ast.KIND_ID_EXPR
	left.Node = &ast.IdExpr{Name: id}

	_, ok = parser.expect(token.DOT)
	// TODO(errors): add proper error
	if !ok {
		return nil, fmt.Errorf("error: expected a dot")
	}

	right, err := parser.parsePrimary(parentScope)
	if err != nil {
		log.Fatal(err)
	}

	n := new(ast.Node)
	n.Kind = ast.KIND_FIELD_ACCESS
	n.Node = &ast.FieldAccess{Left: left, Right: right}
	return n, nil
}

func (parser *Parser) parseForLoop(parentScope *ast.Scope) (*ast.Node, error) {
	_, ok := parser.expect(token.FOR)
	// TODO(errors): add proper error
	if !ok {
		return nil, fmt.Errorf("expected 'for'")
	}

	init, err := parser.parseVar(parentScope)
	if err != nil {
		return nil, err
	}

	_, ok = parser.expect(token.SEMICOLON)
	// TODO(errors): add proper error
	if !ok {
		return nil, fmt.Errorf("expected ';'")
	}

	cond, err := parser.parseExpr(parentScope)
	if err != nil {
		return nil, err
	}

	_, ok = parser.expect(token.SEMICOLON)
	// TODO(errors): add proper error
	if !ok {
		return nil, fmt.Errorf("expected ';'")
	}

	update, err := parser.parseVar(parentScope)
	if err != nil {
		return nil, err
	}

	forScope := ast.NewScope(parentScope)
	block, err := parser.parseBlock(forScope)
	if err != nil {
		return nil, err
	}

	n := new(ast.Node)
	n.Kind = ast.KIND_FOR_LOOP_STMT
	n.Node = &ast.ForLoop{Init: init, Cond: cond, Update: update, Block: block, Scope: forScope}
	return n, nil
}

// Useful for testing
func ParseForLoopFrom(input, filename string) (*ast.ForLoop, error) {
	collector := diagnostics.New()

	src := []byte(input)
	lex := lexer.New(filename, src, collector)
	parser := NewWithLex(lex, collector)

	tempScope := ast.NewScope(nil)
	forLoop, err := parser.parseForLoop(tempScope)
	return forLoop.Node.(*ast.ForLoop), err
}

func ParseWhileLoopFrom(input, filename string) (*ast.WhileLoop, error) {
	collector := diagnostics.New()

	src := []byte(input)
	lex := lexer.New(filename, src, collector)
	parser := NewWithLex(lex, collector)

	// TODO: set scope properly
	tempScope := ast.NewScope(nil)
	whileLoop, err := parser.parseWhileLoop(tempScope)
	return whileLoop.Node.(*ast.WhileLoop), err
}

func (p *Parser) parseWhileLoop(parentScope *ast.Scope) (*ast.Node, error) {
	_, ok := p.expect(token.WHILE)
	if !ok {
		return nil, fmt.Errorf("expected 'while'")
	}

	expr, err := p.parseExpr(parentScope)
	if err != nil {
		return nil, err
	}

	whileScope := ast.NewScope(parentScope)
	block, err := p.parseBlock(whileScope)
	if err != nil {
		return nil, err
	}

	n := new(ast.Node)
	n.Kind = ast.KIND_WHILE_LOOP_STMT
	n.Node = &ast.WhileLoop{Cond: expr, Block: block, Scope: whileScope}
	return n, nil
}
