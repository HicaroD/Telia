package ast

import (
	"fmt"

	"github.com/HicaroD/Telia/internal/lexer/token"
)

type FnDecl struct {
	Scope       *Scope
	Attributes  *Attributes
	Name        *token.Token
	Params      *FieldList
	RetType     *ExprType
	Block       *BlockStmt
	BackendType any // LLVM: *values.Function
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
	Params     *FieldList
	RetType    *ExprType

	BackendType any // LLVM: *values.Function
}

func (proto Proto) String() string { return fmt.Sprintf("PROTO: %s", proto.Name) }

type PkgDecl struct {
	Name *token.Token
}

func (pkg PkgDecl) String() string { return fmt.Sprintf("PKG: %s", pkg.Name) }

type UseDecl struct {
	Path    []string
	Std     bool
	Package bool
}

func (imp UseDecl) String() string {
	return fmt.Sprintf("USE: %s | Std: %v | Package: %v", imp.Path, imp.Std, imp.Package)
}

type Attributes struct {
	LinkName string

	// Specific for extern declaration
	DefaultCallingConvention string
	LinkPrefix               string

	// Specific for prototype declaration
	Linkage string
}
