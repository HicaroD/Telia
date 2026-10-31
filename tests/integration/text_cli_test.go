package integration

import (
	"bytes"
	"context"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/HicaroD/Telia/tests/compiler"
)

func compileTextCLI(t *testing.T, flags ...string) string {
	t.Helper()
	// Compile a private copy so concurrent suites never share an output binary.
	source, err := os.ReadFile("../../examples/text_cli.t")
	if err != nil {
		t.Fatal(err)
	}
	file := filepath.Join(t.TempDir(), "text_cli.t")
	if err := os.WriteFile(file, source, 0600); err != nil {
		t.Fatal(err)
	}
	path, diags := compiler.CompileOnlyWithCFlags(file, flags...)
	if len(diags.Diags) != 0 {
		t.Fatalf("compile: %v", diags.Diags)
	}
	t.Cleanup(func() { removeTestBinary(t, path) })
	return path
}

func checkTextCLI(
	t *testing.T,
	path string,
	args, env []string,
	input, wantOut, wantErr string,
	wantCode int,
) {
	t.Helper()
	stdout, stderr, code := runTextCLI(t, path, args, env, input)
	if code != wantCode || stdout != wantOut || stderr != wantErr {
		t.Fatalf(
			"exit=%d stdout=%q stderr=%q; want exit=%d stdout=%q stderr=%q",
			code, stdout, stderr, wantCode, wantOut, wantErr,
		)
	}
}

func runTextCLI(t *testing.T, path string, args, env []string, input string) (string, string, int) {
	t.Helper()
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	cmd := exec.CommandContext(ctx, path, args...)
	cmd.Env = env
	cmd.Stdin = bytes.NewBufferString(input)
	var stdout, stderr bytes.Buffer
	cmd.Stdout, cmd.Stderr = &stdout, &stderr
	err := cmd.Run()
	code := 0
	if err != nil {
		if exit, ok := err.(*exec.ExitError); ok {
			code = exit.ExitCode()
		} else {
			t.Fatal(err)
		}
	}
	return stdout.String(), stderr.String(), code
}

func TestTextCLITrimsUnicodeLinesFromStdin(t *testing.T) {
	for _, opt := range []string{"-O0", "-O3"} {
		t.Run(opt, func(t *testing.T) {
			path := compileTextCLI(t, opt)
			checkTextCLI(
				t, path, nil, []string{},
				"\u2003 世界 😀 \u00a0\n\t café \r\n", "世界 😀\ncafé\n\n", "", 0,
			)
		})
	}
}

func TestTextCLIFileInputEnvironmentDefaultAndOverride(t *testing.T) {
	file := filepath.Join(t.TempDir(), "世界 😀.txt")
	if err := os.WriteFile(file, []byte("\u2003é😀\u00a0"), 0600); err != nil {
		t.Fatal(err)
	}
	for _, opt := range []string{"-O0", "-O3"} {
		t.Run(opt, func(t *testing.T) {
			path := compileTextCLI(t, opt)
			checkTextCLI(
				t, path, []string{file}, []string{"TELIA_TEXT_REPEAT=2"},
				"ignored", "é😀\né😀\n", "", 0,
			)
			checkTextCLI(
				t, path, []string{file, "1"}, []string{"TELIA_TEXT_REPEAT=bad"},
				"", "é😀\n", "", 0,
			)
		})
	}
}

func TestTextCLICountsCodePointsNotBytesOrGraphemes(t *testing.T) {
	for _, opt := range []string{"-O0", "-O3"} {
		t.Run(opt, func(t *testing.T) {
			path := compileTextCLI(t, opt)
			checkTextCLI(
				t, path, []string{"-", "+2", "true"}, []string{},
				"é😀e\u0301\x00", "5\n5\n", "", 0,
			)
			checkTextCLI(t, path, []string{"-", "1", "false"}, []string{}, " é😀 ", "é😀\n", "", 0)
		})
	}
}

