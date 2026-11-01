# Programs, packages, and imports

Prerequisite: [your first program](../getting-started/first-program.md).

A package groups declarations: functions, types, and constants. Source files use
the `.t` extension. Begin each file with `package`, then put its `use` imports
before other declarations. Import names are local to the importing file.

An import can have a short alias:

```telia
package main

use output "std::io"

fn main() i32 {
    // A line comment explains the next statement.
    output::println("One package, one program")
    return 0
}
```

Expected output:

```text
One package, one program
```

`::` accesses a package member; `.` accesses a value's field. Do not replace
`std::io` with a filesystem path such as `std/io`.

You can build one file with `telia build main.t` or a directory package with
`telia build .`. Use a directory build when the program has declarations in
multiple files. Each file in the package declares the same package name.

Import prefixes identify where to look:

| Prefix | Meaning | Example |
| --- | --- | --- |
| `std` | Standard library | `use "std::strings"` |
| `pkg` | Application packages relative to the build's root directory | `use "pkg::helpers"` |
| `vendor` | Native vendor packages | `use "vendor::sdl"` |

For an application with `main.t` and `helpers/helpers.t`, build the application
directory and import `pkg::helpers` from `main.t`. Keep helpers in their own
package directory rather than mixing package names in one directory.

**Try it:** change the alias from `output` to `console`, including the call.

Reference: [program structure](../reference/programs.md).
Next: [variables and basic types](values.md).
