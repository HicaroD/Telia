package integration

import (
	"strings"
	"testing"

	"github.com/HicaroD/Telia/tests/compiler"
)

func TestPointerCastsRoundTripRawStorageWithoutNativeBindings(t *testing.T) {
	path := loopControlSource(t, `package main
use "std::io"
struct Value {
number i32
}
fn main() i32 {
value := Value{number: 42}
raw := ptr_cast<rawptr>(&value)
typed := ptr_cast<*Value>(raw)
typed.number = 99
readonly := ptr_cast<const *Value>(typed)
if readonly.number != 99 or value.number != 99 {
return 1
}
empty := zero<rawptr>()
if ptr_cast<*Value>(empty) != nil {
return 2
}
io::println("pointers")
return 0
}
`)
	output, diags := compiler.CompileFile(path)
	if len(diags.Diags) != 0 {
		t.Fatalf("pointer conversions should compile and run: %v", diags.Diags)
	}
	if output != "pointers\n" {
		t.Fatalf("want pointer conversion success, got %q", output)
	}
}

func TestPointerCastsRejectConstErasureAndNonPointers(t *testing.T) {
	for _, body := range []string{
		"value i32 := 1\np := ptr_cast<const *i32>(&value)\nq := ptr_cast<*i32>(p)",
		"value i32 := 1\np := ptr_cast<const *i32>(&value)\nq := ptr_cast<rawptr>(p)",
		"value u64 := 1\nq := ptr_cast<*i32>(value)",
		"value i32 := 1\nq := ptr_cast<u64>(&value)",
	} {
		path := loopControlSource(t, "package main\nfn main() i32 {\n"+body+"\nreturn 0\n}\n")
		_, diags := compiler.CompileOnly(path)
		if len(diags.Diags) == 0 || !strings.Contains(diags.Diags[0].Message, "ptr_cast") {
			t.Fatalf("want invalid pointer conversion diagnostic, got %v", diags.Diags)
		}
	}
}

func TestWidenRejectsLossySignednessAndWidthChanges(t *testing.T) {
	for _, body := range []string{
		"value i32 := 1\nresult := widen<i16>(value)",
		"value i32 := 1\nresult := widen<u64>(value)",
		"value u32 := 1\nresult := widen<i32>(value)",
	} {
		path := loopControlSource(t, "package main\nfn main() i32 {\n"+body+"\nreturn 0\n}\n")
		_, diags := compiler.CompileOnly(path)
		if len(diags.Diags) == 0 ||
			!strings.Contains(diags.Diags[0].Message, "preserve every source value") {
			t.Fatalf("want lossy widening diagnostic, got %v", diags.Diags)
		}
	}
}

func TestCheckedNarrowingHandlesBoundariesAndSignedness(t *testing.T) {
	path := loopControlSource(t, `package main
use "std::errors"
use "std::io"
fn checked<T, U>(value T) (U, error) {
return narrow<U>(value)
}
fn main() i32 {
min i64 := -128
max u64 := 255
lo := checked<i64, i8>(min) @fail
hi := narrow<u8>(max) @fail
if lo != -128 or hi != 255 {
return 1
}
negative i64 := -1
zero, err := narrow<u8>(negative)
if err == nil or err.kind != errors::Error::ConversionOutOfRange or zero != 0 {
return 2
}
too_large u64 := 256
zero2, err2 := narrow<u8>(too_large)
if err2 == nil or zero2 != 0 {
return 3
}
signed_overflow u64 := 128
zero3, err3 := narrow<i8>(signed_overflow)
if err3 == nil or zero3 != 0 {
return 4
}
below i64 := -129
zero4, err4 := narrow<i8>(below)
if err4 == nil or zero4 != 0 {
return 5
}
half u64 := 18446744073709551615
wide := (widen<u128>(half) << 64) | widen<u128>(half)
same := narrow<u128>(wide) @fail
zero5, err5 := narrow<i128>(wide)
if same != wide or err5 == nil or zero5 != 0 {
return 6
}
one i128 := 1
half_min := -(one << 126)
signed_min := half_min + half_min
same_min := narrow<i128>(signed_min) @fail
if same_min != signed_min {
return 7
}
io::println("checked")
return 0
}
`)
	output, diags := compiler.CompileFile(path)
	if len(diags.Diags) != 0 {
		t.Fatalf("checked narrowing should compile and run: %v", diags.Diags)
	}
	if output != "checked\n" {
		t.Fatalf("want checked conversion success, got %q", output)
	}
}

func TestExplicitIntegerWideningPreservesValuesWithoutNativeBindings(t *testing.T) {
	path := loopControlSource(t, `package main
fn expand<T>(value T) i64 {
return widen<i64>(value)
}
fn main() i32 {
small i8 := -128
unsigned u32 := 4294967295
if expand<i8>(small) != -128 or widen<i64>(unsigned) != 4294967295 {
return 1
}
large u64 := 18446744073709551615
if widen<u128>(large) != 18446744073709551615 {
return 2
}
return 0
}
`)
	output, diags := compiler.CompileFile(path)
	if len(diags.Diags) != 0 {
		t.Fatalf("widening should compile: %v", diags.Diags)
	}
	if output != "" {
		t.Fatalf("want successful silent execution, got %q", output)
	}
}
