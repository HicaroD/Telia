package ast

// TODO: what If I want to build a library?
// It should not be a main module because it assumes that
// main contains an entrypoint

type Program struct {
	Root *Module
}

func (program Program) astNode() {}

type Module struct {
	Name    string
	Files   []*File
	Modules []*Module
	Scope   *Scope
	IsRoot  bool // If true, "Scope" represents the universe
}

func (module Module) astNode() {}

type File struct {
	Dir  string
	Path string
	Body []Node
}

func (file File) astNode() {}
