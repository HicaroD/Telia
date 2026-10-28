package main

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
)

func TestCompilerDiagnosticProcess(t *testing.T) {
	if os.Getenv("TELIA_DIAGNOSTIC_TEST") != "1" {
		return
	}
	for i, arg := range os.Args {
		if arg == "--" {
			os.Args = append([]string{"telia"}, os.Args[i+1:]...)
			main()
			return
		}
	}
}

func TestCLIReportsMultipleErrorsToStderrAndFails(t *testing.T) {
	dir := t.TempDir()
	path, executable := filepath.Join(dir, "main.t"), filepath.Join(dir, "program")
	if err := os.WriteFile(path, []byte("package main\nfn main() {\n  age i32 := \"twenty\"\n  count i32 := \"many\"\n}\n"), 0600); err != nil {
		t.Fatal(err)
	}
	cmd := exec.Command(
		os.Args[0],
		"-test.run=^TestCompilerDiagnosticProcess$",
		"--",
		"build",
		path,
		"--color=never",
		"-o",
		executable,
	)
	cmd.Env = append(os.Environ(), "TELIA_DIAGNOSTIC_TEST=1")
	var stdout, stderr bytes.Buffer
	cmd.Stdout, cmd.Stderr = &stdout, &stderr
	if err := cmd.Run(); err == nil {
		t.Fatal("invalid compilation exited successfully")
	}
	if stdout.Len() != 0 {
		t.Fatalf("diagnostics leaked to stdout: %q", stdout.String())
	}
	for _, want := range []string{"main.t:3:14", "main.t:4:16", "~~~~~~~", "Compilation failed · 2 errors"} {
		if !strings.Contains(stderr.String(), want) {
			t.Fatalf("missing %q in stderr:\n%s", want, stderr.String())
		}
	}
	if strings.Contains(stderr.String(), "\x1b") ||
		strings.Contains(stderr.String(), "compiler error found") {
		t.Fatalf("unexpected terminal output: %q", stderr.String())
	}
	if _, err := os.Stat(executable); !os.IsNotExist(err) {
		t.Fatal("invalid compilation produced an executable")
	}
}

func TestCLIColorFlagsAndEnvironment(t *testing.T) {
	path := filepath.Join(t.TempDir(), "main.t")
	if err := os.WriteFile(path, []byte("package main\nfn main() {\n unknown()\n}\n"), 0600); err != nil {
		t.Fatal(err)
	}
	for _, tc := range []struct {
		mode          string
		noColor, ansi bool
	}{
		{"auto", false, false}, {"auto", true, false}, {"always", false, true}, {"always", true, true}, {"never", false, false},
	} {
		t.Run(tc.mode+"/"+fmt.Sprint(tc.noColor), func(t *testing.T) {
			cmd := exec.Command(
				os.Args[0],
				"-test.run=^TestCompilerDiagnosticProcess$",
				"--",
				"build",
				path,
				"--color="+tc.mode,
			)
			for _, env := range os.Environ() {
				if !strings.HasPrefix(env, "NO_COLOR=") {
					cmd.Env = append(cmd.Env, env)
				}
			}
			cmd.Env = append(cmd.Env, "TELIA_DIAGNOSTIC_TEST=1")
			if tc.noColor {
				cmd.Env = append(cmd.Env, "NO_COLOR=1")
			}
			out, err := cmd.CombinedOutput()
			if err == nil {
				t.Fatal("invalid compilation succeeded")
			}
			if got := strings.Contains(string(out), "\x1b["); got != tc.ansi {
				t.Fatalf("want ANSI %v, got %q", tc.ansi, out)
			}
			if !strings.Contains(string(out), "Compilation failed · 1 error") {
				t.Fatalf("missing diagnostic summary: %s", out)
			}
		})
	}
}
