package sema

import (
	"fmt"
	"go/constant"
	gotoken "go/token"

	"github.com/HicaroD/Telia/internal/ast"
	"github.com/HicaroD/Telia/internal/diagnostics"
	"github.com/HicaroD/Telia/internal/lexer/token"
	"github.com/HicaroD/Telia/internal/source"
)

func enumTypeMismatch(actual, expected *ast.ExprType, span source.Span) error {
	d := actual.T.(*ast.EnumType).Decl
	return diagnostics.Diag{
		Title:   "Enum type mismatch",
		Message: fmt.Sprintf("expected %s, got enum %s", expected.T, d.Name.Name()),
		Span:    span,
		Help:    "Use a value of the expected enum type; an enum is distinct from its integer representation.",
		Related: []diagnostics.Related{{Span: d.Name.Span, Message: "enum declared here"}},
	}
}

func (s *sema) enumInstance(
	template *ast.EnumDecl,
	args []*ast.ExprType,
	scope *ast.Scope,
	site *token.Token,
) (*ast.EnumDecl, error) {
	if len(args) != len(template.TypeParams) {
		return nil, diagnostics.Diag{
			Title: "Invalid enum type arguments",
			Message: fmt.Sprintf(
				"enum '%s' requires %d type arguments, got %d",
				template.Name.Name(),
				len(template.TypeParams),
				len(args),
			),
			Span: site.Span,
			Help: "Supply exactly the type arguments declared by this enum.",
			Related: []diagnostics.Related{
				{Span: template.Name.Span, Message: "enum declared here"},
			},
		}
	}
	if len(args) == 0 {
		return template, nil
	}
	bindings, suffix, err := s.concreteArguments(
		template.TypeParams,
		args,
		scope,
		map[*ast.TypeAlias]bool{},
	)
	if err != nil {
		return nil, err
	}
	key := declarationKey(template.Owner, template.Name.Name()) + suffix
	if s.generics.enums == nil {
		s.generics.enums = map[string]*ast.EnumDecl{}
	}
	if cached := s.generics.enums[key]; cached != nil {
		return cached, nil
	}
	if err := s.beginInstance(template.Name.Name()+suffix, site); err != nil {
		return nil, err
	}
	d := &ast.EnumDecl{
		Owner:      template.Owner,
		File:       template.File,
		Name:       template.Name,
		Type:       cloneType(template.Type, bindings),
		HasPayload: template.HasPayload,
		CName:      concreteName("enum", key),
	}
	for _, v := range template.Variants {
		copy := &ast.EnumVariant{
			Name: v.Name,
			Expr: (syntaxClone{bindings: bindings}).node(v.Expr, template.Owner.Scope),
		}
		for _, field := range v.Fields {
			copy.Fields = append(
				copy.Fields,
				&ast.StructField{
					Name:  field.Name,
					Index: field.Index,
					Type:  cloneType(field.Type, bindings),
				},
			)
		}
		d.Variants = append(d.Variants, copy)
	}
	s.generics.enums[key] = d
	err = s.resolveEnum(d)
	return d, s.endInstance(err)
}

