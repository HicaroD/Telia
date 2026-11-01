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

func TestContainerMultiResultReassignment(t *testing.T) {
	runStandardLibraryFixture(t, "container_reassignment.t")
}

func TestMapsHeapLifecycle(t *testing.T) {
	output, diags := compiler.CompileFile("testdata/maps.t")
	if diags.HasErrors() {
		t.Fatalf("map fixture should compile: %v", diags.Diags)
	}
	if output != "maps ok\n" {
		t.Fatalf("unexpected output: %q", output)
	}
}

func TestMapsExampleUsesPublishedAPIs(t *testing.T) {
	fixture, err := filepath.Abs("../../examples/maps.t")
	if err != nil {
		t.Fatal(err)
	}
	header, err := filepath.Abs("testdata/allocation_control.h")
	if err != nil {
		t.Fatal(err)
	}
	runStandardLibraryFixture(t, fixture, "-include", header)
}

func TestMapsRejectInvalidKeysAtInstantiation(t *testing.T) {
	for _, tc := range []struct{ key, want string }{
		{"i32[]", "slice"},
		{"Box<i32[]>", "value.value"},
		{"Choice<i32[]>", "item"},
		{"strings::Owned", "owning handle"},
		{"maps::Map<i32, i32>", "owning handle"},
		{"error", "error"},
		{"fn()", "function"},
	} {
		t.Run(tc.key, func(t *testing.T) {
			file := filepath.Join(t.TempDir(), "main.t")
			source := `package main
use "std::maps"
use "std::strings"
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
 invalid := maps::create<` + tc.key + `, i32>()
}
`
			if err := os.WriteFile(file, []byte(source), 0600); err != nil {
				t.Fatal(err)
			}
			path, diags := compiler.CompileOnly(file)
			if !diags.HasErrors() {
				removeTestBinary(t, path)
				t.Fatal("invalid map key compiled")
			}
			for _, diag := range diags.Diags {
				if strings.Contains(diag.Message, tc.want) &&
					strings.Contains(diag.Message, "key") {
					return
				}
			}
			t.Fatalf("want key diagnostic containing %q, got %v", tc.want, diags.Diags)
		})
	}
}

func TestMapsNative(t *testing.T) {
	header, err := filepath.Abs("testdata/maps_control.h")
	if err != nil {
		t.Fatal(err)
	}
	for _, fixture := range []string{"maps.t", "maps_model.t", "maps_failure.t", "maps_keys.t", "maps_ownership.t", "maps_checked.t", "maps_clone.t", "maps_arena.t", "container_reassignment.t", "containers_model.t"} {
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
							t.Fatalf("map compilation: %v", diags.Diags)
						}
						defer removeTestBinary(t, path)
						ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
						defer cancel()
						if out, err := exec.CommandContext(ctx, path).CombinedOutput(); err != nil {
							t.Fatalf("map execution: %v: %s", err, out)
						}
					},
				)
			}
		}
	}
}
