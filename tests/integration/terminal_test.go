package integration

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"testing"
	"time"

	"github.com/HicaroD/Telia/tests/compiler"
)

func TestTerminalControlsEncodeWithoutATerminal(t *testing.T) {
	runStandardLibraryFixture(t, "terminal_controls.t")
}

func TestTerminalExamplesCompileAndEchoPreservesBytes(t *testing.T) {
	for _, example := range []string{"terminal_echo.t", "terminal_raw.t", "terminal_screen.t"} {
		t.Run(example, func(t *testing.T) {
			path, diags := compiler.CompileOnly(filepath.Join("..", "..", "examples", example))
			if len(diags.Diags) != 0 {
				t.Fatalf("compile: %v", diags.Diags)
			}
			defer removeTestBinary(t, path)
			if example != "terminal_echo.t" {
				return
			}
			ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
			defer cancel()
			cmd := exec.CommandContext(ctx, path)
			want := bytes.Repeat([]byte("A\x00é\xff\n"), 8192)
			cmd.Stdin = bytes.NewReader(want)
			output, err := cmd.CombinedOutput()
			if err != nil || !bytes.Equal(output, want) {
				t.Fatalf("echo: want %d bytes, got %d, error %v", len(want), len(output), err)
			}
		})
	}
}

func TestTerminalWindowsConsoleModesInputAndResize(t *testing.T) {
	if runtime.GOOS != "windows" {
		t.Skip("requires an isolated Windows console")
	}
	include, err := filepath.Abs("testdata")
	if err != nil {
		t.Fatal(err)
	}
	runStandardLibraryFixture(t, "terminal_windows.t", "-I", include)
}

func TestTerminalWriteToClosedPipeReturnsAnError(t *testing.T) {
	path, diags := compiler.CompileOnly("testdata/terminal_broken_pipe.t")
	if len(diags.Diags) != 0 {
		t.Fatalf("compile: %v", diags.Diags)
	}
	defer removeTestBinary(t, path)
	reader, writer, err := os.Pipe()
	if err != nil {
		t.Fatal(err)
	}
	if err := reader.Close(); err != nil {
		t.Fatal(err)
	}
	defer func() { _ = writer.Close() }()
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	cmd := exec.CommandContext(ctx, path)
	cmd.Stdout = writer
	var stderr bytes.Buffer
	cmd.Stderr = &stderr
	if err := cmd.Run(); err != nil {
		t.Fatalf(
			"broken pipe must return an error, not terminate the application: %v: %s",
			err,
			&stderr,
		)
	}
}

func TestTerminalRejectsTerminalOperationsOnPipes(t *testing.T) {
	runStandardLibraryFixture(t, "terminal_redirected.t")
}

func TestTerminalWaitRejectsInvalidBorrowedHandle(t *testing.T) {
	runStandardLibraryFixture(t, "terminal_invalid_handle.t")
}

func TestTerminalControlsUseCapabilitiesAndExplicitOverrides(t *testing.T) {
	if runtime.GOOS == "windows" {
		t.Skip("POSIX pseudoterminal harness")
	}
	python, err := exec.LookPath("python3")
	if err != nil {
		t.Skip("python3 is required for the real pseudoterminal harness")
	}
	path, diags := compiler.CompileOnly("testdata/terminal_output.t")
	if len(diags.Diags) != 0 {
		t.Fatalf("compile: %v", diags.Diags)
	}
	defer removeTestBinary(t, path)
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	output, err := exec.CommandContext(ctx, python, "testdata/terminal_output.py", path).
		CombinedOutput()
	if err != nil {
		t.Fatalf("terminal output: %v\n%s", err, output)
	}
}

func TestTerminalWithNativeSanitizers(t *testing.T) {
	if runtime.GOOS == "windows" {
		t.Skip("MinGW-w64 does not provide the required sanitizer runtime")
	}
	sanitizers := os.Getenv("TELIA_TEST_SANITIZERS")
	if sanitizers == "" {
		sanitizers = "undefined"
	}
	for _, fixture := range []string{"terminal_controls.t", "terminal_redirected.t", "terminal_invalid_handle.t"} {
		t.Run(fixture, func(t *testing.T) {
			runStandardLibraryFixture(
				t,
				fixture,
				"-fsanitize="+sanitizers,
				"-fno-sanitize-recover=all",
			)
		})
	}
}

