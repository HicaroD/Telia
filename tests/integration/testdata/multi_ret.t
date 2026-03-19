package main

use libc "std::libc"

fn get(value i32) (i32, i32) {
  return 1, 2 + value
}

fn main() i32 {
  a, b := get(3)
  libc::printf("%d\n", a)
  libc::printf("%d\n", b)

  c i32, d i64 := 10, 20
  libc::printf("%d\n", c)
  libc::printf("%lld\n", d)

  return 0
}
