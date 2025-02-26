// Package ast defines the abstract syntax tree (AST) for a programming language.
package ast

import "fmt"

type NodeKind int

const (
	DECL_START NodeKind = iota // declaration node start delimiter

	KIND_FN_DECL
	KIND_EXTERN_DECL
	KIND_PKG_DECL
	KIND_USE_DECL
	KIND_TYPE_ALIAS_DECL

	DECL_END // declaration node end delimiter

	STMT_START // statement node start delimiter
	KIND_BLOCK_STMT
	KIND_MULTI_VAR_STMT
	KIND_VAR_STMT
	KIND_RETURN_STMT
	KIND_COND_STMT
	KIND_IF_ELIF_STMT
	KIND_ELSE_STMT
	KIND_FOR_LOOP_STMT
	KIND_WHILE_LOOP_STMT

	EXPR_START // expression node start delimiter

	KIND_FN_CALL      // expression and statement
	KIND_FIELD_ACCESS // expression and statement

	STMT_END // statement node end delimiter
	KIND_VOID_EXPR
	KIND_LITERAl_EXPR
	KIND_ID_EXPR
	KIND_UNARY_EXPR
	KIND_BINARY_EXPR
	KIND_TUPLE_EXPR
	EXPR_END // expression node start delimiter

	KIND_FIELD
	KIND_PROTO
)

type Node struct {
	Kind NodeKind
	Node any
}

func (n *Node) IsStmt() bool {
	return n.Kind > STMT_START && n.Kind < STMT_END
}

func (n *Node) IsExpr() bool {
	return n.Kind > EXPR_START && n.Kind < EXPR_END
}

func (n *Node) IsDecl() bool {
	return n.Kind > DECL_START && n.Kind < DECL_END
}

func (n *Node) IsId() bool {
	return n.Kind == KIND_ID_EXPR
}

func (n *Node) IsReturn() bool {
	return n.Kind == KIND_RETURN_STMT
}

func (n *Node) IsVoid() bool {
	return n.Kind == KIND_VOID_EXPR
}

func (n *Node) String() string {
	switch n.Kind {
	case KIND_FN_DECL:
		return "KIND_FN_DECL"
	case KIND_EXTERN_DECL:
		return "KIND_EXTERN_DECL"
	case KIND_PKG_DECL:
		return "KIND_PKG_DECL"
	case KIND_USE_DECL:
		return "KIND_USE_DECL"
	case KIND_TYPE_ALIAS_DECL:
		return "KIND_TYPE_ALIAS_DECL"
	case KIND_BLOCK_STMT:
		return "KIND_BLOCK_STMT"
	case KIND_MULTI_VAR_STMT:
		return "KIND_MULTI_VAR_STMT"
	case KIND_VAR_STMT:
		return "KIND_VAR_STMT"
	case KIND_RETURN_STMT:
		return "KIND_RETURN_STMT"
	case KIND_COND_STMT:
		return "KIND_COND_STMT"
	case KIND_IF_ELIF_STMT:
		return "KIND_IF_ELIF_STMT"
	case KIND_ELSE_STMT:
		return "KIND_ELSE_STMT"
	case KIND_FOR_LOOP_STMT:
		return "KIND_FOR_LOOP_STMT"
	case KIND_WHILE_LOOP_STMT:
		return "KIND_WHILE_LOOP_STMT"
	case KIND_FN_CALL:
		return "KIND_FN_CALL"
	case KIND_FIELD_ACCESS:
		return "KIND_FIELD_ACCESS"
	case KIND_VOID_EXPR:
		return "KIND_VOID_EXPR"
	case KIND_LITERAl_EXPR:
		return "KIND_LITERAl_EXPR"
	case KIND_ID_EXPR:
		return "KIND_ID_EXPR"
	case KIND_UNARY_EXPR:
		return "KIND_UNARY_EXPR"
	case KIND_BINARY_EXPR:
		return "KIND_BINARY_EXPR"
	case KIND_TUPLE_EXPR:
		return "KIND_TUPLE_EXPR"
	case KIND_FIELD:
		return "KIND_FIELD"
	case KIND_PROTO:
		return "KIND_PROTO"
	default:
		return fmt.Sprintf("Unknown Node Kind: %v", n.Kind)
	}
}
