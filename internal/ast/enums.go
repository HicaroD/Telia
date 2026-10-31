package ast

import "github.com/HicaroD/Telia/internal/lexer/token"

// EnumDecl has nominal identity independent of its integer representation.
type EnumDecl struct {
	Foreign    *ExternDecl
	Attributes Attributes
	Owner      *Package
	File       *File
	Name       *token.Token
	Type       *ExprType
	Variants   []*EnumVariant
	Resolved   bool
	Resolving  bool
	CName      string
	HasPayload bool
	IsError    bool
	TypeParams []*token.Token
}

type EnumVariant struct {
	Attributes Attributes
	Name       *token.Token
	Expr       *Node
	Value      *LiteralExpr
	Fields     []*StructField
	Index      int
}

type EnumType struct{ Decl *EnumDecl }

func (e *EnumType) String() string { return e.Decl.Name.Name() }

type MatchExpr struct {
	Value     *Node
	Arms      []*MatchArm
	Type      *ExprType
	ValueType *ExprType
	Enum      *EnumDecl
	IsError   bool
}

type MatchArm struct {
	MatchPattern
	Alternatives []*MatchPattern
	Guard        *Node
	Scope        *Scope
	Result       *Node
	Block        *BlockStmt
}

// Patterns returns the alternatives in source order, including the first pattern.
func (a *MatchArm) Patterns() []*MatchPattern {
	return append([]*MatchPattern{&a.MatchPattern}, a.Alternatives...)
}

type MatchPattern struct {
	PatternType *ExprType
	Literal     *Node
	Binding     *Param
	VariantName *token.Token
	Variant     *EnumVariant
	Enum        *EnumDecl
	Fallback    bool
	OmitRest    bool
	Fields      []*PatternField
}

type YieldStmt struct{ Value *Node }

type PatternField struct {
	Name    *token.Token
	Binding *Param
}
