package ast

type Program struct {
	Body []*Module
}

func (program Program) astNode() {}

type Module struct {
	Name string
	Body []*File
}

func (module Module) astNode() {}

type File struct {
	Name string
	Body []Node
}

func (file File) astNode() {}
