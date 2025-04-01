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
	Root    *Package
	Runtime *Package
}

type Loc struct {
	Name      string
	Dir       string
	Path      string
	IsPackage bool
}

func LocFromPath(fullPath string) (*Loc, error) {
	loc := new(Loc)
	loc.Path = fullPath

	info, err := os.Stat(fullPath)
	// TODO(errors)
	if err != nil {
		// fmt.Println(fullPath)
		return nil, err
	}

	mode := info.Mode()
	loc.IsPackage = mode.IsDir()
	loc.Name = filepath.Base(fullPath)

	if mode.IsDir() {
		loc.Dir = filepath.Base(fullPath)
	} else {
		loc.Dir = filepath.Base(filepath.Dir(fullPath))
	}

	return loc, nil
}

func (l Loc) String() string {
	return fmt.Sprintf(
		"Name: %s | Dir: %s | Path: %s | isPackage: %v",
		l.Name,
		l.Dir,
		l.Path,
		l.IsPackage,
	)
}

type PackageType int

const (
	PACKAGE_STD = iota
	PACKAGE_RUNTIME
	PACKAGE_USER
)

type Package struct {
	Loc       *Loc
	Files     []*File
	Scope     *Scope
	IsRoot    bool // If true, "Scope" represents the universe
	Analyzed  bool
	Processed bool
}

func (p *Package) String() string {
	if p.Loc.Name == "" {
		return "Name: <EMPTY>"
	}
	return fmt.Sprintf("Name: %s", p.Loc.Name)
}

type File struct {
	Loc     *Loc
	PkgName string
	Imports map[string]*UseDecl
	Body    []*Node

	// Helper state
	PkgNameDefined   bool
	IsFirstNode      bool
	AnyDeclNodeFound bool
}
