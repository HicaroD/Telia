package sema

import (
	"crypto/sha256"
	"fmt"
	"strings"

	"github.com/HicaroD/Telia/internal/ast"
	"github.com/HicaroD/Telia/internal/diagnostics"
	"github.com/HicaroD/Telia/internal/lexer/token"
)

type genericState struct {
	enums     map[string]*ast.EnumDecl
	functions map[string]*ast.FnDecl
	structs   map[string]*ast.StructDecl
	resolved  map[*ast.StructDecl]bool
	origins   map[*ast.StructDecl][]string
	bindings  map[string]*ast.ExprType
	chain     []string
	count     int
}

func (s *sema) resolveSignature(fn *ast.FnDecl) error {
	return s.inContext(fn.Owner, fn.File, nil, func() error {
		for _, param := range fn.Params.Fields {
			ty, err := s.resolveType(param.Type, fn.Scope.Parent)
			if err != nil {
				return err
			}
			param.Type = ty
			if err := rejectOpaqueValue(ty, param.Name.Span); err != nil {
				return err
			}
		}
		ty, err := s.resolveType(fn.RetType, fn.Scope.Parent)
		if err != nil {
			return err
		}
		fn.RetType = ty
		if err := rejectOpaqueValue(ty, fn.Name.Span); err != nil {
			return err
		}
		return nil
	})
}

func (s *sema) instantiateFunction(template *ast.FnDecl, args []*ast.ExprType, scope *ast.Scope, site *token.Token) (*ast.FnDecl, error) {
	if len(args) != len(template.TypeParams) {
		return nil, fmt.Errorf("%s: '%s' requires %d explicit type arguments, got %d", site.Pos, template.Name.Name(), len(template.TypeParams), len(args))
	}
	bindings, suffix, err := s.concreteArguments(template.TypeParams, args, scope, map[*ast.TypeAlias]bool{})
	if err != nil {
		return nil, err
	}
	key := declarationKey(template.Owner, template.Name.Name()) + suffix
	if cached := s.generics.functions[key]; cached != nil {
		return cached, nil
	}
	if err := s.beginInstance(template.Name.Name()+suffix, site); err != nil {
		return nil, err
	}
	instance := cloneFunction(template, bindings)
	instance.CName = concreteName("fn", key)
	s.generics.functions[key] = instance
	s.program.Functions = append(s.program.Functions, instance)
	err = s.inContext(template.Owner, template.File, bindings, func() error {
		return s.checkFnDecl(instance, template.Owner.Scope, false, false)
	})
	if err != nil {
		err = fmt.Errorf("%s: in generic function '%s': %w", template.Name.Pos, template.Name.Name(), err)
	}
	return instance, s.endInstance(err)
}

func (s *sema) prepareGenerics(pkg *ast.Package, seen map[*ast.Package]bool) {
	if pkg == nil || seen[pkg] {
		return
	}
	seen[pkg] = true
	for _, file := range pkg.Files {
		for _, imp := range file.Imports {
			s.prepareGenerics(imp.Package, seen)
		}
		for _, node := range file.Body {
			switch d := node.Node.(type) {
			case *ast.StructDecl:
				d.Owner, d.File = pkg, file
				if len(d.TypeParams) == 0 {
					d.CName = concreteName("struct", declarationKey(pkg, d.Name.Name()))
				}
			case *ast.FnDecl:
				d.Owner, d.File = pkg, file
			case *ast.TypeAlias:
				d.Owner, d.File = pkg, file
			}
		}
	}
}

func declarationKey(pkg *ast.Package, name string) string {
	if pkg == nil {
		return name
	}
	return pkg.Loc.Path + "::" + pkg.Loc.Name + "::" + name
}

func concreteName(kind, key string) string {
	return fmt.Sprintf("telia_%s_%x", kind, sha256.Sum256([]byte(key)))
}

