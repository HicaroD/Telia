package scope

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
type Scope[V any] struct {
	Parent *Scope[V]
	Nodes  map[string]V
}

func New[V any](parent *Scope[V]) *Scope[V] {
	return &Scope[V]{Parent: parent, Nodes: map[string]V{}}
}

func (scope *Scope[V]) Insert(name string, element V) error {
	if _, ok := scope.Nodes[name]; ok {
		return ERR_SYMBOL_ALREADY_DEFINED_ON_SCOPE
	}
	scope.Nodes[name] = element
	return nil
}

func (scope *Scope[V]) LookupCurrentScope(name string) (V, error) {
	if node, ok := scope.Nodes[name]; ok {
		return node, nil
	}
	var empty V
	return empty, ERR_SYMBOL_NOT_FOUND_ON_SCOPE
}

func (scope *Scope[V]) LookupAcrossScopes(name string) (V, error) {
	if node, ok := scope.Nodes[name]; ok {
		return node, nil
	}
	if scope.Parent == nil {
		var empty V
		return empty, ERR_SYMBOL_NOT_FOUND_ON_SCOPE
	}
	return scope.Parent.LookupAcrossScopes(name)
}

func (scope Scope[V]) String() string {
	if scope.Parent == nil {
		return fmt.Sprintf("Scope:\nParent: nil\nCurrent: %v\n", scope.Nodes)
	}
	return fmt.Sprintf("Scope:\nParent: %v\nCurrent: %v\n", scope.Parent, scope.Nodes)
}
