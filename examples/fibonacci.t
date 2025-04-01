package main

use libc "std::c"

fn fib(n i32) i32 {
  if n <= 1 {
    return n
  }
  return fib(n - 1) + fib(n - 2)
}

fn get() i32 {
  if true {
    return 1
  }
  return 0
}

fn main() i32 {
  n := 40
  libc::printf("Telia - Fibonacci(%d) = %d\n", n, fib(n))
  return 0
}
