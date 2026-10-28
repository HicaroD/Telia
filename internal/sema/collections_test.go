package sema

import "testing"

func TestCollectionDiagnostics(t *testing.T) {
	for _, tc := range []struct{ name, body, want string }{
		{"negative index", "a i32[1] := .{}\nx := a[-1]", "nonnegative"},
		{"negative length", "a i32[-1] := .{}", "nonnegative integer literal"},
		{"expression length", "a i32[1 + 1] := .{}", "expected ]"},
		{"named length", "n := 2\na i32[n] := .{}", "nonnegative integer literal"},
		{"overflowing length", "a i32[18446744073709551616] := .{}", "invalid array length"},
		{"short literal", "a i32[2] := .{1}", "exactly 2 elements"},
		{"long literal", "a i32[1] := .{1, 2}", "exactly 1 elements"},
		{"untyped literal", "a := .{1}", "explicit array type"},
		{"element type", "a i32[1] := .{\"bad\"}", "cannot use string as i32"},
		{"array identity", "a i32[1] := .{}\nb i32[2] := a", "collection type mismatch"},
		{"index type", "a i32[1] := .{}\ni i32 := 0\nx := a[i]", "requires u64"},
		{"inferred index type", "a i32[1] := .{}\ni := 0\nx := a[i]", "requires u64"},
		{"range type", "a i32[1] := .{}\ni i32 := 0\nx := a[i:1]", "requires u64"},
		{"scalar indexing", "a := 1\nx := a[0]", "requires an array or slice"},
		{"infallible view annotation", "a i32[1] := .{}\nx := a[:] @fail", "does not return an error"},
	} {
		t.Run(tc.name, func(t *testing.T) {
			diags := parseAndCheck("package main\nfn main() {\n" + tc.body + "\n}\n")
			if !containsDiag(diags.Diags, tc.want) {
				t.Fatalf("want %q, got %v", tc.want, diags.Diags)
			}
		})
	}
}

func TestSliceParameterRejectsScalarVariables(t *testing.T) {
	diags := parseAndCheck(
		"package main\nfn take(values i32[]) {}\nfn main() {\nwrong := 1\ntake(wrong)\n}\n",
	)
	if !containsDiag(diags.Diags, "type mismatch") {
		t.Fatalf("expected slice/scalar mismatch, got %v", diags.Diags)
	}
}

func TestGenericCollectionElementRejectsDifferentStructIdentity(t *testing.T) {
	diags := parseAndCheck(`package main
struct A {
    value i32
}
struct B {
    value i32
}
fn store<T>(values T[], value T) {
    values[0] = value
}
fn main() {
    values A[1] := .{}
    wrong := B.{value: 42}
    store<A>(values[:], wrong)
}
`)
	if !containsDiag(diags.Diags, "type mismatch") {
		t.Fatalf("expected distinct struct element types, got %v", diags.Diags)
	}
}

func TestArrayBorrowRequiresAddressableStorage(t *testing.T) {
	diags := parseAndCheck(
		"package main\nfn make_values() i32[1] {\nreturn .{42}\n}\nfn main() {\nx := make_values()[:]\n}\n",
	)
	if !containsDiag(diags.Diags, "borrowing requires addressable array storage") {
		t.Fatalf("expected a Telia diagnostic, got %v", diags.Diags)
	}
}
