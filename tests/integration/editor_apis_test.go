package integration

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"regexp"
	"runtime"
	"strings"
	"testing"
	"time"
)

func runEditorSource(t *testing.T, imports, body string, flags ...string) {
	t.Helper()
	if runtime.GOOS != "windows" {
		sanitizers := os.Getenv("TELIA_TEST_SANITIZERS")
		if sanitizers == "" {
			sanitizers = "undefined"
		}
		flags = append(flags, "-fsanitize="+sanitizers, "-fno-sanitize-recover=all")
	}
	body = strings.ReplaceAll(body, " { return ", " {\nreturn ")
	body = regexp.MustCompile(`(\d+) }`).ReplaceAllString(body, "$1\n}")
	source := "package main\nuse \"std::strings\"\nuse \"std::errors\"\n" + imports + "\nfn main() i32 {\n" + body + "\nreturn 0\n}\n"
	file := filepath.Join(t.TempDir(), "editor.t")
	if err := os.WriteFile(file, []byte(source), 0600); err != nil {
		t.Fatal(err)
	}
	runStandardLibraryFixture(t, file, flags...)
}

func TestEditorMetadataReportsIdentityTypeSizeTimeAndMissing(t *testing.T) {
	dir := t.TempDir()
	file := filepath.Join(dir, "data")
	if err := os.WriteFile(file, []byte("abc"), 0640); err != nil {
		t.Fatal(err)
	}
	stamp := time.Unix(1700000000, 0)
	if err := os.Chtimes(file, stamp, stamp); err != nil {
		t.Fatal(err)
	}
	runEditorSource(t, `use "std::fs"`, fmt.Sprintf(`
info, err := fs::stat(%s)
if err != nil or info.kind != fs::FileType::Regular or info.size != 3 or info.modified_seconds != 1700000000 { return 1 }
again, again_err := fs::stat(%s)
if again_err != nil or not fs::same_file(info, again) { return 2 }
directory, dir_err := fs::stat(%s)
if dir_err != nil or directory.kind != fs::FileType::Directory or fs::same_file(info, directory) { return 3 }
missing, missing_err := fs::stat(%s)
if missing_err.kind != errors::Error::NotFound { return 4 }
`, utf8Literal(file), utf8Literal(file), utf8Literal(dir), utf8Literal(filepath.Join(dir, "missing"))))
}

func TestEditorPathsResolveExistingAndNewFilesAndChangeDirectory(t *testing.T) {
	dir := t.TempDir()
	canonical, err := filepath.EvalSymlinks(dir)
	if err != nil {
		t.Fatal(err)
	}
	runEditorSource(t, "use \"std::path\"\nuse \"std::process\"\nuse \"std::fs\"", fmt.Sprintf(`
if process::chdir(%s) != nil { return 1 }
directory, dir_err := path::canonicalize(".")
if dir_err != nil or strings::view(&directory) != %s { return 2 }
strings::free(&directory)
new_file, new_err := path::canonicalize_for_create("new.bin")
if new_err != nil or strings::view(&new_file) != %s { return 3 }
strings::free(&new_file)
missing, missing_err := path::canonicalize("new.bin")
if missing_err.kind != errors::Error::NotFound { return 4 }
strings::free(&missing)
bad_parent, parent_err := path::canonicalize_for_create("missing/new.bin")
if parent_err.kind != errors::Error::NotFound { return 5 }
strings::free(&bad_parent)
if fs::write_file("written.bin", "local") != nil { return 6 }
`, utf8Literal(dir), utf8Literal(canonical), utf8Literal(filepath.Join(canonical, "new.bin"))))
	data, err := os.ReadFile(filepath.Join(dir, "written.bin"))
	if err != nil || string(data) != "local" {
		t.Fatalf("working directory write: %q, %v", data, err)
	}
}

