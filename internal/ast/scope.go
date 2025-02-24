package ast

import (
	"errors"
	"fmt"
)

var (
	ErrSymbolAlreadyDefinedOnScope = errors.New("symbol already defined on scope")
	ErrSymbolNotFoundOnScope       = errors.New("symbol not found on scope")
)

type Scope struct {
	Parent *Scope
	Nodes  map[string]*Node
}

func NewScope(parent *Scope) *Scope {
	return &Scope{Parent: parent, Nodes: make(map[string]*Node)}
}

func (scope *Scope) Insert(name string, element *Node) error {
	if _, ok := scope.Nodes[name]; ok {
		return ErrSymbolAlreadyDefinedOnScope
	}
	scope.Nodes[name] = element
	return nil
}

func (scope *Scope) LookupCurrentScope(name string) (*Node, error) {
	if node, ok := scope.Nodes[name]; ok {
		return node, nil
	}
	return nil, ErrSymbolNotFoundOnScope
}

func (scope *Scope) LookupAcrossScopes(name string) (*Node, error) {
	if node, ok := scope.Nodes[name]; ok {
		return node, nil
	}
	if scope.Parent == nil {
		return nil, ErrSymbolNotFoundOnScope
	}
	return scope.Parent.LookupAcrossScopes(name)
}

func (scope Scope) String() string {
	if scope.Parent == nil {
		return fmt.Sprintf("Scope:\nParent: nil\nCurrent: %v\n", scope.Nodes)
	}
	return fmt.Sprintf("Scope:\nParent: %v\nCurrent: %v\n", scope.Parent, scope.Nodes)
}
