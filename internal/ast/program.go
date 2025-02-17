package ast

// TODO: what If I want to build a library?
// It should not be a main module because it assumes that
// main contains an entrypoint

type Program struct {
	Root *Package
}

func (program Program) astNode() {}

type Package struct {
	Name     string
	Files    []*File
	Packages []*Package
	Scope    *Scope
	IsRoot   bool // If true, "Scope" represents the universe
}

func (module Package) astNode() {}

type File struct {
	Dir                string
	Path               string
	Body               []Decl
	PackageNameDefined bool
}

func (file File) astNode() {}
