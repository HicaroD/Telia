package sema

import (
	"errors"
	"fmt"
	"log"
	"reflect"
	"slices"
	"strings"

	"github.com/HicaroD/Telia/internal/ast"
	"github.com/HicaroD/Telia/internal/diagnostics"
	"github.com/HicaroD/Telia/internal/lexer/token"
	"github.com/HicaroD/Telia/internal/source"
)

type sema struct {
	generics         genericState
	mainPackageFound bool
	collector        *diagnostics.Collector

	pkg               *ast.Package
	file              *ast.File
	program           *ast.Program // set by Check(); used by registerTupleType
	currentReturnType *ast.ExprType
	currentYield      *yieldContext
	currentFunction   *ast.FnDecl
	vaInitializer     *ast.FnCall
}

func New(collector *diagnostics.Collector) *sema {
	s := new(sema)
	s.pkg = nil
	s.file = nil
	s.mainPackageFound = false
	s.collector = collector
	return s
}

// registerTupleType appends ty to program.TupleTypes if its shape has not
// been seen before. The deduplication key is the tuple's structural signature
// (element kinds joined with ","). This gives codegen a pre-computed,
// deduplicated list of all typedef structs it must emit.
func (s *sema) registerTupleType(ty *ast.ExprType) {
	if s.program == nil {
		return
	}
	tt := ty.T.(*ast.TupleType)
	// Build a canonical key from the element kinds.
	key := tupleKey(tt)
	for _, existing := range s.program.TupleTypes {
		if tupleKey(existing.T.(*ast.TupleType)) == key {
			return // already registered
		}
	}
	s.program.TupleTypes = append(s.program.TupleTypes, ty)
}

// tupleKey returns a stable string key for a TupleType, used for deduplication.
// The key is the comma-joined list of each element's String() representation.
func tupleKey(tt *ast.TupleType) string {
	parts := make([]string, len(tt.Types))
	for i, t := range tt.Types {
		parts[i] = typeKey(t)
	}
	return strings.Join(parts, ",")
}

func (s *sema) Check(program *ast.Program) error {
	s.program = program
	s.generics = genericState{
		structs:   map[string]*ast.StructDecl{},
		functions: map[string]*ast.FnDecl{},
		resolved:  map[*ast.StructDecl]bool{},
		origins:   map[*ast.StructDecl][]string{},
	}
	s.prepareGenerics(program.Root, map[*ast.Package]bool{})
	if err := s.validateTemplates(program.Root, map[*ast.Package]bool{}); err != nil {
		if !errors.Is(err, diagnostics.COMPILER_ERROR_FOUND) {
			return err
		}
	}
	if err := s.checkPackage(program.Root); err != nil {
		s.collector.ReportError(err, source.Span{})
		return diagnostics.COMPILER_ERROR_FOUND
	}
	if s.collector.HasErrors() {
		return diagnostics.COMPILER_ERROR_FOUND
	}
	if err := s.orderValueTypes(); err != nil {
		s.collector.ReportError(err, source.Span{})
		return diagnostics.COMPILER_ERROR_FOUND
	}
	return nil
}

func (s *sema) checkPackage(pkg *ast.Package) error {
	if pkg == nil {
		return nil
	}
	if pkg.Analyzed {
		return nil
	}

	prevPkg := s.pkg
	defer func() { s.pkg = prevPkg }()
	s.pkg = pkg

	err := s.checkPackageFiles(pkg)
	if err != nil {
		return err
	}

	pkg.Analyzed = true
	return nil
}

func (s *sema) checkPackageFiles(pkg *ast.Package) error {
	if pkg.Loc.Name == "main" {
		return fmt.Errorf("package name is not allowed to be 'main'")
	}

	var hasMainMethod, requiresMainMethod, mainPackageFound bool

	for _, file := range pkg.Files {
		prevFile := s.file
		defer func() { s.file = prevFile }()

		s.file = file

		for _, imp := range file.Imports {
			err := s.checkPackage(imp.Package)
			if err != nil {
				s.collector.ReportError(err, source.Span{})
			}
		}

		if mainPackageFound && file.PkgName != "main" {
			return fmt.Errorf("expected package name to be 'main', but got %s\n", file.PkgName)
		}

		if file.PkgName == "main" {
			if s.mainPackageFound {
				return fmt.Errorf("error: main package already defined somewhere else")
			}
			mainPackageFound = true
			requiresMainMethod = true
		}

		if requiresMainMethod && file.PkgName != "main" {
			// TODO(errors)
			return fmt.Errorf("error: expected package name to be 'main'\n")
		}
		if !requiresMainMethod && file.PkgName != pkg.Loc.Name {
			return fmt.Errorf("error: expected package name to be '%s'\n", pkg.Loc.Name)
		}

		fileHasMain, err := s.checkFile(file)
		if err != nil {
			s.collector.ReportError(err, source.Span{})
		}
		hasMainMethod = hasMainMethod || fileHasMain
	}
	s.mainPackageFound = s.mainPackageFound || mainPackageFound

	// TODO(errors)
	if requiresMainMethod && !hasMainMethod {
		return fmt.Errorf("main package requires a 'main' method as entrypoint\n")
	}

	// TODO(errors)
	if !requiresMainMethod && hasMainMethod {
		return fmt.Errorf("'main' method not allowed on non-main packages\n")
	}

	return nil
}

func (s *sema) checkFile(file *ast.File) (bool, error) {
	var foundMain bool
	failed := false

	for _, node := range file.Body {
		if node.Invalid {
			continue
		}
		switch node.Kind {
		case ast.KIND_ENUM_DECL:
			if len(node.Node.(*ast.EnumDecl).TypeParams) != 0 {
				continue
			}
			if err := s.resolveEnum(node.Node.(*ast.EnumDecl)); err != nil {
				s.collector.ReportError(err, node.DeclarationSpan())
				node.Invalid = true
				failed = true
			}
		case ast.KIND_CONST_DECL:
			if err := s.checkConstDecl(node.Node.(*ast.ConstDecl)); err != nil {
				s.collector.ReportError(err, node.DeclarationSpan())
				node.Invalid = true
				failed = true
			}
		case ast.KIND_FN_DECL:
			fnDecl := node.Node.(*ast.FnDecl)
			if len(fnDecl.TypeParams) != 0 {
				continue
			}
			if !foundMain {
				foundMain = fnDecl.Name.Name() == "main"
			}
			err := s.checkFnDecl(fnDecl, s.pkg.Scope, false, false)
			if err != nil {
				s.collector.ReportError(err, fnDecl.Name.Span)
				node.Invalid = true
				failed = true
			}
		case ast.KIND_EXTERN_DECL:
			externDecl := node.Node.(*ast.ExternDecl)
			err := s.checkExternDecl(externDecl)
			if err != nil {
				s.collector.ReportError(err, node.Span)
				failed = true
			}
		case ast.KIND_STRUCT_DECL:
			structDecl := node.Node.(*ast.StructDecl)
			if len(structDecl.TypeParams) != 0 {
				continue
			}
			err := s.resolveStruct(structDecl, nil)
			if err != nil {
				s.collector.ReportError(err, structDecl.Name.Span)
				node.Invalid = true
				failed = true
			}
		case ast.KIND_PKG_DECL:
			continue
		case ast.KIND_TYPE_ALIAS_DECL:
			continue
		default:
			panic(fmt.Sprintf(
				"unimplemented ast node for sema: %s\n",
				reflect.TypeOf(node.Node),
			))
		}
	}

	if failed {
		return foundMain, diagnostics.COMPILER_ERROR_FOUND
	}
	return foundMain, nil
}

func (sema *sema) checkStructDecl(structDecl *ast.StructDecl) error {
	// TODO: optimize it
	// Depending on the number of fields, it might be better to use a slice or even bitmaps
	// to store seen fields
	seenFields := make(map[string]*token.Token, len(structDecl.Fields))
	for _, field := range structDecl.Fields {
		fieldName := field.Name.Name()
		if previous := seenFields[fieldName]; previous != nil {
			return diagnostics.Diag{
				Title: "Duplicate field",
				Message: fmt.Sprintf(
					"duplicate field '%s' in struct '%s'",
					fieldName,
					structDecl.Name.Name(),
				),
				Span: field.Name.Span,
				Related: []diagnostics.Related{
					{Span: previous.Span, Message: "Previously declared here"},
				},
			}
		}
		seenFields[fieldName] = field.Name
	}
	return nil
}

func (sema *sema) checkFnDecl(
	function *ast.FnDecl,
	declScope *ast.Scope,
	fromImportPackage, isArg bool,
) error {
	if err := sema.resolveSignature(function); err != nil {
		return err
	}
	previousReturnType := sema.currentReturnType
	previousFunction := sema.currentFunction
	sema.currentFunction = function
	defer func() { sema.currentFunction = previousFunction }()
	previousYield := sema.currentYield
	sema.currentYield = nil
	sema.currentReturnType = function.RetType
	defer func() { sema.currentReturnType = previousReturnType; sema.currentYield = previousYield }()
	if token.IsAnyBuiltin(function.Name.Name()) {
		return fmt.Errorf("impossible to redefine builtin '%s'", function.Name.Name())
	}

	err := sema.checkFnAttributes(function.Attributes)
	if err != nil {
		return err
	}

	err = sema.checkBlock(
		function.Block,
		function.RetType,
		function.Scope,
		declScope,
		fromImportPackage,
		isArg,
	)

	if !function.RetType.IsVoid() {
		fnReturns := sema.checkNonVoidFunctionReturns(function.Block)
		if !fnReturns {
			return fmt.Errorf(
				"%s non-void function '%s' must always return a value",
				function.Name.Pos,
				function.Name.Name(),
			)
		}
	}

	// Register tuple return types so codegen can emit their typedef structs
	// upfront without a separate discovery pass.
	if function.RetType.Kind == ast.EXPR_TYPE_TUPLE {
		sema.registerTupleType(function.RetType)
	}

	return err
}

func (sema *sema) checkNonVoidFunctionReturns(block *ast.BlockStmt) bool {
	var hasReturn bool
	for _, stmt := range block.Statements {
		switch stmt.Kind {
		case ast.KIND_COND_STMT:
			cond := stmt.Node.(*ast.CondStmt)

			ifReturns := sema.checkNonVoidFunctionReturns(cond.IfStmt.Block)

			elseReturns := false
			if cond.ElseStmt != nil {
				elseReturns = sema.checkNonVoidFunctionReturns(cond.ElseStmt.Block)
			}

			allElifReturns := true
			for _, elif := range cond.ElifStmts {
				elifReturns := sema.checkNonVoidFunctionReturns(elif.Block)
				if !elifReturns {
					allElifReturns = false
					break
				}
			}

			if ifReturns && (cond.ElseStmt != nil && elseReturns) && allElifReturns {
				hasReturn = true
			}
		case ast.KIND_RETURN_STMT:
			hasReturn = true
		}
	}
	return hasReturn
}

func (sema *sema) checkBlockStmt(
	block *ast.BlockStmt,
	declScope *ast.Scope,
	referenceScope *ast.Scope,
	returnTy *ast.ExprType,
	fromImportPackage, isArg bool,
) error {
	failed := false
	for _, stmt := range block.Statements {
		count := len(sema.collector.Diags)
		err := sema.checkStmt(stmt, referenceScope, declScope, returnTy, fromImportPackage, isArg)
		if err != nil {
			// Keep instantiation context attached while unwinding to its call site.
			if len(sema.generics.chain) > 0 {
				return err
			}
			if len(sema.collector.Diags) == count {
				sema.collector.ReportError(err, stmt.Span)
			}
			failed = true
		}
	}
	if failed {
		return diagnostics.COMPILER_ERROR_FOUND
	}
	return nil
}

func (sema *sema) checkExternAttributes(attributes ast.Attributes) error {
	if attributes.Linkage != "" {
		return fmt.Errorf("'linkage' is not a valid attribute for extern declaration\n")
	}

	// TODO: check linkage

	return validateCallingConvention(
		attributes.DefaultCallingConvention,
		attributes.DefaultCallingConventionMacro,
	)
}

func (sema *sema) checkExternDecl(extern *ast.ExternDecl) error {
	if err := sema.checkExternAttributes(extern.Attributes); err != nil {
		return err
	}
	for _, d := range extern.Enums {
		if err := sema.resolveEnum(d); err != nil {
			return err
		}
	}
	for _, d := range extern.Constants {
		if err := sema.checkConstDecl(d); err != nil {
			return err
		}
	}
	for _, d := range extern.Records {
		if extern.Attributes.Header == "" {
			return diagnostics.Diag{
				Title:   "Foreign record requires a header",
				Message: "foreign records must belong to a header-owned extern block",
				Span:    d.Name.Span,
				Help:    "Add @[header=\"path/to/header.h\"] to the extern block.",
			}
		}
		if len(d.TypeParams) != 0 {
			return diagnostics.Diag{
				Title:   "Generic foreign record",
				Message: "header-owned records cannot have type parameters",
				Span:    d.Name.Span,
				Help:    "Declare the concrete C type without type parameters.",
			}
		}
		if err := sema.resolveStruct(d, nil); err != nil {
			return err
		}
	}
	for _, d := range extern.Opaques {
		if extern.Attributes.Header == "" {
			return diagnostics.Diag{
				Title:   "Opaque type requires a header",
				Message: "opaque declarations must belong to a header-owned extern block",
				Span:    d.Name.Span,
				Help:    "Add @[header=\"path/to/header.h\"] to the extern block.",
			}
		}
	}

	for _, proto := range extern.Prototypes {
		err := sema.checkExternPrototype(extern, proto)
		if err != nil {
			return err
		}
	}

	return nil
}

var VALID_FUNCTION_LINKAGES []string = []string{
	"external", "internal", "weak", "link_once",
}

func (sema *sema) checkFnAttributes(attributes ast.Attributes) error {
	if err := validateCallingConvention(attributes.CallingConvention, attributes.CallingConventionMacro); err != nil {
		return err
	}
	linkageFound := false

	if attributes.Linkage != "" {
		if slices.Contains(VALID_FUNCTION_LINKAGES, attributes.Linkage) {
			linkageFound = true
		}

		if !linkageFound {
			return fmt.Errorf("invalid linkage type: %s\n", attributes.Linkage)
		}
	}

	return nil
}

