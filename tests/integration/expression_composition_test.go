package integration

import (
	"strings"
	"testing"

	"github.com/HicaroD/Telia/tests/compiler"
)

func TestIndexedFieldTypeMismatchProducesDiagnostic(t *testing.T) {
	path := loopControlSource(t, `package main
struct Item {
path string
}
fn main() i32 {
items Item[1] := .{Item{path: "text"}}
wrong i32 := items[0].path
return 0
}
`)
	_, diags := compiler.CompileOnly(path)
	if len(diags.Diags) == 0 || !strings.Contains(diags.Diags[0].Message, "type mismatch") {
		t.Fatalf("want field type mismatch diagnostic, got %v", diags.Diags)
	}
}

func TestIndexedFieldsComposeAsValuesAssignmentsAndPointers(t *testing.T) {
	path := loopControlSource(t, `package main
use "std::io"
struct Item {
path string
}
struct Container {
items Item[2]
}
fn source(items Item[]) Item[] {
io::println("source")
return items
}
fn pick_index() u64 {
io::println("index")
return 0
}
fn main() i32 {
items Item[2] := .{Item{path: "first"}, Item{path: "second"}}
io::println(items[0].path)
items[1].path = "changed"
pointer := &items[1].path
io::println(*pointer)
io::println(source(items[:])[pick_index()].path)
container := Container{items: items}
io::println(container.items[1].path)
return 0
}
`)
	output, diags := compiler.CompileFile(path)
	if len(diags.Diags) != 0 {
		t.Fatalf("composed field expressions should compile: %v", diags.Diags)
	}
	if want := "first\nchanged\nsource\nindex\nfirst\nchanged\n"; output != want {
		t.Fatalf("want %q, got %q", want, output)
	}
}

func TestFieldsOfTemporaryValuesAreNotAssignable(t *testing.T) {
	path := loopControlSource(t, `package main
struct Item {
path string
}
fn item() Item {
return Item{path: "temporary"}
}
fn main() i32 {
item().path = "changed"
return 0
}
`)
	_, diags := compiler.CompileOnly(path)
	if len(diags.Diags) == 0 || !strings.Contains(diags.Diags[0].Message, "addressable") {
		t.Fatalf("want non-addressable field diagnostic, got %v", diags.Diags)
	}
}

func TestIndexedFieldsComposeWithGenericReceiversPointersAndCallbacks(t *testing.T) {
	path := loopControlSource(t, `package main
use "std::io"
struct Item {
path string
callback fn() string
}
type ItemPointer = *Item
fn message() string {
return "callback"
}
fn path<T>(items T[]) string {
return items[0].path
}
fn source(items Item[]) Item[] {
return items
}
fn main() i32 {
items Item[1] := .{Item{path: "generic", callback: message}}
io::println(path<Item>(items[:]))
io::println(items[0].callback())
pointers ItemPointer[1] := .{&items[0]}
source(items[:])[0].path = "updated"
io::println(pointers[0].path)
return 0
}
`)
	output, diags := compiler.CompileFile(path)
	if len(diags.Diags) != 0 {
		t.Fatalf("generic/pointer/callback field composition should compile: %v", diags.Diags)
	}
	if want := "generic\ncallback\nupdated\n"; output != want {
		t.Fatalf("want %q, got %q", want, output)
	}
}
