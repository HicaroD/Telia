package ast

import (
	"errors"
	"fmt"
)

var (
	ERR_SYMBOL_ALREADY_DEFINED_ON_SCOPE = errors.New("symbol already defined on scope")
	ERR_SYMBOL_NOT_FOUND_ON_SCOPE       = errors.New("symbol not found on scope")
)

// TODO: since I won't be using the scope somewhere else in the future,
// I could remove the generics and make the code more efficient
type Scope struct {
	Parent *Scope
	Nodes  map[string]Node
}

func NewScope(parent *Scope) *Scope {
	return &Scope{Parent: parent, Nodes: map[string]Node{}}
}

func (scope *Scope) Insert(name string, element Node) error {
	if _, ok := scope.Nodes[name]; ok {
		return ERR_SYMBOL_ALREADY_DEFINED_ON_SCOPE
	}
	scope.Nodes[name] = element
	return nil
}

func (scope *Scope) LookupCurrentScope(name string) (Node, error) {
	if node, ok := scope.Nodes[name]; ok {
		return node, nil
	}
	return nil, ERR_SYMBOL_NOT_FOUND_ON_SCOPE
}

func (scope *Scope) LookupAcrossScopes(name string) (Node, error) {
	if node, ok := scope.Nodes[name]; ok {
		return node, nil
	}
	if scope.Parent == nil {
		return nil, ERR_SYMBOL_NOT_FOUND_ON_SCOPE
	}
	return scope.Parent.LookupAcrossScopes(name)
}

func (scope Scope) String() string {
	if scope.Parent == nil {
		return fmt.Sprintf("Scope:\nParent: nil\nCurrent: %v\n", scope.Nodes)
	}
	return fmt.Sprintf("Scope:\nParent: %v\nCurrent: %v\n", scope.Parent, scope.Nodes)
}
