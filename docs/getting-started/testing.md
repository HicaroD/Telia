# Unit testing

Telia discovers named test blocks in `*_test.t` files alongside the package's
ordinary source files. Use the same package declaration: tests can call
same-package functions directly. No test `main` is needed, even if the package
already defines one.

## Setup and runnable example

Install Telia and a C compiler using the [installation guide](installation.md).
From a source checkout, build and run the included example:

```sh
make build
cd examples/testing/math
../../../telia test
```

For your own package, create `math/math.t`:

```rust
package math

fn add(a i32, b i32) i32 {
    return a + b
}
```

Then create `math/math_test.t`:

```rust
package math

use "std::testing"

fn check_sum(a i32, b i32, expected i32) {
    testing::assert_eq(add(a, b), expected)
}

test "adds positive and negative numbers" {
    check_sum(2, 3, 5)
    check_sum(-2, 3, 1)
}
```

Run `telia test` from `math/`:

```text
PASS adds positive and negative numbers
1 passed, 0 failed
```

Ordinary functions are reusable helpers, not tests. Functions named `test_*`
are not automatically discovered.

## Assertions

Import `std::testing` and use parenthesized calls:

| Call | Passes when |
| --- | --- |
| `testing::assert(condition)` | `condition` is `true` |
| `testing::assert_eq(actual, expected)` | The values are equal |
| `testing::assert_ne(actual, unexpected)` | The values differ |

Equality assertions infer their operand type; no explicit type arguments are
needed. They support booleans, integers (including 128-bit integers), floats,
and strings. Strings compare their contents, including embedded zero bytes;
floating-point equality is exact. Unsupported types and incompatible operands
produce compilation diagnostics. Compare individual fields or use a boolean
condition for other types.

Operands are evaluated once, from left to right. A failed assertion stops the
current test, including calls through helper functions, and remaining tests
continue. Failure details include the test description and assertion source
location. Equality failures show expected and actual values; inequality failures
show unexpected and actual values.

Tests execute sequentially in filename/declaration order in one process. An
assertion failure unwinds Telia helper calls and runs pending `defer` actions in
reverse order before returning to the runner. Cleanup sees the current values
of local variables, just as on a normal scope exit. If an assertion fails during
cleanup, the rest of that action is skipped, remaining actions still run once,
and both the original and cleanup failures are reported. A test with multiple
failed assertions during cleanup counts as one failed test.

Panics, process exits, and native crashes are not recovered; they stop the runner
and make the command fail. Native code is not unwound: a failure in a Telia
callback propagates when control returns from its native caller. This first
version has no setup or cleanup hooks, parallel execution, or benchmarks.

Use direct assertion calls for call-site diagnostics; calling `assert` through
a function value reports its library implementation location. Assertion names
are provided by `std::testing`, not reserved language keywords.

## Discovery, output, and CI

- Only the current package's test files are included; imported packages' test
  files are always excluded. Subdirectories are not searched for tests.
- Test descriptions must be unique within the package. Test blocks outside
  `*_test.t` files produce an error.
- Normal builds exclude `*_test.t` files, including imports. Building a test
  file explicitly is rejected; use `telia test` instead.
- No tests prints `no tests found` and a zero-count summary, then succeeds.
- Successful runs return zero. Assertion failures, compilation failures, and
  runner errors return nonzero, suitable for CI.
- The test executable is generated in a temporary directory and removed after
  execution; the package's normal build output is not overwritten.

`PASS` is green and `FAIL` is red on a terminal. Control colors with:

```sh
telia test --color=auto
telia test --color=always
telia test --color=never
```

The default `auto` disables colors for redirected output and respects
`NO_COLOR`. Explicit `always` overrides `NO_COLOR`, matching compiler diagnostic
behavior. Results and the summary go to stdout; failure details and compiler
diagnostics go to stderr.
