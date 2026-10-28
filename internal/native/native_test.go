package native_test

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/HicaroD/Telia/internal/native"
)

func TestMissingRuntimeArtifactReportsSetupBeforeCompilation(t *testing.T) {
	path := filepath.Join(t.TempDir(), "native.json")
	if err := os.WriteFile(path, []byte(`{"version":"1","targets":{"linux-amd64":{"runtime_files":["libmissing.so"]}}}`), 0644); err != nil {
		t.Fatal(err)
	}
	_, err := native.Resolve(path, "linux-amd64")
	if err == nil || !strings.Contains(err.Error(), "libmissing.so") ||
		!strings.Contains(err.Error(), "setup") ||
		!strings.Contains(err.Error(), "linux-amd64") {
		t.Fatalf("want actionable missing-runtime error, got %v", err)
	}
}

func TestInvalidNativeManifestIsRejected(t *testing.T) {
	for _, tc := range []struct{ name, manifest, want string }{
		{"missing version", `{"targets":{"linux-amd64":{}}}`, "version"},
		{"unknown field", `{"version":"1","targets":{"linux-amd64":{"libaries":[]}}}`, "libaries"},
		{"escaping path", `{"version":"1","targets":{"linux-amd64":{"include_dirs":["../outside"]}}}`, "package-relative"},
		{"unsupported target", `{"version":"1","targets":{"darwin-arm64":{}}}`, "linux-amd64"},
	} {
		t.Run(tc.name, func(t *testing.T) {
			path := filepath.Join(t.TempDir(), "native.json")
			if err := os.WriteFile(path, []byte(tc.manifest), 0644); err != nil {
				t.Fatal(err)
			}
			_, err := native.Resolve(path, "linux-amd64")
			if err == nil || !strings.Contains(err.Error(), tc.want) {
				t.Fatalf("want error containing %q, got %v", tc.want, err)
			}
		})
	}
}

func TestMissingOrInvalidNativeBuildInputsReportSetup(t *testing.T) {
	for _, field := range []string{"include_dirs", "sources", "libraries", "runtime_files"} {
		t.Run(field, func(t *testing.T) {
			dir := t.TempDir()
			path := filepath.Join(dir, "native.json")
			if field == "libraries" || field == "runtime_files" {
				if err := os.Mkdir(filepath.Join(dir, "invalid"), 0755); err != nil {
					t.Fatal(err)
				}
			}
			text := `{"version":"1","targets":{"linux-amd64":{"` + field + `":["invalid"]}}}`
			if err := os.WriteFile(path, []byte(text), 0644); err != nil {
				t.Fatal(err)
			}
			_, err := native.Resolve(path, "linux-amd64")
			if err == nil || !strings.Contains(err.Error(), "setup") ||
				!strings.Contains(err.Error(), "linux-amd64") {
				t.Fatalf("want actionable invalid-input error, got %v", err)
			}
		})
	}
}
