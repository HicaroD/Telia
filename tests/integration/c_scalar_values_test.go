package integration

import (
	"fmt"
	"os"
	"path/filepath"
	"testing"

	"github.com/HicaroD/Telia/tests/compiler"
)

func TestCFloatValuesAcceptTeliaArithmeticWithoutChangingStorageTypes(t *testing.T) {
	dir := t.TempDir()
	header := filepath.Join(dir, "floats.h")
	if err := os.WriteFile(header, []byte(`#define FIXTURE_HALF 0.5f
static inline float fixture_half(float value) { return value / 2.0f; }
static inline double fixture_half_double(double value) { return value / 2.0; }
`), 0600); err != nil {
		t.Fatal(err)
	}
	source := fmt.Sprintf(`package main
use "std::ffi::c"
@[header=%q, default_cc="c"]
extern {
  fn fixture_half(value c::c_float) c::c_float
  fn fixture_half_double(value c::c_double) c::c_double
  @[link_name="FIXTURE_HALF"]
  const Half c::c_float
}
fn to_c(value f32) c::c_float {
  return value / 2.0
}
fn to_telia(value c::c_float) f32 {
  return value / 2.0
}
fn main() i32 {
  x f32 := 16.0
  scale f32 := 2.0
  if fixture_half(x / scale) != 4.0 {
    return 1
  }
  value f32 := fixture_half(16.0) / scale
  if value != 4.0 or to_c(x) != 8.0 or to_telia(8.0) != 4.0 {
    return 2
  }
  native f32[1] := .{8.0}
  foreign c::c_float[1] := .{8.0}
  copied f32 := foreign[0]
  if fixture_half(native[0]) != 4.0 or copied != 8.0 {
    return 3
  }
  half f32 := Half
  if half != 0.5 {
    return 4
  }
  storage c::c_float := 0.0
  pointer := &storage
  *pointer = x / scale
  if storage != 8.0 {
    return 5
  }
  wide f64 := 16.0
  if fixture_half_double(wide / 2.0) != 4.0 {
    return 6
  }
  return 0
}
`, header)
	path := filepath.Join(dir, "main.t")
	if err := os.WriteFile(path, []byte(source), 0600); err != nil {
		t.Fatal(err)
	}
	_, diags := compiler.CompileFile(path)
	if len(diags.Diags) != 0 {
		t.Fatalf("C scalar value interoperability: %v", diags.Diags)
	}
}
