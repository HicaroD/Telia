package integration

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"testing"
	"time"

	"github.com/HicaroD/Telia/tests/compiler"
)

func runFSSource(t *testing.T, body string, flags ...string) {
	t.Helper()
	source := `package main
use "std::fs"
use "std::strings"
use "std::errors"
@[default_cc="c"]
extern {
fn _test_fs_mode(mode i32)
fn _test_fs_handles() i32
fn _test_fs_opens() i32
fn _test_fail_after(count i64)
fn _test_live_allocations() i64
}
fn main() i32 {
` + body + "\nreturn 0\n}\n"
	file := filepath.Join(t.TempDir(), "files.t")
	if err := os.WriteFile(file, []byte(source), 0600); err != nil {
		t.Fatal(err)
	}
	runStandardLibraryFixture(t, file, flags...)
}

func TestFSBinaryRoundTripAndTruncate(t *testing.T) {
	path := filepath.Join(t.TempDir(), "文件-é-😀.bin")
	payload := strings.Repeat("A\x00\xff\r\n世界", 3000)
	for i, data := range []string{payload, "short\x00\xff\n", ""} {
		body := fmt.Sprintf(`
if fs::write_file(%s, %s) != nil {
return %d
}

owner%d, err%d := fs::read_file(%s)
if err%d != nil or strings::view(&owner%d) != %s {
return %d
}
strings::free(&owner%d)
strings::free(&owner%d)
`, utf8Literal(path), utf8Literal(data), i*2+1, i, i, utf8Literal(path), i, i, utf8Literal(data), i*2+2, i, i)
		runFSSource(t, body)
		actual, err := os.ReadFile(path)
		if err != nil || string(actual) != data {
			t.Fatalf(
				"native file bytes differ: got %d bytes, want %d, error %v",
				len(actual),
				len(data),
				err,
			)
		}
	}
}

func TestFSRejectsInvalidPathsAndReportsOpenFailures(t *testing.T) {
	body := ""
	for i, path := range []string{"", "bad\x00path", "\xff", filepath.Join(t.TempDir(), "missing", "file")} {
		kind := []string{"InvalidPath", "EmbeddedZero", "InvalidUTF8", "FileOpenFailed"}[i]
		body += fmt.Sprintf(`
owner%d, err%d := fs::read_file(%s)
if err%d.kind != errors::Error::%s or strings::view(&owner%d) != "" {
return %d
}
strings::free(&owner%d)
write%d := fs::write_file(%s, "bytes")
if write%d.kind != errors::Error::%s {
return %d
}
`, i, i, utf8Literal(path), i, kind, i, i+1, i, i, utf8Literal(path), i, kind, i+10)
	}
	runFSSource(t, body)
}

