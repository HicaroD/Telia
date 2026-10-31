package integration

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/HicaroD/Telia/tests/compiler"
)

func TestEnumValuesPassThroughFunctions(t *testing.T) {
	path := filepath.Join(t.TempDir(), "main.t")
	source := `package main
enum Event {
Quit
Resize
}
fn identity(event Event) Event {
return event
}
fn main() i32 {
event := identity(Event::Quit)
if event != Event::Quit {
return 1
}
if event == Event::Resize {
return 2
}
return 0
}
`
	if err := os.WriteFile(path, []byte(source), 0600); err != nil {
		t.Fatal(err)
	}
	_, diags := compiler.CompileFile(path)
	if len(diags.Diags) != 0 {
		t.Fatalf("enum construction and identity: %v", diags.Diags)
	}
}

func TestEnumTagsUseDeclaredIntegerRepresentation(t *testing.T) {
	path := filepath.Join(t.TempDir(), "main.t")
	source := `package main
const FirstTag u64 = 4294967296
type Tag = u64
enum Event Tag {
Quit = FirstTag
Resize
}
fn main() i32 {
event Event := Event::Resize
if sizeof<Event>() != 8 {
return 1
}
if event == Event::Quit {
return 2
}
return 0
}
`
	if err := os.WriteFile(path, []byte(source), 0600); err != nil {
		t.Fatal(err)
	}
	_, diags := compiler.CompileFile(path)
	if len(diags.Diags) != 0 {
		t.Fatalf("explicit enum tag representation: %v", diags.Diags)
	}
}

func TestImportedEnumsKeepTheirIdentityThroughGenericFunctions(t *testing.T) {
	dir := t.TempDir()
	if err := os.Mkdir(filepath.Join(dir, "events"), 0700); err != nil {
		t.Fatal(err)
	}
	files := map[string]string{
		"events/events.t": "package events\nenum Event u8 {\nQuit = 255\n}\n",
		"main.t": `package main
use "pkg::events"
type Event = events::Event
fn identity<T>(value T) T {
return value
}
fn main() i32 {
event Event := identity<Event>(Event::Quit)
if event != Event::Quit {
return 1
}
if sizeof<Event>() != 1 {
return 2
}
return 0
}
`,
	}
	for name, source := range files {
		if err := os.WriteFile(filepath.Join(dir, name), []byte(source), 0600); err != nil {
			t.Fatal(err)
		}
	}
	_, diags := compiler.CompilePackage(dir)
	if len(diags.Diags) != 0 {
		t.Fatalf("imported enum identity: %v", diags.Diags)
	}
}

func TestEnumPayloadConstructionEvaluatesEachFieldOnce(t *testing.T) {
	path := filepath.Join(t.TempDir(), "main.t")
	source := `package main
enum Event u32 {
Quit = 0
Resize = 1 {
width i32
height i32
}
}
fn next(counter *i32) i32 {
*counter = *counter + 1
return *counter
}
fn identity(event Event) Event {
return event
}
fn main() i32 {
counter i32 := 0
event := identity(Event::Resize{width: next(&counter), height: next(&counter)})
if counter != 2 {
return 1
}
return 0
}
`
	if err := os.WriteFile(path, []byte(source), 0600); err != nil {
		t.Fatal(err)
	}
	_, diags := compiler.CompileFile(path)
	if len(diags.Diags) != 0 {
		t.Fatalf("named enum payload: %v", diags.Diags)
	}
}

func TestGenericEnumsHaveIndependentPayloadTypes(t *testing.T) {
	path := filepath.Join(t.TempDir(), "main.t")
	source := `package main
enum Option<T> {
None
Some {
value T
}
}
fn identity<T>(value T) T {
return value
}
fn main() i32 {
small := identity<Option<u8>>(Option<u8>::Some{value: 42})
large := identity<Option<u64>>(Option<u64>::Some{value: 42})
empty := Option<u8>::None
if sizeof<Option<u8>>() >= sizeof<Option<u64>>() {
return 1
}
return 0
}
`
	if err := os.WriteFile(path, []byte(source), 0600); err != nil {
		t.Fatal(err)
	}
	_, diags := compiler.CompileFile(path)
	if len(diags.Diags) != 0 {
		t.Fatalf("generic enum specialization: %v", diags.Diags)
	}
}

func TestMatchReadsNamedPayloadsAndProducesValues(t *testing.T) {
	path := filepath.Join(t.TempDir(), "main.t")
	source := `package main
enum Event u32 {
Quit = 0
Resize = 1 {
width i32
height i32
}
}
fn width(event Event) i32 {
return match event {
Event::Quit => 0
Event::Resize{width: w, ..} => w
}
}
fn main() i32 {
if width(Event::Resize{width: 800, height: 600}) != 800 {
return 1
}
if width(Event::Quit) != 0 {
return 2
}
return 0
}
`
	if err := os.WriteFile(path, []byte(source), 0600); err != nil {
		t.Fatal(err)
	}
	_, diags := compiler.CompileFile(path)
	if len(diags.Diags) != 0 {
		t.Fatalf("value-producing match: %v", diags.Diags)
	}
}

