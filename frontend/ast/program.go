package ast

import "github.com/HicaroD/Telia/scope"

type Program struct {
	Universe *scope.Scope[Node]
	Body     []*Module
}

func (program Program) astNode() {}

type Module struct {
	Scope *scope.Scope[Node]
	Name  string
	Body  []*File
}

func (module Module) astNode() {}

type File struct {
	Scope *scope.Scope[Node]
	Name  string
	Body  []Node
}

func (file File) astNode() {}
