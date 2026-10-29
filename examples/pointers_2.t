package main

use "std::libc"

fn main() i32 {
  a i32 := 10
  b := &a
  c := &b
  **c = 300
  libc::printf("%d\n", **c)
  return 0
}
