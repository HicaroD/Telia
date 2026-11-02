package integration

import (
	"context"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"testing"
	"time"

	"github.com/HicaroD/Telia/internal/ast"
	"github.com/HicaroD/Telia/internal/diagnostics"
	"github.com/HicaroD/Telia/internal/parser"
	"github.com/HicaroD/Telia/internal/sema"
	"github.com/HicaroD/Telia/tests/compiler"
)

func TestStructuralEquality(t *testing.T) {
	include, err := filepath.Abs("testdata")
	if err != nil {
		t.Fatal(err)
	}
	for _, opt := range []string{"-O0", "-O3"} {
		t.Run(opt, func(t *testing.T) {
			path, diags := compiler.CompileOnlyWithCFlags("testdata/equality.t", opt, "-I", include)
			if diags.HasErrors() {
				t.Fatalf("equality compilation: %v", diags.Diags)
			}
			defer removeTestBinary(t, path)
			if output, err := compiler.RunBinary(path); err != nil {
				t.Fatalf("equality execution: %v: %s", err, output)
			}
		})
	}
}

func TestKeyEligibilityExcludesOnlyStandardOwnerDeclarations(t *testing.T) {
	for _, tc := range []struct{ pkg, name string }{
		{"strings", "Owned"}, {"strings", "COwned"}, {"memory", "Arena"},
		{"lists", "List<i32>"}, {"fs", "File"}, {"process", "Args"},
		{"process", "Child"}, {"terminal", "Handle"}, {"terminal", "SavedMode"},
	} {
		t.Run(tc.pkg+"/"+tc.name, func(t *testing.T) {
			source := "package main\nuse \"std::" + tc.pkg + "\"\nstruct Envelope {\n owner " + tc.pkg + "::" + tc.name + "\n}\nfn main() {}\n"
			p := equalityProgram(t, source)
			ty := equalityStruct(t, p, "Envelope")
			err := ast.EqualityEligibility(ty, true)
			if err == nil || !strings.Contains(err.Error(), "value.owner") ||
				!strings.Contains(err.Error(), "owning handle") {
				t.Fatalf("want nested owner key diagnostic, got %v", err)
			}
			if tc.pkg != "lists" {
				if err := ast.EqualityEligibility(ty, false); err != nil {
					t.Fatalf("owner equality: %v", err)
				}
			}
		})
	}
	p := equalityProgramFiles(t, map[string]string{
		"main.t":            "package main\nuse \"pkg::strings\"\nstruct Envelope {\nowner strings::Owned\n}\nfn main() {}\n",
		"strings/strings.t": "package strings\nstruct Owned {\ntext string\n}\n",
	})
	if err := ast.EqualityEligibility(equalityStruct(t, p, "Envelope"), true); err != nil {
		t.Fatalf("user package and declaration lookalike rejected: %v", err)
	}
}

func TestKeyEligibilityResolvesAliasesGenericsAndAllEnumAlternatives(t *testing.T) {
	p := equalityProgram(t, `package main
use "std::strings"
type Owner = strings::Owned
struct Box<T> {
 value T
}
struct Nested {
 owners Box<Owner[2]>
}
struct Borrowed {
 owner *Owner
}
enum OptionalCallback {
 None
 Some {
  callback fn() i32
 }
}
fn main() {}
`)
	nested := equalityStruct(t, p, "Nested")
	if err := ast.EqualityEligibility(nested, true); err == nil ||
		!strings.Contains(err.Error(), "value.owners.value[element]") {
		t.Fatalf("nested owner: %v", err)
	}
	if err := ast.EqualityEligibility(nested, false); err != nil {
		t.Fatalf("nested owner equality: %v", err)
	}
	if err := ast.EqualityEligibility(equalityStruct(t, p, "Borrowed"), true); err != nil {
		t.Fatalf("pointer to owner: %v", err)
	}
	found := false
	for _, d := range p.Enums {
		if d.Name.Name() == "OptionalCallback" {
			found = true
			ty := &ast.ExprType{Kind: ast.EXPR_TYPE_ENUM, T: &ast.EnumType{Decl: d}}
			if err := ast.EqualityEligibility(ty, true); err == nil ||
				!strings.Contains(err.Error(), "value::Some.callback") {
				t.Fatalf("inactive callback eligibility: %v", err)
			}
			if err := ast.EqualityEligibility(ty, false); err != nil {
				t.Fatalf("callback enum equality: %v", err)
			}
		}
	}
	if !found {
		t.Fatal("missing OptionalCallback enum")
	}
}

func TestStructuralEqualityWithNativeSanitizers(t *testing.T) {
	if runtime.GOOS == "windows" {
		t.Skip("native sanitizer runtime is not configured on Windows")
	}
	include, err := filepath.Abs("testdata")
	if err != nil {
		t.Fatal(err)
	}
	sanitizers := os.Getenv("TELIA_TEST_SANITIZERS")
	if sanitizers == "" {
		sanitizers = "undefined"
	}
	for _, opt := range []string{"-O0", "-O3"} {
		t.Run(opt, func(t *testing.T) {
			path, diags := compiler.CompileOnlyWithCFlags(
				"testdata/equality.t",
				opt,
				"-I",
				include,
				"-fsanitize="+sanitizers,
				"-fno-sanitize-recover=all",
				"-fno-omit-frame-pointer",
			)
			if diags.HasErrors() {
				t.Fatalf("sanitized equality compilation: %v", diags.Diags)
			}
			defer removeTestBinary(t, path)
			ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
			defer cancel()
			if output, err := exec.CommandContext(ctx, path).CombinedOutput(); err != nil {
				t.Fatalf(
					"sanitized equality execution: %v (context: %v): %s",
					err,
					ctx.Err(),
					output,
				)
			}
		})
	}
}

func equalityProgram(t *testing.T, text string) *ast.Program {
	t.Helper()
	return equalityProgramFiles(t, map[string]string{"main.t": text})
}

func equalityProgramFiles(t *testing.T, files map[string]string) *ast.Program {
	t.Helper()
	dir := t.TempDir()
	for name, text := range files {
		path := filepath.Join(dir, name)
		if err := os.MkdirAll(filepath.Dir(path), 0700); err != nil {
			t.Fatal(err)
		}
		if err := os.WriteFile(path, []byte(text), 0600); err != nil {
			t.Fatal(err)
		}
	}
	path := filepath.Join(dir, "main.t")
	loc, err := ast.LocFromPath(path)
	if err != nil {
		t.Fatal(err)
	}
	diags := diagnostics.New()
	p, err := parser.New(diags).ParseFileAsProgram(dir, loc, diags)
	if err != nil {
		t.Fatalf("parse: %v: %v", err, diags.Diags)
	}
	if err := sema.New(diags).Check(p); err != nil || diags.HasErrors() {
		t.Fatalf("check: %v: %v", err, diags.Diags)
	}
	return p
}

func equalityStruct(t *testing.T, p *ast.Program, name string) *ast.ExprType {
	t.Helper()
	for _, d := range p.Structs {
		if d.Name.Name() == name {
			return &ast.ExprType{Kind: ast.EXPR_TYPE_STRUCT, T: &ast.StructType{Decl: d}}
		}
	}
	t.Fatalf("missing %s struct", name)
	return nil
}
