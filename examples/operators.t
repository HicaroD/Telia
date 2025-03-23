package main

use "std::io"

fn main() {
  a i32 := 1 + 1
  b i32 := 1 - 1
  c i32 := 2 * 2
  d i32 := 10 / 2
  io::printf("1+1=%d\n1-1=%d\n2*2=%d\n10/2=%d\n", a, b, c, d)
  return
}
