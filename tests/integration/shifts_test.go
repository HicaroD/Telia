package integration

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/HicaroD/Telia/tests/compiler"
)

func TestShiftConstantsComposeFlagsAtTheirDeclaredWidth(t *testing.T) {
	path := filepath.Join(t.TempDir(), "main.t")
	source := `package main
type Flags = u64
const Fullscreen Flags = 1 << 0
const Resizable Flags = 1 << 5
const High Flags = 1 << 63
const Mask Flags = 1 << 63 | 1
const Wrapped u8 = 255 << 1
const Negative i8 = 64 << 1
const SignedRight i8 = -8 >> 1
fn main() i32 {
options := Fullscreen | Resizable
if options != 33 or High != 9223372036854775808 or Mask != 9223372036854775809 {
return 1
}
if Wrapped != 254 or Negative != -128 or SignedRight != -4 {
return 2
}
if 1 << 1 + 2 != 8 or 8 >> 1 & 3 != 0 {
return 3
}
return 0
}
`
	if err := os.WriteFile(path, []byte(source), 0600); err != nil {
		t.Fatal(err)
	}
	_, diags := compiler.CompileFile(path)
	if len(diags.Diags) != 0 {
		t.Fatalf("fixed-width shift constants: %v", diags.Diags)
	}
}

func TestRuntimeShiftsUseDefinedFixedWidthOperations(t *testing.T) {
	path := filepath.Join(t.TempDir(), "main.t")
	source := `package main
const Minimum i128 = -170141183460469231731687303715884105728
const Full u128 = 340282366920938463463374607431768211455
const High u128 = 170141183460469231731687303715884105728
fn left<T>(value T, count i16) T {
return value << count
}
fn right<T>(value T, count u8) T {
return value >> count
}
fn main() i32 {
small i8 := 64
wide i128 := Minimum
bits u128 := Full
if left<i8>(small, 1) != -128 or right<i8>(-8, 1) != -4 {
return 1
}
if right<i128>(wide, 127) != -1 {
return 2
}
if left<u128>(bits, 127) != High {
return 3
}
if left<i128>(-1, 127) != wide or right<i128>(wide, 0) != wide {
return 4
}
return 0
}
`
	if err := os.WriteFile(path, []byte(source), 0600); err != nil {
		t.Fatal(err)
	}
	exe, diags := compiler.CompileOnlyWithCFlags(
		path,
		"-fsanitize=undefined",
		"-fno-sanitize-recover=undefined",
	)
	if len(diags.Diags) != 0 {
		t.Fatalf("runtime shift lowering: %v", diags.Diags)
	}
	defer removeTestBinary(t, exe)
	if output, err := compiler.RunBinary(exe); err != nil {
		t.Fatalf("runtime shifts: %v: %s", err, output)
	}
}

func TestInvalidDynamicShiftCountsPanic(t *testing.T) {
	for _, count := range []string{"-1", "8", "170141183460469231731687303715884105728"} {
		t.Run(count, func(t *testing.T) {
			path := filepath.Join(t.TempDir(), "main.t")
			countType := "i16"
			if len(count) > 10 {
				countType = "u128"
			}
			source := "package main\nconst Count " + countType + " = " + count + "\nfn main() i32 {\nvalue u8 := 1\ncount := Count\nresult := value << count\nreturn 0\n}\n"
			if err := os.WriteFile(path, []byte(source), 0600); err != nil {
				t.Fatal(err)
			}
			exe, diags := compiler.CompileOnlyWithCFlags(
				path,
				"-fsanitize=undefined",
				"-fno-sanitize-recover=undefined",
			)
			if len(diags.Diags) != 0 {
				t.Fatalf("dynamic shift: %v", diags.Diags)
			}
			defer removeTestBinary(t, exe)
			output, err := compiler.RunBinary(exe)
			if err == nil || !strings.Contains(output, "shift count out of range") {
				t.Fatalf("want shift panic, got %v: %s", err, output)
			}
		})
	}
}
