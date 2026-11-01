# Contributing to Telia

User documentation starts at [the documentation home](../README.md). This section
is for maintaining the compiler, libraries, and documentation.

## Build and check

From the repository root, install Go as required by `go.mod` and a C compiler,
then run:

```sh
make build
go test ./...
python3 scripts/check-docs.py
```

Use `./telia` for the freshly built compiler. `telia env` exposes library and
configuration paths; overrides pointing outside the checkout can select different
standard-library sources. See [CI](ci.md) for the pinned toolchain/environment,
formatting, lint, and workflow details.

## Where behavior lives

| Area | Sources/evidence |
| --- | --- |
| Tokens and scanning | `internal/lexer/` |
| Parsing | `internal/parser/` |
| Types and semantic checking | `internal/ast/`, `internal/sema/` |
| C lowering | `internal/codegen/c/` |
| Native targets and package manifests | `internal/ctarget/`, `internal/native/` |
| CLI commands | `cmd/compiler/` |
| Standard library | `base/std/` |
| Native packages | `base/vendor/` |
| End-to-end language/runtime checks | `tests/integration/` |
| Compilation test helpers | `tests/compiler/` |

Focused tests are useful while editing; full checks and supported-platform CI
remain separate evidence. Never describe a configured platform matrix as a
passing run. See [release maintenance](releases.md) for existing publishing
procedures and [native libraries](../stdlib/native-libraries.md) for packaging.

## Documentation maintenance

Choose the authoritative home before adding a page:

- `getting-started/`: setup and the first successful program.
- `guide/`: ordered teaching chapters, concepts, small examples, mistakes.
- `reference/`: exact language rules and limitations.
- `stdlib/`: application-facing package contracts.
- `cookbook/`: complete task-oriented programs linking to those contracts.
- `contributing/`: maintenance procedures and historical evidence.
- `agents.md`: navigation and verification, not a duplicate specification.

For a public behavior change:

1. Update the authoritative language/API contract, including edge cases, errors,
   effects on failure, ownership, cleanup, and target restrictions.
2. Update related lessons and examples; avoid repeating full contracts in them.
3. Update the package/reference index when adding or moving public material.
4. Add or update compiler/runtime tests. Do not invent expected behavior solely
   from grammar or an unverified design note.
5. Run the link checker and documentation-example tests, then relevant suites.
6. Update links to moved pages rather than leaving duplicate documents at the
   old paths.

### Example conventions

Use a `telia` code fence for a standalone runnable program. Follow it immediately
with `Expected output:` and a `text` fence. `TestDocumentationExamples` compiles
every such program at `-O0` and `-O3` and compares stdout. Keep these examples
deterministic and self-contained: no terminal, network, or fixture-file setup.

Use `rust` fences for existing highlighted fragments; introduce their imports,
surrounding function, or other prerequisites. Identify intentionally invalid
code and the expected diagnostic in prose. Native/interactively configured
examples should link to source files and their existing integration tests rather
than pretending to be self-contained documentation tests.

`scripts/check-docs.py` checks repository-local Markdown links and heading
anchors from the documentation, root entry points, and example readmes. It also
checks that every standard-library source package appears in the package index.
It does not validate external websites or make semantic claims about code fences.

## Historical evidence

[Internal semantic hashing](hashing.md) records the compiler/library key and
checked seed boundary, native adapters, and focused verification commands.

[Historical verification notes](verification-records.md) preserve local
implementation-time platform/sanitizer observations. The
[standard-library CLI coverage audit](stdlib-cli-coverage.md) preserves an
implementation-time acceptance/evidence record. Its local results and pending
checks are historical; use current test runs to assess current changes.
