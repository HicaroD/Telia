package sema

import (
	"fmt"
	"go/constant"
	gotoken "go/token"
	"math"
	"strconv"

	"github.com/HicaroD/Telia/internal/ast"
	"github.com/HicaroD/Telia/internal/diagnostics"
	"github.com/HicaroD/Telia/internal/lexer/token"
	"github.com/HicaroD/Telia/internal/source"
)

func (s *sema) checkConstDecl(d *ast.ConstDecl) error {
	if d.Foreign != nil {
		if d.Foreign.Attributes.Header == "" {
			return diagnostics.Diag{Title: "Foreign constant requires a header", Message: "foreign constants obtain their values from a C header", Span: d.Name.Span, Help: "Add a header attribute to the extern block."}
		}
		ty, err := s.resolveType(d.Type, d.Scope)
		if err != nil {
			return err
		}
		if ty.Kind != ast.EXPR_TYPE_BASIC || (!ty.IsNumeric() && ty.T.(*ast.BasicType).Kind != token.BOOL_TYPE) {
			return diagnostics.Diag{Title: "Invalid foreign constant type", Message: "foreign constants require a scalar type", Span: d.Name.Span, Help: "Declare the header constant with its numeric or boolean type."}
		}
		d.Type = ty
		return nil
	}
	if d.Value != nil {
		return nil
	}
	bindings := s.generics.bindings
	if d.Scope == d.Owner.Scope {
		bindings = nil
	}
	return s.inContext(d.Owner, d.File, bindings, func() error {
		return s.evaluateConstDecl(d)
	})
}

func (s *sema) evaluateConstDecl(d *ast.ConstDecl) error {
	if d.Evaluating {
		return diagnostics.Diag{Title: "Cyclic constant", Message: fmt.Sprintf("cyclic constant '%s'", d.Name.Name()), Span: d.Expr.Span, Help: "Remove the dependency that leads back to this constant.", Related: []diagnostics.Related{{Span: d.Name.Span, Message: "constant declared here"}}}
	}
	d.Evaluating = true
	defer func() { d.Evaluating = false }()
	var err error
	if d.Type != nil {
		d.Type, err = s.resolveType(d.Type, d.Scope)
		if err == nil {
			_, err = s.inferExprTypeWithContext(d.Expr, d.Type, d.Scope, d.Scope, false, false)
		}
	} else {
		d.Type, _, err = s.inferExprTypeWithoutContext(d.Expr, d.Scope, d.Scope, false, false)
	}
	if err != nil {
		return err
	}
	value, err := evaluateConstant(d.Expr)
	if err != nil {
		return err
	}
	if err := checkConstantRange(value, d.Type, d.Expr.Span); err != nil {
		return err
	}
	text := value.ExactString()
	if value.Kind() == constant.String {
		text = constant.StringVal(value)
	} else if value.Kind() == constant.Float {
		f, _ := constant.Float64Val(value)
		if d.Type.Explicit && d.Type.T.(*ast.BasicType).Kind.BitSize() == 32 {
			f32, _ := constant.Float32Val(value)
			f = float64(f32)
		}
		text = strconv.FormatFloat(f, 'g', -1, 64)
	}
	d.Value = &ast.LiteralExpr{Type: d.Type, Value: []byte(text)}
	return nil
}

func checkConstantRange(value constant.Value, ty *ast.ExprType, span source.Span) error {
	if ty.IsInteger() {
		kind := ty.T.(*ast.BasicType).Kind
		bits := uint(kind.BitSize())
		unsigned := kind >= token.UINT_TYPE && kind <= token.U128_TYPE
		if !unsigned {
			bits--
		}
		limit := constant.Shift(constant.MakeInt64(1), gotoken.SHL, bits)
		minimum := constant.MakeInt64(0)
		if !unsigned {
			minimum = constant.UnaryOp(gotoken.SUB, limit, 0)
		}
		if value.Kind() != constant.Int || constant.Compare(value, gotoken.LSS, minimum) || !constant.Compare(value, gotoken.LSS, limit) {
			return diagnostics.Diag{Title: "Constant overflow", Message: fmt.Sprintf("constant value %s is outside the range of %s", value.ExactString(), kind), Span: span, Help: "Choose a wider integer type or a value representable by the declared type."}
		}
	}
	if ty.IsNumeric() && ty.T.(*ast.BasicType).Kind.IsFloat() {
		kind := ty.T.(*ast.BasicType).Kind
		f, _ := constant.Float64Val(value)
		if kind.BitSize() == 32 {
			f32, _ := constant.Float32Val(value)
			f = float64(f32)
		}
		if math.IsInf(f, 0) {
			return diagnostics.Diag{Title: "Constant overflow", Message: fmt.Sprintf("constant value is outside the finite range of %s", kind), Span: span, Help: "Choose a wider floating-point type or a smaller value."}
		}
	}
	return nil
}

