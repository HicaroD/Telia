package integration

import (
	"bytes"
	"context"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/HicaroD/Telia/tests/compiler"
)

func TestFailTerminatesWithoutDeferredCleanup(t *testing.T) {
	fixture, err := os.ReadFile("testdata/fail_cleanup.t")
	if err != nil {
		t.Fatal(err)
	}
	for _, test := range []struct {
		name      string
		operation string
	}{
		{"error_only", "failure() @fail"},
		{"tuple_statement", "tuple_failure() @fail"},
		{"tuple_assignment", "value := tuple_failure() @fail"},
		{"during_cleanup", "defer failure() @fail"},
	} {
		t.Run(test.name, func(t *testing.T) {
			source := strings.Replace(
				string(fixture),
				"        failure() @fail",
				"        "+test.operation,
				1,
			)
			file := filepath.Join(t.TempDir(), "main.t")
			if err := os.WriteFile(file, []byte(source), 0600); err != nil {
				t.Fatal(err)
			}
			path, diags := compiler.CompileOnly(file)
			if len(diags.Diags) != 0 {
				t.Fatalf("compile: %v", diags.Diags)
			}
			defer removeTestBinary(t, path)
			ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
			defer cancel()
			cmd := exec.CommandContext(ctx, path)
			var stdout, stderr bytes.Buffer
			cmd.Stdout, cmd.Stderr = &stdout, &stderr
			err := cmd.Run()
			if exit, ok := err.(*exec.ExitError); !ok || exit.ExitCode() != 1 {
				t.Fatalf("want exit status 1, got %v (stderr: %s)", err, stderr.String())
			}
			if want := "completed scope\n"; stdout.String() != want {
				t.Fatalf("stdout: want %q, got %q", want, stdout.String())
			}
			if want := "runtime panic: failure\n"; stderr.String() != want {
				t.Fatalf("stderr: want %q, got %q", want, stderr.String())
			}
		})
	}
}

func TestSuccessfulFailRunsCleanupNormally(t *testing.T) {
	output, diags := compiler.CompileFile("testdata/fail_cleanup_success.t")
	if len(diags.Diags) != 0 {
		t.Fatalf("compile: %v", diags.Diags)
	}
	if want := "continued\nworker cleanup\ncaller cleanup\n"; output != want {
		t.Fatalf("want %q, got %q", want, output)
	}
}

func TestErrorPropagationRunsCleanupBeforeFatalBoundary(t *testing.T) {
	path, diags := compiler.CompileOnly("testdata/fail_cleanup_propagation.t")
	if len(diags.Diags) != 0 {
		t.Fatalf("compile: %v", diags.Diags)
	}
	defer removeTestBinary(t, path)
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	cmd := exec.CommandContext(ctx, path)
	var stdout, stderr bytes.Buffer
	cmd.Stdout, cmd.Stderr = &stdout, &stderr
	err := cmd.Run()
	if exit, ok := err.(*exec.ExitError); !ok || exit.ExitCode() != 1 {
		t.Fatalf("want exit status 1, got %v: %s", err, stderr.String())
	}
	if want := "inner cleanup\nworker second\nworker first\ncaller cleanup\n"; stdout.String() != want {
		t.Fatalf("stdout: want %q, got %q", want, stdout.String())
	}
	if want := "runtime panic: failure\n"; stderr.String() != want {
		t.Fatalf("stderr: want %q, got %q", want, stderr.String())
	}
}
