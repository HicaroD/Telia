package integration

import (
	"testing"

	"github.com/HicaroD/Telia/tests/compiler"
)

func TestTakeTransfersOwnersAndLeavesSourcesSafeToCleanUp(t *testing.T) {
	path := loopControlSource(t, `package main
use "std::memory"
use "std::strings"
use "std::lists"
use "std::io"
fn transfer() (strings::Owned, error) {
text := strings::clone("transferred") @try
defer strings::free(&text)
return memory::take<strings::Owned>(&text), nil
}
fn run() error {
text := transfer() @try
defer strings::free(&text)
io::println(strings::view(&text))
items := lists::create<i32>()
defer lists::free<i32>(&items)
lists::append<i32>(&items, 42) @try
moved := memory::take<lists::List<i32>>(&items)
defer lists::free<i32>(&moved)
if lists::count<i32>(&items) != 0 or lists::capacity<i32>(&items) != 0 {
return error("source list not empty")
}
value := lists::get<i32>(&moved, 0) @try
if value != 42 {
return error("list contents lost")
}
lists::append<i32>(&items, 7) @try
owners strings::Owned[1] := .{strings::clone("indexed") @try}
defer strings::free(&owners[0])
indexed := memory::take<strings::Owned>(&owners[0])
defer strings::free(&indexed)
io::println(strings::view(&indexed))
arena := memory::arena_create() @try
defer memory::arena_destroy(&arena)
borrowed := strings::clone_in(&arena, "arena") @try
moved_arena := memory::take<memory::Arena>(&arena)
defer memory::arena_destroy(&moved_arena)
io::println(borrowed)
return nil
}
fn main() i32 {
run() @fail
return 0
}
`)
	exe, diags := compiler.CompileOnlyWithCFlags(
		path,
		"-fsanitize=undefined",
		"-fno-sanitize-recover=all",
	)
	if len(diags.Diags) != 0 {
		t.Fatalf("ownership transfer should compile: %v", diags.Diags)
	}
	defer removeTestBinary(t, exe)
	output, err := compiler.RunBinary(exe)
	if want := "transferred\nindexed\narena\n"; err != nil || output != want {
		t.Fatalf("want %q with safe cleanup, got %q, %v", want, output, err)
	}
}
