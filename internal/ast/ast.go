// Package ast defines the abstract syntax tree (AST) for a programming language.
package ast

type Node interface {
	astNode()
}

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
	EXPR_END // expression node start delimiter

	KIND_FIELD
	KIND_PROTO
)

type MyNode struct {
	Kind NodeKind
	Node any
}

func (n *MyNode) IsStmt() bool {
	return n.Kind > STMT_START && n.Kind < STMT_END
}

func (n *MyNode) IsExpr() bool {
	return n.Kind > EXPR_START && n.Kind < EXPR_END
}

func (n *MyNode) IsDecl() bool {
	return n.Kind > DECL_START && n.Kind < DECL_END
}
