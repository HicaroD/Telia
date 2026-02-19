package main

use libc "std::libc"

fn fib(n i32) i32 {
  if n <= 1 {
    return n
  }
  return fib(n - 1) + fib(n - 2)
}

fn main() i32 {
  libc::printf("%d\n", fib(10))
  return 0
}