func typeKey(ty *ast.ExprType) string {
	switch ty.Kind {
	case ast.EXPR_TYPE_OPAQUE:
		d := ty.T.(*ast.OpaqueType).Decl
		path := d.Name.Pos.Filename
		if d.Name.Span.Valid() {
			path = d.Name.Span.File.Path
		}
		return fmt.Sprintf("opaque:%s:%d", path, d.Name.Span.Start)
	case ast.EXPR_TYPE_ENUM:
		d := ty.T.(*ast.EnumType).Decl
		if d.CName != "" {
			return d.CName
		}
		return declarationKey(d.Owner, d.Name.Name())
	case ast.EXPR_TYPE_ARRAY, ast.EXPR_TYPE_SLICE:
		t := ty.T.(*ast.SequenceType)
		return fmt.Sprintf("sequence%d:%d:%s", ty.Kind, t.Length, typeKey(t.Element))
	case ast.EXPR_TYPE_BASIC:
		return ty.T.(*ast.BasicType).Kind.String()
	case ast.EXPR_TYPE_POINTER:
		return "*" + typeKey(ty.T.(*ast.PointerType).Type)
	case ast.EXPR_TYPE_STRUCT:
		d := ty.T.(*ast.StructType).Decl
		if d.CName != "" {
			return d.CName
		}
		return declarationKey(d.Owner, d.Name.Name())
	case ast.EXPR_TYPE_TUPLE:
		var keys []string
		for _, t := range ty.T.(*ast.TupleType).Types {
			keys = append(keys, typeKey(t))
		}
		return "(" + strings.Join(keys, ",") + ")"
	default:
		panic("nonconcrete type in specialization identity")
	}
}

func (s *sema) inContext(pkg *ast.Package, file *ast.File, bindings map[string]*ast.ExprType, work func() error) error {
	oldPkg, oldFile, oldBindings := s.pkg, s.file, s.generics.bindings
	s.pkg, s.file, s.generics.bindings = pkg, file, bindings
	defer func() { s.pkg, s.file, s.generics.bindings = oldPkg, oldFile, oldBindings }()
	return work()
}

func (s *sema) resolveType(ty *ast.ExprType, scope *ast.Scope) (*ast.ExprType, error) {
	return s.resolveTypeAliases(ty, scope, map[*ast.TypeAlias]bool{})
}

