package integration

import (
	"archive/tar"
	"compress/gzip"
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

func vendorLibraryFixture(t *testing.T) (string, string) {
	t.Helper()
	root := t.TempDir()
	pkg := filepath.Join(root, "answer")
	if err := os.MkdirAll(pkg, 0755); err != nil {
		t.Fatal(err)
	}
	write := func(name, text string) {
		t.Helper()
		if err := os.WriteFile(filepath.Join(pkg, name), []byte(text), 0644); err != nil {
			t.Fatal(err)
		}
	}
	write("answer.c", "int vendor_answer(void) { return 42; }\n")
	lib := "libanswer.so"
	flags := []string{"-shared", "-fPIC", "-Wl,-soname,libanswer.so"}
	if runtime.GOOS == "darwin" {
		lib = "libanswer.dylib"
		flags = []string{"-dynamiclib", "-Wl,-install_name,@rpath/" + lib}
	} else if runtime.GOOS == "windows" {
		lib = "answer.dll"
		flags = []string{"-shared"}
	}
	flags = append(flags, filepath.Join(pkg, "answer.c"), "-o", filepath.Join(pkg, lib))
	if out, err := exec.Command("cc", flags...).CombinedOutput(); err != nil {
		t.Fatalf("build fixture: %v\n%s", err, out)
	}
	write(
		"bridge.c",
		"int vendor_answer(void);\nint bridge_answer(void) { return vendor_answer(); }\n",
	)
	write(
		"answer.t",
		"package answer\n@[default_cc=\"c\"]\nextern {\n @[link_name=\"bridge_answer\"]\n fn vendor_answer() i32\n}\n",
	)
	write(
		"native.json",
		`{"version":"1","targets":{"`+runtime.GOOS+"-"+runtime.GOARCH+`":{"sources":["bridge.c"],"libraries":["`+lib+`"],"runtime_files":["`+lib+`"]}}}`,
	)
	return root, pkg
}

// A real shared library exercises import resolution, native linking, and loader
// deployment together, without requiring SDL or network access in ordinary tests.
func TestVendorImportLinksAndDeploysSharedLibrary(t *testing.T) {
	root, pkg := vendorLibraryFixture(t)
	t.Setenv("T_VENDOR_PATH", root)
	source := filepath.Join(t.TempDir(), "main.t")
	if err := os.WriteFile(source, []byte("package main\nuse \"vendor::answer\"\nfn main() i32 {\n return answer::vendor_answer() - 42\n}\n"), 0644); err != nil {
		t.Fatal(err)
	}
	binary, diags := compiler.CompileOnly(source)
	if len(diags.Diags) != 0 {
		t.Fatalf("compile vendor import: %v", diags.Diags)
	}
	t.Cleanup(func() { removeTestBinary(t, binary) })
	if err := os.Rename(pkg, pkg+"-unavailable"); err != nil {
		t.Fatal(err)
	}
	if out, err := exec.Command(binary).CombinedOutput(); err != nil {
		t.Fatalf("run deployed executable: %v\n%s", err, out)
	}
}

func TestVendorRequirementsAreTransitiveAndDeduplicated(t *testing.T) {
	root, _ := vendorLibraryFixture(t)
	wrapper := filepath.Join(root, "wrapper")
	if err := os.MkdirAll(wrapper, 0755); err != nil {
		t.Fatal(err)
	}
	for file, text := range map[string]string{
		"wrapper.t":   "package wrapper\nuse \"vendor::answer\"\nfn value() i32 {\n return answer::vendor_answer()\n}\n",
		"native.json": `{"version":"1","targets":{"` + runtime.GOOS + "-" + runtime.GOARCH + `":{}}}`,
	} {
		if err := os.WriteFile(filepath.Join(wrapper, file), []byte(text), 0644); err != nil {
			t.Fatal(err)
		}
	}
	t.Setenv("T_VENDOR_PATH", root)
	for _, imports := range []string{"", "use \"vendor::answer\"\n"} {
		source := filepath.Join(t.TempDir(), "main.t")
		if err := os.WriteFile(source, []byte("package main\n"+imports+"use \"vendor::wrapper\"\nfn main() i32 {\n return wrapper::value() - 42\n}\n"), 0644); err != nil {
			t.Fatal(err)
		}
		_, diags := compiler.CompileFile(source)
		if len(diags.Diags) != 0 {
			t.Fatalf("transitive native import: %v", diags.Diags)
		}
	}
}

func TestVendorImportCannotEscapeVendorRoot(t *testing.T) {
	root := t.TempDir()
	outside := filepath.Join(root, "outside")
	if err := os.MkdirAll(outside, 0755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(outside, "outside.t"), []byte("package outside\n"), 0644); err != nil {
		t.Fatal(err)
	}
	t.Setenv("T_VENDOR_PATH", filepath.Join(root, "vendor"))
	source := filepath.Join(root, "main.t")
	if err := os.WriteFile(source, []byte("package main\nuse \"vendor::..::outside\"\nfn main() i32 {\n return 0\n}\n"), 0644); err != nil {
		t.Fatal(err)
	}
	_, diags := compiler.CompileOnly(source)
	if len(diags.Diags) == 0 {
		t.Fatal("vendor import escaped its root")
	}
}

func TestVendorRuntimeNameConflictsAreRejected(t *testing.T) {
	root := t.TempDir()
	for _, name := range []string{"left", "right"} {
		dir := filepath.Join(root, name)
		if err := os.MkdirAll(dir, 0755); err != nil {
			t.Fatal(err)
		}
		files := map[string]string{
			name + ".t":     "package " + name + "\nfn value() i32 {\n return 0\n}\n",
			"collision.dll": name,
			"native.json":   `{"version":"1","targets":{"` + runtime.GOOS + "-" + runtime.GOARCH + `":{"runtime_files":["collision.dll"]}}}`,
		}
		for file, text := range files {
			if err := os.WriteFile(filepath.Join(dir, file), []byte(text), 0644); err != nil {
				t.Fatal(err)
			}
		}
	}
	t.Setenv("T_VENDOR_PATH", root)
	source := filepath.Join(t.TempDir(), "main.t")
	if err := os.WriteFile(source, []byte("package main\nuse \"vendor::left\"\nuse \"vendor::right\"\nfn main() i32 {\n return left::value() + right::value()\n}\n"), 0644); err != nil {
		t.Fatal(err)
	}
	binary, diags := compiler.CompileOnly(source)
	if binary != "" {
		t.Cleanup(func() { removeTestBinary(t, binary) })
	}
	if !strings.Contains(fmt.Sprint(diags.Diags), "runtime file name conflict") {
		t.Fatalf("want runtime name conflict, got %v", diags.Diags)
	}
}

func TestVendorIncompatibleArtifactReportsTargetAndSetup(t *testing.T) {
	root := t.TempDir()
	dir := filepath.Join(root, "broken")
	if err := os.MkdirAll(dir, 0755); err != nil {
		t.Fatal(err)
	}
	for file, text := range map[string]string{
		"broken.t":    "package broken\nfn value() i32 {\n return 0\n}\n",
		"wrong.a":     "not a native library\n",
		"native.json": `{"version":"1","targets":{"` + runtime.GOOS + "-" + runtime.GOARCH + `":{"libraries":["wrong.a"]}}}`,
	} {
		if err := os.WriteFile(filepath.Join(dir, file), []byte(text), 0644); err != nil {
			t.Fatal(err)
		}
	}
	t.Setenv("T_VENDOR_PATH", root)
	source := filepath.Join(t.TempDir(), "main.t")
	if err := os.WriteFile(source, []byte("package main\nuse \"vendor::broken\"\nfn main() i32 {\n return broken::value()\n}\n"), 0644); err != nil {
		t.Fatal(err)
	}
	_, diags := compiler.CompileOnly(source)
	message := fmt.Sprint(diags.Diags)
	if !strings.Contains(message, runtime.GOOS+"-"+runtime.GOARCH) ||
		!strings.Contains(message, "setup") {
		t.Fatalf("want target and recovery instructions, got %s", message)
	}
}

func TestSDL3ExampleDrawsAndCleansUp(t *testing.T) {
	if os.Getenv("TELIA_TEST_SDL") != "1" {
		t.Skip("run SDL setup, then set TELIA_TEST_SDL=1")
	}
	root, err := filepath.Abs("../../base/vendor")
	if err != nil {
		t.Fatal(err)
	}
	t.Setenv("T_VENDOR_PATH", root)
	binary, diags := compiler.CompileOnly("../../examples/sdl.t")
	if len(diags.Diags) != 0 {
		t.Fatalf("compile SDL example: %v", diags.Diags)
	}
	t.Cleanup(func() { removeTestBinary(t, binary) })
	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()
	cmd := exec.CommandContext(ctx, binary)
	cmd.Env = append(
		os.Environ(),
		"SDL_VIDEO_DRIVER=dummy",
		"SDL_RENDER_DRIVER=software",
		"TELIA_SDL_EXAMPLE_FRAMES=3",
	)
	if out, err := cmd.CombinedOutput(); err != nil {
		t.Fatalf("SDL example: %v\n%s", err, out)
	}
}

func TestSDL3ExampleStaysOpenUntilQuitByDefault(t *testing.T) {
	if os.Getenv("TELIA_TEST_SDL") != "1" {
		t.Skip("run SDL setup, then set TELIA_TEST_SDL=1")
	}
	root, err := filepath.Abs("../../base/vendor")
	if err != nil {
		t.Fatal(err)
	}
	t.Setenv("T_VENDOR_PATH", root)
	binary, diags := compiler.CompileOnly("../../examples/sdl.t")
	if len(diags.Diags) != 0 {
		t.Fatalf("compile SDL example: %v", diags.Diags)
	}
	t.Cleanup(func() { removeTestBinary(t, binary) })
	cmd := exec.Command(binary)
	for _, env := range os.Environ() {
		if !strings.HasPrefix(env, "TELIA_SDL_EXAMPLE_FRAMES=") {
			cmd.Env = append(cmd.Env, env)
		}
	}
	cmd.Env = append(cmd.Env, "SDL_VIDEO_DRIVER=dummy", "SDL_RENDER_DRIVER=software")
	if err := cmd.Start(); err != nil {
		t.Fatal(err)
	}
	done := make(chan error, 1)
	go func() { done <- cmd.Wait() }()
	select {
	case err := <-done:
		t.Fatalf("example exited without a quit event: %v", err)
	case <-time.After(3 * time.Second):
		_ = cmd.Process.Kill()
		<-done
	}
}

func TestSDL3QuitEventsCrossBindingBoundary(t *testing.T) {
	if os.Getenv("TELIA_TEST_SDL") != "1" {
		t.Skip("run SDL setup, then set TELIA_TEST_SDL=1")
	}
	root, err := filepath.Abs("../../base/vendor")
	if err != nil {
		t.Fatal(err)
	}
	t.Setenv("T_VENDOR_PATH", root)
	dir := t.TempDir()
	helper := filepath.Join(dir, "events.c")
	if err := os.WriteFile(helper, []byte("#include <SDL3/SDL.h>\nbool push_quit(void) { SDL_Event event; SDL_zero(event); event.type = SDL_EVENT_QUIT; return SDL_PushEvent(&event); }\n"), 0644); err != nil {
		t.Fatal(err)
	}
	source := `package main
use "vendor::sdl"
@[default_cc="c"]
extern {
  fn push_quit() bool
}
fn main() i32 {
  sdl::set_main_ready()
  if sdl::init(sdl::INIT_VIDEO) == false {
    return 1
  }
  defer sdl::quit()
  event := sdl::Event{type: 0}
  while sdl::poll_event(&event) {
  }
  if push_quit() == false {
    return 3
  }
  if sdl::poll_event(&event) == false {
    return 4
  }
  if event.type != sdl::EVENT_QUIT {
    return 2
  }
  if sdl::poll_event(&event) {
    return 5
  }
  return 0
}
`
	path := filepath.Join(dir, "main.t")
	if err := os.WriteFile(path, []byte(source), 0644); err != nil {
		t.Fatal(err)
	}
	binary, diags := compiler.CompileOnlyWithCFlags(path, helper)
	if len(diags.Diags) != 0 {
		t.Fatalf("compile SDL event test: %v", diags.Diags)
	}
	t.Cleanup(func() { removeTestBinary(t, binary) })
	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()
	cmd := exec.CommandContext(ctx, binary)
	cmd.Env = append(os.Environ(), "SDL_VIDEO_DRIVER=dummy")
	if out, err := cmd.CombinedOutput(); err != nil {
		t.Fatalf("SDL quit event: %v\n%s", err, out)
	}
}

func TestSDL3ReleaseBundleWorksWithoutConfiguration(t *testing.T) {
	if os.Getenv("TELIA_TEST_SDL") != "1" {
		t.Skip("run SDL setup, then set TELIA_TEST_SDL=1")
	}
	archive := filepath.Join(t.TempDir(), "telia.tar.gz")
	cmd := exec.Command("python3", "../../scripts/package-release.py", "--output", archive)
	if out, err := cmd.CombinedOutput(); err != nil {
		t.Fatalf("package release: %v\n%s", err, out)
	}
	file, err := os.Open(archive)
	if err != nil {
		t.Fatal(err)
	}
	defer func() { _ = file.Close() }() // Read-only archive handle.
	gz, err := gzip.NewReader(file)
	if err != nil {
		t.Fatal(err)
	}
	defer func() { _ = gz.Close() }() // Read-only decompressor.
	bundle := t.TempDir()
	reader := tar.NewReader(gz)
	for {
		header, err := reader.Next()
		if err == io.EOF {
			break
		}
		if err != nil {
			t.Fatal(err)
		}
		if !filepath.IsLocal(header.Name) {
			t.Fatalf("invalid bundle path: %s", header.Name)
		}
		path := filepath.Join(bundle, header.Name)
		if header.Typeflag == tar.TypeDir {
			if err := os.MkdirAll(path, 0755); err != nil {
				t.Fatal(err)
			}
			continue
		}
		if header.Typeflag != tar.TypeReg {
			t.Fatalf("unexpected archive member: %s", header.Name)
		}
		if err := os.MkdirAll(filepath.Dir(path), 0755); err != nil {
			t.Fatal(err)
		}
		data, err := io.ReadAll(reader)
		if err != nil {
			t.Fatal(err)
		}
		if err := os.WriteFile(path, data, os.FileMode(header.Mode)); err != nil {
			t.Fatal(err)
		}
	}
	name := "telia"
	if runtime.GOOS == "windows" {
		name += ".exe"
	}
	source, err := filepath.Abs("../../examples/sdl.t")
	if err != nil {
		t.Fatal(err)
	}
	output := filepath.Join(t.TempDir(), "example")
	cmd = exec.Command(filepath.Join(bundle, name), "build", source, "-o", output, "--color=never")
	cmd.Dir = t.TempDir()
	for _, env := range os.Environ() {
		if !strings.HasPrefix(env, "T_STD_PATH=") && !strings.HasPrefix(env, "T_VENDOR_PATH=") &&
			!strings.HasPrefix(env, "XDG_CONFIG_HOME=") {
			cmd.Env = append(cmd.Env, env)
		}
	}
	cmd.Env = append(cmd.Env, "XDG_CONFIG_HOME="+t.TempDir())
	if out, err := cmd.CombinedOutput(); err != nil {
		t.Fatalf("build with clean release install: %v\n%s", err, out)
	}
	if runtime.GOOS == "windows" {
		output += ".exe"
	}
	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()
	cmd = exec.CommandContext(ctx, output)
	cmd.Env = append(
		os.Environ(),
		"SDL_VIDEO_DRIVER=dummy",
		"SDL_RENDER_DRIVER=software",
		"TELIA_SDL_EXAMPLE_FRAMES=3",
	)
	if out, err := cmd.CombinedOutput(); err != nil {
		t.Fatalf("run bundled SDL example: %v\n%s", err, out)
	}
}
