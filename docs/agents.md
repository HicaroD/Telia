# Agent entry point

Use the same language and API contracts as human readers. This page is a
navigation aid, not another language specification.

## Read by task

| Task | Start here |
| --- | --- |
| Learn syntax or write a small program | [First program](getting-started/first-program.md), [guide](guide/README.md) |
| Determine a language rule | [Reference index](reference/README.md) |
| Use a standard-library API | [Package index](stdlib/README.md) |
| Reason about cleanup or copying | [Memory rules](reference/memory.md), [ownership contracts](stdlib/strings-and-memory.md#ownership-conventions-and-transfer) |
| Handle a failing operation | [Errors](reference/errors.md), then the package's failure contract |
| Build, test, or resolve imports | [Commands](reference/commands.md), [installation](getting-started/installation.md) |
| Call C or use native packages | [Foreign rules](reference/foreign-interop.md), [C types](stdlib/ffi-c.md), [native packages](stdlib/native-libraries.md) |
| Change compiler behavior | [Contributor guide](contributing/README.md) |

## Avoid cross-language assumptions

- Declarations are `name Type := value`, not `let`, `var`, or `name: Type`.
- Imports use `use "std::io"`; package access is `io::println(...)`.
- Logical operators are `and`, `or`, `not`. There are no `++`, `%`, or implicit
  string `+` operations.
- Struct literals are `Type{field: value}`; contextual arrays are `.{...}`.
- `string` and slices borrow storage. Copying owners is shallow, not a move.
- `defer` is block-scoped and reads current variables at cleanup time.
- `@try` propagates with cleanup; `@fail` skips pending cleanup.
- Indices use `u64`; do not assume other integer types convert automatically.
- A source-visible library helper is not necessarily an application API.
- Do not assume traits, methods, closures, maps, automatic destruction, or borrow
  checking exist because similar-looking languages provide them.

## Check evidence

1. Read the relevant prose contract and linked source.
2. Find corresponding compiler/integration tests before inferring edge behavior.
3. If source, prose, and tests disagree, identify the conflict. Do not silently
   redesign the language or present an untested assumption as a rule.
4. Build and run examples against this checkout. Do not claim cross-platform or
   sanitizer success from a local ordinary run.

Useful checks from the repository root:

```sh
python3 scripts/check-docs.py
go test ./tests/integration -run '^TestDocumentationExamples$' -count=1
go test ./...
```

The [grammar outline](reference/grammar.ebnf) is not a complete semantic contract.
Historical audits in `contributing/` record past evidence, not current test results.
When behavior changes, update its authoritative page and runnable examples as
described in [documentation maintenance](contributing/README.md#documentation-maintenance).
