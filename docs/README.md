# Telia documentation

Telia is a compiled, statically typed language with explicit memory management.
These documents describe the code in this repository. They do not propose new
language behavior.

## Learn Telia

Start with [installation](getting-started/installation.md), then
[your first program](getting-started/first-program.md). Follow the
[language guide](guide/README.md) in order; no knowledge of Go or C is required.
Learn [unit testing](getting-started/testing.md) when you begin writing functions.

## Look something up

- [Language reference](reference/README.md): syntax, types, evaluation, memory,
  errors, generics, and C interoperability.
- [Standard library](stdlib/README.md): every application-facing `std` package,
  API contracts, ownership, failures, and platform restrictions.
- [Cookbook](cookbook/README.md): complete programs and task-oriented examples.
- [Compiler commands](reference/commands.md): building, testing, and configuration.

## Work on Telia

- [Agent entry point](agents.md): where to find authoritative rules and evidence.
- [Contributing](contributing/README.md): compiler development, checks, and
  documentation maintenance.

## How to read examples

In the getting-started and guide chapters, `telia` code fences are complete
programs. Save one as `main.t`, build it, and compare its output with the page.
They are also exercised by `TestDocumentationExamples`.
Other code fences may be fragments; their surrounding text supplies context.
Library and reference examples assume the imports and declarations described
on the page. Invalid examples are identified explicitly.

The guide explains concepts. The reference defines their rules. Library pages
define package contracts; the cookbook links to those contracts instead of
defining competing ones. Ownership and platform caveats apply even when an
example omits the surrounding application.
