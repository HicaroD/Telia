package main

use "std::io"

fn pass_by_pointer(a *i32) {
  *a = 20
  return
}

fn main() i32 {
  a i32 := 10
  b := &a
  c := &b

  pass_by_pointer(*c)
  io::printf("%d\n", a)

  return 0
}
