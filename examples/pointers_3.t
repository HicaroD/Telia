package main

use "std::io"

fn get_ref(a **i32) {
  **a = 1000
  return
}

fn main() i32 {
  a i32 := 10
  b := &a
  c := &b
  get_ref(c)
  io::printf("%d\n", **c)
  return 0
}
