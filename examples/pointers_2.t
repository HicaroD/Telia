package main

use "std::io"

fn main() i32 {
  a i32 := 10
  b := &a
  c := &b
  **c = 300
  io::printf("%d\n", **c)
  return 0
}