func TestTextCLIRejectsInvalidOptions(t *testing.T) {
	for _, opt := range []string{"-O0", "-O3"} {
		path := compileTextCLI(t, opt)
		for _, value := range []string{"", "0", "-1", "101", "1x", " 2", "9223372036854775808"} {
			checkTextCLI(
				t, path, []string{"-", value}, []string{},
				"text", "", "repeat must be an integer from 1 to 100\n", 1,
			)
		}
		checkTextCLI(
			t, path, nil, []string{"TELIA_TEXT_REPEAT="},
			"text", "", "repeat must be an integer from 1 to 100\n", 1,
		)
		checkTextCLI(
			t, path, []string{"-", "1", "TRUE"}, []string{},
			"text", "", "counts must be true or false\n", 1,
		)
		checkTextCLI(
			t, path, []string{"-", "1", "false", "extra"}, []string{},
			"text", "", "usage: text_cli [INPUT|- [REPEAT [COUNTS]]]\n", 1,
		)
	}
}

func TestTextCLIValidatesAllInputBeforeOutput(t *testing.T) {
	for _, opt := range []string{"-O0", "-O3"} {
		path := compileTextCLI(t, opt)
		checkTextCLI(t, path, nil, []string{}, "valid\ninvalid\xff", "", "invalid UTF-8 input\n", 1)
	}
}

func TestTextCLIEmptyInputAndOSFailures(t *testing.T) {
	missing := filepath.Join(t.TempDir(), "missing.txt")
	for _, opt := range []string{"-O0", "-O3"} {
		path := compileTextCLI(t, opt)
		checkTextCLI(t, path, nil, []string{}, "", "\n", "", 0)
		checkTextCLI(
			t, path, []string{"-", "100", "true"}, []string{},
			"", strings.Repeat("0\n", 100), "", 0,
		)
		checkTextCLI(t, path, []string{missing}, []string{}, "", "", "input read failed\n", 1)
		checkTextCLI(t, path, []string{""}, []string{}, "", "", "input read failed\n", 1)
	}
}

func TestTextCLIAllocationFailuresAndSanitizers(t *testing.T) {
	file := filepath.Join(t.TempDir(), "世界.txt")
	if err := os.WriteFile(file, []byte(" é😀 "), 0600); err != nil {
		t.Fatal(err)
	}
	for _, opt := range []string{"-O0", "-O3"} {
		flags := processControlFlags(t, opt, "text_cli_control.h")
		path := compileTextCLI(t, flags...)
		for _, tc := range []struct {
			name string
			args []string
			out  string
		}{
			{"environment", nil, "é😀\né😀\n"},
			{"format", []string{"-", "1", "true"}, "2\n"},
			{"file", []string{file, "1", "true"}, "2\n"},
		} {
			t.Run(opt+"/"+tc.name, func(t *testing.T) {
				for limit := 0; limit < 64; limit++ {
					env := []string{
						"TELIA_TEXT_REPEAT=2",
						"TELIA_TEXT_FAIL_AFTER=" + strconv.Itoa(limit),
					}
					out, errout, code := runTextCLI(t, path, tc.args, env, " é😀 ")
					if code == 0 {
						if limit == 0 || out != tc.out || errout != "" {
							t.Fatalf(
								"limit %d: unexpected success: stdout=%q stderr=%q",
								limit, out, errout,
							)
						}
						return
					}
					// The native destructor aborts if any allocation remains live.
					if code != 1 || errout == "" || out != "" {
						t.Fatalf("limit %d: exit=%d stdout=%q stderr=%q", limit, code, out, errout)
					}
				}
				t.Fatal("never reached a successful allocation budget")
			})
		}
		checkTextCLI(t, path, nil, []string{}, "valid\n\xff", "", "invalid UTF-8 input\n", 1)
		checkTextCLI(
			t, path, []string{"-", "bad"}, []string{},
			"", "", "repeat must be an integer from 1 to 100\n", 1,
		)
		checkTextCLI(
			t, path, []string{file + ".missing"}, []string{},
			"", "", "input read failed\n", 1,
		)
		checkTextCLI(
			t, path, nil, []string{"TELIA_TEXT_IO_MODE=3"},
			"input", "", "input read failed\n", 1,
		)
		// Mode 5 prevents both output and the diagnostic: a checked failure, not a panic.
		checkTextCLI(t, path, nil, []string{"TELIA_TEXT_IO_MODE=5"}, "input", "", "", 1)
		// Fail stdout's post-write flush once, leaving diagnostic output usable.
		checkTextCLI(
			t, path, nil, []string{"TELIA_TEXT_IO_MODE=12"},
			"input", "input", "stdout write failed\n", 1,
		)
	}
}
