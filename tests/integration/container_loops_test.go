package integration

import (
	"context"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"testing"
	"time"

	"github.com/HicaroD/Telia/tests/compiler"
)

func TestContainerLoopsDiagnostics(t *testing.T) {
	for _, tc := range []struct{ decl, body, want string }{
		{"", "m := maps::create<i32, i32>()\nfor key in m {}", "requires 2 binding(s)"},
		{"", "s := sets::create<i32>()\nfor key, value in s {}", "requires 1 binding(s)"},
		{"", "a i32[1] := .{1}\nfor key, value in a {}", "requires 1 binding(s)"},
		{"", "s := sets::create<i32>()\nfor a, b, c in s {}", "requires 1 binding(s)"},
		{"", "for key, value in maps::create<i32, i32>() {}", "addressable live storage"},
		{"", "for value in sets::create<i32>() {}", "addressable live storage"},
		{"struct Box {\nvalues sets::Set<i32>\n}\nfn temporary() Box {\nreturn Box{}\n}", "for value in temporary().values {}", "addressable live storage"},
		{"struct Map<K, V> {\nkey K\nvalue V\n}", "m := Map<i32, i32>{}\nfor key, value in m {}", "iteration requires an array"},
		{"struct Set<T> {\nvalue T\n}", "s := Set<i32>{}\nfor value in s {}", "iteration requires an array"},
		{"", "m := maps::create<i32, i32>()\nfor key, key in m {}", "symbol already defined"},
		{"", "m := maps::create<i32, i32>()\nfor key, value in &m {}", "iteration requires an array"},
		{"fn read(m const *maps::Map<i32, i32>) {\nfor key, value in *m {}\n}", "", "mutable live storage"},
	} {
		t.Run(tc.body, func(t *testing.T) {
			path := loopControlSource(
				t,
				"package main\nuse \"std::maps\"\nuse \"std::sets\"\n"+tc.decl+"\nfn main() {\n"+tc.body+"\n}\n",
			)
			_, diags := compiler.CompileOnly(path)
			if !diags.HasErrors() || !strings.Contains(diags.Diags[0].Message, tc.want) ||
				!diags.Diags[0].Span.Valid() {
				t.Fatalf("want located %q, got %v", tc.want, diags.Diags)
			}
		})
	}
}

func TestContainerLoopsNative(t *testing.T) {
	header, err := filepath.Abs("testdata/maps_control.h")
	if err != nil {
		t.Fatal(err)
	}
	for _, opt := range []string{"-O0", "-O3"} {
		t.Run(opt, func(t *testing.T) {
			flags := []string{opt, "-include", header}
			if runtime.GOOS != "windows" {
				flags = append(flags, "-fsanitize=undefined", "-fno-sanitize-recover=all")
			}
			exe, diags := compiler.CompileOnlyWithCFlags("testdata/container_loops.t", flags...)
			if diags.HasErrors() {
				t.Fatalf("native loops: %v", diags.Diags)
			}
			defer removeTestBinary(t, exe)
			ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
			defer cancel()
			if output, err := exec.CommandContext(ctx, exe).CombinedOutput(); err != nil {
				t.Fatalf("execution: %v: %s", err, output)
			}
		})
	}
}

func TestContainerLoopsGenericSpecializations(t *testing.T) {
	path := loopControlSource(t, `package main
use "std::maps"
use "std::sets"
fn size<K, V>(map *maps::Map<K, V>, set *sets::Set<K>) i32 {
count i32 := 0
for key, value in *map {
copy := value
for element in *set {
if key == element {
count = count + 1
}
}
}
return count
}
fn run() error {
map := maps::create<i32, string>()
defer maps::free(&map)
set := sets::create<i32>()
defer sets::free(&set)
maps::insert(&map, 7, "seven") @try
sets::insert(&set, 7) @try
other := maps::create<string, i32>()
defer maps::free(&other)
names := sets::create<string>()
defer sets::free(&names)
maps::insert(&other, "seven", 7) @try
sets::insert(&names, "seven") @try
if size(&map, &set) != 1 or size(&other, &names) != 1 {
return error("specialization")
}
return nil
}
fn main() i32 {
run() @fail
return 0
}
`)
	runLoopControl(t, path)
}

func TestContainerLoopsVisitEntriesAndCopyBindings(t *testing.T) {
	path := loopControlSource(t, `package main
use "std::maps"
use "std::sets"
fn run() error {
scores maps::Map<i32, i32> := maps::create()
defer maps::free(&scores)
names := sets::create<i32>()
defer sets::free(&names)
for i i32 := 1; i <= 3; i = i + 1 {
maps::insert(&scores, i, i * 10) @try
sets::insert(&names, i) @try
}
total i32 := 0
for key, value in scores {
total = total + value
ptr := maps::get_ptr(&scores, key)
*ptr = *ptr + 1
key = 99
value = 99
}
sum i32 := 0
for element in names {
sum = sum + element
element = 99
}
if total != 60 or sum != 6 or maps::count(&scores) != 3 or sets::contains(&names, 99) {
return error("bad copies or visits")
}
updated i32 := 0
for key, value in scores {
updated = updated + value
}
if updated != 63 {
return error("updates not visible")
}
return nil
}
fn main() i32 {
run() @fail
return 0
}
`)
	runLoopControl(t, path)
}
