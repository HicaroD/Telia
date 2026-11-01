# Your first program

First [install Telia and a C compiler](installation.md). Check that `telia help`
works in your terminal. Telia compiles through C, so the C compiler is required
even when you do not write any C.

Save this complete program as `main.t`:

```telia
package main

use "std::io"

fn main() i32 {
    io::println("Hello, Telia!")
    return 0
}
```

Expected output:

```text
Hello, Telia!
```

Build and run on macOS or Linux:

```sh
telia build main.t -o hello
./hello
```

On Windows (PowerShell):

```powershell
telia build main.t -o hello.exe
.\hello.exe
```

`package main` names the program's package. `use` makes another package available;
`io::println` calls its `println` function. `fn` declares a function. `main` is
the entry point, and `return 0` reports successful completion to the operating
system. `i32` is a signed 32-bit integer.

Statements normally end at a newline, not a semicolon. Put opening braces on the
same line as the declaration or condition. Text in double quotes is a borrowed
`string`; this literal has static storage and needs no cleanup. You will learn
about allocated text before using APIs that create it.

If a library cannot be found, run `telia env` and check the
[installation troubleshooting](installation.md#manual-installation-and-troubleshooting).
If you are building this checkout, use `make build` and `./telia` instead of
`telia`; see [contributing](../contributing/README.md).

Next: [programs and packages](../guide/programs.md).
