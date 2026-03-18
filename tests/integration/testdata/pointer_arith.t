package main

use "std::libc"

fn main() i32 {
  x := 42
  p := &x
  *p = 100
  libc::printf("%d\n", *p)
  libc::printf("%d\n", x)
  return 0
}
