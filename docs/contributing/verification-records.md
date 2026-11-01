# Historical verification notes

These observations were recorded in the feature documents before their
reorganization. They are preserved as historical evidence, **not** results from
the documentation change or a claim about the current checkout.

## C interoperability, enums, and bitwise operations

The interoperability implementation tracked in issue #93 recorded separately
compiled native C regressions on macOS arm64 only. Enum fixtures likewise
recorded local compilation/execution on macOS arm64. Bitwise/shift tests recorded
signed/unsigned and 128-bit boundaries under UBSan on that host. None of these
local observations established other-target compatibility.

## Native CLI expansion

The native CLI feature document recorded macOS arm64 success with UBSan at both
optimization levels. ASan+UBSan fixtures timed out before producing output; an
independent C program with an explicitly flushed entry message also timed out
without reaching that message. This was an unresolved local ASan startup gap,
not an ASan pass. The full Go suite, build, formatting, pinned lint, and Python
script tests were recorded as passing locally. Cross-platform/ASan results were
pending CI at that time.

The earlier [standard-library CLI acceptance audit](stdlib-cli-coverage.md)
contains its own detailed evidence and deferred checks. Do not combine results
from different feature deliveries as if they verified one current commit.
