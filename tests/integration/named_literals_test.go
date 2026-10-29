package integration

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/HicaroD/Telia/tests/compiler"
)

func TestNamedLiteralsInControlFlowHeaders(t *testing.T) {
	source := `package main
struct Box<T> {
value T
}
enum Option<T> {
None
Some {
value T
}
}
fn enabled(flag Box<bool>) bool {
return flag.value
}
fn values(box Box<i32>) i32[1] {
return .{box.value}
}
fn main() i32 {
if enabled(Box<bool>{value: true}) == false {
return 1
}
result := match (Option<i32>::Some{value: 42}) {
Option<i32>::Some{value: answer} => answer
Option<i32>::None => 0
}
if result != 42 {
return 2
}
if match Option<i32>::None {
Option<i32>::None => {
flag := Box<bool>{value: true}
yield flag.value
}
Option<i32>::Some{..} => false
} == false {
return 3
}
ready := true
while ready {
if enabled(Box<bool>{value: false}) {
return 4
} elif enabled(Box<bool>{value: true}) {
ready = false
}
}
sum i32 := 0
for value in values(Box<i32>{value: 42}) {
sum = sum + value
}
if sum != 42 {
return 5
}
for i := 0; i < 1; i = result {
if i != 0 {
return 6
}
}
return 0
}
`
	path := filepath.Join(t.TempDir(), "main.t")
	if err := os.WriteFile(path, []byte(source), 0600); err != nil {
		t.Fatal(err)
	}
	if _, diags := compiler.CompileFile(path); len(diags.Diags) != 0 {
		t.Fatalf("named literals and control-flow bodies: %v", diags.Diags)
	}
}
