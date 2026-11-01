# `std::testing`

Source: [`testing.t`](../../base/std/testing/testing.t). Import `std::testing`
from a package's `*_test.t` files.

| Signature | Contract |
| --- | --- |
| `assert(condition bool)` | Fail the current test when false |
| `assert_eq<T>(actual T, expected T)` | Fail when operands differ |
| `assert_ne<T>(actual T, unexpected T)` | Fail when operands are equal |
| `fail()` | Unconditionally fail the current test |

No caller-owned storage is allocated or returned. Failure stops the current
test and runs pending Telia cleanup before the runner proceeds. This is not a
general exception mechanism; panics, process exits, and native crashes are not
recovered. Use these operations in named tests, not as application error handling.

Equality helpers infer types and support booleans, integers, floats, and strings;
string comparisons use bytes, floating-point comparisons are exact. Unsupported
or incompatible types produce compiler diagnostics. Direct assertion calls get
call-site diagnostics; calls through function values report the library location.

See [the testing guide](../getting-started/testing.md) for complete file setup,
expected output, discovery, failure cleanup, and native callback restrictions.
