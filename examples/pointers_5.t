package main

use "std::io"

fn main() i32 {
  a *i32 := nil
  b *i32 := nil
  c := &a
  if *c == b {
    io::puts("ptr equal")
  }
  return 0
}