func (s *sema) resolveEnum(d *ast.EnumDecl) error {
	if d.Resolved {
		return nil
	}
	if d.Foreign != nil {
		if d.Type != nil {
			ty, err := s.resolveType(d.Type, d.Foreign.Scope)
			if err != nil {
				return err
			}
			if !ty.IsInteger() {
				return diagnostics.Diag{
					Title:   "Invalid enum representation",
					Message: "a foreign enum representation annotation must be an integer type",
					Span:    d.Name.Span,
					Help:    "Omit the annotation to use the header's enum type, or name its exact integer representation.",
				}
			}
			d.Type = ty
		}
		if d.Foreign.Attributes.Header == "" {
			return diagnostics.Diag{
				Title:   "Foreign enum requires a header",
				Message: "foreign enums use C header types and tags",
				Span:    d.Name.Span,
				Help:    "Add a header attribute to the extern block.",
			}
		}
		if d.HasPayload || len(d.TypeParams) != 0 {
			return diagnostics.Diag{
				Title:   "Invalid foreign enum",
				Message: "foreign enums cannot have payloads or type parameters",
				Span:    d.Name.Span,
				Help:    "Declare only the named C enum members.",
			}
		}
		d.CName = d.Attributes.LinkName
		if d.CName == "" {
			d.CName = d.Name.Name()
		}
		seen := map[string]*token.Token{}
		for _, v := range d.Variants {
			if previous := seen[v.Name.Name()]; previous != nil {
				return diagnostics.Diag{
					Title:   "Duplicate enum variant",
					Message: "foreign enum member is already declared",
					Span:    v.Name.Span,
					Help:    "Choose a distinct member name.",
					Related: []diagnostics.Related{
						{Span: previous.Span, Message: "previous member"},
					},
				}
			}
			seen[v.Name.Name()] = v.Name
			if v.Expr != nil {
				return diagnostics.Diag{
					Title:   "Foreign enum initializer",
					Message: "the header supplies foreign enum tag values",
					Span:    v.Expr.Span,
					Help:    "Remove the tag initializer and use link_name to select the C member.",
				}
			}
		}
		d.Resolved = true
		s.program.Enums = append(s.program.Enums, d)
		return nil
	}
	if d.Resolving {
		return diagnostics.Diag{
			Title:   "Cyclic enum declaration",
			Message: "enum representation or tags depend on the enum being defined",
			Span:    d.Name.Span,
			Help:    "Use an independent integer representation and constant tag expressions.",
			Related: []diagnostics.Related{{Span: d.Name.Span, Message: "resolution started here"}},
		}
	}
	d.Resolving = true
	defer func() { d.Resolving = false }()
	if d.CName == "" {
		d.CName = concreteName("enum", declarationKey(d.Owner, d.Name.Name()))
	}
	return s.inContext(d.Owner, d.File, nil, func() error {
		var err error
		d.Type, err = s.resolveType(d.Type, d.Owner.Scope)
		if err != nil {
			return err
		}
		if !d.Type.IsInteger() {
			return diagnostics.Diag{
				Title:   "Invalid enum representation",
				Message: "enum tags require an integer type",
				Span:    d.Name.Span,
				Help:    "Use an integer type such as u8 or i32, or an alias to one.",
			}
		}
		d.Type.Explicit = true
		d.Type.T.(*ast.BasicType).Explicit = true
		seen := map[string]*token.Token{}
		tags := map[string]*token.Token{}
		next := constant.MakeInt64(0)
		for index, v := range d.Variants {
			v.Index = index
			if previous := seen[v.Name.Name()]; previous != nil {
				return diagnostics.Diag{
					Title:   "Duplicate enum variant",
					Message: fmt.Sprintf("variant '%s' is already declared", v.Name.Name()),
					Span:    v.Name.Span,
					Help:    "Choose a distinct variant name.",
					Related: []diagnostics.Related{
						{Span: previous.Span, Message: "previous variant"},
					},
				}
			}
			seen[v.Name.Name()] = v.Name
			value := next
			if v.Expr != nil {
				if _, err := s.inferExprTypeWithContext(v.Expr, d.Type, d.Owner.Scope, d.Owner.Scope, false, false); err != nil {
					return err
				}
				value, err = evaluateConstant(v.Expr)
				if err != nil {
					return err
				}
			}
			if err := checkConstantRange(value, d.Type, v.Name.Span); err != nil {
				return err
			}
			if previous := tags[value.ExactString()]; previous != nil {
				return diagnostics.Diag{
					Title: "Duplicate enum tag",
					Message: fmt.Sprintf(
						"tag %s identifies more than one variant",
						value.ExactString(),
					),
					Span: v.Name.Span,
					Help: "Give each variant a distinct integer tag.",
					Related: []diagnostics.Related{
						{Span: previous.Span, Message: "previous variant with this tag"},
					},
				}
			}
			tags[value.ExactString()] = v.Name
			v.Value = &ast.LiteralExpr{Type: d.Type, Value: []byte(value.ExactString())}
			next = constant.BinaryOp(value, gotoken.ADD, constant.MakeInt64(1))
			fields := map[string]*token.Token{}
			for _, field := range v.Fields {
				if previous := fields[field.Name.Name()]; previous != nil {
					return diagnostics.Diag{
						Title:   "Duplicate payload field",
						Message: "payload field is already declared",
						Span:    field.Name.Span,
						Help:    "Give each named payload field a distinct name.",
						Related: []diagnostics.Related{
							{Span: previous.Span, Message: "previous field"},
						},
					}
				}
				fields[field.Name.Name()] = field.Name
				field.Type, err = s.resolveType(field.Type, d.Owner.Scope)
				if err != nil {
					return err
				}
				if !valueType(field.Type) {
					return diagnostics.Diag{
						Title:   "Invalid payload field type",
						Message: "payload fields require value types",
						Span:    field.Name.Span,
						Help:    "Choose a concrete value type for this field.",
					}
				}
			}
		}
		d.Resolved = true
		s.program.Enums = append(s.program.Enums, d)
		return nil
	})
}

