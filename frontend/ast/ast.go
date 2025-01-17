// Package ast defines the abstract syntax tree (AST) for a programming language.
package ast

type Node interface {
	astNode()
}
