package integration

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"testing"
	"time"

	"github.com/HicaroD/Telia/tests/compiler"
)

func TestSetsLifecycle(t *testing.T) {
	output, diags := compiler.CompileFile("testdata/sets.t")
	if diags.HasErrors() {
		t.Fatalf("set fixture: %v", diags.Diags)
	}
	if output != "sets ok\n" {
		t.Fatalf("unexpected output: %q", output)
	}
}

func TestSetsExampleUsesPublishedAPIs(t *testing.T) {
	fixture, err := filepath.Abs("../../examples/sets.t")
	if err != nil {
		t.Fatal(err)
	}
	header, err := filepath.Abs("testdata/allocation_control.h")
	if err != nil {
		t.Fatal(err)
	}
	runStandardLibraryFixture(t, fixture, "-include", header)
}

func TestSetsNative(t *testing.T) {
	header, err := filepath.Abs("testdata/maps_control.h")
	if err != nil {
		t.Fatal(err)
	}
	for _, fixture := range []string{"sets.t", "sets_keys.t", "sets_ownership.t", "sets_model.t", "sets_failure.t", "sets_arena.t"} {
		for _, opt := range []string{"-O0", "-O3"} {
			for _, collision := range []bool{false, true} {
				t.Run(
					fmt.Sprintf("%s/%s/collision=%v", fixture, opt, collision),
					func(t *testing.T) {
						flags := []string{opt, "-include", header}
						if runtime.GOOS != "windows" {
							sanitizers := os.Getenv("TELIA_TEST_SANITIZERS")
							if sanitizers == "" {
								sanitizers = "undefined"
							}
							flags = append(
								flags,
								"-fsanitize="+sanitizers,
								"-fno-sanitize-recover=all",
							)
						}
						if collision {
							flags = append(flags, "-DTELIA_MAP_COLLISIONS=1")
						}
						path, diags := compiler.CompileOnlyWithCFlags("testdata/"+fixture, flags...)
						if diags.HasErrors() {
							t.Fatalf("set compilation: %v", diags.Diags)
						}
						defer removeTestBinary(t, path)
						ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
						defer cancel()
						if out, err := exec.CommandContext(ctx, path).CombinedOutput(); err != nil {
							t.Fatalf("set execution: %v: %s", err, out)
						}
					},
				)
			}
		}
	}
}

func TestSetsRejectInvalidElementsAtInstantiation(t *testing.T) {
	for _, tc := range []struct{ element, want string }{
		{"i32[]", "slice"},
		{"Box<i32[]>", "value.value"},
		{"Choice<i32[]>", "item"},
		{"strings::Owned", "owning handle"},
		{"sets::Set<i32>", "owning handle"},
		{"Box<strings::Owned>", "owning handle"},
		{"error", "error"},
		{"Callback", "function"},
	} {
		t.Run(tc.element, func(t *testing.T) {
			file := filepath.Join(t.TempDir(), "main.t")
			for _, expression := range []string{
				"sets::create<" + tc.element + ">()",
				"zero<sets::Set<" + tc.element + ">>()",
			} {
				source := `package main
use "std::sets"
use "std::strings"
type Callback = fn() i32
struct Box<T> {
 value T
}
enum Choice<T> {
 Empty
 Value {
  item T
 }
}
fn main() {
 invalid := ` + expression + `
}
`
				if err := os.WriteFile(file, []byte(source), 0600); err != nil {
					t.Fatal(err)
				}
				path, diags := compiler.CompileOnly(file)
				if !diags.HasErrors() {
					removeTestBinary(t, path)
					t.Fatal("invalid set element compiled")
				}
				matched := false
				for _, diag := range diags.Diags {
					if strings.Contains(diag.Message, tc.want) &&
						strings.Contains(diag.Message, "set element") {
						matched = true
					}
				}
				if !matched {
					t.Fatalf("want element diagnostic containing %q, got %v", tc.want, diags.Diags)
				}
			}
		})
	}
}
