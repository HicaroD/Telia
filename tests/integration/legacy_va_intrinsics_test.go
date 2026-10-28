package integration

import (
	"strings"
	"testing"
)

func TestLegacyRawPointerVaIntrinsicsAreNotExposedByLibc(t *testing.T) {
	exe, diags := compileRecordFixture(
		t,
		"package main\nuse \"std::libc\"\nfn main() {\nlibc::_va_start(nil)\n}\n",
	)
	if exe != "" {
		defer removeTestBinary(t, exe)
	}
	for _, d := range diags.Diags {
		if d.Span.Valid() && strings.Contains(d.Message, "not defined on scope") {
			return
		}
	}
	t.Fatalf(
		"want unavailable rawptr va intrinsic diagnosed before code generation, got %v",
		diags.Diags,
	)
}