func TestEditorBoundedReadsReturnOnlyPreviewAndDistinguishEOF(t *testing.T) {
	file := filepath.Join(t.TempDir(), "large")
	f, err := os.Create(file)
	if err != nil {
		t.Fatal(err)
	}
	if _, err := f.Write([]byte("A\x00\xffZ")); err != nil {
		t.Fatal(err)
	}
	if err := f.Truncate(1 << 30); err != nil {
		t.Fatal(err)
	}
	if err := f.Close(); err != nil {
		t.Fatal(err)
	}
	runEditorSource(t, `use "std::fs"`, fmt.Sprintf(`
preview, err := fs::read_preview(%s, 3)
if err != nil or strings::view(&preview) != %s { return 1 }
strings::free(&preview)
file, open_err := fs::open_read(%s)
if open_err != nil { return 2 }
buffer u8[4] := .{}
count, eof, read_err := fs::read(&file, buffer[:])
if read_err != nil or eof or count != 4 or buffer[3] != 90 { return 3 }
if fs::close(&file) != nil or fs::close(&file) != nil { return 4 }
closed_count, closed_eof, closed_err := fs::read(&file, buffer[:])
if closed_err.kind != errors::Error::InvalidState { return 5 }
`, utf8Literal(file), utf8Literal("A\x00\xff"), utf8Literal(file)))
}

func TestEditorAtomicSaveReplacesBytesAndPreservesPermissions(t *testing.T) {
	dir := t.TempDir()
	file := filepath.Join(dir, "save")
	if err := os.WriteFile(file, []byte("original"), 0640); err != nil {
		t.Fatal(err)
	}
	check := `if err != nil { return 2 }`
	if runtime.GOOS == "windows" {
		check = `if err.kind != errors::Error::Unsupported { return 2 }`
	}
	runEditorSource(t, `use "std::fs"`, fmt.Sprintf(`
replaced, err := fs::write_atomic(%s, %s)
if not replaced { return 1 }
%s
`, utf8Literal(file), utf8Literal("saved\x00\xff"), check))
	data, err := os.ReadFile(file)
	if err != nil || string(data) != "saved\x00\xff" {
		t.Fatalf("saved bytes: %q, %v", data, err)
	}
	info, err := os.Stat(file)
	if err != nil {
		t.Fatal(err)
	}
	if runtime.GOOS != "windows" && info.Mode().Perm() != 0640 {
		t.Fatalf("permissions: %o", info.Mode().Perm())
	}
	entries, err := os.ReadDir(dir)
	if err != nil || len(entries) != 1 {
		t.Fatalf("temporary file leaked: %v, %v", entries, err)
	}
}

