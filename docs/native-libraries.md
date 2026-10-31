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

`vendor::sdl` is the Telia-facing layer. Its text APIs accept native `string`,
including dynamically built text and slices. Applications do not need to import
`std::ffi::c` or manage temporary C-string buffers. The ABI-faithful generated
bindings are available separately:

```rust
use "vendor::sdl"
use "std::io"

window := sdl::create_window(title, 640, 480, 0) @fail
defer sdl::destroy_window(window)
renderer := sdl::create_renderer(window) @fail
defer sdl::destroy_renderer(renderer)
sdl::render_debug_text(renderer, 10.0, 10.0, message) @fail
sdl::render_debug_i32(renderer, 10.0, 30.0, score) @fail
io::println(sdl::get_error())
```

This snippet assumes SDL is already initialized. In resource-owning applications,
prefer `@catch` and normal returns when failures must run deferred cleanup;
`@fail` terminates without unwinding.

### Native adapter contracts

- Text inputs are copied to temporary NUL-terminated storage and freed before
  returning. Embedded NULs and allocation failures return `error` without calling
  SDL. These conversions currently allocate even for literals; no implicit
  conversion was added to the compiler.
- Reviewed text operations that return SDL `bool` now return `error`; text
  operations that create handles return `(handle, error)`. Other forwarding
  functions retain SDL's result conventions: for example, `poll_event` is a
  predicate, and drawing functions such as `render_clear` still return `bool`.
- Errors identify the failing operation using static, lifetime-independent text.
  For SDL's detailed diagnostic, read `get_error()` promptly. It is a borrowed
  Telia `string` and can be invalidated by another SDL call. Conversion errors
  instead contain the standard string library's specific failure message.
- Borrowed text such as `get_window_title` is not freed by the caller and must not
  outlive its SDL storage. Use `strings::clone` to retain an independent copy.
  Nullable getters return `OptionalString{value, present}`; an absent pointer is
  distinct from a present empty string. `get_hint(name)` returns
  `(OptionalString, error)` because its input conversion can fail.
- Owned results such as `get_clipboard_text` return `(strings::Owned, error)`.
  The adapter copies the text and releases the original with **SDL's** allocator;
  callers release the returned owner with `strings::free`.
- `create_renderer(window)` selects SDL's default renderer (a null name).
  `create_renderer_named(window, name)` passes an explicit name, including an
  empty string. These meanings are not conflated.
- `get_window_size(window)` returns `(Size, error)` with native `i32` width and
  height instead of C output slots.
- `io_from_bytes(bytes)` borrows a nonempty writable `u8[]`; the storage must stay
  valid until the stream is closed. `read_bytes(stream, destination)` bounds the
  read by the destination slice, preserves embedded zeros, and returns the byte
  count. An empty destination returns zero without calling SDL. A short read is
  not itself an error; inspect SDL's IO status as appropriate.
- `create_thread(entry, name, data)` returns `(*Thread, error)` and
  `wait_thread(thread)` returns its `i32` status. Entries have signature
  `@[cc="c"] fn(data rawptr) i32`: no C scalar aliases or `SDLCALL` annotation
  are required. An internal trampoline handles SDL's convention and frees its
  state, including on creation failure. Keep userdata alive until the worker
  finishes; this adds no closures, borrow checking, dispatch, or synchronization.
- Simple scalar arguments/results use native types on supported targets. Header
  records, unions, enums, and opaque handles retain their original identity and
  layout through aliases. **Their pointer fields and callback fields remain
  low-level**; the facade is not a memory-safety layer.

### Raw escape hatch and coverage

```rust
use raw "vendor::sdl::raw"
```

Use the raw layer for APIs without a native adapter, including arbitrary C
varargs, target-dependent wide characters, many callback registrations, retained
text, and complex output/buffer ownership. Existing raw callers can migrate by
changing their import to `use sdl "vendor::sdl::raw"`.