func (sema *sema) checkExternPrototype(extern *ast.ExternDecl, proto *ast.Proto) error {
	if err := checkCVariadicParams(proto.Params); err != nil {
		return err
	}
	if proto.Attributes.CallingConvention == "" && proto.Attributes.CallingConventionMacro == "" {
		proto.Attributes.CallingConvention = extern.Attributes.DefaultCallingConvention
		proto.Attributes.CallingConventionMacro = extern.Attributes.DefaultCallingConventionMacro
		if proto.Attributes.CallingConvention == "" &&
			proto.Attributes.CallingConventionMacro == "" {
			proto.Attributes.CallingConvention = "c"
		}
	}
	err := sema.checkFnAttributes(proto.Attributes)
	if err != nil {
		return err
	}

	symbols := make(map[string]*token.Token, len(proto.Params.Fields))
	for _, param := range proto.Params.Fields {
		// Id
		if previous, found := symbols[param.Name.Name()]; found {
			return diagnostics.Diag{
				Title: "Duplicate parameter", Message: fmt.Sprintf("parameter '%s' is declared more than once in extern function '%s'", param.Name.Name(), proto.Name.Name()),
				Span: param.Name.Span, Help: "Give each parameter a unique name.",
				Related: []diagnostics.Related{
					{Span: previous.Span, Message: "previous parameter declaration"},
				},
			}
		}
		symbols[param.Name.Name()] = param.Name

		// Type
		revealedTy, err := sema.resolveType(param.Type, extern.Scope)
		if err != nil {
			return err
		}
		param.Type = revealedTy
		if err := checkForeignCallbacks(revealedTy, param.Name.Span); err != nil {
			return err
		}
		if revealedTy.Kind == ast.EXPR_TYPE_VA_LIST_OWNER {
			return rejectVaListEscape(revealedTy, param.Name.Span)
		}
		if err := checkForeignPointerDeclarator(revealedTy, param.Name.Span); err != nil {
			return err
		}
		if revealedTy.Kind == ast.EXPR_TYPE_ARRAY || revealedTy.Kind == ast.EXPR_TYPE_SLICE {
			return diagnostics.Diag{
				Title:   "Foreign array parameter",
				Message: "Telia wrapper arrays and slices are not C array parameters",
				Span:    param.Name.Span,
				Help:    "Declare an explicit element pointer, such as *u8, and pass &array[0]; pass the length separately when required.",
			}
		}
		if err := rejectOpaqueValue(revealedTy, param.Name.Span); err != nil {
			return err
		}
	}

	revealedTy, err := sema.resolveType(proto.RetType, extern.Scope)
	if err != nil {
		return err
	}
	proto.RetType = revealedTy
	if err := checkForeignCallbacks(revealedTy, proto.Name.Span); err != nil {
		return err
	}
	if err := checkForeignPointerDeclarator(revealedTy, proto.Name.Span); err != nil {
		return err
	}
	if err := rejectVaListEscape(revealedTy, proto.Name.Span); err != nil {
		return err
	}
	if err := rejectOpaqueValue(revealedTy, proto.Name.Span); err != nil {
		return err
	}

	return nil
}

func (sema *sema) checkBlock(
	block *ast.BlockStmt,
	returnTy *ast.ExprType,
	referenceScope *ast.Scope,
	declScope *ast.Scope,
	fromImportPackage, isArg bool,
) error {
	return sema.checkBlockStmt(block, declScope, referenceScope, returnTy, fromImportPackage, isArg)
}

func (sema *sema) checkStmt(
	stmt *ast.Node,
	referenceScope *ast.Scope,
	declScope *ast.Scope,
	returnTy *ast.ExprType,
	fromImportPackage, isArg bool,
) error {
	switch stmt.Kind {
	case ast.KIND_YIELD_STMT:
		if sema.currentYield == nil {
			return diagnostics.Diag{
				Title:   "Yield outside match",
				Message: "yield is only valid inside a match arm block",
				Span:    stmt.Span,
				Help:    "Use return to return from a function, or yield inside a match arm.",
			}
		}
		expr := stmt.Node.(*ast.YieldStmt).Value
		ty, explicit, err := sema.inferExprTypeWithoutContext(
			expr,
			referenceScope,
			declScope,
			fromImportPackage,
			false,
		)
		if err != nil {
			return err
		}
		sema.currentYield.sites = append(
			sema.currentYield.sites,
			yieldSite{expr: expr, scope: referenceScope, ty: ty, explicit: explicit},
		)
		return nil
	case ast.KIND_CONST_DECL:
		d := stmt.Node.(*ast.ConstDecl)
		d.Scope = referenceScope
		if err := referenceScope.Insert(d.Name.Name(), stmt); err != nil {
			return diagnostics.Diag{
				Title:   "Duplicate declaration",
				Message: fmt.Sprintf("constant '%s' already declared on scope", d.Name.Name()),
				Span:    d.Name.Span,
				Related: []diagnostics.Related{
					{
						Span:    referenceScope.Nodes[d.Name.Name()].DeclarationSpan(),
						Message: "previous declaration",
					},
				},
			}
		}
		if err := sema.checkConstDecl(d); err != nil {
			stmt.Invalid = true
			return err
		}
		return nil
	case ast.KIND_FN_CALL:
		_, err := sema.checkFnCall(
			stmt.Node.(*ast.FnCall),
			referenceScope,
			declScope,
			returnTy,
			fromImportPackage,
			false,
		)
		return err
	case ast.KIND_VAR_STMT:
		err := sema.checkVar(
			stmt.Node.(*ast.VarStmt),
			referenceScope,
			declScope,
			returnTy,
			fromImportPackage,
			isArg,
		)
		return err
	case ast.KIND_COND_STMT:
		err := sema.checkCondStmt(
			stmt.Node.(*ast.CondStmt),
			returnTy,
			referenceScope,
			declScope,
			fromImportPackage,
			isArg,
		)
		return err
	case ast.KIND_RETURN_STMT:
		returnStmt := stmt.Node.(*ast.ReturnStmt)
		_, err := sema.inferExprTypeWithContext(
			returnStmt.Value,
			returnTy,
			referenceScope,
			declScope,
			fromImportPackage,
			false,
		)
		return err
	case ast.KIND_NAMESPACE_ACCESS:
		_, err := sema.checkNamespaceAccess(
			stmt.Node.(*ast.NamespaceAccess),
			referenceScope,
			declScope,
			fromImportPackage,
		)
		return err
	case ast.KIND_FOR_LOOP_STMT:
		err := sema.checkForLoop(
			stmt.Node.(*ast.ForLoop),
			returnTy,
			referenceScope,
			declScope,
			fromImportPackage,
			isArg,
		)
		return err
	case ast.KIND_WHILE_LOOP_STMT:
		err := sema.checkWhileLoop(
			stmt.Node.(*ast.WhileLoop),
			referenceScope,
			declScope,
			returnTy,
			fromImportPackage,
			isArg,
		)
		return err
	case ast.KIND_DEFER_STMT:
		deferStmt := stmt.Node.(*ast.DeferStmt)
		previousYield := sema.currentYield
		sema.currentYield = nil
		defer func() { sema.currentYield = previousYield }()
		err := sema.checkStmt(
			deferStmt.Stmt,
			referenceScope,
			declScope,
			returnTy,
			fromImportPackage,
			isArg,
		)
		return err
	case ast.KIND_BLOCK_STMT:
		blockStmt := stmt.Node.(*ast.BlockStmt)
		blockScope := ast.NewScope(referenceScope)
		return sema.checkBlock(blockStmt, returnTy, blockScope, declScope, fromImportPackage, isArg)
	case ast.KIND_ASSIGNMENT_STMT:
		assignment := stmt.Node.(*ast.AssignmentStmt)
		err := sema.checkAssignment(
			assignment,
			referenceScope,
			declScope,
			fromImportPackage,
			isArg,
		)
		return err
	default:
		return fmt.Errorf("error: unimplemented statement on sema: %v\n", stmt.Kind)
	}
}

func (sema *sema) checkVar(
	variable *ast.VarStmt,
	referenceScope *ast.Scope,
	declScope *ast.Scope,
	returnTy *ast.ExprType,
	fromImportPackage, isArg bool,
) (resultErr error) {
	previousInitializer := sema.vaInitializer
	sema.vaInitializer = nil
	if variable.IsDecl && len(variable.Names) == 1 &&
		variable.Names[0].Kind == ast.KIND_VAR_ID_STMT {
		sema.vaInitializer = terminalVaCall(variable.Expr)
	}
	defer func() { sema.vaInitializer = previousInitializer }()
	defer func() {
		if resultErr == nil || !variable.IsDecl {
			return
		}
		for _, name := range variable.Names {
			if name.Kind != ast.KIND_VAR_ID_STMT {
				continue
			}
			v := name.Node.(*ast.VarIdStmt)
			if v.N != nil && v.N.Node == v {
				v.N.Invalid = true
			}
		}
	}()
	for _, currentVar := range variable.Names {
		var err error

		if currentVar.Kind == ast.KIND_VAR_ID_STMT {
			err = sema.checkNormalVarDecl(
				currentVar.Node.(*ast.VarIdStmt),
				variable.IsDecl,
				referenceScope,
			)
		} else {
			_, _, err = sema.checkFieldAccessVariable(currentVar.Node.(*ast.FieldAccess), referenceScope)
		}

		if err != nil {
			return err
		}
	}

	// Resolve qualified calls before applying tuple/error assignment rules.
	// The backend uses the resolved declaration, so the namespace is no longer needed.
	var qualifiedRetType *ast.ExprType
	if variable.Expr.Kind == ast.KIND_NAMESPACE_ACCESS {
		ns := variable.Expr.Node.(*ast.NamespaceAccess)
		terminal := ns.Right
		for terminal.Kind == ast.KIND_NAMESPACE_ACCESS {
			terminal = terminal.Node.(*ast.NamespaceAccess).Right
		}
		if terminal.Kind == ast.KIND_FN_CALL {
			var err error
			qualifiedRetType, err = sema.checkNamespaceAccess(
				ns,
				referenceScope,
				declScope,
				fromImportPackage,
			)
			if err != nil {
				return err
			}
			variable.Expr = terminal
		}
	}
	switch variable.Expr.Kind {
	case ast.KIND_TUPLE_LITERAL_EXPR:
		tuple := variable.Expr.Node.(*ast.TupleExpr)
		if len(variable.Names) != len(tuple.Exprs) {
			return fmt.Errorf("more variables than expressions\n")
		}
		for i, name := range variable.Names {
			err := sema.checkVarExpr(
				name,
				tuple.Exprs[i],
				referenceScope,
				declScope,
				fromImportPackage,
				isArg,
			)
			if err != nil {
				return err
			}
		}
		return nil
	case ast.KIND_FN_CALL:
		fnCall := variable.Expr.Node.(*ast.FnCall)
		fnRetType := qualifiedRetType
		if fnRetType == nil {
			var err error
			fnRetType, err = sema.checkFnCall(
				fnCall,
				referenceScope,
				declScope,
				returnTy,
				fromImportPackage,
				false,
			)
			if err != nil {
				return err
			}
		}
		if fnRetType.Kind == ast.EXPR_TYPE_C_VA_LIST {
			return diagnostics.Diag{
				Title:   "Borrowed va_list copy",
				Message: "a va_list borrow cannot be stored in ordinary local value storage",
				Span:    variable.Expr.Span,
				Help:    "Use va_borrow(&owner) directly as a call argument, or create independent storage with va_copy.",
			}
		}
		if fnRetType.Kind == ast.EXPR_TYPE_VA_LIST_OWNER {
			if !variable.IsDecl || len(variable.Names) != 1 ||
				variable.Names[0].Kind != ast.KIND_VAR_ID_STMT ||
				sema.vaInitializer != fnCall {
				return ownerCopyDiagnostic(variable.Expr.Span)
			}
			local := variable.Names[0].Node.(*ast.VarIdStmt)
			if local.Type != nil {
				ty, err := sema.resolveType(local.Type, referenceScope)
				if err != nil {
					return err
				}
				if !ty.Equals(fnRetType) {
					return diagnostics.Diag{
						Title:   "va_list owner type mismatch",
						Message: "a va_list initializer requires local owning storage",
						Span:    local.Name.Span,
						Help:    "Use owner := c::va_copy(...) or annotate it with c::VaList from std::ffi::c.",
					}
				}
			}
			local.Type = fnRetType
			return nil
		}

		// @fail and @catch unwrap the error from the tuple — only non-error elements
		// are assigned to the variables. The error is checked at runtime.
		if fnCall.AtOp != nil &&
			(fnCall.AtOp.Kind == ast.AT_OPERATOR_FAIL || fnCall.AtOp.Kind == ast.AT_OPERATOR_CATCH) &&
			fnRetType.Kind == ast.EXPR_TYPE_TUPLE {
			fnRetType = unwrapErrorFromTuple(fnRetType)
		}

		if fnRetType.Kind == ast.EXPR_TYPE_TUPLE {
			err := sema.checkTupleTypeAssignedToVariable(
				variable.Names,
				fnRetType.T.(*ast.TupleType),
				referenceScope,
			)
			return err
		} else if fnCall.AtOp != nil && (fnCall.AtOp.Kind == ast.AT_OPERATOR_FAIL || fnCall.AtOp.Kind == ast.AT_OPERATOR_CATCH) {
			// @fail/@catch on error-only function: there is no non-error value to assign.
			// Use the bare call form instead: fn() @fail / fn() @catch err { ... }.
			if fnRetType.IsError() {
				pos := fnCall.Name.Pos
				atName := "@fail"
				if fnCall.AtOp.Kind == ast.AT_OPERATOR_CATCH {
					atName = "@catch"
				}
				d := diagnostics.Diag{
					Span: fnCall.Name.Span,
					Message: fmt.Sprintf(
						"%s:%d:%d: %s on error-only function '%s' has no value to assign; use '%s() %s' without variable assignment",
						pos.Filename,
						pos.Line,
						pos.Column,
						atName,
						fnCall.Name.Name(),
						fnCall.Name.Name(),
						atName,
					),
				}
				sema.collector.ReportAndSave(d)
				return diagnostics.COMPILER_ERROR_FOUND
			}
			// @fail on single-var with non-tuple, non-error type: set type
			// directly. checkFnCall already validated the @fail annotation.
			if len(variable.Names) != 1 {
				return fmt.Errorf("more variables than expressions\n")
			}
			varId := variable.Names[0].Node.(*ast.VarIdStmt)
			varId.Type = fnRetType
		} else {
			// TODO(errors)
			if len(variable.Names) != 1 {
				return fmt.Errorf("more variables than expressions\n")
			}
			if qualifiedRetType != nil {
				return sema.checkTupleTypeAssignedToVariable(variable.Names, &ast.TupleType{Types: []*ast.ExprType{fnRetType}}, referenceScope)
			}
			err := sema.checkVarExpr(variable.Names[0], variable.Expr, referenceScope, declScope, fromImportPackage, isArg)
			if err != nil {
				return err
			}
		}
	default:
		// TODO(errors)
		if len(variable.Names) != 1 {
			return fmt.Errorf("more variables than expressions\n")
		}
		err := sema.checkVarExpr(
			variable.Names[0],
			variable.Expr,
			referenceScope,
			declScope,
			fromImportPackage,
			isArg,
		)
		if err != nil {
			return err
		}
	}

	return nil
}