// A real child executable, independent of shells, rg, or clipboard installation.
func TestEditorNativeChild(t *testing.T) {
	if len(os.Args) < 4 || os.Args[2] != "--" {
		return
	}
	switch os.Args[3] {
	case "arguments":
		want := []string{"", "two words", `a"b\c\`, "世界-é-😀"}
		if len(os.Args[4:]) != len(want) {
			os.Exit(99)
		}
		for i, value := range want {
			if os.Args[i+4] != value {
				os.Exit(98)
			}
		}
		os.Exit(7)
	case "pipes":
		_, _ = os.Stdout.Write([]byte("ready"))
		_, _ = io.Copy(os.Stdout, os.Stdin)
		_, _ = os.Stderr.Write([]byte("err"))
		os.Exit(7)
	case "hold":
		_, _ = os.Stdout.Write([]byte("ready"))
		time.Sleep(time.Minute)
		os.Exit(0)
	}
	os.Exit(97)
}

func TestEditorSubprocessPreservesArgumentsAndExitStatus(t *testing.T) {
	executable, err := os.Executable()
	if err != nil {
		t.Fatal(err)
	}
	runEditorSource(t, `use "std::process"`, fmt.Sprintf(`
arguments string[7] := .{"-test.run=^TestEditorNativeChild$", "--", "arguments", "", "two words", %s, "世界-é-😀"}
child, err := process::spawn(%s, arguments[:], process::SpawnOptions{})
if err != nil { return 1 }
status, wait_err := process::wait(&child)
if wait_err != nil or status.code != 7 or status.signal != 0 { return 2 }
again, done, poll_err := process::poll(&child)
if poll_err != nil or not done or again.code != 7 { return 3 }
if process::close_child(&child) != nil or process::close_child(&child) != nil { return 4 }
`, utf8Literal(`a"b\c\`), utf8Literal(executable)))
}

func TestEditorSubprocessPipesDistinguishWouldBlockDataAndEOF(t *testing.T) {
	executable, err := os.Executable()
	if err != nil {
		t.Fatal(err)
	}
	runEditorSource(t, `use "std::process"`, fmt.Sprintf(`
arguments string[3] := .{"-test.run=^TestEditorNativeChild$", "--", "pipes"}
child, err := process::spawn(%s, arguments[:], process::SpawnOptions{pipe_stdin: true, pipe_stdout: true, pipe_stderr: true})
if err != nil { return 1 }
defer process::close_child(&child)
buffer u8[5] := .{}
received u64 := 0
attempts u64 := 0
while received < 5 and attempts < 10000000 {
    portion := buffer[received:5] @fail
    count, eof, blocked, read_err := process::read_stdout(&child, portion)
    if read_err != nil or eof { return 2 }
    received = received + count
    attempts = attempts + 1
}

if received != 5 or buffer[0] != 114 or buffer[4] != 121 { return 3 }
empty_count, empty_eof, would_block, empty_err := process::read_stdout(&child, buffer[:])
if empty_err != nil or empty_eof or not would_block or empty_count != 0 { return 4 }
written, blocked_write, write_err := process::write_stdin(&child, %s)
if write_err != nil or blocked_write or written != 3 { return 5 }
if process::close_stdin(&child) != nil { return 6 }
received = 0
attempts = 0
while received < 3 and attempts < 10000000 {
    portion := buffer[received:3] @fail
    count, eof, blocked, read_err := process::read_stdout(&child, portion)
    if read_err != nil { return 7 }
    received = received + count
    attempts = attempts + 1
}
if received != 3 or buffer[0] != 65 or buffer[1] != 0 or buffer[2] != 255 { return 8 }
status, wait_err := process::wait(&child)
if wait_err != nil or status.code != 7 { return 9 }
eof_count, eof, eof_blocked, eof_err := process::read_stdout(&child, buffer[:])
if eof_err != nil or not eof or eof_blocked or eof_count != 0 { return 10 }
stderr_count, stderr_eof, stderr_blocked, stderr_err := process::read_stderr(&child, buffer[:])
if stderr_err != nil or stderr_count != 3 or buffer[0] != 101 { return 11 }
`, utf8Literal(executable), utf8Literal("A\x00\xff")))
}

func TestEditorByteStringCopiesPreserveBinaryAndOwnership(t *testing.T) {
	runEditorSource(t, `use "std::lists"`, fmt.Sprintf(`
owner, err := strings::clone(%s)
if err != nil { return 1 }
bytes, bytes_err := strings::to_bytes(strings::view(&owner))
if bytes_err != nil { return 2 }
strings::free(&owner)
data := lists::view<u8>(&bytes)
if len(data) != 3 or data[0] != 65 or data[1] != 0 or data[2] != 255 { return 3 }
copy, copy_err := strings::from_bytes(data)
if copy_err != nil { return 4 }
data[0] = 90
lists::free<u8>(&bytes)
if strings::view(&copy) != %s { return 5 }
value, access_err := strings::byte_at(strings::view(&copy), 2)
if value != 255 or access_err != nil { return 6 }
missing, bounds_err := strings::byte_at(strings::view(&copy), 3)
if missing != 0 or bounds_err.kind != errors::Error::OutOfBounds { return 7 }
strings::free(&copy)
`, utf8Literal("A\x00\xff"), utf8Literal("A\x00\xff")))
}

func TestEditorUnicodeWidthHandlesCombiningWideAndNonprintable(t *testing.T) {
	body := ""
	for i, tc := range []struct {
		point uint32
		width int
	}{
		{0x41, 1}, {0x4e2d, 2}, {0x301, 0}, {0x200d, 0}, {0x93e, 0},
		{0x1f600, 2}, {0, 0}, {9, -1}, {0x7f, -1}, {0x2028, -1},
		{0x378, -1}, {0x10ffff, -1}, {0xa1, 1}, {0x1160, 0},
		{0x1e4ec, 0}, {0x1fae9, 2},
	} {
		body += fmt.Sprintf(
			"width%d, err%d := unicode::width(%d)\nif err%d != nil or width%d != %d { return %d }\n",
			i,
			i,
			tc.point,
			i,
			i,
			tc.width,
			i+1,
		)
	}
	body += `invalid, invalid_err := unicode::width(55296)
if invalid != 0 or invalid_err.kind != errors::Error::InvalidCodePoint { return 90 }
too_large, large_err := unicode::width(1114112)
if too_large != 0 or large_err.kind != errors::Error::InvalidCodePoint { return 91 }
`
	runEditorSource(t, `use "std::unicode"`, body)
}

func TestEditorMissingFilesAreDistinguishableAcrossReadAPIs(t *testing.T) {
	missing := filepath.Join(t.TempDir(), "absent")
	runEditorSource(t, `use "std::fs"`, fmt.Sprintf(`
all, all_err := fs::read_file(%s)
if all_err.kind != errors::Error::NotFound { return 1 }
strings::free(&all)
preview, preview_err := fs::read_preview(%s, 3)
if preview_err.kind != errors::Error::NotFound { return 2 }
strings::free(&preview)
file, open_err := fs::open_read(%s)
if open_err.kind != errors::Error::NotFound { return 3 }
if fs::close(&file) != nil { return 4 }
`, utf8Literal(missing), utf8Literal(missing), utf8Literal(missing)))
}

func TestEditorAtomicSaveAllocationFailuresDoNotLeakTemporaryFiles(t *testing.T) {
	dir := t.TempDir()
	file := filepath.Join(dir, "save")
	if err := os.WriteFile(file, []byte("original"), 0600); err != nil {
		t.Fatal(err)
	}
	header, err := filepath.Abs("testdata/allocation_control.h")
	if err != nil {
		t.Fatal(err)
	}
	runEditorSource(t, `use "std::fs"
@[default_cc="c"]
extern {
fn _test_fail_after(count i64)
fn _test_live_allocations() i64
}`, fmt.Sprintf(`
attempt i64 := 0
succeeded := false
while attempt < 64 and not succeeded {
    _test_fail_after(attempt)
    replaced, err := fs::write_atomic(%s, "saved")
    _test_fail_after(-1)
    if not replaced and err.kind != errors::Error::AllocationFailed { return 1 }
    if _test_live_allocations() != 0 { return 2 }
    succeeded = replaced
    attempt = attempt + 1
}
if not succeeded { return 3 }
`, utf8Literal(file)), "-include", header)
	entries, err := os.ReadDir(dir)
	if err != nil || len(entries) != 1 {
		t.Fatalf("temporary files leaked during allocation failures: %v, %v", entries, err)
	}
}

func TestEditorAtomicSaveReportsBeforeAndAfterReplacementFailures(t *testing.T) {
	if runtime.GOOS == "windows" {
		t.Skip("POSIX rename and directory-sync failure injection")
	}
	header, err := filepath.Abs("testdata/save_control.h")
	if err != nil {
		t.Fatal(err)
	}
	for _, mode := range []int{1, 2, 3} {
		t.Run(fmt.Sprint(mode), func(t *testing.T) {
			dir := t.TempDir()
			file := filepath.Join(dir, "save")
			kind := "FileSyncFailed"
			if mode == 2 {
				kind = "FileRenameFailed"
			}
			runEditorSource(t, `use "std::fs"
@[default_cc="c"]
extern {
fn _test_save_mode(mode i32)
}`, fmt.Sprintf(`
if fs::write_file(%s, "original") != nil { return 1 }
_test_save_mode(%d)
replaced, err := fs::write_atomic(%s, "saved")
if replaced != %t or err.kind != errors::Error::%s { return 2 }
`, utf8Literal(file), mode, utf8Literal(file), mode == 3, kind), "-include", header)
			want := "original"
			if mode == 3 {
				want = "saved"
			}
			data, err := os.ReadFile(file)
			if err != nil || string(data) != want {
				t.Fatalf("destination after failure: %q, %v", data, err)
			}
			entries, err := os.ReadDir(dir)
			if err != nil || len(entries) != 1 {
				t.Fatalf("temporary file leaked: %v, %v", entries, err)
			}
		})
	}
}

func TestEditorSubprocessBackpressureTerminationAndBrokenPipeAreChecked(t *testing.T) {
	executable, err := os.Executable()
	if err != nil {
		t.Fatal(err)
	}
	statusCheck := "if status.signal != 15 { return 8 }"
	if runtime.GOOS == "windows" {
		statusCheck = "if status.code != 1 { return 8 }"
	}
	runEditorSource(t, `use "std::process"`, fmt.Sprintf(`
arguments string[3] := .{"-test.run=^TestEditorNativeChild$", "--", "hold"}
child, err := process::spawn(%s, arguments[:], process::SpawnOptions{pipe_stdin: true, pipe_stdout: true})
if err != nil { return 1 }
defer process::close_child(&child)
buffer u8[5] := .{}
received u64 := 0
attempts u64 := 0
while received < 5 and attempts < 10000000 {
    portion := buffer[received:5] @fail
    count, eof, blocked, read_err := process::read_stdout(&child, portion)
    if read_err != nil or eof { return 2 }
    received = received + count
    attempts = attempts + 1
}
if received != 5 { return 3 }
pending, done, poll_err := process::poll(&child)
if poll_err != nil or done { return 4 }
blocked := false
attempts = 0
while not blocked and attempts < 1024 {
    count, would_block, write_err := process::write_stdin(&child, %s)
    if write_err != nil { return 5 }
    blocked = would_block
    attempts = attempts + 1
}
if not blocked { return 6 }
if process::terminate(&child) != nil { return 7 }
status, wait_err := process::wait(&child)
%s
if wait_err != nil { return 9 }
count, would_block, broken_err := process::write_stdin(&child, "after exit")
if broken_err.kind != errors::Error::StreamWriteFailed { return 10 }
if process::close_child(&child) != nil or process::close_child(&child) != nil { return 11 }
`, utf8Literal(executable), utf8Literal(strings.Repeat("x", 4096)), statusCheck))
}

func TestEditorSubprocessSpawnFailuresValidateArgumentsAndReleaseAllocations(t *testing.T) {
	header, err := filepath.Abs("testdata/allocation_control.h")
	if err != nil {
		t.Fatal(err)
	}
	executable, err := os.Executable()
	if err != nil {
		t.Fatal(err)
	}
	runEditorSource(t, `use "std::process"
@[default_cc="c"]
extern {
fn _test_fail_after(count i64)
fn _test_live_allocations() i64
}`, fmt.Sprintf(`
bad_arguments string[1] := .{"bad\0argument"}
invalid, invalid_err := process::spawn("unused", bad_arguments[:], process::SpawnOptions{})
if invalid_err.kind != errors::Error::EmbeddedZero { return 1 }
arguments string[3] := .{"-test.run=^TestEditorNativeChild$", "--", "hold"}
missing, missing_err := process::spawn(%s, arguments[:], process::SpawnOptions{pipe_stdin: true, pipe_stdout: true, pipe_stderr: true})
if missing_err.kind != errors::Error::ProcessSpawnFailed { return 2 }
if _test_live_allocations() != 0 { return 3 }
attempt i64 := 0
succeeded := false
while attempt < 32 and not succeeded {
    _test_fail_after(attempt)
    child, err := process::spawn(%s, arguments[:], process::SpawnOptions{pipe_stdin: true, pipe_stdout: true, pipe_stderr: true})
    _test_fail_after(-1)
    if err != nil and err.kind != errors::Error::AllocationFailed { return 4 }
    succeeded = err == nil
    if process::close_child(&child) != nil { return 5 }
    if _test_live_allocations() != 0 { return 6 }
    attempt = attempt + 1
}
if not succeeded { return 7 }
`, utf8Literal(filepath.Join(t.TempDir(), "missing-executable")), utf8Literal(executable)), "-include", header)
}

func TestEditorBoundedReadsDistinguishEmptyBufferAndEOF(t *testing.T) {
	file := filepath.Join(t.TempDir(), "short")
	if err := os.WriteFile(file, []byte("abc"), 0600); err != nil {
		t.Fatal(err)
	}
	runEditorSource(t, `use "std::fs"`, fmt.Sprintf(`
preview, preview_err := fs::read_preview(%s, 0)
if preview_err != nil or strings::len(strings::view(&preview)) != 0 { return 1 }
strings::free(&preview)
file, open_err := fs::open_read(%s)
if open_err != nil { return 2 }
defer fs::close(&file)
empty u8[0] := .{}
zero, zero_eof, zero_err := fs::read(&file, empty[:])
if zero_err != nil or zero != 0 or zero_eof { return 3 }
buffer u8[5] := .{}
count, eof, read_err := fs::read(&file, buffer[:])
if read_err != nil or count != 3 or not eof or buffer[2] != 99 { return 4 }
again, again_eof, again_err := fs::read(&file, buffer[:])
if again_err != nil or again != 0 or not again_eof { return 5 }
`, utf8Literal(file), utf8Literal(file)))
}

func TestEditorFilePrimitivesCreateSyncRenameAndCleanUp(t *testing.T) {
	dir := t.TempDir()
	file := filepath.Join(dir, "new")
	syncCheck := "if sync_err != nil { return 6 }"
	if runtime.GOOS == "windows" {
		syncCheck = "if sync_err.kind != errors::Error::Unsupported { return 6 }"
	}
	runEditorSource(t, `use "std::fs"`, fmt.Sprintf(`
file, name, err := fs::create_temp(%s)
if err != nil { return 1 }
defer strings::free(&name)
if fs::write(&file, "new") != nil or fs::sync(&file) != nil { return 2 }
if fs::close(&file) != nil { return 3 }
if fs::rename(strings::view(&name), %s) != nil { return 4 }
sync_err := fs::sync_parent(%s)
%s
if fs::remove(%s) != nil { return 7 }
remove_again := fs::remove(%s)
if remove_again.kind != errors::Error::NotFound { return 8 }
replaced, save_err := fs::write_atomic(%s, "first save")
if not replaced { return 9 }
`, utf8Literal(file), utf8Literal(file), utf8Literal(file), syncCheck, utf8Literal(file), utf8Literal(file), utf8Literal(file)))
	data, err := os.ReadFile(file)
	if err != nil || string(data) != "first save" {
		t.Fatalf("new-file atomic save: %q, %v", data, err)
	}
	entries, err := os.ReadDir(dir)
	if err != nil || len(entries) != 1 {
		t.Fatalf("same-directory temporary cleanup: %v, %v", entries, err)
	}
}

func TestEditorSymlinksCanonicalizeAndSaveWithoutReplacingTheLink(t *testing.T) {
	if runtime.GOOS == "windows" {
		t.Skip("creating Windows symlinks requires host privileges")
	}
	dir := t.TempDir()
	canonical, err := filepath.EvalSymlinks(dir)
	if err != nil {
		t.Fatal(err)
	}
	target := filepath.Join(dir, "target")
	link := filepath.Join(dir, "link")
	alias := filepath.Join(dir, "alias")
	identitySource := filepath.Join(dir, "identity-source")
	dangling := filepath.Join(dir, "dangling")
	if err := os.WriteFile(target, []byte("original"), 0640); err != nil {
		t.Fatal(err)
	}
	if err := os.Symlink("target", link); err != nil {
		t.Fatal(err)
	}
	if err := os.Symlink("missing", dangling); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(identitySource, []byte("original"), 0640); err != nil {
		t.Fatal(err)
	}
	if err := os.Link(identitySource, alias); err != nil {
		t.Fatal(err)
	}
	runEditorSource(t, "use \"std::fs\"\nuse \"std::path\"", fmt.Sprintf(`
target, target_err := fs::stat(%s)
link, link_err := fs::lstat(%s)
followed, followed_err := fs::stat(%s)
if target_err != nil or link_err != nil or followed_err != nil { return 1 }
if link.kind != fs::FileType::Symlink or fs::same_file(link, target) or not fs::same_file(followed, target) { return 2 }
alias, alias_err := fs::stat(%s)
identity_source, identity_err := fs::stat(%s)
if alias_err != nil or identity_err != nil or not fs::same_file(alias, identity_source) { return 3 }
resolved, resolve_err := path::canonicalize(%s)
if resolve_err != nil or strings::view(&resolved) != %s { return 4 }
strings::free(&resolved)
dangling, dangling_err := path::canonicalize_for_create(%s)
if dangling_err.kind != errors::Error::PathFailed { return 5 }
strings::free(&dangling)
replaced, save_err := fs::write_atomic(%s, "saved")
if save_err != nil or not replaced { return 6 }
`, utf8Literal(target), utf8Literal(link), utf8Literal(link), utf8Literal(alias), utf8Literal(identitySource), utf8Literal(link), utf8Literal(filepath.Join(canonical, "target")), utf8Literal(dangling), utf8Literal(link)))
	name, err := os.Readlink(link)
	if err != nil || name != "target" {
		t.Fatalf("symlink was replaced: %q, %v", name, err)
	}
	data, err := os.ReadFile(target)
	if err != nil || string(data) != "saved" {
		t.Fatalf("symlink target not saved: %q, %v", data, err)
	}
	data, err = os.ReadFile(alias)
	if err != nil || string(data) != "original" {
		t.Fatalf("hard-link identity fixture changed: %q, %v", data, err)
	}
}

func TestEditorByteAndPreviewAllocationFailuresReturnEmptyOwners(t *testing.T) {
	header, err := filepath.Abs("testdata/allocation_control.h")
	if err != nil {
		t.Fatal(err)
	}
	file := filepath.Join(t.TempDir(), "data")
	if err := os.WriteFile(file, []byte("abc"), 0600); err != nil {
		t.Fatal(err)
	}
	runEditorSource(t, "use \"std::fs\"\nuse \"std::lists\"\n"+`@[default_cc="c"]
extern {
fn _test_fail_after(count i64)
fn _test_live_allocations() i64
}`, fmt.Sprintf(`
empty u8[0] := .{}
empty_string, empty_err := strings::from_bytes(empty[:])
if empty_err != nil or strings::view(&empty_string) != "" { return 1 }
strings::free(&empty_string)
empty_bytes, empty_bytes_err := strings::to_bytes("")
if empty_bytes_err != nil or lists::count<u8>(&empty_bytes) != 0 { return 2 }
lists::free<u8>(&empty_bytes)
attempt i64 := 0
succeeded := false
while attempt < 16 and not succeeded {
    _test_fail_after(attempt)
    bytes, err := strings::to_bytes("abc\0def")
    _test_fail_after(-1)
    if err != nil and err.kind != errors::Error::AllocationFailed { return 3 }
    if err != nil and lists::count<u8>(&bytes) != 0 { return 4 }
    succeeded = err == nil
    lists::free<u8>(&bytes)
    if _test_live_allocations() != 0 { return 5 }
    attempt = attempt + 1
}
if not succeeded { return 6 }
attempt = 0
succeeded = false
while attempt < 16 and not succeeded {
    _test_fail_after(attempt)
    preview, err := fs::read_preview(%s, 4)
    _test_fail_after(-1)
    if err != nil and err.kind != errors::Error::AllocationFailed { return 7 }
    if err != nil and strings::view(&preview) != "" { return 8 }
    if err == nil and strings::view(&preview) != "abc" { return 9 }
    succeeded = err == nil
    strings::free(&preview)
    if _test_live_allocations() != 0 { return 10 }
    attempt = attempt + 1
}
if not succeeded { return 11 }
`, utf8Literal(file)), "-include", header)
}

func TestEditorPermissionDeniedIsNotReportedAsNotFound(t *testing.T) {
	if runtime.GOOS == "windows" || os.Geteuid() == 0 {
		t.Skip("requires POSIX mode enforcement for a non-root caller")
	}
	file := filepath.Join(t.TempDir(), "denied")
	if err := os.WriteFile(file, []byte("secret"), 0000); err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() {
		if err := os.Chmod(file, 0600); err != nil {
			t.Error(err)
		}
	})
	runEditorSource(t, `use "std::fs"`, fmt.Sprintf(`
data, err := fs::read_file(%s)
if err.kind != errors::Error::PermissionDenied { return 1 }
strings::free(&data)
preview, preview_err := fs::read_preview(%s, 4)
if preview_err.kind != errors::Error::PermissionDenied { return 2 }
strings::free(&preview)
`, utf8Literal(file), utf8Literal(file)))
}
