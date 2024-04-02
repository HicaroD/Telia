package ast

import "errors"

var (
	SYMBOL_ALREADY_DEFINED_ON_SCOPE error = errors.New("symbol already defined on scope")
	SYMBOL_NOT_FOUND_ON_SCOPE error = errors.New("symbol not found on scope")
)

type Scope struct {
	parent   *Scope
	elements map[string]AstNode
}

func NewScope(parent *Scope) *Scope {
	return &Scope{parent: parent, elements: map[string]AstNode{}}
}

func (scope *Scope) Insert(name string, element AstNode) error {
	if _, ok := scope.elements[name]; ok {
		return SYMBOL_ALREADY_DEFINED_ON_SCOPE
	}
	scope.elements[name] = element
	return nil
}

func (scope *Scope) Lookup(name string) AstNode {
	if node, ok := scope.elements[name]; ok {
		return node
	}
	return nil
}
