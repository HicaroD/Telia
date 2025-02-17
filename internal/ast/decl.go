package ast

import (
	"fmt"

	"github.com/HicaroD/Telia/internal/lexer/token"
)

type FunctionDecl struct {
	Node
	Scope       *Scope
	Name        *token.Token
	Params      *FieldList
	RetType     ExprType
	Block       *BlockStmt
	BackendType any // LLVM: *values.Function
}

func (fnDecl FunctionDecl) String() string {
	return fmt.Sprintf(
		"Scope: %s\nName: %s\nParams: %s\nRetType: %s\nBlock: %s\n",
		fnDecl.Scope,
		fnDecl.Name,
		fnDecl.Params,
		fnDecl.RetType,
		fnDecl.Block,
	)
}
func (fnDecl FunctionDecl) astNode() {}

type ExternDecl struct {
	// Decl
	Attributes  *ExternAttrs
	Scope       *Scope
	Name        *token.Token
	Prototypes  []*Proto
	BackendType any // LLVM: *values.Extern
}

func (extern ExternDecl) String() string {
	return fmt.Sprintf("EXTERN: %s", extern.Name)
}
func (extern ExternDecl) astNode() {}

// TODO: add attribute for prototype, such as link_name and linkage type
type ProtoAttrs struct {
	LinkName string
	Linkage  string
}

type Proto struct {
	Attributes *ProtoAttrs
	Name       *token.Token
	Params     *FieldList
	RetType    ExprType

	BackendType any // LLVM: *values.Function
}

func (proto Proto) String() string { return fmt.Sprintf("PROTO: %s", proto.Name) }
func (proto Proto) astNode()       {}

type ExternAttrs struct {
	DefaultCallingConvention string
	LinkPrefix               string
	LinkName                 string
}

type PkgDecl struct {
	Name *token.Token
}

func (pkg PkgDecl) String() string { return fmt.Sprintf("PKG: %s", pkg.Name) }
func (pkg PkgDecl) astNode()       {}