`coverage.json` lists all discovered functions, native forwarding/adaptation
policies, reasons for raw-only functions and constants, and header-owned types
with their original field types and callback signatures. Raw-only does not
mean unavailable: it means the facade deliberately does not guess ownership,
nullability, buffer lengths, or calling conventions. Coverage is not a promise
that every SDL API has a high-level adapter. Extend the reviewed policies and
add boundary tests when promoting an API.

`examples/sdl.t` opens and raises a window, polls events, draws a rectangle, and
stays open until a quit event. Window, event-pump, and rendering calls run on the
main thread. See [Snake](../examples/snake/README.md) for a playable Telia game.
For automated smoke checks only, set `TELIA_SDL_EXAMPLE_FRAMES` to a positive
frame count to exit automatically. Normal desktop runs should not set
`SDL_VIDEO_DRIVER=dummy` (or `SDL_VIDEODRIVER=dummy`): that driver deliberately
creates no visible window. Rebuild the executable after changing the example.

- Call `set_main_ready()`, then `init(INIT_VIDEO)`; pair successful initialization
  with `quit`. Initialization flags follow SDL rather than always enabling video.
- `create_window(title, width, height, flags)` returns `(*Window, error)` and
  `create_renderer(window)` returns `(*Renderer, error)`. Destroy
  the renderer before its window, then call `quit`. Typed handles are not
  interchangeable, and use after destruction is invalid.
- Pass a native `string` title. SDL copies it; the facade releases its temporary
  conversion automatically.
- Non-text drawing functions return `bool`; false indicates an SDL error. `get_error`
  returns borrowed SDL text; do not free it or assume it survives another SDL call.
- `poll_event(&event)` returns `bool` and writes a header-owned `Event` union.
  Initialize storage with `Event{type: 0}`. Check `event.type` before accessing
  its active member, such as `event.key`. Quit events have type `EVENT_QUIT`.
- Rendering uses SDL-faithful names such as `set_render_draw_color`,
  `render_clear`, `render_fill_rect`, and `render_present`. Rectangles use
  header-owned `FRect` fields `x`, `y`, `w`, and `h`. Ordinary Telia `f32` values
  and arithmetic can be passed to C `float` parameters.
- Raw callback declarations use `SDLCALL`. Keep callbacks and userdata alive until
  SDL stops invoking them. Ownership and threading requirements remain those of
  the corresponding upstream SDL API; the binding does not add automatic cleanup.

### Regenerating bindings

```sh
python3 scripts/generate-sdl.py --include base/vendor/sdl/artifacts/darwin-arm64/include
```

Regeneration requires Clang and the pinned SDL 3.2.28 headers; ordinary builds do
not run the generator. It emits `raw/sdl.t`, the native forwarding layer `sdl.t`,
and an inventory in `coverage.json`. Reviewed conversion policies live in
`scripts/sdl_facade.py`; handwritten native adapters and `raw/adapters.t` are not
overwritten by regeneration. `text.c` contains the small C ABI bridges.
SDL's headers own record layouts and native type identity, including nested
anonymous records accessed through `__typeof__`; no replacement layouts are emitted.

The inventory describes declarations visible on the regeneration host, not every
platform-specific API or macro. The current generated surface includes the Vulkan
integration header and explicit thread-creation macro signatures. Full platform,
type, and macro coverage auditing remains unfinished. Local runtime verification
is macOS ARM64; CI configuration alone is not proof of results on other targets.

Test the real SDL example after setup:

```sh
TELIA_TEST_SDL=1 go test ./tests/integration -run TestSDL3 -count=1 -timeout=30m
```

The tests use SDL's dummy video/software renderer. Ordinary `go test ./...` uses
a locally built shared-library fixture to verify native linking and relocatable
deployment without downloading SDL. A real display smoke check is still useful.

## Adding a native package

Create `base/vendor/<name>` with normal Telia bindings and a `native.json`:

Vendor subpackages inherit the nearest enclosing package's `native.json`, never
searching above the top-level vendor package. An explicit subpackage manifest
takes precedence. Shared manifests are collected once, including C sources,
even when an application imports both a facade and its raw subpackage.

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
