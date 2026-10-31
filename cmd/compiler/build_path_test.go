package main

import (
	"os"
	"os/exec"
	"path/filepath"
	"testing"
)

func TestBuildDirectoryPathsProduceRunnableDefaultOutput(t *testing.T) {
	for _, spelling := range []string{".", "./", "../project", "absolute", "default"} {
		t.Run(spelling, func(t *testing.T) {
			dir := filepath.Join(t.TempDir(), "project")
			if err := os.Mkdir(dir, 0755); err != nil {
				t.Fatal(err)
			}
			if err := os.WriteFile(filepath.Join(dir, "main.t"), []byte("package main\nfn main() i32 {\nreturn 0\n}\n"), 0600); err != nil {
				t.Fatal(err)
			}
			args := []string{"-test.run=^TestCompilerDiagnosticProcess$", "--", "build"}
			if spelling == "absolute" {
				args = append(args, dir)
			} else if spelling != "default" {
				args = append(args, spelling)
			}
			cmd := exec.Command(os.Args[0], args...)
			cmd.Dir = dir
			cmd.Env = testCompilerEnv(t)
			if output, err := cmd.CombinedOutput(); err != nil {
				t.Fatalf("build directory: %v: %s", err, output)
			}
			if output, err := exec.Command(filepath.Join(dir, "project")).CombinedOutput(); err != nil {
				t.Fatalf("run default output: %v: %s", err, output)
			}
		})
	}
}
