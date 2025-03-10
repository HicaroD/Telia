package main

use "std::io"

fn fib(n int) int {
  if n <= 1 {
    return n
  }
  return fib(n - 1) + fib(n - 2)
}

fn main() i32 {
  n := 40
  result := fib(i)
  io::libc::printf("Telia - Fibonacci(%d) = %d\n", n, result)
  return 0
}
