package ast

import (
	"fmt"

	"github.com/HicaroD/Telia/internal/lexer/token"
)

type StructDecl struct {
	Attributes *Attributes
	Name       *token.Token
	Fields     []*StructField
}

func (st StructDecl) String() string {
	return fmt.Sprintf("STRUCT: %v | Fields: %v\n", st.Name.Name(), st.Fields)
}

type StructField struct {
	Name *token.Token
	Type *ExprType
}

type FnDecl struct {
	Scope      *Scope
	Attributes *Attributes
	Name       *token.Token
	Params     *Params
	RetType    *ExprType
	Block      *BlockStmt
}

func (fnDecl FnDecl) String() string {
	return fmt.Sprintf(
		"Scope: %v\nName: %v\nParams: %v\nRetType: %v\nBlock: %v\n",
		fnDecl.Scope,
		fnDecl.Name,
		fnDecl.Params,
		fnDecl.RetType,
		fnDecl.Block,
	)
}

type ExternDecl struct {
	Scope       *Scope
	Attributes  *Attributes
	Name        *token.Token
	Prototypes  []*Proto
	BackendType any // LLVM: *values.Extern
}

func (extern ExternDecl) String() string {
	return fmt.Sprintf("EXTERN: %s", extern.Name)
}

type Proto struct {
	Attributes *Attributes
	Name       *token.Token
	Params     *Params
	RetType    *ExprType
}

func (proto Proto) String() string { return fmt.Sprintf("PROTO: %s", proto.Name) }

type PkgDecl struct {
	Name *token.Token
}

func (pkg PkgDecl) String() string { return fmt.Sprintf("PKG: %s", pkg.Name) }

type UseDecl struct {
	Name    string
	Package *Package
}

func (imp UseDecl) String() string {
	return fmt.Sprintf("USE: %s", imp.Name)
}

type Attributes struct {
	LinkName string
	Linkage  string

	// Specific for extern declaration
	DefaultCallingConvention string
	LinkPrefix               string
	Global                   bool
}

func (a *Attributes) String() string {
	if a == nil {
		return "No attributes"
	}
	return fmt.Sprintf("LinkName: '%s' | Linkage: '%s' | DefaultCC: '%s' | LinkPrefix: '%s'\n", a.LinkName, a.Linkage, a.DefaultCallingConvention, a.LinkPrefix)
}

type AtOperatorKind int

const (
	// @fail
	AT_OPERATOR_FAIL AtOperatorKind = iota
	// @prop
	AT_OPERATOR_PROP
	// @catch <name> {...}
	AT_OPERATOR_CATCH
)

type AtOperator struct {
	Kind AtOperatorKind
	Op   any
}

type CatchAtOperator struct {
	Scope      *Scope
	ErrVarName *token.Token
	Block      *BlockStmt
}