func (s *sema) resolveTypeAliases(ty *ast.ExprType, scope *ast.Scope, aliases map[*ast.TypeAlias]bool) (*ast.ExprType, error) {
	if ty == nil {
		return nil, nil
	}
	switch ty.Kind {
	case ast.EXPR_TYPE_ARRAY, ast.EXPR_TYPE_SLICE:
		t := ty.T.(*ast.SequenceType)
		element, err := s.resolveTypeAliases(t.Element, scope, aliases)
		if err != nil {
			return nil, err
		}
		if !valueType(element) {
			return nil, fmt.Errorf("collection element requires a value type")
		}
		resolved := &ast.ExprType{Kind: ty.Kind, T: &ast.SequenceType{Element: element, Length: t.Length}}
		s.program.SequenceTypes = append(s.program.SequenceTypes, resolved)
		return resolved, nil
	case ast.EXPR_TYPE_ID:
		id := ty.T.(*ast.IdType)
		argumentScope := scope
		if id.Namespace == nil {
			if bound := s.generics.bindings[id.Name.Name()]; bound != nil {
				if len(id.Args) != 0 {
					return nil, fmt.Errorf("%s: type parameter '%s' cannot take type arguments", id.Name.Pos, id.Name.Name())
				}
				return bound, nil
			}
		} else {
			imp, ok := s.file.Imports[id.Namespace.Name()]
			if !ok {
				namespace, err := scope.LookupAcrossScopes(id.Namespace.Name())
				if err != nil || namespace.Kind != ast.KIND_EXTERN_DECL {
					return nil, diagnostics.Diag{Title: "Unknown type namespace", Message: "type namespace is not an imported package or extern block", Span: id.Namespace.Span, Help: "Use the declaring package or named extern block."}
				}
				scope = namespace.Node.(*ast.ExternDecl).Scope
			} else {
				scope = imp.Package.Scope
			}
		}
		sym, err := scope.LookupAcrossScopes(id.Name.Name())
		if err != nil {
			return nil, diagnostics.Diag{Title: "Unknown type", Message: fmt.Sprintf("'%s' type not found on scope", id.Name.Name()), Span: id.Name.Span, Help: "Declare the type or use the correct imported package qualifier."}
		}
		if sym.Invalid {
			return nil, diagnostics.COMPILER_ERROR_FOUND
		}
		switch d := sym.Node.(type) {
		case *ast.OpaqueDecl:
			if len(id.Args) != 0 {
				return nil, diagnostics.Diag{Title: "Invalid opaque type arguments", Message: "header-owned opaque types do not take type arguments", Span: id.Name.Span, Help: "Remove the type arguments.", Related: []diagnostics.Related{{Span: d.Name.Span, Message: "opaque type declared here"}}}
			}
			return &ast.ExprType{Kind: ast.EXPR_TYPE_OPAQUE, T: &ast.OpaqueType{Decl: d}}, nil
		case *ast.EnumDecl:
			d, err = s.enumInstance(d, id.Args, argumentScope, id.Name)
			if err != nil {
				return nil, err
			}
			if err := s.resolveEnum(d); err != nil {
				return nil, err
			}
			return &ast.ExprType{Kind: ast.EXPR_TYPE_ENUM, T: &ast.EnumType{Decl: d}}, nil
		case *ast.TypeAlias:
			if len(id.Args) != 0 {
				return nil, fmt.Errorf("%s: alias '%s' does not take type arguments", id.Name.Pos, id.Name.Name())
			}
			if aliases[d] {
				return nil, fmt.Errorf("%s: cyclic type alias '%s'", id.Name.Pos, id.Name.Name())
			}
			aliases[d] = true
			defer delete(aliases, d)
			var resolved *ast.ExprType
			err := s.inContext(d.Owner, d.File, nil, func() error {
				var err error
				resolved, err = s.resolveTypeAliases(d.Type, d.Owner.Scope, aliases)
				return err
			})
			return resolved, err
		case *ast.StructDecl:
			if len(d.TypeParams) != len(id.Args) {
				return nil, fmt.Errorf("%s: '%s' requires %d type arguments, got %d", id.Name.Pos, id.Name.Name(), len(d.TypeParams), len(id.Args))
			}
			if len(d.TypeParams) != 0 {
				d, err = s.instantiateStruct(d, id.Args, argumentScope, id.Name, aliases)
			} else {
				err = s.resolveStruct(d, nil)
			}
			if err != nil {
				return nil, err
			}
			return &ast.ExprType{Kind: ast.EXPR_TYPE_STRUCT, T: &ast.StructType{Decl: d}}, nil
		default:
			return nil, fmt.Errorf("%s: '%s' is not a type", id.Name.Pos, id.Name.Name())
		}
	case ast.EXPR_TYPE_POINTER:
		inner, err := s.resolveTypeAliases(ty.T.(*ast.PointerType).Type, scope, aliases)
		if err != nil {
			return nil, err
		}
		return ast.NewPointerType(inner), nil
	case ast.EXPR_TYPE_TUPLE:
		result := &ast.TupleType{}
		for _, t := range ty.T.(*ast.TupleType).Types {
			r, err := s.resolveTypeAliases(t, scope, aliases)
			if err != nil {
				return nil, err
			}
			result.Types = append(result.Types, r)
		}
		resolved := &ast.ExprType{Kind: ast.EXPR_TYPE_TUPLE, T: result}
		s.registerTupleType(resolved)
		return resolved, nil
	default:
		return ty, nil
	}
}

