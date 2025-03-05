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

	argLoc     string
	imports    map[string]*ast.Package
	processing map[string]bool
	pkg        *ast.Package
	file       *ast.File
}

func New(collector *diagnostics.Collector) *Parser {
	parser := new(Parser)
	parser.lex = nil
	parser.argLoc = ""
	parser.pkg = nil
	parser.file = nil
	parser.imports = make(map[string]*ast.Package)
	parser.collector = collector
	parser.processing = make(map[string]bool)
	return parser
}

// Useful for testing
func NewWithLex(lex *lexer.Lexer, collector *diagnostics.Collector) *Parser {
	return &Parser{lex: lex, collector: collector}
}

func (p *Parser) ParsePackageAsProgram(argLoc string, loc *ast.Loc) (*ast.Program, error) {
	p.argLoc = argLoc
	root, err := p.parsePackage(loc, true)
	if err != nil {
		return nil, err
	}
	return &ast.Program{Root: root}, nil
}

func (p *Parser) parsePackage(loc *ast.Loc, isRoot bool) (*ast.Package, error) {
	scope := ast.NewScope(nil)

	pkg := new(ast.Package)
	pkg.Loc = loc
	pkg.Scope = scope
	pkg.IsRoot = isRoot
	if isRoot {
		p.imports = make(map[string]*ast.Package)
		p.imports[p.argLoc] = pkg
	}

	err := p.buildPackage(pkg)
	if err != nil {
		return nil, err
	}

	return pkg, nil
}

func (p *Parser) addPackage(std bool, path []string) (string, string, *ast.Package, error) {
	pkgPath := strings.Join(path, "/")

	var prefixPath string
	if std {
		// TODO: get path to std directory from env
		prefixPath = "./std"
	} else {
		prefixPath = p.argLoc
	}

	impName := path[len(path)-1]
	fullPkgPath := filepath.Join(prefixPath, pkgPath)
	if pkg, found := p.imports[fullPkgPath]; found {
		return impName, fullPkgPath, pkg, nil
	}

	// TODO(errors)
	if _, processing := p.processing[fullPkgPath]; processing {
		return "", "", nil, fmt.Errorf("circular import detected: %s", fullPkgPath)
	}
	p.processing[fullPkgPath] = true
	defer delete(p.processing, fullPkgPath)

	loc, err := ast.LocFromPath(fullPkgPath)
	if err != nil {
		return "", "", nil, err
	}

	currentLex := p.lex
	currentPkg := p.pkg
	currentFile := p.file
	defer func() {
		p.lex = currentLex
		p.pkg = currentPkg
		p.file = currentFile
	}()

	pkg, err := p.parsePackage(loc, false)
	if err != nil {
		return "", "", nil, err
	}

	p.imports[fullPkgPath] = pkg
	return impName, fullPkgPath, pkg, nil
}

func (p *Parser) ParseFileAsProgram(argLoc string, loc *ast.Loc, collector *diagnostics.Collector) (*ast.Program, error) {
	p.argLoc = argLoc

	l, err := lexer.NewFromFilePath(loc, collector)
	if err != nil {
		return nil, err
	}

	// Universe scope has a nil parent
	universe := ast.NewScope(nil)
	// TODO: add builtins to universe scope
	packageScope := ast.NewScope(universe)

	pkg := &ast.Package{
		Loc:    loc,
		Scope:  packageScope,
		IsRoot: true,
	}

	file, err := p.parseFile(l, pkg)
	if err != nil {
		return nil, err
	}
	pkg.Files = []*ast.File{file}

	p.pkg = pkg
	program := &ast.Program{Root: pkg}
	return program, nil
}

func (p *Parser) parseFile(lex *lexer.Lexer, pkg *ast.Package) (*ast.File, error) {
	file := &ast.File{
		Loc:              lex.Loc,
		PkgNameDefined:   false,
		Imports:          make(map[string]*ast.UseDecl),
		IsFirstNode:      true,
		AnyDeclNodeFound: false,
	}

	p.lex = lex
	p.pkg = pkg
	p.file = file

	err := p.parseFileDecls(file)
	if err != nil {
		return nil, err
	}

	return file, nil
}