func (sema *sema) checkNormalVarDecl(
	currentVar *ast.VarIdStmt,
	isDecl bool,
	referenceScope *ast.Scope,
) error {
	if isDecl {
		if token.IsAnyBuiltin(currentVar.Name.Name()) {
			return fmt.Errorf("impossible to redefine builtin '%s'", currentVar.Name.Name())
		}

		previous, err := referenceScope.LookupCurrentScope(currentVar.Name.Name())
		if err == nil {
			return diagnostics.Diag{
				Title: "Duplicate declaration",
				Message: fmt.Sprintf(
					"'%s' already declared in the current block",
					currentVar.Name.Name(),
				),
				Span: currentVar.Name.Span,
				Related: []diagnostics.Related{
					{Span: previous.DeclarationSpan(), Message: "Previously declared here"},
				},
			}
		}

		n := new(ast.Node)
		n.Kind = ast.KIND_VAR_ID_STMT
		n.Node = currentVar
		n.Span = currentVar.Name.Span
		currentVar.N = n

		err = referenceScope.Insert(currentVar.Name.Name(), n)
		if err != nil {
			return err
		}
	} else {
		symbol, err := referenceScope.LookupAcrossScopes(currentVar.Name.Name())
		// TODO(errors)
		if err != nil {
			if err == ast.ErrSymbolNotFoundOnScope {
				name := currentVar.Name.Name()
				pos := currentVar.Name.Pos
				return fmt.Errorf("%s '%s' not declared in the current block", pos, name)
			}
			return err
		}
		if symbol.Invalid {
			return diagnostics.COMPILER_ERROR_FOUND
		}
		if symbol.Kind == ast.KIND_CONST_DECL {
			return diagnostics.Diag{
				Title: "Constant reassignment", Message: fmt.Sprintf("cannot reassign constant '%s'", currentVar.Name.Name()),
				Span: currentVar.Name.Span, Help: "Use a variable declared with := if the value needs to change.",
				Related: []diagnostics.Related{{Span: symbol.DeclarationSpan(), Message: "constant declared here"}},
			}
		}
		if currentVar.Type != nil {
			return fmt.Errorf("impossible to define type for variable '%s' since you're reassigning", currentVar.Name.Name())
		}

		currentVar.N = symbol
	}
	return nil
}

func (sema *sema) checkTupleTypeAssignedToVariable(
	variables []*ast.Node,
	tupleTy *ast.TupleType,
	referenceScope *ast.Scope,
) error {
	if len(variables) != len(tupleTy.Types) {
		return fmt.Errorf("expected %d variables, but got %d", len(tupleTy.Types), len(variables))
	}

	for i, ty := range tupleTy.Types {
		n := variables[i]

		if n.Kind == ast.KIND_VAR_ID_STMT {
			currentVar := n.Node.(*ast.VarIdStmt)
			if currentVar.Type == nil {
				currentVar.Type = ty
			} else {
				resolved, err := sema.resolveType(currentVar.Type, referenceScope)
				if err != nil {
					return err
				}
				currentVar.Type = resolved
				if !typeAssignable(ty, resolved) {
					return diagnostics.Diag{
						Title: "Type mismatch", Message: fmt.Sprintf("expected %s, got %s", resolved.T, ty.T),
						Span: currentVar.Name.Span, Help: "Use a destination type compatible with the returned value, or omit the type to infer it.",
					}
				}
			}
		} else {
			currentVar := n.Node.(*ast.FieldAccess)
			_, field, err := sema.getAccessedField(currentVar, referenceScope)
			if err != nil {
				return err
			}
			// TODO(errors)
			if !field.Type.Equals(ty) {
				return fmt.Errorf("type mismatch: expected %v, got %v\n", ty, field.Type)
			}
		}
	}

	return nil
}

func (sema *sema) checkVarExpr(
	variable *ast.Node,
	expr *ast.Node,
	referenceScope, declScope *ast.Scope,
	fromImportPackage, isArg bool,
) error {
	var err error

	switch variable.Kind {
	case ast.KIND_INDEX_EXPR:
		ty, e := sema.inferIndex(
			variable.Node.(*ast.IndexExpr),
			referenceScope,
			declScope,
			fromImportPackage,
		)
		if e != nil {
			return e
		}
		if sema.readOnlyStorage(variable, referenceScope) {
			return constMutation(variable.Span)
		}
		_, err = sema.inferExprTypeWithContext(
			expr,
			ty,
			referenceScope,
			declScope,
			fromImportPackage,
			isArg,
		)
	case ast.KIND_VAR_ID_STMT:
		varId := variable.Node.(*ast.VarIdStmt)
		_, err = sema.checkNormalVarExpr(
			varId,
			expr,
			referenceScope,
			declScope,
			fromImportPackage,
		)
	case ast.KIND_FIELD_ACCESS:
		_, err = sema.checkFieldAccessExpr(
			variable.Node.(*ast.FieldAccess),
			expr,
			referenceScope,
			declScope,
			fromImportPackage,
		)
	case ast.KIND_DEREF_POINTER_EXPR:
		_, err = sema.checkDerefPointerVar(
			variable.Node.(*ast.DerefPointerExpr),
			expr,
			referenceScope,
			declScope,
			fromImportPackage,
			isArg,
		)
	default:
		panic(fmt.Sprintf("unimplemented %v", variable.Kind))
	}

	return err
}

func (s *sema) checkDerefPointerVar(
	deref *ast.DerefPointerExpr,
	expr *ast.Node,
	referenceScope, declScope *ast.Scope,
	fromImportPackage, isArg bool,
) (*ast.ExprType, error) {
	ptr, _, pointerErr := s.inferExprTypeWithoutContext(
		deref.Expr,
		referenceScope,
		declScope,
		fromImportPackage,
		isArg,
	)
	if pointerErr != nil {
		return nil, pointerErr
	}
	if ptr.IsPointer() && ptr.T.(*ast.PointerType).Const {
		return nil, diagnostics.Diag{
			Title:   "Const pointee mutation",
			Message: "cannot write through a pointer to a const pointee",
			Span:    deref.Expr.Span,
			Help:    "Use a mutable pointer only when the API permits writing to this storage.",
		}
	}
	var err error

	varTy, _, err := s.inferDerefPtrExprTypeWithoutContext(
		deref,
		referenceScope,
		declScope,
		fromImportPackage,
		isArg,
	)
	if err != nil {
		return nil, err
	}

	ty, err := s.inferExprTypeWithContext(
		expr,
		varTy,
		referenceScope,
		declScope,
		fromImportPackage,
		isArg,
	)
	if err != nil {
		return nil, err
	}
	if !typeAssignable(ty, varTy) {
		return nil, diagnostics.Diag{
			Title:   "Type mismatch",
			Message: fmt.Sprintf("expected %s, got %s", varTy.T, ty.T),
			Span:    expr.Span,
			Help:    "Use a value compatible with the pointed-to storage type.",
		}
	}

	return ty, nil
}

func (sema *sema) checkNormalVarExpr(
	variable *ast.VarIdStmt,
	expr *ast.Node,
	referenceScope, declScope *ast.Scope,
	fromImportPackage bool,
) (*ast.ExprType, error) {
	if variable.NeedsInference {
		// TODO(errors): need a test for it
		if variable.Type != nil {
			return nil, fmt.Errorf(
				"needs inference, but variable already has a type: %v",
				variable.Type,
			)
		}
		exprTy, _, err := sema.inferExprTypeWithoutContext(
			expr,
			referenceScope,
			declScope,
			fromImportPackage,
			false,
		)
		// TODO(errors)
		if err != nil {
			return nil, err
		}
		if exprTy.Kind == ast.EXPR_TYPE_ARRAY && exprTy.T.(*ast.SequenceType).Foreign {
			return nil, diagnostics.Diag{
				Title:   "Foreign array copy",
				Message: "inline C array storage is not a Telia wrapper-array value",
				Span:    expr.Span,
				Help:    "Copy elements explicitly, borrow a slice, or pass &record.field[0] to a pointer parameter.",
			}
		}
		if exprTy.Kind == ast.EXPR_TYPE_C_VA_LIST {
			return nil, diagnostics.Diag{
				Title:   "Borrowed va_list copy",
				Message: "a borrowed va_list cannot be copied with assignment",
				Span:    expr.Span,
				Help:    "Forward it directly during this call, or create independent owning storage with c::va_copy from std::ffi::c.",
			}
		}
		if exprTy.Kind == ast.EXPR_TYPE_VA_LIST_OWNER &&
			(terminalVaCall(expr) == nil || terminalVaCall(expr) != sema.vaInitializer) {
			return nil, ownerCopyDiagnostic(expr.Span)
		}
		variable.Type = exprTy
		return exprTy, nil
	} else {
		// TODO(errors)
		if variable.Type == nil {
			log.Fatalf("variable does not have a type and it said it does not need inference")
		}

		exprTy, err := sema.inferExprTypeWithContext(expr, variable.Type, referenceScope, declScope, fromImportPackage, false)
		if err == nil && exprTy.Kind == ast.EXPR_TYPE_C_VA_LIST {
			return nil, diagnostics.Diag{Title: "Borrowed va_list copy", Message: "a borrowed va_list cannot be copied or reassigned", Span: expr.Span, Help: "Forward it directly to a c::c_va_list parameter during this call."}
		}
		return exprTy, err
	}
}

func (sema *sema) checkFieldAccessExpr(
	variable *ast.FieldAccess,
	expr *ast.Node,
	referenceScope, declScope *ast.Scope,
	fromImportPackage bool,
) (*ast.ExprType, error) {
	_, field, err := sema.getAccessedField(variable, referenceScope)
	if err != nil {
		return nil, err
	}
	if field.Type.Kind == ast.EXPR_TYPE_ARRAY && field.Type.T.(*ast.SequenceType).Foreign {
		return nil, diagnostics.Diag{
			Title:   "Foreign array copy",
			Message: "C array fields cannot be assigned as whole-array values",
			Span:    variable.Left.Name.Span,
			Help:    "Initialize the field in its record literal or assign individual elements.",
			Related: []diagnostics.Related{
				{Span: field.Name.Span, Message: "C array field declared here"},
			},
		}
	}

	if sema.readOnlyStorage(&ast.Node{Node: variable}, referenceScope) {
		return nil, constMutation(expr.Span)
	}
	exprTy, err := sema.inferExprTypeWithContext(
		expr,
		field.Type,
		referenceScope,
		declScope,
		fromImportPackage,
		false,
	)
	if err != nil {
		return nil, err
	}
	// NOTE: I don't this check is necessary
	if !field.Type.Equals(exprTy) {
		return nil, fmt.Errorf(
			"type mismatch on variable decl, expected %v, got %v",
			field.Type,
			exprTy,
		)
	}
	return exprTy, nil
}

func (sema *sema) getAccessedField(
	fieldAccess *ast.FieldAccess,
	referenceScope *ast.Scope,
	// fromImportPackage bool,
) (*ast.IdExpr, *ast.StructField, error) {
	id := fieldAccess.Left.Name.Name()
	sym, err := referenceScope.LookupAcrossScopes(id)
	if err != nil {
		return nil, nil, err
	}
	if sym.Invalid {
		return nil, nil, diagnostics.COMPILER_ERROR_FOUND
	}

	var ty *ast.ExprType

	switch sym.Kind {
	case ast.KIND_VAR_ID_STMT:
		normalVar := sym.Node.(*ast.VarIdStmt)
		fieldAccess.StructVar = normalVar
		ty = normalVar.Type
	case ast.KIND_PARAM:
		parameter := sym.Node.(*ast.Param)
		fieldAccess.Param = true
		fieldAccess.StructParam = parameter
		ty = parameter.Type
	}

	return sema.resolveFieldAccess(fieldAccess, ty, referenceScope)
}

func (sema *sema) resolveFieldAccess(
	fieldAccess *ast.FieldAccess,
	receiverTy *ast.ExprType,
	referenceScope *ast.Scope,
) (*ast.IdExpr, *ast.StructField, error) {
	fieldID, nestedAccess, err := immediateFieldAccess(fieldAccess)
	if err != nil {
		return nil, nil, err
	}
	base := receiverTy
	readOnly := fieldAccess.ReadOnly
	for base != nil && base.IsPointer() {
		ptr := base.T.(*ast.PointerType)
		readOnly = ptr.Const
		base = ptr.Type
	}
	fieldAccess.ReadOnly = readOnly
	if base != nil && base.Kind == ast.EXPR_TYPE_OPAQUE {
		d := base.T.(*ast.OpaqueType).Decl
		return nil, nil, diagnostics.Diag{
			Title:   "Opaque field access",
			Message: "opaque types do not expose fields",
			Span:    fieldID.Name.Span,
			Help:    "Use a header-declared accessor function instead.",
			Related: []diagnostics.Related{
				{Span: d.Name.Span, Message: "opaque type declared here"},
			},
		}
	}

	structDecl, field, err := sema.resolveFieldOnType(
		receiverTy,
		fieldID.Name.Name(),
		referenceScope,
	)
	if err != nil {
		return nil, nil, err
	}

	fieldAccess.Decl = structDecl
	fieldAccess.AccessedField = field

	if nestedAccess != nil {
		nestedAccess.ReadOnly = readOnly
		return sema.resolveFieldAccess(nestedAccess, field.Type, referenceScope)
	}

	return fieldID, field, nil
}

func immediateFieldAccess(fieldAccess *ast.FieldAccess) (*ast.IdExpr, *ast.FieldAccess, error) {
	switch fieldAccess.Right.Kind {
	case ast.KIND_ID_EXPR:
		return fieldAccess.Right.Node.(*ast.IdExpr), nil, nil
	case ast.KIND_FIELD_ACCESS:
		nested := fieldAccess.Right.Node.(*ast.FieldAccess)
		return nested.Left, nested, nil
	default:
		return nil, nil, fmt.Errorf("other field access type was found during sema")
	}
}

func (sema *sema) resolveFieldOnType(
	ty *ast.ExprType,
	fieldName string,
	referenceScope *ast.Scope,
) (*ast.StructDecl, *ast.StructField, error) {
	switch ty.Kind {
	case ast.EXPR_TYPE_STRUCT:
		structDecl := ty.T.(*ast.StructType).Decl
		field, found := structDecl.FindAttribute(fieldName)
		if !found {
			return nil, nil, fmt.Errorf(
				"attribute '%s' not found on '%s' struct\n",
				fieldName,
				structDecl.Name.Name(),
			)
		}
		return structDecl, field, nil
	case ast.EXPR_TYPE_POINTER:
		ptr := ty.T.(*ast.PointerType)
		switch ptr.Type.Kind {
		case ast.EXPR_TYPE_ID:
			id := ptr.Type.T.(*ast.IdType)
			sym, err := referenceScope.LookupAcrossScopes(id.Name.Name())
			if err != nil {
				return nil, nil, err
			}
			if sym.Kind != ast.KIND_STRUCT_DECL {
				return nil, nil, fmt.Errorf(
					"expected pointee type to be a struct, not %s",
					ty.T,
				)
			}
			structDecl := sym.Node.(*ast.StructDecl)
			field, found := structDecl.FindAttribute(fieldName)
			if !found {
				return nil, nil, fmt.Errorf(
					"attribute '%s' not found on '%s' struct\n",
					fieldName,
					structDecl.Name.Name(),
				)
			}
			return structDecl, field, nil
		case ast.EXPR_TYPE_STRUCT:
			structDecl := ptr.Type.T.(*ast.StructType).Decl
			field, found := structDecl.FindAttribute(fieldName)
			if !found {
				return nil, nil, fmt.Errorf(
					"attribute '%s' not found on '%s' struct\n",
					fieldName,
					structDecl.Name.Name(),
				)
			}
			return structDecl, field, nil
		default:
			return nil, nil, fmt.Errorf(
				"expected pointer to struct type, not %s",
				ty.T,
			)
		}
	case ast.EXPR_TYPE_BASIC:
		basic := ty.T.(*ast.BasicType)
		if basic.Kind == token.ERROR_TYPE {
			if fieldName != "msg" {
				return nil, nil, fmt.Errorf(
					"field '%s' not found on 'error' type, only 'msg' is available",
					fieldName,
				)
			}
			return ast.ErrorStructDecl, ast.ErrorStructDecl.Fields[0], nil
		}
	}

	return nil, nil, fmt.Errorf(
		"expected type to be struct or pointer to struct, but got %s",
		ty.T,
	)
}

