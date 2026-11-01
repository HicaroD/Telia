package integration

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/HicaroD/Telia/tests/compiler"
)

func TestImportedInferenceUsesDeclarationAndCallerContexts(t *testing.T) {
	dir := t.TempDir()
	if err := os.Mkdir(filepath.Join(dir, "helpers"), 0700); err != nil {
		t.Fatal(err)
	}
	files := map[string]string{
		"helpers/helpers.t": `package helpers
struct Box<T> {
value T
}
enum Option<T> {
None
Some {
value T
}
}
fn empty<T>() T {
return T{}
}
fn identity<T>(value T) T {
return value
}
fn unbox<T>(value Box<T>) T {
return value.value
}
`,
		"main.t": `package main
use "pkg::helpers"
fn transform<A, B>(callback fn(value A) B, value A) B {
return callback(value)
}
fn forward<T>(value T) T {
return transform(helpers::identity, value)
}
fn main() i32 {
number i32 := 42
box := helpers::Box{value: number}
value := helpers::unbox(box)
zero helpers::Box<i32> := helpers::empty()
none helpers::Option<i32> := helpers::Option::None
callback fn(value i32) i32 := helpers::identity
explicit := helpers::identity<i32>
transformed := transform(helpers::identity, number)
if callback(explicit(value)) != 42 or zero.value != 0 or transformed != 42 or forward(number) != 42 {
return 1
}
return match none {
helpers::Option::None => 0
helpers::Option::Some{value} => 1
}
}
`,
	}
	for name, text := range files {
		if err := os.WriteFile(filepath.Join(dir, name), []byte(text), 0600); err != nil {
			t.Fatal(err)
		}
	}
	output, diags := compiler.CompilePackage(dir)
	if diags.HasErrors() {
		t.Fatalf("compile/run: %v; output: %s", diags.Diags, output)
	}
}

func TestGenericInferenceExecutesConcreteSpecializations(t *testing.T) {
	path := filepath.Join(t.TempDir(), "main.t")
	err := os.WriteFile(path, []byte(`package main
struct Box<T> {
value T
}
enum Option<T> {
None
Some {
value T
}
}
fn identity<T>(value T) T {
return value
}
fn apply<T>(callback fn(value T) T, value T) T {
return callback(value)
}
fn empty<T>() T {
return T{}
}
fn main() i32 {
number i32 := 42
box := Box<_>{value: number}
item := Option::Some{value: box.value}
callback fn(value i32) i32 := identity<_>
explicit := identity<i32>
value := match item {
Option::None => 0
Option::Some{value} => value
}
zero Box<i32> := empty()
if apply(identity, callback(explicit(value))) != 42 or zero.value != 0 {
return 1
}
return 0
}
`), 0600)
	if err != nil {
		t.Fatal(err)
	}
	output, diags := compiler.CompileFile(path)
	if diags.HasErrors() {
		t.Fatalf("compile/run: %v; output: %s", diags.Diags, output)
	}
}
