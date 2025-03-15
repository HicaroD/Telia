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
  a i32 := 10
  b := &a

  pass_by_value(*b)
  io::printf("pass_by_value: a=%d b=%d\n", a, *b)

  pass_by_pointer(b)
  io::printf("pass_by_pointer: a=%d b=%d\n", a, *b)

  return 0
}
