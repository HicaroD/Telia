package main

import (
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
)

func TestReleaseVersionWithoutConfiguration(t *testing.T) {
	dir := t.TempDir()
	binary := filepath.Join(dir, "telia.exe")
	cmd := exec.Command(
		"go",
		"build",
		"-ldflags=-X main.Version=v0.1.0 -X main.Commit=abc123",
		"-o",
		binary,
		".",
	)
	if out, err := cmd.CombinedOutput(); err != nil {
		t.Fatalf("build compiler: %v\n%s", err, out)
	}
	for _, argument := range []string{"version", "--version"} {
		cmd := exec.Command(binary, argument)
		for _, env := range os.Environ() {
			if !strings.HasPrefix(env, "XDG_CONFIG_HOME=") {
				cmd.Env = append(cmd.Env, env)
			}
		}
		config := filepath.Join(dir, "unused-config")
		cmd.Env = append(cmd.Env, "XDG_CONFIG_HOME="+config)
		out, err := cmd.CombinedOutput()
		if err != nil || string(out) != "telia v0.1.0 (abc123)\n" {
			t.Fatalf("%s: %v, output %q", argument, err, out)
		}
		if _, err := os.Stat(config); !os.IsNotExist(err) {
			t.Fatalf("version reporting touched configuration: %v", err)
		}
	}
}
