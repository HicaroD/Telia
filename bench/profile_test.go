package main

import (
	"fmt"
	"os"
	"path/filepath"
	"testing"

	"github.com/HicaroD/Telia/internal/ast"
)

func TestGeneratedWorkloadPassesFrontend(t *testing.T) {
	for _, n := range []int{1, 100, 1000} {
		t.Run(fmt.Sprint(n), func(t *testing.T) {
			path := filepath.Join(t.TempDir(), "scaling.t")
			if err := os.WriteFile(path, workload(n), 0600); err != nil {
				t.Fatal(err)
			}
			program, loc, parseMS, semaMS, err := frontend(path)
			if err != nil {
				t.Fatal(err)
			}
			if program == nil || loc == nil || parseMS <= 0 || semaMS <= 0 {
				t.Fatal("missing program or phase measurements")
			}
			functions := 0
			for _, node := range program.Root.Files[0].Body {
				if node.Kind == ast.KIND_FN_DECL {
					functions++
				}
			}
			if functions != n+1 {
				t.Fatalf("got %d functions, want %d", functions, n+1)
			}
		})
	}
}

func TestFrontendMissingSource(t *testing.T) {
	if _, _, _, _, err := frontend(filepath.Join(t.TempDir(), "missing.t")); err == nil {
		t.Fatal("expected missing source to fail")
	}
}
