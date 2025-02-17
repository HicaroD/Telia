package ast

import (
	"fmt"

	"github.com/HicaroD/Telia/internal/lexer/token"
)

type Decl interface {
	Node
	declNode()
}

type FunctionDecl struct {
	Decl
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
func (fnDecl FunctionDecl) astNode()  {}
func (fnDecl FunctionDecl) declNode() {}

type ExternDecl struct {
	Decl
	Attributes  *ExternAttrs
	Scope       *Scope
	Name        *token.Token
	Prototypes  []*Proto
	BackendType any // LLVM: *values.Extern
}

func (extern ExternDecl) String() string {
	return fmt.Sprintf("EXTERN: %s", extern.Name)
}
func (extern ExternDecl) astNode()  {}
func (extern ExternDecl) declNode() {}

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
	Decl
	Name *token.Token
}

func (pkg PkgDecl) String() string { return fmt.Sprintf("PKG: %s", pkg.Name) }
func (pkg PkgDecl) astNode()       {}
func (pkg PkgDecl) declNode()      {}
