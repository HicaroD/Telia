package sema

import (
	"strings"
	"testing"

	"github.com/HicaroD/Telia/internal/diagnostics"
)

func TestRedundantTypeArgumentsProduceSafeHints(t *testing.T) {
	for _, tc := range []struct{ name, statement, want string }{
		{"omit", "a := identity<i32>(number)", "identity(number)"},
		{"keep numeric width", "a := identity<i64>(42)", ""},
		{"expected width", "a i64 := identity<i64>(42)", "identity(42)"},
		{"partial", `a := pair<i64, string>(42, "count")`, `pair<i64, _>(42, "count")`},
	} {
		t.Run(tc.name, func(t *testing.T) {
			c := parseAndCheck(`package main
fn identity<T>(value T) T {
return value
}
fn pair<A, B>(first A, second B) B {
return second
}
fn main() {
number i32 := 42
` + tc.statement + "\n}\n")
			if c.HasErrors() {
				t.Fatal(c.Diags)
			}
			var hints []diagnostics.Diag
			for _, d := range c.Diags {
				if d.Title == "Type arguments can be inferred" {
					hints = append(hints, d)
				}
			}
			if tc.want == "" {
				if len(hints) != 0 {
					t.Fatalf("unsafe hint: %v", hints)
				}
			} else if len(hints) != 1 || hints[0].Severity != diagnostics.Hint || !strings.Contains(hints[0].Help, tc.want) {
				t.Fatalf("want hint %q, got %v", tc.want, c.Diags)
			}
		})
	}
}

func TestInferNestedGenericAndCompoundArguments(t *testing.T) {
	requireInferenceOK(t, `struct Box<T> {
value T
}
fn unbox<T>(box Box<T>) T {
return box.value
}
fn dereference<T>(value *T) T {
return *value
}
fn main() {
box := Box<i32>{value: 42}
a := unbox(box)
b := dereference(&a)
}`)
}

func requireInferenceOK(t *testing.T, text string) {
	t.Helper()
	c := parseAndCheck("package main\n" + text)
	for _, d := range c.Diags {
		if d.Severity == diagnostics.Error {
			t.Fatalf("unexpected error: %v", c.Diags)
		}
	}
}

func TestHintsPreserveLegacyLiteralSuffixes(t *testing.T) {
	c := parseAndCheck(`package main
struct Box<T> {
value T
}
fn main() {
number i32 := 42
a := Box<i32>.{value: number}
}`)
	if c.HasErrors() {
		t.Fatal(c.Diags)
	}
	if len(c.Diags) != 1 || !strings.Contains(c.Diags[0].Help, "Box.{value: number}") {
		t.Fatalf("want complete legacy literal rewrite, got %v", c.Diags)
	}
}

func TestInferenceAcceptsPointersToOpaqueTypes(t *testing.T) {
	requireInferenceOK(t, `@[header="opaque.h", default_cc="c"]
extern {
opaque Window
fn get() *Window
}
fn identity<T>(value T) T {
return value
}
fn main() {
window := get()
copy := identity(window)
}`)
}

func TestHintsDoNotUseContextInferredFromTheCandidateItself(t *testing.T) {
	c := parseAndCheck(`package main
fn identity<T>(value T) T {
return value
}
fn main() {
a := identity(identity<i64>(42))
}`)
	if c.HasErrors() {
		t.Fatal(c.Diags)
	}
	if len(c.Diags) != 0 {
		t.Fatalf("the outer type depends on the explicit inner type: %v", c.Diags)
	}
}

func TestHintsDoNotSuggestInterdependentNestedRemovals(t *testing.T) {
	c := parseAndCheck(`package main
struct Box<T> {
value T
}
fn identity<T>(value T) T {
return value
}
fn main() {
a := Box<i64>{value: identity<i64>(42)}
}`)
	if c.HasErrors() {
		t.Fatal(c.Diags)
	}
	if len(c.Diags) != 1 || !strings.Contains(c.Diags[0].Help, "Box{value: identity<i64>(42)}") {
		t.Fatalf(
			"want only the enclosing rewrite; removing both would change i64 to int: %v",
			c.Diags,
		)
	}
}

func TestInferArrayLiteralElementsFromGenericParameterShape(t *testing.T) {
	requireInferenceOK(t, `fn first<T>(values T[2]) T {
return values[0]
}
fn main() {
number i32 := 42
a := first(.{number, number})
b := first(.{1, 2})
}`)
}

func TestCallbackAndValueConstraintsResolveEachOther(t *testing.T) {
	requireInferenceOK(t, `fn identity<T>(value T) T {
return value
}
fn transform<A, B>(callback fn(value A) B, value A) B {
return callback(value)
}
fn main() {
number i32 := 42
a := transform(identity, number)
b := transform(identity<_>, number)
}`)
}

