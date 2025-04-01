package main

use "std::libc"

fn okay() {
  a i32 := 10
  b := &a
  c := *b
  libc::puts("OKAY")
}

fn causes_panic() {
  a *i32 := nil
  b := *a
  libc::puts("NULL DEREF - UNREACHABLE")
}

fn main() i32 {
  okay()
  causes_panic()
  return 0
}
