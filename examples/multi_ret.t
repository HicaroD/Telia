package main

use "std::io"

fn get(value i32) (i32, i32) {
  return 1, 2 + value
}

fn other(a i32, b i32) (i32, i32) {
  a, b = get(1)
  return a, b
}

fn other2(a i32, b i32) (i32, i32) {
  return a, b
}

fn main() i32 {
  a, b := other2(1, 3)
  io::libc::printf("%d %d", a, b)
  return 0
}
