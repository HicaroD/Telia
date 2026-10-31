package config

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestInvalidStandardLibraryOverrideExplainsRecovery(t *testing.T) {
	oldDir, oldEnvs := TELIA_CONFIG_DIR, ENVS
	t.Cleanup(func() { TELIA_CONFIG_DIR, ENVS = oldDir, oldEnvs })
	TELIA_CONFIG_DIR = t.TempDir()
	t.Setenv("T_STD_PATH", filepath.Join(t.TempDir(), "missing"))
	err := SetupEnvFile()
	if err == nil || !strings.Contains(err.Error(), "unset T_STD_PATH") ||
		!strings.Contains(err.Error(), filepath.Join(TELIA_CONFIG_DIR, "env")) {
		t.Fatalf("expected actionable override diagnostic, got %v", err)
	}
}

func TestStandardLibraryOverrideMustBeDirectory(t *testing.T) {
	path := filepath.Join(t.TempDir(), "file")
	if err := os.WriteFile(path, []byte("not a directory"), 0600); err != nil {
		t.Fatal(err)
	}
	err := MapEnvToStruct(map[string]string{"T_STD_PATH": path}, &Envs{})
	if err == nil || !strings.Contains(err.Error(), "expected a directory") {
		t.Fatalf("expected directory diagnostic, got %v", err)
	}
}
