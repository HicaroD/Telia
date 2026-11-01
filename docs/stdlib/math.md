# `std::math`

Source: [`math.t`](../../base/std/math/math.t). Import `std::math`.
These are native C math calls, not checked Telia operations.

| Signature | Meaning |
| --- | --- |
| `cos(x f64) f64`, `sin(x f64) f64`, `tan(x f64) f64` | Trigonometric functions; input is radians |
| `acos(x f64) f64`, `asin(x f64) f64`, `atan(x f64) f64` | Inverse trigonometric functions; result is radians |
| `acosf(x f32) f32` | Single-precision inverse cosine |
| `atan2(y f64, x f64) f64` | Angle determined by both coordinates |
| `pow(y f64, x f64) f64` | First argument raised to the second argument |
| `powf(y f32, x f32) f32` | Single-precision power |

Example expression: `math::cos(0.0)` returns one. Inputs and results are scalar
values; these operations acquire no owned storage and require no cleanup.

Domain/range errors follow the target C math implementation (including NaN,
infinity, and native floating-point behavior); no `error` result, portable errno
wrapper, or exact cross-platform rounding guarantee is provided. Use suitable
input checks and tolerances for floating-point comparisons. The package does
not currently expose a general-purpose checked arithmetic API.
