package integration

import (
	"reflect"
	"strings"
	"testing"
)

func collectionSanitizerFlags(goos string) []string {
	if goos == "windows" {
		return nil
	}
	return []string{"-fsanitize=undefined", "-fno-sanitize-recover=all"}
}

func collectionTextOutput(output string) string {
	return strings.ReplaceAll(output, "\r\n", "\n")
}

func TestCollectionHarnessKeepsWindowsChecksUnsanitized(t *testing.T) {
	if flags := collectionSanitizerFlags("windows"); len(flags) != 0 {
		t.Fatalf("Windows behavior checks must not require sanitizer runtimes: %v", flags)
	}
	for _, goos := range []string{"linux", "darwin"} {
		want := []string{"-fsanitize=undefined", "-fno-sanitize-recover=all"}
		if flags := collectionSanitizerFlags(goos); !reflect.DeepEqual(flags, want) {
			t.Fatalf("%s must retain UBSan: got %v, want %v", goos, flags, want)
		}
	}
}

func TestCollectionHarnessNormalizesOnlyTextLineEndings(t *testing.T) {
	for _, tc := range []struct{ input, want string }{
		{"record\n2\r\n9\r\nnested\nrecord\n", "record\n2\n9\nnested\nrecord\n"},
		{"owned\n0\r\n", "owned\n0\n"},
		{" \x00\rvalue\n\n", " \x00\rvalue\n\n"},
	} {
		if got := collectionTextOutput(tc.input); got != tc.want {
			t.Errorf("text output %q: got %q, want %q", tc.input, got, tc.want)
		}
	}
}
