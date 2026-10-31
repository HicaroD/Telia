package integration

import (
	"bytes"
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

func runIOSource(t *testing.T, body, input string, flags ...string) (string, string) {
	t.Helper()
	file := filepath.Join(t.TempDir(), "io.t")
	source := "package main\nuse \"std::io\"\nuse \"std::strings\"\nuse \"std::errors\"\n" + body
	if err := os.WriteFile(file, []byte(source), 0600); err != nil {
		t.Fatal(err)
	}
	return runIOFile(t, file, input, flags...)
}

func runIOFile(t *testing.T, file, input string, flags ...string) (string, string) {
	t.Helper()
	path, diags := compiler.CompileOnlyWithCFlags(file, flags...)
	if len(diags.Diags) != 0 {
		t.Fatalf("compile: %v", diags.Diags)
	}
	defer removeTestBinary(t, path)
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	cmd := exec.CommandContext(ctx, path)
	cmd.Stdin = bytes.NewBufferString(input)
	var stdout, stderr bytes.Buffer
	cmd.Stdout, cmd.Stderr = &stdout, &stderr
	if err := cmd.Run(); err != nil {
		t.Fatalf("execution: %v; stdout %q; stderr %q", err, stdout.String(), stderr.String())
	}
	return stdout.String(), stderr.String()
}

const ioControls = `@[default_cc="c"]
extern {
fn _test_io_mode(mode i32)
fn _test_io_flushes() i32
fn _test_io_streams_open() i32
fn _test_fail_after(count i64)
fn _test_live_allocations() i64
}
`

func ioControlFlags(t *testing.T, opt string) []string {
	t.Helper()
	header, err := filepath.Abs("testdata/io_control.h")
	if err != nil {
		t.Fatal(err)
	}
	flags := []string{opt, "-include", header}
	if runtime.GOOS != "windows" {
		sanitizers := os.Getenv("TELIA_TEST_SANITIZERS")
		if sanitizers == "" {
			sanitizers = "undefined"
		}
		flags = append(flags, "-fsanitize="+sanitizers, "-fno-sanitize-recover=all")
	}
	return flags
}

func TestIOShortOperationsFlushAndOpenStreams(t *testing.T) {
	input := strings.Repeat("a\x00\xff\r\n\x1a", 2000)
	for _, opt := range []string{"-O0", "-O3"} {
		t.Run(opt, func(t *testing.T) {
			out, errout := runIOSource(t, ioControls+`fn main() i32 {
_test_io_mode(1)
owner, err := io::read_stdin()
if err != nil {
return 1
}
_test_io_mode(1)
if io::write_stdout(strings::view(&owner)) != nil or _test_io_flushes() == 0 {
return 2
}
_test_io_mode(1)
if io::write_stderr(strings::view(&owner)) != nil or _test_io_flushes() == 0 {
return 3
}
strings::free(&owner)
if _test_live_allocations() != 0 or _test_io_streams_open() == 0 {
return 4
}
return 0
}`, input, ioControlFlags(t, opt)...)
			if out != input || errout != input {
				t.Fatal("short operations lost or changed bytes")
			}
		})
	}
}

func TestIOFailuresAndCleanup(t *testing.T) {
	for _, opt := range []string{"-O0", "-O3"} {
		for _, tc := range []struct {
			name  string
			mode  int
			write bool
			kind  string
		}{
			{"partial read", 2, false, "StreamReadFailed"},
			{"read bytes with error", 11, false, "StreamReadFailed"},
			{"read no progress", 3, false, "StreamReadFailed"},
			{"partial write", 4, true, "StreamWriteFailed"},
			{"write no progress", 5, true, "StreamWriteFailed"},
			{"flush", 6, true, "StreamFlushFailed"},
			{"mode read", 8, false, "StreamModeFailed"},
			{"mode write", 8, true, "StreamModeFailed"},
			{"restore read", 9, false, "StreamModeFailed"},
			{"restore write", 9, true, "StreamModeFailed"},
		} {
			if (tc.mode == 8 || tc.mode == 9) && runtime.GOOS != "windows" {
				continue
			}
			t.Run(opt+"/"+tc.name, func(t *testing.T) {
				operation := `owner, err := io::read_stdin()
if strings::view(&owner) != "" {
return 2
}
strings::free(&owner)`
				if tc.write {
					operation = `err := io::write_stdout("0123456789abcdef")`
				}
				body := fmt.Sprintf(`fn main() i32 {
_test_io_mode(%d)
%s
if err.kind != errors::Error::%s {
return 1
}
if _test_live_allocations() != 0 or _test_io_streams_open() == 0 {
return 3
}
return 0
}`, tc.mode, operation, tc.kind)
				out, errout := runIOSource(
					t,
					ioControls+body,
					"0123456789abcdef",
					ioControlFlags(t, opt)...)
				want := ""
				if tc.write && tc.mode != 5 && tc.mode != 8 {
					want = "0123456789abcdef"
					if tc.mode == 4 {
						want = "0123456"
					}
				}
				if out != want || errout != "" {
					t.Fatalf("stdout %q, want %q; stderr %q", out, want, errout)
				}
			})
		}
		for fail := 0; fail < 4; fail++ {
			t.Run(fmt.Sprintf("%s/allocation%d", opt, fail), func(t *testing.T) {
				runIOSource(t, ioControls+fmt.Sprintf(`fn main() i32 {
_test_fail_after(%d)
owner, err := io::read_stdin()
if err.kind != errors::Error::AllocationFailed or strings::view(&owner) != "" {
return 1
}
strings::free(&owner)
if _test_live_allocations() != 0 or _test_io_streams_open() == 0 {
return 2
}
return 0
}`, fail), strings.Repeat("x", 36000), ioControlFlags(t, opt)...)
			})
		}
	}
}

func TestIOEmptyEOFAndEmptyWriteFlush(t *testing.T) {
	for _, opt := range []string{"-O0", "-O3"} {
		t.Run(opt, func(t *testing.T) {
			runIOSource(t, ioControls+`fn main() i32 {
_test_fail_after(0)
owner, err := io::read_stdin()
if err != nil or strings::view(&owner) != "" {
return 1
}
strings::free(&owner)
second, second_err := io::read_stdin()
if second_err != nil or strings::view(&second) != "" {
return 2
}
strings::free(&second)
_test_io_mode(10)
out_err := io::write_stdout("")
err_err := io::write_stderr("")
if out_err.kind != errors::Error::StreamFlushFailed or err_err.kind != errors::Error::StreamFlushFailed {
return 3
}
if _test_live_allocations() != 0 or _test_io_streams_open() == 0 {
return 4
}
return 0
}`, "", ioControlFlags(t, opt)...)
		})
	}
}

func TestIONativeTextModeRestored(t *testing.T) {
	for _, opt := range []string{"-O0", "-O3"} {
		t.Run(opt, func(t *testing.T) {
			out, errout := runIOSource(t, ioControls+`@[default_cc="c"]
extern {
fn _test_io_native_text()
}
fn main() i32 {
_test_io_native_text()
if io::write_stdout("raw\n") != nil or io::write_stderr("raw\n") != nil {
return 1
}
_test_io_native_text()
return 0
}`, "", ioControlFlags(t, opt)...)
			text := "native\n"
			if runtime.GOOS == "windows" {
				text = "native\r\n"
			}
			want := text + "raw\n" + text
			if out != want || errout != want {
				t.Fatalf("stdout %q, stderr %q, want %q", out, errout, want)
			}
		})
	}
}

func TestIOStderrFailures(t *testing.T) {
	for _, opt := range []string{"-O0", "-O3"} {
		for _, mode := range []int{4, 5, 6} {
			t.Run(fmt.Sprintf("%s/%d", opt, mode), func(t *testing.T) {
				kind, want := "StreamWriteFailed", ""
				if mode == 4 {
					want = "0123456"
				}
				if mode == 6 {
					kind, want = "StreamFlushFailed", "0123456789abcdef"
				}
				out, errout := runIOSource(t, ioControls+fmt.Sprintf(`fn main() i32 {
_test_io_mode(%d)
err := io::write_stderr("0123456789abcdef")
if err.kind != errors::Error::%s or _test_io_flushes() == 0 or _test_io_streams_open() == 0 {
return 1
}
return 0
}`, mode, kind), "", ioControlFlags(t, opt)...)
				if out != "" || errout != want {
					t.Fatalf("stdout %q, stderr %q, want %q", out, errout, want)
				}
			})
		}
	}
}

func TestIOBrokenPipeIsChecked(t *testing.T) {
	if runtime.GOOS == "windows" {
		t.Skip("POSIX signal contract")
	}
	for _, opt := range []string{"-O0", "-O3"} {
		for _, pending := range []bool{false, true} {
			setup := ""
			if pending {
				setup = "_test_io_pending_signal()"
			}
			t.Run(fmt.Sprintf("%s/pending=%v", opt, pending), func(t *testing.T) {
				runIOSource(t, fmt.Sprintf(`@[default_cc="c"]
extern {
fn _test_io_broken_pipe() i32
fn _test_io_signal_unchanged() i32
fn _test_io_pending_signal()
}
fn main() i32 {
if _test_io_broken_pipe() == 0 {
return 1
}
%s
err := io::write_stdout("broken pipe")
if err.kind != errors::Error::StreamWriteFailed and err.kind != errors::Error::StreamFlushFailed {
return 2
}
if _test_io_signal_unchanged() == 0 {
return 3
}
return 0
}`, setup), "", ioControlFlags(t, opt)...)
			})
		}
	}
}

func TestIOReadSizeLimit(t *testing.T) {
	for _, opt := range []string{"-O0", "-O3"} {
		t.Run(opt, func(t *testing.T) {
			flags := append(ioControlFlags(t, opt), "-D_TELIA_IO_MAX_SIZE=8192")
			runIOSource(t, ioControls+`fn main() i32 {
owner, err := io::read_stdin()
if err.kind != errors::Error::SizeOverflow or strings::view(&owner) != "" {
strings::free(&owner)
return 1
}
strings::free(&owner)
if _test_live_allocations() != 0 or _test_io_streams_open() == 0 {
return 2
}
return 0
}`, strings.Repeat("x", 8193), flags...)
		})
	}
}

func TestIOStdinThroughEOF(t *testing.T) {
	for _, opt := range []string{"-O0", "-O3"} {
		for _, input := range []string{"", "A\x00\xff\r\n\x1aZ", strings.Repeat("a\x00\xff\r\n\x1a", 6000)} {
			t.Run(fmt.Sprintf("%s/%d", opt, len(input)), func(t *testing.T) {
				out, errout := runIOSource(t, fmt.Sprintf(`fn main() i32 {
owner, err := io::read_stdin()
if err != nil or strings::view(&owner) != %s {
return 1
}
if io::write_stdout(strings::view(&owner)) != nil {
return 2
}
strings::free(&owner)
strings::free(&owner)
return 0
}`, utf8Literal(input)), input, opt)
				if out != input || errout != "" {
					t.Fatalf("output differs: stdout %d bytes, stderr %q", len(out), errout)
				}
			})
		}
	}
}

func TestIOExample(t *testing.T) {
	input := "example\x00\xff\r\n\x1a"
	for _, opt := range []string{"-O0", "-O3"} {
		t.Run(opt, func(t *testing.T) {
			out, errout := runIOFile(t, "../../examples/stdio.t", input, opt)
			if out != input || errout != "" {
				t.Fatalf("stdout %q, stderr %q", out, errout)
			}
		})
	}
}

func TestIOStdoutBytes(t *testing.T) {
	for _, opt := range []string{"-O0", "-O3"} {
		t.Run(opt, func(t *testing.T) {
			out, errout := runIOSource(t, fmt.Sprintf(`fn main() i32 {
if io::write_stdout(%s) != nil {
return 1
}
return 0
}`, utf8Literal("A\x00\r\n\x1a")), "", opt)
			if out != "A\x00\r\n\x1a" || errout != "" {
				t.Fatalf("stdout %q, stderr %q", out, errout)
			}
		})
	}
}

func TestIOSeparateStreamsAndPrintln(t *testing.T) {
	for _, opt := range []string{"-O0", "-O3"} {
		t.Run(opt, func(t *testing.T) {
			out, errout := runIOSource(t, fmt.Sprintf(`fn main() i32 {
if io::write_stdout("before:") != nil {
return 1
}
io::println("line")
if io::write_stderr(%s) != nil {
return 2
}
if io::write_stdout("after") != nil or io::write_stdout("") != nil or io::write_stderr("") != nil {
return 3
}
return 0
}`, utf8Literal("err\x00\xff\r\n\x1a")), "", opt)
			if out != "before:line\nafter" || errout != "err\x00\xff\r\n\x1a" {
				t.Fatalf("stdout %q, stderr %q", out, errout)
			}
		})
	}
}