func (s *sema) concreteArguments(params []*token.Token, args []*ast.ExprType, scope *ast.Scope, aliases map[*ast.TypeAlias]bool) (map[string]*ast.ExprType, string, error) {
	bindings := map[string]*ast.ExprType{}
	var keys []string
	for i, arg := range args {
		ty, err := s.resolveTypeAliases(arg, scope, aliases)
		if err != nil {
			return nil, "", err
		}
		if !valueType(ty) {
			return nil, "", fmt.Errorf("invalid generic type argument %s: requires a concrete value type", ty.T)
		}
		bindings[params[i].Name()] = ty
		keys = append(keys, typeKey(ty))
	}
	return bindings, "<" + strings.Join(keys, ",") + ">", nil
}

func valueType(ty *ast.ExprType) bool {
	if ty.Kind == ast.EXPR_TYPE_OPAQUE {
		return false
	}
	if ty.IsVoid() {
		return false
	}
	if ty.Kind == ast.EXPR_TYPE_BASIC {
		kind := ty.T.(*ast.BasicType).Kind
		return kind <= token.UNTYPED_START || kind >= token.UNTYPED_END
	}
	if ty.Kind == ast.EXPR_TYPE_TUPLE {
		for _, t := range ty.T.(*ast.TupleType).Types {
			if !valueType(t) {
				return false
			}
		}
	}
	return true
}

func (s *sema) beginInstance(label string, site *token.Token) error {
	if len(s.generics.chain) >= 128 {
		return fmt.Errorf("%s: instantiation depth limit (128) exceeded: %s", site.Pos, label)
	}
	if s.generics.count >= 4096 {
		return fmt.Errorf("%s: specialization limit (4096) exceeded: %s", site.Pos, label)
	}
	s.generics.count++
	s.generics.chain = append(s.generics.chain, fmt.Sprintf("%s at %s", label, site.Pos))
	return nil
}

func (s *sema) endInstance(err error) error {
	if err != nil {
		err = fmt.Errorf("%w\n  instantiated %s", err, s.generics.chain[len(s.generics.chain)-1])
	}
	s.generics.chain = s.generics.chain[:len(s.generics.chain)-1]
	return err
}

func (s *sema) instantiateStruct(template *ast.StructDecl, args []*ast.ExprType, scope *ast.Scope, site *token.Token, aliases map[*ast.TypeAlias]bool) (*ast.StructDecl, error) {
	bindings, suffix, err := s.concreteArguments(template.TypeParams, args, scope, aliases)
	if err != nil {
		return nil, err
	}
	key := declarationKey(template.Owner, template.Name.Name()) + suffix
	if cached := s.generics.structs[key]; cached != nil {
		return cached, nil
	}
	if err := s.beginInstance(template.Name.Name()+suffix, site); err != nil {
		return nil, err
	}
	instance := &ast.StructDecl{Name: template.Name, Owner: template.Owner, File: template.File, CName: concreteName("struct", key)}
	for _, field := range template.Fields {
		copy := *field
		instance.Fields = append(instance.Fields, &copy)
	}
	s.generics.structs[key] = instance
	s.generics.origins[instance] = append([]string(nil), s.generics.chain...)
	err = s.resolveStruct(instance, bindings)
	if err != nil {
		err = fmt.Errorf("%s: in generic struct '%s': %w", template.Name.Pos, template.Name.Name(), err)
	}
	return instance, s.endInstance(err)
}

func (s *sema) resolveStruct(d *ast.StructDecl, bindings map[string]*ast.ExprType) error {
	if s.generics.resolved[d] {
		return nil
	}
	s.generics.resolved[d] = true
	s.program.Structs = append(s.program.Structs, d)
	return s.inContext(d.Owner, d.File, bindings, func() error {
		if err := s.checkStructDecl(d); err != nil {
			return err
		}
		for _, field := range d.Fields {
			ty, err := s.resolveType(field.Type, d.Owner.Scope)
			if err != nil {
				return err
			}
			if !valueType(ty) {
				return fmt.Errorf("%s: field requires a value type", field.Name.Pos)
			}
			field.Type = ty
		}
		return nil
	})
}