func TestTerminalWaitDistinguishesTimeoutFromInput(t *testing.T) {
	path, diags := compiler.CompileOnly("testdata/terminal_wait.t")
	if len(diags.Diags) != 0 {
		t.Fatalf("compile: %v", diags.Diags)
	}
	defer removeTestBinary(t, path)
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	cmd := exec.CommandContext(ctx, path)
	input, err := cmd.StdinPipe()
	if err != nil {
		t.Fatal(err)
	}
	defer func() { _ = input.Close() }()
	output, err := cmd.StdoutPipe()
	if err != nil {
		t.Fatal(err)
	}
	var stderr bytes.Buffer
	cmd.Stderr = &stderr
	if err := cmd.Start(); err != nil {
		t.Fatal(err)
	}
	defer func() { _ = cmd.Process.Kill() }()
	ready := make([]byte, 1)
	if _, err := io.ReadFull(output, ready); err != nil || ready[0] != 'S' {
		t.Fatalf("start of timed wait: %q, %v", ready, err)
	}
	start := time.Now()
	if _, err := io.ReadFull(output, ready); err != nil || ready[0] != 'R' {
		t.Fatalf("readiness: %q, %v", ready, err)
	}
	if time.Since(start) < 30*time.Millisecond {
		t.Fatal("timed wait returned before its deadline")
	}
	if _, err := input.Write([]byte{'X'}); err != nil {
		t.Fatal(err)
	}
	if err := cmd.Wait(); err != nil {
		t.Fatalf("wait behavior: %v: %s", err, stderr.String())
	}
}

func TestTerminalControlSequencesAndBufferBoundaries(t *testing.T) {
	cases := []struct{ control, want string }{
		{"CursorUp, first: 3", "\x1b[3A"},
		{"CursorDown, first: 4", "\x1b[4B"},
		{"CursorForward, first: 5", "\x1b[5C"},
		{"CursorBack, first: 6", "\x1b[6D"},
		{"CursorVisible, enabled: true", "\x1b[?25h"},
		{"CursorVisible", "\x1b[?25l"},
		{"ClearScreen", "\x1b[2J"},
		{"ClearLine", "\x1b[2K"},
		{"ResetStyle", "\x1b[0m"},
		{"Bold, enabled: true", "\x1b[1m"},
		{"Bold", "\x1b[22m"},
		{"Foreground, first: 255", "\x1b[38;5;255m"},
		{"Background, first: 0", "\x1b[48;5;0m"},
		{"ForegroundRGB, first: 12, second: 34, third: 56", "\x1b[38;2;12;34;56m"},
		{"BackgroundRGB, first: 255, second: 0, third: 1", "\x1b[48;2;255;0;1m"},
		{"AlternateScreen, enabled: true", "\x1b[?1049h"},
		{"AlternateScreen", "\x1b[?1049l"},
		{"BracketedPaste, enabled: true", "\x1b[?2004h"},
		{"BracketedPaste", "\x1b[?2004l"},
		{"MouseTracking, enabled: true", "\x1b[?1000h\x1b[?1006h"},
		{"MouseTracking", "\x1b[?1000l\x1b[?1006l"},
	}
	var source strings.Builder
	source.WriteString(
		"package main\nuse \"std::terminal\"\nfn main() i32 {\nbuffer := zero<u8[64]>()\n",
	)
	for i, tc := range cases {
		fmt.Fprintf(
			&source,
			"count%d, err%d := terminal::encode_control(terminal::Control{kind: terminal::ControlKind::%s}, buffer[:])\n",
			i,
			i,
			tc.control,
		)
		fmt.Fprintf(
			&source,
			"if err%d != nil or count%d != %d {\nreturn %d\n}\n",
			i,
			i,
			len(tc.want),
			i+1,
		)
		for j, b := range []byte(tc.want) {
			fmt.Fprintf(&source, "if buffer[%d] != %d {\nreturn %d\n}\n", j, b, i+1)
		}
	}
	source.WriteString(`small u8[1] := .{99}
count, err := terminal::encode_control(terminal::Control{kind: terminal::ControlKind::ClearScreen}, small[:])
if err == nil or count != 0 or small[0] != 99 {
return 90
}
invalid, invalid_err := terminal::encode_control(terminal::Control{kind: terminal::ControlKind::CursorPosition}, buffer[:])
if invalid_err == nil or invalid != 0 {
return 91
}
color, color_err := terminal::encode_control(terminal::Control{kind: terminal::ControlKind::Foreground, first: 256}, buffer[:])
if color_err == nil or color != 0 {
return 92
}
return 0
}
`)
	file := filepath.Join(t.TempDir(), "controls.t")
	if err := os.WriteFile(file, []byte(source.String()), 0600); err != nil {
		t.Fatal(err)
	}
	runStandardLibraryFixture(t, file)
}