func TestInferenceUsesSuccessTypeOfHandledCalls(t *testing.T) {
	requireInferenceOK(t, `struct Box<T> {
value T
}
fn create<T>() (T, error) {
return T{}, nil
}
fn main() {
value Box<i32> := create() @fail
}`)
}

func TestExpectedTypesReachIndexedPointerAndParallelAssignments(t *testing.T) {
	for _, statement := range []string{"values[0] = empty()", "*pointer = empty()", "first, second = twins()"} {
		t.Run(statement, func(t *testing.T) {
			requireInferenceOK(t, `struct Box<T> {
value T
}
fn empty<T>() T {
return T{}
}
fn twins<T>() (T, T) {
return T{}, T{}
}
fn main() {
values Box<i32>[1] := .{}
first Box<i32> := Box{}
second Box<i32> := Box{}
pointer := &first
`+statement+"\n}")
		})
	}
}

func TestDiscardedCallsDoNotInheritTheEnclosingReturnType(t *testing.T) {
	c := parseAndCheck(`package main
fn identity<T>(value T) T {
return value
}
fn task() i64 {
identity<i64>(42)
return 0
}
fn main() {
task()
}`)
	if c.HasErrors() {
		t.Fatal(c.Diags)
	}
	for _, d := range c.Diags {
		if d.Title == "Type arguments can be inferred" {
			t.Fatalf("discarded call has no expected result: %v", c.Diags)
		}
	}
}

func TestHintsHandleForwardDeclaredNamedReturnTypes(t *testing.T) {
	requireInferenceOK(t, `struct Box<T> {
value T
}
fn identity<T>(value T) T {
return value
}
fn main() {
a := identity<Box<i32>>(make())
}
fn make() Box<i32> {
return Box<i32>{value: 42}
}`)
}

func TestConcreteFunctionValuesRejectTypeArguments(t *testing.T) {
	c := parseAndCheck(`package main
fn plain(value i32) i32 {
return value
}
fn main() {
callback := plain<i32>
}`)
	if !c.HasErrors() || !containsDiag(c.Diags, "does not take type arguments") {
		t.Fatalf("want non-generic diagnostic, got %v", c.Diags)
	}
}

func TestInferenceThroughArraysSlicesTuplesAndCallbackResults(t *testing.T) {
	requireInferenceOK(t, `fn first<T>(values T[]) T {
return values[0]
}
fn array_first<T>(values T[2]) T {
return values[0]
}
fn numbered<T>(value T) (T, i32) {
return value, 1
}
fn call<T>(callback fn() T) T {
return callback()
}
fn number() i32 {
return 42
}
fn main() {
values i32[2] := .{42, 7}
a := first(values[:])
b := array_first(values)
c, count := numbered(a)
d := call(number)
}`)
}

func TestHintsDoNotUseNestedCallsInferredFromExplicitContext(t *testing.T) {
	c := parseAndCheck(`package main
fn identity<T>(value T) T {
return value
}
fn main() {
a := identity<i64>(identity(42))
}`)
	if c.HasErrors() {
		t.Fatal(c.Diags)
	}
	for _, d := range c.Diags {
		if d.Title == "Type arguments can be inferred" {
			t.Fatalf("removal would change i64 to int: %v", c.Diags)
		}
	}
}

func TestHintsPreserveFunctionTypesAndNestedArgumentCommas(t *testing.T) {
	c := parseAndCheck(`package main
fn choose<A, B>(value B) B {
return value
}
fn main() {
a := choose<fn(left i32, right i32) i32, string>("hello")
}`)
	if c.HasErrors() {
		t.Fatal(c.Diags)
	}
	if len(c.Diags) != 1 ||
		!strings.Contains(c.Diags[0].Help, `choose<fn(left i32, right i32) i32, _>("hello")`) {
		t.Fatalf("want a partial hint preserving the callback type, got %v", c.Diags)
	}
}

func TestGenericCallbacksInsideTemplates(t *testing.T) {
	requireInferenceOK(t, `fn identity<T>(value T) T {
return value
}
fn apply<T>(callback fn(value T) T, value T) T {
return callback(value)
}
fn forward<T>(value T) T {
callback fn(value T) T := identity<_>
return apply(identity, callback(value))
}
fn main() {
a i32 := forward(42)
}`)
}

func TestConcreteCallResultsPrecedeNestedLiteralDefaults(t *testing.T) {
	requireInferenceOK(t, `fn number() i32 {
return 42
}
fn identity<T>(value T) T {
return value
}
fn same<T>(first T, second T) T {
return first
}
fn main() {
a := same(identity(42), number())
b := same(number(), identity(42))
}`)
}

