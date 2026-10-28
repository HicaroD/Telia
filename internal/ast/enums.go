package ast

import "github.com/HicaroD/Telia/internal/lexer/token"

// EnumDecl has nominal identity independent of its integer representation.
type EnumDecl struct {
	Owner      *Package
	File       *File
	Name       *token.Token
	Type       *ExprType
	Variants   []*EnumVariant
	Resolved   bool
	Resolving  bool
	CName      string
	HasPayload bool
	TypeParams []*token.Token
}

type EnumVariant struct {
	Name   *token.Token
	Expr   *Node
	Value  *LiteralExpr
	Fields []*StructField
	Index  int
}

type EnumType struct{ Decl *EnumDecl }

func (e *EnumType) String() string { return e.Decl.Name.Name() }

type MatchExpr struct {
	Value *Node
	Arms  []*MatchArm
	Type  *ExprType
	Enum  *EnumDecl
}

type MatchArm struct {
	PatternType *ExprType
	VariantName *token.Token
	Variant     *EnumVariant
	Fallback    bool
	OmitRest    bool
	Fields      []*PatternField
	Scope       *Scope
	Result      *Node
	Block       *BlockStmt
}

type YieldStmt struct{ Value *Node }

type PatternField struct {
	Name    *token.Token
	Binding *Param
}