func (sema *sema) checkCondStmt(
	condStmt *ast.CondStmt,
	returnTy *ast.ExprType,
	referenceScope *ast.Scope,
	declScope *ast.Scope,
	fromImportPackage, isArg bool,
) error {
	err := sema.checkIfExpr(condStmt.IfStmt.Expr, referenceScope, declScope, fromImportPackage)
	// TODO(errors)
	if err != nil {
		return err
	}

	err = sema.checkBlock(
		condStmt.IfStmt.Block,
		returnTy,
		condStmt.IfStmt.Scope,
		condStmt.IfStmt.Scope,
		fromImportPackage,
		isArg,
	)
	// TODO(errors)
	if err != nil {
		return err
	}

	for i := range condStmt.ElifStmts {
		err := sema.checkIfExpr(
			condStmt.ElifStmts[i].Expr,
			condStmt.ElifStmts[i].Scope,
			declScope,
			fromImportPackage,
		)
		// TODO(errors)
		if err != nil {
			return err
		}
		err = sema.checkBlock(
			condStmt.ElifStmts[i].Block,
			returnTy,
			condStmt.ElifStmts[i].Scope,
			declScope,
			fromImportPackage,
			isArg,
		)
		// TODO(errors)
		if err != nil {
			return err
		}
	}

	if condStmt.ElseStmt != nil {
		err = sema.checkBlock(
			condStmt.ElseStmt.Block,
			returnTy,
			condStmt.ElseStmt.Scope,
			declScope,
			fromImportPackage,
			isArg,
		)
		// TODO(errors)
		if err != nil {
			return err
		}
	}

	return nil
}

func (sema *sema) checkFnCall(
	fnCall *ast.FnCall,
	referenceScope *ast.Scope,
	declScope *ast.Scope,
	returnTy *ast.ExprType,
	fromImportPackage bool,
	isArg bool,
) (*ast.ExprType, error) {
	if fnCall.CalleeExpr != nil {
		ty, _, err := sema.inferExprTypeWithoutContext(
			fnCall.CalleeExpr,
			referenceScope,
			declScope,
			fromImportPackage,
			false,
		)
		if err != nil {
			return nil, err
		}
		if ty.Kind != ast.EXPR_TYPE_FUNCTION {
			return nil, diagnostics.Diag{
				Title:   "Not callable",
				Message: "the field is not a callback value",
				Span:    fnCall.Name.Span,
				Help:    "Call a field whose type is fn(...) Result.",
			}
		}
		f := ty.T.(*ast.FunctionType)
		if err := sema.checkFnCallArgs(fnCall, f.Params, referenceScope, declScope, fromImportPackage); err != nil {
			return nil, err
		}
		return sema.checkCallAtOperator(fnCall, f.Result, referenceScope, fromImportPackage, isArg)
	}
	var symbol *ast.Node
	var err error
	_, localBuiltinLookup := referenceScope.LookupAcrossScopes(fnCall.Name.Name())
	_, declBuiltinLookup := declScope.LookupAcrossScopes(fnCall.Name.Name())
	if isMemoryBuiltin(fnCall.Name.Name()) && localBuiltinLookup != nil &&
		declBuiltinLookup != nil {
		ty, err := sema.checkMemoryBuiltin(fnCall, referenceScope, declScope, fromImportPackage)
		if err != nil {
			return nil, err
		}
		return sema.checkCallAtOperator(fnCall, ty, referenceScope, fromImportPackage, isArg)
	}
	if fnCall.Builtin == "view" || fnCall.Builtin == "slice" {
		ty, err := sema.checkSliceCall(fnCall, referenceScope, declScope, fromImportPackage)
		if err != nil {
			return nil, err
		}
		return sema.checkCallAtOperator(fnCall, ty, referenceScope, fromImportPackage, isArg)
	}
	_, localLookup := referenceScope.LookupAcrossScopes(fnCall.Name.Name())
	_, declLookup := declScope.LookupAcrossScopes(fnCall.Name.Name())
	if fnCall.Name.Name() == "len" && localLookup != nil && declLookup != nil {
		if len(fnCall.Args) != 1 || len(fnCall.TypeArgs) != 0 || fnCall.AtOp != nil {
			return nil, fmt.Errorf("len requires one collection argument")
		}
		ty, _, err := sema.inferExprTypeWithoutContext(
			fnCall.Args[0],
			referenceScope,
			declScope,
			fromImportPackage,
			false,
		)
		if err != nil {
			return nil, err
		}
		if ty.Kind != ast.EXPR_TYPE_ARRAY && ty.Kind != ast.EXPR_TYPE_SLICE {
			return nil, fmt.Errorf("len requires an array or slice")
		}
		fnCall.Builtin, fnCall.BuiltinType = "len", ty
		fnCall.ResultType = concreteBasicType(token.U64_TYPE)
		return fnCall.ResultType, nil
	}

	if isArg {
		symbol, err = referenceScope.LookupAcrossScopes(fnCall.Name.Name())
	} else {
		symbol, err = declScope.LookupAcrossScopes(fnCall.Name.Name())
		if err != nil {
			symbol, err = referenceScope.LookupAcrossScopes(fnCall.Name.Name())
		}
	}

	if err != nil {
		if err == ast.ErrSymbolNotFoundOnScope {
			pos := fnCall.Name.Pos
			functionNotDefined := diagnostics.Diag{
				Span: fnCall.Name.Span,
				Message: fmt.Sprintf(
					"%s:%d:%d: function '%s' not defined on scope",
					pos.Filename,
					pos.Line,
					pos.Column,
					fnCall.Name.Name(),
				),
			}
			sema.collector.ReportAndSave(functionNotDefined)
			return nil, diagnostics.COMPILER_ERROR_FOUND
		}
		return nil, err
	}

	if symbol.Invalid {
		return nil, diagnostics.COMPILER_ERROR_FOUND
	}
	if symbol.Kind == ast.KIND_VA_INTRINSIC {
		return sema.checkVaIntrinsic(fnCall, referenceScope, declScope, fromImportPackage)
	}
	if symbol.Kind == ast.KIND_VAR_ID_STMT || symbol.Kind == ast.KIND_PARAM {
		callee := &ast.IdExpr{Name: fnCall.Name, N: symbol}
		ty, _, err := sema.inferIdExprTypeWithoutContext(callee, referenceScope)
		if err != nil {
			return nil, err
		}
		if ty.Kind == ast.EXPR_TYPE_FUNCTION {
			if len(fnCall.TypeArgs) != 0 {
				return nil, diagnostics.Diag{
					Title:   "Callback type arguments",
					Message: "callback values do not accept type arguments",
					Span:    fnCall.Name.Span,
					Help:    "Call the callback without type arguments.",
				}
			}
			f := ty.T.(*ast.FunctionType)
			if err := sema.checkFnCallArgs(fnCall, f.Params, referenceScope, declScope, fromImportPackage); err != nil {
				return nil, err
			}
			fnCall.Callee = callee
			return sema.checkCallAtOperator(
				fnCall,
				f.Result,
				referenceScope,
				fromImportPackage,
				isArg,
			)
		}
	}
	if symbol.Kind != ast.KIND_FN_DECL && symbol.Kind != ast.KIND_PROTO {
		pos := fnCall.Name.Pos
		notCallable := diagnostics.Diag{
			Span: fnCall.Name.Span,
			Message: fmt.Sprintf(
				"%s:%d:%d: '%s' is not callable",
				pos.Filename,
				pos.Line,
				pos.Column,
				fnCall.Name.Name(),
			),
		}
		sema.collector.ReportAndSave(notCallable)
		return nil, diagnostics.COMPILER_ERROR_FOUND
	}

	// NOTE: It happens when extern is declared without a name, it means prototypes
	// can be accessed globally in the package scope
	var retType *ast.ExprType

	if symbol.Kind == ast.KIND_PROTO {
		if len(fnCall.TypeArgs) != 0 {
			return nil, fmt.Errorf(
				"%s: '%s' does not take type arguments",
				fnCall.Name.Pos,
				fnCall.Name.Name(),
			)
		}
		proto := symbol.Node.(*ast.Proto)
		fnCall.Proto = proto
		err = sema.checkFnCallArgs(
			fnCall,
			proto.Params,
			referenceScope,
			declScope,
			fromImportPackage,
		)
		retType = proto.RetType
	} else {
		fnDecl := symbol.Node.(*ast.FnDecl)
		if len(fnDecl.TypeParams) != 0 {
			fnDecl, err = sema.instantiateFunction(fnDecl, fnCall.TypeArgs, referenceScope, fnCall.Name)
			if err != nil {
				return nil, err
			}
		} else {
			if len(fnCall.TypeArgs) != 0 {
				return nil, fmt.Errorf("%s: '%s' does not take type arguments", fnCall.Name.Pos, fnCall.Name.Name())
			}
			if err := sema.resolveSignature(fnDecl); err != nil {
				return nil, err
			}
		}
		fnCall.Decl = fnDecl
		err = sema.checkFnCallArgs(fnCall, fnDecl.Params, referenceScope, declScope, fromImportPackage)
		retType = fnDecl.RetType
	}

	if err != nil {
		return nil, err
	}

	return sema.checkCallAtOperator(fnCall, retType, referenceScope, fromImportPackage, isArg)
}

func (sema *sema) checkCallAtOperator(
	fnCall *ast.FnCall,
	retType *ast.ExprType,
	referenceScope *ast.Scope,
	fromImportPackage, isArg bool,
) (*ast.ExprType, error) {
	if fnCall.AtOp != nil && fnCall.AtOp.Kind == ast.AT_OPERATOR_FAIL {
		if !returnsError(retType) {
			pos := fnCall.Name.Pos
			d := diagnostics.Diag{
				Span: fnCall.Name.Span,
				Message: fmt.Sprintf(
					"%s:%d:%d: @fail requires error-returning function, but '%s' does not return an error",
					pos.Filename,
					pos.Line,
					pos.Column,
					fnCall.Name.Name(),
				),
			}
			sema.collector.ReportAndSave(d)
			return nil, diagnostics.COMPILER_ERROR_FOUND
		}
	}

	if fnCall.AtOp != nil && fnCall.AtOp.Kind == ast.AT_OPERATOR_CATCH {
		if !returnsError(retType) {
			pos := fnCall.Name.Pos
			d := diagnostics.Diag{
				Span: fnCall.Name.Span,
				Message: fmt.Sprintf(
					"%s:%d:%d: @catch requires error-returning function, but '%s' does not return an error",
					pos.Filename,
					pos.Line,
					pos.Column,
					fnCall.Name.Name(),
				),
			}
			sema.collector.ReportAndSave(d)
			return nil, diagnostics.COMPILER_ERROR_FOUND
		}

		catchOp := fnCall.AtOp.Op.(*ast.CatchAtOperator)
		catchScope := ast.NewScope(referenceScope)

		errVarId := &ast.VarIdStmt{
			Name: catchOp.ErrVarName,
			Type: ast.NewBasicType(token.ERROR_TYPE),
		}
		errNode := &ast.Node{
			Kind: ast.KIND_VAR_ID_STMT,
			Node: errVarId,
		}
		errVarId.N = errNode

		if err := catchScope.Insert(catchOp.ErrVarName.Name(), errNode); err != nil {
			return nil, err
		}

		if err := sema.checkBlock(
			catchOp.Block,
			sema.currentReturnType,
			catchScope,
			sema.pkg.Scope,
			fromImportPackage,
			isArg,
		); err != nil {
			return nil, err
		}
	}

	return retType, nil
}

func (sema *sema) checkFnCallArgs(
	fnCall *ast.FnCall,
	params *ast.Params,
	referenceScope, declScope *ast.Scope,
	fromImportPackage bool,
) error {
	// Calls may be checked again with an expected result type. Restore source
	// arguments before validating; never treat the lowering wrapper as a value.
	var arguments []*ast.Node
	for _, argument := range fnCall.Args {
		if argument.Kind == ast.KIND_VARG_EXPR {
			arguments = append(arguments, argument.Node.(*ast.VarArgsExpr).Args...)
		} else {
			arguments = append(arguments, argument)
		}
	}
	fnCall.Args = arguments
	if params.IsVariadic && len(fnCall.Args) < params.Len {
		pos := fnCall.Name.Pos
		// TODO(errors): show which arguments were passed and which types we
		// were expecting
		notEnoughArguments := diagnostics.Diag{
			Span: fnCall.Name.Span,
			Message: fmt.Sprintf(
				"%s:%d:%d: not enough arguments in call to '%s', expected at least %d arguments",
				pos.Filename,
				pos.Line,
				pos.Column,
				fnCall.Name.Name(),
				params.Len,
			),
		}
		sema.collector.ReportAndSave(notEnoughArguments)
		return diagnostics.COMPILER_ERROR_FOUND
	}

	if !params.IsVariadic && len(fnCall.Args) != params.Len {
		pos := fnCall.Name.Pos
		// TODO(errors): show which arguments were passed and which types we
		// were expecting
		notEnoughArguments := diagnostics.Diag{
			Span: fnCall.Name.Span,
			Message: fmt.Sprintf(
				"%s:%d:%d: not enough arguments in call to '%s', expected %d arguments, but got %d arguments",
				pos.Filename,
				pos.Line,
				pos.Column,
				fnCall.Name.Name(),
				params.Len,
				len(fnCall.Args),
			),
		}
		sema.collector.ReportAndSave(notEnoughArguments)
		return diagnostics.COMPILER_ERROR_FOUND
	}

	requiredArgs := slices.Clone(fnCall.Args[:params.Len])
	variadicArgs := slices.Clone(fnCall.Args[params.Len:])

	for i, requiredArg := range requiredArgs {
		if _, err := sema.inferExprTypeWithContext(requiredArg, params.Fields[i].Type, referenceScope, declScope, fromImportPackage, true); err != nil {
			return err
		}
	}

	if params.IsVariadic {
		variadicParam := params.Fields[params.Len]
		for _, variadicArg := range variadicArgs {
			if variadicParam.Attributes.C {
				// @c marks a C-style variadic (e.g. printf) — the callee accepts
				// any type at the C ABI level. Still resolve the expression (so
				// field accesses, identifiers, etc. are fully wired up) but skip
				// the type constraint check against variadicParam.Type.
				argType, _, err := sema.inferExprTypeWithoutContext(
					variadicArg,
					referenceScope,
					declScope,
					fromImportPackage,
					true,
				)
				if err != nil {
					return err
				}
				if argType.Equals(ast.NewBasicType(token.STRING_TYPE)) {
					return fmt.Errorf(
						"borrowed string in C variadic call requires explicit C string conversion",
					)
				}
				if !cVariadicType(argType) {
					return diagnostics.Diag{
						Title:   "Unsupported C variadic argument",
						Message: "this value does not have a supported C variadic representation",
						Span:    variadicArg.Span,
						Help:    "Pass a C scalar, typed pointer, callback, or payload-free enum; use explicit storage and pointers for aggregates.",
						Related: []diagnostics.Related{
							{
								Span:    variadicParam.Name.Span,
								Message: "C variadic parameter declared here",
							},
						},
					}
				}
			} else {
				if _, err := sema.inferExprTypeWithContext(variadicArg, variadicParam.Type, referenceScope, declScope, fromImportPackage, true); err != nil {
					return err
				}
			}
		}
	}

	fnCall.Args = requiredArgs
	if len(variadicArgs) > 0 {
		varArgs := new(ast.Node)
		varArgs.Kind = ast.KIND_VARG_EXPR
		varArgs.Node = &ast.VarArgsExpr{Args: variadicArgs}
		fnCall.Args = append(fnCall.Args, varArgs)
		fnCall.Variadic = true
	}

	return nil
}

