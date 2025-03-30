package main

use libc "std::c"

fn pass_by_pointer(a *i32) {
  *a = 20
  return
}

fn main() i32 {
  a i32 := 10
  b := &a
  c := &b
  d := &c

  pass_by_pointer(**d)
  libc::printf("%d\n", a)

  return 0
}
