package main

use "std::libc"

fn main() i32 {
  a f64 := 2.0 + 1.0
  libc::printf("%lf\n", a)
  return 0
}
