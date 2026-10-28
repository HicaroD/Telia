# Native vendor packages

SDL3 is the first native integration (SDL **3.2.28**). Native builds use
`use "vendor::sdl"` and automatically collect package-owned native requirements,
including transitive imports. No per-project linker flags are needed.

## Supported targets

| OS | CPU | Toolchain / minimum build environment |
| --- | --- | --- |
| Linux | x86-64 | GCC or Clang; glibc; X11 development headers; CI builds on Ubuntu 22.04 |
| macOS | Intel x86-64 | Apple Clang / Xcode command-line tools; macOS 11+ |
| macOS | Apple Silicon ARM64 | Apple Clang / Xcode command-line tools; macOS 11+ |
| Windows | x86-64 | MSYS2 UCRT64 MinGW-w64 GCC; Windows 10+ |

These are **native builds**, not cross-compilation support. Windows MSVC, Linux
musl, other CPUs, SDL extensions, mobile, and web targets are not supported yet.
Linux Wayland is disabled in the supplied recipe; X11/XWayland is required for
desktop windows. System graphics/audio drivers and OS libraries remain prerequisites.

## Using a bundled Telia release

Release archives created by `scripts/package-release.py` include the compiler,
standard library, SDL bindings, headers, and the target's shared library. Extract
the **whole archive**, keep `base/` beside `telia`, and run:

```sh
telia build examples/sdl.t -o bin/sdl-example
```

No separate SDL installation, setup command, or linker flags are needed. A C
compiler and the platform prerequisites below are still required. An existing
`T_STD_PATH` configuration or `T_VENDOR_PATH` environment override takes precedence
over bundled discovery; remove outdated overrides when moving installations.

Maintainers build these archives after the explicit setup step:

```sh
python3 scripts/package-release.py --output telia-<os>-<arch>.tar.gz
```

CI produces these archives for all four targets. It does not publish GitHub
releases automatically.

## Explicit setup (source builds and release maintainers)

Install Go, Python 3.9+, CMake 3.16+, Ninja, and the target C toolchain. For Linux
install X11 development packages; CI lists the full desktop build prerequisites in
`.github/workflows/native.yml`. On Windows run setup in MSYS2 UCRT64 with its GCC,
Python, CMake, and Ninja on `PATH`. No dependencies are installed automatically.

From the Telia checkout:

```sh
python3 scripts/setup-sdl.py
export T_VENDOR_PATH="$PWD/base/vendor"
telia build examples/sdl.t -o bin/sdl-example
```

The setup script downloads a pinned upstream archive, verifies its SHA-256, and
builds shared SDL3 into `base/vendor/sdl/artifacts/<os>-<arch>`. It never runs during
compilation. Source/build caches live under `.native-cache`; use `--cache` to
change that location. A previously downloaded archive can be supplied with
`--archive /path/to/SDL3-3.2.28.tar.gz` for offline setup. After setup, compiling
requires neither the network nor CMake/Ninja: only the C toolchain and installed
headers/libraries. Compiler artifacts must match the selected OS, CPU, and ABI.

Bindings default to `vendor` alongside the configured `T_STD_PATH` directory;
`T_VENDOR_PATH` overrides that location. Keep Telia's existing standard-library
configuration pointing at `base/std` as well.

The recipe pins source and build options, not a bit-for-bit reproducible binary:
SDK, compiler, and installed system dependencies also affect the result. CI
builds all four targets and uploads native artifacts, headers, and license notices
for reuse. Published release hosting is separate from this setup workflow.

## Runtime deployment

Telia copies package-declared runtime files beside the output executable. Linux
uses `$ORIGIN` lookup and macOS uses `@loader_path` with SDL's `@rpath` install
name. Windows uses its executable-directory DLL lookup. Distribute the executable
**with these files**, SDL's license notice, and any required toolchain runtimes;
copying only the executable is not supported. Do not mix builds of SDL in the same
output directory. Runtime file-name collisions between different packages are
rejected rather than silently overwriting dependencies.

