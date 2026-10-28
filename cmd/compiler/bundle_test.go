package main

import (
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"testing"
)

func TestBundledLibrariesAreFoundWithoutUserConfiguration(t *testing.T) {
	dir := t.TempDir()
	std := filepath.Join(dir, "base", "std")
	if err := os.MkdirAll(std, 0755); err != nil {
		t.Fatal(err)
	}
	std, err := filepath.EvalSymlinks(std)
	if err != nil {
		t.Fatal(err)
	}
	name := "telia"
	if runtime.GOOS == "windows" {
		name += ".exe"
	}
	binary := filepath.Join(dir, name)
	if out, err := exec.Command("go", "build", "-o", binary, ".").CombinedOutput(); err != nil {
		t.Fatalf("build compiler: %v\n%s", err, out)
	}
	cmd := exec.Command(binary, "env")
	for _, env := range os.Environ() {
		if !strings.HasPrefix(env, "XDG_CONFIG_HOME=") && !strings.HasPrefix(env, "T_STD_PATH=") {
			cmd.Env = append(cmd.Env, env)
		}
	}
	cmd.Env = append(cmd.Env, "XDG_CONFIG_HOME="+filepath.Join(t.TempDir(), "config"))
	out, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("bundled compiler: %v\n%s", err, out)
	}
	if !strings.Contains(string(out), "T_STD_PATH='"+std+"'") {
		t.Fatalf("bundle's standard library was not discovered: %s", out)
	}
}
