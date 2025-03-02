package ast

import (
	"fmt"
	"os"
	"path/filepath"
)

// TODO: what If I want to build a library?
// It should not be a main package because it assumes that
// main contains an entrypoint

type Program struct {
	Root *Package
}

type Loc struct {
	Name      string
	Dir       string
	Path      string
	IsPackage bool
}

func LocFromPath(path string) (*Loc, error) {
	loc := new(Loc)

	fullPath, err := filepath.Abs(path)
	// TODO(errors)
	if err != nil {
		return nil, err
	}
	loc.Path = fullPath

	info, err := os.Stat(fullPath)
	// TODO(errors)
	if err != nil {
		return nil, err
	}

	mode := info.Mode()
	if mode.IsDir() {
		loc.Name = filepath.Base(fullPath)
		loc.Dir = filepath.Base(fullPath)
		loc.IsPackage = true
	} else {
		loc.Name = filepath.Base(path)
		loc.Dir = filepath.Base(filepath.Dir(path))
		loc.IsPackage = false
	}

	return loc, nil
}

func (l Loc) String() string {
	return fmt.Sprintf("Name: %s | Dir: %s | Path: %s | isPackage: %v", l.Name, l.Dir, l.Path, l.IsPackage)
}

type Package struct {
	Loc      *Loc
	Files    []*File
	Packages []*Package
	Scope    *Scope
	IsRoot   bool // If true, "Scope" represents the universe
}

func (p *Package) String() string {
	if p.Loc.Name == "" {
		return "Name: <EMPTY>"
	}
	return fmt.Sprintf("Name: %s", p.Loc.Name)
}

type File struct {
	Loc            *Loc
	PkgName        string
	PkgNameDefined bool
	Body           []*Node
}
