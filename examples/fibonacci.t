package main

use "std::io"

fn fib(n i32) i32 {
  if n <= 1 {
    return n
  }
  return fib(n - 1) + fib(n - 2)
}

fn main() i32 {
  n := 40
  io::printf("Telia - Fibonacci(%d) = %d\n", n, fib(n))
  return 0
}
