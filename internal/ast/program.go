package ast

// TODO: what If I want to build a library?
// It should not be a main module because it assumes that
// main contains an entrypoint

type Program struct {
	Root *Package
}

type Package struct {
	Name     string
	Files    []*File
	Packages []*Package
	Scope    *Scope
	IsRoot   bool // If true, "Scope" represents the universe
}

type File struct {
	Dir            string
	Path           string
	Body           []*Node
	PkgName        string
	PkgNameDefined bool
}