func (sema *sema) checkIfExpr(
	expr *ast.Node,
	referenceScope *ast.Scope,
	declScope *ast.Scope,
	fromImportPackage bool,
) error {
	ifExprTy, _, err := sema.inferExprTypeWithoutContext(
		expr,
		referenceScope,
		declScope,
		fromImportPackage,
		false,
	)
	if err != nil {
		return err
	}
	boolType := ast.NewBasicType(token.BOOL_TYPE)
	if !ifExprTy.Equals(boolType) {
		return fmt.Errorf("error: expected bool type on if expr, but got %s\n", ifExprTy.T)
	}
	return nil
}

func (s *sema) inferExprTypeWithContext(
	expr *ast.Node,
	expectedType *ast.ExprType,
	referenceScope *ast.Scope,
	declScope *ast.Scope,
	fromImportPackage bool,
	isArg bool,
) (result *ast.ExprType, resultErr error) {
	defer func() { resultErr = diagnostics.WithSpan(resultErr, expr.Span) }()
	resolved, err := s.resolveType(expectedType, referenceScope)
	if err != nil {
		return nil, err
	}
	*expectedType = *resolved
	switch expr.Kind {
	case ast.KIND_MATCH_EXPR:
		return s.checkMatch(expr.Node.(*ast.MatchExpr), expectedType, referenceScope)
	case ast.KIND_ARRAY_LITERAL_EXPR:
		lit := expr.Node.(*ast.ArrayLiteralExpr)
		if expectedType.Kind != ast.EXPR_TYPE_ARRAY {
			return nil, fmt.Errorf("array literal requires an array type")
		}
		t := expectedType.T.(*ast.SequenceType)
		if len(lit.Values) != 0 && uint64(len(lit.Values)) != t.Length {
			return nil, fmt.Errorf("array literal requires exactly %d elements", t.Length)
		}
		for _, value := range lit.Values {
			if _, err := s.inferExprTypeWithContext(value, t.Element, referenceScope, declScope, fromImportPackage, isArg); err != nil {
				return nil, err
			}
		}
		lit.Type = expectedType
		return expectedType, nil
	case ast.KIND_INDEX_EXPR:
		actual, err := s.inferIndex(
			expr.Node.(*ast.IndexExpr),
			referenceScope,
			declScope,
			fromImportPackage,
		)
		if err != nil {
			return nil, err
		}
		if !typeAssignable(actual, expectedType) {
			return nil, diagnostics.Diag{
				Title:   "Type mismatch",
				Message: fmt.Sprintf("expected %s, got %s", expectedType.T, actual.T),
				Span:    expr.Span,
				Help:    "Use an element value with a compatible type; numeric widths and signedness must match.",
			}
		}
		return actual, nil
	case ast.KIND_LITERAL_EXPR:
		return s.inferLiteralExprTypeWithContext(expr.Node.(*ast.LiteralExpr), expectedType)
	case ast.KIND_ID_EXPR:
		return s.inferIdExprTypeWithContext(expr.Node.(*ast.IdExpr), expectedType, referenceScope)
	case ast.KIND_BINARY_EXPR:
		ty, _, err := s.inferBinaryExprType(
			expr.Node.(*ast.BinExpr),
			expectedType,
			referenceScope,
			declScope,
			fromImportPackage,
		)
		return ty, err
	case ast.KIND_UNARY_EXPR:
		ty, _, err := s.inferUnaryExprType(
			expr.Node.(*ast.UnaryExpr),
			expectedType,
			referenceScope,
			declScope,
			fromImportPackage,
		)
		return ty, err
	case ast.KIND_FN_CALL:
		actual, err := s.inferFnCallExprTypeWithContext(
			expr.Node.(*ast.FnCall),
			referenceScope,
			declScope,
			fromImportPackage,
			isArg,
		)
		if err != nil {
			return nil, err
		}
		if !typeAssignable(actual, expectedType) {
			return nil, fmt.Errorf("type mismatch: expected %s, got %s", expectedType.T, actual.T)
		}
		return actual, nil
	case ast.KIND_VOID_EXPR:
		return s.inferVoidExprTypeWithContext(expectedType)
	case ast.KIND_TUPLE_LITERAL_EXPR:
		return s.inferTupleExprTypeWithContext(
			expr.Node.(*ast.TupleExpr),
			expectedType,
			referenceScope,
			declScope,
			fromImportPackage,
		)
	case ast.KIND_STRUCT_EXPR:
		return s.inferStructLiteralExprWithContext(
			expr.Node.(*ast.StructLiteralExpr),
			expectedType,
			referenceScope,
			declScope,
			fromImportPackage,
		)
	case ast.KIND_FIELD_ACCESS:
		return s.inferStructFieldAccessExprWithContext(
			expr.Node.(*ast.FieldAccess),
			expectedType,
			referenceScope,
		)
	case ast.KIND_DEREF_POINTER_EXPR:
		return s.inferDerefPtrExprTypeWithContext(
			expr.Node.(*ast.DerefPointerExpr),
			expectedType,
			referenceScope,
			declScope,
			fromImportPackage,
			isArg,
		)
	case ast.KIND_ADDRESS_OF_EXPR:
		return s.inferPtrExprWithContext(
			expr.Node.(*ast.AddressOfExpr),
			expectedType,
			referenceScope,
			declScope,
			fromImportPackage,
			isArg,
		)
	case ast.KIND_NULLPTR_EXPR:
		return s.inferNullptrExprTypeWithContext(
			expr.Node.(*ast.NullPtrExpr),
			expectedType,
			// referenceScope,
			// declScope,
			// fromImportPackage,
			// isArg,
		)
	case ast.KIND_NAMESPACE_ACCESS:
		namespace := expr.Node.(*ast.NamespaceAccess)
		if namespace.IsImport && namespace.Right.Kind == ast.KIND_ID_EXPR {
			imp := s.file.Imports[namespace.Left.Name.Name()]
			return s.inferImportedValue(
				namespace.Right.Node.(*ast.IdExpr),
				imp.Package.Scope,
				expectedType,
			)
		}
		actual, err := s.checkNamespaceAccess(
			expr.Node.(*ast.NamespaceAccess),
			referenceScope,
			declScope,
			fromImportPackage,
		)
		if err != nil {
			return nil, err
		}
		if !typeAssignable(actual, expectedType) {
			if actual.Kind == ast.EXPR_TYPE_ENUM {
				return nil, enumTypeMismatch(actual, expectedType, expr.Span)
			}
			return nil, fmt.Errorf("type mismatch: expected %s, got %s", expectedType.T, actual.T)
		}
		return actual, nil
	default:
		panic(fmt.Sprintf("infer expr type with context - unimplemented expr: %v\n", expr))
	}
}

func (sema *sema) inferLiteralExprTypeWithContext(
	literal *ast.LiteralExpr,
	expectedType *ast.ExprType,
) (*ast.ExprType, error) {
	// TODO(errors)
	if literal.Type.Kind != ast.EXPR_TYPE_BASIC {
		return nil, fmt.Errorf("error: expected literal expression type to be basic")
	}
	actualBasicType := literal.Type.T.(*ast.BasicType)
	// A literal may directly name static C storage. Existing string values
	// still require an explicit, owned conversion before crossing the C ABI.
	if actualBasicType.Kind == token.STRING_TYPE && expectedType.Kind == ast.EXPR_TYPE_BASIC &&
		expectedType.T.(*ast.BasicType).Kind == token.CSTRING_TYPE {
		for _, b := range literal.Value {
			if b == 0 {
				return nil, fmt.Errorf("C string literal contains an embedded zero byte")
			}
		}
		literal.Type = expectedType
		return expectedType, nil
	}

	if actualBasicType.Kind == token.ERROR_TYPE {
		actualBasicType.Explicit = true
	}

	switch expectedType.Kind {
	case ast.EXPR_TYPE_BASIC:
		expectedBasicType := expectedType.T.(*ast.BasicType)
		_, err := sema.inferBasicExprTypeWithContext(actualBasicType, expectedBasicType)
		if err != nil {
			return nil, diagnostics.Diag{Title: "Type mismatch", Message: err.Error()}
		}
	case ast.EXPR_TYPE_POINTER:
		// TODO(errors)
		return nil, fmt.Errorf("expected a pointer expression, not a literal")
	}

	return expectedType, nil
}

func (sema *sema) inferIdExprTypeWithContext(
	id *ast.IdExpr,
	expectedType *ast.ExprType,
	referenceScope *ast.Scope,
) (*ast.ExprType, error) {
	symbol, err := referenceScope.LookupAcrossScopes(id.Name.Name())
	// TODO(errors)
	if err != nil {
		return nil, diagnostics.Diag{
			Title: "Undeclared name",
			Message: fmt.Sprintf(
				"symbol not found on scope: '%s' is not declared.",
				id.Name.Name(),
			),
			Span: id.Name.Span,
		}
	}
	id.N = symbol
	if symbol.Invalid {
		return nil, diagnostics.COMPILER_ERROR_FOUND
	}

	var idTy *ast.ExprType

	switch symbol.Kind {
	case ast.KIND_VA_INTRINSIC:
		return nil, vaIntrinsicValueDiagnostic(id.Name.Span)
	case ast.KIND_FN_DECL, ast.KIND_PROTO:
		actual, err := sema.functionValue(id, symbol)
		if err != nil {
			return nil, err
		}
		if !actual.Equals(expectedType) {
			return nil, diagnostics.Diag{
				Title:   "Callback signature mismatch",
				Message: "function does not have the required callback signature",
				Span:    id.Name.Span,
				Help:    "Match the parameter types, return type, and ABI; annotate C callback aliases and implementations with @[cc=\"c\"] or the matching cc_macro.",
				Related: []diagnostics.Related{
					{Span: symbol.DeclarationSpan(), Message: "function declared here"},
				},
			}
		}
		return expectedType, nil
	case ast.KIND_CONST_DECL:
		return sema.inferConstReference(id, symbol.Node.(*ast.ConstDecl), expectedType)
	case ast.KIND_VAR_ID_STMT:
		idTy = symbol.Node.(*ast.VarIdStmt).Type
	case ast.KIND_PARAM:
		param := symbol.Node.(*ast.Param)
		if err := rejectCVariadicValue(param, id.Name.Span); err != nil {
			return nil, err
		}
		idTy = param.Type
	default:
		// TODO(errors)
		return nil, fmt.Errorf("expected to be a variable or parameter, but got %v", symbol.Kind)
	}

	switch idTy.Kind {
	case ast.EXPR_TYPE_VA_LIST_OWNER:
		return nil, ownerCopyDiagnostic(id.Name.Span)
	case ast.EXPR_TYPE_C_VA_LIST:
		if !idTy.Equals(expectedType) {
			return nil, diagnostics.Diag{
				Title:   "va_list type mismatch",
				Message: "a borrowed va_list is not a raw pointer or ordinary value",
				Span:    id.Name.Span,
				Help:    "Forward it to a c::c_va_list parameter during the current call.",
			}
		}
		return expectedType, nil
	case ast.EXPR_TYPE_FUNCTION:
		if !idTy.Equals(expectedType) {
			return nil, diagnostics.Diag{
				Title:   "Callback signature mismatch",
				Message: "callback value has an incompatible signature",
				Span:    id.Name.Span,
				Help:    "Use a callback with matching parameter types, return type, and native or foreign ABI.",
			}
		}
		return expectedType, nil
	case ast.EXPR_TYPE_ENUM:
		if !expectedType.Equals(idTy) {
			return nil, enumTypeMismatch(idTy, expectedType, id.Name.Span)
		}
		return idTy, nil
	case ast.EXPR_TYPE_ARRAY, ast.EXPR_TYPE_SLICE:
		if !idTy.Equals(expectedType) {
			return nil, fmt.Errorf("collection type mismatch")
		}
		return idTy, nil
	case ast.EXPR_TYPE_BASIC:
		actualBasicType := idTy.T.(*ast.BasicType)

		var err error

		switch expectedType.Kind {
		case ast.EXPR_TYPE_BASIC:
			expectedBasicType := expectedType.T.(*ast.BasicType)
			ty, err := sema.inferBasicExprTypeWithContext(actualBasicType, expectedBasicType)
			if err != nil {
				return nil, err
			}
			expectedType.T = ty
		case ast.EXPR_TYPE_POINTER:
			expectedPtrType := expectedType.T.(*ast.PointerType)
			ty, err := sema.inferIdExprTypeWithContext(
				id,
				expectedPtrType.Type,
				referenceScope,
			)
			if err != nil {
				return nil, err
			}
			return ty, nil
		default:
			return nil, fmt.Errorf("type mismatch: expected %s, got %s", expectedType.T, idTy.T)
		}
		return expectedType, err
	case ast.EXPR_TYPE_POINTER:
		if !pointerAssignable(idTy, expectedType) {
			if idTy.T.(*ast.PointerType).Const &&
				(expectedType.Equals(ast.RAWPTR_TYPE) || expectedType.IsPointer() && !expectedType.T.(*ast.PointerType).Const) {
				return nil, diagnostics.Diag{
					Title:   "Const qualifier removed",
					Message: "cannot discard the pointer's immediate-pointee const qualifier",
					Span:    id.Name.Span,
					Help:    "Keep the destination const-qualified or copy into separately owned mutable storage.",
				}
			}
			return nil, fmt.Errorf(
				"type mismatch, expected %s, but got %s\n",
				expectedType.T,
				idTy.T,
			)
		}
		return expectedType, nil
	case ast.EXPR_TYPE_STRUCT:
		if !expectedType.Equals(idTy) {
			return nil, fmt.Errorf("type mismatch: expected %s, got %s", expectedType.T, idTy.T)
		}
		return expectedType, nil
	default:
		// TODO(errors)
		panic("unimplemented infer id expr type with context")
	}
}

func (s *sema) inferDerefPtrExprTypeWithContext(
	deref *ast.DerefPointerExpr,
	expectedType *ast.ExprType,
	referenceScope *ast.Scope,
	declScope *ast.Scope,
	fromImportPackage bool,
	isArg bool,
) (*ast.ExprType, error) {
	innerExpr := deref.Expr

	ty, _, err := s.inferExprTypeWithoutContext(
		innerExpr,
		referenceScope,
		declScope,
		fromImportPackage,
		isArg,
	)
	if err != nil {
		return nil, err
	}
	if !ty.IsPointer() {
		return nil, fmt.Errorf("cannot dereference non-pointer expression: %s\n", ty.T)
	}

	pointeeType := ty.T.(*ast.PointerType)
	if !typeAssignable(pointeeType.Type, expectedType) {
		return nil, fmt.Errorf(
			"type mismatch, expected %s, got %s\n",
			expectedType.T,
			pointeeType.Type.T,
		)
	}
	if err := rejectOpaqueValue(pointeeType.Type, deref.Expr.Span); err != nil {
		return nil, err
	}
	deref.Type = pointeeType.Type
	deref.ReadOnly = pointeeType.Const
	return pointeeType.Type, nil
}

