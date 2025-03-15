package main

use "std::io"

fn pass_by_value(a i32) {
  a = 20
  return
}

fn pass_by_pointer(a *i32) {
  *a = 20
  return
}

fn main() i32 {
  a := 10

  pass_by_value(a)
  io::printf("pass_by_value: %d\n", a)

  pass_by_pointer(&a)
  io::printf("pass_by_pointer: %d\n", a)

  return 0
}
