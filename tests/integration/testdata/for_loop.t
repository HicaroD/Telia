package main

use libc "std::libc"

fn main() i32 {
  sum := 0
  for i := 0; i < 5; i = i + 1 {
    sum = sum + i
  }
  libc::printf("%d\n", sum)
  return 0
}
