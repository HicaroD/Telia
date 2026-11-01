# Cookbook

These examples show complete tasks. For underlying language rules use the
[reference](../reference/README.md); for API contracts use the
[standard library](../stdlib/README.md).

Commands below run from the repository root after `make build`. On Windows,
choose `.exe` output paths and invoke the compiler/executables accordingly.

| Task | Program | Contracts |
| --- | --- | --- |
| Build a Unicode-aware command-line tool | [Walkthrough](command-line-programs.md), [text_cli.t](../../examples/text_cli.t) | Arguments, environment, files, streams, UTF-8, ownership |
| Read and write binary files | [files.t](../../examples/files.t) | [Whole-file I/O](../stdlib/files.md) |
| Copy stdin to stdout | [stdio.t](../../examples/stdio.t) | [Checked standard streams](../stdlib/standard-io.md) |
| Read arguments and environment | [process_args.t](../../examples/process_args.t), [process_env.t](../../examples/process_env.t) | [Process inputs](../stdlib/process.md) |
| Search, trim, and split borrowed text | [string_utilities.t](../../examples/string_utilities.t) | [Strings and cleanup](../stdlib/strings-and-memory.md) |
| Parse and format text | [strconv.t](../../examples/strconv.t) | [Checked conversions](../stdlib/string-conversions.md) |
| Decode and iterate UTF-8 | [utf8.t](../../examples/utf8.t) | [Strict UTF-8](../stdlib/unicode.md) |
| Work with arrays, slices, and lists | [collections.t](../../examples/collections.t) | [Collections](../reference/collections.md) |
| Classify and match errors | [error_enums.t](../../examples/error_enums.t) | [Errors](../reference/errors.md) |
| Reuse code with generic types | [generics.t](../../examples/generics.t) | [Generics](../reference/generics.md) |
| Build an interactive terminal application | [terminal_echo.t](../../examples/terminal_echo.t), [terminal_raw.t](../../examples/terminal_raw.t), [terminal_screen.t](../../examples/terminal_screen.t) | [Terminal lifecycle](../stdlib/terminal.md) |
| Build a graphical native application | [Snake](../../examples/snake/README.md) | [Native packages and deployment](../stdlib/native-libraries.md) |
| Write unit tests | [Testing walkthrough](../getting-started/testing.md), [example package](../../examples/testing/math) | [Testing APIs](../stdlib/testing.md) |

For example:

```sh
./telia build examples/string_utilities.t -o bin/string_utilities
./bin/string_utilities
```

Examples involving native libraries or a controlling terminal have additional
setup described by their linked guide. Not every file in `examples/` is an
introductory or successful program: some are diagnostic experiments. Use this
curated list or the tested guide examples when learning.

For metadata, bounded file reads, canonical paths, atomic saves, child processes,
and display-cell widths, see [native CLI recipes and contracts](../stdlib/native-cli-apis.md).