func (s *sema) literalType(lit *ast.StructLiteralExpr, scope *ast.Scope) (*ast.ExprType, error) {
	if lit.Decl != nil {
		return &ast.ExprType{Kind: ast.EXPR_TYPE_STRUCT, T: &ast.StructType{Decl: lit.Decl}}, nil
	}
	ty, err := s.resolveType(&ast.ExprType{Kind: ast.EXPR_TYPE_ID, T: &ast.IdType{Name: lit.Name, Namespace: lit.Namespace, Args: lit.TypeArgs}}, scope)
	if err != nil {
		return nil, err
	}
	if ty.Kind != ast.EXPR_TYPE_STRUCT {
		if err := rejectOpaqueValue(ty, lit.Name.Span); err != nil {
			return nil, err
		}
		return nil, fmt.Errorf("%s: struct literal requires a struct type", lit.Name.Pos)
	}
	lit.Decl = ty.T.(*ast.StructType).Decl
	return ty, nil
}

// Pointers need only a forward declaration; by-value edges require a complete
// layout and must form a DAG. Use the same order to drive C emission.
func (s *sema) orderValueTypes() error {
	state := map[string]int{}
	var visit func(*ast.ExprType) error
	visit = func(ty *ast.ExprType) error {
		if ty.Kind != ast.EXPR_TYPE_ENUM && ty.Kind != ast.EXPR_TYPE_STRUCT && ty.Kind != ast.EXPR_TYPE_TUPLE && ty.Kind != ast.EXPR_TYPE_ARRAY && ty.Kind != ast.EXPR_TYPE_SLICE {
			return nil
		}
		key := typeKey(ty)
		if state[key] == 2 {
			return nil
		}
		if state[key] == 1 {
			return fmt.Errorf("infinitely sized by-value layout involving %s", ty.T)
		}
		state[key] = 1
		if ty.Kind == ast.EXPR_TYPE_ENUM {
			for _, v := range ty.T.(*ast.EnumType).Decl.Variants {
				for _, field := range v.Fields {
					if err := visit(field.Type); err != nil {
						return err
					}
				}
			}
		} else if ty.Kind == ast.EXPR_TYPE_STRUCT {
			for _, field := range ty.T.(*ast.StructType).Decl.Fields {
				if err := visit(field.Type); err != nil {
					return fmt.Errorf("%s: %w", field.Name.Pos, err)
				}
			}
		} else if ty.Kind == ast.EXPR_TYPE_ARRAY {
			if err := visit(ty.T.(*ast.SequenceType).Element); err != nil {
				return err
			}
		} else if ty.Kind == ast.EXPR_TYPE_TUPLE {
			for _, t := range ty.T.(*ast.TupleType).Types {
				if err := visit(t); err != nil {
					return err
				}
			}
		}
		state[key] = 2
		s.program.ValueTypes = append(s.program.ValueTypes, ty)
		return nil
	}
	for _, d := range s.program.Enums {
		if err := visit(&ast.ExprType{Kind: ast.EXPR_TYPE_ENUM, T: &ast.EnumType{Decl: d}}); err != nil {
			return err
		}
	}
	for _, d := range s.program.Structs {
		if err := visit(&ast.ExprType{Kind: ast.EXPR_TYPE_STRUCT, T: &ast.StructType{Decl: d}}); err != nil {
			for i := len(s.generics.origins[d]) - 1; i >= 0; i-- {
				err = fmt.Errorf("%w\n  instantiated %s", err, s.generics.origins[d][i])
			}
			return err
		}
	}
	for _, ty := range s.program.TupleTypes {
		if err := visit(ty); err != nil {
			return err
		}
	}
	for _, ty := range s.program.SequenceTypes {
		if err := visit(ty); err != nil {
			return err
		}
	}
	return nil
}