func TestInferenceRejectsAmbiguousAndConflictingTypes(t *testing.T) {
	for _, tc := range []struct{ name, body, want string }{
		{"no evidence", "a := absent()", "cannot infer type parameter 'T'"},
		{"placeholder without evidence", "a := absent<_>()", "cannot infer type parameter 'T'"},
		{"nil", "a := identity(nil)", "cannot infer type parameter 'T'"},
		{"conflicting inputs", `a := same(number, "bad")`, "conflicting types"},
		{"conflicting expected result", "a string := identity(number)", "conflicting types"},
		{"missing partial position", "a := pair<i32>(number, number)", "requires 2"},
		{"bare generic callback", "callback := identity", "cannot infer type parameter 'T'"},
		{"empty constructor", "a := Box{}", "cannot infer type parameter 'T'"},
		{"placeholder in annotation", "a Box<_> := Box<i32>{}", "'_' type not found"},
	} {
		t.Run(tc.name, func(t *testing.T) {
			c := parseAndCheck(`package main
struct Box<T> {
value T
}
fn absent<T>() T {
return T{}
}
fn identity<T>(value T) T {
return value
}
fn same<T>(first T, second T) T {
return first
}
fn pair<A, B>(first A, second B) B {
return second
}
fn main() {
number i32 := 42
` + tc.body + "\n}\n")
			if !c.HasErrors() || !containsDiag(c.Diags, tc.want) {
				t.Fatalf("want %q, got %v", tc.want, c.Diags)
			}
		})
	}
}

func TestInferenceHintsCoverConstructorsPatternsAndCallbacks(t *testing.T) {
	c := parseAndCheck(`package main
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
fn main() {
number i32 := 42
box := Box<i32>{value: number}
item := Option<i32>::Some{value: number}
value := match item {
Option<i32>::None => 0
Option<i32>::Some{value} => value
}
callback fn(value i32) i32 := identity<i32>
}`)
	if c.HasErrors() {
		t.Fatal(c.Diags)
	}
	count := 0
	for _, d := range c.Diags {
		if d.Title == "Type arguments can be inferred" {
			count++
		}
	}
	if count != 5 {
		t.Fatalf("want five refactoring hints, got %d: %v", count, c.Diags)
	}
}

func TestInferenceInsideGenericTemplates(t *testing.T) {
	requireInferenceOK(t, `struct Box<T> {
value T
}
enum Option<T> {
None
Some {
value T
}
}
fn wrap<T>(value T) Box<T> {
return Box{value: value}
}
fn some<T>(value T) Option<T> {
return Option<_>::Some{value: value}
}
fn main() {
a := wrap<i32>(42)
b := some<i32>(42)
}`)
}

func TestInferNestedExpressionsRegardlessOfArgumentOrder(t *testing.T) {
	requireInferenceOK(t, `fn identity<T>(value T) T {
return value
}
fn same<T>(first T, second T) T {
return first
}
struct Box<T> {
value T
}
fn empty<T>() T {
return T{}
}
fn main() {
number i32 := 42
a := same(identity(42), number)
b := same(number, identity(42))
box := Box<i32>{value: 42}
c := same(empty(), box)
}`)
}

func TestInferGenericCallbackValues(t *testing.T) {
	requireInferenceOK(t, `fn identity<T>(value T) T {
return value
}
fn apply<T>(callback fn(value T) T, value T) T {
return callback(value)
}
fn main() {
callback fn(value i32) i32 := identity
explicit := identity<i32>
partial fn(value i32) i32 := identity<_>
number i32 := 42
a := apply(identity, number)
b i32 := explicit(callback(partial(a)))
}`)
}

func TestInferEnumConstructionAndPatterns(t *testing.T) {
	requireInferenceOK(t, `enum Option<T> {
None
Some {
value T
}
}
fn absent() Option<i32> {
return Option::None
}
fn main() {
number i32 := 42
a := Option::Some{value: number}
b := Option<_>::Some{value: number}
c Option<i32> := Option::None
result := match a {
Option::None => 0
Option<_>::Some{value} => value
}
}`)
}

func TestInferStructConstructors(t *testing.T) {
	requireInferenceOK(t, `struct Pair<A, B> {
first A
second B
}
fn pair() Pair<string, i32> {
return Pair{first: "count", second: 42}
}
fn main() {
number i32 := 42
a := Pair{first: "count", second: number}
b := Pair<string, _>{first: "count", second: number}
c Pair<string, i32> := Pair{}
d Pair<string, i32> := pair()
}`)
}

func TestInferFunctionArguments(t *testing.T) {
	requireInferenceOK(t, `fn identity<T>(value T) T {
return value
}
fn main() {
number i32 := 42
result i32 := identity(number)
text string := identity("hello")
}`)
}

func TestInferPartialAndExpectedResults(t *testing.T) {
	requireInferenceOK(t, `fn identity<T>(value T) T {
return value
}
fn choose<A, B>(first A, second B) B {
return second
}
struct Box<T> {
value T
}
fn empty<T>() T {
return T{}
}
fn main() {
number i32 := 42
a i32 := identity<_>(number)
b i32 := choose<string, _>("count", number)
c i64 := identity(42)
d Box<i32> := empty()
}`)
}
