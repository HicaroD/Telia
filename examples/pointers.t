package main

use "std::io"

fn pass_by_pointer(a *i32) {
  *a = 20
}

fn main() {
  a i32 := 10
  b := &a
  c := &b
  d := &c

  pass_by_pointer(**d)
  if a == 10 {
    io::println("a is 20")
  }
}
