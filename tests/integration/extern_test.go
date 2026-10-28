package integration

import (
	"path/filepath"
	"testing"

	"github.com/HicaroD/Telia/tests/compiler"
)

func TestExternNestedPointersRoundTripThroughC(t *testing.T) {
	header, err := filepath.Abs("testdata/extern_pointers.h")
	if err != nil {
		t.Fatal(err)
	}
	path, diags := compiler.CompileOnlyWithCFlags(
		"testdata/extern_pointers.t",
		"-include",
		header,
		"-Werror=incompatible-pointer-types",
	)
	if len(diags.Diags) != 0 {
		t.Fatalf("compiling nested extern signature: %v", diags.Diags)
	}
	defer removeTestBinary(t, path)
	output, err := compiler.RunBinary(path)
	if err != nil {
		t.Fatalf("C pointer round trip: %v: %s", err, output)
	}
}