func (s *sema) checkEnumAccess(
	access *ast.NamespaceAccess,
	d *ast.EnumDecl,
	scope *ast.Scope,
) (*ast.ExprType, error) {
	if err := s.resolveEnum(d); err != nil {
		return nil, err
	}
	var name *token.Token
	var values []*ast.StructFieldValue
	switch access.Right.Kind {
	case ast.KIND_ID_EXPR:
		name = access.Right.Node.(*ast.IdExpr).Name
	case ast.KIND_STRUCT_EXPR:
		lit := access.Right.Node.(*ast.StructLiteralExpr)
		name, values = lit.Name, lit.Values
	default:
		return nil, diagnostics.Diag{
			Title:   "Invalid enum construction",
			Message: "expected a named enum variant",
			Span:    access.Left.Name.Span,
			Help:    "Construct a payload-free variant with Enum::Variant.",
			Related: []diagnostics.Related{{Span: d.Name.Span, Message: "enum declared here"}},
		}
	}
	for _, v := range d.Variants {
		if v.Name.Name() == name.Name() {
			seen := map[string]*token.Token{}
			for _, value := range values {
				if previous := seen[value.Name.Name()]; previous != nil {
					return nil, diagnostics.Diag{
						Title:   "Duplicate payload initializer",
						Message: "payload field is initialized twice",
						Span:    value.Name.Span,
						Help:    "Initialize each payload field once.",
						Related: []diagnostics.Related{
							{Span: previous.Span, Message: "previous initializer"},
						},
					}
				}
				seen[value.Name.Name()] = value.Name
				var field *ast.StructField
				for _, f := range v.Fields {
					if f.Name.Name() == value.Name.Name() {
						field = f
						break
					}
				}
				if field == nil {
					return nil, diagnostics.Diag{
						Title: "Unknown payload field",
						Message: fmt.Sprintf(
							"variant '%s' has no field '%s'",
							name.Name(),
							value.Name.Name(),
						),
						Span: value.Name.Span,
						Help: "Use the named fields declared by this variant.",
						Related: []diagnostics.Related{
							{Span: v.Name.Span, Message: "variant declared here"},
						},
					}
				}
				if _, err := s.inferExprTypeWithContext(value.Value, field.Type, scope, scope, false, false); err != nil {
					return nil, err
				}
			}
			for _, field := range v.Fields {
				if seen[field.Name.Name()] == nil {
					return nil, diagnostics.Diag{
						Title:   "Missing payload field",
						Message: fmt.Sprintf("missing initializer for '%s'", field.Name.Name()),
						Span:    name.Span,
						Help:    "Initialize every payload field by name.",
						Related: []diagnostics.Related{
							{Span: field.Name.Span, Message: "required field declared here"},
						},
					}
				}
			}
			access.EnumVariant = v
			access.EnumDecl = d
			return &ast.ExprType{Kind: ast.EXPR_TYPE_ENUM, T: &ast.EnumType{Decl: d}}, nil
		}
	}
	return nil, diagnostics.Diag{
		Title:   "Unknown enum variant",
		Message: fmt.Sprintf("enum '%s' has no variant '%s'", d.Name.Name(), name.Name()),
		Span:    name.Span,
		Help:    "Use a variant declared in this enum.",
		Related: []diagnostics.Related{{Span: d.Name.Span, Message: "enum declared here"}},
	}
}