func (s *sema) inferConstReference(id *ast.IdExpr, d *ast.ConstDecl, expected *ast.ExprType) (*ast.ExprType, error) {
	if err := s.checkConstDecl(d); err != nil {
		return nil, err
	}
	if d.Foreign != nil {
		if expected != nil && !expected.Equals(d.Type) {
			return nil, diagnostics.Diag{Title: "Foreign constant type mismatch", Message: fmt.Sprintf("expected %s, got %s", expected.T, d.Type.T), Span: id.Name.Span, Help: "Use the declared foreign constant type or an explicit conversion.", Related: []diagnostics.Related{{Span: d.Name.Span, Message: "foreign constant declared here"}}}
		}
		return d.Type, nil
	}
	lit := &ast.LiteralExpr{Type: cloneType(d.Type, nil), Value: d.Value.Value}
	if expected != nil {
		if _, err := s.inferLiteralExprTypeWithContext(lit, cloneType(expected, nil)); err != nil {
			return nil, err
		}
	}
	value, err := evaluateConstant(&ast.Node{Kind: ast.KIND_LITERAL_EXPR, Node: lit})
	if err != nil {
		return nil, err
	}
	if err := checkConstantRange(value, lit.Type, id.Name.Span); err != nil {
		return nil, err
	}
	id.Constant = lit
	return lit.Type, nil
}

func evaluateConstant(node *ast.Node) (constant.Value, error) {
	switch node.Kind {
	case ast.KIND_NAMESPACE_ACCESS:
		return evaluateConstant(node.Node.(*ast.NamespaceAccess).Right)
	case ast.KIND_LITERAL_EXPR:
		lit := node.Node.(*ast.LiteralExpr)
		kind := lit.Type.T.(*ast.BasicType).Kind
		switch {
		case kind == token.BOOL_TYPE:
			return constant.MakeBool(string(lit.Value) == "true"), nil
		case kind == token.STRING_TYPE || kind == token.CSTRING_TYPE:
			return constant.MakeString(string(lit.Value)), nil
		case kind.IsInteger():
			return constant.MakeFromLiteral(string(lit.Value), gotoken.INT, 0), nil
		case kind.IsFloat():
			return constant.MakeFromLiteral(string(lit.Value), gotoken.FLOAT, 0), nil
		}
	case ast.KIND_ID_EXPR:
		id := node.Node.(*ast.IdExpr)
		if id.N != nil && id.N.Kind == ast.KIND_CONST_DECL && id.Constant != nil {
			return evaluateConstant(&ast.Node{Kind: ast.KIND_LITERAL_EXPR, Node: id.Constant})
		}
	case ast.KIND_UNARY_EXPR:
		u := node.Node.(*ast.UnaryExpr)
		v, err := evaluateConstant(u.Value)
		if err != nil {
			return nil, err
		}
		op := constantOperator(u.Op)
		var precision uint
		if u.Op == token.TILDE {
			kind := u.Type.T.(*ast.BasicType).Kind
			if kind >= token.UINT_TYPE && kind <= token.U128_TYPE {
				precision = uint(kind.BitSize())
			}
		}
		result := constant.UnaryOp(op, v, precision)
		if err := checkConstantRange(result, u.Type, node.Span); err != nil {
			return nil, err
		}
		return result, nil
	case ast.KIND_BINARY_EXPR:
		b := node.Node.(*ast.BinExpr)
		left, err := evaluateConstant(b.Left)
		if err != nil {
			return nil, err
		}
		right, err := evaluateConstant(b.Right)
		if err != nil {
			return nil, err
		}
		if b.Op == token.SHIFT_LEFT || b.Op == token.SHIFT_RIGHT {
			return evaluateShift(left, right, b)
		}
		op := constantOperator(b.Op)
		if b.Op.IsCmpOp() {
			return constant.MakeBool(constant.Compare(left, op, right)), nil
		}
		if b.Op == token.SLASH {
			if constant.Sign(right) == 0 {
				return nil, diagnostics.Diag{Title: "Division by zero", Message: "division by zero in constant expression", Span: b.Right.Span, Help: "Use a nonzero divisor in this constant expression."}
			}
			if left.Kind() == constant.Int && right.Kind() == constant.Int {
				op = gotoken.QUO_ASSIGN
			}
		}
		result := constant.BinaryOp(left, op, right)
		if err := checkConstantRange(result, b.OperandType, node.Span); err != nil {
			return nil, err
		}
		return result, nil
	}
	return nil, diagnostics.Diag{Title: "Nonconstant initializer", Message: "constant initializer must be a compile-time expression", Span: node.Span, Help: "Use literals and other constants, or declare a runtime variable with :=."}
}

