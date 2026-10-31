package integration

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/HicaroD/Telia/tests/compiler"
)

func runMatchProgram(t *testing.T, source string) {
	t.Helper()
	path := filepath.Join(t.TempDir(), "main.t")
	if err := os.WriteFile(path, []byte(source), 0600); err != nil {
		t.Fatal(err)
	}
	_, diags := compiler.CompileFile(path)
	if len(diags.Diags) != 0 {
		t.Fatalf("match program: %v", diags.Diags)
	}
}

func TestMatchIntegerAlternatives(t *testing.T) {
	runMatchProgram(t, `package main
fn classify(code i32) i32 {
return match code {
1 | 2 => 3
-1 => 4
_ => 0
}
}
fn main() i32 {
if classify(1) != 3 or classify(2) != 3 or classify(-1) != 4 or classify(9) != 0 {
return 1
}
return 0
}
`)
}

func TestMatchGuardsFallThroughAndEvaluateInOrder(t *testing.T) {
	runMatchProgram(t, `package main
fn value(calls *i32) i32 {
*calls = *calls + 1
return 15
}
fn guard(calls *i32, result bool) bool {
*calls = *calls + 1
return result
}
fn main() i32 {
values i32 := 0
guards i32 := 0
result := match value(&values) {
1 | 2 if guard(&guards, true) => 99
15 | 15 if guard(&guards, false) => 98
n if n > 10 and n < 20 and guard(&guards, true) => n + 1
_ if guard(&guards, true) => 97
_ => 0
}
if result != 16 or values != 1 or guards != 2 {
return 1
}
return 0
}
`)
}

func TestMatchEnumPayloadGuardsFallThrough(t *testing.T) {
	runMatchProgram(t, `package main
enum Option<T> {
None
Some {
value T
}
}
fn classify<T>(item Option<T>) T {
return match item {
Option<T>::Some{value} if value > 10 => value
Option<T>::Some{value} => value + 1
Option<T>::None => zero<T>()
}
}
fn main() i32 {
if classify<i32>(Option<i32>::Some{value: 15}) != 15 {
return 1
}
if classify<u64>(Option<u64>::Some{value: 5}) != 6 {
return 2
}
if classify<i32>(Option<i32>::None) != 0 {
return 3
}
return 0
}
`)
}

func TestMatchEnumAlternativesShareBindingsAndGuard(t *testing.T) {
	runMatchProgram(t, `package main
enum Event {
Quit
Idle
Width {
width i32
}
Height {
height i32
}
}
fn size(event Event) i32 {
return match event {
Event::Quit | Event::Idle => 0
Event::Width{width: n} | Event::Height{height: n} if n > 10 => n
Event::Width{width: n} | Event::Height{height: n} => n + 1
}
}
fn main() i32 {
if size(Event::Quit) != 0 or size(Event::Idle) != 0 {
return 1
}
if size(Event::Width{width: 20}) != 20 or size(Event::Height{height: 30}) != 30 {
return 2
}
if size(Event::Width{width: 5}) != 6 or size(Event::Height{height: 7}) != 8 {
return 3
}
return 0
}
`)
}

func TestMatchGuardsStayAtTheirEvaluationSite(t *testing.T) {
	runMatchProgram(t, `package main
fn guard(calls *i32) bool {
*calls = *calls + 1
return true
}
fn main() i32 {
calls i32 := 0
skip := true or match 1 {
n if guard(&calls) => false
_ => false
}
iterations i32 := 0
while match iterations {
n | n if n < 3 and guard(&calls) => true
_ => false
} {
iterations = iterations + 1
}
if not skip or iterations != 3 or calls != 3 {
return 1
}
return 0
}
`)
}

func TestMatchIntegerPatternsRespectSignedAndWideTypes(t *testing.T) {
	runMatchProgram(t, `package main
fn small(code i8) i32 {
return match code {
-128 | 127 => 1
_ => 0
}
}
fn wide(code u64) i32 {
return match code {
18446744073709551615 => 2
_ => 0
}
}
fn main() i32 {
if small(-128) != 1 or small(127) != 1 or small(0) != 0 {
return 1
}
if wide(18446744073709551615) != 2 or wide(0) != 0 {
return 2
}
return 0
}
`)
}
