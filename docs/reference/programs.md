# Program structure

## Files and declarations

Telia source files end in `.t`. A file starts with `package Name`, followed by
`use` imports and declarations. Imports must precede other declarations.
Package declarations include functions, structs, enums, type aliases, constants,
and `extern` blocks. Ordinary local variables belong in function bodies.
Named `test` blocks are allowed only in `*_test.t` files.

A directory build collects the package's source files. A single-file build
starts with that file; do not rely on it to include neighboring files. Ordinary
builds exclude test files. Imported packages' tests are not included by `test`.

Use `fn main() i32` with `return 0` for executable examples with an explicit
successful exit status. Non-entry functions may omit their result type when
they return no value.

## Imports and qualification

```rust
use "std::io"
use text "std::strings"
use "pkg::helpers"
use "vendor::sdl3"
```

An import path uses `::` separators. Its prefix selects the standard library,
the application package tree rooted at the build's input directory, or the
configured vendor tree. Nested imports such as `std::unicode::utf8` bind the
last component by default. An explicit alias replaces that local name.

Imports are file-local. Declarations live in the package's scope, but each file
must import the packages it uses. Imported names use `package::Name`; fields use
`value.field`. A named `extern` namespace adds a qualifier, such as
`libc::strings::strlen`.

Do not infer Go-style visibility from capitalization: the language has no
`pub`/`private` declarations here. Source accessibility is not an API stability
promise; library pages distinguish application APIs from implementation helpers.

## Lexical conventions

- Use ASCII letters, digits, and `_` for portable identifiers; do not start an
  identifier with a digit. Keywords and built-in type words are reserved.
- Newlines normally separate statements. Semicolons separate the three clauses
  of a `for`; they are not a general replacement for statement newlines.
- Use `//` for line comments and `/* ... */` for non-nesting block comments.
  Newlines inside block comments can still separate statements.
- Strings use double quotes and may contain UTF-8 text. Supported escapes are
  `\n`, `\t`, `\\`, `\"`, and `\0`. Do not assume C's other escapes work.
- Use decimal integer literals and decimal floating-point literals; signed
  values use unary `-`. See [constants](constants-and-bitwise.md) for evaluation
  and overflow rules, and [literals](literals.md) for aggregates.

Diagnostics use source locations. An informational refactoring hint does not
prevent a build; semantic or backend errors do.

See [commands](commands.md) for package builds, paths, and test discovery.