func (s *sema) inferDerefPtrExprTypeWithoutContext(
	deref *ast.DerefPointerExpr,
	referenceScope *ast.Scope,
	declScope *ast.Scope,
	fromImportPackage bool,
	isArg bool,
) (*ast.ExprType, bool, error) {
	innerExpr := deref.Expr

	ty, _, err := s.inferExprTypeWithoutContext(
		innerExpr,
		referenceScope,
		declScope,
		fromImportPackage,
		isArg,
	)
	if err != nil {
		return nil, false, err
	}
	if !ty.IsPointer() {
		return nil, false, fmt.Errorf("cannot dereference non-pointer type: %s\n", ty.T)
	}

	pointeeType := ty.T.(*ast.PointerType)
	if err := rejectOpaqueValue(pointeeType.Type, deref.Expr.Span); err != nil {
		return nil, false, err
	}
	deref.Type = pointeeType.Type
	deref.ReadOnly = pointeeType.Const
	return pointeeType.Type, !pointeeType.Type.Explicit, nil
}

func (s *sema) inferBinaryExprType(
	binary *ast.BinExpr,
	expectedType *ast.ExprType,
	referenceScope *ast.Scope,
	declScope *ast.Scope,
	fromImportPackage bool,
) (*ast.ExprType, bool, error) {
	if binary.Op == token.SHIFT_LEFT || binary.Op == token.SHIFT_RIGHT {
		return s.inferShift(binary, expectedType, referenceScope, declScope, fromImportPackage)
	}
	lhs, rhs, ctx, err := s.ensureBinaryOperatorsAreTheSame(
		binary,
		expectedType,
		referenceScope,
		declScope,
		fromImportPackage,
	)
	if err != nil {
		return nil, false, err
	}

	if s.isErrorNilComparison(lhs, rhs) {
		return ast.NewBasicType(token.BOOL_TYPE), ctx, nil
	}
	if lhs.Kind == ast.EXPR_TYPE_POINTER &&
		(binary.Op == token.EQUAL_EQUAL || binary.Op == token.BANG_EQUAL) {
		return ast.NewBasicType(token.BOOL_TYPE), ctx, nil
	}
	if lhs.Kind == ast.EXPR_TYPE_FUNCTION && rhs.Kind == ast.EXPR_TYPE_FUNCTION &&
		(binary.Op == token.EQUAL_EQUAL || binary.Op == token.BANG_EQUAL) {
		return ast.NewBasicType(token.BOOL_TYPE), ctx, nil
	}

	commonType := lhs // since they are the same
	binary.OperandType = commonType
	if commonType.Kind == ast.EXPR_TYPE_ENUM {
		d := commonType.T.(*ast.EnumType).Decl
		if d.HasPayload {
			return nil, false, diagnostics.Diag{
				Title:   "Payload enum comparison",
				Message: "payload enums do not have implicit equality",
				Span:    binary.Left.Span,
				Help:    "Use match to inspect the variant and compare the fields explicitly.",
				Related: []diagnostics.Related{
					{Span: d.Name.Span, Message: "payload enum declared here"},
				},
			}
		}
		if binary.Op != token.EQUAL_EQUAL && binary.Op != token.BANG_EQUAL {
			return nil, false, diagnostics.Diag{
				Title:   "Invalid enum operator",
				Message: "enum values only support == and !=",
				Help:    "Use integer aliases and constants for arithmetic or flags.",
			}
		}
		result := ast.NewBasicType(token.BOOL_TYPE)
		if expectedType != nil && !result.Equals(expectedType) {
			return nil, false, fmt.Errorf("type mismatch: expected %s, got bool", expectedType.T)
		}
		return result, ctx, nil
	}

	// TODO(errors)
	validation, exists := ast.BinaryOperators[binary.Op]
	if !exists {
		return nil, false, fmt.Errorf("invalid operator: %s", binary.Op)
	}

	valid := false
	for _, validType := range validation.ValidTypes {
		if operatorTypeMatches(validType, commonType) {
			valid = true
			break
		}
	}
	if !valid {
		if binary.Op.IsBitwiseOp() {
			return nil, false, diagnostics.Diag{
				Title: "Invalid bitwise operands",
				Message: fmt.Sprintf(
					"invalid operand types for %s: %s and %s",
					binary.Op,
					lhs.T,
					rhs.T,
				),
				Help: "Bitwise operators require integers of the same type; use and/or for boolean logic.",
			}
		}
		return nil, false, fmt.Errorf(
			"invalid operand types for %s: %s and %s\n",
			binary.Op,
			lhs.T,
			rhs.T,
		)
	}

	var resultType *ast.ExprType
	if validation.Handler != nil {
		operands := []*ast.ExprType{lhs, rhs}
		resultType, err = validation.Handler(operands)
		if err != nil {
			return nil, false, err
		}
	} else {
		resultType = validation.ResultType
		if resultType == nil {
			resultType = commonType
		}
	}

	if expectedType != nil && !typeAssignable(resultType, expectedType) {
		// If the result is a non-explicit basic type (untyped literal result) and
		// the expected type is also basic, try the same compatibility+coercion path
		// that single literals use in inferBasicExprTypeWithContext. This allows
		// expressions like `a f64 := 1.5 + 0.5` or `a i32 := 1 + 2` to work.
		if resultType.Kind == ast.EXPR_TYPE_BASIC && expectedType.Kind == ast.EXPR_TYPE_BASIC {
			resultBasic := resultType.T.(*ast.BasicType)
			expectedBasic := expectedType.T.(*ast.BasicType)
			if !resultBasic.Explicit && resultBasic.IsCompatibleWith(expectedBasic) {
				resultBasic.Kind = expectedBasic.Kind
				return expectedType, ctx, nil
			}
		}
		return nil, false, diagnostics.Diag{
			Title: "Type mismatch", Message: fmt.Sprintf("expected %s, got %s", expectedType.T, resultType.T),
			Help: "Use numeric values with matching widths and signedness; equal-width C scalar values are compatible with Telia values.",
		}
	}

	return resultType, ctx, nil
}

func (s *sema) ensureBinaryOperatorsAreTheSame(
	binary *ast.BinExpr,
	expected *ast.ExprType,
	referenceScope *ast.Scope,
	declScope *ast.Scope,
	fromImportPackage bool,
) (*ast.ExprType, *ast.ExprType, bool, error) {
	var lhs, rhs *ast.ExprType
	var ctx bool
	var err error

	var lhsHasContext, rhsHasContext bool
	contextualInteger := expected != nil && expected.IsInteger() &&
		(binary.Op.IsBitwiseOp() || binary.Op == token.PLUS || binary.Op == token.MINUS || binary.Op == token.STAR || binary.Op == token.SLASH)
	if contextualInteger {
		lhs, err = s.inferExprTypeWithContext(
			binary.Left,
			expected,
			referenceScope,
			declScope,
			fromImportPackage,
			false,
		)
		lhsHasContext = true
	} else {
		lhs, lhsHasContext, err = s.inferExprTypeWithoutContext(
			binary.Left,
			referenceScope,
			declScope,
			fromImportPackage,
			false,
		)
	}
	if err != nil {
		return nil, nil, ctx, err
	}

	if contextualInteger {
		rhs, err = s.inferExprTypeWithContext(
			binary.Right,
			expected,
			referenceScope,
			declScope,
			fromImportPackage,
			false,
		)
		rhsHasContext = true
	} else {
		rhs, rhsHasContext, err = s.inferExprTypeWithoutContext(
			binary.Right,
			referenceScope,
			declScope,
			fromImportPackage,
			false,
		)
	}
	if err != nil {
		return nil, nil, ctx, err
	}

	if lhsHasContext && !rhsHasContext {
		rhsTypeWithContext, err := s.inferExprTypeWithContext(
			binary.Right,
			lhs,
			referenceScope,
			declScope,
			fromImportPackage,
			false,
		)
		if err != nil {
			return nil, nil, ctx, err
		}
		rhs = rhsTypeWithContext
		ctx = true
	} else if !lhsHasContext && rhsHasContext {
		lhsTypeWithContext, err := s.inferExprTypeWithContext(
			binary.Left,
			rhs,
			referenceScope,
			declScope,
			fromImportPackage,
			false,
		)
		// TODO(errors)
		if err != nil {
			return nil, nil, false, err
		}
		lhs = lhsTypeWithContext
		ctx = true
	}

	if s.isErrorNilComparison(lhs, rhs) {
		return lhs, rhs, ctx, nil
	}
	if lhs.IsNumeric() && rhs.IsNumeric() &&
		lhs.T.(*ast.BasicType).Kind == rhs.T.(*ast.BasicType).Kind {
		return lhs, lhs, ctx, nil
	}

	if !lhs.Equals(rhs) {
		// Allow one side to be an untyped literal whose basic type is compatible
		// with the other side's explicit type — e.g. `f32_var + 1.5` where lhs
		// is F32_TYPE (explicit) and rhs comes from the literal `1.5` (FLOAT_TYPE,
		// non-explicit). Only literals are safe to coerce in-place; coercing a
		// named variable's type would corrupt the declaration's shared type node.
		if lhs.Kind == ast.EXPR_TYPE_BASIC && rhs.Kind == ast.EXPR_TYPE_BASIC {
			lhsBasic := lhs.T.(*ast.BasicType)
			rhsBasic := rhs.T.(*ast.BasicType)
			if lhsBasic.Explicit && !rhsBasic.Explicit &&
				binary.Right.Kind == ast.KIND_LITERAL_EXPR &&
				rhsBasic.IsCompatibleWith(lhsBasic) {
				rhsBasic.Kind = lhsBasic.Kind
				return lhs, lhs, ctx, nil
			}
			if rhsBasic.Explicit && !lhsBasic.Explicit &&
				binary.Left.Kind == ast.KIND_LITERAL_EXPR &&
				lhsBasic.IsCompatibleWith(rhsBasic) {
				lhsBasic.Kind = rhsBasic.Kind
				return rhs, rhs, ctx, nil
			}
		}
		if binary.Op.IsBitwiseOp() {
			return nil, nil, ctx, diagnostics.Diag{
				Title: "Invalid bitwise operands",
				Message: fmt.Sprintf(
					"bitwise operands must have the same integer type, got %s and %s",
					lhs.T,
					rhs.T,
				),
				Help: "Convert the operands explicitly to the same integer type.",
			}
		}
		return nil, nil, ctx, fmt.Errorf("invalid operands types: %s and %s\n", lhs.T, rhs.T)
	}
	return lhs, rhs, ctx, nil
}

// returnsError reports whether retType is an error type directly or is a tuple
// whose last element is an error type.
func returnsError(retType *ast.ExprType) bool {
	if retType.IsError() {
		return true
	}
	if retType.Kind == ast.EXPR_TYPE_TUPLE {
		tt := retType.T.(*ast.TupleType)
		if len(tt.Types) > 0 {
			return tt.Types[len(tt.Types)-1].IsError()
		}
	}
	return false
}

// unwrapErrorFromTuple returns a new ExprType with the last element (error)
// removed from the tuple. If the tuple has only one non-error element after
// unwrapping, it returns a basic type instead of a 1-element tuple.
func unwrapErrorFromTuple(ty *ast.ExprType) *ast.ExprType {
	tt := ty.T.(*ast.TupleType)
	nonErrTypes := tt.Types[:len(tt.Types)-1]
	if len(nonErrTypes) == 1 {
		return nonErrTypes[0]
	}
	return &ast.ExprType{
		Kind: ast.EXPR_TYPE_TUPLE,
		T:    &ast.TupleType{Types: nonErrTypes},
	}
}

// basicKind returns the token kind for a basic type, or a negative value
// if the type is nil, has no inner type, or is not EXPR_TYPE_BASIC.
// For EXPR_TYPE_POINTER, it looks through to the pointee's basic kind
// (used to detect untyped nullptr in error-nil comparisons).
func basicKind(ty *ast.ExprType) token.Kind {
	if ty == nil || ty.T == nil {
		return -1
	}
	if ty.Kind == ast.EXPR_TYPE_BASIC {
		return ty.T.(*ast.BasicType).Kind
	}
	if ty.Kind == ast.EXPR_TYPE_POINTER {
		ptr := ty.T.(*ast.PointerType)
		if ptr.Type != nil && ptr.Type.Kind == ast.EXPR_TYPE_BASIC {
			return ptr.Type.T.(*ast.BasicType).Kind
		}
	}
	return -1
}

// isErrorNilComparison returns true when comparing an error value with nil.
// Handles both orderings: err != nil and nil != err.
func (s *sema) isErrorNilComparison(lhs, rhs *ast.ExprType) bool {
	l := basicKind(lhs)
	r := basicKind(rhs)
	return (l == token.ERROR_TYPE && r == token.UNTYPED_NULLPTR) ||
		(r == token.ERROR_TYPE && l == token.UNTYPED_NULLPTR)
}

func (sema *sema) inferUnaryExprType(
	unary *ast.UnaryExpr,
	expectedType *ast.ExprType,
	referenceScope *ast.Scope,
	declScope *ast.Scope,
	fromImportPackage bool,
) (*ast.ExprType, bool, error) {
	var operandType *ast.ExprType
	var err error
	hasContext := expectedType != nil

	if expectedType != nil {
		operandType, err = sema.inferExprTypeWithContext(
			unary.Value,
			expectedType,
			referenceScope,
			declScope,
			fromImportPackage,
			false,
		)
		if err != nil {
			return nil, false, err
		}
	} else {
		operandType, hasContext, err = sema.inferExprTypeWithoutContext(unary.Value, referenceScope, declScope, fromImportPackage, false)
		if err != nil {
			return nil, false, err
		}
	}

	validation, exists := ast.UnaryOperators[unary.Op]
	if !exists {
		return nil, false, fmt.Errorf("invalid unary operator: %s\n", unary.Op)
	}

	valid := slices.ContainsFunc(
		validation.ValidTypes,
		func(valid *ast.ExprType) bool { return operatorTypeMatches(valid, operandType) },
	)

	if !valid {
		if unary.Op == token.TILDE {
			return nil, false, diagnostics.Diag{
				Title:   "Invalid bitwise operands",
				Message: fmt.Sprintf("bitwise NOT requires an integer, got %s", operandType.T),
				Help:    "Use ~ with an integer operand; use not for boolean negation.",
			}
		}
		return nil, false, fmt.Errorf(
			"invalid operand type %s for operator %s\n",
			operandType.T,
			unary.Op,
		)
	}

	resultType := validation.ResultType
	if resultType == nil {
		resultType = operandType
	}

	if expectedType != nil && resultType.Kind != expectedType.Kind {
		return nil, false, fmt.Errorf(
			"type mismatch: expected %v, got %v\n",
			expectedType,
			resultType,
		)
	}

	unary.Type = resultType
	return resultType, hasContext, nil
}

