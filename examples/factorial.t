package main

use "std::io"

fn factorial(n i32) i32 {
  if n == 1 {
    return n
  }
  return n * factorial(n - 1)
}

fn main() i32 {
  result := factorial(6)
  io::printf("%d\n", result)
  return 0
}
