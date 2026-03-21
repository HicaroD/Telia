package c

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/HicaroD/Telia/config"
	"github.com/HicaroD/Telia/internal/ast"
)

const preamble = `#include <stdint.h>
#include <stdbool.h>
#include <string.h>
#include <stdio.h>
#include <stdlib.h>

static inline void _check_nil_pointer_deref(void *ptr) {
    if (ptr == NULL) {
        fprintf(stderr, "runtime panic: null pointer deref\n");
        exit(1);
    }
}

typedef struct { char *msg; } _Error;

`

type CCodegen struct {
	buf          strings.Builder
	loc          *ast.Loc
	program      *ast.Program
	currentPkg   *ast.Package
	currentRetTy *ast.ExprType
	tmpCnt       int
	exePath      string
}

func NewCG(loc *ast.Loc, program *ast.Program) *CCodegen {
	return &CCodegen{
		loc:     loc,
		program: program,
	}
}

func (c *CCodegen) Generate(buildType config.BuildOptimizationType) error {
	compiler, err := detectCCompiler()
	if err != nil {
		return err
	}

	dir, err := os.MkdirTemp("", "build")
	if err != nil {
		return err
	}

	filenameNoExt := strings.TrimSuffix(filepath.Base(c.loc.Name), filepath.Ext(c.loc.Name))
	cFilePath := filepath.Join(dir, filenameNoExt+".c")
	exePath := filepath.Join(dir, filenameNoExt)

	// Emit into buf
	c.buf.Reset()
	c.buf.WriteString(preamble)

	// Reset processed flags so each Generate() call is clean
	resetProcessed(c.program.Root)

	// Emit tuple typedef structs before any function declarations that
	// reference them. program.TupleTypes was populated and deduplicated by
	// sema.Check(), so we emit exactly the right set, exactly once.
	for _, ty := range c.program.TupleTypes {
		c.emitTupleTypedef(ty)
	}

	// Two-pass emission: declarations then bodies, DFS over import graph
	c.generatePackage(c.program.Root)

	// Write .c file
	cFile, err := os.Create(cFilePath)
	if err != nil {
		return err
	}
	_, err = cFile.WriteString(c.buf.String())
	cFile.Close()
	if err != nil {
		return err
	}

	// Determine optimisation level
	var optFlag string
	switch buildType {
	case config.BUILD_OPT_RELEASE:
		optFlag = "-O3"
	case config.BUILD_OPT_DEBUG:
		optFlag = "-O0"
	default:
		return fmt.Errorf("unknown build type: %s", buildType)
	}

	// Compile
	cmd := exec.Command(compiler, optFlag, "-o", exePath, cFilePath, "-lm")
	if out, err := cmd.CombinedOutput(); err != nil {
		return fmt.Errorf("C compiler error:\n%s", string(out))
	}

	c.exePath = exePath

	if !config.DEV {
		if err := os.RemoveAll(dir); err != nil {
			return err
		}
	}

	return nil
}

func (c *CCodegen) ExePath() string {
	return c.exePath
}

// generatePackage walks the import graph depth-first and emits all declarations
// then all bodies for each package, deduplicating via pkg.Processed.
func (c *CCodegen) generatePackage(pkg *ast.Package) {
	if pkg == nil || pkg.Processed {
		return
	}

	// Save and restore currentPkg so nested calls don't clobber it.
	prev := c.currentPkg
	c.currentPkg = pkg
	defer func() { c.currentPkg = prev }()

	// First pass: recurse into imports, then emit forward declarations.
	for _, file := range pkg.Files {
		for _, imp := range file.Imports {
			c.generatePackage(imp.Package)
		}
		c.emitDeclarations(file)
	}

	// Second pass: emit function bodies.
	for _, file := range pkg.Files {
		c.emitBodies(file)
	}

	pkg.Processed = true
}

// emitDeclarations emits extern proto forward-decls and fn forward-decls for a file.
func (c *CCodegen) emitDeclarations(file *ast.File) {
	for _, node := range file.Body {
		switch node.Kind {
		case ast.KIND_EXTERN_DECL:
			c.emitExternDecl(node.Node.(*ast.ExternDecl))
		case ast.KIND_FN_DECL:
			c.emitFnForwardDecl(node.Node.(*ast.FnDecl))
		case ast.KIND_STRUCT_DECL:
			c.emitStructDecl(node.Node.(*ast.StructDecl))
		}
	}
}

// emitBodies emits function bodies for a file.
func (c *CCodegen) emitBodies(file *ast.File) {
	for _, node := range file.Body {
		if node.Kind == ast.KIND_FN_DECL {
			c.emitFnBody(node.Node.(*ast.FnDecl))
		}
	}
}

// mangledName returns the C-safe name for a Telia function.
// Functions in the root (main) package keep their bare name.
// All other packages get <pkg>__<fn>.
func (c *CCodegen) mangledName(fn *ast.FnDecl) string {
	if c.currentPkg == c.program.Root {
		return fn.Name.Name()
	}
	return c.currentPkg.Loc.Name + "__" + fn.Name.Name()
}

// nextTmp returns a unique temporary variable name.
func (c *CCodegen) nextTmp() string {
	name := fmt.Sprintf("_t%d", c.tmpCnt)
	c.tmpCnt++
	return name
}

// resetProcessed recursively clears pkg.Processed so Generate can be called again.
func resetProcessed(pkg *ast.Package) {
	if pkg == nil || !pkg.Processed {
		return
	}
	pkg.Processed = false
	for _, file := range pkg.Files {
		for _, imp := range file.Imports {
			resetProcessed(imp.Package)
		}
	}
}

// detectCCompiler tries cc, gcc, then clang in order and returns the path
// of the first one found on PATH. Returns an error if none are available.
func detectCCompiler() (string, error) {
	for _, candidate := range []string{"cc", "gcc", "clang"} {
		if path, err := exec.LookPath(candidate); err == nil {
			return path, nil
		}
	}
	return "", fmt.Errorf("no C compiler found on PATH; install gcc or clang")
}