func (sema *sema) inferFnCallExprTypeWithContext(
	fnCall *ast.FnCall,
	referenceScope *ast.Scope,
	declScope *ast.Scope,
	fromImportPackage bool,
	isArg bool,
) (*ast.ExprType, error) {
	fnRetType, err := sema.checkFnCall(
		fnCall,
		referenceScope,
		declScope,
		nil,
		fromImportPackage,
		isArg,
	)
	return fnRetType, err
}

func (sema *sema) inferVoidExprTypeWithContext(expectedType *ast.ExprType) (*ast.ExprType, error) {
	// TODO(errors)
	if !expectedType.IsVoid() {
		return nil, fmt.Errorf("expected return type to be '%v'", expectedType)
	}
	return expectedType, nil
}

func (sema *sema) inferTupleExprTypeWithContext(
	tuple *ast.TupleExpr,
	expectedType *ast.ExprType,
	referenceScope *ast.Scope,
	declScope *ast.Scope,
	fromImportPackage bool,
) (*ast.ExprType, error) {
	// TODO(errors)
	if expectedType.Kind != ast.EXPR_TYPE_TUPLE {
		return nil, fmt.Errorf("expected to be %v, got tuple\n", expectedType)
	}
	expectedTuple := expectedType.T.(*ast.TupleType)

	// TODO(errors)
	if len(expectedTuple.Types) != len(tuple.Exprs) {
		return nil, fmt.Errorf(
			"expected %d expressions, but got %d\n",
			len(expectedTuple.Types),
			len(tuple.Exprs),
		)
	}

	tupleTy := new(ast.ExprType)
	tupleTy.Kind = ast.EXPR_TYPE_TUPLE

	types := make([]*ast.ExprType, 0)
	for i, expr := range tuple.Exprs {
		ty, err := sema.inferExprTypeWithContext(
			expr,
			expectedTuple.Types[i],
			referenceScope,
			declScope,
			fromImportPackage,
			false,
		)
		if err != nil {
			return nil, err
		}
		types = append(types, ty)
	}

	ty := &ast.TupleType{Types: types}
	tupleTy.T = ty
	tuple.Type = ty
	// Register this tuple shape so codegen has the full set upfront.
	sema.registerTupleType(tupleTy)
	return tupleTy, nil
}

func (sema *sema) inferStructLiteralExprWithContext(
	stLit *ast.StructLiteralExpr,
	expectedType *ast.ExprType,
	referenceScope *ast.Scope,
	declScope *ast.Scope,
	fromImportPackage bool,
) (*ast.ExprType, error) {
	actual, err := sema.literalType(stLit, referenceScope)
	if err != nil {
		return nil, err
	}
	if !actual.Equals(expectedType) {
		return nil, fmt.Errorf(
			"struct type mismatch, expected %s, got %s\n",
			expectedType.T,
			stLit.Name.Name(),
		)
	}
	st := actual.T.(*ast.StructType)
	if st.Decl.IsUnion && len(stLit.Values) != 1 {
		return nil, diagnostics.Diag{
			Title:   "Union member selection",
			Message: "a union literal must initialize exactly one member",
			Span:    stLit.Name.Span,
			Help:    "Write Union{member: value} with one member.",
			Related: []diagnostics.Related{
				{Span: st.Decl.Name.Span, Message: "union declared here"},
			},
		}
	}
	seen := make(map[string]bool)
	for _, value := range stLit.Values {
		name := value.Name.Name()
		if seen[name] {
			return nil, fmt.Errorf("duplicate field %s in struct literal", name)
		}
		seen[name] = true
		var fieldType *ast.ExprType
		for i, field := range st.Decl.Fields {
			if name == field.Name.Name() {
				value.Index = i
				fieldType = field.Type
				break
			}
		}
		if fieldType == nil {
			return nil, fmt.Errorf("field %s does not exist in the struct", name)
		}
		_, err := sema.inferExprTypeWithContext(
			value.Value,
			fieldType,
			referenceScope,
			declScope,
			fromImportPackage,
			false,
		)
		if err != nil {
			return nil, err
		}
	}

	return expectedType, nil
}

func (s *sema) inferStructFieldAccessExprWithContext(
	fieldAccess *ast.FieldAccess,
	expectedType *ast.ExprType,
	referenceScope *ast.Scope,
) (*ast.ExprType, error) {
	_, field, err := s.checkFieldAccessVariable(
		fieldAccess,
		referenceScope,
	)
	if err != nil {
		return nil, err
	}
	if !typeAssignable(field.Type, expectedType) {
		return nil, fmt.Errorf(
			"%s type mismatch on struct field access, expected %s, got %s\n",
			fieldAccess.Left.Name.Pos,
			expectedType.T,
			field.Type.T,
		)
	}
	return expectedType, nil
}

func (s *sema) inferPtrExprWithContext(
	ptr *ast.AddressOfExpr,
	expectedType *ast.ExprType,
	referenceScope *ast.Scope,
	declScope *ast.Scope,
	fromImportPackage bool,
	isArg bool,
) (*ast.ExprType, error) {
	if !expectedType.IsPointer() {
		return nil, fmt.Errorf("expected type needs to be a pointer, not %s\n", expectedType.T)
	}

	expectedPtrType := expectedType.T.(*ast.PointerType)
	// Address-of must retain the storage's actual type. Contextual value
	// conversion here could otherwise turn &f32 into a pointer to C float.
	underlyingPtrExprType, _, err := s.inferExprTypeWithoutContext(
		ptr.Expr,
		referenceScope,
		declScope,
		fromImportPackage,
		isArg,
	)
	if err != nil {
		return nil, err
	}
	if !ptr.Expr.IsAddressable() {
		return nil, diagnostics.Diag{
			Title:   "Not addressable",
			Message: "address-of requires addressable storage",
			Span:    ptr.Expr.Span,
			Help:    "Store the value in a variable before taking its address; constants have no mutable storage.",
		}
	}
	if err := rejectVaListEscape(underlyingPtrExprType, ptr.Expr.Span); err != nil {
		return nil, err
	}

	t := new(ast.ExprType)
	t.Kind = ast.EXPR_TYPE_POINTER
	readOnly := s.readOnlyStorage(ptr.Expr, referenceScope)
	t.T = &ast.PointerType{Type: underlyingPtrExprType, Const: readOnly}
	if readOnly && !expectedPtrType.Const {
		return nil, constMutation(ptr.Expr.Span)
	}

	if !pointerAssignable(t, expectedType) {
		return nil, diagnostics.Diag{
			Title: "Pointer type mismatch", Message: fmt.Sprintf("expected %s, got %s", expectedType.T, t.T), Span: ptr.Expr.Span,
			Help: "Match the storage's exact pointee type and qualifiers; scalar value compatibility does not convert pointers.",
		}
	}
	return t, err
}

func (s *sema) inferNullptrExprTypeWithContext(
	nullptr *ast.NullPtrExpr,
	expectedType *ast.ExprType,
	// referenceScope *ast.Scope,
	// declScope *ast.Scope,
	// fromImportPackage bool,
	// isArg bool,
) (*ast.ExprType, error) {
	isCString := expectedType.Kind == ast.EXPR_TYPE_BASIC &&
		expectedType.T.(*ast.BasicType).Kind == token.CSTRING_TYPE
	if !expectedType.IsPointer() && expectedType.Kind != ast.EXPR_TYPE_FUNCTION &&
		!expectedType.Equals(ast.RAWPTR_TYPE) &&
		!expectedType.IsError() &&
		!isCString {
		return nil, fmt.Errorf("unable to assign nil to non-pointer type")
	}
	nullptr.Type = expectedType
	return nullptr.Type, nil
}

func (s *sema) inferBasicExprTypeWithContext(
	actual *ast.BasicType,
	expected *ast.BasicType,
) (*ast.BasicType, error) {
	if actual.Kind.IsNumeric() && actual.Kind == expected.Kind {
		return expected, nil
	}
	if actual.Kind == token.CSTRING_TYPE && expected.Kind == token.CSTRING_TYPE {
		if actual.Const && !expected.Const {
			return nil, diagnostics.Diag{
				Title:   "Const qualifier removed",
				Message: "cannot convert const C text to writable C text",
				Help:    "Keep c::ConstCString from std::ffi::c, or explicitly copy into owned writable storage.",
			}
		}
		return expected, nil
	}
	if !actual.Explicit || !expected.Explicit {
		if !actual.IsCompatibleWith(expected) {
			return nil, fmt.Errorf("cannot use %s as %s\n", actual, expected)
		}

		// TODO: deal with overflow for numeric types
		actual.Kind = expected.Kind
		return actual, nil
	}
	if !actual.Equals(expected) {
		return nil, fmt.Errorf("type mismatch: expected %s, got %s\n", expected, actual)
	}
	return actual, nil
}

func (sema *sema) inferExprTypeWithoutContext(
	expr *ast.Node,
	referenceScope *ast.Scope,
	declScope *ast.Scope,
	fromImportPackage bool,
	isArg bool,
) (result *ast.ExprType, isLiteral bool, resultErr error) {
	defer func() { resultErr = diagnostics.WithSpan(resultErr, expr.Span) }()
	switch expr.Kind {
	case ast.KIND_MATCH_EXPR:
		ty, err := sema.checkMatch(expr.Node.(*ast.MatchExpr), nil, referenceScope)
		return ty, true, err
	case ast.KIND_LITERAL_EXPR:
		return sema.inferLiteralExprTypeWithoutContext(expr.Node.(*ast.LiteralExpr))
	case ast.KIND_ARRAY_LITERAL_EXPR:
		return nil, false, fmt.Errorf("array literal needs an explicit array type")
	case ast.KIND_INDEX_EXPR:
		ty, err := sema.inferIndex(
			expr.Node.(*ast.IndexExpr),
			referenceScope,
			declScope,
			fromImportPackage,
		)
		return ty, false, err
	case ast.KIND_ID_EXPR:
		return sema.inferIdExprTypeWithoutContext(expr.Node.(*ast.IdExpr), referenceScope)
	case ast.KIND_UNARY_EXPR:
		return sema.inferUnaryExprType(
			expr.Node.(*ast.UnaryExpr),
			nil,
			referenceScope,
			declScope,
			fromImportPackage,
		)
	case ast.KIND_BINARY_EXPR:
		return sema.inferBinaryExprType(
			expr.Node.(*ast.BinExpr),
			nil,
			referenceScope,
			declScope,
			fromImportPackage,
		)
	case ast.KIND_FN_CALL:
		return sema.inferFnCallExprTypeWithoutContext(
			expr.Node.(*ast.FnCall),
			referenceScope,
			declScope,
			fromImportPackage,
			isArg,
		)
	case ast.KIND_NAMESPACE_ACCESS:
		return sema.inferNamespaceAccessExprTypeWithoutContext(
			expr.Node.(*ast.NamespaceAccess),
			referenceScope,
			declScope,
			fromImportPackage,
		)
	case ast.KIND_STRUCT_EXPR:
		return sema.inferStructLiteralExprWithoutContext(
			expr.Node.(*ast.StructLiteralExpr),
			referenceScope,
			declScope,
			fromImportPackage,
		)
	case ast.KIND_FIELD_ACCESS:
		return sema.inferStructFieldAccessExprWithoutContext(
			expr.Node.(*ast.FieldAccess),
			referenceScope,
		)
	case ast.KIND_DEREF_POINTER_EXPR:
		return sema.inferDerefPtrExprTypeWithoutContext(
			expr.Node.(*ast.DerefPointerExpr),
			referenceScope,
			declScope,
			fromImportPackage,
			isArg,
		)
	case ast.KIND_ADDRESS_OF_EXPR:
		return sema.inferAddressOfExprWithoutContext(
			expr.Node.(*ast.AddressOfExpr),
			referenceScope,
			declScope,
			fromImportPackage,
			isArg,
		)
	case ast.KIND_NULLPTR_EXPR:
		nullptr := &ast.ExprType{
			Kind: ast.EXPR_TYPE_POINTER,
			T: &ast.PointerType{
				Type: &ast.ExprType{
					Kind: ast.EXPR_TYPE_BASIC,
					T:    &ast.BasicType{Kind: token.UNTYPED_NULLPTR},
				},
			},
		}
		return nullptr, false, nil
	default:
		log.Fatalf("unimplemented expression: %s\n", expr.Node)
		return nil, false, nil
	}
}

func (sema *sema) inferLiteralExprTypeWithoutContext(
	literal *ast.LiteralExpr,
) (*ast.ExprType, bool, error) {
	// TODO(errors)
	if literal.Type.Kind != ast.EXPR_TYPE_BASIC {
		return nil, false, fmt.Errorf("type mismatch for literal expression without context")
	}

	actualBasicType := literal.Type.T.(*ast.BasicType)
	switch actualBasicType.Kind {
	case token.UNTYPED_NULLPTR:
		nullptr := &ast.ExprType{
			Kind: ast.EXPR_TYPE_POINTER,
			T: &ast.PointerType{
				Type: &ast.ExprType{
					Kind: ast.EXPR_TYPE_BASIC,
					T:    &ast.BasicType{Kind: token.UNTYPED_NULLPTR},
				},
			},
		}
		literal.Type = nullptr
		return literal.Type, false, nil
	case token.ERROR_TYPE:
		basic := literal.Type.T.(*ast.BasicType)
		basic.Explicit = true
		return literal.Type, true, nil
	}
	return literal.Type, false, nil
}

func (sema *sema) inferIdExprTypeWithoutContext(
	id *ast.IdExpr,
	scope *ast.Scope,
) (*ast.ExprType, bool, error) {
	variableName := id.Name.Name()
	variable, err := scope.LookupAcrossScopes(variableName)
	// TODO(errors)
	if err != nil {
		return nil, false, diagnostics.Diag{
			Title:   "Undeclared name",
			Message: fmt.Sprintf("symbol not found on scope: '%s' is not declared.", variableName),
			Span:    id.Name.Span,
		}
	}

	id.N = variable
	if variable.Invalid {
		return nil, false, diagnostics.COMPILER_ERROR_FOUND
	}

	switch variable.Kind {
	case ast.KIND_VA_INTRINSIC:
		return nil, false, vaIntrinsicValueDiagnostic(id.Name.Span)
	case ast.KIND_FN_DECL, ast.KIND_PROTO:
		ty, err := sema.functionValue(id, variable)
		return ty, true, err
	case ast.KIND_CONST_DECL:
		d := variable.Node.(*ast.ConstDecl)
		ty, err := sema.inferConstReference(id, d, nil)
		return ty, d.Type != nil && d.Type.Explicit, err
	case ast.KIND_VAR_ID_STMT:
		ty := variable.Node.(*ast.VarIdStmt).Type
		return ty, ty.Explicit, nil
	case ast.KIND_PARAM:
		param := variable.Node.(*ast.Param)
		if err := rejectCVariadicValue(param, id.Name.Span); err != nil {
			return nil, false, err
		}
		return param.Type, true, nil
	default:
		// TODO(errors)
		return nil, false, fmt.Errorf("'%s' is not a variable or parameter", variableName)
	}
}

func (sema *sema) inferFnCallExprTypeWithoutContext(
	fnCall *ast.FnCall,
	referenceScope *ast.Scope,
	declScope *ast.Scope,
	fromImportPackage bool,
	isArg bool,
) (*ast.ExprType, bool, error) {
	fnRetType, err := sema.checkFnCall(
		fnCall,
		referenceScope,
		declScope,
		nil,
		fromImportPackage,
		isArg,
	)
	return fnRetType, true, err
}

