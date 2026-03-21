package main

use libc "std::libc"

fn mightFail() (i32, error) {
  return 99, nil
}

fn main() i32 {
  val := mightFail() @fail
  libc::printf("%d\n", val)
  return 0
}
