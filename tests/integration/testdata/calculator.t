package main

use libc "std::libc"

fn add(a i32, b i32) i32 {
  return a + b
}

fn sub(a i32, b i32) i32 {
  return a - b
}

fn mul(a i32, b i32) i32 {
  return a * b
}

fn main() i32 {
  libc::printf("%d\n", add(5, 3))
  libc::printf("%d\n", sub(10, 4))
  libc::printf("%d\n", mul(6, 7))
  return 0
}