func (p *Parser) parseFileDecls(file *ast.File) error {
	var decls []*ast.Node

	for {
		node, eof, err := p.next(file)
		if err != nil {
			return err
		}

		if eof {
			break
		}

		if node == nil {
			continue
		}

		decls = append(decls, node)
	}

	file.Body = decls
	return nil
}

func (p *Parser) buildPackage(pkg *ast.Package) error {
	return p.processPackageFiles(pkg.Loc.Path, func(entry os.DirEntry, fullPath string) error {
		if filepath.Ext(entry.Name()) == ".t" {
			loc, err := ast.LocFromPath(fullPath)
			if err != nil {
				return err
			}
			lex, err := lexer.NewFromFilePath(loc, p.collector)
			if err != nil {
				return err
			}
			file, err := p.parseFile(lex, pkg)
			if err != nil {
				return err
			}
			pkg.Files = append(pkg.Files, file)
		}
		return nil
	})
}

func (p *Parser) processPackageFiles(path string, handler func(entry os.DirEntry, fullPath string) error) error {
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

func (p *Parser) next(file *ast.File) (*ast.Node, bool, error) {
	var err error
	var attributes *ast.Attributes
	var attributesFound bool
	var eof bool

peekAgain:
	p.skipNewLines()

	tok := p.lex.Peek()
	switch tok.Kind {
	case token.EOF:
		eof = true
		return nil, eof, nil
	case token.SHARP:
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

	switch tok.Kind {
	case token.PACKAGE:
		if !file.IsFirstNode {
			return nil, eof, fmt.Errorf("expected package declaration as first node")
		}
		pkgName, pkgDecl, err := p.parsePkgDecl()
		if err != nil {
			return nil, false, err
		}
		if file.PkgNameDefined {
			// TODO(errors)
			return nil, eof, fmt.Errorf("redeclaration of package name\n")
		}
		file.PkgName = pkgName
		file.PkgNameDefined = true
		return pkgDecl, eof, err
	case token.USE:
		if file.AnyDeclNodeFound {
			return nil, eof, fmt.Errorf("use statements should be placed at the top of file")
		}
		imp, err := p.parseUse()
		return imp, eof, err
	case token.TYPE:
		_, err := p.parseTypeAlias()
		file.AnyDeclNodeFound = true
		return nil, eof, err
	case token.EXTERN:
		externDecl, err := p.parseExternDecl(attributes)
		file.AnyDeclNodeFound = true
		return externDecl, eof, err
	case token.FN:
		fnDecl, err := p.parseFnDecl(attributes)
		file.AnyDeclNodeFound = true
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
	loc := new(ast.Loc)
	loc.Name = filename
	lex := lexer.New(loc, src, collector)
	parser := NewWithLex(lex, collector)

	exprAst, err := parser.parseSingleExpr(nil)
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

	if !p.lex.NextIs(token.ID) {
		if externDecl.Attributes == nil {
			externDecl.Attributes = new(ast.Attributes)
		}
		externDecl.Attributes.Global = true
	} else {
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

	newline, ok := p.expect(token.NEWLINE)
	if !ok {
		pos := openCurly.Pos
		expectedOpenCurly := diagnostics.Diag{
			Message: fmt.Sprintf(
				"%s:%d:%d: expected new line after {, not %s",
				pos.Filename,
				pos.Line,
				pos.Column,
				newline.Kind,
			),
		}
		p.collector.ReportAndSave(expectedOpenCurly)
		return nil, diagnostics.COMPILER_ERROR_FOUND
	}

	var err error
	var prototypes []*ast.Proto
	for {
		p.skipNewLines()

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

	if externDecl.Attributes.Global {
		externDecl.Scope = p.pkg.Scope
	} else {
		externScope := ast.NewScope(p.pkg.Scope)
		externDecl.Scope = externScope
	}

	for i, prototype := range prototypes {
		n := new(ast.Node)
		n.Kind = ast.KIND_PROTO
		n.Node = prototype

		err := externDecl.Scope.Insert(prototype.Name.Name(), n)
		if err != nil {
			if err == ast.ErrSymbolAlreadyDefinedOnScope {
				if externDecl.Attributes.Global {
					pos := prototypes[i].Name.Pos
					prototypeRedeclaration := diagnostics.Diag{
						Message: fmt.Sprintf(
							"%s:%d:%d: prototype '%s' already declared in the package scope",
							pos.Filename,
							pos.Line,
							pos.Column,
							prototype.Name.Name(),
						),
					}
					p.collector.ReportAndSave(prototypeRedeclaration)
					return nil, diagnostics.COMPILER_ERROR_FOUND
				} else {
					name := externDecl.Name
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
			}
			return nil, err
		}
	}

	n := new(ast.Node)
	n.Kind = ast.KIND_EXTERN_DECL
	n.Node = externDecl

	if !externDecl.Attributes.Global {
		name := externDecl.Name
		err = p.pkg.Scope.Insert(name.Name(), n)
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
	}

	return n, nil
}

func (p *Parser) skipNewLines() {
	for {
		newline := p.skipSingleNewLine()
		if !newline {
			break
		}
	}
}

func (p *Parser) parsePkgDecl() (string, *ast.Node, error) {
	pkg, ok := p.expect(token.PACKAGE)
	// TODO(errors)
	if !ok {
		return "", nil, fmt.Errorf("expected 'package' keyword, not %s\n", pkg.Kind.String())
	}

	name, ok := p.expect(token.ID)
	// TODO(errors)
	if !ok {
		return "", nil, fmt.Errorf("expected package name, not %s\n", name.Kind.String())
	}

	newline, ok := p.expect(token.NEWLINE)
	// TODO(errors)
	if !ok {
		return "", nil, fmt.Errorf("expected new line, not %s\n", newline.Kind.String())
	}

	node := new(ast.Node)
	node.Kind = ast.KIND_PKG_DECL
	node.Node = &ast.PkgDecl{Name: name}
	return name.Name(), node, nil
}

func (p *Parser) parseUse() (*ast.Node, error) {
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

	var std bool
	switch parts[0] {
	case "std":
		std = true
	case "pkg":
		std = false
	default:
		// TODO(errors)
		return nil, fmt.Errorf("error: invalid use string prefix")
	}

	newline, ok := p.expect(token.NEWLINE)
	// TODO(errors)
	if !ok {
		return nil, fmt.Errorf("expected new line, not %s\n", newline.Kind.String())
	}

	impName, fullPath, pkg, err := p.addPackage(std, parts[1:])
	if err != nil {
		return nil, err
	}

	// NOTE: by default, import name is the last name of the import
	// TODO: add custom import name

	_, found := p.file.Imports[fullPath]
	// TODO(errors)
	if found {
		return nil, fmt.Errorf("name conflict for import: %s\n", impName)
	}
	p.file.Imports[impName] = &ast.UseDecl{Name: impName, Package: pkg}
	return nil, nil
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
		return nil, fmt.Errorf("expected equal sign, not %s\n", name.Kind.String())
	}

	ty, err := p.parseExprType()
	// TODO(errors)
	if err != nil {
		return nil, fmt.Errorf("expected valid alias type for '%s'\n", name.Name())
	}

	newline, ok := p.expect(token.NEWLINE)
	// TODO(errors)
	if !ok {
		return nil, fmt.Errorf("expected new line, not %s\n", newline.Kind.String())
	}

	node := new(ast.Node)
	node.Kind = ast.KIND_TYPE_ALIAS_DECL

	alias := new(ast.TypeAlias)
	alias.Name = name
	alias.Type = ty

	node.Node = alias

	err = p.pkg.Scope.Insert(name.Name(), node)
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

	params, err := p.parseFnParams(name, nil, true)
	if err != nil {
		return nil, err
	}
	prototype.Params = params

	returnType, err := p.parseReturnType( /*isPrototype=*/ true)
	if err != nil {
		return nil, err
	}
	prototype.RetType = returnType

	newline, ok := p.expect(token.NEWLINE)
	if !ok {
		pos := newline.Pos
		expectedNewLine := diagnostics.Diag{
			Message: fmt.Sprintf(
				"%s:%d:%d: expected new line at the end of prototype, not %s",
				pos.Filename,
				pos.Line,
				pos.Column,
				newline.Kind,
			),
		}
		p.collector.ReportAndSave(expectedNewLine)
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

	fnScope := ast.NewScope(p.pkg.Scope)
	fnDecl.Scope = fnScope

	params, err := p.parseFnParams(name, fnScope, false)
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

	err = p.pkg.Scope.Insert(name.Name(), n)
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
	loc := new(ast.Loc)
	loc.Name = filename
	lex := lexer.New(loc, src, collector)
	parser := NewWithLex(lex, collector)
	parser.pkg.Scope = scope

	fnDecl, err := parser.parseFnDecl(nil)
	if err != nil {
		return nil, err
	}

	return fnDecl.Node.(*ast.FnDecl), nil
}

func (p *Parser) parseFnParams(functionName *token.Token, scope *ast.Scope, isPrototype bool) (*ast.Params, error) {
	var params []*ast.Param
	var length int
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

		param := new(ast.Param)

		attributes := new(ast.ParamAttributes)

		if p.lex.NextIs(token.AT) {
			for p.lex.NextIs(token.AT) {
				p.lex.Skip() // @

				attributeName, ok := p.expect(token.ID)
				if !ok {
					return nil, fmt.Errorf("expected parameter attribute name, not %s\n")
				}

				switch attributeName.Name() {
				case "for_c":
					// TODO(errors)
					if !isPrototype {
						return nil, fmt.Errorf("@for_c only allowed on prototypes\n")
					}
					// TODO(errors)
					if attributes.ForC {
						return nil, fmt.Errorf("cannot redeclare @for_c attribute\n")
					}
					attributes.ForC = true
				case "const":
					// TODO(errors)
					if attributes.Const {
						return nil, fmt.Errorf("cannot redeclare @const attribute\n")
					}
					attributes.Const = true
				}
			}
		}
		param.Attributes = attributes

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

		if p.lex.NextIs(token.DOT_DOT_DOT) {
			isVariadic = true
			param.Variadic = true
			p.lex.Skip() // ...
		}

		if !param.Variadic && param.Attributes.ForC {
			return nil, fmt.Errorf("@for_c attribute only allowed on variadic arguments")
		}

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

		if isVariadic {
			tok := p.lex.Peek()
			pos := tok.Pos

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
		}

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
		length++
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

	// disconsider the variadic argument
	if isVariadic {
		length--
	}

	return &ast.Params{
		Open:       openParen,
		Fields:     params,
		Len:        length,
		Close:      closeParen,
		IsVariadic: isVariadic,
	}, nil
}

func (p *Parser) parseReturnType(isPrototype bool) (*ast.ExprType, error) {
	ty := new(ast.ExprType)

	if (isPrototype && p.lex.NextIs(token.NEWLINE)) ||
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
		symbol, err := p.pkg.Scope.LookupAcrossScopes(tok.Name())
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
	case token.OPEN_PAREN:
		tupleType, err := p.parseTupleExpr()
		if err != nil {
			return nil, err
		}
		t.Kind = ast.EXPR_TYPE_TUPLE
		t.T = tupleType
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

func (p *Parser) parseTupleExpr() (*ast.TupleType, error) {
	_, ok := p.expect(token.OPEN_PAREN)
	if !ok {
		return nil, fmt.Errorf("expected (")
	}

	types := make([]*ast.ExprType, 0)
	for {
		if p.lex.NextIs(token.CLOSE_PAREN) {
			break
		}

		ty, err := p.parseExprType()
		if err != nil {
			return nil, err
		}
		types = append(types, ty)

		if p.lex.NextIs(token.COMMA) {
			p.lex.Skip() // ,
			continue
		}
	}

	_, ok = p.expect(token.CLOSE_PAREN)
	if !ok {
		return nil, fmt.Errorf("expected )")
	}

	return &ast.TupleType{Types: types}, nil
}

func (p *Parser) parseStmt(block *ast.BlockStmt, parentScope *ast.Scope, allowDefer bool) (*ast.Node, error) {
	n := new(ast.Node)
	endsWithNewLine := false

	tok := p.lex.Peek()
	switch tok.Kind {
	case token.RETURN:
		endsWithNewLine = true

		p.lex.Skip()

		n.Kind = ast.KIND_RETURN_STMT
		returnStmt := new(ast.ReturnStmt)
		returnStmt.Return = tok
		returnStmt.Value = &ast.Node{Kind: ast.KIND_VOID_EXPR, Node: nil}
		n.Node = returnStmt

		if p.lex.NextIs(token.NEWLINE) {
			break
		}

		returnValue, err := p.parseAnyExpr([]token.Kind{token.NEWLINE}, parentScope)
		if err != nil {
			tok := p.lex.Peek()
			pos := tok.Pos
			expectedNewline := diagnostics.Diag{
				Message: fmt.Sprintf(
					"%s:%d:%d: expected expression or new line, not %s",
					pos.Filename,
					pos.Line,
					pos.Column,
					tok.Kind,
				),
			}
			p.collector.ReportAndSave(expectedNewline)
			return nil, diagnostics.COMPILER_ERROR_FOUND
		}
		returnStmt.Value = returnValue
	case token.ID:
		endsWithNewLine = true
		idStmt, err := p.ParseIdStmt(parentScope)
		if err != nil {
			return nil, err
		}
		n = idStmt
	case token.IF:
		condStmt, err := p.parseCondStmt(parentScope)
		if err != nil {
			return nil, err
		}
		n = condStmt
	case token.FOR:
		forLoop, err := p.parseForLoop(parentScope)
		if err != nil {
			return nil, err
		}
		n = forLoop
	case token.WHILE:
		whileLoop, err := p.parseWhileLoop(parentScope)
		if err != nil {
			return nil, err
		}
		n = whileLoop
	case token.DEFER:
		// TODO(errors)
		if !allowDefer {
			return nil, fmt.Errorf("invalid nested defer statement")
		}
		deferStmt, err := p.parseDefer(block, parentScope)
		if err != nil {
			return nil, err
		}
		n = deferStmt
	default:
		log.Fatalf("unimplemented: %s\n", tok.Kind)
		return nil, nil
	}

	if endsWithNewLine {
		_, ok := p.expect(token.NEWLINE)
		if !ok {
			tok := p.lex.Peek()
			pos := tok.Pos
			expectedNewline := diagnostics.Diag{
				Message: fmt.Sprintf(
					"%s:%d:%d: expected new line at the end of statement, not %s",
					pos.Filename,
					pos.Line,
					pos.Column,
					tok.Kind,
				),
			}
			p.collector.ReportAndSave(expectedNewline)
			return nil, diagnostics.COMPILER_ERROR_FOUND
		}
	}

	return n, nil
}

func (p *Parser) parseBlock(parentScope *ast.Scope) (*ast.BlockStmt, error) {
	block := new(ast.BlockStmt)
	block.DeferStack = make([]*ast.DeferStmt, 0)

	openCurly, ok := p.expect(token.OPEN_CURLY)
	if !ok {
		return nil, fmt.Errorf("expected '{', but got %s", openCurly)
	}
	block.OpenCurly = openCurly.Pos

	var statements []*ast.Node

	for {
		p.skipNewLines()

		tok := p.lex.Peek()
		if tok.Kind == token.CLOSE_CURLY {
			break
		}

		stmt, err := p.parseStmt(block, parentScope, true)
		if err != nil {
			return nil, err
		}
		if stmt == nil {
			continue
		}

		statements = append(statements, stmt)
		if stmt.Kind == ast.KIND_RETURN_STMT {
			block.FoundReturn = true
		}
	}
	block.Statements = statements

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
	block.CloseCurly = closeCurly.Pos
	return block, nil
}

func (p *Parser) ParseIdStmt(parentScope *ast.Scope) (*ast.Node, error) {
	aheadId := p.lex.Peek1()
	switch aheadId.Kind {
	case token.OPEN_PAREN:
		fnCall, err := p.parseFnCall(parentScope)
		return fnCall, err
	case token.COLON_COLON:
		namespaceAccessing, err := p.parseNamespaceAccess(parentScope)
		return namespaceAccessing, err
	case token.DOT:
		return nil, fmt.Errorf("unimplemented field access with dot")
	default:
		return p.parseVar(parentScope)
	}
}

func (p *Parser) parseVar(parentScope *ast.Scope) (*ast.Node, error) {
	variables := make([]*ast.VarId, 0)
	isDecl := false

VarDecl:
	for {
		name, ok := p.expect(token.ID)
		// TODO(errors): add proper error
		if !ok {
			return nil, fmt.Errorf("expected ID")
		}

		variable := &ast.VarId{
			Name:           name,
			NeedsInference: true,
			Type:           nil,
			BackendType:    nil,
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

	expr, err := p.parseAnyExpr([]token.Kind{token.NEWLINE, token.AT, token.SEMICOLON, token.OPEN_CURLY}, parentScope)
	if err != nil {
		return nil, err
	}

	n := new(ast.Node)
	n.Kind = ast.KIND_VAR_STMT
	n.Node = &ast.VarStmt{IsDecl: isDecl, Names: variables, Expr: expr}
	return n, nil
}

// Useful for testing
func parseVarFrom(filename, input string) (*ast.VarId, error) {
	collector := diagnostics.New()

	src := []byte(input)
	loc := new(ast.Loc)
	loc.Name = filename
	lex := lexer.New(loc, src, collector)
	parser := NewWithLex(lex, collector)

	tmpScope := ast.NewScope(nil)
	stmt, err := parser.ParseIdStmt(tmpScope)
	if err != nil {
		return nil, err
	}
	return stmt.Node.(*ast.VarId), nil
}

func (p *Parser) parseCondStmt(parentScope *ast.Scope) (*ast.Node, error) {
	ifCond, err := p.parseIfCond(parentScope)
	if err != nil {
		return nil, err
	}

	p.skipSingleNewLine()

	// TODO: validate proximity of if and elifs
	elifConds, err := p.parseElifConds(parentScope)
	if err != nil {
		return nil, err
	}

	isNewLineBeforeElse := p.skipSingleNewLine()

	elseCond, err := p.parseElseCond(parentScope)
	if err != nil {
		return nil, err
	}

	isNewLine := p.skipSingleNewLine()
	if !isNewLine && elseCond != nil {
		return nil, fmt.Errorf("expected new line at the end of else block")
	}

	if isNewLineBeforeElse && p.lex.NextIs(token.NEWLINE) {
		return nil, fmt.Errorf("invalid isolated else, considering removing the line between if/elif and else")
	}

	n := new(ast.Node)
	n.Kind = ast.KIND_COND_STMT
	n.Node = &ast.CondStmt{IfStmt: ifCond, ElifStmts: elifConds, ElseStmt: elseCond}
	return n, nil
}

func (p *Parser) skipSingleNewLine() bool {
	next := p.lex.Peek().Kind
	isNewLine := next == token.NEWLINE
	if isNewLine {
		p.lex.Skip()
	}
	return isNewLine
}

func (p *Parser) parseIfCond(parentScope *ast.Scope) (*ast.IfElifCond, error) {
	ifToken, ok := p.expect(token.IF)
	// TODO(errors): add proper error
	if !ok {
		return nil, fmt.Errorf("expected 'if'")
	}

	ifExpr, err := p.parseSingleExpr(parentScope)
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
		elifExpr, err := p.parseSingleExpr(parentScope)
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

func (p *Parser) parseAnyExpr(possibleEnds []token.Kind, parentScope *ast.Scope) (*ast.Node, error) {
	exprs, err := p.parseExprList(possibleEnds, parentScope)
	if err != nil {
		return nil, err
	}

	if len(exprs) == 1 {
		return exprs[0], nil
	}

	n := new(ast.Node)
	n.Kind = ast.KIND_TUPLE_EXPR
	n.Node = &ast.TupleExpr{Exprs: exprs}
	return n, nil
}

func (p *Parser) parseSingleExpr(parentScope *ast.Scope) (*ast.Node, error) {
	return p.parseLogical(parentScope)
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
		case token.DOT, token.COLON_COLON:
			fieldAccess, err := p.parseNamespaceAccess(parentScope)
			return fieldAccess, err
		}

		p.lex.Skip()

		n.Kind = ast.KIND_ID_EXPR
		n.Node = idExpr
		return n, nil
	case token.OPEN_PAREN:
		p.lex.Skip() // (
		expr, err := p.parseSingleExpr(parentScope)
		if err != nil {
			return nil, err
		}
		_, ok := p.expect(token.CLOSE_PAREN)
		if !ok {
			return nil, fmt.Errorf("expected closing parenthesis")
		}
		return expr, nil
	default:
		if tok.Kind.IsUntyped() {
			p.lex.Skip()

			n.Kind = ast.KIND_LITERAL_EXPR
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

		expr, err := p.parseSingleExpr(parentScope)
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

	var atOp *ast.AtOperator
	if parser.lex.NextIs(token.AT) {
		atOp, err = parser.parseAtOperator(parentScope)
		if err != nil {
			return nil, err
		}
	}

	n := new(ast.Node)
	n.Kind = ast.KIND_FN_CALL
	n.Node = &ast.FnCall{Name: name, Args: args, AtOp: atOp}
	return n, nil
}

func (p *Parser) parseNamespaceAccess(parentScope *ast.Scope) (*ast.Node, error) {
	id, ok := p.expect(token.ID)
	// TODO(errors): add proper error
	if !ok {
		return nil, fmt.Errorf("error: expected ID")
	}

	_, lookupErr := parentScope.LookupAcrossScopes(id.Name())
	_, isImport := p.file.Imports[id.Name()]

	p.lex.Skip() // ::

	right, err := p.parsePrimary(parentScope)
	if err != nil {
		log.Fatal(err)
	}

	n := new(ast.Node)
	n.Kind = ast.KIND_NAMESPACE_ACCESS

	access := new(ast.NamespaceAccess)
	// NOTE: it is an import when I found it in the import list and there is no
	// symbol with the same name in the parent scope
	access.IsImport = isImport && lookupErr != nil
	access.Left = &ast.IdExpr{Name: id}
	access.Right = right

	n.Node = access
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

	cond, err := parser.parseSingleExpr(parentScope)
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
	loc := new(ast.Loc)
	loc.Name = filename
	lex := lexer.New(loc, src, collector)
	parser := NewWithLex(lex, collector)

	tempScope := ast.NewScope(nil)
	forLoop, err := parser.parseForLoop(tempScope)
	return forLoop.Node.(*ast.ForLoop), err
}

func ParseWhileLoopFrom(input, filename string) (*ast.WhileLoop, error) {
	collector := diagnostics.New()

	src := []byte(input)
	loc := new(ast.Loc)
	loc.Name = filename
	lex := lexer.New(loc, src, collector)
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

	expr, err := p.parseSingleExpr(parentScope)
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

func (p *Parser) parseDefer(block *ast.BlockStmt, parentScope *ast.Scope) (*ast.Node, error) {
	_, ok := p.expect(token.DEFER)
	if !ok {
		return nil, fmt.Errorf("expected 'while'")
	}

	stmt, err := p.parseStmt(block, parentScope, false)
	if err != nil {
		return nil, err
	}

	def := &ast.DeferStmt{Stmt: stmt, Skip: block.FoundReturn}
	block.DeferStack = append(block.DeferStack, def)

	n := new(ast.Node)
	n.Kind = ast.KIND_DEFER_STMT
	n.Node = def
	return n, nil
}