func (sema *sema) inferNamespaceAccessExprTypeWithoutContext(
	namespaceAccess *ast.NamespaceAccess,
	referenceScope *ast.Scope,
	declScope *ast.Scope,
	fromImportPackage bool,
) (*ast.ExprType, bool, error) {
	ty, err := sema.checkNamespaceAccess(
		namespaceAccess,
		referenceScope,
		declScope,
		fromImportPackage,
	)
	if err == nil && namespaceAccess.IsImport && namespaceAccess.Right.Kind == ast.KIND_ID_EXPR {
		return ty, ty.Explicit, nil
	}
	return ty, true, err
}

func (sema *sema) inferStructLiteralExprWithoutContext(
	stLit *ast.StructLiteralExpr,
	referenceScope *ast.Scope,
	declScope *ast.Scope,
	fromImportPackage bool,
) (*ast.ExprType, bool, error) {
	ty, err := sema.literalType(stLit, declScope)
	if err != nil {
		return nil, false, err
	}
	inferredTy, err := sema.inferStructLiteralExprWithContext(
		stLit,
		ty,
		referenceScope,
		declScope,
		fromImportPackage,
	)
	return inferredTy, true, err
}

func (sema *sema) inferStructFieldAccessExprWithoutContext(
	fieldAccess *ast.FieldAccess,
	referenceScope *ast.Scope,
) (*ast.ExprType, bool, error) {
	_, field, err := sema.checkFieldAccessVariable(
		fieldAccess,
		referenceScope,
	)
	if err != nil {
		return nil, false, err
	}
	return field.Type, true, err
}

func (sema *sema) inferAddressOfExprWithoutContext(
	ptr *ast.AddressOfExpr,
	referenceScope, declScope *ast.Scope,
	fromImportPackage bool,
	isArg bool,
) (*ast.ExprType, bool, error) {
	ty, _, err := sema.inferExprTypeWithoutContext(
		ptr.Expr,
		referenceScope,
		declScope,
		fromImportPackage,
		isArg,
	)
	if err != nil {
		return nil, false, err
	}
	if err := rejectVaListEscape(ty, ptr.Expr.Span); err != nil {
		return nil, false, err
	}
	if !ptr.Expr.IsAddressable() {
		return nil, false, diagnostics.Diag{
			Title:   "Not addressable",
			Message: "address-of requires addressable storage",
			Span:    ptr.Expr.Span,
			Help:    "Store the value in a variable before taking its address; constants have no mutable storage.",
		}
	}
	result := ast.NewPointerType(ty)
	if err := checkForeignPointerDeclarator(result, ptr.Expr.Span); err != nil &&
		ty.Kind == ast.EXPR_TYPE_ARRAY &&
		ty.T.(*ast.SequenceType).Foreign {
		return nil, false, err
	}
	result.T.(*ast.PointerType).Const = sema.readOnlyStorage(ptr.Expr, referenceScope)
	result.Explicit = true
	return result, true, nil
}

func (sema *sema) checkNamespaceAccess(
	namespaceAccess *ast.NamespaceAccess,
	referenceScope *ast.Scope,
	declScope *ast.Scope,
	fromImportPackage bool,
) (*ast.ExprType, error) {
	if namespaceAccess.IsImport {
		imp := namespaceAccess.Left.Name.Name()
		useDecl, ok := sema.file.Imports[imp]
		if !ok {
			panic("import not found")
		}
		return sema.checkImportAccess(
			namespaceAccess.Right,
			referenceScope,
			useDecl.Package.Scope,
			fromImportPackage,
		)
	}

	var left *ast.Node
	var err error

	if fromImportPackage {
		left, err = declScope.LookupAcrossScopes(namespaceAccess.Left.Name.Name())
		if err != nil {
			return nil, err
		}
	} else {
		left, err = referenceScope.LookupAcrossScopes(namespaceAccess.Left.Name.Name())
		if err != nil {
			return nil, err
		}
	}

	namespaceAccess.Left.N = left
	if left.Invalid {
		return nil, diagnostics.COMPILER_ERROR_FOUND
	}
	if d, ok := left.Node.(*ast.EnumDecl); ok {
		var err error
		d, err = sema.enumInstance(
			d,
			namespaceAccess.TypeArgs,
			referenceScope,
			namespaceAccess.Left.Name,
		)
		if err != nil {
			return nil, err
		}
		return sema.checkEnumAccess(namespaceAccess, d, referenceScope)
	}

	// TODO(errors)
	if left.Kind != ast.KIND_EXTERN_DECL {
		return nil, fmt.Errorf("invalid access")
	}

	externDecl := left.Node.(*ast.ExternDecl)
	if namespaceAccess.Right.Kind == ast.KIND_NAMESPACE_ACCESS {
		return sema.checkNamespaceAccess(
			namespaceAccess.Right.Node.(*ast.NamespaceAccess),
			referenceScope,
			externDecl.Scope,
			true,
		)
	}
	if namespaceAccess.Right.Kind == ast.KIND_ID_EXPR {
		return sema.inferImportedValue(
			namespaceAccess.Right.Node.(*ast.IdExpr),
			externDecl.Scope,
			nil,
		)
	}
	if namespaceAccess.Right.Kind != ast.KIND_FN_CALL {
		return nil, fmt.Errorf("expected prototype call")
	}
	protoCall, err := sema.checkPrototypeCall(
		externDecl,
		namespaceAccess.Right.Node.(*ast.FnCall),
		referenceScope,
		declScope,
		fromImportPackage,
	)
	if err != nil {
		return nil, err
	}
	return protoCall.RetType, nil
}

func (sema *sema) checkFieldAccessVariable(
	fieldAccess *ast.FieldAccess,
	referenceScope *ast.Scope,
) (*ast.IdExpr, *ast.StructField, error) {
	id, field, err := sema.getAccessedField(fieldAccess, referenceScope)
	return id, field, err
}

func (sema *sema) checkImportAccess(
	node *ast.Node,
	referenceScope *ast.Scope,
	declScope *ast.Scope,
	fromImportPackage bool,
) (*ast.ExprType, error) {
	switch node.Kind {
	case ast.KIND_ID_EXPR:
		return sema.inferImportedValue(node.Node.(*ast.IdExpr), declScope, nil)
	case ast.KIND_FN_CALL:
		fnCall := node.Node.(*ast.FnCall)
		fnRetType, err := sema.checkFnCall(
			fnCall,
			referenceScope,
			declScope,
			nil,
			fromImportPackage,
			false,
		)
		return fnRetType, err
	case ast.KIND_NAMESPACE_ACCESS:
		namespaceAccess := node.Node.(*ast.NamespaceAccess)
		ty, err := sema.checkNamespaceAccess(namespaceAccess, referenceScope, declScope, true)
		return ty, err
	case ast.KIND_EXTERN_DECL:
		return nil, fmt.Errorf("nothing do with a extern declaration, try accessing a prototype")
	case ast.KIND_STRUCT_EXPR:
		stLit := node.Node.(*ast.StructLiteralExpr)
		ty, _, err := sema.inferStructLiteralExprWithoutContext(
			stLit,
			referenceScope,
			declScope,
			fromImportPackage,
		)
		return ty, err
	default:
		return nil, fmt.Errorf("unimplemented symbol to import: %v\n", node.Kind)
	}
}

func (sema *sema) checkPrototypeCall(
	extern *ast.ExternDecl,
	prototypeCall *ast.FnCall,
	referenceScope *ast.Scope,
	declScope *ast.Scope,
	fromImportPackage bool,
) (*ast.Proto, error) {
	if len(prototypeCall.TypeArgs) != 0 {
		return nil, fmt.Errorf(
			"%s: '%s' does not take type arguments",
			prototypeCall.Name.Pos,
			prototypeCall.Name.Name(),
		)
	}
	symbol, err := extern.Scope.LookupCurrentScope(prototypeCall.Name.Name())
	if err != nil {
		pos := prototypeCall.Name.Pos
		prototypeNotFound := diagnostics.Diag{
			Span: prototypeCall.Name.Span,
			Message: fmt.Sprintf(
				"%s:%d:%d: function '%s' not declared on extern '%s'",
				pos.Filename,
				pos.Line,
				pos.Column,
				prototypeCall.Name.Name(),
				extern.Name.Name(),
			),
		}
		sema.collector.ReportAndSave(prototypeNotFound)
		return nil, diagnostics.COMPILER_ERROR_FOUND
	}
	if symbol.Kind != ast.KIND_PROTO {
		return nil, fmt.Errorf("expected prototype function, not %v\n", symbol.Kind)
	}

	prototype := symbol.Node.(*ast.Proto)
	prototypeCall.Proto = prototype
	err = sema.checkFnCallArgs(
		prototypeCall,
		prototype.Params,
		referenceScope,
		declScope,
		fromImportPackage,
	)
	return prototype, err
}

func (sema *sema) checkForLoop(
	forLoop *ast.ForLoop,
	returnTy *ast.ExprType,
	referenceScope *ast.Scope,
	declScope *ast.Scope,
	fromImportPackage, isArg bool,
) error {
	if forLoop.Range != nil {
		ty, _, err := sema.inferExprTypeWithoutContext(
			forLoop.Range,
			referenceScope,
			declScope,
			fromImportPackage,
			false,
		)
		if err != nil {
			return err
		}
		if ty.Kind != ast.EXPR_TYPE_ARRAY && ty.Kind != ast.EXPR_TYPE_SLICE {
			return fmt.Errorf("iteration requires an array or slice")
		}
		forLoop.RangeType = ty
		forLoop.Binding.Type = ty.T.(*ast.SequenceType).Element
		binding := &ast.Node{Kind: ast.KIND_VAR_ID_STMT, Node: forLoop.Binding}
		forLoop.Binding.N = binding
		if err := forLoop.Scope.Insert(forLoop.Binding.Name.Name(), binding); err != nil {
			return err
		}
		return sema.checkBlock(
			forLoop.Block,
			returnTy,
			forLoop.Scope,
			declScope,
			fromImportPackage,
			isArg,
		)
	}
	err := sema.checkStmt(
		forLoop.Init,
		referenceScope,
		declScope,
		returnTy,
		fromImportPackage,
		isArg,
	)
	if err != nil {
		return err
	}

	err = sema.checkIfExpr(forLoop.Cond, referenceScope, declScope, fromImportPackage)
	if err != nil {
		return err
	}

	err = sema.checkStmt(
		forLoop.Update,
		referenceScope,
		declScope,
		returnTy,
		fromImportPackage,
		isArg,
	)
	if err != nil {
		return err
	}

	err = sema.checkBlock(
		forLoop.Block,
		returnTy,
		referenceScope,
		declScope,
		fromImportPackage,
		isArg,
	)
	return err
}

func (sema *sema) checkWhileLoop(
	whileLoop *ast.WhileLoop,
	referenceScope *ast.Scope,
	declScope *ast.Scope,
	returnTy *ast.ExprType,
	fromImportPackage, isArg bool,
) error {
	err := sema.checkIfExpr(whileLoop.Cond, referenceScope, declScope, fromImportPackage)
	if err != nil {
		return err
	}
	err = sema.checkBlock(
		whileLoop.Block,
		returnTy,
		referenceScope,
		declScope,
		fromImportPackage,
		isArg,
	)
	return err
}

func (s *sema) checkAssignment(
	assignment *ast.AssignmentStmt,
	referenceScope *ast.Scope,
	declScope *ast.Scope,
	fromImportPackage, isArg bool,
) error {
	for _, target := range assignment.Targets {
		// TODO(errors)
		if !s.isAssignable(target) {
			return fmt.Errorf("target of type %v is not assignable", target.Kind)
		}

		var err error
		switch target.Kind {
		case ast.KIND_INDEX_EXPR:
			if assignment.Decl {
				return fmt.Errorf("index expressions cannot declare variables")
			}
			_, err = s.inferIndex(
				target.Node.(*ast.IndexExpr),
				referenceScope,
				declScope,
				fromImportPackage,
			)
		case ast.KIND_VAR_ID_STMT:
			err = s.checkNormalVarDecl(
				target.Node.(*ast.VarIdStmt),
				assignment.Decl,
				referenceScope,
			)
		case ast.KIND_FIELD_ACCESS:
			_, _, err = s.checkFieldAccessVariable(
				target.Node.(*ast.FieldAccess),
				referenceScope,
			)
		case ast.KIND_DEREF_POINTER_EXPR:
			// TODO(errors)
			if assignment.Decl {
				return fmt.Errorf(
					"not allowed to create new variable with pointer dereference operator",
				)
			}
			innerDeref := s.getBehindDeref(target.Node.(*ast.DerefPointerExpr))

			switch innerDeref.Kind {
			case ast.KIND_ID_EXPR:
				id := innerDeref.Node.(*ast.IdExpr)

				symbol, err := referenceScope.LookupAcrossScopes(id.Name.Name())
				if err != nil {
					return err
				}

				if symbol.Invalid {
					return diagnostics.COMPILER_ERROR_FOUND
				}
			default:
				// TODO(errors)
				return fmt.Errorf(
					"not allowed to use dereference operator for expresion of type %v",
					innerDeref.Kind,
				)
			}
		default:
			panic(fmt.Sprintf("unimplemented %d", target.Kind))
		}

		if err != nil {
			return err
		}
	}

	t := 0
	for _, expr := range assignment.Values {
		switch expr.Kind {
		case ast.KIND_FN_CALL:
			fnCall := expr.Node.(*ast.FnCall)
			fnRetType, err := s.checkFnCall(
				fnCall,
				referenceScope,
				declScope,
				nil,
				fromImportPackage,
				false,
			)
			if err != nil {
				return err
			}

			if fnRetType.Kind == ast.EXPR_TYPE_TUPLE {
				tupleType := fnRetType.T.(*ast.TupleType)
				affectedVariables := assignment.Targets[t : t+len(tupleType.Types)]
				if err := s.checkTupleTypeAssignedToVariable(
					affectedVariables,
					tupleType,
					referenceScope,
				); err != nil {
					return err
				}
				t += len(affectedVariables)
			} else {
				err := s.checkVarExpr(assignment.Targets[t], expr, referenceScope, declScope, fromImportPackage, isArg)
				if err != nil {
					return err
				}
				t++
			}
		default:
			err := s.checkVarExpr(
				assignment.Targets[t],
				expr,
				referenceScope,
				declScope,
				fromImportPackage,
				isArg,
			)
			if err != nil {
				return err
			}

			t++
		}
	}

	return nil
}

func (s *sema) getBehindDeref(deref *ast.DerefPointerExpr) *ast.Node {
	e := deref.Expr
	for e.Kind == ast.KIND_DEREF_POINTER_EXPR {
		innerDeref := e.Node.(*ast.DerefPointerExpr)
		e = innerDeref.Expr
	}
	return e
}

func (s *sema) isAssignable(lhs *ast.Node) bool {
	switch lhs.Kind {
	case ast.KIND_VAR_ID_STMT,
		ast.KIND_FIELD_ACCESS,
		ast.KIND_DEREF_POINTER_EXPR,
		ast.KIND_INDEX_EXPR:
		return true
	default:
		return false
	}
}