The minimum runtime glibc version follows the environment used to build the Linux
artifact. Building on a newer distribution does not guarantee compatibility with
an older one. This workflow does not bundle OS drivers or automatically install
system dependencies. Missing native artifacts report the target and setup action;
the native linker diagnoses incompatible architectures/ABIs.

## SDL API and ownership

`examples/sdl.t` opens and raises a window, polls events, draws a rectangle, and
stays open until a quit event. Calls must run on the main thread.
For automated smoke checks only, set `TELIA_SDL_EXAMPLE_FRAMES` to a positive
frame count to exit automatically. Normal desktop runs should not set
`SDL_VIDEO_DRIVER=dummy` (or `SDL_VIDEODRIVER=dummy`): that driver deliberately
creates no visible window. Rebuild the executable after changing the example.

- `init` initializes video; pair successful initialization with `quit`.
- `create_window` and `create_renderer` return nullable, opaque `rawptr` handles.
  Destroy the renderer before its window, then call `quit`. Handles are not
  interchangeable despite sharing the Telia type; use after destruction is invalid.
- Pass a NUL-terminated `cstring` title; `std::strings::to_c` creates owned storage.
  SDL copies the title. Free the conversion with `free_c`.
- Drawing functions return `bool`; false indicates an SDL error. `get_error`
  returns borrowed SDL text; do not free it or assume it survives another SDL call.
- `poll_event` returns an SDL event type, or zero when the queue is empty.
  `is_quit_event` recognizes quit events. Event payloads, callbacks, and full SDL
  API coverage are not exposed. The package-local C adapter owns the `SDL_Event`
  union and `SDL_FRect` layout; Telia does not reproduce those layouts.

Test the real SDL example after setup:

```sh
TELIA_TEST_SDL=1 go test ./tests/integration -run TestSDL3 -count=1
```

The tests use SDL's dummy video/software renderer. Ordinary `go test ./...` uses
a locally built shared-library fixture to verify native linking and relocatable
deployment without downloading SDL. A real display smoke check is still useful.

## Adding a native package

Create `base/vendor/<name>` with normal Telia bindings and a `native.json`:

```json
{
  "version": "1.2.3",
  "targets": {
    "linux-amd64": {
      "include_dirs": ["artifacts/linux-amd64/include"],
      "sources": ["adapter.c"],
      "libraries": ["artifacts/linux-amd64/lib/libexample.so.1"],
      "runtime_files": ["artifacts/linux-amd64/lib/libexample.so.1"],
      "compile_flags": [],
      "link_flags": []
    }
  }
}
```

All file paths are package-relative. `sources` are compiled alongside generated C;
`include_dirs` provide header lookup for those adapters. Libraries follow objects
on the linker command line; use `link_flags` for system libraries/frameworks.
Shared libraries need relocatable SONAME/install names. Declare every non-system
runtime dependency in `runtime_files`, including versioned names required by the
loader. Package metadata is trusted build input, like a build script.

Use a small C adapter when C header types, macros, unions, or calling conventions
cannot safely be expressed as Telia extern prototypes. Include the upstream
headers in that adapter, not in generated Telia C with conflicting declarations.
Add a pinned, checksum-verified setup recipe, license notices, an example, and
compile/link/runtime CI coverage for every advertised target.

For an optional system-library override, set `T_VENDOR_PATH` to an alternate
bindings tree with a `native.json` that uses system `compile_flags` and
`link_flags` (for example `-I/usr/local/include`, `-L/usr/local/lib`, `-lSDL3`)
instead of packaged artifacts. Supply the same binding/adapter sources there.
That mode depends on the system library's compatible version and runtime loader
configuration; Telia does not install or deploy system libraries automatically.

raylib, SQLite, and libcurl remain subsequent integrations; they are not supplied
by this SDL3 milestone.