func TestMatchBlocksYieldCopiedFieldsAndRunCleanup(t *testing.T) {
	path := filepath.Join(t.TempDir(), "main.t")
	source := `package main
enum Option<T> {
None
Some {
value T
}
}
fn increment(counter *i32) {
*counter = *counter + 1
}
fn main() i32 {
item := Option<i32>::Some{value: 42}
counter i32 := 0
result := match item {
Option<i32>::None => 0
Option<i32>::Some{value} => {
defer increment(&counter)
value = value + 1
if value == 43 {
yield value
} else {
yield 99
}
}
}
original := match item {
Option<i32>::None => 0
Option<i32>::Some{value} => value
}
if result != 43 or original != 42 or counter != 1 {
return 1
}
return 0
}
`
	if err := os.WriteFile(path, []byte(source), 0600); err != nil {
		t.Fatal(err)
	}
	_, diags := compiler.CompileFile(path)
	if len(diags.Diags) != 0 {
		t.Fatalf("yield and copied pattern bindings: %v", diags.Diags)
	}
}

func TestGenericFunctionsMatchGenericEnums(t *testing.T) {
	path := filepath.Join(t.TempDir(), "main.t")
	source := `package main
enum Option<T> {
None
Some {
value T
}
}
fn unwrap<T>(item Option<T>) T {
return match item {
Option<T>::None => zero<T>()
Option<T>::Some{value} => {
yield value
}
}
}
fn main() i32 {
if unwrap<u8>(Option<u8>::Some{value: 42}) != 42 {
return 1
}
if unwrap<u64>(Option<u64>::Some{value: 9000000000}) != 9000000000 {
return 2
}
return 0
}
`
	if err := os.WriteFile(path, []byte(source), 0600); err != nil {
		t.Fatal(err)
	}
	_, diags := compiler.CompileFile(path)
	if len(diags.Diags) != 0 {
		t.Fatalf("generic match specializations: %v", diags.Diags)
	}
}

func TestGenericFunctionsConstructGenericEnums(t *testing.T) {
	path := filepath.Join(t.TempDir(), "main.t")
	source := `package main
enum Option<T> {
None
Some {
value T
}
}
fn some<T>(value T) Option<T> {
return Option<T>::Some{value: value}
}
fn main() i32 {
item := some<u64>(9000000000)
value := match item {
Option<u64>::None => 0
Option<u64>::Some{value} => value
}
if value != 9000000000 {
return 1
}
return 0
}
`
	if err := os.WriteFile(path, []byte(source), 0600); err != nil {
		t.Fatal(err)
	}
	_, diags := compiler.CompileFile(path)
	if len(diags.Diags) != 0 {
		t.Fatalf("generic enum factory: %v", diags.Diags)
	}
}

func TestMatchPreservesShortCircuitEvaluation(t *testing.T) {
	path := filepath.Join(t.TempDir(), "main.t")
	source := `package main
enum Event {
Quit
}
fn main() i32 {
counter i32 := 0
value := true or match Event::Quit {
Event::Quit => {
counter = counter + 1
yield false
}
}
if counter != 0 or not value {
return 1
}
return 0
}
`
	if err := os.WriteFile(path, []byte(source), 0600); err != nil {
		t.Fatal(err)
	}
	_, diags := compiler.CompileFile(path)
	if len(diags.Diags) != 0 {
		t.Fatalf("short-circuited match must not execute: %v", diags.Diags)
	}
}

func TestMatchReevaluatesLoopConditionsAndOnlyRunsSelectedArms(t *testing.T) {
	path := filepath.Join(t.TempDir(), "main.t")
	source := `package main
enum Event {
Quit
Resize
}
fn event(counter *i32) Event {
*counter = *counter + 1
return Event::Quit
}
fn main() i32 {
evaluations i32 := 0
iterations i32 := 0
while match event(&evaluations) {
Event::Quit => iterations < 3
Event::Resize => {
iterations = 99
yield false
}
} {
iterations = iterations + 1
}
if iterations != 3 or evaluations != 4 {
return 1
}
return 0
}
`
	if err := os.WriteFile(path, []byte(source), 0600); err != nil {
		t.Fatal(err)
	}
	_, diags := compiler.CompileFile(path)
	if len(diags.Diags) != 0 {
		t.Fatalf("match condition evaluation: %v", diags.Diags)
	}
}