func TestFSShortIOFailuresAndCleanup(t *testing.T) {
	header, err := filepath.Abs("testdata/fs_control.h")
	if err != nil {
		t.Fatal(err)
	}
	path := filepath.Join(t.TempDir(), "data.bin")
	payload := strings.Repeat("a\x00\xff\r\n", 6000)
	if err := os.WriteFile(path, []byte(payload), 0600); err != nil {
		t.Fatal(err)
	}
	flags := []string{"-include", header}
	if runtime.GOOS != "windows" {
		sanitizers := os.Getenv("TELIA_TEST_SANITIZERS")
		if sanitizers == "" {
			sanitizers = "undefined"
		}
		flags = append(flags, "-fsanitize="+sanitizers, "-fno-sanitize-recover=all")
	}
	t.Run("short IO", func(t *testing.T) {
		runFSSource(t, fmt.Sprintf(`
_test_fs_mode(1)
if fs::write_file(%s, %s) != nil {
return 1
}

owner, err := fs::read_file(%s)
if err != nil or strings::view(&owner) != %s {
return 2
}
strings::free(&owner)
if _test_fs_handles() != 0 or _test_live_allocations() != 0 {
return 3
}
`, utf8Literal(path), utf8Literal(payload), utf8Literal(path), utf8Literal(payload)), flags...)
	})
	for _, tc := range []struct {
		name  string
		mode  int
		write bool
		kind  string
	}{
		{"read and close failure", 2, false, "FileReadFailed"},
		{"write and close failure", 3, true, "FileWriteFailed"},
		{"read close failure", 4, false, "FileCloseFailed"},
		{"write close failure", 4, true, "FileCloseFailed"},
		{"read no progress", 5, false, "FileReadFailed"},
		{"write no progress", 6, true, "FileWriteFailed"},
		{"read open failure", 7, false, "FileOpenFailed"},
		{"write open failure", 7, true, "FileOpenFailed"},
	} {
		t.Run(tc.name, func(t *testing.T) {
			if err := os.WriteFile(path, []byte(payload), 0600); err != nil {
				t.Fatal(err)
			}
			operation := fmt.Sprintf(
				"owner, err := fs::read_file(%s)\nif strings::view(&owner) != \"\" {\nreturn 8\n}\nstrings::free(&owner)\n",
				utf8Literal(path),
			)
			if tc.write {
				operation = fmt.Sprintf(
					"err := fs::write_file(%s, %s)\n",
					utf8Literal(path),
					utf8Literal(payload),
				)
			}
			body := fmt.Sprintf(
				"_test_fs_mode(%d)\n%s\nif err.kind != errors::Error::%s {\nreturn 1\n}\nif _test_fs_handles() != 0 or _test_live_allocations() != 0 {\nreturn 2\n}\n",
				tc.mode,
				operation,
				tc.kind,
			)
			runFSSource(t, body, flags...)
			if tc.mode == 3 {
				data, err := os.ReadFile(path)
				if err != nil || string(data) != payload[:7] {
					t.Fatalf("partial write: %q, %v", data, err)
				}
			}
		})
	}

	t.Run("read nonseekable input", func(t *testing.T) {
		runFSSource(t, fmt.Sprintf(`
_test_fs_mode(8)
owner, err := fs::read_file("pipe-fixture")
if err != nil or strings::view(&owner) != %s {
return 1
}
strings::free(&owner)
if _test_fs_handles() != 0 or _test_live_allocations() != 0 {
return 2
}
`, utf8Literal("A\x00\xff\r\n")), flags...)
	})
	t.Run("allocation failure at every growth", func(t *testing.T) {
		if err := os.WriteFile(path, []byte(payload), 0600); err != nil {
			t.Fatal(err)
		}
		body := ""
		for i := 0; i < 4; i++ {
			body += fmt.Sprintf(`
_test_fail_after(%d)
owner%d, err%d := fs::read_file(%s)
if err%d.kind != errors::Error::AllocationFailed or strings::view(&owner%d) != "" {
return %d
}
strings::free(&owner%d)
if _test_fs_handles() != 0 or _test_live_allocations() != 0 {
return 9
}
`, i, i, i, utf8Literal(path), i, i, i+1, i)
		}
		body += fmt.Sprintf(`
opens := _test_fs_opens()
_test_fail_after(0)
write_err := fs::write_file(%s, "replacement")
if write_err.kind != errors::Error::AllocationFailed or _test_fs_opens() != opens {
return 10
}
if _test_fs_handles() != 0 or _test_live_allocations() != 0 {
return 11
}
_test_fail_after(-1)
`, utf8Literal(path))
		runFSSource(t, body, flags...)
		data, err := os.ReadFile(path)
		if err != nil || string(data) != payload {
			t.Fatalf(
				"allocation failure must precede truncation: got %d bytes, error %v",
				len(data),
				err,
			)
		}
	})
	t.Run("validation precedes OS calls", func(t *testing.T) {
		runFSSource(t, `
owner, err := fs::read_file("")
write_err := fs::write_file("bad\0path", "")
if err == nil or write_err == nil or _test_fs_opens() != 0 {
return 1
}
`, flags...)
	})
}

func TestFSExampleUsesExplicitTextValidation(t *testing.T) {
	path, diags := compiler.CompileOnly("../../examples/files.t")
	if len(diags.Diags) != 0 {
		t.Fatalf("compile example: %v", diags.Diags)
	}
	defer removeTestBinary(t, path)
	absolute, err := filepath.Abs(path)
	if err != nil {
		t.Fatal(err)
	}
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	cmd := exec.CommandContext(ctx, absolute)
	cmd.Dir = t.TempDir()
	output, err := cmd.CombinedOutput()
	if err != nil || string(output) != "File bytes validated as UTF-8\n" {
		t.Fatalf("example: %q, %v", output, err)
	}
	data, err := os.ReadFile(filepath.Join(cmd.Dir, "telia-file-example.bin"))
	if err != nil || string(data) != "Hello, 世界!\n\x00" {
		t.Fatalf("example bytes: %q, %v", data, err)
	}
}
