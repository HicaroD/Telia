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

_Static_assert(sizeof(size_t) == sizeof(uint64_t), "Telia libc bindings require 64-bit size_t");

static inline void _check_nil_pointer_deref(void *ptr) {
    if (ptr == NULL) {
        fprintf(stderr, "runtime panic: null pointer deref\n");
        exit(1);
    }
}

typedef struct { const unsigned char *data; size_t len; } _String;

static inline bool _telia_string_equal(_String a, _String b) {
    return a.len == b.len && (a.len == 0 || memcmp(a.data, b.data, a.len) == 0);
}

static inline void _telia_println(_String text) {
    if (text.len) fwrite(text.data, 1, text.len, stdout);
    fputc('\n', stdout);
}

static inline void _panic(_String msg) {
    fputs("runtime panic: ", stderr);
    if (msg.len) fwrite(msg.data, 1, msg.len, stderr);
    fputc('\n', stderr);
    exit(1);
}

typedef struct { _String msg; } _Error;

`

type CCodegen struct {
	buf          strings.Builder
	loc          *ast.Loc
	program      *ast.Program
	currentPkg   *ast.Package
	currentRetTy *ast.ExprType
	tmpCnt       int
	exePath      string
	outputPath   string
	deferScopes  [][]*ast.DeferStmt
	cFlags       []string
}

func NewCG(loc *ast.Loc, program *ast.Program) *CCodegen {
	return &CCodegen{
		loc:     loc,
		program: program,
	}
}

func (c *CCodegen) SetOutputPath(path string) {
	c.outputPath = path
}

// SetCFlags configures additional native compiler flags (e.g. sanitizers).
func (c *CCodegen) SetCFlags(flags ...string) {
	c.cFlags = append([]string(nil), flags...)
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
	c.buf.WriteString(runtimeSource)

	// Reset processed flags so each Generate() call is clean
	resetProcessed(c.program.Root)
	for _, st := range c.program.Structs {
		name := emitCType(&ast.ExprType{Kind: ast.EXPR_TYPE_STRUCT, T: &ast.StructType{Decl: st}})
		c.buf.WriteString(fmt.Sprintf("typedef struct %s %s;\n", name, name))
	}
	for _, ty := range c.program.TupleTypes {
		name := tupleTypedefName(ty)
		c.buf.WriteString(fmt.Sprintf("typedef struct %s %s;\n", name, name))
	}
	for _, ty := range c.program.ValueTypes {
		if ty.Kind == ast.EXPR_TYPE_ARRAY || ty.Kind == ast.EXPR_TYPE_SLICE {
			name := emitCType(ty)
			c.buf.WriteString(fmt.Sprintf("typedef struct %s %s;\n", name, name))
		}
	}
	for _, ty := range c.program.ValueTypes {
		if ty.Kind == ast.EXPR_TYPE_STRUCT {
			c.emitStructDecl(ty.T.(*ast.StructType).Decl)
		} else if ty.Kind == ast.EXPR_TYPE_ARRAY || ty.Kind == ast.EXPR_TYPE_SLICE {
			t := ty.T.(*ast.SequenceType)
			name := emitCType(ty)
			if ty.Kind == ast.EXPR_TYPE_SLICE {
				c.buf.WriteString(fmt.Sprintf("struct %s { %s *data; uint64_t len; };\n", name, emitCType(t.Element)))
			} else if t.Length == 0 {
				c.buf.WriteString(fmt.Sprintf("struct %s { unsigned char _empty; };\n", name))
			} else {
				c.buf.WriteString(fmt.Sprintf("struct %s { %s data[%d]; };\n", name, emitCType(t.Element), t.Length))
			}
		} else {
			c.emitTupleTypedef(ty)
		}
	}

	// Two-pass emission: declarations then bodies, DFS over import graph
	for _, fn := range c.program.Functions {
		c.emitFnForwardDecl(fn)
	}
	c.generatePackage(c.program.Root)
	for _, fn := range c.program.Functions {
		c.emitFnBody(fn)
	}

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

	// Compile into temp dir
	args := append([]string{optFlag}, c.cFlags...)
	args = append(args, "-o", exePath, cFilePath, "-lm")
	cmd := exec.Command(compiler, args...)
	if out, err := cmd.CombinedOutput(); err != nil {
		return fmt.Errorf("C compiler error:\n%s", string(out))
	}

	// Determine output path
	var outputPath string
	if c.outputPath != "" {
		if filepath.IsAbs(c.outputPath) {
			outputPath = c.outputPath
		} else {
			cwd, err := os.Getwd()
			if err != nil {
				return fmt.Errorf("failed to get working directory: %w", err)
			}
			outputPath = filepath.Join(cwd, c.outputPath)
		}
	} else {
		cwd, err := os.Getwd()
		if err != nil {
			return fmt.Errorf("failed to get working directory: %w", err)
		}
		outputPath = filepath.Join(cwd, filenameNoExt)
	}

	// Ensure output directory exists
	outDir := filepath.Dir(outputPath)
	if err := os.MkdirAll(outDir, 0755); err != nil {
		return fmt.Errorf("failed to create output directory %s: %w", outDir, err)
	}

	// Move executable to output path
	if err := os.Rename(exePath, outputPath); err != nil {
		// Cross-device rename fallback: copy then remove
		data, readErr := os.ReadFile(exePath)
		if readErr != nil {
			return fmt.Errorf("failed to read compiled binary: %w", readErr)
		}
		if writeErr := os.WriteFile(outputPath, data, 0755); writeErr != nil {
			return fmt.Errorf("failed to write binary to %s: %w", outputPath, writeErr)
		}
		os.Remove(exePath)
	}

	c.exePath = outputPath

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
			fn := node.Node.(*ast.FnDecl)
			if len(fn.TypeParams) == 0 {
				c.emitFnForwardDecl(fn)
			}
		}
	}
}

// emitBodies emits function bodies for a file.
func (c *CCodegen) emitBodies(file *ast.File) {
	for _, node := range file.Body {
		if node.Kind == ast.KIND_FN_DECL {
			fn := node.Node.(*ast.FnDecl)
			if len(fn.TypeParams) == 0 {
				c.emitFnBody(fn)
			}
		}
	}
}

// mangledName returns the C-safe name for a Telia function.
// Functions in the root (main) package keep their bare name.
// All other packages get <pkg>__<fn>.
func (c *CCodegen) mangledName(fn *ast.FnDecl) string {
	if fn.CName != "" {
		return fn.CName
	}
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
