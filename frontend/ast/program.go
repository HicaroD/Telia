package ast

type Program struct {
	Universe *Scope
	Body     []*Module
}

func (program Program) astNode() {}

type Module struct {
	Scope *Scope
	Name  string
	Body  []*File
}

func (module Module) astNode() {}

type File struct {
	Scope *Scope
	Name  string
	Body  []Node
}

func (file File) astNode() {}
