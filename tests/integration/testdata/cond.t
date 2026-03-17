package main

use libc "std::libc"

fn sign(n i32) i32 {
  if n < 0 {
    return -1
  } elif n == 0 {
    return 0
  } else {
    return 1
  }
}

fn main() i32 {
  libc::printf("%d\n", sign(-5))
  libc::printf("%d\n", sign(0))
  libc::printf("%d\n", sign(3))
  return 0
}
