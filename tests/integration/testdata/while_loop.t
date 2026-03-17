package main

use libc "std::libc"

fn main() i32 {
  n := 1
  while n <= 3 {
    libc::printf("%d\n", n)
    n = n + 1
  }
  return 0
}
