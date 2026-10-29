package main

use "std::libc"

fn get_ref(a **i32) {
  **a = 1000
  return
}

fn main() i32 {
  a := 10
  b := &a
  c := &b
  get_ref(c)
  libc::printf("%d\n", **c)
  return 0
}
