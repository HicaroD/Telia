package scope

import (
	"errors"
)

var (
	SYMBOL_ALREADY_DEFINED_ON_SCOPE error = errors.New("symbol already defined on scope")
	SYMBOL_NOT_FOUND_ON_SCOPE       error = errors.New("symbol not found on scope")
)

type Scope[V any] struct {
	parent *Scope[V]
	nodes  map[string]V
}

func New[V any](parent *Scope[V]) *Scope[V] {
	return &Scope[V]{parent: parent, nodes: map[string]V{}}
}

func (scope *Scope[V]) Insert(name string, element V) error {
	if _, ok := scope.nodes[name]; ok {
		return SYMBOL_ALREADY_DEFINED_ON_SCOPE
	}
	scope.nodes[name] = element
	return nil
}

func (scope *Scope[V]) Lookup(name string) (V, error) {
	if node, ok := scope.nodes[name]; ok {
		return node, nil
	}
	if scope.parent == nil {
		// HACK
		var empty V
		return empty, SYMBOL_NOT_FOUND_ON_SCOPE
	}
	return scope.parent.Lookup(name)
}
