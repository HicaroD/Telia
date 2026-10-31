package integration

import (
	"context"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/HicaroD/Telia/tests/compiler"
)

func loopControlSource(t *testing.T, source string) string {
	t.Helper()
	path := filepath.Join(t.TempDir(), "main.t")
	if err := os.WriteFile(path, []byte(source), 0600); err != nil {
		t.Fatal(err)
	}
	return path
}

func runLoopControl(t *testing.T, path string) string {
	t.Helper()
	exe, diags := compiler.CompileOnly(path)
	if len(diags.Diags) != 0 {
		t.Fatalf("loop control program: %v", diags.Diags)
	}
	defer removeTestBinary(t, exe)
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	output, err := exec.CommandContext(ctx, exe).CombinedOutput()
	if err != nil {
		t.Fatalf("execution: %v: %s", err, output)
	}
	return string(output)
}

func TestBreakExitsNearestLoop(t *testing.T) {
	path := loopControlSource(t, `package main
fn main() i32 {
count i32 := 0
for i := 0; i < 3; i = i + 1 {
while true {
count = count + 1
break
count = count + 100
}
}
if count != 3 {
return 1
}
return 0
}
`)
	runLoopControl(t, path)
}

func TestContinueRunsForUpdateAndSkipsRemainingBody(t *testing.T) {
	path := loopControlSource(t, `package main
fn main() i32 {
sum i32 := 0
for i := 0; i < 5; i = i + 1 {
if i < 3 {
continue
}
sum = sum + i
}
if sum != 7 {
return 1
}
j i32 := 0
while j < 3 {
j = j + 1
continue
sum = sum + 100
}
if sum != 7 or j != 3 {
return 2
}
return 0
}
`)
	runLoopControl(t, path)
}

func TestLoopControlRequiresEnclosingLoop(t *testing.T) {
	for _, keyword := range []string{"break", "continue"} {
		t.Run(keyword, func(t *testing.T) {
			path := loopControlSource(t, "package main\nfn main() {\n"+keyword+"\n}\n")
			_, diags := compiler.CompileOnly(path)
			if len(diags.Diags) != 1 ||
				!strings.Contains(diags.Diags[0].Message, "requires an enclosing loop") ||
				!diags.Diags[0].Span.Valid() {
				t.Fatalf("want source-located loop context diagnostic, got %v", diags.Diags)
			}
		})
	}
}

func TestLoopControlUnwindsOnlyExitedScopesInLIFOOrder(t *testing.T) {
	path := loopControlSource(t, `package main
use "std::io"
fn main() i32 {
defer io::println("outer")
for i := 0; i < 2; i = i + 1 {
defer io::println("body first")
defer io::println("body last")
if i == 0 {
defer io::println("continue")
continue
}
{
defer io::println("break")
break
}
defer io::println("unreached")
}
io::println("after")
return 0
}
`)
	output := runLoopControl(t, path)
	if want := "continue\nbody last\nbody first\nbreak\nbody last\nbody first\nafter\nouter\n"; output != want {
		t.Fatalf("want %q, got %q", want, output)
	}
}

func TestDeferredCleanupCannotJumpToEnclosingLoop(t *testing.T) {
	for _, keyword := range []string{"break", "continue"} {
		t.Run(keyword, func(t *testing.T) {
			path := loopControlSource(
				t,
				"package main\nfn main() {\nwhile true {\ndefer "+keyword+"\nbreak\n}\n}\n",
			)
			_, diags := compiler.CompileOnly(path)
			if len(diags.Diags) != 1 ||
				!strings.Contains(diags.Diags[0].Message, "requires an enclosing loop") {
				t.Fatalf("want deferred jump diagnostic, got %v", diags.Diags)
			}
		})
	}
}

func TestLoopControlInMatchArmsUnwindsAndTargetsLoop(t *testing.T) {
	path := loopControlSource(t, `package main
use "std::io"
fn main() i32 {
for i := 0; i < 4; i = i + 1 {
defer io::println("body")
value := match i {
0 => {
defer io::println("skip")
continue
}
2 => {
defer io::println("stop")
break
}
_ => 42
}
if value != 42 {
return 1
}
io::println("value")
}
return 0
}
`)
	output := runLoopControl(t, path)
	if want := "skip\nbody\nvalue\nbody\nstop\nbody\n"; output != want {
		t.Fatalf("want %q, got %q", want, output)
	}
}

func TestLoopControlInGenericArrayAndSliceIteration(t *testing.T) {
	path := loopControlSource(t, `package main
fn sum<T>(values T[], skipped T, stop T) T {
total := zero<T>()
for value in values {
if value == skipped {
continue
}
if value == stop {
break
}
total = total + value
}
return total
}
fn main() i32 {
values i32[5] := .{1, 2, 3, 4, 5}
total i32 := 0
for value in values {
if value == 2 {
continue
}
if value == 4 {
break
}
total = total + value
}
wide u64[5] := .{1, 2, 3, 4, 5}
if total != 4 or sum<i32>(values[:], 2, 4) != 4 or sum<u64>(wide[:], 2, 4) != 4 {
return 1
}
return 0
}
`)
	runLoopControl(t, path)
}

func TestDeferredCleanupCanControlItsOwnLoops(t *testing.T) {
	path := loopControlSource(t, `package main
use "std::io"
fn main() i32 {
defer {
for i := 0; i < 3; i = i + 1 {
if i == 0 {
continue
}
io::println("cleanup")
break
}
}
while true {
break
}
return 0
}
`)
	output := runLoopControl(t, path)
	if want := "cleanup\n"; output != want {
		t.Fatalf("want %q, got %q", want, output)
	}
}

func TestLoopControlCleansUpVaListOwnersBeforeScopeExit(t *testing.T) {
	path := loopControlSource(t, `package main
use "std::ffi::c"
fn inspect(marker i32, args @c ...) i32 {
outer := c::va_start(args)
defer c::va_end(&outer)
for i := 0; i < 3; i = i + 1 {
owner := c::va_copy(c::va_borrow(&outer))
defer c::va_end(&owner)
if c::va_arg<c::c_int>(&owner) != 42 {
return 1
}
if i == 0 {
continue
}
break
}
if c::va_arg<c::c_int>(&outer) != 42 {
return 2
}
return 0
}
fn main() i32 {
return inspect(0, 42)
}
`)
	runLoopControl(t, path)
}

func TestLoopControlStillChecksUncleanedVaListOwners(t *testing.T) {
	for _, keyword := range []string{"break", "continue"} {
		t.Run(keyword, func(t *testing.T) {
			runVaLifetimeFailure(t,
				"while true {\nowner := c::va_copy(args)\n"+keyword+"\n}\nreturn 0",
				"va_list owner left scope without va_end")
		})
	}
}

func TestLoopContextDoesNotLeakIntoGenericCallee(t *testing.T) {
	path := loopControlSource(t, `package main
fn invalid<T>() T {
break
return zero<T>()
}
fn main() i32 {
while true {
value := invalid<i32>()
break
}
return 0
}
`)
	_, diags := compiler.CompileOnly(path)
	if len(diags.Diags) == 0 ||
		!strings.Contains(diags.Diags[0].Message, "requires an enclosing loop") {
		t.Fatalf("want callee loop context diagnostic, got %v", diags.Diags)
	}
}
