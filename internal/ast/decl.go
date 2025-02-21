package ast

import (
	"fmt"

	"github.com/HicaroD/Telia/internal/lexer/token"
)

type FnDecl struct {
	Scope       *Scope
	Name        *token.Token
	Params      *FieldList
	RetType     *ExprType
	Block       *BlockStmt
	BackendType any // LLVM: *values.Function
}

func (fnDecl FnDecl) String() string {
	return fmt.Sprintf(
		"Scope: %s\nName: %s\nParams: %s\nRetType: %s\nBlock: %s\n",
		fnDecl.Scope,
		fnDecl.Name,
		fnDecl.Params,
		fnDecl.RetType,
		fnDecl.Block,
	)
}

type ExternDecl struct {
	Attributes  *ExternAttrs
	Scope       *Scope
	Name        *token.Token
	Prototypes  []*Proto
	BackendType any // LLVM: *values.Extern
}

func (extern ExternDecl) String() string {
	return fmt.Sprintf("EXTERN: %s", extern.Name)
}

type ProtoAttrs struct {
	LinkName string
	Linkage  string
}

type Proto struct {
	Attributes *ProtoAttrs
	Name       *token.Token
	Params     *FieldList
	RetType    *ExprType

	BackendType any // LLVM: *values.Function
}

func (proto Proto) String() string { return fmt.Sprintf("PROTO: %s", proto.Name) }

type ExternAttrs struct {
	DefaultCallingConvention string
	LinkPrefix               string
	LinkName                 string
}

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
