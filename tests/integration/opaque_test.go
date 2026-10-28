package integration

import (
	"fmt"
	"os"
	"path/filepath"
	"testing"

	"github.com/HicaroD/Telia/tests/compiler"
)

func TestHeaderOwnedOpaquePointersUseTheActualCType(t *testing.T) {
	header, err := filepath.Abs("testdata/opaque.h")
	if err != nil {
		t.Fatal(err)
	}
	fixture, err := filepath.Abs("testdata/opaque.c")
	if err != nil {
		t.Fatal(err)
	}
	path := filepath.Join(t.TempDir(), "main.t")
	source := fmt.Sprintf(`package main
@[header=%q, default_cc="c"]
extern {
@[link_name="FixtureWindow"]
opaque Window
@[link_name="fixture_window_create"]
fn create(value u32) *Window
@[link_name="fixture_window_value"]
fn value(window *Window) u32
@[link_name="fixture_window_destroy"]
fn destroy(window *Window)
}
fn main() i32 {
window := create(42)
if value(window) != 42 {
return 1
}
destroy(window)
if value(window) != 0 {
return 2
}
return 0
}
`, header)
	if err := os.WriteFile(path, []byte(source), 0600); err != nil {
		t.Fatal(err)
	}
	exe, diags := compiler.CompileOnlyWithCFlags(
		path,
		fixture,
		"-Werror=incompatible-pointer-types",
		"-Werror=implicit-function-declaration",
	)
	if len(diags.Diags) != 0 {
		t.Fatalf("header-owned opaque pointer: %v", diags.Diags)
	}
	defer removeTestBinary(t, exe)
	if output, err := compiler.RunBinary(exe); err != nil {
		t.Fatalf("opaque C fixture: %v: %s", err, output)
	}
}