func TestTerminalRawModeRestoresExactNativeSettings(t *testing.T) {
	runTerminalPTY(t, "terminal_raw.t", "raw")
}

func TestTerminalErrorPropagationRunsDeferredRestoration(t *testing.T) {
	runTerminalPTY(t, "terminal_error_cleanup.t", "cleanup")
}

func TestTerminalWaitReportsResizeWithoutConsumingInput(t *testing.T) {
	runTerminalPTY(t, "terminal_resize.t", "resize")
}

func TestTerminalOpensAndClosesOwnedControllingTerminal(t *testing.T) {
	runTerminalPTY(t, "terminal_controlling.t", "controlling")
}

func runTerminalPTY(t *testing.T, fixture, mode string) {
	t.Helper()
	if runtime.GOOS == "windows" {
		t.Skip("POSIX pseudoterminal harness")
	}
	python, err := exec.LookPath("python3")
	if err != nil {
		t.Skip("python3 is required for the real pseudoterminal harness")
	}
	path, diags := compiler.CompileOnly("testdata/" + fixture)
	if len(diags.Diags) != 0 {
		t.Fatalf("compile: %v", diags.Diags)
	}
	defer removeTestBinary(t, path)
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	output, err := exec.CommandContext(ctx, python, "testdata/terminal_pty.py", path, mode).
		CombinedOutput()
	if err != nil {
		t.Fatalf("pseudoterminal: %v\n%s", err, output)
	}
}

func TestTerminalClosingBorrowedHandlesPreservesNativeStream(t *testing.T) {
	output, diags := compiler.CompileFile("testdata/terminal_ownership.t")
	if len(diags.Diags) != 0 {
		t.Fatalf("compile: %v", diags.Diags)
	}
	if output != "ok" {
		t.Fatalf("want borrowed stream still usable, got %q", output)
	}
}

func TestTerminalReadsRedirectedInputUntilEOF(t *testing.T) {
	path, diags := compiler.CompileOnly("testdata/terminal_read.t")
	if len(diags.Diags) != 0 {
		t.Fatalf("compile: %v", diags.Diags)
	}
	defer removeTestBinary(t, path)
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	cmd := exec.CommandContext(ctx, path)
	want := []byte("A\x00B\n\xc3\xa9\xff")
	cmd.Stdin = bytes.NewReader(want)
	output, err := cmd.CombinedOutput()
	if err != nil || !bytes.Equal(output, want) {
		t.Fatalf("want %q, got %q, error %v", want, output, err)
	}
}

func TestTerminalWritesPreserveRedirectedBytes(t *testing.T) {
	output, diags := compiler.CompileFile("testdata/terminal_write.t")
	if len(diags.Diags) != 0 {
		t.Fatalf("compile: %v", diags.Diags)
	}
	if want := "A\x00B\n\xc3\xa9"; output != want {
		t.Fatalf("want bytes %q, got %q", want, output)
	}
}

func TestTerminalTextWritesPreserveEmbeddedZerosAndUTF8(t *testing.T) {
	output, diags := compiler.CompileFile("testdata/terminal_text.t")
	if len(diags.Diags) != 0 {
		t.Fatalf("compile: %v", diags.Diags)
	}
	if want := "A\x00é\n"; output != want {
		t.Fatalf("want bytes %q, got %q", want, output)
	}
}
