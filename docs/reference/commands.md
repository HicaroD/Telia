# Compiler commands

Commands use `telia` from `PATH`. In a source checkout, substitute `./telia`
after `make build`. A C compiler must be installed.

## Build

```sh
telia build
telia build main.t -o hello
telia build path/to/package -release -o application
telia build path/to/package -debug --color=never
```

The path defaults to the current directory. A directory builds a package; a
file starts from that file. `-o` selects the executable path. Without it, the
compiler derives a source name in the current directory. On Windows use an
appropriate `.exe` output name.

Debug is the default; `-debug` requests it explicitly. `-release` selects the
release optimization configuration. These flags do not remove language bounds
checks or make invalid pointer use safe.

## Test

Run `telia test` from the package directory. It includes that package's named
tests, not tests in subdirectories or dependencies. There is no recursive
all-packages test command documented here. See [testing](../getting-started/testing.md)
for assertions, cleanup, discovery, output, and exit statuses.

## Diagnostics and configuration

- `telia help` lists command usage.
- `telia env` prints effective library paths and the configuration location.
- `--color=auto|always|never` controls build/test diagnostic colors; `auto`
  respects redirected output and `NO_COLOR`.
- `T_STD_PATH` and `T_VENDOR_PATH` override library discovery. A configured
  standard path can also override bundled sources. Keep compiler and library
  sources consistent; inspect `telia env` when imports fail.

For installation layouts and native-library deployment, see
[installation](../getting-started/installation.md) and
[native packages](../stdlib/native-libraries.md).