func (s *sema) inferImportedValue(id *ast.IdExpr, scope *ast.Scope, expected *ast.ExprType) (*ast.ExprType, error) {
	symbol, err := scope.LookupCurrentScope(id.Name.Name())
	if err == nil && symbol.Kind == ast.KIND_VA_INTRINSIC {
		return nil, vaIntrinsicValueDiagnostic(id.Name.Span)
	}
	if err == nil && (symbol.Kind == ast.KIND_FN_DECL || symbol.Kind == ast.KIND_PROTO) {
		id.N = symbol
		if symbol.Invalid {
			return nil, diagnostics.COMPILER_ERROR_FOUND
		}
		ty, err := s.functionValue(id, symbol)
		if err != nil {
			return nil, err
		}
		if expected != nil && !ty.Equals(expected) {
			return nil, diagnostics.Diag{Title: "Callback signature mismatch", Message: "imported function does not have the required callback signature", Span: id.Name.Span, Help: "Match its parameter types, return type, and calling convention.", Related: []diagnostics.Related{{Span: symbol.DeclarationSpan(), Message: "function declared here"}}}
		}
		return ty, nil
	}
	if err != nil || symbol.Kind != ast.KIND_CONST_DECL {
		return nil, diagnostics.Diag{Title: "Unknown imported constant", Message: fmt.Sprintf("package has no constant '%s'", id.Name.Name()), Span: id.Name.Span}
	}
	id.N = symbol
	if symbol.Invalid {
		return nil, diagnostics.COMPILER_ERROR_FOUND
	}
	return s.inferConstReference(id, symbol.Node.(*ast.ConstDecl), expected)
}

func constantOperator(op token.Kind) gotoken.Token {
	switch op {
	case token.PLUS:
		return gotoken.ADD
	case token.MINUS:
		return gotoken.SUB
	case token.STAR:
		return gotoken.MUL
	case token.SLASH:
		return gotoken.QUO
	case token.AMPERSAND:
		return gotoken.AND
	case token.PIPE:
		return gotoken.OR
	case token.CARET, token.TILDE:
		return gotoken.XOR
	case token.NOT:
		return gotoken.NOT
	case token.AND:
		return gotoken.LAND
	case token.OR:
		return gotoken.LOR
	case token.EQUAL_EQUAL:
		return gotoken.EQL
	case token.BANG_EQUAL:
		return gotoken.NEQ
	case token.LESS:
		return gotoken.LSS
	case token.LESS_EQ:
		return gotoken.LEQ
	case token.GREATER:
		return gotoken.GTR
	case token.GREATER_EQ:
		return gotoken.GEQ
	default:
		panic("constant operator must have passed semantic validation")
	}
}
